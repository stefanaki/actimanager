package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodCPUBindingResourceStatus string

const (
	StatusBindingPending         PodCPUBindingResourceStatus = "Pending"
	StatusInvalidCPUSet          PodCPUBindingResourceStatus = "InvalidCPUSet"
	StatusPodNotFound            PodCPUBindingResourceStatus = "PodNotFound"
	StatusNodeTopologyNotFound   PodCPUBindingResourceStatus = "NodeTopologyNotFound"
	StatusApplied                PodCPUBindingResourceStatus = "Applied"
	StatusFailed                 PodCPUBindingResourceStatus = "Failed"
	StatusCPUSetAllocationFailed PodCPUBindingResourceStatus = "CPUSetAllocationFailed"
	StatusValidated              PodCPUBindingResourceStatus = "Validated"
)

var FinalizerPodCPUBinding = GroupVersion.Group + "/pod-cpu-binding-finalizer"
var FinalizerCPUBoundPod = GroupVersion.Group + "/cpu-bound-pod"

var AnnotationExclusivenessLevel = GroupVersion.Group + "/exclusiveness-level"

type ResourceLevel string

const (
	ResourceLevelNone   ResourceLevel = "None"
	ResourceLevelCPU    ResourceLevel = "CPU"
	ResourceLevelCore   ResourceLevel = "Core"
	ResourceLevelSocket ResourceLevel = "Socket"
	ResourceLevelNUMA   ResourceLevel = "NUMA"
)

// PodCPUBindingSpec defines the CPU set on which a pod is bound,
// as well as the level of exclusiveness of the resources it needs
type PodCPUBindingSpec struct {
	// +kubebuilder:validation:Required
	PodName string `json:"podName"`

	// +kubebuilder:validation:Required
	CPUSet []CPU `json:"cpuSet"`

	// +kubebuilder:validation:Enum=None;CPU;Core;Socket;NUMA
	ExclusivenessLevel ResourceLevel `json:"exclusivenessLevel"`
}

// PodCPUBindingStatus defines the observed state of PodCPUBinding
type PodCPUBindingStatus struct {
	// +kubebuilder:validation:Enum=Applied;Pending;PodNotFound;InvalidCPUSet;Collision;Failed;CPUSetAllocationFailed;Validated
	ResourceStatus PodCPUBindingResourceStatus `json:"resourceStatus"`
	NodeName       string                      `json:"nodeName"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=pcb
// +kubebuilder:printcolumn:name="Pod Name",type=string,JSONPath=`.spec.podName`
// +kubebuilder:printcolumn:name="Exclusiveness Level",type=string,JSONPath=`.spec.exclusivenessLevel`
// +kubebuilder:printcolumn:name="Resource Status",type=string,JSONPath=`.status.resourceStatus`

// PodCPUBinding is the Schema for the podcpubindings API
type PodCPUBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodCPUBindingSpec   `json:"spec"`
	Status PodCPUBindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodCPUBindingList contains a list of PodCPUBinding
type PodCPUBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodCPUBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodCPUBinding{}, &PodCPUBindingList{})
}
