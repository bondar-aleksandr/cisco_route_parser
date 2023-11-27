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
	RouteParser_Upload_FullMethodName          = "/route_parser.routeParser/Upload"
	RouteParser_RouteLookupByIP_FullMethodName = "/route_parser.routeParser/RouteLookupByIP"
)

// RouteParserClient is the client API for RouteParser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RouteParserClient interface {
	Upload(ctx context.Context, opts ...grpc.CallOption) (RouteParser_UploadClient, error)
	RouteLookupByIP(ctx context.Context, in *RouteLookupByIPRequest, opts ...grpc.CallOption) (RouteParser_RouteLookupByIPClient, error)
}

type routeParserClient struct {
	cc grpc.ClientConnInterface
}

func NewRouteParserClient(cc grpc.ClientConnInterface) RouteParserClient {
	return &routeParserClient{cc}
}

func (c *routeParserClient) Upload(ctx context.Context, opts ...grpc.CallOption) (RouteParser_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &RouteParser_ServiceDesc.Streams[0], RouteParser_Upload_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &routeParserUploadClient{stream}
	return x, nil
}

type RouteParser_UploadClient interface {
	Send(*FileUploadRequest) error
	CloseAndRecv() (*FileUploadResponse, error)
	grpc.ClientStream
}

type routeParserUploadClient struct {
	grpc.ClientStream
}

func (x *routeParserUploadClient) Send(m *FileUploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *routeParserUploadClient) CloseAndRecv() (*FileUploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileUploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *routeParserClient) RouteLookupByIP(ctx context.Context, in *RouteLookupByIPRequest, opts ...grpc.CallOption) (RouteParser_RouteLookupByIPClient, error) {
	stream, err := c.cc.NewStream(ctx, &RouteParser_ServiceDesc.Streams[1], RouteParser_RouteLookupByIP_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &routeParserRouteLookupByIPClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RouteParser_RouteLookupByIPClient interface {
	Recv() (*RouteLookupByIPResponse, error)
	grpc.ClientStream
}

type routeParserRouteLookupByIPClient struct {
	grpc.ClientStream
}

func (x *routeParserRouteLookupByIPClient) Recv() (*RouteLookupByIPResponse, error) {
	m := new(RouteLookupByIPResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouteParserServer is the server API for RouteParser service.
// All implementations must embed UnimplementedRouteParserServer
// for forward compatibility
type RouteParserServer interface {
	Upload(RouteParser_UploadServer) error
	RouteLookupByIP(*RouteLookupByIPRequest, RouteParser_RouteLookupByIPServer) error
	mustEmbedUnimplementedRouteParserServer()
}

// UnimplementedRouteParserServer must be embedded to have forward compatible implementations.
type UnimplementedRouteParserServer struct {
}

func (UnimplementedRouteParserServer) Upload(RouteParser_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedRouteParserServer) RouteLookupByIP(*RouteLookupByIPRequest, RouteParser_RouteLookupByIPServer) error {
	return status.Errorf(codes.Unimplemented, "method RouteLookupByIP not implemented")
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

func _RouteParser_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RouteParserServer).Upload(&routeParserUploadServer{stream})
}

type RouteParser_UploadServer interface {
	SendAndClose(*FileUploadResponse) error
	Recv() (*FileUploadRequest, error)
	grpc.ServerStream
}

type routeParserUploadServer struct {
	grpc.ServerStream
}

func (x *routeParserUploadServer) SendAndClose(m *FileUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *routeParserUploadServer) Recv() (*FileUploadRequest, error) {
	m := new(FileUploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RouteParser_RouteLookupByIP_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RouteLookupByIPRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RouteParserServer).RouteLookupByIP(m, &routeParserRouteLookupByIPServer{stream})
}

type RouteParser_RouteLookupByIPServer interface {
	Send(*RouteLookupByIPResponse) error
	grpc.ServerStream
}

type routeParserRouteLookupByIPServer struct {
	grpc.ServerStream
}

func (x *routeParserRouteLookupByIPServer) Send(m *RouteLookupByIPResponse) error {
	return x.ServerStream.SendMsg(m)
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
			StreamName:    "Upload",
			Handler:       _RouteParser_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RouteLookupByIP",
			Handler:       _RouteParser_RouteLookupByIP_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "route_parser.proto",
}
