package actischeduler

import (
	"context"
	"flag"
	"fmt"
	"math"

	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/podcpubinding"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nct "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ActiScheduler is responsible for filtering and scoring nodes based on CPU topology and CPU bindings.
// It ensures that pods are scheduled on nodes with feasible CPUs and calculates the score based on the locality of the feasible CPUs.
// This plugin is used to optimize CPU resource allocation in a Kubernetes cluster.

// Name is the name of the ActiScheduler plugin.
const Name string = "ActiScheduler"

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

// ActiScheduler is the implementation of the ActiScheduler plugin.
// It embeds the Kubernetes client, framework handle, and logger.
type ActiScheduler struct {
	client.Client
	handle framework.Handle
	logger klog.Logger
}

// ActiSchedulerStateData represents the state data for the ActiScheduler.
type ActiSchedulerStateData struct {
	// NodeCpuTopologies stores the CPU topology of each node in the cluster.
	NodeCpuTopologies map[string]v1alpha1.NodeCpuTopology

	// PodCpuBindings stores the CPU bindings for each node in the cluster.
	PodCpuBindings map[string][]v1alpha1.PodCpuBinding

	// FeasibleCpus stores the feasible CPUs for each node in the cluster.
	FeasibleCpus map[string]v1alpha1.CpuTopology
}

var _ framework.PreFilterPlugin = &ActiScheduler{}
var _ framework.FilterPlugin = &ActiScheduler{}
var _ framework.ScorePlugin = &ActiScheduler{}
var _ framework.ScoreExtensions = &ActiScheduler{}
var _ framework.PostBindPlugin = &ActiScheduler{}

// Name returns the name of the ActiScheduler plugin.
func (a *ActiScheduler) Name() string {
	return Name
}

// New creates a new instance of ActiScheduler plugin.
// It initializes the ActiScheduler with the provided context, runtime object, and framework handle.
// It returns the ActiScheduler plugin and an error if any.
func New(ctx context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	kubeconfig := h.KubeConfig()
	kubeconfig.ContentType = "application/json"

	c, err := client.New(kubeconfig, client.Options{Scheme: scheme})

	if err != nil {
		return nil, err
	}

	klogFlags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(klogFlags)
	l := klog.NewKlogr().WithName("actischeduler")

	return &ActiScheduler{Client: c, handle: h, logger: l}, nil
}

// PreFilter is responsible for performing pre-filtering operations on a pod before filtering the nodes.
// It lists the NodeCpuTopologies and PodCpuBindings, and populates the state data with the feasible CPUs for each node.
// It returns the pre-filter result with the set of node names.
func (a *ActiScheduler) PreFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod) (*framework.PreFilterResult, *framework.Status) {
	// logger := a.logger.WithName("pre-filter").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
	topologies := &v1alpha1.NodeCpuTopologyList{}
	bindings := &v1alpha1.PodCpuBindingList{}
	nodes := sets.Set[string]{}

	stateData := &ActiSchedulerStateData{
		NodeCpuTopologies: make(map[string]v1alpha1.NodeCpuTopology),
		PodCpuBindings:    make(map[string][]v1alpha1.PodCpuBinding),
		FeasibleCpus:      make(map[string]v1alpha1.CpuTopology),
	}

	// List NodeCpuTopologies and PodCpuBindings
	if err := a.List(ctx, topologies); err != nil {
		return nil, framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not list CPU topologies: %v", pod.Namespace, pod.Name, err))
	}

	if err := a.List(ctx, bindings); err != nil {
		return nil, framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not list CPU bindings: %v", pod.Namespace, pod.Name, err))
	}

	for _, topology := range topologies.Items {
		// Store topology in state and update NodeNames set for Filter plugin
		nodeName := topology.Spec.NodeName
		nodes = nodes.Insert(nodeName)
		stateData.NodeCpuTopologies[nodeName] = topology

		// nodeFeasibleCpus is initially a copy of the node's topology
		nodeFeasibleCpus := topology.Spec.Topology
		for _, binding := range bindings.Items {
			// For each applied PodCpuBinding on current topology,
			// get all exclusive CPUs based on the ExclusivenessLevel
			// and exclude them from nodeFeasibleCpus
			if binding.Status.ResourceStatus == v1alpha1.StatusApplied && binding.Status.NodeName == nodeName {
				stateData.PodCpuBindings[nodeName] = append(stateData.PodCpuBindings[nodeName], binding)

				// For every CPU binding, get all exclusive CPUs and remove them from the topology
				for _, binding := range bindings.Items {
					for exclusiveCpu := range pcbutils.GetExclusiveCpusOfCpuBinding(&binding, &topology) {
						cpuId, coreId, socketId, numaId := nct.GetCpuParentInfo(&topology, exclusiveCpu)
						// Delete CPU with key cpuId from nodeFeasibleCpus topology
						delete(nodeFeasibleCpus.Sockets[socketId].Cores[coreId].Cpus, cpuId)
						delete(nodeFeasibleCpus.NumaNodes[numaId].Cpus, cpuId)
					}
				}
			}
		}

		// Store current node's feasible CPUs on cycle state
		stateData.FeasibleCpus[nodeName] = nodeFeasibleCpus
	}

	state.Write(framework.StateKey(Name), stateData)
	return &framework.PreFilterResult{NodeNames: nodes}, nil
}

