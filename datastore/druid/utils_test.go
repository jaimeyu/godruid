package druid

import (
	"fmt"
	"testing"

	"github.com/accedian/adh-gather/models"
	"github.com/stretchr/testify/assert"
)

func TestReformatSLASummary(t *testing.T) {
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

	formattedResponse, err := reformatSLASummary(responseStr)
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
