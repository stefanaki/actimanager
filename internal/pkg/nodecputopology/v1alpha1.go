package nodecputopology

import (
	"strconv"

	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
)

// V1Alpha1 returns the v1alpha1 representation of the NodeCpuTopology.
// It converts the NodeCpuTopology struct into the v1alpha1.CpuTopology struct,
// mapping the Sockets, Cores, and Cpus accordingly.
// The returned v1alpha1.CpuTopology contains the converted data.
func (t *NodeCpuTopology) V1Alpha1() v1alpha1.CpuTopology {
	topology := v1alpha1.CpuTopology{
		Sockets:   make(map[string]v1alpha1.Socket),
		NumaNodes: make(map[string]v1alpha1.NumaNode),
	}

	for socketId, socket := range t.Sockets {
		s := v1alpha1.Socket{
			Cores: make(map[string]v1alpha1.Core),
		}

		for coreId, core := range socket.Cores {
			c := v1alpha1.Core{
				Cpus: make(map[string]v1alpha1.Cpu),
			}

			for cpuId := range core.Cpus {
				c.Cpus[strconv.Itoa(cpuId)] = v1alpha1.Cpu{CpuId: cpuId}
			}

			s.Cores[strconv.Itoa(coreId)] = c
		}

		topology.Sockets[strconv.Itoa(socketId)] = s
	}

	for numaNodeId, numaNode := range t.NumaNodes {
		n := v1alpha1.NumaNode{Cpus: make(map[string]v1alpha1.Cpu)}

		for cpuId := range numaNode.Cpus {
			n.Cpus[strconv.Itoa(cpuId)] = v1alpha1.Cpu{CpuId: cpuId}
		}

		topology.NumaNodes[strconv.Itoa(numaNodeId)] = n
	}

	return topology
}

// NewV1Alpha1CpuTopologyFromLscpuOutput creates a v1alpha1.CpuTopology object
// from the output of the `lscpu -p=socket,node,core,cpu` command.
// It parses the lscpu output and returns the v1alpha1.CpuTopology object along with any error encountered during parsing.
func NewV1Alpha1CpuTopologyFromLscpuOutput(lscpuOutput string) (v1alpha1.CpuTopology, error) {
	topology := &NodeCpuTopology{}
	err := ParseNodeCpuTopology(topology, lscpuOutput)
	return topology.V1Alpha1(), err
}

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
func GetCpuParentInfo(topology *v1alpha1.NodeCpuTopology, targetCpuId int) (string, string, string, string) {
	numaId := "-1"
	for numaNodeId, numaNode := range topology.Spec.Topology.NumaNodes {
		for _, cpu := range numaNode.Cpus {
			if cpu.CpuId == targetCpuId {
				numaId = numaNodeId
			}
		}
	}

	for socketId, socket := range topology.Spec.Topology.Sockets {
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
	var cpus []int
	for _, socket := range topology.Sockets {
		for coreId, core := range socket.Cores {
			if coreId == targetCoreId {
				for _, cpu := range core.Cpus {
					cpus = append(cpus, cpu.CpuId)
				}
				return cpus
			}
		}
	}
	return cpus
}

// GetAllCpusInSocket returns a map of CPU IDs belonging to the specified socket in the given NodeCpuTopology.
func GetAllCpusInSocket(topology *v1alpha1.CpuTopology, targetSocketId string) []int {
	var cpus []int
	for socketId, socket := range topology.Sockets {
		if socketId == targetSocketId {
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

// GetAllCpusInNuma returns a map of CPU IDs belonging to the specified NUMA node in the given NodeCpuTopology.
func GetAllCpusInNuma(topology *v1alpha1.CpuTopology, targetNumaId string) []int {
	var cpus []int
	for numaNodeId, numaNode := range topology.NumaNodes {
		if numaNodeId == targetNumaId {
			for _, cpu := range numaNode.Cpus {
				cpus = append(cpus, cpu.CpuId)
			}

			return cpus
		}
	}
	return cpus
}

// GetNumaNodesOfCpuSet returns a map of NumaNodes that contain the given CPUs in the provided CpuTopology.
func GetNumaNodesOfCpuSet(cpus []v1alpha1.Cpu, topology v1alpha1.CpuTopology) map[string]v1alpha1.NumaNode {
	numaNodes := make(map[string]v1alpha1.NumaNode)

	for numaNodeId, numaNode := range topology.NumaNodes {
		for _, cpuInNuma := range numaNode.Cpus {
			for _, cpu := range cpus {
				if cpuInNuma.CpuId == cpu.CpuId {
					numaNodes[numaNodeId] = numaNode
					break
				}
			}
		}
	}

	return numaNodes
}

func GetTotalCpusCount(topology v1alpha1.CpuTopology) int {
	count := 0
	for _, socket := range topology.Sockets {
		for _, core := range socket.Cores {
			for range core.Cpus {
				count++
			}
		}
	}
	return count
}
