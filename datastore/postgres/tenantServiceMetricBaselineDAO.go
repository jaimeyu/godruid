package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	TmpMetricBaselineTableName = "tmp_metric_baselines"
	MetricBaselineTableName    = "metric_baselines"
	bulkUpdateSQL              = `UPDATE metric_baselines set baselines = tmp_metric_baselines.baselines, 
						last_modified_timestamp = tmp_metric_baselines.last_modified_timestamp 
						from tmp_metric_baselines where metric_baselines.tenant_id = tmp_metric_baselines.tenant_id
						and metric_baselines.monitored_object_id = tmp_metric_baselines.monitored_object_id and metric_baselines.hour_of_week = tmp_metric_baselines.hour_of_week 
						returning tmp_metric_baselines.tenant_id, tmp_metric_baselines.monitored_object_id, tmp_metric_baselines.hour_of_week, tmp_metric_baselines.baselines`
)

var (
	upsertSQL          = fmt.Sprintf("INSERT INTO %s (tenant_id, monitored_object_id, hour_of_week, baselines, created_timestamp, last_modified_timestamp) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (tenant_id, monitored_object_id, hour_of_week) DO UPDATE SET baselines = EXCLUDED.baselines, last_modified_timestamp = EXCLUDED.last_modified_timestamp", MetricBaselineTableName)
	getByPrimaryKeySQL = fmt.Sprintf("SELECT baselines FROM %s WHERE tenant_id = $1 and monitored_object_id = $2 and hour_of_week = $3", MetricBaselineTableName)
)

type TenantMetricBaselinePostgresDAO struct {
	DB        *sql.DB
	batchSize int64
}

// CreateUserServiceDAO - creates an instance of the User Service datastore that has been
// implemented using a Postgres DB.
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

			// Try to create the table
			tableInput, err := ioutil.ReadFile(schemaDir + "/initMetricBaselines.sql")
			if err != nil {
				return nil, fmt.Errorf("Unable to locate metric baseline table schema: %s", err.Error())
			}
			tableString := string(tableInput)
			logger.Log.Debugf("Setting up table: %s", tableString)

			_, err = result.DB.Exec(tableString)
			if err != nil {
				return nil, fmt.Errorf("Unable create metric baseline table: %s", err.Error())
			}

			tempConn.Close()

		}
	}

	log.Infof("Metric Baseline datastore is POSTGRES located at %s:%d", host, port)

	return result, nil
}

// CreateMetricBaseline - CouchDB implementation of CreateMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) CreateMetricBaseline(metricBaselineReq *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	// logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(metricBaselineReq))
	// metricBaselineReq.ID = ds.GenerateID(metricBaselineReq, string(tenmod.TenantMetricBaselineType))
	// tenantID := ds.PrependToDataID(metricBaselineReq.TenantID, string(admmod.TenantType))

	// // Make sure there is no existing record for this id:
	// dataID := generateMetricBaselineID(metricBaselineReq.MonitoredObjectID, metricBaselineReq.HourOfWeek)
	// existing, err := tsd.GetMetricBaseline(metricBaselineReq.TenantID, dataID)
	// if err != nil {
	// 	if !strings.Contains(err.Error(), ds.NotFoundStr) {
	// 		return nil, fmt.Errorf("Unable to create %s. Receieved this error when checking for existing %s record: %s", tenmod.TenantMetricBaselineStr, tenmod.TenantMetricBaselineStr, err.Error())
	// 	}
	// }
	// if existing != nil {
	// 	return nil, fmt.Errorf(ds.ConflictStr)
	// }

	// dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	// dataContainer := &tenmod.MetricBaseline{}
	// if err := createDataInCouch(dbName, metricBaselineReq, dataContainer, string(tenmod.TenantMetricBaselineType), tenmod.TenantMetricBaselineStr); err != nil {
	// 	return nil, err
	// }

	// logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	// return dataContainer, nil
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

