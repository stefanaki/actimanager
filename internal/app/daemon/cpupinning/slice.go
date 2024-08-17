package cpupinning

import (
	cgroupsctrl "cslab.ece.ntua.gr/actimanager/internal/pkg/cgroups"
	"fmt"
	"strings"
)

// SliceName returns path to container cgroup leaf slice in cgroupfs.
func SliceName(c ContainerInfo, r ContainerRuntime, d cgroupsctrl.CgroupsDriver) string {
	if r == Kind {
		return sliceNameKind(c)
	}
	if d == cgroupsctrl.DriverSystemd {
		return sliceNameDockerContainerdWithSystemd(c, r)
	}
	return sliceNameDockerContainerdWithCgroupfs(c, r)
}

func sliceNameKind(c ContainerInfo) string {
	podType := [3]string{"", "besteffort/", "burstable/"}
	return fmt.Sprintf(
		"kubelet/kubepods/%spod%s/%s",
		podType[QosFromResources(c.Resources)],
		c.PID,
		strings.ReplaceAll(c.CID, "containerd://", ""),
	)
}

func sliceNameDockerContainerdWithSystemd(c ContainerInfo, r ContainerRuntime) string {
	sliceType := [3]string{"", "kubepods-besteffort.slice/", "kubepods-burstable.slice/"}
	podType := [3]string{"", "-besteffort", "-burstable"}
	runtimeTypePrefix := [2]string{"docker", "cri-containerd"}
	runtimeURLPrefix := [2]string{"docker://", "containerd://"}
	return fmt.Sprintf(
		"/kubepods.slice/%skubepods%s-pod%s.slice/%s-%s.scope",
		sliceType[QosFromResources(c.Resources)],
		podType[QosFromResources(c.Resources)],
		strings.ReplaceAll(c.PID, "-", "_"),
		runtimeTypePrefix[r],
		strings.ReplaceAll(c.CID, runtimeURLPrefix[r], ""),
	)
}

func sliceNameDockerContainerdWithCgroupfs(c ContainerInfo, r ContainerRuntime) string {
	sliceType := [3]string{"", "besteffort/", "burstable/"}
	runtimeURLPrefix := [2]string{"docker://", "containerd://"}
	return fmt.Sprintf(
		"/kubepods/%spod%s/%s",
		sliceType[QosFromResources(c.Resources)],
		c.PID,
		strings.ReplaceAll(c.CID, runtimeURLPrefix[r], ""),
	)
}
