// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gathergrpc/commonModels.proto

/*
Package gathergrpc is a generated protocol buffer package.

It is generated from these files:
	gathergrpc/commonModels.proto
	gathergrpc/adminModels.proto
	gathergrpc/tenantModels.proto
	gathergrpc/metricModels.proto
	gathergrpc/gather.proto

It has these top-level messages:
	JSONAPIObject
	Data
	Error
	Links
	Resource
	Relationships
	TenantDescriptor
	TenantDescriptorRequest
	TenantDescriptorResponse
	TenantDescriptorListResponse
	AdminUser
	AdminUserRequest
	AdminUserResponse
	AdminUserListResponse
	IngestionDictionaryData
	IngestionDictionary
	TenantDomain
	TenantDomainRequest
	TenantDomainResponse
	TenantDomainListResponse
	TenantDomainIdRequest
	TenantIngestionProfile
	TenantIngestionProfileRequest
	TenantIngestionProfileResponse
	TenantIngestionProfileIdRequest
	TenantThresholdProfile
	TenantThresholdProfileRequest
	TenantThresholdProfileResponse
	TenantThresholdProfileIdRequest
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
	MonitoredObjectCountByDomainRequest
	MonitoredObjectCountByDomainResponse
	MonitoredObjectList
	ThresholdCrossing
	FormattedThresholdCrossing
	ThresholdCrossingResponse
	ThresholdCrossingRequest
	HistogramRequest
	HistogramBuckets
	HistogramResult
	Histogram
	HistogramResponse
	ThresholdCrossingByMonitoredObjectResponse
*/
package gathergrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/any"

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

type JSONAPIObject struct {
	Data     []*Data           `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
	Errors   []*Error          `protobuf:"bytes,2,rep,name=errors" json:"errors,omitempty"`
	Metadata map[string]string `protobuf:"bytes,3,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Jsonapi  map[string]string `protobuf:"bytes,4,rep,name=jsonapi" json:"jsonapi,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Links    *Links            `protobuf:"bytes,5,opt,name=links" json:"links,omitempty"`
	Included []*Resource       `protobuf:"bytes,6,rep,name=included" json:"included,omitempty"`
}

func (m *JSONAPIObject) Reset()                    { *m = JSONAPIObject{} }
func (m *JSONAPIObject) String() string            { return proto.CompactTextString(m) }
func (*JSONAPIObject) ProtoMessage()               {}
func (*JSONAPIObject) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *JSONAPIObject) GetData() []*Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *JSONAPIObject) GetErrors() []*Error {
	if m != nil {
		return m.Errors
	}
	return nil
}

func (m *JSONAPIObject) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *JSONAPIObject) GetJsonapi() map[string]string {
	if m != nil {
		return m.Jsonapi
	}
	return nil
}

func (m *JSONAPIObject) GetLinks() *Links {
	if m != nil {
		return m.Links
	}
	return nil
}

func (m *JSONAPIObject) GetIncluded() []*Resource {
	if m != nil {
		return m.Included
	}
	return nil
}

type Data struct {
	Id         string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type       string               `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Attributes *google_protobuf.Any `protobuf:"bytes,3,opt,name=attributes" json:"attributes,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Data) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Data) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Data) GetAttributes() *google_protobuf.Any {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type Error struct {
	Id       string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Links    *Links            `protobuf:"bytes,2,opt,name=links" json:"links,omitempty"`
	Status   string            `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
	Code     string            `protobuf:"bytes,4,opt,name=code" json:"code,omitempty"`
	Title    string            `protobuf:"bytes,5,opt,name=title" json:"title,omitempty"`
	Detail   string            `protobuf:"bytes,6,opt,name=detail" json:"detail,omitempty"`
	Metadata map[string]string `protobuf:"bytes,7,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Source   map[string]string `protobuf:"bytes,8,rep,name=source" json:"source,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Error) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Error) GetLinks() *Links {
	if m != nil {
		return m.Links
	}
	return nil
}

