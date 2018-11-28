package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/swagmodels"
	_ "github.com/lib/pq"
	"github.com/prometheus/common/log"
)

const (
	metricBaselineType                 = string(tenmod.TenantMetricBaselineType)
	duplicateKey                       = "duplicate key"
	metricBaselineTableNameTemplateSQL = "metric_baselines_%s"

	createMetricBaselineDBTemplateSQL = `create TABLE IF NOT EXISTS %s (
		tenant_id varchar(256) NOT NULL,
		monitored_object_id varchar(256) NOT NULL,
		hour_of_week int NOT NULL,
		baselines jsonb,
		created_timestamp bigint NOT NULL default 0,
		last_modified_timestamp bigint NOT NULL,
		last_reset_timestamp bigint NOT NULL default 0,
		PRIMARY KEY (tenant_id, monitored_object_id, hour_of_week)
	);`
	deleteMetricBaselineDBTemplateSQL = `DROP TABLE %s;`

	wherePrimaryKeySelectorSQL                = "WHERE tenant_id = $1 and monitored_object_id = $2 and hour_of_week = $3"
	whereTeantAndMonitoredObjectIDSelectorSQL = "WHERE tenant_id = $1 and monitored_object_id = $2"
	insertSQL                                 = "INSERT INTO %s (tenant_id, monitored_object_id, hour_of_week, baselines, created_timestamp, last_modified_timestamp) VALUES ($1, $2, $3, $4, $5, $6)"
	upsertSQL                                 = "INSERT INTO %s (tenant_id, monitored_object_id, hour_of_week, baselines, created_timestamp, last_modified_timestamp) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (tenant_id, monitored_object_id, hour_of_week) DO UPDATE SET baselines = EXCLUDED.baselines, last_modified_timestamp = EXCLUDED.last_modified_timestamp"
	getBaselinesByPrimaryKeySQL               = "SELECT baselines FROM %s " + wherePrimaryKeySelectorSQL
	getAllByPrimaryKeySQL                     = "SELECT * FROM %s " + wherePrimaryKeySelectorSQL
	updateSQL                                 = "UPDATE %s SET baselines = $1::jsonb, last_modified_timestamp = $2::bigint, last_reset_timestamp = $3::bigint WHERE tenant_id = $4 and monitored_object_id = $5 and hour_of_week = $6"
	deleteSQL                                 = "DELETE FROM %s " + wherePrimaryKeySelectorSQL
	getBaselinesByMonitoredObjectSQL          = "SELECT * FROM %s " + whereTeantAndMonitoredObjectIDSelectorSQL
	deleteBaselinesByMonitoredObjectSQL       = "DELETE FROM %s " + whereTeantAndMonitoredObjectIDSelectorSQL
	resetBaselinesByMonitoredObjectSQL        = "UPDATE %s SET baselines = $1::jsonb, last_modified_timestamp = $2::bigint, last_reset_timestamp = $3::bigint WHERE tenant_id = $4 and monitored_object_id = $5"
)

type TenantMetricBaselinePostgresDAO struct {
	DB        *sql.DB
	batchSize int64
}

