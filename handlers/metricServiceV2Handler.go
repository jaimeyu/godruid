package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	metmod "github.com/accedian/adh-gather/models/metrics"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/metrics_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/getlantern/deepcopy"
	"github.com/go-openapi/runtime/middleware"
	"github.com/manyminds/api2go/jsonapi"
	uuid "github.com/satori/go.uuid"
)

// HandleGetThresholdCrossingByMonitoredObjectTopNV2 - Retrieves threshold profile based Top N for distinct monitored objects
func HandleGetThresholdCrossingByMonitoredObjectTopNV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore) func(params metrics_service_v2.GetThresholdCrossingByMonitoredObjectTopNV2Params) middleware.Responder {
	return func(params metrics_service_v2.GetThresholdCrossingByMonitoredObjectTopNV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetThresholdCrossingByMonitoredObjectTopNV2(allowedRoles, metricsDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return metrics_service_v2.NewGetThresholdCrossingByMonitoredObjectTopNV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return metrics_service_v2.NewGetThresholdCrossingByMonitoredObjectTopNV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return metrics_service_v2.NewGetThresholdCrossingByMonitoredObjectTopNV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return metrics_service_v2.NewGetThresholdCrossingByMonitoredObjectTopNV2NotFound().WithPayload(errorMessage)
		default:
			return metrics_service_v2.NewGetThresholdCrossingByMonitoredObjectTopNV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGenerateSLAReportV2 - Retrieves SLA report for the specified tenant with the provided parameters
func HandleGenerateSLAReportV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore) func(params metrics_service_v2.GenerateSLAReportV2Params) middleware.Responder {
	return func(params metrics_service_v2.GenerateSLAReportV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGenerateSLAReportFromParamsV2(allowedRoles, metricsDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return metrics_service_v2.NewGenerateSLAReportV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GenerateSLAReportStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return metrics_service_v2.NewGenerateSLAReportV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return metrics_service_v2.NewGenerateSLAReportV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return metrics_service_v2.NewGenerateSLAReportV2NotFound().WithPayload(errorMessage)
		default:
			return metrics_service_v2.NewGenerateSLAReportV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetRawMetricsV2 - Retrieves raw metrics for the specified tenant with the provided parameters
func HandleGetRawMetricsV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore) func(params metrics_service_v2.GetFilteredRawMetricsV2Params) middleware.Responder {
	return func(params metrics_service_v2.GetFilteredRawMetricsV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetRawMetricsV2(allowedRoles, metricsDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return metrics_service_v2.NewGetFilteredRawMetricsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetRawMetricStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return metrics_service_v2.NewGetFilteredRawMetricsV2Forbidden().WithPayload(errorMessage)
		default:
			return metrics_service_v2.NewGetFilteredRawMetricsV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetAggregateMetricsV2 - Retrieves metrics in aggregation for the specified tenant with the provided parameters
func HandleGetAggregateMetricsV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore) func(params metrics_service_v2.QueryAggregateMetricsV2Params) middleware.Responder {
	return func(params metrics_service_v2.QueryAggregateMetricsV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetAggregateMetricsV2(allowedRoles, metricsDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return metrics_service_v2.NewQueryAggregateMetricsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return metrics_service_v2.NewQueryAggregateMetricsV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return metrics_service_v2.NewQueryAggregateMetricsV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return metrics_service_v2.NewQueryAggregateMetricsV2NotFound().WithPayload(errorMessage)
		default:
			return metrics_service_v2.NewQueryAggregateMetricsV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetHistogramV2 - Retrieves a histogram for the specified tenant with the provided parameters
func HandleGetHistogramV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore) func(params metrics_service_v2.GetHistogramV2Params) middleware.Responder {
	return func(params metrics_service_v2.GetHistogramV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetHistogramV2(allowedRoles, metricsDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return metrics_service_v2.NewGetHistogramV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetHistogramObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return metrics_service_v2.NewGetHistogramV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return metrics_service_v2.NewGetHistogramV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return metrics_service_v2.NewGetHistogramV2NotFound().WithPayload(errorMessage)
		default:
			return metrics_service_v2.NewGetHistogramV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetTopNForMetricV2 - Retrieves the top n for the specified tenant with the provided parameters
func HandleGetTopNForMetricV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore) func(params metrics_service_v2.GetTopNForMetricV2Params) middleware.Responder {
	return func(params metrics_service_v2.GetTopNForMetricV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetTopNForMetricV2(allowedRoles, metricsDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return metrics_service_v2.NewGetTopNForMetricV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetTopNReqStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return metrics_service_v2.NewGetTopNForMetricV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return metrics_service_v2.NewGetTopNForMetricV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return metrics_service_v2.NewGetTopNForMetricV2NotFound().WithPayload(errorMessage)
		default:
			return metrics_service_v2.NewGetTopNForMetricV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetThresholdCrossingV2 - Retrieves the threshold crossings for the specified tenant with the provided parameters
func HandleGetThresholdCrossingV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore) func(params metrics_service_v2.QueryThresholdCrossingV2Params) middleware.Responder {
	return func(params metrics_service_v2.QueryThresholdCrossingV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetThresholdCrossingV2(allowedRoles, metricsDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return metrics_service_v2.NewQueryThresholdCrossingV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return metrics_service_v2.NewQueryThresholdCrossingV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return metrics_service_v2.NewQueryThresholdCrossingV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return metrics_service_v2.NewQueryThresholdCrossingV2NotFound().WithPayload(errorMessage)
		default:
			return metrics_service_v2.NewQueryThresholdCrossingV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

func doGetHistogramV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, params metrics_service_v2.GetHistogramV2Params) (time.Time, int, *swagmodels.JSONAPIHistogramResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Retrieving %s for %s %s", datastore.HistogramStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.MetricAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", datastore.HistogramStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest := metmod.Histogram{}
	err = jsonapi.Unmarshal(requestBytes, &daoRequest)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	logger.Log.Debugf("Retrieving monitored objects for %s associated with meta criteria: %v", datastore.HistogramStr, daoRequest.Meta)

	// Retrieve monitored objects associated with the metadata
	metaMOs, err := tenantDB.GetFilteredMonitoredObjectList(tenantID, daoRequest.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved the following monitored objects for %s request based on meta criteria %v: %v", datastore.HistogramStr, daoRequest.Meta, metaMOs)
	}

	daoRequest.TenantID = tenantID

	// Issue request to DAO Layer
	queryReport, queryKeySpec, err := metricsDB.GetHistogram(&daoRequest, metaMOs)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.HistogramStr, err.Error())
	}

	var rendered map[string]interface{}
	if daoRequest.Meta != nil {
		rendered, err = renderHistogramTimeseriesMetrics("histograms", uuid.NewV4().String(), params.Body.Data.Attributes, queryKeySpec, map[string]interface{}{"monitoredObjectIds": metaMOs}, queryReport)
	} else {
		rendered, err = renderHistogramTimeseriesMetrics("histograms", uuid.NewV4().String(), params.Body.Data.Attributes, queryKeySpec, nil, queryReport)
	}
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to render %s report: %s", datastore.HistogramStr, err.Error())
	}

	rr, err := json.Marshal(rendered)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.HistogramStr, err.Error())
	}

	converted := swagmodels.JSONAPIHistogramResponse{}
	err = json.Unmarshal(rr, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.HistogramStr, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %s %s", datastore.HistogramStr, models.AsJSONString(converted))
	} else {
		logger.Log.Infof("Retrieved %d entries for %s", len(queryReport), datastore.HistogramStr)
	}
	reportAPICompletionState(startTime, http.StatusOK, mon.GetHistogramObjStr, mon.APICompleted, mon.MetricAPICompleted)

	return startTime, http.StatusOK, &converted, nil
}

func doGetTopNForMetricV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, params metrics_service_v2.GetTopNForMetricV2Params) (time.Time, int, *swagmodels.JSONAPITopNForMetricResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Retrieving %s for %s %s", datastore.TopNForMetricStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.MetricAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", datastore.TopNForMetricStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest := metmod.TopNForMetric{}
	err = jsonapi.Unmarshal(requestBytes, &daoRequest)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	if daoRequest.Meta != nil && daoRequest.MonitoredObjects != nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Request for %s cannot contain both meta filter and monitored object Id filter", datastore.TopNForMetricStr)
	}

	logger.Log.Debugf("Retrieving monitored objects for %s associated with meta criteria: %v", datastore.TopNForMetricStr, daoRequest.Meta)

	// Retrieve monitored objects associated with the metadata
	metaMOs, err := tenantDB.GetFilteredMonitoredObjectList(tenantID, daoRequest.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved the following monitored objects for %s request based on meta criteria %v: %v", datastore.TopNForMetricStr, daoRequest.Meta, metaMOs)
	}

	daoRequest.TenantID = tenantID

	// Issue request to DAO Layer
	queryReport, err := metricsDB.GetTopNForMetric(&daoRequest, metaMOs)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.TopNForMetricStr, err.Error())
	}

	swagMid := params.Body.Data.Attributes.Metric
	mID := metmod.MetricIdentifierFilter{Vendor: *swagMid.Vendor,
		ObjectType: swagMid.ObjectType,
		Metric:     *swagMid.Metric,
		Direction:  swagMid.Direction}

	rendered, err := renderTopNMetrics(params.Body.Data.Attributes, mID, queryReport, uuid.NewV4().String(), "topNForMetrics")
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to render %s report: %s", datastore.TopNForMetricStr, err.Error())
	}

	rr, err := json.Marshal(rendered)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.TopNForMetricStr, err.Error())
	}

	converted := swagmodels.JSONAPITopNForMetricResponse{}
	err = json.Unmarshal(rr, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.TopNForMetricStr, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %s %s", datastore.TopNForMetricStr, models.AsJSONString(converted))
	} else {
		logger.Log.Infof("Retrieved %d entries for %s", len(queryReport), datastore.TopNForMetricStr)
	}
	reportAPICompletionState(startTime, http.StatusOK, mon.GetTopNReqStr, mon.APICompleted, mon.MetricAPICompleted)

	return startTime, http.StatusOK, &converted, nil
}

func doGetThresholdCrossingV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, params metrics_service_v2.QueryThresholdCrossingV2Params) (time.Time, int, *swagmodels.JSONAPIThresholdCrossingResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Retrieving %s for %s %s", datastore.ThresholdCrossingStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.MetricAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", datastore.ThresholdCrossingStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest := metmod.ThresholdCrossing{}
	err = jsonapi.Unmarshal(requestBytes, &daoRequest)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest.TenantID = tenantID

	logger.Log.Debugf("Retrieving threshold profile for %s with id %s and tenant %s", datastore.ThresholdCrossingStr, daoRequest.ThresholdProfileID, tenantID)
	// Retrieve threshold profile associated with the tenant
	thresholdProfile, err := tenantDB.GetTenantThresholdProfile(tenantID, daoRequest.ThresholdProfileID)
	if err != nil {
		return startTime, http.StatusNotFound, nil, err
	}
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	logger.Log.Debugf("Retrieving monitored objects for %s associated with meta criteria: %v", datastore.ThresholdCrossingStr, daoRequest.Meta)
	// Retrieve monitored objects associated with the metadata
	metaMOs, err := tenantDB.GetFilteredMonitoredObjectList(tenantID, daoRequest.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved the following monitored objects for %s request based on meta criteria %v: %v", datastore.ThresholdCrossingStr, daoRequest.Meta, metaMOs)
	}

	// Issue request to DAO Layer
	queryReport, err := metricsDB.QueryThresholdCrossing(&daoRequest, &pbTP, metaMOs)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.ThresholdCrossingStr, err.Error())
	}

	// TODO REMOVE THE RE-RENDER ONCE THE QUERY PROPERLY RETURNS THE VIOLATIONS AGAINS THE V2 VERSION OF METRIC IDENTIFIER
	rendered := rerenderThresholdCrossingV2(daoRequest.Metrics, renderThresholdCrossingV2(params.Body.Data.Attributes, queryReport))

	rr, err := json.Marshal(rendered)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.ThresholdCrossingStr, err.Error())
	}

	converted := swagmodels.JSONAPIThresholdCrossingResponse{}
	err = json.Unmarshal(rr, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.ThresholdCrossingStr, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %s %s", datastore.ThresholdCrossingStr, models.AsJSONString(converted))
	} else {
		logger.Log.Infof("Retrieved %d entries for %s", len(queryReport), datastore.ThresholdCrossingStr)
	}
	reportAPICompletionState(startTime, http.StatusOK, mon.QueryThresholdCrossingStr, mon.APICompleted, mon.MetricAPICompleted)

	return startTime, http.StatusOK, &converted, nil
}

