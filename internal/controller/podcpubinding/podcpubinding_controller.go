package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	nctutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/podcpubinding"
	"fmt"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// PodCpuBindingReconciler reconciles a PodCpuBinding object
type PodCpuBindingReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var eventFilters = builder.WithPredicates(predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		oldObj := e.ObjectOld.(*v1alpha1.PodCpuBinding)
		newObj := e.ObjectNew.(*v1alpha1.PodCpuBinding)

		specChanged := !reflect.DeepEqual(oldObj.Spec, newObj.Spec)
		statusBindingPending := newObj.Status.ResourceStatus == v1alpha1.StatusBindingPending
		isDeleted := !newObj.DeletionTimestamp.IsZero()

		return specChanged || statusBindingPending || isDeleted
	},
})

// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=podcpubindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=podcpubindings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=podcpubindings/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PodCpuBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("pcb-controller")

	// Get PodCpuBinding CR
	cpuBinding := &v1alpha1.PodCpuBinding{}

	// Handle delete
	err := r.Get(ctx, req.NamespacedName, cpuBinding)
	if errors.IsNotFound(err) {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// If CR is not deleted, add the finalizer
	if cpuBinding.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(cpuBinding, v1alpha1.FinalizerPodCpuBinding) {
			controllerutil.AddFinalizer(cpuBinding, v1alpha1.FinalizerPodCpuBinding)
			if err := r.Update(ctx, cpuBinding); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %v", err.Error())
			}
		}
	} else {
		// If CR is deleted and contains the finalizer,
		// execute the finalizer that removes the CPU pinning
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

	// If CR previously had status `Applied`, remove the CPU pinning first before reconciling
	if cpuBinding.Status.ResourceStatus == v1alpha1.StatusApplied {
		pod, err := r.getPod(ctx, types.NamespacedName{
			Namespace: cpuBinding.Namespace,
			Name:      cpuBinding.Spec.PodName,
		})
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error getting pod: %v", err.Error())
		}
		if err := r.removeCpuPinning(ctx, pod); err != nil {
			return ctrl.Result{}, fmt.Errorf("error removing cpu pinning: %v", err.Error())
		}
		r.Recorder.Eventf(cpuBinding, corev1.EventTypeNormal, "CpuPinningRemoved", "CPU pinning removed from pod %s", cpuBinding.Spec.PodName)
	}

	// Initialize CR
	if cpuBinding.Status.ResourceStatus == "" ||
		cpuBinding.Status.ResourceStatus == v1alpha1.StatusInvalidCpuSet ||
		cpuBinding.Status.ResourceStatus == v1alpha1.StatusPodNotFound ||
		cpuBinding.Status.ResourceStatus == v1alpha1.StatusNodeTopologyNotFound ||
		cpuBinding.Status.ResourceStatus == v1alpha1.StatusFailed ||
		cpuBinding.Status.ResourceStatus == v1alpha1.StatusCpuSetAllocationFailed {

		cpuBinding.Status.ResourceStatus = v1alpha1.StatusBindingPending
		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
		}
		r.Recorder.Eventf(cpuBinding, corev1.EventTypeNormal, string(v1alpha1.StatusBindingPending), "CPU binding for pod %s/%s is pending", cpuBinding.Namespace, cpuBinding.Spec.PodName)
		return ctrl.Result{}, nil
	}

	// Validate Resource
	pod := &corev1.Pod{}
	topology := &v1alpha1.NodeCpuTopology{}
	ok, status, message, err := r.validateResource(ctx, cpuBinding, topology, pod)
	if !ok {
		if err != nil {
			return ctrl.Result{}, err
		}
		if status != "" {
			cpuBinding.Status.ResourceStatus = status
			if err := r.Status().Update(ctx, cpuBinding); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
			}
			r.Recorder.Event(cpuBinding, corev1.EventTypeWarning, string(status), message)
			return ctrl.Result{}, nil
		}
	}

	// Check if all containers are ready
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 2}, nil
		}
	}

	// Handle reconcilation
	// Apply CPU pinning
	cpuSet := pcbutils.ConvertCpuSliceToIntSlice(cpuBinding.Spec.CpuSet)
	memSet := nctutils.GetNumaNodesOfCpuSet(cpuSet, topology.Spec.Topology)
	err = r.applyCpuPinning(ctx, cpuSet, memSet, pod)
	if err != nil {
		r.Recorder.Eventf(cpuBinding, corev1.EventTypeWarning, string(v1alpha1.StatusFailed), "Failed to apply CPU pinning to pod %s: %v", cpuBinding.Spec.PodName, err)
		cpuBinding.Status.ResourceStatus = v1alpha1.StatusFailed
	}
	cpuBinding.Status.ResourceStatus = v1alpha1.StatusApplied
	cpuBinding.Status.NodeName = pod.Spec.NodeName

	if err := r.Status().Update(ctx, cpuBinding); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating status: %v", err.Error())
	}

	r.Recorder.Eventf(cpuBinding, corev1.EventTypeNormal, string(v1alpha1.StatusApplied), "CPUSet %v, MemSet %v applied to pod %s/%s", cpuSet, memSet, pod.Namespace, pod.Name)
	controllerutil.AddFinalizer(pod, v1alpha1.FinalizerCpuBoundPod)
	if err := r.Update(ctx, pod); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to add finalizer to pod: %v", err.Error())
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

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1alpha1.PodCpuBinding{}, "status.nodeName", func(rawObj client.Object) []string {
		podCpuBinding := rawObj.(*v1alpha1.PodCpuBinding)
		return []string{podCpuBinding.Status.NodeName}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1alpha1.PodCpuBinding{}, "status.resourceStatus", func(rawObj client.Object) []string {
		podCpuBinding := rawObj.(*v1alpha1.PodCpuBinding)
		return []string{string(podCpuBinding.Status.ResourceStatus)}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PodCpuBinding{}, eventFilters).
		WithOptions(controller.Options{MaxConcurrentReconciles: 3}).
		Complete(r)
}