// CreateTenantMetricBaselinePostgresDAO - creates the Postgres DB impl object
func CreateTenantMetricBaselinePostgresDAO() (*TenantMetricBaselinePostgresDAO, error) {
	result := new(TenantMetricBaselinePostgresDAO)

	cfg := gather.GetConfig()
	host := cfg.GetString(gather.CK_args_metricbaselines_ip.String())
	port := cfg.GetInt(gather.CK_args_metricbaselines_port.String())
	user := cfg.GetString(gather.CK_args_metricbaselines_user.String())
	password := cfg.GetString(gather.CK_args_metricbaselines_password.String())
	dbname := cfg.GetString(gather.CK_args_metricbaselines_dbname.String())
	schemaDir := cfg.GetString(gather.CK_args_metricbaselines_schemadir.String())

	result.batchSize = int64(cfg.GetInt(gather.CK_server_datastore_batchsize.String()))

	postgresConnInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	result.DB, err = sql.Open("postgres", postgresConnInfo)
	if err != nil {
		return nil, fmt.Errorf("Unable to open Postgres DB: %s", err.Error())
	}

	// Try to ping the DB to ensure connection is up:
	if err = result.DB.Ping(); err != nil {
		// Creatre the DB if it does not exist
		if strings.Contains(err.Error(), "does not exist") {
			// Create a connection that is not for the DB directly
			tempConnStr := fmt.Sprintf("host=%s port=%d user=%s "+
				"password=%s sslmode=disable",
				host, port, user, password)
			tempConn, err := sql.Open("postgres", tempConnStr)
			if err != nil {
				return nil, fmt.Errorf("Unable to connect to POSTGRES to create DB %s: %s", dbname, err.Error())
			}

			// Try to create the DB
			input, err := ioutil.ReadFile(schemaDir + "/createGatherDB.sql")
			if err != nil {
				return nil, fmt.Errorf("Unable to locate datahub DB schema: %s", err.Error())
			}
			dbString := string(input)
			logger.Log.Debugf("Setting up schema: %s", dbString)

			_, err = tempConn.Exec(dbString)
			if err != nil {
				return nil, fmt.Errorf("Unable create datahub DB: %s", err.Error())
			}

			// Make sure the DB connection works:
			if err = result.DB.Ping(); err != nil {
				return nil, fmt.Errorf("Unable to ping Postgres DB: %s", err.Error())
			}

			tempConn.Close()

		}
	}

	log.Infof("Metric Baseline datastore is POSTGRES located at %s:%d", host, port)

	return result, nil
}

// CreateMetricBaseline - CouchDB implementation of CreateMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) CreateMetricBaseline(metricBaselineReq *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Creating %s for Tenant %s for %s %s for hour of week %d", tenmod.TenantMetricBaselineStr, metricBaselineReq.TenantID, tenmod.TenantMonitoredObjectStr, metricBaselineReq.MonitoredObjectID, metricBaselineReq.HourOfWeek)

	sqlStatement, err := mbdb.DB.Prepare(addTableNameToSQL(insertSQL, metricBaselineReq.TenantID))
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_create")
		return nil, fmt.Errorf("Unable to create insert metric baseline statement template: %s", err.Error())
	}
	defer sqlStatement.Close()

	createTime := datastore.MakeTimestamp()
	_, err = sqlStatement.Exec(metricBaselineReq.TenantID, metricBaselineReq.MonitoredObjectID, metricBaselineReq.HourOfWeek, models.AsJSONString(metricBaselineReq.Baselines), createTime, createTime)
	if err != nil {
		if strings.Contains(err.Error(), duplicateKey) {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "409", "met_bsln_create")
			return nil, fmt.Errorf(datastore.ConflictErrorStr)
		}
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_create")
		return nil, fmt.Errorf("Unable to insert metric baseline: %s", err.Error())
	}

	// fill in the missing values of the response
	metricBaselineReq.ID = datastore.GetDataIDFromFullID(datastore.GenerateID(metricBaselineReq, metricBaselineType))
	metricBaselineReq.Datatype = metricBaselineType
	metricBaselineReq.REV = fmt.Sprintf("%d", createTime)
	metricBaselineReq.CreatedTimestamp = createTime
	metricBaselineReq.LastModifiedTimestamp = createTime

	logger.Log.Debugf("Completed baseline insert for Tenant %s Monitored Object %s Hour Of Week %d", metricBaselineReq.TenantID, metricBaselineReq.TenantID, metricBaselineReq.HourOfWeek)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_create")
	return metricBaselineReq, nil
}

