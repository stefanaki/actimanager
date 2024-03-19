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

type CpuPinningController struct {
	cgroupsController   cgroupsctrl.CgroupsController
	containerRuntime    ContainerRuntime
	podCpuBindingClient *client.PodCpuBindingClient
	podClient           *client.PodClient
	cpuTopology         v1alpha1.CpuTopology
	nodeName            string
	logger              logr.Logger
}

// NewCpuPinningController returns a reference to a new CpuPinningController instance
func NewCpuPinningController(containerRuntime ContainerRuntime,
	cgroupsDriver cgroupsctrl.CgroupsDriver, cgroupsPath string,
	podCpuBindingClient *client.PodCpuBindingClient,
	podClient *client.PodClient,
	cpuTopology v1alpha1.CpuTopology,
	nodeName string,
	logger logr.Logger) (*CpuPinningController, error) {

	cgroupsController, err := cgroupsctrl.NewCgroupsController(cgroupsDriver, cgroupsPath, logger)
	if err != nil {
		return nil, fmt.Errorf("could create cgroups controller: %v", err)
	}

	c := CpuPinningController{
		containerRuntime:    containerRuntime,
		cgroupsController:   cgroupsController,
		podCpuBindingClient: podCpuBindingClient,
		podClient:           podClient,
		cpuTopology:         cpuTopology,
		nodeName:            nodeName,
		logger:              logger.WithName("cpu-pinning"),
	}

	return &c, nil
}

// UpdateCpuSet updates the cpu set of a given child process.
func (c CpuPinningController) UpdateCpuSet(container ContainerInfo, cSet string, memSet string, quota *int64, shares, period *uint64) error {
	runtimeURLPrefix := [2]string{"docker://", "containerd://"}
	if c.containerRuntime == Kind || c.containerRuntime != Kind &&
		strings.Contains(container.CID, runtimeURLPrefix[c.containerRuntime]) {
		slice := SliceName(container, c.containerRuntime, c.cgroupsController.CgroupsDriver)
		// c.logger.V(2).Info("allocating cgroup", "cgroupPath", c.cgroupsController.CgroupsPath, "slicePath", slice, "cpuSet", cSet, "memSet", memSet)
		return c.cgroupsController.UpdateCpuSet(slice, cSet, memSet, quota, shares, period)
	}

	return nil
}

// Apply updates the CPU set of the container, reconciling with the CPU bindings of other pods.
func (c CpuPinningController) Apply(pod *cpupinning.Pod, cpuSet string, memSet string) error {
	if err := c.reconcilePodsWithSharedResources(pod, false); err != nil {
		return fmt.Errorf("failed to reconcile pods with shared resources: %v", err)
	}
	for _, container := range pod.Containers {
		info := ContainerInfo{
			CID:  container.Id,
			PID:  pod.Id,
			Name: container.Name,
			Resources: ResourceInfo{
				RequestedCpus:   int64(container.Resources.RequestedCpus),
				LimitCpus:       int64(container.Resources.LimitCpus),
				RequestedMemory: container.Resources.RequestedMemory,
				LimitMemory:     container.Resources.LimitMemory,
			},
		}
		quota := MilliCPUToQuota(info.Resources.LimitCpus, QuotaPeriod)
		shares := MilliCPUToShares(info.Resources.LimitCpus)
		period := uint64(QuotaPeriod)
		if err := c.UpdateCpuSet(info, cpuSet, memSet, &quota, &shares, &period); err != nil {
			return fmt.Errorf("failed to update cpu set for container %s in pod %s/%s: %v", container.Name, pod.Namespace, pod.Name, err)
		}
		c.logger.Info("CPUSet updated", "pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "container", container.Name, "cpuSet", cpuSet, "memSet", memSet, "quota", quota, "shares", shares, "period", period)
	}
	return nil
}

// Remove updates the CPU set of the container to all the available CPUs.
func (c CpuPinningController) Remove(pod *cpupinning.Pod) error {
	err := c.reconcilePodsWithSharedResources(pod, true)
	if err != nil {
		return fmt.Errorf("failed to reconcile pods with shared resources: %v", err)
	}
	return nil
}

func (c CpuPinningController) reconcilePodsWithSharedResources(pod *cpupinning.Pod, rm bool) error {
	sharedCpus := c.cpuTopology.DeepCopy().Cpus
	cpuBindings, err := c.podCpuBindingClient.PodCpuBindingsForNode(c.nodeName)
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
		for cpu := range pcbutils.GetExclusiveCpusOfCpuBinding(&binding, &c.cpuTopology) {
			for i, sharedCpu := range sharedCpus {
				if sharedCpu == cpu {
					sharedCpus = append(sharedCpus[:i], sharedCpus[i+1:]...)
					break
				}
			}
		}
	}
	c.logger.Info("--- shared cpus ---", "cpus", sharedCpus)
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
				cpus := ConvertIntSliceToString(sharedCpus)
				mems := ConvertIntSliceToString(nctutils.GetNumaNodesOfCpuSet(sharedCpus, c.cpuTopology))
				resources := cpupinning.ParseContainerResources(container.Name, &pod)
				info := ContainerInfo{
					CID:  container.ContainerID,
					PID:  string(pod.ObjectMeta.UID),
					Name: container.Name,
					Resources: ResourceInfo{
						RequestedCpus:   int64(resources.RequestedCpus),
						LimitCpus:       int64(resources.LimitCpus),
						RequestedMemory: resources.RequestedMemory,
						LimitMemory:     resources.LimitMemory,
					},
				}
				quota := MilliCPUToQuota(info.Resources.LimitCpus, QuotaPeriod)
				shares := MilliCPUToShares(info.Resources.LimitCpus)
				period := uint64(QuotaPeriod)
				err := c.UpdateCpuSet(info, cpus, mems, &quota, &shares, &period)
				c.logger.Info("CPUSet updated", "pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name), "container", container.Name, "cpuSet", cpus, "memSet", mems, "quota", quota, "shares", shares, "period", period)
				if err != nil {
					return fmt.Errorf("failed to update cpu set: %v", err)
				}
			}
		}
	}
	return nil
}
