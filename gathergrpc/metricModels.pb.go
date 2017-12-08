// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gathergrpc/metricModels.proto

package gathergrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/struct"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ThresholdCrossing struct {
	Timestamp string             `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Result    map[string]float32 `protobuf:"bytes,2,rep,name=result" json:"result,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"fixed32,2,opt,name=value"`
}

func (m *ThresholdCrossing) Reset()                    { *m = ThresholdCrossing{} }
func (m *ThresholdCrossing) String() string            { return proto.CompactTextString(m) }
func (*ThresholdCrossing) ProtoMessage()               {}
func (*ThresholdCrossing) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *ThresholdCrossing) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *ThresholdCrossing) GetResult() map[string]float32 {
	if m != nil {
		return m.Result
	}
	return nil
}

type FormattedThresholdCrossing struct {
	Timestamp string                   `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Result    *google_protobuf1.Struct `protobuf:"bytes,2,opt,name=result" json:"result,omitempty"`
}

func (m *FormattedThresholdCrossing) Reset()                    { *m = FormattedThresholdCrossing{} }
func (m *FormattedThresholdCrossing) String() string            { return proto.CompactTextString(m) }
func (*FormattedThresholdCrossing) ProtoMessage()               {}
func (*FormattedThresholdCrossing) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *FormattedThresholdCrossing) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *FormattedThresholdCrossing) GetResult() *google_protobuf1.Struct {
	if m != nil {
		return m.Result
	}
	return nil
}

type ThresholdCrossingResponse struct {
	Data []*FormattedThresholdCrossing `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *ThresholdCrossingResponse) Reset()                    { *m = ThresholdCrossingResponse{} }
func (m *ThresholdCrossingResponse) String() string            { return proto.CompactTextString(m) }
func (*ThresholdCrossingResponse) ProtoMessage()               {}
func (*ThresholdCrossingResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *ThresholdCrossingResponse) GetData() []*FormattedThresholdCrossing {
	if m != nil {
		return m.Data
	}
	return nil
}

type ThresholdCrossingRequest struct {
	// ISO-8601 Intervals
	Interval string `protobuf:"bytes,1,opt,name=interval" json:"interval,omitempty"`
	Tenant   string `protobuf:"bytes,2,opt,name=tenant" json:"tenant,omitempty"`
	Domain   string `protobuf:"bytes,3,opt,name=domain" json:"domain,omitempty"`
	// ISO-8601 period combination
	Granularity        string `protobuf:"bytes,4,opt,name=granularity" json:"granularity,omitempty"`
	ObjectType         string `protobuf:"bytes,5,opt,name=objectType" json:"objectType,omitempty"`
	Direction          string `protobuf:"bytes,6,opt,name=direction" json:"direction,omitempty"`
	Metric             string `protobuf:"bytes,7,opt,name=metric" json:"metric,omitempty"`
	ThresholdProfileId string `protobuf:"bytes,8,opt,name=thresholdProfileId" json:"thresholdProfileId,omitempty"`
}

func (m *ThresholdCrossingRequest) Reset()                    { *m = ThresholdCrossingRequest{} }
func (m *ThresholdCrossingRequest) String() string            { return proto.CompactTextString(m) }
func (*ThresholdCrossingRequest) ProtoMessage()               {}
func (*ThresholdCrossingRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *ThresholdCrossingRequest) GetInterval() string {
	if m != nil {
		return m.Interval
	}
	return ""
}

func (m *ThresholdCrossingRequest) GetTenant() string {
	if m != nil {
		return m.Tenant
	}
	return ""
}

func (m *ThresholdCrossingRequest) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *ThresholdCrossingRequest) GetGranularity() string {
	if m != nil {
		return m.Granularity
	}
	return ""
}

func (m *ThresholdCrossingRequest) GetObjectType() string {
	if m != nil {
		return m.ObjectType
	}
	return ""
}

func (m *ThresholdCrossingRequest) GetDirection() string {
	if m != nil {
		return m.Direction
	}
	return ""
}

func (m *ThresholdCrossingRequest) GetMetric() string {
	if m != nil {
		return m.Metric
	}
	return ""
}

func (m *ThresholdCrossingRequest) GetThresholdProfileId() string {
	if m != nil {
		return m.ThresholdProfileId
	}
	return ""
}

