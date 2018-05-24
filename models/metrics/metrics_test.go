package metrics

import (
	"encoding/json"
	"testing"

	"github.com/accedian/adh-gather/logger"
	"github.com/stretchr/testify/assert"
)

var (
	defaultMetricsRequestBytes = []byte(`{
		 "tenantId": "8501f157-b7f5-41c3-aaba-c75e0566c54c",
		 "domainIds": [
		   "7c3d3280-628c-c778-92ec-4e9b83fcbb4d",
		   "41f1b537-b7f5-41c3-a1b2-a75e1536c54e"
		 ],
		 "interval": "2018-04-08T14:00:00/2018-04-09T15:00:00",
		 "granularity": "PT1H",
		 "timeout": 30000,
		 "aggregation": {
			 "name": "avg"
		 },
		 "metrics": [
		   {
		     "vendor": "accedian-twamp",
		     "objectType": "twamp-pe",
		     "name": "delayP95",
		     "direction": 0
		   }
		 ]
	}`)
)

func TestAggregateMetricsRequestSerialization(t *testing.T) {

	actual := AggregateMetricsAPIRequest{}

	expected := &AggregateMetricsAPIRequest{
		TenantID:  "8501f157-b7f5-41c3-aaba-c75e0566c54c",
		DomainIDs: []string{"7c3d3280-628c-c778-92ec-4e9b83fcbb4d", "41f1b537-b7f5-41c3-a1b2-a75e1536c54e"},
		Interval:  "2018-04-08T14:00:00/2018-04-09T15:00:00",
		Aggregation: AggregationSpec{
			Name: "avg",
		},
	}

	if err := json.Unmarshal(defaultMetricsRequestBytes, &actual); err != nil {
		logger.Log.Debugf("Unable to umarshal ingestion profile: %s", err.Error())
	}

	assert.Equal(t, expected.TenantID, actual.TenantID)
	assert.Equal(t, expected.DomainIDs, actual.DomainIDs)
	assert.Equal(t, expected.Interval, actual.Interval)
	assert.Equal(t, expected.Aggregation.Name, actual.Aggregation.Name)

}
