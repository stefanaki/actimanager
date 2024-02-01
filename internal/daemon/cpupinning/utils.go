package cpupinning

import (
	cgroupsctrl "cslab.ece.ntua.gr/actimanager/internal/pkg/cgroups"
	"k8s.io/klog/v2"
)

// QoSFromLimit returns QoS class based on resource requests and limits of a pod.
func QoSFromLimit(limitCpu, requestCpu int32, limitMemory, requestMemory string) QoS {
	if (limitCpu > 0 || requestCpu > 0) || (limitMemory != "0" || requestMemory != "0") {
		if limitCpu > 0 && requestCpu > 0 && limitCpu == requestCpu &&
			limitMemory != "0" && requestMemory != "0" && limitMemory == requestMemory {
			return Guaranteed
		}
		return Burstable
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
		klog.Fatalf("unknown cgroups driver %s", driver)
	}
	return val
}
