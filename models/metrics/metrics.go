package metrics

import (
	"errors"
	"time"
)

type RawMetrics struct {
	TenantID string `json:"tenantId"`
	// ISO-8601 Intervals
	Interval    string              `json:"interval,omitempty"`
	Granularity string              `json:"granularity,omitempty"`
	Directions  []string            `json:"directions,omitempty"`
	Metrics     []string            `json:"metrics,omitempty"`
	ObjectType  string              `json:"objectType,omitempty"`
	Meta        map[string][]string `json:"meta,omitempty"`
	Timeout     int32               `json:"timeout,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (rm *RawMetrics) GetID() string {
	return "1"
}

// SetID - required implementation for jsonapi unmarshalling
func (rm *RawMetrics) SetID(s string) error {
	return nil
}

// DEPRECATED
type HistogramV1 struct {
	TenantID string              `json:"tenantId"`
	Meta     map[string][]string `json:"meta,omitempty"`
	// ISO-8601 Intervals
	Interval string `json:"interval,omitempty"`
	// ISO-8601 period combination
	Granularity          string                  `json:"granularity,omitempty"`
	MetricBucketRequests []MetricBucketRequestV1 `json:"metrics,omitempty"`
	// in Milliseconds
	Timeout int32 `json:"timeout,omitempty"`
}

type Histogram struct {
	TenantID string              `json:"tenantId"`
	Meta     map[string][]string `json:"meta,omitempty"`
	// ISO-8601 Intervals
	Interval string `json:"interval,omitempty"`
	// ISO-8601 period combination
	Granularity          string                `json:"granularity,omitempty"`
	MetricBucketRequests []MetricBucketRequest `json:"metrics,omitempty"`
	// in Milliseconds
	Timeout int32 `json:"timeout,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (h *Histogram) GetID() string {
	return "1"
}

// SetID - required implementation for jsonapi unmarshalling
func (h *Histogram) SetID(s string) error {
	return nil
}

type MetricBucketRequest struct {
	MetricIdentifierFilter
	Buckets []MetricBucket `json:"buckets"`
}

// DEPRECATED
type MetricBucketRequestV1 struct {
	Vendor     string           `json:"vendor,omitempty"`
	ObjectType string           `json:"objectType,omitempty"`
	Direction  string           `json:"direction"`
	Name       string           `json:"name"`
	Buckets    []MetricBucketV1 `json:"buckets"`
}

type MetricBucket struct {
	Lower *MetricBucketBoundarySpec `json:"lower"`
	Upper *MetricBucketBoundarySpec `json:"upper"`
}

type MetricBucketBoundarySpec struct {
	Value  float32 `json:"value"`
	Strict bool    `json:"strict"`
}

// DEPRECATED
type MetricBucketV1 struct {
	Index      string  `json:"index"`
	LowerBound float64 `json:"lower"`
	UpperBound float64 `json:"upper"`
}

type HistogramReport struct {
	ReportCompletionTime string                     `json:"reportCompletionTime"`
	TenantID             string                     `json:"tenantId"`
	Meta                 map[string][]string        `json:"meta"`
	ReportTimeRange      string                     `json:"reportTimeRange"`
	TimeSeriesResult     []HistogramTimeSeriesEntry `json:"timeSeriesResult"`
}

// DEPRECATED
type HistogramReportV1 struct {
	ReportCompletionTime string                       `json:"reportCompletionTime"`
	TenantID             string                       `json:"tenantId"`
	Meta                 map[string][]string          `json:"meta"`
	ReportTimeRange      string                       `json:"reportTimeRange"`
	TimeSeriesResult     []HistogramTimeSeriesEntryV1 `json:"timeSeriesResult"`
}

type HistogramTimeSeriesEntry struct {
	Timestamp string         `json:"timestamp"`
	Result    []MetricResult `json:"result"`
}

// DEPRECATED
type HistogramTimeSeriesEntryV1 struct {
	Timestamp string           `json:"timestamp"`
	Result    []MetricResultV1 `json:"result"`
}

type MetricResult struct {
	MetricIdentifier
	Results []BucketResult `json:"result"`
}

