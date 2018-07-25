package druid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/models/metrics"
	"github.com/accedian/godruid"
	uuid "github.com/satori/go.uuid"
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

func convertHistogramCustomResponse(tenantId string, meta map[string]string, interval string, rawResponse string) (map[string]interface{}, error) {

	const (
		HistogramCustomReport = "customHistogramReports"
		AttrData              = "data"
		AttrTimestamp         = "timestamp"
		AttrResult            = "result"
		KeyDelim              = "."

		IndexVendor      = 1
		IndexObjectType  = 2
		IndexMetricName  = 3
		IndexDirection   = 4
		IndexBucketIndex = 5
	)

	fieldsRegex := regexp.MustCompile(`(?P<Vendor>.+?)\.(?P<ObjectType>.+?)\.(?P<MetricName>.+?)\.(?P<Direction>.+?).(?P<Index>.+)`)
	metrickeyRegex := regexp.MustCompile(`(?P<Vendor>.+?)\.(?P<ObjectType>.+?)\.(?P<MetricName>.+?)\.(?P<Direction>.+)`)

	// Temporary hack to put the payload in a format understandable by the json library
	jsonResponse, err := gabs.ParseJSON([]byte(fmt.Sprintf(`{"data":%s}`, rawResponse)))
	if err != nil {
		return nil, err
	}

	hcReport := metrics.HistogramCustomReport{}
	timeSlices := make([]metrics.HistogramCustomTimeSeriesEntry, 0)

	// Process each time slice in the raw druid response
	rawTimeslices, _ := jsonResponse.S(AttrData).Children()
	for _, rawTimeslice := range rawTimeslices {
		// Create a new timeslice for the current set of time that we are processing from the druid response
		timeslice := metrics.HistogramCustomTimeSeriesEntry{Timestamp: rawTimeslice.S(AttrTimestamp).Data().(string)}

		// Process each bucket response for each metric in the time slice
		rawResultMap, _ := rawTimeslice.S(AttrResult).ChildrenMap()
		resultMap := make(map[string][]metrics.BucketResult)
		for rawkey, value := range rawResultMap {

			// Build up the key as vendor.objecttype.metricname.direction. This allows us to figure out which buckets belong to which metric since
			// there can be multiple in the original request
			fields := fieldsRegex.FindStringSubmatch(rawkey)
			mapkey := fields[IndexVendor] + KeyDelim + fields[IndexObjectType] + KeyDelim + fields[IndexMetricName] + KeyDelim + fields[IndexDirection]

			bucketResult := metrics.BucketResult{Index: fields[IndexBucketIndex], Count: int(value.Data().(float64))}

			// If a set of result buckets does not already exist for the metric we need to create one and add it to our current time slice
			metricBucket, found := resultMap[mapkey]
			if !found {
				metricBucket = make([]metrics.BucketResult, 0)
			}
			metricBucket = append(metricBucket, bucketResult)
			resultMap[mapkey] = metricBucket
		}

		metricResults := make([]metrics.MetricResult, 0)
		for k, m := range resultMap {
			keyfields := metrickeyRegex.FindStringSubmatch(k)
			metricResults = append(metricResults, metrics.MetricResult{Vendor: keyfields[IndexVendor],
				ObjectType: keyfields[IndexObjectType],
				Name:       keyfields[IndexMetricName],
				Direction:  keyfields[IndexDirection],
				Results:    m})
		}
		timeslice.Result = metricResults
		timeSlices = append(timeSlices, timeslice)
	}

	hcReport.TimeSeriesResult = timeSlices
	hcReport.TenantID = tenantId
	hcReport.Meta = meta
	hcReport.ReportTimeRange = interval
	hcReport.ReportCompletionTime = time.Now().UTC().String()

	uuid := uuid.NewV4()
	rr := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         uuid.String(),
			"type":       HistogramCustomReport,
			"attributes": hcReport,
		},
	}

	return rr, nil
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
			direction := parts[1]
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
				obj.SetP(v, direction+"."+lastParts)
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

