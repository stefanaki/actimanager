package nodecputopology

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var NodeCpuTopologyParseError = errors.New("could not parse node's CPU topology")

// ParseNodeCpuTopology uses the output of `lscpu -p=socket,node,core,cpu` command
// to populate a `NodeCpuTopology` object of the CPU topology of the Kubernetes node
func ParseNodeCpuTopology(lscpuOutput string) (*NodeCpuTopology, error) {
	t := &NodeCpuTopology{
		Sockets:   make(map[int]*Socket),
		NumaNodes: make(map[int]*NumaNode),
		ListCpus:  make([]int, 0),
	}

	for _, lsLine := range strings.Split(strings.TrimSuffix(lscpuOutput, "\n"), "\n") {
		if strings.HasPrefix(lsLine, "#") {
			continue
		}

		values := strings.Split(lsLine, ",")
		if len(values) != 4 {
			fmt.Printf("Invalid format for socket,node,core,cpu: %s\n", lsLine)
			return nil, NodeCpuTopologyParseError
		}

		socketId, err := strconv.Atoi(values[0])
		if err != nil {
			fmt.Printf("Could not parse socket ID: %v\n", err.Error())
			return nil, NodeCpuTopologyParseError
		}

		nodeId, err := strconv.Atoi(values[1])
		if err != nil {
			fmt.Printf("Could not parse node ID: %v\n", err.Error())
			return nil, NodeCpuTopologyParseError
		}

		coreId, err := strconv.Atoi(values[2])
		if err != nil {
			fmt.Printf("Could not parse core ID: %v\n", err.Error())
			return nil, NodeCpuTopologyParseError
		}

		cpuId, err := strconv.Atoi(values[3])
		if err != nil {
			fmt.Printf("Could not parse cpu ID: %v\n", err.Error())
			return nil, NodeCpuTopologyParseError
		}

		socket, exists := t.Sockets[socketId]
		if !exists {
			socket = &Socket{
				SocketId: socketId,
				Cores:    make(map[int]*Core),
				ListCpus: make([]int, 0),
			}
			t.Sockets[socketId] = socket
		}

		numaNode, exists := t.NumaNodes[nodeId]
		if !exists {
			numaNode = &NumaNode{
				NumaNodeId: nodeId,
				Cpus:       make([]*Cpu, 0),
				ListCpus:   make([]int, 0),
			}
			t.NumaNodes[nodeId] = numaNode
		}

		core, exists := t.Sockets[socketId].Cores[coreId]
		if !exists {
			core = &Core{
				CoreId:   coreId,
				Cpus:     make(map[int]*Cpu),
				ListCpus: make([]int, 0),
			}
			t.Sockets[socketId].Cores[coreId] = core
		}

		cpu, exists := t.Sockets[socketId].Cores[coreId].Cpus[cpuId]
		if !exists {
			cpu = &Cpu{CpuId: cpuId}
			t.Sockets[socketId].Cores[coreId].Cpus[cpuId] = cpu
			t.Sockets[socketId].Cores[coreId].ListCpus = append(t.Sockets[socketId].Cores[coreId].ListCpus, cpuId)
			t.Sockets[socketId].ListCpus = append(t.Sockets[socketId].ListCpus, cpuId)
		}

		t.NumaNodes[nodeId].Cpus = append(t.NumaNodes[nodeId].Cpus, cpu)
		t.NumaNodes[nodeId].ListCpus = append(t.NumaNodes[nodeId].ListCpus, cpuId)
		t.ListCpus = append(t.ListCpus, cpuId)
	}

	return t, nil
}

func PrintTopology(topology *NodeCpuTopology) {
	fmt.Println("NodeCpuTopology:")
	fmt.Println("Sockets:")
	for _, socket := range topology.Sockets {
		fmt.Printf("\tSocket ID: %d\n", socket.SocketId)
		fmt.Printf("ListCpus: %v\n", socket.ListCpus)
		for _, core := range socket.Cores {
			fmt.Printf("\t\tCore ID: %d\n", core.CoreId)
			fmt.Printf("\t\tListCpus: %v\n", core.ListCpus)
			for _, cpu := range core.Cpus {
				fmt.Printf("\t\t\tCPU: %d\n", cpu.CpuId)
			}
		}
	}
	fmt.Println("NUMA:")
	for _, numa := range topology.NumaNodes {
		fmt.Printf("\tNUMA ID: %d\n", numa.NumaNodeId)
		fmt.Printf("ListCpus: %v\n", numa.ListCpus)
		fmt.Printf("\t")
		for _, cpu := range numa.Cpus {
			fmt.Printf("%d,", cpu.CpuId)
		}
		fmt.Printf("\n")
	}
}
