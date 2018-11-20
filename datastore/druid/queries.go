package druid

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/godruid"
	"github.com/satori/go.uuid"
)

type TimeBucket int

const (
	TimeZoneUTC     = "UTC"
	GranularityAll  = "all"
	GranularityNone = "none"
	HourOfDay       = 0
	DayOfWeek       = 1
)

// Constants for Threshold Crossing Baseline functionality
const (
	DruidBaselineIdPrefix = "bl_"
	DruidFilterIdSuffix   = "_baseline"
)

// Map used to reference the mathematical functions to calculate baselines for the different types of provisioned thresholds
var baselineFunctions = map[tenmod.ThresholdType]string{
	tenmod.ThresholdPercentageBaseline: "(%s/%s)-1",
	tenmod.ThresholdStaticBaseline:     "%s-%s",
}

var knownEventNames = []string{"critical", "major", "minor", "warn", "info"}

// FilterHelper - helper function to select correct druid filter based on
// a given event and metric
func FilterHelper(metric string, event map[string]string) (*godruid.Filter, error) {

	upperStrict, err := strconv.ParseBool(event["upperStrict"])
	if err != nil && event["upperStrict"] != "" {
		return nil, fmt.Errorf("Invalid value for 'upperStrict' : %v. Must be a boolean", upperStrict)
	}

	lowerStrict, err := strconv.ParseBool(event["lowerStrict"])
	if err != nil && event["lowerStrict"] != "" {
		return nil, fmt.Errorf("Invalid value for 'lowerStrict' : %v. Must be a boolean", lowerStrict)
	}

	lowerLimit, err := strconv.ParseFloat(event["lowerLimit"], 32)
	if err != nil && event["lowerLimit"] != "" {
		return nil, fmt.Errorf("Invalid value for 'lowerLimit' : %v. Must be a number", lowerLimit)
	}

	upperLimit, err := strconv.ParseFloat(event["upperLimit"], 32)
	if err != nil && event["upperLimit"] != "" {
		return nil, fmt.Errorf("Invalid value for 'upperLimit' : %v. Must be a number", upperLimit)
	}

	return FilterLimitSelectorHelper(metric, lowerLimit, lowerStrict, upperLimit, upperStrict), nil
}

// Retrieves the appropriate druid bounded filter based on the defined lower and upper bounds for a given metric.
// NOTE: This should not be used for metrics that need to take negative numbers into account since lower/upper
// bound values of 0 are considered to be infinite lower or infinite upper respectively
// Arguments:
//   metric - the name of the metric we want to apply the filter to
//   lowerLimit - the lower limit of the bounded filter that we want to build. A value of 0 assumes no lower bound
//   lowerStrict - a value of false assumes the comparison <=. A value of true assumes the comparison <
//   upperLimit - the upper limit of the bounded filter that we want to build. A value of 0 assumes no upper bound
//   upperStrict - a value of false assumes the comparison >=. A value of true assumes the comparison >
func FilterLimitSelectorHelper(metric string, lowerLimit float64, lowerStrict bool, upperLimit float64, upperStrict bool) *godruid.Filter {
	// Builds a filter that behaves as lowerLimit <[=] val <[=] upperLimit
	if upperLimit != 0 && lowerLimit != 0 {
		return godruid.FilterLowerUpperBound(metric, godruid.NUMERIC, float32(lowerLimit), lowerStrict, float32(upperLimit), upperStrict)
	}

	// Builds a filter that behaves as val <[=] upperLimit
	if upperLimit != 0 {
		return godruid.FilterUpperBound(metric, godruid.NUMERIC, float32(upperLimit), upperStrict)
	}

	// Builds a filter that behaves as lowerLimit <[=] val
	if lowerLimit != 0 {
		return godruid.FilterLowerBound(metric, godruid.NUMERIC, float32(lowerLimit), lowerStrict)
	}

	return nil
}

// Special function necessary to properly handle when either the lower or upper bound value is 0. This case needs to be handled specially since
// the json api marshaller nulls out values of 0 since they are the default and druid cannot handle having a null filter passed to it
func correctedBoundaryValues(boundarySpec *metrics.MetricBucketBoundarySpec, isUpper bool) *metrics.MetricBucketBoundarySpec {
	if boundarySpec != nil {
		// If the boundary spec is not a value of 0 then just return the original one
		if boundarySpec.Value == 0 {
			if isUpper {
				// If we are the upper bound and strictness is set to true then we want the upperbound to be the smallest possible negative number to 0 without including 0
				// ie: values <= -1.401298464324817070923729583289916131280e-45
				if boundarySpec.Strict {
					return &metrics.MetricBucketBoundarySpec{Value: -math.SmallestNonzeroFloat32, Strict: false}
				}
				// Otherwise we want the value of 0 to be included in our filter but we cannot include the value of 0 in the filter spec so use the smallest
				// possible positive number and use non-strict to make sure that it isn't included but 0 is
				// ie: values < 1.401298464324817070923729583289916131280e-45
				return &metrics.MetricBucketBoundarySpec{Value: math.SmallestNonzeroFloat32, Strict: true}
			}

			// If we are the lower bound and strictness is set to true then we want the lowerbound to be the smallest possible positive number to 0 without including 0
			// ie: 1.401298464324817070923729583289916131280e-45 <= values
			if boundarySpec.Strict {
				return &metrics.MetricBucketBoundarySpec{Value: math.SmallestNonzeroFloat32, Strict: false}
			}
			// Otherwise we want the value of 0 to be included in our filter but we cannot include the value of 0 in the filter spec so use the smallest
			// possible negative number and use non-strict to make sure that it isn't included but 0 is
			// ie: -1.401298464324817070923729583289916131280e-45 < values
			return &metrics.MetricBucketBoundarySpec{Value: -math.SmallestNonzeroFloat32, Strict: true}
		}
	}
	return boundarySpec
}

// Retrieves the appropriate druid bounded filter based on the defined lower and upper bounds for a given metric.
// Arguments:
//   metric - the name of the metric we want to apply the filter to
//   lowerSpec - the specification for the lower limit of the bounded filter that we want to build. A nil value assumes no lower bound.
//				 A value of true for strict indicates that we are using < for comparison as opposed to <=
//   upperSpec - the specification for the upper limit of the bounded filter that we want to build. A nil value assumes no upper bound.
//				 A value of true for strict indicates that we are using > for comparison as opposed to >=
func BoundarySpecFilterLimitSelectorHelper(metric string, lowerSpec *metrics.MetricBucketBoundarySpec, upperSpec *metrics.MetricBucketBoundarySpec) *godruid.Filter {

	correctedUpper := correctedBoundaryValues(upperSpec, true)
	correctedLower := correctedBoundaryValues(lowerSpec, false)
	// Builds a filter that behaves as lower value <[=] val <[=] upper value where strict=true indicates that "or equal" should be applicable to each boundary
	if upperSpec != nil && lowerSpec != nil {
		// This is to handle a special case. The godruid client treats 0 as an absence of a value but it could be a legitimate scenario where we want to return a bucket of just 0 value metrics
		if upperSpec.Value == 0 && lowerSpec.Value == 0 {
			return godruid.FilterSelector(metric, 0)
		}

		return godruid.FilterLowerUpperBound(metric, godruid.NUMERIC, correctedLower.Value, correctedLower.Strict, correctedUpper.Value, correctedUpper.Strict)
	}

	// Builds a filter that behaves as val <[=] upper value where strict=true indicates that "or equal" should be applicable to the upper limit boundary
	if upperSpec != nil {
		return godruid.FilterUpperBound(metric, godruid.NUMERIC, correctedUpper.Value, correctedUpper.Strict)
	}

	// Builds a filter that behaves as lower value <[=] val where strict=true indicates that "or equal" should be applicable to the lower limit boundary
	if lowerSpec != nil {
		return godruid.FilterLowerBound(metric, godruid.NUMERIC, correctedLower.Value, correctedLower.Strict)
	}

	return nil
}

func HistogramQuery(tenant string, metaMOs []string, dataSource string, interval string, granularity string, timeout int32, metrics []metrics.MetricBucketRequest) (*godruid.QueryTimeseries, *db.QueryKeySpec, error) {

	var aggregations []godruid.Aggregation

	querySpec := db.QueryKeySpec{}

	for _, met := range metrics {

		querySpecID := querySpec.AddKeySpec(map[string]interface{}{"vendor": met.Vendor, "objectType": met.ObjectType, "metric": met.Metric, "direction": met.Direction})

		for i, bucket := range met.Buckets {
			metIndex := i

			name := querySpecID + db.QueryDelimeter + strconv.Itoa(metIndex)

			filter := BoundarySpecFilterLimitSelectorHelper(met.Metric, bucket.Lower, bucket.Upper)

			aggregation := godruid.AggFiltered(
				godruid.FilterAnd(
					filter,
					buildInFilter("objectType", met.ObjectType),
					buildInFilter("direction", met.Direction),
				),
				&godruid.Aggregation{
					Type: "count",
					Name: name,
				},
			)
			aggregations = append(aggregations, aggregation)
		}
	}

	return &godruid.QueryTimeseries{
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
		),
		Aggregations: aggregations,
		Intervals:    []string{interval}}, &querySpec, nil
}

