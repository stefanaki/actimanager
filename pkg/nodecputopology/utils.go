package nodecputopology

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

var NodeCpuTopologyParseError = errors.New("Could not parse node's CPU topology")

// lscpu runs the `lscpu` command in parsable format
func lscpu(ext string) (string, error) {
	return ExecCommand(exec.Command("lscpu", ext))
}

func ExecCommand(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return string(stdout.Bytes()), nil
}

func PrintTopology(topology *NodeCpuTopology) {
	fmt.Println("NodeCpuTopology:")
	for nodeID, numaNode := range topology.NumaNodes {
		fmt.Printf("\tNumaNode ID: %d\n", nodeID)
		for socketID, socket := range numaNode.Sockets {
			fmt.Printf("\t\tSocket ID: %d\n", socketID)
			for coreID, core := range socket.Cores {
				fmt.Printf("\t\t\tCore ID: %d\n", coreID)
				fmt.Printf("\t\t\t\tThreads: %d\n", core.Threads)
			}
		}
	}
}
