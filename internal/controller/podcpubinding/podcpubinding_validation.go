package podcpubinding

import (
	"context"
	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/podcpubinding"
	"fmt"

	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nct "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// validatePodName checks if the specified pod exists in the namespace
func (r *PodCpuBindingReconciler) validatePodName(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding, pod *corev1.Pod) (bool, v1alpha1.PodCpuBindingResourceStatus, error) {
	err := r.Get(ctx, client.ObjectKey{Name: cpuBinding.Spec.PodName, Namespace: cpuBinding.ObjectMeta.Namespace}, pod)

	if errors.IsNotFound(err) {
		return false, v1alpha1.StatusPodNotFound, nil
	}

	if err != nil {
		return false, "", fmt.Errorf("error getting pod: %v", err.Error())
	}

	return true, "", nil
}

// validateTopology checks if the node topology for the specified pod's
// node is available and if the specified CPU set is valid
func (r *PodCpuBindingReconciler) validateTopology(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding,
	topology *v1alpha1.NodeCpuTopology, pod *corev1.Pod) (bool, v1alpha1.PodCpuBindingResourceStatus, error) {

	// Get NodeCpuTopology of node
	topologies := &v1alpha1.NodeCpuTopologyList{}
	err := r.List(ctx,
		topologies,
		client.MatchingFields{"spec.nodeName": pod.Spec.NodeName})

	if errors.IsNotFound(err) {
		return false, v1alpha1.StatusNodeTopologyNotFound, nil
	}

	if err != nil {
		return false, "", fmt.Errorf("error listing CPU topologies: %v", err.Error())
	}

	*topology = topologies.Items[0]

	// Check if specified cpuset is available in the node topology
	if !nct.IsCpuSetInTopology(&topology.Spec.Topology, cpuBinding.Spec.CpuSet) {
		return false, v1alpha1.StatusInvalidCpuSet, nil
	}

	return true, "", nil
}

// validateExclusivenessLevel checks if the specified CPU binding has
// an exclusive CPU set based on the specified exclusiveness level
func (r *PodCpuBindingReconciler) validateExclusivenessLevel(ctx context.Context,
	cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology,
	namespacedName types.NamespacedName, nodeName string) (bool, v1alpha1.PodCpuBindingResourceStatus, error) {

	// println("validating exclusiveness lvl of cpu binding", cpuBinding.Name)

	unfeasibleCpus := make(map[int]struct{})
	podCpuBindingList := &v1alpha1.PodCpuBindingList{}

	err := r.List(ctx, podCpuBindingList,
		client.MatchingFields{"status.nodeName": nodeName},
		client.MatchingFields{"status.resourceStatus": string(v1alpha1.StatusApplied)})

	if err != nil {
		return false, "", fmt.Errorf("failed to list PodCpuBindings: %v", err.Error())
	}

	for _, pcb := range podCpuBindingList.Items {
		// println("checking pcb " + pcb.Name)
		if (pcb.Namespace == namespacedName.Namespace && pcb.Name == namespacedName.Name) ||
			pcb.Status.ResourceStatus != v1alpha1.StatusApplied {
			// println("this pcb should not be checked", pcb.Name)
			continue
		}

		exclusiveCpus := pcbutils.GetExclusiveCpusOfCpuBinding(&pcb, topology)
		for c := range exclusiveCpus {
			// println("found exclusive cpu for", pcb.Name, c)
			unfeasibleCpus[c] = struct{}{}
		}
	}

	if hasUnfeasibleCpus(cpuBinding, topology, unfeasibleCpus) {
		// println("cant pin")
		return false, v1alpha1.StatusCpuSetAllocationFailed, nil
	} else {
		// println("can pin")
	}

	return true, "", nil
}

// hasUnfeasibleCpus checks if there are unfeasible CPUs for a given CPU binding
func hasUnfeasibleCpus(cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology, unfeasibleCpus map[int]struct{}) bool {
	exclusiveCpus := pcbutils.GetExclusiveCpusOfCpuBinding(cpuBinding, topology)

	for cpu := range exclusiveCpus {
		if _, exists := unfeasibleCpus[cpu]; exists {
			return true
		}
	}
	return false
}
