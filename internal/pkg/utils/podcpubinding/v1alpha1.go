package podcpubinding

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	nct "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
)

// GetExclusiveCpusOfCpuBinding returns the exclusive CPUs
// for a given CPU binding based on its exclusiveness level
func GetExclusiveCpusOfCpuBinding(cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology) map[int]struct{} {
	exclusiveCpus := make(map[int]struct{})

	for _, cpu := range cpuBinding.Spec.CpuSet {
		_, coreId, socketId, numaId := nct.GetCpuParentInfo(&topology.Spec.Topology, cpu.CpuId)

		switch cpuBinding.Spec.ExclusivenessLevel {
		case "Cpu":
			exclusiveCpus[cpu.CpuId] = struct{}{}
		case "Core":
			for _, c := range nct.GetAllCpusInCore(&topology.Spec.Topology, coreId) {
				exclusiveCpus[c] = struct{}{}
			}
		case "Socket":
			for _, c := range nct.GetAllCpusInSocket(&topology.Spec.Topology, socketId) {
				exclusiveCpus[c] = struct{}{}
			}
		case "Numa":
			for _, c := range nct.GetAllCpusInNuma(&topology.Spec.Topology, numaId) {
				exclusiveCpus[c] = struct{}{}
			}
		default:

		}
	}

	return exclusiveCpus
}
