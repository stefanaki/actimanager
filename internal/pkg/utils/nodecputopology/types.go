package nodecputopology

// Cpu is a logical CPU core of the parent core
type Cpu struct {
	CpuId int `json:"cpuId"`
}

// Core is a physical CPU core of the parent socket
type Core struct {
	CoreId   int          `json:"coreId"`
	Cpus     map[int]*Cpu `json:"cpus"`
	ListCpus []int        `json:"listCpus"`
}

// Socket is a CPU socket of the Kubernetes node
type Socket struct {
	SocketId int           `json:"socketId"`
	Cores    map[int]*Core `json:"cores"`
	ListCpus []int         `json:"listCpus"`
}

// NumaNode is a NUMA node of the Kubernetes node
type NumaNode struct {
	NumaNodeId int    `json:"numaNodeId"`
	Cpus       []*Cpu `json:"cpus"`
	ListCpus   []int  `json:"listCpus"`
}

// NodeCpuTopology represents the hierarchical topology of the CPU of a Kubernetes node
type NodeCpuTopology struct {
	Sockets   map[int]*Socket   `json:"sockets"`
	NumaNodes map[int]*NumaNode `json:"numaNodes"`
	ListCpus  []int             `json:"listCpus"`
}
