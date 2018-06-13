package metrics

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

type ThresholdCrossingTimeSeriesEntry struct {
	Timestamp string                            `json:"timestamp"`
	Result    ThresholdCrossingTimeSeriesResult `json:"result"`
}

type ThresholdCrossingTimeSeriesResult struct {
	TotalDuration          int64                             `json:"totalDuration"`
	TotalViolationCount    int32                             `json:"totalViolationCount"`
	TotalViolationDuration int64                             `json:"totalViolationDuration"`
	ByMetric               []*ThresholdCrossingMetricResult  `json:"byMetric"`
	BySeverity             map[string]map[string]interface{} `json:"bySeverity"`
}

type ThresholdCrossingMetricResult struct {
	ObjectType    string                            `json:"objectType"`
	Direction     string                            `json:"direction"`
	Metric        string                            `json:"metric"`
	Vendor        string                            `json:"vendor"`
	TotalDuration float64                           `json:"totalDuration"`
	BySeverity    map[string]map[string]interface{} `json:"bySeverity"`
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

type ThresholdCrossingRequest struct {
	TenantID           string             `json:"tenantId"`
	DomainIDs          []string           `json:"domainIds,omitempty"`
	Interval           string             `json:"interval,omitempty"`
	Granularity        string             `json:"granularity,omitempty"`
	ThresholdProfileID string             `json:"thresholdProfileId,omitempty"`
	MetricWhitelist    []MetricIdentifier `json:"metricWhitelist,omitempty"`
	Timeout            int32              `json:"timeout,omitempty"`
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
