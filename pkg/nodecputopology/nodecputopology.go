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
	if topology.NumaNodes == nil {
		topology.NumaNodes = make(map[int]*NumaNode)
	}

	for _, lsLine := range strings.Split(strings.TrimSuffix(lscpuOutput, "\n"), "\n") {
		if strings.HasPrefix(lsLine, "#") {
			continue
		}

		values := strings.Split(lsLine, ",")
		if len(values) != 4 {
			fmt.Printf("Invalid format for node,socket,core line: %s\n", lsLine)
			return NodeCpuTopologyParseError
		}

		nodeId, err := strconv.Atoi(values[0])
		if err != nil {
			fmt.Printf("Could not parse node ID: %v\n", err.Error())
			return NodeCpuTopologyParseError
		}

		socketId, err := strconv.Atoi(values[1])
		if err != nil {
			fmt.Printf("Could not parse socket ID: %v\n", err.Error())
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

		existingNumaNode, exists := topology.NumaNodes[nodeId]
		if !exists {
			existingNumaNode = &NumaNode{Id: nodeId, Sockets: make(map[int]*Socket)}
			topology.NumaNodes[nodeId] = existingNumaNode
		}

		existingSocket, exists := topology.NumaNodes[nodeId].Sockets[socketId]
		if !exists {
			existingSocket = &Socket{Id: socketId, Cores: make(map[int]*Core)}
			topology.NumaNodes[nodeId].Sockets[socketId] = existingSocket
		}

		existingCore, exists := topology.NumaNodes[nodeId].Sockets[socketId].Cores[coreId]
		if !exists {
			existingCore = &Core{Id: coreId, Cpus: make(map[int]*Cpu)}
			topology.NumaNodes[nodeId].Sockets[socketId].Cores[coreId] = existingCore
		}

		existingCpu, exists := topology.NumaNodes[nodeId].Sockets[socketId].Cores[coreId].Cpus[cpuId]
		if !exists {
			existingCpu = &Cpu{Id: cpuId}
			topology.NumaNodes[nodeId].Sockets[socketId].Cores[coreId].Cpus[cpuId] = existingCpu
		}
	}

	return nil
}

func PrintTopology(topology *NodeCpuTopology) {
	fmt.Println("NodeCpuTopology:")
	for nodeID, numaNode := range topology.NumaNodes {
		fmt.Printf("\tNumaNode ID: %d\n", nodeID)
		for socketID, socket := range numaNode.Sockets {
			fmt.Printf("\t\tSocket ID: %d\n", socketID)
			for coreID, core := range socket.Cores {
				fmt.Printf("\t\t\tCore ID: %d\n", coreID)
				for cpu := range core.Cpus {
					fmt.Printf("\t\t\t\tCPU: %d", cpu)
				}
			}
		}
	}
}
