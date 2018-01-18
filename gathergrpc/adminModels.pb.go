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
type TenantDescriptorData struct {
	Datatype              string    `protobuf:"bytes,2,opt,name=datatype" json:"datatype,omitempty"`
	Name                  string    `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	UrlSubdomain          string    `protobuf:"bytes,4,opt,name=urlSubdomain" json:"urlSubdomain,omitempty"`
	State                 UserState `protobuf:"varint,5,opt,name=state,enum=gathergrpc.UserState" json:"state,omitempty"`
	CreatedTimestamp      int64     `protobuf:"varint,6,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64     `protobuf:"varint,7,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *TenantDescriptorData) Reset()                    { *m = TenantDescriptorData{} }
func (m *TenantDescriptorData) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptorData) ProtoMessage()               {}
func (*TenantDescriptorData) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *TenantDescriptorData) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *TenantDescriptorData) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TenantDescriptorData) GetUrlSubdomain() string {
	if m != nil {
		return m.UrlSubdomain
	}
	return ""
}

func (m *TenantDescriptorData) GetState() UserState {
	if m != nil {
		return m.State
	}
	return UserState_USER_UNKNOWN
}

func (m *TenantDescriptorData) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *TenantDescriptorData) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

// TenantDescriptor - wrapper for passing TenantDescriptor
// data as a request to the service.
type TenantDescriptor struct {
	XId  string                `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string                `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *TenantDescriptorData `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *TenantDescriptor) Reset()                    { *m = TenantDescriptor{} }
func (m *TenantDescriptor) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptor) ProtoMessage()               {}
func (*TenantDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *TenantDescriptor) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *TenantDescriptor) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *TenantDescriptor) GetData() *TenantDescriptorData {
	if m != nil {
		return m.Data
	}
	return nil
}

// Wrapper message to provide a response in the form of
// a container of multiple TenantDescriptor objects.
type TenantDescriptorList struct {
	Data []*TenantDescriptor `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *TenantDescriptorList) Reset()                    { *m = TenantDescriptorList{} }
func (m *TenantDescriptorList) String() string            { return proto.CompactTextString(m) }
func (*TenantDescriptorList) ProtoMessage()               {}
func (*TenantDescriptorList) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *TenantDescriptorList) GetData() []*TenantDescriptor {
	if m != nil {
		return m.Data
	}
	return nil
}

// User data for an Adminstrative User.
type AdminUserData struct {
	Datatype              string    `protobuf:"bytes,2,opt,name=datatype" json:"datatype,omitempty"`
	Username              string    `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	Password              string    `protobuf:"bytes,4,opt,name=password" json:"password,omitempty"`
	SendOnboardingEmail   bool      `protobuf:"varint,5,opt,name=sendOnboardingEmail" json:"sendOnboardingEmail,omitempty"`
	OnboardingToken       string    `protobuf:"bytes,6,opt,name=onboardingToken" json:"onboardingToken,omitempty"`
	UserVerified          bool      `protobuf:"varint,7,opt,name=userVerified" json:"userVerified,omitempty"`
	State                 UserState `protobuf:"varint,8,opt,name=state,enum=gathergrpc.UserState" json:"state,omitempty"`
	CreatedTimestamp      int64     `protobuf:"varint,9,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64     `protobuf:"varint,10,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *AdminUserData) Reset()                    { *m = AdminUserData{} }
func (m *AdminUserData) String() string            { return proto.CompactTextString(m) }
func (*AdminUserData) ProtoMessage()               {}
func (*AdminUserData) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *AdminUserData) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *AdminUserData) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AdminUserData) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *AdminUserData) GetSendOnboardingEmail() bool {
	if m != nil {
		return m.SendOnboardingEmail
	}
	return false
}

func (m *AdminUserData) GetOnboardingToken() string {
	if m != nil {
		return m.OnboardingToken
	}
	return ""
}

func (m *AdminUserData) GetUserVerified() bool {
	if m != nil {
		return m.UserVerified
	}
	return false
}

