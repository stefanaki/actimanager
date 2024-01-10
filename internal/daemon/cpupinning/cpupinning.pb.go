// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: internal/daemon/cpupinning/cpupinning.proto

package cpupinning

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ResponseStatus int32

const (
	ResponseStatus_SUCCESSFUL ResponseStatus = 0
	ResponseStatus_ERROR      ResponseStatus = 1
)

// Enum value maps for ResponseStatus.
var (
	ResponseStatus_name = map[int32]string{
		0: "SUCCESSFUL",
		1: "ERROR",
	}
	ResponseStatus_value = map[string]int32{
		"SUCCESSFUL": 0,
		"ERROR":      1,
	}
)

func (x ResponseStatus) Enum() *ResponseStatus {
	p := new(ResponseStatus)
	*p = x
	return p
}

func (x ResponseStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_daemon_cpupinning_cpupinning_proto_enumTypes[0].Descriptor()
}

func (ResponseStatus) Type() protoreflect.EnumType {
	return &file_internal_daemon_cpupinning_cpupinning_proto_enumTypes[0]
}

func (x ResponseStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseStatus.Descriptor instead.
func (ResponseStatus) EnumDescriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{0}
}

type ErrorType int32

const (
	ErrorType_UNKNOWN           ErrorType = 0
	ErrorType_INVALID_ARGUMENT  ErrorType = 1
	ErrorType_PERMISSION_DENIED ErrorType = 2
	ErrorType_INTERNAL_ERROR    ErrorType = 3
)

// Enum value maps for ErrorType.
var (
	ErrorType_name = map[int32]string{
		0: "UNKNOWN",
		1: "INVALID_ARGUMENT",
		2: "PERMISSION_DENIED",
		3: "INTERNAL_ERROR",
	}
	ErrorType_value = map[string]int32{
		"UNKNOWN":           0,
		"INVALID_ARGUMENT":  1,
		"PERMISSION_DENIED": 2,
		"INTERNAL_ERROR":    3,
	}
)

func (x ErrorType) Enum() *ErrorType {
	p := new(ErrorType)
	*p = x
	return p
}

func (x ErrorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_daemon_cpupinning_cpupinning_proto_enumTypes[1].Descriptor()
}

func (ErrorType) Type() protoreflect.EnumType {
	return &file_internal_daemon_cpupinning_cpupinning_proto_enumTypes[1]
}

func (x ErrorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorType.Descriptor instead.
func (ErrorType) EnumDescriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{1}
}

type CpuSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cpu []int32 `protobuf:"varint,1,rep,packed,name=cpu,proto3" json:"cpu,omitempty"`
}

func (x *CpuSet) Reset() {
	*x = CpuSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CpuSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CpuSet) ProtoMessage() {}

func (x *CpuSet) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CpuSet.ProtoReflect.Descriptor instead.
func (*CpuSet) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{0}
}

func (x *CpuSet) GetCpu() []int32 {
	if x != nil {
		return x.Cpu
	}
	return nil
}

type ResourceInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestedCpus   int32  `protobuf:"varint,1,opt,name=requestedCpus,proto3" json:"requestedCpus,omitempty"`
	LimitCpus       int32  `protobuf:"varint,2,opt,name=limitCpus,proto3" json:"limitCpus,omitempty"`
	RequestedMemory string `protobuf:"bytes,3,opt,name=requestedMemory,proto3" json:"requestedMemory,omitempty"`
	LimitMemory     string `protobuf:"bytes,4,opt,name=limitMemory,proto3" json:"limitMemory,omitempty"`
}

func (x *ResourceInfo) Reset() {
	*x = ResourceInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceInfo) ProtoMessage() {}

func (x *ResourceInfo) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceInfo.ProtoReflect.Descriptor instead.
func (*ResourceInfo) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{1}
}

func (x *ResourceInfo) GetRequestedCpus() int32 {
	if x != nil {
		return x.RequestedCpus
	}
	return 0
}

func (x *ResourceInfo) GetLimitCpus() int32 {
	if x != nil {
		return x.LimitCpus
	}
	return 0
}

func (x *ResourceInfo) GetRequestedMemory() string {
	if x != nil {
		return x.RequestedMemory
	}
	return ""
}

func (x *ResourceInfo) GetLimitMemory() string {
	if x != nil {
		return x.LimitMemory
	}
	return ""
}

type Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string        `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Resources *ResourceInfo `protobuf:"bytes,3,opt,name=resources,proto3" json:"resources,omitempty"`
}

func (x *Container) Reset() {
	*x = Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Container) ProtoMessage() {}

func (x *Container) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Container.ProtoReflect.Descriptor instead.
func (*Container) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{2}
}

func (x *Container) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Container) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Container) GetResources() *ResourceInfo {
	if x != nil {
		return x.Resources
	}
	return nil
}

type Pod struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Namespace  string       `protobuf:"bytes,3,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Containers []*Container `protobuf:"bytes,5,rep,name=containers,proto3" json:"containers,omitempty"`
}

func (x *Pod) Reset() {
	*x = Pod{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pod) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pod) ProtoMessage() {}

