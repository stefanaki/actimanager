package podcpubinding

import (
	"context"
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nctv1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology/v1alpha1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// validatePodName checks if the specified pod exists in the namespace
func (r *PodCpuBindingReconciler) validatePodName(ctx context.Context, cpuBinding *v1alpha1.PodCpuBinding, pod *corev1.Pod) (bool, string, error) {
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
	topology *v1alpha1.NodeCpuTopology, pod *corev1.Pod) (bool, string, error) {

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
	if !nctv1alpha1.IsCpuSetInTopology(&topology.Spec.Topology, cpuBinding.Spec.CpuSet) {
		return false, v1alpha1.StatusInvalidCpuSet, nil
	}

	return true, "", nil
}

// validateExclusivenessLevel checks if the specified CPU binding has
// an exclusive CPU set based on the specified exclusiveness level
func (r *PodCpuBindingReconciler) validateExclusivenessLevel(ctx context.Context,
	cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology,
	namespacedName types.NamespacedName, nodeName string) (bool, string, error) {

	// println("validating exclusiveness lvl of cpu binding", cpuBinding.Name)

	unfeasibleCpus := make(map[int]struct{})
	podCpuBindingList := &v1alpha1.PodCpuBindingList{}

	err := r.List(ctx, podCpuBindingList)

	if err != nil {
		return false, "", fmt.Errorf("failed to list PodCpuBindings: %v", err.Error())
	}

	for _, pcb := range podCpuBindingList.Items {
		// println("checking pcb " + pcb.Name)
		if pcb.Status.NodeName != nodeName ||
			(pcb.Namespace == namespacedName.Namespace && pcb.Name == namespacedName.Name) ||
			pcb.Status.ResourceStatus != v1alpha1.StatusApplied {
			// println("this pcb should not be checked", pcb.Name)
			continue
		}

		exclusiveCpus := getExclusiveCpusOfCpuBinding(&pcb, topology)
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

// getExclusiveCpusOfCpuBinding returns the exclusive CPUs
// for a given CPU binding based on its exclusiveness level
func getExclusiveCpusOfCpuBinding(cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology) map[int]struct{} {
	exclusiveCpus := make(map[int]struct{})

	// println("get exclusive cpus for pcb", cpuBinding.Name)

	for _, cpu := range cpuBinding.Spec.CpuSet {
		_, coreId, socketId, numaId := nctv1alpha1.GetCpuParentInfo(topology, cpu.CpuId)

		switch cpuBinding.Spec.ExclusivenessLevel {
		case "Cpu":
			exclusiveCpus[cpu.CpuId] = struct{}{}
		case "Core":
			for _, c := range nctv1alpha1.GetAllCpusInCore(topology, coreId) {
				exclusiveCpus[c] = struct{}{}
			}
		case "Socket":
			for _, c := range nctv1alpha1.GetAllCpusInSocket(topology, socketId) {
				exclusiveCpus[c] = struct{}{}
			}
		case "Numa":
			// println("pod", cpuBinding.Name, "is numa")
			for _, c := range nctv1alpha1.GetAllCpusInNuma(topology, numaId) {
				exclusiveCpus[c] = struct{}{}
			}
		default:

		}
	}

	// println("end exclusive cpus for pcb", cpuBinding.Name)
	//for c := range exclusiveCpus {
	//	log.Printf("%v", c)
	//}
	//log.Printf("\n")

	return exclusiveCpus
}

// hasUnfeasibleCpus checks if there are unfeasible CPUs for a given CPU binding
func hasUnfeasibleCpus(cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology, unfeasibleCpus map[int]struct{}) bool {
	exclusiveCpus := getExclusiveCpusOfCpuBinding(cpuBinding, topology)

	for cpu := range exclusiveCpus {
		if _, exists := unfeasibleCpus[cpu]; exists {
			return true
		}
	}
	return false
}
