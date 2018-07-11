package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/accedian/adh-gather/models/metrics"

	"github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/metrics_service"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
)

func populateThresholdCrossingRequestSwag(params metrics_service.GetThresholdCrossingParams) *pb.ThresholdCrossingRequest {

	thresholdCrossingReq := pb.ThresholdCrossingRequest{
		Direction:          params.Direction,
		Domain:             params.Domain,
		Granularity:        *params.Granularity,
		Interval:           params.Interval,
		Metric:             params.Metric,
		ObjectType:         params.ObjectType,
		Tenant:             params.Tenant,
		ThresholdProfileId: params.ThresholdProfileID,
		Vendor:             params.Vendor,
		Timeout:            *params.Timeout,
	}

	if thresholdCrossingReq.Timeout == 0 {
		thresholdCrossingReq.Timeout = 5000 // default value
	}

	if len(thresholdCrossingReq.Granularity) == 0 {
		thresholdCrossingReq.Granularity = "PT1H"
	}

	return &thresholdCrossingReq
}

func populateThresholdCrossingRequestForMonObjSwag(params metrics_service.GetThresholdCrossingByMonitoredObjectParams) *pb.ThresholdCrossingRequest {

	thresholdCrossingReq := pb.ThresholdCrossingRequest{
		Direction:          params.Direction,
		Domain:             params.Domain,
		Granularity:        *params.Granularity,
		Interval:           params.Interval,
		Metric:             params.Metric,
		ObjectType:         params.ObjectType,
		Tenant:             params.Tenant,
		ThresholdProfileId: params.ThresholdProfileID,
		Vendor:             params.Vendor,
		Timeout:            *params.Timeout,
	}

	if thresholdCrossingReq.Timeout == 0 {
		thresholdCrossingReq.Timeout = 5000 // default value
	}

	if len(thresholdCrossingReq.Granularity) == 0 {
		thresholdCrossingReq.Granularity = "PT1H"
	}

	return &thresholdCrossingReq
}

func populateThresholdCrossingTopNRequestSwag(params metrics_service.GetThresholdCrossingByMonitoredObjectTopNParams) (*metrics.ThresholdCrossingTopNRequest, error) {

	thresholdCrossingReq := metrics.ThresholdCrossingTopNRequest{
		Direction:          *params.Direction,
		Domain:             params.Domain,
		Granularity:        *params.Granularity,
		Interval:           params.Interval,
		Metric:             params.Metric,
		ObjectType:         params.ObjectType,
		TenantID:           params.TenantID,
		ThresholdProfileID: params.ThresholdProfileID,
		Vendor:             params.Vendor,
		Timeout:            *params.Timeout,
		NumResults:         *params.NumResults,
	}

	if thresholdCrossingReq.Timeout == 0 {
		thresholdCrossingReq.Timeout = 5000 // default value
	}

	if thresholdCrossingReq.NumResults == 0 {
		thresholdCrossingReq.NumResults = 10 // default value
	}

	if len(thresholdCrossingReq.Granularity) == 0 {
		thresholdCrossingReq.Granularity = "PT1H"
	}

	if len(thresholdCrossingReq.Vendor) == 0 {
		err := fmt.Errorf("vendor is required")
		return nil, err
	}

	if len(thresholdCrossingReq.ObjectType) == 0 {
		err := fmt.Errorf("objectType is required")
		return nil, err
	}

	return &thresholdCrossingReq, nil
}

func populateSLAReportRequestSwag(params metrics_service.GenSLAReportParams) *metrics.SLAReportRequest {

	request := metrics.SLAReportRequest{
		TenantID:           params.Tenant,
		Interval:           params.Interval,
		Domain:             params.Domain,
		ThresholdProfileID: params.ThresholdProfileID,
		Granularity:        *params.Granularity,
		Timezone:           *params.Timezone,
		Timeout:            *params.Timeout,
	}

	if request.Timeout == 0 {
		request.Timeout = 5000 // default value
	}

	if len(request.Granularity) == 0 {
		request.Granularity = "PT1H"
	}

	return &request
}

