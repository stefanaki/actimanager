package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeCPUTopologyResourceStatus string

const (
	StatusNeedsSync      NodeCPUTopologyResourceStatus = "NeedsSync"
	StatusNodeNotFound   NodeCPUTopologyResourceStatus = "NodeNotFound"
	StatusFresh          NodeCPUTopologyResourceStatus = "Fresh"
	StatusTopologyFailed NodeCPUTopologyResourceStatus = "Failed"
)

// NodeCPUTopologySpec defines the desired state of NodeCPUTopology
type NodeCPUTopologySpec struct {
	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Optional
	Topology CPUTopology `json:"topology"`
}

// NodeCPUTopologyStatus defines the observed state of NodeCPUTopology
type NodeCPUTopologyStatus struct {
	//+kubebuilder:validation:Required
	ResourceStatus NodeCPUTopologyResourceStatus `json:"resourceStatus"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=nct
// +kubebuilder:printcolumn:name="Node",type=string,JSONPath=`.spec.nodeName`
// +kubebuilder:printcolumn:name="Resource Status",type=string,JSONPath=`.status.resourceStatus`

// NodeCPUTopology is the Schema for the nodecputopologies API
type NodeCPUTopology struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeCPUTopologySpec   `json:"spec,omitempty"`
	Status NodeCPUTopologyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeCPUTopologyList contains a list of NodeCPUTopology
type NodeCPUTopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeCPUTopology `json:"items"`
}

// CPUTopology represents the hierarchical topology of the CPU of a Kubernetes node
type CPUTopology struct {
	//+kubebuilder:validation:Optional
	Sockets map[string]Socket `json:"sockets"`
	//+kubebuilder:validation:Optional
	NUMANodes map[string]NUMANode `json:"numaNodes"`
	//+kubebuilder:validation:Optional
	CPUs []int `json:"cpus"`
}

// Socket is a CPU socket of the Kubernetes node
type Socket struct {
	Cores map[string]Core `json:"cores"`
	CPUs  []int           `json:"cpus"`
}

// NUMANode is a NUMA node of the Kubernetes node
type NUMANode struct {
	CPUs []int `json:"cpus"`
}

// Core is a physical CPU core of the parent socket
type Core struct {
	CPUs []int `json:"cpus"`
}

// CPU is a logical CPU core of the parent core
type CPU struct {
	CPUID int `json:"cpuID"`
}

func init() {
	SchemeBuilder.Register(&NodeCPUTopology{}, &NodeCPUTopologyList{})
}
