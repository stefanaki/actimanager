package workloadaware

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/config"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/client"
	clientset "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/clientset/versioned"
	nctutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/v1/resource"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"math"
	"slices"
	"time"
)

const Name string = "WorkloadAware"

type WorkloadAware struct {
	args   *config.WorkloadAwareArgs
	handle framework.Handle
	client *clientset.Clientset
	logger klog.Logger
}

type State struct {
	WorkloadType    string
	PodRequests     corev1.ResourceList
	PodLimits       corev1.ResourceList
	Topologies      map[string]v1alpha1.CPUTopology
	AllocatableCPUs map[string]NodeAllocatableCPUs
}

var _ framework.PreFilterPlugin = &WorkloadAware{}
var _ framework.FilterPlugin = &WorkloadAware{}
var _ framework.ScorePlugin = &WorkloadAware{}
var _ framework.BindPlugin = &WorkloadAware{}

func (w *WorkloadAware) Name() string {
	return Name
}

func New(ctx context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	args, err := parseArgs(obj)
	if err != nil {
		return nil, err
	}
	csLabClient, err := client.NewCSLabClient()
	if err != nil {
		return nil, err
	}
	l := klog.NewKlogr().WithName(Name)
	l.Info("args", "policy", args.Policy, "features", args.Features)
	return &WorkloadAware{args: args, handle: h, client: csLabClient, logger: l}, err
}

func (w *WorkloadAware) PreFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod) (*framework.PreFilterResult, *framework.Status) {
	logger := w.logger.WithName("pre-filter").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
	nodes := sets.Set[string]{}
	workloadType, ok := pod.Labels[config.LabelWorkloadType]
	if !ok {
		logger.Info("Application type not specified, assuming best effort")
		workloadType = config.WorkloadTypeBestEffort
	}

	// Wait for previous PodCPUBindings to be validated
	time.Sleep(4 * time.Second)

	stateData := &State{
		WorkloadType:    workloadType,
		PodRequests:     resource.PodRequests(pod, resource.PodResourcesOptions{}),
		PodLimits:       resource.PodLimits(pod, resource.PodResourcesOptions{}),
		Topologies:      make(map[string]v1alpha1.CPUTopology),
		AllocatableCPUs: make(map[string]NodeAllocatableCPUs),
	}

	topologies, err := w.client.CslabV1alpha1().NodeCPUTopologies().List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error(err, "failed to get node cpu topologies")
		return nil, framework.NewStatus(framework.Error, "failed to get node cpu topologies")
	}
	cpuBindings, err := w.client.CslabV1alpha1().PodCPUBindings("").List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error(err, "failed to get pod cpu bindings")
		return nil, framework.NewStatus(framework.Error, "failed to get pod cpu bindings")
	}

	for _, t := range topologies.Items {
		if t.Status.ResourceStatus != v1alpha1.StatusFresh {
			continue
		}
		nodes.Insert(t.Spec.NodeName)
		stateData.AllocatableCPUs[t.Spec.NodeName] = allocatableCPUsForNode(t.Spec.NodeName, &t.Spec.Topology, cpuBindings, workloadType)
		stateData.Topologies[t.Spec.NodeName] = t.Spec.Topology
	}

	state.Write(framework.StateKey(Name), stateData)
	return &framework.PreFilterResult{NodeNames: nodes}, nil
}