// UpdateMetricBaseline - CouchDB implementation of UpdateMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) UpdateMetricBaseline(metricBaselineReq *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Updating %s for Tenant %s for %s %s for hour of week %d", tenmod.TenantMetricBaselineStr, metricBaselineReq.TenantID, tenmod.TenantMonitoredObjectStr, metricBaselineReq.MonitoredObjectID, metricBaselineReq.HourOfWeek)

	// Start a txn
	txn, err := mbdb.DB.Begin()
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_update")
		return nil, fmt.Errorf("Unable to create update transaction: %s", err.Error())
	}

	existing, err := mbdb.GetMetricBaseline(metricBaselineReq.TenantID, metricBaselineReq.ID)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_update")
		return nil, fmt.Errorf("Unable to update insert metric baseline: %s", err.Error())
	}

	sqlStatement, err := txn.Prepare(addTableNameToSQL(updateSQL, metricBaselineReq.TenantID))
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_update")
		return nil, fmt.Errorf("Unable to create update metric baseline statement template: %s", err.Error())
	}
	defer sqlStatement.Close()

	// Make sure the revision is correct
	if existing.LastModifiedTimestamp != metricBaselineReq.LastModifiedTimestamp {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "409", "met_bsln_update")
		return nil, fmt.Errorf("Unable to update metric baseline: incorrect revision %d, was expecting %d", metricBaselineReq.LastModifiedTimestamp, existing.LastModifiedTimestamp)
	}

	// Update existing for change
	modTime := datastore.MakeTimestamp()
	existing.Baselines = metricBaselineReq.Baselines
	if len(existing.Baselines) == 0 {
		existing.LastResetTimestamp = modTime
	}
	existing.LastModifiedTimestamp = modTime

	_, err = sqlStatement.Exec(models.AsJSONString(existing.Baselines), existing.LastModifiedTimestamp, existing.LastResetTimestamp, metricBaselineReq.TenantID, metricBaselineReq.MonitoredObjectID, metricBaselineReq.HourOfWeek)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_update")
		return nil, fmt.Errorf("Unable to update metric baseline: %s", err.Error())
	}

	// Commit the transaction
	err = txn.Commit()
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_update")
		return nil, fmt.Errorf("Unable to commit bulk update transaction: %s", err.Error())
	}

	// fill in the missing values of the response
	existing.ID = datastore.GetDataIDFromFullID(datastore.GenerateID(metricBaselineReq, metricBaselineType))
	existing.Datatype = metricBaselineType
	existing.REV = fmt.Sprintf("%d", existing.LastModifiedTimestamp)

	logger.Log.Debugf("Completed baseline update for Tenant %s Monitored Object %s Hour Of Week %d", metricBaselineReq.TenantID, metricBaselineReq.TenantID, metricBaselineReq.HourOfWeek)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_update")
	return existing, nil
}

// GetMetricBaseline - CouchDB implementation of GetMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) GetMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Fetching %s for Tenant %s for ID %s", tenmod.TenantMetricBaselineStr, tenantID, dataID)

	sqlStatement, err := mbdb.DB.Prepare(addTableNameToSQL(getAllByPrimaryKeySQL, tenantID))
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_get")
		return nil, fmt.Errorf("Unable to create fetch metric baseline statement template: %s", err.Error())
	}
	defer sqlStatement.Close()

	idDelimeter := strings.LastIndex(dataID, "_")
	if idDelimeter == -1 {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_get")
		return nil, fmt.Errorf("Unable to fetch Metric Baseline: Invalid ID format %s", dataID)
	}
	monObjID := dataID[0:idDelimeter]
	hourOfWeek, err := strconv.ParseInt(dataID[idDelimeter+1:], 10, 64)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_get")
		return nil, fmt.Errorf("Unable to fetch Metric Baseline: Hour of Week not an integer %s", dataID)
	}
	row := sqlStatement.QueryRow(tenantID, monObjID, hourOfWeek)

	var baselineResultContainer []byte
	result := tenmod.MetricBaseline{}
	err = row.Scan(&result.TenantID, &result.MonitoredObjectID, &result.HourOfWeek, &baselineResultContainer, &result.CreatedTimestamp, &result.LastModifiedTimestamp, &result.LastResetTimestamp)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "404", "met_bsln_get")
			return nil, fmt.Errorf(datastore.NotFoundStr)
		}
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_get")
		return nil, fmt.Errorf("Unable to read query result: %s", err.Error())
	}

	// fill in the missing values of the response
	err = json.Unmarshal(baselineResultContainer, &result.Baselines)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_get")
		return nil, fmt.Errorf("Unable to convert query result: %s", err.Error())
	}

	result.ID = datastore.GetDataIDFromFullID(datastore.GenerateID(result, metricBaselineType))
	result.Datatype = metricBaselineType
	result.REV = fmt.Sprintf("%d", result.LastModifiedTimestamp)

	logger.Log.Debugf("Retrieved Metric Baseline for Tenant %s Id %s", tenantID, dataID)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_get")
	return &result, nil
}

