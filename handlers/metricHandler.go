package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/models/metrics"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
)

const (
	contentType        = "Content-Type"
	jsonAPIContentType = "application/vnd.api+json"
)

type MetricServiceHandler struct {
	druidDB  db.DruidDatastore
	tenantDB db.TenantServiceDatastore
	routes   []server.Route
}

func CreateMetricServiceHandler() *MetricServiceHandler {
	result := new(MetricServiceHandler)

	ddb := druid.NewDruidDatasctoreClient()
	result.druidDB = ddb

	tdb, err := GetTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminServiceRESTHandler: %s", err.Error())
	}
	result.tenantDB = tdb

	result.routes = []server.Route{

		server.Route{
			Name:        "QueryThresholdCrossing",
			Method:      "POST",
			Pattern:     "/api/v1/threshold-crossing",
			HandlerFunc: result.QueryThresholdCrossing,
		},

		server.Route{
			Name:        "GetThresholdCrossingByMonitoredObject",
			Method:      "POST",
			Pattern:     "/api/v1/threshold-crossing-by-monitored-object",
			HandlerFunc: result.GetThresholdCrossingByMonitoredObject,
		},

		server.Route{
			Name:        "GetThresholdCrossingByMonitoredObjectTopN",
			Method:      "POST",
			Pattern:     "/api/v1/threshold-crossing-by-monitored-object-top-n",
			HandlerFunc: result.GetThresholdCrossingByMonitoredObjectTopN,
		},

		server.Route{
			Name:        "GenSLAReport",
			Method:      "POST",
			Pattern:     "/api/v1/generate-sla-report",
			HandlerFunc: result.GetSLAReport,
		},

		server.Route{
			Name:        "GetHistogram",
			Method:      "POST",
			Pattern:     "/api/v1/histogram",
			HandlerFunc: result.GetHistogram,
		},

		server.Route{
			Name:        "GetRawMetrics",
			Method:      "GET",
			Pattern:     "/api/v1/raw-metrics",
			HandlerFunc: result.GetRawMetrics,
		},

		server.Route{
			Name:        "QueryAggregatedMetrics",
			Method:      "POST",
			Pattern:     "/api/v1/aggregated-metrics",
			HandlerFunc: result.QueryAggregatedMetrics,
		},

		server.Route{
			Name:        "GetTopNFor",
			Method:      "POST",
			Pattern:     "/api/v1/topn-metrics",
			HandlerFunc: result.GetTopNFor,
		},
	}

	return result
}

// RegisterAPIHandlers - will bind any REST API routes defined in this service
// to the passed in request multiplexor.
func (msh *MetricServiceHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range msh.routes {
		logger.Log.Debugf("Registering endpoint: %v", route)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

func populateRawMetricsRequest(queryParams url.Values) *pb.RawMetricsRequest {
	rmr := pb.RawMetricsRequest{
		Direction:         toStringSplice(queryParams.Get("direction")),
		Interval:          queryParams.Get("interval"),
		Metric:            toStringSplice(queryParams.Get("metric")),
		Tenant:            queryParams.Get("tenant"),
		ObjectType:        queryParams.Get("objectType"),
		MonitoredObjectId: toStringSplice(queryParams.Get("monitoredObjectId")),
		Granularity:       queryParams.Get("granularity"),
	}

	timeout, err := strconv.Atoi(queryParams.Get("timeout"))
	if err == nil {
		rmr.Timeout = int32(timeout)
	} else {
		rmr.Timeout = 0
	}

	return &rmr
}

func (msh *MetricServiceHandler) QueryThresholdCrossing(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.QueryThresholdCrossingStr, msg, http.StatusBadRequest)
		return
	}
	request := metrics.ThresholdCrossingRequest{}
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.QueryThresholdCrossingStr, msg, http.StatusBadRequest)
		return
	}

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(request.TenantID, request.ThresholdProfileID)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query: %s. Error: %s", models.AsJSONString(request), err.Error())
		reportError(w, startTime, "404", mon.QueryThresholdCrossingStr, msg, http.StatusNotFound)
		return
	}
	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.QueryThresholdCrossingStr, err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossStr, msg, http.StatusNotFound)
		return
	}

	logger.Log.Infof("Retrieving %s for: %v", db.QueryThresholdCrossingStr, request)

	result, err := msh.druidDB.QueryThresholdCrossing(&request, &pbTP)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve  Threshold Crossing Metrics. %s:", err.Error())
		reportError(w, startTime, "500", mon.QueryThresholdCrossingStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal Threshold Crossing response. %s:", err.Error())
		reportError(w, startTime, "500", mon.QueryThresholdCrossingStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.QueryThresholdCrossingStr, request)
	trackAPIMetrics(startTime, "200", mon.QueryThresholdCrossingStr)
	fmt.Fprintf(w, string(res))
}