// UpdateMetricBaseline - CouchDB implementation of UpdateMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) UpdateMetricBaseline(metricBaselineReq *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	// logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(metricBaselineReq))
	// metricBaselineReq.ID = ds.PrependToDataID(metricBaselineReq.ID, string(tenmod.TenantMetricBaselineType))
	// tenantID := ds.PrependToDataID(metricBaselineReq.TenantID, string(admmod.TenantType))

	// dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	// dataContainer := &tenmod.MetricBaseline{}
	// if err := updateDataInCouch(dbName, metricBaselineReq, dataContainer, string(tenmod.TenantMetricBaselineType), tenmod.TenantMetricBaselineStr); err != nil {
	// 	return nil, err
	// }
	// logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	// return dataContainer, nil
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

// GetMetricBaseline - CouchDB implementation of GetMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) GetMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	// logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantMetricBaselineStr, dataID)
	// dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMetricBaselineType))
	// tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	// dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	// dataContainer := tenmod.MetricBaseline{}
	// if err := getDataFromCouch(dbName, dataID, &dataContainer, tenmod.TenantMetricBaselineStr); err != nil {
	// 	return nil, err
	// }
	// logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	// return &dataContainer, nil
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

// DeleteMetricBaseline - CouchDB implementation of DeleteMetricBaseline
func (mbdb *TenantMetricBaselinePostgresDAO) DeleteMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	// logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantMetricBaselineStr, dataID)
	// dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMetricBaselineType))
	// tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	// dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	// dataContainer := tenmod.MetricBaseline{}
	// if err := deleteDataFromCouch(dbName, dataID, &dataContainer, tenmod.TenantMetricBaselineStr); err != nil {
	// 	return nil, err
	// }
	// logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	// return &dataContainer, nil
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

func (mbdb *TenantMetricBaselinePostgresDAO) UpdateMetricBaselineForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32, baselineData *tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	// if logger.IsDebugEnabled() {
	// 	logger.Log.Debugf("Updating %s for %s %s for %s %s for hour of week %d", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)
	// }

	// dataID := ds.GetDataIDFromFullID(generateMetricBaselineID(monObjID, hourOfWeek))
	// existing, err := tsd.GetMetricBaseline(tenantID, dataID)
	// if err != nil {
	// 	if !strings.Contains(err.Error(), ds.NotFoundStr) {
	// 		// Error was something permanent, return it
	// 		return nil, err
	// 	}

	// 	// Error was that the Baseline does not exist for this Monitored Object, let's create it
	// 	createObj := tenmod.MetricBaseline{
	// 		MonitoredObjectID: monObjID,
	// 		TenantID:          tenantID,
	// 		HourOfWeek:        hourOfWeek,
	// 		Baselines:         []*tenmod.MetricBaselineData{baselineData},
	// 	}

	// 	return tsd.CreateMetricBaseline(&createObj)
	// }

	// existing.MergeBaseline(baselineData)
	// existing.ID = ds.GetDataIDFromFullID(existing.ID)

	// updated, err := tsd.UpdateMetricBaseline(existing)
	// if err != nil {
	// 	return nil, err
	// }

	// logger.Log.Debugf("Updated %s for %s %s %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID)
	// return updated, nil
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

