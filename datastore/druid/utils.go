package druid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/models/metrics"
	"github.com/accedian/godruid"
)

// Format a ThresholdCrossing object into something the UI can consume
func reformatThresholdCrossingResponse(thresholdCrossing []*pb.ThresholdCrossing) (map[string]interface{}, error) {
	res := gabs.New()
	_, err := res.Array("data")

	if err != nil {
		return nil, fmt.Errorf("Error formatting Threshold Crossing JSON. Err: %s", err)
	}
	for _, tc := range thresholdCrossing {
		obj := gabs.New()
		obj.SetP(tc.GetTimestamp(), "timestamp")
		for k, v := range tc.Result {
			obj.SetP(v, "result."+k)
		}
		res.ArrayAppend(obj.Data(), "data")
	}

	dataContainer := map[string]interface{}{}
	err = json.Unmarshal(res.Bytes(), &dataContainer)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("Reformatted threshold crossing data: %v", dataContainer)
	return dataContainer, nil
}

func reformatThresholdCrossingByMonitoredObjectResponse(thresholdCrossing []ThresholdCrossingByMonitoredObjectResponse) (map[string]interface{}, error) {
	res := gabs.New()
	for _, tc := range thresholdCrossing {
		monObjId := tc.Event["monitoredObjectId"]
		monObj := ""
		if monObjId != nil {
			monObj = monObjId.(string)
		}
		if !res.ExistsP("result." + monObj) {
			_, err := res.ArrayP("result." + monObj)
			if err != nil {
				return nil, fmt.Errorf("Error formatting Threshold Crossing By Monitored Object JSON. Err: %s", err)
			}
		}

		obj := gabs.New()
		obj.SetP(tc.Timestamp, "timestamp")
		for k, v := range tc.Event {
			obj.SetP(v, k)
		}
		res.ArrayAppendP(obj.Data(), "result."+monObj)

	}

	dataContainer := map[string]interface{}{}
	if err := json.Unmarshal(res.Bytes(), &dataContainer); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Reformatted threshold crossing by mon obj data: %v", dataContainer)
	return dataContainer, nil
}

func reformatRawMetricsResponse(rawMetrics []RawMetricsResponse) (map[string]interface{}, error) {
	res := gabs.New()
	var hasData bool
	for _, r := range rawMetrics {

		obj := gabs.New()
		var monObj string
		for k, v := range r.Result {

			parts := strings.Split(k, ".")
			monObj = parts[0]
			lastParts := parts[len(parts)-1]

			switch v.(type) {
			case float32:
				hasData = true
			case string:
				hasData = !strings.Contains(v.(string), "Infinity")
			default:
				hasData = true
			}
			if !strings.Contains(lastParts, "temporary") && hasData {
				obj.SetP(v, lastParts)
			}
		}

		if !res.ExistsP("result." + monObj) {
			_, err := res.ArrayP("result." + monObj)
			if err != nil {
				return nil, fmt.Errorf("Error formatting RawMetric JSON. Err: %s", err)
			}
		}

		if hasData {
			obj.SetP(r.Timestamp, "timestamp")
			res.ArrayAppendP(obj.Data(), "result."+monObj)
		}

	}

	dataContainer := map[string]interface{}{}
	if err := json.Unmarshal(res.Bytes(), &dataContainer); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Reformatted raw metrics data: %v", dataContainer)
	return dataContainer, nil
}

// convert a query object to string, mainly for debugging purposes
func queryToString(query godruid.Query, debug bool) string {
	var reqJson []byte
	var err error

	if debug {
		reqJson, err = json.MarshalIndent(query, "", "  ")
	} else {
		reqJson, err = json.Marshal(query)
	}

	if err != nil {
		return ""
	}

	return string(reqJson)
}

// Check to see if a value is in a slice
func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

type druidTimeSeriesEntry struct {
	Timestamp string
	Result    map[string]interface{}
}