// DeleteMetricBaseline - CouchDB implementation of DeleteMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) DeleteMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Deleting %s for Tenant %s for ID %s", tenmod.TenantMetricBaselineStr, tenantID, dataID)

	existing, err := mbdb.GetMetricBaseline(tenantID, dataID)
	if err != nil {
		return nil, err
	}

	sqlStatement, err := mbdb.DB.Prepare(addTableNameToSQL(deleteSQL, tenantID))
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_del")
		return nil, fmt.Errorf("Unable to create delete metric baseline statement template: %s", err.Error())
	}
	defer sqlStatement.Close()

	_, err = sqlStatement.Exec(tenantID, existing.MonitoredObjectID, existing.HourOfWeek)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_del")
		return nil, fmt.Errorf("Unable to convert query result: %s", err.Error())
	}

	logger.Log.Debugf("Retrieved Metric Baseline for Tenant %s Id %s", tenantID, dataID)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_del")
	return existing, nil
}

func (mbdb *TenantMetricBaselinePostgresDAO) UpdateMetricBaselineForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32, baselineData *tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Updating %s for Tenant %s for %s %s for hour of week %d for a single entry", tenmod.TenantMetricBaselineStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)

	existing, err := mbdb.GetMetricBaseline(tenantID, datastore.GetDataIDFromFullID(datastore.GenerateMetricBaselineID(monObjID, hourOfWeek)))
	if err != nil {
		if !strings.Contains(err.Error(), datastore.NotFoundStr) {
			// Error was something permanent, return it
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_single_val_update")
			return nil, fmt.Errorf("Unable to update metric baseline single entry: %s", err.Error())
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			HourOfWeek:        hourOfWeek,
			Baselines:         []*tenmod.MetricBaselineData{baselineData},
		}

		return mbdb.CreateMetricBaseline(&createObj)
	}

	existing.MergeBaseline(baselineData)

	result, err := mbdb.UpdateMetricBaseline(existing)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_single_val_update")
		return nil, err
	}

	// fill in the missing values of the response
	result.ID = datastore.GetDataIDFromFullID(datastore.GenerateID(existing, metricBaselineType))
	result.Datatype = metricBaselineType
	result.REV = fmt.Sprintf("%d", existing.LastModifiedTimestamp)

	logger.Log.Debugf("Completed baseline update for Tenant %s Monitored Object %s Hour Of Week %d for a single entry", result.TenantID, result.TenantID, result.HourOfWeek)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_single_val_update")
	return result, nil
}