// This is a temporary method that is only used to re-group the violations under the v2 version of the metric identifier.
// We are essentially doing a post aggregation that we should be doing in druid. Once the underlying druid query is modified to
// properly build the query by the v2 version of the metric identifier this should go away
func rerenderThresholdCrossingV2(metricFilters []metmod.MetricIdentifierFilter, report map[string]interface{}) map[string]interface{} {
	results := report["data"].(map[string]interface{})["attributes"].(map[string]interface{})["result"].(map[string]interface{})

	if len(results) == 0 {
		return report
	}

	retrievedMetrics := results["metric"].([]*metmod.MetricViolationsTimeSeries)

	metricMapCache := make(map[string]*metmod.MetricViolationsTimeSeries)

	for _, m := range retrievedMetrics {
		metricIdentifierName := fmt.Sprintf("%s.%s.%s.%s", m.Vendor, m.ObjectType, m.Direction, m.Metric)
		metricMapCache[metricIdentifierName] = m
	}

	txfEntries := make([]map[string]interface{}, 0)

	for _, m := range metricFilters {
		eventMap := make(map[string][]metmod.MetricViolationSummaryType)

		for _, o := range m.ObjectType {
			for _, d := range m.Direction {
				metricIdentifierName := fmt.Sprintf("%s.%s.%s.%s", m.Vendor, o, d, m.Metric)

				if cachedMetricData, ok := metricMapCache[metricIdentifierName]; ok {
					if len(cachedMetricData.Critical) != 0 {
						event, ok := eventMap["critical"]
						if !ok {
							event = make([]metmod.MetricViolationSummaryType, 0)
						}
						event = mergeCategory(event, cachedMetricData.Critical)
						eventMap["critical"] = event
					}
					if len(cachedMetricData.Major) != 0 {
						event, ok := eventMap["major"]
						if !ok {
							event = make([]metmod.MetricViolationSummaryType, 0)
						}
						event = mergeCategory(event, cachedMetricData.Major)
						eventMap["major"] = event
					}
					if len(cachedMetricData.Minor) != 0 {
						event, ok := eventMap["minor"]
						if !ok {
							event = make([]metmod.MetricViolationSummaryType, 0)
						}
						event = mergeCategory(event, cachedMetricData.Minor)
						eventMap["minor"] = event
					}
					if len(cachedMetricData.Warning) != 0 {
						event, ok := eventMap["warning"]
						if !ok {
							event = make([]metmod.MetricViolationSummaryType, 0)
						}
						event = mergeCategory(event, cachedMetricData.Warning)
						eventMap["warning"] = event
					}
					if len(cachedMetricData.SLA) != 0 {
						event, ok := eventMap["sla"]
						if !ok {
							event = make([]metmod.MetricViolationSummaryType, 0)
						}
						event = mergeCategory(event, cachedMetricData.SLA)
						eventMap["sla"] = event
					}
				}

			}
		}

		finishedMetricFilter := make(map[string]interface{})
		finishedMetricFilter["vendor"] = m.Vendor
		finishedMetricFilter["objectType"] = m.ObjectType
		finishedMetricFilter["direction"] = m.Direction
		finishedMetricFilter["metric"] = m.Metric

		for k, v := range eventMap {
			finishedMetricFilter[k] = v
		}

		txfEntries = append(txfEntries, finishedMetricFilter)
	}

	delete(results, "metric")
	results["metric"] = txfEntries
	fmt.Printf("%v", report)

	return report
}

