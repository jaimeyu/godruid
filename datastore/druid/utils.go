package druid

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

func reformatThresholdCrossingResponse(thresholdCrossing []*pb.ThresholdCrossing) {
	res := gabs.New()
	res.ArrayOfSize(len(thresholdCrossing), "data")
	dataElements, err := res.S("data").Children()
	// TODO NEED TO FIX THIS
	if err != nil {
		fmt.Println("ERR-->", err)
	}
	for i, tc := range thresholdCrossing {
		_, err := dataElements[i].SetP(tc.GetTimestamp(), "timestamp")

		if err != nil {
			fmt.Println("ERR-->", err)
		}

		for k, v := range tc.Result {
			dataElements[i].SetP(v, "result."+k)
		}
		res.ArrayAppend(dataElements[i].Data(), "data")
	}

	fmt.Println("RES---->", res)
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
