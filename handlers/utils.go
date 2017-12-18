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
	moMap.MonitoredObjectMap = make(map[string]*pb.TenantIngestionProfile_MetricMap)
	moMap.MonitoredObjectMap["pe"] = &metricMap
	moMap.MonitoredObjectMap["sl"] = &metricMap
	moMap.MonitoredObjectMap["sf"] = &metricMap
	vendorMap := make(map[string]*pb.TenantIngestionProfile_MonitoredObjectMap)
	vendorMap["accedian"] = &moMap
	ingPrf.VendorMap = vendorMap

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
	thrPrf.Thresholds = &pb.TenantThresholdProfile_VendorMap{}
	// TODO: Add in the hardcoded defaults here.
	thrPrf.CreatedTimestamp = time.Now().Unix()
	thrPrf.LastModifiedTimestamp = thrPrf.GetCreatedTimestamp()

	return &thrPrf
}

func getDBFieldFromRequest(r *http.Request, urlPart int32) string {
	urlParts := strings.Split(r.URL.Path, "/")
	return urlParts[urlPart]
}