func HistogramQueryV1(tenant string, metaMOs []string, dataSource string, interval string, granularity string, timeout int32, metrics []map[string]interface{}) (*godruid.QueryTimeseries, error) {

	var aggregations []godruid.Aggregation

	for _, met := range metrics {

		metName := met["metric"].(string)
		metVendor := met["vendor"].(string)
		metDirection := met["direction"].(string)
		metObjectType := met["objectType"].(string)

		for i, bucket := range met["buckets"].([]interface{}) {
			bucketMap := bucket.(map[string]interface{})
			metUpper := bucketMap["upper"].(float64)
			metLower := bucketMap["lower"].(float64)
			metIndex := i

			name := fmt.Sprintf("%s.%s.%s.%s.%d", metVendor, metObjectType, metName, metDirection, metIndex)

			filter := FilterLimitSelectorHelper(metName, metLower, false, metUpper, true)

			aggregation := godruid.AggFiltered(
				godruid.FilterAnd(
					filter,
					godruid.FilterSelector("objectType", metObjectType),
					godruid.FilterSelector("direction", metDirection),
				),
				&godruid.Aggregation{
					Type: "count",
					Name: name,
				},
			)
			aggregations = append(aggregations, aggregation)
		}
	}

	return &godruid.QueryTimeseries{
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
		),
		Aggregations: aggregations,
		Intervals:    []string{interval}}, nil
}

// Returns the associated baseline column based on the provided metric name
func BaselineColumn(metric string) string {
	// May need to put a check here to see if it is actually a metric that even has a baseline
	return DruidBaselineIdPrefix + metric
}

// Maps the type of threshold configured for a specific event severity in the event attributes map. Defaults to the standard type
func ThresholdCrossingType(eventAttrs map[string]string) tenmod.ThresholdType {
	if tType, ok := eventAttrs["eventType"]; ok {
		switch tType {
		case string(tenmod.ThresholdPercentageBaseline):
			return tenmod.ThresholdPercentageBaseline
		case string(tenmod.ThresholdStaticBaseline):
			return tenmod.ThresholdStaticBaseline
		case string(tenmod.ThresholdStandard):
			return tenmod.ThresholdStandard
		}
	}

	return tenmod.ThresholdStandard
}

// Determines whether the provided threshold type is baseline-based or a standard threshold profile
func IsBaselineType(thresholdType tenmod.ThresholdType) bool {
	for ttype, _ := range baselineFunctions {
		if ttype == thresholdType {
			return true
		}
	}
	return false
}

// Creates a threshold crossing filter based on the type of threshold passed in (baseline-based or standard)
// Arguments:
//   fieldNamePrefix - the identifying prefix needed for baseline-based thresholds to uniquely name the virtual column
//   metric - the name of the metric that we are building the filter for
//   eventMap - the provisioned attribute value pairs for the threshold we are processing
//	 baselineColumns - the set of virtual columns that have been built up to ensure that there are no duplicate columns
// Returns:
//	 *godruid.Filter - a filter that represents the provisioned threshold
//	 *godruid.VirtualColumn - a virtual column that represents the baseline calculation for the metric. This is only applicable if the threshold is baseline-based
//	 error - if any issue is encountered during processing
func BuildThresholdCrossingFilter(fieldNamePrefix string, metric string, eventMap *tenmod.ThrPrfEventAttrMap, baselineColumns *[]godruid.VirtualColumn) (*godruid.Filter, error) {
	thresholdCrossingType := ThresholdCrossingType(eventMap.EventAttrMap)
	filterDim := metric

	var divByZeroFilter *godruid.Filter

	// Add a baseline function in a virtual column for threshold crossings that are baseline based
	if IsBaselineType(thresholdCrossingType) {
		bcol := BaselineColumn(metric)
		filterDim = fieldNamePrefix + DruidFilterIdSuffix

		// Ensure that this virtual column has not already been built up from the provided set since we cannot have dups
		found := false
		for _, bc := range *baselineColumns {
			if bc.Name == filterDim {
				found = true
				break
			}
		}
		if !found {
			*baselineColumns = append(*baselineColumns, godruid.NewVirtualColumn(filterDim, fmt.Sprintf(baselineFunctions[thresholdCrossingType], metric, bcol), godruid.VirtualColumnDouble))
		}

		if thresholdCrossingType == tenmod.ThresholdPercentageBaseline {
			divByZeroFilter = godruid.FilterNot(godruid.FilterSelector(bcol, "0"))
		}
	}

	boundedFilter, err := FilterHelper(filterDim, eventMap.EventAttrMap)
	if err != nil {
		return nil, err
	}
	thresholdFilter := boundedFilter

	if divByZeroFilter != nil {
		// To guard against division by zero.
		// The division by zero filter MUST BE FIRST
		thresholdFilter = godruid.FilterAnd(
			divByZeroFilter,
			boundedFilter,
		)
	}

	return thresholdFilter, nil
}

func ThresholdViolationsQuery(tenant string, dataSource string, metaMOs []string, granularity string, interval string, metricWhitelist []metrics.MetricIdentifierFilter, thresholdProfile *tenmod.ThresholdProfile, timeout int32) (*godruid.QueryTimeseries, error) {

	// columns used for baseline calculations
	var baselineCalculationCols []godruid.VirtualColumn
	// all of the aggregations
	var aggregations []godruid.Aggregation

	type objectTypeDirectionFilters struct {
		BaseFilter              *godruid.Filter
		ThresholdFiltersByEvent map[string][]*godruid.Filter
		ThresholdFilterList     []*godruid.Filter
	}

	sortedVendorMapKeys := getSortedKeySlice(reflect.ValueOf(thresholdProfile.Thresholds.VendorMap).MapKeys())
	for _, vk := range sortedVendorMapKeys {
		v := thresholdProfile.Thresholds.VendorMap[vk]

		sortedMOTypeMapKeys := getSortedKeySlice(reflect.ValueOf(v.MonitoredObjectTypeMap).MapKeys())
		for _, tk := range sortedMOTypeMapKeys {
			t := v.MonitoredObjectTypeMap[tk]
			// This is for de-duping violation duration for metrics that are violated at the same time for the same object.
			perDirectionFilters := make(map[string]*objectTypeDirectionFilters)

			sortedMetricTypeMapKeys := getSortedKeySlice(reflect.ValueOf(t.MetricMap).MapKeys())
			for _, mk := range sortedMetricTypeMapKeys {
				m := t.MetricMap[mk]
				sortedDirectionTypeMapKeys := getSortedKeySlice(reflect.ValueOf(m.DirectionMap).MapKeys())
				for _, dk := range sortedDirectionTypeMapKeys {
					d := m.DirectionMap[dk]

					// skip metrics that are not on the whitelist (if one was provided)
					if !inWhitelist(metricWhitelist, vk, tk, mk, dk) {
						continue
					}

					aggNamePrefix := buildMetricAggPrefix(vk, tk, mk, dk)

					// create a base filter for this objectType and direction (druid doesn't store vendor)
					objectTypeAndDirectionFilter := godruid.FilterAnd(
						godruid.FilterSelector("objectType", tk),
						godruid.FilterSelector("direction", dk),
					)

					// process the provisioned events (severities) and create aggregations
					sortedEventTypeMapKeys := getSortedKeySlice(reflect.ValueOf(d.EventMap).MapKeys())
					for _, ek := range sortedEventTypeMapKeys {
						e := d.EventMap[ek]

						thresholdFilter, err := BuildThresholdCrossingFilter(aggNamePrefix, mk, e, &baselineCalculationCols)
						if err != nil {
							return nil, err
						}

						// store the threshold filter in a map - needed for de-duping later
						dirFilters, ok := perDirectionFilters[vk+"."+tk+"."+dk]
						if !ok {
							perDirectionFilters[vk+"."+tk+"."+dk] = &objectTypeDirectionFilters{
								BaseFilter:              objectTypeAndDirectionFilter,
								ThresholdFilterList:     []*godruid.Filter{thresholdFilter},
								ThresholdFiltersByEvent: map[string][]*godruid.Filter{ek: []*godruid.Filter{thresholdFilter}},
							}
						} else {
							dirFilters.ThresholdFilterList = append(dirFilters.ThresholdFilterList, thresholdFilter)
							filterList, ok := dirFilters.ThresholdFiltersByEvent[ek]
							if !ok {
								dirFilters.ThresholdFiltersByEvent[ek] = []*godruid.Filter{thresholdFilter}
							} else {
								dirFilters.ThresholdFiltersByEvent[ek] = append(filterList, thresholdFilter)
							}
						}

						aggNameEventPrefix := aggNamePrefix + "." + ek

						// Count number of times the metric was violated
						violationCountAggName := aggNameEventPrefix + ".violationCount"
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type: "count",
								Name: violationCountAggName,
							},
						))

						// Sum the duration while this metric was in violation.
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      aggNameEventPrefix + ".violationDuration",
							},
						))

					}
				}
			}

			// Duration de-dupping aggregations are created here.
			if len(perDirectionFilters) > 0 {

				for k, v := range perDirectionFilters {

					if len(v.ThresholdFilterList) > 0 {
						// Sum the violation duration per vendor/objecttype/direction
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								v.BaseFilter,
								godruid.FilterOr(v.ThresholdFilterList...),
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      k + ".violationDuration",
							},
						))
					}

					if len(v.ThresholdFiltersByEvent) > 0 {
						// Sum the violation duration per vendor/objecttype/direction/event

						// This is done in the fixed order of knownEventNames for additional de-dup.
						// It is possible to have 2 metrics in the same record that are both violated but
						// for different events.  If 1 metric violates critical and the other violates
						// minor we want the duration violation counted against critical only.
						processed := []*godruid.Filter{}
						for _, eventName := range knownEventNames {
							tf, ok := v.ThresholdFiltersByEvent[eventName]
							if !ok {
								continue
							}

							// Here's where we build a filter containing the threshold conditions
							// In order to de-dup (as describe above), we must exclude any filters
							// we previous processed.
							var filter *godruid.Filter
							if len(processed) == 0 {
								// i.e. here we'd be building the critical filter
								filter = godruid.FilterAnd(
									v.BaseFilter,
									godruid.FilterOr(tf...),
								)
							} else {
								// here we'd build, for example, the major filter and exclude critical filters
								filter = godruid.FilterAnd(
									v.BaseFilter,
									godruid.FilterOr(tf...),
									godruid.FilterNot(godruid.FilterOr(processed...)),
								)
							}
							aggregations = append(aggregations, godruid.AggFiltered(
								filter,
								&godruid.Aggregation{
									Type:      "longSum",
									FieldName: "duration",
									Name:      "__event." + eventName + "." + k + ".violationDuration",
								},
							))

							// add the threshold filters to the processed list so we can exclude them for the next event
							processed = append(processed, tf...)

						}

					}

				}

			}
		}
	}

	return &godruid.QueryTimeseries{
		QueryType:   godruid.TIMESERIES,
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
		),
		Aggregations:   aggregations,
		VirtualColumns: baselineCalculationCols,
		Intervals:      []string{interval}}, nil
}

