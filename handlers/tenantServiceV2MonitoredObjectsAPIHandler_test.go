package handlers_test

import (
	"encoding/json"
	"math/rand"
	"strings"
	"testing"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/stretchr/testify/assert"

	"github.com/icrowley/fake"

	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

var (
	monitoredObjectUrl     = "http://deployment.test.cool/api/v2/monitored-objects"
	monitoredObjectBulkUrl = "http://deployment.test.cool/api/v2/bulk/insert/monitored-objects"

	monitoredObjectTypeString = "monitoredObjects"
)

func TestMonitoredObjectCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing MonitoredObjects
	existing := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMonitoredObjectV2Params{Body: generateRandomTenantMonitoredObjectCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMonitoredObjectV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ActuatorName)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ActuatorType)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ReflectorName)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ReflectorType)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ObjectName)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ObjectType)
	assert.Equal(t, *castedCreate.Payload.Data.ID, castedCreate.Payload.Data.Attributes.ObjectID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMonitoredObjectV2Params{MonObjID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMonitoredObjectV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newAName := fake.CharactersN(16)
	newOName := fake.CharactersN(16)
	newRName := fake.CharactersN(16)
	updateRequestBody := generateMonitoredObjectUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newOName, &newAName, &newRName)
	updated := handlers.HandleUpdateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params{MonObjID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMonitoredObjectV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newAName, castedUpdate.Payload.Data.Attributes.ActuatorName)
	assert.Equal(t, newOName, castedUpdate.Payload.Data.Attributes.ObjectName)
	assert.Equal(t, newRName, castedUpdate.Payload.Data.Attributes.ReflectorName)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.ReflectorType, castedUpdate.Payload.Data.Attributes.ReflectorType)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.ObjectType, castedUpdate.Payload.Data.Attributes.ObjectType)

	// Delete the record
	deleted := handlers.HandleDeleteMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params{MonObjID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMonitoredObjectV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing MonitoredObjects
	existing = handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestMonitoredObjectNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get MonitoredObject
	fetched := handlers.HandleGetMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMonitoredObjectV2Params{MonObjID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMonitoredObjectV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete MonitoredObject
	deleted := handlers.HandleDeleteMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params{MonObjID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, monitoredObjectUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMonitoredObjectV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch MonitoredObject
	updateRequest := generateMonitoredObjectUpdateRequest(notFoundID, "reviosionstuff", nil, nil, nil)
	updated := handlers.HandleUpdateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params{MonObjID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, monitoredObjectUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMonitoredObjectV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestMonitoredObjectBadRequestV2(t *testing.T) {

	// CreateMonitoredObject
	created := handlers.HandleCreateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMonitoredObjectV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, monitoredObjectUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMonitoredObjectV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update MonitoredObject
	updated := handlers.HandleUpdateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params{MonObjID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, monitoredObjectUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMonitoredObjectV2BadRequest)
	assert.NotNil(t, castedUpdate)

	// Bulk Insert
	insertReq := &swagmodels.BulkMonitoredObjectCreateRequest{
		Data: []*swagmodels.MonitoredObjectCreate{},
	}
	inserted := handlers.HandleBulkCreateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params{Body: insertReq, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantAdmin, monitoredObjectBulkUrl, "POST")})
	castedInsert := inserted.(*tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2BadRequest)
	assert.NotNil(t, castedInsert)

	// Bulk update
	updateReq := &swagmodels.BulkMonitoredObjectUpdateRequest{
		Data: []*swagmodels.MonitoredObjectUpdate{},
	}
	bulkUpdate := handlers.HandleBulkUpdateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Params{Body: updateReq, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantAdmin, monitoredObjectBulkUrl, "POST")})
	castedBulkUpdate := bulkUpdate.(*tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2BadRequest)
	assert.NotNil(t, castedBulkUpdate)
}

func TestMonitoredObjectConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing MonitoredObject
	existing := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantMonitoredObjectCreationRequest()
	created := handlers.HandleCreateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMonitoredObjectV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMonitoredObjectV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ActuatorName)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ActuatorType)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ReflectorName)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ReflectorType)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ObjectName)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ObjectType)
	assert.Equal(t, *castedCreate.Payload.Data.ID, castedCreate.Payload.Data.Attributes.ObjectID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again
	createdConflict := handlers.HandleCreateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMonitoredObjectV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "POST")})
	castedCreateConflict := createdConflict.(*tenant_provisioning_service_v2.CreateMonitoredObjectV2Conflict)
	assert.NotNil(t, castedCreateConflict)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateMonitoredObjectUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil, nil)
	updated := handlers.HandleUpdateMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params{MonObjID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMonitoredObjectV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params{MonObjID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMonitoredObjectV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestMonitoredObjectAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, monitoredObjectUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetMonitoredObjectV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMonitoredObjectV2Params{MonObjID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, monitoredObjectUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMonitoredObjectV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateMonitoredObjectV2Params{Body: generateRandomTenantMonitoredObjectCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, monitoredObjectUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMonitoredObjectV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateMonitoredObjectV2Params{Body: generateRandomTenantMonitoredObjectCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, monitoredObjectUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateMonitoredObjectV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params{MonObjID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, monitoredObjectUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMonitoredObjectV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params{MonObjID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, monitoredObjectUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateMonitoredObjectV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params{MonObjID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, monitoredObjectUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMonitoredObjectV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteMonitoredObjectV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params{MonObjID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, monitoredObjectUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteMonitoredObjectV2Forbidden)
	assert.NotNil(t, castedDelete)

	// Bulk Insert - SkylightAdmin and TenantAdmin Only
	inserted := handlers.HandleBulkCreateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, monitoredObjectBulkUrl, "POST")})
	castedInsert := inserted.(*tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Forbidden)
	assert.NotNil(t, castedInsert)

	inserted = handlers.HandleBulkCreateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, monitoredObjectBulkUrl, "POST")})
	castedInsert = inserted.(*tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Forbidden)
	assert.NotNil(t, castedInsert)

	// Bulk Update - SkylightAdmin and TenantAdmin Only
	bulkUpdate := handlers.HandleBulkUpdateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, monitoredObjectBulkUrl, "PUT")})
	castedBulkUpdate := bulkUpdate.(*tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Forbidden)
	assert.NotNil(t, castedBulkUpdate)

	bulkUpdate = handlers.HandleBulkUpdateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, monitoredObjectBulkUrl, "PUT")})
	castedBulkUpdate = bulkUpdate.(*tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Forbidden)
	assert.NotNil(t, castedBulkUpdate)
	// Bulk Patch
	bulkPatch := handlers.HandleBulkPatchMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkPatchMonitoredObjectsV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, monitoredObjectBulkUrl, "PUT")})
	castedBulkPatch := bulkPatch.(*tenant_provisioning_service_v2.BulkPatchMonitoredObjectsV2Forbidden)
	assert.NotNil(t, castedBulkPatch)

}
func TestGetAllMonitoredObjectsV2(t *testing.T) {

	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no MOs to start:
	existing := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantUser)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2NotFound)
	assert.NotNil(t, castedResponse)

	// Create some MOs
	monitoredObjects := generateSliceOfRandomMOs(createdTenant.ID, 20)
	_, err = tenantDB.BulkInsertMonitoredObjects(createdTenant.ID, monitoredObjects)
	assert.Nil(t, err)

	// Fetch all of the MOs in one batch
	successFetch := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantUser, monitoredObjectUrl, "GET")})
	castedSuccess := successFetch.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedSuccess)
	assert.Equal(t, len(monitoredObjects), len(castedSuccess.Payload.Data), "Should have same number of MOs in DB as in result")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksFirst], "Should have a 'first' link")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksSelf], "Should have a 'self' link")
	assert.Equal(t, castedSuccess.Payload.Links[linksSelf], castedSuccess.Payload.Links[linksFirst], "The 'self' and 'first' links should be the same")
	assert.Empty(t, castedSuccess.Payload.Links[linksPrev], "Should not have a 'prev' link")
	assert.Empty(t, castedSuccess.Payload.Links[linksNext], "Should not have a 'next' link")

	// Store the last offset for later
	lastOffset := castedSuccess.Payload.Data[len(castedSuccess.Payload.Data)-1].Attributes.ObjectName

	// Fetch a known page size from start
	pageSize := int64(5)
	successFetch = handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{Limit: &pageSize, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantUser, monitoredObjectUrl, "GET")})
	castedSuccess = successFetch.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedSuccess)
	assert.Equal(t, pageSize, int64(len(castedSuccess.Payload.Data)), "Should have same number of MOs in DB as the page size")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksFirst], "Should have a 'first' link")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksSelf], "Should have a 'self' link")
	assert.Equal(t, castedSuccess.Payload.Links[linksSelf], castedSuccess.Payload.Links[linksFirst], "The 'self' and 'first' links should be the same")
	assert.Empty(t, castedSuccess.Payload.Links[linksPrev], "Should not have a 'prev' link")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksNext], "Should have a 'next' link")

	// Fetch the next page based on an offset
	pageOffset := getOffsetFromURL(castedSuccess.Payload.Links[linksNext])
	successFetch = handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{StartKey: &pageOffset, Limit: &pageSize, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantUser, monitoredObjectUrl, "GET")})
	castedSuccess = successFetch.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedSuccess)
	assert.Equal(t, pageSize, int64(len(castedSuccess.Payload.Data)), "Should have same number of MOs as the page size")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksFirst], "Should have a 'first' link")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksSelf], "Should have a 'self' link")
	assert.NotEqual(t, castedSuccess.Payload.Links[linksSelf], castedSuccess.Payload.Links[linksFirst], "The 'self' and 'first' links should not be the same")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksPrev], "Should have a 'prev' link")
	assert.NotEmpty(t, castedSuccess.Payload.Links[linksNext], "Should have a 'next' link")

	// lookup something after the last page
	oneAfterLastOffset := lastOffset + "zzzzzzz"
	notFoundFetch := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{StartKey: &oneAfterLastOffset, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantUser, monitoredObjectUrl, "GET")})
	castedNotFound := notFoundFetch.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2NotFound)
	assert.NotNil(t, castedNotFound)

	// Try to fetch with a role that is not allowed
	frbdnFetch := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{StartKey: &oneAfterLastOffset, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleUnknown, monitoredObjectUrl, "GET")})
	castedFrbdn := frbdnFetch.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Forbidden)
	assert.NotNil(t, castedFrbdn)
}

