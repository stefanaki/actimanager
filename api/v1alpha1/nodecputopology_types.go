package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeCpuTopologyResourceStatus string
type NodeCpuTopologyJobStatus string

const (
	StatusNeedsSync    NodeCpuTopologyResourceStatus = "NeedsSync"
	StatusNodeNotFound NodeCpuTopologyResourceStatus = "NodeNotFound"
	StatusFresh        NodeCpuTopologyResourceStatus = "Fresh"
	StatusJobNone      NodeCpuTopologyJobStatus      = "None"
	StatusJobPending   NodeCpuTopologyJobStatus      = "Pending"
	StatusJobCompleted NodeCpuTopologyJobStatus      = "Completed"
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
	InitJobStatus NodeCpuTopologyJobStatus `json:"initJobStatus"`
	InitJobName   string                   `json:"initJobName"`

	//+kubebuilder:validation:Required
	ResourceStatus NodeCpuTopologyResourceStatus `json:"resourceStatus"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=nct
// +kubebuilder:printcolumn:name="Node",type=string,JSONPath=`.spec.nodeName`
// +kubebuilder:printcolumn:name="Resource Status",type=string,JSONPath=`.status.resourceStatus`
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
	//+kubebuilder:validation:Optional
	Sockets map[string]Socket `json:"sockets"`
	//+kubebuilder:validation:Optional
	NumaNodes map[string]NumaNode `json:"numaNodes"`
}

// Socket is a CPU socket of the Kubernetes node
type Socket struct {
	Cores map[string]Core `json:"cores"`
}

// NumaNode is a NUMA node of the Kubernetes node
type NumaNode struct {
	Cpus map[string]Cpu `json:"cpus"`
}

// Core is a physical CPU core of the parent socket
type Core struct {
	Cpus map[string]Cpu `json:"cpus"`
}

// Cpu is a logical CPU core of the parent core
type Cpu struct {
	CpuId int `json:"cpuId"`
}

func init() {
	SchemeBuilder.Register(&NodeCpuTopology{}, &NodeCpuTopologyList{})
}