func (w *WorkloadAware) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	node := nodeInfo.Node()
	// logger := w.logger.WithName("filter").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", node.Name)

	stateData, err := w.getState(state)
	if err != nil {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("could not get state data: %v", err))
	}

	cpuRequests, cpuLimits := stateData.PodRequests.Cpu().MilliValue(), stateData.PodLimits.Cpu().MilliValue()
	allocatable := stateData.AllocatableCPUs[node.Name]
	numAllocatable := int64(len(allocatable) * 1000)
	noun := "cpus"

	featurePhysicalCores := slices.Contains(w.args.Features, config.FeaturePhysicalCores)
	featureMemoryBoundExclusiveCores := slices.Contains(w.args.Features, config.FeatureMemoryBoundExclusiveSockets)

	if featurePhysicalCores || stateData.WorkloadType == config.WorkloadTypeCPUBound {
		noun = "cores"
		numAllocatable = int64(len(allocatableCores(allocatable, true)) * 1000)
	}
	if featureMemoryBoundExclusiveCores && stateData.WorkloadType == config.WorkloadTypeMemoryBound {
		noun = "sockets"
		numAllocatable = int64(len(allocatableSockets(allocatable, true)) * 1000)
	}

	if numAllocatable == 0 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("all %s are reserved for other pods", noun))
	}
	if cpuRequests != 0 && numAllocatable < cpuRequests {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("not enough %s, request: %dm, available: %dm", noun, cpuRequests, numAllocatable))
	}
	if cpuLimits != 0 && numAllocatable < cpuLimits {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("not enough %s, limit: %dm, available: %dm", noun, cpuLimits, numAllocatable))
	}

	return framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) Score(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (int64, *framework.Status) {
	logger := w.logger.WithName("score").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", nodeName)

	stateData, err := w.getState(state)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("could not get state data: %v", err))
	}
	allocatable := stateData.AllocatableCPUs[nodeName]
	topology := stateData.Topologies[nodeName]

	featurePhysicalCores := slices.Contains(w.args.Features, config.FeaturePhysicalCores)
	featureMemoryBoundExclusiveSockets := slices.Contains(w.args.Features, config.FeatureMemoryBoundExclusiveSockets)

	score := int64(0)
	switch stateData.WorkloadType {
	case config.WorkloadTypeMemoryBound:
		// MemoryBound workloads require as many memory nodes (sockets) as possible
		// We have assured that all scored nodes can fit the Pod's requirements,
		// so we just need to count the number of sockets
		numAllocatableSockets := len(allocatableSockets(allocatable, featureMemoryBoundExclusiveSockets))
		numAllSockets := len(stateData.Topologies[nodeName].Sockets)
		score = int64(math.Ceil(float64(numAllocatableSockets)/float64(numAllSockets)*10000)) + int64(numAllocatableSockets)
	case config.WorkloadTypeCPUBound:
		// CPUBound workloads place threads on different, non-utilized cores, to avoid interference
		numAllocatableCores := len(allocatableCores(allocatable, true))
		numAllCores := len(nctutils.CoresInTopology(&topology))
		score = int64(math.Ceil(float64(numAllocatableCores)/float64(numAllCores))*10000) + int64(numAllocatableCores)
	case config.WorkloadTypeIOBound:
		// IOBound workloads place threads on the same physical core
		numAllocatableCores := len(allocatableCores(allocatable, featurePhysicalCores))
		numAllCores := len(nctutils.CoresInTopology(&topology))
		score = int64(math.Ceil(float64(numAllocatableCores)/float64(numAllCores)*10000)) + int64(numAllocatableCores)
	case config.WorkloadTypeBestEffort:
		// BestEffort workloads place every thread on the same logical CPU
		// So we need the number of all logical CPUs
		if featurePhysicalCores {
			numAllocatableCores := len(allocatableCores(allocatable, true))
			numAllCores := len(nctutils.CoresInTopology(&topology))
			score = int64(math.Ceil(float64(numAllocatableCores)/float64(numAllCores))*10000) + int64(numAllocatableCores)
		} else {
			numAllocatableCPUs := len(allocatable)
			numAllCPUs := len(stateData.Topologies[nodeName].CPUs)
			score = int64(math.Ceil(float64(numAllocatableCPUs)/float64(numAllCPUs))*10000) + int64(numAllocatableCPUs)
		}
	}

	switch w.args.Policy {
	case config.PolicyMaximumUtilization:
		// Score is a metric of the capacity
		// In MaximumUtilization we should favor nodes with the least capacity in resources
		score = math.MaxInt64 - score
	case config.PolicyBalanced:
		// Balanced policy does not require any further action
	}

	logger.Info("scored node", "score", score)

	return score, framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	// logger := w.logger.WithName("normalize").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))

	// Find highest and lowest scores.
	var highest int64 = -math.MaxInt64
	var lowest int64 = math.MaxInt64
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	// Transform the highest to lowest score range to fit the framework's min to max node score range.
	oldRange := highest - lowest
	newRange := framework.MaxNodeScore - framework.MinNodeScore
	for i, nodeScore := range scores {
		if oldRange == 0 {
			scores[i].Score = framework.MinNodeScore
		} else {
			scores[i].Score = ((nodeScore.Score - lowest) * newRange / oldRange) + framework.MinNodeScore
		}
	}

	// logger.Info("normalized scores", "scores", scores)
	return framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) Bind(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) *framework.Status {
	logger := w.logger.WithName("bind").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", nodeName)

	// Bind the Pod to the Node
	podBinding := &corev1.Binding{
		ObjectMeta: metav1.ObjectMeta{Namespace: pod.Namespace, Name: pod.Name, UID: pod.UID},
		Target:     corev1.ObjectReference{Kind: "Node", Name: nodeName},
	}
	err := w.handle.ClientSet().CoreV1().Pods(podBinding.Namespace).Bind(ctx, podBinding, metav1.CreateOptions{})
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("failed to bind pod: %v", err))
	}

	// Create a PodCPUBinding resource to bind the Pod's threads to the Node's CPUs
	stateData, err := w.getState(state)
	if err != nil {
		w.logger.Error(err, "could not get state data")
		return framework.NewStatus(framework.Error, fmt.Sprintf("could not get state data: %v", err))
	}

	var cpuSet []v1alpha1.CPU
	var exclusivenessLevel = v1alpha1.ResourceLevelCPU

	var featureFullCores = slices.Contains(w.args.Features, config.FeaturePhysicalCores)
	var featureMemoryBoundExclusiveSockets = slices.Contains(w.args.Features, config.FeatureMemoryBoundExclusiveSockets)
	var featureBestEffortSharedCPUs = slices.Contains(w.args.Features, config.FeatureBestEffortSharedCPUs)

	if featureFullCores {
		exclusivenessLevel = v1alpha1.ResourceLevelCore
	}
	if featureMemoryBoundExclusiveSockets && stateData.WorkloadType == config.WorkloadTypeMemoryBound {
		exclusivenessLevel = v1alpha1.ResourceLevelSocket
	}
	if featureBestEffortSharedCPUs && stateData.WorkloadType == config.WorkloadTypeBestEffort {
		exclusivenessLevel = v1alpha1.ResourceLevelNone
	}

	switch stateData.WorkloadType {
	case config.WorkloadTypeMemoryBound:
		cpuSet = cpuSetForMemoryBound(stateData, nodeName, featureMemoryBoundExclusiveSockets)
	case config.WorkloadTypeCPUBound:
		cpuSet = cpuSetForCPUBound(stateData, nodeName, featureFullCores)
	case config.WorkloadTypeIOBound:
		cpuSet = cpuSetForIOBound(stateData, nodeName, featureFullCores)
	case config.WorkloadTypeBestEffort:
		cpuSet = cpuSetForBestEffort(stateData, nodeName, featureFullCores)
	}

	cpuBinding := v1alpha1.PodCPUBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-pcb", pod.Name),
			Namespace: pod.Namespace,
		},
		Spec: v1alpha1.PodCPUBindingSpec{
			PodName:            pod.Name,
			CPUSet:             cpuSet,
			ExclusivenessLevel: exclusivenessLevel,
		},
	}
	_, err = w.client.CslabV1alpha1().PodCPUBindings(pod.Namespace).Create(ctx, &cpuBinding, metav1.CreateOptions{})
	if err != nil {
		w.logger.Error(err, "failed to create pod cpu binding")
		return framework.NewStatus(framework.Error, fmt.Sprintf("failed to create pod cpu binding: %v", err))
	}

	logger.Info("pod scheduled", "pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", nodeName, "cpuSet", cpuBinding.Spec.CPUSet)

	return framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func (w *WorkloadAware) ScoreExtensions() framework.ScoreExtensions {
	return w
}

func (s State) Clone() framework.StateData {
	return s
}
