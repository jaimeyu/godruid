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
	BulkOperationResult
	BulkOperationResponse
	TenantDescriptorData
	TenantDescriptor
	TenantDescriptorList
	AdminUserData
	AdminUser
	AdminUserList
	IngestionDictionaryData
	IngestionDictionary
	ValidTypesRequest
	ValidTypesData
	ValidTypes
	TenantDomainData
	TenantDomain
	TenantDomainList
	TenantDomainIdRequest
	TenantIngestionProfileData
	TenantIngestionProfile
	TenantIngestionProfileIdRequest
	TenantThresholdProfileData
	TenantThresholdProfile
	TenantThresholdProfileIdRequest
	TenantThresholdProfileList
	TenantUserData
	TenantUser
	TenantUserList
	TenantUserIdRequest
	MonitoredObjectData
	MonitoredObject
	MonitoredObjectList
	MonitoredObjectIdRequest
	MonitoredObjectCountByDomainRequest
	MonitoredObjectCountByDomainResponse
	MonitoredObjectSet
	TenantMonitoredObjectSet
	TenantMeta
	TenantMetadata
	TenantMetaIdRequest
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

type BulkOperationResult struct {
	Ok     bool   `protobuf:"varint,1,opt,name=ok" json:"ok,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Rev    string `protobuf:"bytes,3,opt,name=rev" json:"rev,omitempty"`
	Error  string `protobuf:"bytes,4,opt,name=error" json:"error,omitempty"`
	Reason string `protobuf:"bytes,5,opt,name=reason" json:"reason,omitempty"`
}

func (m *BulkOperationResult) Reset()                    { *m = BulkOperationResult{} }
func (m *BulkOperationResult) String() string            { return proto.CompactTextString(m) }
func (*BulkOperationResult) ProtoMessage()               {}
func (*BulkOperationResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *BulkOperationResult) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *BulkOperationResult) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BulkOperationResult) GetRev() string {
	if m != nil {
		return m.Rev
	}
	return ""
}

func (m *BulkOperationResult) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *BulkOperationResult) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type BulkOperationResponse struct {
	Results []*BulkOperationResult `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}

