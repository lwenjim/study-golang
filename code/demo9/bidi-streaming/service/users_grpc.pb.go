// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.1
// source: users.proto

package service

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
	Users_GetHep_FullMethodName = "/Users/GetHep"
)

// UsersClient is the client API for Users service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersClient interface {
	GetHep(ctx context.Context, opts ...grpc.CallOption) (Users_GetHepClient, error)
}

type usersClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersClient(cc grpc.ClientConnInterface) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) GetHep(ctx context.Context, opts ...grpc.CallOption) (Users_GetHepClient, error) {
	stream, err := c.cc.NewStream(ctx, &Users_ServiceDesc.Streams[0], Users_GetHep_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &usersGetHepClient{stream}
	return x, nil
}

type Users_GetHepClient interface {
	Send(*UserHelpRequest) error
	Recv() (*UserHelpReply, error)
	grpc.ClientStream
}

type usersGetHepClient struct {
	grpc.ClientStream
}

func (x *usersGetHepClient) Send(m *UserHelpRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *usersGetHepClient) Recv() (*UserHelpReply, error) {
	m := new(UserHelpReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UsersServer is the server API for Users service.
// All implementations must embed UnimplementedUsersServer
// for forward compatibility
type UsersServer interface {
	GetHep(Users_GetHepServer) error
	mustEmbedUnimplementedUsersServer()
}

// UnimplementedUsersServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServer struct {
}

func (UnimplementedUsersServer) GetHep(Users_GetHepServer) error {
	return status.Errorf(codes.Unimplemented, "method GetHep not implemented")
}
func (UnimplementedUsersServer) mustEmbedUnimplementedUsersServer() {}

// UnsafeUsersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServer will
// result in compilation errors.
type UnsafeUsersServer interface {
	mustEmbedUnimplementedUsersServer()
}

func RegisterUsersServer(s grpc.ServiceRegistrar, srv UsersServer) {
	s.RegisterService(&Users_ServiceDesc, srv)
}

func _Users_GetHep_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UsersServer).GetHep(&usersGetHepServer{stream})
}

type Users_GetHepServer interface {
	Send(*UserHelpReply) error
	Recv() (*UserHelpRequest, error)
	grpc.ServerStream
}

type usersGetHepServer struct {
	grpc.ServerStream
}

func (x *usersGetHepServer) Send(m *UserHelpReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *usersGetHepServer) Recv() (*UserHelpRequest, error) {
	m := new(UserHelpRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Users_ServiceDesc is the grpc.ServiceDesc for Users service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Users_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Users",
	HandlerType: (*UsersServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetHep",
			Handler:       _Users_GetHep_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "users.proto",
}
