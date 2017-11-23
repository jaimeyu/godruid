// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gathergrpc/commonModels.proto

/*
Package gathergrpc is a generated protocol buffer package.

It is generated from these files:
	gathergrpc/commonModels.proto
	gathergrpc/adminModels.proto
	gathergrpc/tenantModels.proto
	gathergrpc/gather.proto

It has these top-level messages:
	TenantDescriptor
	TenantDescriptorRequest
	TenantDescriptorResponse
	TenantDescriptorListResponse
	AdminUser
	AdminUserRequest
	AdminUserResponse
	AdminUserListResponse
	TenantDomain
	TenantDomainRequest
	TenantDomainResponse
	TenantDomainListResponse
	TenantDomainIdRequest
	TenantIngestionProfile
	TenantIngestionProfileRequest
	TenantIngestionProfileResponse
	TenantIngestionProfileIdRequest
	TenantUser
	TenantUserRequest
	TenantUserResponse
	TenantUserListResponse
	TenantUserIdRequest
	MonitoredObject
	MonitoredObjectRequest
	MonitoredObjectResponse
	MonitoredObjectListResponse
	MonitoredObjectIdRequest
*/
package gathergrpc

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

// Enumeration of User state.
type UserState int32

const (
	UserState_USER_UNKNOWN   UserState = 0
	UserState_INVITED        UserState = 1
	UserState_ACTIVE         UserState = 2
	UserState_SUSPENDED      UserState = 3
	UserState_PENDING_DELETE UserState = 4
)

var UserState_name = map[int32]string{
	0: "USER_UNKNOWN",
	1: "INVITED",
	2: "ACTIVE",
	3: "SUSPENDED",
	4: "PENDING_DELETE",
}
var UserState_value = map[string]int32{
	"USER_UNKNOWN":   0,
	"INVITED":        1,
	"ACTIVE":         2,
	"SUSPENDED":      3,
	"PENDING_DELETE": 4,
}

func (x UserState) String() string {
	return proto.EnumName(UserState_name, int32(x))
}
func (UserState) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterEnum("gathergrpc.UserState", UserState_name, UserState_value)
}

func init() { proto.RegisterFile("gathergrpc/commonModels.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 153 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4d, 0x4f, 0x2c, 0xc9,
	0x48, 0x2d, 0x4a, 0x2f, 0x2a, 0x48, 0xd6, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xf3, 0xcd, 0x4f,
	0x49, 0xcd, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x48, 0x6b, 0x45, 0x72,
	0x71, 0x86, 0x16, 0xa7, 0x16, 0x05, 0x97, 0x24, 0x96, 0xa4, 0x0a, 0x09, 0x70, 0xf1, 0x84, 0x06,
	0xbb, 0x06, 0xc5, 0x87, 0xfa, 0x79, 0xfb, 0xf9, 0x87, 0xfb, 0x09, 0x30, 0x08, 0x71, 0x73, 0xb1,
	0x7b, 0xfa, 0x85, 0x79, 0x86, 0xb8, 0xba, 0x08, 0x30, 0x0a, 0x71, 0x71, 0xb1, 0x39, 0x3a, 0x87,
	0x78, 0x86, 0xb9, 0x0a, 0x30, 0x09, 0xf1, 0x72, 0x71, 0x06, 0x87, 0x06, 0x07, 0xb8, 0xfa, 0xb9,
	0xb8, 0xba, 0x08, 0x30, 0x0b, 0x09, 0x71, 0xf1, 0x81, 0xd8, 0x9e, 0x7e, 0xee, 0xf1, 0x2e, 0xae,
	0x3e, 0xae, 0x21, 0xae, 0x02, 0x2c, 0x49, 0x6c, 0x60, 0xdb, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x4f, 0x8c, 0x50, 0xe9, 0x8e, 0x00, 0x00, 0x00,
}
