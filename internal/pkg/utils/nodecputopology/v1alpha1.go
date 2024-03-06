package nodecputopology

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"golang.org/x/exp/maps"
	"sort"
	"strconv"
)

// IsCpuSetInTopology checks if a given CPU set is present in the provided topology.
func IsCpuSetInTopology(topology *v1alpha1.CpuTopology, cpuSet []v1alpha1.Cpu) bool {
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

// GetCpuParentInfo returns the CPU parent information for the given CPU ID in the provided NodeCpuTopology.
// It searches for the CPU ID in the NumaNodes and Sockets of the topology and returns the corresponding CPU ID, Core ID, Socket ID, and Numa ID.
// If the CPU ID is not found, it returns -1 for all the values.
func GetCpuParentInfo(topology *v1alpha1.CpuTopology, targetCpuId int) (string, string, string, string) {
	numaId := "-1"
	for numaNodeId, numaNode := range topology.NumaNodes {
		for _, cpu := range numaNode.Cpus {
			if cpu.CpuId == targetCpuId {
				numaId = numaNodeId
			}
		}
	}
	for socketId, socket := range topology.Sockets {
		for coreId, core := range socket.Cores {
			for cid, cpu := range core.Cpus {
				if cpu.CpuId == targetCpuId {
					return cid, coreId, socketId, numaId
				}
			}
		}
	}
	return "-1", "-1", "-1", "-1"
}

// GetAllCpusInCore returns a map of CPU IDs belonging to the specified core ID in the given NodeCpuTopology.
func GetAllCpusInCore(topology *v1alpha1.CpuTopology, targetCoreId string) []int {
	for _, socket := range topology.Sockets {
		for coreId, core := range socket.Cores {
			if coreId == targetCoreId {
				return core.ListCpus
			}
		}
	}
	return []int{}
}

// GetAllCpusInSocket returns a map of CPU IDs belonging to the specified socket in the given NodeCpuTopology.
func GetAllCpusInSocket(topology *v1alpha1.CpuTopology, targetSocketId string) []int {
	for socketId, socket := range topology.Sockets {
		if socketId == targetSocketId {
			return socket.ListCpus
		}
	}
	return []int{}
}

// GetAllCpusInNuma returns a map of CPU IDs belonging to the specified NUMA node in the given NodeCpuTopology.
func GetAllCpusInNuma(topology *v1alpha1.CpuTopology, targetNumaId string) []int {
	for numaNodeId, numaNode := range topology.NumaNodes {
		if numaNodeId == targetNumaId {
			return numaNode.ListCpus
		}
	}
	return []int{}
}

// GetNumaNodesOfCpuSet returns a map of NumaNodes that contain the given CPUs in the provided CpuTopology.
func GetNumaNodesOfCpuSet(cpus []int, topology v1alpha1.CpuTopology) []int {
	numaNodes := make(map[int]struct{})
	for numaNodeId, numaNode := range topology.NumaNodes {
		for _, cpuInNuma := range numaNode.Cpus {
			for _, cpu := range cpus {
				if cpuInNuma.CpuId == cpu {
					id, _ := strconv.Atoi(numaNodeId)
					numaNodes[id] = struct{}{}
					break
				}
			}
		}
	}
	nodeSlice := maps.Keys(numaNodes)
	sort.Ints(nodeSlice)
	return nodeSlice
}

func GetTotalCpusCount(topology v1alpha1.CpuTopology) int {
	return len(topology.ListCpus)
}

func DeleteCpuFromTopology(topology *v1alpha1.CpuTopology, cpuId int) {
	for _, socket := range topology.Sockets {
		for _, core := range socket.Cores {
			for i, cpu := range core.ListCpus {
				if cpu == cpuId {
					delete(core.Cpus, strconv.Itoa(cpuId))
					core.ListCpus = append(core.ListCpus[:i], core.ListCpus[i+1:]...)
					break
				}
			}
		}
	}
	for _, numaNode := range topology.NumaNodes {
		for i, cpu := range numaNode.ListCpus {
			if cpu == cpuId {
				delete(numaNode.Cpus, strconv.Itoa(cpuId))
				numaNode.ListCpus = append(numaNode.ListCpus[:i], numaNode.ListCpus[i+1:]...)
				break
			}
		}
	}
	for i, cpu := range topology.ListCpus {
		if cpu == cpuId {
			topology.ListCpus = append(topology.ListCpus[:i], topology.ListCpus[i+1:]...)
			break
		}
	}
}

func GetAvailableResources(exclusivenessLevel string, feasibleCpus v1alpha1.CpuTopology, topology v1alpha1.CpuTopology) []int {
	switch exclusivenessLevel {
	case "Cpu":
		return feasibleCpus.ListCpus
	case "Core":
		cores := make([]int, 0)
		for socketId, socket := range feasibleCpus.Sockets {
			for coreId, core := range socket.Cores {
				if len(topology.Sockets[socketId].Cores[coreId].Cpus) > len(core.Cpus) {
					continue
				}
				id, _ := strconv.Atoi(coreId)
				cores = append(cores, id)
			}
		}
		return cores
	case "Socket":
		sockets := make([]int, 0)
		for socketId := range feasibleCpus.Sockets {
			socketCpus := GetAllCpusInSocket(&feasibleCpus, socketId)
			if len(topology.Sockets[socketId].Cores) > len(socketCpus) {
				continue
			}
			id, _ := strconv.Atoi(socketId)
			sockets = append(sockets, id)
		}
		return sockets
	case "Numa":
		numas := make([]int, 0)
		for numaId := range feasibleCpus.NumaNodes {
			numaCpus := GetAllCpusInNuma(&feasibleCpus, numaId)
			if len(topology.NumaNodes[numaId].Cpus) > len(numaCpus) {
				continue
			}
			id, _ := strconv.Atoi(numaId)
			numas = append(numas, id)
		}
		return numas
	}
	return []int{}
}
