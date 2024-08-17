package nodecputopology

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const NodeRoleControlPlane = "node-role.kubernetes.io/control-plane"

// CreateInitialNodeCPUTopologies creates initial NodeCPUTopology resources for each cluster node.
// It lists the cluster nodes and existing topologies, and creates a new topology for each node that doesn't have one.
// The new topology is created with the name "<node-name>-cputopology".
func (r *NodeCPUTopologyReconciler) CreateInitialNodeCPUTopologies(ctx context.Context) error {
	nodes := &v1.NodeList{}
	topologies := &v1alpha1.NodeCPUTopologyList{}
	if err := r.List(ctx, nodes); err != nil {
		return fmt.Errorf("error listing cluster nodes: %v", err.Error())
	}
	if err := r.List(ctx, topologies); err != nil {
		return fmt.Errorf("error listing topologies: %v", err.Error())
	}
	for _, n := range nodes.Items {
		skip := false
		for _, t := range topologies.Items {
			if t.Spec.NodeName == n.Name {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		newTopology := &v1alpha1.NodeCPUTopology{
			ObjectMeta: metav1.ObjectMeta{Name: n.Name + "-cputopology"},
			Spec: v1alpha1.NodeCPUTopologySpec{
				NodeName: n.Name,
			},
		}
		if _, controlPlane := n.Labels[NodeRoleControlPlane]; controlPlane {
			newTopology.Labels = map[string]string{NodeRoleControlPlane: ""}
		}
		if err := r.Create(ctx, newTopology); err != nil {
			return fmt.Errorf("error creating new topology: %v", err.Error())
		}
	}
	return nil
}
