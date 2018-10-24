package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/models/common"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/manyminds/api2go/jsonapi"
)

type httpErrorString string

const (
	// Error search strings
	notFound httpErrorString = "status 404 - not found"
	conflict httpErrorString = "already exists"

	// Custom Error types
	errorMarshal = -100

	// API Prefix values
	apiV1Prefix      = "/api/v1/"
	tenantsAPIPrefix = "tenants/{tenantID}/"

	startKeyQueryParamStr   = "start_key"
	limitQueryParamStr      = "limit"
	descendingQueryParamStr = "descending"

	linkFirstStr = "first"
	linkLastStr  = "last"
	linkPrevStr  = "prev"
	linkNextStr  = "next"
	linkSelfStr  = "self"
)

var (
	defaultIngestionProfileMetricNames = []string{
		"delayMax", "delayP95", "delayPHi", "delayVarP95", "delayVarPHi",
		"jitterMax", "jitterP95", "jitterPHi", "packetsLost", "packetsLostPct",
		"lostBurstMax", "packetsReceived"}
	defaultIngestionProfileFlowmeterMetricNames = []string{
		"throughputAvg", "throughputMax", "throughputMin", "bytesReceived", "packetsReceived"}
	defaultThresholdProfileShell *tenmod.ThresholdProfile
	defaultThresholdsBytes       = []byte(`{
		"thresholds": {
			"vendorMap": {
				"accedian-flowmeter": {
					"metricMap": {
						"throughputAvg": {
							"eventAttrMap": {
								"critical": "25000000",
								"enabled": "true",
								"major": "20000000",
								"minor": "18000000"
							}
						}
					},
					"monitoredObjectTypeMap": {
						"flowmeter": {
							"metricMap": {
								"throughputAvg": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "25000000",
														"lowerStrict": "true",
														"unit": "bps"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "20000000",
														"lowerStrict": "true",
														"unit": "bps",
														"upperLimit": "25000000",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "18000000",
														"lowerStrict": "true",
														"unit": "bps",
														"upperLimit": "20000000"
													}
												}
											}
										}
									}
								}
							}
						}
					}
				},
				"accedian-twamp": {
					"metricMap": {
						"delayP95": {
							"eventAttrMap": {
								"critical": "100000",
								"enabled": "true",
								"major": "95000",
								"minor": "92500"
							}
						},
						"jitterP95": {
							"eventAttrMap": {
								"critical": "30000",
								"enabled": "true",
								"major": "20000",
								"minor": "15000"
							}
						},
						"packetsLostPct": {
							"eventAttrMap": {
								"critical": "0.8",
								"enabled": "true",
								"major": "0.3",
								"minor": "0.1"
							}
						}
					},
					"monitoredObjectTypeMap": {
						"twamp-sf": {
							"metricMap": {
								"delayP95": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "100000",
														"lowerStrict": "true",
														"unit": "ms"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "95000",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "100000",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "92500",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "95000"
													}
												}
											}
										}
									}
								},
								"jitterP95": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "30000",
														"lowerStrict": "true",
														"unit": "ms"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "20000",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "30000",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "15000",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "20000"
													}
												}
											}
										}
									}
								},
								"packetsLostPct": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "0.8",
														"lowerStrict": "true",
														"unit": "pct"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "0.3",
														"lowerStrict": "true",
														"unit": "pct",
														"upperLimit": "0.8",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "0.1",
														"lowerStrict": "true",
														"unit": "pct",
														"upperLimit": "0.3"
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`)

	defaultIngstionProfileShell  *tenmod.IngestionProfile
	defaultIngestionMetricsBytes = []byte(`{"metrics": {
		"vendorMap": {
		  "accedian-flowmeter": {
			"monitoredObjectTypeMap": {
			  "flowmeter": {
				"metricMap": {
				  "bytesReceived": true,
				  "packetsReceived": true,
				  "throughputAvg": true,
				  "throughputMax": true,
				  "throughputMin": true
				}
			  }
			}
		  },
		  "accedian-twamp": {
			"monitoredObjectTypeMap": {
			  "twamp-pe": {
				"metricMap": {
				  "delayMax": true,
				  "delayP95": true,
				  "delayPHi": true,
				  "delayVarP95": true,
				  "delayVarPHi": true,
				  "jitterMax": true,
				  "jitterP95": true,
				  "jitterPHi": true,
				  "lostBurstMax": true,
				  "packetsLost": true,
				  "packetsLostPct": true,
				  "packetsReceived": true
				}
			  },
			  "twamp-sf": {
				"metricMap": {
				  "delayMax": true,
				  "delayP95": true,
				  "delayPHi": true,
				  "delayVarP95": true,
				  "delayVarPHi": true,
				  "jitterMax": true,
				  "jitterP95": true,
				  "jitterPHi": true,
				  "lostBurstMax": true,
				  "packetsLost": true,
				  "packetsLostPct": true,
				  "packetsReceived": true
				}
			  },
			  "twamp-sl": {
				"metricMap": {
				  "delayMax": true,
				  "delayP95": true,
				  "delayPHi": true,
				  "delayVarP95": true,
				  "delayVarPHi": true,
				  "jitterMax": true,
				  "jitterP95": true,
				  "jitterPHi": true,
				  "lostBurstMax": true,
				  "packetsLost": true,
				  "packetsLostPct": true,
				  "packetsReceived": true
				}
			  }
			}
		  }
		}
	  }}`)
)

