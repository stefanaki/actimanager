package cpupinning

type ResourceInfo struct {
	RequestedCpus   int64
	LimitCpus       int64
	RequestedMemory string
	LimitMemory     string
}

// ContainerInfo Represents a container in the Daemon.
type ContainerInfo struct {
	CID       string
	PID       string
	Name      string
	Resources ResourceInfo
}

// QoS pod and containers quality of service type.
type QoS int

// QoS classes as defined in K8s.
const (
	Guaranteed QoS = iota
	BestEffort
	Burstable
)

// ContainerRuntime represents different CRI used by k8s.
type ContainerRuntime int

// Supported runtimes.
const (
	Docker ContainerRuntime = iota
	ContainerdRunc
	Kind
)

func (cr ContainerRuntime) String() string {
	return []string{
		"Docker",
		"Containerd+Runc",
		"Kind",
	}[cr]
}