type HistogramRequest struct {
	// ISO-8601 Intervals
	Interval string `protobuf:"bytes,1,opt,name=interval" json:"interval,omitempty"`
	Tenant   string `protobuf:"bytes,2,opt,name=tenant" json:"tenant,omitempty"`
	Domain   string `protobuf:"bytes,3,opt,name=domain" json:"domain,omitempty"`
	// ISO-8601 period combination
	Granularity        string `protobuf:"bytes,4,opt,name=granularity" json:"granularity,omitempty"`
	Direction          string `protobuf:"bytes,5,opt,name=direction" json:"direction,omitempty"`
	Metric             string `protobuf:"bytes,6,opt,name=metric" json:"metric,omitempty"`
	GranularityBuckets int32  `protobuf:"varint,7,opt,name=granularityBuckets" json:"granularityBuckets,omitempty"`
	Resolution         int32  `protobuf:"varint,8,opt,name=resolution" json:"resolution,omitempty"`
}

func (m *HistogramRequest) Reset()                    { *m = HistogramRequest{} }
func (m *HistogramRequest) String() string            { return proto.CompactTextString(m) }
func (*HistogramRequest) ProtoMessage()               {}
func (*HistogramRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *HistogramRequest) GetInterval() string {
	if m != nil {
		return m.Interval
	}
	return ""
}

func (m *HistogramRequest) GetTenant() string {
	if m != nil {
		return m.Tenant
	}
	return ""
}

func (m *HistogramRequest) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *HistogramRequest) GetGranularity() string {
	if m != nil {
		return m.Granularity
	}
	return ""
}

func (m *HistogramRequest) GetDirection() string {
	if m != nil {
		return m.Direction
	}
	return ""
}

func (m *HistogramRequest) GetMetric() string {
	if m != nil {
		return m.Metric
	}
	return ""
}

func (m *HistogramRequest) GetGranularityBuckets() int32 {
	if m != nil {
		return m.GranularityBuckets
	}
	return 0
}

func (m *HistogramRequest) GetResolution() int32 {
	if m != nil {
		return m.Resolution
	}
	return 0
}

type HistogramBuckets struct {
	Breaks []float32 `protobuf:"fixed32,1,rep,packed,name=breaks" json:"breaks,omitempty"`
	Counts []float32 `protobuf:"fixed32,2,rep,packed,name=counts" json:"counts,omitempty"`
}

func (m *HistogramBuckets) Reset()                    { *m = HistogramBuckets{} }
func (m *HistogramBuckets) String() string            { return proto.CompactTextString(m) }
func (*HistogramBuckets) ProtoMessage()               {}
func (*HistogramBuckets) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *HistogramBuckets) GetBreaks() []float32 {
	if m != nil {
		return m.Breaks
	}
	return nil
}

func (m *HistogramBuckets) GetCounts() []float32 {
	if m != nil {
		return m.Counts
	}
	return nil
}

type HistogramResult struct {
	ThresholdBuckets *HistogramBuckets `protobuf:"bytes,1,opt,name=thresholdBuckets" json:"thresholdBuckets,omitempty"`
}

func (m *HistogramResult) Reset()                    { *m = HistogramResult{} }
func (m *HistogramResult) String() string            { return proto.CompactTextString(m) }
func (*HistogramResult) ProtoMessage()               {}
func (*HistogramResult) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *HistogramResult) GetThresholdBuckets() *HistogramBuckets {
	if m != nil {
		return m.ThresholdBuckets
	}
	return nil
}

type Histogram struct {
	Timestamp string           `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Result    *HistogramResult `protobuf:"bytes,2,opt,name=result" json:"result,omitempty"`
}

func (m *Histogram) Reset()                    { *m = Histogram{} }
func (m *Histogram) String() string            { return proto.CompactTextString(m) }
func (*Histogram) ProtoMessage()               {}
func (*Histogram) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

func (m *Histogram) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *Histogram) GetResult() *HistogramResult {
	if m != nil {
		return m.Result
	}
	return nil
}

type HistogramResponse struct {
	Data []*Histogram `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *HistogramResponse) Reset()                    { *m = HistogramResponse{} }
