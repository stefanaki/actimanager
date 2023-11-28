package nodecputopology

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nodecputopologyv1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/utils"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func (r *NodeCpuTopologyReconciler) createInitJob(topology *cslabecentuagrv1alpha1.NodeCpuTopology, ctx context.Context, logger *logr.Logger) (string, error) {
	jobName, job := lscpuJob(topology.Spec.NodeName)

	err := r.Client.Create(ctx, job)
	if err != nil {
		logger.Error(err, "Could not dispatch lscpu job")
		return "", err
	}

	return jobName, nil
}

func (r *NodeCpuTopologyReconciler) parseCompletedPod(topology *cslabecentuagrv1alpha1.NodeCpuTopology, ctx context.Context, logger *logr.Logger) (cslabecentuagrv1alpha1.CpuTopology, error) {
	podList := &corev1.PodList{}
	cpuTopology := cslabecentuagrv1alpha1.CpuTopology{}
	var err error

	if err := r.List(ctx, podList, client.MatchingLabels{"job-name": topology.Status.InitJobName}); err != nil {
		logger.Error(err, "Error getting retrieving job", err.Error())
		return cpuTopology, err
	}

	if len(podList.Items) > 0 {
		podLogs, _ := utils.GetPodLogs(podList.Items[0], ctx)

		cpuTopology, err = nodecputopologyv1alpha1.NodeCpuTopologyV1Alpha1(podLogs)

		if err != nil {
			logger.Error(err, "Error parsing cpu topology")
			return cpuTopology, err
		}
	}

	return cpuTopology, nil
}

// isJobCompleted checks if job with name InitJobName has been completed
func (r *NodeCpuTopologyReconciler) isJobCompleted(topology *cslabecentuagrv1alpha1.NodeCpuTopology, ctx context.Context) (bool, error) {
	job := &batchv1.Job{}
	err := r.Get(ctx, client.ObjectKey{Name: topology.Status.InitJobName, Namespace: "actimanager-system"}, job)
	return job.Status.Succeeded > 0, err
}

// deleteJob deletes the job that was created by the reconciler
func (r *NodeCpuTopologyReconciler) deleteJob(ctx context.Context, jobName string) error {
	job := &batchv1.Job{}
	deletePropagationBackground := metav1.DeletePropagationBackground

	if err := r.Get(ctx, client.ObjectKey{Name: jobName, Namespace: "actimanager-system"}, job); err != nil {
		return fmt.Errorf("could not retrieve job: %v", err.Error())
	}

	if err := r.Client.Delete(ctx, job, &client.DeleteOptions{PropagationPolicy: &deletePropagationBackground}); err != nil {
		return fmt.Errorf("could not delete job: %v", err.Error())
	}

	return nil
}

// lscpuJob generates a Kubernetes Job for running the 'lscpu' command on a specific node
func lscpuJob(node string) (string, *batchv1.Job) {
	jobName := "lscpu-job-" + node + strings.Split(uuid.New().String(), "-")[0]

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: "actimanager-system",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeName: node,
					Containers: []corev1.Container{
						{
							Name:  "lscpu-container",
							Image: "actions/lscpu",
							Command: []string{
								"lscpu",
							},
							Args: []string{"-p=socket,node,core,cpu"},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	return jobName, job
}