func ThresholdViolationsQueryV1(tenant string, dataSource string, metaMOs []string, granularity string, interval string, metricWhitelist []metrics.MetricIdentifierV1, thresholdProfile *pb.TenantThresholdProfileData, timeout int32) (*godruid.QueryTimeseries, error) {
	// all of the aggregations
	var aggregations []godruid.Aggregation
	// all of the post aggregations
	var postAggregations []godruid.PostAggregation
	// the names of the aggregations that are computing violation counts
	var violationCountAggs []string
	// the names of the aggregations that are computing violation counts, grouped by eventName
	violationCountAggsByEvent := map[string][]string{}
	// the names of the aggregations that are computing de-duped duration sum
	var durationAggs []string
	// the names of the aggregations that are computing de-duped violation duration sum
	var violationDurationAggs []string
	// the names of the aggregations that are computing de-duped violation duration sum, grouped by eventName
	violationDurationAggsByEvent := map[string][]string{}

	type objectTypeDirectionFilters struct {
		BaseFilter              *godruid.Filter
		ThresholdFiltersByEvent map[string][]*godruid.Filter
		ThresholdFilterList     []*godruid.Filter
	}

	sortedVendorMapKeys := getSortedKeySlice(reflect.ValueOf(thresholdProfile.GetThresholds().GetVendorMap()).MapKeys())
	for _, vendorKey := range sortedVendorMapKeys {
		vk := vendorKey
		v := thresholdProfile.GetThresholds().GetVendorMap()[vk]
		// for vk, v := range thresholdProfile.GetThresholds().GetVendorMap() {
		sortedMOTypeMapKeys := getSortedKeySlice(reflect.ValueOf(v.GetMonitoredObjectTypeMap()).MapKeys())
		for _, typeKey := range sortedMOTypeMapKeys {
			tk := typeKey
			t := v.GetMonitoredObjectTypeMap()[tk]
			// for tk, t := range v.GetMonitoredObjectTypeMap() {
			// This is for de-duping violation duration for metrics that are violated at the same time for the same object.
			perDirectionFilters := make(map[string]*objectTypeDirectionFilters)

			sortedMetricTypeMapKeys := getSortedKeySlice(reflect.ValueOf(t.GetMetricMap()).MapKeys())
			for _, metricTypeKey := range sortedMetricTypeMapKeys {
				mk := metricTypeKey
				m := t.GetMetricMap()[mk]
				// for mk, m := range t.GetMetricMap() {
				sortedDirectionTypeMapKeys := getSortedKeySlice(reflect.ValueOf(m.GetDirectionMap()).MapKeys())
				for _, directionTypeKey := range sortedDirectionTypeMapKeys {
					dk := directionTypeKey
					d := m.GetDirectionMap()[dk]
					// for dk, d := range m.GetDirectionMap() {

					// skip metrics that are not on the whitelist (if one was provided)
					if !inWhitelistV1(metricWhitelist, vk, tk, mk, dk) {
						continue
					}

					// create a base filter for this objectType and direction (druid doesn't store vendor)
					objectTypeAndDirectionFilter := godruid.FilterAnd(
						godruid.FilterSelector("objectType", tk),
						godruid.FilterSelector("direction", dk),
					)

					aggNamePrefix := buildMetricAggPrefix(vk, tk, mk, dk)
					// create an aggregation to sum the total duration this metric was measured.
					aggregations = append(aggregations, godruid.AggFiltered(
						objectTypeAndDirectionFilter,
						&godruid.Aggregation{
							Type:      "longSum",
							FieldName: "duration",
							Name:      aggNamePrefix + ".totalDuration",
						},
					))

					// process the provisioned events (severities) and create aggregations
					sortedEventTypeMapKeys := getSortedKeySlice(reflect.ValueOf(d.GetEventMap()).MapKeys())
					for _, eventTypeKey := range sortedEventTypeMapKeys {
						ek := eventTypeKey
						e := d.GetEventMap()[ek]
						// for ek, e := range d.GetEventMap() {

						thresholdFilter, err := FilterHelper(mk, e.EventAttrMap)
						if err != nil {
							return nil, err
						}

						// store the threshold filter in a map - needed for de-duping later
						dirFilters, ok := perDirectionFilters[vk+"."+tk+"."+dk]
						if !ok {
							perDirectionFilters[vk+"."+tk+"."+dk] = &objectTypeDirectionFilters{
								BaseFilter:              objectTypeAndDirectionFilter,
								ThresholdFilterList:     []*godruid.Filter{thresholdFilter},
								ThresholdFiltersByEvent: map[string][]*godruid.Filter{ek: []*godruid.Filter{thresholdFilter}},
							}
						} else {
							dirFilters.ThresholdFilterList = append(dirFilters.ThresholdFilterList, thresholdFilter)
							filterList, ok := dirFilters.ThresholdFiltersByEvent[ek]
							if !ok {
								dirFilters.ThresholdFiltersByEvent[ek] = []*godruid.Filter{thresholdFilter}
							} else {
								dirFilters.ThresholdFiltersByEvent[ek] = append(filterList, thresholdFilter)
							}
						}

						aggNameEventPrefix := aggNamePrefix + "." + ek

						// Count number of times the metric was violated
						violationCountAggName := aggNameEventPrefix + ".violationCount"
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type: "count",
								Name: violationCountAggName,
							},
						))
						violationCountAggs = append(violationCountAggs, violationCountAggName)

						aggs, ok := violationCountAggsByEvent[ek]
						if !ok {
							violationCountAggsByEvent[ek] = []string{violationCountAggName}
						} else {
							violationCountAggsByEvent[ek] = append(aggs, violationCountAggName)
						}

						// Sum the duration while this metric was in violation.
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      aggNameEventPrefix + ".violationDuration",
							},
						))

					}
				}
			}

			// Duration de-dupping aggregations are created here.
			if len(perDirectionFilters) > 0 {

				for k, v := range perDirectionFilters {

					// An aggregation to sum the duration for a vendor/objectType/direction
					aggregations = append(aggregations, godruid.AggFiltered(
						v.BaseFilter,
						&godruid.Aggregation{
							Type:      "longSum",
							FieldName: "duration",
							Name:      k + ".totalDuration",
						},
					))
					durationAggs = append(durationAggs, k+".totalDuration")

					if len(v.ThresholdFilterList) > 0 {
						// Sum the violation duration per vendor/objecttype/direction
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								v.BaseFilter,
								godruid.FilterOr(v.ThresholdFilterList...),
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      k + ".violationDuration",
							},
						))
						violationDurationAggs = append(violationDurationAggs, k+".violationDuration")
					}

					if len(v.ThresholdFiltersByEvent) > 0 {
						// Sum the violation duration per vendor/objecttype/direction/event

						// This is done in the fixed order of knownEventNames for additional de-dup.
						// It is possible to have 2 metrics in the same record that are both violated but
						// for different events.  If 1 metric violates critical and the other violates
						// minor we want the duration violation counted against critical only.
						processed := []*godruid.Filter{}
						for _, eventName := range knownEventNames {
							tf, ok := v.ThresholdFiltersByEvent[eventName]
							if !ok {
								continue
							}

							// Here's where we build a filter containing the threshold conditions
							// In order to de-dup (as describe above), we must exclude any filters
							// we previous processed.
							var filter *godruid.Filter
							if len(processed) == 0 {
								// i.e. here we'd be building the critical filter
								filter = godruid.FilterAnd(
									v.BaseFilter,
									godruid.FilterOr(tf...),
								)
							} else {
								// here we'd build, for example, the major filter and exclude critical filters
								filter = godruid.FilterAnd(
									v.BaseFilter,
									godruid.FilterOr(tf...),
									godruid.FilterNot(godruid.FilterOr(processed...)),
								)
							}
							aggregations = append(aggregations, godruid.AggFiltered(
								filter,
								&godruid.Aggregation{
									Type:      "longSum",
									FieldName: "duration",
									Name:      "__event." + eventName + "." + k + ".violationDuration",
								},
							))

							aggs, ok := violationDurationAggsByEvent[eventName]
							if !ok {
								violationDurationAggsByEvent[eventName] = []string{"__event." + eventName + "." + k + ".violationDuration"}
							} else {
								violationDurationAggsByEvent[eventName] = append(aggs, "__event."+eventName+"."+k+".violationDuration")
							}
							// add the threshold filters to the processed list so we can exclude them for the next event
							processed = append(processed, tf...)

						}

					}

				}

			}
		}
	}

	if len(violationCountAggs) > 0 {
		// Sum the violation count per metric to get an overall total.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalViolationCount",
			"+",
			buildPostAggregationFields(violationCountAggs)))
	}
	if len(violationCountAggsByEvent) > 0 {
		for ek, v := range violationCountAggsByEvent {
			postAggregations = append(postAggregations, godruid.PostAggArithmetic(
				buildTopLevelEventAgg(ek, "totalViolationCount"),
				"+",
				buildPostAggregationFields(v)))
		}
	}
	if len(violationDurationAggs) > 0 {
		// Sum the violation duration per metric to get an overall violation duration.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalViolationDuration",
			"+",
			buildPostAggregationFields(violationDurationAggs)))
	}
	if len(violationDurationAggsByEvent) > 0 {
		for ek, v := range violationDurationAggsByEvent {
			postAggregations = append(postAggregations, godruid.PostAggArithmetic(
				buildTopLevelEventAgg(ek, "totalViolationDuration"),
				"+",
				buildPostAggregationFields(v)))
		}

	}
	if len(durationAggs) > 0 {
		// Sum the total duration per metric to get an overall total duration.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalDuration",
			"+",
			buildPostAggregationFields(durationAggs)))
	}

	return &godruid.QueryTimeseries{
		QueryType:   godruid.TIMESERIES,
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
		),
		Aggregations:     aggregations,
		PostAggregations: postAggregations,
		Intervals:        []string{interval}}, nil
}

