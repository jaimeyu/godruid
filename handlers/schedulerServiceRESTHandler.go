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
			Name:        "CreateScheduleConfig",
			Method:      "POST",
			Pattern:     apiV1Prefix + "schedules",
			HandlerFunc: result.CreateScheduleConfig,
		},
		server.Route{
			Name:        "UpdateScheduleConfig",
			Method:      "PUT",
			Pattern:     apiV1Prefix + "schedules",
			HandlerFunc: result.UpdateScheduleConfig,
		},
		server.Route{
			Name:        "GetScheduleConfig",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/schedules/{configID}",
			HandlerFunc: result.GetScheduleConfig,
		},
		server.Route{
			Name:        "DeleteScheduleConfig",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/schedules/{configID}",
			HandlerFunc: result.DeleteScheduleConfig,
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

func (ssh *SchedulerServiceRESTHandler) CreateScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	data := metmod.SLAScheduleConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", metmod.SLAScheduleConfigStr, models.AsJSONString(&data))

	result, err := ssh.schedulerDB.CreateScheduleConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", metmod.SLAScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.CreateScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateScheduleConfigStr, metmod.SLAScheduleConfigStr, "Created")
}

func (ssh *SchedulerServiceRESTHandler) UpdateScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	data := metmod.SLAScheduleConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", metmod.SLAScheduleConfigStr, models.AsJSONString(&data))

	result, err := ssh.schedulerDB.UpdateScheduleConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", metmod.SLAScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateScheduleConfigStr, metmod.SLAScheduleConfigStr, "Updated")
}

func (ssh *SchedulerServiceRESTHandler) GetScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)
	configID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", metmod.SLAScheduleConfigStr, configID)

	result, err := ssh.schedulerDB.GetScheduleConfig(tenantID, configID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", metmod.SLAScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.GetSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetSLAReportStr, metmod.SLAScheduleConfigStr, "Retrieved")
}

func (ssh *SchedulerServiceRESTHandler) DeleteScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)
	configID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", metmod.SLAScheduleConfigStr, configID)

	result, err := ssh.schedulerDB.DeleteScheduleConfig(tenantID, configID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", metmod.SLAScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteScheduleConfigStr, metmod.SLAScheduleConfigStr, "Deleted")
}
