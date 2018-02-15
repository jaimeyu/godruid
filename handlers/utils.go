package handlers

import (
	"net/http"
	"strings"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

var (
	defaultIngestionProfileMetricNames = []string{
		"delayMax", "delayP95", "delayPHi", "delayVarP95", "delayVarPHi",
		"jitterMax", "jitterP95", "jitterPHi", "packetsLost", "packetsLostPct",
		"lostBurstMax", "packetsReceived"}
	defaultIngestionProfileFlowmeterMetricNames = []string{
		"throughputAvg", "throughputMax", "throughputMin", "bytesReceived", "packetsReceived"}
)

func createDefaultTenantIngPrf(tenantID string) *pb.TenantIngestionProfileData {
	ingPrf := pb.TenantIngestionProfileData{}
	ingPrf.TenantId = tenantID
	ingPrf.Datatype = string(db.TenantIngestionProfileType)
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
	thrPrf.Datatype = string(db.TenantThresholdProfileType)
	thrPrf.Name = "Default"

	thrPrf.Thresholds = createDefaultThreshold()

	thrPrf.CreatedTimestamp = db.MakeTimestamp()
	thrPrf.LastModifiedTimestamp = thrPrf.GetCreatedTimestamp()

	return &thrPrf
}

func createDefaultTenantMeta(tenantID string, defaultThresholdProfile string, tenantName string) *pb.TenantMeta {
	result := pb.TenantMeta{}

	result.TenantId = tenantID
	result.Datatype = string(db.TenantMetaType)
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
				MonitoredObjectTypeMap: map[string]*pb.TenantThresholdProfileData_MetricMap{
					string(TwampPE): &pb.TenantThresholdProfileData_MetricMap{
						MetricMap: map[string]*pb.TenantThresholdProfileData_DirectionMap{
							"delayP95": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "135000",
													"lowerStrict": "true",
													"upperLimit":  "140000",
													"unit":        "ms",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "140000",
													"lowerStrict": "true",
													"upperLimit":  "170000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "170000",
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
													"lowerLimit":  "35000",
													"lowerStrict": "true",
													"upperLimit":  "40000",
													"unit":        "ms",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "40000",
													"lowerStrict": "true",
													"upperLimit":  "45000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "45000",
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
													"lowerLimit":  "0.20",
													"lowerStrict": "true",
													"upperLimit":  "0.30",
													"unit":        "pct",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "0.30",
													"lowerStrict": "true",
													"upperLimit":  "0.35",
													"upperStrict": "false",
													"unit":        "pct",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "0.35",
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
				MonitoredObjectTypeMap: map[string]*pb.TenantThresholdProfileData_MetricMap{
					string(Flowmeter): &pb.TenantThresholdProfileData_MetricMap{
						MetricMap: map[string]*pb.TenantThresholdProfileData_DirectionMap{
							"throughputAvg": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "35000000",
													"lowerStrict": "true",
													"upperLimit":  "50000000",
													"unit":        "bps",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "50000000",
													"lowerStrict": "true",
													"upperLimit":  "70000000",
													"upperStrict": "false",
													"unit":        "bps",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "70000000",
													"lowerStrict": "true",
													"unit":        "bps",
												},
											},
										},
									},
								},
							},
							"throughputMax": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "35000000",
													"lowerStrict": "true",
													"upperLimit":  "50000000",
													"unit":        "bps",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "50000000",
													"lowerStrict": "true",
													"upperLimit":  "80000000",
													"upperStrict": "false",
													"unit":        "bps",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "80000000",
													"lowerStrict": "true",
													"unit":        "bps",
												},
											},
										},
									},
								},
							},
							"throughputMin": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "35000000",
													"lowerStrict": "true",
													"upperLimit":  "50000000",
													"unit":        "bps",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "50000000",
													"lowerStrict": "true",
													"upperLimit":  "80000000",
													"upperStrict": "false",
													"unit":        "bps",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "80000000",
													"lowerStrict": "true",
													"unit":        "bps",
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
		},
	}
}

func getDBFieldFromRequest(r *http.Request, urlPart int32) string {
	urlParts := strings.Split(r.URL.Path, "/")
	return urlParts[urlPart]
}