func SLAViolationsQuery(tenant string, dataSource string, metaMOs []string, granularity string, interval string, thresholdProfile *tenmod.ThresholdProfile, timeout int32) (*godruid.QueryTimeseries, metrics.DruidViolationsMap, error) {
	var aggregations []godruid.Aggregation
	var postAggregations []godruid.PostAggregation
	var violationCountAggs []string
	var totalDurationAggs []string
	var violationDurationAggs []string
	var objectDirectionFilters []*godruid.Filter
	var baselineCalculationCols []godruid.VirtualColumn

	type objectTypeDirectionFilters struct {
		BaseFilter       *godruid.Filter
		ThresholdFilters []*godruid.Filter
	}

	responseSchemaMap := make(metrics.DruidViolationsMap)
	sortedVendorMapKeys := getSortedKeySlice(reflect.ValueOf(thresholdProfile.Thresholds.VendorMap).MapKeys())
	for _, vendorKey := range sortedVendorMapKeys {
		vk := vendorKey
		v := thresholdProfile.Thresholds.VendorMap[vk]

		sortedMOTypeMapKeys := getSortedKeySlice(reflect.ValueOf(v.MonitoredObjectTypeMap).MapKeys())
		for _, typeKey := range sortedMOTypeMapKeys {
			tk := typeKey
			t := v.MonitoredObjectTypeMap[tk]
			// @TODO: HEY! We need to remove this!
			if tk != "twamp-sf" {
				continue
			}

			perDirectionFilters := make(map[string]*objectTypeDirectionFilters)

			sortedMetricTypeMapKeys := getSortedKeySlice(reflect.ValueOf(t.MetricMap).MapKeys())
			for _, metricTypeKey := range sortedMetricTypeMapKeys {
				mk := metricTypeKey
				m := t.MetricMap[mk]
				sortedDirectionTypeMapKeys := getSortedKeySlice(reflect.ValueOf(m.DirectionMap).MapKeys())
				for _, directionTypeKey := range sortedDirectionTypeMapKeys {
					dk := directionTypeKey
					d := m.DirectionMap[dk]
					sortedEventTypeMapKeys := getSortedKeySlice(reflect.ValueOf(d.EventMap).MapKeys())
					for _, eventTypeKey := range sortedEventTypeMapKeys {
						ek := eventTypeKey
						e := d.EventMap[ek]
						if ek != "sla" {
							continue
						}

						objectTypeAndDirectionFilter := godruid.FilterAnd(
							godruid.FilterSelector("objectType", tk),
							godruid.FilterSelector("direction", dk),
						)

						aggNamePrefix := vk + "." + tk + "." + mk + "." + ek + "." + dk

						thresholdFilter, err := BuildThresholdCrossingFilter(aggNamePrefix, mk, e, &baselineCalculationCols)
						if err != nil {
							return nil, nil, err
						}

						dirFilters, ok := perDirectionFilters[vk+"."+tk+"."+dk]
						if !ok {
							perDirectionFilters[vk+"."+tk+"."+dk] = &objectTypeDirectionFilters{objectTypeAndDirectionFilter, []*godruid.Filter{thresholdFilter}}
						} else {
							dirFilters.ThresholdFilters = append(dirFilters.ThresholdFilters, thresholdFilter)
						}

						// Count violations for this metric
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type: "count",
								Name: aggNamePrefix + ".violationCount",
							},
						))

						responseSchemaMap.AddMetric(aggNamePrefix+".violationCount",
							mk, aggNamePrefix, "violationCount", vendorKey, tk, dk)

						violationCountAggs = append(violationCountAggs, aggNamePrefix+".violationCount")

						// Sum the total duration this metric was measured.
						aggregations = append(aggregations, godruid.AggFiltered(
							objectTypeAndDirectionFilter,
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      aggNamePrefix + ".totalDuration",
							},
						))

						responseSchemaMap.AddMetric(aggNamePrefix+".totalDuration",
							mk, aggNamePrefix, "totalDuration", vendorKey, tk, dk)

						// Sum the duration while this metric was in violation.
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      aggNamePrefix + ".violationDuration",
							},
						))

						responseSchemaMap.AddMetric(aggNamePrefix+".violationDuration",
							mk, aggNamePrefix, "violationDuration", vendorKey, tk, dk)

					}
				}
			}

			if len(perDirectionFilters) > 0 {

				// Sum the duration per vendor/objecttype/direction
				for k, v := range perDirectionFilters {
					objectDirectionFilters = append(objectDirectionFilters, v.BaseFilter)

					aggregations = append(aggregations, godruid.AggFiltered(
						v.BaseFilter,
						&godruid.Aggregation{
							Type:      "longSum",
							FieldName: "duration",
							Name:      k + ".totalDuration",
						},
					))
					totalDurationAggs = append(totalDurationAggs, k+".totalDuration")

					if len(v.ThresholdFilters) > 0 {
						// Sum the violation duration per vendor/objecttype/direction
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								v.BaseFilter,
								godruid.FilterOr(v.ThresholdFilters...),
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      k + ".violationDuration",
							},
						))
						violationDurationAggs = append(violationDurationAggs, k+".violationDuration")
					}

				}

			}
		}
	}

	// Count the monitored objects
	if len(objectDirectionFilters) > 0 {

		aggregations = append(aggregations, godruid.AggFiltered(
			godruid.FilterOr(
				objectDirectionFilters...,
			),
			&godruid.Aggregation{
				Type:       "cardinality",
				Name:       "objectCount",
				FieldNames: []string{"monitoredObjectId"},
				ByRow:      false,
				Round:      true,
			},
		))

		responseSchemaMap.AddMetric("objectCount",
			"objectCount", "objectCount", "SLA_Summary", "objectCount", "objectCount", "objectCount")

	}

	if len(violationCountAggs) > 0 {
		// Sum the violation count per metric to get an overall total.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalViolationCount",
			"+",
			buildPostAggregationFields(violationCountAggs)))
	}

	responseSchemaMap.AddMetric("totalViolationCount",
		"totalViolationCount", "totalViolationCount", "SLA_Summary", "totalViolationCount", "totalViolationCount", "totalViolationCount")

	if len(violationDurationAggs) > 0 {
		// Sum the violation duration per metric to get an overal violation duration.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalViolationDuration",
			"+",
			buildPostAggregationFields(violationDurationAggs)))
	}

	responseSchemaMap.AddMetric("totalViolationDuration",
		"totalViolationDuration", "totalViolationDuration", "SLA_Summary", "totalViolationDuration", "totalViolationDuration", "totalViolationCount")

	if len(totalDurationAggs) > 0 {
		// Sum the total duration per metric to get an overall total duration.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalDuration",
			"+",
			buildPostAggregationFields(totalDurationAggs)))
	}

	responseSchemaMap.AddMetric("totalDuration",
		"totalDuration", "totalDuration", "SLA_Summary", "totalDuration", "totalDuration", "totalDuration")

	return &godruid.QueryTimeseries{
			QueryType:   godruid.TIMESERIES,
			DataSource:  dataSource,
			Granularity: toGranularity(granularity),
			Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
			Filter: godruid.FilterAnd(
				godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
				cleanFilter(),
				BuildMonitoredObjectFilter(tenant, metaMOs),
			),
			VirtualColumns:   baselineCalculationCols,
			Aggregations:     aggregations,
			PostAggregations: postAggregations,
			Intervals:        []string{interval}},
		responseSchemaMap,
		nil
}

