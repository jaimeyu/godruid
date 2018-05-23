package metrics

import "errors"

const (
	ReportScheduleConfigType = "reportScheduleConfigType"
	ReportScheduleConfigStr  = "Report Schedule Config"
)

type ReportScheduleConfig struct {
	ID       string `json:"_id"`
	REV      string `json:"_rev"`
	TenantID string `json:"tenantId"`

	// Report parameters
	DatePeriodDays     string   `json:"datePeriodDays,omitempty"`
	Domain             []string `json:"domain,omitempty"`
	ThresholdProfileID string   `json:"thresholdProfileId,omitempty"`
	Granularity        string   `json:"granularity,omitempty"`
	Timeout            int32    `json:"timeout,omitempty"`
	Timezone           string   `json:"timezone,omitempty"`

	// Scheduling Execution timing
	ReportName string `json:"reportName"`
	Minute     string `json:"minute"`
	Hour       string `json:"hour"`
	DayOfMonth string `json:"dayMonth"`
	Month      string `json:"month"`
	DayOfWeek  string `json:"dayWeek"`
}

// GetID - required implementation for jsonapi marshalling
func (ssc *ReportScheduleConfig) GetID() string {
	return ssc.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (ssc *ReportScheduleConfig) SetID(s string) error {
	ssc.ID = s
	return nil
}

func (ssc *ReportScheduleConfig) Validate(isUpdate bool) error {
	if len(ssc.TenantID) == 0 {
		return errors.New("Invalid Report Schedule Config request: must provide a Tenant ID")
	}

	return nil
}
