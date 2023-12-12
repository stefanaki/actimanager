package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	if cpuBinding.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(cpuBinding, v1alpha1.FinalizerPodCpuBinding) {
			controllerutil.AddFinalizer(cpuBinding, v1alpha1.FinalizerPodCpuBinding)

			if err := r.Update(ctx, cpuBinding); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %v", err.Error())
			}

		}
	} else {
		if controllerutil.ContainsFinalizer(cpuBinding, v1alpha1.FinalizerPodCpuBinding) {
			if err := r.PodCpuBindingFinalizer(ctx, cpuBinding, logger); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(cpuBinding, v1alpha1.FinalizerPodCpuBinding)
			if err := r.Update(ctx, cpuBinding); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	// Initialize CR
	if needsReconciliation(cpuBinding) {
		if cpuBinding.Status.Status == v1alpha1.StatusApplied {
			// TODO Remove CPU pinning and set status to Pending
		}

		cpuBinding.Status.Status = v1alpha1.StatusBindingPending
		cpuBinding.Status.LastSpec = cpuBinding.Spec
		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
		}

		return ctrl.Result{}, nil
	}

	if cpuBinding.Status.Status == v1alpha1.StatusApplied {
		return ctrl.Result{}, nil
	}

	if cpuBinding.Status.Status == v1alpha1.StatusInvalidCpuSet ||
		cpuBinding.Status.Status == v1alpha1.StatusPodNotFound ||
		cpuBinding.Status.Status == v1alpha1.StatusNodeTopologyNotFound ||
		cpuBinding.Status.Status == v1alpha1.StatusFailed ||
		cpuBinding.Status.Status == v1alpha1.StatusCpuSetAllocationFailed {
		return ctrl.Result{}, nil
	}

	// Validate Pod name
	pod := &corev1.Pod{}
	ok, status, err := r.validatePodName(ctx, cpuBinding, pod)

	if !ok {
		if err != nil {
			return ctrl.Result{}, err
		}

		if status != "" {
			cpuBinding.Status.Status = status

			if err := r.Status().Update(ctx, cpuBinding); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
			}

			return ctrl.Result{}, nil
		}
	}

	// Validate Topology
	topology := &v1alpha1.NodeCpuTopology{}
	ok, status, err = r.validateTopology(ctx, cpuBinding, topology, pod)

	if !ok {
		if err != nil {
			return ctrl.Result{}, err
		}

		if status != "" {
			cpuBinding.Status.Status = status

			if err := r.Status().Update(ctx, cpuBinding); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
			}

			return ctrl.Result{}, nil
		}
	}

	ok, status, err = r.validateExclusivenessLevel(ctx, cpuBinding, topology, req.NamespacedName)
	if !ok {
		if err != nil {
			return ctrl.Result{}, err
		}

		if status != "" {
			cpuBinding.Status.Status = status

			if err := r.Status().Update(ctx, cpuBinding); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
			}

			return ctrl.Result{}, nil
		}
	}

	// Check if all containers are ready
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// Handle reconcilation
	// Apply CPU pinning
	err = r.applyCpuPinning(ctx, cpuBinding.Spec.CpuSet, pod, logger)
	if err != nil {
		cpuBinding.Status.Status = v1alpha1.StatusFailed
	}
	cpuBinding.Status.Status = v1alpha1.StatusApplied
	cpuBinding.Status.NodeName = pod.Spec.NodeName

	if err := r.Status().Update(ctx, cpuBinding); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
	}

	return ctrl.Result{}, nil
}

func needsReconciliation(cpuBinding *v1alpha1.PodCpuBinding) bool {
	return !reflect.DeepEqual(cpuBinding.Spec, cpuBinding.Status.LastSpec)
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