func populateHistogramRequestSwag(params metrics_service.GetHistogramParams) *pb.HistogramRequest {

	histogramRequest := pb.HistogramRequest{
		Direction:          *params.Direction,
		Domain:             *params.Domain,
		Granularity:        *params.Granularity,
		Interval:           *params.Interval,
		Metric:             *params.Metric,
		Tenant:             *params.Tenant,
		Vendor:             *params.Vendor,
		GranularityBuckets: *params.GranularityBuckets,
		Resolution:         *params.Resolution,
	}

	if histogramRequest.Timeout == 0 {
		histogramRequest.Timeout = 5000 // default value
	}

	return &histogramRequest
}

func populateRawMetricsRequestSwag(params metrics_service.GetRawMetricsParams) *pb.RawMetricsRequest {
	rmr := pb.RawMetricsRequest{
		Direction:         toStringSplice(*params.Direction),
		Interval:          *params.Interval,
		Metric:            params.Metric,
		Tenant:            *params.Tenant,
		ObjectType:        *params.ObjectType,
		MonitoredObjectId: params.MonitoredObjectID,
		Granularity:       *params.Granularity,
		Timeout:           *params.Timeout,
	}

	if rmr.Timeout == 0 {
		rmr.Timeout = 5000 // default value
	}

	return &rmr
}

func validateDomainsSwag(tenantId string, domains []string, tenantDB datastore.TenantServiceDatastore) error {
	if domains == nil || len(domains) == 0 {
		return nil
	}
	for _, dom := range domains {
		if _, err := tenantDB.GetTenantDomain(tenantId, dom); err != nil {
			return err
		}
	}
	return nil
}

// HandleGetThresholdCrossing - retrieves threshold crossing details
func HandleGetThresholdCrossing(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) func(params metrics_service.GetThresholdCrossingParams) middleware.Responder {
	return func(params metrics_service.GetThresholdCrossingParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Fetching %s for Tenant %s using %s %s", datastore.ThresholdCrossingStr, params.Tenant, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewGetThresholdCrossingForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", datastore.ThresholdCrossingStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Turn the query Params into the request object:
		thresholdCrossingReq := populateThresholdCrossingRequestSwag(params)
		tenantID := thresholdCrossingReq.Tenant

		thresholdProfile, err := tenantDB.GetTenantThresholdProfile(tenantID, thresholdCrossingReq.ThresholdProfileId)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error()), startTime, http.StatusNotFound, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert to PB type...will remove this when we remove the PB handling
		pbTP := pb.TenantThresholdProfile{}
		if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
			return metrics_service.NewGetThresholdCrossingInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert request to fetch %s: %s", datastore.ThresholdCrossingStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		if err = validateDomainsSwag(thresholdCrossingReq.Tenant, thresholdCrossingReq.Domain, tenantDB); err != nil {
			return metrics_service.NewGetThresholdCrossingNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error()), startTime, http.StatusNotFound, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		result, err := druidDB.GetThresholdCrossing(thresholdCrossingReq, &pbTP)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve Threshold Crossing. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.ThresholdCrossingStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Retrieved %s for Tenant %s using %s %s", datastore.ThresholdCrossingStr, params.Tenant, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)
		return metrics_service.NewGetThresholdCrossingOK().WithPayload(&converted)
	}
}