func (m *Error) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Error) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Error) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Error) GetDetail() string {
	if m != nil {
		return m.Detail
	}
	return ""
}

func (m *Error) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Error) GetSource() map[string]string {
	if m != nil {
		return m.Source
	}
	return nil
}

// Links technically allow any valid json string as a key,
// we obviously can't do this, so we have to stick to the following.
type Links struct {
	Related *Links_Related `protobuf:"bytes,1,opt,name=related" json:"related,omitempty"`
	Self    string         `protobuf:"bytes,2,opt,name=self" json:"self,omitempty"`
	First   string         `protobuf:"bytes,3,opt,name=first" json:"first,omitempty"`
	Next    string         `protobuf:"bytes,4,opt,name=next" json:"next,omitempty"`
	Prev    string         `protobuf:"bytes,5,opt,name=prev" json:"prev,omitempty"`
	Last    string         `protobuf:"bytes,6,opt,name=last" json:"last,omitempty"`
	About   string         `protobuf:"bytes,7,opt,name=about" json:"about,omitempty"`
	Article string         `protobuf:"bytes,8,opt,name=article" json:"article,omitempty"`
}

func (m *Links) Reset()                    { *m = Links{} }
func (m *Links) String() string            { return proto.CompactTextString(m) }
func (*Links) ProtoMessage()               {}
func (*Links) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Links) GetRelated() *Links_Related {
	if m != nil {
		return m.Related
	}
	return nil
}

func (m *Links) GetSelf() string {
	if m != nil {
		return m.Self
	}
	return ""
}

func (m *Links) GetFirst() string {
	if m != nil {
		return m.First
	}
	return ""
}

func (m *Links) GetNext() string {
	if m != nil {
		return m.Next
	}
	return ""
}

func (m *Links) GetPrev() string {
	if m != nil {
		return m.Prev
	}
	return ""
}

func (m *Links) GetLast() string {
	if m != nil {
		return m.Last
	}
	return ""
}

func (m *Links) GetAbout() string {
	if m != nil {
		return m.About
	}
	return ""
}

func (m *Links) GetArticle() string {
	if m != nil {
		return m.Article
	}
	return ""
}

type Links_Related struct {
	Href     string            `protobuf:"bytes,1,opt,name=href" json:"href,omitempty"`
	Metadata map[string]string `protobuf:"bytes,2,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Links_Related) Reset()                    { *m = Links_Related{} }
func (m *Links_Related) String() string            { return proto.CompactTextString(m) }
func (*Links_Related) ProtoMessage()               {}
func (*Links_Related) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

func (m *Links_Related) GetHref() string {
	if m != nil {
		return m.Href
	}
	return ""
}

func (m *Links_Related) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Resource struct {
	Id            string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type          string            `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Attributes    map[string]string `protobuf:"bytes,3,rep,name=attributes" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Links         *Links            `protobuf:"bytes,4,opt,name=links" json:"links,omitempty"`
	Metadata      map[string]string `protobuf:"bytes,5,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Relationships *Relationships    `protobuf:"bytes,6,opt,name=relationships" json:"relationships,omitempty"`
}

func (m *Resource) Reset()                    { *m = Resource{} }
func (m *Resource) String() string            { return proto.CompactTextString(m) }
func (*Resource) ProtoMessage()               {}
func (*Resource) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Resource) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Resource) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Resource) GetAttributes() map[string]string {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *Resource) GetLinks() *Links {
	if m != nil {
		return m.Links
	}
	return nil
}

func (m *Resource) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Resource) GetRelationships() *Relationships {
	if m != nil {
		return m.Relationships
	}
	return nil
}