func (m *AdminUserData) GetState() UserState {
	if m != nil {
		return m.State
	}
	return UserState_USER_UNKNOWN
}

func (m *AdminUserData) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *AdminUserData) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

// AdminUser - wrapper for passing AdminUser
// data as a request to the service.
type AdminUser struct {
	XId  string         `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string         `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *AdminUserData `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *AdminUser) Reset()                    { *m = AdminUser{} }
func (m *AdminUser) String() string            { return proto.CompactTextString(m) }
func (*AdminUser) ProtoMessage()               {}
func (*AdminUser) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *AdminUser) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *AdminUser) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *AdminUser) GetData() *AdminUserData {
	if m != nil {
		return m.Data
	}
	return nil
}

// Wrapper message to provide a response in the form of
// a container of multiple AdminUser objects.
type AdminUserList struct {
	Data []*AdminUser `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *AdminUserList) Reset()                    { *m = AdminUserList{} }
func (m *AdminUserList) String() string            { return proto.CompactTextString(m) }
func (*AdminUserList) ProtoMessage()               {}
func (*AdminUserList) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *AdminUserList) GetData() []*AdminUser {
	if m != nil {
		return m.Data
	}
	return nil
}

// Stores the available values of the metrics that may be ingested by ADH.
type IngestionDictionaryData struct {
	Datatype              string                                        `protobuf:"bytes,2,opt,name=datatype" json:"datatype,omitempty"`
	Metrics               map[string]*IngestionDictionaryData_MetricMap `protobuf:"bytes,3,rep,name=metrics" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	CreatedTimestamp      int64                                         `protobuf:"varint,4,opt,name=createdTimestamp" json:"createdTimestamp,omitempty"`
	LastModifiedTimestamp int64                                         `protobuf:"varint,5,opt,name=lastModifiedTimestamp" json:"lastModifiedTimestamp,omitempty"`
}

func (m *IngestionDictionaryData) Reset()                    { *m = IngestionDictionaryData{} }
func (m *IngestionDictionaryData) String() string            { return proto.CompactTextString(m) }
func (*IngestionDictionaryData) ProtoMessage()               {}
func (*IngestionDictionaryData) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *IngestionDictionaryData) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

func (m *IngestionDictionaryData) GetMetrics() map[string]*IngestionDictionaryData_MetricMap {
	if m != nil {
		return m.Metrics
	}
	return nil
}

func (m *IngestionDictionaryData) GetCreatedTimestamp() int64 {
	if m != nil {
		return m.CreatedTimestamp
	}
	return 0
}

func (m *IngestionDictionaryData) GetLastModifiedTimestamp() int64 {
	if m != nil {
		return m.LastModifiedTimestamp
	}
	return 0
}

type IngestionDictionaryData_UIData struct {
	Group    string `protobuf:"bytes,1,opt,name=group" json:"group,omitempty"`
	Position string `protobuf:"bytes,2,opt,name=position" json:"position,omitempty"`
}

func (m *IngestionDictionaryData_UIData) Reset()         { *m = IngestionDictionaryData_UIData{} }
func (m *IngestionDictionaryData_UIData) String() string { return proto.CompactTextString(m) }
func (*IngestionDictionaryData_UIData) ProtoMessage()    {}
func (*IngestionDictionaryData_UIData) Descriptor() ([]byte, []int) {
	return fileDescriptor1, []int{6, 0}
}

func (m *IngestionDictionaryData_UIData) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *IngestionDictionaryData_UIData) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

type IngestionDictionaryData_MonitoredObjectType struct {
	Key         string   `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	RawMetricId string   `protobuf:"bytes,2,opt,name=rawMetricId" json:"rawMetricId,omitempty"`
	Unit        string   `protobuf:"bytes,3,opt,name=unit" json:"unit,omitempty"`
	Directions  []string `protobuf:"bytes,4,rep,name=directions" json:"directions,omitempty"`
}

