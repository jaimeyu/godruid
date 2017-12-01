// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gathergrpc/tenantModels.proto

package gathergrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MonitoredObject_MonitoredObjectType int32

const (
	MonitoredObject_MO_UNKNOWN MonitoredObject_MonitoredObjectType = 0
	MonitoredObject_TWAMP      MonitoredObject_MonitoredObjectType = 1
)

var MonitoredObject_MonitoredObjectType_name = map[int32]string{
	0: "MO_UNKNOWN",
	1: "TWAMP",
}
var MonitoredObject_MonitoredObjectType_value = map[string]int32{
	"MO_UNKNOWN": 0,
	"TWAMP":      1,
}

func (x MonitoredObject_MonitoredObjectType) String() string {
	return proto.EnumName(MonitoredObject_MonitoredObjectType_name, int32(x))
}
func (MonitoredObject_MonitoredObjectType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor2, []int{14, 0}
}

type MonitoredObject_DeviceType int32

const (
	MonitoredObject_DT_UNKNOWN    MonitoredObject_DeviceType = 0
	MonitoredObject_ACCEDIAN_NID  MonitoredObject_DeviceType = 1
	MonitoredObject_ACCEDIAN_VNID MonitoredObject_DeviceType = 2
)

var MonitoredObject_DeviceType_name = map[int32]string{
	0: "DT_UNKNOWN",
	1: "ACCEDIAN_NID",
	2: "ACCEDIAN_VNID",
}
var MonitoredObject_DeviceType_value = map[string]int32{
	"DT_UNKNOWN":    0,
	"ACCEDIAN_NID":  1,
	"ACCEDIAN_VNID": 2,
}

func (x MonitoredObject_DeviceType) String() string {
	return proto.EnumName(MonitoredObject_DeviceType_name, int32(x))
}
func (MonitoredObject_DeviceType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor2, []int{14, 1}
}

// TenantDomain - model for a Domain for a single Tenant.
type TenantDomain struct {
	TenantId              string `protobuf:"bytes,1,opt,name=tenantId" json:"tenantId,omitempty"`
	Datatype              string `protobuf:"bytes,2,opt,name=datatype" json:"datatype,omitempty"`
	Name                  string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	Color                 string `protobuf:"bytes,4,opt,name=color" json:"color,omitempty"`
	CreatedTimestamp      int64  `protobuf:"varint,5,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64  `protobuf:"varint,6,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *TenantDomain) Reset()                    { *m = TenantDomain{} }
func (m *TenantDomain) String() string            { return proto.CompactTextString(m) }
func (*TenantDomain) ProtoMessage()               {}
func (*TenantDomain) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *TenantDomain) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *TenantDomain) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *TenantDomain) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TenantDomain) GetColor() string {
	if m != nil {
		return m.Color
	}
	return ""
}

func (m *TenantDomain) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *TenantDomain) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

// TenantDomainRequest - wrapper for requests that involve a Tenant Domain
type TenantDomainRequest struct {
	XId  string        `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string        `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantDomain `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantDomainRequest) Reset()                    { *m = TenantDomainRequest{} }
func (m *TenantDomainRequest) String() string            { return proto.CompactTextString(m) }
func (*TenantDomainRequest) ProtoMessage()               {}
func (*TenantDomainRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *TenantDomainRequest) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantDomainRequest) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantDomainRequest) GetData() *TenantDomain {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantDomainResponse - wrapper for responses that involve a Tenant Domain
type TenantDomainResponse struct {
	XId  string        `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string        `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantDomain `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantDomainResponse) Reset()                    { *m = TenantDomainResponse{} }
