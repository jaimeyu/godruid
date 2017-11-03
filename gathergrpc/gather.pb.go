// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gathergrpc/gather.proto

/*
Package gathergrpc is a generated protocol buffer package.

It is generated from these files:
	gathergrpc/gather.proto

It has these top-level messages:
	TenantDescriptor
	AdminUser
	AdminUserList
*/
package gathergrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/wrappers"
import google_protobuf2 "github.com/golang/protobuf/ptypes/empty"

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

// Enumeration of User state.
type UserState int32

const (
	UserState_UNKNOWN        UserState = 0
	UserState_INVITED        UserState = 1
	UserState_ACTIVE         UserState = 2
	UserState_SUSPENDED      UserState = 3
	UserState_PENDING_DELETE UserState = 4
)

var UserState_name = map[int32]string{
	0: "UNKNOWN",
	1: "INVITED",
	2: "ACTIVE",
	3: "SUSPENDED",
	4: "PENDING_DELETE",
}
var UserState_value = map[string]int32{
	"UNKNOWN":        0,
	"INVITED":        1,
	"ACTIVE":         2,
	"SUSPENDED":      3,
	"PENDING_DELETE": 4,
}

func (x UserState) String() string {
	return proto.EnumName(UserState_name, int32(x))
}
func (UserState) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Provides the metadata of a Tenant. This information is used to
// idetify/describe a full Tenant entity.
type TenantDescriptor struct {
	Id                    string    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Rev                   string    `protobuf:"bytes,2,opt,name=rev" json:"rev,omitempty"`
	Name                  string    `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	UrlSubdomain          string    `protobuf:"bytes,4,opt,name=urlSubdomain" json:"urlSubdomain,omitempty"`
	State                 UserState `protobuf:"varint,5,opt,name=state,enum=gathergrpc.UserState" json:"state,omitempty"`
	CreatedTimestamp      int64     `protobuf:"varint,6,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64     `protobuf:"varint,7,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *TenantDescriptor) Reset()                    { *m = TenantDescriptor{} }