// UpdateMetricBaselineForHourOfWeekWithCollection - couchDB implementation of UpdateMetricBaselineForHourOfWeekWithCollection
func (mbdb *TenantMetricBaselinePostgresDAO) UpdateMetricBaselineForHourOfWeekWithCollection(tenantID string, monObjID string, hourOfWeek int32, baselineDataCollection []*tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Updating %s for Tenant %s for %s %s for hour of week %d for multiple entries", tenmod.TenantMetricBaselineStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)

	existing, err := mbdb.GetMetricBaseline(tenantID, datastore.GetDataIDFromFullID(datastore.GenerateMetricBaselineID(monObjID, hourOfWeek)))
	if err != nil {
		if !strings.Contains(err.Error(), datastore.NotFoundStr) {
			// Error was something permanent, return it
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_multi_val_update")
			return nil, fmt.Errorf("Unable to update metric baseline multiple entries: %s", err.Error())
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			HourOfWeek:        hourOfWeek,
			Baselines:         baselineDataCollection,
		}

		return mbdb.CreateMetricBaseline(&createObj)
	}

	existing.MergeBaselines(baselineDataCollection)

	result, err := mbdb.UpdateMetricBaseline(existing)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_multi_val_update")
		return nil, err
	}

	// fill in the missing values of the response
	result.ID = datastore.GetDataIDFromFullID(datastore.GenerateID(existing, metricBaselineType))
	result.Datatype = metricBaselineType
	result.REV = fmt.Sprintf("%d", existing.LastModifiedTimestamp)

	logger.Log.Debugf("Completed baseline update for Tenant %s Monitored Object %s Hour Of Week %d for multiple entries", result.TenantID, result.TenantID, result.HourOfWeek)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_multi_val_update")
	return result, nil
}

func (mbdb *TenantMetricBaselinePostgresDAO) GetMetricBaselineForMonitoredObjectForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32) ([]*tenmod.MetricBaselineData, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Retrieving %ss for Tenant %s for %s %s for hour of week %s", tenmod.TenantMetricBaselineStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)

	sqlStatement, err := mbdb.DB.Prepare(addTableNameToSQL(getBaselinesByPrimaryKeySQL, tenantID))
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_hrwk_getall")
		return nil, fmt.Errorf("Unable to create get metric baseline statement template: %s", err.Error())
	}
	defer sqlStatement.Close()
	row := sqlStatement.QueryRow(tenantID, monObjID, hourOfWeek)

	var baselineBytes []byte
	err = row.Scan(&baselineBytes)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "404", "met_bsln_hrwk_getall")
			return nil, fmt.Errorf(datastore.NotFoundStr)
		}
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_hrwk_getall")
		return nil, fmt.Errorf("Unable to read query result: %s", err.Error())
	}

	result := []*tenmod.MetricBaselineData{}
	err = json.Unmarshal(baselineBytes, &result)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_hrwk_getall")
		return nil, fmt.Errorf("Unable to convert query result: %s", err.Error())
	}

	logger.Log.Debugf("Completed baseline fetch for Tenant %s Monitored Object %s Hour Of Week %d", tenantID, monObjID, hourOfWeek)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_hrwk_getall")
	return result, nil
}

// GetMetricBaselinesFor - note that this function will return results that are not stored in the DB as new "empty" items so that they can be populated
// in a subsequent bulk PUT call.
func (mbdb *TenantMetricBaselinePostgresDAO) GetMetricBaselinesFor(tenantID string, moIDToHourOfWeekMap map[string][]int32, addNotFoundValuesInResponse bool) ([]*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED FOR POSTGRES")
}