// DEPRECATED
type MetricResultV1 struct {
	Vendor     string         `json:"vendor,omitempty"`
	ObjectType string         `json:"objectType,omitempty"`
	Direction  string         `json:"direction"`
	Name       string         `json:"name"`
	Results    []BucketResult `json:"result"`
}

type BucketResult struct {
	Index string `json:"index"`
	Count int    `json:"count"`
}

type ReportSummary struct {
	TotalDuration          int64       `json:"totalDuration"`
	TotalViolationCount    int32       `json:"totalViolationCount"`
	TotalViolationDuration int64       `json:"totalViolationDuration"`
	SLACompliancePercent   float32     `json:"slaCompliancePercent"`
	ObjectCount            int32       `json:"objectCount"`
	PerMetricSummary       interface{} `json:"perMetricSummary"`
}

type TimeSeriesEntry struct {
	Timestamp string           `json:"timestamp"`
	Result    TimeSeriesResult `json:"result"`
}

type TimeSeriesResult struct {
	TotalDuration          int64       `json:"totalDuration"`
	TotalViolationCount    int32       `json:"totalViolationCount"`
	TotalViolationDuration int64       `json:"totalViolationDuration"`
	PerMetricResult        interface{} `json:"perMetricResult"`
}

type ThresholdCrossingTimeSeriesEntry struct {
	Timestamp string                            `json:"timestamp"`
	Result    ThresholdCrossingTimeSeriesResult `json:"result"`
}

type ThresholdCrossingTimeSeriesResult struct {
	ByMetric   []*ThresholdCrossingMetricResult  `json:"byMetric"`
	BySeverity map[string]map[string]interface{} `json:"bySeverity"`
}

type ThresholdCrossingMetricResult struct {
	ObjectType    string                            `json:"objectType"`
	Direction     string                            `json:"direction"`
	Metric        string                            `json:"metric"`
	Vendor        string                            `json:"vendor"`
	TotalDuration float64                           `json:"totalDuration"`
	BySeverity    map[string]map[string]interface{} `json:"bySeverity"`
}

