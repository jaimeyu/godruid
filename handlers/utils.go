package handlers

import (
	"net/http"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

var (
	defaultIngestionProfileMetricNames = []string{
		"delayMax", "delayP95", "delayPHi", "delayVarP95", "delayVarPHi",
		"jitterMax", "jitterP95", "jitterPHi", "packetsLost", "packetsLostPct",
		"lostBurstMax", "packetsReceived"}
)

func createDefaultTenantIngPrf(tenantID string) *pb.TenantIngestionProfile {
	ingPrf := pb.TenantIngestionProfile{}
	ingPrf.TenantId = tenantID
	ingPrf.Datatype = string(db.TenantIngestionProfileType)
	ingPrf.CreatedTimestamp = time.Now().Unix()
	ingPrf.LastModifiedTimestamp = ingPrf.GetCreatedTimestamp()

	// Default Values for the metrics:
	moMap := pb.TenantIngestionProfile_MonitoredObjectMap{}
	metricMap := pb.TenantIngestionProfile_MetricMap{}
	metricMap.MetricMap = createMetricMap(defaultIngestionProfileMetricNames...)
	moMap.MonitoredObjectTypeMap = make(map[string]*pb.TenantIngestionProfile_MetricMap)
	moMap.MonitoredObjectTypeMap["pe"] = &metricMap
	moMap.MonitoredObjectTypeMap["sl"] = &metricMap
	moMap.MonitoredObjectTypeMap["sf"] = &metricMap
	metrics := make(map[string]*pb.TenantIngestionProfile_MonitoredObjectMap)
	metrics["accedian"] = &moMap
	vendorMap := &pb.TenantIngestionProfile_VendorMap{}
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

func createDefaultTenantThresholdPrf(tenantID string) *pb.TenantThresholdProfile {
	thrPrf := pb.TenantThresholdProfile{}

	thrPrf.TenantId = tenantID
	thrPrf.Datatype = string(db.TenantThresholdProfileType)
	thrPrf.Name = "Default"

	thrPrf.Thresholds = createDefaultThreshold()

	thrPrf.CreatedTimestamp = time.Now().Unix()
	thrPrf.LastModifiedTimestamp = thrPrf.GetCreatedTimestamp()

	return &thrPrf
}

func createDefaultTenantMeta(tenantID string, defaultThresholdProfile string, tenantName string) *pb.TenantMetaData {
	result := pb.TenantMetaData{}

	result.TenantId = tenantID
	result.Datatype = string(db.TenantMetaType)
	result.DefaultThresholdProfile = defaultThresholdProfile
	result.TenantName = tenantName

	result.CreatedTimestamp = time.Now().Unix()
	result.LastModifiedTimestamp = result.GetCreatedTimestamp()

	return &result
}

func createDefaultThreshold() *pb.TenantThresholdProfile_VendorMap {
	return &pb.TenantThresholdProfile_VendorMap{
		VendorMap: map[string]*pb.TenantThresholdProfile_MonitoredObjectTypeMap{
			"accedian": &pb.TenantThresholdProfile_MonitoredObjectTypeMap{
				MonitoredObjectTypeMap: map[string]*pb.TenantThresholdProfile_MetricMap{
					"pe": &pb.TenantThresholdProfile_MetricMap{
						MetricMap: map[string]*pb.TenantThresholdProfile_DirectionMap{
							"delayP95": &pb.TenantThresholdProfile_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfile_EventMap{
									"0": &pb.TenantThresholdProfile_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfile_EventAttrMap{
											"minor": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"upperLimit": "50000",
													"unit":       "ms",
												},
											},
											"major": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"lowerLimit":  "200000",
													"lowerStrict": "true",
													"upperLimit":  "300000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"lowerLimit":  "300000",
													"lowerStrict": "true",
													"unit":        "ms",
												},
											},
										},
									},
								},
							},
							"jitterP95": &pb.TenantThresholdProfile_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfile_EventMap{
									"0": &pb.TenantThresholdProfile_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfile_EventAttrMap{
											"minor": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"upperLimit": "1000",
													"unit":       "ms",
												},
											},
											"major": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"lowerLimit":  "3500",
													"lowerStrict": "true",
													"upperLimit":  "5000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"lowerLimit":  "5000",
													"lowerStrict": "true",
													"unit":        "ms",
												},
											},
										},
									},
								},
							},
							"packetsLostPct": &pb.TenantThresholdProfile_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfile_EventMap{
									"0": &pb.TenantThresholdProfile_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfile_EventAttrMap{
											"minor": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"upperLimit": "0.5",
													"unit":       "%",
												},
											},
											"major": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"lowerLimit":  "0.5",
													"lowerStrict": "true",
													"upperLimit":  "1.0",
													"upperStrict": "false",
													"unit":        "%",
												},
											},
											"critical": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"lowerLimit":  "1.0",
													"lowerStrict": "true",
													"unit":        "%",
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
