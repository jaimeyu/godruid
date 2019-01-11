package handlers_test

import (
	"math/rand"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/models/common"
	"github.com/accedian/adh-gather/models/metrics"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
)

var (
	ReportUrl = "http://deployment.test.cool/api/v2/reports"
)

func TestReportInstanceFetchV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing reports
	existing := handlers.HandleGetAllSLAReportsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllSLAReportsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllSLAReportsV2NotFound)
	assert.NotNil(t, castedResponse)

	singleFetchNotFound := handlers.HandleGetSLAReportV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetSLAReportV2Params{ReportID: "notFound", HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportUrl, "GET")})
	singleFetchResponseNotFound := singleFetchNotFound.(*tenant_provisioning_service_v2.GetSLAReportV2NotFound)
	assert.NotNil(t, singleFetchResponseNotFound)

	// Create some reports
	report1 := generateRandomReport(*castedCreateTeant.Payload.Data.ID)
	report2 := generateRandomReport(*castedCreateTeant.Payload.Data.ID)
	report3 := generateRandomReport(*castedCreateTeant.Payload.Data.ID)

	_, err := tenantDB.CreateSLAReport(report1)
	assert.Nil(t, err)
	_, err = tenantDB.CreateSLAReport(report2)
	assert.Nil(t, err)
	createdReport3, err := tenantDB.CreateSLAReport(report3)
	assert.Nil(t, err)

	// Should now have 3 reports
	existingFetch := handlers.HandleGetAllSLAReportsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllSLAReportsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportUrl, "GET")})
	castedResponseFetch := existingFetch.(*tenant_provisioning_service_v2.GetAllSLAReportsV2OK)
	assert.NotNil(t, castedResponseFetch)
	assert.Equal(t, 3, len(castedResponseFetch.Payload.Data))

	// Fetch by id
	singleFetch := handlers.HandleGetSLAReportV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetSLAReportV2Params{ReportID: createdReport3.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportUrl, "GET")})
	singleFetchResponse := singleFetch.(*tenant_provisioning_service_v2.GetSLAReportV2OK)
	assert.NotNil(t, singleFetchResponse)
	assert.Equal(t, createdReport3.ID, singleFetchResponse.Payload.Data[0].ID)
}

func TestSLAReportAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllSLAReportsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllSLAReportsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ReportUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllSLAReportsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetch := handlers.HandleGetSLAReportV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetSLAReportV2Params{ReportID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ReportUrl, "GET")})
	castedFetch := fetch.(*tenant_provisioning_service_v2.GetSLAReportV2Forbidden)
	assert.NotNil(t, castedFetch)
}

func generateRandomReport(tenantID string) *metrics.SLAReport {

	return &metrics.SLAReport{
		ReportCompletionTime: fake.CharactersN(7),
		ReportScheduleConfig: fake.CharactersN(16),
		ReportTimeRange:      fake.CharactersN(12),
		TenantID:             tenantID,
		ReportSummary: metrics.ReportSummary{
			ObjectCount:            int32(rand.Intn(500)),
			SLACompliancePercent:   rand.Float32(),
			TotalDuration:          rand.Int63(),
			TotalViolationCount:    int32(rand.Intn(500)),
			TotalViolationDuration: rand.Int63(),
		},
	}
}