func mergeCategory(mergee []metmod.MetricViolationSummaryType, merger []*metmod.MetricViolationSummaryType) []metmod.MetricViolationSummaryType {

	addedEntry := false

	for _, from := range merger {
		found := false
		for _, to := range mergee {
			if to["timestamp"] == (*from)["timestamp"] {
				mergeTimeseriesEntry(to, (*from))
				found = true
			}
		}
		if !found {
			mergee = append(mergee, (*from))
			addedEntry = true // Optimization so we only sort everything if we know that we added a new entry to the tail of the event order
		}
	}

	if addedEntry {
		sortEntries(mergee)
	}

	return mergee
}

func sortEntries(mergee []metmod.MetricViolationSummaryType) {
	sort.Slice(mergee, func(i, j int) bool {
		t1, err := time.Parse(time.RFC3339, mergee[i]["timestamp"].(string))
		if err != nil {
			return false
		}
		t2, err := time.Parse(time.RFC3339, mergee[j]["timestamp"].(string))
		if err != nil {
			return false
		}
		return t1.Before(t2)
	})
}

func mergeTimeseriesEntry(mergee map[string]interface{}, merger map[string]interface{}) {

	if mergee["timestamp"] != merger["timestamp"] {
		return
	}

	for k, v := range merger {
		if k == "timestamp" {
			continue
		}

		if val, ok := mergee[k]; ok {
			mergee[k] = val.(float64) + v.(float64)
		}
	}
}

