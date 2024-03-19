package nodecputopology

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	pbtopo "cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/topology"
	"strconv"
)

func TopologyToV1Alpha1(topologyResponse *pbtopo.TopologyResponse) *v1alpha1.CpuTopology {
	cpuTopology := &v1alpha1.CpuTopology{
		Sockets:   make(map[string]v1alpha1.Socket),
		NumaNodes: make(map[string]v1alpha1.NumaNode),
		Cpus:      make([]int, 0),
	}
	for _, cpu := range topologyResponse.Cpus {
		cpuTopology.Cpus = append(cpuTopology.Cpus, int(cpu))
	}
	for _, socket := range topologyResponse.Sockets {
		socketIdStr := strconv.Itoa(int(socket.Id))
		s := v1alpha1.Socket{
			Cores: make(map[string]v1alpha1.Core),
			Cpus:  make([]int, 0),
		}
		for _, core := range socket.Cores {
			coreIdStr := strconv.Itoa(int(core.Id))
			c := v1alpha1.Core{
				Cpus: make([]int, 0),
			}
			for _, cpu := range core.Cpus {
				c.Cpus = append(c.Cpus, int(cpu))
				s.Cpus = append(s.Cpus, int(cpu))
			}
			s.Cores[coreIdStr] = c
		}
		cpuTopology.Sockets[socketIdStr] = s
	}
	for _, numaNode := range topologyResponse.NumaNodes {
		numaNodeIdStr := strconv.Itoa(int(numaNode.Id))
		n := v1alpha1.NumaNode{
			Cpus: make([]int, 0),
		}
		for _, cpu := range numaNode.Cpus {
			n.Cpus = append(n.Cpus, int(cpu))
		}
		cpuTopology.NumaNodes[numaNodeIdStr] = n
	}
	return cpuTopology
}
