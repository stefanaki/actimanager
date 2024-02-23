package podisolation

import (
	"context"
	"flag"
	"fmt"
	"math"
	"strconv"

	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	nctutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/podcpubinding"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PodIsolation is responsible for filtering and scoring nodes based on CPU topology and CPU bindings.
// It ensures that pods are scheduled on nodes with feasible CPUs and calculates the score based on the locality of the feasible CPUs.
// This plugin is used to optimize CPU resource allocation in a Kubernetes cluster.

// Name is the name of the PodIsolation plugin.
const Name string = "PodIsolation"

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

// PodIsolation is the implementation of the PodIsolation plugin.
// It embeds the Kubernetes client, framework handle, and logger.
type PodIsolation struct {
	client.Client
	handle framework.Handle
	logger klog.Logger
}

// PodIsolationStateData represents the state data for the PodIsolation.
type PodIsolationStateData struct {
	// ExclusivenessLevel is the level of exclusiveness of the resources needed by the pod.
	ExclusivenessLevel string

	// NodeCpuTopologies stores the CPU topology of each node in the cluster.
	NodeCpuTopologies map[string]v1alpha1.NodeCpuTopology

	// PodCpuBindings stores the CPU bindings for each node in the cluster.
	PodCpuBindings map[string][]v1alpha1.PodCpuBinding

	// FeasibleCpus stores the feasible CPUs for each node in the cluster.
	FeasibleCpus map[string]v1alpha1.CpuTopology

	// PodCpuRequests is the total CPU requests of the pod.
	PodCpuRequests int64
}

var _ framework.PreFilterPlugin = &PodIsolation{}
var _ framework.FilterPlugin = &PodIsolation{}
var _ framework.ScorePlugin = &PodIsolation{}
var _ framework.ScoreExtensions = &PodIsolation{}
var _ framework.PostBindPlugin = &PodIsolation{}

// Name returns the name of the PodIsolation plugin.
func (p *PodIsolation) Name() string {
	return Name
}

// New creates a new instance of PodIsolation plugin.
// It initializes the PodIsolation with the provided context, runtime object, and framework handle.
// It returns the PodIsolation plugin and an error if any.
func New(ctx context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	kubeconfig := h.KubeConfig()
	kubeconfig.ContentType = "application/json"

	c, err := client.New(kubeconfig, client.Options{Scheme: scheme})

	if err != nil {
		return nil, err
	}

	klogFlags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(klogFlags)
	l := klog.NewKlogr().WithName("podisolation")

	return &PodIsolation{Client: c, handle: h, logger: l}, nil
}

// PreFilter is a method of the PodIsolation struct that implements the PreFilter interface of the Kubernetes scheduler framework.
// It is responsible for performing pre-filtering operations on a pod before it is scheduled to a node.
// The method takes the context, cycle state, and pod as input parameters and returns the pre-filter result and status.
// The pre-filtering operations include calculating the CPU requests of the pod, listing the node CPU topologies and pod CPU bindings,
// and determining the feasible CPUs for each node based on the exclusiveness level of the pod.
// The method also stores the pre-filtering state data in the cycle state for later use by other scheduling plugins.
func (p *PodIsolation) PreFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod) (*framework.PreFilterResult, *framework.Status) {
	// logger := p.logger.WithName("pre-filter").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
	topologies := &v1alpha1.NodeCpuTopologyList{}
	bindings := &v1alpha1.PodCpuBindingList{}
	nodes := sets.Set[string]{}

	exclusivenessLevel, ok := pod.Annotations[v1alpha1.AnnotationExclusivenessLevel]
	if !ok {
		exclusivenessLevel = "None"
	}

	podCpuRequests := int64(0)
	for _, container := range pod.Spec.Containers {
		podCpuRequests += container.Resources.Requests.Cpu().Value()
	}

	stateData := &PodIsolationStateData{
		ExclusivenessLevel: exclusivenessLevel,
		PodCpuRequests:     podCpuRequests,
		NodeCpuTopologies:  make(map[string]v1alpha1.NodeCpuTopology),
		PodCpuBindings:     make(map[string][]v1alpha1.PodCpuBinding),
		FeasibleCpus:       make(map[string]v1alpha1.CpuTopology),
	}

	// List NodeCpuTopologies and PodCpuBindings
	if err := p.List(ctx, topologies); err != nil {
		return nil, framework.NewStatus(framework.Error, fmt.Sprintf("%s/%s: could not list CPU topologies: %v", pod.Namespace, pod.Name, err))
	}

	if err := p.List(ctx, bindings); err != nil {
		return nil, framework.NewStatus(framework.Error, fmt.Sprintf("%s/%s: could not list CPU bindings: %v", pod.Namespace, pod.Name, err))
	}

	for _, topology := range topologies.Items {
		// Store topology in state and update NodeNames set for Filter plugin
		nodeName := topology.Spec.NodeName
		nodes = nodes.Insert(nodeName)
		stateData.NodeCpuTopologies[nodeName] = topology
		// nodeFeasibleCpus is initially p copy of the node's topology
		nodeFeasibleCpus := topology.Spec.Topology.DeepCopy()
		for _, binding := range bindings.Items {
			if !(binding.Status.ResourceStatus == v1alpha1.StatusApplied && binding.Status.NodeName == nodeName) {
				continue
			}
			// For each applied PodCpuBinding on current topology,
			// get all exclusive CPUs based on the ExclusivenessLevel
			// and exclude them from nodeFeasibleCpus
			stateData.PodCpuBindings[nodeName] = append(stateData.PodCpuBindings[nodeName], binding)
			// For every CPU binding, get all exclusive CPUs and remove them from the topology
			for exclusiveCpu := range pcbutils.GetExclusiveCpusOfCpuBinding(&binding, &topology) {
				// Delete CPU with key cpuId from nodeFeasibleCpus topology
				nctutils.DeleteCpuFromTopology(nodeFeasibleCpus, exclusiveCpu)
			}
		}
		// Store current node's feasible CPUs on cycle state
		stateData.FeasibleCpus[nodeName] = *nodeFeasibleCpus
	}
	state.Write(framework.StateKey(Name), stateData)
	return &framework.PreFilterResult{NodeNames: nodes}, nil
}