func (mbdb *TenantMetricBaselinePostgresDAO) BulkUpdateMetricBaselines(tenantID string, entries []*swagmodels.MetricBaselineBulkUpdateRequestDataAttributesItems0) ([]*common.BulkOperationResult, error) {
	methodStartTime := time.Now()
	logger.Log.Debugf("Bulk updating %s for Tenant %s", tenmod.TenantMetricBaselineStr, tenantID)

	numMbs := len(entries)
	if int64(numMbs) > mbdb.batchSize {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_bulk_update")
		return nil, fmt.Errorf("Too many Monitored Objects in bulk request. Limit is %d but request contains %d", mbdb.batchSize, numMbs)
	}

	// Start a txn
	txn, err := mbdb.DB.Begin()
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_bulk_update")
		return nil, fmt.Errorf("Unable to create bulk update transaction: %s", err.Error())
	}

	// Use upsert to add the records to the DB
	currentTime := datastore.MakeTimestamp()
	for _, val := range entries {
		sqlStatement, err := txn.Prepare(addTableNameToSQL(upsertSQL, tenantID))
		if err != nil {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_bulk_update")
			txn.Rollback()
			sqlStatement.Close()
			return nil, fmt.Errorf("Unable to create upsert statement template: %s", err.Error())
		}
		_, err = sqlStatement.Exec(tenantID, val.MonitoredObjectID, val.HourOfWeek, models.AsJSONString(val.Baselines), currentTime, currentTime)
		if err != nil {
			txn.Rollback()
			sqlStatement.Close()
			return nil, fmt.Errorf("Unable to upsert record for Tenant %s Monitored Object %s Hour of Week %d: %s", tenantID, val.MonitoredObjectID, val.HourOfWeek, err.Error())
		}

		sqlStatement.Close()
	}

	// Commit the transaction
	err = txn.Commit()
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_bulk_update")
		return nil, fmt.Errorf("Unable to commit bulk update transaction: %s", err.Error())
	}

	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_bulk_update")
	logger.Log.Debugf("Successfully completed bulk update for Tenant %s", tenantID)
	return []*common.BulkOperationResult{}, nil
}

type MetricBaselineRowMapper struct {
	TenantID              string `json:"tenant"`
	HourOfWeek            int32  `json:"hourOfWeek"`
	Baselines             []byte `json:"baselines"`
	MonitoredObjectID     string `json:"monitoredObjectId"`
	LastModifiedTimestamp int64  `json:"lastModifiedTimestamp"`
	CreatedTimestamp      int64  `json:"createdTimestamp"`
	LastResetTimestamp    int64  `json:"lastResetimestamp"`
}

func (mbdb *TenantMetricBaselinePostgresDAO) BulkUpdateMetricBaselinesFromList(tenantID string, baselineUpdateList []*tenmod.MetricBaseline) ([]*common.BulkOperationResult, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

func (mbdb *TenantMetricBaselinePostgresDAO) GetMetricBaselineForMonitoredObject(tenantID string, monObjID string) ([]*tenmod.MetricBaseline, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Fetching all %s for Tenant %s for Monitored Object ID %s", tenmod.TenantMetricBaselineStr, tenantID, monObjID)

	sqlStatement, err := mbdb.DB.Prepare(addTableNameToSQL(getBaselinesByMonitoredObjectSQL, tenantID))
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_by_monobj_get")
		return nil, fmt.Errorf("Unable to create fetch metric baseline by monitored object statement template: %s", err.Error())
	}
	defer sqlStatement.Close()

	rows, err := sqlStatement.Query(tenantID, monObjID)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_by_monobj_get")
		return nil, fmt.Errorf("Unable to fetch metric baseline data for monitored object %s statement template: %s", monObjID, err.Error())
	}
	defer rows.Close()

	resultList := []*tenmod.MetricBaseline{}

	for rows.Next() {
		var baselineResultContainer []byte
		result := tenmod.MetricBaseline{}
		err = rows.Scan(&result.TenantID, &result.MonitoredObjectID, &result.HourOfWeek, &baselineResultContainer, &result.CreatedTimestamp, &result.LastModifiedTimestamp, &result.LastResetTimestamp)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "404", "met_bsln_by_monobj_get")
				return nil, fmt.Errorf(datastore.NotFoundStr)
			}
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_by_monobj_get")
			return nil, fmt.Errorf("Unable to read query result: %s", err.Error())
		}

		// fill in the missing values of the response
		err = json.Unmarshal(baselineResultContainer, &result.Baselines)
		if err != nil {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_by_monobj_get")
			return nil, fmt.Errorf("Unable to convert query result: %s", err.Error())
		}

		result.ID = datastore.GetDataIDFromFullID(datastore.GenerateID(result, metricBaselineType))
		result.Datatype = metricBaselineType
		result.REV = fmt.Sprintf("%d", result.LastModifiedTimestamp)

		resultList = append(resultList, &result)
	}

	if len(resultList) == 0 {
		return nil, fmt.Errorf(datastore.NotFoundStr)
	}

	logger.Log.Debugf("Retrieved %d Metric Baselines for Tenant %s Monitored Object ID %s", len(resultList), tenantID, monObjID)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_by_monobj_get")
	return resultList, nil
}

