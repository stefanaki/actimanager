package v1alpha1

import (
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/pkg/nodecputopology"
)

func convertToV1Alpha1(t *nodecputopology.NodeCpuTopology) v1alpha1.CpuTopology {
	var topology v1alpha1.CpuTopology

	for _, socket := range t.Sockets {
		s := v1alpha1.Socket{
			SocketId:  socket.SocketId,
			NumaNodes: make([]v1alpha1.NumaNode, 0),
		}

		for _, numaNode := range socket.NumaNodes {
			n := v1alpha1.NumaNode{
				NumaNodeId: numaNode.NumaNodeId,
				Cores:      make([]v1alpha1.Core, 0),
			}

			for _, core := range numaNode.Cores {
				c := v1alpha1.Core{
					CoreId: core.CoreId,
					Cpus:   make([]v1alpha1.Cpu, 0),
				}

				for _, cpu := range core.Cpus {
					c.Cpus = append(c.Cpus, v1alpha1.Cpu{CpuId: cpu.CpuId})
				}
				n.Cores = append(n.Cores, c)
			}
			s.NumaNodes = append(s.NumaNodes, n)
		}
		topology.Sockets = append(topology.Sockets, s)
	}

	return topology
}

func NodeCpuTopologyV1Alpha1(lscpuOutput string) (v1alpha1.CpuTopology, error) {
	topology := &nodecputopology.NodeCpuTopology{}
	err := nodecputopology.ParseNodeCpuTopology(topology, lscpuOutput)
	return convertToV1Alpha1(topology), err
}
