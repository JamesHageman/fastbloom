// Code generated by protoc-gen-go. DO NOT EDIT.
// source: filter.proto

package filter

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Filter struct {
	M                    uint64   `protobuf:"varint,1,opt,name=m,proto3" json:"m,omitempty"`
	K                    uint64   `protobuf:"varint,2,opt,name=k,proto3" json:"k,omitempty"`
	Data                 []uint32 `protobuf:"varint,3,rep,packed,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()    {}
func (*Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f5303cab7a20d6f, []int{0}
}

func (m *Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Filter.Unmarshal(m, b)
}
func (m *Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Filter.Marshal(b, m, deterministic)
}
func (m *Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter.Merge(m, src)
}
func (m *Filter) XXX_Size() int {
	return xxx_messageInfo_Filter.Size(m)
}
func (m *Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Filter proto.InternalMessageInfo

func (m *Filter) GetM() uint64 {
	if m != nil {
		return m.M
	}
	return 0
}

func (m *Filter) GetK() uint64 {
	if m != nil {
		return m.K
	}
	return 0
}

func (m *Filter) GetData() []uint32 {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Filter)(nil), "LockFreeFilter")
}

func init() { proto.RegisterFile("filter.proto", fileDescriptor_1f5303cab7a20d6f) }

var fileDescriptor_1f5303cab7a20d6f = []byte{
	// 88 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xcb, 0xcc, 0x29,
	0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xb2, 0xe0, 0x62, 0x73, 0x03, 0xf3, 0x85,
	0x78, 0xb8, 0x18, 0x73, 0x25, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82, 0x18, 0x73, 0x41, 0xbc, 0x6c,
	0x09, 0x26, 0x08, 0x2f, 0x5b, 0x48, 0x88, 0x8b, 0x25, 0x25, 0xb1, 0x24, 0x51, 0x82, 0x59, 0x81,
	0x59, 0x83, 0x37, 0x08, 0xcc, 0x4e, 0x62, 0x03, 0x1b, 0x60, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff,
	0xd5, 0xa0, 0xb3, 0x88, 0x50, 0x00, 0x00, 0x00,
}