func (x *Pod) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pod.ProtoReflect.Descriptor instead.
func (*Pod) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{3}
}

func (x *Pod) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Pod) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Pod) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Pod) GetContainers() []*Container {
	if x != nil {
		return x.Containers
	}
	return nil
}

type ApplyPinningRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pod    *Pod    `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
	CpuSet *CpuSet `protobuf:"bytes,2,opt,name=cpuSet,proto3" json:"cpuSet,omitempty"`
}

func (x *ApplyPinningRequest) Reset() {
	*x = ApplyPinningRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplyPinningRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplyPinningRequest) ProtoMessage() {}

func (x *ApplyPinningRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplyPinningRequest.ProtoReflect.Descriptor instead.
func (*ApplyPinningRequest) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{4}
}

func (x *ApplyPinningRequest) GetPod() *Pod {
	if x != nil {
		return x.Pod
	}
	return nil
}

func (x *ApplyPinningRequest) GetCpuSet() *CpuSet {
	if x != nil {
		return x.CpuSet
	}
	return nil
}

type RemovePinningRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pod *Pod `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
}

func (x *RemovePinningRequest) Reset() {
	*x = RemovePinningRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemovePinningRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemovePinningRequest) ProtoMessage() {}

func (x *RemovePinningRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemovePinningRequest.ProtoReflect.Descriptor instead.
func (*RemovePinningRequest) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{5}
}

func (x *RemovePinningRequest) GetPod() *Pod {
	if x != nil {
		return x.Pod
	}
	return nil
}

type UpdatePinningRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pod    *Pod    `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
	CpuSet *CpuSet `protobuf:"bytes,2,opt,name=cpuSet,proto3" json:"cpuSet,omitempty"`
}

func (x *UpdatePinningRequest) Reset() {
	*x = UpdatePinningRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePinningRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePinningRequest) ProtoMessage() {}

func (x *UpdatePinningRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePinningRequest.ProtoReflect.Descriptor instead.
func (*UpdatePinningRequest) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{6}
}

func (x *UpdatePinningRequest) GetPod() *Pod {
	if x != nil {
		return x.Pod
	}
	return nil
}

func (x *UpdatePinningRequest) GetCpuSet() *CpuSet {
	if x != nil {
		return x.CpuSet
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status    ResponseStatus `protobuf:"varint,1,opt,name=status,proto3,enum=ResponseStatus" json:"status,omitempty"`
	Message   *string        `protobuf:"bytes,2,opt,name=message,proto3,oneof" json:"message,omitempty"`
	ErrorType *ErrorType     `protobuf:"varint,3,opt,name=errorType,proto3,enum=ErrorType,oneof" json:"errorType,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP(), []int{7}
}

func (x *Response) GetStatus() ResponseStatus {
	if x != nil {
		return x.Status
	}
	return ResponseStatus_SUCCESSFUL
}

func (x *Response) GetMessage() string {
	if x != nil && x.Message != nil {
		return *x.Message
	}
	return ""
}

func (x *Response) GetErrorType() ErrorType {
	if x != nil && x.ErrorType != nil {
		return *x.ErrorType
	}
	return ErrorType_UNKNOWN
}

var File_internal_daemon_cpupinning_cpupinning_proto protoreflect.FileDescriptor

var file_internal_daemon_cpupinning_cpupinning_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f,
	0x6e, 0x2f, 0x63, 0x70, 0x75, 0x70, 0x69, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x70, 0x75,
	0x70, 0x69, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1a, 0x0a,
	0x06, 0x43, 0x70, 0x75, 0x53, 0x65, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x05, 0x52, 0x03, 0x63, 0x70, 0x75, 0x22, 0x9e, 0x01, 0x0a, 0x0c, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x43, 0x70, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x43, 0x70, 0x75, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x43, 0x70, 0x75, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x43, 0x70, 0x75, 0x73, 0x12, 0x28,
	0x0a, 0x0f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x6d, 0x6f, 0x72,
	0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x22, 0x5c, 0x0a, 0x09, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x09, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x22, 0x73, 0x0a, 0x03, 0x50, 0x6f, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x12, 0x2a, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73, 0x22, 0x4e, 0x0a,
	0x13, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x50, 0x69, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x04, 0x2e, 0x50, 0x6f, 0x64, 0x52, 0x03, 0x70, 0x6f, 0x64, 0x12, 0x1f, 0x0a, 0x06,
	0x63, 0x70, 0x75, 0x53, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x43,
	0x70, 0x75, 0x53, 0x65, 0x74, 0x52, 0x06, 0x63, 0x70, 0x75, 0x53, 0x65, 0x74, 0x22, 0x2e, 0x0a,
	0x14, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x69, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x04, 0x2e, 0x50, 0x6f, 0x64, 0x52, 0x03, 0x70, 0x6f, 0x64, 0x22, 0x4f, 0x0a,
	0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x69, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x04, 0x2e, 0x50, 0x6f, 0x64, 0x52, 0x03, 0x70, 0x6f, 0x64, 0x12, 0x1f, 0x0a,
	0x06, 0x63, 0x70, 0x75, 0x53, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e,
	0x43, 0x70, 0x75, 0x53, 0x65, 0x74, 0x52, 0x06, 0x63, 0x70, 0x75, 0x53, 0x65, 0x74, 0x22, 0x9b,
	0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x48, 0x01, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x88,
	0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x0c,
	0x0a, 0x0a, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x2a, 0x2b, 0x0a, 0x0e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e,
	0x0a, 0x0a, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x46, 0x55, 0x4c, 0x10, 0x00, 0x12, 0x09,
	0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x2a, 0x59, 0x0a, 0x09, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57,
	0x4e, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41,
	0x52, 0x47, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x45, 0x52,
	0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4e, 0x49, 0x45, 0x44, 0x10, 0x02,
	0x12, 0x12, 0x0a, 0x0e, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0x03, 0x32, 0x70, 0x0a, 0x0a, 0x43, 0x70, 0x75, 0x50, 0x69, 0x6e, 0x6e, 0x69,
	0x6e, 0x67, 0x12, 0x2f, 0x0a, 0x0c, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x50, 0x69, 0x6e, 0x6e, 0x69,
	0x6e, 0x67, 0x12, 0x14, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x50, 0x69, 0x6e, 0x6e, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x0d, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x69, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x12, 0x15, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x69, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x2f, 0x63, 0x70, 0x75, 0x70, 0x69,
	0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_daemon_cpupinning_cpupinning_proto_rawDescOnce sync.Once
	file_internal_daemon_cpupinning_cpupinning_proto_rawDescData = file_internal_daemon_cpupinning_cpupinning_proto_rawDesc
)

