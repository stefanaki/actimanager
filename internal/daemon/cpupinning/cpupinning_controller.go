package cpupinning

import (
	"fmt"
	"github.com/containerd/cgroups"
	cgroupsv2 "github.com/containerd/cgroups/v2"
	"github.com/go-logr/logr"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ResourceNotSet is used as default resource allocation in CgroupController.UpdateCPUSet invocations.
const ResourceNotSet = ""

// CGroupDriver stores cgroup driver used by kubelet.
type CGroupDriver int

// CGroup drivers as defined in kubelet.
const (
	DriverSystemd CGroupDriver = iota
	DriverCgroupfs
)

// ContainerInfo Represents a container in the Daemon.
type ContainerInfo struct {
	CID  string
	PID  string
	Name string
	Cpus int
	QS   QoS
}

// QoS pod and containers quality of service type.
type QoS int

// QoS classes as defined in K8s.
const (
	Guaranteed QoS = iota
	BestEffort
	Burstable
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

type CpuPinningController struct {
	containerRuntime ContainerRuntime
	cgroupDriver     CGroupDriver
	cgroupPath       string
	availableCpus    string
	logger           logr.Logger
}

func NewCpuPinningController(containerRuntime ContainerRuntime,
	cgroupDriver CGroupDriver, cgroupPath string,
	logger logr.Logger) (*CpuPinningController, error) {

	var (
		cpuSetFilePath string
		cpuSetFileName string
	)

	if cgroups.Mode() != cgroups.Unified {
		cpuSetFilePath = cgroupPath + "/cpuset"
		cpuSetFileName = "cpuset.cpus"
	} else {
		cpuSetFilePath = cgroupPath
		cpuSetFileName = "cpuset.cpus.effective"
	}

	cpuSet, err := os.ReadFile(filepath.Join(cpuSetFilePath, cpuSetFileName))
	if err != nil {
		logger.Error(err, "could not get cpuset from file system")
		os.Exit(1)
	}

	c := CpuPinningController{
		containerRuntime: containerRuntime,
		cgroupDriver:     cgroupDriver,
		cgroupPath:       cgroupPath,
		availableCpus:    strings.Trim(strings.Trim(string(cpuSet), " "), "\n"),
		logger:           logger.WithName("cpu-pinning"),
	}

	return &c, nil
}

// SliceName returns path to container cgroup leaf slice in cgroupfs.
func SliceName(c ContainerInfo, r ContainerRuntime, d CGroupDriver) string {
	if r == Kind {
		return sliceNameKind(c)
	}
	if d == DriverSystemd {
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
func (c CpuPinningController) UpdateCPUSet(pPath string, container ContainerInfo, cSet string, memSet string) error {
	runtimeURLPrefix := [2]string{"docker://", "containerd://"}
	if c.containerRuntime == Kind || c.containerRuntime != Kind &&
		strings.Contains(container.CID, runtimeURLPrefix[c.containerRuntime]) {
		slice := SliceName(container, c.containerRuntime, c.cgroupDriver)
		c.logger.V(2).Info("allocating cgroup", "cgroupPath", pPath, "slicePath", slice, "cpuSet", cSet, "memSet", memSet)

		if cgroups.Mode() == cgroups.Unified {
			return c.updateCgroupsV2(pPath, slice, cSet, memSet)
		}
		return c.updateCgroupsV1(pPath, slice, cSet, memSet)
	}

	return nil
}

func (c CpuPinningController) updateCgroupsV1(pPath, slice, cSet, memSet string) error {
	ctrl := cgroups.NewCpuset(pPath)
	err := ctrl.Update(slice, &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Cpus: cSet,
			Mems: memSet,
		},
	})
	// if we set the memory pinning we should enable memory_migrate in cgroups v1
	if err == nil && memSet != "" {
		migratePath := path.Join(pPath, "cpuset", slice, "cpuset.memory_migrate")
		err = os.WriteFile(migratePath, []byte("1"), os.FileMode(0))
	}
	return err
}

func (c CpuPinningController) updateCgroupsV2(pPath, slice, cSet, memSet string) error {
	res := cgroupsv2.Resources{CPU: &cgroupsv2.CPU{
		Cpus: cSet,
		Mems: memSet,
	}}
	_, err := cgroupsv2.NewManager(pPath, slice, &res)
	// memory migration in cgroups v2 is always enabled, no need to set it as in cgroupsv1
	return err
}

func (c CpuPinningController) Apply(container *ContainerInfo, cSet string) error {
	return c.UpdateCPUSet(c.cgroupPath, *container, cSet, ResourceNotSet)
}

func (c CpuPinningController) Remove(container *ContainerInfo) error {
	return c.UpdateCPUSet(c.cgroupPath, *container, c.availableCpus, ResourceNotSet)
}
