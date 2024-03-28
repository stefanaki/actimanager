package podcpubinding

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	nct "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
)

// ExclusiveCPUsOfCPUBinding returns the exclusive CPUs
// for a given CPU binding based on its exclusiveness level
func ExclusiveCPUsOfCPUBinding(cpuBinding *v1alpha1.PodCPUBinding, topology *v1alpha1.CPUTopology) map[int]struct{} {
	exclusiveCPUs := make(map[int]struct{})
	for _, cpu := range cpuBinding.Spec.CPUSet {
		_, coreID, socketID, numaID := nct.CPUParentInfo(topology, cpu.CPUID)
		switch cpuBinding.Spec.ExclusivenessLevel {
		case v1alpha1.ResourceLevelCPU:
			exclusiveCPUs[cpu.CPUID] = struct{}{}
		case v1alpha1.ResourceLevelCore:
			for _, c := range nct.CPUsInCore(topology, coreID) {
				exclusiveCPUs[c] = struct{}{}
			}
		case v1alpha1.ResourceLevelSocket:
			for _, c := range nct.CPUsInSocket(topology, socketID) {
				exclusiveCPUs[c] = struct{}{}
			}
		case v1alpha1.ResourceLevelNUMA:
			for _, c := range nct.CPUsInNUMA(topology, numaID) {
				exclusiveCPUs[c] = struct{}{}
			}
		default:
		}
	}
	return exclusiveCPUs
}

func CPUsOfCPUBinding(cpuBinding *v1alpha1.PodCPUBinding) map[int]struct{} {
	cpus := make(map[int]struct{})
	for _, cpu := range cpuBinding.Spec.CPUSet {
		cpus[cpu.CPUID] = struct{}{}
	}
	return cpus
}

func CPUSliceToIntSlice(cpuSlice []v1alpha1.CPU) []int {
	intSlice := make([]int, len(cpuSlice))
	for i, cpu := range cpuSlice {
		intSlice[i] = cpu.CPUID
	}
	return intSlice
}

func IntSliceToCPUSlice(intSlice []int) []v1alpha1.CPU {
	cpuSlice := make([]v1alpha1.CPU, len(intSlice))
	for i, cpuID := range intSlice {
		cpuSlice[i] = v1alpha1.CPU{CPUID: cpuID}
	}
	return cpuSlice
}