func (m *IngestionDictionaryData_MonitoredObjectType) Reset() {
	*m = IngestionDictionaryData_MonitoredObjectType{}
}
func (m *IngestionDictionaryData_MonitoredObjectType) String() string {
	return proto.CompactTextString(m)
}
func (*IngestionDictionaryData_MonitoredObjectType) ProtoMessage() {}
func (*IngestionDictionaryData_MonitoredObjectType) Descriptor() ([]byte, []int) {
	return fileDescriptor1, []int{6, 1}
}

func (m *IngestionDictionaryData_MonitoredObjectType) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *IngestionDictionaryData_MonitoredObjectType) GetRawMetricId() string {
	if m != nil {
		return m.RawMetricId
	}
	return ""
}

func (m *IngestionDictionaryData_MonitoredObjectType) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *IngestionDictionaryData_MonitoredObjectType) GetDirections() []string {
	if m != nil {
		return m.Directions
	}
	return nil
}

type IngestionDictionaryData_MetricDefinition struct {
	MonitoredObjectTypes []*IngestionDictionaryData_MonitoredObjectType `protobuf:"bytes,1,rep,name=monitoredObjectTypes" json:"monitoredObjectTypes,omitempty"`
	Ui                   *IngestionDictionaryData_UIData                `protobuf:"bytes,2,opt,name=ui" json:"ui,omitempty"`
}

func (m *IngestionDictionaryData_MetricDefinition) Reset() {
	*m = IngestionDictionaryData_MetricDefinition{}
}
func (m *IngestionDictionaryData_MetricDefinition) String() string { return proto.CompactTextString(m) }
func (*IngestionDictionaryData_MetricDefinition) ProtoMessage()    {}
func (*IngestionDictionaryData_MetricDefinition) Descriptor() ([]byte, []int) {
	return fileDescriptor1, []int{6, 2}
}

func (m *IngestionDictionaryData_MetricDefinition) GetMonitoredObjectTypes() []*IngestionDictionaryData_MonitoredObjectType {
	if m != nil {
		return m.MonitoredObjectTypes
	}
	return nil
}

func (m *IngestionDictionaryData_MetricDefinition) GetUi() *IngestionDictionaryData_UIData {
	if m != nil {
		return m.Ui
	}
	return nil
}

type IngestionDictionaryData_MetricMap struct {
	MetricMap map[string]*IngestionDictionaryData_MetricDefinition `protobuf:"bytes,1,rep,name=metricMap" json:"metricMap,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *IngestionDictionaryData_MetricMap) Reset()         { *m = IngestionDictionaryData_MetricMap{} }
func (m *IngestionDictionaryData_MetricMap) String() string { return proto.CompactTextString(m) }
func (*IngestionDictionaryData_MetricMap) ProtoMessage()    {}
func (*IngestionDictionaryData_MetricMap) Descriptor() ([]byte, []int) {
	return fileDescriptor1, []int{6, 3}
}

func (m *IngestionDictionaryData_MetricMap) GetMetricMap() map[string]*IngestionDictionaryData_MetricDefinition {
	if m != nil {
		return m.MetricMap
	}
	return nil
}

// Wrapper fos the object that stores the available values of the metrics that may be ingested by ADH.
type IngestionDictionary struct {
	XId  string                   `protobuf:"bytes,1,opt,name=_id,json=Id" json:"_id,omitempty"`
	XRev string                   `protobuf:"bytes,2,opt,name=_rev,json=Rev" json:"_rev,omitempty"`
	Data *IngestionDictionaryData `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *IngestionDictionary) Reset()                    { *m = IngestionDictionary{} }
func (m *IngestionDictionary) String() string            { return proto.CompactTextString(m) }
func (*IngestionDictionary) ProtoMessage()               {}
func (*IngestionDictionary) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *IngestionDictionary) GetXId() string {
	if m != nil {
		return m.XId
	}
	return ""
}

