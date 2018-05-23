package metrics

import "errors"

const (
	SLAScheduleConfigStr = "SLA Schedule Config"
)

type SLAReportRequest struct {
	TenantID string `json:"tenantId"`
	// ISO-8601 Intervals
	Interval string   `json:"interval,omitempty"`
	Domain   []string `json:"domain,omitempty"`
	// ISO-8601 period combination
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`
	Granularity        string `json:"granularity,omitempty"`
	Timeout            int32  `json:"timeout,omitempty"`
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

type SLAScheduleConfig struct {
	ID       string `json:"_id"`
	REV      string `json:"_rev"`
	TenantID string `json:"tenantId"`

	// Report parameters
	DatePeriodDays     int
	Domain             []string `json:"domain,omitempty"`
	ThresholdProfileID string   `json:"thresholdProfileId,omitempty"`
	Granularity        string   `json:"granularity,omitempty"`
	Timeout            int32    `json:"timeout,omitempty"`
	//Timezone           int      //TODO

	// Scheduling Execution timing
	ReportName string `json:"reportName"`
	Minute     string `json:"minute"`
	Hour       string `json:"hour"`
	DayOfMonth string `json:"dayMonth"`
	Month      string `json:"month"`
	DayOfWeek  string `json:"dayWeek"`
}

// GetID - required implementation for jsonapi marshalling
func (ssc *SLAScheduleConfig) GetID() string {
	return ssc.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (ssc *SLAScheduleConfig) SetID(s string) error {
	ssc.ID = s
	return nil
}

func (ssc *SLAScheduleConfig) Validate(isUpdate bool) error {
	if len(ssc.TenantID) == 0 {
		return errors.New("Invalid SLA Config request: must provide a Tenant ID")
	}

	return nil
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
