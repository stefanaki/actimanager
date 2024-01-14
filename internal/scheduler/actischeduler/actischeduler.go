package actischeduler

import (
	"context"
	"fmt"

	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const Name string = "ActiScheduler"

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

type ActiScheduler struct {
	client.Client
	handle framework.Handle
}

var _ framework.FilterPlugin = &ActiScheduler{}
var _ framework.ScorePlugin = &ActiScheduler{}
var _ framework.ScoreExtensions = &ActiScheduler{}

func (a *ActiScheduler) Name() string {
	return Name
}

func New(ctx context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	c, err := client.New(h.KubeConfig(), client.Options{Scheme: scheme})

	if err != nil {
		return nil, err
	}

	return &ActiScheduler{Client: c, handle: h}, nil
}

func (a *ActiScheduler) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	node := nodeInfo.Node()
	topologies := &v1alpha1.NodeCpuTopologyList{}
	bindings := &v1alpha1.PodCpuBindingList{}

	err := a.List(ctx, topologies)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not list CPU topologies while fitering node %s: %v", pod.Namespace, pod.Name, node.Name, err))
	}

	err = a.List(ctx, bindings)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not list CPU bindings while fitering node %s: %v", pod.Namespace, pod.Name, node.Name, err))
	}

	var topology *v1alpha1.NodeCpuTopology = nil
	for _, nct := range topologies.Items {
		if nct.Spec.NodeName == node.Name {
			topology = &nct
			break
		}
	}

	if topology == nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("scheduling pod %s/%s: could not find CPU topology while fitering node %s: %v", pod.Namespace, pod.Name, node.Name, err))
	}

	fmt.Printf("scheduling pod %s/%s: found topology %s\n", pod.Namespace, pod.Name, topology.Name)

	var nodeBindings = make([]*v1alpha1.PodCpuBinding, 0)
	for _, pcb := range bindings.Items {
		if pcb.Status.NodeName == node.Name && pcb.Status.ResourceStatus == v1alpha1.StatusApplied {
			nodeBindings = append(nodeBindings, &pcb)
		}
	}

	for _, b := range nodeBindings {
		fmt.Printf("found node binding: %s/%s/%s\n", b.Status.NodeName, b.Namespace, b.Name)
	}

	return framework.NewStatus(framework.Success)
}

func (a *ActiScheduler) Score(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) (int64, *framework.Status) {
	println("I THINK ITS WORKING1...")
	return 0, &framework.Status{}
}

func (a *ActiScheduler) NormalizeScore(ctx context.Context, state *framework.CycleState, p *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	return &framework.Status{}
}

func (a *ActiScheduler) ScoreExtensions() framework.ScoreExtensions {
	return a
}