type Relationships struct {
	Links    *Links               `protobuf:"bytes,1,opt,name=links" json:"links,omitempty"`
	Data     *google_protobuf.Any `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
	Metadata map[string]string    `protobuf:"bytes,3,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Relationships) Reset()                    { *m = Relationships{} }
func (m *Relationships) String() string            { return proto.CompactTextString(m) }
func (*Relationships) ProtoMessage()               {}
func (*Relationships) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Relationships) GetLinks() *Links {
	if m != nil {
		return m.Links
	}
	return nil
}

func (m *Relationships) GetData() *google_protobuf.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Relationships) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func init() {
	proto.RegisterType((*JSONAPIObject)(nil), "gathergrpc.JSONAPIObject")
	proto.RegisterType((*Data)(nil), "gathergrpc.Data")
	proto.RegisterType((*Error)(nil), "gathergrpc.Error")
	proto.RegisterType((*Links)(nil), "gathergrpc.Links")
	proto.RegisterType((*Links_Related)(nil), "gathergrpc.Links.Related")
	proto.RegisterType((*Resource)(nil), "gathergrpc.Resource")
	proto.RegisterType((*Relationships)(nil), "gathergrpc.Relationships")
	proto.RegisterEnum("gathergrpc.UserState", UserState_name, UserState_value)
}

