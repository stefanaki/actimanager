package podcpubinding

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	nct "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology"
	"errors"
)

// GetExclusiveCpusOfCpuBinding returns the exclusive CPUs
// for a given CPU binding based on its exclusiveness level
func GetExclusiveCpusOfCpuBinding(cpuBinding *v1alpha1.PodCpuBinding, topology *v1alpha1.NodeCpuTopology) map[int]struct{} {
	exclusiveCpus := make(map[int]struct{})

	for _, cpu := range cpuBinding.Spec.CpuSet {
		_, coreId, socketId, numaId := nct.GetCpuParentInfo(topology, cpu.CpuId)

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

func CalculateCpuSetForPod(requests int64, level string, feasibleCpus v1alpha1.CpuTopology, topology v1alpha1.CpuTopology) ([]int, error) {
	switch level {
	case "Cpu":
		return CalculateCpuSetForCpu(requests, feasibleCpus, topology)
	case "Core":
		return CalculateCpuSetForCore(requests, feasibleCpus, topology)
	case "Socket":
		return CalculateCpuSetForSocket(requests, feasibleCpus, topology)
	case "Numa":
		return CalculateCpuSetForNuma(requests, feasibleCpus, topology)
	default:
		return CalculateCpuSetForNone(feasibleCpus)
	}
}

func CalculateCpuSetForNone(cpus v1alpha1.CpuTopology) ([]int, error) {
	res := make([]int, 0)
	for _, socket := range cpus.Sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.Cpus {
				res = append(res, cpu.CpuId)
			}
		}
	}
	if len(res) == 0 {
		return nil, errors.New("No shared CPUs found")
	}
	return res, nil
}

func CalculateCpuSetForCpu(requests int64, cpus v1alpha1.CpuTopology, topology v1alpha1.CpuTopology) ([]int, error) {
	res := make([]int, 0)
	for _, socket := range cpus.Sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.Cpus {
				res = append(res, cpu.CpuId)
			}
		}
	}

	if int64(len(res)) < requests {
		return nil, errors.New("Not enough shared CPUs")
	}

	return res[0:requests], nil
}

func CalculateCpuSetForCore(requests int64, feasibleCpus v1alpha1.CpuTopology, topology v1alpha1.CpuTopology) ([]int, error) {
	cores := make([]v1alpha1.Core, 0)

	for socketId, socket := range feasibleCpus.Sockets {
		for coreId, core := range socket.Cores {
			if len(topology.Sockets[socketId].Cores[coreId].Cpus) > len(core.Cpus) {
				continue
			}
			cores = append(cores, core)
		}
	}
	cpus := make([]int, 0)
	for _, core := range cores {
		for _, cpu := range core.Cpus {
			cpus = append(cpus, cpu.CpuId)
		}
	}

	if int64(len(cpus)) < requests {
		return nil, errors.New("not enough cores")
	}

	return cpus[0:requests], nil
}

func CalculateCpuSetForSocket(requests int64, feasibleCpus v1alpha1.CpuTopology, topology v1alpha1.CpuTopology) ([]int, error) {
	sockets := make([]v1alpha1.Socket, 0)
	for socketId, socket := range feasibleCpus.Sockets {
		socketCpus := nct.GetAllCpusInSocket(&feasibleCpus, socketId)
		if len(topology.Sockets[socketId].Cores) > len(socketCpus) {
			continue
		}
		sockets = append(sockets, socket)
	}
	cpus := make([]int, 0)
	for _, socket := range sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.Cpus {
				cpus = append(cpus, cpu.CpuId)
			}
		}
	}
	if int64(len(cpus)) < requests {
		return nil, errors.New("not enough sockets")
	}
	return cpus[0:requests], nil
}

func CalculateCpuSetForNuma(requests int64, feasibleCpus v1alpha1.CpuTopology, topology v1alpha1.CpuTopology) ([]int, error) {
	numas := make([]v1alpha1.NumaNode, 0)
	for numaId, numa := range feasibleCpus.NumaNodes {
		numaCpus := nct.GetAllCpusInNuma(&feasibleCpus, numaId)
		if len(topology.NumaNodes[numaId].Cpus) > len(numaCpus) {
			continue
		}
		numas = append(numas, numa)
	}
	cpus := make([]int, 0)
	for _, numa := range numas {
		for _, cpu := range numa.Cpus {
			cpus = append(cpus, cpu.CpuId)
		}
	}
	if int64(len(cpus)) < requests {
		return nil, errors.New("not enough numas")
	}
	return cpus[0:requests], nil
}
