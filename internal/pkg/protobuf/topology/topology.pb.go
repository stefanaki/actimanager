// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.3
// source: internal/pkg/protobuf/topology/topology.proto

package topology

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NUMANode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Cpus []int64 `protobuf:"varint,2,rep,packed,name=cpus,proto3" json:"cpus,omitempty"`
}

func (x *NUMANode) Reset() {
	*x = NUMANode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NUMANode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NUMANode) ProtoMessage() {}

func (x *NUMANode) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NUMANode.ProtoReflect.Descriptor instead.
func (*NUMANode) Descriptor() ([]byte, []int) {
	return file_internal_pkg_protobuf_topology_topology_proto_rawDescGZIP(), []int{0}
}

func (x *NUMANode) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *NUMANode) GetCpus() []int64 {
	if x != nil {
		return x.Cpus
	}
	return nil
}

type Socket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64            `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Cores map[string]*Core `protobuf:"bytes,2,rep,name=cores,proto3" json:"cores,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Socket) Reset() {
	*x = Socket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Socket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Socket) ProtoMessage() {}

func (x *Socket) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Socket.ProtoReflect.Descriptor instead.
func (*Socket) Descriptor() ([]byte, []int) {
	return file_internal_pkg_protobuf_topology_topology_proto_rawDescGZIP(), []int{1}
}

func (x *Socket) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Socket) GetCores() map[string]*Core {
	if x != nil {
		return x.Cores
	}
	return nil
}

type Core struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Cpus []int64 `protobuf:"varint,2,rep,packed,name=cpus,proto3" json:"cpus,omitempty"`
}

func (x *Core) Reset() {
	*x = Core{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Core) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Core) ProtoMessage() {}

func (x *Core) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Core.ProtoReflect.Descriptor instead.
func (*Core) Descriptor() ([]byte, []int) {
	return file_internal_pkg_protobuf_topology_topology_proto_rawDescGZIP(), []int{2}
}

func (x *Core) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Core) GetCpus() []int64 {
	if x != nil {
		return x.Cpus
	}
	return nil
}

type TopologyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumaNodes map[string]*NUMANode `protobuf:"bytes,1,rep,name=numaNodes,proto3" json:"numaNodes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Sockets   map[string]*Socket   `protobuf:"bytes,2,rep,name=sockets,proto3" json:"sockets,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Cpus      []int64              `protobuf:"varint,3,rep,packed,name=cpus,proto3" json:"cpus,omitempty"`
}

func (x *TopologyResponse) Reset() {
	*x = TopologyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopologyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopologyResponse) ProtoMessage() {}

func (x *TopologyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_protobuf_topology_topology_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopologyResponse.ProtoReflect.Descriptor instead.
func (*TopologyResponse) Descriptor() ([]byte, []int) {
	return file_internal_pkg_protobuf_topology_topology_proto_rawDescGZIP(), []int{3}
}

func (x *TopologyResponse) GetNumaNodes() map[string]*NUMANode {
	if x != nil {
		return x.NumaNodes
	}
	return nil
}

func (x *TopologyResponse) GetSockets() map[string]*Socket {
	if x != nil {
		return x.Sockets
	}
	return nil
}

func (x *TopologyResponse) GetCpus() []int64 {
	if x != nil {
		return x.Cpus
	}
	return nil
}

var File_internal_pkg_protobuf_topology_topology_proto protoreflect.FileDescriptor

var file_internal_pkg_protobuf_topology_topology_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79,
	0x2f, 0x74, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2e, 0x0a, 0x08,
	0x4e, 0x55, 0x4d, 0x41, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x70, 0x75, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x04, 0x63, 0x70, 0x75, 0x73, 0x22, 0x83, 0x01, 0x0a,
	0x06, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x28, 0x0a, 0x05, 0x63, 0x6f, 0x72, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e,
	0x43, 0x6f, 0x72, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x63, 0x6f, 0x72, 0x65,
	0x73, 0x1a, 0x3f, 0x0a, 0x0a, 0x43, 0x6f, 0x72, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x1b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x05, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x2a, 0x0a, 0x04, 0x43, 0x6f, 0x72, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x70,
	0x75, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x04, 0x63, 0x70, 0x75, 0x73, 0x22, 0xae,
	0x02, 0x0a, 0x10, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x6e, 0x75, 0x6d, 0x61, 0x4e, 0x6f, 0x64, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4e, 0x75, 0x6d, 0x61, 0x4e, 0x6f,
	0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x6e, 0x75, 0x6d, 0x61, 0x4e, 0x6f,
	0x64, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x07, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x70, 0x75, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x03, 0x52, 0x04, 0x63, 0x70, 0x75,
	0x73, 0x1a, 0x47, 0x0a, 0x0e, 0x4e, 0x75, 0x6d, 0x61, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x4e, 0x55, 0x4d, 0x41, 0x4e, 0x6f, 0x64, 0x65, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x43, 0x0a, 0x0c, 0x53, 0x6f,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1d, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x53, 0x6f,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32,
	0x46, 0x0a, 0x08, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x12, 0x3a, 0x0a, 0x0b, 0x47,
	0x65, 0x74, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x11, 0x2e, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x2f, 0x74, 0x6f, 0x70, 0x6f,
	0x6c, 0x6f, 0x67, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_pkg_protobuf_topology_topology_proto_rawDescOnce sync.Once
	file_internal_pkg_protobuf_topology_topology_proto_rawDescData = file_internal_pkg_protobuf_topology_topology_proto_rawDesc
)

