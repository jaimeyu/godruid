package tenant

import (
	"testing"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/stretchr/testify/assert"
)

func (runner *TenantServiceDatastoreTestRunner) RunTenantMetricBaselineCRUD(t *testing.T) {
	mbDB := runner.baselineDB
	const COMPANY1 = "MetricBaselineCompany"
	const SUBDOMAIN1 = "subdom1"
	const MONOBJ1 = "MONOBJ1"
	const MONOBJ2 = "MONOBJ2"
	const DELAYP95 = "delayP95"
	const JITTER95 = "jitterP95"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Try to fetch a record even though none exist:
	fail, err := mbDB.GetMetricBaseline(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	fail, err = mbDB.GetMetricBaseline(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	failArray, err := mbDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, "someID", int32(1))
	assert.NotNil(t, err)
	assert.Nil(t, failArray)

	// Try to Update a record that does not exist:
	fail, err = mbDB.UpdateMetricBaseline(&tenmod.MetricBaseline{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to update an hour in a record that does not exist, should creatre the record
	upsert, err := mbDB.UpdateMetricBaselineForHourOfWeek(TENANT, MONOBJ1, 150, &tenmod.MetricBaselineData{Metric: DELAYP95, Direction: "0"})
	assert.Nil(t, err)
	assert.NotNil(t, upsert)
	assert.NotEmpty(t, upsert.ID)
	assert.NotEmpty(t, upsert.REV)
	assert.Equal(t, string(tenmod.TenantMetricBaselineType), upsert.Datatype)
	assert.Equal(t, DELAYP95, upsert.Baselines[0].Metric, "Metric not the same")
	assert.Equal(t, "0", upsert.Baselines[0].Direction, "Direction not the same")
	assert.Equal(t, int32(150), upsert.HourOfWeek, "HourofWeek not the same")
	assert.Equal(t, TENANT, upsert.TenantID, "Tenant ID not the same")
	assert.True(t, upsert.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, upsert.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Try to create the same record, should fail
	createObj := tenmod.MetricBaseline{
		Datatype:          string(tenmod.TenantMetricBaselineType),
		TenantID:          TENANT,
		MonitoredObjectID: MONOBJ1,
		HourOfWeek:        150,
	}
	failedCreate, err := mbDB.CreateMetricBaseline(&createObj)
	assert.NotNil(t, err)
	assert.Nil(t, failedCreate)

	// Create a record successfully
	metricBaseline1 := tenmod.MetricBaseline{
		Datatype:          string(tenmod.TenantMetricBaselineType),
		TenantID:          TENANT,
		MonitoredObjectID: MONOBJ2,
		Baselines:         []*tenmod.MetricBaselineData{},
	}
	created, err := mbDB.CreateMetricBaseline(&metricBaseline1)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantMetricBaselineType), created.Datatype)
	assert.Equal(t, TENANT, upsert.TenantID, "Tenant ID not the same")
	assert.Empty(t, created.Baselines, "Baseline array should be empty")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := mbDB.GetMetricBaseline(TENANT, created.ID)
	assert.Nil(t, err)
	assert.ElementsMatch(t, created.Baselines, fetched.Baselines, "The retrieved record should have the same baselines as the created record")
	assert.Equal(t, created.MonitoredObjectID, fetched.MonitoredObjectID, "Monitored object not the same")

	logger.Log.Warningf("BEFORE UPDATE: %s", models.AsJSONString(upsert))

	// Add new baseline data to existing record
	updated, err := mbDB.UpdateMetricBaselineForHourOfWeek(TENANT, upsert.MonitoredObjectID, 150, &tenmod.MetricBaselineData{Metric: JITTER95, Direction: "1"})
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, upsert.ID, updated.ID, "ID not the same")
	assert.Equal(t, TENANT, updated.TenantID, "Tenant ID not the same")
	assert.Equal(t, 2, len(updated.Baselines), "Baseline array should have 2 elements")

	logger.Log.Warningf("AFTER UPDATE ONE: %s", models.AsJSONString(updated))

	// add another bit of data for hour 150
	anotherUpdate, err := mbDB.UpdateMetricBaselineForHourOfWeek(TENANT, upsert.MonitoredObjectID, 150, &tenmod.MetricBaselineData{Metric: DELAYP95, Direction: "1"})
	assert.Nil(t, err)
	logger.Log.Warningf("AFTER UPDATE TWO: %s", models.AsJSONString(anotherUpdate))

	// Get baselines for an hour of the week - success
	baselineArray, err := mbDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, upsert.MonitoredObjectID, 150)
	assert.Nil(t, err)
	assert.NotNil(t, baselineArray)
	assert.Equal(t, 3, len(baselineArray))

	baselineArray, err = mbDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, upsert.MonitoredObjectID, 250)
	assert.NotNil(t, err)
	assert.Nil(t, baselineArray)

	// Update an entire record
	fetched, err = mbDB.GetMetricBaseline(TENANT, upsert.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, 3, len(fetched.Baselines))

	fetched.Baselines = []*tenmod.MetricBaselineData{}
	fetched.ID = ds.GetDataIDFromFullID(fetched.ID)
	replaced, err := mbDB.UpdateMetricBaseline(fetched)
	assert.Nil(t, err)
	assert.NotNil(t, replaced)
	assert.Equal(t, 0, len(replaced.Baselines))

	fetched, err = mbDB.GetMetricBaseline(TENANT, MONOBJ1)
	assert.NotNil(t, err)
	assert.Nil(t, fetched)

	// Now delete a metric baseline successfully by ID
	deleted, err := mbDB.DeleteMetricBaseline(TENANT, ds.GetDataIDFromFullID(created.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)

	fetched, err = mbDB.GetMetricBaseline(TENANT, created.MonitoredObjectID)
	assert.NotNil(t, err)
	assert.Nil(t, fetched)
}
