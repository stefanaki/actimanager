package cpupinning

import (
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	"cslab.ece.ntua.gr/actimanager/internal/daemon/client"
	cgroupsctrl "cslab.ece.ntua.gr/actimanager/internal/pkg/cgroups"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/protobuf/cpupinning"
	nctutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/nodecputopology"
	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/podcpubinding"
	"fmt"
	"github.com/go-logr/logr"
	"strings"
)

type CPUPinningController struct {
	cgroupsController   cgroupsctrl.CgroupsController
	containerRuntime    ContainerRuntime
	podCPUBindingClient *client.PodCPUBindingClient
	podClient           *client.PodClient
	cpuTopology         v1alpha1.CPUTopology
	nodeName            string
	logger              logr.Logger
}

// NewCPUPinningController returns a reference to a new CPUPinningController instance
func NewCPUPinningController(containerRuntime ContainerRuntime,
	cgroupsDriver cgroupsctrl.CgroupsDriver, cgroupsPath string,
	podCPUBindingClient *client.PodCPUBindingClient,
	podClient *client.PodClient,
	cpuTopology v1alpha1.CPUTopology,
	nodeName string,
	logger logr.Logger) (*CPUPinningController, error) {

	cgroupsController, err := cgroupsctrl.NewCgroupsController(cgroupsDriver, cgroupsPath, logger)
	if err != nil {
		return nil, fmt.Errorf("could create cgroups controller: %v", err)
	}

	c := CPUPinningController{
		containerRuntime:    containerRuntime,
		cgroupsController:   cgroupsController,
		podCPUBindingClient: podCPUBindingClient,
		podClient:           podClient,
		cpuTopology:         cpuTopology,
		nodeName:            nodeName,
		logger:              logger.WithName("cpu-pinning"),
	}

	return &c, nil
}

// UpdateCPUSet updates the cpu set of a given child process.
func (c CPUPinningController) UpdateCPUSet(container ContainerInfo, cSet string, memSet string, quota *int64, shares, period *uint64) error {
	runtimeURLPrefix := [2]string{"docker://", "containerd://"}
	if c.containerRuntime == Kind || c.containerRuntime != Kind &&
		strings.Contains(container.CID, runtimeURLPrefix[c.containerRuntime]) {
		slice := SliceName(container, c.containerRuntime, c.cgroupsController.CgroupsDriver)
		// c.logger.V(2).Info("allocating cgroup", "cgroupPath", c.cgroupsController.CgroupsPath, "slicePath", slice, "cpuSet", cSet, "memSet", memSet)
		return c.cgroupsController.UpdateCPUSet(slice, cSet, memSet, quota, shares, period)
	}

	return nil
}

// Apply updates the CPU set of the container, reconciling with the CPU bindings of other pods.
func (c CPUPinningController) Apply(pod *cpupinning.Pod, cpuSet string, memSet string) error {
	if err := c.reconcilePodsWithSharedResources(pod, false); err != nil {
		return fmt.Errorf("failed to reconcile pods with shared resources: %v", err)
	}
	for _, container := range pod.Containers {
		info := ContainerInfo{
			CID:  container.Id,
			PID:  pod.Id,
			Name: container.Name,
			Resources: ResourceInfo{
				RequestedCPUs:   int64(container.Resources.RequestedCPUs),
				LimitCPUs:       int64(container.Resources.LimitCPUs),
				RequestedMemory: container.Resources.RequestedMemory,
				LimitMemory:     container.Resources.LimitMemory,
			},
		}
		quota := MilliCPUToQuota(info.Resources.LimitCPUs, QuotaPeriod)
		shares := MilliCPUToShares(info.Resources.LimitCPUs)
		period := uint64(QuotaPeriod)
		if err := c.UpdateCPUSet(info, cpuSet, memSet, &quota, &shares, &period); err != nil {
			return fmt.Errorf("failed to update cpu set for container %s in pod %s/%s: %v", container.Name, pod.Namespace, pod.Name, err)
		}
		c.logger.Info("CPUSet updated", "pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "container", container.Name, "cpuSet", cpuSet, "memSet", memSet, "quota", quota, "shares", shares, "period", period)
	}
	return nil
}

// Remove updates the CPU set of the container to all the available CPUs.
func (c CPUPinningController) Remove(pod *cpupinning.Pod) error {
	err := c.reconcilePodsWithSharedResources(pod, true)
	if err != nil {
		return fmt.Errorf("failed to reconcile pods with shared resources: %v", err)
	}
	return nil
}

func (c CPUPinningController) reconcilePodsWithSharedResources(pod *cpupinning.Pod, rm bool) error {
	sharedCPUs := c.cpuTopology.DeepCopy().CPUs
	cpuBindings, err := c.podCPUBindingClient.PodCPUBindingsForNode(c.nodeName)
	if err != nil {
		return fmt.Errorf("failed to get pod cpu bindings: %v", err)
	}
	cpuBoundPods := make(map[string]struct{})
	for _, binding := range cpuBindings {
		cpuBoundPods[fmt.Sprintf("%s/%s", binding.Namespace, binding.Spec.PodName)] = struct{}{}
		if binding.Status.ResourceStatus != v1alpha1.StatusApplied &&
			binding.Status.ResourceStatus != v1alpha1.StatusBindingPending &&
			binding.Status.ResourceStatus != v1alpha1.StatusValidated {
			continue
		}
		if rm && pod.Namespace == binding.Namespace && pod.Name == binding.Spec.PodName {
			continue
		}
		for cpu := range pcbutils.GetExclusiveCPUsOfCPUBinding(&binding, &c.cpuTopology) {
			for i, sharedCPU := range sharedCPUs {
				if sharedCPU == cpu {
					sharedCPUs = append(sharedCPUs[:i], sharedCPUs[i+1:]...)
					break
				}
			}
		}
	}
	c.logger.Info("--- shared cpus ---", "cpus", sharedCPUs)
	pods, err := c.podClient.PodsForNode(c.nodeName)
	if err != nil {
		return fmt.Errorf("failed to get pods for node: %v", err)
	}
	for _, pod := range pods {
		if _, ok := cpuBoundPods[fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)]; !ok {
			// Pod is not bound to any CPUs, so we can use all the calculated shared CPUs
			for _, container := range pod.Status.ContainerStatuses {
				if !container.Ready {
					continue
				}
				cpus := ConvertIntSliceToString(sharedCPUs)
				mems := ConvertIntSliceToString(nctutils.GetNUMANodesOfCPUSet(sharedCPUs, c.cpuTopology))
				resources := cpupinning.ParseContainerResources(container.Name, &pod)
				info := ContainerInfo{
					CID:  container.ContainerID,
					PID:  string(pod.ObjectMeta.UID),
					Name: container.Name,
					Resources: ResourceInfo{
						RequestedCPUs:   int64(resources.RequestedCPUs),
						LimitCPUs:       int64(resources.LimitCPUs),
						RequestedMemory: resources.RequestedMemory,
						LimitMemory:     resources.LimitMemory,
					},
				}
				quota := MilliCPUToQuota(info.Resources.LimitCPUs, QuotaPeriod)
				shares := MilliCPUToShares(info.Resources.LimitCPUs)
				period := uint64(QuotaPeriod)
				err := c.UpdateCPUSet(info, cpus, mems, &quota, &shares, &period)
				c.logger.Info("CPUSet updated", "pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "container", container.Name, "cpuSet", cpus, "memSet", mems, "quota", quota, "shares", shares, "period", period)
				if err != nil {
					return fmt.Errorf("failed to update cpu set: %v", err)
				}
			}
		}
	}
	return nil
}
