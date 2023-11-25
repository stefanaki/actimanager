package nodecputopology

import (
	"context"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *NodeCpuTopologyReconciler) CreateInitialNodeCpuTopologies(ctx context.Context) error {
	nodes := &v1.NodeList{}
	topologies := &cslabecentuagrv1alpha1.NodeCpuTopologyList{}

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

		newTopology := &cslabecentuagrv1alpha1.NodeCpuTopology{
			ObjectMeta: metav1.ObjectMeta{Name: n.Name + "-cputopology"},
			Spec: cslabecentuagrv1alpha1.NodeCpuTopologySpec{
				NodeName: n.Name,
			},
		}

		if err := r.Create(ctx, newTopology); err != nil {
			return fmt.Errorf("error creating new topology: %v", err.Error())
		}
	}

	return nil
}
