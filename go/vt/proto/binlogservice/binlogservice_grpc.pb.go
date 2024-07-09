// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.3
// source: binlogservice.proto

package binlogservice

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	binlogdata "mdibaiee/vitess/oracle/go/vt/proto/binlogdata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UpdateStreamClient is the client API for UpdateStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UpdateStreamClient interface {
	// StreamKeyRange returns the binlog transactions related to
	// the specified Keyrange.
	StreamKeyRange(ctx context.Context, in *binlogdata.StreamKeyRangeRequest, opts ...grpc.CallOption) (UpdateStream_StreamKeyRangeClient, error)
	// StreamTables returns the binlog transactions related to
	// the specified Tables.
	StreamTables(ctx context.Context, in *binlogdata.StreamTablesRequest, opts ...grpc.CallOption) (UpdateStream_StreamTablesClient, error)
}

type updateStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewUpdateStreamClient(cc grpc.ClientConnInterface) UpdateStreamClient {
	return &updateStreamClient{cc}
}

func (c *updateStreamClient) StreamKeyRange(ctx context.Context, in *binlogdata.StreamKeyRangeRequest, opts ...grpc.CallOption) (UpdateStream_StreamKeyRangeClient, error) {
	stream, err := c.cc.NewStream(ctx, &UpdateStream_ServiceDesc.Streams[0], "/binlogservice.UpdateStream/StreamKeyRange", opts...)
	if err != nil {
		return nil, err
	}
	x := &updateStreamStreamKeyRangeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UpdateStream_StreamKeyRangeClient interface {
	Recv() (*binlogdata.StreamKeyRangeResponse, error)
	grpc.ClientStream
}

type updateStreamStreamKeyRangeClient struct {
	grpc.ClientStream
}

func (x *updateStreamStreamKeyRangeClient) Recv() (*binlogdata.StreamKeyRangeResponse, error) {
	m := new(binlogdata.StreamKeyRangeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *updateStreamClient) StreamTables(ctx context.Context, in *binlogdata.StreamTablesRequest, opts ...grpc.CallOption) (UpdateStream_StreamTablesClient, error) {
	stream, err := c.cc.NewStream(ctx, &UpdateStream_ServiceDesc.Streams[1], "/binlogservice.UpdateStream/StreamTables", opts...)
	if err != nil {
		return nil, err
	}
	x := &updateStreamStreamTablesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UpdateStream_StreamTablesClient interface {
	Recv() (*binlogdata.StreamTablesResponse, error)
	grpc.ClientStream
}

type updateStreamStreamTablesClient struct {
	grpc.ClientStream
}

func (x *updateStreamStreamTablesClient) Recv() (*binlogdata.StreamTablesResponse, error) {
	m := new(binlogdata.StreamTablesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UpdateStreamServer is the server API for UpdateStream service.
// All implementations must embed UnimplementedUpdateStreamServer
// for forward compatibility
type UpdateStreamServer interface {
	// StreamKeyRange returns the binlog transactions related to
	// the specified Keyrange.
	StreamKeyRange(*binlogdata.StreamKeyRangeRequest, UpdateStream_StreamKeyRangeServer) error
	// StreamTables returns the binlog transactions related to
	// the specified Tables.
	StreamTables(*binlogdata.StreamTablesRequest, UpdateStream_StreamTablesServer) error
	mustEmbedUnimplementedUpdateStreamServer()
}

// UnimplementedUpdateStreamServer must be embedded to have forward compatible implementations.
type UnimplementedUpdateStreamServer struct {
}

func (UnimplementedUpdateStreamServer) StreamKeyRange(*binlogdata.StreamKeyRangeRequest, UpdateStream_StreamKeyRangeServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamKeyRange not implemented")
}
func (UnimplementedUpdateStreamServer) StreamTables(*binlogdata.StreamTablesRequest, UpdateStream_StreamTablesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamTables not implemented")
}
func (UnimplementedUpdateStreamServer) mustEmbedUnimplementedUpdateStreamServer() {}

// UnsafeUpdateStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UpdateStreamServer will
// result in compilation errors.
type UnsafeUpdateStreamServer interface {
	mustEmbedUnimplementedUpdateStreamServer()
}

func RegisterUpdateStreamServer(s grpc.ServiceRegistrar, srv UpdateStreamServer) {
	s.RegisterService(&UpdateStream_ServiceDesc, srv)
}

func _UpdateStream_StreamKeyRange_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(binlogdata.StreamKeyRangeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UpdateStreamServer).StreamKeyRange(m, &updateStreamStreamKeyRangeServer{stream})
}

type UpdateStream_StreamKeyRangeServer interface {
	Send(*binlogdata.StreamKeyRangeResponse) error
	grpc.ServerStream
}

type updateStreamStreamKeyRangeServer struct {
	grpc.ServerStream
}

func (x *updateStreamStreamKeyRangeServer) Send(m *binlogdata.StreamKeyRangeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _UpdateStream_StreamTables_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(binlogdata.StreamTablesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UpdateStreamServer).StreamTables(m, &updateStreamStreamTablesServer{stream})
}

type UpdateStream_StreamTablesServer interface {
	Send(*binlogdata.StreamTablesResponse) error
	grpc.ServerStream
}

type updateStreamStreamTablesServer struct {
	grpc.ServerStream
}

func (x *updateStreamStreamTablesServer) Send(m *binlogdata.StreamTablesResponse) error {
	return x.ServerStream.SendMsg(m)
}

// UpdateStream_ServiceDesc is the grpc.ServiceDesc for UpdateStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UpdateStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "binlogservice.UpdateStream",
	HandlerType: (*UpdateStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamKeyRange",
			Handler:       _UpdateStream_StreamKeyRange_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamTables",
			Handler:       _UpdateStream_StreamTables_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "binlogservice.proto",
}