func renderThresholdCrossingV2(config interface{}, report map[string]interface{}) map[string]interface{} {

	topLevelResult := report["results"]

	metrics := topLevelResult.([]metmod.ThresholdCrossingTimeSeriesEntry)

	listThresholdCrossings := make(map[string]*metmod.MetricViolationsTimeSeries)

	// var totals []map[string]interface{}
	totals := make(map[string][]*metmod.MetricViolationSummaryType)
	totals["summary"] = make([]*metmod.MetricViolationSummaryType, 0)

	for _, tsEntry := range metrics {
		ts := tsEntry.Timestamp
		results := tsEntry.Result

		if (float64(results.TotalDuration) == 0) &&
			(float64(results.TotalViolationCount) == 0) &&
			(float64(results.TotalViolationDuration) == 0) {

		} else {

			totals["summary"] = append(totals["summary"], &metmod.MetricViolationSummaryType{
				"timestamp":              ts,
				"totalDuration":          float64(results.TotalDuration),
				"totalViolationCount":    float64(results.TotalViolationCount),
				"totalViolationDuration": float64(results.TotalViolationDuration),
			})
		}

		totals["critical"] = addValToList(totals["critical"], results.BySeverity["critical"], ts)

		totals["major"] = addValToList(totals["major"], results.BySeverity["major"], ts)

		totals["minor"] = addValToList(totals["minor"], results.BySeverity["minor"], ts)

		totals["warning"] = addValToList(totals["warning"], results.BySeverity["warning"], ts)

		totals["sla"] = addValToList(totals["sla"], results.BySeverity["sla"], ts)

		for _, metric := range results.ByMetric {

			logger.Log.Debugf("metric: %s", models.AsJSONString(metric))
			hash := buildHash("metric", metric.Metric, metric.ObjectType, metric.Vendor, metric.Direction)
			c := listThresholdCrossings[hash]
			if c == nil {
				c = &metmod.MetricViolationsTimeSeries{
					Direction:  metric.Direction,
					Metric:     metric.Metric,
					Vendor:     metric.Vendor,
					ObjectType: metric.ObjectType,
				}

			}

			// Need a null check for BySeverity[xxxxx]
			if metric.BySeverity["sla"] != nil {
				c.SLA = addValToList(c.SLA, metric.BySeverity["sla"], ts)
			}

			if metric.BySeverity["major"] != nil {
				c.Major = addValToList(c.Major, metric.BySeverity["major"], ts)
			}
			if metric.BySeverity["minor"] != nil {
				c.Minor = addValToList(c.Minor, metric.BySeverity["minor"], ts)
			}
			if metric.BySeverity["warning"] != nil {
				c.Warning = addValToList(c.Warning, metric.BySeverity["warning"], ts)
			}
			if metric.BySeverity["critical"] != nil {
				c.Critical = addValToList(c.Critical, metric.BySeverity["critical"], ts)
			}

			listThresholdCrossings[hash] = c

		}

	}

	swaggerResponse := make(map[string]interface{})
	swaggerResponse["summary"] = totals["summary"]

	swaggerMetrics := make([]*metmod.MetricViolationsTimeSeries, 0)
	for _, metric := range listThresholdCrossings {
		swaggerMetrics = append(swaggerMetrics, metric)
	}
	swaggerResponse["metric"] = swaggerMetrics

	rawResponse := make(map[string]interface{})
	rawResponse["config"] = config
	rawResponse["result"] = swaggerResponse

	return wrapJsonAPIObject(rawResponse, uuid.NewV4().String(), "thresholdCrossings")
}

func doGetAggregateMetricsV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, params metrics_service_v2.QueryAggregateMetricsV2Params) (time.Time, int, *swagmodels.JSONAPIAggregationResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Retrieving %s for %s %s", datastore.AggMetricsStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.MetricAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", datastore.AggMetricsStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest := metmod.AggregateMetrics{}
	err = jsonapi.Unmarshal(requestBytes, &daoRequest)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	if daoRequest.Meta != nil && daoRequest.MonitoredObjects != nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Request for %s cannot contain both meta filter and monitored object Id filter", datastore.AggMetricsStr)
	}

	logger.Log.Debugf("Retrieving monitored objects for %s associated with meta criteria: %v", datastore.AggMetricsStr, daoRequest.Meta)

	// Retrieve monitored objects associated with the metadata
	metaMOs, err := tenantDB.GetFilteredMonitoredObjectList(tenantID, daoRequest.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved the following monitored objects for %s request based on meta criteria %v: %v", datastore.AggMetricsStr, daoRequest.Meta, metaMOs)
	}

	daoRequest.TenantID = tenantID

	// Issue request to DAO Layer
	queryReport, queryKeySpec, err := metricsDB.GetAggregatedMetrics(&daoRequest, metaMOs)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.AggMetricsStr, err.Error())
	}

	var rendered map[string]interface{}
	if daoRequest.Meta != nil {
		rendered, err = renderTimeseriesMetrics("aggregateMetrics", uuid.NewV4().String(), params.Body.Data.Attributes, queryKeySpec, map[string]interface{}{"monitoredObjectIds": metaMOs}, queryReport)
	} else {
		rendered, err = renderTimeseriesMetrics("aggregateMetrics", uuid.NewV4().String(), params.Body.Data.Attributes, queryKeySpec, nil, queryReport)
	}
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to render timeseries metrics for %s: %s", datastore.AggMetricsStr, err.Error())
	}

	rr, err := json.Marshal(rendered)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.AggMetricsStr, err.Error())
	}

	converted := swagmodels.JSONAPIAggregationResponse{}
	err = json.Unmarshal(rr, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.AggMetricsStr, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %s %s", datastore.AggMetricsStr, models.AsJSONString(converted))
	} else {
		logger.Log.Infof("Retrieved %d entries for %s", len(queryReport), datastore.AggMetricsStr)
	}
	reportAPICompletionState(startTime, http.StatusOK, mon.QueryAggregatedMetricsStr, mon.APICompleted, mon.MetricAPICompleted)

	return startTime, http.StatusOK, &converted, nil
}

