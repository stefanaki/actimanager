package v1alpha1

import cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/v1alpha1"

func IsCpuSetInTopology(topology *cslabecentuagrv1alpha1.CpuTopology, cpuSet []cslabecentuagrv1alpha1.Cpu) bool {
	remaining := len(cpuSet)

	for _, socket := range topology.Sockets {
		for _, core := range socket.Cores {
			for _, cpu := range core.Cpus {
				for _, inputCpu := range cpuSet {
					if cpu.CpuId == inputCpu.CpuId {
						remaining--
						break
					}
				}
			}
		}
	}

	return remaining == 0
}
