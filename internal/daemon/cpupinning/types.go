package cpupinning

const (
	MinShares      = 2      // Minimum shares allowed by CFS.
	MaxShares      = 262144 // Maximum shares allowed by CFS.
	SharesPerCPU   = 1024   // Shares per CPU.
	MilliCPUToCPU  = 1000   // MilliCPU to CPU conversion factor.
	QuotaPeriod    = 100000 // Quota period in microseconds.
	MinQuotaPeriod = 1000   // Minimum quota period allowed by CFS.
)

// ResourceInfo represents the information about the requested resources for CPU pinning.
type ResourceInfo struct {
	RequestedCPUs   int64  // The number of CPUs requested for pinning.
	LimitCPUs       int64  // The maximum number of CPUs allowed for pinning.
	RequestedMemory string // The amount of memory requested for pinning.
	LimitMemory     string // The maximum amount of memory allowed for pinning.
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
