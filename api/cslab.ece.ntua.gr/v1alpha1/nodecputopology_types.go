package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeCpuTopologyResourceStatus string

const (
	StatusNeedsSync      NodeCpuTopologyResourceStatus = "NeedsSync"
	StatusNodeNotFound   NodeCpuTopologyResourceStatus = "NodeNotFound"
	StatusFresh          NodeCpuTopologyResourceStatus = "Fresh"
	StatusTopologyFailed NodeCpuTopologyResourceStatus = "Failed"
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
	ResourceStatus NodeCpuTopologyResourceStatus `json:"resourceStatus"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=nct
// +kubebuilder:printcolumn:name="Node",type=string,JSONPath=`.spec.nodeName`
// +kubebuilder:printcolumn:name="Resource Status",type=string,JSONPath=`.status.resourceStatus`

// NodeCpuTopology is the Schema for the nodecputopologies API
type NodeCpuTopology struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeCpuTopologySpec   `json:"spec,omitempty"`
	Status NodeCpuTopologyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

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
	//+kubebuilder:validation:Optional
	ListCpus []int `json:"listCpus"`
}

// Socket is a CPU socket of the Kubernetes node
type Socket struct {
	Cores    map[string]Core `json:"cores"`
	ListCpus []int           `json:"listCpus"`
}

// NumaNode is a NUMA node of the Kubernetes node
type NumaNode struct {
	Cpus     map[string]Cpu `json:"cpus"`
	ListCpus []int          `json:"listCpus"`
}

// Core is a physical CPU core of the parent socket
type Core struct {
	Cpus     map[string]Cpu `json:"cpus"`
	ListCpus []int          `json:"listCpus"`
}

// Cpu is a logical CPU core of the parent core
type Cpu struct {
	CpuId int `json:"cpuId"`
}

func init() {
	SchemeBuilder.Register(&NodeCpuTopology{}, &NodeCpuTopologyList{})
}
