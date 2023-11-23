package nodecputopology

import (
	"github.com/google/uuid"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func LscpuJobTemplate(node string) (string, *batchv1.Job) {
	jobName := "lscpu-job-" + node + strings.Split(uuid.New().String(), "-")[0]

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: "default",
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
							Args: []string{"-p=node,socket,core,cpu"},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	return jobName, job
}
