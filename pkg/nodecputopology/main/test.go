package main

import "cslab.ece.ntua.gr/actimanager/pkg/nodecputopology"

func main() {
	var topology nodecputopology.NodeCpuTopology

	nodecputopology.RetrieveNodeCpuTopology(&topology)

	nodecputopology.PrintTopology(&topology)
}
