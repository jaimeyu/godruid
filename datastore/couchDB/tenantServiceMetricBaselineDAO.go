package couchDB

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	couchdb "github.com/leesper/couchdb-golang"
)

const (
	metricBaselinesByNameIndex         = "_design/baselineIndex/_view/byName"
	metricBaselineBulkFetchNotFoundStr = "not_found"
)

// CreateMetricBaseline - CouchDB implementation of CreateMetricBaseline
func (tsd *TenantServiceDatastoreCouchDB) CreateMetricBaseline(metricBaselineReq *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(metricBaselineReq))
	metricBaselineReq.ID = ds.GenerateID(metricBaselineReq, string(tenmod.TenantMetricBaselineType))
	tenantID := ds.PrependToDataID(metricBaselineReq.TenantID, string(admmod.TenantType))

	// Make sure there is no existing record for this id:
	existing, err := tsd.GetMetricBaseline(metricBaselineReq.TenantID, metricBaselineReq.MonitoredObjectID)
	if err != nil {
		if !strings.Contains(err.Error(), ds.NotFoundStr) {
			return nil, fmt.Errorf("Unable to create %s. Receieved this error when checking for existing %s record: %s", tenmod.TenantMetricBaselineStr, tenmod.TenantMetricBaselineStr, err.Error())
		}
	}
	if existing != nil {
		return nil, fmt.Errorf(ds.ConflictStr)
	}

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	dataContainer := &tenmod.MetricBaseline{}
	if err := createDataInCouch(dbName, metricBaselineReq, dataContainer, string(tenmod.TenantMetricBaselineType), tenmod.TenantMetricBaselineStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateMetricBaseline - CouchDB implementation of UpdateMetricBaseline
func (tsd *TenantServiceDatastoreCouchDB) UpdateMetricBaseline(metricBaselineReq *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(metricBaselineReq))
	metricBaselineReq.ID = ds.PrependToDataID(metricBaselineReq.ID, string(tenmod.TenantMetricBaselineType))
	tenantID := ds.PrependToDataID(metricBaselineReq.TenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	dataContainer := &tenmod.MetricBaseline{}
	if err := updateDataInCouch(dbName, metricBaselineReq, dataContainer, string(tenmod.TenantMetricBaselineType), tenmod.TenantMetricBaselineStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetMetricBaseline - CouchDB implementation of GetMetricBaseline
func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantMetricBaselineStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMetricBaselineType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	dataContainer := tenmod.MetricBaseline{}
	if err := getDataFromCouch(dbName, dataID, &dataContainer, tenmod.TenantMetricBaselineStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteMetricBaseline - CouchDB implementation of DeleteMetricBaseline
func (tsd *TenantServiceDatastoreCouchDB) DeleteMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantMetricBaselineStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMetricBaselineType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	dataContainer := tenmod.MetricBaseline{}
	if err := deleteDataFromCouch(dbName, dataID, &dataContainer, tenmod.TenantMetricBaselineStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllMetricBaselines - CouchDB implementation of GetAllMetricBaselines
func (tsd *TenantServiceDatastoreCouchDB) GetAllMetricBaselines(tenantID string) ([]*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantMetricBaselineStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	res := make([]*tenmod.MetricBaseline, 0)
	if err := getAllOfTypeFromCouchAndFlatten(dbName, string(tenmod.TenantMetricBaselineType), tenmod.TenantMetricBaselineStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantMetricBaselineStr)
	return res, nil
}

// GetAllMetricBaselinesByPage - CouchDB implementation of GetAllMetricBaselinesByPage
func (tsd *TenantServiceDatastoreCouchDB) GetAllMetricBaselinesByPage(tenantID string, startKey string, limit int64) ([]*tenmod.MetricBaseline, *common.PaginationOffsets, error) {
	//logger.Log.Debugf("Fetching next %d %ss from startKey %s\n", limit, tenmod.TenantMetricBaselineStr, startKey)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	db, err := getDatabase(dbName)
	if err != nil {
		return nil, nil, err
	}

	res := make([]*tenmod.MetricBaseline, 0)

	// Need to retrieve 1 more than the asking size to be able to give back a startKey for the next page
	var batchSize int64
	if limit <= 0 || limit > int64(tsd.batchSize) {
		batchSize = tsd.batchSize
		logger.Log.Warningf("Provided limit %d is outside of range [1 - %d]. Using value %d in query", limit, batchSize, batchSize)
	} else {
		batchSize = limit
	}

	// Get 1 more object than the real response so that we can have the start key of the next page
	batchPlus1 := batchSize + 1

	params := generatePaginationQueryParams(startKey, batchPlus1, true, false)
	fetchResponse, err := getByDocIDWithQueryParams(metricBaselinesByNameIndex, tenmod.TenantMetricBaselineStr, db, &params)
	if err != nil {
		return nil, nil, err
	}

	if fetchResponse["rows"] == nil {
		return nil, nil, fmt.Errorf(ds.NotFoundStr)
	}

	castedRows := fetchResponse["rows"].([]interface{})
	if len(castedRows) == 0 {
		return nil, nil, fmt.Errorf(ds.NotFoundStr)
	}

	// Convert interface results to map results
	rows := []map[string]interface{}{}
	for _, obj := range castedRows {
		castedObj := obj.(map[string]interface{})
		genericDoc := castedObj["doc"].(map[string]interface{})
		rows = append(rows, genericDoc)
	}

	convertCouchDataArrayToFlattenedArray(rows, &res, tenmod.TenantMetricBaselineStr)

	nextPageStartKey := ""
	if int64(len(res)) == batchPlus1 {
		// Have an extra item, need to remove it and store the key for the next page
		nextPageStartKey = res[batchSize].MonitoredObjectID
		res = res[:batchSize]
	}

	paginationOffsets := common.PaginationOffsets{
		Self: startKey,
		Next: nextPageStartKey,
	}

	// Try to retrieve the previous page as well to get the previous start key
	prevPageParams := generatePaginationQueryParams(res[0].MonitoredObjectID, batchPlus1, true, true)
	prevPageResponse, err := getByDocIDWithQueryParams(metricBaselinesByNameIndex, tenmod.TenantMetricBaselineStr, db, &prevPageParams)
	if err == nil {
		// Try to get previous page details
		if prevPageResponse["rows"] != nil {
			prevPageRows := prevPageResponse["rows"].([]interface{})

			// There will always be 1 result at this point for the start of the current page, only add previous page key if there are actually records on the prev page
			if len(prevPageRows) > 1 {
				lastRow := prevPageRows[len(prevPageRows)-1].(map[string]interface{})
				paginationOffsets.Prev = lastRow["key"].(string)
			}
		}
	}

	// logger.Log.Debugf("Retrieved %d %ss from startKey %s\n", len(res), tenmod.TenantMetricBaselineStr, startKey)
	return res, &paginationOffsets, nil
}

func (tsd *TenantServiceDatastoreCouchDB) UpdateMetricBaselineForHourOfWeek(tenantID string, monObjID string, baselineData *tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Updating %s for %s %s for %s %s for %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, models.AsJSONString(baselineData.HourOfWeek))
	}

	existing, err := tsd.GetMetricBaseline(tenantID, monObjID)
	if err != nil {
		if !strings.Contains(err.Error(), ds.NotFoundStr) {
			// Error was something permanent, return it
			return nil, err
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			Baselines:         []*tenmod.MetricBaselineData{baselineData},
		}

		return tsd.CreateMetricBaseline(&createObj)
	}

	existing.MergeBaseline(baselineData)
	existing.ID = ds.GetDataIDFromFullID(existing.ID)

	updated, err := tsd.UpdateMetricBaseline(existing)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Updated %s for %s %s %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID)
	return updated, nil
}

// UpdateMetricBaselineForHourOfWeekWithCollection - couchDB implementation of UpdateMetricBaselineForHourOfWeekWithCollection
func (tsd *TenantServiceDatastoreCouchDB) UpdateMetricBaselineForHourOfWeekWithCollection(tenantID string, monObjID string, baselineDataCollection []*tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Updating %s for %s %s for %s %s for %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, models.AsJSONString(baselineDataCollection))
	}

	existing, err := tsd.GetMetricBaseline(tenantID, monObjID)
	if err != nil {
		if !strings.Contains(err.Error(), ds.NotFoundStr) {
			// Error was something permanent, return it
			return nil, err
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			Baselines:         baselineDataCollection,
		}

		return tsd.CreateMetricBaseline(&createObj)
	}

	existing.MergeBaselines(baselineDataCollection)
	existing.ID = ds.GetDataIDFromFullID(existing.ID)

	updated, err := tsd.UpdateMetricBaseline(existing)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Updated %s for %s %s %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID)
	return updated, nil
}