func TestBulkInsertAndUpdateMonitoredObjectsV2(t *testing.T) {
	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no MOs to start:
	existing := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantUser)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2NotFound)
	assert.NotNil(t, castedResponse)

	moCount := 5

	// Insert a some MOs
	createBody := generateRandomBulkInsertMORequest(moCount)
	inserted := handlers.HandleBulkCreateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params{Body: createBody, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantAdmin, monitoredObjectBulkUrl, "POST")})
	castedInsert := inserted.(*tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2OK)
	assert.NotNil(t, castedInsert)
	assert.Equal(t, 5, len(castedInsert.Payload.Data))
	for _, res := range castedInsert.Payload.Data {
		assert.NotEmpty(t, res.Attributes.ID)
		assert.True(t, res.Attributes.Ok)
	}

	// Try to insert the same objects again - should fail each part of the request
	inserted = handlers.HandleBulkCreateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params{Body: createBody, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantAdmin, monitoredObjectBulkUrl, "POST")})
	castedInsert = inserted.(*tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2OK)
	assert.NotNil(t, castedInsert)
	assert.Equal(t, moCount, len(castedInsert.Payload.Data))
	for _, res := range castedInsert.Payload.Data {
		assert.NotEmpty(t, res.Attributes.ID)
		assert.False(t, res.Attributes.Ok)
	}

	// Get all the existing MOs
	fetchList := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, moCount, len(castedFetchList.Payload.Data))

	// Update Some names to be known values
	knownName1 := "REALCHANGE"
	knownName2 := "EQUALLYrealChange"
	bulkUpdateBody := generateBulkUpdateRequest(castedFetchList.Payload.Data, knownName1, knownName2)
	updated := handlers.HandleBulkUpdateMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Params{Body: bulkUpdateBody, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantAdmin, monitoredObjectBulkUrl, "POST")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2OK)
	assert.NotNil(t, castedUpdate)
	assert.Equal(t, moCount, len(castedUpdate.Payload.Data))
	for _, res := range castedUpdate.Payload.Data {
		assert.NotEmpty(t, res.Attributes.ID)
		assert.True(t, res.Attributes.Ok)
	}

	// Retrieve the records again and make sure that the updates were handled properly
	fetchList2 := handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedFetchList2 := fetchList2.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedFetchList2)
	assert.Equal(t, moCount, len(castedFetchList2.Payload.Data))

	knownName1Count := 0
	knownName2Count := 0
	for _, res := range castedFetchList2.Payload.Data {
		if res.Attributes.ActuatorName == knownName1 {
			knownName1Count++
		}
		if res.Attributes.ReflectorName == knownName2 {
			knownName2Count++
		}
	}

	assert.Equal(t, moCount, knownName1Count+knownName2Count)
	// Get all the existing MOs
	fetchList = handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedFetchList = fetchList.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, moCount, len(castedFetchList.Payload.Data))

	// Patch Some names to be known values
	knownName1 = "Patched1"
	knownName2 = "Patchey2"
	bulkPatchBody := generateBulkPatchRequest(castedFetchList.Payload.Data, knownName1, knownName2)
	updated = handlers.HandleBulkPatchMonitoredObjectsV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.BulkPatchMonitoredObjectsV2Params{Body: bulkPatchBody, HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleTenantAdmin, monitoredObjectBulkUrl, "POST")})
	castedPatched := updated.(*tenant_provisioning_service_v2.BulkPatchMonitoredObjectsV2OK)
	assert.NotNil(t, castedPatched)
	assert.Equal(t, moCount, len(castedUpdate.Payload.Data))
	for _, res := range castedUpdate.Payload.Data {
		assert.NotEmpty(t, res.Attributes.ID)
		assert.True(t, res.Attributes.Ok)
	}

	// Retrieve the records again and make sure that the updates were handled properly
	fetchList2 = handlers.HandleGetAllMonitoredObjectsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params{HTTPRequest: createHttpRequestWithParams(createdTenant.ID, handlers.UserRoleSkylight, monitoredObjectUrl, "GET")})
	castedFetchList2 = fetchList2.(*tenant_provisioning_service_v2.GetAllMonitoredObjectsV2OK)
	assert.NotNil(t, castedFetchList2)
	assert.Equal(t, moCount, len(castedFetchList2.Payload.Data))

	knownName1Count = 0
	knownName2Count = 0
	for _, res := range castedFetchList2.Payload.Data {
		if res.Attributes.ActuatorName == knownName1 {
			knownName1Count++
		}
		if res.Attributes.ReflectorName == knownName2 {
			knownName2Count++
		}
	}

	assert.Equal(t, moCount, knownName1Count+knownName2Count)
}

