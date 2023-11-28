package v1alpha1

import (
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	nodecputopology2 "cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology"
)

func convertToV1Alpha1(t *nodecputopology2.NodeCpuTopology) v1alpha1.CpuTopology {
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

func NodeCpuTopologyV1Alpha1(lscpuOutput string) (v1alpha1.CpuTopology, error) {
	topology := &nodecputopology2.NodeCpuTopology{}
	err := nodecputopology2.ParseNodeCpuTopology(topology, lscpuOutput)
	return convertToV1Alpha1(topology), err
}
