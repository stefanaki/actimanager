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

// allocatableCPUsForNode computes the allocatable CPUs for a given node based on its
// topology and current CPU bindings.
func allocatableCPUsForNode(
	nodeName string,
	topology *v1alpha1.CPUTopology,
	bindings []*v1alpha1.PodCPUBinding,
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

	for _, binding := range bindings {
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
			cpus := pcbutils.CPUsOfCPUBinding(binding)
			for cpu := range cpus {
				x := res[cpu]
				x.Shared = true
				res[cpu] = x
			}
		} else {
			cpus := pcbutils.ExclusiveCPUsOfCPUBinding(binding, topology)
			for cpu := range cpus {
				delete(res, cpu)
			}
		}
	}

	return res
}

// allocatableSockets identifies the sockets with fully allocatable CPUs based on
// the given NodeAllocatableCPUs map and a flag indicating whether to filter fully allocatable sockets.
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

// allocatableCores identifies the cores with fully allocatable CPUs based on
// the given NodeAllocatableCPUs map and a flag indicating whether to filter fully allocatable cores.
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

// allCPUsAllocatableInCore checks if all CPUs in a given core are allocatable
// based on the provided NodeAllocatableCPUs map.
func allCPUsAllocatableInCore(core *v1alpha1.Core, allocatable NodeAllocatableCPUs) bool {
	for _, cpu := range core.CPUs {
		if _, ok := allocatable[cpu]; !ok {
			return false
		}
	}
	return true
}

// allCPUsAllocatableInSocket checks if all CPUs in a given socket are allocatable
// based on the provided NodeAllocatableCPUs map.
func allCPUsAllocatableInSocket(socket *v1alpha1.Socket, allocatable NodeAllocatableCPUs) bool {
	for _, cpu := range socket.CPUs {
		if _, ok := allocatable[cpu]; !ok {
			return false
		}
	}
	return true
}

// cpuFeasible checks if a specific CPU is feasible for allocation based on
// its availability in the allocatable CPUs map and absence in the current cpuSet.
func cpuFeasible(cpu int, allocatable NodeAllocatableCPUs, cpuSet map[int]struct{}) bool {
	_, allocatableCPU := allocatable[cpu]
	_, cpuAlreadyUsed := cpuSet[cpu]
	return allocatableCPU && !cpuAlreadyUsed
}

// done checks if the CPU allocation process is complete based on the number of
// allocated CPUs and the requested CPU resources.
func done(cpuSet map[int]struct{}, requests int64) bool {
	return int64(len(cpuSet)*1000) >= requests
}

func cpuSetForMemoryBound(state *State, nodeName string, fullCores, fullSockets bool) []v1alpha1.CPU {
	cpuSet := make(map[int]struct{})
	allocatable := state.AllocatableCPUs[nodeName]
	cpuRequests := state.PodRequests.Cpu().MilliValue()
	seenCores := make(map[string]struct{})

	// First try to allocate one thread per socket
	for _, socket := range state.Topologies[nodeName].Sockets {
		socketAllocated := false
		// Check if whole socket is allocatable
		socketAllocatable := true
		if fullSockets {
			socketAllocatable = allCPUsAllocatableInSocket(&socket, allocatable)
		}
		if !socketAllocatable {
			continue
		}
		for coreID, core := range socket.Cores {
			// Check if whole physical core is allocatable
			coreAllocatable := true
			if fullCores {
				coreAllocatable = allCPUsAllocatableInCore(&core, allocatable)
			}
			if !coreAllocatable {
				continue
			}
			// Pick thread from physical core
			for _, cpu := range core.CPUs {
				if !cpuFeasible(cpu, allocatable, cpuSet) {
					continue
				}
				cpuSet[cpu] = struct{}{}
				seenCores[coreID] = struct{}{}
				socketAllocated = true
				break
			}
			if socketAllocated {
				break
			}
		}
		if done(cpuSet, cpuRequests) {
			break
		}
	}
	// Pick threads from already utilized sockets
	for !done(cpuSet, cpuRequests) {
		for _, socket := range state.Topologies[nodeName].Sockets {
			for coreID, core := range socket.Cores {
				// Check if whole physical core is allocatable
				coreAllocatable := true
				if fullCores {
					coreAllocatable = allCPUsAllocatableInCore(&core, allocatable)
				}
				if !coreAllocatable {
					continue
				}
				// Pick thread from physical core
				coreAllocated := false
				for _, cpu := range core.CPUs {
					alloc, ok := allocatable[cpu]
					if !ok {
						continue
					}
					// Skip seen cores
					if _, ok := seenCores[alloc.CoreID]; ok {
						continue
					}
					if _, ok := cpuSet[cpu]; ok {
						continue
					}
					seenCores[coreID] = struct{}{}
					cpuSet[cpu] = struct{}{}
					coreAllocated = true
					break
				}
				if coreAllocated || done(cpuSet, cpuRequests) {
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
					coreAllocatable = allCPUsAllocatableInCore(&core, allocatable)
				}
				if !coreAllocatable {
					continue
				}
				for _, cpu := range core.CPUs {
					if !cpuFeasible(cpu, allocatable, cpuSet) {
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
					coreAllocatable = allCPUsAllocatableInCore(&core, allocatable)
				}
				if !coreAllocatable {
					continue
				}
				for _, cpu := range core.CPUs {
					if !cpuFeasible(cpu, allocatable, cpuSet) {
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
					coreAllocatable = allCPUsAllocatableInCore(&core, allocatable)
				}
				if !coreAllocatable {
					continue
				}

				var nonSharedCPUS []int
				var sharedCPUs []int

				for _, cpu := range core.CPUs {
					if !cpuFeasible(cpu, allocatable, cpuSet) {
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