func init() { proto.RegisterFile("gathergrpc/commonModels.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 769 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xef, 0x4e, 0xdb, 0x48,
	0x10, 0x3f, 0xdb, 0xf9, 0x3b, 0x21, 0x9c, 0x6f, 0x85, 0x4e, 0x26, 0x12, 0x3a, 0x14, 0xa1, 0x83,
	0xbb, 0x0f, 0xe6, 0x04, 0x77, 0xd2, 0x1d, 0xe8, 0xda, 0x46, 0xc4, 0xaa, 0x42, 0x21, 0xa0, 0x0d,
	0xa1, 0xea, 0x27, 0xba, 0x89, 0x37, 0xc1, 0x60, 0xec, 0x68, 0xbd, 0x41, 0xcd, 0x0b, 0xf4, 0x21,
	0xfa, 0xa5, 0x0f, 0xd2, 0xd7, 0xe9, 0x23, 0xf4, 0x01, 0xaa, 0xdd, 0xb5, 0x13, 0x27, 0x21, 0x4d,
	0x11, 0xea, 0xb7, 0x99, 0x9d, 0xdf, 0xfc, 0xd9, 0xf9, 0xcd, 0xec, 0xc2, 0x46, 0x9f, 0xf0, 0x6b,
	0xca, 0xfa, 0x6c, 0xd0, 0xdd, 0xed, 0x86, 0x77, 0x77, 0x61, 0x70, 0x1a, 0xba, 0xd4, 0x8f, 0xec,
	0x01, 0x0b, 0x79, 0x88, 0x60, 0x62, 0xae, 0xac, 0xf7, 0xc3, 0xb0, 0xef, 0xd3, 0x5d, 0x69, 0xe9,
	0x0c, 0x7b, 0xbb, 0x24, 0x18, 0x29, 0x58, 0xf5, 0x93, 0x01, 0xe5, 0xe3, 0xd6, 0x59, 0xb3, 0x76,
	0xde, 0x38, 0xeb, 0xdc, 0xd0, 0x2e, 0x47, 0x5b, 0x90, 0x71, 0x09, 0x27, 0x96, 0xb6, 0x69, 0xec,
	0x94, 0xf6, 0x4c, 0x7b, 0x12, 0xc7, 0xae, 0x13, 0x4e, 0xb0, 0xb4, 0xa2, 0x3f, 0x20, 0x47, 0x19,
	0x0b, 0x59, 0x64, 0xe9, 0x12, 0xf7, 0x4b, 0x1a, 0xe7, 0x08, 0x0b, 0x8e, 0x01, 0xe8, 0x08, 0x0a,
	0x77, 0x94, 0x13, 0x19, 0xd4, 0x90, 0xe0, 0xed, 0x34, 0x78, 0x2a, 0xbb, 0x7d, 0x1a, 0x23, 0x9d,
	0x80, 0xb3, 0x11, 0x1e, 0x3b, 0xa2, 0x17, 0x90, 0xbf, 0x89, 0xc2, 0x80, 0x0c, 0x3c, 0x2b, 0x23,
	0x63, 0xfc, 0xbe, 0x38, 0xc6, 0xb1, 0x02, 0xaa, 0x10, 0x89, 0x1b, 0xda, 0x86, 0xac, 0xef, 0x05,
	0xb7, 0x91, 0x95, 0xdd, 0xd4, 0x66, 0x0b, 0x3e, 0x11, 0x06, 0xac, 0xec, 0xe8, 0x2f, 0x28, 0x78,
	0x41, 0xd7, 0x1f, 0xba, 0xd4, 0xb5, 0x72, 0x32, 0xd7, 0x5a, 0x1a, 0x8b, 0x69, 0x14, 0x0e, 0x59,
	0x97, 0xe2, 0x31, 0xaa, 0x72, 0x08, 0xe5, 0xa9, 0xba, 0x91, 0x09, 0xc6, 0x2d, 0x1d, 0x59, 0xda,
	0xa6, 0xb6, 0x53, 0xc4, 0x42, 0x44, 0x6b, 0x90, 0xbd, 0x27, 0xfe, 0x90, 0x5a, 0xba, 0x3c, 0x53,
	0xca, 0x81, 0xfe, 0xaf, 0x56, 0x39, 0x80, 0x95, 0x74, 0xc1, 0x8f, 0xf1, 0xad, 0xbe, 0x85, 0x8c,
	0xe0, 0x04, 0xad, 0x82, 0xee, 0xb9, 0xb1, 0x8b, 0xee, 0xb9, 0x08, 0x41, 0x86, 0x8f, 0x06, 0x89,
	0x83, 0x94, 0xd1, 0xdf, 0x00, 0x84, 0x73, 0xe6, 0x75, 0x86, 0x9c, 0x46, 0x96, 0x21, 0x9b, 0xb0,
	0x66, 0xab, 0xc9, 0xb0, 0x93, 0xc9, 0xb0, 0x6b, 0xc1, 0x08, 0xa7, 0x70, 0xd5, 0xf7, 0x06, 0x64,
	0x25, 0x9d, 0x73, 0x39, 0xc6, 0xfd, 0xd4, 0x97, 0xf4, 0xf3, 0x57, 0xc8, 0x45, 0x9c, 0xf0, 0xa1,
	0x4a, 0x5a, 0xc4, 0xb1, 0x26, 0x8a, 0xec, 0x86, 0x2e, 0xb5, 0x32, 0xaa, 0x48, 0x21, 0x8b, 0xab,
	0x72, 0x8f, 0xfb, 0x54, 0x92, 0x54, 0xc4, 0x4a, 0x11, 0x11, 0x5c, 0xca, 0x89, 0xe7, 0x5b, 0x39,
	0x15, 0x41, 0x69, 0xe8, 0x30, 0x35, 0x59, 0x79, 0xc9, 0xd4, 0x6f, 0x73, 0x63, 0xb8, 0x70, 0xa2,
	0xfe, 0x81, 0x9c, 0x22, 0xd2, 0x2a, 0x48, 0xd7, 0x8d, 0x79, 0xd7, 0x96, 0xb4, 0x2b, 0xc7, 0x18,
	0xfc, 0x34, 0xae, 0xff, 0x83, 0x52, 0x2a, 0xe6, 0xa3, 0xa8, 0xfe, 0xa2, 0x43, 0x56, 0xb6, 0x15,
	0xed, 0x43, 0x9e, 0x51, 0x9f, 0x70, 0xaa, 0xd8, 0x28, 0xed, 0xad, 0xcf, 0xb5, 0xde, 0xc6, 0x0a,
	0x80, 0x13, 0xa4, 0x68, 0x76, 0x44, 0xfd, 0x5e, 0x32, 0x11, 0x42, 0x16, 0xc9, 0x7a, 0x1e, 0x8b,
	0x78, 0xcc, 0x8b, 0x52, 0x04, 0x32, 0xa0, 0xef, 0x78, 0x42, 0x8b, 0x90, 0xc5, 0xd9, 0x80, 0xd1,
	0xfb, 0x98, 0x15, 0x29, 0x8b, 0x33, 0x9f, 0x44, 0x3c, 0xa6, 0x44, 0xca, 0x22, 0x22, 0xe9, 0x84,
	0x43, 0x6e, 0xe5, 0x55, 0x44, 0xa9, 0x20, 0x0b, 0xf2, 0x84, 0x71, 0xaf, 0xeb, 0x8b, 0x56, 0x8b,
	0xf3, 0x44, 0xad, 0x7c, 0xd4, 0x20, 0x8f, 0x27, 0x15, 0x5e, 0x33, 0xda, 0x8b, 0xbb, 0x21, 0xe5,
	0xa9, 0xa7, 0x43, 0x9f, 0x7f, 0x3a, 0xa6, 0xee, 0xba, 0x88, 0xe8, 0x27, 0x31, 0x56, 0xfd, 0x60,
	0x40, 0x21, 0xd9, 0xf8, 0xef, 0x5a, 0xb3, 0xfa, 0xcc, 0x9a, 0x89, 0xa2, 0xb7, 0x1e, 0x7a, 0x3f,
	0xec, 0xda, 0x18, 0xa6, 0x2a, 0x4e, 0xf9, 0x4d, 0x96, 0x2b, 0xb3, 0x64, 0xb9, 0x9e, 0xa5, 0x3a,
	0x94, 0x95, 0xc9, 0xaa, 0x0f, 0x26, 0x5b, 0xb4, 0x05, 0xcf, 0xa1, 0x2c, 0x47, 0xc4, 0x0b, 0x83,
	0xe8, 0xda, 0x1b, 0x44, 0x92, 0xce, 0x99, 0x91, 0xc2, 0x69, 0x00, 0x9e, 0xc6, 0x57, 0xfe, 0x87,
	0x9f, 0x67, 0x2e, 0xf2, 0xa8, 0x8d, 0x78, 0x12, 0x39, 0x9f, 0x35, 0x28, 0x4f, 0x15, 0x37, 0xe9,
	0x9b, 0xb6, 0xa4, 0x6f, 0x3b, 0xf1, 0x2f, 0xa7, 0x7f, 0xe3, 0x1d, 0x54, 0x3f, 0xdd, 0x92, 0xef,
	0x6b, 0x2a, 0xff, 0x0f, 0x99, 0xc1, 0x3f, 0xdf, 0x40, 0xb1, 0x1d, 0x51, 0xd6, 0xe2, 0x84, 0x53,
	0x64, 0xc2, 0x4a, 0xbb, 0xe5, 0xe0, 0xab, 0x76, 0xf3, 0x55, 0xf3, 0xec, 0x75, 0xd3, 0xfc, 0x09,
	0x95, 0x20, 0xdf, 0x68, 0x5e, 0x36, 0x2e, 0x9c, 0xba, 0xa9, 0x21, 0x80, 0x5c, 0xed, 0xe8, 0xa2,
	0x71, 0xe9, 0x98, 0x3a, 0x2a, 0x43, 0xb1, 0xd5, 0x6e, 0x9d, 0x3b, 0xcd, 0xba, 0x53, 0x37, 0x0d,
	0x84, 0x60, 0x55, 0xc8, 0x8d, 0xe6, 0xcb, 0xab, 0xba, 0x73, 0xe2, 0x5c, 0x38, 0x66, 0xa6, 0x93,
	0x93, 0x17, 0xde, 0xff, 0x1a, 0x00, 0x00, 0xff, 0xff, 0x43, 0xc4, 0xae, 0x86, 0x4d, 0x08, 0x00,
	0x00,
}