// DEPRECATE - THIS IS ONLY KEPT IN ORDER TO NOT DISRUPT COLTS USE OF THE API. REMOVE AS SOON AS WE CAN GET THEM ONTO V2 AGGREGATE
func doGetRawMetricsV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, params metrics_service_v2.GetFilteredRawMetricsV2Params) (time.Time, int, map[string]interface{}, error) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	// Turn the query Params into the request object:
	request := &metmod.RawMetrics{}
	err = json.Unmarshal(requestBytes, request)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.RawMetricStr, err.Error())
	}

	logger.Log.Infof("Retrieving %s for: %v", datastore.RawMetricStr, request)

	var metaMOs []string

	if len(request.Meta) != 0 {
		logger.Log.Debugf("Retrieving monitored objects by meta data for request: %v", request)
		metaMOs, err = tenantDB.GetFilteredMonitoredObjectList(request.TenantID, request.Meta)
	} else {
		logger.Log.Debugf("Retrieving all monitored objects for request: %v", request)
		metaMOs, err = tenantDB.GetAllMonitoredObjectsIDs(request.TenantID)
	}

	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve monitored object list for meta data for %s request: %s", datastore.RawMetricStr, err.Error())
	}

	result, err := metricsDB.GetFilteredRawMetrics(request, metaMOs)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.RawMetricStr, err.Error())
	}

	// Convert the res to byte[]
	res, err := json.Marshal(result)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to marshal %s response: %s", datastore.RawMetricStr, err.Error())
	}

	converted := make(map[string]interface{})
	err = json.Unmarshal(res, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to generate %s response: %s", datastore.RawMetricStr, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %s %s", datastore.RawMetricStr, string(res))
	} else {
		logger.Log.Infof("Retrieved %d entries for %s", len(result), datastore.RawMetricStr)
	}
	reportAPICompletionState(startTime, http.StatusOK, mon.GetRawMetricStr, mon.APICompleted, mon.MetricAPICompleted)

	return startTime, http.StatusOK, converted, nil
}

func renderTopNByBuckets(raw [][]byte, schema metmod.DruidViolationsMap) (*metmod.MetricViolationsAsTimeSeries, error) {

	// Render vars
	prerender := make(metmod.DruidResponse2TimeSeriesMap)
	for _, resp := range raw {
		var testmap metmod.DruidTopNResponse
		err := json.Unmarshal(resp, &testmap)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshal :%s", err.Error())
		}
		if len(testmap) == 0 {
			continue
		}
		if len(testmap[0].Result) == 0 {
			continue
		}

		for i, ia := range testmap[0].Result {

			for k, val := range ia {

				if schema[k] == nil {
					continue
				}

				prerender.Put(schema[k].Name, fmt.Sprintf("%d", i), schema[k], val)

			}
		}
	}
	// flatten the maps
	var response metmod.MetricViolationsAsTimeSeries

	response.PerMetricResult = prerender //append(response.PerMetricResult, *metricRes)

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Rendered:%s", models.AsJSONString(response))
	}
	return &response, nil

}
func getSLATopNByBuckets(druidDB datastore.MetricsDatastore, request *metmod.SLAReportRequest, granularity int, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string, sla bool) ([][]byte, metmod.DruidViolationsMap, error) {

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}
	// responseSchemaMap := make(metrics.DruidViolationsMap)
	var responses [][]byte
	responsSchemas := make(metmod.DruidViolationsMap)

	for vk, v := range thresholdProfile.Data.GetThresholds().GetVendorMap() {
		for tk, t := range v.GetMonitoredObjectTypeMap() {

			for mk, m := range t.GetMetricMap() {
				for dk, d := range m.GetDirectionMap() {
					for ek, e := range d.GetEventMap() {
						if ek != "sla" {
							continue
						}

						resp, schema, err := druidDB.GetTopNTimeByBuckets(request, granularity,
							vk, tk, mk, dk, "sla", e, metaMOs)
						if err != nil {
							return nil, nil, fmt.Errorf("Issue getting time buckets:%s", err.Error())
						}

						responses = append(responses, resp)
						responsSchemas.Merge(schema)
					}
				}
			}
		}
	}
	return responses, responsSchemas, nil
}

