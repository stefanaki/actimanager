package nodecputopology

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getNodeByTopologyNodeName return the node with name specified in the NodeName field of the spec
func (r *NodeCpuTopologyReconciler) getNodeByTopologyNodeName(
	topology *cslabecentuagrv1alpha1.NodeCpuTopology,
	ctx context.Context) (*corev1.Node, error) {

	nodeName := topology.Spec.NodeName
	targetNode := &corev1.Node{}
	err := r.Get(ctx, client.ObjectKey{Name: nodeName}, targetNode)
	return targetNode, err
}
