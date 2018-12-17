package druid_test

import (
	"testing"

	"github.com/accedian/adh-gather/datastore/druid"
	"github.com/accedian/adh-gather/logger"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/godruid"
	"github.com/stretchr/testify/assert"
)

func TestFilterHelper(t *testing.T) {
	filter, err := druid.FilterHelper("foo", lowerEvent.EventAttrMap)
	if err != nil {
		logger.Log.Error("Filter helper error: ", err)
	}

	assert.Equal(t, &godruid.Filter{
		Type:        "bound",
		Dimension:   "foo",
		Ordering:    "numeric",
		Lower:       10000,
		LowerStrict: true,
	}, filter)

	filter, err = druid.FilterHelper("foo", upperEvent.EventAttrMap)
	if err != nil {
		logger.Log.Error("Filter helper error: ", err)
	}

	assert.Equal(t, &godruid.Filter{
		Type:        "bound",
		Dimension:   "foo",
		Ordering:    "numeric",
		Upper:       10000,
		UpperStrict: true,
	}, filter)

	filter, err = druid.FilterHelper("foo", bothEvent.EventAttrMap)
	if err != nil {
		logger.Log.Error("Filter helper error: ", err)
	}

	assert.Equal(t, &godruid.Filter{
		Type:        "bound",
		Dimension:   "foo",
		Ordering:    "numeric",
		Lower:       10000,
		LowerStrict: true,
		Upper:       10000,
		UpperStrict: true,
	}, filter)
}

func TestBuildMonitoredObjectFilterNoEntries(t *testing.T) {
	if druid.BuildMonitoredObjectFilter(nil) != nil {
		t.Errorf("Expecting nil filter to be returned")
	}
}

func TestBuildMonitoredObjectFilterOk(t *testing.T) {
	testMOs := []string{"mon1", "mon2"}

	rFilter := druid.BuildMonitoredObjectFilter(testMOs)

	if rFilter == nil || len(rFilter.Fields) != 0 {
		t.Errorf("Expecting only the monitored object filter to be returned but got %d", len(rFilter.Fields))
	}
}

func TestBaselineColumn(t *testing.T) {
	assert.Equal(t, "bl_delayAvg", druid.BaselineColumn("delayAvg"))
}

func TestThresholdCrossingTypeDefault(t *testing.T) {
	assert.Equal(t, tenmod.ThresholdStandard, druid.ThresholdCrossingType(map[string]string{}))
}

func TestThresholdCrossingTypeExplicitTypes(t *testing.T) {
	assert.Equal(t, tenmod.ThresholdStandard, druid.ThresholdCrossingType(map[string]string{"eventType": "standard"}))
	assert.Equal(t, tenmod.ThresholdPercentageBaseline, druid.ThresholdCrossingType(map[string]string{"eventType": "baseline_percentage"}))
	assert.Equal(t, tenmod.ThresholdStaticBaseline, druid.ThresholdCrossingType(map[string]string{"eventType": "baseline_static"}))
}

func TestIsBaselineTypes(t *testing.T) {
	assert.True(t, druid.IsBaselineType(tenmod.ThresholdStaticBaseline))
	assert.True(t, druid.IsBaselineType(tenmod.ThresholdPercentageBaseline))
	assert.False(t, druid.IsBaselineType(tenmod.ThresholdStandard))
}

var metric1 = "delayP95"

var lowerEvent = &tenmod.ThrPrfEventAttrMap{
	EventAttrMap: map[string]string{
		"lowerLimit":  "10000",
		"lowerStrict": "true",
	},
}

var upperEvent = &tenmod.ThrPrfEventAttrMap{
	EventAttrMap: map[string]string{
		"upperLimit":  "10000",
		"upperStrict": "true",
	},
}

var bothEvent = &tenmod.ThrPrfEventAttrMap{
	EventAttrMap: map[string]string{
		"upperLimit":  "10000",
		"upperStrict": "true",
		"lowerLimit":  "10000",
		"lowerStrict": "true",
	},
}

var aggTotalCount = godruid.AggCount("total")

var testThresholdCrossing1 = &godruid.QueryTimeseries{
	DataSource:  "druidTableName",
	Granularity: godruid.GranPeriod("PT1H", druid.TimeZoneUTC, ""),
	Context:     map[string]interface{}{"timeout": 60000},
	Aggregations: []godruid.Aggregation{
		*aggTotalCount,
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterUpperBound(metric1, "numeric", 30000, false),
				godruid.FilterSelector("objectType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.minor.0",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerUpperBound(metric1, "numeric", 30000, true, 50000, false),
				godruid.FilterSelector("objectType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.major.0",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerBound(metric1, "numeric", 50000, true),
				godruid.FilterSelector("objectType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.critical.0",
			},
		),
	},
	Intervals: []string{"1900-11-02/2100-01-01"},
}
