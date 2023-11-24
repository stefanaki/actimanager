package nodecputopology

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nodecputopologyv1alpha1 "cslab.ece.ntua.gr/actimanager/pkg/nodecputopology/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func (r *NodeCpuTopologyReconciler) createInitJob(topology *cslabecentuagrv1alpha1.NodeCpuTopology, ctx context.Context, logger *logr.Logger) (string, error) {
	jobName, job := LscpuJobTemplate(topology.Spec.NodeName)

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

func LscpuJobTemplate(node string) (string, *batchv1.Job) {
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
