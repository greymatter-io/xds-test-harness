// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package adapter

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

// AdapterClient is the client API for Adapter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdapterClient interface {
	SetState(ctx context.Context, in *Snapshot, opts ...grpc.CallOption) (*SetStateResponse, error)
	ClearState(ctx context.Context, in *ClearRequest, opts ...grpc.CallOption) (*ClearResponse, error)
}

type adapterClient struct {
	cc grpc.ClientConnInterface
}

func NewAdapterClient(cc grpc.ClientConnInterface) AdapterClient {
	return &adapterClient{cc}
}

func (c *adapterClient) SetState(ctx context.Context, in *Snapshot, opts ...grpc.CallOption) (*SetStateResponse, error) {
	out := new(SetStateResponse)
	err := c.cc.Invoke(ctx, "/adapter.Adapter/SetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adapterClient) ClearState(ctx context.Context, in *ClearRequest, opts ...grpc.CallOption) (*ClearResponse, error) {
	out := new(ClearResponse)
	err := c.cc.Invoke(ctx, "/adapter.Adapter/ClearState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdapterServer is the server API for Adapter service.
// All implementations must embed UnimplementedAdapterServer
// for forward compatibility
type AdapterServer interface {
	SetState(context.Context, *Snapshot) (*SetStateResponse, error)
	ClearState(context.Context, *ClearRequest) (*ClearResponse, error)
	mustEmbedUnimplementedAdapterServer()
}

// UnimplementedAdapterServer must be embedded to have forward compatible implementations.
type UnimplementedAdapterServer struct {
}

func (UnimplementedAdapterServer) SetState(context.Context, *Snapshot) (*SetStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetState not implemented")
}
func (UnimplementedAdapterServer) ClearState(context.Context, *ClearRequest) (*ClearResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearState not implemented")
}
func (UnimplementedAdapterServer) mustEmbedUnimplementedAdapterServer() {}

// UnsafeAdapterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdapterServer will
// result in compilation errors.
type UnsafeAdapterServer interface {
	mustEmbedUnimplementedAdapterServer()
}

func RegisterAdapterServer(s grpc.ServiceRegistrar, srv AdapterServer) {
	s.RegisterService(&Adapter_ServiceDesc, srv)
}

func _Adapter_SetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Snapshot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdapterServer).SetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/adapter.Adapter/SetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdapterServer).SetState(ctx, req.(*Snapshot))
	}
	return interceptor(ctx, in, info, handler)
}

func _Adapter_ClearState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdapterServer).ClearState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/adapter.Adapter/ClearState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdapterServer).ClearState(ctx, req.(*ClearRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Adapter_ServiceDesc is the grpc.ServiceDesc for Adapter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Adapter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "adapter.Adapter",
	HandlerType: (*AdapterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetState",
			Handler:    _Adapter_SetState_Handler,
		},
		{
			MethodName: "ClearState",
			Handler:    _Adapter_ClearState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/adapter/adapter.proto",
}