// DoGenerateSLAReportV2 - Expose this API so the scheduler can use it.
func DoGenerateSLAReportV2(druidDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, daoRequest metmod.SLAReportRequest) (time.Time, int, map[string]interface{}, error) {
	tenantID := daoRequest.TenantID
	startTime := time.Now()

	logger.Log.Debugf("Retrieving monitored objects for %s associated with meta criteria: %v", datastore.SLAReportStr, daoRequest.Meta)

	// Retrieve monitored objects associated with the metadata
	metaMOs, err := tenantDB.GetFilteredMonitoredObjectList(tenantID, daoRequest.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved the following monitored objects for %s request based on meta criteria %v: %v", datastore.SLAReportStr, daoRequest.Meta, metaMOs)
	}

	daoRequest.TenantID = tenantID

	// Retrieve threshold profile associated with the tenant
	thresholdProfile, err := tenantDB.GetTenantThresholdProfile(tenantID, daoRequest.ThresholdProfileID)
	if err != nil {
		return startTime, http.StatusNotFound, nil, fmt.Errorf("Could not get ThresholdProfile (%s) for tenant ID %s", daoRequest.ThresholdProfileID, tenantID)
	}
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	/**********************/

	/* Generate the SLA Report here */

	/* Broken up sla report */
	slaReportRequest := daoRequest

	/*** SLA violations per hour & day of week buckets ***/

	dayOfWeekBucket, dayOfWeekSchema, err := getSLATopNByBuckets(druidDB, &slaReportRequest, datastore.DayOfWeek, &pbTP, metaMOs, true)
	dayOfWeekRendered, err := renderTopNByBuckets(dayOfWeekBucket, dayOfWeekSchema)
	if err != nil {
		logger.Log.Errorf("Could not render day of week:%s", err.Error())
	}

	hourOfDayBucket, hourOfDaySchema, err := getSLATopNByBuckets(druidDB, &slaReportRequest, datastore.HourOfDay, &pbTP, metaMOs, true)
	hourOfDayRendered, err := renderTopNByBuckets(hourOfDayBucket, hourOfDaySchema)
	if err != nil {
		logger.Log.Errorf("Could not render day of week:%s", err.Error())
	}

	/*** SLA Violations per time bucket ***/
	slaViolationsGranular, schemaGranular, err := druidDB.GetSLAViolationsQueryWithGranularity(&slaReportRequest, &pbTP, metaMOs)

	slaViolationsGranularSummary, err := ReformatGetSLAViolationsQueryWithGranularityV2(slaViolationsGranular, schemaGranular, false)

	/*** SLA Violations summary ****/
	var slaReportRequestWithGranularityAll metmod.SLAReportRequest
	err = deepcopy.Copy(&slaReportRequestWithGranularityAll, &slaReportRequest)
	if err != nil {
		logger.Log.Errorf("Could not copy request: %s", err.Error())
	}
	slaReportRequestWithGranularityAll.Granularity = "all"
	slaViolationsAll, slaViolationsAllSchema, err := druidDB.GetSLAViolationsQueryWithGranularity(&slaReportRequestWithGranularityAll, &pbTP, metaMOs)

	slaViolationsAllSummary, err := ReformatGetSLAViolationsQueryAllGranularityV2(slaViolationsAll, slaViolationsAllSchema, true)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	// Convert from map stucture to an array
	aggregatedResponse := make(metmod.DruidResponse2TimeSeriesMap)
	for key, val := range slaViolationsAllSummary.PerMetricResult {
		if aggregatedResponse[key] == nil {
			aggregatedResponse[key] = val
		}
		totals := utilMetricViolationSummaryTypeMap2Array(val.InternalSeries)
		aggregatedResponse[key].Totals = totals[0]
		aggregatedResponse[key].InternalSeries = nil
	}

	for key, val := range slaViolationsGranularSummary.PerMetricResult {
		if aggregatedResponse[key] == nil {
			aggregatedResponse[key] = val
		}
		aggregatedResponse[key].ByGranularity = utilMetricViolationSummaryTypeMap2Array(val.InternalSeries)
		aggregatedResponse[key].InternalSeries = nil
	}

	for key, val := range hourOfDayRendered.PerMetricResult {
		if aggregatedResponse[key] == nil {
			aggregatedResponse[key] = val
		}
		aggregatedResponse[key].ByHourPerDay = utilMetricViolationSummaryTypeMap2Array(val.InternalSeries)
		aggregatedResponse[key].InternalSeries = nil
	}

	for key, val := range dayOfWeekRendered.PerMetricResult {
		if aggregatedResponse[key] == nil {
			aggregatedResponse[key] = val
		}
		aggregatedResponse[key].ByDayPerWeek = utilMetricViolationSummaryTypeMap2Array(val.InternalSeries)
		aggregatedResponse[key].InternalSeries = nil
	}

	// Special hack that I want to remove in V3. The front end should be able to do this calculatiohn
	// It does a SLA Percent count
	if slaViolationsAllSummary.Summary == nil {
		slaViolationsAllSummary.Summary = make(metmod.MetricViolationSummaryType)
	}
	slaViolationsAllSummary.Summary["slaCompliancePercent"] = 0
	if (slaViolationsAllSummary.Summary["totalDuration"] != nil) &&
		(slaViolationsAllSummary.Summary["totalViolationDuration"] != nil) {

		slaViolationsAllSummary.Summary["slaCompliancePercent"] = (slaViolationsAllSummary.Summary["totalDuration"].(float64) - slaViolationsAllSummary.Summary["totalViolationDuration"].(float64)) / slaViolationsAllSummary.Summary["totalDuration"].(float64) * 100
	}

	// Granular data
	slaViolationsAllSummary.Summary["byGranularity"] = utilMetricViolationsSummaryAsTimeSeriesEntryMap2Array(slaViolationsGranularSummary.SummaryResult)
	// Convert to universal response structure
	slaReport := metmod.SLAReportV2{
		ID:     uuid.NewV4().String(),
		Config: daoRequest,
		Result: metmod.SLAReportV2Result{
			Summary: slaViolationsAllSummary.Summary,
			Metric:  aggregatedResponse.ToArray(),
		},
	}

	logger.Log.Debugf("    Result: %s ", models.AsJSONString(slaReport))

	// Wrap for JSON API
	jsonapi := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         slaReport.ID,
			"type":       slaReport.GetName(),
			"attributes": slaReport}}

	return startTime, http.StatusOK, jsonapi, nil

}

func convertSLAMapToJSONAPI(input map[string]interface{}) (*swagmodels.JSONAPISLAReportResponse, error) {

	rr, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert %s data to map: %s", datastore.SLAReportStr, err.Error())
	}

	converted := swagmodels.JSONAPISLAReportResponse{}
	err = json.Unmarshal(rr, &converted)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.SLAReportStr, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Infof("Retrieved %s %s", datastore.SLAReportStr, models.AsJSONString(converted))
	}

	return &converted, nil
}

func doGenerateSLAReportFromParamsV2(allowedRoles []string, druidDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, params metrics_service_v2.GenerateSLAReportV2Params) (time.Time, int, *swagmodels.JSONAPISLAReportResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Retrieving %s for %s %s", datastore.SLAReportStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.MetricAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", datastore.SLAReportStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest := metmod.SLAReportRequest{}
	err = jsonapi.Unmarshal(requestBytes, &daoRequest)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest.TenantID = tenantID

	startTime, responseCode, mapresponse, err := DoGenerateSLAReportV2(druidDB, tenantDB, daoRequest)
	if err != nil {
		return startTime, responseCode, nil, err
	}

	response, err := convertSLAMapToJSONAPI(mapresponse)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GenerateSLAReportStr, mon.APICompleted, mon.MetricAPICompleted)

	return startTime, http.StatusOK, response, nil
}

func renderSLAReportV2(config interface{}, report interface{}) map[string]interface{} {

	rawResponse := make(map[string]interface{})
	rawResponse["config"] = config
	rawResponse["result"] = report

	return wrapJsonAPIObject(rawResponse, uuid.NewV4().String(), "slaReports")
}

