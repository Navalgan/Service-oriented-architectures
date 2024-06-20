// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: statistic/proto/statistic.proto

package statistic_v1

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

// StatisticClient is the client API for Statistic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatisticClient interface {
	GetPostStatByID(ctx context.Context, in *PostIDRequest, opts ...grpc.CallOption) (*PostStatResponse, error)
	GetTopPosts(ctx context.Context, in *TopPostsRequest, opts ...grpc.CallOption) (*TopPostsResponse, error)
	GetTopUsers(ctx context.Context, in *TopUsersRequest, opts ...grpc.CallOption) (*TopUsersResponse, error)
}

type statisticClient struct {
	cc grpc.ClientConnInterface
}

func NewStatisticClient(cc grpc.ClientConnInterface) StatisticClient {
	return &statisticClient{cc}
}

func (c *statisticClient) GetPostStatByID(ctx context.Context, in *PostIDRequest, opts ...grpc.CallOption) (*PostStatResponse, error) {
	out := new(PostStatResponse)
	err := c.cc.Invoke(ctx, "/statistic.Statistic/GetPostStatByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) GetTopPosts(ctx context.Context, in *TopPostsRequest, opts ...grpc.CallOption) (*TopPostsResponse, error) {
	out := new(TopPostsResponse)
	err := c.cc.Invoke(ctx, "/statistic.Statistic/GetTopPosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) GetTopUsers(ctx context.Context, in *TopUsersRequest, opts ...grpc.CallOption) (*TopUsersResponse, error) {
	out := new(TopUsersResponse)
	err := c.cc.Invoke(ctx, "/statistic.Statistic/GetTopUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatisticServer is the server API for Statistic service.
// All implementations must embed UnimplementedStatisticServer
// for forward compatibility
type StatisticServer interface {
	GetPostStatByID(context.Context, *PostIDRequest) (*PostStatResponse, error)
	GetTopPosts(context.Context, *TopPostsRequest) (*TopPostsResponse, error)
	GetTopUsers(context.Context, *TopUsersRequest) (*TopUsersResponse, error)
	mustEmbedUnimplementedStatisticServer()
}

// UnimplementedStatisticServer must be embedded to have forward compatible implementations.
type UnimplementedStatisticServer struct {
}

func (UnimplementedStatisticServer) GetPostStatByID(context.Context, *PostIDRequest) (*PostStatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostStatByID not implemented")
}
func (UnimplementedStatisticServer) GetTopPosts(context.Context, *TopPostsRequest) (*TopPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopPosts not implemented")
}
func (UnimplementedStatisticServer) GetTopUsers(context.Context, *TopUsersRequest) (*TopUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopUsers not implemented")
}
func (UnimplementedStatisticServer) mustEmbedUnimplementedStatisticServer() {}

// UnsafeStatisticServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatisticServer will
// result in compilation errors.
type UnsafeStatisticServer interface {
	mustEmbedUnimplementedStatisticServer()
}

func RegisterStatisticServer(s grpc.ServiceRegistrar, srv StatisticServer) {
	s.RegisterService(&Statistic_ServiceDesc, srv)
}

func _Statistic_GetPostStatByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetPostStatByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statistic.Statistic/GetPostStatByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetPostStatByID(ctx, req.(*PostIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_GetTopPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetTopPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statistic.Statistic/GetTopPosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetTopPosts(ctx, req.(*TopPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_GetTopUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetTopUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statistic.Statistic/GetTopUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetTopUsers(ctx, req.(*TopUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Statistic_ServiceDesc is the grpc.ServiceDesc for Statistic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Statistic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "statistic.Statistic",
	HandlerType: (*StatisticServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPostStatByID",
			Handler:    _Statistic_GetPostStatByID_Handler,
		},
		{
			MethodName: "GetTopPosts",
			Handler:    _Statistic_GetTopPosts_Handler,
		},
		{
			MethodName: "GetTopUsers",
			Handler:    _Statistic_GetTopUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statistic/proto/statistic.proto",
}
