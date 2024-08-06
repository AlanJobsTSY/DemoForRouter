// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.0
// - protoc             v5.27.2
// source: minigame_router.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MiniGameRouter_RegisterService_FullMethodName = "/proto.MiniGameRouter/RegisterService"
	MiniGameRouter_DiscoverService_FullMethodName = "/proto.MiniGameRouter/DiscoverService"
	MiniGameRouter_SayHello_FullMethodName        = "/proto.MiniGameRouter/SayHello"
)

// MiniGameRouterClient is the client API for MiniGameRouter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MiniGameRouterClient interface {
	RegisterService(ctx context.Context, in *RegisterServiceRequest, opts ...grpc.CallOption) (*RegisterServiceResponse, error)
	DiscoverService(ctx context.Context, in *DiscoverServiceRequest, opts ...grpc.CallOption) (*DiscoverServiceResponse, error)
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
}

type miniGameRouterClient struct {
	cc grpc.ClientConnInterface
}

func NewMiniGameRouterClient(cc grpc.ClientConnInterface) MiniGameRouterClient {
	return &miniGameRouterClient{cc}
}

func (c *miniGameRouterClient) RegisterService(ctx context.Context, in *RegisterServiceRequest, opts ...grpc.CallOption) (*RegisterServiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterServiceResponse)
	err := c.cc.Invoke(ctx, MiniGameRouter_RegisterService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *miniGameRouterClient) DiscoverService(ctx context.Context, in *DiscoverServiceRequest, opts ...grpc.CallOption) (*DiscoverServiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DiscoverServiceResponse)
	err := c.cc.Invoke(ctx, MiniGameRouter_DiscoverService_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *miniGameRouterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, MiniGameRouter_SayHello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MiniGameRouterServer is the server API for MiniGameRouter service.
// All implementations must embed UnimplementedMiniGameRouterServer
// for forward compatibility.
type MiniGameRouterServer interface {
	RegisterService(context.Context, *RegisterServiceRequest) (*RegisterServiceResponse, error)
	DiscoverService(context.Context, *DiscoverServiceRequest) (*DiscoverServiceResponse, error)
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	mustEmbedUnimplementedMiniGameRouterServer()
}

// UnimplementedMiniGameRouterServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMiniGameRouterServer struct{}

func (UnimplementedMiniGameRouterServer) RegisterService(context.Context, *RegisterServiceRequest) (*RegisterServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterService not implemented")
}
func (UnimplementedMiniGameRouterServer) DiscoverService(context.Context, *DiscoverServiceRequest) (*DiscoverServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiscoverService not implemented")
}
func (UnimplementedMiniGameRouterServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedMiniGameRouterServer) mustEmbedUnimplementedMiniGameRouterServer() {}
func (UnimplementedMiniGameRouterServer) testEmbeddedByValue()                        {}

// UnsafeMiniGameRouterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MiniGameRouterServer will
// result in compilation errors.
type UnsafeMiniGameRouterServer interface {
	mustEmbedUnimplementedMiniGameRouterServer()
}

func RegisterMiniGameRouterServer(s grpc.ServiceRegistrar, srv MiniGameRouterServer) {
	// If the following call pancis, it indicates UnimplementedMiniGameRouterServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MiniGameRouter_ServiceDesc, srv)
}

func _MiniGameRouter_RegisterService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniGameRouterServer).RegisterService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniGameRouter_RegisterService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniGameRouterServer).RegisterService(ctx, req.(*RegisterServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MiniGameRouter_DiscoverService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoverServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniGameRouterServer).DiscoverService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniGameRouter_DiscoverService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniGameRouterServer).DiscoverService(ctx, req.(*DiscoverServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MiniGameRouter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniGameRouterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniGameRouter_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniGameRouterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MiniGameRouter_ServiceDesc is the grpc.ServiceDesc for MiniGameRouter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MiniGameRouter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MiniGameRouter",
	HandlerType: (*MiniGameRouterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterService",
			Handler:    _MiniGameRouter_RegisterService_Handler,
		},
		{
			MethodName: "DiscoverService",
			Handler:    _MiniGameRouter_DiscoverService_Handler,
		},
		{
			MethodName: "SayHello",
			Handler:    _MiniGameRouter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "minigame_router.proto",
}