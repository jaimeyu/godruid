// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/restapi/operations"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service"
	"github.com/accedian/adh-gather/restapi/operations/metrics_service"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service"
)

//go:generate swagger generate server --target .. --name gather --spec ../files/swagger.yml --model-package swagmodels --exclude-main --exclude-spec

func configureFlags(api *operations.GatherAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GatherAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TxtProducer = runtime.TextProducer()

	handlers.InitializeAuthHelper()

	// Create DAO objects to handle data retrieval as needed.
	adminDB, err := handlers.GetAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate Admin Service DAO: %s", err.Error())
	}

	tenantDB, err := handlers.GetTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate Tenant Service DAO: %s", err.Error())
	}

	api.TenantProvisioningServiceBulkInsertMonitoredObjectHandler = tenant_provisioning_service.BulkInsertMonitoredObjectHandlerFunc(func(params tenant_provisioning_service.BulkInsertMonitoredObjectParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.BulkInsertMonitoredObject has not yet been implemented")
	})
	api.TenantProvisioningServiceBulkUpdateMonitoredObjectHandler = tenant_provisioning_service.BulkUpdateMonitoredObjectHandlerFunc(func(params tenant_provisioning_service.BulkUpdateMonitoredObjectParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.BulkUpdateMonitoredObject has not yet been implemented")
	})

	api.AdminProvisioningServiceCreateIngestionDictionaryHandler = admin_provisioning_service.CreateIngestionDictionaryHandlerFunc(handlers.HandleCreateIngestionDictionary(handlers.SkylightAdminRoleOnly, adminDB))
	api.TenantProvisioningServiceCreateTenantIngestionProfileHandler = tenant_provisioning_service.CreateTenantIngestionProfileHandlerFunc(handlers.HandleCreateTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateReportScheduleConfigHandler = tenant_provisioning_service.CreateReportScheduleConfigHandlerFunc(func(params tenant_provisioning_service.CreateReportScheduleConfigParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.CreateReportScheduleConfig has not yet been implemented")
	})
	api.AdminProvisioningServiceCreateTenantHandler = admin_provisioning_service.CreateTenantHandlerFunc(handlers.HandleCreateTenant(handlers.SkylightAdminRoleOnly, adminDB, tenantDB))
	api.TenantProvisioningServiceCreateTenantConnectorConfigHandler = tenant_provisioning_service.CreateTenantConnectorConfigHandlerFunc(handlers.HandleCreateTenantConnectorConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantConnectorInstanceHandler = tenant_provisioning_service.CreateTenantConnectorInstanceHandlerFunc(handlers.HandleCreateTenantConnectorInstance(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantDomainHandler = tenant_provisioning_service.CreateTenantDomainHandlerFunc(handlers.HandleCreateTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantMetadataHandler = tenant_provisioning_service.CreateTenantMetadataHandlerFunc(func(params tenant_provisioning_service.CreateTenantMetadataParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.CreateTenantMetadata has not yet been implemented")
	})
	api.TenantProvisioningServiceCreateTenantMonitoredObjectHandler = tenant_provisioning_service.CreateTenantMonitoredObjectHandlerFunc(handlers.HandleCreateTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantThresholdProfileHandler = tenant_provisioning_service.CreateTenantThresholdProfileHandlerFunc(handlers.HandleCreateTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceCreateTenantUserHandler = tenant_provisioning_service.CreateTenantUserHandlerFunc(func(params tenant_provisioning_service.CreateTenantUserParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.CreateTenantUser has not yet been implemented")
	})
	api.AdminProvisioningServiceCreateValidTypesHandler = admin_provisioning_service.CreateValidTypesHandlerFunc(handlers.HandleCreateValidTypes(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceDeleteIngestionDictionaryHandler = admin_provisioning_service.DeleteIngestionDictionaryHandlerFunc(handlers.HandleDeleteIngestionDictionary(handlers.SkylightAdminRoleOnly, adminDB))
	api.TenantProvisioningServiceDeleteReportScheduleConfigHandler = tenant_provisioning_service.DeleteReportScheduleConfigHandlerFunc(func(params tenant_provisioning_service.DeleteReportScheduleConfigParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.DeleteReportScheduleConfig has not yet been implemented")
	})
	api.AdminProvisioningServiceDeleteTenantHandler = admin_provisioning_service.DeleteTenantHandlerFunc(handlers.HandleDeleteTenant(handlers.SkylightAdminRoleOnly, adminDB))
	api.TenantProvisioningServiceDeleteTenantConnectorConfigHandler = tenant_provisioning_service.DeleteTenantConnectorConfigHandlerFunc(handlers.HandleDeleteTenantConnectorConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantConnectorInstanceHandler = tenant_provisioning_service.DeleteTenantConnectorInstanceHandlerFunc(handlers.HandleDeleteTenantConnectorInstance(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantDomainHandler = tenant_provisioning_service.DeleteTenantDomainHandlerFunc(handlers.HandleDeleteTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantIngestionProfileHandler = tenant_provisioning_service.DeleteTenantIngestionProfileHandlerFunc(handlers.HandleDeleteTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantMetadataHandler = tenant_provisioning_service.DeleteTenantMetadataHandlerFunc(func(params tenant_provisioning_service.DeleteTenantMetadataParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.DeleteTenantMetadata has not yet been implemented")
	})
	api.TenantProvisioningServiceDeleteTenantMonitoredObjectHandler = tenant_provisioning_service.DeleteTenantMonitoredObjectHandlerFunc(handlers.HandleDeleteTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantThresholdProfileHandler = tenant_provisioning_service.DeleteTenantThresholdProfileHandlerFunc(handlers.HandleDeleteTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceDeleteTenantUserHandler = tenant_provisioning_service.DeleteTenantUserHandlerFunc(func(params tenant_provisioning_service.DeleteTenantUserParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.DeleteTenantUser has not yet been implemented")
	})
	api.AdminProvisioningServiceDeleteValidTypesHandler = admin_provisioning_service.DeleteValidTypesHandlerFunc(handlers.HandleDeleteValidTypes(handlers.SkylightAdminRoleOnly, adminDB))
	api.MetricsServiceGenSLAReportHandler = metrics_service.GenSLAReportHandlerFunc(func(params metrics_service.GenSLAReportParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.GenSLAReport has not yet been implemented")
	})
	api.TenantProvisioningServiceGetActiveTenantIngestionProfileHandler = tenant_provisioning_service.GetActiveTenantIngestionProfileHandlerFunc(handlers.HandleGetActiveTenantIngestionProfile(handlers.AllRoles, tenantDB))

	api.TenantProvisioningServiceGetAllReportScheduleConfigHandler = tenant_provisioning_service.GetAllReportScheduleConfigHandlerFunc(func(params tenant_provisioning_service.GetAllReportScheduleConfigParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetAllReportScheduleConfig has not yet been implemented")
	})
	api.TenantProvisioningServiceGetAllSLAReportsHandler = tenant_provisioning_service.GetAllSLAReportsHandlerFunc(func(params tenant_provisioning_service.GetAllSLAReportsParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetAllSLAReports has not yet been implemented")
	})
	api.TenantProvisioningServiceGetAllTenantConnectorConfigsHandler = tenant_provisioning_service.GetAllTenantConnectorConfigsHandlerFunc(handlers.HandleGetAllTenantConnectorConfigs(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantConnectorInstancesHandler = tenant_provisioning_service.GetAllTenantConnectorInstancesHandlerFunc(handlers.HandleGetAllTenantConnectorInstances(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantDomainsHandler = tenant_provisioning_service.GetAllTenantDomainsHandlerFunc(handlers.HandleGetAllTenantDomains(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantMonitoredObjectsHandler = tenant_provisioning_service.GetAllTenantMonitoredObjectsHandlerFunc(handlers.HandleGetAllTenantMonitoredObjects(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantThresholdProfilesHandler = tenant_provisioning_service.GetAllTenantThresholdProfilesHandlerFunc(handlers.HandleGetAllTenantThresholdProfiles(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetAllTenantUsersHandler = tenant_provisioning_service.GetAllTenantUsersHandlerFunc(func(params tenant_provisioning_service.GetAllTenantUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetAllTenantUsers has not yet been implemented")
	})
	api.AdminProvisioningServiceGetAllTenantsHandler = admin_provisioning_service.GetAllTenantsHandlerFunc(handlers.HandleGetAllTenants(handlers.SkylightAdminRoleOnly, adminDB))
	api.TenantProvisioningServiceGetDomainToMonitoredObjectMapHandler = tenant_provisioning_service.GetDomainToMonitoredObjectMapHandlerFunc(func(params tenant_provisioning_service.GetDomainToMonitoredObjectMapParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetDomainToMonitoredObjectMap has not yet been implemented")
	})
	api.MetricsServiceGetHistogramHandler = metrics_service.GetHistogramHandlerFunc(func(params metrics_service.GetHistogramParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.GetHistogram has not yet been implemented")
	})
	api.AdminProvisioningServiceGetIngestionDictionaryHandler = admin_provisioning_service.GetIngestionDictionaryHandlerFunc(handlers.HandleGetIngestionDictionary(handlers.SkylightAdminRoleOnly, adminDB))
	api.MetricsServiceGetRawMetricsHandler = metrics_service.GetRawMetricsHandlerFunc(func(params metrics_service.GetRawMetricsParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.GetRawMetrics has not yet been implemented")
	})
	api.TenantProvisioningServiceGetReportScheduleConfigHandler = tenant_provisioning_service.GetReportScheduleConfigHandlerFunc(func(params tenant_provisioning_service.GetReportScheduleConfigParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetReportScheduleConfig has not yet been implemented")
	})
	api.TenantProvisioningServiceGetSLAReportHandler = tenant_provisioning_service.GetSLAReportHandlerFunc(func(params tenant_provisioning_service.GetSLAReportParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetSLAReport has not yet been implemented")
	})
	api.AdminProvisioningServiceGetTenantHandler = admin_provisioning_service.GetTenantHandlerFunc(handlers.HandleGetTenant(handlers.SkylightAdminRoleOnly, adminDB))

	api.TenantProvisioningServiceGetTenantConnectorConfigHandler = tenant_provisioning_service.GetTenantConnectorConfigHandlerFunc(handlers.HandleGetTenantConnectorConfig(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantConnectorInstanceHandler = tenant_provisioning_service.GetTenantConnectorInstanceHandlerFunc(handlers.HandleGetTenantConnectorInstance(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantDomainHandler = tenant_provisioning_service.GetTenantDomainHandlerFunc(handlers.HandleGetTenantDomain(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantIngestionProfileHandler = tenant_provisioning_service.GetTenantIngestionProfileHandlerFunc(handlers.HandleGetTenantIngestionProfile(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantMetadataHandler = tenant_provisioning_service.GetTenantMetadataHandlerFunc(func(params tenant_provisioning_service.GetTenantMetadataParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetTenantMetadata has not yet been implemented")
	})
	api.TenantProvisioningServiceGetTenantMonitoredObjectHandler = tenant_provisioning_service.GetTenantMonitoredObjectHandlerFunc(handlers.HandleGetTenantMonitoredObject(handlers.AllRoles, tenantDB))
	api.AdminProvisioningServiceGetTenantSummaryByAliasHandler = admin_provisioning_service.GetTenantSummaryByAliasHandlerFunc(handlers.HandleGetTenantSummaryByAlias(adminDB))
	api.TenantProvisioningServiceGetTenantThresholdProfileHandler = tenant_provisioning_service.GetTenantThresholdProfileHandlerFunc(handlers.HandleGetTenantThresholdProfile(handlers.AllRoles, tenantDB))
	api.TenantProvisioningServiceGetTenantUserHandler = tenant_provisioning_service.GetTenantUserHandlerFunc(func(params tenant_provisioning_service.GetTenantUserParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.GetTenantUser has not yet been implemented")
	})
	api.MetricsServiceGetThresholdCrossingHandler = metrics_service.GetThresholdCrossingHandlerFunc(func(params metrics_service.GetThresholdCrossingParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.GetThresholdCrossing has not yet been implemented")
	})
	api.MetricsServiceGetThresholdCrossingByMonitoredObjectHandler = metrics_service.GetThresholdCrossingByMonitoredObjectHandlerFunc(func(params metrics_service.GetThresholdCrossingByMonitoredObjectParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.GetThresholdCrossingByMonitoredObject has not yet been implemented")
	})
	api.MetricsServiceGetThresholdCrossingByMonitoredObjectTopNHandler = metrics_service.GetThresholdCrossingByMonitoredObjectTopNHandlerFunc(func(params metrics_service.GetThresholdCrossingByMonitoredObjectTopNParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.GetThresholdCrossingByMonitoredObjectTopN has not yet been implemented")
	})
	api.MetricsServiceGetTopNForMetricHandler = metrics_service.GetTopNForMetricHandlerFunc(func(params metrics_service.GetTopNForMetricParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.GetTopNForMetric has not yet been implemented")
	})
	api.AdminProvisioningServiceGetValidTypesHandler = admin_provisioning_service.GetValidTypesHandlerFunc(handlers.HandleGetValidTypes(handlers.SkylightAdminRoleOnly, adminDB))
	api.AdminProvisioningServiceGetTenantIDByAliasHandler = admin_provisioning_service.GetTenantIDByAliasHandlerFunc(handlers.HandleGetTenantIDByAlias(adminDB))
	api.TenantProvisioningServicePatchTenantMetadataHandler = tenant_provisioning_service.PatchTenantMetadataHandlerFunc(func(params tenant_provisioning_service.PatchTenantMetadataParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.PatchTenantMetadata has not yet been implemented")
	})
	api.AdminProvisioningServicePatchTenantHandler = admin_provisioning_service.PatchTenantHandlerFunc(handlers.HandlePatchTenant(handlers.SkylightAdminRoleOnly, adminDB))
	api.TenantProvisioningServicePatchTenantDomainHandler = tenant_provisioning_service.PatchTenantDomainHandlerFunc(handlers.HandlePatchTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServicePatchTenantIngestionProfileHandler = tenant_provisioning_service.PatchTenantIngestionProfileHandlerFunc(handlers.HandlePatchTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServicePatchTenantMonitoredObjectHandler = tenant_provisioning_service.PatchTenantMonitoredObjectHandlerFunc(handlers.HandlePatchTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServicePatchTenantThresholdProfileHandler = tenant_provisioning_service.PatchTenantThresholdProfileHandlerFunc(handlers.HandlePatchTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))

	api.MetricsServiceQueryAggregatedMetricsHandler = metrics_service.QueryAggregatedMetricsHandlerFunc(func(params metrics_service.QueryAggregatedMetricsParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.QueryAggregatedMetrics has not yet been implemented")
	})
	api.MetricsServiceQueryThresholdCrossingHandler = metrics_service.QueryThresholdCrossingHandlerFunc(func(params metrics_service.QueryThresholdCrossingParams) middleware.Responder {
		return middleware.NotImplemented("operation metrics_service.QueryThresholdCrossing has not yet been implemented")
	})

	api.AdminProvisioningServiceUpdateIngestionDictionaryHandler = admin_provisioning_service.UpdateIngestionDictionaryHandlerFunc(handlers.HandleUpdateIngestionDictionary(handlers.SkylightAdminRoleOnly, adminDB))
	api.TenantProvisioningServiceUpdateReportScheduleConfigHandler = tenant_provisioning_service.UpdateReportScheduleConfigHandlerFunc(func(params tenant_provisioning_service.UpdateReportScheduleConfigParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.UpdateReportScheduleConfig has not yet been implemented")
	})

	api.TenantProvisioningServiceUpdateTenantConnectorConfigHandler = tenant_provisioning_service.UpdateTenantConnectorConfigHandlerFunc(handlers.HandleUpdateTenantConnectorConfig(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantConnectorInstanceHandler = tenant_provisioning_service.UpdateTenantConnectorInstanceHandlerFunc(handlers.HandleUpdateTenantConnectorInstance(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.AdminProvisioningServiceUpdateValidTypesHandler = admin_provisioning_service.UpdateValidTypesHandlerFunc(handlers.HandleUpdateValidTypes(handlers.SkylightAdminRoleOnly, adminDB))

	// TODO: calls that will be removed, but just moving them here for now until it is certain we will not use them
	// ======================= START OF CALLS TO REMOVE ===========================================================
	api.AdminProvisioningServiceCreateAdminUserHandler = admin_provisioning_service.CreateAdminUserHandlerFunc(func(params admin_provisioning_service.CreateAdminUserParams) middleware.Responder {
		return middleware.NotImplemented("operation admin_provisioning_service.CreateAdminUser has not yet been implemented")
	})
	api.AdminProvisioningServiceDeleteAdminUserHandler = admin_provisioning_service.DeleteAdminUserHandlerFunc(func(params admin_provisioning_service.DeleteAdminUserParams) middleware.Responder {
		return middleware.NotImplemented("operation admin_provisioning_service.DeleteAdminUser has not yet been implemented")
	})
	api.AdminProvisioningServiceGetAdminUserHandler = admin_provisioning_service.GetAdminUserHandlerFunc(func(params admin_provisioning_service.GetAdminUserParams) middleware.Responder {
		return middleware.NotImplemented("operation admin_provisioning_service.GetAdminUser has not yet been implemented")
	})
	api.AdminProvisioningServiceGetAllAdminUsersHandler = admin_provisioning_service.GetAllAdminUsersHandlerFunc(func(params admin_provisioning_service.GetAllAdminUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation admin_provisioning_service.GetAllAdminUsers has not yet been implemented")
	})
	api.TenantProvisioningServicePatchTenantUserHandler = tenant_provisioning_service.PatchTenantUserHandlerFunc(func(params tenant_provisioning_service.PatchTenantUserParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.PatchTenantUser has not yet been implemented")
	})
	api.AdminProvisioningServiceUpdateAdminUserHandler = admin_provisioning_service.UpdateAdminUserHandlerFunc(func(params admin_provisioning_service.UpdateAdminUserParams) middleware.Responder {
		return middleware.NotImplemented("operation admin_provisioning_service.UpdateAdminUser has not yet been implemented")
	})
	api.AdminProvisioningServiceUpdateTenantHandler = admin_provisioning_service.UpdateTenantHandlerFunc(func(params admin_provisioning_service.UpdateTenantParams) middleware.Responder {
		return middleware.NotImplemented("operation admin_provisioning_service.UpdateTenant has not yet been implemented")
	})
	api.TenantProvisioningServiceUpdateTenantDomainHandler = tenant_provisioning_service.UpdateTenantDomainHandlerFunc(handlers.HandleUpdateTenantDomain(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantIngestionProfileHandler = tenant_provisioning_service.UpdateTenantIngestionProfileHandlerFunc(handlers.HandleUpdateTenantIngestionProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantMetadataHandler = tenant_provisioning_service.UpdateTenantMetadataHandlerFunc(func(params tenant_provisioning_service.UpdateTenantMetadataParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.UpdateTenantMetadata has not yet been implemented")
	})
	api.TenantProvisioningServiceUpdateTenantMonitoredObjectHandler = tenant_provisioning_service.UpdateTenantMonitoredObjectHandlerFunc(handlers.HandleUpdateTenantMonitoredObject(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantThresholdProfileHandler = tenant_provisioning_service.UpdateTenantThresholdProfileHandlerFunc(handlers.HandleUpdateTenantThresholdProfile(handlers.SkylightAndTenantAdminRoles, tenantDB))
	api.TenantProvisioningServiceUpdateTenantUserHandler = tenant_provisioning_service.UpdateTenantUserHandlerFunc(func(params tenant_provisioning_service.UpdateTenantUserParams) middleware.Responder {
		return middleware.NotImplemented("operation tenant_provisioning_service.UpdateTenantUser has not yet been implemented")
	})
	// ============================================ END OF CALLS TO BE REMOVED ==================================================================

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
