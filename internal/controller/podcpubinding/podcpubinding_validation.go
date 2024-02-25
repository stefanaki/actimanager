package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	nct "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/podcpubinding"
	"fmt"
	"golang.org/x/exp/maps"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *PodCpuBindingReconciler) validateResource(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology, pod *corev1.Pod) (bool, v1alpha1.PodCpuBindingResourceStatus, string, error) {
	// Validate pod name
	ok, status, message, err := r.validatePodName(ctx, cpuBinding, pod)
	if !ok {
		return ok, status, message, err
	}
	// Validate topology
	ok, status, message, err = r.validateTopology(ctx, cpuBinding, topology, pod)
	if !ok {
		return ok, status, message, err
	}
	// Validate exclusiveness level
	ok, status, message, err = r.validateExclusivenessLevel(ctx, cpuBinding, topology, types.NamespacedName{Name: cpuBinding.Name, Namespace: cpuBinding.Namespace}, pod.Spec.NodeName)
	if !ok {
		return ok, status, message, err
	}
	return true, "", "", nil
}

// validatePodName checks if the specified pod exists in the namespace
func (r *PodCpuBindingReconciler) validatePodName(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding, pod *corev1.Pod) (bool, v1alpha1.PodCpuBindingResourceStatus, string, error) {
	err := r.Get(ctx, client.ObjectKey{Name: cpuBinding.Spec.PodName, Namespace: cpuBinding.ObjectMeta.Namespace}, pod)
	if errors.IsNotFound(err) {
		return false, v1alpha1.StatusPodNotFound, fmt.Sprintf("pod %s/%s not found", cpuBinding.Namespace, cpuBinding.Spec.PodName), nil
	}
	if err != nil {
		return false, "", "", fmt.Errorf("error getting pod: %v", err)
	}
	return true, "", "", nil
}

// validateTopology checks if the node topology for the specified pod's
// node is available and if the specified CPU set is valid
func (r *PodCpuBindingReconciler) validateTopology(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding,
	topology *v1alpha1.NodeCpuTopology, pod *corev1.Pod) (bool, v1alpha1.PodCpuBindingResourceStatus, string, error) {
	// Get NodeCpuTopology of node
	topologies := &v1alpha1.NodeCpuTopologyList{}
	err := r.List(ctx,
		topologies,
		client.MatchingFields{"spec.nodeName": pod.Spec.NodeName})

	if errors.IsNotFound(err) {
		return false, v1alpha1.StatusNodeTopologyNotFound, fmt.Sprintf("topology for node %v not found", pod.Spec.NodeName), nil
	}
	if err != nil {
		return false, "", "", fmt.Errorf("error listing CPU topologies: %v", err.Error())
	}
	*topology = topologies.Items[0]
	// Check if specified cpuset is available in the node topology
	if !nct.IsCpuSetInTopology(&topology.Spec.Topology, cpuBinding.Spec.CpuSet) {
		return false,
			v1alpha1.StatusInvalidCpuSet,
			fmt.Sprintf("CPUs %v do not exist in node %v", pcbutils.ConvertCpuSliceToIntSlice(cpuBinding.Spec.CpuSet),
				topology.Name),
			nil
	}
	return true, "", "", nil
}

// validateExclusivenessLevel checks if the specified CPU binding has
// an exclusive CPU set based on the specified exclusiveness level
func (r *PodCpuBindingReconciler) validateExclusivenessLevel(ctx context.Context,
	cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology,
	namespacedName types.NamespacedName, nodeName string) (bool, v1alpha1.PodCpuBindingResourceStatus, string, error) {
	unfeasibleCpus := make(map[int]struct{})
	podCpuBindingList := &v1alpha1.PodCpuBindingList{}
	err := r.List(ctx, podCpuBindingList,
		client.MatchingFields{"status.nodeName": nodeName},
		client.MatchingFields{"status.resourceStatus": string(v1alpha1.StatusApplied)})

	if err != nil {
		return false, "", "", fmt.Errorf("failed to list PodCpuBindings: %v", err.Error())
	}
	for _, pcb := range podCpuBindingList.Items {
		if (pcb.Namespace == namespacedName.Namespace && pcb.Name == namespacedName.Name) ||
			pcb.Status.ResourceStatus != v1alpha1.StatusApplied || pcb.Status.NodeName != nodeName {
			continue
		}
		exclusiveCpus := pcbutils.GetExclusiveCpusOfCpuBinding(&pcb, topology)
		for c := range exclusiveCpus {
			unfeasibleCpus[c] = struct{}{}
		}
	}
	isUnfeasible, cpus := hasUnfeasibleCpus(cpuBinding, topology, unfeasibleCpus)
	if isUnfeasible {
		return false,
			v1alpha1.StatusCpuSetAllocationFailed,
			fmt.Sprintf("CPUs %v are already allocated to another pod", cpus),
			nil
	}
	return true, "", "", nil
}

// hasUnfeasibleCpus checks if there are unfeasible CPUs for a given CPU binding
func hasUnfeasibleCpus(cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology, unfeasibleCpus map[int]struct{}) (bool, []int) {
	exclusiveCpus := pcbutils.GetExclusiveCpusOfCpuBinding(cpuBinding, topology)
	cpus := make(map[int]struct{})
	for cpu := range exclusiveCpus {
		if _, exists := unfeasibleCpus[cpu]; exists {
			cpus[cpu] = struct{}{}
		}
	}
	if len(cpus) > 0 {
		return true, maps.Keys(cpus)
	}
	return false, nil
}
