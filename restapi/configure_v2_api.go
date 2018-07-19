package restapi

import (
	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/restapi/operations"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
)

func configureAdminServiceV2API(api *operations.GatherAPI, adminDB datastore.AdminServiceDatastore, tenantDB datastore.TenantServiceDatastore) {

}

func configureTenantServiceV2API(api *operations.GatherAPI, tenantDB datastore.TenantServiceDatastore) {
	api.TenantProvisioningServiceV2GetDataCleaningProfileHandler = tenant_provisioning_service_v2.GetDataCleaningProfileHandlerFunc(handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2GetDataCleaningProfilesHandler = tenant_provisioning_service_v2.GetDataCleaningProfilesHandlerFunc(handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceV2DeleteDataCleaningProfileHandler = tenant_provisioning_service_v2.DeleteDataCleaningProfileHandlerFunc(handlers.HandleDeleteDataCleaningProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceV2CreateDataCleaningProfileHandler = tenant_provisioning_service_v2.CreateDataCleaningProfileHandlerFunc(handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAdminRoleOnly, tenantDB))
	api.TenantProvisioningServiceV2UpdateDataCleaningProfileHandler = tenant_provisioning_service_v2.UpdateDataCleaningProfileHandlerFunc(handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB))
}

func configureMetricServiceV2API(api *operations.GatherAPI, druidDB datastore.DruidDatastore) {
}
