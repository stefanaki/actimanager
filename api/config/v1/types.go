package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkloadAwarePolicy string

const (
	PolicyMaximumUtilization WorkloadAwarePolicy = "PolicyMaximumUtilization"
	PolicyBalanced           WorkloadAwarePolicy = "PolicyBalanced"
)

type Feature string

const (
	FeaturePhysicalCores               Feature = "PhysicalCores"
	FeatureMemoryBoundExclusiveSockets Feature = "MemoryBoundExclusiveSockets"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkloadAwareArgs defines the scheduling parameters for WorkloadAware plugin.
type WorkloadAwareArgs struct {
	metav1.TypeMeta `json:",inline"`

	Policy   *WorkloadAwarePolicy `json:"policy,omitempty"`
	Features *[]Feature           `json:"features,omitempty"`
}
