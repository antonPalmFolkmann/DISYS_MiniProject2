// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ChatService

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

// ChittyChatServiceClient is the client API for ChittyChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatServiceClient interface {
	Publish(ctx context.Context, in *PublishMessageRequest, opts ...grpc.CallOption) (*PublishMessageReply, error)
	BroadCast(ctx context.Context, in *BroadCastRequest, opts ...grpc.CallOption) (*BroadCastReply, error)
}

type chittyChatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatServiceClient(cc grpc.ClientConnInterface) ChittyChatServiceClient {
	return &chittyChatServiceClient{cc}
}

func (c *chittyChatServiceClient) Publish(ctx context.Context, in *PublishMessageRequest, opts ...grpc.CallOption) (*PublishMessageReply, error) {
	out := new(PublishMessageReply)
	err := c.cc.Invoke(ctx, "/chat.ChittyChatService/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittyChatServiceClient) BroadCast(ctx context.Context, in *BroadCastRequest, opts ...grpc.CallOption) (*BroadCastReply, error) {
	out := new(BroadCastReply)
	err := c.cc.Invoke(ctx, "/chat.ChittyChatService/BroadCast", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittyChatServiceServer is the server API for ChittyChatService service.
// All implementations must embed UnimplementedChittyChatServiceServer
// for forward compatibility
type ChittyChatServiceServer interface {
	Publish(context.Context, *PublishMessageRequest) (*PublishMessageReply, error)
	BroadCast(context.Context, *BroadCastRequest) (*BroadCastReply, error)
	mustEmbedUnimplementedChittyChatServiceServer()
}

// UnimplementedChittyChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChittyChatServiceServer struct {
}

func (UnimplementedChittyChatServiceServer) Publish(context.Context, *PublishMessageRequest) (*PublishMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedChittyChatServiceServer) BroadCast(context.Context, *BroadCastRequest) (*BroadCastReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadCast not implemented")
}
func (UnimplementedChittyChatServiceServer) mustEmbedUnimplementedChittyChatServiceServer() {}

// UnsafeChittyChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServiceServer will
// result in compilation errors.
type UnsafeChittyChatServiceServer interface {
	mustEmbedUnimplementedChittyChatServiceServer()
}

func RegisterChittyChatServiceServer(s grpc.ServiceRegistrar, srv ChittyChatServiceServer) {
	s.RegisterService(&ChittyChatService_ServiceDesc, srv)
}

func _ChittyChatService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChittyChatService/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServiceServer).Publish(ctx, req.(*PublishMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChittyChatService_BroadCast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BroadCastRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServiceServer).BroadCast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChittyChatService/BroadCast",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServiceServer).BroadCast(ctx, req.(*BroadCastRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChittyChatService_ServiceDesc is the grpc.ServiceDesc for ChittyChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChittyChatService",
	HandlerType: (*ChittyChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _ChittyChatService_Publish_Handler,
		},
		{
			MethodName: "BroadCast",
			Handler:    _ChittyChatService_BroadCast_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ChatService/Chat.proto",
}