func generateRandomTenantMonitoredObjectCreationRequest() *swagmodels.MonitoredObjectCreateRequest {
	oName := fake.CharactersN(12)
	oType := getRandomObjectType()
	aName := fake.CharactersN(12)
	aType := getRandomDeviceType()
	rName := fake.CharactersN(12)
	rType := getRandomDeviceType()
	moID := strings.Join([]string{oName, aName, rName}, "-")

	return &swagmodels.MonitoredObjectCreateRequest{
		Data: &swagmodels.MonitoredObjectCreate{
			Type: &monitoredObjectTypeString,
			Attributes: &swagmodels.MonitoredObjectCreateAttributes{
				ObjectName:    oName,
				ObjectType:    oType,
				ActuatorName:  aName,
				ActuatorType:  aType,
				ReflectorName: rName,
				ReflectorType: rType,
				ObjectID:      &moID,
			},
		},
	}
}

func generateMonitoredObjectUpdateRequest(id string, rev string, oName *string, aName *string, rName *string) *swagmodels.MonitoredObjectUpdateRequest {
	result := &swagmodels.MonitoredObjectUpdateRequest{
		Data: &swagmodels.MonitoredObjectUpdate{
			Type:       &monitoredObjectTypeString,
			ID:         &id,
			Attributes: &swagmodels.MonitoredObjectUpdateAttributes{Rev: &rev},
		},
	}

	if oName != nil {
		result.Data.Attributes.ObjectName = oName
	}
	if aName != nil {
		result.Data.Attributes.ActuatorName = aName
	}
	if rName != nil {
		result.Data.Attributes.ReflectorName = rName
	}

	return result
}

