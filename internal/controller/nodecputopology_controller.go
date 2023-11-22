package controller

import (
	"context"
	"encoding/json"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
)

// NodeCpuTopologyReconciler reconciles a NodeCpuTopology object
type NodeCpuTopologyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/finalizers,verbs=update

func (r *NodeCpuTopologyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("controller")
	// nodeName := "minikube"

	topology := &cslabecentuagrv1alpha1.NodeCpuTopology{}
	if err := r.Get(ctx, req.NamespacedName, topology); err != nil {
		logger.V(5).Info("Error listing NodeCpuTopology resources:" + err.Error())
		return ctrl.Result{}, err
	}

	data, _ := json.Marshal(topology)
	logger.V(1).Info(string(data))
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeCpuTopologyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cslabecentuagrv1alpha1.NodeCpuTopology{}).
		Complete(r)
}
