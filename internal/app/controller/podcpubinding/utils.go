package podcpubinding

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getNode retrieves the node with the given nodeName from the Kubernetes cluster.
func (r *PodCPUBindingReconciler) getNode(ctx context.Context, nodeName string) (*corev1.Node, error) {
	node := &corev1.Node{}

	err := r.Get(ctx, client.ObjectKey{
		Name: nodeName,
	}, node)

	return node, err
}

// getPod retrieves the pod with the given podNamespacedName from the Kubernetes cluster.
func (r *PodCPUBindingReconciler) getPod(ctx context.Context, podNamespacedName types.NamespacedName) (*corev1.Pod, error) {
	pod := &corev1.Pod{}

	err := r.Get(ctx, client.ObjectKey{
		Name:      podNamespacedName.Name,
		Namespace: podNamespacedName.Namespace,
	}, pod)

	return pod, err
}

// getNodeAddress retrieves the IP address of the node with the given nodeName from the Kubernetes cluster.
func (r *PodCPUBindingReconciler) getNodeAddress(ctx context.Context, nodeName string) (string, error) {
	node, err := r.getNode(ctx, nodeName)
	if err != nil {
		return "", fmt.Errorf("failed to get node by name: %v", err.Error())
	}

	nodeAddress := ""
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			nodeAddress = address.Address
			break
		}
	}

	if nodeAddress == "" {
		return "", fmt.Errorf("failed to get IP address of node " + nodeName)
	}

	return nodeAddress, nil
}