func (m *BulkOperationResponse) Reset()                    { *m = BulkOperationResponse{} }
func (m *BulkOperationResponse) String() string            { return proto.CompactTextString(m) }
func (*BulkOperationResponse) ProtoMessage()               {}
func (*BulkOperationResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *BulkOperationResponse) GetResults() []*BulkOperationResult {
	if m != nil {
		return m.Results
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
	proto.RegisterType((*BulkOperationResult)(nil), "gathergrpc.BulkOperationResult")
	proto.RegisterType((*BulkOperationResponse)(nil), "gathergrpc.BulkOperationResponse")
	proto.RegisterEnum("gathergrpc.UserState", UserState_name, UserState_value)
}

func init() { proto.RegisterFile("gathergrpc/commonModels.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 852 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0xdd, 0x8e, 0xdb, 0x44,
	0x14, 0xc6, 0x76, 0x7e, 0xcf, 0x36, 0xc5, 0x0c, 0x0b, 0x72, 0x23, 0x55, 0xac, 0xa2, 0x8a, 0x2e,
	0x5c, 0x78, 0x51, 0x0b, 0x12, 0x6d, 0xc5, 0xcf, 0xd2, 0x58, 0x28, 0xa5, 0xcd, 0x56, 0x93, 0x4d,
	0x11, 0x57, 0x65, 0x12, 0x4f, 0xb2, 0x6e, 0xbc, 0x1e, 0x6b, 0x66, 0x5c, 0x91, 0x17, 0xe0, 0x21,
	0xb8, 0xe1, 0x41, 0x78, 0x1d, 0x1e, 0x81, 0x07, 0x40, 0x33, 0x63, 0x27, 0x76, 0xb2, 0x21, 0xac,
	0x56, 0xbd, 0x3b, 0x67, 0xe6, 0x3b, 0x3f, 0x3e, 0xdf, 0x39, 0x67, 0x0c, 0x77, 0xe7, 0x44, 0x5e,
	0x50, 0x3e, 0xe7, 0xe9, 0xf4, 0x64, 0xca, 0x2e, 0x2f, 0x59, 0xf2, 0x82, 0x85, 0x34, 0x16, 0x7e,
	0xca, 0x99, 0x64, 0x08, 0xd6, 0xd7, 0xdd, 0x3b, 0x73, 0xc6, 0xe6, 0x31, 0x3d, 0xd1, 0x37, 0x93,
	0x6c, 0x76, 0x42, 0x92, 0xa5, 0x81, 0xf5, 0xfe, 0x72, 0xa0, 0xf3, 0x6c, 0x74, 0x36, 0x3c, 0x7d,
	0x39, 0x38, 0x9b, 0xbc, 0xa1, 0x53, 0x89, 0xee, 0x41, 0x2d, 0x24, 0x92, 0x78, 0xd6, 0x91, 0x73,
	0x7c, 0xf0, 0xc0, 0xf5, 0xd7, 0x7e, 0xfc, 0x3e, 0x91, 0x04, 0xeb, 0x5b, 0xf4, 0x19, 0x34, 0x28,
	0xe7, 0x8c, 0x0b, 0xcf, 0xd6, 0xb8, 0x0f, 0xca, 0xb8, 0x40, 0xdd, 0xe0, 0x1c, 0x80, 0x9e, 0x42,
	0xeb, 0x92, 0x4a, 0xa2, 0x9d, 0x3a, 0x1a, 0x7c, 0xbf, 0x0c, 0xae, 0x44, 0xf7, 0x5f, 0xe4, 0xc8,
	0x20, 0x91, 0x7c, 0x89, 0x57, 0x86, 0xe8, 0x7b, 0x68, 0xbe, 0x11, 0x2c, 0x21, 0x69, 0xe4, 0xd5,
	0xb4, 0x8f, 0x4f, 0x77, 0xfb, 0x78, 0x66, 0x80, 0xc6, 0x45, 0x61, 0x86, 0xee, 0x43, 0x3d, 0x8e,
	0x92, 0x85, 0xf0, 0xea, 0x47, 0xd6, 0x66, 0xc2, 0xcf, 0xd5, 0x05, 0x36, 0xf7, 0xe8, 0x0b, 0x68,
	0x45, 0xc9, 0x34, 0xce, 0x42, 0x1a, 0x7a, 0x0d, 0x1d, 0xeb, 0xb0, 0x8c, 0xc5, 0x54, 0xb0, 0x8c,
	0x4f, 0x29, 0x5e, 0xa1, 0xba, 0x4f, 0xa0, 0x53, 0xc9, 0x1b, 0xb9, 0xe0, 0x2c, 0xe8, 0xd2, 0xb3,
	0x8e, 0xac, 0xe3, 0x36, 0x56, 0x22, 0x3a, 0x84, 0xfa, 0x5b, 0x12, 0x67, 0xd4, 0xb3, 0xf5, 0x99,
	0x51, 0x1e, 0xdb, 0x5f, 0x5b, 0xdd, 0xc7, 0x70, 0xab, 0x9c, 0xf0, 0x75, 0x6c, 0x7b, 0xbf, 0x42,
	0x4d, 0x71, 0x82, 0x6e, 0x83, 0x1d, 0x85, 0xb9, 0x89, 0x1d, 0x85, 0x08, 0x41, 0x4d, 0x2e, 0xd3,
	0xc2, 0x40, 0xcb, 0xe8, 0x4b, 0x00, 0x22, 0x25, 0x8f, 0x26, 0x99, 0xa4, 0xc2, 0x73, 0x74, 0x11,
	0x0e, 0x7d, 0xd3, 0x19, 0x7e, 0xd1, 0x19, 0xfe, 0x69, 0xb2, 0xc4, 0x25, 0x5c, 0xef, 0x77, 0x07,
	0xea, 0x9a, 0xce, 0xad, 0x18, 0xab, 0x7a, 0xda, 0x7b, 0xea, 0xf9, 0x31, 0x34, 0x84, 0x24, 0x32,
	0x33, 0x41, 0xdb, 0x38, 0xd7, 0x54, 0x92, 0x53, 0x16, 0x52, 0xaf, 0x66, 0x92, 0x54, 0xb2, 0xfa,
	0x54, 0x19, 0xc9, 0x98, 0x6a, 0x92, 0xda, 0xd8, 0x28, 0xca, 0x43, 0x48, 0x25, 0x89, 0x62, 0xaf,
	0x61, 0x3c, 0x18, 0x0d, 0x3d, 0x29, 0x75, 0x56, 0x53, 0x33, 0xf5, 0xc9, 0x56, 0x1b, 0xee, 0xec,
	0xa8, 0xaf, 0xa0, 0x61, 0x88, 0xf4, 0x5a, 0xda, 0xf4, 0xee, 0xb6, 0xe9, 0x48, 0xdf, 0x1b, 0xc3,
	0x1c, 0x7c, 0x33, 0xae, 0x1f, 0xc1, 0x41, 0xc9, 0xe7, 0xb5, 0xa8, 0xfe, 0xc7, 0x86, 0xba, 0x2e,
	0x2b, 0x7a, 0x08, 0x4d, 0x4e, 0x63, 0x22, 0xa9, 0x61, 0xe3, 0xe0, 0xc1, 0x9d, 0xad, 0xd2, 0xfb,
	0xd8, 0x00, 0x70, 0x81, 0x54, 0xc5, 0x16, 0x34, 0x9e, 0x15, 0x1d, 0xa1, 0x64, 0x15, 0x6c, 0x16,
	0x71, 0x21, 0x73, 0x5e, 0x8c, 0xa2, 0x90, 0x09, 0xfd, 0x4d, 0x16, 0xb4, 0x28, 0x59, 0x9d, 0xa5,
	0x9c, 0xbe, 0xcd, 0x59, 0xd1, 0xb2, 0x3a, 0x8b, 0x89, 0x90, 0x39, 0x25, 0x5a, 0x56, 0x1e, 0xc9,
	0x84, 0x65, 0xd2, 0x6b, 0x1a, 0x8f, 0x5a, 0x41, 0x1e, 0x34, 0x09, 0x97, 0xd1, 0x34, 0x56, 0xa5,
	0x56, 0xe7, 0x85, 0xda, 0xfd, 0xd3, 0x82, 0x26, 0x5e, 0x67, 0x78, 0xc1, 0xe9, 0x2c, 0xaf, 0x86,
	0x96, 0x2b, 0xab, 0xc3, 0xde, 0x5e, 0x1d, 0x95, 0x6f, 0xdd, 0x45, 0xf4, 0x8d, 0x18, 0xeb, 0xfd,
	0xe1, 0x40, 0xab, 0x98, 0xf8, 0xff, 0x35, 0x66, 0xfd, 0x8d, 0x31, 0x53, 0x49, 0xdf, 0xbb, 0x6a,
	0x7f, 0xf8, 0xa7, 0x2b, 0x98, 0xc9, 0xb8, 0x64, 0xb7, 0x1e, 0xae, 0xda, 0x9e, 0xe1, 0xfa, 0xb6,
	0x54, 0xa1, 0xba, 0x0e, 0xd6, 0xbb, 0x32, 0xd8, 0xae, 0x29, 0xf8, 0x0e, 0x3a, 0xba, 0x45, 0x22,
	0x96, 0x88, 0x8b, 0x28, 0x15, 0x9a, 0xce, 0x8d, 0x96, 0xc2, 0x65, 0x00, 0xae, 0xe2, 0xbb, 0xdf,
	0xc0, 0xfb, 0x1b, 0x1f, 0x72, 0xad, 0x89, 0xb8, 0x11, 0x39, 0x7f, 0x5b, 0xd0, 0xa9, 0x24, 0xb7,
	0xae, 0x9b, 0xb5, 0xa7, 0x6e, 0xc7, 0xf9, 0x2b, 0x67, 0xff, 0xc7, 0x1e, 0x34, 0x2f, 0xdd, 0x9e,
	0xe7, 0xab, 0x12, 0xff, 0xdd, 0xf4, 0x60, 0x06, 0x1f, 0xfe, 0x90, 0xc5, 0x8b, 0xb3, 0x94, 0x72,
	0x1d, 0x0a, 0x53, 0x91, 0xc5, 0x52, 0x75, 0x23, 0x5b, 0x68, 0x0f, 0x2d, 0x6c, 0xb3, 0x45, 0xde,
	0x9d, 0xf6, 0xaa, 0x3b, 0x5d, 0x70, 0xd4, 0xcc, 0x9a, 0xe1, 0x56, 0xa2, 0x0a, 0xa1, 0xdf, 0xe4,
	0x7c, 0xb6, 0x8d, 0xa2, 0xb6, 0x2b, 0xa7, 0x44, 0xb0, 0x24, 0x1f, 0xef, 0x5c, 0xeb, 0x61, 0xf8,
	0x68, 0x33, 0x6c, 0xca, 0x12, 0x41, 0xd1, 0x23, 0xb5, 0x80, 0x54, 0x0a, 0x22, 0xff, 0x49, 0xa8,
	0x6c, 0xdd, 0x2b, 0x52, 0xc5, 0x05, 0xfe, 0xf3, 0x5f, 0xa0, 0x3d, 0x16, 0x94, 0x8f, 0x24, 0x91,
	0x14, 0xb9, 0x70, 0x6b, 0x3c, 0x0a, 0xf0, 0xeb, 0xf1, 0xf0, 0xa7, 0xe1, 0xd9, 0xcf, 0x43, 0xf7,
	0x3d, 0x74, 0x00, 0xcd, 0xc1, 0xf0, 0xd5, 0xe0, 0x3c, 0xe8, 0xbb, 0x16, 0x02, 0x68, 0x9c, 0x3e,
	0x3d, 0x1f, 0xbc, 0x0a, 0x5c, 0x1b, 0x75, 0xa0, 0x3d, 0x1a, 0x8f, 0x5e, 0x06, 0xc3, 0x7e, 0xd0,
	0x77, 0x1d, 0x84, 0xe0, 0xb6, 0x92, 0x07, 0xc3, 0x1f, 0x5f, 0xf7, 0x83, 0xe7, 0xc1, 0x79, 0xe0,
	0xd6, 0x26, 0x0d, 0xcd, 0xdd, 0xc3, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x9c, 0x0d, 0xe2, 0x29,
	0x18, 0x09, 0x00, 0x00,
}
