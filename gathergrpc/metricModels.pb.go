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

type ThresholdCrossingByMonitoredObjectResponse struct {
	Result *google_protobuf1.Struct `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *ThresholdCrossingByMonitoredObjectResponse) Reset() {
	*m = ThresholdCrossingByMonitoredObjectResponse{}
}
func (m *ThresholdCrossingByMonitoredObjectResponse) String() string {
	return proto.CompactTextString(m)
}
func (*ThresholdCrossingByMonitoredObjectResponse) ProtoMessage() {}
func (*ThresholdCrossingByMonitoredObjectResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor3, []int{9}
}

func (m *ThresholdCrossingByMonitoredObjectResponse) GetResult() *google_protobuf1.Struct {
	if m != nil {
		return m.Result
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
	proto.RegisterType((*ThresholdCrossingByMonitoredObjectResponse)(nil), "gathergrpc.ThresholdCrossingByMonitoredObjectResponse")
}

func init() { proto.RegisterFile("gathergrpc/metricModels.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 562 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x95, 0x9d, 0xc6, 0x24, 0x93, 0x03, 0xe9, 0x0a, 0x8a, 0x09, 0x01, 0x45, 0x3e, 0xa0, 0x94,
	0x83, 0x23, 0xb5, 0x17, 0xe8, 0x01, 0x89, 0x20, 0x50, 0x39, 0x54, 0x20, 0x53, 0x89, 0x03, 0x02,
	0x69, 0x63, 0x4f, 0x1d, 0x13, 0xdb, 0x1b, 0x76, 0xc7, 0x95, 0xf2, 0x03, 0x1c, 0xf9, 0x15, 0x7e,
	0x11, 0x79, 0xd7, 0x71, 0xdc, 0x24, 0x15, 0x70, 0xe2, 0xe6, 0x79, 0x33, 0x3b, 0x6f, 0xe6, 0x69,
	0x9e, 0xe1, 0x71, 0xcc, 0x69, 0x8e, 0x32, 0x96, 0xcb, 0x70, 0x92, 0x21, 0xc9, 0x24, 0xbc, 0x10,
	0x11, 0xa6, 0xca, 0x5f, 0x4a, 0x41, 0x82, 0xc1, 0x26, 0x3d, 0x18, 0xc6, 0x42, 0xc4, 0x29, 0x4e,
	0x74, 0x66, 0x56, 0x5c, 0x4d, 0x14, 0xc9, 0x22, 0x24, 0x53, 0xe9, 0xfd, 0xb2, 0xe0, 0xf0, 0x72,
	0x2e, 0x51, 0xcd, 0x45, 0x1a, 0xbd, 0x96, 0x42, 0xa9, 0x24, 0x8f, 0xd9, 0x10, 0xba, 0x94, 0x64,
	0xa8, 0x88, 0x67, 0x4b, 0xd7, 0x1a, 0x59, 0xe3, 0x6e, 0xb0, 0x01, 0xd8, 0x2b, 0x70, 0x24, 0xaa,
	0x22, 0x25, 0xd7, 0x1e, 0xb5, 0xc6, 0xbd, 0x93, 0x63, 0x7f, 0x43, 0xe7, 0xef, 0x34, 0xf3, 0x03,
	0x5d, 0xfb, 0x26, 0x27, 0xb9, 0x0a, 0xaa, 0x87, 0x83, 0x17, 0xd0, 0x6b, 0xc0, 0xac, 0x0f, 0xad,
	0x05, 0xae, 0x2a, 0xa6, 0xf2, 0x93, 0xdd, 0x83, 0xf6, 0x35, 0x4f, 0x0b, 0x74, 0xed, 0x91, 0x35,
	0xb6, 0x03, 0x13, 0x9c, 0xd9, 0xcf, 0x2d, 0x6f, 0x01, 0x83, 0xb7, 0x42, 0x66, 0x9c, 0x08, 0xa3,
	0x7f, 0x9d, 0x7c, 0xd2, 0x98, 0xdc, 0x1a, 0xf7, 0x4e, 0x1e, 0xf8, 0x46, 0x1c, 0x7f, 0x2d, 0x8e,
	0xff, 0x51, 0x8b, 0xb3, 0x9e, 0xd3, 0xfb, 0x04, 0x0f, 0x77, 0x38, 0x02, 0x54, 0x4b, 0x91, 0x2b,
	0x64, 0x67, 0x70, 0x10, 0x71, 0xe2, 0xae, 0xa5, 0x55, 0x78, 0xda, 0x54, 0xe1, 0xf6, 0x09, 0x03,
	0xfd, 0xc6, 0xfb, 0x69, 0x83, 0xbb, 0xa7, 0xf3, 0xf7, 0x02, 0x15, 0xb1, 0x01, 0x74, 0x92, 0x9c,
	0x50, 0x5e, 0xf3, 0xb4, 0xda, 0xa1, 0x8e, 0xd9, 0x11, 0x38, 0x84, 0x39, 0xcf, 0xcd, 0x0a, 0xdd,
	0xa0, 0x8a, 0x4a, 0x3c, 0x12, 0x19, 0x4f, 0x72, 0xb7, 0x65, 0x70, 0x13, 0xb1, 0x11, 0xf4, 0x62,
	0xc9, 0xf3, 0x22, 0xe5, 0x32, 0xa1, 0x95, 0x7b, 0xa0, 0x93, 0x4d, 0x88, 0x3d, 0x01, 0x10, 0xb3,
	0x6f, 0x18, 0xd2, 0xe5, 0x6a, 0x89, 0x6e, 0x5b, 0x17, 0x34, 0x90, 0x52, 0xd2, 0x28, 0x91, 0x18,
	0x52, 0x22, 0x72, 0xd7, 0x31, 0x92, 0xd6, 0x40, 0xc9, 0x6b, 0x0e, 0xd0, 0xbd, 0x63, 0x78, 0x4d,
	0xc4, 0x7c, 0x60, 0xb4, 0xde, 0xef, 0x83, 0x14, 0x57, 0x49, 0x8a, 0xef, 0x22, 0xb7, 0xa3, 0x6b,
	0xf6, 0x64, 0xbc, 0x1f, 0x36, 0xf4, 0xcf, 0x13, 0x45, 0x22, 0x96, 0x3c, 0xfb, 0x3f, 0x42, 0xdc,
	0x58, 0xb4, 0x7d, 0xfb, 0xa2, 0xce, 0xf6, 0xa2, 0x8d, 0x26, 0xd3, 0x22, 0x5c, 0x20, 0x29, 0x2d,
	0x46, 0x3b, 0xd8, 0x93, 0x29, 0xe5, 0x96, 0xa8, 0x44, 0x5a, 0x68, 0x9a, 0x8e, 0xae, 0x6b, 0x20,
	0xde, 0xb4, 0xa1, 0xc3, 0xfa, 0xcd, 0x11, 0x38, 0x33, 0x89, 0x7c, 0xa1, 0xf4, 0xad, 0xd9, 0x41,
	0x15, 0x95, 0x78, 0x28, 0x8a, 0x9c, 0x94, 0x76, 0xa2, 0x1d, 0x54, 0x91, 0xf7, 0x19, 0xee, 0x36,
	0xb4, 0x2c, 0x2f, 0x99, 0x9d, 0x43, 0xbf, 0x56, 0x7d, 0x3d, 0xa4, 0xa5, 0x4d, 0x30, 0x6c, 0x1e,
	0xee, 0x36, 0x75, 0xb0, 0xf3, 0xca, 0xfb, 0x0a, 0xdd, 0xba, 0xea, 0x0f, 0x7e, 0x3b, 0xdd, 0xf2,
	0xdb, 0xa3, 0xbd, 0x54, 0x66, 0xc2, 0xda, 0x73, 0x2f, 0xe1, 0xb0, 0x99, 0x32, 0x5e, 0x3b, 0xbe,
	0xe1, 0xb5, 0xfb, 0xfb, 0xfb, 0x18, 0x6b, 0x7d, 0x81, 0x67, 0x3b, 0xce, 0x9a, 0xae, 0x2e, 0x44,
	0x9e, 0x90, 0x90, 0x18, 0xbd, 0xd7, 0x97, 0x5d, 0x37, 0xde, 0xfc, 0x12, 0xac, 0xbf, 0xfa, 0x25,
	0xcc, 0x1c, 0x9d, 0x38, 0xfd, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x33, 0x81, 0x11, 0x28, 0x83, 0x05,
	0x00, 0x00,
}