func (m *HistogramResponse) String() string            { return proto.CompactTextString(m) }
func (*HistogramResponse) ProtoMessage()               {}
func (*HistogramResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

func (m *HistogramResponse) GetData() []*Histogram {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*ThresholdCrossing)(nil), "gathergrpc.ThresholdCrossing")
	proto.RegisterType((*FormattedThresholdCrossing)(nil), "gathergrpc.FormattedThresholdCrossing")
	proto.RegisterType((*ThresholdCrossingResponse)(nil), "gathergrpc.ThresholdCrossingResponse")
	proto.RegisterType((*ThresholdCrossingRequest)(nil), "gathergrpc.ThresholdCrossingRequest")
	proto.RegisterType((*HistogramRequest)(nil), "gathergrpc.HistogramRequest")
	proto.RegisterType((*HistogramBuckets)(nil), "gathergrpc.HistogramBuckets")
	proto.RegisterType((*HistogramResult)(nil), "gathergrpc.HistogramResult")
	proto.RegisterType((*Histogram)(nil), "gathergrpc.Histogram")
	proto.RegisterType((*HistogramResponse)(nil), "gathergrpc.HistogramResponse")
}

func init() { proto.RegisterFile("gathergrpc/metricModels.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 538 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x53, 0x41, 0x6f, 0xd3, 0x4c,
	0x10, 0x95, 0xdd, 0xc6, 0x5f, 0x33, 0x39, 0x7c, 0xe9, 0x0a, 0x8a, 0x09, 0x01, 0x45, 0x3e, 0xa0,
	0xf4, 0xe2, 0x48, 0xed, 0x05, 0x7a, 0x40, 0xa2, 0x08, 0x54, 0x0e, 0x48, 0x68, 0xa9, 0xc4, 0x01,
	0x09, 0x69, 0x63, 0x4f, 0x1d, 0x13, 0xdb, 0x1b, 0x76, 0xc7, 0x95, 0xf2, 0x07, 0x38, 0xf2, 0x57,
	0xf8, 0x8b, 0xc8, 0xbb, 0x4e, 0xb2, 0x4d, 0x52, 0x21, 0x4e, 0xdc, 0x76, 0xde, 0xcc, 0xce, 0x9b,
	0x79, 0x9a, 0x07, 0x4f, 0x33, 0x41, 0x33, 0x54, 0x99, 0x5a, 0x24, 0x93, 0x12, 0x49, 0xe5, 0xc9,
	0x07, 0x99, 0x62, 0xa1, 0xe3, 0x85, 0x92, 0x24, 0x19, 0x6c, 0xd2, 0x83, 0x61, 0x26, 0x65, 0x56,
	0xe0, 0xc4, 0x64, 0xa6, 0xf5, 0xcd, 0x44, 0x93, 0xaa, 0x13, 0xb2, 0x95, 0xd1, 0x2f, 0x0f, 0x8e,
	0xaf, 0x67, 0x0a, 0xf5, 0x4c, 0x16, 0xe9, 0x1b, 0x25, 0xb5, 0xce, 0xab, 0x8c, 0x0d, 0xa1, 0x4b,
	0x79, 0x89, 0x9a, 0x44, 0xb9, 0x08, 0xbd, 0x91, 0x37, 0xee, 0xf2, 0x0d, 0xc0, 0x5e, 0x43, 0xa0,
	0x50, 0xd7, 0x05, 0x85, 0xfe, 0xe8, 0x60, 0xdc, 0x3b, 0x3b, 0x8d, 0x37, 0x74, 0xf1, 0x4e, 0xb3,
	0x98, 0x9b, 0xda, 0xb7, 0x15, 0xa9, 0x25, 0x6f, 0x3f, 0x0e, 0x5e, 0x42, 0xcf, 0x81, 0x59, 0x1f,
	0x0e, 0xe6, 0xb8, 0x6c, 0x99, 0x9a, 0x27, 0x7b, 0x00, 0x9d, 0x5b, 0x51, 0xd4, 0x18, 0xfa, 0x23,
	0x6f, 0xec, 0x73, 0x1b, 0x5c, 0xf8, 0x2f, 0xbc, 0x68, 0x0e, 0x83, 0x77, 0x52, 0x95, 0x82, 0x08,
	0xd3, 0xbf, 0x9d, 0x7c, 0xe2, 0x4c, 0xee, 0x8d, 0x7b, 0x67, 0x8f, 0x62, 0x2b, 0x4e, 0xbc, 0x12,
	0x27, 0xfe, 0x64, 0xc4, 0x59, 0xcd, 0x19, 0x7d, 0x86, 0xc7, 0x3b, 0x1c, 0x1c, 0xf5, 0x42, 0x56,
	0x1a, 0xd9, 0x05, 0x1c, 0xa6, 0x82, 0x44, 0xe8, 0x19, 0x15, 0x9e, 0xbb, 0x2a, 0xdc, 0x3f, 0x21,
	0x37, 0x7f, 0xa2, 0x9f, 0x3e, 0x84, 0x7b, 0x3a, 0x7f, 0xaf, 0x51, 0x13, 0x1b, 0xc0, 0x51, 0x5e,
	0x11, 0xaa, 0x5b, 0x51, 0xb4, 0x3b, 0xac, 0x63, 0x76, 0x02, 0x01, 0x61, 0x25, 0x2a, 0xbb, 0x42,
	0x97, 0xb7, 0x51, 0x83, 0xa7, 0xb2, 0x14, 0x79, 0x15, 0x1e, 0x58, 0xdc, 0x46, 0x6c, 0x04, 0xbd,
	0x4c, 0x89, 0xaa, 0x2e, 0x84, 0xca, 0x69, 0x19, 0x1e, 0x9a, 0xa4, 0x0b, 0xb1, 0x67, 0x00, 0x72,
	0xfa, 0x0d, 0x13, 0xba, 0x5e, 0x2e, 0x30, 0xec, 0x98, 0x02, 0x07, 0x69, 0x24, 0x4d, 0x73, 0x85,
	0x09, 0xe5, 0xb2, 0x0a, 0x03, 0x2b, 0xe9, 0x1a, 0x68, 0x78, 0xed, 0x01, 0x86, 0xff, 0x59, 0x5e,
	0x1b, 0xb1, 0x18, 0x18, 0xad, 0xf6, 0xfb, 0xa8, 0xe4, 0x4d, 0x5e, 0xe0, 0xfb, 0x34, 0x3c, 0x32,
	0x35, 0x7b, 0x32, 0xd1, 0x0f, 0x1f, 0xfa, 0x57, 0xb9, 0x26, 0x99, 0x29, 0x51, 0xfe, 0x1b, 0x21,
	0xee, 0x2c, 0xda, 0xb9, 0x7f, 0xd1, 0x60, 0x7b, 0x51, 0xa7, 0xc9, 0x65, 0x9d, 0xcc, 0x91, 0xb4,
	0x11, 0xa3, 0xc3, 0xf7, 0x64, 0x1a, 0xb9, 0x15, 0x6a, 0x59, 0xd4, 0x86, 0xe6, 0xc8, 0xd4, 0x39,
	0x48, 0x74, 0xe9, 0xe8, 0xb0, 0xfa, 0x73, 0x02, 0xc1, 0x54, 0xa1, 0x98, 0x6b, 0x73, 0x6b, 0x3e,
	0x6f, 0xa3, 0x06, 0x4f, 0x64, 0x5d, 0x91, 0x36, 0x4e, 0xf4, 0x79, 0x1b, 0x45, 0x5f, 0xe0, 0x7f,
	0x47, 0xcb, 0xe6, 0x92, 0xd9, 0x15, 0xf4, 0xd7, 0xaa, 0xaf, 0x86, 0xf4, 0x8c, 0x09, 0x86, 0xee,
	0xe1, 0x6e, 0x53, 0xf3, 0x9d, 0x5f, 0xd1, 0x57, 0xe8, 0xae, 0xab, 0xfe, 0xe0, 0xb7, 0xf3, 0x2d,
	0xbf, 0x3d, 0xd9, 0x4b, 0x65, 0x27, 0x5c, 0x7b, 0xee, 0x15, 0x1c, 0xbb, 0x29, 0xeb, 0xb5, 0xd3,
	0x3b, 0x5e, 0x7b, 0xb8, 0xbf, 0x8f, 0x29, 0x99, 0x06, 0xc6, 0xcc, 0xe7, 0xbf, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x93, 0x83, 0x76, 0xad, 0x24, 0x05, 0x00, 0x00,
}