func (mbdb *TenantMetricBaselinePostgresDAO) DeleteMetricBaselineForMonitoredObject(tenantID string, monObjID string, reset bool) error {
	methodStartTime := time.Now()
	sqlStatementStr := deleteBaselinesByMonitoredObjectSQL
	operation := "Delete"
	if reset {
		operation = "Reset"
		sqlStatementStr = resetBaselinesByMonitoredObjectSQL
	}

	logger.Log.Debugf("%s all %s for Tenant %s for Monitored Object ID %s", operation, tenmod.TenantMetricBaselineStr, tenantID, monObjID)
	sqlStatement, err := mbdb.DB.Prepare(addTableNameToSQL(sqlStatementStr, tenantID))
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_by_monobj_delete")
		return fmt.Errorf("Unable to create %s metric baseline by monitored object statement template: %s", operation, err.Error())
	}
	defer sqlStatement.Close()

	if reset {
		updateTime := datastore.MakeTimestamp()
		_, err := sqlStatement.Exec("[]", updateTime, updateTime, tenantID, monObjID)
		if err != nil {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_by_monobj_delete")
			return fmt.Errorf("Unable to %s metric baseline data for monitored object %s statement template: %s", operation, monObjID, err.Error())
		}
	} else {
		_, err := sqlStatement.Exec(tenantID, monObjID)
		if err != nil {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_by_monobj_delete")
			return fmt.Errorf("Unable to %s metric baseline data for monitored object %s statement template: %s", operation, monObjID, err.Error())
		}
	}

	logger.Log.Debugf("%s of all Metric Baselines for Tenant %s Monitored Object ID %s complete", operation, tenantID, monObjID)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_by_monobj_get")
	return nil
}

func (mbdb *TenantMetricBaselinePostgresDAO) CreateMetricBaselineDB(tenantID string) error {
	methodStartTime := time.Now()

	logger.Log.Debugf("Creating Metric Baselines table for Tenant %s", tenantID)

	sqlStr := addTableNameToSQL(createMetricBaselineDBTemplateSQL, tenantID)
	sqlStatement, err := mbdb.DB.Prepare(sqlStr)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_tbl_create")
		return fmt.Errorf("Unable to create baseline table creation statement template: %s", err.Error())
	}
	defer sqlStatement.Close()
	_, err = sqlStatement.Exec()
	if err != nil {
		return err
	}

	logger.Log.Debugf("Created Metric Baselines table for Tenant %s", tenantID)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_tbl_create")
	return nil
}

func (mbdb *TenantMetricBaselinePostgresDAO) DeleteMetricBaselineDB(tenantID string) error {
	methodStartTime := time.Now()

	logger.Log.Debugf("Deleting Metric Baselines table for Tenant %s", tenantID)

	sqlStr := addTableNameToSQL(deleteMetricBaselineDBTemplateSQL, tenantID)
	sqlStatement, err := mbdb.DB.Prepare(sqlStr)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_tbl_create")
		return fmt.Errorf("Unable to create baseline table creation statement template: %s", err.Error())
	}
	defer sqlStatement.Close()
	_, err = sqlStatement.Exec()
	if err != nil {
		return err
	}

	logger.Log.Debugf("Deleted Metric Baselines table for Tenant %s", tenantID)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_tbl_create")
	return nil
}

func getMetricsBaselineTableNameSQL(tenantID string) string {
	tenantIDTransformed := strings.Replace(tenantID, "-", "_", -1)
	return fmt.Sprintf(metricBaselineTableNameTemplateSQL, tenantIDTransformed)
}

func addTableNameToSQL(sqlStr string, tenantID string) string {
	return fmt.Sprintf(sqlStr, getMetricsBaselineTableNameSQL(tenantID))
}
