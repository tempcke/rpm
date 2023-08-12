// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: rpm.proto

package pb

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

// RPMClient is the client API for RPM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RPMClient interface {
	StoreProperty(ctx context.Context, in *StorePropertyReq, opts ...grpc.CallOption) (*StorePropertyRes, error)
	GetProperty(ctx context.Context, in *GetPropertyReq, opts ...grpc.CallOption) (*GetPropertyRes, error)
	RemoveProperty(ctx context.Context, in *RemovePropertyReq, opts ...grpc.CallOption) (*RemovePropertyRes, error)
	ListProperties(ctx context.Context, in *ListPropertiesReq, opts ...grpc.CallOption) (RPM_ListPropertiesClient, error)
	StoreTenant(ctx context.Context, in *StoreTenantReq, opts ...grpc.CallOption) (*StoreTenantRes, error)
	GetTenant(ctx context.Context, in *GetTenantReq, opts ...grpc.CallOption) (*GetTenantRes, error)
	ListTenants(ctx context.Context, in *ListTenantsReq, opts ...grpc.CallOption) (RPM_ListTenantsClient, error)
}

type rPMClient struct {
	cc grpc.ClientConnInterface
}

func NewRPMClient(cc grpc.ClientConnInterface) RPMClient {
	return &rPMClient{cc}
}

