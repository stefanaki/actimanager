package v1alpha1

import (
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
)

func IsCpuSetInTopology(topology *cslabecentuagrv1alpha1.CpuTopology, cpuSet []cslabecentuagrv1alpha1.Cpu) bool {
	remaining := len(cpuSet)

	for _, socket := range topology.Sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.Cpus {
				for _, inputCpu := range cpuSet {
					if cpu.CpuId == inputCpu.CpuId {
						remaining--
						break
					}
				}
			}
		}
	}

	return remaining == 0
}

func GetCpuParentInfo(topology *cslabecentuagrv1alpha1.NodeCpuTopology, cpuId int) (int, int, int, int) {
	numaId := -1
	for _, numa := range topology.Spec.Topology.NumaNodes {
		for _, cpu := range numa.Cpus {
			if cpu.CpuId == cpuId {
				numaId = numa.NumaNodeId
			}
		}
	}

	for _, socket := range topology.Spec.Topology.Sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.Cpus {
				if cpu.CpuId == cpuId {
					return cpu.CpuId, core.CoreId, socket.SocketId, numaId
				}
			}
		}
	}

	return -1, -1, -1, -1
}

func GetAllCpusInCore(topology *cslabecentuagrv1alpha1.NodeCpuTopology, coreId int) []int {
	var cpus []int
	for _, socket := range topology.Spec.Topology.Sockets {
		for _, core := range socket.Cores {
			if core.CoreId == coreId {
				for _, cpu := range core.Cpus {
					cpus = append(cpus, cpu.CpuId)
				}
				return cpus
			}
		}
	}
	return cpus
}

func GetAllCpusInSocket(topology *cslabecentuagrv1alpha1.NodeCpuTopology, socketId int) []int {
	var cpus []int
	for _, socket := range topology.Spec.Topology.Sockets {
		if socket.SocketId == socketId {
			for _, core := range socket.Cores {
				for _, cpu := range core.Cpus {
					cpus = append(cpus, cpu.CpuId)
				}
			}
			return cpus
		}
	}
	return cpus
}

func GetAllCpusInNuma(topology *cslabecentuagrv1alpha1.NodeCpuTopology, numaId int) []int {
	var cpus []int
	for _, numaNode := range topology.Spec.Topology.NumaNodes {
		if numaNode.NumaNodeId == numaId {
			for _, cpu := range numaNode.Cpus {
				cpus = append(cpus, cpu.CpuId)
			}

			return cpus
		}
	}
	return cpus
}