func checkError(err error, errorType httpErrorString) bool {
	if strings.Contains(err.Error(), string(errorType)) {
		return true
	}

	return false
}

func createDefaultTenantIngPrf(tenantID string) *tenmod.IngestionProfile {
	if defaultIngstionProfileShell == nil {
		defaultIngstionProfileShell = &tenmod.IngestionProfile{}
		if err := json.Unmarshal(defaultIngestionMetricsBytes, &defaultIngstionProfileShell); err != nil {
			logger.Log.Debugf("Unable to construct Default Ingestion Profile from bytes: %s", err.Error())
		}
	}

	logger.Log.Debugf("Will use: %s", models.AsJSONString(defaultIngstionProfileShell))

	ingPrf := tenmod.IngestionProfile{}
	ingPrf.TenantID = tenantID
	ingPrf.Datatype = string(tenmod.TenantIngestionProfileType)
	ingPrf.CreatedTimestamp = db.MakeTimestamp()
	ingPrf.LastModifiedTimestamp = ingPrf.CreatedTimestamp

	ingPrf.Metrics = defaultIngstionProfileShell.Metrics

	logger.Log.Debugf("Generated: %s", models.AsJSONString(ingPrf))

	return &ingPrf
}

func createDefaultTenantThresholdPrf(tenantID string) *tenmod.ThresholdProfile {
	if defaultThresholdProfileShell == nil {
		defaultThresholdProfileShell = &tenmod.ThresholdProfile{}
		if err := json.Unmarshal(defaultThresholdsBytes, &defaultThresholdProfileShell); err != nil {
			logger.Log.Debugf("Unable to construct Default Ingestion Dictionary from file: %s", err.Error())
		}
		logger.Log.Debugf("The defualt thresholds used will be: %s", models.AsJSONString(defaultThresholdProfileShell))
	}
	thrPrf := tenmod.ThresholdProfile{}

	thrPrf.TenantID = tenantID
	thrPrf.Datatype = string(tenmod.TenantThresholdProfileType)
	thrPrf.Name = "Default"
	thrPrf.Thresholds = defaultThresholdProfileShell.Thresholds

	thrPrf.CreatedTimestamp = db.MakeTimestamp()
	thrPrf.LastModifiedTimestamp = thrPrf.CreatedTimestamp

	return &thrPrf
}

func createDefaultTenantMeta(tenantID string, defaultThresholdProfile string, tenantName string) *tenmod.Metadata {
	result := tenmod.Metadata{}

	result.TenantID = tenantID
	result.Datatype = string(tenmod.TenantMetaType)
	result.TenantName = tenantName

	result.CreatedTimestamp = db.MakeTimestamp()
	result.LastModifiedTimestamp = result.CreatedTimestamp

	return &result
}

func getDBFieldFromRequest(r *http.Request, urlPart int32) string {
	urlParts := strings.Split(r.URL.Path, "/")
	return urlParts[urlPart]
}

func reportInternalError(startTime time.Time, code string, objType string, msg string) {
	trackAPIMetrics(startTime, code, objType)
	logger.Log.Error(msg)
}
func reportError(w http.ResponseWriter, startTime time.Time, code string, objType string, msg string, responseCode int) {
	trackAPIMetrics(startTime, code, objType)
	logger.Log.Error(msg)
	http.Error(w, fmt.Sprintf(msg), responseCode)
}

