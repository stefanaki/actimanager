package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nodecputopologyv1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology/v1alpha1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *PodCpuBindingReconciler) validatePodCpuBinding(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding, pod *corev1.Pod) error {
	err := r.Get(ctx, client.ObjectKey{Name: cpuBinding.Spec.PodName, Namespace: cpuBinding.ObjectMeta.Namespace}, pod)

	if errors.IsNotFound(err) {
		cpuBinding.Status.Status = v1alpha1.StatusPodNotFound
		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return fmt.Errorf("error updating cpu binding status: %v", err.Error())
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting pod: %v", err.Error())
	}

	topologies := &v1alpha1.NodeCpuTopologyList{}
	err = r.List(ctx,
		topologies,
		client.MatchingFields{"spec.nodeName": pod.Spec.NodeName})

	if errors.IsNotFound(err) {
		cpuBinding.Status.Status = v1alpha1.StatusNodeTopologyNotFound
		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return fmt.Errorf("error updating cpu binding status: %v", err.Error())
		}
	}

	if err != nil {
		return fmt.Errorf("error listing CPU topologies: %v", err.Error())
	}

	targetTopology := topologies.Items[0]

	// Check if specified cpuset is available in the node topology
	if !nodecputopologyv1alpha1.IsCpuSetInTopology(&targetTopology.Spec.Topology, cpuBinding.Spec.CpuSet) {
		cpuBinding.Status.Status = v1alpha1.StatusInvalidCpuSet

		if err := r.Status().Update(ctx, cpuBinding); err != nil {
			return fmt.Errorf("error updating cpu binding status: %v", err.Error())
		}

		return nil
	}

	return nil
}
