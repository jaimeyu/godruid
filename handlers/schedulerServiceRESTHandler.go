package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	metmod "github.com/accedian/adh-gather/models/metrics"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
)

type SchedulerServiceRESTHandler struct {
	schedulerDB db.SchedulerServiceDatastore
	routes      []server.Route
}

func CreateSchedulerServiceRESTHandler() *SchedulerServiceRESTHandler {
	result := new(SchedulerServiceRESTHandler)

	// Setup the DB implementation based on configuration
	sdb, err := getSchedulerServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate SchedulerServiceRESTHandler: %s", err.Error())
	}
	result.schedulerDB = sdb

	result.routes = []server.Route{
		server.Route{
			Name:        "CreateReportScheduleConfig",
			Method:      "POST",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs",
			HandlerFunc: result.CreateReportScheduleConfig,
		},
		server.Route{
			Name:        "UpdateReportScheduleConfig",
			Method:      "PUT",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs",
			HandlerFunc: result.UpdateReportScheduleConfig,
		},
		server.Route{
			Name:        "GetReportScheduleConfig",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs/{configID}",
			HandlerFunc: result.GetReportScheduleConfig,
		},
		server.Route{
			Name:        "DeleteReportScheduleConfig",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs/{configID}",
			HandlerFunc: result.DeleteReportScheduleConfig,
		},
		server.Route{
			Name:        "GetAllReportScheduleConfigs",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-config-list",
			HandlerFunc: result.GetAllReportScheduleConfigs,
		},
	}

	return result
}

func (ssh *SchedulerServiceRESTHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range ssh.routes {
		logger.Log.Debugf("Registering endpoint: %v", route)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

func getSchedulerServiceDatastore() (db.SchedulerServiceDatastore, error) {
	cfg := gather.GetConfig()
	dbType := gather.DBImpl(cfg.GetInt(gather.CK_args_schedulerdb_impl.String()))
	switch dbType {
	case gather.COUCH:
		logger.Log.Debug("SchedulerService DB is using CouchDB Implementation")
		return couchDB.CreateSchedulerServiceDAO()
	case gather.MEM:
		logger.Log.Debug("SchedulerService DB is using InMemory Implementation")
		// TODO return inMemory.CreateSchedulerServiceDAO()
	}

	return nil, errors.New("No DB implementation provided for Scheduler Service. Check configuration")
}

func (ssh *SchedulerServiceRESTHandler) CreateReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	data := metmod.ReportScheduleConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", metmod.ReportScheduleConfigStr, models.AsJSONString(&data))

	result, err := ssh.schedulerDB.CreateReportScheduleConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.CreateReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Created")
}

func (ssh *SchedulerServiceRESTHandler) UpdateReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	data := metmod.ReportScheduleConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", metmod.ReportScheduleConfigStr, models.AsJSONString(&data))

	result, err := ssh.schedulerDB.UpdateReportScheduleConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Updated")
}

func (ssh *SchedulerServiceRESTHandler) GetReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)
	configID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", metmod.ReportScheduleConfigStr, configID)

	result, err := ssh.schedulerDB.GetReportScheduleConfig(tenantID, configID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.GetSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetSLAReportStr, metmod.ReportScheduleConfigStr, "Retrieved")
}

func (ssh *SchedulerServiceRESTHandler) GetAllReportScheduleConfigs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s list for Tenant %s", metmod.ReportScheduleConfigStr, tenantID)

	result, err := ssh.schedulerDB.GetAllReportScheduleConfigs(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllTenantUserStr, metmod.ReportScheduleConfigStr, "Retrieved list of")
}

func (ssh *SchedulerServiceRESTHandler) DeleteReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)
	configID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", metmod.ReportScheduleConfigStr, configID)

	result, err := ssh.schedulerDB.DeleteReportScheduleConfig(tenantID, configID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Deleted")
}