// Filter filters the given pod based on the available resources on the node.
// It checks if there are feasible CPUs available on the node to schedule the pod.
// If there are no feasible CPUs or there are not enough allocatable resources of
// the exclusiveness level type needed by the pod, it returns an Unschedulable status.
// Otherwise, it returns a Success status.
func (p *PodIsolation) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	logger := p.logger.WithName("filter").WithValues(fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), nodeInfo.Node().Name)
	node := nodeInfo.Node()

	// Read cycle state
	data, err := state.Read(framework.StateKey(Name))
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("%s/%s: could not read state data while filtering node %s: %v", pod.Namespace, pod.Name, node.Name, err))
	}

	stateData, ok := data.(*PodIsolationStateData)
	if !ok {
		return framework.NewStatus(framework.Error, fmt.Sprintf("%s/%s: could not cast state data while filtering node %s", pod.Namespace, pod.Name, node.Name))
	}

	feasibleCpus := stateData.FeasibleCpus[node.Name]
	topology := stateData.NodeCpuTopologies[node.Name].Spec.Topology
	exclusivenessLevel := stateData.ExclusivenessLevel
	totalCpus := int64(nctutils.GetTotalCpusCount(feasibleCpus))

	if totalCpus == 0 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("%s/%s: node has no unreserved CPUs", pod.Namespace, pod.Name))
	}

	availableResources := nctutils.GetAvailableResources(exclusivenessLevel, feasibleCpus, topology)

	// Check if there are enough allocatable resources of the exclusiveness level type needed by the pod
	res := 0
	for _, r := range availableResources {
		var cpus []int
		id := strconv.Itoa(r)
		switch exclusivenessLevel {
		case "Core":
			cpus = nctutils.GetAllCpusInCore(&feasibleCpus, id)
		case "Socket":
			cpus = nctutils.GetAllCpusInSocket(&feasibleCpus, id)
		case "Numa":
			cpus = nctutils.GetAllCpusInNuma(&feasibleCpus, id)
		default:
			res = int(totalCpus)
			break
		}
		res += len(cpus)
		if res >= int(stateData.PodCpuRequests) {
			break
		}
	}

	if res < int(stateData.PodCpuRequests) {
		p.logger.Info("Unschedulable", "resource", exclusivenessLevel, "available", availableResources, "requested", stateData.PodCpuRequests, "res", res)
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("%s/%s: no available resources found on node %s for resource %s", pod.Namespace, pod.Name, node.Name, exclusivenessLevel))
	}

	logger.Info("Schedulable", "resource", exclusivenessLevel, "available", availableResources)
	return framework.NewStatus(framework.Success)
}