// UpdateMetricBaselineForHourOfWeekWithCollection - couchDB implementation of UpdateMetricBaselineForHourOfWeekWithCollection
func (mbdb *TenantMetricBaselinePostgresDAO) UpdateMetricBaselineForHourOfWeekWithCollection(tenantID string, monObjID string, hourOfWeek int32, baselineDataCollection []*tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	// if logger.IsDebugEnabled() {
	// 	logger.Log.Debugf("Updating %s for %s %s for %s %s for hour of week %d with multiple values", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)
	// }

	// dataID := ds.GetDataIDFromFullID(generateMetricBaselineID(monObjID, hourOfWeek))
	// existing, err := tsd.GetMetricBaseline(tenantID, dataID)
	// if err != nil {
	// 	if !strings.Contains(err.Error(), ds.NotFoundStr) {
	// 		// Error was something permanent, return it
	// 		return nil, err
	// 	}

	// 	// Error was that the Baseline does not exist for this Monitored Object, let's create it
	// 	createObj := tenmod.MetricBaseline{
	// 		MonitoredObjectID: monObjID,
	// 		TenantID:          tenantID,
	// 		HourOfWeek:        hourOfWeek,
	// 		Baselines:         baselineDataCollection,
	// 	}

	// 	return tsd.CreateMetricBaseline(&createObj)
	// }

	// existing.MergeBaselines(baselineDataCollection)
	// existing.ID = ds.GetDataIDFromFullID(existing.ID)

	// updated, err := tsd.UpdateMetricBaseline(existing)
	// if err != nil {
	// 	return nil, err
	// }

	// logger.Log.Debugf("Updated %s for %s %s %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID)
	// return updated, nil
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

func (mbdb *TenantMetricBaselinePostgresDAO) GetMetricBaselineForMonitoredObjectForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32) ([]*tenmod.MetricBaselineData, error) {
	methodStartTime := time.Now()

	logger.Log.Debugf("Retrieving %ss for Tenant %s for %s %s for hour of week %s", tenmod.TenantMetricBaselineStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)

	sqlStatement, err := mbdb.DB.Prepare(getByPrimaryKeySQL)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_hrwk_getall")
		return nil, fmt.Errorf("Unable to create get metric baseline statement template: %s", err)
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
		return nil, fmt.Errorf("Unable to read query result: %s", err)
	}

	result := []*tenmod.MetricBaselineData{}
	err = json.Unmarshal(baselineBytes, &result)
	if err != nil {
		monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_hrwk_getall")
		return nil, fmt.Errorf("Unable to convert query result: %s", err)
	}

	logger.Log.Debugf("Completed baseline fetch for Tenant %s Monitored Object %s Hour Of Week %d", tenantID, monObjID, hourOfWeek)
	monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "200", "met_bsln_hrwk_getall")
	return result, nil
}

