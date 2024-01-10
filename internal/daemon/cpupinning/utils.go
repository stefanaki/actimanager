package cpupinning

import (
	cgroupsctrl "cslab.ece.ntua.gr/actimanager/internal/pkg/cgroups"
	"k8s.io/klog/v2"
)

// QoSFromLimit returns QoS class based on limits set on pod cpu.
func QoSFromLimit[T int | int32 | int64](limitCpu, requestCpu T) QoS {
	if limitCpu > 0 || requestCpu > 0 {
		if limitCpu == requestCpu {
			return Guaranteed
		}
		if requestCpu < limitCpu {
			return Burstable
		}
	}
	return BestEffort
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
		klog.Fatalf("unknown cgroups1 driver %s", driver)
	}
	return val
}
