package nodecputopology

import "cslab.ece.ntua.gr/actimanager/api/v1alpha1"

// V1Alpha1 returns the v1alpha1 representation of the NodeCpuTopology.
// It converts the NodeCpuTopology struct into the v1alpha1.CpuTopology struct,
// mapping the Sockets, Cores, and Cpus accordingly.
// The returned v1alpha1.CpuTopology contains the converted data.
func (t *NodeCpuTopology) V1Alpha1() v1alpha1.CpuTopology {
	var topology v1alpha1.CpuTopology

	for _, socket := range t.Sockets {
		s := v1alpha1.Socket{
			SocketId: socket.SocketId,
			Cores:    make([]v1alpha1.Core, 0),
		}

		for _, core := range socket.Cores {
			c := v1alpha1.Core{
				CoreId: core.CoreId,
				Cpus:   make([]v1alpha1.Cpu, 0),
			}

			for _, cpu := range core.Cpus {
				c.Cpus = append(c.Cpus, v1alpha1.Cpu{CpuId: cpu.CpuId})
			}

			s.Cores = append(s.Cores, c)
		}
		topology.Sockets = append(topology.Sockets, s)
	}

	for _, numa := range t.NumaNodes {
		n := v1alpha1.NumaNode{NumaNodeId: numa.NumaNodeId, Cpus: make([]v1alpha1.Cpu, 0)}
		for _, cpu := range numa.Cpus {
			n.Cpus = append(n.Cpus, v1alpha1.Cpu{CpuId: cpu.CpuId})
		}

		topology.NumaNodes = append(topology.NumaNodes, n)
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
func GetCpuParentInfo(topology *v1alpha1.NodeCpuTopology, cpuId int) (int, int, int, int) {
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

// GetAllCpusInCore returns a list of CPU IDs belonging to the specified core ID in the given NodeCpuTopology.
func GetAllCpusInCore(topology *v1alpha1.NodeCpuTopology, coreId int) []int {
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

// GetAllCpusInSocket returns a list of CPU IDs belonging to the specified socket in the given NodeCpuTopology.
func GetAllCpusInSocket(topology *v1alpha1.NodeCpuTopology, socketId int) []int {
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

// GetAllCpusInNuma returns a list of CPU IDs belonging to the specified NUMA node in the given NodeCpuTopology.
func GetAllCpusInNuma(topology *v1alpha1.NodeCpuTopology, numaId int) []int {
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

// GetNumaNodesOfCpuSet returns the list of NumaNodes that contain the given CPUs in the provided CpuTopology.
func GetNumaNodesOfCpuSet(cpus []v1alpha1.Cpu, topology v1alpha1.CpuTopology) []v1alpha1.NumaNode {
	numaNodesMap := make(map[int]struct{})
	var numaNodes []v1alpha1.NumaNode

	for _, numaNode := range topology.NumaNodes {
		for _, cpuInNuma := range numaNode.Cpus {
			for _, cpu := range cpus {
				if cpuInNuma.CpuId == cpu.CpuId {
					if _, exists := numaNodesMap[numaNode.NumaNodeId]; !exists {
						numaNodesMap[numaNode.NumaNodeId] = struct{}{}
						numaNodes = append(numaNodes, numaNode)
					}
					break
				}
			}
		}
	}

	return numaNodes
}