// GetMetricBaselinesFor - note that this function will return results that are not stored in the DB as new "empty" items so that they can be populated
// in a subsequent bulk PUT call.
func (mbdb *TenantMetricBaselinePostgresDAO) GetMetricBaselinesFor(tenantID string, moIDToHourOfWeekMap map[string][]int32, addNotFoundValuesInResponse bool) ([]*tenmod.MetricBaseline, error) {
	// logger.Log.Debugf("Bulk retrieving %ss for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID)

	// tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	// dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	// resource, err := couchdb.NewResource(dbName, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// numEntries := 0
	// keyMap := []string{}
	// for moID, hourOfWeekList := range moIDToHourOfWeekMap {
	// 	for _, hour := range hourOfWeekList {
	// 		numEntries++
	// 		keyMap = append(keyMap, generateMetricBaselineID(moID, hour))
	// 	}
	// }
	// if int64(numEntries) > tsd.batchSize {
	// 	return nil, fmt.Errorf("Too many Monitored Objects in bulk request. Limit is %d but request contains %d", tsd.batchSize, numEntries)
	// }

	// mbTypeString := string(tenmod.TenantMetricBaselineType)
	// addTenantId := ds.GetDataIDFromFullID(tenantID)

	// // Build request to fetch existing records
	// qp := &url.Values{}
	// qp.Add("include_docs", "true")
	// requestBody := map[string]interface{}{}
	// requestBody["keys"] = keyMap

	// fetchedGetData, err := fetchAllDocsWithQP(qp, requestBody, resource)
	// if err != nil {
	// 	return nil, err
	// }

	// // Need to convert retrieved records:
	// rows, ok := fetchedGetData["rows"].([]interface{})
	// if !ok {
	// 	return nil, fmt.Errorf("Unable to convert bulk fetch response: response from fetch has invalid format %s", models.AsJSONString(rows))
	// }
	// convertedRetrievedBaselines := []*tenmod.MetricBaseline{}
	// for _, row := range rows {
	// 	value, ok := row.(map[string]interface{})
	// 	if !ok {
	// 		return nil, fmt.Errorf("Unable to convert single row to generic object %s", models.AsJSONString(value))
	// 	}
	// 	// Make sure that any records that are not found are added to the response as new empty records
	// 	errorStr, okError := value["error"].(string)
	// 	fetchedBaseline, okBaseline := value["doc"].(map[string]interface{})

	// 	// Handle weird case when there is no eror and no data returned:
	// 	if !okBaseline && errorStr == "" && addNotFoundValuesInResponse {
	// 		if !addMetricBaselineRecordToBulkFetchResponse(mbTypeString, addTenantId, value, &convertedRetrievedBaselines) {
	// 			return nil, fmt.Errorf("Can't create record for: %s", models.AsJSONString(row))
	// 		}
	// 		continue
	// 	}

	// 	if okError && addNotFoundValuesInResponse {
	// 		if errorStr == metricBaselineBulkFetchNotFoundStr {
	// 			if !addMetricBaselineRecordToBulkFetchResponse(mbTypeString, addTenantId, value, &convertedRetrievedBaselines) {
	// 				return nil, fmt.Errorf("Can't create record for: %s", models.AsJSONString(row))
	// 			}
	// 			continue
	// 		}

	// 		return nil, fmt.Errorf("Can't find record for: %s", models.AsJSONString(row))
	// 	}

	// 	// Handle the retrieved record
	// 	if !okBaseline {
	// 		return nil, fmt.Errorf("Unable to convert bulk fetch single value response to baseline: %s", err.Error())
	// 	}
	// 	stripPrefixFromID(fetchedBaseline)

	// 	// Convert a fetched doc to the proper type for merging
	// 	dataContainer := tenmod.MetricBaseline{}
	// 	if err = convertGenericCouchDataToObject(fetchedBaseline, &dataContainer, tenmod.TenantMetricBaselineStr); err != nil {
	// 		return nil, err
	// 	}

	// 	convertedRetrievedBaselines = append(convertedRetrievedBaselines, &dataContainer)
	// }

	// // Return the converted baseline data
	// logger.Log.Debugf("Completed bulk retrieval of %ss for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID)
	// return convertedRetrievedBaselines, nil
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

// func addMetricBaselineRecordToBulkFetchResponse(mbTypeString string, tenantID string, value map[string]interface{}, baselineArray *[]*tenmod.MetricBaseline) bool {
// 	moID, ok := value["key"].(string)
// 	if !ok {
// 		return ok
// 	}
// 	ts := ds.MakeTimestamp()
// 	strippedMOID := ds.GetDataIDFromFullID(moID)
// 	strippedMOIDParts := strings.Split(strippedMOID, "_")
// 	how, _ := strconv.ParseInt(strippedMOIDParts[1], 10, 32)
// 	addObject := tenmod.MetricBaseline{
// 		ID:                    strippedMOID,
// 		Datatype:              mbTypeString,
// 		TenantID:              tenantID,
// 		MonitoredObjectID:     strippedMOIDParts[0],
// 		HourOfWeek:            int32(how),
// 		Baselines:             []*tenmod.MetricBaselineData{},
// 		CreatedTimestamp:      ts,
// 		LastModifiedTimestamp: ts,
// 	}
// 	*baselineArray = append(*baselineArray, &addObject)

// 	return true
// }

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
		sqlStatement, err := txn.Prepare(upsertSQL)
		if err != nil {
			monitoring.TrackPostgresTimeMetricInSeconds(monitoring.PostgresAPIMethodDurationType, methodStartTime, "500", "met_bsln_bulk_update")
			txn.Rollback()
			sqlStatement.Close()
			return nil, fmt.Errorf("Unable to create upsert statement template: %s", err)
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
		return nil, fmt.Errorf("Unable to commit bulk update transaction: %s", err)
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
}

func (mbdb *TenantMetricBaselinePostgresDAO) BulkUpdateMetricBaselinesFromList(tenantID string, baselineUpdateList []*tenmod.MetricBaseline) ([]*common.BulkOperationResult, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}

// func generateMetricBaselineID(monObjID string, hourOfWeek int32) string {
// 	return fmt.Sprintf("%s_2_%s_%d", string(tenmod.TenantMetricBaselineType), monObjID, hourOfWeek)
// }
