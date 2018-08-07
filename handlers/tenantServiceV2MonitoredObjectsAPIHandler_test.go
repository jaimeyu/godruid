package handlers_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/icrowley/fake"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/stretchr/testify/assert"

	tenmod "github.com/accedian/adh-gather/models/tenant"
)

var (
	monitoredObjectUrl = "http://deployment.test.cool/api/v2/monitored-objects"
)

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
