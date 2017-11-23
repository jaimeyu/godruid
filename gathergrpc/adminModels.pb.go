// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gathergrpc/adminModels.proto

package gathergrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Provides the metadata of a Tenant. This information is used to
// idetify/describe a full Tenant entity.
type TenantDescriptor struct {
	Datatype              string    `protobuf:"bytes,3,opt,name=datatype" json:"datatype,omitempty"`
	Name                  string    `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	UrlSubdomain          string    `protobuf:"bytes,5,opt,name=urlSubdomain" json:"urlSubdomain,omitempty"`
	State                 UserState `protobuf:"varint,6,opt,name=state,enum=gathergrpc.UserState" json:"state,omitempty"`
	CreatedTimestamp      int64     `protobuf:"varint,7,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64     `protobuf:"varint,8,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *TenantDescriptor) Reset()                    { *m = TenantDescriptor{} }
func (m *TenantDescriptor) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptor) ProtoMessage()               {}
func (*TenantDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *TenantDescriptor) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *TenantDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TenantDescriptor) GetUrlSubdomain() string {
	if m != nil {
		return m.UrlSubdomain
	}
	return ""
}

func (m *TenantDescriptor) GetState() UserState {
	if m != nil {
		return m.State
	}
	return UserState_UNKNOWN
}

func (m *TenantDescriptor) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *TenantDescriptor) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

// TenantDescriptorRequest - wrapper for passing TenantDescriptor
// data as a request to the service.
type TenantDescriptorRequest struct {
	XId  string            `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string            `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantDescriptor `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantDescriptorRequest) Reset()                    { *m = TenantDescriptorRequest{} }
func (m *TenantDescriptorRequest) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptorRequest) ProtoMessage()               {}
func (*TenantDescriptorRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *TenantDescriptorRequest) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantDescriptorRequest) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantDescriptorRequest) GetData() *TenantDescriptor {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantDescriptorResponse - wrapper for passing TenantDescriptor
// data as a response from the service.
type TenantDescriptorResponse struct {
	XId  string            `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string            `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantDescriptor `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantDescriptorResponse) Reset()                    { *m = TenantDescriptorResponse{} }
func (m *TenantDescriptorResponse) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptorResponse) ProtoMessage()               {}
func (*TenantDescriptorResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *TenantDescriptorResponse) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantDescriptorResponse) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantDescriptorResponse) GetData() *TenantDescriptor {
	if m != nil {
		return m.Data
	}
	return nil
}

// Wrapper message to provide a response in the form of
// a container of multiple TenantDescriptor objects.
type TenantDescriptorListResponse struct {
	Data []*TenantDescriptorResponse `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *TenantDescriptorListResponse) Reset()                    { *m = TenantDescriptorListResponse{} }
func (m *TenantDescriptorListResponse) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptorListResponse) ProtoMessage()               {}
func (*TenantDescriptorListResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *TenantDescriptorListResponse) GetData() []*TenantDescriptorResponse {
	if m != nil {
		return m.Data
	}
	return nil
}

// User data for an Adminstrative User.
type AdminUser struct {
	Datatype              string    `protobuf:"bytes,3,opt,name=datatype" json:"datatype,omitempty"`
	Username              string    `protobuf:"bytes,4,opt,name=username" json:"username,omitempty"`
	Password              string    `protobuf:"bytes,5,opt,name=password" json:"password,omitempty"`
	SendOnboardingEmail   bool      `protobuf:"varint,6,opt,name=sendOnboardingEmail" json:"sendOnboardingEmail,omitempty"`
	OnboardingToken       string    `protobuf:"bytes,7,opt,name=onboardingToken" json:"onboardingToken,omitempty"`
	UserVerified          bool      `protobuf:"varint,8,opt,name=userVerified" json:"userVerified,omitempty"`
	State                 UserState `protobuf:"varint,9,opt,name=state,enum=gathergrpc.UserState" json:"state,omitempty"`
	CreatedTimestamp      int64     `protobuf:"varint,10,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64     `protobuf:"varint,11,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *AdminUser) Reset()                    { *m = AdminUser{} }
func (m *AdminUser) String() string            { return proto.CompactTextString(m) }
func (*AdminUser) ProtoMessage()               {}
func (*AdminUser) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *AdminUser) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *AdminUser) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AdminUser) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *AdminUser) GetSendOnboardingEmail() bool {
	if m != nil {
		return m.SendOnboardingEmail
	}
	return false
}

func (m *AdminUser) GetOnboardingToken() string {
	if m != nil {
		return m.OnboardingToken
	}
	return ""
}

