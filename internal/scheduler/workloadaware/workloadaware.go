package workloadaware

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/config"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/client"
	clientset "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/clientset/versioned"
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
)

// WorkloadAware is a Scheduler Plugin that takes into account the workload type of a Pod
// to make scheduling and resource binding decisions, based on the following classification:
// 	- MemoryBound: Workloads that require as many memory nodes (sockets) as possible
// 	- CPUBound: Workloads with execution time that depends on the available CPU resources
// 	- IOBound: Workloads that have threads with high IO wait time
// 	- BestEffort: Workloads that place every thread on the same logical CPU (oversubscription)

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
var _ framework.ReservePlugin = &WorkloadAware{}

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
		stateData.AllocatableCPUs[t.Spec.NodeName] = allocatableCPUsForNode(t.Spec.NodeName, &t.Spec.Topology, cpuBindings)
		stateData.Topologies[t.Spec.NodeName] = t.Spec.Topology
	}

	state.Write(framework.StateKey(Name), stateData)
	return &framework.PreFilterResult{NodeNames: nodes}, nil
}

func (w *WorkloadAware) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	node := nodeInfo.Node()
	logger := w.logger.WithName("filter").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", node.Name)

	stateData, err := w.getState(state)
	if err != nil {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("could not get state data: %v", err))
	}

	featurePhysicalCores := slices.Contains(w.args.Features, config.FeaturePhysicalCores)
	cpuRequests, cpuLimits := stateData.PodRequests.Cpu().MilliValue(), stateData.PodLimits.Cpu().MilliValue()
	allocatable := stateData.AllocatableCPUs[node.Name]
	numAllocatableCPUs := int64(len(allocatable) * 1000)
	if featurePhysicalCores || stateData.WorkloadType == config.WorkloadTypeCPUBound {
		numAllocatableCPUs = int64(len(allocatableCores(allocatable, true)) * 1000)
	}

	if numAllocatableCPUs == 0 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("all cpus are reserved for other pods"))
	}
	if cpuRequests != 0 && numAllocatableCPUs < cpuRequests {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("not enough cpus, request: %dm, available: %dm", cpuRequests, numAllocatableCPUs))
	}
	if cpuLimits != 0 && numAllocatableCPUs < cpuLimits {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("not enough cpus, limit: %dm, available: %dm", cpuLimits, numAllocatableCPUs))
	}

	logger.Info("filter passed", "cpuRequests", cpuRequests, "cpuLimits", cpuLimits, "numAllocatableCPUs", numAllocatableCPUs)
	return framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) Score(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (int64, *framework.Status) {
	stateData, err := w.getState(state)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("could not get state data: %v", err))
	}
	allocatable := stateData.AllocatableCPUs[nodeName]
	featurePhysicalCores := slices.Contains(w.args.Features, config.FeaturePhysicalCores)

	score := int64(0)
	switch stateData.WorkloadType {
	case config.WorkloadTypeMemoryBound:
		// MemoryBound workloads require as many memory nodes (sockets) as possible
		// We have assured that all scored nodes can fit the Pod's requirements,
		// so we just need to count the number of sockets
		score = int64(len(allocatableSockets(allocatable))*1000 + len(allocatableCores(allocatable, featurePhysicalCores))*100 + len(allocatable)*10)
	case config.WorkloadTypeCPUBound:
		// CPUBound workloads place threads on different, non-utilized cores, to avoid interference
		score = int64(len(allocatableCores(allocatable, true))*100 + len(allocatable)*10)
	case config.WorkloadTypeIOBound:
		// IOBound workloads place threads on the same physical core
		score = int64(len(allocatableCores(allocatable, featurePhysicalCores))*100 + len(allocatable)*10)
	case config.WorkloadTypeBestEffort:
		// BestEffort workloads place every thread on the same logical CPU
		// So we need the number of all logical CPUs
		if featurePhysicalCores {
			score = int64(len(allocatableCores(allocatable, true))*100 + len(allocatable)*10)
		} else {
			score = int64(len(allocatable) * 10)
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

	return score, framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	logger := w.logger.WithName("normalize").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))

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

	logger.Info("normalized scores", "scores", scores)
	return framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) Reserve(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) *framework.Status {
	logger := w.logger.WithName("reserve").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", nodeName)
	stateData, err := w.getState(state)
	if err != nil {
		w.logger.Error(err, "could not get state data")
		return framework.NewStatus(framework.Error, fmt.Sprintf("could not get state data: %v", err))
	}

	var fullCores = false
	var cpuSet []v1alpha1.CPU
	var exclusivenessLevel = v1alpha1.ResourceLevelCPU

	if slices.Contains(w.args.Features, config.FeaturePhysicalCores) {
		fullCores = true
		exclusivenessLevel = v1alpha1.ResourceLevelCore
	}
	if stateData.WorkloadType == config.WorkloadTypeMemoryBound &&
		slices.Contains(w.args.Features, config.FeatureMemoryBoundExclusiveSockets) {
		exclusivenessLevel = v1alpha1.ResourceLevelSocket
	}

	cpuSet = cpuSetForWorkloadType(stateData, nodeName, stateData.WorkloadType, fullCores)
	logger.Info("cpuSet", "cpuSet", cpuSet)

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

	logger.Info("created pod cpu binding", "pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", nodeName, "cpuSet", cpuBinding.Spec.CPUSet)
	return framework.NewStatus(framework.Success)
}

func (w *WorkloadAware) Unreserve(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) {
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
