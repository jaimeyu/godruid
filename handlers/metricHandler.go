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
			Name:        "GetThresholdCrossingByMonitoredObjectTopN",
			Method:      "GET",
			Pattern:     "/api/v1/threshold-crossing-by-monitored-object-top-n",
			HandlerFunc: result.GetThresholdCrossingByMonitoredObjectTopN,
		},

		server.Route{
			Name:        "GenSLAReport",
			Method:      "GET",
			Pattern:     "/api/v1/generate-sla-report",
			HandlerFunc: result.GetSLAReport,
		},

		server.Route{
			Name:        "GetHistogram",
			Method:      "GET",
			Pattern:     "/api/v1/histogram",
			HandlerFunc: result.GetHistogram,
		},

		server.Route{
			Name:        "GetHistogramCustom",
			Method:      "POST",
			Pattern:     "/api/v1/histogram-custom",
			HandlerFunc: result.GetHistogramCustom,
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

	if len(thresholdCrossingReq.Granularity) == 0 {
		thresholdCrossingReq.Granularity = "PT1H"
	}

	return &thresholdCrossingReq
}

func populateThresholdCrossingTopNRequest(queryParams url.Values) (*metrics.ThresholdCrossingTopNRequest, error) {

	thresholdCrossingReq := metrics.ThresholdCrossingTopNRequest{
		Direction:          queryParams.Get("direction"),
		Domain:             toStringSplice(queryParams.Get("domain")),
		Granularity:        queryParams.Get("granularity"),
		Interval:           queryParams.Get("interval"),
		Metric:             queryParams.Get("metric"),
		ObjectType:         queryParams.Get("objectType"),
		TenantID:           queryParams.Get("tenantId"),
		ThresholdProfileID: queryParams.Get("thresholdProfileId"),
		Vendor:             queryParams.Get("vendor"),
	}

	timeout, err := strconv.Atoi(queryParams.Get("timeout"))
	if err == nil {
		thresholdCrossingReq.Timeout = int32(timeout)
	} else {
		thresholdCrossingReq.Timeout = 5000 // default value
	}

	numResults, err := strconv.Atoi(queryParams.Get("numResults"))
	if err == nil {
		thresholdCrossingReq.NumResults = int32(numResults)
	} else {
		thresholdCrossingReq.NumResults = 10 // default value
	}

	if len(thresholdCrossingReq.Granularity) == 0 {
		thresholdCrossingReq.Granularity = "PT1H"
	}

	if len(thresholdCrossingReq.Vendor) == 0 {
		err = fmt.Errorf("vendor is required")
		return nil, err
	}

	if len(thresholdCrossingReq.ObjectType) == 0 {
		err = fmt.Errorf("objectType is required")
		return nil, err
	}

	return &thresholdCrossingReq, nil
}

func populateSLAReportRequest(queryParams url.Values) *metrics.SLAReportRequest {

	request := metrics.SLAReportRequest{
		TenantID:           queryParams.Get("tenant"),
		Interval:           queryParams.Get("interval"),
		Domain:             toStringSplice(queryParams.Get("domains")),
		ThresholdProfileID: queryParams.Get("thresholdProfileId"),
		Granularity:        queryParams.Get("granularity"),
		Timezone:           queryParams.Get("timezone"),
	}

	timeout, err := strconv.Atoi(queryParams.Get("timeout"))
	if err == nil {
		request.Timeout = int32(timeout)
	} else {
		request.Timeout = 0
	}

	if len(request.Granularity) == 0 {
		request.Granularity = "PT1H"
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

	if err = msh.validateDomains(thresholdCrossingReq.Tenant, thresholdCrossingReq.Domain); err != nil {
		msg := fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossStr, msg, http.StatusNotFound)
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

	if err := msh.validateDomains(slaReportRequest.TenantID, slaReportRequest.Domain); err != nil {
		msg := fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", models.AsJSONString(slaReportRequest), err.Error())
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

	// Turn the query Params into the request object:
	queryParams := r.URL.Query()
	slaReportRequest := populateSLAReportRequest(queryParams)
	logger.Log.Infof("Retrieving %s for: %v", db.SLAReportStr, models.AsJSONString(slaReportRequest))

	tenantID := slaReportRequest.TenantID

	thresholdProfile, err := msh.tenantDB.GetTenantThresholdProfile(tenantID, slaReportRequest.ThresholdProfileID)
	if err != nil {
		msg := fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", models.AsJSONString(slaReportRequest), err.Error())
		reportError(w, startTime, "404", mon.GenerateSLAReportStr, msg, http.StatusNotFound)
		return
	}

	if err = msh.validateDomains(slaReportRequest.TenantID, slaReportRequest.Domain); err != nil {
		msg := fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", models.AsJSONString(slaReportRequest), err.Error())
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

	result, err := msh.druidDB.GetSLAReport(slaReportRequest, &pbTP)
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
		reportError(w, startTime, "500", mon.GetThrCrossByMonObjStr, msg, http.StatusNotFound)
		return
	}

	if err = msh.validateDomains(thresholdCrossingReq.Tenant, thresholdCrossingReq.Domain); err != nil {
		msg := fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossByMonObjStr, msg, http.StatusNotFound)
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

// GetThresholdCrossingByMonitoredObjectTopN - Retrieves the TopN Threshold crossings for a given threshold profile,
// interval, tenant, domain, and groups by monitoredObjectID
func (msh *MetricServiceHandler) GetThresholdCrossingByMonitoredObjectTopN(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Turn the query Params into the request object:
	queryParams := r.URL.Query()
	thresholdCrossingReq, err := populateThresholdCrossingTopNRequest(queryParams)
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

	if err = msh.validateDomains(thresholdCrossingReq.TenantID, thresholdCrossingReq.Domain); err != nil {
		msg := fmt.Sprintf("Unable find domain for given query parameters: %+v. Error: %s", thresholdCrossingReq, err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossByMonObjTopNStr, msg, http.StatusNotFound)
		return
	}

	result, err := msh.druidDB.GetThresholdCrossingByMonitoredObjectTopN(thresholdCrossingReq, &pbTP)
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

// GetHistogram - Retrieve bucket data from druid
func (msh *MetricServiceHandler) GetHistogramCustom(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		//TODO FIX THE STRIIIIIIIIIIING
		reportError(w, startTime, "400", mon.CreateTenantStr, msg, http.StatusBadRequest)
		return
	}

	// Turn the query Params into the request object:
	hcRequest := &metrics.HistogramCustomRequest{}
	err = json.Unmarshal(requestBytes, hcRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Custom Histogram. %s:", err.Error())
		reportError(w, startTime, "500", mon.GetHistogramObjStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Retrieving %s for: %v", db.HistogramStr, hcRequest)

	result, err := msh.druidDB.GetHistogramCustom(hcRequest)
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

	if err = msh.validateDomains(request.TenantID, request.DomainIDs); err != nil {
		msg := fmt.Sprintf("Unable find domain for given request: %s. Error: %s", models.AsJSONString(request), err.Error())
		reportError(w, startTime, "404", mon.GetThrCrossStr, msg, http.StatusNotFound)
		return
	}

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

func (msh *MetricServiceHandler) validateDomains(tenantId string, domains []string) error {
	if domains == nil || len(domains) == 0 {
		return nil
	}
	for _, dom := range domains {
		if _, err := msh.tenantDB.GetTenantDomain(tenantId, dom); err != nil {
			return err
		}
	}
	return nil
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