func SLATimeBucketQuery(tenant string, dataSource string, metaMOs []string, timeBucket int, timeZone string, vendor, objectType, metric, direction, event string, eventAttr *tenmod.ThrPrfEventAttrMap, granularity string, interval string, timeout int32) (*godruid.QueryTopN, metrics.DruidViolationsMap, error) {
	var aggregations []godruid.Aggregation
	var dimension godruid.DimSpec
	var baselineCalculationCols []godruid.VirtualColumn
	schema := make(metrics.DruidViolationsMap)
	prefix := vendor + "." + objectType + "." + metric + "." + event + "." + direction

	threshold := 0
	if timeBucket == DayOfWeek {
		threshold = 7
		schema.AddMetric(prefix+".dayOfWeek", metric, prefix, "dayOfWeek", vendor, objectType, direction)

		dimension = godruid.TimeExtractionDimensionSpec{
			Type:       "extraction",
			Dimension:  "__time",
			OutputName: prefix + ".dayOfWeek",
			ExtractionFunction: godruid.TimeExtractionFn{
				Type:     "timeFormat",
				Format:   "e",
				TimeZone: timeZone,
				Locale:   "en",
			},
		}
	} else if timeBucket == HourOfDay {
		schema.AddMetric(prefix+".hourOfDay", metric, prefix, "hourOfDay", vendor, objectType, direction)
		threshold = 24
		dimension = godruid.TimeExtractionDimensionSpec{
			Type:       "extraction",
			Dimension:  "__time",
			OutputName: prefix + ".hourOfDay",
			ExtractionFunction: godruid.TimeExtractionFn{
				Type:     "timeFormat",
				Format:   "HH",
				TimeZone: timeZone,
				Locale:   "en",
			},
		}
	} else {
		return nil, nil, fmt.Errorf("Invalid value for 'timeBucket' : %v", timeBucket)
	}

	thresholdFilter, err := BuildThresholdCrossingFilter(prefix, metric, eventAttr, &baselineCalculationCols)
	if err != nil {
		return nil, nil, err
	}

	// Count violations for this metric
	countName := vendor + "." + objectType + "." + metric + "." + event + "." + direction + ".violationCount"

	schema.AddMetric(countName, metric, prefix, "violationCount", vendor, objectType, direction)

	aggregations = append(aggregations, godruid.Aggregation{
		Type: "count",
		Name: countName,
	})

	sort.Slice(aggregations, func(i, j int) bool {
		return aggregations[i].Aggregator.Name < aggregations[j].Aggregator.Name
	})

	return &godruid.QueryTopN{
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
			thresholdFilter,
			godruid.FilterSelector("objectType", objectType),
			godruid.FilterSelector("direction", direction),
		),
		Metric:         countName,
		Dimension:      dimension,
		Threshold:      threshold,
		Aggregations:   aggregations,
		VirtualColumns: baselineCalculationCols,
		Intervals:      []string{interval},
	}, schema, nil

}

// ThresholdCrossingByMonitoredObjectQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile. Groups results my monitored object ID.
func ThresholdCrossingByMonitoredObjectTopNQuery(tenant string, dataSource string, metaMOs []string, metric metrics.MetricIdentifierFilter, granularity string, interval string, thresholdProfile *tenmod.ThresholdProfile, timeout int32, numResults int32) (*godruid.QueryTopN, error) {

	var aggregations []godruid.Aggregation
	var postAggregations []godruid.PostAggregation
	// columns used for baseline calculations
	var baselineCalculationCols []godruid.VirtualColumn

	var eventWeights = make(map[string]float32)
	eventWeights["minor"] = 0.0001
	eventWeights["major"] = 0.001
	eventWeights["critical"] = 1

	aggregations = append(aggregations, godruid.AggCount("total"))

	vendorMap := thresholdProfile.Thresholds.VendorMap
	// NOTE: We make the assumption here, for now, that all of the object types that are passed in will be grouped in some way in the threshold profile and therefore will
	// have the same event thresholds for each object type. This may not be the case in the future but for now it is.
	events := vendorMap[metric.Vendor].MonitoredObjectTypeMap[metric.ObjectType[0]].MetricMap[metric.Metric].DirectionMap[metric.Direction[0]].EventMap

	sortedEventTypeMapKeys := getSortedKeySlice(reflect.ValueOf(events).MapKeys())
	for _, eventTypeKey := range sortedEventTypeMapKeys {
		ek := eventTypeKey
		e := events[ek]

		name := ek
		filter, err := BuildThresholdCrossingFilter(name, metric.Metric, e, &baselineCalculationCols)
		if err != nil {
			return nil, err
		}

		aggregation := godruid.AggFiltered(
			filter,
			&godruid.Aggregation{
				Type: "count",
				Name: name,
			},
		)

		postAggregation := godruid.PostAggArithmetic("", "*", []godruid.PostAggregation{
			godruid.PostAggConstant("", eventWeights[ek]),
			godruid.PostAggFieldAccessor(ek),
		})

		postAggregations = append(postAggregations, postAggregation)
		aggregations = append(aggregations, aggregation)
	}

	var scoredPostAggregation []godruid.PostAggregation

	switch len(postAggregations) {
	case 0:
		return nil, fmt.Errorf("No threshold crossing events configured with threshold profile %s for the requested metrics", thresholdProfile.Name)
	case 1: // There should not be any post aggregation if we have less than 2 events that we care about since we need 2 entries for the aggregation operator so we deal with this by adding 0
		postAggregations = append(postAggregations, godruid.PostAggConstant("", 0))
		scoredPostAggregation = []godruid.PostAggregation{
			godruid.PostAggArithmetic("scored", "+", postAggregations),
		}
	default:
		scoredPostAggregation = []godruid.PostAggregation{
			godruid.PostAggArithmetic("scored", "+", postAggregations),
		}
	}

	return &godruid.QueryTopN{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
			buildInFilter("objectType", metric.ObjectType),
			buildInFilter("direction", metric.Direction),
		),
		PostAggregations: scoredPostAggregation,
		VirtualColumns:   baselineCalculationCols,
		Intervals:        []string{interval},
		Metric:           "scored",
		Threshold:        int(numResults),
		Dimension:        "monitoredObjectId",
	}, nil
}

// ThresholdCrossingByMonitoredObjectQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile. Groups results my monitored object ID.
func ThresholdCrossingByMonitoredObjectTopNQueryV1(tenant string, dataSource string, metaMOs []string, metric string, granularity string, interval string, objectType string, direction string, thresholdProfile *pb.TenantThresholdProfileData, vendor string, timeout int32, numResults int32) (*godruid.QueryTopN, error) {

	var aggregations []godruid.Aggregation
	var postAggregations []godruid.PostAggregation

	var eventWeights = make(map[string]float32)
	eventWeights["minor"] = 0.0001
	eventWeights["major"] = 0.001
	eventWeights["critical"] = 1

	aggregations = append(aggregations, godruid.AggCount("total"))

	vendorMap := thresholdProfile.GetThresholds().GetVendorMap()
	events := vendorMap[vendor].GetMonitoredObjectTypeMap()[objectType].GetMetricMap()[metric].GetDirectionMap()[direction].GetEventMap()

	sortedEventTypeMapKeys := getSortedKeySlice(reflect.ValueOf(events).MapKeys())
	for _, eventTypeKey := range sortedEventTypeMapKeys {
		ek := eventTypeKey
		e := events[ek]
		// for ek, e := range events {
		name := ek
		filter, err := FilterHelper(metric, e.EventAttrMap)
		if err != nil {
			return nil, err
		}
		aggregation := godruid.AggFiltered(
			godruid.FilterAnd(
				filter,
			),
			&godruid.Aggregation{
				Type: "count",
				Name: name,
			},
		)

		postAggregation := godruid.PostAggArithmetic("", "*", []godruid.PostAggregation{
			godruid.PostAggConstant("", eventWeights[ek]),
			godruid.PostAggFieldAccessor(ek),
		})

		postAggregations = append(postAggregations, postAggregation)
		aggregations = append(aggregations, aggregation)
	}

	var scoredPostAggregation []godruid.PostAggregation
	if len(postAggregations) == 0 {
		scoredPostAggregation = []godruid.PostAggregation{}
	} else {
		scoredPostAggregation = []godruid.PostAggregation{
			godruid.PostAggArithmetic("scored", "+", postAggregations),
		}
	}

	return &godruid.QueryTopN{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
			godruid.FilterSelector("objectType", objectType),
			godruid.FilterSelector("direction", direction),
		),
		PostAggregations: scoredPostAggregation,
		Intervals:        []string{interval},
		Metric:           "scored",
		Threshold:        int(numResults),
		Dimension:        "monitoredObjectId",
	}, nil
}

//RawMetricsQuery  - Query that returns a raw metric values
func RawMetricsQuery(tenant string, dataSource string, metrics []string, interval string, objectType string, directions []string, monitoredObjects []string, timeout int32, granularity string, cleanOnly bool) (*godruid.QueryTimeseries, error) {

	var aggregations []godruid.Aggregation

	for _, monObj := range monitoredObjects {
		for _, metric := range metrics {
			for _, direction := range directions {
				aggregationMax := godruid.AggFiltered(
					godruid.FilterAnd(
						godruid.FilterSelector("monitoredObjectId", monObj),
						godruid.FilterSelector("direction", direction),
						godruid.FilterSelector("objectType", objectType),
					),
					&godruid.Aggregation{
						Type:      "doubleMax",
						Name:      monObj + "." + objectType + "." + direction + "." + metric,
						FieldName: metric,
					},
				)
				aggregations = append(aggregations, aggregationMax)
			}
		}
	}

	var queryFilter *godruid.Filter
	if cleanOnly {
		queryFilter = godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
		)
	} else {
		queryFilter = godruid.FilterSelector("tenantId", strings.ToLower(tenant))
	}
	return &godruid.QueryTimeseries{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Aggregations: aggregations,
		Filter:       queryFilter,
		Intervals:    []string{interval},
	}, nil
}

//RawMetricsQueryV1  - Query that returns a raw metric values
//DEPRECATED: TERMINATE ONCE V1 IS NOT SUPPORTED
func RawMetricsQueryV1(tenant string, dataSource string, metrics []string, interval string, objectType string, directions []string, monitoredObjects []string, timeout int32, granularity string, cleanOnly bool) (*godruid.QueryTimeseries, error) {

	var aggregations []godruid.Aggregation

	for _, monObj := range monitoredObjects {
		for _, metric := range metrics {
			for _, direction := range directions {
				aggregationMax := godruid.AggFiltered(
					godruid.FilterAnd(
						godruid.FilterSelector("monitoredObjectId", monObj),
						godruid.FilterSelector("direction", direction),
					),
					&godruid.Aggregation{
						Type:      "doubleMax",
						Name:      monObj + "." + direction + "." + metric,
						FieldName: metric,
					},
				)
				aggregations = append(aggregations, aggregationMax)
			}
		}
	}

	var queryFilter *godruid.Filter
	if cleanOnly {
		queryFilter = godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			godruid.FilterSelector("objectType", objectType),
		)
	} else {
		godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			godruid.FilterSelector("objectType", objectType),
		)
	}
	return &godruid.QueryTimeseries{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Aggregations: aggregations,
		Filter:       queryFilter,
		Intervals:    []string{interval},
	}, nil
}

//AggMetricsQuery  - Query that returns a aggregated metric values
func AggMetricsQuery(tenant string, dataSource string, interval string, monitoredObjectIds []string, aggregateOnMeta bool, aggregationFunc string, metrics []metrics.MetricIdentifierFilter, ignoreCleaning bool, timeout int32, granularity string) (godruid.Query, *PostProcessor, *db.QueryKeySpec, error) {

	var aggregations []godruid.Aggregation
	var pp PostProcessor
	postAggs := []godruid.PostAggregation{}

	keyToDrop := []string{}
	countKeys := map[string][]string{}

	var moidAggregateList []string

	queryspec := db.QueryKeySpec{}

	// If the initial request contained meta data or no explicit monitored objects were provided
	// then we want to aggregate our response on all monitored objects that fit that metadata set or all monitored objects respectively
	// Otherwise we need to break down our query into individual metric queries based on monitored object ID
	if aggregateOnMeta || len(monitoredObjectIds) == 0 {
		moidAggregateList = []string{""} // Place an empty value in order to loop through at least once to add the aggregation query
	} else {
		moidAggregateList = monitoredObjectIds
	}
	for _, mo := range moidAggregateList {
		for _, metric := range metrics {
			var querySpecID string
			// If the initial request does not contain explicit monitored object Ids then we do not need to separate our query out by monitored object ID
			if aggregateOnMeta || len(monitoredObjectIds) == 0 {
				querySpecID = queryspec.AddKeySpec(map[string]interface{}{"vendor": metric.Vendor, "objectType": metric.ObjectType, "direction": metric.Direction, "metric": metric.Metric})
			} else {
				querySpecID = queryspec.AddKeySpec(map[string]interface{}{"monitoredObjectIds": []string{mo}, "vendor": metric.Vendor, "objectType": metric.ObjectType, "direction": metric.Direction, "metric": metric.Metric})
			}
			countName := querySpecID + db.QueryDelimeter + "count"
			keyToDrop = append(keyToDrop, countName)
			countKeys[countName] = []string{querySpecID}
			aggregations = append(aggregations, buildMetricAggregation("count", mo, &metric, countName))
			if aggregationFunc == "max" {
				aggregations = append(aggregations, buildMetricAggregation("doubleMax", mo, &metric, querySpecID))

			} else if aggregationFunc == "min" {
				aggregations = append(aggregations, buildMetricAggregation("doubleMin", mo, &metric, querySpecID))

			} else if aggregationFunc == "avg" {

				aggregations = append(aggregations, buildMetricAggregation("doubleSum", mo, &metric, querySpecID+db.QueryDelimeter+"sum"))

				keyToDrop = append(keyToDrop, querySpecID+db.QueryDelimeter+"sum")

				postAgg := godruid.PostAggArithmetic(
					querySpecID,
					"/",
					[]godruid.PostAggregation{godruid.PostAggFieldAccessor(querySpecID + db.QueryDelimeter + "sum"), godruid.PostAggFieldAccessor(querySpecID + db.QueryDelimeter + "count")},
				)
				postAggs = append(postAggs, postAgg)
			} else {
				return nil, nil, nil, fmt.Errorf("Invalid value for 'aggregation' : %v", aggregationFunc)
			}

		}
	}
	// Drop the intermediate sum and count aggregations after the response returns from druid.
	// There doesn't seem to be an option in the druid query to do this.
	pp = DropKeysPostprocessor{
		keysToDrop: keyToDrop,
		countKeys:  countKeys,
	}

	var cleanFilterRef *godruid.Filter
	if !ignoreCleaning {
		cleanFilterRef = cleanFilter()
	}

	return &godruid.QueryTimeseries{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilterRef,
			BuildMonitoredObjectFilter(tenant, monitoredObjectIds),
		),
		Intervals:        []string{interval},
		PostAggregations: postAggs,
	}, &pp, &queryspec, nil
}

