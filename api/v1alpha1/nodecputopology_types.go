package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeCpuTopologySpec defines the desired state of NodeCpuTopology
type NodeCpuTopologySpec struct {
	NodeName string      `json:"nodeName"`
	Topology CpuTopology `json:"topology"`
}

// NodeCpuTopologyStatus defines the observed state of NodeCpuTopology
type NodeCpuTopologyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// NodeCpuTopology is the Schema for the nodecputopologies API
type NodeCpuTopology struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeCpuTopologySpec   `json:"spec,omitempty"`
	Status NodeCpuTopologyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeCpuTopologyList contains a list of NodeCpuTopology
type NodeCpuTopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeCpuTopology `json:"items"`
}

// CpuTopology represents the hierarchical topology of the CPU of a Kubernetes node
type CpuTopology struct {
	NumaNodes []NumaNode `json:"numaNodes"`
}

// NumaNode is a NUMA node of the Kubernetes node
type NumaNode struct {
	Id      int      `json:"id"`
	Sockets []Socket `json:"sockets"`
}

// Socket is a CPU socket of the parent NUMA node
type Socket struct {
	Id    int    `json:"id"`
	Cores []Core `json:"cores"`
}

// Core is a physical CPU core of the parent socket
type Core struct {
	Id      int `json:"id"`
	Threads int `json:"threads"`
}

func init() {
	SchemeBuilder.Register(&NodeCpuTopology{}, &NodeCpuTopologyList{})
}