// DEPRECATED
type ThresholdCrossingTopNV1 struct {
	Metric   MetricIdentifierV1 `json:"metric"`
	TenantID string             `json:"tenantId"`
	// ISO-8601 Intervals
	Interval string              `json:"interval,omitempty"`
	Meta     map[string][]string `json:"meta,omitempty"`
	// ISO-8601 period combination
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`
	Granularity        string `json:"granularity,omitempty"`
	Timeout            int32  `json:"timeout,omitempty"`
	NumResults         int32  `json:"numResults,omitempty"`
}

type ThresholdCrossingTopN struct {
	Metric   MetricIdentifierFilter `json:"metric"`
	TenantID string                 `json:"tenantId"`
	// ISO-8601 Intervals
	Interval string              `json:"interval,omitempty"`
	Meta     map[string][]string `json:"meta,omitempty"`
	// ISO-8601 period combination
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`
	Granularity        string `json:"granularity,omitempty"`
	Timeout            int32  `json:"timeout,omitempty"`
	NumResults         int32  `json:"numResults,omitempty"`
	// Indicates whether the results should be in ascending or descending order
	Sorted string `json:"sorted,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (tctn *ThresholdCrossingTopN) GetID() string {
	return "1"
}

// SetID - required implementation for jsonapi unmarshalling
func (tctn *ThresholdCrossingTopN) SetID(s string) error {
	return nil
}

func (tctn *ThresholdCrossingTopN) GetName() string {
	return "thresholdCrossingByMOTopNs"
}

// DEPRECATED
type TopNForMetricV1 struct {
	Meta map[string][]string `json:"meta,omitempty"`
	// List of monitored objects (optional)
	MonitoredObjects []string `json:"monitoredObjects,omitempty"`

	// Required Time range for the requestin ISO 8601 format for intervals
	Interval string `json:"interval,omitempty"`
	// Rquired Vendor (to avoid overlaps, eg: flowmeter does not have Jitter values
	// so if you do a min TopN then you'll just get a list of 0s)
	TenantID string `json:"tenant,omitempty"`
	// Timeout for the request
	Timeout int32 `json:"timeout,omitempty"`
	// Number of Results (default is 10)
	NumResult int32 `json:"NumResults,omitempty"`

	// Operation - 'avg', 'min', 'max'
	Aggregator string `json:"aggregator,omitempty"`
	// Metric that we are apply Aggregation to
	Metric MetricIdentifierV1 `json:"metric,omitempty"`

	// Metrics that are related and interesting BUT are NOT part of the post aggregation
	MetricsView []MetricAggregation `json:"metricsView,omitempty"`
}

type TopNForMetric struct {
	Meta map[string][]string `json:"meta,omitempty"`
	// List of monitored objects (optional)
	MonitoredObjects []string `json:"monitoredObjects,omitempty"`

	// Required Time range for the requestin ISO 8601 format for intervals
	Interval string `json:"interval,omitempty"`
	// Rquired Vendor (to avoid overlaps, eg: flowmeter does not have Jitter values
	// so if you do a min TopN then you'll just get a list of 0s)
	TenantID string `json:"tenant,omitempty"`
	// Timeout for the request
	Timeout int32 `json:"timeout,omitempty"`
	// Number of Results (default is 10)
	NumResult int32 `json:"NumResults,omitempty"`

	// Operation - 'avg', 'min', 'max'
	Aggregator string `json:"aggregator,omitempty"`
	// Metric that we are apply Aggregation to
	Metric MetricIdentifierFilter `json:"metric,omitempty"`

	// Metrics that are related and interesting BUT are NOT part of the post aggregation
	MetricsView []MetricAggregation `json:"metricsView,omitempty"`

	// Indicates whether the results should be in ascending or descending order
	Sorted string `json:"sorted,omitempty"`
}

// DEPRECATED - Remove once V1 is removed
func (tpn *TopNForMetricV1) Validate() (*TopNForMetricV1, error) {
	req := tpn
	if req.Timeout == 0 {
		req.Timeout = 5000
	}

	if tpn.NumResult == 0 {
		tpn.NumResult = 10
	}

	if len(req.TenantID) == 0 {
		return nil, errors.New("Tenant must not be empty.")
	}

	if len(req.Interval) == 0 {
		return nil, errors.New("Interval must not be empty")
	}

	return req, nil
}

// GetID - required implementation for jsonapi marshalling
func (tm *TopNForMetric) GetID() string {
	return "1"
}

// SetID - required implementation for jsonapi unmarshalling
func (tm *TopNForMetric) SetID(s string) error {
	return nil
}

type MetricAggregation struct {
	// Metric name eg jitterP95
	Metric string `json:"metric,omitempty"`
	// Operation - 'sum', 'count', 'min', 'max'
	Aggregator string `json:"aggregator,omitempty"`
	// Name for this Aggregation (must be unique)
	Name string `json:"name,omitempty"`
}

func (tpn *TopNForMetric) Validate() (*TopNForMetric, error) {
	req := tpn
	if req.Timeout == 0 {
		req.Timeout = 5000
	}

	if tpn.NumResult == 0 {
		tpn.NumResult = 10
	}

	if len(req.TenantID) == 0 {
		return nil, errors.New("Tenant must not be empty.")
	}

	if len(req.Interval) == 0 {
		return nil, errors.New("Interval must not be empty")
	}

	return req, nil
}

type ThresholdCrossing struct {
	TenantID           string                   `json:"tenantId"`
	Meta               map[string][]string      `json:"meta,omitempty"`
	Interval           string                   `json:"interval,omitempty"`
	Granularity        string                   `json:"granularity,omitempty"`
	ThresholdProfileID string                   `json:"thresholdProfileId,omitempty"`
	Metrics            []MetricIdentifierFilter `json:"metrics,omitempty"`
	Timeout            int32                    `json:"timeout,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (tcr *ThresholdCrossing) GetID() string {
	return "1"
}

// SetID - required implementation for jsonapi unmarshalling
func (tcr *ThresholdCrossing) SetID(s string) error {
	return nil
}

// DEPRECATED
type ThresholdCrossingV1 struct {
	TenantID           string               `json:"tenantId"`
	Meta               map[string][]string  `json:"meta,omitempty"`
	Interval           string               `json:"interval,omitempty"`
	Granularity        string               `json:"granularity,omitempty"`
	ThresholdProfileID string               `json:"thresholdProfileId,omitempty"`
	Metrics            []MetricIdentifierV1 `json:"metrics,omitempty"`
	Timeout            int32                `json:"timeout,omitempty"`
}

