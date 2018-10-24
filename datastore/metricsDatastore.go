package datastore

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/adh-gather/swagmodels"
)

const (
	// ThresholdCrossingStr - common name of the ThresholdCrossingStr data type for use in logs.
	ThresholdCrossingStr = "Threshold Crossing"

	// TopNThresholdCrossingByMonitoredObjectStr - common name for use in logs.
	TopNThresholdCrossingByMonitoredObjectStr = "TopN Threshold Crossing by Monitored Object"

	// HistogramStr - common name for use in logs.
	HistogramStr = "Histogram"

	// RawMetricString - common name for use in logs.
	RawMetricStr = "Raw Metric"

	// AggMetricsStr - common name for use in logs.
	AggMetricsStr = "Agg Metric"

	// SLAReport - common name for use in logs.
	SLAReportStr = "SLA Report"

	// TopNForMetricStr - common name for use in logs
	TopNForMetricStr = "Top-N report"

	DataCleaningStr = "Data Cleaning History"
	HourOfDay       = 0
	DayOfWeek       = 1
)

const QueryDelimeter = "|"

type QueryKeySpec struct {
	KeySpecMap map[string]map[string]interface{}
}

func (qks *QueryKeySpec) AddKeySpec(keyspec map[string]interface{}) string {

	if len(qks.KeySpecMap) == 0 {
		qks.KeySpecMap = make(map[string]map[string]interface{})
	}

	// Sort any string arrays so that we always get the same hashed name
	for _, v := range keyspec {

		rType := reflect.TypeOf(v)

		switch rType.Kind() {
		case reflect.Slice:
			if rType.Elem().Kind() == reflect.String {
				sorted := v.([]string)
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i] < sorted[j]
				})
			}
		case reflect.Array:
			if rType.Elem().Kind() == reflect.String {
				sorted := v.([]string)
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i] < sorted[j]
				})
			}
		default:
		}

	}

	bytes, _ := json.Marshal(keyspec)
	id := fmt.Sprintf("%x", md5.Sum(bytes))
	qks.KeySpecMap[id] = keyspec

	return id
}

type MetricsDatastore interface {

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	QueryThresholdCrossing(request *metrics.ThresholdCrossing, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) (map[string]interface{}, error)

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	// DEPRECATED: DELETE WHEN V1 IS TERMINATED
	QueryThresholdCrossingV1(request *metrics.ThresholdCrossingV1, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) (map[string]interface{}, error)

	GetSLAReportV1(request *metrics.SLAReportRequest, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) (*metrics.SLAReport, error)
	GetSLAViolationsQueryAllGranularity(request *metrics.SLAReportRequest, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) ([]byte, metrics.DruidViolationsMap, error)
	GetSLAViolationsQueryWithGranularity(request *metrics.SLAReportRequest, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) ([]byte, metrics.DruidViolationsMap, error)
	GetSLATimeSeries(request *metrics.SLAReportRequest, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) (map[string]interface{}, map[string]interface{}, error)
	GetTopNTimeByBuckets(request *metrics.SLAReportRequest, extractFn int, vendor, objType, metric, direction, event string, eventAttr *pb.TenantThresholdProfileData_EventAttrMap,
		metaMOs []string) ([]byte, metrics.DruidViolationsMap, error)

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	// Uses TopN query.
	GetThresholdCrossingByMonitoredObjectTopN(request *metrics.ThresholdCrossingTopN, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) ([]metrics.TopNEntryResponse, error)

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	// Uses TopN query.
	// DEPRECATED: DELETE WHEN V1 IS TERMINATED
	GetThresholdCrossingByMonitoredObjectTopNV1(request *metrics.ThresholdCrossingTopNV1, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) (map[string]interface{}, error)

	// Returns the count for a set of specified metrics in set of specified buckets
	GetHistogram(request *metrics.Histogram, metaMOs []string) ([]metrics.TimeseriesEntryResponse, *QueryKeySpec, error)
	// Returns the count for a set of specified metrics in set of specified buckets
	// DEPRECATED: DELETE WHEN V1 IS TERMINATED
	GetHistogramV1(request *metrics.HistogramV1, metaMOs []string) (map[string]interface{}, error)

	// Returns raw metrics from metrics datastore
	// DEPRECATED: DELETE WHEN COLT IS NO LONGER USING
	GetFilteredRawMetrics(request *metrics.RawMetrics, metaMOs []string) (map[string]interface{}, error)
	// Returns raw metrics from metrics datastore
	// DEPRECATED: DELETE WHEN V1 IS TERMINATED
	GetRawMetricsV1(request *pb.RawMetricsRequest) (map[string]interface{}, error)

	// Get aggregated metrics from metrics datastore
	GetAggregatedMetrics(request *metrics.AggregateMetrics, metaMOs []string) ([]metrics.TimeseriesEntryResponse, *QueryKeySpec, error)
	// Get aggregated metrics from metrics datastore
	// DEPRECATED: DELETE WHEN V1 IS TERMINATED
	GetAggregatedMetricsV1(request *metrics.AggregateMetricsV1, metaMOs []string) (map[string]interface{}, error)

	// Get top N for specific metric from metrics datastore
	GetTopNForMetric(metric *metrics.TopNForMetric, metaMOs []string) ([]metrics.TopNEntryResponse, error)
	// Get top N for specific metric from metrics datastore
	// DEPRECATED: DELETE WHEN V1 IS TERMINATED
	GetTopNForMetricV1(metric *metrics.TopNForMetricV1, metaMOs []string) (map[string]interface{}, error)

	// Adds a monitored object to a metrics datastore look up
	AddMonitoredObjectToLookup(tenantID string, monitoredObjects []*tenmod.MonitoredObject, datatype string) error

	GetDataCleaningHistory(tenantID string, monitoredObjectID string, interval string) ([]*swagmodels.DataCleaningTransition, error)
}
