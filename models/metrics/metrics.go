package metrics

type SLAReportRequest struct {
	TenantID string `json:"tenantId"`
	// ISO-8601 Intervals
	Interval string   `json:"interval,omitempty"`
	Domain   []string `json:"domain,omitempty"`
	// ISO-8601 period combination
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`
	Granularity        string `json:"granularity,omitempty"`
	Timeout            int32  `json:"timeout,omitempty"`
	Timezone           string `json:"timezone,omitempty"`
}

type SLAReport struct {
	ReportInstanceID     string            `json:"reportInstanceId"`
	ReportCompletionTime string            `json:"reportCompletionTime"`
	TenantID             string            `json:"tenantId"`
	ReportTimeRange      string            `json:"reportTimeRange"`
	SLASummary           SLASummary        `json:"slaSummary"`
	TimeSeriesResult     []TimeSeriesEntry `json:"timeSeriesResult"`
	ByHourOfDayResult    interface{}       `json:"byHourOfDayResult"`
	ByDayOfWeekResult    interface{}       `json:"byDayOfWeekResult"`
}

type SLASummary struct {
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
