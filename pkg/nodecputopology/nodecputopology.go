package nodecputopology

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var NodeCpuTopologyParseError = errors.New("could not parse node's CPU topology")

// ParseNodeCpuTopology uses the output of `lscpu` command to populate a
// `NodeCpuTopology` object of the CPU topology of the Kubernetes node
func ParseNodeCpuTopology(topology *NodeCpuTopology, lscpuOutput string) error {
	if topology.Sockets == nil {
		topology.Sockets = make(map[int]*Socket)
	}
	if topology.NumaNodes == nil {
		topology.NumaNodes = make(map[int]*NumaNode)
	}

	for _, lsLine := range strings.Split(strings.TrimSuffix(lscpuOutput, "\n"), "\n") {
		if strings.HasPrefix(lsLine, "#") {
			continue
		}

		values := strings.Split(lsLine, ",")
		if len(values) != 4 {
			fmt.Printf("Invalid format for socket,node,core,cpu: %s\n", lsLine)
			return NodeCpuTopologyParseError
		}

		socketId, err := strconv.Atoi(values[0])
		if err != nil {
			fmt.Printf("Could not parse socket ID: %v\n", err.Error())
			return NodeCpuTopologyParseError
		}

		nodeId, err := strconv.Atoi(values[1])
		if err != nil {
			fmt.Printf("Could not parse node ID: %v\n", err.Error())
			return NodeCpuTopologyParseError
		}

		coreId, err := strconv.Atoi(values[2])
		if err != nil {
			fmt.Printf("Could not parse core ID: %v\n", err.Error())
			return NodeCpuTopologyParseError
		}

		cpuId, err := strconv.Atoi(values[3])
		if err != nil {
			fmt.Printf("Could not parse cpu ID: %v\n", err.Error())
			return NodeCpuTopologyParseError
		}

		existingSocket, exists := topology.Sockets[socketId]
		if !exists {
			existingSocket = &Socket{SocketId: socketId, Cores: make(map[int]*Core)}
			topology.Sockets[socketId] = existingSocket
		}

		existingNumaNode, exists := topology.NumaNodes[nodeId]
		if !exists {
			existingNumaNode = &NumaNode{NumaNodeId: nodeId, Cpus: make([]*Cpu, 0)}
			topology.NumaNodes[nodeId] = existingNumaNode
		}

		existingCore, exists := topology.Sockets[socketId].Cores[coreId]
		if !exists {
			existingCore = &Core{CoreId: coreId, Cpus: make(map[int]*Cpu)}
			topology.Sockets[socketId].Cores[coreId] = existingCore
		}

		existingCpu, exists := topology.Sockets[socketId].Cores[coreId].Cpus[cpuId]
		if !exists {
			existingCpu = &Cpu{CpuId: cpuId}
			topology.Sockets[socketId].Cores[coreId].Cpus[cpuId] = existingCpu
		}

		topology.NumaNodes[nodeId].Cpus = append(topology.NumaNodes[nodeId].Cpus, existingCpu)
	}

	return nil
}

func PrintTopology(topology *NodeCpuTopology) {
	fmt.Println("NodeCpuTopology:")
	fmt.Println("Sockets:")
	for _, socket := range topology.Sockets {
		fmt.Printf("\tSocket ID: %d\n", socket.SocketId)
		for _, core := range socket.Cores {
			fmt.Printf("\t\tCore ID: %d\n", core.CoreId)
			for _, cpu := range core.Cpus {
				fmt.Printf("\t\t\tCPU: %d\n", cpu.CpuId)
			}
		}
	}
	fmt.Println("NUMA:")
	for _, numa := range topology.NumaNodes {
		fmt.Printf("\tNUMA ID: %d\n", numa.NumaNodeId)
		fmt.Printf("\t")
		for _, cpu := range numa.Cpus {
			fmt.Printf("%d,", cpu.CpuId)
		}
		fmt.Printf("\n")
	}
}
