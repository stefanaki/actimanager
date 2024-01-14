package cpupinning

import (
	cgroupsctrl "cslab.ece.ntua.gr/actimanager/internal/pkg/cgroups"
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/go-logr/logr"
	"os"
	"path/filepath"
	"strings"
)

type CpuPinningController struct {
	cgroupsController  cgroupsctrl.CgroupsController
	containerRuntime   ContainerRuntime
	availableCpus      string
	availableNumaNodes string
	logger             logr.Logger
}

// NewCpuPinningController returns a reference to a new CpuPinningController instance
func NewCpuPinningController(containerRuntime ContainerRuntime,
	cgroupsDriver cgroupsctrl.CgroupsDriver, cgroupsPath string,
	logger logr.Logger) (*CpuPinningController, error) {

	var (
		cpuSetFilePath string
		cpuSetFileName string
		memSetFileName string
	)

	if cgroups.Mode() != cgroups.Unified {
		cpuSetFilePath = cgroupsPath + "/cpuset"
		cpuSetFileName = "cpuset.cpus"
		memSetFileName = "cpuset.mems"
	} else {
		cpuSetFilePath = cgroupsPath
		cpuSetFileName = "cpuset.cpus.effective"
		memSetFileName = "cpuset.mems.effective"
	}

	cpuSet, err := os.ReadFile(filepath.Join(cpuSetFilePath, cpuSetFileName))
	if err != nil {
		return nil, fmt.Errorf("could not get cpuset from file system: %v", err)
	}

	memSet, err := os.ReadFile(filepath.Join(cpuSetFilePath, memSetFileName))
	if err != nil {
		return nil, fmt.Errorf("could not get memset from file system: %v", err)
	}

	cgroupsController, err := cgroupsctrl.NewCgroupsController(cgroupsDriver, cgroupsPath, logger)
	if err != nil {
		return nil, fmt.Errorf("could create cgroups controller: %v", err)
	}

	c := CpuPinningController{
		containerRuntime:   containerRuntime,
		cgroupsController:  cgroupsController,
		availableCpus:      strings.Trim(strings.Trim(string(cpuSet), " "), "\n"),
		availableNumaNodes: strings.Trim(strings.Trim(string(memSet), " "), "\n"),
		logger:             logger.WithName("cpu-pinning"),
	}

	return &c, nil
}

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
		podType[c.QS],
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
		sliceType[c.QS],
		podType[c.QS],
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
		sliceType[c.QS],
		c.PID,
		strings.ReplaceAll(c.CID, runtimeURLPrefix[r], ""),
	)
}

// UpdateCPUSet updates the cpu set of a given child process.
func (c CpuPinningController) UpdateCPUSet(container ContainerInfo, cSet string, memSet string) error {
	runtimeURLPrefix := [2]string{"docker://", "containerd://"}
	if c.containerRuntime == Kind || c.containerRuntime != Kind &&
		strings.Contains(container.CID, runtimeURLPrefix[c.containerRuntime]) {
		slice := SliceName(container, c.containerRuntime, c.cgroupsController.CgroupsDriver)
		c.logger.V(2).Info("allocating cgroup", "cgroupPath", c.cgroupsController.CgroupsPath, "slicePath", slice, "cpuSet", cSet, "memSet", memSet)

		return c.cgroupsController.UpdateCpuSet(slice, cSet, memSet)
	}

	return nil
}

// Apply updates the CPU set of the container.
func (c CpuPinningController) Apply(container ContainerInfo, cpuSet string, memSet string) error {
	return c.UpdateCPUSet(container, cpuSet, memSet)
}

// Remove updates the CPU set of the container to all the available CPUs.
func (c CpuPinningController) Remove(container ContainerInfo) error {
	return c.UpdateCPUSet(container, c.availableCpus, c.availableNumaNodes)
}
