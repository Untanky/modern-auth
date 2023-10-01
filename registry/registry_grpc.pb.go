// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: registry.proto

package registry

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

// RegistryClient is the client API for Registry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistryClient interface {
	Register(ctx context.Context, in *RegistrationInfo, opts ...grpc.CallOption) (*RegistrationResponse, error)
	Unregister(ctx context.Context, in *RegistrationResponse, opts ...grpc.CallOption) (*Empty, error)
	Subscribe(ctx context.Context, in *EndpointRequest, opts ...grpc.CallOption) (Registry_SubscribeClient, error)
}

type registryClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistryClient(cc grpc.ClientConnInterface) RegistryClient {
	return &registryClient{cc}
}

func (c *registryClient) Register(ctx context.Context, in *RegistrationInfo, opts ...grpc.CallOption) (*RegistrationResponse, error) {
	out := new(RegistrationResponse)
	err := c.cc.Invoke(ctx, "/registry.Registry/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) Unregister(ctx context.Context, in *RegistrationResponse, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/registry.Registry/Unregister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) Subscribe(ctx context.Context, in *EndpointRequest, opts ...grpc.CallOption) (Registry_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Registry_ServiceDesc.Streams[0], "/registry.Registry/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &registrySubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Registry_SubscribeClient interface {
	Recv() (*EndpointResponse, error)
	grpc.ClientStream
}

type registrySubscribeClient struct {
	grpc.ClientStream
}

func (x *registrySubscribeClient) Recv() (*EndpointResponse, error) {
	m := new(EndpointResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RegistryServer is the server API for Registry service.
// All implementations must embed UnimplementedRegistryServer
// for forward compatibility
type RegistryServer interface {
	Register(context.Context, *RegistrationInfo) (*RegistrationResponse, error)
	Unregister(context.Context, *RegistrationResponse) (*Empty, error)
	Subscribe(*EndpointRequest, Registry_SubscribeServer) error
	mustEmbedUnimplementedRegistryServer()
}

// UnimplementedRegistryServer must be embedded to have forward compatible implementations.
type UnimplementedRegistryServer struct {
}

func (UnimplementedRegistryServer) Register(context.Context, *RegistrationInfo) (*RegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRegistryServer) Unregister(context.Context, *RegistrationResponse) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unregister not implemented")
}
func (UnimplementedRegistryServer) Subscribe(*EndpointRequest, Registry_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedRegistryServer) mustEmbedUnimplementedRegistryServer() {}

// UnsafeRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistryServer will
// result in compilation errors.
type UnsafeRegistryServer interface {
	mustEmbedUnimplementedRegistryServer()
}

func RegisterRegistryServer(s grpc.ServiceRegistrar, srv RegistryServer) {
	s.RegisterService(&Registry_ServiceDesc, srv)
}

func _Registry_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegistrationInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.Registry/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).Register(ctx, req.(*RegistrationInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_Unregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegistrationResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).Unregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.Registry/Unregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).Unregister(ctx, req.(*RegistrationResponse))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EndpointRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RegistryServer).Subscribe(m, &registrySubscribeServer{stream})
}

type Registry_SubscribeServer interface {
	Send(*EndpointResponse) error
	grpc.ServerStream
}

type registrySubscribeServer struct {
	grpc.ServerStream
}

func (x *registrySubscribeServer) Send(m *EndpointResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Registry_ServiceDesc is the grpc.ServiceDesc for Registry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Registry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "registry.Registry",
	HandlerType: (*RegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Registry_Register_Handler,
		},
		{
			MethodName: "Unregister",
			Handler:    _Registry_Unregister_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Registry_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "registry.proto",
}