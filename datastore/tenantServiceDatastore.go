package datastore

import (
	"github.com/accedian/adh-gather/models/common"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// TenantServiceDatastore - interface which provides the functionality
// of the TenantService Datastore.
type TenantServiceDatastore interface {
	CreateTenantUser(*tenmod.User) (*tenmod.User, error)
	UpdateTenantUser(*tenmod.User) (*tenmod.User, error)
	DeleteTenantUser(tenantID string, userID string) (*tenmod.User, error)
	GetTenantUser(tenantID string, userID string) (*tenmod.User, error)
	GetAllTenantUsers(string) ([]*tenmod.User, error)

	CreateTenantDomain(*tenmod.Domain) (*tenmod.Domain, error)
	UpdateTenantDomain(*tenmod.Domain) (*tenmod.Domain, error)
	DeleteTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error)
	GetTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error)
	GetAllTenantDomains(string) ([]*tenmod.Domain, error)

	CreateTenantConnectorConfig(*tenmod.ConnectorConfig) (*tenmod.ConnectorConfig, error)
	UpdateTenantConnectorConfig(*tenmod.ConnectorConfig) (*tenmod.ConnectorConfig, error)
	DeleteTenantConnectorConfig(tenantID string, dataID string) (*tenmod.ConnectorConfig, error)
	GetTenantConnectorConfig(tenantID string, dataID string) (*tenmod.ConnectorConfig, error)
	GetAllTenantConnectorConfigs(tenantID, zone string) ([]*tenmod.ConnectorConfig, error)
	GetAllAvailableTenantConnectorConfigs(tenantID, zone string) ([]*tenmod.ConnectorConfig, error)
	GetAllTenantConnectorConfigsByInstanceID(tenantID, instanceID string) ([]*tenmod.ConnectorConfig, error)
	GetConnectorConfigUpdateChan() chan *tenmod.ConnectorConfig

	CreateTenantConnectorInstance(*tenmod.ConnectorInstance) (*tenmod.ConnectorInstance, error)
	UpdateTenantConnectorInstance(*tenmod.ConnectorInstance) (*tenmod.ConnectorInstance, error)
	DeleteTenantConnectorInstance(tenantID string, dataID string) (*tenmod.ConnectorInstance, error)
	GetTenantConnectorInstance(tenantID string, dataID string) (*tenmod.ConnectorInstance, error)
	GetAllTenantConnectorInstances(tenantID string) ([]*tenmod.ConnectorInstance, error)

	CreateTenantIngestionProfile(*tenmod.IngestionProfile) (*tenmod.IngestionProfile, error)
	UpdateTenantIngestionProfile(*tenmod.IngestionProfile) (*tenmod.IngestionProfile, error)
	GetTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error)
	DeleteTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error)
	GetActiveTenantIngestionProfile(tenantID string) (*tenmod.IngestionProfile, error)

	CreateTenantThresholdProfile(*tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error)
	UpdateTenantThresholdProfile(*tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error)
	GetTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error)
	DeleteTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error)
	GetAllTenantThresholdProfile(tenantID string) ([]*tenmod.ThresholdProfile, error)

	CreateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error)
	UpdateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error)
	GetMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error)
	DeleteMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error)
	GetAllMonitoredObjects(tenantID string) ([]*tenmod.MonitoredObject, error)
	GetMonitoredObjectToDomainMap(moByDomReq *tenmod.MonitoredObjectCountByDomainRequest) (*tenmod.MonitoredObjectCountByDomainResponse, error)
	BulkInsertMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error)
	BulkUpdateMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error)
	GetAllMonitoredObjectsInIDList(tenantID string, idList []string) ([]*tenmod.MonitoredObject, error)
	GetAllMonitoredObjectsByPage(tenantID string, startKey string, limit int64) ([]*tenmod.MonitoredObject, *common.PaginationOffsets, error)

	CreateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error)
	UpdateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error)
	DeleteTenantMeta(tenantID string) (*tenmod.Metadata, error)
	GetTenantMeta(tenantID string) (*tenmod.Metadata, error)

	CreateReportScheduleConfig(config *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error)
	UpdateReportScheduleConfig(config *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error)
	DeleteReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error)
	GetReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error)
	GetAllReportScheduleConfigs(tenantID string) ([]*metmod.ReportScheduleConfig, error)

	CreateSLAReport(slaReport *metmod.SLAReport) (*metmod.SLAReport, error)
	DeleteSLAReport(tenantID string, slaReportID string) (*metmod.SLAReport, error)
	GetSLAReport(tenantID string, slaReportID string) (*metmod.SLAReport, error)
	GetAllSLAReports(tenantID string) ([]*metmod.SLAReport, error)

	// For testing
	CreateDashboard(dashboard *tenmod.Dashboard) (*tenmod.Dashboard, error)
	// For testing
	DeleteDashboard(tenantID string, dataID string) (*tenmod.Dashboard, error)
	HasDashboardsWithDomain(tenantID string, domainID string) (bool, error)

	CreateTenantDataCleaningProfile(dcp *tenmod.DataCleaningProfile) (*tenmod.DataCleaningProfile, error)
	UpdateTenantDataCleaningProfile(dcp *tenmod.DataCleaningProfile) (*tenmod.DataCleaningProfile, error)
	GetTenantDataCleaningProfile(tenantID string, dataID string) (*tenmod.DataCleaningProfile, error)
	DeleteTenantDataCleaningProfile(tenantID string, dataID string) (*tenmod.DataCleaningProfile, error)
	GetAllTenantDataCleaningProfiles(tenantID string) ([]*tenmod.DataCleaningProfile, error)
}
