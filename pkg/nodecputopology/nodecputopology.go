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
			existingSocket = &Socket{SocketId: socketId, NumaNodes: make(map[int]*NumaNode)}
			topology.Sockets[socketId] = existingSocket
		}

		existingNumaNode, exists := topology.Sockets[socketId].NumaNodes[nodeId]
		if !exists {
			existingNumaNode = &NumaNode{NumaNodeId: nodeId, Cores: make(map[int]*Core)}
			topology.Sockets[socketId].NumaNodes[nodeId] = existingNumaNode
		}

		existingCore, exists := topology.Sockets[socketId].NumaNodes[nodeId].Cores[coreId]
		if !exists {
			existingCore = &Core{CoreId: coreId, Cpus: make(map[int]*Cpu)}
			topology.Sockets[socketId].NumaNodes[nodeId].Cores[coreId] = existingCore
		}

		existingCpu, exists := topology.Sockets[socketId].NumaNodes[nodeId].Cores[coreId].Cpus[cpuId]
		if !exists {
			existingCpu = &Cpu{CpuId: cpuId}
			topology.Sockets[socketId].NumaNodes[nodeId].Cores[coreId].Cpus[cpuId] = existingCpu
		}
	}

	return nil
}

func PrintTopology(topology *NodeCpuTopology) {
	fmt.Println("NodeCpuTopology:")
	for _, socket := range topology.Sockets {
		fmt.Printf("\tSocket ID: %d\n", socket.SocketId)
		for _, numaNode := range socket.NumaNodes {
			fmt.Printf("\t\tNumaNode ID: %d\n", numaNode.NumaNodeId)
			for _, core := range numaNode.Cores {
				fmt.Printf("\t\t\tCore ID: %d\n", core.CoreId)
				for _, cpu := range core.Cpus {
					fmt.Printf("\t\t\t\tCPU: %d\n", cpu.CpuId)
				}
			}
		}
	}
}