// Filter filters the given pod based on the available resources on the node.
// It checks if there are feasible CPUs available on the node to schedule the pod.
// If there are no feasible CPUs or if the requested CPU resources exceed the available feasible CPUs,
// it returns an error status indicating that the pod is unschedulable.
// Otherwise, it returns a Success status.
func (a *ActiScheduler) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	logger := a.logger.WithName("filter").WithValues(fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), nodeInfo.Node().Name)
	node := nodeInfo.Node()

	// Read cycle state
	data, err := state.Read(framework.StateKey(Name))
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not read state data while filtering node %s: %v", pod.Namespace, pod.Name, node.Name, err))
	}

	stateData, ok := data.(*ActiSchedulerStateData)
	if !ok {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not cast state data while filtering node %s", pod.Namespace, pod.Name, node.Name))
	}

	feasibleCpus := stateData.FeasibleCpus[node.Name]
	totalCpus := int64(nct.GetTotalCpusCount(feasibleCpus))
	if totalCpus == 0 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("scheduling pod %s/%s: no feasible CPUs found on node %s", pod.Namespace, pod.Name, node.Name))
	}

	// Get pod's CPU requests and compare them to the feasible CPUs of current node
	podCpuRequests := int64(0)
	for _, container := range pod.Spec.Containers {
		podCpuRequests += container.Resources.Requests.Cpu().Value()
	}

	// If pod's CPU requests exceed the number of available CPUs
	// mark pod as Unschedulable on current node
	if podCpuRequests > totalCpus {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("scheduling pod %s/%s: not enough feasible CPUs found on node %s", pod.Namespace, pod.Name, node.Name))
	}

	logger.Info("Schedulable", "totalCpus", totalCpus, "pod reqs", podCpuRequests)

	return framework.NewStatus(framework.Success)
}

// Score calculates the score for a pod on a specific node based on the locality of the feasible CPUs.
func (a *ActiScheduler) Score(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (int64, *framework.Status) {
	logger := a.logger.WithName("score").WithValues(fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), nodeName)

	data, err := state.Read(framework.StateKey(Name))
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not read state data while scoring node %s: %v", pod.Namespace, pod.Name, nodeName, err))
	}

	stateData, ok := data.(*ActiSchedulerStateData)
	if !ok {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not cast state data while scoring node %s", pod.Namespace, pod.Name, nodeName))
	}

	// feasibleCpus is a v1alpha1.CpuTopology that contains `only` the CPUs
	// that are not exclusive to any pod
	feasibleCpus := stateData.FeasibleCpus[nodeName]
	score := int64(0)

	podCpuRequests := int64(0)
	for _, container := range pod.Spec.Containers {
		podCpuRequests += container.Resources.Requests.Cpu().Value()
	}

	for socketId, socket := range feasibleCpus.Sockets {
		for _, core := range socket.Cores {
			if int64(len(core.Cpus)) >= podCpuRequests {
				score += 10000
			}
		}

		if int64(len(nct.GetAllCpusInSocket(&feasibleCpus, socketId))) >= podCpuRequests {
			score += 10
		}
	}

	for numaId := range feasibleCpus.NumaNodes {
		if int64(len(nct.GetAllCpusInNuma(&feasibleCpus, numaId))) >= podCpuRequests {
			score += 100
		}
	}

	logger.Info("score", "score", score)

	return score, framework.NewStatus(framework.Success)
}

