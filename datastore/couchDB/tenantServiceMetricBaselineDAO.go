package couchDB

import (
	"fmt"
	"strings"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

const (
	metricBaselinesByNameIndex = "_design/baselineIndex/_view/byName"
)

// CreateMetricBaseline - CouchDB implementation of CreateMetricBaseline
func (tsd *TenantServiceDatastoreCouchDB) CreateMetricBaseline(metricBaselineReq *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMetricBaselineStr, models.AsJSONString(metricBaselineReq))
	metricBaselineReq.ID = ds.GenerateID(metricBaselineReq, string(tenmod.TenantMetricBaselineType))
	tenantID := ds.PrependToDataID(metricBaselineReq.TenantID, string(admmod.TenantType))

	// Make sure there is no existing record for this id:
	existing, err := tsd.GetMetricBaselineForMonitoredObject(metricBaselineReq.TenantID, metricBaselineReq.MonitoredObjectID)
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

	existing, err := tsd.GetMetricBaselineForMonitoredObject(tenantID, monObjID)
	if err != nil {
		if !strings.Contains(err.Error(), ds.NotFoundStr) {
			// Error was something permanent, return it
			return nil, err
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			Baselines:         []*tenmod.MetricBaselineData{},
		}

		created, err := tsd.CreateMetricBaseline(&createObj)
		if err != nil {
			return nil, fmt.Errorf("Error trying to create %s as one did not exist to update: %s", tenmod.TenantMetricBaselineStr, err.Error())
		}

		// Set existing to be the returned value from the creation call
		existing = created
	}

	didUpdateBaseline := false
	for index, baseline := range existing.Baselines {
		if baseline.HourOfWeek == baselineData.HourOfWeek && baseline.Metric == baselineData.Metric && baseline.Direction == baselineData.Direction {
			// The baseline is already being tracked, update it
			existing.Baselines[index] = baselineData
			didUpdateBaseline = true
			break
		}
	}

	// If the baseline was not already being tracked, add it
	if !didUpdateBaseline {
		existing.Baselines = append(existing.Baselines, baselineData)
	}

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

	existing, err := tsd.GetMetricBaselineForMonitoredObject(tenantID, monObjID)
	if err != nil {
		if !strings.Contains(err.Error(), ds.NotFoundStr) {
			// Error was something permanent, return it
			return nil, err
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			Baselines:         []*tenmod.MetricBaselineData{},
		}

		created, err := tsd.CreateMetricBaseline(&createObj)
		if err != nil {
			return nil, fmt.Errorf("Error trying to create %s as one did not exist to update: %s", tenmod.TenantMetricBaselineStr, err.Error())
		}

		// Set existing to be the returned value from the creation call
		existing = created
	}

	for _, baselineFromCollection := range baselineDataCollection {
		didUpdateBaseline := false
		for index, baseline := range existing.Baselines {
			if baseline.HourOfWeek == baselineFromCollection.HourOfWeek && baseline.Metric == baselineFromCollection.Metric && baseline.Direction == baselineFromCollection.Direction {
				// The baseline is already being tracked, update it
				existing.Baselines[index] = baselineFromCollection
				didUpdateBaseline = true
				break
			}
		}

		// If the baseline was not already being tracked, add it
		if !didUpdateBaseline {
			existing.Baselines = append(existing.Baselines, baselineFromCollection)
		}
	}
	existing.ID = ds.GetDataIDFromFullID(existing.ID)

	updated, err := tsd.UpdateMetricBaseline(existing)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Updated %s for %s %s %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID)
	return updated, nil
}

func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaselineForMonitoredObject(tenantID string, monObjID string) (*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Retrieving %s for %s %s for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	db, err := getDatabase(dbName)
	if err != nil {
		return nil, err
	}

	resultList, err := getAllOfAny(string(tenmod.TenantMetricBaselineType), "monitoredObjectId", monObjID, tenmod.TenantMetricBaselineStr, db)
	if err != nil {
		return nil, err
	}

	if len(resultList) == 0 {
		return nil, fmt.Errorf(ds.NotFoundStr)
	}

	res := tenmod.MetricBaseline{}
	if err = convertGenericCouchDataToObject(resultList[0], &res, string(tenmod.TenantMetricBaselineType)); err != nil {
		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %s %s", tenmod.TenantMetricBaselineStr, models.AsJSONString(res))
	}

	res.ID = ds.GetDataIDFromFullID(res.ID)

	return &res, nil
}

func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaselineForMonitoredObjectForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32) ([]*tenmod.MetricBaselineData, error) {
	logger.Log.Debugf("Retrieving %ss for %s %s for %s %s for hour of week %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)

	existing, err := tsd.GetMetricBaselineForMonitoredObject(tenantID, monObjID)
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
