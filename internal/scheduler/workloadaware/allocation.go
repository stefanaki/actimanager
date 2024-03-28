package workloadaware

import (
	"cslab.ece.ntua.gr/actimanager/api/config"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/podcpubinding"
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
	"strconv"
)

type AllocatableCPU struct {
	CoreID          string
	SocketID        string
	NumCPUsInCore   int
	NumCPUsInSocket int
	Shared          bool
}

type NodeAllocatableCPUs map[int]AllocatableCPU

func allocatableCPUsForNode(nodeName string, t *v1alpha1.CPUTopology, b *v1alpha1.PodCPUBindingList) NodeAllocatableCPUs {
	res := make(map[int]AllocatableCPU)
	for socketID, socket := range t.Sockets {
		for coreID, core := range socket.Cores {
			for _, cpuID := range core.CPUs {
				res[cpuID] = AllocatableCPU{
					CoreID:          coreID,
					SocketID:        socketID,
					NumCPUsInCore:   len(core.CPUs),
					NumCPUsInSocket: len(socket.CPUs),
					Shared:          false,
				}
			}
		}
	}

	for _, binding := range b.Items {
		if binding.Status.NodeName != nodeName {
			continue
		}
		if binding.Status.ResourceStatus != "" &&
			binding.Status.ResourceStatus != v1alpha1.StatusApplied &&
			binding.Status.ResourceStatus != v1alpha1.StatusValidated &&
			binding.Status.ResourceStatus != v1alpha1.StatusBindingPending {
			continue
		}

		if binding.Spec.ExclusivenessLevel == v1alpha1.ResourceLevelNone {
			cpus := pcbutils.CPUsOfCPUBinding(&binding)
			for cpu := range cpus {
				x := res[cpu]
				x.Shared = true
				res[cpu] = x
			}
		} else {
			cpus := pcbutils.ExclusiveCPUsOfCPUBinding(&binding, t)
			for cpu := range cpus {
				delete(res, cpu)
			}
		}
	}

	return res
}

func allocatableSockets(allocatable NodeAllocatableCPUs) []int {
	sockets := make(map[int]struct{})
	for _, cpu := range allocatable {
		socketID, _ := strconv.Atoi(cpu.SocketID)
		sockets[socketID] = struct{}{}
	}
	return maps.Keys(sockets)
}

func allocatableCores(allocatable NodeAllocatableCPUs, full bool) []int {
	cores := make(map[int]int)
	coreNumCPUs := make(map[int]int)
	for _, cpu := range allocatable {
		coreID, _ := strconv.Atoi(cpu.CoreID)
		cores[coreID] = cores[coreID] + 1
		if _, ok := coreNumCPUs[coreID]; !ok {
			coreNumCPUs[coreID] = cpu.NumCPUsInCore
		}
	}
	if full {
		for core, count := range cores {
			if count < coreNumCPUs[core] {
				delete(cores, core)
			}
		}
	}
	return maps.Keys(cores)
}

// Spaghetti alert!
func cpuSetForWorkloadType(state *State, nodeName string, workloadType string, fullCores bool) []v1alpha1.CPU {
	cpuSet := make(map[int]struct{})
	allocatable := state.AllocatableCPUs[nodeName]

	cpuRequests := state.PodRequests.Cpu().MilliValue()

	for int64(len(cpuSet)*1000) < cpuRequests {
		for _, socket := range state.Topologies[nodeName].Sockets {
			for _, core := range socket.Cores {
				coreAllocatable := true
				// If workload is CPUBound or the PhysicalCores feature is enabled, check if the core is fully allocatable
				if workloadType == config.WorkloadTypeCPUBound || fullCores {
					for _, cpu := range core.CPUs {
						if _, ok := allocatable[cpu]; !ok {
							coreAllocatable = false
							break
						}
					}
				}
				if coreAllocatable {
					for _, cpu := range core.CPUs {
						if _, ok := allocatable[cpu]; !ok {
							fmt.Printf("CPU %d is not allocatable\n", cpu)
							continue
						}
						if _, ok := cpuSet[cpu]; !ok {
							cpuSet[cpu] = struct{}{}
							if workloadType != config.WorkloadTypeIOBound {
								break
							}
							if int64(len(cpuSet)*1000) >= cpuRequests {
								break
							}
							if workloadType == config.WorkloadTypeBestEffort && len(cpuSet) >= 1 {
								break
							}
						}
					}
				}
				if int64(len(cpuSet)*1000) >= cpuRequests {
					break
				}
				if workloadType == config.WorkloadTypeBestEffort && len(cpuSet) >= 1 {
					break
				}
			}
			if int64(len(cpuSet)*1000) >= cpuRequests {
				break
			}
			if workloadType == config.WorkloadTypeBestEffort && len(cpuSet) >= 1 {
				break
			}
		}
		if workloadType == config.WorkloadTypeBestEffort && len(cpuSet) >= 1 {
			break
		}
	}

	res := maps.Keys(cpuSet)
	slices.Sort(res)
	return pcbutils.IntSliceToCPUSlice(res)
}
