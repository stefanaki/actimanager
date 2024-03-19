package cpupinning

import (
	cgroupsctrl "cslab.ece.ntua.gr/actimanager/internal/pkg/cgroups"
	"fmt"
	"k8s.io/klog/v2"
	"strings"
)

const (
	MinShares      = 2
	MaxShares      = 262144
	SharesPerCPU   = 1024
	MilliCPUToCPU  = 1000
	QuotaPeriod    = 100000
	MinQuotaPeriod = 1000
)

// QosFromResources returns QoS class based on resource requests and limits of a pod.
func QosFromResources(resources ResourceInfo) QoS {
	limitCPU := resources.LimitCPUs
	requestCPU := resources.RequestedCPUs
	limitMemory := resources.LimitMemory
	requestMemory := resources.RequestedMemory
	if (limitCPU > 0 || requestCPU > 0) || (limitMemory != "0" || requestMemory != "0") {
		if limitCPU > 0 && requestCPU > 0 && limitCPU == requestCPU &&
			limitMemory != "0" && requestMemory != "0" && limitMemory == requestMemory {
			return Guaranteed
		}
		return Burstable
	}
	return BestEffort
}

// MilliCPUToQuota converts milliCPU to CFS quota and period values.
// Input parameters and resulting value is number of microseconds.
func MilliCPUToQuota(milliCPU int64, period int64) int64 {
	// CFS quota is measured in two values:
	//  - cfs_period_us=100ms (the amount of time to measure usage across given by period)
	//  - cfs_quota=20ms (the amount of cpu time allowed to be used across a period)
	// so in the above example, you are limited to 20% of a single CPU
	// for multi-cpu environments, you just scale equivalent amounts
	// see https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt for details
	if milliCPU == 0 {
		return 0
	}
	// we then convert your milliCPU to a value normalized over a period
	quota := (milliCPU * period) / MilliCPUToCPU
	// quota needs to be a minimum of 1ms.
	if quota < MinQuotaPeriod {
		quota = MinQuotaPeriod
	}
	return quota
}

// MilliCPUToShares converts the milliCPU to CFS shares.
func MilliCPUToShares(milliCPU int64) uint64 {
	if milliCPU == 0 {
		// Docker converts zero milliCPU to unset, which maps to kernel default
		// for unset: 1024. Return 2 here to really match kernel default for
		// zero milliCPU.
		return MinShares
	}
	// Conceptually (milliCPU / milliCPUToCPU) * sharesPerCPU, but factored to improve rounding.
	shares := (milliCPU * SharesPerCPU) / MilliCPUToCPU
	if shares < MinShares {
		return MinShares
	}
	if shares > MaxShares {
		return MaxShares
	}
	return uint64(shares)
}

func ParseRuntime(runtime string) ContainerRuntime {
	val, ok := map[string]ContainerRuntime{
		"containerd": ContainerdRunc,
		"kind":       Kind,
		"docker":     Docker,
	}[runtime]
	if !ok {
		klog.Fatalf("unknown runtime %s", runtime)
	}
	return val
}

func ParseCgroupsDriver(driver string) cgroupsctrl.CgroupsDriver {
	val, ok := map[string]cgroupsctrl.CgroupsDriver{
		"systemd":  cgroupsctrl.DriverSystemd,
		"cgroupfs": cgroupsctrl.DriverCgroupfs,
	}[driver]
	if !ok {
		klog.Fatalf("unknown cgroups driver %s", driver)
	}
	return val
}

func ConvertIntSliceToString[T int | int32 | int64](slice []T) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(slice)), ","), "[]")
}
