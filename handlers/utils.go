package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/manyminds/api2go/jsonapi"
)

type httpErrorString string

const (
	notFound httpErrorString = "status 404 - not found"
)

const (
	errorMarshal = -100
)

var (
	defaultIngestionProfileMetricNames = []string{
		"delayMax", "delayP95", "delayPHi", "delayVarP95", "delayVarPHi",
		"jitterMax", "jitterP95", "jitterPHi", "packetsLost", "packetsLostPct",
		"lostBurstMax", "packetsReceived"}
	defaultIngestionProfileFlowmeterMetricNames = []string{
		"throughputAvg", "throughputMax", "throughputMin", "bytesReceived", "packetsReceived"}
)

func checkError(err error, errorType httpErrorString) bool {
	if strings.Contains(err.Error(), string(errorType)) {
		return true
	}

	return false
}

func createDefaultTenantIngPrf(tenantID string) *pb.TenantIngestionProfileData {
	ingPrf := pb.TenantIngestionProfileData{}
	ingPrf.TenantId = tenantID
	ingPrf.Datatype = string(tenmod.TenantIngestionProfileType)
	ingPrf.CreatedTimestamp = db.MakeTimestamp()
	ingPrf.LastModifiedTimestamp = ingPrf.GetCreatedTimestamp()

	// Default Values for the metrics:
	moMap := pb.TenantIngestionProfileData_MonitoredObjectMap{}
	metricMap := pb.TenantIngestionProfileData_MetricMap{}
	metricMap.MetricMap = createMetricMap(defaultIngestionProfileMetricNames...)
	moMap.MonitoredObjectTypeMap = make(map[string]*pb.TenantIngestionProfileData_MetricMap)
	moMap.MonitoredObjectTypeMap[string(TwampPE)] = &metricMap
	moMap.MonitoredObjectTypeMap[string(TwampSL)] = &metricMap
	moMap.MonitoredObjectTypeMap[string(TwampSF)] = &metricMap
	metrics := make(map[string]*pb.TenantIngestionProfileData_MonitoredObjectMap)
	metrics[string(AccedianTwamp)] = &moMap

	// Add flowmeter metrics:
	flowMOMap := pb.TenantIngestionProfileData_MonitoredObjectMap{}
	flowMetricMap := pb.TenantIngestionProfileData_MetricMap{}
	flowMetricMap.MetricMap = createMetricMap(defaultIngestionProfileFlowmeterMetricNames...)
	flowMOMap.MonitoredObjectTypeMap = make(map[string]*pb.TenantIngestionProfileData_MetricMap)
	flowMOMap.MonitoredObjectTypeMap[string(Flowmeter)] = &flowMetricMap
	metrics[string(AccedianFlowmeter)] = &flowMOMap

	vendorMap := &pb.TenantIngestionProfileData_VendorMap{}
	vendorMap.VendorMap = metrics
	ingPrf.Metrics = vendorMap

	return &ingPrf
}

func createMetricMap(metricNames ...string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range metricNames {
		result[s] = true
	}

	return result
}

func createDefaultTenantThresholdPrf(tenantID string) *pb.TenantThresholdProfileData {
	thrPrf := pb.TenantThresholdProfileData{}

	thrPrf.TenantId = tenantID
	thrPrf.Datatype = string(tenmod.TenantThresholdProfileType)
	thrPrf.Name = "Default"

	thrPrf.Thresholds = createDefaultThreshold()

	thrPrf.CreatedTimestamp = db.MakeTimestamp()
	thrPrf.LastModifiedTimestamp = thrPrf.GetCreatedTimestamp()

	return &thrPrf
}

func createDefaultTenantMeta(tenantID string, defaultThresholdProfile string, tenantName string) *pb.TenantMeta {
	result := pb.TenantMeta{}

	result.TenantId = tenantID
	result.Datatype = string(tenmod.TenantMetaType)
	result.DefaultThresholdProfile = defaultThresholdProfile
	result.TenantName = tenantName

	result.CreatedTimestamp = db.MakeTimestamp()
	result.LastModifiedTimestamp = result.GetCreatedTimestamp()

	return &result
}

