package nodecputopology

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	pbtopo "cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/topology"
	"strconv"
)

func TopologyToV1Alpha1(topologyResponse *pbtopo.TopologyResponse) *v1alpha1.CPUTopology {
	cpuTopology := &v1alpha1.CPUTopology{
		Sockets:   make(map[string]v1alpha1.Socket),
		NUMANodes: make(map[string]v1alpha1.NUMANode),
		CPUs:      make([]int, 0),
	}
	for _, cpu := range topologyResponse.Cpus {
		cpuTopology.CPUs = append(cpuTopology.CPUs, int(cpu))
	}
	for _, socket := range topologyResponse.Sockets {
		socketIdStr := strconv.Itoa(int(socket.Id))
		s := v1alpha1.Socket{
			Cores: make(map[string]v1alpha1.Core),
			CPUs:  make([]int, 0),
		}
		for _, core := range socket.Cores {
			coreIdStr := strconv.Itoa(int(core.Id))
			c := v1alpha1.Core{
				CPUs: make([]int, 0),
			}
			for _, cpu := range core.Cpus {
				c.CPUs = append(c.CPUs, int(cpu))
				s.CPUs = append(s.CPUs, int(cpu))
			}
			s.Cores[coreIdStr] = c
		}
		cpuTopology.Sockets[socketIdStr] = s
	}
	for _, numaNode := range topologyResponse.NumaNodes {
		numaNodeIdStr := strconv.Itoa(int(numaNode.Id))
		n := v1alpha1.NUMANode{
			CPUs: make([]int, 0),
		}
		for _, cpu := range numaNode.Cpus {
			n.CPUs = append(n.CPUs, int(cpu))
		}
		cpuTopology.NUMANodes[numaNodeIdStr] = n
	}
	return cpuTopology
}
