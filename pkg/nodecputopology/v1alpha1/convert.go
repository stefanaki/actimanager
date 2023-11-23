package v1alpha1

import (
	"cslab.ece.ntua.gr/actimanager/api/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/pkg/nodecputopology"
)

func convertToV1Alpha1(internalTopology *nodecputopology.NodeCpuTopology) v1alpha1.CpuTopology {
	var topology v1alpha1.CpuTopology

	for _, numaNode := range internalTopology.NumaNodes {
		n := v1alpha1.NumaNode{
			Id:      numaNode.Id,
			Sockets: make([]v1alpha1.Socket, 0),
		}
		for _, socket := range numaNode.Sockets {
			s := v1alpha1.Socket{
				Id:    socket.Id,
				Cores: make([]v1alpha1.Core, 0),
			}
			for _, core := range socket.Cores {
				c := v1alpha1.Core{
					Id:      core.Id,
					Threads: make([]v1alpha1.Thread, 0),
				}
				for _, thread := range core.Threads {
					t := v1alpha1.Thread{Id: thread.Id}
					c.Threads = append(c.Threads, t)
				}
				s.Cores = append(s.Cores, c)
			}
			n.Sockets = append(n.Sockets, s)
		}
		topology.NumaNodes = append(topology.NumaNodes, n)
	}
	return topology
}

func NodeCpuTopologyV1Alpha1(lscpuOutput string) (v1alpha1.CpuTopology, error) {
	topology := &nodecputopology.NodeCpuTopology{}
	err := nodecputopology.ParseNodeCpuTopology(topology, lscpuOutput)
	nodecputopology.PrintTopology(topology)
	return convertToV1Alpha1(topology), err
}
