package metrics

import "errors"
import admmod "github.com/accedian/adh-gather/models/admin"
import cron "github.com/robfig/cron"
import "fmt"

const (
	ReportScheduleConfigType = "reportScheduleConfig"
	ReportScheduleConfigStr  = "Report Schedule Config"
)

type ReportScheduleConfig struct {
	// Meta information
	ID                    string `json:"_id"`
	REV                   string `json:"_rev"`
	Datatype              string `json:"datatype"`
	CreatedTimestamp      int64  `json:"createdTimestamp"`
	LastModifiedTimestamp int64  `json:"lastModifiedTimestamp"`
	TenantID              string `json:"tenantId"`
	Active                bool   `json:"active"`

	// Report parameters
	TimeRangeDuration string              `json:"timeRangeDuration"`
	Meta              map[string][]string `json:"meta,omitempty"`
	ThresholdProfile  string              `json:"thresholdProfile,omitempty"`
	// Used for UI to help classify SLA Reports
	Name       string `json:"name,omitempty"`
	ReportType string `json:"reportType,omitempty"`

	// ISO8601 duration
	Granularity string `json:"granularity,omitempty"`
	// in Milliseconds
	Timeout int32 `json:"timeout,omitempty"`

	// The time format is in cron, see https://godoc.org/github.com/robfig/cron
	// We may be able to change to a single spec string that is literally just the cron spec.
	// eg: schedule string <- "0 0 * * * *"
	Second     string `json:"-"` // Do not expose, only used for testing
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

	if len(ssc.ThresholdProfile) == 0 {
		return errors.New("Invalid Report Schedule Config request: must provide a threshold profile")
	}

	if ssc.Second != "0" {
		ssc.Second = "0"
	}

	if ssc.Timeout < 5000 {
		ssc.Timeout = 5000
	}

	if ssc.Minute == "*" {
		// Avoid letting the scheduler run every 1 minute.
		return errors.New("Minutes cannot be wildcard.")
	}

	s := ssc
	s.Second = "0"

	spec := fmt.Sprintf("%s %s %s %s %s %s", s.Second, s.Minute, s.Hour, s.DayOfMonth, s.Month, s.DayOfWeek)
	//spec := fmt.Sprintf("0 * * * * *") // runs every 1 minute schedule
	//spec := fmt.Sprintf("0 0 3 * * 0") // runs every week on sunday at 3am EDT
	//loc := time.LoadLocation("Europe/London")
	//logger.Log.Infof("Testing cron job %s spec to '%s'", s.Name, spec)

	_, err := cron.Parse(spec)
	if err != nil {
		msg := fmt.Sprintf("Failure parsing cron spec: '%s'", err)
		return errors.New(msg)
	}

	return nil
}

// Interface to access the Schedule Configs from the DB
// Note current implementation uses the TENANT DB for configs storage!
type ScheduleDB interface {
	GetAllReportScheduleConfigs(tenantID string) ([]*ReportScheduleConfig, error)
	CreateSLAReport(report *SLAReport) (*SLAReport, error)
	DeleteReportScheduleConfig(tenantID string, configID string) (*ReportScheduleConfig, error)
}

// This interface may go away. We have a dependency to tenants' ID and in order
// to get a list of tenant IDs, we need to pull them from the AdminTenant table...
// So this is the interface to let us do this.
type AdminInterface interface {
	GetAllTenantDescriptors() ([]*admmod.Tenant, error)
}

// This is interface for the metric service handler.
type MetricServiceHandler interface {
	GetInternalSLAReport(request *SLAReportRequest) (*SLAReport, error)
}
