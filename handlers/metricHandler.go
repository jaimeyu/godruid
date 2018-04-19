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

func CreateMetricServiceHandler(grpcServiceHandler *GRPCServiceHandler) *MetricServiceHandler {
	result := new(MetricServiceHandler)

	ddb := druid.NewDruidDatasctoreClient()
	result.druidDB = ddb

	tdb, err := getTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminServiceRESTHandler: %s", err.Error())
	}
	result.tenantDB = tdb

	result.routes = []server.Route{
		server.Route{
			Name:        "GetThresholdCrossing",
			Method:      "GET",
			Pattern:     "/api/v1/threshold-crossing",
			HandlerFunc: result.GetThresholdCrossing,
		},

		server.Route{
			Name:        "GetThresholdCrossingByMonitoredObject",
			Method:      "GET",
			Pattern:     "/api/v1/threshold-crossing-by-monitored-object",
			HandlerFunc: result.GetThresholdCrossingByMonitoredObject,
		},

		server.Route{
			Name:        "GetSLAReport",
			Method:      "GET",
			Pattern:     "/api/v1/sla-report",
			HandlerFunc: result.GetSLAReport,
		},

		server.Route{
			Name:        "GetHistogram",
			Method:      "GET",
			Pattern:     "/api/v1/histogram",
			HandlerFunc: result.GetHistogram,
		},

		server.Route{
			Name:        "GetRawMetrics",
			Method:      "GET",
			Pattern:     "/api/v1/raw-metrics",
			HandlerFunc: result.GetRawMetrics,
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

func populateThresholdCrossingRequest(queryParams url.Values) *pb.ThresholdCrossingRequest {

	thresholdCrossingReq := pb.ThresholdCrossingRequest{
		Direction:          toStringSplice(queryParams.Get("direction")),
		Domain:             toStringSplice(queryParams.Get("domain")),
		Granularity:        queryParams.Get("granularity"),
		Interval:           queryParams.Get("interval"),
		Metric:             toStringSplice(queryParams.Get("metric")),
		ObjectType:         toStringSplice(queryParams.Get("objectType")),
		Tenant:             queryParams.Get("tenant"),
		ThresholdProfileId: queryParams.Get("thresholdProfileId"),
		Vendor:             toStringSplice(queryParams.Get("vendor")),
	}

	timeout, err := strconv.Atoi(queryParams.Get("timeout"))
	if err == nil {
		thresholdCrossingReq.Timeout = int32(timeout)
	} else {
		thresholdCrossingReq.Timeout = 0
	}

	return &thresholdCrossingReq
}

func populateSLAReportRequest(queryParams url.Values) *metrics.SLAReportRequest {

	request := metrics.SLAReportRequest{
		TenantID:           queryParams.Get("tenant"),
		Interval:           queryParams.Get("interval"),
		Domain:             toStringSplice(queryParams.Get("domain")),
		ThresholdProfileID: queryParams.Get("thresholdProfileId"),
	}

	timeout, err := strconv.Atoi(queryParams.Get("timeout"))
	if err == nil {
		request.Timeout = int32(timeout)
	} else {
		request.Timeout = 0
	}

	return &request
}

func populateHistogramRequest(queryParams url.Values) *pb.HistogramRequest {

	histogramRequest := pb.HistogramRequest{
		Direction:   queryParams.Get("direction"),
		Domain:      queryParams.Get("domain"),
		Granularity: queryParams.Get("granularity"),
		Interval:    queryParams.Get("interval"),
		Metric:      queryParams.Get("metric"),
		Tenant:      queryParams.Get("tenant"),
		Vendor:      queryParams.Get("vendor"),
	}

	timeout, err := strconv.Atoi(queryParams.Get("timeout"))
	if err == nil {
		histogramRequest.Timeout = int32(timeout)
	} else {
		histogramRequest.Timeout = 0
	}

	resolution, err := strconv.Atoi(queryParams.Get("resolution"))
	if err == nil {
		histogramRequest.Resolution = int32(resolution)
	} else {
		histogramRequest.Resolution = 0
	}

	granularityBuckets, err := strconv.Atoi(queryParams.Get("granularityBuckets"))
	if err == nil {
		histogramRequest.GranularityBuckets = int32(granularityBuckets)
	} else {
		histogramRequest.GranularityBuckets = 0
	}

	return &histogramRequest
}

// string interval = 1;
// string tenant = 2;
// string direction = 3;
// string metric = 4;
// string objectType = 5;
// string monitoredObjectId = 6;
// int32  timeout = 10;
// string  granularity = 11;

func populateRawMetricsRequest(queryParams url.Values) *pb.RawMetricsRequest {
	rmr := pb.RawMetricsRequest{
		Direction:         queryParams.Get("direction"),
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

// GetThresholdCrossing - Retrieves the Threshold crossings for a given threshold profile,
// interval, tenant, domain
func (msh *MetricServiceHandler) GetThresholdCrossing(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Turn the query Params into the request object:
	queryParams := r.URL.Query()
	thresholdCrossingReq := populateThresholdCrossingRequest(queryParams)
	logger.Log.Infof("Retrieving %s for: %v", db.ThresholdCrossingStr, thresholdCrossingReq)

	tenantID := thresholdCrossingReq.Tenant

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, thresholdCrossingReq.ThresholdProfileId)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossStr, msg, http.StatusNotFound)
		return
	}

	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.ThresholdCrossingStr, err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetThresholdCrossing(thresholdCrossingReq, &pbTP)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Threshold Crossing. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal Threshold Crossing. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.ThresholdCrossingStr, thresholdCrossingReq)
	trackAPIMetrics(startTime, "200", mon.GetThrCrossStr)
	fmt.Fprintf(w, string(res))
}

