package actischeduler

import (
	"context"
	"flag"
	"fmt"
	"math"

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
	c, err := client.New(h.KubeConfig(), client.Options{Scheme: scheme})

	if err != nil {
		return nil, err
	}

	klogFlags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(klogFlags)
	l := klog.NewKlogr().WithName("actischeduler")

	return &ActiScheduler{Client: c, handle: h, logger: l}, nil
}

// Filter filters the available CPUs on a node based on the CPU topology and CPU bindings.
// It checks the CPU topology and the CPU bindings of the node to determine the feasible CPUs.
// Feasible CPUs are those that are within the CPU topology of the node and are not already exclusively bound to other pods.
// If there are no feasible CPUs, it returns an Unschedulable status.
// If the requested CPU resources of the pod exceed the number of feasible CPUs, it returns an Unschedulable status.
// Otherwise, it returns a Success status.
func (a *ActiScheduler) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	logger := a.logger.WithName("filter")
	node := nodeInfo.Node()
	topologies := &v1alpha1.NodeCpuTopologyList{}
	bindings := &v1alpha1.PodCpuBindingList{}

	// List NodeCpuTopologies and PodCpuBindings
	if err := a.List(ctx, topologies); err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not list CPU topologies while filtering node %s: %v", pod.Namespace, pod.Name, node.Name, err))
	}

	if err := a.List(ctx, bindings); err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not list CPU bindings while filtering node %s: %v", pod.Namespace, pod.Name, node.Name, err))
	}

	// Find the topology of the node
	var topology *v1alpha1.NodeCpuTopology
	for _, n := range topologies.Items {
		if n.Spec.NodeName == node.Name {
			topology = &n
			break
		}
	}

	if topology == nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not find CPU topology while filtering node %s", pod.Namespace, pod.Name, node.Name))
	}

	logger.Info("scheduling pod %s/%s: found topology %s\n", pod.Namespace, pod.Name, topology.Name)

	// Find the bindings of the node
	var nodeBindings []*v1alpha1.PodCpuBinding
	for _, pcb := range bindings.Items {
		if pcb.Status.NodeName == node.Name && pcb.Status.ResourceStatus == v1alpha1.StatusApplied {
			nodeBindings = append(nodeBindings, &pcb)
		}
	}

	feasibleCpus := make(map[int]struct{})
	for _, socket := range topology.Spec.Topology.Sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.Cpus {
				feasibleCpus[cpu.CpuId] = struct{}{}
			}
		}
	}

	for _, b := range nodeBindings {
		for _, cpu := range b.Spec.CpuSet {
			_, coreId, socketId, numaId := nct.GetCpuParentInfo(topology, cpu.CpuId)

			switch b.Spec.ExclusivenessLevel {
			case "Cpu":
				delete(feasibleCpus, cpu.CpuId)
			case "Core":
				for _, cpu := range nct.GetAllCpusInCore(topology, coreId) {
					delete(feasibleCpus, cpu)
				}
			case "Socket":
				for _, cpu := range nct.GetAllCpusInSocket(topology, socketId) {
					delete(feasibleCpus, cpu)
				}
			case "Numa":
				for _, cpu := range nct.GetAllCpusInNuma(topology, numaId) {
					delete(feasibleCpus, cpu)
				}
			}
		}
	}

	if len(feasibleCpus) == 0 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("scheduling pod %s/%s: no feasible CPUs found on node %s", pod.Namespace, pod.Name, node.Name))
	}

	podCpuRequestsMilli := int64(0)
	for _, container := range pod.Spec.Containers {
		podCpuRequestsMilli += container.Resources.Requests.Cpu().MilliValue()
	}

	if podCpuRequestsMilli > int64(len(feasibleCpus))*1000 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("scheduling pod %s/%s: not enough feasible CPUs found on node %s", pod.Namespace, pod.Name, node.Name))
	}

	return framework.NewStatus(framework.Success)
}

// Score calculates the score for a pod on a specific node based on the locality of the feasible CPUs.
func (a *ActiScheduler) Score(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) (int64, *framework.Status) {
	return 0, &framework.Status{}
}

func (a *ActiScheduler) NormalizeScore(ctx context.Context, state *framework.CycleState, p *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
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

	return framework.NewStatus(framework.Success)
}

func (a *ActiScheduler) PostBind(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) {

}

func (a *ActiScheduler) ScoreExtensions() framework.ScoreExtensions {
	return a
}