// func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaselineForMonitoredObject(tenantID string, monObjID string) (*tenmod.MetricBaseline, error) {
// 	logger.Log.Debugf("Retrieving %s for %s %s for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID)
// 	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

// 	// Retrieve just the subset of values.
// 	requestBody := map[string]interface{}{}
// 	requestBody["keys"] = []string{strings.ToLower(monObjID)}

// 	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))

// 	fetchResponse, err := fetchDesignDocumentResults(requestBody, dbName, metricBaselineByMOIDIndex)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows := fetchResponse["rows"].([]interface{})
// 	if rows == nil || len(rows) == 0 {
// 		return nil, fmt.Errorf(ds.NotFoundStr)
// 	}
// 	obj := rows[0].(map[string]interface{})
// 	value := obj["value"].(map[string]interface{})
// 	logger.Log.Debugf("Retrieved %s", models.AsJSONString(value))

// 	response := tenmod.MetricBaseline{}
// 	stripPrefixFromID(value)

// 	// Marshal the response from the datastore to bytes so that it
// 	// can be Marshalled back to the proper type.
// 	if err = convertGenericCouchDataToObject(value, &response, string(tenmod.TenantMetricBaselineType)); err != nil {
// 		return nil, err
// 	}

// 	if logger.IsDebugEnabled() {
// 		logger.Log.Debugf("Retrieved %s %s", tenmod.TenantMetricBaselineStr, models.AsJSONString(response))
// 	}

