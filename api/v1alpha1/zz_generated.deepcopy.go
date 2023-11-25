//go:build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Core) DeepCopyInto(out *Core) {
	*out = *in
	if in.Cpus != nil {
		in, out := &in.Cpus, &out.Cpus
		*out = make([]Cpu, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Core.
func (in *Core) DeepCopy() *Core {
	if in == nil {
		return nil
	}
	out := new(Core)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cpu) DeepCopyInto(out *Cpu) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cpu.
func (in *Cpu) DeepCopy() *Cpu {
	if in == nil {
		return nil
	}
	out := new(Cpu)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CpuTopology) DeepCopyInto(out *CpuTopology) {
	*out = *in
	if in.Sockets != nil {
		in, out := &in.Sockets, &out.Sockets
		*out = make([]Socket, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NumaNodes != nil {
		in, out := &in.NumaNodes, &out.NumaNodes
		*out = make([]NumaNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CpuTopology.
func (in *CpuTopology) DeepCopy() *CpuTopology {
	if in == nil {
		return nil
	}
	out := new(CpuTopology)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeCpuTopology) DeepCopyInto(out *NodeCpuTopology) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeCpuTopology.
func (in *NodeCpuTopology) DeepCopy() *NodeCpuTopology {
	if in == nil {
		return nil
	}
	out := new(NodeCpuTopology)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeCpuTopology) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeCpuTopologyList) DeepCopyInto(out *NodeCpuTopologyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeCpuTopology, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeCpuTopologyList.
func (in *NodeCpuTopologyList) DeepCopy() *NodeCpuTopologyList {
	if in == nil {
		return nil
	}
	out := new(NodeCpuTopologyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeCpuTopologyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeCpuTopologySpec) DeepCopyInto(out *NodeCpuTopologySpec) {
	*out = *in
	in.Topology.DeepCopyInto(&out.Topology)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeCpuTopologySpec.
func (in *NodeCpuTopologySpec) DeepCopy() *NodeCpuTopologySpec {
	if in == nil {
		return nil
	}
	out := new(NodeCpuTopologySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeCpuTopologyStatus) DeepCopyInto(out *NodeCpuTopologyStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeCpuTopologyStatus.
func (in *NodeCpuTopologyStatus) DeepCopy() *NodeCpuTopologyStatus {
	if in == nil {
		return nil
	}
	out := new(NodeCpuTopologyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NumaNode) DeepCopyInto(out *NumaNode) {
	*out = *in
	if in.Cpus != nil {
		in, out := &in.Cpus, &out.Cpus
		*out = make([]Cpu, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NumaNode.
func (in *NumaNode) DeepCopy() *NumaNode {
	if in == nil {
		return nil
	}
	out := new(NumaNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Socket) DeepCopyInto(out *Socket) {
	*out = *in
	if in.Cores != nil {
		in, out := &in.Cores, &out.Cores
		*out = make([]Core, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Socket.
func (in *Socket) DeepCopy() *Socket {
	if in == nil {
		return nil
	}
	out := new(Socket)
	in.DeepCopyInto(out)
	return out
}