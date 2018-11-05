package handlers

import (
	"encoding/json"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/messaging"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

const (
	metricBaselineRequestLogStr = "metric baseline request"
	metricBaselineRequestTopic  = "baseline-report"
)

type MetricBaselineProvisioner struct {
	requestReader *messaging.KafkaConsumer
	tenantDB      datastore.TenantServiceDatastore
}

func CreateMetricBaselineProvisioner(db datastore.TenantServiceDatastore) *MetricBaselineProvisioner {
	result := MetricBaselineProvisioner{}
	result.tenantDB = db
	result.requestReader = messaging.CreateKafkaReader(metricBaselineRequestTopic, "0")

	logger.Log.Infof("Starting Metric Baseline Provisioner for topic: %s", metricBaselineRequestTopic)

	// Start the message readers
	go func() {
		for {
			result.requestReader.ReadMessage(result.handleMetricBaselineProvisioningRequest)
		}
	}()

	return &result
}

func (mbp *MetricBaselineProvisioner) handleMetricBaselineProvisioningRequest(requestBytes []byte) bool {
	requestObj := &tenmod.MetricBaseline{}
	err := json.Unmarshal(requestBytes, requestObj)
	if err != nil {
		logger.Log.Errorf("Unable to read %s data: %s", metricBaselineRequestLogStr, err.Error())
		return true
	}

	requestObjStr := models.AsJSONString(requestObj)
	logger.Log.Infof("Received %s: %s", metricBaselineRequestLogStr, requestObjStr)

	_, err = mbp.tenantDB.UpdateMetricBaselineForHourOfWeekWithCollection(requestObj.TenantID, requestObj.MonitoredObjectID, requestObj.Baselines)
	if err != nil {
		logger.Log.Errorf("Error updating %s for %s %s for %s %s for baseline data %s: %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, requestObj.TenantID, tenmod.TenantMonitoredObjectStr, requestObj.MonitoredObjectID, models.AsJSONString(requestObj.Baselines), err.Error())
	}

	logger.Log.Infof("Completed %s: %s", metricBaselineRequestLogStr, requestObjStr)
	return true
}