func (m *TenantDomainResponse) String() string            { return proto.CompactTextString(m) }
func (*TenantDomainResponse) ProtoMessage()               {}
func (*TenantDomainResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *TenantDomainResponse) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantDomainResponse) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantDomainResponse) GetData() *TenantDomain {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantDomainListResponse - a wrapper for a list of TenantDomain objects that
// are returned as a response to a request..
type TenantDomainListResponse struct {
	Data []*TenantDomainResponse `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *TenantDomainListResponse) Reset()                    { *m = TenantDomainListResponse{} }
func (m *TenantDomainListResponse) String() string            { return proto.CompactTextString(m) }
func (*TenantDomainListResponse) ProtoMessage()               {}
func (*TenantDomainListResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *TenantDomainListResponse) GetData() []*TenantDomainResponse {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantDomainIdRequest - wrapper for requests that involve a Tenant Domain,
// but only require the domainID to complete the request.
type TenantDomainIdRequest struct {
	TenantId string `protobuf:"bytes,1,opt,name=tenantId" json:"tenantId,omitempty"`
	DomainId string `protobuf:"bytes,2,opt,name=domainId" json:"domainId,omitempty"`
}

func (m *TenantDomainIdRequest) Reset()                    { *m = TenantDomainIdRequest{} }
func (m *TenantDomainIdRequest) String() string            { return proto.CompactTextString(m) }
func (*TenantDomainIdRequest) ProtoMessage()               {}
func (*TenantDomainIdRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *TenantDomainIdRequest) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *TenantDomainIdRequest) GetDomainId() string {
	if m != nil {
		return m.DomainId
	}
	return ""
}

// TenantIngestionProfile - model for the singleton object that
// governs what data is displayed for a Tenant.
type TenantIngestionProfile struct {
	TenantId              string              `protobuf:"bytes,1,opt,name=tenantId" json:"tenantId,omitempty"`
	Datatype              string              `protobuf:"bytes,2,opt,name=datatype" json:"datatype,omitempty"`
	ScpUsername           string              `protobuf:"bytes,3,opt,name=scpUsername" json:"scpUsername,omitempty"`
	ScpPassword           string              `protobuf:"bytes,4,opt,name=scpPassword" json:"scpPassword,omitempty"`
	CreatedTimestamp      int64               `protobuf:"varint,5,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64               `protobuf:"varint,6,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
	ThresholdProfile      []*ThresholdProfile `protobuf:"bytes,7,rep,name=thresholdProfile" json:"thresholdProfile,omitempty"`
}

func (m *TenantIngestionProfile) Reset()                    { *m = TenantIngestionProfile{} }
func (m *TenantIngestionProfile) String() string            { return proto.CompactTextString(m) }
func (*TenantIngestionProfile) ProtoMessage()               {}
func (*TenantIngestionProfile) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *TenantIngestionProfile) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *TenantIngestionProfile) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *TenantIngestionProfile) GetScpUsername() string {
	if m != nil {
		return m.ScpUsername
	}
	return ""
}

func (m *TenantIngestionProfile) GetScpPassword() string {
	if m != nil {
		return m.ScpPassword
	}
	return ""
}

func (m *TenantIngestionProfile) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *TenantIngestionProfile) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

func (m *TenantIngestionProfile) GetThresholdProfile() []*ThresholdProfile {
	if m != nil {
		return m.ThresholdProfile
	}
	return nil
}

// TenantIngestionProfileRequest - wrapper for requests that involve the
// Tenant Ingestion Profile
type TenantIngestionProfileRequest struct {
	XId  string                  `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string                  `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantIngestionProfile `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantIngestionProfileRequest) Reset()                    { *m = TenantIngestionProfileRequest{} }
func (m *TenantIngestionProfileRequest) String() string            { return proto.CompactTextString(m) }
func (*TenantIngestionProfileRequest) ProtoMessage()               {}
func (*TenantIngestionProfileRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

func (m *TenantIngestionProfileRequest) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantIngestionProfileRequest) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantIngestionProfileRequest) GetData() *TenantIngestionProfile {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantIngestionProfileResponse - wrapper to provide a Tenant Ingestion Profile
// in a response.
type TenantIngestionProfileResponse struct {
	XId  string                  `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string                  `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantIngestionProfile `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantIngestionProfileResponse) Reset()                    { *m = TenantIngestionProfileResponse{} }
func (m *TenantIngestionProfileResponse) String() string            { return proto.CompactTextString(m) }
func (*TenantIngestionProfileResponse) ProtoMessage()               {}
func (*TenantIngestionProfileResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{7} }

func (m *TenantIngestionProfileResponse) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantIngestionProfileResponse) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantIngestionProfileResponse) GetData() *TenantIngestionProfile {
	if m != nil {
		return m.Data
	}
	return nil
}

type TenantIngestionProfileIdRequest struct {
	TenantId           string `protobuf:"bytes,1,opt,name=tenantId" json:"tenantId,omitempty"`
	IngestionProfileId string `protobuf:"bytes,2,opt,name=ingestionProfileId" json:"ingestionProfileId,omitempty"`
}