func (m *AdminUser) GetUserVerified() bool {
	if m != nil {
		return m.UserVerified
	}
	return false
}

func (m *AdminUser) GetState() UserState {
	if m != nil {
		return m.State
	}
	return UserState_UNKNOWN
}

func (m *AdminUser) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *AdminUser) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

// AdminUserRequest - wrapper for passing AdminUser
// data as a request to the service.
type AdminUserRequest struct {
	XId  string     `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string     `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *AdminUser `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *AdminUserRequest) Reset()                    { *m = AdminUserRequest{} }
func (m *AdminUserRequest) String() string            { return proto.CompactTextString(m) }
func (*AdminUserRequest) ProtoMessage()               {}
func (*AdminUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *AdminUserRequest) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *AdminUserRequest) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *AdminUserRequest) GetData() *AdminUser {
	if m != nil {
		return m.Data
	}
	return nil
}

// AdminUserResponse - wrapper for passing AdminUser
// data as a response from the service.
type AdminUserResponse struct {
	XId  string     `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string     `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *AdminUser `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *AdminUserResponse) Reset()                    { *m = AdminUserResponse{} }
func (m *AdminUserResponse) String() string            { return proto.CompactTextString(m) }
func (*AdminUserResponse) ProtoMessage()               {}
func (*AdminUserResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *AdminUserResponse) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *AdminUserResponse) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *AdminUserResponse) GetData() *AdminUser {
	if m != nil {
		return m.Data
	}
	return nil
}

