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
		ListCpus:  make([]int, 0),
	}
	for _, cpu := range topologyResponse.Cpus {
		cpuTopology.ListCpus = append(cpuTopology.ListCpus, int(cpu))
	}
	for _, socket := range topologyResponse.Sockets {
		socketIdStr := strconv.Itoa(int(socket.Id))
		s := v1alpha1.Socket{
			Cores:    make(map[string]v1alpha1.Core),
			ListCpus: make([]int, 0),
		}
		for _, core := range socket.Cores {
			coreIdStr := strconv.Itoa(int(core.Id))
			c := v1alpha1.Core{
				Cpus:     make(map[string]v1alpha1.Cpu),
				ListCpus: make([]int, 0),
			}
			for _, cpu := range core.Cpus {
				c.ListCpus = append(c.ListCpus, int(cpu))
				s.ListCpus = append(s.ListCpus, int(cpu))
				c.Cpus[strconv.Itoa(int(cpu))] = v1alpha1.Cpu{
					CpuId: int(cpu),
				}
			}
			s.Cores[coreIdStr] = c
		}
		cpuTopology.Sockets[socketIdStr] = s
	}
	for _, numaNode := range topologyResponse.NumaNodes {
		numaNodeIdStr := strconv.Itoa(int(numaNode.Id))
		n := v1alpha1.NumaNode{
			Cpus:     make(map[string]v1alpha1.Cpu),
			ListCpus: make([]int, 0),
		}
		for _, cpu := range numaNode.Cpus {
			n.ListCpus = append(n.ListCpus, int(cpu))
			n.Cpus[strconv.Itoa(int(cpu))] = v1alpha1.Cpu{
				CpuId: int(cpu),
			}
		}
		cpuTopology.NumaNodes[numaNodeIdStr] = n
	}
	return cpuTopology
}
