package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/manyminds/api2go/jsonapi"
)

type httpErrorString string

const (
	// Error search strings
	notFound httpErrorString = "status 404 - not found"

	// Custom Error types
	errorMarshal = -100

	// API Prefix values
	apiV1Prefix      = "/api/v1/"
	tenantsAPIPrefix = "tenants/{tenantID}/"
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
						"twamp-pe": {
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
	result.DefaultThresholdProfile = defaultThresholdProfile
	result.TenantName = tenantName

	result.CreatedTimestamp = db.MakeTimestamp()
	result.LastModifiedTimestamp = result.CreatedTimestamp

	return &result
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
