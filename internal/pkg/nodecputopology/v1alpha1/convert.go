package v1alpha1

import (
	apiv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/nodecputopology"
)

func convertToV1Alpha1(t *nodecputopology.NodeCpuTopology) apiv1alpha1.CpuTopology {
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
	topology := &nodecputopology.NodeCpuTopology{}
	err := nodecputopology.ParseNodeCpuTopology(topology, lscpuOutput)
	return convertToV1Alpha1(topology), err
}
