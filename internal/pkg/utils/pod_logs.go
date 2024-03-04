package utils

import (
	"bytes"
	"context"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/client"
	"io"
	corev1 "k8s.io/api/core/v1"
	"log"
)

func GetPodLogs(pod corev1.Pod, ctx context.Context) (string, error) {
	clientset, err := client.NewClient()

	if err != nil {
		log.Printf("error getting clientset: %v", err.Error())
		return "", err
	}

	podLogOpts := corev1.PodLogOptions{}

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)

	if err != nil {
		log.Printf("error getting log stream: %v", err.Error())
		return "", err
	}

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)

	if err != nil {
		log.Printf("error copying logs to buffer: %v", err)
		return "", err
	}

	str := buf.String()

	return str, nil
}