func doGetThresholdCrossingByMonitoredObjectTopNV2(allowedRoles []string, metricsDB datastore.MetricsDatastore, tenantDB datastore.TenantServiceDatastore, params metrics_service_v2.GetThresholdCrossingByMonitoredObjectTopNV2Params) (time.Time, int, *swagmodels.JSONAPIThresholdCrossingByMOTopNResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Retrieving %s for %s %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.MetricAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	daoRequest := metmod.ThresholdCrossingTopN{}
	err = jsonapi.Unmarshal(requestBytes, &daoRequest)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	logger.Log.Debugf("Retrieving monitored objects for %s associated with meta criteria: %v", datastore.TopNThresholdCrossingByMonitoredObjectStr, daoRequest.Meta)

	// Retrieve monitored objects associated with the metadata
	metaMOs, err := tenantDB.GetFilteredMonitoredObjectList(tenantID, daoRequest.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved the following monitored objects for %s request based on meta criteria %v: %v", datastore.TopNThresholdCrossingByMonitoredObjectStr, daoRequest.Meta, metaMOs)
	}

	daoRequest.TenantID = tenantID
	daoRequest.Granularity = "all" // For the v2 queries we only care about a single bucket for the query

	logger.Log.Debugf("Retrieving threshold profile for %s with id %s and tenant %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, daoRequest.ThresholdProfileID, tenantID)
	// Retrieve threshold profile associated with the tenant
	thresholdProfile, err := tenantDB.GetTenantThresholdProfile(tenantID, daoRequest.ThresholdProfileID)
	if err != nil {
		return startTime, http.StatusNotFound, nil, err
	}
	pbTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(thresholdProfile, &pbTP); err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	// Issue request to DAO Layer
	queryReport, err := metricsDB.GetThresholdCrossingByMonitoredObjectTopN(&daoRequest, &pbTP, metaMOs)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, err.Error())
	}

	rendered, err := renderTopNMetrics(params.Body.Data.Attributes, daoRequest.Metric, queryReport, uuid.NewV4().String(), "thresholdCrossingByMOTopNs")
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to render %s report: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, err.Error())
	}

	rr, err := json.Marshal(rendered)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, err.Error())
	}

	converted := swagmodels.JSONAPIThresholdCrossingByMOTopNResponse{}
	err = json.Unmarshal(rr, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Retrieved %s %s", datastore.TopNThresholdCrossingByMonitoredObjectStr, models.AsJSONString(converted))
	} else {
		logger.Log.Infof("Retrieved %d entries for %s", len(queryReport), datastore.TopNThresholdCrossingByMonitoredObjectStr)
	}
	reportAPICompletionState(startTime, http.StatusOK, mon.GetThrCrossByMonObjTopNStr, mon.APICompleted, mon.MetricAPICompleted)

	return startTime, http.StatusOK, &converted, nil
}

func renderTimeseriesMetrics(reportType string, reportID string, config interface{}, queryKeySpec *datastore.QueryKeySpec, staticKeyEntries map[string]interface{}, reportEntries []metmod.TimeseriesEntryResponse) (map[string]interface{}, error) {

	metricIdentifierMap := make(map[string][]map[string]interface{})

	if reportEntries != nil {
		for _, rEntry := range reportEntries {

			rTimestamp := rEntry.Timestamp
			for accessorKey, v := range rEntry.Result {
				hasData := false

				// Initialize an empty map if one does not exist for the current metric identifier
				if _, ok := metricIdentifierMap[accessorKey]; !ok {
					metricIdentifierMap[accessorKey] = make([]map[string]interface{}, 0)
				}

				var value float64
				switch v.(type) {
				case float32:
					hasData = true
					value = float64(v.(float32))
				case float64:
					hasData = true
					value = v.(float64)
				case int:
					hasData = true
					value = float64(v.(int))
				case string:
					hasData = false
				default:
					hasData = true
				}
				if hasData {
					metricIdentifierTimeseries := metricIdentifierMap[accessorKey]
					metricIdentifierTimeseries = append(metricIdentifierTimeseries, map[string]interface{}{"timestamp": rTimestamp, "value": value})
					metricIdentifierMap[accessorKey] = metricIdentifierTimeseries
				}
			}
		}
	}

	reportResponse := make([]map[string]interface{}, len(metricIdentifierMap))

	i := 0
	for accessor, series := range metricIdentifierMap {

		keyMap := queryKeySpec.KeySpecMap[accessor]

		metricMap := make(map[string]interface{})
		for qk, qv := range keyMap {
			metricMap[qk] = qv
		}

		// Loop over the static key entries that must be added to each key entry in the key response
		for sk, sv := range staticKeyEntries {
			metricMap[sk] = sv
		}

		metricMap["series"] = series
		reportResponse[i] = metricMap
		i++
	}
	return renderV2Report(config, reportResponse, reportID, reportType)
}

func renderHistogramTimeseriesMetrics(reportType string, reportID string, config interface{}, queryKeySpec *datastore.QueryKeySpec, staticKeyEntries map[string]interface{}, reportEntries []metmod.TimeseriesEntryResponse) (map[string]interface{}, error) {

	type HistogramEntry struct {
		Timestamp string
		Buckets   map[string]interface{}
	}

	metricIdentifierMap := make(map[string][]HistogramEntry)

	if reportEntries != nil {
		for _, rEntry := range reportEntries {

			rTimestamp := rEntry.Timestamp
			for k, v := range rEntry.Result {
				// Expecting a key structure with the query spec ID with the order suffixed to it
				parts := strings.Split(k, datastore.QueryDelimeter)

				// The index should not be part of the key as the index is only used to preserve the order of the histogram response
				accessorKey := parts[0]
				index := parts[1]

				// Initialize an empty map if one does not exist for the current metric identifier
				if _, ok := metricIdentifierMap[accessorKey]; !ok {
					metricIdentifierMap[accessorKey] = make([]HistogramEntry, 0)
				}

				metricIdentifierTimeseries := metricIdentifierMap[accessorKey]
				if len(metricIdentifierTimeseries) == 0 {
					// This is the first time we've seem this metric identifier so create an empty entry for it
					metricIdentifierTimeseries = append(metricIdentifierTimeseries, HistogramEntry{Timestamp: rTimestamp, Buckets: make(map[string]interface{})})
				}
				// Get the latest timeseries entry for the metric identifier. If it doesn't match that means we need to create a new entry since
				// we are in a different timeblock
				currentTimeseries := metricIdentifierTimeseries[len(metricIdentifierTimeseries)-1]
				if rTimestamp != currentTimeseries.Timestamp {
					currentTimeseries = HistogramEntry{Timestamp: rTimestamp, Buckets: make(map[string]interface{})}
					metricIdentifierTimeseries = append(metricIdentifierTimeseries, currentTimeseries)
				}

				currentTimeseries.Buckets[index] = v
				metricIdentifierMap[accessorKey] = metricIdentifierTimeseries
			}
		}
	}

	reportResponse := make([]map[string]interface{}, len(metricIdentifierMap))

	i := 0
	for accessor, series := range metricIdentifierMap {

		orderedSeries := make([]map[string]interface{}, 0)
		for _, he := range series {
			histogramEntryMap := make(map[string]interface{})
			orderedHistogramEntries := make([]interface{}, len(he.Buckets))
			// Here we use the index value that arrived in the druid response to properly order which the bucket the count belongs to
			for index, count := range he.Buckets {
				iconv, err := strconv.Atoi(index)
				if err != nil {
					return nil, err
				}
				orderedHistogramEntries[iconv] = count
			}
			histogramEntryMap["timestamp"] = he.Timestamp
			histogramEntryMap["values"] = orderedHistogramEntries

			orderedSeries = append(orderedSeries, histogramEntryMap)
		}

		keyMap := queryKeySpec.KeySpecMap[accessor]

		metricMap := make(map[string]interface{})
		for qk, qv := range keyMap {
			metricMap[qk] = qv
		}

		// Loop over the static key entries that must be added to each key entry in the key response
		for sk, sv := range staticKeyEntries {
			metricMap[sk] = sv
		}

		metricMap["series"] = orderedSeries
		reportResponse[i] = metricMap
		i++
	}
	return renderV2Report(config, reportResponse, reportID, reportType)
}

