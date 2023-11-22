package nodecputopology

// Thread is a logical CPU core of the parent core
type Thread struct {
	Id int `json:"id"`
}

// Core is a physical CPU core of the parent socket
type Core struct {
	Id      int `json:"id"`
	Threads int `json:"threads"`
}

// Socket is a CPU socket of the parent NUMA node
type Socket struct {
	Id    int           `json:"id"`
	Cores map[int]*Core `json:"cores"`
}

// NumaNode is a NUMA node of the Kubernetes node
type NumaNode struct {
	Id      int             `json:"id"`
	Sockets map[int]*Socket `json:"sockets"`
}

// NodeCpuTopology represents the hierarchical topology of the CPU of a Kubernetes node
type NodeCpuTopology struct {
	NumaNodes map[int]*NumaNode `json:"numaNodes"`
}
