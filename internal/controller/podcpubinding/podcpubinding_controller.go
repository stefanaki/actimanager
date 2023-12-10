package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nodecputopologyv1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology/v1alpha1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
	cpuBinding := &v1alpha1.PodCpuBinding{}

	// Handle delete
	err := r.Get(ctx, req.NamespacedName, cpuBinding)
	if errors.IsNotFound(err) {
		logger.Info("Deleted PodCpuBinding")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error reconciling PodCpuBinding: %v", err.Error())
	}

	// Initialize CR
	if needsUpdate(cpuBinding) {
		cpuBinding.Status.Status = v1alpha1.StatusBindingPending
		cpuBinding.Status.LastSpec = cpuBinding.Spec
		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
		}

		return ctrl.Result{}, nil
	}

	// Validate CR
	pod := &corev1.Pod{}
	err = r.Get(ctx, client.ObjectKey{Name: cpuBinding.Spec.PodName, Namespace: cpuBinding.ObjectMeta.Namespace}, pod)

	if errors.IsNotFound(err) {
		cpuBinding.Status.Status = v1alpha1.StatusPodNotFound
		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating cpu binding status: %v", err.Error())
		}

		return ctrl.Result{}, nil
	}

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting pod: %v", err.Error())
	}

	// Check if all containers are ready
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// Get NodeCpuTopology of node
	topologies := &v1alpha1.NodeCpuTopologyList{}
	err = r.List(ctx,
		topologies,
		client.MatchingFields{"spec.nodeName": pod.Spec.NodeName})

	if errors.IsNotFound(err) {
		cpuBinding.Status.Status = v1alpha1.StatusNodeTopologyNotFound
		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating cpu binding status: %v", err.Error())
		}
	}

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing CPU topologies: %v", err.Error())
	}

	targetTopology := topologies.Items[0]

	// Check if specified cpuset is available in the node topology
	if !nodecputopologyv1alpha1.IsCpuSetInTopology(&targetTopology.Spec.Topology, cpuBinding.Spec.CpuSet) {
		cpuBinding.Status.Status = v1alpha1.StatusInvalidCpuSet

		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating cpu binding status: %v", err.Error())
		}

		return ctrl.Result{}, nil
	}

	// Handle reconcilation
	switch cpuBinding.Status.Status {
	case v1alpha1.StatusBindingPending:
		// Apply CPU pinning
		err = r.applyCpuPinning(ctx, cpuBinding.Spec.CpuSet, pod, logger)
		if err != nil {
			cpuBinding.Status.Status = v1alpha1.StatusFailed
		}
		cpuBinding.Status.Status = v1alpha1.StatusApplied

		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
		}

	case v1alpha1.StatusApplied:
		if needsDelete(cpuBinding) {
			// Remove CPU pinning and delete CR
			err := r.removeCpuPinning(ctx, pod, logger)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("error removing CPU pinning: %v", err.Error())
			}

			err = r.Delete(ctx, cpuBinding)

			if err != nil {
				return ctrl.Result{}, fmt.Errorf("error deleting PodCpuBinding: %v", err.Error())
			}

			return ctrl.Result{}, nil
		}
	default:
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodCpuBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1alpha1.NodeCpuTopology{}, "spec.nodeName", func(rawObj client.Object) []string {
		topology := rawObj.(*v1alpha1.NodeCpuTopology)
		return []string{topology.Spec.NodeName}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PodCpuBinding{}).
		Complete(r)
}