func renderTopNMetrics(config interface{}, metricIdentifier metmod.MetricIdentifierFilter, report []metmod.TopNEntryResponse, ID string, reportType string) (map[string]interface{}, error) {

	renderedReport := make([]map[string]interface{}, len(report))

	for i, r := range report {
		renderedReport[i] = map[string]interface{}{"monitoredObjectIds": []string{r.MonitoredObjectId},
			"vendor":     metricIdentifier.Vendor,
			"objectType": metricIdentifier.ObjectType,
			"metric":     metricIdentifier.Metric,
			"direction":  metricIdentifier.Direction,
			"result":     r.Result}
	}

	return renderV2Report(config, renderedReport, ID, reportType)
}

func renderV2Report(config interface{}, report interface{}, ID string, reportType string) (map[string]interface{}, error) {
	rawResponse, err := models.ConvertObj2Map(config)
	if err != nil {
		return nil, err
	}
	rawResponse["result"] = report

	return wrapJsonAPIObject(rawResponse, ID, reportType), nil
}

/* Transformation utilities */
// ReformatGetSLAViolationsQueryWithGranularityV2 - Converts SLA report to an interface
func ReformatGetSLAViolationsQueryWithGranularityV2(druidResponse []byte, schema metmod.DruidViolationsMap, isSLAMode bool) (*metmod.MetricViolationsAsTimeSeries, error) {
	// logger.Log.Debugf("Response from druid for  %s", string(druidResponse))

	v2Model := metmod.MetricViolationsAsTimeSeries{}
	if v2Model.SummaryResult == nil {
		v2Model.SummaryResult = make(map[string]*metmod.MetricViolationsSummaryAsTimeSeriesEntry)
	}

	entries := []*metmod.DruidTimeSeriesResponse{}
	if err := json.Unmarshal(druidResponse, &entries); err != nil {
		return nil, err
	}

	// Is this correct if there are no SLA violations, there will be no entries?
	if len(entries) < 1 {
		return &v2Model, nil
	}
	tmpViolations := make(metmod.DruidResponse2TimeSeriesMap)

	for _, entry := range entries {

		ts := entry.Timestamp
		// For a summary, we expect only 1 entry in the druid results so just use the first entry.
		for k, v := range entry.Result {

			if v != 0 {
				/* Check for schema based responses */

				schemaEntry := schema[k]

				if schemaEntry != nil {

					if schemaEntry.Type == "SLA_Summary" {
						cs := v2Model.SummaryResult[ts]

						if cs == nil {
							cs = &metmod.MetricViolationsSummaryAsTimeSeriesEntry{
								"Timestamp": ts,
							}
						}
						(*cs)[schemaEntry.Name] = v
						v2Model.SummaryResult[ts] = cs

						continue
					}
					// Add the value to the metric
					tmpViolations.Put(schemaEntry.Name, ts, schemaEntry, v)
					// Add the timestamp to the metric
					tsobj := tmpViolations[schemaEntry.Name].InternalSeries[ts]
					(*tsobj)["timestamp"] = &ts
				}
			}
		}
	}
	v2Model.PerMetricResult = tmpViolations

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("V2 v2Model Formatted result for %v", models.AsJSONString(v2Model))
	}
	return &v2Model, nil
}

// ReformatGetSLAViolationsQueryAllGranularityV2 - Converts SLA Violationn to a MetricViolationsAsSummary
func ReformatGetSLAViolationsQueryAllGranularityV2(druidResponse []byte, schema metmod.DruidViolationsMap, isSLAMode bool) (*metmod.MetricViolationsAsSummary, error) {
	// logger.Log.Debugf("Response from druid for %s", string(druidResponse))

	v2Model := metmod.MetricViolationsAsSummary{}

	entries := []*metmod.DruidTimeSeriesResponse{}
	if err := json.Unmarshal(druidResponse, &entries); err != nil {
		return nil, err
	}

	// Is this correct if there are no SLA violations, there will be no entries?
	if len(entries) < 1 {
		return &v2Model, nil
	}

	prerender := make(metmod.DruidResponse2TimeSeriesMap)

	// For a summary, we expect only 1 entry in the druid results so just use the first entry.
	for k, v := range entries[0].Result {

		/* Check for schema based responses */

		if v != 0 {
			schemaEntry := schema[k]
			if schemaEntry != nil {

				if schemaEntry.Type == "SLA_Summary" {
					if v2Model.Summary == nil {
						v2Model.Summary = metmod.MetricViolationSummaryType{}
					}
					v2Model.Summary[schemaEntry.Name] = v
					continue
				}

				// We only expect 1 entry per metric so we cheat here and set the subkey to "0"
				prerender.Put(schema[k].Name, "0", schemaEntry, v)

			}
		}

	}

	v2Model.PerMetricResult = prerender

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("V2 v2Model Formatted result for: %s", models.AsJSONString(v2Model))
	}
	return &v2Model, nil
}