func (m *TenantIngestionProfileIdRequest) Reset()                    { *m = TenantIngestionProfileIdRequest{} }
func (m *TenantIngestionProfileIdRequest) String() string            { return proto.CompactTextString(m) }
func (*TenantIngestionProfileIdRequest) ProtoMessage()               {}
func (*TenantIngestionProfileIdRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{8} }

func (m *TenantIngestionProfileIdRequest) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *TenantIngestionProfileIdRequest) GetIngestionProfileId() string {
	if m != nil {
		return m.IngestionProfileId
	}
	return ""
}

// TenantUser - model for a User that is scoped to a single Tenant.
type TenantUser struct {
	TenantId              string    `protobuf:"bytes,1,opt,name=tenantId" json:"tenantId,omitempty"`
	Datatype              string    `protobuf:"bytes,2,opt,name=datatype" json:"datatype,omitempty"`
	Username              string    `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	Password              string    `protobuf:"bytes,4,opt,name=password" json:"password,omitempty"`
	SendOnboardingEmail   bool      `protobuf:"varint,5,opt,name=sendOnboardingEmail" json:"sendOnboardingEmail,omitempty"`
	OnboardingToken       string    `protobuf:"bytes,6,opt,name=onboardingToken" json:"onboardingToken,omitempty"`
	UserVerified          bool      `protobuf:"varint,7,opt,name=userVerified" json:"userVerified,omitempty"`
	State                 UserState `protobuf:"varint,8,opt,name=state,enum=gathergrpc.UserState" json:"state,omitempty"`
	Domains               []string  `protobuf:"bytes,9,rep,name=domains" json:"domains,omitempty"`
	CreatedTimestamp      int64     `protobuf:"varint,10,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64     `protobuf:"varint,11,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *TenantUser) Reset()                    { *m = TenantUser{} }
func (m *TenantUser) String() string            { return proto.CompactTextString(m) }
func (*TenantUser) ProtoMessage()               {}
func (*TenantUser) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{9} }

func (m *TenantUser) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *TenantUser) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *TenantUser) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *TenantUser) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *TenantUser) GetSendOnboardingEmail() bool {
	if m != nil {
		return m.SendOnboardingEmail
	}
	return false
}

func (m *TenantUser) GetOnboardingToken() string {
	if m != nil {
		return m.OnboardingToken
	}
	return ""
}

func (m *TenantUser) GetUserVerified() bool {
	if m != nil {
		return m.UserVerified
	}
	return false
}

func (m *TenantUser) GetState() UserState {
	if m != nil {
		return m.State
	}
	return UserState_USER_UNKNOWN
}

func (m *TenantUser) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

func (m *TenantUser) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *TenantUser) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

// TenantUserRequest - wrapper for requests that involve a User that
// is scoped to a single Tenant.
type TenantUserRequest struct {
	XId  string      `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string      `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantUser `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantUserRequest) Reset()                    { *m = TenantUserRequest{} }
func (m *TenantUserRequest) String() string            { return proto.CompactTextString(m) }
func (*TenantUserRequest) ProtoMessage()               {}
func (*TenantUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{10} }

func (m *TenantUserRequest) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantUserRequest) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantUserRequest) GetData() *TenantUser {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantUserResponse - wrapper for responses to requests that involve a User that
// is scoped to a single Tenant.
type TenantUserResponse struct {
	XId  string      `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string      `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantUser `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantUserResponse) Reset()                    { *m = TenantUserResponse{} }
func (m *TenantUserResponse) String() string            { return proto.CompactTextString(m) }
func (*TenantUserResponse) ProtoMessage()               {}
func (*TenantUserResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{11} }

