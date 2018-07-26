package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

const (
	// ThresholdCrossingStr - common name of the ThresholdCrossingStr data type for use in logs.
	ThresholdCrossingStr = "Threshold Crossing"

	QueryThresholdCrossingStr = "Threshold Crossing Query"

	// ThresholdCrossingByMonitoredObjectStr - common name for use in logs.
	ThresholdCrossingByMonitoredObjectStr = "Threshold Crossing by Monitored Object"

	// ThresholdCrossingByMonitoredObjectStr - common name for use in logs.
	TopNThresholdCrossingByMonitoredObjectStr = "TopN Threshold Crossing by Monitored Object"

	// HistogramStr - common name for use in logs.
	HistogramStr = "Histogram"

	// HistogramCustomStr - common name for use in logs.
	HistogramCustomStr = "HistogramCustom"

	// RawMetricString - common name for use in logs.
	RawMetricStr = "Raw Metric"

	// AggMetricsStr - common name for use in logs.
	AggMetricsStr = "Agg Metric"

	// SLAReport - common name for use in logs.
	SLAReportStr = "SLA Report"

	// TopNForMetricString - common name for use in logs
	TopNForMetricString = "Top-N report"
)

type DruidDatastore interface {

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	GetThresholdCrossing(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error)

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	QueryThresholdCrossing(request *metrics.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error)

	GetSLAReport(request *metrics.SLAReportRequest, thresholdProfile *pb.TenantThresholdProfile) (*metrics.SLAReport, error)

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	GetThresholdCrossingByMonitoredObject(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error)

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	// Uses TopN query.
	GetThresholdCrossingByMonitoredObjectTopN(request *metrics.ThresholdCrossingTopNRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error)

	// Returns the min,max,avg,median for a given metric
	GetHistogram(request *pb.HistogramRequest) (map[string]interface{}, error)

	// Returns the count for a set of specified metrics in set of specified buckets
	GetHistogramCustom(request *metrics.HistogramCustomRequest) (map[string]interface{}, error)

	// Returns raw metrics from druid
	GetRawMetrics(request *pb.RawMetricsRequest) (map[string]interface{}, error)

	// Get aggregated metrics from druid
	GetAggregatedMetrics(request *metrics.AggregateMetricsAPIRequest) (map[string]interface{}, error)
	GetTopNForMetric(metric *metrics.TopNForMetric) (map[string]interface{}, error)

	// Adds a monitored object to a druid look up
	AddMonitoredObjectToLookup(tenantID string, monitoredObjects []*tenmod.MonitoredObject, datatype string) error
}
