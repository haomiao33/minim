// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: proto/OfflineService.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	OfflineService_Push_FullMethodName = "/myservice.OfflineService/Push"
)

// OfflineServiceClient is the client API for OfflineService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OfflineServiceClient interface {
	// 推送
	Push(ctx context.Context, in *OfflinePushRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type offlineServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOfflineServiceClient(cc grpc.ClientConnInterface) OfflineServiceClient {
	return &offlineServiceClient{cc}
}

func (c *offlineServiceClient) Push(ctx context.Context, in *OfflinePushRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, OfflineService_Push_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OfflineServiceServer is the server API for OfflineService service.
// All implementations must embed UnimplementedOfflineServiceServer
// for forward compatibility.
type OfflineServiceServer interface {
	// 推送
	Push(context.Context, *OfflinePushRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedOfflineServiceServer()
}

// UnimplementedOfflineServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOfflineServiceServer struct{}

func (UnimplementedOfflineServiceServer) Push(context.Context, *OfflinePushRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedOfflineServiceServer) mustEmbedUnimplementedOfflineServiceServer() {}
func (UnimplementedOfflineServiceServer) testEmbeddedByValue()                        {}

// UnsafeOfflineServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OfflineServiceServer will
// result in compilation errors.
type UnsafeOfflineServiceServer interface {
	mustEmbedUnimplementedOfflineServiceServer()
}

func RegisterOfflineServiceServer(s grpc.ServiceRegistrar, srv OfflineServiceServer) {
	// If the following call pancis, it indicates UnimplementedOfflineServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OfflineService_ServiceDesc, srv)
}

func _OfflineService_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OfflinePushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OfflineServiceServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OfflineService_Push_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OfflineServiceServer).Push(ctx, req.(*OfflinePushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OfflineService_ServiceDesc is the grpc.ServiceDesc for OfflineService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OfflineService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "myservice.OfflineService",
	HandlerType: (*OfflineServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _OfflineService_Push_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/OfflineService.proto",
}
