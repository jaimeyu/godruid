package handlers

import (
	"encoding/json"
	"time"

	"github.com/accedian/adh-gather/gather"

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
	cfg := gather.GetConfig()

	result := MetricBaselineProvisioner{}
	result.tenantDB = db
	result.requestReader = messaging.CreateKafkaReaderWithSyncTime(metricBaselineRequestTopic, "0", 5*time.Minute)

	logger.Log.Infof("Starting Metric Baseline Provisioner for topic: %s", metricBaselineRequestTopic)

	// Start the message readers
	go func() {
		numJobs := cfg.GetInt(gather.CK_args_metricbaselines_maxnumjobs.String())
		jobs := make(chan []byte, numJobs)

		numWorkers := cfg.GetInt(gather.CK_args_metricbaselines_numworkers.String())
		for w := 1; w <= numWorkers; w++ {
			go result.metricBaselineProvsioningWorker(w, jobs)
		}

		for {
			msgBytes, err := result.requestReader.ReadMessageWithoutExplicitOffsetManagement()
			if err == nil {
				jobs <- msgBytes
			}
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

	logger.Log.Infof("Received %s for Monitored Object %s with %d metric baseline values", metricBaselineRequestLogStr, requestObj.MonitoredObjectID, len(requestObj.Baselines))
	if logger.IsDebugEnabled() {
		logger.Log.Infof("Received %s: %s", metricBaselineRequestLogStr, models.AsJSONString(requestObj))
	}

	_, err = mbp.tenantDB.UpdateMetricBaselineForHourOfWeekWithCollection(requestObj.TenantID, requestObj.MonitoredObjectID, requestObj.Baselines)
	if err != nil {
		logger.Log.Errorf("Error updating %s for %s %s for %s %s for baseline data %s: %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, requestObj.TenantID, tenmod.TenantMonitoredObjectStr, requestObj.MonitoredObjectID, models.AsJSONString(requestObj.Baselines), err.Error())
	}

	logger.Log.Infof("Completed %s for Monitored Object %s with %d metric baseline values", metricBaselineRequestLogStr, requestObj.MonitoredObjectID, len(requestObj.Baselines))
	if logger.IsDebugEnabled() {
		logger.Log.Infof("Completed %s: %s", metricBaselineRequestLogStr, models.AsJSONString(requestObj))
	}
	return true
}

func (mbp *MetricBaselineProvisioner) metricBaselineProvsioningWorker(id int, jobs <-chan []byte) {
	for j := range jobs {
		logger.Log.Debugf("Metric Baseline provisioning worker %d started  job: %s", id, string(j))
		result := mbp.handleMetricBaselineProvisioningRequest(j)
		logger.Log.Debugf("Metric Baseline provisioning worker %d completed job: %s with result %t", id, string(j), result)
	}
}
