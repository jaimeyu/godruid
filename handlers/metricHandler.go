package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
)

const (
	contentTytpe       = "Content-Type"
	jsonAPIContentType = "application/vnd.api+json"
)

type MetricServiceHandler struct {
	druidDB db.DruidDatastore
	routes  []server.Route
	gsh     *GRPCServiceHandler
}

func CreateMetricServiceHandler(grpcServiceHandler *GRPCServiceHandler) *MetricServiceHandler {
	result := new(MetricServiceHandler)

	db := druid.NewDruidDatasctoreClient()

	result.druidDB = db

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

	result.gsh = grpcServiceHandler

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
		Direction:          queryParams.Get("direction"),
		Domain:             queryParams.Get("domain"),
		Granularity:        queryParams.Get("granularity"),
		Interval:           queryParams.Get("interval"),
		Metric:             queryParams.Get("metric"),
		ObjectType:         queryParams.Get("objectType"),
		Tenant:             queryParams.Get("tenant"),
		ThresholdProfileId: queryParams.Get("thresholdProfileId"),
		Vendor:             queryParams.Get("vendor"),
	}

	timeout, err := strconv.Atoi(queryParams.Get("timeout"))
	if err == nil {
		thresholdCrossingReq.Timeout = int32(timeout)
	} else {
		thresholdCrossingReq.Timeout = 0
	}

	return &thresholdCrossingReq
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

func populateRawMetricsRequest(queryParams url.Values) *pb.RawMetricsRequest {
	rmr := pb.RawMetricsRequest{
		Direction:         queryParams.Get("direction"),
		Interval:          queryParams.Get("interval"),
		Metric:            queryParams.Get("metric"),
		Tenant:            queryParams.Get("tenant"),
		ObjectType:        queryParams.Get("objectType"),
		MonitoredObjectId: queryParams.Get("monitoredObjectId"),
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

	thresholdProfile, err := msh.gsh.GetTenantThresholdProfile(nil, &pb.TenantThresholdProfileIdRequest{
		TenantId:           tenantID,
		ThresholdProfileId: thresholdCrossingReq.ThresholdProfileId,
	})

	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetThresholdCrossing(thresholdCrossingReq, thresholdProfile)
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

	w.Header().Set(contentTytpe, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.ThresholdCrossingStr, thresholdCrossingReq)
	trackAPIMetrics(startTime, "200", mon.GetThrCrossStr)
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

	thresholdProfile, err := msh.gsh.GetTenantThresholdProfile(nil, &pb.TenantThresholdProfileIdRequest{
		TenantId:           tenantID,
		ThresholdProfileId: thresholdCrossingReq.ThresholdProfileId,
	})

	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossByMonObjStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetThresholdCrossingByMonitoredObject(thresholdCrossingReq, thresholdProfile)
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

	w.Header().Set(contentTytpe, jsonAPIContentType)
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

	w.Header().Set(contentTytpe, jsonAPIContentType)
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

	w.Header().Set(contentTytpe, jsonAPIContentType)
	logger.Log.Infof("Completed %s fetch for: %v", db.RawMetricStr, rawMetricReq)
	trackAPIMetrics(startTime, "200", mon.GetRawMetricStr)
	fmt.Fprintf(w, string(res))
}