//AggMetricsQuery  - Query that returns a aggregated metric values
//DEPRECATED: TERMINATE ONCE V1 IS NOT SUPPORTED
func AggMetricsQueryV1(tenant string, dataSource string, interval string, metaMOs []string, aggregationFunc metrics.AggregationSpecV1, metrics []metrics.MetricIdentifierV1, timeout int32, granularity string) (*godruid.QueryTimeseries, *PostProcessorV1, error) {

	var aggregations []godruid.Aggregation
	var pp PostProcessorV1
	postAggs := []godruid.PostAggregation{}

	keyToDrop := []string{}
	countKeys := map[string][]string{}

	for _, metric := range metrics {
		countName := metric.Name + "Count"
		keyToDrop = append(keyToDrop, countName)
		countKeys[countName] = []string{metric.Name}
		aggregations = append(aggregations, buildMetricAggregationV1("count", &metric, countName))
		if aggregationFunc.Name == "max" {
			aggregations = append(aggregations, buildMetricAggregationV1("doubleMax", &metric))

		} else if aggregationFunc.Name == "min" {
			aggregations = append(aggregations, buildMetricAggregationV1("doubleMin", &metric))

		} else if aggregationFunc.Name == "avg" {

			aggregations = append(aggregations, buildMetricAggregationV1("doubleSum", &metric, metric.Name+"Sum"))

			keyToDrop = append(keyToDrop, metric.Name+"Sum")

			postAgg := godruid.PostAggArithmetic(
				metric.Name,
				"/",
				[]godruid.PostAggregation{godruid.PostAggFieldAccessor(metric.Name + "Sum"), godruid.PostAggFieldAccessor(metric.Name + "Count")},
			)
			postAggs = append(postAggs, postAgg)
		} else {
			return nil, nil, fmt.Errorf("Invalid value for 'aggregation' : %v", aggregationFunc)
		}

	}

	// Drop the intermediate sum and count aggregations after the response returns from druid.
	// There doesn't seem to be an option in the druid query to do this.
	pp = DropKeysPostprocessorV1{
		keysToDrop: keyToDrop,
		countKeys:  countKeys,
	}

	return &godruid.QueryTimeseries{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			cleanFilter(),
			BuildMonitoredObjectFilter(tenant, metaMOs),
		),
		Intervals:        []string{interval},
		PostAggregations: postAggs,
	}, &pp, nil
}

func buildMetricAggregation(aggType string, monitoredObjectID string, metric *metrics.MetricIdentifierFilter, aggName string) godruid.Aggregation {

	filterSet := []*godruid.Filter{
		buildInFilter("objectType", metric.ObjectType),
		buildInFilter("direction", metric.Direction)}

	// Add the monitored object ID as a filter if it has been provided
	if len(monitoredObjectID) != 0 {
		filterSet = append(filterSet, godruid.FilterSelector("monitoredObjectId", monitoredObjectID))
	}

	return godruid.AggFiltered(
		godruid.FilterAnd(filterSet...),
		&godruid.Aggregation{
			Type:      aggType,
			Name:      aggName,
			FieldName: metric.Metric,
		})

}

// DEPRECATED
func buildMetricAggregationV1(aggType string, metric *metrics.MetricIdentifierV1, name ...string) godruid.Aggregation {
	var aggName string
	if len(name) == 0 {
		aggName = metric.Name
	} else {
		aggName = name[0]
	}
	return godruid.AggFiltered(
		godruid.FilterAnd(
			godruid.FilterSelector("objectType", metric.ObjectType),
			godruid.FilterSelector("direction", metric.Direction),
		),
		&godruid.Aggregation{
			Type:      aggType,
			Name:      aggName,
			FieldName: metric.Name,
		})

}

// BuildMonitoredObjectFilter - Builds a monitored object filter for druid
func BuildMonitoredObjectFilter(tenantID string, mos []string) *godruid.Filter {
	// It is important to draw a distinction between a nil set of monitored objects and an empty set of monitored objects. A nil
	// set means that the query does not care which monitored object it belongs to. An empty set means that no monitored objects
	// were found as a result of potential pre-filtering activities
	if mos == nil {
		return nil
	}

	return godruid.FilterAnd(
		godruid.FilterSelector("tenantId", strings.ToLower(tenantID)),
		buildInFilter("monitoredObjectId", mos))
}

func buildInFilter(dimension string, values []string) *godruid.Filter {
	if values == nil {
		return nil
	}

	var fValues []string

	if len(values) == 0 {
		fValues = []string{"[\"\"]"}
	} else {
		// Sort the entries to ensure that we will generate the same filter for the same set of input entries
		sort.Slice(values, func(i, j int) bool {
			return values[i] < values[j]
		})
		fValues = values
	}

	return &godruid.Filter{
		Type:      "in",
		Dimension: dimension,
		Values:    fValues,
	}
}

func buildOrFilter(dimensionName string, values []string) *godruid.Filter {
	if len(values) < 1 {
		return nil
	}
	filters := make([]*godruid.Filter, len(values))
	for i, value := range values {
		filters[i] = godruid.FilterSelector(dimensionName, value)
	}
	return godruid.FilterOr(filters...)

}

func buildPostAggregationFields(fieldNames []string) []godruid.PostAggregation {

	fields := make([]godruid.PostAggregation, len(fieldNames))
	for i, name := range fieldNames {
		fields[i] = godruid.PostAggFieldAccessor(name)
	}

	if len(fields) == 1 {
		// If we are only summing 1 value, satisfy the post aggregator by adding 0.
		fields = append(fields, godruid.PostAggConstant("const", 0))
	}
	return fields
}

func toGranularity(granularityStr string) godruid.Granlarity {
	if strings.ToLower(granularityStr) == GranularityAll {
		return godruid.GranAll
	} else if strings.ToLower(granularityStr) == GranularityNone {
		return godruid.GranNone
	}
	return godruid.GranPeriod(granularityStr, TimeZoneUTC, "")
}

const (
	op_sum   = "sum"
	op_max   = "max"
	op_min   = "min"
	op_count = "count"
	op_avg   = "avg"
)

func buildMetricAggregator(metricsView []metrics.MetricAggregation) []godruid.Aggregation {

	var aggregations []godruid.Aggregation

	if metricsView == nil {
		return nil
	}
	for _, input := range metricsView {

		switch input.Aggregator {
		case op_sum:
			aggregations = append(aggregations, godruid.AggDoubleSum(input.Name, input.Metric))
		case op_max:
			aggregations = append(aggregations, godruid.AggDoubleMax(input.Name, input.Metric))
		case op_min:
			aggregations = append(aggregations, godruid.AggDoubleMin(input.Name, input.Metric))
		case op_count:
			aggregations = append(aggregations, godruid.AggCount(input.Name))
		}
	}

	return aggregations
}

