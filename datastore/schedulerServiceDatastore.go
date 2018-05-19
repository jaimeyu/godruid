package datastore

import metmod "github.com/accedian/adh-gather/models/metrics"

type SchedulerServiceDatastore interface {
	CreateSchedulerConfig(tenantID string, config *metmod.SLASchedulerConfig) (*metmod.SLASchedulerConfig, error)
	UpdateSchedulerConfig(tenantID string, config *metmod.SLASchedulerConfig) (*metmod.SLASchedulerConfig, error)
	DeleteSchedulerConfig(tenantID string, configID string) (*metmod.SLASchedulerConfig, error)
	GetSchedulerConfig(tenantID string, config *metmod.SLASchedulerConfig) (*metmod.SLASchedulerConfig, error)
}