func (msh *MetricServiceHandler) GetInternalSLAReport(slaReportRequest *metrics.SLAReportRequest) (*metrics.SLAReport, error) {
	startTime := time.Now()

	// Turn the query Params into the request object:
	logger.Log.Debugf("Retrieving %s for: %v", db.SLAReportStr, models.AsJSONString(slaReportRequest))

	tenantID := slaReportRequest.TenantID

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, slaReportRequest.ThresholdProfileID)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", models.AsJSONString(slaReportRequest), err.Error())
		reportInternalError(startTime, "404", mon.GetSLAReportStr, msg)
		return nil, err
	}

	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.SLAReportStr, err.Error())
		reportInternalError(startTime, "500", mon.GetSLAReportStr, msg)
		return nil, err
	}

	report, err := msh.druidDB.GetSLAReport(slaReportRequest, &pbTP)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve SLA Report. %s:", err.Error())
		reportInternalError(startTime, "500", mon.GetSLAReportStr, msg)
		return nil, err
	}

	report.ReportScheduleConfig = slaReportRequest.SlaScheduleConfig
	report.TenantID = slaReportRequest.TenantID
	logger.Log.Debugf("Completed %s fetch for: %+v, report %+v", db.SLAReportStr, models.AsJSONString(slaReportRequest), report)

	// STORE into DB the generated SLA report
	trackAPIMetrics(startTime, "200", mon.GetSLAReportStr)
	return report, nil
}

func (msh *MetricServiceHandler) GetSLAReport(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GenerateSLAReportStr, msg, http.StatusBadRequest)
		return
	}
	slaReportRequest := metrics.SLAReportRequest{}
	if err := json.Unmarshal(requestBytes, &slaReportRequest); err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GenerateSLAReportStr, msg, http.StatusBadRequest)
		return
	}
	logger.Log.Infof("Retrieving %s for: %v", db.SLAReportStr, models.AsJSONString(slaReportRequest))

	tenantID := slaReportRequest.TenantID

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, slaReportRequest.ThresholdProfileID)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", models.AsJSONString(slaReportRequest), err.Error())
		reportError(w, startTime, "404", mon.GenerateSLAReportStr, msg, http.StatusNotFound)
		return
	}

	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.SLAReportStr, err.Error())
		reportError(w, startTime, "500", mon.GenerateSLAReportStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetSLAReport(&slaReportRequest, &pbTP)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve SLA Report. %s:", err.Error())
		reportError(w, startTime, "500", mon.GenerateSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := jsonapi.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal SLA Report. %s:", err.Error())
		reportError(w, startTime, "500", mon.GenerateSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.SLAReportStr, models.AsJSONString(slaReportRequest))
	trackAPIMetrics(startTime, "200", mon.GenerateSLAReportStr)
	fmt.Fprintf(w, string(res))
}

// GetThresholdCrossingByMonitoredObject - Retrieves the Threshold crossings for a given threshold profile,
// interval, tenant, domain, and groups by monitoredObjectID
func (msh *MetricServiceHandler) GetThresholdCrossingByMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GetThrCrossByMonObjStr, msg, http.StatusBadRequest)
		return
	}
	thresholdCrossingReq := metrics.ThresholdCrossingRequest{}
	if err := json.Unmarshal(requestBytes, &thresholdCrossingReq); err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GetThrCrossByMonObjStr, msg, http.StatusBadRequest)
		return
	}

	if len(thresholdCrossingReq.Granularity) == 0 {
		thresholdCrossingReq.Granularity = "PT1H"
	}
	logger.Log.Infof("Retrieving %s for: %v", db.ThresholdCrossingByMonitoredObjectStr, thresholdCrossingReq)

	tenantID := thresholdCrossingReq.TenantID

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, thresholdCrossingReq.ThresholdProfileID)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile: %s. Error: %s", thresholdCrossingReq.ThresholdProfileID, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossByMonObjStr, msg, http.StatusNotFound)
		return
	}

	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.ThresholdCrossingStr, err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossByMonObjStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetThresholdCrossingByMonitoredObject(&thresholdCrossingReq, &pbTP)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Threshold Crossing By Monitored Object. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossByMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal Threshold Crossing by Monitored Object. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossByMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.ThresholdCrossingByMonitoredObjectStr, thresholdCrossingReq)
	trackAPIMetrics(startTime, "200", mon.GetThrCrossByMonObjStr)
	fmt.Fprintf(w, string(res))
}

