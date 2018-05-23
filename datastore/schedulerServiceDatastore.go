package datastore

import metmod "github.com/accedian/adh-gather/models/metrics"

type SchedulerServiceDatastore interface {
	CreateReportScheduleConfig(config *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error)
	UpdateReportScheduleConfig(config *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error)
	DeleteReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error)
	GetReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error)
	GetAllReportScheduleConfigs(tenantID string) ([]*metmod.ReportScheduleConfig, error)
}