func (c *rPMClient) StoreProperty(ctx context.Context, in *StorePropertyReq, opts ...grpc.CallOption) (*StorePropertyRes, error) {
	out := new(StorePropertyRes)
	err := c.cc.Invoke(ctx, "/rpmpb.RPM/StoreProperty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPMClient) GetProperty(ctx context.Context, in *GetPropertyReq, opts ...grpc.CallOption) (*GetPropertyRes, error) {
	out := new(GetPropertyRes)
	err := c.cc.Invoke(ctx, "/rpmpb.RPM/GetProperty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPMClient) RemoveProperty(ctx context.Context, in *RemovePropertyReq, opts ...grpc.CallOption) (*RemovePropertyRes, error) {
	out := new(RemovePropertyRes)
	err := c.cc.Invoke(ctx, "/rpmpb.RPM/RemoveProperty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPMClient) ListProperties(ctx context.Context, in *ListPropertiesReq, opts ...grpc.CallOption) (RPM_ListPropertiesClient, error) {
	stream, err := c.cc.NewStream(ctx, &RPM_ServiceDesc.Streams[0], "/rpmpb.RPM/ListProperties", opts...)
	if err != nil {
		return nil, err
	}
	x := &rPMListPropertiesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RPM_ListPropertiesClient interface {
	Recv() (*Property, error)
	grpc.ClientStream
}

type rPMListPropertiesClient struct {
	grpc.ClientStream
}

func (x *rPMListPropertiesClient) Recv() (*Property, error) {
	m := new(Property)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rPMClient) StoreTenant(ctx context.Context, in *StoreTenantReq, opts ...grpc.CallOption) (*StoreTenantRes, error) {
	out := new(StoreTenantRes)
	err := c.cc.Invoke(ctx, "/rpmpb.RPM/StoreTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPMClient) GetTenant(ctx context.Context, in *GetTenantReq, opts ...grpc.CallOption) (*GetTenantRes, error) {
	out := new(GetTenantRes)
	err := c.cc.Invoke(ctx, "/rpmpb.RPM/GetTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPMClient) ListTenants(ctx context.Context, in *ListTenantsReq, opts ...grpc.CallOption) (RPM_ListTenantsClient, error) {
	stream, err := c.cc.NewStream(ctx, &RPM_ServiceDesc.Streams[1], "/rpmpb.RPM/ListTenants", opts...)
	if err != nil {
		return nil, err
	}
	x := &rPMListTenantsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RPM_ListTenantsClient interface {
	Recv() (*Tenant, error)
	grpc.ClientStream
}

type rPMListTenantsClient struct {
	grpc.ClientStream
}

func (x *rPMListTenantsClient) Recv() (*Tenant, error) {
	m := new(Tenant)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RPMServer is the server API for RPM service.
// All implementations must embed UnimplementedRPMServer
// for forward compatibility
type RPMServer interface {
	StoreProperty(context.Context, *StorePropertyReq) (*StorePropertyRes, error)
	GetProperty(context.Context, *GetPropertyReq) (*GetPropertyRes, error)
	RemoveProperty(context.Context, *RemovePropertyReq) (*RemovePropertyRes, error)
	ListProperties(*ListPropertiesReq, RPM_ListPropertiesServer) error
	StoreTenant(context.Context, *StoreTenantReq) (*StoreTenantRes, error)
	GetTenant(context.Context, *GetTenantReq) (*GetTenantRes, error)
	ListTenants(*ListTenantsReq, RPM_ListTenantsServer) error
	mustEmbedUnimplementedRPMServer()
}

// UnimplementedRPMServer must be embedded to have forward compatible implementations.
type UnimplementedRPMServer struct {
}

func (UnimplementedRPMServer) StoreProperty(context.Context, *StorePropertyReq) (*StorePropertyRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreProperty not implemented")
}
func (UnimplementedRPMServer) GetProperty(context.Context, *GetPropertyReq) (*GetPropertyRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProperty not implemented")
}
func (UnimplementedRPMServer) RemoveProperty(context.Context, *RemovePropertyReq) (*RemovePropertyRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveProperty not implemented")
}
func (UnimplementedRPMServer) ListProperties(*ListPropertiesReq, RPM_ListPropertiesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListProperties not implemented")
}
func (UnimplementedRPMServer) StoreTenant(context.Context, *StoreTenantReq) (*StoreTenantRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreTenant not implemented")
}
func (UnimplementedRPMServer) GetTenant(context.Context, *GetTenantReq) (*GetTenantRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTenant not implemented")
}
func (UnimplementedRPMServer) ListTenants(*ListTenantsReq, RPM_ListTenantsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListTenants not implemented")
}
func (UnimplementedRPMServer) mustEmbedUnimplementedRPMServer() {}

// UnsafeRPMServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RPMServer will
// result in compilation errors.
type UnsafeRPMServer interface {
	mustEmbedUnimplementedRPMServer()
}

func RegisterRPMServer(s grpc.ServiceRegistrar, srv RPMServer) {
	s.RegisterService(&RPM_ServiceDesc, srv)
}

func _RPM_StoreProperty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StorePropertyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPMServer).StoreProperty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpmpb.RPM/StoreProperty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPMServer).StoreProperty(ctx, req.(*StorePropertyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPM_GetProperty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPropertyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPMServer).GetProperty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpmpb.RPM/GetProperty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPMServer).GetProperty(ctx, req.(*GetPropertyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPM_RemoveProperty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemovePropertyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPMServer).RemoveProperty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpmpb.RPM/RemoveProperty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPMServer).RemoveProperty(ctx, req.(*RemovePropertyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPM_ListProperties_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListPropertiesReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RPMServer).ListProperties(m, &rPMListPropertiesServer{stream})
}

type RPM_ListPropertiesServer interface {
	Send(*Property) error
	grpc.ServerStream
}

type rPMListPropertiesServer struct {
	grpc.ServerStream
}

func (x *rPMListPropertiesServer) Send(m *Property) error {
	return x.ServerStream.SendMsg(m)
}

func _RPM_StoreTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreTenantReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPMServer).StoreTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpmpb.RPM/StoreTenant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPMServer).StoreTenant(ctx, req.(*StoreTenantReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPM_GetTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTenantReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPMServer).GetTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpmpb.RPM/GetTenant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPMServer).GetTenant(ctx, req.(*GetTenantReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPM_ListTenants_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListTenantsReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RPMServer).ListTenants(m, &rPMListTenantsServer{stream})
}

type RPM_ListTenantsServer interface {
	Send(*Tenant) error
	grpc.ServerStream
}

type rPMListTenantsServer struct {
	grpc.ServerStream
}

func (x *rPMListTenantsServer) Send(m *Tenant) error {
	return x.ServerStream.SendMsg(m)
}

// RPM_ServiceDesc is the grpc.ServiceDesc for RPM service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RPM_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpmpb.RPM",
	HandlerType: (*RPMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreProperty",
			Handler:    _RPM_StoreProperty_Handler,
		},
		{
			MethodName: "GetProperty",
			Handler:    _RPM_GetProperty_Handler,
		},
		{
			MethodName: "RemoveProperty",
			Handler:    _RPM_RemoveProperty_Handler,
		},
		{
			MethodName: "StoreTenant",
			Handler:    _RPM_StoreTenant_Handler,
		},
		{
			MethodName: "GetTenant",
			Handler:    _RPM_GetTenant_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListProperties",
			Handler:       _RPM_ListProperties_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListTenants",
			Handler:       _RPM_ListTenants_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpm.proto",
}
