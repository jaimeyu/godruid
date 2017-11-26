package handlers

import (
	"context"
	"fmt"

	"github.com/accedian/adh-gather/datastore"
	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

type MetricServiceHandler struct {
	druidDB db.DruidDatastore
}

// GetThresholdCrossing
func (msh *MetricServiceHandler) GetThresholdCrossing(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest) (*pb.ThresholdCrossingResponse, error) {

	// Issue request to DAO Layer to Get the requested Admin User

	res, err := msh.druidDB.GetNumberOfThesholdViolations("", "")

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.AdminUserStr, err.Error())
	}

	log.Prinln("SUCCESS")

	return result, nil
}