func (msh *MetricServiceHandler) GetSLAReport(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Turn the query Params into the request object:
	queryParams := r.URL.Query()
	slaReqportRequest := populateSLAReportRequest(queryParams)
	logger.Log.Infof("Retrieving %s for: %v", db.SLAReportStr, models.AsJSONString(slaReqportRequest))

	tenantID := slaReqportRequest.TenantID

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, slaReqportRequest.ThresholdProfileID)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", models.AsJSONString(slaReqportRequest), err.Error())
		reportError(w, startTime, "404", mon.GetSLAReportStr, msg, http.StatusNotFound)
		return
	}

	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.SLAReportStr, err.Error())
		reportError(w, startTime, "500", mon.GetSLAReportStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetSLAReport(slaReqportRequest, &pbTP)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve SLA Report. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal SLA Report. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.SLAReportStr, models.AsJSONString(slaReqportRequest))
	trackAPIMetrics(startTime, "200", mon.GetSLAReportStr)
	fmt.Fprintf(w, string(res))
}

// GetThresholdCrossingByMonitoredObject - Retrieves the Threshold crossings for a given threshold profile,
// interval, tenant, domain, and groups by monitoredObjectID
func (msh *MetricServiceHandler) GetThresholdCrossingByMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Turn the query Params into the request object:
	queryParams := r.URL.Query()
	thresholdCrossingReq := populateThresholdCrossingRequest(queryParams)
	logger.Log.Infof("Retrieving %s for: %v", db.ThresholdCrossingByMonitoredObjectStr, thresholdCrossingReq)

	tenantID := thresholdCrossingReq.Tenant

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, thresholdCrossingReq.ThresholdProfileId)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossByMonObjStr, msg, http.StatusNotFound)
		return
	}

	// Convert to PB type...will remove this when we remove the PB handling
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to fetch %s: %s", db.ThresholdCrossingStr, err.Error())
		reportError(w, startTime, "500", mon.GetThrCrossStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetThresholdCrossingByMonitoredObject(thresholdCrossingReq, &pbTP)
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

// GetHistogram - Retrieve bucket data from druid
func (msh *MetricServiceHandler) GetHistogram(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Turn the query Params into the request object:
	queryParams := r.URL.Query()
	histogramReq := populateHistogramRequest(queryParams)
	logger.Log.Infof("Retrieving %s for: %v", db.HistogramStr, histogramReq)

	result, err := msh.druidDB.GetHistogram(histogramReq)
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
	logger.Log.Infof("Completed %s fetch for: %v", db.HistogramStr, histogramReq)
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

func toStringSplice(paramCSV string) []string {
	if len(paramCSV) < 1 {
		return nil
	}

	return strings.Split(paramCSV, ",")
}