func (m *IngestionDictionary) GetXRev() string {
	if m != nil {
		return m.XRev
	}
	return ""
}

func (m *IngestionDictionary) GetData() *IngestionDictionaryData {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*TenantDescriptorData)(nil), "gathergrpc.TenantDescriptorData")
	proto.RegisterType((*TenantDescriptor)(nil), "gathergrpc.TenantDescriptor")
	proto.RegisterType((*TenantDescriptorList)(nil), "gathergrpc.TenantDescriptorList")
	proto.RegisterType((*AdminUserData)(nil), "gathergrpc.AdminUserData")
	proto.RegisterType((*AdminUser)(nil), "gathergrpc.AdminUser")
	proto.RegisterType((*AdminUserList)(nil), "gathergrpc.AdminUserList")
	proto.RegisterType((*IngestionDictionaryData)(nil), "gathergrpc.IngestionDictionaryData")
	proto.RegisterType((*IngestionDictionaryData_UIData)(nil), "gathergrpc.IngestionDictionaryData.UIData")
	proto.RegisterType((*IngestionDictionaryData_MonitoredObjectType)(nil), "gathergrpc.IngestionDictionaryData.MonitoredObjectType")
	proto.RegisterType((*IngestionDictionaryData_MetricDefinition)(nil), "gathergrpc.IngestionDictionaryData.MetricDefinition")
	proto.RegisterType((*IngestionDictionaryData_MetricMap)(nil), "gathergrpc.IngestionDictionaryData.MetricMap")
	proto.RegisterType((*IngestionDictionary)(nil), "gathergrpc.IngestionDictionary")
}

