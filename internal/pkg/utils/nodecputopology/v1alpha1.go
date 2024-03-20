package nodecputopology

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"golang.org/x/exp/maps"
	"sort"
	"strconv"
)

// IsCPUSetInTopology checks if a given CPU set is present in the provided topology.
func IsCPUSetInTopology(topology *v1alpha1.CPUTopology, cpuSet []v1alpha1.CPU) bool {
	remaining := len(cpuSet)
	for _, socket := range topology.Sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.CPUs {
				for _, inputCPU := range cpuSet {
					if cpu == inputCPU.CPUID {
						remaining--
						break
					}
				}
			}
		}
	}
	return remaining == 0
}

// GetCPUParentInfo returns the CPU parent information for the given CPU ID in the provided NodeCPUTopology.
// It searches for the CPU ID in the NUMANodes and Sockets of the topology and returns the corresponding CPU ID, Core ID, Socket ID, and NUMA ID.
// If the CPU ID is not found, it returns -1 for all the values.
func GetCPUParentInfo(topology *v1alpha1.CPUTopology, targetCPUID int) (string, string, string, string) {
	numaID := "-1"
	for numaNodeID, numaNode := range topology.NUMANodes {
		for _, cpu := range numaNode.CPUs {
			if cpu == targetCPUID {
				numaID = numaNodeID
			}
		}
	}
	for socketID, socket := range topology.Sockets {
		for coreID, core := range socket.Cores {
			for _, cpu := range core.CPUs {
				if cpu == targetCPUID {
					return strconv.Itoa(cpu), coreID, socketID, numaID
				}
			}
		}
	}
	return "-1", "-1", "-1", "-1"
}

// GetAllCPUsInCore returns a map of CPU IDs belonging to the specified core ID in the given NodeCPUTopology.
func GetAllCPUsInCore(topology *v1alpha1.CPUTopology, targetCoreID string) []int {
	for _, socket := range topology.Sockets {
		for coreID, core := range socket.Cores {
			if coreID == targetCoreID {
				return core.CPUs
			}
		}
	}
	return []int{}
}

// GetAllCPUsInSocket returns a map of CPU IDs belonging to the specified socket in the given NodeCPUTopology.
func GetAllCPUsInSocket(topology *v1alpha1.CPUTopology, targetSocketID string) []int {
	for socketID, socket := range topology.Sockets {
		if socketID == targetSocketID {
			return socket.CPUs
		}
	}
	return []int{}
}

// GetAllCPUsInNUMA returns a map of CPU IDs belonging to the specified NUMA node in the given NodeCPUTopology.
func GetAllCPUsInNUMA(topology *v1alpha1.CPUTopology, targetNUMAID string) []int {
	for numaNodeID, numaNode := range topology.NUMANodes {
		if numaNodeID == targetNUMAID {
			return numaNode.CPUs
		}
	}
	return []int{}
}

// GetNUMANodesOfCPUSet returns a map of NUMANodes that contain the given CPUs in the provided CPUTopology.
func GetNUMANodesOfCPUSet(cpus []int, topology v1alpha1.CPUTopology) []int {
	numaNodes := make(map[int]struct{})
	for numaNodeID, numaNode := range topology.NUMANodes {
		for _, cpuInNUMA := range numaNode.CPUs {
			for _, cpu := range cpus {
				if cpuInNUMA == cpu {
					id, _ := strconv.Atoi(numaNodeID)
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

func GetTotalCPUsCount(topology v1alpha1.CPUTopology) int {
	return len(topology.CPUs)
}

func DeleteCPUFromTopology(topology *v1alpha1.CPUTopology, cpuID int) {
	for _, socket := range topology.Sockets {
		for _, core := range socket.Cores {
			for i, cpu := range core.CPUs {
				if cpu == cpuID {
					core.CPUs = append(core.CPUs[:i], core.CPUs[i+1:]...)
					break
				}
			}
		}
	}
	for _, numaNode := range topology.NUMANodes {
		for i, cpu := range numaNode.CPUs {
			if cpu == cpuID {
				numaNode.CPUs = append(numaNode.CPUs[:i], numaNode.CPUs[i+1:]...)
				break
			}
		}
	}
	for i, cpu := range topology.CPUs {
		if cpu == cpuID {
			topology.CPUs = append(topology.CPUs[:i], topology.CPUs[i+1:]...)
			break
		}
	}
}

func GetAvailableResources(level v1alpha1.ResourceLevel, feasibleCPUs v1alpha1.CPUTopology, topology v1alpha1.CPUTopology) []int {
	switch level {
	case v1alpha1.ResourceLevelCPU:
		return feasibleCPUs.CPUs
	case v1alpha1.ResourceLevelCore:
		cores := make([]int, 0)
		for socketID, socket := range feasibleCPUs.Sockets {
			for coreID, core := range socket.Cores {
				if len(topology.Sockets[socketID].Cores[coreID].CPUs) > len(core.CPUs) {
					continue
				}
				id, _ := strconv.Atoi(coreID)
				cores = append(cores, id)
			}
		}
		return cores
	case v1alpha1.ResourceLevelSocket:
		sockets := make([]int, 0)
		for socketID := range feasibleCPUs.Sockets {
			socketCPUs := GetAllCPUsInSocket(&feasibleCPUs, socketID)
			if len(topology.Sockets[socketID].Cores) > len(socketCPUs) {
				continue
			}
			id, _ := strconv.Atoi(socketID)
			sockets = append(sockets, id)
		}
		return sockets
	case v1alpha1.ResourceLevelNUMA:
		numas := make([]int, 0)
		for numaID := range feasibleCPUs.NUMANodes {
			numaCPUs := GetAllCPUsInNUMA(&feasibleCPUs, numaID)
			if len(topology.NUMANodes[numaID].CPUs) > len(numaCPUs) {
				continue
			}
			id, _ := strconv.Atoi(numaID)
			numas = append(numas, id)
		}
		return numas
	}
	return []int{}
}
