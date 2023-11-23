package utils

import (
	"bytes"
	"context"
	"cslab.ece.ntua.gr/actimanager/pkg/client"
	"github.com/pkg/errors"
	"io"
	corev1 "k8s.io/api/core/v1"
)

func GetPodLogs(pod corev1.Pod, ctx context.Context) (string, error) {
	clientset, err := client.GetClientSet()

	if err != nil {
		errors.Errorf("error getting clientset: %v", err)
		return "", err
	}

	podLogOpts := corev1.PodLogOptions{}

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)

	if err != nil {
		errors.Errorf("error getting log stream: %v", err)
		return "", err
	}

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)

	if err != nil {
		errors.Errorf("error copying logs to buffer: %v", err)
		return "", err
	}

	str := buf.String()

	return str, nil
}