// NormalizeScore normalizes the scores of nodes based on the highest and lowest scores in the given list.
// It transforms the highest to lowest score range to fit the framework's minimum to maximum node score range.
// The normalized scores are updated in the provided scores list.
// The function returns a framework.Status indicating the success of the normalization process.
func (a *ActiScheduler) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	logger := a.logger.WithName("normalize").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))

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

func (a *ActiScheduler) PostBind(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) {
	logger := a.logger.WithName("post-bind").WithValues("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "node", nodeName)

	data, err := state.Read(framework.StateKey(Name))
	if err != nil {
		logger.Error(err, "could not read state data while post-binding node")
		return
	}

	stateData, ok := data.(*ActiSchedulerStateData)
	if !ok {
		logger.Error(err, "could not cast state data while post-binding node")
		return
	}

	feasibleCpus := stateData.FeasibleCpus[nodeName]
	podCpuRequests := int64(0)
	for _, container := range pod.Spec.Containers {
		podCpuRequests += container.Resources.Requests.Cpu().Value()
	}

	cpuSet := make([]v1alpha1.Cpu, 0)
	for _, socket := range feasibleCpus.Sockets {
		for _, core := range socket.Cores {
			for cpu := range core.Cpus {
				cpuSet = append(cpuSet, core.Cpus[cpu])
				break
			}
		}
	}

	// if you dont find a core with the needed cpu reqeuests
	// then you need to find a socket with the needed cpu requests
	// and then a numa node with the needed cpu requests

	if len(cpuSet) == 0 {
		for socketId, socket := range feasibleCpus.Sockets {
			if int64(len(nct.GetAllCpusInSocket(&feasibleCpus, socketId))) >= podCpuRequests {
				for _, core := range socket.Cores { // Fix: Iterate over the cores of the socket
					for cpu := range core.Cpus { // Fix: Iterate over the CPUs of the core
						cpuSet = append(cpuSet, core.Cpus[cpu])
						break
					}
				}
			}
		}
	}

	if len(cpuSet) == 0 {
		for numaId := range feasibleCpus.NumaNodes {
			if int64(len(nct.GetAllCpusInNuma(&feasibleCpus, numaId))) >= podCpuRequests {
				for cpu := range feasibleCpus.NumaNodes[numaId].Cpus {
					cpuSet = append(cpuSet, feasibleCpus.NumaNodes[numaId].Cpus[cpu])
					break
				}
			}
		}
	}

	// create the podcpubinding object
	cpuBinding := &v1alpha1.PodCpuBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-binding", pod.Name),
			Namespace: pod.Namespace,
		},
		Spec: v1alpha1.PodCpuBindingSpec{
			PodName:            pod.Name,
			CpuSet:             []v1alpha1.Cpu{{CpuId: 15}},
			ExclusivenessLevel: "Cpu",
		},
	}

	//logger.Info("cpuset for new pod", "cpus", cpuSet)

	// save the podcpubinding object to the api server
	if err := a.Create(ctx, cpuBinding); err != nil {
		logger.Error(err, "could not create podcpubinding object")
		return
	}
}

func (a *ActiScheduler) ScoreExtensions() framework.ScoreExtensions {
	return a
}

func (a *ActiScheduler) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func (s *ActiSchedulerStateData) Clone() framework.StateData {
	return s
}
