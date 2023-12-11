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

func (r *PodCpuBindingReconciler) validateExclusivenessLevel(ctx context.Context,
	cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology,
	namespacedName types.NamespacedName) (bool, string, error) {

	unfeasibleCpus := make(map[int]struct{})
	podCpuBindingList := &v1alpha1.PodCpuBindingList{}

	err := r.List(ctx, podCpuBindingList)

	if err != nil {
		return false, "", fmt.Errorf("failed to list PodCpuBindings: %v", err.Error())
	}

	for _, pcb := range podCpuBindingList.Items {
		if pcb.Status.NodeName != cpuBinding.Status.NodeName ||
			(pcb.Namespace == namespacedName.Namespace && pcb.Name == namespacedName.Name) ||
			pcb.Status.Status != v1alpha1.StatusApplied {
			continue
		}

		switch pcb.Spec.ExclusivenessLevel {
		case "Cpu":
			for _, cpu := range pcb.Spec.CpuSet {
				unfeasibleCpus[cpu.CpuId] = struct{}{}
			}
		case "Core":
			for _, cpu := range pcb.Spec.CpuSet {
				_, coreId, _, _ := nctv1alpha1.GetCpuParents(topology, cpu.CpuId)
				for _, c := range nctv1alpha1.GetAllCpusInCore(topology, coreId) {
					unfeasibleCpus[c] = struct{}{}
				}
			}
		case "Socket":
			for _, cpu := range pcb.Spec.CpuSet {
				_, _, socketId, _ := nctv1alpha1.GetCpuParents(topology, cpu.CpuId)

				for _, c := range nctv1alpha1.GetAllCpusInSocket(topology, socketId) {
					unfeasibleCpus[c] = struct{}{}
				}
			}
		case "Numa":
			for _, cpu := range pcb.Spec.CpuSet {
				_, _, _, numaId := nctv1alpha1.GetCpuParents(topology, cpu.CpuId)

				for _, c := range nctv1alpha1.GetAllCpusInNuma(topology, numaId) {
					unfeasibleCpus[c] = struct{}{}
				}
			}
		default:
			// Exclusiveness Level: None
		}
	}

	println("UNFEASIBLE CPUS: ", namespacedName.Name, namespacedName.Namespace)
	for c := range unfeasibleCpus {
		fmt.Printf("%v,", c)
	}

	if hasUnfeasibleCpus(cpuBinding.Spec.CpuSet, unfeasibleCpus) {
		return false, v1alpha1.StatusCpuSetAllocationFailed, nil
	}

	return true, "", nil
}

func hasUnfeasibleCpus(cpuSet []v1alpha1.Cpu, unfeasibleCpus map[int]struct{}) bool {
	for _, cpu := range cpuSet {
		if _, exists := unfeasibleCpus[cpu.CpuId]; exists {
			return true
		}
	}
	return false
}
