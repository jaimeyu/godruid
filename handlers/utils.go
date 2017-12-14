package handlers

import (
	"net/http"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

func createDefaultTenantIngPrf(tenantID string) *pb.TenantIngestionProfile {
	ingPrf := pb.TenantIngestionProfile{}
	ingPrf.ScpUsername = "default"
	ingPrf.ScpPassword = "password"
	ingPrf.TenantId = tenantID
	ingPrf.Datatype = string(db.TenantIngestionProfileType)
	ingPrf.CreatedTimestamp = time.Now().Unix()
	ingPrf.LastModifiedTimestamp = ingPrf.GetCreatedTimestamp()

	return &ingPrf
}

func createDefaultTenantThresholdPrf(tenantID string) *pb.TenantThresholdProfile {
	thrPrf := pb.TenantThresholdProfile{}

	thrPrf.TenantId = tenantID
	thrPrf.Datatype = string(db.TenantThresholdProfileType)
	thrPrf.Thresholds = []*pb.TenantThreshold{}
	// TODO: Add in the hardcoded defaults here.
	thrPrf.CreatedTimestamp = time.Now().Unix()
	thrPrf.LastModifiedTimestamp = thrPrf.GetCreatedTimestamp()

	return &thrPrf
}

func getDBFieldFromRequest(r *http.Request, urlPart int32) string {
	urlParts := strings.Split(r.URL.Path, "/")
	return urlParts[urlPart]
}