func (m *TenantDescriptor) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptor) ProtoMessage()               {}
func (*TenantDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TenantDescriptor) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TenantDescriptor) GetRev() string {
	if m != nil {
		return m.Rev
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

// User data for an Adminstrative User.
type AdminUser struct {
	Id                    string    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Rev                   string    `protobuf:"bytes,2,opt,name=rev" json:"rev,omitempty"`
	Username              string    `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	Password              string    `protobuf:"bytes,4,opt,name=password" json:"password,omitempty"`
	SendOnboardingEmail   bool      `protobuf:"varint,5,opt,name=sendOnboardingEmail" json:"sendOnboardingEmail,omitempty"`
	OnboardingToken       string    `protobuf:"bytes,6,opt,name=onboardingToken" json:"onboardingToken,omitempty"`
	UserVerified          bool      `protobuf:"varint,7,opt,name=userVerified" json:"userVerified,omitempty"`
	State                 UserState `protobuf:"varint,8,opt,name=state,enum=gathergrpc.UserState" json:"state,omitempty"`
	CreatedTimestamp      int64     `protobuf:"varint,9,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64     `protobuf:"varint,10,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *AdminUser) Reset()                    { *m = AdminUser{} }
func (m *AdminUser) String() string            { return proto.CompactTextString(m) }
func (*AdminUser) ProtoMessage()               {}
func (*AdminUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AdminUser) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AdminUser) GetRev() string {
	if m != nil {
		return m.Rev
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

// Wrapper message to provide a response in the form of
// a container of multiple AdminUser objects.
type AdminUserList struct {
	List []*AdminUser `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *AdminUserList) Reset()                    { *m = AdminUserList{} }
func (m *AdminUserList) String() string            { return proto.CompactTextString(m) }
func (*AdminUserList) ProtoMessage()               {}
func (*AdminUserList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AdminUserList) GetList() []*AdminUser {
	if m != nil {
		return m.List
	}
	return nil
}

func init() {
	proto.RegisterType((*TenantDescriptor)(nil), "gathergrpc.TenantDescriptor")
	proto.RegisterType((*AdminUser)(nil), "gathergrpc.AdminUser")
	proto.RegisterType((*AdminUserList)(nil), "gathergrpc.AdminUserList")
	proto.RegisterEnum("gathergrpc.UserState", UserState_name, UserState_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AdminProvisioningService service

type AdminProvisioningServiceClient interface {
	// Create a User with Administrative access.
	CreateAdminUser(ctx context.Context, in *AdminUser, opts ...grpc.CallOption) (*AdminUser, error)
	// Update a User with Administrative access.
	UpdateAdminUser(ctx context.Context, in *AdminUser, opts ...grpc.CallOption) (*AdminUser, error)
	// Delete a User with Administrative access.
	DeleteAdminUser(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*AdminUser, error)
	// Retrieve and Administrative User by id.
	GetAdminUser(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*AdminUser, error)
	// Retrieve all Administrative Users.
	GetAllAdminUsers(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*AdminUserList, error)
	// Creates a Tenant and returns a TenantDescriptor which provides
	// metadata for the newly created Tenant.
	CreateTenant(ctx context.Context, in *TenantDescriptor, opts ...grpc.CallOption) (*TenantDescriptor, error)
	// Updates a TenantDescriptor, which provides metadata
	// for the specified Tenant.
	UpdateTenantDescriptor(ctx context.Context, in *TenantDescriptor, opts ...grpc.CallOption) (*TenantDescriptor, error)
	// Deletes a Tenant and returns a TenantDescriptor which provides
	// metadata for the now deleted Tenant.
	DeleteTenant(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*TenantDescriptor, error)
	// Retrieves the metadata of a single Tenant by id.
	GetTenantDescriptor(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*TenantDescriptor, error)
}

type adminProvisioningServiceClient struct {
	cc *grpc.ClientConn
}

func NewAdminProvisioningServiceClient(cc *grpc.ClientConn) AdminProvisioningServiceClient {
	return &adminProvisioningServiceClient{cc}
}

func (c *adminProvisioningServiceClient) CreateAdminUser(ctx context.Context, in *AdminUser, opts ...grpc.CallOption) (*AdminUser, error) {
	out := new(AdminUser)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/CreateAdminUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) UpdateAdminUser(ctx context.Context, in *AdminUser, opts ...grpc.CallOption) (*AdminUser, error) {
	out := new(AdminUser)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/UpdateAdminUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) DeleteAdminUser(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*AdminUser, error) {
	out := new(AdminUser)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/DeleteAdminUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) GetAdminUser(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*AdminUser, error) {
	out := new(AdminUser)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/GetAdminUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) GetAllAdminUsers(ctx context.Context, in *google_protobuf2.Empty, opts ...grpc.CallOption) (*AdminUserList, error) {
	out := new(AdminUserList)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/GetAllAdminUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) CreateTenant(ctx context.Context, in *TenantDescriptor, opts ...grpc.CallOption) (*TenantDescriptor, error) {
	out := new(TenantDescriptor)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/CreateTenant", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) UpdateTenantDescriptor(ctx context.Context, in *TenantDescriptor, opts ...grpc.CallOption) (*TenantDescriptor, error) {
	out := new(TenantDescriptor)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/UpdateTenantDescriptor", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) DeleteTenant(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*TenantDescriptor, error) {
	out := new(TenantDescriptor)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/DeleteTenant", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminProvisioningServiceClient) GetTenantDescriptor(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*TenantDescriptor, error) {
	out := new(TenantDescriptor)
	err := grpc.Invoke(ctx, "/gathergrpc.AdminProvisioningService/GetTenantDescriptor", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AdminProvisioningService service

type AdminProvisioningServiceServer interface {
	// Create a User with Administrative access.
	CreateAdminUser(context.Context, *AdminUser) (*AdminUser, error)
	// Update a User with Administrative access.
	UpdateAdminUser(context.Context, *AdminUser) (*AdminUser, error)
	// Delete a User with Administrative access.
	DeleteAdminUser(context.Context, *google_protobuf1.StringValue) (*AdminUser, error)
	// Retrieve and Administrative User by id.
	GetAdminUser(context.Context, *google_protobuf1.StringValue) (*AdminUser, error)
	// Retrieve all Administrative Users.
	GetAllAdminUsers(context.Context, *google_protobuf2.Empty) (*AdminUserList, error)
	// Creates a Tenant and returns a TenantDescriptor which provides
	// metadata for the newly created Tenant.
	CreateTenant(context.Context, *TenantDescriptor) (*TenantDescriptor, error)
	// Updates a TenantDescriptor, which provides metadata
	// for the specified Tenant.
	UpdateTenantDescriptor(context.Context, *TenantDescriptor) (*TenantDescriptor, error)
	// Deletes a Tenant and returns a TenantDescriptor which provides
	// metadata for the now deleted Tenant.
	DeleteTenant(context.Context, *google_protobuf1.StringValue) (*TenantDescriptor, error)
	// Retrieves the metadata of a single Tenant by id.
	GetTenantDescriptor(context.Context, *google_protobuf1.StringValue) (*TenantDescriptor, error)
}

func RegisterAdminProvisioningServiceServer(s *grpc.Server, srv AdminProvisioningServiceServer) {
	s.RegisterService(&_AdminProvisioningService_serviceDesc, srv)
}

func _AdminProvisioningService_CreateAdminUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdminUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).CreateAdminUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/CreateAdminUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).CreateAdminUser(ctx, req.(*AdminUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_UpdateAdminUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdminUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).UpdateAdminUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/UpdateAdminUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).UpdateAdminUser(ctx, req.(*AdminUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_DeleteAdminUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).DeleteAdminUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/DeleteAdminUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).DeleteAdminUser(ctx, req.(*google_protobuf1.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_GetAdminUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).GetAdminUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/GetAdminUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).GetAdminUser(ctx, req.(*google_protobuf1.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_GetAllAdminUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf2.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).GetAllAdminUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/GetAllAdminUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).GetAllAdminUsers(ctx, req.(*google_protobuf2.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_CreateTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).CreateTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/CreateTenant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).CreateTenant(ctx, req.(*TenantDescriptor))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_UpdateTenantDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).UpdateTenantDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/UpdateTenantDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).UpdateTenantDescriptor(ctx, req.(*TenantDescriptor))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_DeleteTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).DeleteTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/DeleteTenant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).DeleteTenant(ctx, req.(*google_protobuf1.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminProvisioningService_GetTenantDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminProvisioningServiceServer).GetTenantDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gathergrpc.AdminProvisioningService/GetTenantDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminProvisioningServiceServer).GetTenantDescriptor(ctx, req.(*google_protobuf1.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdminProvisioningService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gathergrpc.AdminProvisioningService",
	HandlerType: (*AdminProvisioningServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAdminUser",
			Handler:    _AdminProvisioningService_CreateAdminUser_Handler,
		},
		{
			MethodName: "UpdateAdminUser",
			Handler:    _AdminProvisioningService_UpdateAdminUser_Handler,
		},
		{
			MethodName: "DeleteAdminUser",
			Handler:    _AdminProvisioningService_DeleteAdminUser_Handler,
		},
		{
			MethodName: "GetAdminUser",
			Handler:    _AdminProvisioningService_GetAdminUser_Handler,
		},
		{
			MethodName: "GetAllAdminUsers",
			Handler:    _AdminProvisioningService_GetAllAdminUsers_Handler,
		},
		{
			MethodName: "CreateTenant",
			Handler:    _AdminProvisioningService_CreateTenant_Handler,
		},
		{
			MethodName: "UpdateTenantDescriptor",
			Handler:    _AdminProvisioningService_UpdateTenantDescriptor_Handler,
		},
		{
			MethodName: "DeleteTenant",
			Handler:    _AdminProvisioningService_DeleteTenant_Handler,
		},
		{
			MethodName: "GetTenantDescriptor",
			Handler:    _AdminProvisioningService_GetTenantDescriptor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gathergrpc/gather.proto",
}

func init() { proto.RegisterFile("gathergrpc/gather.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 707 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xdd, 0x6e, 0x1a, 0x39,
	0x14, 0xc7, 0x77, 0x80, 0x24, 0x70, 0x42, 0xc8, 0xac, 0x51, 0x92, 0xc9, 0xe4, 0x43, 0x88, 0xbd,
	0x61, 0x59, 0x09, 0x76, 0xb3, 0x7b, 0x95, 0xbb, 0x28, 0x8c, 0x22, 0xb4, 0x29, 0x89, 0xf8, 0xaa,
	0x7a, 0xd5, 0x1a, 0xc6, 0xa1, 0x56, 0x67, 0xec, 0x91, 0x6d, 0x88, 0xaa, 0xaa, 0x37, 0x7d, 0x85,
	0xbe, 0x4b, 0x5f, 0xa4, 0xaf, 0xd0, 0x47, 0xe8, 0x5d, 0x6f, 0x2a, 0x7b, 0x02, 0x84, 0x8f, 0xa0,
	0x36, 0xe9, 0x9d, 0x7d, 0xfe, 0xc7, 0xe7, 0x37, 0xff, 0x63, 0xcf, 0x81, 0xbd, 0x01, 0x56, 0xaf,
	0x89, 0x18, 0x88, 0xa8, 0x5f, 0x8d, 0x97, 0x95, 0x48, 0x70, 0xc5, 0x11, 0x4c, 0x05, 0xf7, 0x70,
	0xc0, 0xf9, 0x20, 0x20, 0x55, 0x1c, 0xd1, 0x2a, 0x66, 0x8c, 0x2b, 0xac, 0x28, 0x67, 0x32, 0xce,
	0x74, 0x8f, 0xef, 0x54, 0xb3, 0xeb, 0x0d, 0x6f, 0xaa, 0xb7, 0x02, 0x47, 0x11, 0x11, 0x63, 0xfd,
	0x60, 0x5e, 0x27, 0x61, 0xa4, 0xde, 0xc6, 0x62, 0xf1, 0x9b, 0x05, 0x76, 0x9b, 0x30, 0xcc, 0x54,
	0x8d, 0xc8, 0xbe, 0xa0, 0x91, 0xe2, 0x02, 0xe5, 0x20, 0x41, 0x7d, 0xc7, 0x2a, 0x58, 0xa5, 0x4c,
	0x33, 0x41, 0x7d, 0x64, 0x43, 0x52, 0x90, 0x91, 0x93, 0x30, 0x01, 0xbd, 0x44, 0x08, 0x52, 0x0c,
	0x87, 0xc4, 0x49, 0x9a, 0x90, 0x59, 0xa3, 0x22, 0x64, 0x87, 0x22, 0x68, 0x0d, 0x7b, 0x3e, 0x0f,
	0x31, 0x65, 0x4e, 0xca, 0x68, 0x33, 0x31, 0xf4, 0x17, 0xac, 0x49, 0x85, 0x15, 0x71, 0xd6, 0x0a,
	0x56, 0x29, 0x77, 0xb2, 0x53, 0x99, 0xba, 0xac, 0x74, 0x24, 0x11, 0x2d, 0x2d, 0x36, 0xe3, 0x1c,
	0x54, 0x06, 0xbb, 0x2f, 0x08, 0x56, 0xc4, 0x6f, 0xd3, 0x90, 0x48, 0x85, 0xc3, 0xc8, 0x59, 0x2f,
	0x58, 0xa5, 0x64, 0x73, 0x21, 0x8e, 0xfe, 0x83, 0x9d, 0x00, 0x4b, 0xf5, 0x8c, 0xfb, 0xf4, 0x86,
	0xde, 0x3f, 0xb0, 0x61, 0x0e, 0x2c, 0x17, 0x8b, 0x5f, 0x13, 0x90, 0x39, 0xf3, 0x43, 0xca, 0x34,
	0xfb, 0x07, 0x6c, 0xbb, 0x90, 0x1e, 0x4a, 0x22, 0xee, 0x59, 0x9f, 0xec, 0xb5, 0x16, 0x61, 0x29,
	0x6f, 0xb9, 0xf0, 0xef, 0xac, 0x4f, 0xf6, 0xe8, 0x6f, 0xc8, 0x4b, 0xc2, 0xfc, 0x2b, 0xd6, 0xe3,
	0x58, 0xf8, 0x94, 0x0d, 0xbc, 0x10, 0xd3, 0xc0, 0x34, 0x21, 0xdd, 0x5c, 0x26, 0xa1, 0x12, 0x6c,
	0xf3, 0x49, 0xa8, 0xcd, 0xdf, 0x10, 0x66, 0xac, 0x67, 0x9a, 0xf3, 0x61, 0xd3, 0x76, 0x49, 0x44,
	0x97, 0x08, 0x63, 0xce, 0x18, 0x4e, 0x37, 0x67, 0x62, 0xd3, 0xb6, 0xa7, 0x1f, 0xd9, 0xf6, 0xcc,
	0xcf, 0xb6, 0x1d, 0x56, 0xb5, 0xfd, 0x14, 0xb6, 0x26, 0x5d, 0xbf, 0xa4, 0x52, 0xa1, 0x3f, 0x21,
	0x15, 0x50, 0xa9, 0x1c, 0xab, 0x90, 0x2c, 0x6d, 0xce, 0x7e, 0xde, 0x24, 0xb1, 0x69, 0x52, 0xca,
	0x6d, 0xc8, 0x4c, 0xbe, 0x18, 0x6d, 0xc2, 0x46, 0xa7, 0xf1, 0x7f, 0xe3, 0xea, 0x79, 0xc3, 0xfe,
	0x4d, 0x6f, 0xea, 0x8d, 0x6e, 0xbd, 0xed, 0xd5, 0x6c, 0x0b, 0x01, 0xac, 0x9f, 0x9d, 0xb7, 0xeb,
	0x5d, 0xcf, 0x4e, 0xa0, 0x2d, 0xc8, 0xb4, 0x3a, 0xad, 0x6b, 0xaf, 0x51, 0xf3, 0x6a, 0x76, 0x12,
	0x21, 0xc8, 0xe9, 0x75, 0xbd, 0x71, 0xf1, 0xb2, 0xe6, 0x5d, 0x7a, 0x6d, 0xcf, 0x4e, 0x9d, 0x7c,
	0xda, 0x00, 0xc7, 0x90, 0xae, 0x05, 0x1f, 0x51, 0x49, 0x39, 0xa3, 0x6c, 0xd0, 0x22, 0x62, 0x44,
	0xfb, 0x04, 0xbd, 0x80, 0xed, 0x73, 0x63, 0x7c, 0xfa, 0x54, 0x96, 0x7f, 0xa2, 0xbb, 0x3c, 0x5c,
	0x74, 0x3e, 0x7c, 0xfe, 0xf2, 0x31, 0x81, 0x8a, 0x5b, 0xe6, 0x0f, 0x1e, 0xfd, 0x53, 0xc5, 0x5a,
	0x3a, 0xb5, 0xca, 0xba, 0x74, 0x27, 0xf2, 0x9f, 0x5e, 0xda, 0x5d, 0x2c, 0x4d, 0x60, 0xbb, 0x46,
	0x02, 0x72, 0xbf, 0xf4, 0x61, 0x25, 0x1e, 0x05, 0x95, 0xf1, 0x28, 0xa8, 0xb4, 0x94, 0xa0, 0x6c,
	0xd0, 0xc5, 0xc1, 0x90, 0x3c, 0x44, 0x38, 0x32, 0x84, 0xbd, 0xf2, 0xce, 0x0c, 0xa1, 0xfa, 0x6e,
	0xa4, 0x0f, 0xbd, 0x47, 0x3d, 0xc8, 0x5e, 0x10, 0xf5, 0x6b, 0x18, 0xe8, 0x01, 0xc6, 0x2b, 0xb0,
	0x35, 0x23, 0x08, 0x26, 0x27, 0x24, 0xda, 0x5d, 0xe0, 0x78, 0x7a, 0xac, 0xb9, 0xfb, 0x4b, 0x09,
	0xfa, 0x95, 0x15, 0xf7, 0x0d, 0x25, 0x8f, 0x7e, 0x9f, 0xa5, 0xe0, 0x20, 0x40, 0x04, 0xb2, 0xf1,
	0x15, 0xc7, 0xb3, 0x50, 0xbb, 0x98, 0x56, 0x99, 0x9f, 0x8f, 0xee, 0x4a, 0x75, 0x8c, 0x29, 0xe6,
	0xc6, 0x18, 0x65, 0x32, 0xf4, 0x9d, 0x28, 0xd8, 0x8d, 0xaf, 0x7b, 0x61, 0xe4, 0x3e, 0x05, 0x78,
	0x6c, 0x80, 0x8e, 0x9b, 0x9f, 0x05, 0x56, 0x43, 0xa2, 0xb0, 0xa6, 0xf6, 0x20, 0x1b, 0xbf, 0x84,
	0xa9, 0xb9, 0x15, 0x57, 0xb4, 0x9a, 0xb5, 0x6b, 0x58, 0x76, 0x79, 0xce, 0x1c, 0x1a, 0x41, 0xfe,
	0x82, 0xa8, 0x65, 0xb6, 0x1e, 0x8d, 0xfa, 0xc3, 0xa0, 0x8e, 0xd0, 0xc1, 0x12, 0x5b, 0xe3, 0xa7,
	0xd1, 0x5b, 0x37, 0x85, 0xff, 0xfd, 0x1e, 0x00, 0x00, 0xff, 0xff, 0x50, 0x0c, 0x9a, 0x2c, 0x48,
	0x07, 0x00, 0x00,
}
