package nodecputopology

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	nctutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// NodeCPUTopologyReconciler reconciles a NodeCPUTopology object
type NodeCPUTopologyReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var eventFilters = builder.WithPredicates(predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		oldObj := e.ObjectOld.(*v1alpha1.NodeCPUTopology)
		newObj := e.ObjectNew.(*v1alpha1.NodeCPUTopology)

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
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *NodeCPUTopologyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// logger := log.FromContext(ctx).WithName("nct-controller")

	// Get NodeCPUTopology CR
	topology := &v1alpha1.NodeCPUTopology{}

	// Handle delete
	err := r.Get(ctx, req.NamespacedName, topology)
	if errors.IsNotFound(err) {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Initialize CR
	if topology.Status.ResourceStatus == "" ||
		topology.Status.ResourceStatus == v1alpha1.StatusNodeNotFound ||
		topology.Status.ResourceStatus == v1alpha1.StatusTopologyFailed {
		topology.Status.ResourceStatus = v1alpha1.StatusNeedsSync
		if err := r.Status().Update(ctx, topology); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating resource: %v", err)
		}
		return ctrl.Result{}, nil
	}

	// Check if specified NodeName is a valid name of a node
	node := &corev1.Node{}
	if err := r.Get(ctx, client.ObjectKey{Name: topology.Spec.NodeName}, node); err != nil {
		topology.Status.ResourceStatus = v1alpha1.StatusNodeNotFound
		if err := r.Status().Update(ctx, topology); err != nil {
			return ctrl.Result{}, fmt.Errorf("could not update status: %v", err)
		}
		r.Recorder.Eventf(topology, corev1.EventTypeWarning, string(v1alpha1.StatusNodeNotFound), "Node %s not found", topology.Spec.NodeName)
		return ctrl.Result{}, nil
	}

	// Handle reconciliation
	topologyResponse, err := r.getTopology(ctx, node)
	if err != nil {
		r.Recorder.Eventf(topology, corev1.EventTypeWarning, string(v1alpha1.StatusTopologyFailed), "Failed to get topology: %v", err)
		return ctrl.Result{}, err
	}
	cpuTopology := nctutils.TopologyToV1Alpha1(topologyResponse)
	topology.Spec.Topology = *cpuTopology
	topology.Spec.NodeName = node.Name
	if err := r.Update(ctx, topology); err != nil {
		r.Recorder.Eventf(topology, corev1.EventTypeWarning, string(v1alpha1.StatusTopologyFailed), "Failed to update topology: %v", err)
		return ctrl.Result{}, err
	}
	topology.Status.ResourceStatus = v1alpha1.StatusFresh
	if err := r.Status().Update(ctx, topology); err != nil {
		r.Recorder.Eventf(topology, corev1.EventTypeWarning, string(v1alpha1.StatusTopologyFailed), "Failed to update status: %v", err)
		return ctrl.Result{}, err
	}

	r.Recorder.Eventf(topology, corev1.EventTypeNormal, string(v1alpha1.StatusFresh), "Topology is up to date, CPUs: %v", topology.Spec.Topology.CPUs)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeCPUTopologyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NodeCPUTopology{}, eventFilters).
		Complete(r)
}