func file_internal_pkg_protobuf_topology_topology_proto_rawDescGZIP() []byte {
	file_internal_pkg_protobuf_topology_topology_proto_rawDescOnce.Do(func() {
		file_internal_pkg_protobuf_topology_topology_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_pkg_protobuf_topology_topology_proto_rawDescData)
	})
	return file_internal_pkg_protobuf_topology_topology_proto_rawDescData
}

var file_internal_pkg_protobuf_topology_topology_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_pkg_protobuf_topology_topology_proto_goTypes = []interface{}{
	(*NUMANode)(nil),         // 0: NUMANode
	(*Socket)(nil),           // 1: Socket
	(*Core)(nil),             // 2: Core
	(*TopologyResponse)(nil), // 3: TopologyResponse
	nil,                      // 4: Socket.CoresEntry
	nil,                      // 5: TopologyResponse.NumaNodesEntry
	nil,                      // 6: TopologyResponse.SocketsEntry
	(*emptypb.Empty)(nil),    // 7: google.protobuf.Empty
}
var file_internal_pkg_protobuf_topology_topology_proto_depIdxs = []int32{
	4, // 0: Socket.cores:type_name -> Socket.CoresEntry
	5, // 1: TopologyResponse.numaNodes:type_name -> TopologyResponse.NumaNodesEntry
	6, // 2: TopologyResponse.sockets:type_name -> TopologyResponse.SocketsEntry
	2, // 3: Socket.CoresEntry.value:type_name -> Core
	0, // 4: TopologyResponse.NumaNodesEntry.value:type_name -> NUMANode
	1, // 5: TopologyResponse.SocketsEntry.value:type_name -> Socket
	7, // 6: Topology.GetTopology:input_type -> google.protobuf.Empty
	3, // 7: Topology.GetTopology:output_type -> TopologyResponse
	7, // [7:8] is the sub-list for method output_type
	6, // [6:7] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_internal_pkg_protobuf_topology_topology_proto_init() }
func file_internal_pkg_protobuf_topology_topology_proto_init() {
	if File_internal_pkg_protobuf_topology_topology_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_pkg_protobuf_topology_topology_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NUMANode); i {
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
		file_internal_pkg_protobuf_topology_topology_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Socket); i {
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
		file_internal_pkg_protobuf_topology_topology_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Core); i {
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
		file_internal_pkg_protobuf_topology_topology_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopologyResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_pkg_protobuf_topology_topology_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_pkg_protobuf_topology_topology_proto_goTypes,
		DependencyIndexes: file_internal_pkg_protobuf_topology_topology_proto_depIdxs,
		MessageInfos:      file_internal_pkg_protobuf_topology_topology_proto_msgTypes,
	}.Build()
	File_internal_pkg_protobuf_topology_topology_proto = out.File
	file_internal_pkg_protobuf_topology_topology_proto_rawDesc = nil
	file_internal_pkg_protobuf_topology_topology_proto_goTypes = nil
	file_internal_pkg_protobuf_topology_topology_proto_depIdxs = nil
}
