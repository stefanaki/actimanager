package cpuaware

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const Name string = "CpuAware"

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(batchv1.AddToScheme(scheme))
	utilruntime.Must(cslabecentuagrv1alpha1.AddToScheme(scheme))
}

type CpuAware struct {
	handle framework.Handle
	client client.Client
}

func (c *CpuAware) Name() string {
	return Name
}

func (c *CpuAware) Score(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) (int64, *framework.Status) {
	println("I THINK ITS WORKING...")
	return 0, &framework.Status{}
}

func (c *CpuAware) NormalizeScore(ctx context.Context, state *framework.CycleState, p *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	return &framework.Status{}
}

func (c *CpuAware) ScoreExtensions() framework.ScoreExtensions {
	return c
}

func New(ctx context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	c, err := client.New(h.KubeConfig(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	println("cpu aware plugin initialized")

	test := &cslabecentuagrv1alpha1.NodeCpuTopologyList{}
	errr := c.List(ctx, test)

	if errr != nil {
		println("FATAL")
	} else {
		for _, nct := range test.Items {
			println(nct.Name, nct.Spec.NodeName)
		}
	}

	return &CpuAware{handle: h, client: c}, nil
}

var _ = framework.ScorePlugin(&CpuAware{})
