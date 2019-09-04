// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Plugin ...
//
// These are all info for a plugin.
type Plugin struct {
	// name ...
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// identifier ...
	Identifier string `protobuf:"bytes,2,opt,name=identifier,proto3" json:"identifier,omitempty"`
	// version ...
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	// path ...
	Path                 string   `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Plugin) Reset()         { *m = Plugin{} }
func (m *Plugin) String() string { return proto.CompactTextString(m) }
func (*Plugin) ProtoMessage()    {}
func (*Plugin) Descriptor() ([]byte, []int) {
	return fileDescriptor_server_d2010764b3ce5b41, []int{0}
}
func (m *Plugin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Plugin.Unmarshal(m, b)
}
func (m *Plugin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Plugin.Marshal(b, m, deterministic)
}
func (dst *Plugin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plugin.Merge(dst, src)
}
func (m *Plugin) XXX_Size() int {
	return xxx_messageInfo_Plugin.Size(m)
}
func (m *Plugin) XXX_DiscardUnknown() {
	xxx_messageInfo_Plugin.DiscardUnknown(m)
}

var xxx_messageInfo_Plugin proto.InternalMessageInfo

func (m *Plugin) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Plugin) GetIdentifier() string {
	if m != nil {
		return m.Identifier
	}
	return ""
}

func (m *Plugin) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Plugin) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

// ListDomains ...
//
//
// Listing the plugins in the server.
type ListDomains struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListDomains) Reset()         { *m = ListDomains{} }
func (m *ListDomains) String() string { return proto.CompactTextString(m) }
func (*ListDomains) ProtoMessage()    {}
func (*ListDomains) Descriptor() ([]byte, []int) {
	return fileDescriptor_server_d2010764b3ce5b41, []int{1}
}
func (m *ListDomains) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListDomains.Unmarshal(m, b)
}
func (m *ListDomains) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListDomains.Marshal(b, m, deterministic)
}
func (dst *ListDomains) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListDomains.Merge(dst, src)
}
func (m *ListDomains) XXX_Size() int {
	return xxx_messageInfo_ListDomains.Size(m)
}
func (m *ListDomains) XXX_DiscardUnknown() {
	xxx_messageInfo_ListDomains.DiscardUnknown(m)
}

var xxx_messageInfo_ListDomains proto.InternalMessageInfo

// Request ...
type ListDomains_Request struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListDomains_Request) Reset()         { *m = ListDomains_Request{} }
func (m *ListDomains_Request) String() string { return proto.CompactTextString(m) }
func (*ListDomains_Request) ProtoMessage()    {}
func (*ListDomains_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_server_d2010764b3ce5b41, []int{1, 0}
}
func (m *ListDomains_Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListDomains_Request.Unmarshal(m, b)
}
func (m *ListDomains_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListDomains_Request.Marshal(b, m, deterministic)
}
func (dst *ListDomains_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListDomains_Request.Merge(dst, src)
}
func (m *ListDomains_Request) XXX_Size() int {
	return xxx_messageInfo_ListDomains_Request.Size(m)
}
func (m *ListDomains_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_ListDomains_Request.DiscardUnknown(m)
}

var xxx_messageInfo_ListDomains_Request proto.InternalMessageInfo

// Response ...
type ListDomains_Response struct {
	Plugins              []*Plugin `protobuf:"bytes,1,rep,name=plugins,proto3" json:"plugins,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ListDomains_Response) Reset()         { *m = ListDomains_Response{} }
func (m *ListDomains_Response) String() string { return proto.CompactTextString(m) }
func (*ListDomains_Response) ProtoMessage()    {}
func (*ListDomains_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_server_d2010764b3ce5b41, []int{1, 1}
}
func (m *ListDomains_Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListDomains_Response.Unmarshal(m, b)
}
func (m *ListDomains_Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListDomains_Response.Marshal(b, m, deterministic)
}
func (dst *ListDomains_Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListDomains_Response.Merge(dst, src)
}
func (m *ListDomains_Response) XXX_Size() int {
	return xxx_messageInfo_ListDomains_Response.Size(m)
}
func (m *ListDomains_Response) XXX_DiscardUnknown() {
	xxx_messageInfo_ListDomains_Response.DiscardUnknown(m)
}

var xxx_messageInfo_ListDomains_Response proto.InternalMessageInfo

func (m *ListDomains_Response) GetPlugins() []*Plugin {
	if m != nil {
		return m.Plugins
	}
	return nil
}

func init() {
	proto.RegisterType((*Plugin)(nil), "proto.Plugin")
	proto.RegisterType((*ListDomains)(nil), "proto.ListDomains")
	proto.RegisterType((*ListDomains_Request)(nil), "proto.ListDomains.Request")
	proto.RegisterType((*ListDomains_Response)(nil), "proto.ListDomains.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AutobotClient is the client API for Autobot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AutobotClient interface {
	// ListPlugins is listing plugins.
	ListPlugins(ctx context.Context, in *ListDomains_Request, opts ...grpc.CallOption) (*ListDomains_Response, error)
}

type autobotClient struct {
	cc *grpc.ClientConn
}

func NewAutobotClient(cc *grpc.ClientConn) AutobotClient {
	return &autobotClient{cc}
}

func (c *autobotClient) ListPlugins(ctx context.Context, in *ListDomains_Request, opts ...grpc.CallOption) (*ListDomains_Response, error) {
	out := new(ListDomains_Response)
	err := c.cc.Invoke(ctx, "/proto.Autobot/ListPlugins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AutobotServer is the server API for Autobot service.
type AutobotServer interface {
	// ListPlugins is listing plugins.
	ListPlugins(context.Context, *ListDomains_Request) (*ListDomains_Response, error)
}

func RegisterAutobotServer(s *grpc.Server, srv AutobotServer) {
	s.RegisterService(&_Autobot_serviceDesc, srv)
}

func _Autobot_ListPlugins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDomains_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AutobotServer).ListPlugins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Autobot/ListPlugins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AutobotServer).ListPlugins(ctx, req.(*ListDomains_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Autobot_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Autobot",
	HandlerType: (*AutobotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPlugins",
			Handler:    _Autobot_ListPlugins_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}

func init() { proto.RegisterFile("server.proto", fileDescriptor_server_d2010764b3ce5b41) }

var fileDescriptor_server_d2010764b3ce5b41 = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8e, 0x41, 0x4f, 0x84, 0x30,
	0x10, 0x85, 0x83, 0xbb, 0x6e, 0x65, 0xd0, 0xcb, 0x9c, 0x1a, 0x4c, 0x0c, 0xe1, 0x22, 0x27, 0x0e,
	0xf0, 0x0b, 0x4c, 0x8c, 0x27, 0x4d, 0xb4, 0xff, 0x00, 0xe2, 0xa8, 0x35, 0xd2, 0xd6, 0x4e, 0xe1,
	0xf7, 0x1b, 0x5a, 0x49, 0x38, 0xec, 0xa9, 0xaf, 0xf3, 0xcd, 0xbc, 0xf7, 0xe0, 0x9a, 0xc9, 0x2f,
	0xe4, 0x5b, 0xe7, 0x6d, 0xb0, 0x78, 0x19, 0x9f, 0xfa, 0x1b, 0x4e, 0xaf, 0x3f, 0xf3, 0xa7, 0x36,
	0x88, 0x70, 0x34, 0xc3, 0x44, 0x32, 0xab, 0xb2, 0x26, 0x57, 0x51, 0xe3, 0x1d, 0x80, 0x7e, 0x27,
	0x13, 0xf4, 0x87, 0x26, 0x2f, 0x2f, 0x22, 0xd9, 0x4d, 0x50, 0x82, 0x58, 0xc8, 0xb3, 0xb6, 0x46,
	0x1e, 0x22, 0xdc, 0xbe, 0xab, 0x9b, 0x1b, 0xc2, 0x97, 0x3c, 0x26, 0xb7, 0x55, 0xd7, 0x2f, 0x50,
	0x3c, 0x6b, 0x0e, 0x8f, 0x76, 0x1a, 0xb4, 0xe1, 0x32, 0x07, 0xa1, 0xe8, 0x77, 0x26, 0x0e, 0x65,
	0x0f, 0x57, 0x8a, 0xd8, 0x59, 0xc3, 0x84, 0xf7, 0x20, 0x5c, 0x6c, 0xc4, 0x32, 0xab, 0x0e, 0x4d,
	0xd1, 0xdd, 0xa4, 0xc6, 0x6d, 0xea, 0xa9, 0x36, 0xda, 0xbd, 0x81, 0x78, 0x98, 0x83, 0x1d, 0x6d,
	0xc0, 0xa7, 0xe4, 0x9c, 0x36, 0x18, 0xcb, 0xff, 0x8b, 0x5d, 0x5a, 0xbb, 0x45, 0xdd, 0x9e, 0x65,
	0x29, 0x7b, 0x3c, 0x45, 0xd6, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0xfe, 0xaf, 0x9b, 0x62, 0x2b,
	0x01, 0x00, 0x00,
}
