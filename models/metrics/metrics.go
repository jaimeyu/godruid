package metrics

import "errors"

const (
	ReportType = "reports"
	ReportStr  = "Report"
)

type SLAReportRequest struct {
	SlaScheduleConfig string `json:"slaScheduleConfigId"`
	TenantID          string `json:"tenantId"`
	// ISO-8601 Intervals
	Interval string   `json:"interval,omitempty"`
	Domain   []string `json:"domain,omitempty"`
	// ISO-8601 period combination
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`
	Granularity        string `json:"granularity,omitempty"`
	// in Milliseconds
	Timeout  int32  `json:"timeout,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (sr *SLAReport) GetID() string {
	return sr.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (sr *SLAReport) SetID(s string) error {
	sr.ID = s
	return nil
}

func (sr *SLAReport) GetName() string {
	return ReportType
}

type SLAReport struct {
	ID                   string            `json:"_id"`
	REV                  string            `json:"_rev"`
	ReportCompletionTime string            `json:"reportCompletionTime"`
	TenantID             string            `json:"tenantId"`
	ReportTimeRange      string            `json:"reportTimeRange"`
	ReportSummary        ReportSummary     `json:"reportSummary"`
	TimeSeriesResult     []TimeSeriesEntry `json:"timeSeriesResult"`
	ByHourOfDayResult    interface{}       `json:"byHourOfDayResult"`
	ByDayOfWeekResult    interface{}       `json:"byDayOfWeekResult"`
	ReportScheduleConfig string            `json:"reportScheduleConfig"`
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

type ThresholdCrossingTopNRequest struct {
	ObjectType string `json:"objectType"`
	Direction  string `json:"direction"`
	Metric     string `json:"metric"`
	Vendor     string `json:"vendor"`
	TenantID   string `json:"tenantId"`
	// ISO-8601 Intervals
	Interval string   `json:"interval,omitempty"`
	Domain   []string `json:"domain,omitempty"`
	// ISO-8601 period combination
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`
	Granularity        string `json:"granularity,omitempty"`
	Timeout            int32  `json:"timeout,omitempty"`
	NumResults         int32  `json:"numResults,omitempty"`
}

type TopNForMetric struct {
	// One of the two must be populated for the request to be valid, domains or monitoredObjects.
	// But if both are given, then the behaviour will be the query will be based on a subset of monitoredObjects that belong to the domains.
	// List of domains (optional)
	Domains []string `json:"domains,omitempty"`
	// List of monitored objects (optional)
	MonitoredObjects []string `json:"monitoredObjects,omitempty"`

	// Required Time range for the requestin ISO 8601 format for intervals
	Interval string `json:"interval,,omitempty"`
	// Rquired Vendor (to avoid overlaps, eg: flowmeter does not have Jitter values
	// so if you do a min TopN then you'll just get a list of 0s)
	TenantID string `json:"tenant,omitempty"`
	// Timeout for the request
	Timeout int32 `json:"timeout,omitempty"`
	// Number of Results (default is 10)
	NumResult int32 `json:"NumResults,omitempty"`

	// Operation - 'avg', 'min', 'max'
	Aggregator string `json:"aggregator,omitempty"`
	// Name of the metric for the aggregation
	Aggregation string `json:"aggregation,omitempty"`
	// Metric that we are apply Aggregation to
	Metric MetricIdentifier `json:"metric,omitempty"`

	// Metrics that are related and interesting BUT are NOT part of the post aggregation
	MetricsView []MetricAggregation `json:"metricsView,omitempty"`
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

	if len(req.Domains) == len(req.MonitoredObjects) && len(req.Domains) == 0 {
		return nil, errors.New("Either Domain or/and Monitored Objects list must not be empty.")
	}

	if len(req.TenantID) == 0 {
		return nil, errors.New("Tenant must not be empty.")
	}

	if len(req.Interval) == 0 {
		return nil, errors.New("Interval must not be empty")
	}

	return req, nil
}

type AggregateMetricsAPIRequest struct {
	TenantID    string             `json:"tenantId"`
	DomainIDs   []string           `json:"domainIds,omitempty"`
	Interval    string             `json:"interval,omitempty"`
	Granularity string             `json:"granularity,omitempty"`
	Timeout     int32              `json:"timeout,omitempty"`
	Aggregation AggregationSpec    `json:"aggregation"`
	Metrics     []MetricIdentifier `json:"metrics,omitempty"`
}

type AggregationSpec struct {
	Name string `json:"name"`
}

type MetricIdentifier struct {
	Vendor     string `json:"vendor"`
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
	Direction  int32  `json:"direction"`
}
