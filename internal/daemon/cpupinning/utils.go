package cpupinning

import "k8s.io/klog/v2"

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

func ParseCGroupDriver(driver string) CGroupDriver {
	val, ok := map[string]CGroupDriver{
		"systemd":  DriverSystemd,
		"cgroupfs": DriverCgroupfs,
	}[driver]
	if !ok {
		klog.Fatalf("unknown cgroup driver %s", driver)
	}
	return val
}
