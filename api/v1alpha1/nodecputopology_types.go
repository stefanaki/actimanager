package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeCpuTopologySpec defines the desired state of NodeCpuTopology
type NodeCpuTopologySpec struct {
	NodeName string `json:"nodeName"`

	//+kubebuilder:validation:Optional
	Topology CpuTopology `json:"topology"`
}

// NodeCpuTopologyStatus defines the observed state of NodeCpuTopology
type NodeCpuTopologyStatus struct {
	//+kubebuilder:validation:Required
	InitJobStatus string `json:"initJobStatus"`
	InitJobName   string `json:"initJobName"`

	//+kubebuilder:validation:Required
	Status string `json:"status"`

	//+kubebuilder:validation:Required
	LastNodeName string `json:"lastNodeName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// +kubebuilder:printcolumn:name="Node",type=string,JSONPath=`.spec.nodeName`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Job Status",type=string,JSONPath=`.status.initJobStatus`
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
	Sockets []Socket `json:"sockets"`
}

// Socket is a CPU socket of the Kubernetes node
type Socket struct {
	SocketId  int        `json:"socketId"`
	NumaNodes []NumaNode `json:"numaNodes"`
}

// NumaNode is a NUMA node of the parent socket
type NumaNode struct {
	NumaNodeId int    `json:"numaNodeId"`
	Cores      []Core `json:"cores"`
}

// Core is a physical CPU core of the parent NUMA node
type Core struct {
	CoreId int   `json:"coreId"`
	Cpus   []Cpu `json:"cpus"`
}

// Cpu is a logical CPU core of the parent core
type Cpu struct {
	CpuId int `json:"cpuId"`
}

func init() {
	SchemeBuilder.Register(&NodeCpuTopology{}, &NodeCpuTopologyList{})
}