func (m *TenantUserResponse) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantUserResponse) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantUserResponse) GetData() *TenantUser {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantUserListResponse - a wrapper to handle requests that return a
// list of TenantUser objects.
type TenantUserListResponse struct {
	Data []*TenantUserResponse `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *TenantUserListResponse) Reset()                    { *m = TenantUserListResponse{} }
func (m *TenantUserListResponse) String() string            { return proto.CompactTextString(m) }
func (*TenantUserListResponse) ProtoMessage()               {}
func (*TenantUserListResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{12} }

func (m *TenantUserListResponse) GetData() []*TenantUserResponse {
	if m != nil {
		return m.Data
	}
	return nil
}

// TenantUserIdRequest - wrapper for requests that involve a Tenant User,
// but only require the userID to complete the request.
type TenantUserIdRequest struct {
	TenantId string `protobuf:"bytes,1,opt,name=tenantId" json:"tenantId,omitempty"`
	UserId   string `protobuf:"bytes,2,opt,name=userId" json:"userId,omitempty"`
}

func (m *TenantUserIdRequest) Reset()                    { *m = TenantUserIdRequest{} }
func (m *TenantUserIdRequest) String() string            { return proto.CompactTextString(m) }
func (*TenantUserIdRequest) ProtoMessage()               {}
func (*TenantUserIdRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{13} }

func (m *TenantUserIdRequest) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *TenantUserIdRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

// MonitoredObject - describes a unique device/object which is reporting
// data from the network.
type MonitoredObject struct {
	Id                    string                              `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	TenantId              string                              `protobuf:"bytes,2,opt,name=tenantId" json:"tenantId,omitempty"`
	Datatype              string                              `protobuf:"bytes,3,opt,name=datatype" json:"datatype,omitempty"`
	ActuatorType          MonitoredObject_DeviceType          `protobuf:"varint,4,opt,name=actuatorType,enum=gathergrpc.MonitoredObject_DeviceType" json:"actuatorType,omitempty"`
	ActuatorName          string                              `protobuf:"bytes,5,opt,name=actuatorName" json:"actuatorName,omitempty"`
	ReflectorType         MonitoredObject_DeviceType          `protobuf:"varint,6,opt,name=reflectorType,enum=gathergrpc.MonitoredObject_DeviceType" json:"reflectorType,omitempty"`
	ReflectorName         string                              `protobuf:"bytes,7,opt,name=reflectorName" json:"reflectorName,omitempty"`
	ObjectName            string                              `protobuf:"bytes,8,opt,name=objectName" json:"objectName,omitempty"`
	ObjectType            MonitoredObject_MonitoredObjectType `protobuf:"varint,9,opt,name=objectType,enum=gathergrpc.MonitoredObject_MonitoredObjectType" json:"objectType,omitempty"`
	CreatedTimestamp      int64                               `protobuf:"varint,10,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64                               `protobuf:"varint,11,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *MonitoredObject) Reset()                    { *m = MonitoredObject{} }
func (m *MonitoredObject) String() string            { return proto.CompactTextString(m) }
func (*MonitoredObject) ProtoMessage()               {}
func (*MonitoredObject) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{14} }

func (m *MonitoredObject) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *MonitoredObject) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *MonitoredObject) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *MonitoredObject) GetActuatorType() MonitoredObject_DeviceType {
	if m != nil {
		return m.ActuatorType
	}
	return MonitoredObject_DT_UNKNOWN
}

func (m *MonitoredObject) GetActuatorName() string {
	if m != nil {
		return m.ActuatorName
	}
	return ""
}

func (m *MonitoredObject) GetReflectorType() MonitoredObject_DeviceType {
	if m != nil {
		return m.ReflectorType
	}
	return MonitoredObject_DT_UNKNOWN
}

func (m *MonitoredObject) GetReflectorName() string {
	if m != nil {
		return m.ReflectorName
	}
	return ""
}

func (m *MonitoredObject) GetObjectName() string {
	if m != nil {
		return m.ObjectName
	}
	return ""
}

func (m *MonitoredObject) GetObjectType() MonitoredObject_MonitoredObjectType {
	if m != nil {
		return m.ObjectType
	}
	return MonitoredObject_MO_UNKNOWN
}

func (m *MonitoredObject) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *MonitoredObject) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

// MonitoredObjectRequest - Wrapper for requests involving a MonitoredObject that are
// scoped to a single Tenant.
type MonitoredObjectRequest struct {
	XId  string           `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string           `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *MonitoredObject `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *MonitoredObjectRequest) Reset()                    { *m = MonitoredObjectRequest{} }
func (m *MonitoredObjectRequest) String() string            { return proto.CompactTextString(m) }
func (*MonitoredObjectRequest) ProtoMessage()               {}
func (*MonitoredObjectRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{15} }

func (m *MonitoredObjectRequest) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *MonitoredObjectRequest) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *MonitoredObjectRequest) GetData() *MonitoredObject {
	if m != nil {
		return m.Data
	}
	return nil
}

// MonitoredObjectRequest - Wrapper for responses involving a MonitoredObject that are
// scoped to a single Tenant.
type MonitoredObjectResponse struct {
	XId  string           `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string           `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *MonitoredObject `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *MonitoredObjectResponse) Reset()                    { *m = MonitoredObjectResponse{} }
