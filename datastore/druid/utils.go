package druid

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

func reformatThresholdCrossingResponse(thresholdCrossing []*pb.ThresholdCrossing) ([]byte, error) {
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

	return res.Bytes(), nil
}

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