// HandleQueryThresholdCrossing - query for threshold crossing data
func HandleQueryThresholdCrossing(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) func(params metrics_service.QueryThresholdCrossingParams) middleware.Responder {
	return func(params metrics_service.QueryThresholdCrossingParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Issuing %s for Tenant %s using %s %s", datastore.QueryThresholdCrossingStr, params.Body.TenantID, tenmod.TenantThresholdProfileStr, params.Body.ThresholdProfileID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewQueryThresholdCrossingForbidden().WithPayload(reportAPIError(fmt.Sprintf("%s operation not authorized for role: %s", datastore.QueryThresholdCrossingStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return metrics_service.NewQueryThresholdCrossingBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}
		request := metrics.ThresholdCrossingRequest{}
		if err := json.Unmarshal(requestBytes, &request); err != nil {
			return metrics_service.NewQueryThresholdCrossingBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		thresholdProfile, err := tenantDB.GetTenantThresholdProfile(*params.Body.TenantID, params.Body.ThresholdProfileID)
		if err != nil {
			return metrics_service.NewQueryThresholdCrossingNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable to find threshold profile %s. Error: %s", request.ThresholdProfileID, err.Error()), startTime, http.StatusNotFound, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}
		// Convert to PB type...will remove this when we remove the PB handling
		pbTP := pb.TenantThresholdProfile{}
		if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
			return metrics_service.NewQueryThresholdCrossingInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert request to fetch %s: %s", datastore.QueryThresholdCrossingStr, err.Error()), startTime, http.StatusInternalServerError, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		if err = validateDomainsSwag(*params.Body.TenantID, params.Body.DomainIds, tenantDB); err != nil {
			return metrics_service.NewQueryThresholdCrossingNotFound().WithPayload(reportAPIError(fmt.Sprintf("Error looking up domains %s. Error: %s", models.AsJSONString(request.DomainIDs), err.Error()), startTime, http.StatusNotFound, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		result, err := druidDB.QueryThresholdCrossing(&request, &pbTP)
		if err != nil {
			return metrics_service.NewQueryThresholdCrossingInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve  Threshold Crossing Metrics. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.MetricResultsResponseObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewQueryThresholdCrossingInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.ThresholdCrossingStr, err.Error()), startTime, http.StatusInternalServerError, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Retrieved %s for Tenant %s using %s %s", datastore.ThresholdCrossingStr, request.TenantID, tenmod.TenantThresholdProfileStr, request.ThresholdProfileID)
		return metrics_service.NewQueryThresholdCrossingOK().WithPayload(&converted)
	}
}

// HandleGenSLAReport - fetch an internal SLA report
func HandleGenSLAReport(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) func(params metrics_service.GenSLAReportParams) middleware.Responder {
	return func(params metrics_service.GenSLAReportParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Issuing %s for Tenant %s using %s %s", datastore.SLAReportStr, params.Tenant, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewGenSLAReportForbidden().WithPayload(reportAPIError(fmt.Sprintf("Generate %s operation not authorized for role: %s", datastore.SLAReportStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GenSLAReportStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Turn the query Params into the request object:
		slaReportRequest := populateSLAReportRequestSwag(params)
		tenantID := slaReportRequest.TenantID

		thresholdProfile, err := tenantDB.GetTenantThresholdProfile(tenantID, slaReportRequest.ThresholdProfileID)
		if err != nil {
			return metrics_service.NewGenSLAReportNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", models.AsJSONString(slaReportRequest), err.Error()), startTime, http.StatusNotFound, mon.GenSLAReportStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		if err = validateDomainsSwag(slaReportRequest.TenantID, slaReportRequest.Domain, tenantDB); err != nil {
			return metrics_service.NewGenSLAReportNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", models.AsJSONString(slaReportRequest), err.Error()), startTime, http.StatusNotFound, mon.GenSLAReportStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert to PB type...will remove this when we remove the PB handling
		pbTP := pb.TenantThresholdProfile{}
		if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
			return metrics_service.NewGenSLAReportInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert request to fetch %s: %s", datastore.SLAReportStr, err.Error()), startTime, http.StatusInternalServerError, mon.GenSLAReportStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		result, err := druidDB.GetSLAReport(slaReportRequest, &pbTP)
		if err != nil {
			return metrics_service.NewGenSLAReportInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve SLA Report. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.GenSLAReportStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewGenSLAReportInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.SLAReportStr, err.Error()), startTime, http.StatusInternalServerError, mon.GenSLAReportStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetThrCrossStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Generated %s for Tenant %s using %s %s", datastore.SLAReportStr, params.Tenant, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)
		return metrics_service.NewGenSLAReportOK().WithPayload(&converted)
	}
}

// HandleGetThresholdCrossingByMonitoredObject - fetch all domains for a tenant
func HandleGetThresholdCrossingByMonitoredObject(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) func(params metrics_service.GetThresholdCrossingByMonitoredObjectParams) middleware.Responder {
	return func(params metrics_service.GetThresholdCrossingByMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Issuing %s for Tenant %s using %s %s", datastore.ThresholdCrossingByMonitoredObjectStr, params.Tenant, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Generate %s operation not authorized for role: %s", datastore.ThresholdCrossingByMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetThrCrossByMonObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Turn the query Params into the request object:
		thresholdCrossingReq := populateThresholdCrossingRequestForMonObjSwag(params)
		tenantID := thresholdCrossingReq.Tenant

		thresholdProfile, err := tenantDB.GetTenantThresholdProfile(tenantID, thresholdCrossingReq.ThresholdProfileId)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err.Error()), startTime, http.StatusNotFound, mon.GetThrCrossByMonObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert to PB type...will remove this when we remove the PB handling
		pbTP := pb.TenantThresholdProfile{}
		if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert request to fetch %s: %s", datastore.ThresholdCrossingByMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossByMonObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		if err = validateDomainsSwag(thresholdCrossingReq.Tenant, thresholdCrossingReq.Domain, tenantDB); err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", models.AsJSONString(thresholdCrossingReq.Domain), err.Error()), startTime, http.StatusNotFound, mon.GetThrCrossByMonObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		result, err := druidDB.GetThresholdCrossingByMonitoredObject(thresholdCrossingReq, &pbTP)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve Threshold Crossing By Monitored Object. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossByMonObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.ThresholdCrossingByMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossByMonObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetThrCrossByMonObjStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Generated %s for Tenant %s using %s %s", datastore.ThresholdCrossingByMonitoredObjectStr, params.Tenant, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)
		return metrics_service.NewGetThresholdCrossingByMonitoredObjectOK().WithPayload(&converted)
	}
}

// HandleGetThresholdCrossingByMonitoredObjectTopN - delete a domain for a tenant
func HandleGetThresholdCrossingByMonitoredObjectTopN(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) func(params metrics_service.GetThresholdCrossingByMonitoredObjectTopNParams) middleware.Responder {
	return func(params metrics_service.GetThresholdCrossingByMonitoredObjectTopNParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Issuing %s for Tenant %s using %s %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, params.TenantID, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectTopNForbidden().WithPayload(reportAPIError(fmt.Sprintf("Generate %s operation not authorized for role: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Turn the query Params into the request object:
		thresholdCrossingReq, _ := populateThresholdCrossingTopNRequestSwag(params)
		tenantID := thresholdCrossingReq.TenantID

		thresholdProfile, err := tenantDB.GetTenantThresholdProfile(tenantID, thresholdCrossingReq.ThresholdProfileID)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectTopNNotFound().WithPayload(reportAPIError(fmt.Sprintf("Generate %s operation not authorized for role: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusNotFound, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert to PB type...will remove this when we remove the PB handling
		pbTP := pb.TenantThresholdProfile{}
		if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectTopNInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert request to fetch %s: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		if err = validateDomainsSwag(thresholdCrossingReq.TenantID, thresholdCrossingReq.Domain, tenantDB); err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectTopNNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable find domain for given query parameters: %s. Error: %s", models.AsJSONString(thresholdCrossingReq.Domain), err.Error()), startTime, http.StatusNotFound, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		result, err := druidDB.GetThresholdCrossingByMonitoredObjectTopN(thresholdCrossingReq, &pbTP)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectTopNInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve Threshold Crossing By Monitored Object. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewGetThresholdCrossingByMonitoredObjectTopNInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Generated %s for Tenant %s using %s %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, params.TenantID, tenmod.TenantThresholdProfileStr, params.ThresholdProfileID)
		return metrics_service.NewGetThresholdCrossingByMonitoredObjectTopNOK().WithPayload(&converted)
	}
}

// HandleGetHistogram - get a metric histogramfor a tenant
func HandleGetHistogram(allowedRoles []string, druidDB datastore.DruidDatastore) func(params metrics_service.GetHistogramParams) middleware.Responder {
	return func(params metrics_service.GetHistogramParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Fetching %s for Tenant %s for Metric %s", datastore.HistogramStr, params.Tenant, params.Metric)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewGetHistogramForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", datastore.HistogramStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetHistogramObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		histogramReq := populateHistogramRequestSwag(params)

		result, err := druidDB.GetHistogram(histogramReq)
		if err != nil {
			return metrics_service.NewGetHistogramInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve Histogram. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.GetHistogramObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewGetHistogramInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.HistogramStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetHistogramObjStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetHistogramObjStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Generated %s for Tenant %s using metric %s", datastore.HistogramStr, params.Tenant, params.Metric)
		return metrics_service.NewGetHistogramOK().WithPayload(&converted)
	}
}

// HandleGetRawMetrics - get a metric histogramfor a tenant
func HandleGetRawMetrics(allowedRoles []string, druidDB datastore.DruidDatastore) func(params metrics_service.GetRawMetricsParams) middleware.Responder {
	return func(params metrics_service.GetRawMetricsParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Fetching %s for Tenant %s for Metric %s", datastore.RawMetricStr, params.Tenant, params.Metric)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewGetRawMetricsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", datastore.RawMetricStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetRawMetricStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Turn the query Params into the request object:
		rawMetricReq := populateRawMetricsRequestSwag(params)

		result, err := druidDB.GetRawMetrics(rawMetricReq)
		if err != nil {
			return metrics_service.NewGetRawMetricsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve Raw Metrics. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.GetRawMetricStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewGetRawMetricsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.RawMetricStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetRawMetricStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetRawMetricStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Generated %s for Tenant %s using metric %s", datastore.RawMetricStr, params.Tenant, params.Metric)
		return metrics_service.NewGetRawMetricsOK().WithPayload(&converted)
	}
}

// HandleQueryAggregatedMetrics -
func HandleQueryAggregatedMetrics(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) func(params metrics_service.QueryAggregatedMetricsParams) middleware.Responder {
	return func(params metrics_service.QueryAggregatedMetricsParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Fetching %s for Tenant %s", datastore.AggMetricsStr, params.Body.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewQueryAggregatedMetricsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", datastore.AggMetricsStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return metrics_service.NewQueryAggregatedMetricsBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted))
		}
		request := metrics.AggregateMetricsAPIRequest{}
		if err := json.Unmarshal(requestBytes, &request); err != nil {
			return metrics_service.NewQueryAggregatedMetricsBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		if err = validateDomainsSwag(request.TenantID, request.DomainIDs, tenantDB); err != nil {
			return metrics_service.NewQueryAggregatedMetricsNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable find domain for given request: %s. Error: %s", models.AsJSONString(request), err.Error()), startTime, http.StatusNotFound, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		result, err := druidDB.GetAggregatedMetrics(&request)
		if err != nil {
			return metrics_service.NewQueryAggregatedMetricsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve Aggregated Metrics. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewQueryAggregatedMetricsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.AggMetricsStr, err.Error()), startTime, http.StatusInternalServerError, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Generated %s for Tenant %s", datastore.AggMetricsStr, params.Body.TenantID)
		return metrics_service.NewGetRawMetricsOK().WithPayload(&converted)
	}
}

// HandleGetTopNFor -
func HandleGetTopNFor(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, druidDB datastore.DruidDatastore) func(params metrics_service.GetTopNForMetricParams) middleware.Responder {
	return func(params metrics_service.GetTopNForMetricParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.MetricAPIRecieved)
		logger.Log.Infof("Fetching %s for Tenant %s", datastore.TopNForMetricString, params.Body.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return metrics_service.NewGetTopNForMetricForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", datastore.TopNForMetricString, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return metrics_service.NewGetTopNForMetricBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted))
		}
		request := metrics.TopNForMetric{}
		if err := json.Unmarshal(requestBytes, &request); err != nil {
			return metrics_service.NewGetTopNForMetricBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		if _, err = request.Validate(); err != nil {
			return metrics_service.NewGetTopNForMetricBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		topNreq := request
		logger.Log.Infof("Fetching data for TopN request: %+v", topNreq)

		if err = validateDomainsSwag(topNreq.TenantID, topNreq.Domains, tenantDB); err != nil {
			return metrics_service.NewGetTopNForMetricNotFound().WithPayload(reportAPIError(fmt.Sprintf("Unable find domain for given query parameters: %+v. Error: %s", topNreq, err.Error()), startTime, http.StatusNotFound, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		result, err := druidDB.GetTopNForMetric(&topNreq)
		if err != nil {
			return metrics_service.NewGetTopNForMetricInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve Top N response. %s:", err.Error()), startTime, http.StatusInternalServerError, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		// Convert the res to byte[]
		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return metrics_service.NewGetTopNForMetricInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", datastore.TopNForMetricString, err.Error()), startTime, http.StatusInternalServerError, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted)
		logger.Log.Infof("Generated %s for Tenant %s", datastore.TopNForMetricString, params.Body.TenantID)
		return metrics_service.NewGetRawMetricsOK().WithPayload(&converted)
	}
}