func createDefaultThreshold() *pb.TenantThresholdProfileData_VendorMap {
	return &pb.TenantThresholdProfileData_VendorMap{
		VendorMap: map[string]*pb.TenantThresholdProfileData_MonitoredObjectTypeMap{
			string(AccedianTwamp): &pb.TenantThresholdProfileData_MonitoredObjectTypeMap{
				MetricMap: map[string]*pb.TenantThresholdProfileData_UIEventAttrMap{
					"delayP95": &pb.TenantThresholdProfileData_UIEventAttrMap{
						EventAttrMap: map[string]string{
							"enabled":  "true",
							"minor":    "92500",
							"major":    "95000",
							"critical": "100000",
						},
					},
					"jitterP95": &pb.TenantThresholdProfileData_UIEventAttrMap{
						EventAttrMap: map[string]string{
							"enabled":  "true",
							"minor":    "15000",
							"major":    "20000",
							"critical": "30000",
						},
					},
					"packetsLostPct": &pb.TenantThresholdProfileData_UIEventAttrMap{
						EventAttrMap: map[string]string{
							"enabled":  "true",
							"minor":    "0.1",
							"major":    "0.3",
							"critical": "0.8",
						},
					},
				},
				MonitoredObjectTypeMap: map[string]*pb.TenantThresholdProfileData_MetricMap{
					string(TwampPE): &pb.TenantThresholdProfileData_MetricMap{
						MetricMap: map[string]*pb.TenantThresholdProfileData_DirectionMap{
							"delayP95": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "92500",
													"lowerStrict": "true",
													"upperLimit":  "95000",
													"unit":        "ms",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "95000",
													"lowerStrict": "true",
													"upperLimit":  "100000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "100000",
													"lowerStrict": "true",
													"unit":        "ms",
												},
											},
										},
									},
								},
							},
							"jitterP95": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "15000",
													"lowerStrict": "true",
													"upperLimit":  "20000",
													"unit":        "ms",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "20000",
													"lowerStrict": "true",
													"upperLimit":  "30000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "30000",
													"lowerStrict": "true",
													"unit":        "ms",
												},
											},
										},
									},
								},
							},
							"packetsLostPct": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "0.1",
													"lowerStrict": "true",
													"upperLimit":  "0.3",
													"unit":        "pct",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "0.3",
													"lowerStrict": "true",
													"upperLimit":  "0.8",
													"upperStrict": "false",
													"unit":        "pct",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "0.8",
													"lowerStrict": "true",
													"unit":        "pct",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			string(AccedianFlowmeter): &pb.TenantThresholdProfileData_MonitoredObjectTypeMap{
				MetricMap: map[string]*pb.TenantThresholdProfileData_UIEventAttrMap{
					"throughputAvg": &pb.TenantThresholdProfileData_UIEventAttrMap{
						EventAttrMap: map[string]string{
							"enabled":  "true",
							"minor":    "18000000",
							"major":    "20000000",
							"critical": "25000000",
						},
					},
					// Removing these items for MWC, leaving them commented out in case
					// there is a desire for them later.
					// "throughputMax": &pb.TenantThresholdProfileData_UIEventAttrMap{
					// 	EventAttrMap: map[string]string{
					// 		"enabled":  "true",
					// 		"minor":    "16500000",
					// 		"major":    "17500000",
					// 		"critical": "20000000",
					// 	},
					// },
					// "throughputMin": &pb.TenantThresholdProfileData_UIEventAttrMap{
					// 	EventAttrMap: map[string]string{
					// 		"enabled":  "true",
					// 		"minor":    "16500000",
					// 		"major":    "17500000",
					// 		"critical": "20000000",
					// 	},
					// },
				},
				MonitoredObjectTypeMap: map[string]*pb.TenantThresholdProfileData_MetricMap{
					string(Flowmeter): &pb.TenantThresholdProfileData_MetricMap{
						MetricMap: map[string]*pb.TenantThresholdProfileData_DirectionMap{
							"throughputAvg": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "18000000",
													"lowerStrict": "true",
													"upperLimit":  "20000000",
													"unit":        "bps",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "20000000",
													"lowerStrict": "true",
													"upperLimit":  "25000000",
													"upperStrict": "false",
													"unit":        "bps",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "25000000",
													"lowerStrict": "true",
													"unit":        "bps",
												},
											},
										},
									},
								},
							},
							// Removing these items for MWC, leaving them commented out in case
							// there is a desire for them later.
							// "throughputMax": &pb.TenantThresholdProfileData_DirectionMap{
							// 	DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
							// 		"0": &pb.TenantThresholdProfileData_EventMap{
							// 			EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
							// 				"minor": &pb.TenantThresholdProfileData_EventAttrMap{
							// 					map[string]string{
							// 						"lowerLimit":  "16500000",
							// 						"lowerStrict": "true",
							// 						"upperLimit":  "17500000",
							// 						"unit":        "bps",
							// 					},
							// 				},
							// 				"major": &pb.TenantThresholdProfileData_EventAttrMap{
							// 					map[string]string{
							// 						"lowerLimit":  "17500000",
							// 						"lowerStrict": "true",
							// 						"upperLimit":  "20000000",
							// 						"upperStrict": "false",
							// 						"unit":        "bps",
							// 					},
							// 				},
							// 				"critical": &pb.TenantThresholdProfileData_EventAttrMap{
							// 					map[string]string{
							// 						"lowerLimit":  "20000000",
							// 						"lowerStrict": "true",
							// 						"unit":        "bps",
							// 					},
							// 				},
							// 			},
							// 		},
							// 	},
							// },
							// "throughputMin": &pb.TenantThresholdProfileData_DirectionMap{
							// 	DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
							// 		"0": &pb.TenantThresholdProfileData_EventMap{
							// 			EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
							// 				"minor": &pb.TenantThresholdProfileData_EventAttrMap{
							// 					map[string]string{
							// 						"lowerLimit":  "16500000",
							// 						"lowerStrict": "true",
							// 						"upperLimit":  "17500000",
							// 						"unit":        "bps",
							// 					},
							// 				},
							// 				"major": &pb.TenantThresholdProfileData_EventAttrMap{
							// 					map[string]string{
							// 						"lowerLimit":  "17500000",
							// 						"lowerStrict": "true",
							// 						"upperLimit":  "20000000",
							// 						"upperStrict": "false",
							// 						"unit":        "bps",
							// 					},
							// 				},
							// 				"critical": &pb.TenantThresholdProfileData_EventAttrMap{
							// 					map[string]string{
							// 						"lowerLimit":  "20000000",
							// 						"lowerStrict": "true",
							// 						"unit":        "bps",
							// 					},
							// 				},
							// 			},
							// 		},
							// 	},
							// },
						},
					},
				},
			},
		},
	}
}

func getDBFieldFromRequest(r *http.Request, urlPart int32) string {
	urlParts := strings.Split(r.URL.Path, "/")
	return urlParts[urlPart]
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