// GetTopNForMetricAvg - Provides TopN for certain metrics.
func GetTopNForMetric(dataSource string, request *metrics.TopNForMetric, timeout int32, metaMOs []string) (*godruid.QueryTopN, error) {

	var aggregations []godruid.Aggregation
	var postAggregations godruid.PostAggregation
	var scoredPostAggregation []godruid.PostAggregation

	// Create the labels for the average operation (for some reason,
	// druid has no native idea of average but it does for SUM)
	const (
		sumLbl    = "topn_sum"
		countLbl  = "count"
		opLbl     = "value"
		metricLbl = "metric"
		typeLbl   = "type"
	)

	typeInvertedLbl := "inverted"
	// Only operate on the first item
	metric := request.Metric

	// Metric order and sort on
	selectedMetric := map[string]interface{}{metricLbl: opLbl}
	// Create the Filters

	var molist []string

	// We can only filter on the monitored objects associated with the request meta or the monitored objects
	// explicitly asked for in the request but not both.
	// If none are provided then the query defaults to checking all monitored objects with the metric filters
	if metaMOs != nil {
		molist = metaMOs
	} else {
		molist = request.MonitoredObjects
	}

	filterOn := godruid.FilterAnd(
		godruid.FilterSelector("tenantId", strings.ToLower(request.TenantID)),
		cleanFilter(),
		buildInFilter("objectType", metric.ObjectType),
		buildInFilter("direction", request.Metric.Direction),
		BuildMonitoredObjectFilter(request.TenantID, molist),
	)

	// Create the aggregations

	// We need the total COUNT for the average op
	//aggregations = append(aggregations, godruid.AggCount(countFilterLbl))

	// Build the metricView. This isn't part of the average operation but it
	// helps the caller have more information about the object druid finds.
	// Eg: For Top N average for DelayP95, what is its MAX/MIN DelayP95 & Delay & Max dropped Packets and Max Jitter.
	// We can even show the average for the metricsView but it requires more POST Aggregations. May a future feature.
	listOfMetricViewAgg := buildMetricAggregator(request.MetricsView)
	if listOfMetricViewAgg != nil {
		aggregations = append(aggregations, listOfMetricViewAgg...)
	}

	switch request.Aggregator {
	case op_max:
		aggregations = append(aggregations, godruid.AggDoubleMax(opLbl, metric.Metric))
		break
	case op_min:
		aggregations = append(aggregations, godruid.AggDoubleMin(opLbl, metric.Metric))
		selectedMetric[typeLbl] = typeInvertedLbl
		break
	default:
		// We need the SUM to do the average operation
		aggregations = append(aggregations, godruid.AggDoubleSum(sumLbl, metric.Metric))
		// Makes sure we don't pass in a 0 into a division operation (not necessary actually,
		// testing shows druid doesn't segfault on a divide by zero operation and returns 0 as a result).
		aggroFilter := godruid.FilterNot(godruid.FilterSelector(metric.Metric, 0))
		aggroCount := godruid.AggCount(countLbl)
		aggroFunc := godruid.AggFiltered(aggroFilter, &aggroCount)
		aggregations = append(aggregations, aggroFunc)
		// Post Aggregation is where the Average operation is executed
		postAggregations.Fields = append(postAggregations.Fields, godruid.PostAggFieldAccessor(sumLbl))
		postAggregations.Fields = append(postAggregations.Fields, godruid.PostAggFieldAccessor(countLbl))
		// We can actually define more operations here if we really wanted to.
		scoredPostAggregation = []godruid.PostAggregation{
			godruid.PostAggArithmetic(opLbl, "/", postAggregations.Fields),
		}

	}
	return &godruid.QueryTopN{
		QueryType:        godruid.TOPN,
		DataSource:       dataSource,
		Granularity:      godruid.GranAll,
		Context:          map[string]interface{}{"timeout": timeout, "queryId": uuid.NewV4().String()},
		Aggregations:     aggregations,
		Filter:           filterOn,
		PostAggregations: scoredPostAggregation,
		Intervals:        []string{request.Interval},
		// !! LOOK HERE. Metric is used to tell Druid which METRIC we want to sort by.
		// Because the average operation is a POST AGGREGATION and not an existing column,
		// the metric name here must match the POST AGGREGATION name for it to sort.
		// Default is to sort in descending order (use `"type":"inverted"` to reverse the order)
		Metric:    selectedMetric,
		Threshold: int(request.NumResult),
		Dimension: "monitoredObjectId",
	}, nil
}

// DEPRECATED - TERMINATE ONCE V1 IS REMOVED
// GetTopNForMetricAvg - Provides TopN for certain metrics.
func GetTopNForMetricV1(dataSource string, request *metrics.TopNForMetricV1, metaMOs []string) (*godruid.QueryTopN, error) {

	var aggregations []godruid.Aggregation
	var postAggregations godruid.PostAggregation
	var scoredPostAggregation []godruid.PostAggregation

	// Create the labels for the average operation (for some reason,
	// druid has no native idea of average but it does for SUM)
	const (
		sumLbl    = "topn_sum"
		countLbl  = "topn_count"
		opLbl     = "result"
		metricLbl = "metric"
		typeLbl   = "type"
	)

	typeInvertedLbl := "inverted"
	// Only operate on the first item
	metric := request.Metric

	// Metric order and sort on
	selectedMetric := map[string]interface{}{metricLbl: opLbl}
	// Create the Filters
	// TODO: I think we may need to specify DIRECTION for monitored objects.
	// TODO: The existing MetricIdentifier model does not work for directions since
	// it isn't well defined and we can't use it to specify more than 1 direction.

	var filterOn *godruid.Filter

	filterOn = godruid.FilterAnd(
		godruid.FilterSelector("tenantId", strings.ToLower(request.TenantID)),
		cleanFilter(),
		godruid.FilterSelector("objectType", metric.ObjectType),
	)

	// Prefer the domains list
	if len(request.MonitoredObjects) == 0 {
		filterOn = godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(request.TenantID)),
			cleanFilter(),
			godruid.FilterSelector("objectType", metric.ObjectType),
			BuildMonitoredObjectFilter(request.TenantID, metaMOs),
		)
	} else {
		monObjFilter := BuildMonitoredObjectFilter(request.TenantID, request.MonitoredObjects)
		filterOn = godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(request.TenantID)),
			cleanFilter(),
			godruid.FilterSelector("objectType", metric.ObjectType),
			monObjFilter,
		)
	}

	// Create the aggregations

	// We need the total COUNT for the average op
	//aggregations = append(aggregations, godruid.AggCount(countFilterLbl))

	// Build the metricView. This isn't part of the average operation but it
	// helps the caller have more information about the object druid finds.
	// Eg: For Top N average for DelayP95, what is its MAX/MIN DelayP95 & Delay & Max dropped Packets and Max Jitter.
	// We can even show the average for the metricsView but it requires more POST Aggregations. May a future feature.
	listOfMetricViewAgg := buildMetricAggregator(request.MetricsView)
	if listOfMetricViewAgg != nil {
		aggregations = append(aggregations, listOfMetricViewAgg...)
	}

	switch request.Aggregator {
	case op_max:
		aggregations = append(aggregations, godruid.AggDoubleMax(opLbl, metric.Name))
		break
	case op_min:
		aggregations = append(aggregations, godruid.AggDoubleMin(opLbl, metric.Name))
		selectedMetric[typeLbl] = typeInvertedLbl
		break
	default:
		// We need the SUM to do the average operation
		aggregations = append(aggregations, godruid.AggDoubleSum(sumLbl, metric.Name))
		// Makes sure we don't pass in a 0 into a division operation (not necessary actually,
		// testing shows druid doesn't segfault on a divide by zero operation and returns 0 as a result).
		aggroFilter := godruid.FilterNot(godruid.FilterSelector(metric.Name, 0))
		aggroCount := godruid.AggCount(countLbl)
		aggroFunc := godruid.AggFiltered(aggroFilter, &aggroCount)
		aggregations = append(aggregations, aggroFunc)
		// Post Aggregation is where the Average operation is executed
		postAggregations.Fields = append(postAggregations.Fields, godruid.PostAggFieldAccessor(sumLbl))
		postAggregations.Fields = append(postAggregations.Fields, godruid.PostAggFieldAccessor(countLbl))
		// We can actually define more operations here if we really wanted to.
		scoredPostAggregation = []godruid.PostAggregation{
			godruid.PostAggArithmetic(opLbl, "/", postAggregations.Fields),
		}

	}
	return &godruid.QueryTopN{
		QueryType:        godruid.TOPN,
		DataSource:       dataSource,
		Granularity:      godruid.GranAll,
		Context:          map[string]interface{}{"timeout": request.Timeout, "queryId": uuid.NewV4().String()},
		Aggregations:     aggregations,
		Filter:           filterOn,
		PostAggregations: scoredPostAggregation,
		Intervals:        []string{request.Interval},
		// !! LOOK HERE. Metric is used to tell Druid which METRIC we want to sort by.
		// Because the average operation is a POST AGGREGATION and not an existing column,
		// the metric name here must match the POST AGGREGATION name for it to sort.
		// Default is to sort in descending order (use `"type":"inverted"` to reverse the order)
		Metric:    selectedMetric,
		Threshold: int(request.NumResult),
		Dimension: "monitoredObjectId",
	}, nil
}
func inWhitelist(whitelist []metrics.MetricIdentifierFilter, vendor, objectType, metricName, direction string) bool {
	if whitelist == nil || len(whitelist) == 0 {
		return true
	}

	for _, mi := range whitelist {
		if vendor == mi.Vendor && contains(mi.ObjectType, objectType) && metricName == mi.Metric && contains(mi.Direction, direction) {
			return true
		}
	}
	return false
}

// DEPRECATED
func inWhitelistV1(whitelist []metrics.MetricIdentifierV1, vendor, objectType, metricName, direction string) bool {
	if whitelist == nil || len(whitelist) == 0 {
		return true
	}

	for _, mi := range whitelist {
		if vendor == mi.Vendor && objectType == mi.ObjectType && metricName == mi.Name && direction == fmt.Sprint(mi.Direction) {
			return true
		}
	}
	return false
}

func cleanFilter() *godruid.Filter {
	// CleanStatus is null (for old data) OR cleanStatus > -1.
	return godruid.FilterOr(
		godruid.FilterSelector("cleanStatus", ""),
		godruid.FilterLowerBound("cleanStatus", godruid.NUMERIC, -1, true),
	)
}

// getSortedKeySlice - function used to make sure our keys are sorted before we build druid queries which prevents us from ending up with different
// queries for the same data which results in us missing the cache on subsequent requests
func getSortedKeySlice(originalSlice []reflect.Value) []string {

	keys := make([]string, len(originalSlice))
	for i, val := range originalSlice {
		keys[i] = val.String()
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}
