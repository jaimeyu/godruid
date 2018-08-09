package restapi

import (
	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/restapi/operations"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
)

func configureAdminServiceV2API(api *operations.GatherAPI, adminDB datastore.AdminServiceDatastore, tenantDB datastore.TenantServiceDatastore) {
	api.AdminProvisioningServiceV2CreateTenantV2Handler = admin_provisioning_service_v2.CreateTenantV2HandlerFunc(handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB))
	api.AdminProvisioningServiceV2PatchTenantV2Handler = admin_provisioning_service_v2.PatchTenantV2HandlerFunc(handlers.HandlePatchTenantV2(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceV2GetTenantV2Handler = admin_provisioning_service_v2.GetTenantV2HandlerFunc(handlers.HandleGetTenantV2(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceV2GetAllTenantsV2Handler = admin_provisioning_service_v2.GetAllTenantsV2HandlerFunc(handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceV2DeleteTenantV2Handler = admin_provisioning_service_v2.DeleteTenantV2HandlerFunc(handlers.HandleDeleteTenantV2(handlers.SkylightAdminRoleOnly, adminDB))

	api.AdminProvisioningServiceV2GetTenantIDByAliasV2Handler = admin_provisioning_service_v2.GetTenantIDByAliasV2HandlerFunc(handlers.HandleGetTenantIDByAliasV2(adminDB))
	api.AdminProvisioningServiceV2GetTenantSummaryByAliasV2Handler = admin_provisioning_service_v2.GetTenantSummaryByAliasV2HandlerFunc(handlers.HandleGetTenantSummaryByAliasV2(adminDB))

	api.AdminProvisioningServiceV2GetIngestionDictionaryV2Handler = admin_provisioning_service_v2.GetIngestionDictionaryV2HandlerFunc(handlers.HandleGetIngestionDictionaryV2(handlers.AllRoles, adminDB))
	api.AdminProvisioningServiceV2GetValidTypesV2Handler = admin_provisioning_service_v2.GetValidTypesV2HandlerFunc(handlers.HandleGetValidTypesV2(handlers.AllRoles, adminDB))
}

func configureTenantServiceV2API(api *operations.GatherAPI, tenantDB datastore.TenantServiceDatastore) {
	api.TenantProvisioningServiceV2GetDataCleaningProfileHandler = tenant_provisioning_service_v2.GetDataCleaningProfileHandlerFunc(handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2GetDataCleaningProfilesHandler = tenant_provisioning_service_v2.GetDataCleaningProfilesHandlerFunc(handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteDataCleaningProfileHandler = tenant_provisioning_service_v2.DeleteDataCleaningProfileHandlerFunc(handlers.HandleDeleteDataCleaningProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceV2CreateDataCleaningProfileHandler = tenant_provisioning_service_v2.CreateDataCleaningProfileHandlerFunc(handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceV2UpdateDataCleaningProfileHandler = tenant_provisioning_service_v2.UpdateDataCleaningProfileHandlerFunc(handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))

	api.TenantProvisioningServiceV2GetAllMonitoredObjectsV2Handler = tenant_provisioning_service_v2.GetAllMonitoredObjectsV2HandlerFunc(handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2GetMonitoredObjectV2Handler = tenant_provisioning_service_v2.GetMonitoredObjectV2HandlerFunc(handlers.HandleGetMonitoredObjectV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2CreateMonitoredObjectV2Handler = tenant_provisioning_service_v2.CreateMonitoredObjectV2HandlerFunc(handlers.HandleCreateMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2UpdateMonitoredObjectV2Handler = tenant_provisioning_service_v2.UpdateMonitoredObjectV2HandlerFunc(handlers.HandleUpdateMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteMonitoredObjectV2Handler = tenant_provisioning_service_v2.DeleteMonitoredObjectV2HandlerFunc(handlers.HandleDeleteMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2BulkInsertMonitoredObjectsV2Handler = tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2HandlerFunc(handlers.HandleBulkCreateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2BulkUpdateMonitoredObjectsV2Handler = tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2HandlerFunc(handlers.HandleBulkUpdateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB))

	api.TenantProvisioningServiceV2GetAllConnectorConfigsV2Handler = tenant_provisioning_service_v2.GetAllConnectorConfigsV2HandlerFunc(handlers.HandleGetAllConnectorConfigsV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2GetConnectorConfigV2Handler = tenant_provisioning_service_v2.GetConnectorConfigV2HandlerFunc(handlers.HandleGetConnectorConfigV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2CreateConnectorConfigV2Handler = tenant_provisioning_service_v2.CreateConnectorConfigV2HandlerFunc(handlers.HandleCreateConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2UpdateConnectorConfigV2Handler = tenant_provisioning_service_v2.UpdateConnectorConfigV2HandlerFunc(handlers.HandleUpdateConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteConnectorConfigV2Handler = tenant_provisioning_service_v2.DeleteConnectorConfigV2HandlerFunc(handlers.HandleDeleteConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB))

	api.TenantProvisioningServiceV2GetAllConnectorInstancesV2Handler = tenant_provisioning_service_v2.GetAllConnectorInstancesV2HandlerFunc(handlers.HandleGetAllConnectorInstancesV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2GetConnectorInstanceV2Handler = tenant_provisioning_service_v2.GetConnectorInstanceV2HandlerFunc(handlers.HandleGetConnectorInstanceV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2CreateConnectorInstanceV2Handler = tenant_provisioning_service_v2.CreateConnectorInstanceV2HandlerFunc(handlers.HandleCreateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2UpdateConnectorInstanceV2Handler = tenant_provisioning_service_v2.UpdateConnectorInstanceV2HandlerFunc(handlers.HandleUpdateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteConnectorInstanceV2Handler = tenant_provisioning_service_v2.DeleteConnectorInstanceV2HandlerFunc(handlers.HandleDeleteConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB))

	api.TenantProvisioningServiceV2GetAllIngestionProfilesV2Handler = tenant_provisioning_service_v2.GetAllIngestionProfilesV2HandlerFunc(handlers.HandleGetAllIngestionProfilesV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2GetIngestionProfileV2Handler = tenant_provisioning_service_v2.GetIngestionProfileV2HandlerFunc(handlers.HandleGetIngestionProfileV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2CreateIngestionProfileV2Handler = tenant_provisioning_service_v2.CreateIngestionProfileV2HandlerFunc(handlers.HandleCreateIngestionProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceV2UpdateIngestionProfileV2Handler = tenant_provisioning_service_v2.UpdateIngestionProfileV2HandlerFunc(handlers.HandleUpdateIngestionProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteIngestionProfileV2Handler = tenant_provisioning_service_v2.DeleteIngestionProfileV2HandlerFunc(handlers.HandleDeleteIngestionProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))

	api.TenantProvisioningServiceV2GetAllThresholdProfilesV2Handler = tenant_provisioning_service_v2.GetAllThresholdProfilesV2HandlerFunc(handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2GetThresholdProfileV2Handler = tenant_provisioning_service_v2.GetThresholdProfileV2HandlerFunc(handlers.HandleGetThresholdProfileV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2CreateThresholdProfileV2Handler = tenant_provisioning_service_v2.CreateThresholdProfileV2HandlerFunc(handlers.HandleCreateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2UpdateThresholdProfileV2Handler = tenant_provisioning_service_v2.UpdateThresholdProfileV2HandlerFunc(handlers.HandleUpdateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteThresholdProfileV2Handler = tenant_provisioning_service_v2.DeleteThresholdProfileV2HandlerFunc(handlers.HandleDeleteThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))

	api.TenantProvisioningServiceV2GetAllReportScheduleConfigsV2Handler = tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2HandlerFunc(handlers.HandleGetAllReportScheduleConfigsV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2GetReportScheduleConfigV2Handler = tenant_provisioning_service_v2.GetReportScheduleConfigV2HandlerFunc(handlers.HandleGetReportScheduleConfigV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2CreateReportScheduleConfigV2Handler = tenant_provisioning_service_v2.CreateReportScheduleConfigV2HandlerFunc(handlers.HandleCreateReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2UpdateReportScheduleConfigV2Handler = tenant_provisioning_service_v2.UpdateReportScheduleConfigV2HandlerFunc(handlers.HandleUpdateReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteReportScheduleConfigV2Handler = tenant_provisioning_service_v2.DeleteReportScheduleConfigV2HandlerFunc(handlers.HandleDeleteReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB))

	api.TenantProvisioningServiceV2GetAllSLAReportsV2Handler = tenant_provisioning_service_v2.GetAllSLAReportsV2HandlerFunc(handlers.HandleGetAllSLAReportsV2(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceV2GetSLAReportV2Handler = tenant_provisioning_service_v2.GetSLAReportV2HandlerFunc(handlers.HandleGetSLAReportV2(handlers.AllRoles, tenantDB))
}

func configureMetricServiceV2API(api *operations.GatherAPI, druidDB datastore.DruidDatastore) {
}
