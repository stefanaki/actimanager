package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StatusBindingPending       string = "Pending"
	StatusInvalidCpuSet        string = "InvalidCpuSet"
	StatusPodNotFound          string = "PodNotFound"
	StatusNodeTopologyNotFound string = "NodeTopologyNotFound"
	StatusApplied              string = "Applied"
	StatusFailed               string = "Failed"
)

const (
	ActionUpdateAnnotationKey string = "action-update"
	ActionDeleteAnnotationKey string = "action-delete"
)

// PodCpuBindingSpec defines the CPU set on which a pod is bound,
// as well as the level of exclusiveness of the resources it needs
type PodCpuBindingSpec struct {
	// +kubebuilder:validation:Required
	PodName string `json:"podName"`

	// +kubebuilder:validation:Required
	CpuSet []Cpu `json:"cpuSet"`

	// +kubebuilder:validation:Enum=None;Cpu;Core;Socket;Numa
	// +kubebuilder:default:Cpu
	ExclusivenessLevel string `json:"exclusivenessLevel"`
}

// PodCpuBindingStatus defines the observed state of PodCpuBinding
type PodCpuBindingStatus struct {
	// +kubebuilder:validation:Enum=Applied;Pending;PodNotFound;InvalidCpuSet;Collision;Failed
	Status   string            `json:"status"`
	LastSpec PodCpuBindingSpec `json:"lastSpec"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=pcb
// +kubebuilder:printcolumn:name="Pod Name",type=string,JSONPath=`.spec.podName`
// +kubebuilder:printcolumn:name="Exclusiveness Level",type=string,JSONPath=`.spec.exclusivenessLevel`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`

// PodCpuBinding is the Schema for the podcpubindings API
type PodCpuBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodCpuBindingSpec   `json:"spec"`
	Status PodCpuBindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodCpuBindingList contains a list of PodCpuBinding
type PodCpuBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodCpuBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodCpuBinding{}, &PodCpuBindingList{})
}