// Score calculates the score for a given pod on a specific node based on the pod isolation criteria.
// It takes the context, cycle state, pod, and node name as input parameters.
// It returns the score as an int64 value and a framework.Status indicating the success or failure of the scoring process.
func (p *PodIsolation) Score(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (int64, *framework.Status) {
	logger := p.logger.WithName("score").WithValues(fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), nodeName)

	data, err := state.Read(framework.StateKey(Name))
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("%s/%s: could not read state data while scoring node %s: %v", pod.Namespace, pod.Name, nodeName, err))
	}
	stateData, ok := data.(*PodIsolationStateData)
	if !ok {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("%s/%s: could not cast state data while scoring node %s", pod.Namespace, pod.Name, nodeName))
	}

	feasibleCpus := stateData.FeasibleCpus[nodeName]
	topology := stateData.NodeCpuTopologies[nodeName].Spec.Topology
	exclusivenessLevel := stateData.ExclusivenessLevel

	availableResources := nctutils.GetAvailableResources(exclusivenessLevel, feasibleCpus, topology)
	score := int64(len(availableResources) * 100)

	logger.Info("scored", "score", score, "node", nodeName)
	return score, framework.NewStatus(framework.Success)
}

// NormalizeScore normalizes the scores of nodes based on the highest and lowest scores in the given list.
// It transforms the highest to lowest score range to fit the framework's minimum to maximum node score range.
// The normalized scores are updated in the provided scores list.
// The function returns a framework.Status indicating the success of the normalization process.
func (p *PodIsolation) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	logger := p.logger.WithName("normalize").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))

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

	for _, score := range scores {
		logger.Info("normalized", "node", score.Name, "score", score.Score)
	}

	return framework.NewStatus(framework.Success)
}

// PostBind is a method of the PodIsolation struct that is called after a pod is bound to a node.
// It assigns CPU resources to the pod based on the specified exclusiveness level and the available resources on the node.
// The method takes the context, cycle state, pod, and nodeName as parameters.
// It returns an error if there is any issue creating the PodCpuBinding object.
func (p *PodIsolation) PostBind(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) {
	logger := p.logger.WithName("post-bind").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", nodeName)

	data, err := state.Read(framework.StateKey(Name))
	if err != nil {
		logger.Error(err, "could not read state data while post-binding node")
		return
	}
	stateData, ok := data.(*PodIsolationStateData)
	if !ok {
		logger.Error(err, "could not cast state data while post-binding node")
		return
	}

	cpus := make([]int, 0)
	feasibleCpus := stateData.FeasibleCpus[nodeName]
	topology := stateData.NodeCpuTopologies[nodeName].Spec.Topology
	exclusivenessLevel := stateData.ExclusivenessLevel
	requestedCpus := stateData.PodCpuRequests
	availableResources := nctutils.GetAvailableResources(exclusivenessLevel, feasibleCpus, topology)

	for _, r := range availableResources {
		id := strconv.Itoa(r)
		switch exclusivenessLevel {
		case "Cpu":
			cpus = append(cpus, r)
			requestedCpus--
		case "Core":
			c := nctutils.GetAllCpusInCore(&feasibleCpus, id)
			cpus = append(cpus, c...)
			requestedCpus -= int64(len(c))
		case "Socket":
			c := nctutils.GetAllCpusInSocket(&feasibleCpus, id)
			cpus = append(cpus, c...)
			requestedCpus -= int64(len(c))
		case "Numa":
			c := nctutils.GetAllCpusInNuma(&feasibleCpus, id)
			cpus = append(cpus, c...)
			requestedCpus -= int64(len(c))
		}
		if requestedCpus <= 0 {
			break
		}
	}

	cpuSet := make([]v1alpha1.Cpu, 0)
	for i := 0; i < int(stateData.PodCpuRequests); i++ {
		cpuSet = append(cpuSet, v1alpha1.Cpu{CpuId: cpus[i]})
	}

	cpuBinding := &v1alpha1.PodCpuBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", pod.Name, nodeName),
			Namespace: pod.Namespace,
		},
		Spec: v1alpha1.PodCpuBindingSpec{
			PodName:            pod.Name,
			CpuSet:             cpuSet,
			ExclusivenessLevel: exclusivenessLevel,
		},
	}

	if err := p.Create(ctx, cpuBinding); err != nil {
		logger.Error(err, "could not create podcpubinding object")
		return
	}
}

func (p *PodIsolation) ScoreExtensions() framework.ScoreExtensions {
	return p
}

func (p *PodIsolation) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func (s *PodIsolationStateData) Clone() framework.StateData {
	return s
}