func init() { proto.RegisterFile("gathergrpc/adminModels.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 731 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0xcd, 0x4e, 0xdb, 0x4a,
	0x14, 0x96, 0xe3, 0x04, 0x92, 0x13, 0x2e, 0xe4, 0x0e, 0xa0, 0xeb, 0x6b, 0x71, 0xaf, 0xa2, 0x74,
	0x93, 0x52, 0x91, 0xa2, 0x14, 0x89, 0x2a, 0xea, 0xa6, 0x6a, 0x90, 0x1a, 0xd4, 0x08, 0xc9, 0x84,
	0x2e, 0xba, 0x28, 0x9a, 0x64, 0x86, 0x74, 0x4a, 0x3c, 0x63, 0xcd, 0x8c, 0x41, 0x59, 0xf4, 0x71,
	0xfa, 0x14, 0x7d, 0x94, 0xbe, 0x45, 0x77, 0xdd, 0x55, 0x1e, 0x3b, 0x89, 0x4d, 0x5c, 0x08, 0x5d,
	0x79, 0xce, 0x39, 0x73, 0xbe, 0xf3, 0xf3, 0x7d, 0x1a, 0xc3, 0xde, 0x18, 0xeb, 0x4f, 0x54, 0x8e,
	0x65, 0x30, 0x7a, 0x8e, 0x89, 0xcf, 0x78, 0x5f, 0x10, 0x3a, 0x51, 0xad, 0x40, 0x0a, 0x2d, 0x10,
	0x2c, 0xa2, 0xee, 0x7f, 0xa9, 0x9b, 0x23, 0xe1, 0xfb, 0x22, 0x73, 0xb5, 0xf1, 0xd3, 0x82, 0x9d,
	0x01, 0xe5, 0x98, 0xeb, 0x2e, 0x55, 0x23, 0xc9, 0x02, 0x2d, 0x64, 0x17, 0x6b, 0x8c, 0x5c, 0x28,
	0x13, 0xac, 0xb1, 0x9e, 0x06, 0xd4, 0x29, 0xd4, 0xad, 0x66, 0xc5, 0x9b, 0xdb, 0x08, 0x41, 0x91,
	0x63, 0x9f, 0x3a, 0xb6, 0xf1, 0x9b, 0x33, 0x6a, 0xc0, 0x46, 0x28, 0x27, 0xe7, 0xe1, 0x90, 0x08,
	0x1f, 0x33, 0xee, 0x14, 0x4d, 0x2c, 0xe3, 0x43, 0xcf, 0xa0, 0xa4, 0x34, 0xd6, 0xd4, 0x29, 0xd5,
	0xad, 0xe6, 0x66, 0x7b, 0xb7, 0xb5, 0xe8, 0xad, 0x75, 0xa1, 0xa8, 0x3c, 0x8f, 0x82, 0x5e, 0x7c,
	0x07, 0xed, 0x43, 0x6d, 0x24, 0x29, 0xd6, 0x94, 0x0c, 0x98, 0x4f, 0x95, 0xc6, 0x7e, 0xe0, 0xac,
	0xd5, 0xad, 0xa6, 0xed, 0x2d, 0xf9, 0xd1, 0x11, 0xec, 0x4e, 0xb0, 0xd2, 0x7d, 0x41, 0xd8, 0x15,
	0x4b, 0x27, 0xac, 0x9b, 0x84, 0xfc, 0x60, 0x63, 0x02, 0xb5, 0xbb, 0xa3, 0xa3, 0x2d, 0xb0, 0x2f,
	0x19, 0x71, 0x2c, 0xd3, 0x7d, 0xa1, 0x47, 0xd0, 0xdf, 0x50, 0xbc, 0x94, 0xf4, 0x26, 0xd9, 0x81,
	0xed, 0xd1, 0x1b, 0x74, 0x04, 0xc5, 0x68, 0x15, 0x66, 0xfc, 0x6a, 0xbb, 0x9e, 0x9e, 0x22, 0x6f,
	0x95, 0x9e, 0xb9, 0xdd, 0x78, 0xbb, 0xbc, 0xe8, 0x77, 0x4c, 0x69, 0x74, 0x98, 0xa0, 0x59, 0x75,
	0xbb, 0x59, 0x6d, 0xef, 0xdd, 0x87, 0x96, 0x20, 0xfd, 0x28, 0xc0, 0x5f, 0xaf, 0x23, 0xd2, 0xa3,
	0x9d, 0x3d, 0x48, 0x96, 0x0b, 0xe5, 0x50, 0x51, 0x99, 0x22, 0x6c, 0x6e, 0x47, 0xb1, 0x00, 0x2b,
	0x75, 0x2b, 0x24, 0x49, 0x08, 0x9b, 0xdb, 0xe8, 0x10, 0xb6, 0x15, 0xe5, 0xe4, 0x8c, 0x0f, 0x05,
	0x96, 0x84, 0xf1, 0xf1, 0x89, 0x8f, 0xd9, 0xc4, 0x50, 0x57, 0xf6, 0xf2, 0x42, 0xa8, 0x09, 0x5b,
	0x62, 0xee, 0x1a, 0x88, 0x6b, 0xca, 0x0d, 0x61, 0x15, 0xef, 0xae, 0xdb, 0x88, 0x45, 0x51, 0xf9,
	0x9e, 0x4a, 0x43, 0x89, 0xa1, 0xa9, 0xec, 0x65, 0x7c, 0x0b, 0xb1, 0x94, 0xff, 0x50, 0x2c, 0x95,
	0xc7, 0x8a, 0x05, 0xee, 0x13, 0xcb, 0x47, 0xa8, 0xcc, 0x77, 0xbe, 0x92, 0x4a, 0x0e, 0x32, 0x2a,
	0xf9, 0x37, 0xdd, 0x7e, 0x86, 0xbc, 0x84, 0xd4, 0x4e, 0x8a, 0x53, 0xa3, 0x8b, 0xa7, 0x19, 0x5d,
	0xec, 0xe6, 0xe6, 0x27, 0xb9, 0x5f, 0xd7, 0xe1, 0x9f, 0x1e, 0x1f, 0x53, 0xa5, 0x99, 0xe0, 0x5d,
	0x36, 0x8a, 0x3e, 0x58, 0x4e, 0x1f, 0x94, 0xc6, 0x29, 0xac, 0xfb, 0x54, 0x4b, 0x36, 0x52, 0x8e,
	0x6d, 0xaa, 0x1c, 0xa6, 0xab, 0xfc, 0x06, 0xb1, 0xd5, 0x8f, 0x53, 0x4e, 0xb8, 0x96, 0x53, 0x6f,
	0x06, 0x90, 0xcb, 0x40, 0xf1, 0xb1, 0x0c, 0x94, 0xee, 0x61, 0xc0, 0xed, 0xc0, 0xda, 0x45, 0xcf,
	0xcc, 0xb4, 0x03, 0xa5, 0xb1, 0x14, 0x61, 0x90, 0x10, 0x10, 0x1b, 0x46, 0xcc, 0x42, 0xb1, 0xa8,
	0xd5, 0xd9, 0xa4, 0x33, 0xdb, 0xfd, 0x02, 0xdb, 0x7d, 0xc1, 0x99, 0x16, 0x92, 0x92, 0xb3, 0xe1,
	0x67, 0x3a, 0xd2, 0x83, 0x68, 0x01, 0x35, 0xb0, 0xaf, 0xe9, 0x34, 0x81, 0x89, 0x8e, 0xa8, 0x0e,
	0x55, 0x89, 0x6f, 0xe3, 0x11, 0x7b, 0x24, 0xc1, 0x49, 0xbb, 0xa2, 0xc7, 0x2f, 0xe4, 0x4c, 0xcf,
	0x1e, 0xbf, 0xe8, 0x8c, 0xfe, 0x07, 0x20, 0x4c, 0x52, 0xb3, 0x26, 0xe5, 0x14, 0xeb, 0x76, 0xb3,
	0xe2, 0xa5, 0x3c, 0xee, 0x37, 0x0b, 0x6a, 0x31, 0x40, 0x97, 0x5e, 0x31, 0x6e, 0x7a, 0x42, 0xd7,
	0xb0, 0xe3, 0x2f, 0xf7, 0xa4, 0x12, 0xc2, 0x8f, 0x57, 0xa2, 0x62, 0x39, 0xdf, 0xcb, 0x05, 0x45,
	0x1d, 0x28, 0x84, 0xcc, 0x8c, 0x53, 0x6d, 0xef, 0xaf, 0x02, 0x1d, 0xaf, 0xda, 0x2b, 0x84, 0xcc,
	0xfd, 0x6e, 0x41, 0x25, 0xee, 0xbe, 0x8f, 0x03, 0xf4, 0x01, 0x2a, 0xfe, 0xcc, 0x48, 0x7a, 0x7d,
	0xb5, 0xba, 0x6c, 0xfa, 0x38, 0x58, 0x9c, 0x62, 0x09, 0x2d, 0xe0, 0x5c, 0x09, 0x9b, 0xd9, 0x60,
	0x0e, 0x43, 0xa7, 0x50, 0xba, 0xc1, 0x93, 0x90, 0x26, 0xc3, 0x1c, 0xad, 0x5e, 0x7b, 0xb1, 0x7b,
	0x2f, 0x86, 0xe8, 0x14, 0x5e, 0x5a, 0x2e, 0x83, 0x8d, 0xb4, 0xa2, 0x73, 0x2a, 0xbe, 0xc9, 0x56,
	0x3c, 0x78, 0xd4, 0xb4, 0xa9, 0x52, 0x0d, 0x09, 0xdb, 0x39, 0xf7, 0x57, 0x7a, 0x4d, 0x8e, 0x33,
	0xaf, 0xc9, 0x93, 0x15, 0x5a, 0x88, 0xdf, 0x86, 0xe1, 0x9a, 0xf9, 0xcf, 0xbf, 0xf8, 0x15, 0x00,
	0x00, 0xff, 0xff, 0x91, 0x98, 0x13, 0x90, 0x32, 0x08, 0x00, 0x00,
}