// GetThresholdCrossingByMonitoredObjectTopN - Retrieves the TopN Threshold crossings for a given threshold profile,
// interval, tenant, domain, and groups by monitoredObjectID
func (msh *MetricServiceHandler) GetThresholdCrossingByMonitoredObjectTopN(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusBadRequest)
		return
	}
	thresholdCrossingReq := metrics.ThresholdCrossingTopNRequest{}
	if err := json.Unmarshal(requestBytes, &thresholdCrossingReq); err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusBadRequest)
		return
	}

	if thresholdCrossingReq.NumResults == 0 {
		thresholdCrossingReq.NumResults = 10 // default value
	}

	if len(thresholdCrossingReq.Granularity) == 0 {
		thresholdCrossingReq.Granularity = "PT1H"
	}

	if len(thresholdCrossingReq.Vendor) == 0 {
		msg := generateErrorMessage(http.StatusBadRequest, "vendor is required")
		reportError(w, startTime, "400", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusBadRequest)
		return
	}

	if len(thresholdCrossingReq.ObjectType) == 0 {
		msg := generateErrorMessage(http.StatusBadRequest, "object type is required")
		reportError(w, startTime, "400", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusBadRequest)
		return
	}

	if err != nil {
		reportError(w, startTime, "602", mon.GetThrCrossByMonObjTopNStr, err.Error(), 602)
		return
	}
	logger.Log.Infof("Retrieving %s for: %v", db.TopNThresholdCrossingByMonitoredObjectStr, thresholdCrossingReq)

	tenantID := thresholdCrossingReq.TenantID

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, thresholdCrossingReq.ThresholdProfileID)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %+v. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusNotFound)
		return
	}

	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.ThresholdCrossingStr, err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusNotFound)
		return
	}

	if err = validateMetricForThresholdProfile(thresholdCrossingReq.Vendor, thresholdCrossingReq.ObjectType, thresholdCrossingReq.Metric, &pbTP); err != nil {
		reportError(w, startTime, "404", mon.GetThrCrossByMonObjTopNStr, err.Error(), http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetThresholdCrossingByMonitoredObjectTopN(&thresholdCrossingReq, &pbTP)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Threshold Crossing By Monitored Object. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal Threshold Crossing by Monitored Object. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.TopNThresholdCrossingByMonitoredObjectStr, thresholdCrossingReq)
	trackAPIMetrics(startTime, "200", mon.GetThrCrossByMonObjTopNStr)
	fmt.Fprintf(w, string(res))
}

// GetHistogram - Retrieve bucket data from druid
func (msh *MetricServiceHandler) GetHistogram(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GetHistogramObjStr, msg, http.StatusBadRequest)
		return
	}

	// Turn the query Params into the request object:
	hcRequest := &metrics.HistogramRequest{}
	err = json.Unmarshal(requestBytes, hcRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Histogram. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetHistogramObjStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Retrieving %s for: %v", db.HistogramStr, hcRequest)

	result, err := msh.druidDB.GetHistogram(hcRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Histogram. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetHistogramObjStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal Histogram response. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetHistogramObjStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.HistogramStr, hcRequest)
	trackAPIMetrics(startTime, "200", mon.GetHistogramObjStr)
	fmt.Fprintf(w, string(res))
}

