package v1alpha1

import (
	apiv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nct "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology"
)

func ConvertToV1Alpha1(t *nct.NodeCpuTopology) apiv1alpha1.CpuTopology {
	var topology apiv1alpha1.CpuTopology

	for _, socket := range t.Sockets {
		s := apiv1alpha1.Socket{
			SocketId: socket.SocketId,
			Cores:    make([]apiv1alpha1.Core, 0),
		}

		for _, core := range socket.Cores {
			c := apiv1alpha1.Core{
				CoreId: core.CoreId,
				Cpus:   make([]apiv1alpha1.Cpu, 0),
			}

			for _, cpu := range core.Cpus {
				c.Cpus = append(c.Cpus, apiv1alpha1.Cpu{CpuId: cpu.CpuId})
			}

			s.Cores = append(s.Cores, c)
		}
		topology.Sockets = append(topology.Sockets, s)
	}

	for _, numa := range t.NumaNodes {
		n := apiv1alpha1.NumaNode{NumaNodeId: numa.NumaNodeId, Cpus: make([]apiv1alpha1.Cpu, 0)}
		for _, cpu := range numa.Cpus {
			n.Cpus = append(n.Cpus, apiv1alpha1.Cpu{CpuId: cpu.CpuId})
		}

		topology.NumaNodes = append(topology.NumaNodes, n)
	}

	return topology
}

func NodeCpuTopologyV1Alpha1(lscpuOutput string) (apiv1alpha1.CpuTopology, error) {
	topology := &nct.NodeCpuTopology{}
	err := nct.ParseNodeCpuTopology(topology, lscpuOutput)
	return ConvertToV1Alpha1(topology), err
}

func IsCpuSetInTopology(topology *apiv1alpha1.CpuTopology, cpuSet []apiv1alpha1.Cpu) bool {
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

func GetCpuParentInfo(topology *apiv1alpha1.NodeCpuTopology, cpuId int) (int, int, int, int) {
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

func GetAllCpusInCore(topology *apiv1alpha1.NodeCpuTopology, coreId int) []int {
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

func GetAllCpusInSocket(topology *apiv1alpha1.NodeCpuTopology, socketId int) []int {
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

func GetAllCpusInNuma(topology *apiv1alpha1.NodeCpuTopology, numaId int) []int {
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

func GetNumaNodesOfCpuSet(cpus []apiv1alpha1.Cpu, topology apiv1alpha1.CpuTopology) []apiv1alpha1.NumaNode {
	numaNodesMap := make(map[int]struct{})
	var numaNodes []apiv1alpha1.NumaNode

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