func getRequestBytes(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalRequest(r *http.Request, data interface{}, isUpdate bool) error {
	if err := unmarshalData(r, data); err != nil {
		return err
	}

	// Validate the request
	return validateRESTObject(data, isUpdate)
}

func unmarshalData(r *http.Request, data interface{}) error {
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(requestBytes, &data)
	if err != nil {
		return err
	}

	return nil
}

func sendSuccessResponse(result interface{}, w http.ResponseWriter, startTime time.Time, monLogStr string, objTypeStr string, opTypeString string) {
	// Convert the res to byte[]
	res, err := jsonapi.Marshal(result)
	if err != nil {
		msg := generateErrorMessage(errorMarshal, err.Error())
		reportError(w, startTime, "500", monLogStr, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, jsonAPIContentType)
	logger.Log.Infof("%s %s: %s", opTypeString, objTypeStr, models.AsJSONString(result))
	trackAPIMetrics(startTime, "200", monLogStr)
	fmt.Fprintf(w, string(res))
}

func generateErrorMessage(errCode int, errMsg string) string {
	switch errCode {
	case http.StatusBadRequest:
		return fmt.Sprintf("Unable to read request: %s", errMsg)
	case errorMarshal:
		return fmt.Sprintf("Unable to marshal response: %s", errMsg)
	default:
		return errMsg
	}
}

/*
User roles as defined by Skylight AAA
    SkylightAdmin UserRole = "skylight-admin"
    TenantAdmin   UserRole = "tenant-admin"
    TenantUser    UserRole = "tenant-user"
    UnknownRole   UserRole = "unknown"
*/
const (
	UserRoleSkylight    = "skylight-admin"
	UserRoleTenantAdmin = "tenant-admin"
	UserRoleTenantUser  = "tenant-user"
	UserRoleSystem      = "system"
	UserRoleUnknown     = "unknown"
)

// X-Forward strings that will come from skylight AAA
/*
X-Forwarded-User-Id   (format string)
X-Forwarded-User-Username  (format string)
X-Forwarded-User-Roles   (format string)
X-Forwarded-Tenant-Id   (format string)
*/
const (
	XFwdUserId    = "X-Forwarded-User-Id"
	XFwdUserName  = "X-Forwarded-Username"
	XFwdUserRoles = "X-Forwarded-User-Roles"
	XFwdTenantId  = "X-Forwarded-Tenant-Id"
)

// RequestUserAuth - AAA will forward us information about the requester and this struct will hold the info
type RequestUserAuth struct {
	UserID   string
	UserName string
	// Roles are CSV
	UserRoles []string
	TenantID  string
}

// ExtractHeaderToUserAuthRequest - Converts a header into a requestUserAuth struct
func ExtractHeaderToUserAuthRequest(h http.Header) (*RequestUserAuth, error) {
	logger.Log.Debugf("Received Headers: %s", models.AsJSONString(h))
	roles := h.Get(XFwdUserRoles)
	lRoles := strings.Split(roles, ",")
	req := RequestUserAuth{
		UserID:    h.Get(XFwdUserId),
		UserRoles: lRoles,
		UserName:  h.Get(XFwdUserName),
		TenantID:  h.Get(XFwdTenantId),
	}

	return &req, nil
}

// GetAuthorizationToggle - Check if we need to check the header for authorizations
func GetAuthorizationToggle() bool {
	cfg := gather.GetConfig()
	authAAA := cfg.GetBool(gather.CK_args_authorizationAAA.String())
	logger.Log.Debugf("AAA Auth is enabled? %t", authAAA)

	return authAAA
}

// GetChangeNotificationsToggle - Check if we need to send notifications for certain model changes
func GetChangeNotificationsToggle() bool {
	cfg := gather.GetConfig()
	chgNtf := cfg.GetBool("changeNotifications")
	logger.Log.Debugf("Change Notifications are enabled? %t", chgNtf)

	return chgNtf
}

// RoleAccessControl - Checks if the user-role from AAA is allowed to access this endpoint
func RoleAccessControl(header http.Header, allowedRoles []string) bool {
	// if auth is disabled, let the calls go through
	if GetAuthorizationToggle() == false {
		return true
	}

	user, err := ExtractHeaderToUserAuthRequest(header)
	if err != nil {
		logger.Log.Error("Error parsing header's x-forwards")
		return false
	}

	if allowedRoles == nil {
		logger.Log.Error("Allowed roles is nil, this cannot be.")
		return false
	}

	if len(allowedRoles) == 0 {
		logger.Log.Error("No allowed roles for this endpoint, contact admin for more info")
		return false
	}

	// We currenly only support 1 allowed role, this may change in the future
	allowedRole := allowedRoles[0]
	if allowedRole == UserRoleSystem {
		// Always allow the "system" level auth access to the APIs
		logger.Log.Debugf("Access role %s provided. Access Granted", allowedRole)
		return true
	}

	// Otherwise, handle the roles
	for _, userRole := range user.UserRoles {
		for _, provisionnedRole := range allowedRoles {
			if userRole == provisionnedRole {
				logger.Log.Debugf("Request from %s matches allowed access: %s", userRole, provisionnedRole)
				return true
			}
		}
	}

	logger.Log.Debugf("Request from %s doesn't match any allowed access roles: User Roles{%s}, Provisionned Access Roles {%s}",
		user.UserRoles,
		allowedRoles)
	return false
}

// BuildRouteHandlerWithRAC - To simplify maintainance, this function adds Role Access Control to existing http.serve functions
func BuildRouteHandlerWithRAC(allow []string, fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	functor := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		user, _ := ExtractHeaderToUserAuthRequest(r.Header)
		if RoleAccessControl(r.Header, allow) == false {
			logger.Log.Errorf("User role is not allowed to access endpoint")
			msg := fmt.Sprintf("%s (role:%s) is not allowed to access this endpoint %s.", user.UserName, user.UserRoles, r.URL.Path)
			reportError(w, startTime, "401", "Build Route Handler", msg, http.StatusUnauthorized)

			return
		}
		fn(w, r)

	}

	return functor
}

// BuildRouteHandlerWithRAC - To simplify maintainance, this function adds Role Access Control to existing http.serve functions
func BuildRouteHandlerWithRACSystemCall(allow []string, fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	functor := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		user, _ := ExtractHeaderToUserAuthRequest(r.Header)
		if RoleAccessControl(r.Header, allow) == false {
			logger.Log.Errorf("User role is not allowed to access endpoint")
			msg := fmt.Sprintf("%s (role:%s) is not allowed to access this endpoint %s.", user.UserName, user.UserRoles, r.URL.Path)
			reportError(w, startTime, "401", "Build Route Handler", msg, http.StatusUnauthorized)

			return
		}
		fn(w, r)

	}

	return functor
}

