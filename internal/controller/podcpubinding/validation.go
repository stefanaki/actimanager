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

func (r *PodCPUBindingReconciler) validateResource(ctx context.Context, cpuBinding *v1alpha1.PodCPUBinding, topology *v1alpha1.NodeCPUTopology, pod *corev1.Pod) (bool, v1alpha1.PodCPUBindingResourceStatus, string, error) {
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
func (r *PodCPUBindingReconciler) validatePodName(ctx context.Context, cpuBinding *v1alpha1.PodCPUBinding, pod *corev1.Pod) (bool, v1alpha1.PodCPUBindingResourceStatus, string, error) {
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
func (r *PodCPUBindingReconciler) validateTopology(ctx context.Context, cpuBinding *v1alpha1.PodCPUBinding,
	topology *v1alpha1.NodeCPUTopology, pod *corev1.Pod) (bool, v1alpha1.PodCPUBindingResourceStatus, string, error) {
	// Get NodeCPUTopology of node
	topologies := &v1alpha1.NodeCPUTopologyList{}
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
	if !nct.IsCPUSetInTopology(&topology.Spec.Topology, cpuBinding.Spec.CPUSet) {
		return false,
			v1alpha1.StatusInvalidCPUSet,
			fmt.Sprintf("CPUs %v do not exist in node %v", pcbutils.ConvertCPUSliceToIntSlice(cpuBinding.Spec.CPUSet),
				topology.Name),
			nil
	}
	return true, "", "", nil
}

// validateExclusivenessLevel checks if the specified CPU binding has
// an exclusive CPU set based on the specified exclusiveness level
func (r *PodCPUBindingReconciler) validateExclusivenessLevel(ctx context.Context,
	cpuBinding *v1alpha1.PodCPUBinding, topology *v1alpha1.NodeCPUTopology,
	namespacedName types.NamespacedName, nodeName string) (bool, v1alpha1.PodCPUBindingResourceStatus, string, error) {
	unfeasibleCPUs := make(map[int]struct{})
	podCPUBindingList := &v1alpha1.PodCPUBindingList{}
	err := r.List(ctx, podCPUBindingList,
		client.MatchingFields{"status.nodeName": nodeName},
		client.MatchingFields{"status.resourceStatus": string(v1alpha1.StatusApplied)})

	if err != nil {
		return false, "", "", fmt.Errorf("failed to list PodCPUBindings: %v", err.Error())
	}
	for _, pcb := range podCPUBindingList.Items {
		if (pcb.Namespace == namespacedName.Namespace && pcb.Name == namespacedName.Name) ||
			pcb.Status.ResourceStatus != v1alpha1.StatusApplied || pcb.Status.NodeName != nodeName {
			continue
		}
		exclusiveCPUs := pcbutils.GetExclusiveCPUsOfCPUBinding(&pcb, &topology.Spec.Topology)
		for c := range exclusiveCPUs {
			unfeasibleCPUs[c] = struct{}{}
		}
	}
	isUnfeasible, cpus := hasUnfeasibleCPUs(cpuBinding, topology, unfeasibleCPUs)
	if isUnfeasible {
		return false,
			v1alpha1.StatusCPUSetAllocationFailed,
			fmt.Sprintf("CPUs %v are already allocated to another pod", cpus),
			nil
	}
	return true, "", "", nil
}

// hasUnfeasibleCPUs checks if there are unfeasible CPUs for a given CPU binding
func hasUnfeasibleCPUs(cpuBinding *v1alpha1.PodCPUBinding, topology *v1alpha1.NodeCPUTopology, unfeasibleCPUs map[int]struct{}) (bool, []int) {
	var requestedCPUs map[int]struct{}
	if cpuBinding.Spec.ExclusivenessLevel == "None" {
		requestedCPUs = pcbutils.GetCPUsOfCPUBinding(cpuBinding)
	} else {
		requestedCPUs = pcbutils.GetExclusiveCPUsOfCPUBinding(cpuBinding, &topology.Spec.Topology)
	}
	cpus := make(map[int]struct{})
	for cpu := range requestedCPUs {
		if _, exists := unfeasibleCPUs[cpu]; exists {
			cpus[cpu] = struct{}{}
		}
	}
	if len(cpus) > 0 {
		return true, maps.Keys(cpus)
	}
	return false, nil
}
