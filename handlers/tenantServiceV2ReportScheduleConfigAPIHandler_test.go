package handlers_test

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/models/common"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
)

var (
	ReportScheduleConfigUrl = "http://deployment.test.cool/api/v2/report-schedule-configs"

	reportScehduleConfigTypeString = "reportScheduleConfigs"
)

func TestReportScheduleConfigCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing ReportScheduleConfigs
	existing := handlers.HandleGetAllReportScheduleConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params{Body: generateRandomTenantReportScheduleConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateReportScheduleConfigV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Relationships)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Hour)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Minute)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.DayMonth)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.DayWeek)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetReportScheduleConfigV2Params{ConfigID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetReportScheduleConfigV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllReportScheduleConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateReportScheduleConfigUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil)
	updated := handlers.HandleUpdateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params{ConfigID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateReportScheduleConfigV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Name)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.Hour, castedUpdate.Payload.Data.Attributes.Hour)
	assert.Equal(t, castedCreate.Payload.Data.Relationships.ThresholdProfile, castedUpdate.Payload.Data.Relationships.ThresholdProfile)

	// Delete the record
	deleted := handlers.HandleDeleteReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params{ConfigID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteReportScheduleConfigV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing ReportScheduleConfigs
	existing = handlers.HandleGetAllReportScheduleConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestReportScheduleConfigNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get ReportScheduleConfig
	fetched := handlers.HandleGetReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetReportScheduleConfigV2Params{ConfigID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ReportScheduleConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetReportScheduleConfigV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete ReportScheduleConfig
	deleted := handlers.HandleDeleteReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params{ConfigID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ReportScheduleConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteReportScheduleConfigV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch ReportScheduleConfig
	updateRequest := generateReportScheduleConfigUpdateRequest(notFoundID, "reviosionstuff", nil, nil)
	updated := handlers.HandleUpdateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params{ConfigID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ReportScheduleConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateReportScheduleConfigV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestReportScheduleConfigBadRequestV2(t *testing.T) {

	// CreateReportScheduleConfig
	created := handlers.HandleCreateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ReportScheduleConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateReportScheduleConfigV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update ReportScheduleConfig
	updated := handlers.HandleUpdateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params{ConfigID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ReportScheduleConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateReportScheduleConfigV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestReportScheduleConfigConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing ReportScheduleConfig
	existing := handlers.HandleGetAllReportScheduleConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantReportScheduleConfigCreationRequest()
	created := handlers.HandleCreateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateReportScheduleConfigV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Relationships.ThresholdProfile)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Hour)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Minute)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.DayMonth)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.DayWeek)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateReportScheduleConfigV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateReportScheduleConfigUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil)
	updated := handlers.HandleUpdateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params{ConfigID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params{ConfigID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteReportScheduleConfigV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestReportScheduleConfigAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllReportScheduleConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ReportScheduleConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetReportScheduleConfigV2Params{ConfigID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ReportScheduleConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetReportScheduleConfigV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params{Body: generateRandomTenantReportScheduleConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ReportScheduleConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateReportScheduleConfigV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params{Body: generateRandomTenantReportScheduleConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, ReportScheduleConfigUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateReportScheduleConfigV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params{ConfigID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ReportScheduleConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params{ConfigID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, ReportScheduleConfigUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params{ConfigID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ReportScheduleConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteReportScheduleConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params{ConfigID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, ReportScheduleConfigUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantReportScheduleConfigCreationRequest() *swagmodels.ReportScheduleConfigCreateRequest {
	name := fake.CharactersN(12)
	thresh := fake.CharactersN(12)
	dayMonth := "1"
	hour := "12"
	minute := "14"
	dayWeek := "5"
	month := "4"

	return &swagmodels.ReportScheduleConfigCreateRequest{
		Data: &swagmodels.ReportScheduleConfigCreateRequestData{
			Type: &reportScehduleConfigTypeString,
			Attributes: &swagmodels.ReportScheduleConfigCreateRequestDataAttributes{
				Name:     &name,
				DayMonth: dayMonth,
				DayWeek:  dayWeek,
				Hour:     hour,
				Minute:   minute,
				Month:    month,
			},
			Relationships: &swagmodels.ReportScheduleConfigRelationships{
				ThresholdProfile: &swagmodels.JSONAPISingleRelationship{
					Data: &swagmodels.JSONAPIRelationshipData{
						ID:   thresh,
						Type: "thresholdProfiles",
					},
				},
			},
		},
	}
}

func generateReportScheduleConfigUpdateRequest(id string, rev string, name *string, thresh *string) *swagmodels.ReportScheduleConfigUpdateRequest {
	result := &swagmodels.ReportScheduleConfigUpdateRequest{
		Data: &swagmodels.ReportScheduleConfigUpdateRequestData{
			Type:       &reportScehduleConfigTypeString,
			ID:         &id,
			Attributes: &swagmodels.ReportScheduleConfigUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if name != nil {
		result.Data.Attributes.Name = *name
	}

	if thresh != nil {
		result.Data.Relationships.ThresholdProfile.Data.ID = *thresh
	}

	return result
}