// Wrapper message to provide a response in the form of
// a container of multiple AdminUser objects.
type AdminUserListResponse struct {
	Data []*AdminUserResponse `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *AdminUserListResponse) Reset()                    { *m = AdminUserListResponse{} }
func (m *AdminUserListResponse) String() string            { return proto.CompactTextString(m) }
func (*AdminUserListResponse) ProtoMessage()               {}
func (*AdminUserListResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *AdminUserListResponse) GetData() []*AdminUserResponse {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*TenantDescriptor)(nil), "gathergrpc.TenantDescriptor")
	proto.RegisterType((*TenantDescriptorRequest)(nil), "gathergrpc.TenantDescriptorRequest")
	proto.RegisterType((*TenantDescriptorResponse)(nil), "gathergrpc.TenantDescriptorResponse")
	proto.RegisterType((*TenantDescriptorListResponse)(nil), "gathergrpc.TenantDescriptorListResponse")
	proto.RegisterType((*AdminUser)(nil), "gathergrpc.AdminUser")
	proto.RegisterType((*AdminUserRequest)(nil), "gathergrpc.AdminUserRequest")
	proto.RegisterType((*AdminUserResponse)(nil), "gathergrpc.AdminUserResponse")
	proto.RegisterType((*AdminUserListResponse)(nil), "gathergrpc.AdminUserListResponse")
}

func init() { proto.RegisterFile("gathergrpc/adminModels.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 463 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x94, 0xcf, 0x8b, 0x13, 0x31,
	0x14, 0xc7, 0x99, 0x4e, 0x77, 0x9d, 0xbe, 0x8a, 0xdb, 0x8d, 0x14, 0x87, 0xd2, 0x85, 0x52, 0x3c,
	0x8c, 0x0a, 0x75, 0xad, 0x1e, 0xbc, 0x0a, 0x7a, 0x50, 0x5c, 0x84, 0x6c, 0x15, 0x6f, 0x4b, 0xda,
	0x3c, 0x6b, 0xb0, 0xf9, 0x61, 0x92, 0xae, 0xf8, 0x4f, 0x7b, 0xf3, 0x2e, 0x93, 0x96, 0xe9, 0xa4,
	0x2d, 0x85, 0x0a, 0x7b, 0x6b, 0xbe, 0xdf, 0x97, 0xef, 0xeb, 0x7b, 0xf9, 0x30, 0xd0, 0x9f, 0x33,
	0xff, 0x1d, 0xed, 0xdc, 0x9a, 0xd9, 0x73, 0xc6, 0xa5, 0x50, 0x57, 0x9a, 0xe3, 0xc2, 0x8d, 0x8c,
	0xd5, 0x5e, 0x13, 0xd8, 0xb8, 0xbd, 0x8b, 0x5a, 0xe5, 0x4c, 0x4b, 0xa9, 0xa3, 0xd2, 0xe1, 0xdf,
	0x04, 0x3a, 0x13, 0x54, 0x4c, 0xf9, 0xb7, 0xe8, 0x66, 0x56, 0x18, 0xaf, 0x2d, 0xe9, 0x41, 0xc6,
	0x99, 0x67, 0xfe, 0xb7, 0xc1, 0x3c, 0x1d, 0x24, 0x45, 0x8b, 0x56, 0x67, 0x42, 0xa0, 0xa9, 0x98,
	0xc4, 0xbc, 0x19, 0xf4, 0xf0, 0x9b, 0x0c, 0xe1, 0xfe, 0xd2, 0x2e, 0xae, 0x97, 0x53, 0xae, 0x25,
	0x13, 0x2a, 0x3f, 0x09, 0x5e, 0xa4, 0x91, 0x67, 0x70, 0xe2, 0x3c, 0xf3, 0x98, 0x9f, 0x0e, 0x92,
	0xe2, 0xc1, 0xb8, 0x3b, 0xda, 0xfc, 0xaf, 0xd1, 0x67, 0x87, 0xf6, 0xba, 0x34, 0xe9, 0xaa, 0x86,
	0x3c, 0x85, 0xce, 0xcc, 0x22, 0xf3, 0xc8, 0x27, 0x42, 0xa2, 0xf3, 0x4c, 0x9a, 0xfc, 0xde, 0x20,
	0x29, 0x52, 0xba, 0xa3, 0x93, 0x57, 0xd0, 0x5d, 0x30, 0xe7, 0xaf, 0x34, 0x17, 0xdf, 0x44, 0xfd,
	0x42, 0x16, 0x2e, 0xec, 0x37, 0x87, 0x1a, 0x1e, 0x6d, 0x8f, 0x4d, 0xf1, 0xe7, 0x12, 0x9d, 0x27,
	0x67, 0x90, 0xde, 0x08, 0x9e, 0x27, 0x61, 0x88, 0xc6, 0x7b, 0x4e, 0xce, 0xa1, 0x79, 0x63, 0xf1,
	0x36, 0x6f, 0x04, 0x25, 0xa5, 0x78, 0x4b, 0x2e, 0xa1, 0x59, 0x6e, 0x24, 0x6c, 0xa7, 0x3d, 0xee,
	0xd7, 0x87, 0xd9, 0x89, 0x0d, 0x95, 0x43, 0x03, 0xf9, 0x6e, 0x43, 0x67, 0xb4, 0x72, 0x78, 0x47,
	0x1d, 0xbf, 0x42, 0x7f, 0xdb, 0xf9, 0x28, 0x9c, 0xaf, 0xba, 0xbe, 0x5e, 0x27, 0x26, 0x83, 0xb4,
	0x68, 0x8f, 0x1f, 0x1f, 0x4c, 0x5c, 0xdf, 0x59, 0x27, 0xff, 0x69, 0x40, 0xeb, 0x4d, 0x49, 0x5d,
	0xf9, 0x70, 0x07, 0x69, 0xe9, 0x41, 0xb6, 0x74, 0x68, 0x6b, 0xc4, 0x54, 0xe7, 0xd2, 0x33, 0xcc,
	0xb9, 0x5f, 0xda, 0xf2, 0x35, 0x31, 0xd5, 0x99, 0x5c, 0xc2, 0x43, 0x87, 0x8a, 0x7f, 0x52, 0x53,
	0xcd, 0x2c, 0x17, 0x6a, 0xfe, 0x4e, 0x32, 0xb1, 0x08, 0xec, 0x64, 0x74, 0x9f, 0x45, 0x0a, 0x38,
	0xd3, 0x95, 0x34, 0xd1, 0x3f, 0x50, 0x05, 0x62, 0x5a, 0x74, 0x5b, 0x0e, 0xb4, 0x3a, 0xb4, 0x5f,
	0xd0, 0x06, 0x26, 0x02, 0x27, 0x19, 0x8d, 0xb4, 0x0d, 0xad, 0xad, 0xff, 0xa4, 0x15, 0x8e, 0xa5,
	0xb5, 0x7d, 0x88, 0x56, 0x06, 0x9d, 0x6a, 0xdf, 0xc7, 0x60, 0xfa, 0x24, 0x82, 0x26, 0x9a, 0x62,
	0x93, 0xb7, 0x7a, 0xd3, 0x29, 0x9c, 0xd7, 0x5a, 0x1c, 0x01, 0xe6, 0x11, 0x3d, 0x3e, 0x40, 0xb7,
	0x92, 0x22, 0x14, 0x5f, 0x44, 0x28, 0x5e, 0xec, 0xcf, 0x88, 0x18, 0x9c, 0x9e, 0x86, 0xef, 0xd7,
	0xcb, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x90, 0xae, 0x3e, 0x3d, 0x0a, 0x05, 0x00, 0x00,
}
