package druid

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/accedian/adh-gather/models"
	"github.com/stretchr/testify/assert"
)

func TestReformatReportSummary(t *testing.T) {
	responseStr := []byte(` [ {
		   "timestamp" : "2018-04-18T00:00:00.000Z",
		   "result" : {
		     "totalDuration" : 1.1616E8,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 1938,
		     "objectCount" : 3,
		     "totalViolationDuration" : 5.814E7,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.totalDuration" : 116160000,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationDuration" : 58140000,
		     "totalViolationCount" : 1938.0
		   }
		 } ]`)

	formattedResponse, err := reformatReportSummary(responseStr)
	assert.Nil(t, err)
	fmt.Printf("%v\n", models.AsJSONString(formattedResponse))
	assert.Equal(t, int64(58140000), formattedResponse.TotalViolationDuration)
	assert.Equal(t, int32(1938), formattedResponse.TotalViolationCount)
	assert.Equal(t, int32(3), formattedResponse.ObjectCount)

}

func TestReformatSLATimeSeries(t *testing.T) {
	responseStr := []byte(`[ {
		   "timestamp" : "2018-04-17T00:00:00.000Z",
		   "result" : {
		     "totalDuration" : 2700000.0,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 46,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationDuration" : 1380000,
		     "totalViolationCount" : 46.0,
		     "totalViolationDuration" : 1380000.0,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.totalDuration" : 2700000
		   }
		 }, {
		   "timestamp" : "2018-04-17T00:15:00.000Z",
		   "result" : {
		     "totalDuration" : 2700000.0,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 44,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationDuration" : 1320000,
		     "totalViolationCount" : 44.0,
		     "totalViolationDuration" : 1320000.0,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.totalDuration" : 2700000
		   }
		 }, {
		   "timestamp" : "2018-04-17T00:30:00.000Z",
		   "result" : {
		     "totalDuration" : 2700000.0,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 46,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.violationDuration" : 1380000,
		     "totalViolationCount" : 46.0,
		     "totalViolationDuration" : 1380000.0,
		     "accedian-twamp.twamp-pe.delayP95.sla.0.totalDuration" : 2700000
		   }
		 }]`)

	formattedResponse, err := reformatSLATimeSeries(responseStr)
	assert.Nil(t, err)
	fmt.Printf("%v\n", models.AsJSONString(formattedResponse))
	assert.Equal(t, 3, len(formattedResponse))
}

func TestFormatSLABucketResponse(t *testing.T) {
	responseStr1 := []byte(`[ {
		"timestamp" : "2018-04-17T00:00:00.000Z",
		"result" : [ {
		  "hourOfDay" : "03",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 180
		}, { 
		  "hourOfDay" : "04",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 180
		}, {
		  "hourOfDay" : "05",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 180
		}, {
		  "hourOfDay" : "06",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 180
		}, {
		  "hourOfDay" : "07",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 180
		}, {
		  "hourOfDay" : "08",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 180
		}, {
		  "hourOfDay" : "00",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 164
		}, {
		  "hourOfDay" : "02",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 142
		}, {
		  "hourOfDay" : "09",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 128
		}, {
		  "hourOfDay" : "01",
		  "accedian-twamp.twamp-pe.delayP95.sla.0.violationCount" : 124
		} ]		
	 } ] `)

	responseStr2 := []byte(`[ {
		"timestamp" : "2018-04-17T00:00:00.000Z",
		"result" : [ {
		  "hourOfDay" : "03",
		  "accedian-twamp.twamp-pe.jitterP95.sla.0.violationCount" : 22
		}, { 
		  "hourOfDay" : "04",
		  "accedian-twamp.twamp-pe.jitterP95.sla.0.violationCount" : 19
		}, {
		  "hourOfDay" : "07",
		  "accedian-twamp.twamp-pe.jitterP95.sla.0.violationCount" : 87
		}, {
		  "hourOfDay" : "00",
		  "accedian-twamp.twamp-pe.jitterP95.sla.0.violationCount" : 44
		}, {
		  "hourOfDay" : "02",
		  "accedian-twamp.twamp-pe.jitterP95.sla.0.violationCount" : 56
		}]		
	 } ] `)

	formattedResponse, err := reformatSLABucketResponse(responseStr1, nil)
	assert.Nil(t, err)
	fmt.Printf("%v\n", models.AsJSONString(formattedResponse))

	formattedResponse, err = reformatSLABucketResponse(responseStr2, formattedResponse)
	assert.Nil(t, err)
	fmt.Printf("%v\n", models.AsJSONString(formattedResponse))
}

func TestPostprocessResults(t *testing.T) {
	//throughputAvgCount:0 delayP95Sum:2.53765403e+08 throughputAvgSum:0
	//jitterP95Sum:3.229835e+06 jitterP95:728.096257889991 delayP95:57205.906898106405
	//throughputAvg:0 jitterP95Count:4436 delayP95Count:4436]

	resultStr := []byte(`
		[{
      "Timestamp": "2018-05-22T21:00:00.000Z",
      "Result": {
        "delayP95": 57205.906898106405,
        "jitterP95": 728.096257889991,
				"throughputAvg": 0,
				"throughputAvgCount": 0,
				"throughputAvgSum": 0
      }
    },
    {
      "Timestamp": "2018-05-22T22:00:00.000Z",
      "Result": {
        "delayP95": 50485.39688888889,
        "jitterP95": 514.4817777777778,
				"throughputAvg": 0,
				"throughputAvgCount": 0,
				"throughputAvgSum": 0
			      }
		}]`)

	response := make([]AggMetricsResponse, 0)
	err := json.Unmarshal(resultStr, &response)
	assert.Nil(t, err)

	pp := DropKeysPostprocessor{
		keysToDrop: []string{"throughputAvgCount", "throughputAvgSum"},
		countKeys:  map[string][]string{"throughputAvgCount": []string{"throughputAvg"}},
	}

	response = pp.Apply(response)
	_, ok := response[0].Result["throughputAvgCount"]
	assert.False(t, ok)
	_, ok = response[0].Result["throughputAvg"]
	assert.False(t, ok)
	_, ok = response[0].Result["throughputAvgSum"]
	assert.False(t, ok)
}
