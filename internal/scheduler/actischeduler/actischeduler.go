package actischeduler

import (
	"context"
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
	// node := nodeInfo.Node()
	cpuTopologies := &v1alpha1.NodeCpuTopologyList{}
	podCpuBindings := &v1alpha1.PodCpuBindingList{}

	err := a.List(ctx, cpuTopologies)
	if err != nil {
		println("handle error later")
		println(err.Error())
	}

	err = a.List(ctx, podCpuBindings)
	if err != nil {
		println("handle error later")
		println(err.Error())
	}

	return &framework.Status{}
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
