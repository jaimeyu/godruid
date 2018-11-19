package couchDB

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

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
	dataID := generateMetricBaselineID(metricBaselineReq.MonitoredObjectID, metricBaselineReq.HourOfWeek)
	existing, err := tsd.GetMetricBaseline(metricBaselineReq.TenantID, dataID)
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
	if err := updateDataInCouchInBatchMode(dbName, metricBaselineReq, dataContainer, string(tenmod.TenantMetricBaselineType), tenmod.TenantMetricBaselineStr); err != nil {
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

func (tsd *TenantServiceDatastoreCouchDB) UpdateMetricBaselineForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32, baselineData *tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Updating %s for %s %s for %s %s for hour of week %d", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)
	}

	dataID := ds.GetDataIDFromFullID(generateMetricBaselineID(monObjID, hourOfWeek))
	existing, err := tsd.GetMetricBaseline(tenantID, dataID)
	if err != nil {
		if !strings.Contains(err.Error(), ds.NotFoundStr) {
			// Error was something permanent, return it
			return nil, err
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			HourOfWeek:        hourOfWeek,
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
func (tsd *TenantServiceDatastoreCouchDB) UpdateMetricBaselineForHourOfWeekWithCollection(tenantID string, monObjID string, hourOfWeek int32, baselineDataCollection []*tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Updating %s for %s %s for %s %s for hour of week %d with multiple values", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)
	}

	dataID := ds.GetDataIDFromFullID(generateMetricBaselineID(monObjID, hourOfWeek))
	existing, err := tsd.GetMetricBaseline(tenantID, dataID)
	if err != nil {
		if !strings.Contains(err.Error(), ds.NotFoundStr) {
			// Error was something permanent, return it
			return nil, err
		}

		// Error was that the Baseline does not exist for this Monitored Object, let's create it
		createObj := tenmod.MetricBaseline{
			MonitoredObjectID: monObjID,
			TenantID:          tenantID,
			HourOfWeek:        hourOfWeek,
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

func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaselineForMonitoredObjectForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32) ([]*tenmod.MetricBaselineData, error) {
	logger.Log.Debugf("Retrieving %ss for %s %s for %s %s for hour of week %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)

	dataID := ds.GetDataIDFromFullID(generateMetricBaselineID(monObjID, hourOfWeek))
	existing, err := tsd.GetMetricBaseline(tenantID, dataID)
	if err != nil {
		return nil, err
	}

	res := existing.Baselines

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %d %ss for %s %s for %s %s for hour of week %s", len(res), tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, monObjID, hourOfWeek)
	}

	return res, nil
}

// GetMetricBaselinesFor - note that this function will return results that are not stored in the DB as new "empty" items so that they can be populated
// in a subsequent bulk PUT call.
func (tsd *TenantServiceDatastoreCouchDB) GetMetricBaselinesFor(tenantID string, moIDToHourOfWeekMap map[string][]int32, addNotFoundValuesInResponse bool) ([]*tenmod.MetricBaseline, error) {
	logger.Log.Debugf("Bulk retrieving %ss for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID)

	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, metricBaselineDBSuffix))
	resource, err := couchdb.NewResource(dbName, nil)
	if err != nil {
		return nil, err
	}

	numEntries := 0
	keyMap := []string{}
	for moID, hourOfWeekList := range moIDToHourOfWeekMap {
		for _, hour := range hourOfWeekList {
			numEntries++
			keyMap = append(keyMap, generateMetricBaselineID(moID, hour))
		}
	}
	if int64(numEntries) > tsd.batchSize {
		return nil, fmt.Errorf("Too many Monitored Objects in bulk request. Limit is %d but request contains %d", tsd.batchSize, numEntries)
	}

	mbTypeString := string(tenmod.TenantMetricBaselineType)

	// Build request to fetch existing records
	qp := &url.Values{}
	qp.Add("include_docs", "true")
	requestBody := map[string]interface{}{}
	requestBody["keys"] = keyMap

	fetchedGetData, err := fetchAllDocsWithQP(qp, requestBody, resource)
	if err != nil {
		return nil, err
	}

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
				strippedMOIDParts := strings.Split(strippedMOID, "_")
				how, _ := strconv.ParseInt(strippedMOIDParts[1], 10, 32)
				addObject := tenmod.MetricBaseline{
					ID:                    strippedMOID,
					Datatype:              mbTypeString,
					TenantID:              tenantID,
					MonitoredObjectID:     strippedMOIDParts[0],
					HourOfWeek:            int32(how),
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

		// Convert a fetched doc to the proper type for merging
		dataContainer := tenmod.MetricBaseline{}
		if err = convertGenericCouchDataToObject(fetchedBaseline, &dataContainer, tenmod.TenantMetricBaselineStr); err != nil {
			return nil, err
		}

		convertedRetrievedBaselines = append(convertedRetrievedBaselines, &dataContainer)
	}

	// Return the converted baseline data
	logger.Log.Debugf("Completed bulk retrieval of %ss for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID)
	return convertedRetrievedBaselines, nil
}

func (tsd *TenantServiceDatastoreCouchDB) BulkUpdateMetricBaselines(tenantID string, baselineUpdateList []*tenmod.MetricBaseline) ([]*common.BulkOperationResult, error) {
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

	fetchedData, err := performBulkUpdateInBatchMode(body, resource)
	if err != nil {
		return nil, err
	}

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

	logger.Log.Debugf("Bulk update of %s complete\n", tenmod.TenantMetricBaselineStr)
	return res, nil
}

func generateMetricBaselineID(monObjID string, hourOfWeek int32) string {
	return fmt.Sprintf("%s_2_%s_%d", string(tenmod.TenantMetricBaselineType), monObjID, hourOfWeek)
}
