// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: admin_l4_commands.proto

package admin_commands_message

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
	L4Commands_CheckHealth_FullMethodName = "/admin_L4_commands.L4Commands/CheckHealth"
)

// L4CommandsClient is the client API for L4Commands service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type L4CommandsClient interface {
	CheckHealth(ctx context.Context, in *L4HealthCheckRequest, opts ...grpc.CallOption) (*L4HealthCheckResponse, error)
}

type l4CommandsClient struct {
	cc grpc.ClientConnInterface
}

func NewL4CommandsClient(cc grpc.ClientConnInterface) L4CommandsClient {
	return &l4CommandsClient{cc}
}

func (c *l4CommandsClient) CheckHealth(ctx context.Context, in *L4HealthCheckRequest, opts ...grpc.CallOption) (*L4HealthCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(L4HealthCheckResponse)
	err := c.cc.Invoke(ctx, L4Commands_CheckHealth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// L4CommandsServer is the server API for L4Commands service.
// All implementations must embed UnimplementedL4CommandsServer
// for forward compatibility.
type L4CommandsServer interface {
	CheckHealth(context.Context, *L4HealthCheckRequest) (*L4HealthCheckResponse, error)
	mustEmbedUnimplementedL4CommandsServer()
}

// UnimplementedL4CommandsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedL4CommandsServer struct{}

func (UnimplementedL4CommandsServer) CheckHealth(context.Context, *L4HealthCheckRequest) (*L4HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHealth not implemented")
}
func (UnimplementedL4CommandsServer) mustEmbedUnimplementedL4CommandsServer() {}
func (UnimplementedL4CommandsServer) testEmbeddedByValue()                    {}

// UnsafeL4CommandsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to L4CommandsServer will
// result in compilation errors.
type UnsafeL4CommandsServer interface {
	mustEmbedUnimplementedL4CommandsServer()
}

func RegisterL4CommandsServer(s grpc.ServiceRegistrar, srv L4CommandsServer) {
	// If the following call pancis, it indicates UnimplementedL4CommandsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&L4Commands_ServiceDesc, srv)
}

func _L4Commands_CheckHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(L4HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(L4CommandsServer).CheckHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: L4Commands_CheckHealth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(L4CommandsServer).CheckHealth(ctx, req.(*L4HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// L4Commands_ServiceDesc is the grpc.ServiceDesc for L4Commands service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var L4Commands_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "admin_L4_commands.L4Commands",
	HandlerType: (*L4CommandsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckHealth",
			Handler:    _L4Commands_CheckHealth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin_l4_commands.proto",
}
