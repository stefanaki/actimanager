// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: cpu_pinning_daemon.proto

package daemon

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	CpuPinningDaemon_ApplyPinning_FullMethodName  = "/CpuPinningDaemon/ApplyPinning"
	CpuPinningDaemon_RemovePinning_FullMethodName = "/CpuPinningDaemon/RemovePinning"
	CpuPinningDaemon_UpdatePinning_FullMethodName = "/CpuPinningDaemon/UpdatePinning"
)

// CpuPinningDaemonClient is the client API for CpuPinningDaemon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CpuPinningDaemonClient interface {
	ApplyPinning(ctx context.Context, in *ApplyPinningRequest, opts ...grpc.CallOption) (*ApplyPinningResponse, error)
	RemovePinning(ctx context.Context, in *RemovePinningRequest, opts ...grpc.CallOption) (*RemovePinningResponse, error)
	UpdatePinning(ctx context.Context, in *UpdatePinningRequest, opts ...grpc.CallOption) (*UpdatePinningResponse, error)
}

type cpuPinningDaemonClient struct {
	cc grpc.ClientConnInterface
}

func NewCpuPinningDaemonClient(cc grpc.ClientConnInterface) CpuPinningDaemonClient {
	return &cpuPinningDaemonClient{cc}
}

func (c *cpuPinningDaemonClient) ApplyPinning(ctx context.Context, in *ApplyPinningRequest, opts ...grpc.CallOption) (*ApplyPinningResponse, error) {
	out := new(ApplyPinningResponse)
	err := c.cc.Invoke(ctx, CpuPinningDaemon_ApplyPinning_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cpuPinningDaemonClient) RemovePinning(ctx context.Context, in *RemovePinningRequest, opts ...grpc.CallOption) (*RemovePinningResponse, error) {
	out := new(RemovePinningResponse)
	err := c.cc.Invoke(ctx, CpuPinningDaemon_RemovePinning_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cpuPinningDaemonClient) UpdatePinning(ctx context.Context, in *UpdatePinningRequest, opts ...grpc.CallOption) (*UpdatePinningResponse, error) {
	out := new(UpdatePinningResponse)
	err := c.cc.Invoke(ctx, CpuPinningDaemon_UpdatePinning_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CpuPinningDaemonServer is the server API for CpuPinningDaemon service.
// All implementations must embed UnimplementedCpuPinningDaemonServer
// for forward compatibility
type CpuPinningDaemonServer interface {
	ApplyPinning(context.Context, *ApplyPinningRequest) (*ApplyPinningResponse, error)
	RemovePinning(context.Context, *RemovePinningRequest) (*RemovePinningResponse, error)
	UpdatePinning(context.Context, *UpdatePinningRequest) (*UpdatePinningResponse, error)
	mustEmbedUnimplementedCpuPinningDaemonServer()
}

// UnimplementedCpuPinningDaemonServer must be embedded to have forward compatible implementations.
type UnimplementedCpuPinningDaemonServer struct {
}

func (UnimplementedCpuPinningDaemonServer) ApplyPinning(context.Context, *ApplyPinningRequest) (*ApplyPinningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplyPinning not implemented")
}
func (UnimplementedCpuPinningDaemonServer) RemovePinning(context.Context, *RemovePinningRequest) (*RemovePinningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemovePinning not implemented")
}
func (UnimplementedCpuPinningDaemonServer) UpdatePinning(context.Context, *UpdatePinningRequest) (*UpdatePinningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePinning not implemented")
}
func (UnimplementedCpuPinningDaemonServer) mustEmbedUnimplementedCpuPinningDaemonServer() {}

// UnsafeCpuPinningDaemonServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CpuPinningDaemonServer will
// result in compilation errors.
type UnsafeCpuPinningDaemonServer interface {
	mustEmbedUnimplementedCpuPinningDaemonServer()
}

func RegisterCpuPinningDaemonServer(s grpc.ServiceRegistrar, srv CpuPinningDaemonServer) {
	s.RegisterService(&CpuPinningDaemon_ServiceDesc, srv)
}

func _CpuPinningDaemon_ApplyPinning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyPinningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CpuPinningDaemonServer).ApplyPinning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CpuPinningDaemon_ApplyPinning_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CpuPinningDaemonServer).ApplyPinning(ctx, req.(*ApplyPinningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CpuPinningDaemon_RemovePinning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemovePinningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CpuPinningDaemonServer).RemovePinning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CpuPinningDaemon_RemovePinning_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CpuPinningDaemonServer).RemovePinning(ctx, req.(*RemovePinningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CpuPinningDaemon_UpdatePinning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePinningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CpuPinningDaemonServer).UpdatePinning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CpuPinningDaemon_UpdatePinning_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CpuPinningDaemonServer).UpdatePinning(ctx, req.(*UpdatePinningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CpuPinningDaemon_ServiceDesc is the grpc.ServiceDesc for CpuPinningDaemon service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CpuPinningDaemon_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CpuPinningDaemon",
	HandlerType: (*CpuPinningDaemonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ApplyPinning",
			Handler:    _CpuPinningDaemon_ApplyPinning_Handler,
		},
		{
			MethodName: "RemovePinning",
			Handler:    _CpuPinningDaemon_RemovePinning_Handler,
		},
		{
			MethodName: "UpdatePinning",
			Handler:    _CpuPinningDaemon_UpdatePinning_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cpu_pinning_daemon.proto",
}