func getOffsetFromURL(url string) string {
	if len(url) == 0 {
		return ""
	}

	parts := strings.Split(url, "start_key=")
	if len(parts) >= 2 {
		startKey := strings.Split(parts[1], "&")[0]
		return startKey
	}

	return ""
}

func generateSliceOfRandomMOs(tenantID string, numObjects int) []*tenmod.MonitoredObject {
	res := make([]*tenmod.MonitoredObject, 0, numObjects)
	for i := 0; i < numObjects; i++ {
		res = append(res, generateRandomMO(tenantID))
	}

	return res
}

func generateRandomMO(tenantID string) *tenmod.MonitoredObject {
	aName := fake.CharactersN(10)
	rName := fake.CharactersN(10)
	oName := fake.CharactersN(10)
	moID := strings.Join([]string{aName, rName, oName}, "-")
	return &tenmod.MonitoredObject{
		ActuatorName:      aName,
		ActuatorType:      getRandomDeviceType(),
		ObjectName:        oName,
		ObjectType:        getRandomObjectType(),
		ReflectorName:     rName,
		ReflectorType:     getRandomDeviceType(),
		MonitoredObjectID: moID,
		TenantID:          tenantID,
	}
}

func getRandomDeviceType() string {
	index := rand.Intn(len(deviceTypes))
	return deviceTypes[index]
}

func getRandomObjectType() string {
	index := rand.Intn(len(objectTypes))
	return objectTypes[index]
}

func generateRandomBulkInsertMORequest(numObjects int) *swagmodels.BulkMonitoredObjectCreateRequest {

	data := []*swagmodels.MonitoredObjectCreate{}

	for i := 0; i < numObjects; i++ {
		addObject := generateRandomTenantMonitoredObjectCreationRequest()
		data = append(data, addObject.Data)
	}

	return &swagmodels.BulkMonitoredObjectCreateRequest{
		Data: data,
	}
}

func generateBulkPatchRequest(objectsToUpdate []*swagmodels.MonitoredObject, knownName1 string, knownName2 string) *swagmodels.BulkMonitoredObjectPatchRequest {
	data := []*swagmodels.MonitoredObjectPatch{}

	for i, obj := range objectsToUpdate {
		dataBytes, _ := json.Marshal(obj)
		addObject := swagmodels.MonitoredObjectPatch{}
		json.Unmarshal(dataBytes, &addObject)

		if i%2 == 0 {
			addObject.Attributes.ActuatorName = knownName1
		} else {
			addObject.Attributes.ReflectorName = knownName2
		}

		data = append(data, &addObject)
	}

	return &swagmodels.BulkMonitoredObjectPatchRequest{
		Data: data,
	}
}

func generateBulkUpdateRequest(objectsToUpdate []*swagmodels.MonitoredObject, knownName1 string, knownName2 string) *swagmodels.BulkMonitoredObjectUpdateRequest {
	data := []*swagmodels.MonitoredObjectUpdate{}

	for i, obj := range objectsToUpdate {
		dataBytes, _ := json.Marshal(obj)
		addObject := swagmodels.MonitoredObjectUpdate{}
		json.Unmarshal(dataBytes, &addObject)

		if i%2 == 0 {
			addObject.Attributes.ActuatorName = &knownName1
		} else {
			addObject.Attributes.ReflectorName = &knownName2
		}

		data = append(data, &addObject)
	}

	return &swagmodels.BulkMonitoredObjectUpdateRequest{
		Data: data,
	}
}