func file_internal_daemon_cpupinning_cpupinning_proto_rawDescGZIP() []byte {
	file_internal_daemon_cpupinning_cpupinning_proto_rawDescOnce.Do(func() {
		file_internal_daemon_cpupinning_cpupinning_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_daemon_cpupinning_cpupinning_proto_rawDescData)
	})
	return file_internal_daemon_cpupinning_cpupinning_proto_rawDescData
}

var file_internal_daemon_cpupinning_cpupinning_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_internal_daemon_cpupinning_cpupinning_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_internal_daemon_cpupinning_cpupinning_proto_goTypes = []interface{}{
	(ResponseStatus)(0),          // 0: ResponseStatus
	(ErrorType)(0),               // 1: ErrorType
	(*CpuSet)(nil),               // 2: CpuSet
	(*ResourceInfo)(nil),         // 3: ResourceInfo
	(*Container)(nil),            // 4: Container
	(*Pod)(nil),                  // 5: Pod
	(*ApplyPinningRequest)(nil),  // 6: ApplyPinningRequest
	(*RemovePinningRequest)(nil), // 7: RemovePinningRequest
	(*UpdatePinningRequest)(nil), // 8: UpdatePinningRequest
	(*Response)(nil),             // 9: Response
}
var file_internal_daemon_cpupinning_cpupinning_proto_depIdxs = []int32{
	3,  // 0: Container.resources:type_name -> ResourceInfo
	4,  // 1: Pod.containers:type_name -> Container
	5,  // 2: ApplyPinningRequest.pod:type_name -> Pod
	2,  // 3: ApplyPinningRequest.cpuSet:type_name -> CpuSet
	5,  // 4: RemovePinningRequest.pod:type_name -> Pod
	5,  // 5: UpdatePinningRequest.pod:type_name -> Pod
	2,  // 6: UpdatePinningRequest.cpuSet:type_name -> CpuSet
	0,  // 7: Response.status:type_name -> ResponseStatus
	1,  // 8: Response.errorType:type_name -> ErrorType
	6,  // 9: CpuPinning.ApplyPinning:input_type -> ApplyPinningRequest
	7,  // 10: CpuPinning.RemovePinning:input_type -> RemovePinningRequest
	9,  // 11: CpuPinning.ApplyPinning:output_type -> Response
	9,  // 12: CpuPinning.RemovePinning:output_type -> Response
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_internal_daemon_cpupinning_cpupinning_proto_init() }
func file_internal_daemon_cpupinning_cpupinning_proto_init() {
	if File_internal_daemon_cpupinning_cpupinning_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CpuSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Container); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pod); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplyPinningRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemovePinningRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePinningRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_internal_daemon_cpupinning_cpupinning_proto_msgTypes[7].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_daemon_cpupinning_cpupinning_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_daemon_cpupinning_cpupinning_proto_goTypes,
		DependencyIndexes: file_internal_daemon_cpupinning_cpupinning_proto_depIdxs,
		EnumInfos:         file_internal_daemon_cpupinning_cpupinning_proto_enumTypes,
		MessageInfos:      file_internal_daemon_cpupinning_cpupinning_proto_msgTypes,
	}.Build()
	File_internal_daemon_cpupinning_cpupinning_proto = out.File
	file_internal_daemon_cpupinning_cpupinning_proto_rawDesc = nil
	file_internal_daemon_cpupinning_cpupinning_proto_goTypes = nil
	file_internal_daemon_cpupinning_cpupinning_proto_depIdxs = nil
}