// DEPRECATED
type AggregateMetricsV1 struct {
	TenantID    string               `json:"tenantId"`
	Meta        map[string][]string  `json:"meta,omitempty"`
	Interval    string               `json:"interval,omitempty"`
	Granularity string               `json:"granularity,omitempty"`
	Timeout     int32                `json:"timeout,omitempty"`
	Aggregation AggregationSpecV1    `json:"aggregation"`
	Metrics     []MetricIdentifierV1 `json:"metrics,omitempty"`
}

type AggregateMetrics struct {
	TenantID         string                   `json:"tenantId"`
	Meta             map[string][]string      `json:"meta,omitempty"`
	Interval         string                   `json:"interval,omitempty"`
	Granularity      string                   `json:"granularity,omitempty"`
	Timeout          int32                    `json:"timeout,omitempty"`
	Aggregation      string                   `json:"aggregation"`
	Metrics          []MetricIdentifierFilter `json:"metrics,omitempty"`
	MonitoredObjects []string                 `json:"monitoredObjects,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (am *AggregateMetrics) GetID() string {
	return "1"
}

// SetID - required implementation for jsonapi unmarshalling
func (am *AggregateMetrics) SetID(s string) error {
	return nil
}

// DEPRECATED
type AggregationSpecV1 struct {
	Name string `json:"name"`
}

type MetricIdentifier struct {
	Vendor     string `json:"vendor"`
	ObjectType string `json:"objectType"`
	Metric     string `json:"metric"`
	Direction  string `json:"direction"`
}

type MetricIdentifierFilter struct {
	Vendor     string   `json:"vendor"`
	ObjectType []string `json:"objectType"`
	Metric     string   `json:"metric"`
	Direction  []string `json:"direction"`
}

// DEPRECATED
type MetricIdentifierV1 struct {
	Vendor     string `json:"vendor"`
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
	Direction  int32  `json:"direction"`
}

type TimeseriesEntryResponse struct {
	Timestamp string                 `json:"timestamp"`
	Result    map[string]interface{} `json:"result"`
}

type TimeseriesEntryResponseV1 struct {
	Timestamp string                 `json:"Timestamp"`
	Result    map[string]interface{} `json:"Result"`
}

type TopNEntryResponse struct {
	MonitoredObjectId string
	Result            map[string]interface{}
}

// DruidTimeSeriesResponse - Druid's Time Series response entry. Results are numbers.
type DruidTimeSeriesResponse struct {
	Timestamp string             `json:"timestamp"`
	Result    map[string]float64 `json:"result"`
}

// DruidTopNResponseEntry - Druid's TopN response. Note that unlike timeseries, results can be strings or numbers.
type DruidTopNResponseEntry map[string]interface{}

// DruidTopNResponse - TopN responses are organized by time buckets and each bucket has a topN series
type DruidTopNResponse []struct {
	Timestamp time.Time                `json:"timestamp"`
	Result    []DruidTopNResponseEntry `json:"result"`
}

// DruidViolationsMap - metadata for druid responses.
// We found an issue where druid would return a time series but we had a
// hard time relating the output to the query columns.
// eg: {
// timestamp: 123,
// value0: 123,
// value2: 432,
// }
// The problem from above is that it isn't easy to corrolate value0 to a specific  metric.
// So DruidMetricViolationResponse structure is a way to do that by forcing us to define each metric with a key
// which we can map to the results.
type DruidViolationsMap map[string]*MetricViolationsTimeSeries

// AddMetric - Utility to help simplify adding a metric to track in the druid response.
func (dmv *DruidViolationsMap) AddMetric(key, metric, name, ctype, vendor, objectType, direction string) *DruidViolationsMap {
	if dmv == nil {
		x := make(DruidViolationsMap)
		dmv = &x
	}

	obj := MetricViolationsTimeSeries{
		Metric:     metric,
		Name:       name,
		Type:       ctype,
		Vendor:     vendor,
		ObjectType: objectType,
		Direction:  direction,
	}
	(*dmv)[key] = &obj
	return dmv
}

// Merge - Used to merge two DruidViolationsMap together.
func (dmv *DruidViolationsMap) Merge(src DruidViolationsMap) *DruidViolationsMap {

	for k, v := range src {
		(*dmv)[k] = v
	}
	return dmv
}

// MetricViolationsTimeSeries - When doing SLA Violation queries, we have a naming scheme and need a way to correlate the results with the query
type MetricViolationsTimeSeries struct {
	Name           string                                 `json:"-,omitempty"`
	Type           string                                 `json:"-,omitempty"`
	Vendor         string                                 `json:"vendor"`
	ObjectType     string                                 `json:"objectType"`
	Metric         string                                 `json:"metric"`
	Direction      string                                 `json:"direction"`
	InternalSeries map[string]*MetricViolationSummaryType `json:"-,omitempty"`
	Totals         *MetricViolationSummaryType            `json:"total,omitempty"`
	ByGranularity  []*MetricViolationSummaryType          `json:"byGranularity,omitempty"`
	ByHourPerDay   []*MetricViolationSummaryType          `json:"byHourOfDay,omitempty"`
	ByDayPerWeek   []*MetricViolationSummaryType          `json:"byDayOfWeek,omitempty"`
	Critical       []*MetricViolationSummaryType          `json:"critical,omitempty"`
	Major          []*MetricViolationSummaryType          `json:"major,omitempty"`
	Minor          []*MetricViolationSummaryType          `json:"minor,omitempty"`
	Warning        []*MetricViolationSummaryType          `json:"warning,omitempty"`
	SLA            []*MetricViolationSummaryType          `json:"sla,omitempty"`
}

// DruidResponse2TimeSeriesMap - map of MetricViolationsTimeSeries which makes our life easier to manage monitored objects
type DruidResponse2TimeSeriesMap map[string]*MetricViolationsTimeSeries

// ToArray - converts a map to an array
func (drts *DruidResponse2TimeSeriesMap) ToArray() []*MetricViolationsTimeSeries {
	var res []*MetricViolationsTimeSeries
	for _, v := range *drts {
		res = append(res, v)
	}
	return res
}

// Put - Adds a metric value to a timeseries
func (drts *DruidResponse2TimeSeriesMap) Put(key string, subkey string, schemaEntry *MetricViolationsTimeSeries, value interface{}) (DruidResponse2TimeSeriesMap, error) {

	prerender := *drts
	if prerender[key] == nil {
		prerender[key] = &MetricViolationsTimeSeries{
			Direction:      schemaEntry.Direction,
			ObjectType:     schemaEntry.ObjectType,
			Vendor:         schemaEntry.Vendor,
			Metric:         schemaEntry.Metric,
			InternalSeries: make(map[string]*MetricViolationSummaryType),
		}
	}
	c := prerender[key]

	if c.InternalSeries[subkey] == nil {
		c.InternalSeries[subkey] = &MetricViolationSummaryType{}
	}

	cs := *(c.InternalSeries[subkey])
	// logger.Log.Debugf("k:%s val:%v",
	// key, value)

	cs[schemaEntry.Type] = value

	// prerender[key] = c

	// drts = &prerender

	return prerender, nil
}

type MetricViolationsSummaryAsTimeSeriesEntry map[string]interface{}

// MetricViolationsAsTimeSeries - SLA violations that have specific granularities will return a timeseries
type MetricViolationsAsTimeSeries struct {
	SummaryResult   map[string]*MetricViolationsSummaryAsTimeSeriesEntry `json:"totals,omitempty"`
	PerMetricResult DruidResponse2TimeSeriesMap                          `json:"byMetric"`
}

// MetricViolationsAsSummary - Similar to MetricViolationsAsTimeSeries but with a exlicit single entry
type MetricViolationsAsSummary struct {
	Summary         MetricViolationSummaryType  `json:"summary"`
	PerMetricResult DruidResponse2TimeSeriesMap `json:"metric"`
}

// MetricViolationSummaryType - Report summaries are a bit weird so rather than trying to map it, we just collect it into an interface.
type MetricViolationSummaryType map[string]interface{}
