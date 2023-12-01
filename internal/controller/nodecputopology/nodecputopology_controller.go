package nodecputopology

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NodeCpuTopologyReconciler reconciles a NodeCpuTopology object
type NodeCpuTopologyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

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
	topology := &cslabecentuagrv1alpha1.NodeCpuTopology{}

	// Handle delete
	err := r.Get(ctx, req.NamespacedName, topology)
	if errors.IsNotFound(err) {
		logger.Info("Deleted NodeCpuTopology")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	// Initialize CR
	if topology.Spec.NodeName != topology.Status.LastNodeName {
		topology.Status.Status = "NeedsSync"
		topology.Status.LastNodeName = topology.Spec.NodeName
		topology.Status.InitJobStatus = "None"
		if err := r.Status().Update(ctx, topology); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating resource: %v", err)
		}
		return ctrl.Result{}, nil
	}

	// Validate CR
	if topology.Status.Status == "NodeNotFound" ||
		topology.Status.Status == "Fresh" {
		return ctrl.Result{}, nil
	}

	// Check if specified NodeName is a valid name of a node
	if err := r.Get(ctx, client.ObjectKey{Name: topology.Spec.NodeName}, &corev1.Node{}); err != nil {
		logger.Info("Node with specified name not found: " + topology.Spec.NodeName)
		topology.Status.Status = "NodeNotFound"
		topology.Status.InitJobStatus = "None"

		if err := r.Status().Update(ctx, topology); err != nil {
			return ctrl.Result{}, fmt.Errorf("could not update status: %v", err)
		}

		return ctrl.Result{}, nil
	}

	// Handle reconcilation
	switch topology.Status.Status {
	case "NeedsSync":
		// If Status is empty or NeedsSync, initiate job
		switch topology.Status.InitJobStatus {
		case "None":
			logger.Info("Dispatch init job for NodeCpuBinding")

			jobName, err := r.createInitJob(topology, ctx, &logger)

			topology.Status.InitJobStatus = "Pending"
			topology.Status.InitJobName = jobName

			if err = r.Status().Update(ctx, topology); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating resource: %v", err)
			}

			return ctrl.Result{}, nil
		case "Pending":
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

			topology.Status.Status = "Fresh"
			topology.Status.InitJobStatus = "Completed"
			if err = r.Status().Update(ctx, topology); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating NodeCpuTopology status: %v", err)
			}

			logger.Info("NodeCpuTopology for node " + topology.Spec.NodeName + " initialized successfully")

			if err := r.deleteJob(ctx, topology.Status.InitJobName); err != nil {
				return ctrl.Result{}, fmt.Errorf("error updating NodeCpuTopology status: %v", err)
			} else {
				logger.Info("Job " + topology.Status.InitJobName + " deleted successfully")
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
		For(&cslabecentuagrv1alpha1.NodeCpuTopology{}).
		Complete(r)
}