// reportAPIError - Used to document API errors both in logging and in the Metrics reporting tool.
func reportAPIError(msg string, startTime time.Time, code int, objType string, counterMetrics ...mon.MetricCounterType) string {
	logger.Log.Errorf(msg)
	reportAPICompletionState(startTime, code, objType, counterMetrics...)
	return msg
}

// reportAPICompletionState - Used to document API completion state both in logging and in the Metrics reporting tool.
func reportAPICompletionState(startTime time.Time, code int, objType string, counterMetrics ...mon.MetricCounterType) {
	incrementAPICounters(counterMetrics...)
	trackAPIMetricsByHttpCode(startTime, code, objType)
}

// incrementAPICounters - updates API call counters in the metric service
func incrementAPICounters(counterMetrics ...mon.MetricCounterType) {
	for _, counter := range counterMetrics {
		mon.IncrementCounter(counter)
	}
}

// trackAPIMetrics - updates API call durations in the metric service
func trackAPIMetricsByHttpCode(startTime time.Time, code int, objType string) {
	codeStr := strconv.Itoa(code)
	mon.TrackAPITimeMetricInSeconds(startTime, codeStr, objType)
}

func convertToJsonapiObject(obj interface{}, dataContainer interface{}) error {

	// Marshal this object into the appropriate format
	jsonapiBytes, err := jsonapi.Marshal(obj)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonapiBytes, dataContainer)
	if err != nil {
		return err
	}

	return nil
}

func wrapJsonAPIObject(obj interface{}, id string, otype string) map[string]interface{} {

	return map[string]interface{}{
		"data": map[string]interface{}{
			"id":         id,
			"type":       otype,
			"attributes": obj}}
}

