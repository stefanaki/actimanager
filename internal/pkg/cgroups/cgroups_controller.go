package cgroups

import (
	"os"
	"path"

	"github.com/containerd/cgroups"
	cgroupsv2 "github.com/containerd/cgroups/v2"
	"github.com/go-logr/logr"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// CgroupsDriver represents the cgroups driver used by the host.
type CgroupsDriver int

// Supported cgroups drivers.
const (
	DriverSystemd CgroupsDriver = iota
	DriverCgroupfs
)

// ResourceNotSet is used as default resource allocation in CgroupsController.
const ResourceNotSet = ""

// CgroupsController represents a controller for managing cgroups.
type CgroupsController struct {
	CgroupsDriver CgroupsDriver
	CgroupsPath   string
	Logger        logr.Logger
}

// NewCgroupsController creates a new instance of CgroupsController.
func NewCgroupsController(cgroupsDriver CgroupsDriver, cgroupsPath string, logger logr.Logger) (CgroupsController, error) {
	return CgroupsController{
		CgroupsDriver: cgroupsDriver,
		CgroupsPath:   cgroupsPath,
		Logger:        logger.WithName("cgroups-controller"),
	}, nil
}

// UpdateCpuSet updates the resources of a slice.
func (c *CgroupsController) UpdateCpuSet(slice, cpuSet, memSet string) error {
	if cgroups.Mode() == cgroups.Unified {
		return c.updateCpuSetV2(slice, cpuSet, memSet)
	}
	return c.updateCpuSetV1(slice, cpuSet, memSet)
}

// updateCpuSetV1 updates cgroups for v1 mode.
func (c *CgroupsController) updateCpuSetV1(slice, cpuSet, memSet string) error {
	ctrl := cgroups.NewCpuset(c.CgroupsPath)
	err := ctrl.Update(slice, &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Cpus: cpuSet,
			Mems: memSet,
		},
	})

	// Enable memory migration in cgroups v1 if memory set is specified.
	if err == nil && memSet != "" {
		migratePath := path.Join(c.CgroupsPath, "cpuset", slice, "cpuset.memory_migrate")
		err = os.WriteFile(migratePath, []byte("1"), os.ModePerm)
	}

	return err
}

// updateCpuSetV2 updates cgroups for v2 (unified) mode.
func (c *CgroupsController) updateCpuSetV2(slice, cpuSet, memSet string) error {
	res := cgroupsv2.Resources{CPU: &cgroupsv2.CPU{
		Cpus: cpuSet,
		Mems: memSet,
	}}
	_, err := cgroupsv2.NewManager(c.CgroupsPath, slice, &res)
	// Memory migration in cgroups v2 is always enabled, no need to set it.
	return err
}