func (m *MonitoredObjectResponse) String() string            { return proto.CompactTextString(m) }
func (*MonitoredObjectResponse) ProtoMessage()               {}
func (*MonitoredObjectResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{16} }

func (m *MonitoredObjectResponse) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *MonitoredObjectResponse) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *MonitoredObjectResponse) GetData() *MonitoredObject {
	if m != nil {
		return m.Data
	}
	return nil
}

// MonitoredObjectListResponse - Wrapper for requests which return a list of MonitoredObjects
// scoped to a single Tenant.
type MonitoredObjectListResponse struct {
	Data []*MonitoredObjectResponse `protobuf:"bytes,3,rep,name=data" json:"data,omitempty"`
}

func (m *MonitoredObjectListResponse) Reset()                    { *m = MonitoredObjectListResponse{} }
func (m *MonitoredObjectListResponse) String() string            { return proto.CompactTextString(m) }
func (*MonitoredObjectListResponse) ProtoMessage()               {}
func (*MonitoredObjectListResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{17} }

func (m *MonitoredObjectListResponse) GetData() []*MonitoredObjectResponse {
	if m != nil {
		return m.Data
	}
	return nil
}

type MonitoredObjectIdRequest struct {
	TenantId          string `protobuf:"bytes,1,opt,name=tenantId" json:"tenantId,omitempty"`
	MonitoredObjectId string `protobuf:"bytes,2,opt,name=monitoredObjectId" json:"monitoredObjectId,omitempty"`
}

func (m *MonitoredObjectIdRequest) Reset()                    { *m = MonitoredObjectIdRequest{} }
func (m *MonitoredObjectIdRequest) String() string            { return proto.CompactTextString(m) }
func (*MonitoredObjectIdRequest) ProtoMessage()               {}
func (*MonitoredObjectIdRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{18} }

func (m *MonitoredObjectIdRequest) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *MonitoredObjectIdRequest) GetMonitoredObjectId() string {
	if m != nil {
		return m.MonitoredObjectId
	}
	return ""
}

func init() {
	proto.RegisterType((*TenantDomain)(nil), "gathergrpc.TenantDomain")
	proto.RegisterType((*TenantDomainRequest)(nil), "gathergrpc.TenantDomainRequest")
	proto.RegisterType((*TenantDomainResponse)(nil), "gathergrpc.TenantDomainResponse")
	proto.RegisterType((*TenantDomainListResponse)(nil), "gathergrpc.TenantDomainListResponse")
	proto.RegisterType((*TenantDomainIdRequest)(nil), "gathergrpc.TenantDomainIdRequest")
	proto.RegisterType((*TenantIngestionProfile)(nil), "gathergrpc.TenantIngestionProfile")
	proto.RegisterType((*TenantIngestionProfileRequest)(nil), "gathergrpc.TenantIngestionProfileRequest")
	proto.RegisterType((*TenantIngestionProfileResponse)(nil), "gathergrpc.TenantIngestionProfileResponse")
	proto.RegisterType((*TenantIngestionProfileIdRequest)(nil), "gathergrpc.TenantIngestionProfileIdRequest")
	proto.RegisterType((*TenantUser)(nil), "gathergrpc.TenantUser")
	proto.RegisterType((*TenantUserRequest)(nil), "gathergrpc.TenantUserRequest")
	proto.RegisterType((*TenantUserResponse)(nil), "gathergrpc.TenantUserResponse")
	proto.RegisterType((*TenantUserListResponse)(nil), "gathergrpc.TenantUserListResponse")
	proto.RegisterType((*TenantUserIdRequest)(nil), "gathergrpc.TenantUserIdRequest")
	proto.RegisterType((*MonitoredObject)(nil), "gathergrpc.MonitoredObject")
	proto.RegisterType((*MonitoredObjectRequest)(nil), "gathergrpc.MonitoredObjectRequest")
	proto.RegisterType((*MonitoredObjectResponse)(nil), "gathergrpc.MonitoredObjectResponse")
	proto.RegisterType((*MonitoredObjectListResponse)(nil), "gathergrpc.MonitoredObjectListResponse")
	proto.RegisterType((*MonitoredObjectIdRequest)(nil), "gathergrpc.MonitoredObjectIdRequest")
	proto.RegisterEnum("gathergrpc.MonitoredObject_MonitoredObjectType", MonitoredObject_MonitoredObjectType_name, MonitoredObject_MonitoredObjectType_value)
	proto.RegisterEnum("gathergrpc.MonitoredObject_DeviceType", MonitoredObject_DeviceType_name, MonitoredObject_DeviceType_value)
}

