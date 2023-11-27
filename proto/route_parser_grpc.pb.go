// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.20.2
// source: route_parser.proto

package proto

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
	RouteParser_FileTransfer_FullMethodName = "/route_parser.routeParser/FileTransfer"
)

// RouteParserClient is the client API for RouteParser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RouteParserClient interface {
	FileTransfer(ctx context.Context, opts ...grpc.CallOption) (RouteParser_FileTransferClient, error)
}

type routeParserClient struct {
	cc grpc.ClientConnInterface
}

func NewRouteParserClient(cc grpc.ClientConnInterface) RouteParserClient {
	return &routeParserClient{cc}
}

func (c *routeParserClient) FileTransfer(ctx context.Context, opts ...grpc.CallOption) (RouteParser_FileTransferClient, error) {
	stream, err := c.cc.NewStream(ctx, &RouteParser_ServiceDesc.Streams[0], RouteParser_FileTransfer_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &routeParserFileTransferClient{stream}
	return x, nil
}

type RouteParser_FileTransferClient interface {
	Send(*FileUploadRequest) error
	CloseAndRecv() (*FileUploadResponse, error)
	grpc.ClientStream
}

type routeParserFileTransferClient struct {
	grpc.ClientStream
}

func (x *routeParserFileTransferClient) Send(m *FileUploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *routeParserFileTransferClient) CloseAndRecv() (*FileUploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileUploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouteParserServer is the server API for RouteParser service.
// All implementations must embed UnimplementedRouteParserServer
// for forward compatibility
type RouteParserServer interface {
	FileTransfer(RouteParser_FileTransferServer) error
	mustEmbedUnimplementedRouteParserServer()
}

// UnimplementedRouteParserServer must be embedded to have forward compatible implementations.
type UnimplementedRouteParserServer struct {
}

func (UnimplementedRouteParserServer) FileTransfer(RouteParser_FileTransferServer) error {
	return status.Errorf(codes.Unimplemented, "method FileTransfer not implemented")
}
func (UnimplementedRouteParserServer) mustEmbedUnimplementedRouteParserServer() {}

// UnsafeRouteParserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RouteParserServer will
// result in compilation errors.
type UnsafeRouteParserServer interface {
	mustEmbedUnimplementedRouteParserServer()
}

func RegisterRouteParserServer(s grpc.ServiceRegistrar, srv RouteParserServer) {
	s.RegisterService(&RouteParser_ServiceDesc, srv)
}

func _RouteParser_FileTransfer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RouteParserServer).FileTransfer(&routeParserFileTransferServer{stream})
}

type RouteParser_FileTransferServer interface {
	SendAndClose(*FileUploadResponse) error
	Recv() (*FileUploadRequest, error)
	grpc.ServerStream
}

type routeParserFileTransferServer struct {
	grpc.ServerStream
}

func (x *routeParserFileTransferServer) SendAndClose(m *FileUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *routeParserFileTransferServer) Recv() (*FileUploadRequest, error) {
	m := new(FileUploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouteParser_ServiceDesc is the grpc.ServiceDesc for RouteParser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RouteParser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "route_parser.routeParser",
	HandlerType: (*RouteParserServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FileTransfer",
			Handler:       _RouteParser_FileTransfer_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "route_parser.proto",
}