// GetRawMetrics - Retrieve raw metric data from druid
func (msh *MetricServiceHandler) GetRawMetrics(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Turn the query Params into the request object:
	queryParams := r.URL.Query()
	rawMetricReq := populateRawMetricsRequest(queryParams)
	logger.Log.Infof("Retrieving %s for: %v", db.RawMetricStr, rawMetricReq)

	result, err := msh.druidDB.GetRawMetrics(rawMetricReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Raw Metrics. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetRawMetricStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal Raw Metrics response. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetRawMetricStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.RawMetricStr, rawMetricReq)
	trackAPIMetrics(startTime, "200", mon.GetRawMetricStr)
	fmt.Fprintf(w, string(res))
}

func (msh *MetricServiceHandler) QueryAggregatedMetrics(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.QueryAggregatedMetricsStr, msg, http.StatusBadRequest)
		return
	}
	request := metrics.AggregateMetricsAPIRequest{}
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.QueryAggregatedMetricsStr, msg, http.StatusBadRequest)
		return
	}
	logger.Log.Infof("Retrieving %s for: %v", db.AggMetricsStr, request)

	result, err := msh.druidDB.GetAggregatedMetrics(&request)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Aggregated Metrics. %s:", err.Error())
		reportError(w, startTime, "500", mon.QueryAggregatedMetricsStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal Aggregated Metrics response. %s:", err.Error())
		reportError(w, startTime, "500", mon.QueryAggregatedMetricsStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.AggMetricsStr, request)
	trackAPIMetrics(startTime, "200", mon.QueryAggregatedMetricsStr)
	fmt.Fprintf(w, string(res))
}

func toStringSplice(paramCSV string) []string {
	if len(paramCSV) < 1 {
		return nil
	}

	return strings.Split(paramCSV, ",")
}

func validateMetricForThresholdProfile(vendor, objectType, metric string, thresholdProfile *pb.TenantThresholdProfile) error {
	vendorMap := thresholdProfile.Data.GetThresholds().GetVendorMap()
	vendorEntry, ok := vendorMap[vendor]
	if !ok {
		return fmt.Errorf("Vendor %s not found in threshold profile with ID %s", vendor, thresholdProfile.GetXId())
	}

	objectTypeEntry, ok := vendorEntry.GetMonitoredObjectTypeMap()[objectType]
	if !ok {
		return fmt.Errorf("Object type %s for vendor %s not found in threshold profile with ID %s", objectType, vendor, thresholdProfile.GetXId())
	}

	_, ok = objectTypeEntry.GetMetricMap()[metric]
	if !ok {
		return fmt.Errorf("Metric %s for vendor %s and object type %s not found in threshold profile with ID %s", metric, vendor, objectType, thresholdProfile.GetXId())
	}

	return nil
}

func (msh *MetricServiceHandler) GetTopNFor(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.QueryAggregatedMetricsStr, msg, http.StatusBadRequest)
		return
	}
	request := metrics.TopNForMetric{}
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.QueryAggregatedMetricsStr, msg, http.StatusBadRequest)
		return
	}
	logger.Log.Infof("Retrieving %s for: %v", "top n req", request)

	if _, err = request.Validate(); err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.GetTopNReqStr, msg, http.StatusBadRequest)
		return
	}

	topNreq := request
	logger.Log.Infof("Fetching data for TopN request: %+v", topNreq)

	result, err := msh.druidDB.GetTopNForMetric(&topNreq)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Top N response. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetTopNReqStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal TOP N. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetTopNReqStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.TopNForMetricString, topNreq)
	trackAPIMetrics(startTime, "200", mon.GetTopNReqStr)
	fmt.Fprintf(w, string(res))

}
