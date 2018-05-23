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
	DatePeriodDays     int      `json:"datePeriodDays,omitempty"`
	Domain             []string `json:"domain,omitempty"`
	ThresholdProfileID string   `json:"thresholdProfileId,omitempty"`
	Granularity        string   `json:"granularity,omitempty"`
	Timeout            int32    `json:"timeout,omitempty"`
	//Timezone           int      //TODO

	// Scheduling Execution timing
	ReportName string `json:"reportName"`
	Minute     int    `json:"minute"`
	Hour       int    `json:"hour"`
	DayOfMonth int    `json:"dayMonth"`
	Month      int    `json:"month"`
	DayOfWeek  int    `json:"dayWeek"`
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
