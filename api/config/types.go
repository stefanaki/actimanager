package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const LabelWorkloadType = "cslab.ece.ntua.gr/workload-type"

const WorkloadTypeMemoryBound = "memory-bound"
const WorkloadTypeCPUBound = "cpu-bound"
const WorkloadTypeIOBound = "io-bound"
const WorkloadTypeBestEffort = "best-effort"

type WorkloadAwarePolicy string

const (
	PolicyMaximumUtilization WorkloadAwarePolicy = "MaximumUtilization"
	PolicyBalanced           WorkloadAwarePolicy = "Balanced"
)

type Feature string

var (
	FeatureMemoryBoundExclusiveSockets Feature = "MemoryBoundExclusiveSockets"
	FeaturePhysicalCores               Feature = "PhysicalCores"
	FeatureBestEffortSharedCPUs        Feature = "BestEffortSharedCPUs"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkloadAwareArgs holds arguments used to configure WorkloadAware plugin.
type WorkloadAwareArgs struct {
	metav1.TypeMeta

	Policy   WorkloadAwarePolicy
	Features []Feature
}
