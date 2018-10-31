package handlers_test

import (
	"math/rand"
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
	MetricBaselineUrl                  = "http://deployment.test.cool/api/v2/metric-baselines"
	MetricBaselineByMonitoredObjectUrl = "http://deployment.test.cool/api/v2/metric-baselines/by-monitored-object"

	MetricBaselineTypeString = "metricBaselines"
)

func TestMetricBaselineCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	created := handlers.HandleCreateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetricBaselineV2Params{Body: generateRandomTenantMetricBaselineCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMetricBaselineV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.MonitoredObjectID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Baselines)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetricBaselineV2Params{MetricBaselineID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMetricBaselineV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Make sure we can retrieve this record by monitorerd object:
	fetchByMO := handlers.HandleGetMetricBaselineByMonitoredObjectIDV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDV2Params{MonitoredObjectID: castedCreate.Payload.Data.Attributes.MonitoredObjectID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "GET")})
	castedFetchByMO := fetchByMO.(*tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDV2OK)
	assert.NotNil(t, castedFetchByMO)
	assert.Equal(t, *castedCreate.Payload.Data.ID, *castedFetchByMO.Payload.Data.ID)

	// Get metric baselines for an hour of the week for a MO
	fetchByMOForHour := handlers.HandleGetMetricBaselineByMonitoredObjectIdForHourOfWeekV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDForHourOfWeekV2Params{MonitoredObjectID: castedCreate.Payload.Data.Attributes.MonitoredObjectID, HourOfWeek: *castedCreate.Payload.Data.Attributes.Baselines[0].HourOfWeek,
		HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "GET")})
	castedFetchByMOForHour := fetchByMOForHour.(*tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDForHourOfWeekV2OK)
	assert.NotNil(t, castedFetchByMOForHour)
	assert.NotEmpty(t, castedFetchByMOForHour.Payload.Data.Attributes.Baselines)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.MonitoredObjectID, castedFetchByMOForHour.Payload.Data.Attributes.MonitoredObjectID)

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateMetricBaselineUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, castedCreate.Payload.Data.Attributes.Baselines)
	updated := handlers.HandleUpdateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetricBaselineV2Params{MetricBaselineID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetricBaselineV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, castedUpdate.Payload.Data.Attributes.MonitoredObjectID)
	assert.ElementsMatch(t, castedCreate.Payload.Data.Attributes.Baselines, castedUpdate.Payload.Data.Attributes.Baselines)

	// Update baseline for hour of week
	newBaseline := generateRandomMetricBaselineData()
	requestBody := swagmodels.MetricBaselineUpdateHourRequest{
		Data: &swagmodels.MetricBaselineUpdateHourRequestData{
			Attributes: newBaseline,
			Type:       &MetricBaselineTypeString,
		},
	}
	updatedBaseline := handlers.HandleUpdateMetricBaselineForHourOfWeekV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetricBaselineForHourOfWeekV2Params{MonitoredObjectID: castedUpdate.Payload.Data.Attributes.MonitoredObjectID, Body: &requestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "PATCH")})
	castedUpdateBaseline := updatedBaseline.(*tenant_provisioning_service_v2.UpdateMetricBaselineForHourOfWeekV2OK)
	assert.NotNil(t, castedUpdateBaseline)
	assert.Equal(t, 3, len(castedUpdateBaseline.Payload.Data.Attributes.Baselines))
	assert.Equal(t, castedUpdate.Payload.Data.Attributes.MonitoredObjectID, castedUpdateBaseline.Payload.Data.Attributes.MonitoredObjectID)

	// Delete the record
	deleted := handlers.HandleDeleteMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetricBaselineV2Params{MetricBaselineID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetricBaselineV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdateBaseline.Payload.Data, castedDelete.Payload.Data)

}

func TestMetricBaselineNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get MetricBaseline
	fetched := handlers.HandleGetMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetricBaselineV2Params{MetricBaselineID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetricBaselineUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMetricBaselineV2NotFound)
	assert.NotNil(t, castedFetch)

	// Get by MO
	fetchedByMO := handlers.HandleGetMetricBaselineByMonitoredObjectIDV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDV2Params{MonitoredObjectID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetricBaselineUrl, "GET")})
	castedFetchByMO := fetchedByMO.(*tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDV2NotFound)
	assert.NotNil(t, castedFetchByMO)

	// By MO for Hour
	fetchByMOForHour := handlers.HandleGetMetricBaselineByMonitoredObjectIdForHourOfWeekV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDForHourOfWeekV2Params{MonitoredObjectID: notFoundID, HourOfWeek: int32(4),
		HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetricBaselineUrl, "GET")})
	castedFetchByMOForHour := fetchByMOForHour.(*tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDForHourOfWeekV2NotFound)
	assert.NotNil(t, castedFetchByMOForHour)
	// Delete MetricBaseline
	deleted := handlers.HandleDeleteMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetricBaselineV2Params{MetricBaselineID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetricBaselineUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetricBaselineV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch MetricBaseline
	updateRequest := generateMetricBaselineUpdateRequest(notFoundID, "reviosionstuff", nil, nil)
	updated := handlers.HandleUpdateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetricBaselineV2Params{MetricBaselineID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetricBaselineUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetricBaselineV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestMetricBaselineBadRequestV2(t *testing.T) {

	// CreateMetricBaseline
	created := handlers.HandleCreateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetricBaselineV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetricBaselineUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMetricBaselineV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update MetricBaseline
	updated := handlers.HandleUpdateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetricBaselineV2Params{MetricBaselineID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetricBaselineUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetricBaselineV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestMetricBaselineConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	createReqBody := generateRandomTenantMetricBaselineCreationRequest()
	created := handlers.HandleCreateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetricBaselineV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMetricBaselineV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.MonitoredObjectID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Baselines)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should fail as only 1 baseline per MO
	createdConflict := handlers.HandleCreateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetricBaselineV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateMetricBaselineV2Conflict)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateMetricBaselineUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil)
	updated := handlers.HandleUpdateMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetricBaselineV2Params{MetricBaselineID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetricBaselineV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetricBaselineV2Params{MetricBaselineID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetricBaselineUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetricBaselineV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestMetricBaselineAPIsProtectedByAuthV2(t *testing.T) {

	fakeID := fake.CharactersN(20)

	fakeTenantID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetMetricBaselineV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetricBaselineV2Params{MetricBaselineID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, MetricBaselineUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMetricBaselineV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateMetricBaselineV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetricBaselineV2Params{Body: generateRandomTenantMetricBaselineCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, MetricBaselineUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMetricBaselineV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateMetricBaselineV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetricBaselineV2Params{Body: generateRandomTenantMetricBaselineCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, MetricBaselineUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateMetricBaselineV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateMetricBaselineV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetricBaselineV2Params{MetricBaselineID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, MetricBaselineUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetricBaselineV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateMetricBaselineV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetricBaselineV2Params{MetricBaselineID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, MetricBaselineUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateMetricBaselineV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteMetricBaselineV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetricBaselineV2Params{MetricBaselineID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, MetricBaselineUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetricBaselineV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteMetricBaselineV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetricBaselineV2Params{MetricBaselineID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, MetricBaselineUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteMetricBaselineV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantMetricBaselineCreationRequest() *swagmodels.MetricBaselineCreateRequest {
	monObjID := fake.CharactersN(6)
	baselines := generateRandomMetricBaselineDataArray(2)
	return &swagmodels.MetricBaselineCreateRequest{
		Data: &swagmodels.MetricBaselineCreateRequestData{
			Type: &MetricBaselineTypeString,
			Attributes: &swagmodels.MetricBaselineCreateRequestDataAttributes{
				MonitoredObjectID: monObjID,
				Baselines:         baselines,
			},
		},
	}
}

func generateMetricBaselineUpdateRequest(id string, rev string, moID *string, baselines []*swagmodels.MetricBaselineData) *swagmodels.MetricBaselineUpdateRequest {
	result := &swagmodels.MetricBaselineUpdateRequest{
		Data: &swagmodels.MetricBaselineUpdateRequestData{
			Type:       &MetricBaselineTypeString,
			ID:         &id,
			Attributes: &swagmodels.MetricBaselineUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if moID != nil {
		result.Data.Attributes.MonitoredObjectID = *moID
	}

	if baselines != nil {
		result.Data.Attributes.Baselines = baselines
	}

	return result
}

func generateRandomMetricBaselineData() *swagmodels.MetricBaselineData {
	hourOfWeek := int32(rand.Intn(3))
	return &swagmodels.MetricBaselineData{Metric: fake.CharactersN(6), Direction: fake.CharactersN(1), HourOfWeek: &hourOfWeek}
}

func generateRandomMetricBaselineDataArray(count int) []*swagmodels.MetricBaselineData {
	result := []*swagmodels.MetricBaselineData{}

	for i := 0; i < count; i++ {
		result = append(result, generateRandomMetricBaselineData())
	}

	return result
}
