package nodecputopology

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// NodeCpuTopologyReconciler reconciles a NodeCpuTopology object
type NodeCpuTopologyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var eventFilters = builder.WithPredicates(predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		oldObj := e.ObjectOld.(*v1alpha1.NodeCpuTopology)
		newObj := e.ObjectNew.(*v1alpha1.NodeCpuTopology)

		nodeNameChanged := oldObj.Spec.NodeName != newObj.Spec.NodeName
		statusNeedsSync := newObj.Status.ResourceStatus == v1alpha1.StatusNeedsSync

		return nodeNameChanged || statusNeedsSync
	},
})

// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=pods/log,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *NodeCpuTopologyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("nct-controller")

	// Get NodeCpuTopology CR
	topology := &v1alpha1.NodeCpuTopology{}

	// Handle delete
	err := r.Get(ctx, req.NamespacedName, topology)
	if errors.IsNotFound(err) {
		logger.Info("Deleted NodeCpuTopology")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Initialize CR
	if topology.Status.ResourceStatus == "" || topology.Status.ResourceStatus == v1alpha1.StatusNodeNotFound {
		topology.Status.ResourceStatus = v1alpha1.StatusNeedsSync
		topology.Status.InitJobStatus = v1alpha1.StatusJobNone
		if err := r.Status().Update(ctx, topology); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating resource: %v", err)
		}
		return ctrl.Result{}, nil
	}

	// Check if specified NodeName is a valid name of a node
	if err := r.Get(ctx, client.ObjectKey{Name: topology.Spec.NodeName}, &corev1.Node{}); err != nil {
		logger.Info("Node not found", "nodeName", topology.Spec.NodeName)
		topology.Status.ResourceStatus = v1alpha1.StatusNodeNotFound
		topology.Status.InitJobStatus = v1alpha1.StatusJobNone

		if err := r.Status().Update(ctx, topology); err != nil {
			return ctrl.Result{}, fmt.Errorf("could not update status: %v", err)
		}

		return ctrl.Result{}, nil
	}

	// Handle reconcilation
	switch topology.Status.ResourceStatus {
	case v1alpha1.StatusNeedsSync:
		// If ResourceStatus is empty or NeedsSync, initiate job
		switch topology.Status.InitJobStatus {
		case v1alpha1.StatusJobNone:
			logger.Info("Dispatch init job for NodeCpuTopology")

			jobName, err := r.createInitJob(topology, ctx, &logger)

			topology.Status.InitJobStatus = v1alpha1.StatusJobPending
			topology.Status.InitJobName = jobName

			if err = r.Status().Update(ctx, topology); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating resource: %v", err)
			}

			return ctrl.Result{}, nil
		case v1alpha1.StatusJobPending:
			// While InitJobStatus is Pending, requeue the CR until the pod completes
			isCompleted, err := r.isJobCompleted(topology, ctx)
			if err != nil {
				return ctrl.Result{}, err
			}

			if !isCompleted {
				return ctrl.Result{Requeue: true}, nil
			}

			cpuTopology, err := r.parseCompletedPod(topology, ctx, &logger)

			if err != nil {
				return ctrl.Result{}, fmt.Errorf("error getting cpu topology: %v", err)
			}

			topology.Spec.Topology = cpuTopology
			if err := r.Update(ctx, topology); err != nil {
				return ctrl.Result{Requeue: true}, fmt.Errorf("error updating NodeCpuTopology spec: %v", err)
			}

			topology.Status.ResourceStatus = v1alpha1.StatusFresh
			topology.Status.InitJobStatus = v1alpha1.StatusJobCompleted
			if err = r.Status().Update(ctx, topology); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating NodeCpuTopology status: %v", err)
			}

			logger.Info("NodeCpuTopology initialized successfully", "name", topology.Name)

			if err := r.deleteJob(ctx, topology.Status.InitJobName); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating NodeCpuTopology status: %v", err)
			} else {
				logger.Info("Job deleted", "jobName", topology.Status.InitJobName)
			}

			return ctrl.Result{}, nil

		default:
			return ctrl.Result{}, nil
		}
	default:
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeCpuTopologyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NodeCpuTopology{}, eventFilters).
		Complete(r)
}
