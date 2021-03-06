// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chatMessaging.proto

package nan0chat

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type User struct {
	UserId               int64    `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	UserName             string   `protobuf:"bytes,2,opt,name=userName,proto3" json:"userName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_chatMessaging_94440c24f04fc7f9, []int{0}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *User) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

type ChatMessage struct {
	UserId               int64    `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	MessageId            int64    `protobuf:"varint,4,opt,name=messageId,proto3" json:"messageId,omitempty"`
	Time                 int64    `protobuf:"varint,5,opt,name=time,proto3" json:"time,omitempty"`
	Message              string   `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChatMessage) Reset()         { *m = ChatMessage{} }
func (m *ChatMessage) String() string { return proto.CompactTextString(m) }
func (*ChatMessage) ProtoMessage()    {}
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_chatMessaging_94440c24f04fc7f9, []int{1}
}
func (m *ChatMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChatMessage.Unmarshal(m, b)
}
func (m *ChatMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChatMessage.Marshal(b, m, deterministic)
}
func (dst *ChatMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChatMessage.Merge(dst, src)
}
func (m *ChatMessage) XXX_Size() int {
	return xxx_messageInfo_ChatMessage.Size(m)
}
func (m *ChatMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ChatMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ChatMessage proto.InternalMessageInfo

func (m *ChatMessage) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *ChatMessage) GetMessageId() int64 {
	if m != nil {
		return m.MessageId
	}
	return 0
}

func (m *ChatMessage) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *ChatMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "nan0chat.User")
	proto.RegisterType((*ChatMessage)(nil), "nan0chat.ChatMessage")
}

func init() { proto.RegisterFile("chatMessaging.proto", fileDescriptor_chatMessaging_94440c24f04fc7f9) }

var fileDescriptor_chatMessaging_94440c24f04fc7f9 = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0xce, 0x48, 0x2c,
	0xf1, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0xcf, 0xcc, 0x4b, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0xc8, 0x4b, 0xcc, 0x33, 0x00, 0x49, 0x28, 0x59, 0x71, 0xb1, 0x84, 0x16, 0xa7, 0x16, 0x09,
	0x89, 0x71, 0xb1, 0x95, 0x16, 0xa7, 0x16, 0x79, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x30, 0x07,
	0x41, 0x79, 0x42, 0x52, 0x5c, 0x1c, 0x20, 0x96, 0x5f, 0x62, 0x6e, 0xaa, 0x04, 0x93, 0x02, 0xa3,
	0x06, 0x67, 0x10, 0x9c, 0xaf, 0x54, 0xc8, 0xc5, 0xed, 0x0c, 0x37, 0x3c, 0x15, 0xc9, 0x08, 0x66,
	0x14, 0x23, 0x64, 0xb8, 0x38, 0x73, 0x21, 0x4a, 0x3c, 0x53, 0x24, 0x58, 0xc0, 0x52, 0x08, 0x01,
	0x21, 0x21, 0x2e, 0x96, 0x92, 0xcc, 0xdc, 0x54, 0x09, 0x56, 0xb0, 0x04, 0x98, 0x2d, 0x24, 0xc1,
	0xc5, 0x0e, 0x55, 0x20, 0xc1, 0x06, 0xb6, 0x13, 0xc6, 0x75, 0xe2, 0x8a, 0x82, 0x3b, 0x3d, 0x89,
	0x0d, 0xec, 0x17, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc7, 0xfc, 0x61, 0xe0, 0xe2, 0x00,
	0x00, 0x00,
}