func wrapJsonAPIObjectAsArray(obj interface{}, id string, otype string) map[string]interface{} {

	return map[string]interface{}{
		"data": []map[string]interface{}{
			map[string]interface{}{
				"id":         id,
				"type":       otype,
				"attributes": obj}}}
}

// authorizeRequest - Does the initial setup of a REST handler function, including logging, incrementing API counters for monitoring and tracking the
// initialization time of the call.
// Returns:
//  - isAuthorized (bool) -> indicates if the request was authorized based on the logged in userr's role.
//  - startTime -> the time at which this API call was initiated
func authorizeRequest(initMsg string, req *http.Request, allowedRoles []string, countersToIncrement ...mon.MetricCounterType) (bool, time.Time) {
	startTime := time.Now()
	incrementAPICounters(countersToIncrement...)

	logger.Log.Info(initMsg)

	return isRequestAuthorized(req, allowedRoles), startTime
}

// convertRequestBodyToDBModel - converts a generic object into a know type. Useful for converting a REST request body into a DB model.
func convertRequestBodyToDBModel(requestBody interface{}, dataContainer interface{}) error {
	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	err = jsonapi.Unmarshal(requestBytes, dataContainer)
	if err != nil {
		return err
	}

	return nil
}

func checkForNotFound(s string) bool {
	return strings.Contains(s, string(notFound))
}

// generateLinks - creates the "links" section to be used in a jsonapi response object
func generateLinks(urlBase string, paginationOffsets *common.PaginationOffsets, limit int64) map[string]string {
	links := map[string]string{}

	links[linkFirstStr] = fmt.Sprintf("%s?%s=%d", urlBase, limitQueryParamStr, limit)
	links[linkSelfStr] = fmt.Sprintf("%s?%s=%d", urlBase, limitQueryParamStr, limit)

	if len(paginationOffsets.Self) != 0 {
		links[linkSelfStr] = fmt.Sprintf("%s?%s=%s&%s=%d", urlBase, startKeyQueryParamStr, paginationOffsets.Self, limitQueryParamStr, limit)
	}

	if len(paginationOffsets.Next) != 0 {
		links[linkNextStr] = fmt.Sprintf("%s?%s=%s&%s=%d", urlBase, startKeyQueryParamStr, paginationOffsets.Next, limitQueryParamStr, limit)
	}

	if len(paginationOffsets.Prev) != 0 {
		links[linkPrevStr] = fmt.Sprintf("%s?%s=%s&%s=%d", urlBase, startKeyQueryParamStr, paginationOffsets.Prev, limitQueryParamStr, limit)
	}

	return links
}

func getIdsFromRelationshipData(relationships *swagmodels.JSONAPIRelationship) []string {
	result := []string{}

	for _, val := range relationships.Data {
		result = append(result, val.ID)
	}

	return result
}

func addValToList(list []*metmod.MetricViolationSummaryType, entry map[string]interface{}, ts string) []*metmod.MetricViolationSummaryType {

	if entry == nil {
		entry = make(map[string]interface{})
	}

	if entry["violationDuration"] != nil || entry["violationCount"] != nil {

		// var check float64
		_, ok := entry["violationDuration"].(float64)
		if !ok {
			return list
		}

		_, ok = entry["violationDuration"].(float64)
		if !ok {
			return list
		}

		tmp := metmod.MetricViolationSummaryType{
			"violationCount":    entry["violationCount"],
			"violationDuration": entry["violationDuration"],
			"timestamp":         ts,
		}

		list = append(list, &tmp)
	}

	return list
}

func buildHash(args ...string) string {
	str := ""
	for _, w := range args {
		str = str + "." + w
	}
	return str
}

func utilMetricViolationSummaryTypeMap2Array(src map[string]*metmod.MetricViolationSummaryType) []*metmod.MetricViolationSummaryType {
	var res []*metmod.MetricViolationSummaryType
	for _, v := range src {
		res = append(res, v)
	}
	return res
}

func utilMetricViolationsSummaryAsTimeSeriesEntryMap2Array(src map[string]*metmod.MetricViolationsSummaryAsTimeSeriesEntry) []*metmod.MetricViolationsSummaryAsTimeSeriesEntry {
	var res []*metmod.MetricViolationsSummaryAsTimeSeriesEntry
	for _, v := range src {
		res = append(res, v)
	}
	return res
}
