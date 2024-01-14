package nodecputopology

import (
	"context"
	"fmt"

	apiv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateInitialNodeCpuTopologies creates initial NodeCpuTopology resources for each cluster node.
// It lists the cluster nodes and existing topologies, and creates a new topology for each node that doesn't have one.
// The new topology is created with the name "<node-name>-cputopology".
func (r *NodeCpuTopologyReconciler) CreateInitialNodeCpuTopologies(ctx context.Context) error {
	nodes := &v1.NodeList{}
	topologies := &apiv1alpha1.NodeCpuTopologyList{}

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

		newTopology := &apiv1alpha1.NodeCpuTopology{
			ObjectMeta: metav1.ObjectMeta{Name: n.Name + "-cputopology"},
			Spec: apiv1alpha1.NodeCpuTopologySpec{
				NodeName: n.Name,
			},
		}

		if err := r.Create(ctx, newTopology); err != nil {
			return fmt.Errorf("error creating new topology: %v", err.Error())
		}
	}

	return nil
}