func reformatReportSummary(druidResponse []byte) (*metrics.ReportSummary, error) {
	logger.Log.Debugf("Response from druid for %s: %s", db.SLAReportStr, string(druidResponse))
	entries := []*druidTimeSeriesEntry{}
	if err := json.Unmarshal(druidResponse, &entries); err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return &metrics.ReportSummary{}, nil
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

	summary := metrics.ReportSummary{}
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

func reformatThresholdCrossingTimeSeries(druidResponse []byte) ([]metrics.ThresholdCrossingTimeSeriesEntry, error) {
	logger.Log.Debugf("Response from druid for %s: %s", db.SLAReportStr, string(druidResponse))
	entries := []*druidTimeSeriesEntry{}
	if err := json.Unmarshal(druidResponse, &entries); err != nil {
		return nil, err
	}

	res := make([]metrics.ThresholdCrossingTimeSeriesEntry, len(entries))
	for i, tc := range entries {

		perMetricResultsMap := map[string]*metrics.ThresholdCrossingMetricResult{}
		perMetricResultsList := []*metrics.ThresholdCrossingMetricResult{}
		byEventMap := make(map[string]map[string]interface{})

		obj := gabs.New()
		for k, v := range tc.Result {
			if isMetricAggName(k) {
				aggrName := parseMetricName(k)
				metricResult, ok := perMetricResultsMap[aggrName.fqMetricName]
				if !ok {
					metricResult = &metrics.ThresholdCrossingMetricResult{
						Vendor:     aggrName.vendor,
						ObjectType: aggrName.objectType,
						Metric:     aggrName.metricName,
						Direction:  aggrName.direction,
						BySeverity: make(map[string]map[string]interface{}),
					}
					perMetricResultsMap[aggrName.fqMetricName] = metricResult
					perMetricResultsList = append(perMetricResultsList, metricResult)
				}

				if len(aggrName.severity) > 0 {
					sevEntry, ok := metricResult.BySeverity[aggrName.severity]
					if !ok {
						sevEntry = make(map[string]interface{})
						metricResult.BySeverity[aggrName.severity] = sevEntry
					}
					sevEntry[aggrName.agg] = v
				} else if aggrName.agg == "totalDuration" {
					metricResult.TotalDuration = v.(float64)
				}

			} else if isTopLevelEventAgg(k) {
				aggrName := parseTopLevelEventAgg(k)
				sevEntry, ok := byEventMap[aggrName.severity]
				if !ok {
					sevEntry = make(map[string]interface{})
					byEventMap[aggrName.severity] = sevEntry
				}
				sevEntry[aggrName.agg] = v
			} else {
				obj.SetP(v, k)
			}
		}
		timeseriesEntryResult := metrics.ThresholdCrossingTimeSeriesResult{}
		if err := json.Unmarshal(obj.Bytes(), &timeseriesEntryResult); err != nil {
			return nil, err
		}
		timeseriesEntryResult.ByMetric = perMetricResultsList
		timeseriesEntryResult.BySeverity = byEventMap

		res[i] = metrics.ThresholdCrossingTimeSeriesEntry{
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

// For postprocessing metrics
type PostProcessor interface {
	Apply(input []AggMetricsResponse) []AggMetricsResponse
}

var (
	NOOP_POSTPROCESSOR = NoopPostProcessor{}
)

type NoopPostProcessor struct{}

func (pp NoopPostProcessor) Apply(input []AggMetricsResponse) []AggMetricsResponse {
	logger.Log.Debugf("NoopPostProcessor.apply called")
	return input
}

type DropKeysPostprocessor struct {
	keysToDrop []string
	countKeys  map[string][]string
}

func (pp DropKeysPostprocessor) Apply(input []AggMetricsResponse) []AggMetricsResponse {
	logger.Log.Debugf("DropKeysPostprocessor.apply called with %v, %v, %v", pp.keysToDrop, pp.countKeys, input)
	if len(pp.keysToDrop) > 0 {
		for _, v := range input {
			for countKey, vals := range pp.countKeys {
				if countVal, ok := v.Result[countKey]; ok {

					if intVal, ok := countVal.(float64); ok && intVal == 0 {
						for _, m := range vals {
							delete(v.Result, m)
						}
					}
				}

			}
			for _, k := range pp.keysToDrop {
				delete(v.Result, k)
			}

		}
	}

	return input
}

// The following functions and structs are to ensure proper building and parsing of aggregation names
type parsedMetricName struct {
	fqMetricName string
	vendor       string
	objectType   string
	metricName   string
	direction    string
	severity     string
	agg          string
}

func parseMetricName(name string) parsedMetricName {
	// The aggregation name is __met.<vendor>|<objectType>|<metricName>|<direction>.<severity>.<aggregation>
	// or  __met.<vendor>|<objectType>|<metricName>|<direction>.<aggregation>

	tokens := strings.Split(name, ".")
	tokens2 := strings.Split(tokens[1], "|")
	res := parsedMetricName{
		fqMetricName: tokens[1],
		vendor:       tokens2[0],
		objectType:   tokens2[1],
		metricName:   tokens2[2],
		direction:    tokens2[3],
	}
	if len(tokens) > 3 {
		res.severity = tokens[2]
		res.agg = tokens[3]
	} else if len(tokens) > 2 {
		res.agg = tokens[2]
	}
	return res

}

func buildMetricAggPrefix(vendor, objectType, metricName, direction string) string {
	return "__met." + vendor + "|" + objectType + "|" + metricName + "|" + direction
}

func isMetricAggName(name string) bool {
	return strings.HasPrefix(name, "__met.")
}

type topLevelEventAgg struct {
	severity string
	agg      string
}

func parseTopLevelEventAgg(name string) topLevelEventAgg {
	// __eventtop.<severity>.<aggName>
	tokens := strings.Split(name, ".")
	return topLevelEventAgg{
		severity: tokens[1],
		agg:      tokens[2],
	}

}
func buildTopLevelEventAgg(severity, agg string) string {
	return "__eventtop." + severity + "." + agg
}
func isTopLevelEventAgg(name string) bool {
	return strings.HasPrefix(name, "__eventtop")
}
