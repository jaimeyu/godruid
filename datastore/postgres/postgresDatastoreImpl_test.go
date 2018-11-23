package postgres

import (
	"log"
	"testing"

	"github.com/icrowley/fake"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/adh-gather/monitoring"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	adminDBName = "adh-admin"
)

var (
	baselineDB *TenantMetricBaselinePostgresDAO
)

func setupPostgresDB() *TenantMetricBaselinePostgresDAO {
	// Configure the test AdminService DAO to use the newly started couch docker image
	cfg := gather.LoadConfig("../../config/adh-gather-test.yml", viper.New())
	cfg.Set("ingDict", "../../files/defaultIngestionDictionary.json")
	cfg.Set(gather.CK_args_metricbaselines_schemadir.String(), "schema")

	monitoring.InitMetrics()

	// Before the tests run, setup the adh-admin db
	var err error
	baselineDB, err = CreateTenantMetricBaselinePostgresDAO()
	if err != nil {
		log.Fatalf("Unabke to create metric baseline db: %s", err.Error())
	}

	return baselineDB
}

func TestCouchDBImplMain(t *testing.T) {
	mbdb := setupPostgresDB()
	defer clearPostgres(mbdb)

	runTenantMetricBaselineCRUD(t)
}

func clearPostgres(dbImpl *TenantMetricBaselinePostgresDAO) {
	_, err := dbImpl.DB.Exec("DELETE FROM metric_baselines")
	if err != nil {
		log.Fatalf("Could not delete DB data after test: %s", err.Error())
	}
}

func runTenantMetricBaselineCRUD(t *testing.T) {
	const COMPANY1 = "MetricBaselineCompany"
	const SUBDOMAIN1 = "subdom1"
	const MONOBJ1 = "MONOBJ1"
	const MONOBJ2 = "MONOBJ2"
	const DELAYP95 = "delayP95"
	const JITTER95 = "jitterP95"

	TENANT := fake.Brand()

	// Try to fetch a record even though none exist:
	fail, err := baselineDB.GetMetricBaseline(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	fail, err = baselineDB.GetMetricBaseline(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	failArray, err := baselineDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, "someID", int32(1))
	assert.NotNil(t, err)
	assert.Nil(t, failArray)

	// Try to Update a record that does not exist:
	fail, err = baselineDB.UpdateMetricBaseline(&tenmod.MetricBaseline{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to update an hour in a record that does not exist, should creatre the record
	upsert, err := baselineDB.UpdateMetricBaselineForHourOfWeek(TENANT, MONOBJ1, 150, &tenmod.MetricBaselineData{Metric: DELAYP95, Direction: "0"})
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
	failedCreate, err := baselineDB.CreateMetricBaseline(&createObj)
	assert.NotNil(t, err)
	assert.Nil(t, failedCreate)

	// Create a record successfully
	metricBaseline1 := tenmod.MetricBaseline{
		Datatype:          string(tenmod.TenantMetricBaselineType),
		TenantID:          TENANT,
		MonitoredObjectID: MONOBJ2,
		Baselines:         []*tenmod.MetricBaselineData{},
	}
	created, err := baselineDB.CreateMetricBaseline(&metricBaseline1)
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
	fetched, err := baselineDB.GetMetricBaseline(TENANT, created.ID)
	assert.Nil(t, err)
	assert.ElementsMatch(t, created.Baselines, fetched.Baselines, "The retrieved record should have the same baselines as the created record")
	assert.Equal(t, created.MonitoredObjectID, fetched.MonitoredObjectID, "Monitored object not the same")

	logger.Log.Warningf("BEFORE UPDATE: %s", models.AsJSONString(upsert))

	// Add new baseline data to existing record
	updated, err := baselineDB.UpdateMetricBaselineForHourOfWeek(TENANT, upsert.MonitoredObjectID, 150, &tenmod.MetricBaselineData{Metric: JITTER95, Direction: "1"})
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, upsert.ID, updated.ID, "ID not the same")
	assert.Equal(t, TENANT, updated.TenantID, "Tenant ID not the same")
	assert.Equal(t, 2, len(updated.Baselines), "Baseline array should have 2 elements")

	logger.Log.Warningf("AFTER UPDATE ONE: %s", models.AsJSONString(updated))

	// add another bit of data for hour 150
	anotherUpdate, err := baselineDB.UpdateMetricBaselineForHourOfWeek(TENANT, upsert.MonitoredObjectID, 150, &tenmod.MetricBaselineData{Metric: DELAYP95, Direction: "1"})
	assert.Nil(t, err)
	logger.Log.Warningf("AFTER UPDATE TWO: %s", models.AsJSONString(anotherUpdate))

	// Get baselines for an hour of the week - success
	baselineArray, err := baselineDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, upsert.MonitoredObjectID, 150)
	assert.Nil(t, err)
	assert.NotNil(t, baselineArray)
	assert.Equal(t, 3, len(baselineArray))

	baselineArray, err = baselineDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, upsert.MonitoredObjectID, 250)
	assert.NotNil(t, err)
	assert.Nil(t, baselineArray)

	// Update an entire record
	fetched, err = baselineDB.GetMetricBaseline(TENANT, upsert.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, 3, len(fetched.Baselines))

	fetched.Baselines = []*tenmod.MetricBaselineData{}
	fetched.ID = datastore.GetDataIDFromFullID(fetched.ID)
	replaced, err := baselineDB.UpdateMetricBaseline(fetched)
	assert.Nil(t, err)
	assert.NotNil(t, replaced)
	assert.Equal(t, 0, len(replaced.Baselines))

	fetched, err = baselineDB.GetMetricBaseline(TENANT, MONOBJ1)
	assert.NotNil(t, err)
	assert.Nil(t, fetched)

	// Now delete a metric baseline successfully by ID
	deleted, err := baselineDB.DeleteMetricBaseline(TENANT, datastore.GetDataIDFromFullID(created.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)

	fetched, err = baselineDB.GetMetricBaseline(TENANT, created.MonitoredObjectID)
	assert.NotNil(t, err)
	assert.Nil(t, fetched)
}
