package couchDB

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	metmod "github.com/accedian/adh-gather/models/metrics"
)

// TODO we should create a generic couch DB object at this point since a bunch of our libs are using it
type SchedulerServiceDatastoreCouchDB struct {
	server string
	cfg    config.Provider
}

func CreateSchedulerServiceDAO() (*SchedulerServiceDatastoreCouchDB, error) {
	result := new(SchedulerServiceDatastoreCouchDB)
	result.cfg = gather.GetConfig()

	provDBURL := fmt.Sprintf("%s:%d",
		result.cfg.GetString(gather.CK_server_datastore_ip.String()),
		result.cfg.GetInt(gather.CK_server_datastore_port.String()))
	logger.Log.Debugf("Scheduler Service CouchDB URL is: %s", provDBURL)
	result.server = provDBURL

	return result, nil
}

func (ssd *SchedulerServiceDatastoreCouchDB) CreateScheduleConfig(slaConfig *metmod.SLAScheduleConfig) (*metmod.SLAScheduleConfig, error) {
	logger.Log.Debugf("Creating %s: %v\n", metmod.SLAScheduleConfigStr, models.AsJSONString(slaConfig))
	slaConfig.ID = ds.GenerateID(slaConfig, "slaConfigType") //TODO Create an actual type
	tenantID := ds.PrependToDataID(slaConfig.TenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(ssd.server, tenantID)

	dataContainer := &metmod.SLAScheduleConfig{}
	if err := createDataInCouch(tenantDBName, slaConfig, dataContainer, "slaConfigType", metmod.SLAScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", metmod.SLAScheduleConfigStr, models.AsJSONString(slaConfig))
	return dataContainer, nil

}
func (ssd *SchedulerServiceDatastoreCouchDB) UpdateScheduleConfig(slaConfig *metmod.SLAScheduleConfig) (*metmod.SLAScheduleConfig, error) {
	logger.Log.Debugf("Updating %s: %v\n", metmod.SLAScheduleConfigStr, models.AsJSONString(slaConfig))
	slaConfig.ID = ds.PrependToDataID(slaConfig.ID, "slaConfigType") //TODO Create an actual type
	tenantID := ds.PrependToDataID(slaConfig.TenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(ssd.server, tenantID)

	dataContainer := &metmod.SLAScheduleConfig{}
	if err := updateDataInCouch(tenantDBName, slaConfig, dataContainer, "slaConfigType", metmod.SLAScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", metmod.SLAScheduleConfigStr, models.AsJSONString(slaConfig))
	return dataContainer, nil
}
func (ssd *SchedulerServiceDatastoreCouchDB) DeleteScheduleConfig(tenantID string, configID string) (*metmod.SLAScheduleConfig, error) {
	logger.Log.Debugf("Deleting %s %s\n", metmod.SLAScheduleConfigStr, configID)
	configID = ds.PrependToDataID(configID, "slaConfigType")
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(ssd.server, tenantID)
	dataContainer := &metmod.SLAScheduleConfig{}
	if err := deleteDataFromCouch(tenantDBName, configID, &dataContainer, metmod.SLAScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", metmod.SLAScheduleConfigStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}
func (ssd *SchedulerServiceDatastoreCouchDB) GetScheduleConfig(tenantID string, configID string) (*metmod.SLAScheduleConfig, error) {
	logger.Log.Debugf("Fetching %s: %s\n", metmod.SLAScheduleConfigStr, configID)
	configID = ds.PrependToDataID(configID, "slaConfigType")
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(ssd.server, tenantID)
	dataContainer := &metmod.SLAScheduleConfig{}
	if err := getDataFromCouch(tenantDBName, configID, &dataContainer, metmod.SLAScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", metmod.SLAScheduleConfigStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}
