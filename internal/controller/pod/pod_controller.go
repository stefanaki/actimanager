package pod

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var eventFilters = builder.WithPredicates(predicate.Funcs{
	CreateFunc: func(e event.CreateEvent) bool { return false },
	DeleteFunc: func(e event.DeleteEvent) bool { return false },
	UpdateFunc: func(e event.UpdateEvent) bool {
		oldPod, _ := e.ObjectNew.(*corev1.Pod)
		newPod, _ := e.ObjectOld.(*corev1.Pod)

		podCompleted := newPod.Status.Phase == corev1.PodSucceeded || oldPod.Status.Phase == corev1.PodFailed
		podDeleted := newPod.GetDeletionTimestamp() != oldPod.GetDeletionTimestamp()

		return podCompleted || podDeleted
	},
})

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// The controller is responsible for cleaning up PodCPUBindings when a Pod is deleted
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("pod-watcher")

	// Get deleted pod
	pod := &corev1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if controllerutil.ContainsFinalizer(pod, v1alpha1.FinalizerCPUBoundPod) {
		// Handle deleted pod
		cpuBindings := &v1alpha1.PodCPUBindingList{}
		if err := r.List(ctx, cpuBindings, &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector("spec.podName", req.NamespacedName.Name),
			Namespace:     pod.Namespace,
		}); err != nil {
			logger.Info("error listing cpu bindings", "error", err.Error())
			return ctrl.Result{}, fmt.Errorf("error listing CPU bindings: %v", err.Error())
		}

		for _, cpuBinding := range cpuBindings.Items {
			if err := r.Delete(ctx, &cpuBinding); err != nil {
				logger.Error(err, "could not delete cpu bindings for deleted pod", "pod", req.Name)
				return ctrl.Result{}, err
			}
		}

		controllerutil.RemoveFinalizer(pod, v1alpha1.FinalizerCPUBoundPod)
		if err := r.Update(ctx, pod); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1alpha1.PodCPUBinding{}, "spec.podName", func(rawObj client.Object) []string {
		cpuBinding := rawObj.(*v1alpha1.PodCPUBinding)
		return []string{cpuBinding.Spec.PodName}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}, eventFilters).
		Complete(r)
}
