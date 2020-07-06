// Code generated by protoc-gen-go. DO NOT EDIT.
// source: authenticator.proto

package approzium_authenticator_protos

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

type ClientLanguage int32

const (
	ClientLanguage_LANGUAGE_NOT_PROVIDED ClientLanguage = 0
	ClientLanguage_PYTHON                ClientLanguage = 1
	ClientLanguage_GO                    ClientLanguage = 2
)

var ClientLanguage_name = map[int32]string{
	0: "LANGUAGE_NOT_PROVIDED",
	1: "PYTHON",
	2: "GO",
}

var ClientLanguage_value = map[string]int32{
	"LANGUAGE_NOT_PROVIDED": 0,
	"PYTHON":                1,
	"GO":                    2,
}

func (x ClientLanguage) String() string {
	return proto.EnumName(ClientLanguage_name, int32(x))
}

func (ClientLanguage) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e86ec39f7c35dea3, []int{0}
}

type AWSIdentity struct {
	SignedGetCallerIdentity string   `protobuf:"bytes,1,opt,name=signed_get_caller_identity,json=signedGetCallerIdentity,proto3" json:"signed_get_caller_identity,omitempty"`
	ClaimedIamArn           string   `protobuf:"bytes,2,opt,name=claimed_iam_arn,json=claimedIamArn,proto3" json:"claimed_iam_arn,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *AWSIdentity) Reset()         { *m = AWSIdentity{} }
func (m *AWSIdentity) String() string { return proto.CompactTextString(m) }
func (*AWSIdentity) ProtoMessage()    {}
func (*AWSIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_e86ec39f7c35dea3, []int{0}
}

func (m *AWSIdentity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AWSIdentity.Unmarshal(m, b)
}
func (m *AWSIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AWSIdentity.Marshal(b, m, deterministic)
}
func (m *AWSIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AWSIdentity.Merge(m, src)
}
func (m *AWSIdentity) XXX_Size() int {
	return xxx_messageInfo_AWSIdentity.Size(m)
}
func (m *AWSIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_AWSIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_AWSIdentity proto.InternalMessageInfo

func (m *AWSIdentity) GetSignedGetCallerIdentity() string {
	if m != nil {
		return m.SignedGetCallerIdentity
	}
	return ""
}

func (m *AWSIdentity) GetClaimedIamArn() string {
	if m != nil {
		return m.ClaimedIamArn
	}
	return ""
}

type PasswordRequest struct {
	ClientLanguage       ClientLanguage `protobuf:"varint,1,opt,name=client_language,json=clientLanguage,proto3,enum=approzium.authenticator.protos.ClientLanguage" json:"client_language,omitempty"`
	Dbhost               string         `protobuf:"bytes,2,opt,name=dbhost,proto3" json:"dbhost,omitempty"`
	Dbport               string         `protobuf:"bytes,3,opt,name=dbport,proto3" json:"dbport,omitempty"`
	Dbuser               string         `protobuf:"bytes,4,opt,name=dbuser,proto3" json:"dbuser,omitempty"`
	Aws                  *AWSIdentity   `protobuf:"bytes,5,opt,name=aws,proto3" json:"aws,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PasswordRequest) Reset()         { *m = PasswordRequest{} }
func (m *PasswordRequest) String() string { return proto.CompactTextString(m) }
func (*PasswordRequest) ProtoMessage()    {}
func (*PasswordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e86ec39f7c35dea3, []int{1}
}

func (m *PasswordRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PasswordRequest.Unmarshal(m, b)
}
func (m *PasswordRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PasswordRequest.Marshal(b, m, deterministic)
}
func (m *PasswordRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PasswordRequest.Merge(m, src)
}
func (m *PasswordRequest) XXX_Size() int {
	return xxx_messageInfo_PasswordRequest.Size(m)
}
func (m *PasswordRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PasswordRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PasswordRequest proto.InternalMessageInfo

func (m *PasswordRequest) GetClientLanguage() ClientLanguage {
	if m != nil {
		return m.ClientLanguage
	}
	return ClientLanguage_LANGUAGE_NOT_PROVIDED
}

func (m *PasswordRequest) GetDbhost() string {
	if m != nil {
		return m.Dbhost
	}
	return ""
}

func (m *PasswordRequest) GetDbport() string {
	if m != nil {
		return m.Dbport
	}
	return ""
}

func (m *PasswordRequest) GetDbuser() string {
	if m != nil {
		return m.Dbuser
	}
	return ""
}

func (m *PasswordRequest) GetAws() *AWSIdentity {
	if m != nil {
		return m.Aws
	}
	return nil
}

type PGMD5HashRequest struct {
	PwdRequest           *PasswordRequest `protobuf:"bytes,1,opt,name=pwd_request,json=pwdRequest,proto3" json:"pwd_request,omitempty"`
	Salt                 []byte           `protobuf:"bytes,2,opt,name=salt,proto3" json:"salt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *PGMD5HashRequest) Reset()         { *m = PGMD5HashRequest{} }
func (m *PGMD5HashRequest) String() string { return proto.CompactTextString(m) }
func (*PGMD5HashRequest) ProtoMessage()    {}
func (*PGMD5HashRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e86ec39f7c35dea3, []int{2}
}

func (m *PGMD5HashRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PGMD5HashRequest.Unmarshal(m, b)
}
func (m *PGMD5HashRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PGMD5HashRequest.Marshal(b, m, deterministic)
}
func (m *PGMD5HashRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PGMD5HashRequest.Merge(m, src)
}
func (m *PGMD5HashRequest) XXX_Size() int {
	return xxx_messageInfo_PGMD5HashRequest.Size(m)
}
func (m *PGMD5HashRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PGMD5HashRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PGMD5HashRequest proto.InternalMessageInfo

func (m *PGMD5HashRequest) GetPwdRequest() *PasswordRequest {
	if m != nil {
		return m.PwdRequest
	}
	return nil
}

func (m *PGMD5HashRequest) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

type PGSHA256HashRequest struct {
	PwdRequest           *PasswordRequest `protobuf:"bytes,1,opt,name=pwd_request,json=pwdRequest,proto3" json:"pwd_request,omitempty"`
	Salt                 string           `protobuf:"bytes,2,opt,name=salt,proto3" json:"salt,omitempty"`
	Iterations           uint32           `protobuf:"varint,3,opt,name=iterations,proto3" json:"iterations,omitempty"`
	AuthenticationMsg    string           `protobuf:"bytes,4,opt,name=authentication_msg,json=authenticationMsg,proto3" json:"authentication_msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *PGSHA256HashRequest) Reset()         { *m = PGSHA256HashRequest{} }
func (m *PGSHA256HashRequest) String() string { return proto.CompactTextString(m) }
func (*PGSHA256HashRequest) ProtoMessage()    {}
func (*PGSHA256HashRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e86ec39f7c35dea3, []int{3}
}

func (m *PGSHA256HashRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PGSHA256HashRequest.Unmarshal(m, b)
}
func (m *PGSHA256HashRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PGSHA256HashRequest.Marshal(b, m, deterministic)
}
func (m *PGSHA256HashRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PGSHA256HashRequest.Merge(m, src)
}
func (m *PGSHA256HashRequest) XXX_Size() int {
	return xxx_messageInfo_PGSHA256HashRequest.Size(m)
}
func (m *PGSHA256HashRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PGSHA256HashRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PGSHA256HashRequest proto.InternalMessageInfo

func (m *PGSHA256HashRequest) GetPwdRequest() *PasswordRequest {
	if m != nil {
		return m.PwdRequest
	}
	return nil
}

func (m *PGSHA256HashRequest) GetSalt() string {
	if m != nil {
		return m.Salt
	}
	return ""
}

func (m *PGSHA256HashRequest) GetIterations() uint32 {
	if m != nil {
		return m.Iterations
	}
	return 0
}

func (m *PGSHA256HashRequest) GetAuthenticationMsg() string {
	if m != nil {
		return m.AuthenticationMsg
	}
	return ""
}

type PGMD5Response struct {
	Hash                 string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Requestid            string   `protobuf:"bytes,2,opt,name=requestid,proto3" json:"requestid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PGMD5Response) Reset()         { *m = PGMD5Response{} }
func (m *PGMD5Response) String() string { return proto.CompactTextString(m) }
func (*PGMD5Response) ProtoMessage()    {}
func (*PGMD5Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_e86ec39f7c35dea3, []int{4}
}

func (m *PGMD5Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PGMD5Response.Unmarshal(m, b)
}
func (m *PGMD5Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PGMD5Response.Marshal(b, m, deterministic)
}
func (m *PGMD5Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PGMD5Response.Merge(m, src)
}
func (m *PGMD5Response) XXX_Size() int {
	return xxx_messageInfo_PGMD5Response.Size(m)
}
func (m *PGMD5Response) XXX_DiscardUnknown() {
	xxx_messageInfo_PGMD5Response.DiscardUnknown(m)
}

var xxx_messageInfo_PGMD5Response proto.InternalMessageInfo

func (m *PGMD5Response) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *PGMD5Response) GetRequestid() string {
	if m != nil {
		return m.Requestid
	}
	return ""
}

type PGSHA256Response struct {
	Cproof               string   `protobuf:"bytes,1,opt,name=cproof,proto3" json:"cproof,omitempty"`
	Sproof               string   `protobuf:"bytes,2,opt,name=sproof,proto3" json:"sproof,omitempty"`
	Requestid            string   `protobuf:"bytes,3,opt,name=requestid,proto3" json:"requestid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PGSHA256Response) Reset()         { *m = PGSHA256Response{} }
func (m *PGSHA256Response) String() string { return proto.CompactTextString(m) }
func (*PGSHA256Response) ProtoMessage()    {}
func (*PGSHA256Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_e86ec39f7c35dea3, []int{5}
}

func (m *PGSHA256Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PGSHA256Response.Unmarshal(m, b)
}
func (m *PGSHA256Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PGSHA256Response.Marshal(b, m, deterministic)
}
func (m *PGSHA256Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PGSHA256Response.Merge(m, src)
}
func (m *PGSHA256Response) XXX_Size() int {
	return xxx_messageInfo_PGSHA256Response.Size(m)
}
func (m *PGSHA256Response) XXX_DiscardUnknown() {
	xxx_messageInfo_PGSHA256Response.DiscardUnknown(m)
}

var xxx_messageInfo_PGSHA256Response proto.InternalMessageInfo

func (m *PGSHA256Response) GetCproof() string {
	if m != nil {
		return m.Cproof
	}
	return ""
}

func (m *PGSHA256Response) GetSproof() string {
	if m != nil {
		return m.Sproof
	}
	return ""
}

func (m *PGSHA256Response) GetRequestid() string {
	if m != nil {
		return m.Requestid
	}
	return ""
}

func init() {
	proto.RegisterEnum("approzium.authenticator.protos.ClientLanguage", ClientLanguage_name, ClientLanguage_value)
	proto.RegisterType((*AWSIdentity)(nil), "approzium.authenticator.protos.AWSIdentity")
	proto.RegisterType((*PasswordRequest)(nil), "approzium.authenticator.protos.PasswordRequest")
	proto.RegisterType((*PGMD5HashRequest)(nil), "approzium.authenticator.protos.PGMD5HashRequest")
	proto.RegisterType((*PGSHA256HashRequest)(nil), "approzium.authenticator.protos.PGSHA256HashRequest")
	proto.RegisterType((*PGMD5Response)(nil), "approzium.authenticator.protos.PGMD5Response")
	proto.RegisterType((*PGSHA256Response)(nil), "approzium.authenticator.protos.PGSHA256Response")
}

func init() { proto.RegisterFile("authenticator.proto", fileDescriptor_e86ec39f7c35dea3) }

var fileDescriptor_e86ec39f7c35dea3 = []byte{
	// 528 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x52, 0xdd, 0x6e, 0xd3, 0x30,
	0x18, 0x5d, 0xda, 0x51, 0xa9, 0x5f, 0xd7, 0x1f, 0x3c, 0x31, 0x4a, 0x85, 0xa6, 0x29, 0x17, 0x68,
	0x02, 0xad, 0x4c, 0x9d, 0xc6, 0x0d, 0x42, 0x28, 0x5a, 0xa7, 0xb4, 0xd2, 0xd6, 0x46, 0xd9, 0x60,
	0xe2, 0xca, 0x78, 0x8d, 0x49, 0x2d, 0x35, 0x71, 0x66, 0x3b, 0x2a, 0xec, 0x11, 0x79, 0x1b, 0x5e,
	0x00, 0xa1, 0x38, 0xee, 0x2f, 0x88, 0xee, 0x6a, 0x77, 0x39, 0xe7, 0xcb, 0xf1, 0xf1, 0x77, 0x7c,
	0x60, 0x97, 0xa4, 0x6a, 0x4c, 0x63, 0xc5, 0x46, 0x44, 0x71, 0xd1, 0x4e, 0x04, 0x57, 0x1c, 0xed,
	0x93, 0x24, 0x11, 0xfc, 0x9e, 0xa5, 0x51, 0xfb, 0x1f, 0x63, 0x69, 0x0b, 0xa8, 0x38, 0x37, 0x57,
	0xfd, 0x20, 0x1b, 0xa8, 0x1f, 0xe8, 0x3d, 0xb4, 0x24, 0x0b, 0x63, 0x1a, 0xe0, 0x90, 0x2a, 0x3c,
	0x22, 0x93, 0x09, 0x15, 0x98, 0x99, 0x69, 0xd3, 0x3a, 0xb0, 0x0e, 0xcb, 0xfe, 0xf3, 0xfc, 0x0f,
	0x97, 0xaa, 0x33, 0x3d, 0x9f, 0x8b, 0x5f, 0x41, 0x7d, 0x34, 0x21, 0x2c, 0xa2, 0x01, 0x66, 0x24,
	0xc2, 0x44, 0xc4, 0xcd, 0x82, 0x56, 0x54, 0x0d, 0xdd, 0x27, 0x91, 0x23, 0x62, 0xfb, 0x97, 0x05,
	0x75, 0x8f, 0x48, 0x39, 0xe5, 0x22, 0xf0, 0xe9, 0x5d, 0x4a, 0xa5, 0x42, 0x37, 0x99, 0x96, 0xd1,
	0x58, 0xe1, 0x09, 0x89, 0xc3, 0x94, 0x84, 0x54, 0xbb, 0xd5, 0x3a, 0xed, 0xf6, 0xff, 0x37, 0x68,
	0x9f, 0x69, 0xd9, 0x85, 0x51, 0xf9, 0xb5, 0xd1, 0x0a, 0x46, 0x7b, 0x50, 0x0a, 0x6e, 0xc7, 0x5c,
	0x2a, 0x73, 0x17, 0x83, 0x72, 0x3e, 0xe1, 0x42, 0x35, 0x8b, 0x33, 0x3e, 0x43, 0x39, 0x9f, 0x4a,
	0x2a, 0x9a, 0xdb, 0x33, 0x3e, 0x43, 0xe8, 0x03, 0x14, 0xc9, 0x54, 0x36, 0x9f, 0x1c, 0x58, 0x87,
	0x95, 0xce, 0x9b, 0x4d, 0x97, 0x5a, 0xca, 0xd4, 0xcf, 0x74, 0xf6, 0x77, 0x68, 0x78, 0xee, 0x65,
	0xf7, 0xb4, 0x47, 0xe4, 0x78, 0xb6, 0xb3, 0x07, 0x95, 0x64, 0x1a, 0x60, 0x91, 0x43, 0xbd, 0x6f,
	0xa5, 0xf3, 0x76, 0xd3, 0xd1, 0x6b, 0xc9, 0xf9, 0x90, 0x4c, 0xe7, 0x29, 0x22, 0xd8, 0x96, 0x64,
	0x92, 0xaf, 0xba, 0xe3, 0xeb, 0x6f, 0xfb, 0xa7, 0x05, 0xbb, 0x9e, 0x7b, 0xd5, 0x73, 0x3a, 0xa7,
	0xef, 0x1e, 0xcf, 0xbd, 0x9c, 0xbb, 0xa3, 0x7d, 0x00, 0xa6, 0xa8, 0x20, 0x8a, 0xf1, 0x58, 0xea,
	0xa8, 0xab, 0xfe, 0x12, 0x83, 0x8e, 0x00, 0x2d, 0xf9, 0x30, 0x1e, 0xe3, 0x48, 0x86, 0x26, 0xfa,
	0xa7, 0xab, 0x93, 0x4b, 0x19, 0xda, 0x0e, 0x54, 0x75, 0x8c, 0x3e, 0x95, 0x09, 0x8f, 0x25, 0xcd,
	0x3c, 0xc7, 0x44, 0x8e, 0x4d, 0x35, 0xf5, 0x37, 0x7a, 0x09, 0x65, 0xb3, 0x15, 0x0b, 0xcc, 0x65,
	0x16, 0x84, 0xfd, 0x35, 0x7b, 0x89, 0x3c, 0x8e, 0xf9, 0x29, 0x7b, 0x50, 0x1a, 0x25, 0x82, 0xf3,
	0x6f, 0xe6, 0x1c, 0x83, 0x32, 0x5e, 0xe6, 0xbc, 0x29, 0x4f, 0x8e, 0x56, 0x1d, 0x8a, 0x6b, 0x0e,
	0xaf, 0x3f, 0x42, 0x6d, 0xb5, 0x94, 0xe8, 0x05, 0x3c, 0xbb, 0x70, 0x06, 0xee, 0x27, 0xc7, 0x3d,
	0xc7, 0x83, 0xe1, 0x35, 0xf6, 0xfc, 0xe1, 0xe7, 0x7e, 0xf7, 0xbc, 0xdb, 0xd8, 0x42, 0x00, 0x25,
	0xef, 0xcb, 0x75, 0x6f, 0x38, 0x68, 0x58, 0xa8, 0x04, 0x05, 0x77, 0xd8, 0x28, 0x74, 0x7e, 0x5b,
	0x50, 0x75, 0x96, 0xd3, 0x47, 0x77, 0xb0, 0xe3, 0x52, 0x35, 0x6f, 0x10, 0x3a, 0xde, 0xf8, 0x4e,
	0x6b, 0x65, 0x6b, 0x1d, 0x3d, 0x48, 0x31, 0x4b, 0xc4, 0xde, 0x42, 0xf7, 0x50, 0xd7, 0x96, 0x8b,
	0xe6, 0xa0, 0x93, 0xcd, 0x67, 0xfc, 0xd5, 0xb3, 0xd6, 0xf1, 0x43, 0x45, 0x0b, 0xef, 0xdb, 0x92,
	0x1e, 0x9d, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0x32, 0x89, 0xa2, 0xe2, 0xd3, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthenticatorClient is the client API for Authenticator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthenticatorClient interface {
	GetPGMD5Hash(ctx context.Context, in *PGMD5HashRequest, opts ...grpc.CallOption) (*PGMD5Response, error)
	GetPGSHA256Hash(ctx context.Context, in *PGSHA256HashRequest, opts ...grpc.CallOption) (*PGSHA256Response, error)
}

type authenticatorClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticatorClient(cc grpc.ClientConnInterface) AuthenticatorClient {
	return &authenticatorClient{cc}
}

func (c *authenticatorClient) GetPGMD5Hash(ctx context.Context, in *PGMD5HashRequest, opts ...grpc.CallOption) (*PGMD5Response, error) {
	out := new(PGMD5Response)
	err := c.cc.Invoke(ctx, "/approzium.authenticator.protos.Authenticator/GetPGMD5Hash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticatorClient) GetPGSHA256Hash(ctx context.Context, in *PGSHA256HashRequest, opts ...grpc.CallOption) (*PGSHA256Response, error) {
	out := new(PGSHA256Response)
	err := c.cc.Invoke(ctx, "/approzium.authenticator.protos.Authenticator/GetPGSHA256Hash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticatorServer is the server API for Authenticator service.
type AuthenticatorServer interface {
	GetPGMD5Hash(context.Context, *PGMD5HashRequest) (*PGMD5Response, error)
	GetPGSHA256Hash(context.Context, *PGSHA256HashRequest) (*PGSHA256Response, error)
}

// UnimplementedAuthenticatorServer can be embedded to have forward compatible implementations.
type UnimplementedAuthenticatorServer struct {
}

func (*UnimplementedAuthenticatorServer) GetPGMD5Hash(ctx context.Context, req *PGMD5HashRequest) (*PGMD5Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPGMD5Hash not implemented")
}
func (*UnimplementedAuthenticatorServer) GetPGSHA256Hash(ctx context.Context, req *PGSHA256HashRequest) (*PGSHA256Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPGSHA256Hash not implemented")
}

func RegisterAuthenticatorServer(s *grpc.Server, srv AuthenticatorServer) {
	s.RegisterService(&_Authenticator_serviceDesc, srv)
}

func _Authenticator_GetPGMD5Hash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PGMD5HashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticatorServer).GetPGMD5Hash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/approzium.authenticator.protos.Authenticator/GetPGMD5Hash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticatorServer).GetPGMD5Hash(ctx, req.(*PGMD5HashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authenticator_GetPGSHA256Hash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PGSHA256HashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticatorServer).GetPGSHA256Hash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/approzium.authenticator.protos.Authenticator/GetPGSHA256Hash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticatorServer).GetPGSHA256Hash(ctx, req.(*PGSHA256HashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Authenticator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "approzium.authenticator.protos.Authenticator",
	HandlerType: (*AuthenticatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPGMD5Hash",
			Handler:    _Authenticator_GetPGMD5Hash_Handler,
		},
		{
			MethodName: "GetPGSHA256Hash",
			Handler:    _Authenticator_GetPGSHA256Hash_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authenticator.proto",
}
