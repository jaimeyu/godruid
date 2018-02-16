package druid

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
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
	for _, e := range rawMetrics[0].Result.Events {
		monObj := e.Event["monitoredObjectId"].(string)
		if !res.ExistsP("result." + monObj) {
			_, err := res.ArrayP("result." + monObj)
			if err != nil {
				return nil, fmt.Errorf("Error formatting RawMetric JSON. Err: %s", err)
			}
		}
		obj := gabs.New()
		for k, v := range e.Event {
			if k != "monitoredObjectId" {
				obj.SetP(v, k)
			}
		}
		res.ArrayAppendP(obj.Data(), "result."+monObj)
	}

	dataContainer := map[string]interface{}{}
	if err := json.Unmarshal(res.Bytes(), &dataContainer); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Reformatted raw metrics data: %v", dataContainer)
	return dataContainer, nil
}

// Convert a query object to string, mainly for debugging purposes
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
