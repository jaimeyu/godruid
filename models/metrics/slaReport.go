package metrics

import (
	"errors"
	"fmt"

	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/manyminds/api2go/jsonapi"
	cron "github.com/robfig/cron"
)

const (
	// ReportStr - Used for logging
	ReportStr = "Report"
)
const (
	ReportScheduleConfigType = "reportScheduleConfig"
	ReportScheduleConfigStr  = "Report Schedule Config"
)

// SLAReportRequest - This is the model to hold the SLA Report request for immediate consumption & scheduled SLA reports
type SLAReportRequest struct {
	SLAScheduleConfig  string              `json:"slaScheduleConfigId,omitempty"`
	TenantID           string              `json:"tenantId"`
	Interval           string              `json:"interval,omitempty"` // ISO-8601 Intervals
	Meta               map[string][]string `json:"meta,omitempty"`
	ThresholdProfileID string              `json:"thresholdProfileId,omitempty"` // ISO-8601 period combination
	Granularity        string              `json:"granularity,omitempty"`
	Timeout            int32               `json:"timeout,omitempty"` // in Milliseconds

	Timezone string `json:"timezone,omitempty"`

	// v2 option, to specify a white list of metrics within the thresholdprofile to execute.
	// This reduces potential druid processing time.
	Metrics        []string `json:"metrics,omitempty"`
	IgnoreCleaning bool     `json:"ignoreCleaning,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (sr *SLAReportRequest) GetID() string {
	return "1"
}

// SetID - required implementation for jsonapi unmarshalling
func (sr *SLAReportRequest) SetID(s string) error {
	return nil
}

// GetName -Returns the report's type
func (sr *SLAReportRequest) GetName() string {
	return "slaReports"
}

// SLAReport - SLA Response structure for v1
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

// GetID - required implementation for jsonapi marshalling
func (sr *SLAReport) GetID() string {
	return sr.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (sr *SLAReport) SetID(s string) error {
	sr.ID = s
	return nil
}

// GetName - Returns the report's type
func (sr *SLAReport) GetName() string {
	return ReportStr
}

type SLAReportV2Result struct {
	Summary MetricViolationSummaryType    `json:"summary,omitempty"`
	Metric  []*MetricViolationsTimeSeries `json:"metric"`
}

// SLAReportV2 - SLA Response structure for v1
type SLAReportV2 struct {
	ID     string            `json:"-"`
	REV    string            `json:"-"`
	Config interface{}       `json:"config,omitempty"`
	Result SLAReportV2Result `json:"result,omitempty"`
}

// GetID - required implementation for jsonapi marshalling
func (sr *SLAReportV2) GetID() string {
	return sr.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (sr *SLAReportV2) SetID(s string) error {
	sr.ID = s
	return nil
}

// GetName - Returns the report's type
func (sr *SLAReportV2) GetName() string {
	return ReportStr
}

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

var (
	reportScheduleConfigTPRelationshipType = "thresholdProfiles"
	reportScheduleConfigTPRelationshipName = "thresholdProfile"
)

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (ssc *ReportScheduleConfig) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: reportScheduleConfigTPRelationshipType,
			Name: reportScheduleConfigTPRelationshipName,
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (ssc *ReportScheduleConfig) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	result = append(result, jsonapi.ReferenceID{
		ID:   ssc.ThresholdProfile,
		Type: reportScheduleConfigTPRelationshipType,
		Name: reportScheduleConfigTPRelationshipName,
	})

	return result
}

// SetToOneReferenceID - satisfy the unmarshalling of relationships that point to only one reference ID
func (ssc *ReportScheduleConfig) SetToOneReferenceID(name, ID string) error {
	if name == reportScheduleConfigTPRelationshipName {
		ssc.ThresholdProfile = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

// Validate - validates the model so its correct
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
	GetInternalSLAReportV1(request *SLAReportRequest) (*SLAReport, error)
}
