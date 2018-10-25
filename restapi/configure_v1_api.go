package restapi

import (
	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/restapi/operations"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service"
	middleware "github.com/go-openapi/runtime/middleware"
)

func configureAdminServiceV1API(api *operations.GatherAPI, adminDB datastore.AdminServiceDatastore, tenantDB datastore.TenantServiceDatastore) {

	api.AdminProvisioningServiceCreateTenantHandler = admin_provisioning_service.CreateTenantHandlerFunc(handlers.HandleCreateTenant(handlers.SkylightAdminRoleOnly, adminDB, tenantDB))
	api.AdminProvisioningServiceDeleteTenantHandler = admin_provisioning_service.DeleteTenantHandlerFunc(handlers.HandleDeleteTenant(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceGetAllTenantsHandler = admin_provisioning_service.GetAllTenantsHandlerFunc(handlers.HandleGetAllTenants(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceGetIngestionDictionaryHandler = admin_provisioning_service.GetIngestionDictionaryHandlerFunc(handlers.HandleGetIngestionDictionary(handlers.AllRoles, adminDB))
	api.AdminProvisioningServiceGetTenantHandler = admin_provisioning_service.GetTenantHandlerFunc(handlers.HandleGetTenant(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceGetTenantSummaryByAliasHandler = admin_provisioning_service.GetTenantSummaryByAliasHandlerFunc(handlers.HandleGetTenantSummaryByAlias(adminDB))
	api.AdminProvisioningServiceGetValidTypesHandler = admin_provisioning_service.GetValidTypesHandlerFunc(handlers.HandleGetValidTypes(handlers.AllRoles, adminDB))
	api.AdminProvisioningServiceGetTenantIDByAliasHandler = admin_provisioning_service.GetTenantIDByAliasHandlerFunc(handlers.HandleGetTenantIDByAlias(adminDB))
	api.AdminProvisioningServicePatchTenantHandler = admin_provisioning_service.PatchTenantHandlerFunc(handlers.HandlePatchTenant(handlers.SkylightAdminRoleOnly, adminDB))
}

func configureTenantServiceV1API(api *operations.GatherAPI, tenantDB datastore.TenantServiceDatastore) {
	api.TenantProvisioningServiceBulkInsertMonitoredObjectHandler = tenant_provisioning_service.BulkInsertMonitoredObjectHandlerFunc(handlers.HandleBulkInsertMonitoredObjects(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceBulkUpdateMonitoredObjectHandler = tenant_provisioning_service.BulkUpdateMonitoredObjectHandlerFunc(handlers.HandleBulkUpdateMonitoredObjects(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantIngestionProfileHandler = tenant_provisioning_service.CreateTenantIngestionProfileHandlerFunc(handlers.HandleCreateTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateReportScheduleConfigHandler = tenant_provisioning_service.CreateReportScheduleConfigHandlerFunc(handlers.HandleCreateReportScheduleConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantConnectorConfigHandler = tenant_provisioning_service.CreateTenantConnectorConfigHandlerFunc(handlers.HandleCreateTenantConnectorConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantConnectorInstanceHandler = tenant_provisioning_service.CreateTenantConnectorInstanceHandlerFunc(handlers.HandleCreateTenantConnectorInstance(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantDomainHandler = tenant_provisioning_service.CreateTenantDomainHandlerFunc(handlers.HandleCreateTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantMetadataHandler = tenant_provisioning_service.CreateTenantMetadataHandlerFunc(handlers.HandleCreateTenantMetadata(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceCreateTenantMonitoredObjectHandler = tenant_provisioning_service.CreateTenantMonitoredObjectHandlerFunc(handlers.HandleCreateTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantThresholdProfileHandler = tenant_provisioning_service.CreateTenantThresholdProfileHandlerFunc(handlers.HandleCreateTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteReportScheduleConfigHandler = tenant_provisioning_service.DeleteReportScheduleConfigHandlerFunc(handlers.HandleDeleteReportScheduleConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantConnectorConfigHandler = tenant_provisioning_service.DeleteTenantConnectorConfigHandlerFunc(handlers.HandleDeleteTenantConnectorConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantConnectorInstanceHandler = tenant_provisioning_service.DeleteTenantConnectorInstanceHandlerFunc(handlers.HandleDeleteTenantConnectorInstance(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantDomainHandler = tenant_provisioning_service.DeleteTenantDomainHandlerFunc(handlers.HandleDeleteTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantIngestionProfileHandler = tenant_provisioning_service.DeleteTenantIngestionProfileHandlerFunc(handlers.HandleDeleteTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantMetadataHandler = tenant_provisioning_service.DeleteTenantMetadataHandlerFunc(handlers.HandleDeleteTenantMetadata(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceDeleteTenantMonitoredObjectHandler = tenant_provisioning_service.DeleteTenantMonitoredObjectHandlerFunc(handlers.HandleDeleteTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantThresholdProfileHandler = tenant_provisioning_service.DeleteTenantThresholdProfileHandlerFunc(handlers.HandleDeleteTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceGetActiveTenantIngestionProfileHandler = tenant_provisioning_service.GetActiveTenantIngestionProfileHandlerFunc(handlers.HandleGetActiveTenantIngestionProfile(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllReportScheduleConfigHandler = tenant_provisioning_service.GetAllReportScheduleConfigHandlerFunc(handlers.HandleGetAllReportScheduleConfigs(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllSLAReportsHandler = tenant_provisioning_service.GetAllSLAReportsHandlerFunc(handlers.HandleGetAllSLAReports(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantConnectorConfigsHandler = tenant_provisioning_service.GetAllTenantConnectorConfigsHandlerFunc(handlers.HandleGetAllTenantConnectorConfigs(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantConnectorInstancesHandler = tenant_provisioning_service.GetAllTenantConnectorInstancesHandlerFunc(handlers.HandleGetAllTenantConnectorInstances(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantDomainsHandler = tenant_provisioning_service.GetAllTenantDomainsHandlerFunc(handlers.HandleGetAllTenantDomains(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantMonitoredObjectsHandler = tenant_provisioning_service.GetAllTenantMonitoredObjectsHandlerFunc(handlers.HandleGetAllTenantMonitoredObjects(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantThresholdProfilesHandler = tenant_provisioning_service.GetAllTenantThresholdProfilesHandlerFunc(handlers.HandleGetAllTenantThresholdProfiles(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetDomainToMonitoredObjectMapHandler = tenant_provisioning_service.GetDomainToMonitoredObjectMapHandlerFunc(handlers.HandleGetDomainToMonitoredObjectMap(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetReportScheduleConfigHandler = tenant_provisioning_service.GetReportScheduleConfigHandlerFunc(handlers.HandleGetReportScheduleConfig(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetSLAReportHandler = tenant_provisioning_service.GetSLAReportHandlerFunc(handlers.HandleGetSLAReport(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantConnectorConfigHandler = tenant_provisioning_service.GetTenantConnectorConfigHandlerFunc(handlers.HandleGetTenantConnectorConfig(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantConnectorInstanceHandler = tenant_provisioning_service.GetTenantConnectorInstanceHandlerFunc(handlers.HandleGetTenantConnectorInstance(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantDomainHandler = tenant_provisioning_service.GetTenantDomainHandlerFunc(handlers.HandleGetTenantDomain(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantIngestionProfileHandler = tenant_provisioning_service.GetTenantIngestionProfileHandlerFunc(handlers.HandleGetTenantIngestionProfile(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantMetadataHandler = tenant_provisioning_service.GetTenantMetadataHandlerFunc(handlers.HandleGetTenantMetadata(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantMonitoredObjectHandler = tenant_provisioning_service.GetTenantMonitoredObjectHandlerFunc(handlers.HandleGetTenantMonitoredObject(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantThresholdProfileHandler = tenant_provisioning_service.GetTenantThresholdProfileHandlerFunc(handlers.HandleGetTenantThresholdProfile(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServicePatchTenantMetadataHandler = tenant_provisioning_service.PatchTenantMetadataHandlerFunc(handlers.HandlePatchTenantMetadata(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServicePatchTenantDomainHandler = tenant_provisioning_service.PatchTenantDomainHandlerFunc(handlers.HandlePatchTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServicePatchTenantIngestionProfileHandler = tenant_provisioning_service.PatchTenantIngestionProfileHandlerFunc(handlers.HandlePatchTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServicePatchTenantMonitoredObjectHandler = tenant_provisioning_service.PatchTenantMonitoredObjectHandlerFunc(handlers.HandlePatchTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServicePatchTenantThresholdProfileHandler = tenant_provisioning_service.PatchTenantThresholdProfileHandlerFunc(handlers.HandlePatchTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateReportScheduleConfigHandler = tenant_provisioning_service.UpdateReportScheduleConfigHandlerFunc(handlers.HandleUpdateReportScheduleConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantConnectorConfigHandler = tenant_provisioning_service.UpdateTenantConnectorConfigHandlerFunc(handlers.HandleUpdateTenantConnectorConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantConnectorInstanceHandler = tenant_provisioning_service.UpdateTenantConnectorInstanceHandlerFunc(handlers.HandleUpdateTenantConnectorInstance(handlers.SkylightAndTenantAdminRoles, tenantDB))

}

func configurev1APIThatWeMayRemove(api *operations.GatherAPI, tenantDB datastore.TenantServiceDatastore) {
	api.AdminProvisioningServiceUpdateTenantHandler = admin_provisioning_service.UpdateTenantHandlerFunc(func(params admin_provisioning_service.UpdateTenantParams) middleware.Responder {
		return middleware.NotImplemented("operation admin_provisioning_service.UpdateTenant has not yet been implemented")
	})
	api.TenantProvisioningServiceUpdateTenantDomainHandler = tenant_provisioning_service.UpdateTenantDomainHandlerFunc(handlers.HandleUpdateTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantIngestionProfileHandler = tenant_provisioning_service.UpdateTenantIngestionProfileHandlerFunc(handlers.HandleUpdateTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantMetadataHandler = tenant_provisioning_service.UpdateTenantMetadataHandlerFunc(handlers.HandleUpdateTenantMetadata(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceUpdateTenantMonitoredObjectHandler = tenant_provisioning_service.UpdateTenantMonitoredObjectHandlerFunc(handlers.HandleUpdateTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantThresholdProfileHandler = tenant_provisioning_service.UpdateTenantThresholdProfileHandlerFunc(handlers.HandleUpdateTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))

}
