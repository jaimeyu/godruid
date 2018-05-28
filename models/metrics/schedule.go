package metrics

import "errors"

const (
	ReportScheduleConfigType = "reportScheduleConfig"
	ReportScheduleConfigStr  = "Report Schedule Config"
)

type ReportScheduleConfig struct {
	// Meta information
	ID                    string `json:"_id"`
	REV                   string `json:"_rev"`
	CreatedTimestamp      int64  `json:"createdTimestamp"`
	LastModifiedTimestamp int64  `json:"lastModifiedTimestamp"`
	TenantID              string `json:"tenantId"`

	// Report parameters
	ReportName         string   `json:"reportName"`
	TimeRangeDuration  int64    `json:"timeRangeDuration,omitempty"` // In nanoseconds
	DomainIds          []string `json:"domainIds,omitempty"`
	ThresholdProfileID string   `json:"thresholdProfileId,omitempty"`
	Granularity        string   `json:"granularity,omitempty"`
	Timeout            int32    `json:"timeout,omitempty"`  // In milliseconds
	Timezone           string   `json:"timezone,omitempty"` // See timezone strings defined as part of the IANA database https://www.iana.org/time-zones

	// Scheduling Execution timing
	// The values here are specified as strings to align to a crontab-like format
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
