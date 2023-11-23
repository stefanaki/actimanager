package nodecputopology

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nodecputopologyv1alpha1 "cslab.ece.ntua.gr/actimanager/pkg/nodecputopology/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/pkg/utils"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var eventFilters = builder.WithPredicates(predicate.Funcs{GenericFunc: func(e event.GenericEvent) bool {
	switch object := e.Object.(type) {
	case *cslabecentuagrv1alpha1.NodeCpuTopology:
		return object.Status.InitJobStatus != "Failed"
	default:
		return false
	}
}})

// NodeCpuTopologyReconciler reconciles a NodeCpuTopology object
type NodeCpuTopologyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/finalizers,verbs=update

func (r *NodeCpuTopologyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("nct-controller")

	topology := &cslabecentuagrv1alpha1.NodeCpuTopology{}

	err := r.Get(ctx, req.NamespacedName, topology)
	if errors.IsNotFound(err) {
		logger.Info("Could not find NodeCpuTopology")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err != nil {
		return reconcile.Result{}, err
	}

	if topology.Status.InitJobStatus == "Completed" {
		return ctrl.Result{}, nil
	}

	_, err = r.getNodeByTopologyNodeName(topology, &ctx)

	if err != nil {
		logger.Info("Node with specified nodeName not found: " + topology.Spec.NodeName)
		topology.Status.InitJobStatus = "Failed"

		if err := r.Status().Update(ctx, topology); err != nil {
			logger.Error(err, "could not update status")
			return ctrl.Result{Requeue: true}, err
		}

		return ctrl.Result{}, nil
	}

	if topology.Status.InitJobStatus == "" {
		logger.Info("Dispatch init job for new NodeCpuBinding")

		err := r.createInitJob(topology, &ctx, &logger)

		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}

		return ctrl.Result{Requeue: true}, nil
	} else if topology.Status.InitJobStatus == "Pending" {
		job := &batchv1.Job{}

		err := r.Get(ctx, client.ObjectKey{Name: topology.Status.InitJobName, Namespace: "actimanager-system"}, job)

		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}

		if job.Status.Succeeded > 0 {
			err := r.parseCompletedPod(topology, &ctx, &logger)

			if err != nil {
				logger.Error(err, "Get cpu topology")
				return ctrl.Result{Requeue: true}, nil
			}
		}

		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeCpuTopologyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cslabecentuagrv1alpha1.NodeCpuTopology{}, eventFilters).
		Complete(r)
}

func (r *NodeCpuTopologyReconciler) getNodeByTopologyNodeName(
	topology *cslabecentuagrv1alpha1.NodeCpuTopology,
	ctx *context.Context) (*corev1.Node, error) {
	nodeName := topology.Spec.NodeName
	targetNode := &corev1.Node{}

	err := r.Get(*ctx, client.ObjectKey{Name: nodeName}, targetNode)

	return targetNode, err
}

func (r *NodeCpuTopologyReconciler) createInitJob(topology *cslabecentuagrv1alpha1.NodeCpuTopology, ctx *context.Context, logger *logr.Logger) error {
	jobName, job := LscpuJobTemplate(topology.Spec.NodeName)

	err := r.Client.Create(*ctx, job)
	if err != nil {
		logger.Error(err, "Could not dispatch lscpu job")
		return err
	}

	topology.Status.InitJobStatus = "Pending"
	topology.Status.InitJobName = jobName
	err = r.Status().Update(*ctx, topology)
	if err != nil {
		logger.Info("Error updating resource:" + err.Error())
		return err
	}

	return nil
}

func (r *NodeCpuTopologyReconciler) parseCompletedPod(topology *cslabecentuagrv1alpha1.NodeCpuTopology, ctx *context.Context, logger *logr.Logger) error {
	topology.Status.InitJobStatus = "Completed"

	if err := r.Status().Update(*ctx, topology); err != nil {
		logger.Error(err, "Error updating topology status to Completed", err.Error())
		return err
	}

	podList := &corev1.PodList{}

	if err := r.List(*ctx, podList, client.MatchingLabels{"job-name": topology.Status.InitJobName}); err != nil {
		logger.Error(err, "Error getting retrieving job", err.Error())
		return err
	}

	if len(podList.Items) > 0 {
		podLogs, _ := utils.GetPodLogs(podList.Items[0], ctx)

		cpuTopology, err := nodecputopologyv1alpha1.NodeCpuTopologyV1Alpha1(podLogs)

		if err != nil {
			logger.Error(err, "Error parsing cpu topology")
			return err
		}

		topology.Spec.Topology = cpuTopology
		if err = r.Update(*ctx, topology); err != nil {
			logger.Error(err, "Error updating NodeCpuTopology")
			return err
		}
	}

	return nil
}
