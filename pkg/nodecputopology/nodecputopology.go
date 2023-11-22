package nodecputopology

import (
	"fmt"
	"strconv"
	"strings"
)

// RetrieveNodeCpuTopology uses `lscpu` internally to populate a
// `NodeCpuTopology` object of the CPU topology of the Kubernetes node
func RetrieveNodeCpuTopology(topology *NodeCpuTopology) error {
	if topology.NumaNodes == nil {
		topology.NumaNodes = make(map[int]*NumaNode)
	}

	out, err := lscpu("-p=node,socket,core")
	if err != nil {
		fmt.Printf("Could not get NUMA nodes: %v\n", err.Error())
		return NodeCpuTopologyParseError
	}

	for _, lsLine := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if strings.HasPrefix(lsLine, "#") {
			continue
		}

		values := strings.Split(lsLine, ",")
		if len(values) != 3 {
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
			existingCore = &Core{Id: coreId, Threads: 1} // threads per core not implemented yet
			topology.NumaNodes[nodeId].Sockets[socketId].Cores[coreId] = existingCore
		}
	}

	return nil
}
