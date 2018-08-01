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

func configureTenantServiceV2API(api *operations.GatherAPI, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) {
	api.TenantProvisioningServiceV2GetDataCleaningProfileHandler = tenant_provisioning_service_v2.GetDataCleaningProfileHandlerFunc(handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2GetDataCleaningProfilesHandler = tenant_provisioning_service_v2.GetDataCleaningProfilesHandlerFunc(handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteDataCleaningProfileHandler = tenant_provisioning_service_v2.DeleteDataCleaningProfileHandlerFunc(handlers.HandleDeleteDataCleaningProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceV2CreateDataCleaningProfileHandler = tenant_provisioning_service_v2.CreateDataCleaningProfileHandlerFunc(handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceV2UpdateDataCleaningProfileHandler = tenant_provisioning_service_v2.UpdateDataCleaningProfileHandlerFunc(handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2GetDataCleaningHistoryHandler = tenant_provisioning_service_v2.GetDataCleaningHistoryHandlerFunc(handlers.HandleGetDataCleaningHistoryV2(handlers.SkylightAndTenantAdminRoles, druidDB))

	api.TenantProvisioningServiceV2GetAllMonitoredObjectsV2Handler = tenant_provisioning_service_v2.GetAllMonitoredObjectsV2HandlerFunc(handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB))
}

func configureMetricServiceV2API(api *operations.GatherAPI, druidDB datastore.DruidDatastore) {
}
