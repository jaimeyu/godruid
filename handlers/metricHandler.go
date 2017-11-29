package handlers

import (
	"context"
	"fmt"

	"github.com/accedian/adh-gather/datastore"
	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

type MetricServiceHandler struct {
	druidDB db.DruidDatastore
}

func CreateMetricServiceHandler() *MetricServiceHandler {
	result := new(MetricServiceHandler)

	db := druid.NewDruidDatasctoreClient()

	result.druidDB = db

	return result
}

// GetThresholdCrossing
func (msh *MetricServiceHandler) GetThresholdCrossing(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest) (*pb.JSONAPIObject, error) {

	// Issue request to DAO Layer to Get the requested Admin User

	fmt.Println(ctx)

	result, err := msh.druidDB.GetThresholdCrossing("", "")

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.AdminUserStr, err.Error())
	}

	return result, nil
}