func reformatSLASummary(druidResponse []byte) (*metrics.SLASummary, error) {
	logger.Log.Debugf("Response from druid for %s: %s", db.SLAReportStr, string(druidResponse))
	entries := []*druidTimeSeriesEntry{}
	if err := json.Unmarshal(druidResponse, &entries); err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return &metrics.SLASummary{}, nil
	}

	// For a summary, we expect only 1 entry in the druid results so just use the first entry.
	obj := gabs.New()
	for k, v := range entries[0].Result {
		if strings.Contains(k, ".sla.") {
			obj.SetP(v, "perMetricSummary."+k)
		} else {
			obj.SetP(v, k)
		}
	}

	summary := metrics.SLASummary{}
	if err := json.Unmarshal(obj.Bytes(), &summary); err != nil {
		return nil, err
	}
	if summary.TotalDuration > 0 {
		summary.SLACompliancePercent = (float32(summary.TotalDuration) - float32(summary.TotalViolationDuration)) * 100.0 / float32(summary.TotalDuration)
	}

	logger.Log.Debugf("Formatted result for %s: %v", db.SLAReportStr, models.AsJSONString(summary))
	return &summary, nil
}

func reformatSLATimeSeries(druidResponse []byte) ([]metrics.TimeSeriesEntry, error) {
	logger.Log.Debugf("Response from druid for %s: %s", db.SLAReportStr, string(druidResponse))
	entries := []*druidTimeSeriesEntry{}
	if err := json.Unmarshal(druidResponse, &entries); err != nil {
		return nil, err
	}

	res := make([]metrics.TimeSeriesEntry, len(entries))
	for i, tc := range entries {

		obj := gabs.New()
		for k, v := range tc.Result {
			if strings.Contains(k, ".sla.") {
				obj.SetP(v, "PerMetricResult."+k)
			} else {
				obj.SetP(v, k)
			}
		}
		timeseriesEntryResult := metrics.TimeSeriesResult{}
		if err := json.Unmarshal(obj.Bytes(), &timeseriesEntryResult); err != nil {
			return nil, err
		}

		res[i] = metrics.TimeSeriesEntry{
			Timestamp: tc.Timestamp,
			Result:    timeseriesEntryResult,
		}
	}

	logger.Log.Debugf("Formatted result for %s: %v", db.SLAReportStr, models.AsJSONString(res))
	return res, nil
}

type druidTopNEntry struct {
	Timestamp string
	Result    []map[string]interface{}
}

func reformatSLABucketResponse(druidResponse []byte, resultMap map[string]interface{}) (map[string]interface{}, error) {
	logger.Log.Debugf("Response from druid for %s: %s", db.SLAReportStr, string(druidResponse))
	entries := []*druidTopNEntry{}
	if err := json.Unmarshal(druidResponse, &entries); err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return nil, nil
	}

	// There should be max 1 entry in the response array
	formattedJSON, err := reformatBucketResponse(entries[0].Result, resultMap)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("Formatted result for %s: %v", db.SLAReportStr, models.AsJSONString(formattedJSON))
	return formattedJSON, nil
}

func reformatBucketResponse(buckets []map[string]interface{}, resultMap map[string]interface{}) (map[string]interface{}, error) {

	if resultMap == nil {
		resultMap = make(map[string]interface{}, len(buckets))
	}

	for _, result := range buckets {
		bucketValue := gabs.New()
		var bucketName string
		for k, v := range result {
			if _, ok := v.(string); ok {
				bucketName = v.(string)
			} else {
				bucketValue.SetP(v, k)
			}
		}

		if existingBucketValue, ok := resultMap[bucketName]; !ok {
			resultMap[bucketName] = bucketValue.Data()
		} else {
			merge, _ := gabs.Consume(existingBucketValue)
			merge.Merge(bucketValue)
			resultMap[bucketName] = merge.Data()
		}
	}

	return resultMap, nil
}

func buildLookupName(dimType, tenantID, dimValue string) string {
	return strings.ToLower(dimType + "|" + tenantID + "|" + dimValue)
}

func buildLookupNamePrefix(dimType, tenantID string) string {
	return strings.ToLower(dimType + "|" + tenantID)
}

func sendRequest(method string, httpClient *http.Client, endpoint, authToken string, req []byte) (result []byte, err error) {

	ep := endpoint + "?pretty"

	var reqBody io.Reader
	if req != nil {
		reqBody = bytes.NewBuffer(req)
	} else {
		reqBody = http.NoBody
	}
	request, err := http.NewRequest(method, ep, reqBody)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		request.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := httpClient.Do(request)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return
	}

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("%s: %s", resp.Status, string(result))
	}

	return
}
