package podcpubinding

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
)

// PodCpuBindingReconciler reconciles a PodCpuBinding object
type PodCpuBindingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=podcpubindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=podcpubindings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=podcpubindings/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PodCpuBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("pcb-controller")

	// Get PodCpuBinding CR
	cpuBinding := &cslabecentuagrv1alpha1.PodCpuBinding{}

	// Handle delete
	err := r.Get(ctx, req.NamespacedName, cpuBinding)
	if errors.IsNotFound(err) {
		logger.Info("Deleted PodCpuBinding")
		// TODO Add logic for notifying node agent
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error reconciling PodCpuBinding: %v", err.Error())
	}

	// Initialize CR

	// Validate CR

	// Handle reconcilation

	// TODO Notify node agent for binding creation
	// wait for response and update CR status

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodCpuBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cslabecentuagrv1alpha1.PodCpuBinding{}).
		Complete(r)
}
