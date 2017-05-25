// Code generated by protoc-gen-gogo.
// source: server.proto
// DO NOT EDIT!

/*
Package proto_rpc is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
*/
package proto_rpc

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for API service

type APIClient interface {
}

type aPIClient struct {
	cc *grpc.ClientConn
}

func NewAPIClient(cc *grpc.ClientConn) APIClient {
	return &aPIClient{cc}
}

// Server API for API service

type APIServer interface {
}

func RegisterAPIServer(s *grpc.Server, srv APIServer) {
	s.RegisterService(&_API_serviceDesc, srv)
}

var _API_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto_rpc.API",
	HandlerType: (*APIServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "server.proto",
}

func init() { proto.RegisterFile("server.proto", fileDescriptorServer) }

var fileDescriptorServer = []byte{
	// 88 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x04, 0x53, 0xf1, 0x45, 0x05, 0xc9,
	0x52, 0xba, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0xe9, 0xf9, 0xe9,
	0xf9, 0xfa, 0x60, 0xa9, 0xa4, 0xd2, 0x34, 0x30, 0x0f, 0xcc, 0x01, 0xb3, 0x20, 0x3a, 0x8d, 0x58,
	0xb9, 0x98, 0x1d, 0x03, 0x3c, 0x93, 0xd8, 0xc0, 0x3c, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x8e, 0x60, 0x75, 0x72, 0x57, 0x00, 0x00, 0x00,
}