// 	return &response, nil
// }

func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaselineForMonitoredObjectForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32) ([]*tenmod.MetricBaselineData, error) {
	logger.Log.Debugf("Retrieving %ss for %s %s for %s %s for hour of week %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)

	existing, err := tsd.GetMetricBaseline(tenantID, monObjID)
	if err != nil {
		return nil, err
	}

	res := []*tenmod.MetricBaselineData{}
	for _, baseline := range existing.Baselines {
		if baseline.HourOfWeek == hourOfWeek {
			res = append(res, baseline)
		}
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %d %ss for %s %s for %s %s for hour of week %s", len(res), tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)
	}

	return res, nil
}

// GetMetricBaselinesForMOsIn - note that this function will return results that are not stored in the DB as new "empty" items so that they can be populated
// in a subsequent bulk PUT call.
func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaselinesForMOsIn(tenantID string, moIDList []string, addNotFoundValuesInResponse bool) ([]*tenmod.MetricBaseline, error) {
	startTime := time.Now()

	logger.Log.Debugf("Bulk retrieving %ss for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID)

	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	resource, err := couchdb.NewResource(dbName, nil)
	if err != nil {
		return nil, err
	}

	numMos := len(moIDList)
	if int64(numMos) > tsd.batchSize {
		return nil, fmt.Errorf("Too many Monitored Objects in bulk request. Limit is %d but request contains %d", tsd.batchSize, numMos)
	}

	mbTypeString := string(tenmod.TenantMetricBaselineType)

	// Build request to fetch existing records
	qp := &url.Values{}
	qp.Add("include_docs", "true")
	requestBody := map[string]interface{}{}
	updatedKeys := []string{}
	for _, moID := range moIDList {
		updatedKeys = append(updatedKeys, ds.PrependToDataID(moID, mbTypeString))
	}
	requestBody["keys"] = updatedKeys

	durationTillFetchBodyComplete := time.Since(startTime).Seconds()
	logger.Log.Warningf("DAO TIME UNTIL FETCH BODY READY: %f", durationTillFetchBodyComplete)

	fetchedGetData, err := fetchAllDocsWithQP(qp, requestBody, resource)
	if err != nil {
		return nil, err
	}

	durationTillFetchComplete := time.Since(startTime).Seconds()
	logger.Log.Warningf("DAO TIME UNTIL FETCH COMPLETE: %f", durationTillFetchComplete)

	logger.Log.Debugf("Fetch All response: %s", models.AsJSONString(fetchedGetData))

	// Need to convert retrieved records:
	rows, ok := fetchedGetData["rows"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to convert bulk fetch response: response from fetch has invalid format %s", models.AsJSONString(rows))
	}
	convertedRetrievedBaselines := []*tenmod.MetricBaseline{}
	for _, row := range rows {
		value, ok := row.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unable to convert single row to generic object %s", models.AsJSONString(value))
		}
		// Make sure that any records that are not found are added to the response as new empty records
		errorStr, okError := value["error"].(string)
		fetchedBaseline, okBaseline := value["doc"].(map[string]interface{})
		if okError && addNotFoundValuesInResponse {
			if errorStr == metricBaselineBulkFetchNotFoundStr {
				// This was an error but it was just that the record was not found, create a new value to be added to the result
				moID, ok := value["key"].(string)
				if !ok {
					return nil, fmt.Errorf("Can't create record for: %s", models.AsJSONString(row))
				}
				ts := ds.MakeTimestamp()
				strippedMOID := ds.GetDataIDFromFullID(moID)
				addObject := tenmod.MetricBaseline{
					ID:                    strippedMOID,
					Datatype:              mbTypeString,
					TenantID:              tenantID,
					MonitoredObjectID:     strippedMOID,
					Baselines:             []*tenmod.MetricBaselineData{},
					CreatedTimestamp:      ts,
					LastModifiedTimestamp: ts,
				}
				convertedRetrievedBaselines = append(convertedRetrievedBaselines, &addObject)
				continue
			}

			return nil, fmt.Errorf("Can't find record for: %s", models.AsJSONString(row))
		}

		// Handle the retrieved record
		if !okBaseline {
			return nil, fmt.Errorf("Unable to convert bulk fetch single value response to baseline: %s", err.Error())
		}
		stripPrefixFromID(fetchedBaseline)

		// Conver a fetched doc to the proper type for merging
		dataContainer := tenmod.MetricBaseline{}
		if err = convertGenericCouchDataToObject(fetchedBaseline, &dataContainer, tenmod.TenantMetricBaselineStr); err != nil {
			return nil, err
		}

		convertedRetrievedBaselines = append(convertedRetrievedBaselines, &dataContainer)
	}

	durationTillFetchMethodComplete := time.Since(startTime).Seconds()
	logger.Log.Warningf("DAO TIME UNTIL FETCH METHOD COMPLETE: %f", durationTillFetchMethodComplete)

	// Return the converted baseline data
	logger.Log.Debugf("Completed bulk retrieval of %ss for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID)
	return convertedRetrievedBaselines, nil
}

func (tsd *TenantServiceDatastoreCouchDB) BulkUpdateMetricBaselines(tenantID string, baselineUpdateList []*tenmod.MetricBaseline) ([]*common.BulkOperationResult, error) {
	startTime := time.Now()

	logger.Log.Debugf("Bulk updating %s", tenmod.TenantMetricBaselineStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	resource, err := couchdb.NewResource(dbName, nil)
	if err != nil {
		return nil, err
	}

	numMbs := len(baselineUpdateList)
	if int64(numMbs) > tsd.batchSize {
		return nil, fmt.Errorf("Too many Monitored Objects in bulk request. Limit is %d but request contains %d", tsd.batchSize, numMbs)
	}

	// Iterate over the collection and populate necessary fields
	data := make([]map[string]interface{}, 0)
	for _, mb := range baselineUpdateList {
		mb.ID = ds.PrependToDataID(mb.ID, string(tenmod.TenantMetricBaselineType))
		mb.Baselines = []*tenmod.MetricBaselineData{}
		genericMB, err := convertDataToCouchDbSupportedModel(mb)
		if err != nil {
			return nil, err
		}

		dataType := string(tenmod.TenantMetricBaselineType)
		dataProp := genericMB["data"].(map[string]interface{})
		dataProp["datatype"] = dataType
		dataProp["lastModifiedTimestamp"] = ds.MakeTimestamp()

		data = append(data, genericMB)
	}
	body := map[string]interface{}{
		"docs": data}

	durationTillUpdateBodyComplete := time.Since(startTime).Seconds()
	logger.Log.Warningf("DAO TIME UNTIL UPDATE BODY COMPLETE: %f", durationTillUpdateBodyComplete)

	fetchedData, err := performBulkUpdate(body, resource)
	if err != nil {
		return nil, err
	}

	durationTillUpdateCallComplete := time.Since(startTime).Seconds()
	logger.Log.Warningf("DAO TIME UNTIL UPDATE CALL COMPLETE: %f", durationTillUpdateCallComplete)

	// Populate the response
	res := make([]*common.BulkOperationResult, 0)
	for _, fetched := range fetchedData {
		newObj := common.BulkOperationResult{}
		if fetched["id"] != nil {
			newObj.ID = fetched["id"].(string)
		}
		if fetched["rev"] != nil {
			newObj.REV = fetched["rev"].(string)
		}
		if fetched["reason"] != nil {
			newObj.REASON = fetched["reason"].(string)
		}
		if fetched["error"] != nil {
			newObj.ERROR = fetched["error"].(string)
		}
		if fetched["ok"] != nil {
			newObj.OK = fetched["ok"].(bool)
		}

		newObj.ID = ds.GetDataIDFromFullID(newObj.ID)
		res = append(res, &newObj)
	}

	durationTillUpdateMethodComplete := time.Since(startTime).Seconds()
	logger.Log.Warningf("DAO TIME UNTIL UPDATE METHOD COMPLETE: %f", durationTillUpdateMethodComplete)

	logger.Log.Debugf("Bulk update of %s complete\n", tenmod.TenantMetricBaselineStr)
	return res, nil
}
