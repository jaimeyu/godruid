package datastore

import metmod "github.com/accedian/adh-gather/models/metrics"

type SchedulerServiceDatastore interface {
	CreateScheduleConfig(config *metmod.SLAScheduleConfig) (*metmod.SLAScheduleConfig, error)
	UpdateScheduleConfig(config *metmod.SLAScheduleConfig) (*metmod.SLAScheduleConfig, error)
	DeleteScheduleConfig(tenantID string, configID string) (*metmod.SLAScheduleConfig, error)
	GetScheduleConfig(tenantID string, configID string) (*metmod.SLAScheduleConfig, error)
}