func init() { proto.RegisterFile("gathergrpc/tenantModels.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 880 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x57, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0xc6, 0x71, 0xd2, 0x26, 0xa7, 0xd9, 0x36, 0x99, 0x6e, 0x83, 0xd5, 0x65, 0x4b, 0x64, 0x10,
	0x8a, 0x96, 0x55, 0xba, 0x0a, 0x2b, 0xb8, 0x8e, 0x36, 0x2b, 0x11, 0x68, 0x7e, 0x64, 0xb2, 0xdd,
	0xcb, 0xc8, 0xb5, 0xa7, 0x89, 0x21, 0xf6, 0x04, 0xcf, 0xa4, 0x08, 0xc1, 0x2d, 0x8f, 0xc2, 0x7b,
	0xf0, 0x0c, 0x3c, 0x11, 0x9a, 0xf1, 0x4f, 0x66, 0xec, 0x34, 0xd4, 0x95, 0x76, 0xef, 0x7a, 0xce,
	0xf9, 0xe6, 0xfb, 0xe6, 0x8c, 0xbf, 0x73, 0xaa, 0xc0, 0xf3, 0x85, 0xcd, 0x96, 0x38, 0x5c, 0x84,
	0x6b, 0xe7, 0x92, 0xe1, 0xc0, 0x0e, 0xd8, 0x88, 0xb8, 0x78, 0x45, 0xbb, 0xeb, 0x90, 0x30, 0x82,
	0x60, 0x5b, 0x3e, 0x97, 0xa1, 0x0e, 0xf1, 0x7d, 0x12, 0xc8, 0x50, 0xf3, 0x5f, 0x0d, 0xea, 0x33,
	0xc1, 0x30, 0x20, 0xbe, 0xed, 0x05, 0xe8, 0x1c, 0xaa, 0x11, 0xe3, 0xd0, 0x35, 0xb4, 0xb6, 0xd6,
	0xa9, 0x59, 0x69, 0xcc, 0x6b, 0xae, 0xcd, 0x6c, 0xf6, 0xfb, 0x1a, 0x1b, 0xa5, 0xa8, 0x96, 0xc4,
	0x08, 0x41, 0x39, 0xb0, 0x7d, 0x6c, 0xe8, 0x22, 0x2f, 0xfe, 0x46, 0x4f, 0xa1, 0xe2, 0x90, 0x15,
	0x09, 0x8d, 0xb2, 0x48, 0x46, 0x01, 0x7a, 0x01, 0x0d, 0x27, 0xc4, 0x36, 0xc3, 0xee, 0xcc, 0xf3,
	0x31, 0x65, 0xb6, 0xbf, 0x36, 0x2a, 0x6d, 0xad, 0xa3, 0x5b, 0xb9, 0x3c, 0x7a, 0x0d, 0x67, 0x2b,
	0x9b, 0xf2, 0xee, 0xbc, 0x5b, 0x4f, 0x3e, 0x70, 0x20, 0x0e, 0xec, 0x2e, 0x9a, 0x0b, 0x38, 0x95,
	0x7b, 0xb2, 0xf0, 0xaf, 0x1b, 0x4c, 0x19, 0x3a, 0x01, 0x7d, 0xee, 0x25, 0x5d, 0x95, 0x86, 0x2e,
	0x6a, 0x42, 0x79, 0x1e, 0xe2, 0xbb, 0xb8, 0x17, 0xdd, 0xc2, 0x77, 0xe8, 0x25, 0x94, 0x79, 0x4b,
	0xa2, 0x8d, 0xa3, 0x9e, 0xd1, 0xdd, 0xbe, 0x5e, 0x57, 0xa1, 0x14, 0x28, 0x73, 0x09, 0x4f, 0x55,
	0x21, 0xba, 0x26, 0x01, 0xc5, 0x1f, 0x40, 0x69, 0x0a, 0x86, 0x9c, 0xbd, 0xf2, 0x28, 0x4b, 0xd5,
	0x5e, 0xc7, 0x4c, 0x5a, 0x5b, 0xef, 0x1c, 0xf5, 0xda, 0xf7, 0x32, 0xc5, 0xf8, 0x98, 0x71, 0x02,
	0x67, 0x72, 0x75, 0xe8, 0x26, 0xcf, 0xf4, 0x7f, 0x0e, 0x88, 0xe1, 0xa9, 0x03, 0xe2, 0xd8, 0xfc,
	0xa7, 0x04, 0xad, 0x88, 0x71, 0x18, 0x2c, 0x30, 0x65, 0x1e, 0x09, 0xa6, 0x21, 0xb9, 0xf5, 0x56,
	0xf8, 0xd1, 0xa6, 0x6a, 0xc3, 0x11, 0x75, 0xd6, 0xef, 0x28, 0x0e, 0x25, 0x6f, 0xc9, 0xa9, 0x18,
	0x31, 0xb5, 0x29, 0xfd, 0x8d, 0x84, 0x6e, 0x6c, 0x34, 0x39, 0xf5, 0xe1, 0xed, 0x86, 0xbe, 0x87,
	0x06, 0x5b, 0x86, 0x98, 0x2e, 0xc9, 0xca, 0x8d, 0x3b, 0x36, 0x0e, 0xc5, 0xb7, 0xf8, 0x4c, 0xf9,
	0x16, 0x19, 0x8c, 0x95, 0x3b, 0x65, 0xfe, 0x01, 0xcf, 0x77, 0xbf, 0x60, 0x11, 0x0b, 0x7f, 0xab,
	0x18, 0xcb, 0xcc, 0xdb, 0x21, 0x47, 0x1e, 0x19, 0xe2, 0x4f, 0xb8, 0xb8, 0x4f, 0xbc, 0x80, 0xad,
	0x1f, 0xab, 0xee, 0xc3, 0xe7, 0xbb, 0xeb, 0x0f, 0x33, 0x66, 0x17, 0x90, 0x97, 0x3b, 0x18, 0xdf,
	0x6b, 0x47, 0xc5, 0xfc, 0x5b, 0x07, 0x88, 0xf4, 0xb8, 0x95, 0x1e, 0x6d, 0xd0, 0x73, 0xa8, 0x6e,
	0x54, 0x77, 0xa6, 0x31, 0xaf, 0xad, 0x55, 0x5f, 0xa6, 0x31, 0x7a, 0x05, 0xa7, 0x14, 0x07, 0xee,
	0x24, 0xb8, 0x21, 0x76, 0xe8, 0x7a, 0xc1, 0xe2, 0xad, 0x6f, 0x7b, 0x2b, 0xe1, 0xcb, 0xaa, 0xb5,
	0xab, 0x84, 0x3a, 0x70, 0x42, 0xd2, 0xd4, 0x8c, 0xfc, 0x82, 0x03, 0x61, 0xca, 0x9a, 0x95, 0x4d,
	0x23, 0x13, 0xea, 0xfc, 0x0e, 0xd7, 0x38, 0x14, 0x3e, 0x35, 0x0e, 0x05, 0xa9, 0x92, 0x43, 0x5f,
	0x43, 0x85, 0x32, 0x9b, 0x61, 0xa3, 0xda, 0xd6, 0x3a, 0xc7, 0xbd, 0x33, 0xf9, 0x33, 0xf1, 0x07,
	0xf9, 0x89, 0x17, 0xad, 0x08, 0x83, 0x0c, 0x38, 0x8c, 0x86, 0x9c, 0x1a, 0xb5, 0xb6, 0xde, 0xa9,
	0x59, 0x49, 0xb8, 0x73, 0xb6, 0xa0, 0xe8, 0x6c, 0x1d, 0xed, 0x5b, 0xe5, 0x0e, 0x34, 0xb7, 0x9f,
	0xa9, 0xc8, 0x14, 0xbc, 0x50, 0x7c, 0xd8, 0xca, 0xfb, 0x50, 0x10, 0x46, 0xde, 0x73, 0x01, 0xc9,
	0x22, 0x05, 0xdc, 0x5e, 0x44, 0xe5, 0x2a, 0x59, 0x8f, 0x3c, 0xa7, 0x2c, 0xf0, 0x9e, 0xb2, 0xc0,
	0x2f, 0xee, 0x61, 0x51, 0xd7, 0xf7, 0x30, 0xf9, 0x1f, 0xc7, 0x6b, 0x0f, 0x9b, 0x91, 0x16, 0x1c,
	0x6c, 0x04, 0x38, 0xee, 0x20, 0x8e, 0xcc, 0xbf, 0x2a, 0x70, 0x32, 0x22, 0x81, 0xc7, 0x48, 0x88,
	0xdd, 0xc9, 0xcd, 0xcf, 0xd8, 0x61, 0xe8, 0x18, 0x4a, 0xdb, 0xde, 0x3d, 0x57, 0xe1, 0x2d, 0xed,
	0x19, 0x10, 0x3d, 0x33, 0x20, 0x3f, 0x40, 0xdd, 0x76, 0xd8, 0xc6, 0x66, 0x24, 0x9c, 0xf1, 0x7a,
	0x59, 0xf8, 0xed, 0x2b, 0xb9, 0xc5, 0x8c, 0x74, 0x77, 0x80, 0xef, 0x3c, 0x07, 0x73, 0xb4, 0xa5,
	0x9c, 0xe5, 0xc6, 0x4e, 0xe2, 0x31, 0x1f, 0xb8, 0x8a, 0xd0, 0x52, 0x72, 0xe8, 0x0a, 0x9e, 0x84,
	0xf8, 0x76, 0x85, 0x9d, 0x44, 0xf0, 0xa0, 0x90, 0xa0, 0x7a, 0x18, 0x7d, 0x29, 0xb1, 0x09, 0xc9,
	0x43, 0x21, 0xa9, 0x26, 0xd1, 0x05, 0x00, 0x11, 0x4c, 0x02, 0x52, 0x15, 0x10, 0x29, 0x83, 0x26,
	0x49, 0x5d, 0x5c, 0xa8, 0x26, 0x2e, 0x74, 0xb9, 0xef, 0x42, 0x99, 0x58, 0xdc, 0x4c, 0xa2, 0xf8,
	0x08, 0x63, 0xf7, 0x0a, 0x4e, 0x77, 0x5c, 0x02, 0x1d, 0x03, 0x8c, 0x26, 0xf3, 0x77, 0xe3, 0x1f,
	0xc7, 0x93, 0xf7, 0xe3, 0xc6, 0x27, 0xa8, 0x06, 0x95, 0xd9, 0xfb, 0xfe, 0x68, 0xda, 0xd0, 0xcc,
	0x3e, 0xc0, 0xf6, 0x1d, 0x39, 0x70, 0x30, 0x93, 0x80, 0x0d, 0xa8, 0xf7, 0xdf, 0xbc, 0x79, 0x3b,
	0x18, 0xf6, 0xc7, 0xf3, 0xf1, 0x70, 0xd0, 0xd0, 0x50, 0x13, 0x9e, 0xa4, 0x99, 0x6b, 0x9e, 0x2a,
	0x99, 0x3e, 0xb4, 0x32, 0xa2, 0x45, 0x06, 0xfe, 0x52, 0x19, 0xc5, 0x67, 0x7b, 0xde, 0x37, 0x9e,
	0xa0, 0x00, 0x3e, 0xcd, 0xc9, 0x15, 0x18, 0xfd, 0xc2, 0x7a, 0xd7, 0xf0, 0x2c, 0x53, 0x50, 0x96,
	0xc0, 0x77, 0x29, 0x1f, 0x5f, 0x02, 0x5f, 0xec, 0xe3, 0x53, 0x37, 0x81, 0x0b, 0x46, 0x06, 0xf0,
	0xb0, 0x75, 0xf0, 0x12, 0x9a, 0x7e, 0xf6, 0x5c, 0xdc, 0x60, 0xbe, 0x70, 0x73, 0x20, 0x7e, 0x2f,
	0x7c, 0xf3, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x57, 0xe5, 0x3b, 0x9b, 0x7b, 0x0c, 0x00, 0x00,
}
