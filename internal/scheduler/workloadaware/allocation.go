package workloadaware

import (
	"cslab.ece.ntua.gr/actimanager/api/config"
	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	pcbutils "cslab.ece.ntua.gr/actimanager/internal/pkg/utils/podcpubinding"
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

func allocatableCPUsForNode(
	nodeName string,
	topology *v1alpha1.CPUTopology,
	bindings *v1alpha1.PodCPUBindingList,
	workloadType string,
) NodeAllocatableCPUs {
	res := make(map[int]AllocatableCPU)
	for socketID, socket := range topology.Sockets {
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

	for _, binding := range bindings.Items {
		if binding.Status.NodeName != nodeName {
			continue
		}
		if binding.Status.ResourceStatus != v1alpha1.StatusApplied &&
			binding.Status.ResourceStatus != v1alpha1.StatusValidated {
			continue
		}
		if binding.Spec.ExclusivenessLevel == v1alpha1.ResourceLevelNone && workloadType != config.WorkloadTypeBestEffort {
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
			cpus := pcbutils.ExclusiveCPUsOfCPUBinding(&binding, topology)
			for cpu := range cpus {
				delete(res, cpu)
			}
		}
	}

	return res
}

func allocatableSockets(allocatable NodeAllocatableCPUs, full bool) []int {
	sockets := make(map[int]int)
	socketNumCPUs := make(map[int]int)
	for _, cpu := range allocatable {
		socketID, _ := strconv.Atoi(cpu.SocketID)
		sockets[socketID] = sockets[socketID] + 1
		if _, ok := socketNumCPUs[socketID]; !ok {
			socketNumCPUs[socketID] = cpu.NumCPUsInSocket
		}
	}
	if full {
		for socket, count := range sockets {
			if count < socketNumCPUs[socket] {
				delete(sockets, socket)
			}
		}
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

func cpuSetForMemoryBound(state *State, nodeName string, fullSockets bool) []v1alpha1.CPU {
	cpuSet := make(map[int]struct{})
	allocatable := state.AllocatableCPUs[nodeName]
	cpuRequests := state.PodRequests.Cpu().MilliValue()

	for !done(cpuSet, cpuRequests) {
		for _, socket := range state.Topologies[nodeName].Sockets {
			socketAllocatable := true
			if fullSockets {
				for _, cpu := range socket.CPUs {
					if _, ok := allocatable[cpu]; !ok {
						socketAllocatable = false
						break
					}
				}
			}
			if !socketAllocatable {
				continue
			}

			socketAllocated := false
			for _, core := range socket.Cores {
				for _, cpu := range core.CPUs {
					if _, ok := allocatable[cpu]; !ok {
						continue
					}
					if _, ok := cpuSet[cpu]; ok {
						continue
					}
					cpuSet[cpu] = struct{}{}
					socketAllocated = true
					break
				}
				if socketAllocated || done(cpuSet, cpuRequests) {
					break
				}
			}
			if done(cpuSet, cpuRequests) {
				break
			}
		}
	}

	res := maps.Keys(cpuSet)
	slices.Sort(res)
	return pcbutils.IntSliceToCPUSlice(res)
}

func cpuSetForCPUBound(state *State, nodeName string, fullCores bool) []v1alpha1.CPU {
	cpuSet := make(map[int]struct{})
	allocatable := state.AllocatableCPUs[nodeName]
	cpuRequests := state.PodRequests.Cpu().MilliValue()

	for !done(cpuSet, cpuRequests) {
		for _, socket := range state.Topologies[nodeName].Sockets {
			for _, core := range socket.Cores {
				coreAllocatable := true
				if fullCores {
					for _, cpu := range core.CPUs {
						if _, ok := allocatable[cpu]; !ok {
							coreAllocatable = false
							break
						}
					}
				}
				if !coreAllocatable {
					continue
				}
				for _, cpu := range core.CPUs {
					if _, ok := allocatable[cpu]; !ok {
						continue
					}
					if _, ok := cpuSet[cpu]; ok {
						continue
					}
					cpuSet[cpu] = struct{}{}
					break
				}
				if done(cpuSet, cpuRequests) {
					break
				}
			}
			if done(cpuSet, cpuRequests) {
				break
			}
		}
	}

	res := maps.Keys(cpuSet)
	slices.Sort(res)
	return pcbutils.IntSliceToCPUSlice(res)
}

func cpuSetForIOBound(state *State, nodeName string, fullCores bool) []v1alpha1.CPU {
	cpuSet := make(map[int]struct{})
	allocatable := state.AllocatableCPUs[nodeName]
	cpuRequests := state.PodRequests.Cpu().MilliValue()

	for !done(cpuSet, cpuRequests) {
		for _, socket := range state.Topologies[nodeName].Sockets {
			for _, core := range socket.Cores {
				coreAllocatable := true
				if fullCores {
					for _, cpu := range core.CPUs {
						if _, ok := allocatable[cpu]; !ok {
							coreAllocatable = false
							break
						}
					}
				}
				if !coreAllocatable {
					continue
				}
				for _, cpu := range core.CPUs {
					if _, ok := allocatable[cpu]; !ok {
						continue
					}
					if _, ok := cpuSet[cpu]; ok {
						continue
					}
					cpuSet[cpu] = struct{}{}
					if done(cpuSet, cpuRequests) {
						break
					}
				}
				if done(cpuSet, cpuRequests) {
					break
				}
			}
			if done(cpuSet, cpuRequests) {
				break
			}
		}
	}

	res := maps.Keys(cpuSet)
	slices.Sort(res)
	return pcbutils.IntSliceToCPUSlice(res)
}

func cpuSetForBestEffort(state *State, nodeName string, fullCores bool) []v1alpha1.CPU {
	cpuSet := make(map[int]struct{})
	allocatable := state.AllocatableCPUs[nodeName]
	cpuRequests := int64(1000)

	for !done(cpuSet, cpuRequests) {
		for _, socket := range state.Topologies[nodeName].Sockets {
			for _, core := range socket.Cores {
				coreAllocatable := true
				if fullCores {
					for _, cpu := range core.CPUs {
						if _, ok := allocatable[cpu]; !ok {
							coreAllocatable = false
							break
						}
					}
				}
				if !coreAllocatable {
					continue
				}

				var nonSharedCPUS []int
				var sharedCPUs []int

				for _, cpu := range core.CPUs {
					if _, ok := allocatable[cpu]; !ok {
						continue
					}
					if _, ok := cpuSet[cpu]; ok {
						continue
					}
					// Split shared and non-shared CPUs
					if allocatable[cpu].Shared {
						sharedCPUs = append(sharedCPUs, cpu)
					} else {
						nonSharedCPUS = append(nonSharedCPUS, cpu)
					}
				}
				// Allocate non-shared CPUs first
				for _, cpu := range nonSharedCPUS {
					cpuSet[cpu] = struct{}{}
					if done(cpuSet, cpuRequests) {
						break
					}
				}
				if done(cpuSet, cpuRequests) {
					break
				}
				// Allocate shared CPUs next if needed
				for _, cpu := range sharedCPUs {
					cpuSet[cpu] = struct{}{}
					if done(cpuSet, cpuRequests) {
						break
					}
				}

				if done(cpuSet, cpuRequests) {
					break
				}
			}
			if done(cpuSet, cpuRequests) {
				break
			}
		}
	}

	res := maps.Keys(cpuSet)
	slices.Sort(res)
	return pcbutils.IntSliceToCPUSlice(res)

}

func done(cpuSet map[int]struct{}, requests int64) bool {
	return int64(len(cpuSet)*1000) >= requests
}
