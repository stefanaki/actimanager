package nodecputopology

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nodecputopologyv1alpha1 "cslab.ece.ntua.gr/actimanager/pkg/nodecputopology/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/pkg/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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

//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cslab.ece.ntua.gr,resources=nodecputopologies/finalizers,verbs=update

func (r *NodeCpuTopologyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("controller")

	topology := &cslabecentuagrv1alpha1.NodeCpuTopology{}
	if err := r.Get(ctx, req.NamespacedName, topology); err != nil {
		logger.V(5).Info("Error listing NodeCpuTopology resources:" + err.Error())
		return ctrl.Result{}, err
	}
	nodeName := topology.Spec.NodeName

	if topology.Status.InitJobStatus == "" {
		jobName, job := LscpuJobTemplate(nodeName)

		err := r.Client.Create(ctx, job)
		if err != nil {
			logger.V(5).Error(err, "Could not dispatch lscpu job")
		}

		topology.Status.InitJobStatus = "Pending"
		topology.Status.InitJobName = jobName
		err = r.Status().Update(ctx, topology)
		if err != nil {
			logger.V(5).Info("Error updating resource:" + err.Error())
		}

		return ctrl.Result{Requeue: true}, nil
	} else if topology.Status.InitJobStatus == "Pending" {
		job := &batchv1.Job{}

		err := r.Get(ctx, client.ObjectKey{Name: topology.Status.InitJobName, Namespace: "default"}, job)

		if err != nil {
			println("again...", err.Error())

		}

		if job.Status.Succeeded > 0 {
			topology.Status.InitJobStatus = "Completed"
			r.Status().Update(ctx, topology)

			podList := &corev1.PodList{}

			err := r.List(ctx, podList, client.MatchingLabels{"job-name": topology.Status.InitJobName})
			if err != nil {
				println("mala")
			}

			if len(podList.Items) > 0 {
				podLogs, _ := utils.GetPodLogs(podList.Items[0], context.TODO())

				cpuTopology, _ := nodecputopologyv1alpha1.NodeCpuTopologyV1Alpha1(podLogs)

				topology.Spec.Topology = cpuTopology
				r.Update(ctx, topology)
			}
		}

		return ctrl.Result{Requeue: true}, nil
	} else if topology.Status.InitJobStatus == "Completed" {

		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeCpuTopologyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cslabecentuagrv1alpha1.NodeCpuTopology{}).
		Complete(r)
}
