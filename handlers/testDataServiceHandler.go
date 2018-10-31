package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mholt/archiver"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	cr "crypto/rand"
	"crypto/x509"
	"encoding/pem"
	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/datastore/inMemory"
	"github.com/accedian/adh-gather/gather"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/server"
	wr "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	unableToReadRequestStr     = "Unable to read Test Data content"
	missingTenantNameStr       = "Missing tenantName field"
	missingTenantIDStr         = "Missing tenantID field"
	missingDomainSetStr        = "Missing domainSet field"
	missingMonObjID            = "Missing Monitored Object id field"
	missingMonObjActuatorName  = "Missing Monitored Object actuatorName field"
	missingMonObjReflectorName = "Missing Monitored Object reflectorName field"
	missingMonObjObjectName    = "Missing Monitored Object objectName field"
	moIDStr                    = "id"
	moActuatorNameStr          = "actuatorName"
	moReflectorNameStr         = "reflectorName"
	moObjectNameStr            = "objectName"
	domainSLAReportName        = "Weekly Domain SLA Report"

	millisecondsPerWeek                         = 1000 * 60 * 60 * 24 * 7
	totalMonitoredObjectCountForDomainSLAReport = 25
	domainSlaReportBucketCount                  = 150
	domainSlaReportBucketDurationInMilliseconds = millisecondsPerWeek / domainSlaReportBucketCount
	domainSlaReportBucketDurationInMinutes      = domainSlaReportBucketDurationInMilliseconds / (1000 * 60)

	populateTestDataStr                 = "populate_test_data"
	populateTestDataBulkRandomizedMOStr = "populate_test_data_bulk_randomize_MO"
	populateTestDataIntoDruidStr        = "populate_test_data_into_druid"
	migrateMOMetadatadStr               = "migrate_mo_metadata"
	purgeDBStr                          = "purge_db"
	generateSLAReportStr                = "gen_sla_report"
	getDocsByTypeStr                    = "get_docs_by_type"
	insertTenViewsStr                   = "insert_tenant_views"
	signCSRStr                          = "sign_csr"
	downloadRRStr                       = "download_roadrunner"

	stringGeneratorCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// TestDataServiceHandler - handler for all APIs related to test data provisioning.
type TestDataServiceHandler struct {
	routes []server.Route
	// adminDB  db.AdminServiceDatastore
	tenantDB db.TenantServiceDatastore
	pouchDB  db.PouchDBPluginServiceDatastore
	grpcSH   *GRPCServiceHandler
	testDB   db.TestDataServiceDatastore
}

// CreateTestDataServiceHandler - generates a TestDataServiceHandler to handle all test
// data provisioning related APIs.
func CreateTestDataServiceHandler() *TestDataServiceHandler {
	result := new(TestDataServiceHandler)

	// Seteup the DB implementation based on configuration
	// db, err := getPouchDBPluginServiceDatastore()
	// if err != nil {
	// 	logger.Log.Fatalf("Unable to instantiate PouchDBPluginServiceHandler: %s", err.Error())
	// }
	// result.pouchPluginDB = db

	result.routes = []server.Route{
		server.Route{
			Name:        "PopulateTestData",
			Method:      "POST",
			Pattern:     "/test-data",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.PopulateTestData),
		},

		server.Route{
			Name:        "PopulateTestDataBulkRandomizedMO",
			Method:      "POST",
			Pattern:     "/test-data/bulkRandomizedMO",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.PopulateTestDataBulkRandomizedMO),
		},

		server.Route{
			Name:        "PurgeDB",
			Method:      "DELETE",
			Pattern:     "/test-data/{dbname}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.PurgeDB),
		},

		server.Route{
			Name:        "GenerateHistoricalDomainSLAReports",
			Method:      "POST",
			Pattern:     "/test-data/domain-sla-reports",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.GenerateHistoricalDomainSLAReports),
		},

		server.Route{
			Name:        "GetAllDocsByType",
			Method:      "GET",
			Pattern:     "/test-data/{dbname}/{datatype}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.GetAllDocsByType),
		},

		server.Route{
			Name:        "InsertTenantViews",
			Method:      "PUT",
			Pattern:     "/test-data/tenant-views/{dbname}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.InsertTenantViews),
		},
		server.Route{
			Name:    "PopulateTestDataIntoDruid",
			Method:  "POST",
			Pattern: "/test-data/populate-druid",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.
				PopulateTestDataIntoDruid),
		},
		server.Route{
			Name:    "MigrateMetadata",
			Method:  "POST",
			Pattern: "/test-data/MigrateMetadata",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.
				MigrateMetadata),
		},
		server.Route{
			Name:        "SignCSR",
			Method:      "POST",
			Pattern:     "/distribution/sign-csr",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.SignCSR),
		},
		server.Route{
			Name:        "DownloadRoadrunner",
			Method:      "GET",
			Pattern:     "/distribution/download-roadrunner",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.DownloadRoadrunner),
		},
	}

	// Wire up the datastore impls
	// admindb, err := getAdminServiceDatastore()
	// if err != nil {
	// 	logger.Log.Fatalf("Unable to instantiate TestDataServiceHandler: %s", err.Error())
	// }
	// result.adminDB = admindb

	tenantdb, err := GetTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TestDataServiceHandler: %s", err.Error())
	}
	result.tenantDB = tenantdb

	pouchdb, err := getPouchDBPluginServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TestDataServiceHandler: %s", err.Error())
	}
	result.pouchDB = pouchdb

	testdb, err := getTestDataServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TestDataServiceHandler: %s", err.Error())
	}
	result.testDB = testdb

	result.grpcSH = CreateCoordinator()

	return result
}

func getTestDataServiceDatastore() (db.TestDataServiceDatastore, error) {
	cfg := gather.GetConfig()
	dbType := gather.DBImpl(cfg.GetInt(gather.CK_args_testdatadb_impl.String()))
	switch dbType {
	case gather.COUCH:
		logger.Log.Debug("TestDataService DB is using CouchDB Implementation")
		return couchDB.CreateTestDataServiceDAO()
	case gather.MEM:
		logger.Log.Debug("TestDataService DB is using InMemory Implementation")
		return inMemory.CreateTestDataServiceDAO()
	}

	return nil, errors.New("No DB implementation provided for Admin Service. Check configuration")
}

// RegisterAPIHandlers - connects the endpoints of the multiplexor to the functions that will
// handle the calls.
func (tsh *TestDataServiceHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range tsh.routes {
		logger.Log.Debugf("Registering endpoint: %v", route)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

// PopulateTestData - adds test data to the deployment.
func (tsh *TestDataServiceHandler) PopulateTestData(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the contents of the request:
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s content: %s", "Test Data", err.Error())
		reportError(w, startTime, "400", populateTestDataStr, msg, http.StatusBadRequest)
		return
	}

	err = validatePopulateTestDataRequest(requestBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to validate %s content: %s", "Test Data", err.Error())
		reportError(w, startTime, "400", populateTestDataStr, msg, http.StatusBadRequest)
		return
	}

	// Create a Tenant metadata using the provided name:
	tenantName := requestBody["tenantName"].(string)
	desc, err := generateTenantDescriptor(tenantName)
	if err != nil {
		msg := fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error())
		reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
		return
	}
	tenantDescriptor, err := tsh.grpcSH.CreateTenant(nil, desc)
	if err != nil {
		msg := fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error())
		reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
		return
	}

	dataIDForTenant := tenantDescriptor.GetXId()

	// Fetch the tenant Meta so that teh default threshold profile id is available:
	tenantMeta, err := tsh.grpcSH.GetTenantMeta(nil, &wr.StringValue{Value: tenantDescriptor.GetXId()})
	if err != nil {
		msg := fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error())
		reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
		return
	}

	// Add in the Domain Objects.
	tenantDomains := requestBody["domainSet"].([]interface{})
	createdDomainIDSet := []string{}
	for _, domainName := range tenantDomains {
		dom, err := generateTenantDomain(domainName.(string), dataIDForTenant, tenantMeta.GetData().GetDefaultThresholdProfile())
		if err != nil {
			msg := fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error())
			reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
			return
		}
		domain, err := tsh.grpcSH.CreateTenantDomain(nil, dom)
		if err != nil {
			msg := fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error())
			reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
			return
		}

		// For the IDs used as references inside other objects, need to strip off the 'domain_2_'
		// as this is just relational pouch adaption:
		domainIDParts := strings.Split(domain.GetXId(), "_")
		createdDomainIDSet = append(createdDomainIDSet, domainIDParts[len(domainIDParts)-1])
	}

	// If there are Monitored Objects, add them as well and map them to Domains
	if requestBody["monObjSet"] != nil {
		monitoredObjects := requestBody["monObjSet"].([]interface{})
		if len(monitoredObjects) != 0 {
			// There are monitored objects, try to provision them.
			for _, monObj := range monitoredObjects {
				obj := monObj.(map[string]interface{})
				if err = validateMonitoredObject(obj); err != nil {
					logger.Log.Errorf("Not creating Monitored Object %v: %s", monObj, err.Error())
				}

				mo, err := generateMonitoredObject(obj[moIDStr].(string), dataIDForTenant, obj[moActuatorNameStr].(string), obj[moReflectorNameStr].(string), obj[moObjectNameStr].(string), createdDomainIDSet)
				if err != nil {
					msg := fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error())
					reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
					return
				}
				_, err = tsh.grpcSH.CreateMonitoredObject(nil, mo)
				if err != nil {
					msg := fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error())
					reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// All operations successful, just return the Tenant Descriptor as success:
	response, err := json.Marshal(tenantDescriptor)
	if err != nil {
		msg := fmt.Sprintf("Error generating response: %s", err.Error())
		reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
		return
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", populateTestDataStr)
	fmt.Fprintf(w, string(response))
}

// PurgeDB - deletes all records from a DB but does not destroy the DB.
func (tsh *TestDataServiceHandler) PurgeDB(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)

	datatypeFilter := r.URL.Query().Get("datatype")

	if len(datatypeFilter) != 0 {
		// doing the delete by datatype filter, not all docs.
		allDocsByType, err := tsh.testDB.GetAllDocsByDatatype(dbName, datatypeFilter)
		if err != nil {
			msg := fmt.Sprintf("Unable to purge DB %s of documents of type %s: %s", dbName, datatypeFilter, err.Error())
			reportError(w, startTime, "500", purgeDBStr, msg, http.StatusInternalServerError)
			return
		}

		if len(allDocsByType) != 0 {
			docsToDelete := make([]map[string]interface{}, 0)
			for _, doc := range allDocsByType {
				docID := doc["_id"].(string)
				docRev := doc["_rev"].(string)
				docsToDelete = append(docsToDelete, map[string]interface{}{"_id": docID, "_rev": docRev, "_deleted": true})
			}

			deleteBody := map[string]interface{}{"docs": docsToDelete}
			logger.Log.Debugf("Attempting to delete the following from DB %s: %v", dbName, docsToDelete)

			_, err = tsh.pouchDB.BulkDBUpdate(dbName, deleteBody)
			if err != nil {
				msg := fmt.Sprintf("Unable to purge DB %s: %s", dbName, err.Error())
				reportError(w, startTime, "500", purgeDBStr, msg, http.StatusInternalServerError)
				return
			}
		}

		// Purge complete, send back success details:
		fmt.Fprintf(w, fmt.Sprintf("Successfully purged all %s records from DB %s", datatypeFilter, dbName))
		mon.TrackAPITimeMetricInSeconds(startTime, "200", populateTestDataStr)
		return
	}

	// Get a list of all documents from the DB:
	docs, err := tsh.pouchDB.GetAllDBDocs(dbName, map[string]interface{}{})
	if err != nil {
		msg := fmt.Sprintf("Unable to purge DB %s: %s", dbName, err.Error())
		reportError(w, startTime, "500", purgeDBStr, msg, http.StatusInternalServerError)
		return
	}

	if docs["rows"] != nil {
		// There are documents to delete, build up a bulk delete request body for all of them
		docList := docs["rows"].([]interface{})
		docsToDelete := make([]map[string]interface{}, 0)
		for _, doc := range docList {
			docObj := doc.(map[string]interface{})
			docID := docObj["id"].(string)
			docRev := docObj["value"].(map[string]interface{})["rev"].(string)
			docsToDelete = append(docsToDelete, map[string]interface{}{"_id": docID, "_rev": docRev, "_deleted": true})
		}

		deleteBody := map[string]interface{}{"docs": docsToDelete}
		logger.Log.Debugf("Attempting to delete the following from DB %s: %v", dbName, docsToDelete)

		_, err = tsh.pouchDB.BulkDBUpdate(dbName, deleteBody)
		if err != nil {
			msg := fmt.Sprintf("Unable to purge DB %s: %s", dbName, err.Error())
			reportError(w, startTime, "500", purgeDBStr, msg, http.StatusInternalServerError)
			return
		}
	}

	// Purge complete, send back success details:
	mon.TrackAPITimeMetricInSeconds(startTime, "200", purgeDBStr)
	fmt.Fprintf(w, "Successfully purged all records from DB "+dbName)
}

// GenerateHistoricalDomainSLAReports - generates weekly SLA reports for a Tenant based on the Domains provisioned for that Tenant.
func (tsh *TestDataServiceHandler) GenerateHistoricalDomainSLAReports(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	// Get the contents of the request:
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s content: %s", "Generate SLA Report", err.Error())
		reportError(w, startTime, "400", generateSLAReportStr, msg, http.StatusBadRequest)
		return
	}

	err = validateGenerateDomainSLAReportRequest(requestBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to validate %s request: %s", "Generate SLA Report", err.Error())
		reportError(w, startTime, "400", generateSLAReportStr, msg, http.StatusBadRequest)
		return
	}

	// Create a Tenant metadata using the provided name:
	tenantID := requestBody["tenantId"].(string)
	numReportsPerDomainObj := requestBody["numReportsPerDomain"]
	numReportsPerDomain := float64(5)
	if numReportsPerDomainObj != nil {
		numReportsPerDomain = numReportsPerDomainObj.(float64)
	}

	// Get the list of Domain Objects provisioned for this tenant:
	domainSetForTenant, err := tsh.grpcSH.GetAllTenantDomains(nil, &wr.StringValue{Value: tenantID})
	if err != nil {
		msg := fmt.Sprintf("Unable to generate %s content for Tenant %s: No Domain data provisioned", db.DomainSlaReportStr, tenantID)
		reportError(w, startTime, "500", generateSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	// Generate teh data for each domain for the number of reports and insert them as a bulk update:
	rand.Seed(time.Now().UTC().UnixNano())
	docsToInsert := make([]map[string]interface{}, 0)
	endReportTS := time.Now().Truncate(24*time.Hour).UnixNano() / 1000000 // Truncate a current timestamp to the beginning of this day.
	startReportTS := endReportTS - millisecondsPerWeek
	for i := float64(0); i < numReportsPerDomain; i++ {
		for _, domain := range domainSetForTenant.GetData() {
			// Add a new Report to the list of documents to insert
			slaReport, err := generateSLADomainReport(domain, startReportTS, endReportTS)
			if err != nil {
				msg := fmt.Sprintf("Unable to generate %s content for Tenant %s: %s", db.DomainSlaReportStr, tenantID, err.Error())
				reportError(w, startTime, "500", generateSLAReportStr, msg, http.StatusInternalServerError)
				return
			}
			docsToInsert = append(docsToInsert, slaReport)
		}

		// Update the timestamps for the previous week
		endReportTS = startReportTS
		startReportTS = startReportTS - millisecondsPerWeek
	}

	insertBody := map[string]interface{}{"docs": docsToInsert}
	// logger.Log.Debugf("Attempting to insert %d %ss for Tenant %s", len(docsToInsert), db.DomainSlaReportStr, tenantID)

	// Prepend the tenant id with the known prefix otherwise the tenant DB will not be found.
	tenantIDForUpdate := db.PrependToDataID(tenantID, string(admmod.TenantType))
	_, err = tsh.pouchDB.BulkDBUpdate(tenantIDForUpdate, insertBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to insert %ss for Tenant %s: %s", db.DomainSlaReportStr, tenantID, err.Error())
		reportError(w, startTime, "500", generateSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", generateSLAReportStr)
	fmt.Fprintf(w, fmt.Sprintf("Successfully generated all %ss for tenant %s", db.DomainSlaReportStr, tenantID))
}

// GetAllDocsByType - retrieves all docs by type.
func (tsh *TestDataServiceHandler) GetAllDocsByType(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)
	docType := getDBFieldFromRequest(r, 3)

	if len(dbName) == 0 || len(docType) == 0 {
		msg := fmt.Sprintf("Unable to retrieve documents without a DB name and Document type")
		reportError(w, startTime, "500", getDocsByTypeStr, msg, http.StatusInternalServerError)
		return
	}

	allDocsByType, err := tsh.testDB.GetAllDocsByDatatype(dbName, docType)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve documents of type %s from DB %s: %s", docType, dbName, err.Error())
		reportError(w, startTime, "500", getDocsByTypeStr, msg, http.StatusInternalServerError)
		return
	}

	// All operations successful, just return the Tenant Descriptor as success:
	response, err := json.Marshal(allDocsByType)
	if err != nil {
		msg := fmt.Sprintf("Error generating response: %s", err.Error())
		reportError(w, startTime, "500", getDocsByTypeStr, msg, http.StatusInternalServerError)
		return
	}
	mon.TrackAPITimeMetricInSeconds(startTime, "200", getDocsByTypeStr)
	fmt.Fprintf(w, string(response))
}

// InsertTenantViews - inserts tenant views.
func (tsh *TestDataServiceHandler) InsertTenantViews(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 3)

	if len(dbName) == 0 {
		msg := fmt.Sprintf("Unable to insert documents without a DB name")
		reportError(w, startTime, "500", insertTenViewsStr, msg, http.StatusInternalServerError)
		return
	}

	result := tsh.testDB.InsertTenantViews(dbName)

	mon.TrackAPITimeMetricInSeconds(startTime, "200", insertTenViewsStr)
	fmt.Fprintf(w, string(result))
}

// PopulateTestDataBulkRandomizedMO - Generates a set of monitored objects with randomized values according to query parameters provided in the incoming rest request
// Supported query parameters are:
//			count: the number of desires monitored objects to be generated. Defaults to 1
//			batchSize: the number of monitored objects that should be placed in a batch to be sent to the db. Defaults to 1000
//			tenant: the ID of the tenant to associate the monitored objects with
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TestDataServiceHandler) PopulateTestDataBulkRandomizedMO(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	queryParams := r.URL.Query()

	var (
		moRequestCount  uint64 = 1
		moRequestTenant string
		batchSize       uint64 = 1000
		err             error
	)

	// Defines the number of monitored objects that we want to create. Defaults to 1
	if queryParams["count"] != nil {
		moRequestCount, err = strconv.ParseUint(queryParams["count"][0], 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Unacceptable value provided for monitored object count: %s", err.Error())
			reportError(w, startTime, "400", populateTestDataBulkRandomizedMOStr, msg, http.StatusBadRequest)
			return
		}
	}

	// Defines the number of monitored objects we want to place in a batch request to the datastore
	if queryParams["batchSize"] != nil {
		batchSize, err = strconv.ParseUint(queryParams["batchSize"][0], 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Unacceptable value provided for monitored object batch size: %s", err.Error())
			reportError(w, startTime, "400", populateTestDataBulkRandomizedMOStr, msg, http.StatusBadRequest)
			return
		}
	}

	// Defines the tenant to place the monitored objects against. This is a required field
	if len(queryParams["tenant"]) == 0 {
		msg := fmt.Sprintf("Tenant ID must be provided")
		reportError(w, startTime, "400", populateTestDataBulkRandomizedMOStr, msg, http.StatusBadRequest)
		return
	} else {
		moRequestTenant = queryParams["tenant"][0]
	}

	// Retrieve all the domains associated with the tenant in order to associate a random subset of them against the monitored object
	domainSet, err := tsh.tenantDB.GetAllTenantDomains(moRequestTenant)
	if err != nil {
		msg := fmt.Sprintf("(Unable to retrieve domain set for tenant %s: %s", moRequestTenant, err.Error())
		reportError(w, startTime, "500", populateTestDataBulkRandomizedMOStr, msg, http.StatusInternalServerError)
		return
	}
	domainSetIDs := make([]string, len(domainSet))
	for i, d := range domainSet {
		domainSetIDs[i] = d.ID
	}

	// Create our initial empty bucket of MOs
	moData := make([]*tenmod.MonitoredObject, min(batchSize, moRequestCount))
	expectedSize := batchSize

	for i := uint64(0); i < moRequestCount; i++ {
		// We have hit the limit of our bucket size
		if i%expectedSize == 0 && i != 0 {
			// Actually attempt to insert the batch of MOs to the datastore
			_, err = tsh.tenantDB.BulkInsertMonitoredObjects(moRequestTenant, moData)
			if err != nil {
				msg := fmt.Sprintf("Unable to provision monitored object content for tenant %s: %s", moRequestTenant, err.Error())
				reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
			}
			logger.Log.Debugf("Finished batch: %d", i)

			// Only bother creating a new array if we know that we are not at the end of our request count size
			if i != (moRequestCount - 1) {
				// We may potentially adjust our expected size if we have less MOs that need to be generated that the desired batch size
				expectedSize = min(batchSize, (moRequestCount - i))
				moData = make([]*tenmod.MonitoredObject, expectedSize)
				// We know that we will be looping one more time so set the first value of the array
				moData[0] = generateRandomMonitoredObject(moRequestTenant, domainSetIDs)
			}
		} else {
			moData[i%expectedSize] = generateRandomMonitoredObject(moRequestTenant, domainSetIDs)
		}
	}

	// Insert any remaining MOs to the datastore
	if len(moData) != 0 {
		_, err = tsh.tenantDB.BulkInsertMonitoredObjects(moRequestTenant, moData)
		if err != nil {
			msg := fmt.Sprintf("Unable to provision monitored object content for tenant %s: %s", moRequestTenant, err.Error())
			reportError(w, startTime, "500", populateTestDataStr, msg, http.StatusInternalServerError)
		}
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", populateTestDataBulkRandomizedMOStr)
	fmt.Fprintf(w, "Success")
}

func generateSLADomainReport(domain *pb.TenantDomain, reportStartTS int64, reportEndTS int64) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	typeName := string(db.DomainSlaReportType)
	uuid := uuid.NewV4()
	result["_id"] = typeName + db.PouchDBIdBridgeStr + uuid.String()

	resultContent := map[string]interface{}{}
	resultContent["datatype"] = typeName
	resultContent["reportName"] = domainSLAReportName
	resultContent["domain"] = domain.XId
	resultContent["objectCount"] = 25
	resultContent["thresholdName"] = "Default"

	reportRange := map[string]interface{}{}
	reportRange["start"] = reportStartTS
	reportRange["end"] = reportEndTS
	resultContent["range"] = reportRange

	// Set compliance rate somewhere between 94% and 96%
	slaComplianceRate := 940 + rand.Intn(960-940)

	values := make([]map[string]interface{}, 0)
	for i := 0; i < domainSlaReportBucketCount; i++ {
		// Make a value for each timestamp in the range of the report
		values = append(values, generateSlaRangeValue(reportStartTS, slaComplianceRate))

		// Increment the bucket timestamp
		reportStartTS = reportStartTS + domainSlaReportBucketDurationInMilliseconds
	}

	resultContent["buckets"] = values
	result["data"] = resultContent

	return result, nil
}

func generateSlaRangeValue(timestamp int64, complianceRate int) map[string]interface{} {
	result := map[string]interface{}{}

	result["timestamp"] = timestamp

	metrics := map[string]interface{}{}
	metrics["jitterP95"] = generateValueInRangeIfPassesTest(1, 80, complianceRate)
	metrics["delayP95"] = generateValueInRangeIfPassesTest(1, 80, complianceRate)
	metrics["packetsLostPct"] = generateValueInRangeIfPassesTest(1, 80, complianceRate)
	result["metrics"] = metrics

	return result
}

func generateValueInRangeIfPassesTest(minValue int, maxValue int, testRate int) int {
	roll := rand.Intn(1000)
	if roll <= testRate {
		// failed the test, just return 0
		return 0
	}

	// Weight the result so that it skews towards 0
	roll = rand.Intn(100)
	if roll < 80 {
		maxValue = 20
	} else if roll < 90 {
		maxValue = 60
	}

	// Generate a number for return in the range of desired values.
	return minValue + rand.Intn(maxValue-minValue)
}

func validatePopulateTestDataRequest(requestBody map[string]interface{}) error {
	if requestBody["tenantName"] == nil || len(requestBody["tenantName"].(string)) == 0 {
		return fmt.Errorf("%s: %s", unableToReadRequestStr, missingTenantNameStr)
	}

	if requestBody["domainSet"] == nil {
		return fmt.Errorf("%s: %s", unableToReadRequestStr, missingDomainSetStr)
	}

	return nil
}

func validateGenerateDomainSLAReportRequest(requestBody map[string]interface{}) error {
	if requestBody["tenantId"] == nil || len(requestBody["tenantId"].(string)) == 0 {
		return fmt.Errorf("%s: %s", unableToReadRequestStr, missingTenantIDStr)
	}

	return nil
}

func validateMonitoredObject(monObj map[string]interface{}) error {
	if monObj[moIDStr] == nil || len(monObj[moIDStr].(string)) == 0 {
		return fmt.Errorf("%s: %s", unableToReadRequestStr, missingMonObjID)
	}
	if monObj[moActuatorNameStr] == nil || len(monObj[moActuatorNameStr].(string)) == 0 {
		return fmt.Errorf("%s: %s", unableToReadRequestStr, missingMonObjActuatorName)
	}
	if monObj[moReflectorNameStr] == nil || len(monObj[moReflectorNameStr].(string)) == 0 {
		return fmt.Errorf("%s: %s", unableToReadRequestStr, missingMonObjReflectorName)
	}
	if monObj[moObjectNameStr] == nil || len(monObj[moObjectNameStr].(string)) == 0 {
		return fmt.Errorf("%s: %s", unableToReadRequestStr, missingMonObjObjectName)
	}

	return nil
}

func generateTenantDescriptor(name string) (*pb.TenantDescriptor, error) {
	result := pb.TenantDescriptor{}

	tenantStr := string(admmod.TenantType)
	result.Data = &pb.TenantDescriptorData{}
	result.Data.Name = name
	result.Data.UrlSubdomain = strings.ToLower(name) + ".npav.accedian.net"
	result.Data.State = 2
	result.Data.Datatype = tenantStr
	result.Data.CreatedTimestamp = db.MakeTimestamp()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	return &result, nil
}

func generateTenantUser(name string, tenantID string) (*pb.TenantUser, error) {
	result := pb.TenantUser{}

	tenantUserStr := string(tenmod.TenantUserType)

	result.Data = &pb.TenantUserData{}
	result.Data.Datatype = tenantUserStr
	result.Data.Username = "admin@" + strings.ToLower(name) + ".com"
	result.Data.TenantId = tenantID
	result.Data.UserVerified = true
	result.Data.State = 2
	result.Data.SendOnboardingEmail = false
	result.Data.Password = "admin"
	result.Data.OnboardingToken = "anonboardingtokenforthedefaulttenantuser"
	result.Data.CreatedTimestamp = db.MakeTimestamp()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	result.XId = db.GenerateID(result.Data, tenantUserStr)

	return &result, nil
}

func generateTenantDomain(name string, tenantID string, defaultThreshPrf string) (*pb.TenantDomain, error) {
	result := pb.TenantDomain{}

	tenantDomainStr := string(tenmod.TenantDomainType)
	result.Data = &pb.TenantDomainData{}
	result.Data.Datatype = tenantDomainStr
	result.Data.Name = name
	result.Data.TenantId = tenantID
	result.Data.Color = "#4EC5C1"

	// Only use the hash portion of the id for the reference since the thresholdProfile_2_ is all
	// relational pouch tagging.
	result.Data.CreatedTimestamp = db.MakeTimestamp()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	result.XId = db.GenerateID(result.Data, tenantDomainStr)

	return &result, nil
}

func generateMonitoredObject(id string, tenantID string, actuatorName string, reflectorName string, objectName string, domainIDSet []string) (*pb.MonitoredObject, error) {
	result := pb.MonitoredObject{}

	tenantMonObjStr := string(tenmod.TenantMonitoredObjectType)
	result.Data = &pb.MonitoredObjectData{}
	result.Data.Datatype = tenantMonObjStr

	result.Data.Id = id
	result.Data.TenantId = tenantID
	result.Data.ActuatorName = actuatorName
	result.Data.ActuatorType = string(tenmod.AccedianVNID)
	result.Data.ReflectorName = reflectorName
	result.Data.ReflectorType = string(tenmod.AccedianVNID)
	result.Data.ObjectName = objectName
	result.Data.ObjectType = string(tenmod.TwampPE)

	// To provision the DomainSet, need to obtain a subset of the passed in domain set.
	result.Data.DomainSet = generateRandomStringArray(domainIDSet)

	result.XId = db.GenerateID(result.Data, tenantMonObjStr)

	return &result, nil
}

// Generates monitored object with random field values associated with a specific tenant and provided domainSet
// Params:
//		tenantID: the tenant to associate this monitored object with
//		domainSet: the domains to consider for random association with this monitored object
// Returns:
//		A monitored object populate with random field values
func generateRandomMonitoredObject(tenantID string, domainSet []string) *tenmod.MonitoredObject {
	result := tenmod.MonitoredObject{DomainSet: generateRandomStringArray(domainSet)}

	// Generate basic field values randomly and associate the appropriate tenant with the MO
	result.TenantID = tenantID
	result.ActuatorName = generateRandomString(10)
	result.ActuatorType = string(tenmod.AccedianVNID)
	result.ReflectorName = generateRandomString(10)
	result.ReflectorType = string(tenmod.AccedianVNID)
	result.ObjectName = generateRandomEnodeB()
	result.ObjectType = string(tenmod.TwampPE)
	//result.MonitoredObjectID = strings.Join([]string{result.ObjectName, result.ActuatorName, result.ReflectorName, generateRandomString(10)}, "-")
	result.MonitoredObjectID = strings.Join([]string{result.ObjectName}, "-")

	// Generate random meta data
	result.Meta = generateRandomMeta()

	return &result
}

// Generates random meta data for monitored objects
func generateRandomMeta() map[string]string {
	num := rand.Intn(8)
	// Generate random meta data
	regions := []string{"london", "tokyo", "toronto", "montreal", "vancouver", "calgary", "regina", "new york", "las vegas", "boston", "chicago", "winnepeg", "washington", "hamilton", "paris", "lyon"}
	colors := []string{"black", "white", "orange", "blue", "green", "red", "purple", "gold", "yellow", "brown", "aqua"}
	superheroes := []string{"superman", "batman", "ironman", "spider-man", "ant-man", "aquaman", "wonderwoman", "batwoman", "birdman", "catwoman", "invisible woman"}
	meta := make(map[string]string)

	// The following inserts monitored objects with blanks, so schemas are not always the same across all objects.
	rI := num % len(regions)
	if rI != 0 {
		meta["region"] = regions[rI]
	}

	rI = num % len(colors)
	if rI != 0 {
		meta["colors"] = colors[num%len(colors)]
	}

	rI = num % len(superheroes)
	if rI != 0 {
		meta["superheroes"] = superheroes[num%len(superheroes)]
	}

	return meta
}

// Compares two uint64 and returns the smallest one
// Params:
//		a: the first number to compare
//		b: the second number to compare
// Returns:
//		the smaller of the two provided numbers
func min(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// Generates a random subset of the provided string array
// Params:
//		stringset: the full set of strings that we wish to take a subset of
// Returns:
//		a random subset of the provided string array
func generateRandomStringArray(stringSet []string) []string {
	// If the passed in set is nil then return a nil value
	if stringSet == nil {
		return nil
	}
	setLength := len(stringSet)
	// If we have zero length array then pass back a zero length slice
	if setLength == 0 {
		return stringSet[:0]
	}
	take := rand.Intn(setLength + 1)
	// If we are meant to only take a single string then choose an arbitrary one
	if take == 1 {
		randVal := rand.Intn(setLength)
		return stringSet[randVal : randVal+1]
	}

	// Otherwise return a subset from 0 to the randomly chosen size of the new slice
	return stringSet[:take]
}

// Generate a string with arbitrary characters based on the provided length
// Params:
//		length: the size of the desired randomized string
// Returns:
//		a randomly generated string of the provided length
func generateRandomString(length int) string {
	generated := make([]byte, length)
	for i := range generated {
		generated[i] = stringGeneratorCharset[rand.Int63()%int64(len(stringGeneratorCharset))]
	}
	return string(generated)
}

// Generate a string that uses a simulated BYT enode name.
func generateRandomEnodeB() string {
	offset := rand.Int()
	//generated[i] = stringGeneratorCharset[rand.Int63()%int64(len(stringGeneratorCharset))]
	// Enode B template
	// E84717_WST_H_ENB_MON_ZIPEC_AF22_admin_17-01-18
	regions := []string{"WST", "EST", "NOE", "SWT", "CTA", "IDF", "MED"}
	vendors := []string{"H", "E", "C"}
	// E84717_WST_H_ENB_MON_ZIPEC_AF22_admin_17-01-18
	objname_template := "E%d_%s_%s_ENB_MON_ZIPEC_AF22_admin_%s"
	t := time.Now()
	generatedName := fmt.Sprintf(objname_template, offset, regions[offset%len(regions)], vendors[offset%len(vendors)], t.Format("2006-01-02"))
	return generatedName
}

// Constants for the kafka simulator
const (
	kafkaTopic     = "npav-ts-metrics"
	broker         = "kafka:9092"
	timestampWord  = "{{TIMESTAMP}}"
	tenantWord     = "{{TENANTID}}"
	monObjIDWord   = "{{MONOBJID}}"
	monObjNameWord = "{{MONOBJNAME}}"
	// {"timestamp":1532979563973,"tenantId":"fc76af94-5804-450a-a922-d8957146321b","monitoredObjectId":"1473880047000-2026","sessionId":"2026","delayMin":7153,"delayMax":7385,"delayAvg":7222,"delayStdDevAvg":61,"delayP25":7175,"delayP50":7192,"delayP75":7282,"delayP95":7339,"delayPLo":7343,"delayPMi":7357,"delayPHi":7382,"delayVarMax":232,"delayVarAvg":69,"delayVarP25":22,"delayVarP50":39,"delayVarP75":129,"delayVarP95":186,"delayVarPLo":190,"delayVarPMi":204,"delayVarPHi":229,"jitterMin":0,"jitterMax":208,"jitterAvg":43,"jitterStdDev":48,"jitterP25":9,"jitterP50":23,"jitterP75":63,"jitterP95":140,"jitterPLo":154,"jitterPMi":166,"jitterPHi":181,"packetsLost":0,"packetsLostPct":0.0,"packetsMisordered":0,"packetsDuplicated":0,"packetsTooLate":0,"periodsLost":0,"lostBurstMin":0,"lostBurstMax":0,"packetsReceived":300,"bytesReceived":38400,"ipTOSMax":0,"ipTOSMin":0,"ttlMin":241,"ttlMax":241,"vlanPBitMin":0,"vlanPBitMax":0,"mos":4.409286022186279,"rValue":9.32E7,"deviceId":"demo1VCX-e5","direction":1,"objectType":"twamp-sf","throughputMin":0,"throughputMax":0,"throughputAvg":0,"duration":30000,"packetsSent":0,"domains":[],"monitoredObjectName":"1473880047000-2026","cleanStatus":1,"failedRules":[],"errorCode":0,"objectVendor":"accedian-twamp"}
	fauxDatatemplate     = `{"timestamp":{{TIMESTAMP}},"tenantId":"{{TENANTID}}","monitoredObjectId":"{{MONOBJID}}","sessionId":"2026","delayMin":7153,"delayMax":7385,"delayAvg":7222,"delayStdDevAvg":61,"delayP25":7175,"delayP50":7192,"delayP75":7282,"delayP95":7339,"delayPLo":7343,"delayPMi":7357,"delayPHi":7382,"delayVarMax":232,"delayVarAvg":69,"delayVarP25":22,"delayVarP50":39,"delayVarP75":129,"delayVarP95":186,"delayVarPLo":190,"delayVarPMi":204,"delayVarPHi":229,"jitterMin":0,"jitterMax":208,"jitterAvg":43,"jitterStdDev":48,"jitterP25":9,"jitterP50":23,"jitterP75":63,"jitterP95":140,"jitterPLo":154,"jitterPMi":166,"jitterPHi":181,"packetsLost":0,"packetsLostPct":0,"packetsMisordered":0,"packetsDuplicated":0,"packetsTooLate":0,"periodsLost":0,"lostBurstMin":0,"lostBurstMax":0,"packetsReceived":300,"bytesReceived":38400,"ipTOSMax":0,"ipTOSMin":0,"ttlMin":241,"ttlMax":241,"vlanPBitMin":0,"vlanPBitMax":0,"mos":4.409286022186279,"rValue":93200000,"deviceId":"demo1VCX-e5","direction":1,"objectType":"twamp-sf","throughputMin":0,"throughputMax":0,"throughputAvg":0,"duration":30000,"packetsSent":0,"domains":[],"monitoredObjectName":"{{MONOBJNAME}}","cleanStatus":1,"failedRules":[],"errorCode":0,"objectVendor":"accedian-twamp"}`
	fauxDataFailtemplate = `{"timestamp":{{TIMESTAMP}},"tenantId":"{{TENANTID}}","monitoredObjectId":"{{MONOBJID}}","sessionId":"2026","delayMin":60000,"delayMax":60000,"delayAvg":60000,"delayStdDevAvg":61,"delayP25":60000,"delayP50":60000,"delayP75":60000,"delayP95":60000,"delayPLo":60000,"delayPMi":60000,"delayPHi":60000,"delayVarMax":60000,"delayVarAvg":60000,"delayVarP25":60000,"delayVarP50":60000,"delayVarP75":60000,"delayVarP95":60000,"delayVarPLo":60000,"delayVarPMi":60000,"delayVarPHi":60000,"jitterMin":60000,"jitterMax":60000,"jitterAvg":60000,"jitterStdDev":60000,"jitterP25":60000,"jitterP50":60000,"jitterP75":60000,"jitterP95":60000,"jitterPLo":60000,"jitterPMi":60000,"jitterPHi":60000,"packetsLost":60000,"packetsLostPct":90,"packetsMisordered":0,"packetsDuplicated":0,"packetsTooLate":0,"periodsLost":0,"lostBurstMin":0,"lostBurstMax":0,"packetsReceived":60000,"bytesReceived":3840000,"ipTOSMax":0,"ipTOSMin":0,"ttlMin":241,"ttlMax":241,"vlanPBitMin":0,"vlanPBitMax":0,"mos":4.409286022186279,"rValue":93200000,"deviceId":"demo1VCX-e5","direction":1,"objectType":"twamp-sf","throughputMin":0,"throughputMax":0,"throughputAvg":0,"duration":30000,"packetsSent":0,"domains":[],"monitoredObjectName":"{{MONOBJNAME}}","cleanStatus":1,"failedRules":[],"errorCode":0,"objectVendor":"accedian-twamp"}`
)

// MigrateMetadata - Bulk operation to fix issues with metadata or to force new rules on metadata
func (tsh *TestDataServiceHandler) MigrateMetadata(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	queryParams := r.URL.Query()

	var (
		moRequestTenant string
		err             error
	)

	// Defines the tenant to place the monitored objects against. This is a required field
	if len(queryParams["tenant"]) == 0 {
		msg := fmt.Sprintf("Tenant ID must be provided")
		reportError(w, startTime, "400", migrateMOMetadatadStr, msg, http.StatusBadRequest)
		return
	} else {
		moRequestTenant = queryParams["tenant"][0]
	}

	mos, err := tsh.tenantDB.GetAllMonitoredObjectsIDs(moRequestTenant)
	if err != nil {
		msg := fmt.Sprintf("Could not get all MOs. %s", err.Error())
		reportError(w, startTime, "400", migrateMOMetadatadStr, msg, http.StatusBadRequest)
	}

	metas := make(map[string]string)
	for _, name := range mos {
		mo, err := tsh.tenantDB.GetMonitoredObject(moRequestTenant, name)
		if err != nil {
			msg := fmt.Sprintf("Could not get MO %s %s. %s", moRequestTenant, name, err.Error())
			reportError(w, startTime, "400", migrateMOMetadatadStr, msg, http.StatusBadRequest)
		}
		mo.Validate(true)
		_, err = tsh.tenantDB.UpdateMonitoredObject(mo)
		if err != nil {
			msg := fmt.Sprintf("Could not update MO %s %s. %s", moRequestTenant, mo.ID, err.Error())
			reportError(w, startTime, "400", migrateMOMetadatadStr, msg, http.StatusBadRequest)
		}
		for k, v := range mo.Meta {
			metas[k] = v
		}
	}
	err = tsh.tenantDB.UpdateMonitoredObjectMetadataViews(moRequestTenant, metas)
	if err != nil {
		msg := fmt.Sprintf("Could not update metadata views %s %v. %s", moRequestTenant, metas, err.Error())
		reportError(w, startTime, "400", migrateMOMetadatadStr, msg, http.StatusBadRequest)
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", migrateMOMetadatadStr)
	fmt.Fprintf(w, "Success")
}

// PopulateTestDataIntoDruid - Populate Test data produces data but druid does not contain any corresponding
// data. So this endpoint allows us to populate druid with some test data so we can make queries to it
// Supported query parameters are:
//			minutes: the number of desires minutes
//			tenant: the ID of the tenant to associate the monitored objects with
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TestDataServiceHandler) PopulateTestDataIntoDruid(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	queryParams := r.URL.Query()

	var (
		moRequestTenant string
		minutes         uint64 = 60
		err             error
	)

	// Defines the number of monitored objects that we want to create. Defaults to 1
	if queryParams["minutes"] != nil {
		minutes, err = strconv.ParseUint(queryParams["minutes"][0], 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Unacceptable value provided for monitored object count: %s", err.Error())
			reportError(w, startTime, "400", populateTestDataIntoDruidStr, msg, http.StatusBadRequest)
			return
		}
	}

	// Defines the tenant to place the monitored objects against. This is a required field
	if len(queryParams["tenant"]) == 0 {
		msg := fmt.Sprintf("Tenant ID must be provided")
		reportError(w, startTime, "400", populateTestDataIntoDruidStr, msg, http.StatusBadRequest)
		return
	} else {
		moRequestTenant = queryParams["tenant"][0]
	}

	err = tsh.PopulateDruidWithFauxData(moRequestTenant, minutes)
	if err != nil {
		msg := fmt.Sprintf("Could not populate with fake data. %s", err.Error())
		reportError(w, startTime, "400", populateTestDataIntoDruidStr, msg, http.StatusBadRequest)
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", populateTestDataIntoDruidStr)
	fmt.Fprintf(w, "Success")
}

// PopulateDruidWithFauxData - Populates druid with data so we can query it
func (tsh *TestDataServiceHandler) PopulateDruidWithFauxData(tenantID string, minutes uint64) error {

	kafkaProducer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        kafkaTopic,
		RequiredAcks: 0,
		Async:        true,
		Balancer:     &kafka.LeastBytes{},
	})
	defer func() {
		kafkaProducer.Close()
	}()

	listOfMonObjs, err := tsh.tenantDB.GetAllMonitoredObjectsIDs(tenantID)
	if err != nil {
		return fmt.Errorf("Could not get all monitored objects :%s", err.Error())
	}

	// for debugging!
	//istOfMonObjs = listOfMonObjs[0:2]

	ts := db.MakeTimestamp()
	ts = ts - ((60 * 1000) * int64(minutes)) - 24000

	logger.Log.Debugf("Sending FAKED OUT MO that has failure data")
	faux := "debug_mo_failure_0000"
	logger.Log.Debugf("Sending data populating for MO: %s", faux)
	_, err = generateAndSendKafkaMsg(kafkaProducer, ts, tenantID, faux, fauxDataFailtemplate)
	if err != nil {
		logger.Log.Errorf("Could not send to kafka %s", err.Error())
		return err
	}

	logger.Log.Debugf("Starting loop to send data over kafka to druid")
	for _, mo := range listOfMonObjs {

		ts = db.MakeTimestamp()
		ts = ts - ((60 * 1000) * int64(minutes)) - 24000
		logger.Log.Debugf("Starting data populating for MO: %s", mo)
		// Make sure to send one for the defined number of minutes
		for i := uint64(0); i < minutes; i++ {
			ts = ts + (60 * 1000) // Add a minute
			_, err = generateAndSendKafkaMsg(kafkaProducer, ts, tenantID, mo, fauxDatatemplate)
			if err != nil {
				logger.Log.Errorf("Could not send to kafka %s", err.Error())
				return err
			}
		}
	}

	return nil
}

// SignCSR - Sign CSR and return client cert
func (tsh *TestDataServiceHandler) SignCSR(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	logger.Log.Debugf("Received CSR request")

	caPublicKeyFile, err := ioutil.ReadFile("/run/secrets/tls_ca_crt")
	if err != nil {
		msg := fmt.Sprintf("Unable to find local ca.crt: %s ", err.Error())
		reportError(w, startTime, "500", signCSRStr, msg, http.StatusInternalServerError)
		return
	}
	pemBlock, _ := pem.Decode(caPublicKeyFile)
	if pemBlock == nil {
		msg := fmt.Sprintf("Could not decode ca.crt public key: %s ", err.Error())
		reportError(w, startTime, "500", signCSRStr, msg, http.StatusInternalServerError)
		return
	}
	caCRT, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		msg := fmt.Sprintf("Could not parse ca.crt pemblock: %s ", err.Error())
		reportError(w, startTime, "500", signCSRStr, msg, http.StatusInternalServerError)
		return
	}

	//      private key
	caPrivateKeyFile, err := ioutil.ReadFile("/run/secrets/tls_ca_key")

	if err != nil {
		msg := fmt.Sprintf("Unable to find ca.key private key: %s ", err.Error())
		reportError(w, startTime, "500", signCSRStr, msg, http.StatusInternalServerError)
		return
	}
	pemBlock, _ = pem.Decode(caPrivateKeyFile)
	if pemBlock == nil {
		msg := fmt.Sprintf("Could not decode ca.key private key: %s ", err.Error())
		reportError(w, startTime, "500", signCSRStr, msg, http.StatusInternalServerError)
		return
	}

	caPrivateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		msg := fmt.Sprintf("Could not parse ca.key: %s ", err.Error())
		reportError(w, startTime, "500", signCSRStr, msg, http.StatusInternalServerError)
		return
	}

	csrBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("Could not read CSR POST body: %s ", err.Error())
		reportError(w, startTime, "400", signCSRStr, msg, http.StatusBadRequest)
		return
	}

	pemBlock, _ = pem.Decode(csrBytes)
	if pemBlock == nil {
		msg := fmt.Sprintf("Could not decode CSR: %s ", err.Error())
		reportError(w, startTime, "400", signCSRStr, msg, http.StatusBadRequest)
		return
	}
	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		msg := fmt.Sprintf("Could not parse CSR: %s ", err.Error())
		reportError(w, startTime, "400", signCSRStr, msg, http.StatusBadRequest)
		return
	}
	if err = clientCSR.CheckSignature(); err != nil {
		msg := fmt.Sprintf("Invalid CSR signature: %s ", err.Error())
		reportError(w, startTime, "400", signCSRStr, msg, http.StatusBadRequest)
		return
	}

	// create client certificate template
	clientCRTTemplate := x509.Certificate{
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,

		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		PublicKey:          clientCSR.PublicKey,

		SerialNumber: big.NewInt(2),
		Issuer:       caCRT.Subject,
		Subject:      clientCSR.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(720 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// create client certificate from template and CA public key
	clientCRTRaw, err := x509.CreateCertificate(cr.Reader, &clientCRTTemplate, caCRT, clientCSR.PublicKey, caPrivateKey)

	if err != nil {
		msg := fmt.Sprintf("Could not create x509 client certificate: %s ", err.Error())
		reportError(w, startTime, "500", signCSRStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Debugf("Successfully generate client cert, sending to client.")

	w.Write(clientCRTRaw)
}

func dockerLogin() (string, error) {
	type GoogleToken struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	meta := "http://metadata.google.internal/computeMetadata/v1"
	svcAcc := meta + "/instance/service-accounts/default/token"

	req, _ := http.NewRequest("GET", svcAcc, nil)
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var token GoogleToken
	json.Unmarshal(body, &token)

	return token.AccessToken, nil
}

type ManifestObject struct {
	Mediatype string
	Size      int
	Digest    string
}

type ManifestV2 struct {
	SchemaVersion int
	MediaType     string
	Config        ManifestObject
	Layers        []ManifestObject
}

type ManifestV1 struct {
	RepoTags []string
	Config   string
	Layers   []string
}

func convertManifest(manifest *ManifestV2, repo string) ([]ManifestV1, error) {
	var layers []string

	for _, l := range manifest.Layers {
		name := strings.Split(l.Digest, ":")[1] + ".tar.gz"
		layers = append(layers, name)
	}

	return []ManifestV1{
		ManifestV1{
			Config:   manifest.Config.Digest,
			Layers:   layers,
			RepoTags: []string{repo},
		},
	}, nil
}

func (tsh *TestDataServiceHandler) writeConnectorConfigs(archiveDir string, tenantID string, zone string) error {
	cfg := gather.GetConfig()
	configs, err := tsh.tenantDB.GetAllTenantConnectorConfigs(tenantID, zone)
	if err != nil {
		return err
	}

	config := configs[0]
	envTemplate := `export FILE_DIR=%s
                        export VERSION=%s`
	env := fmt.Sprintf(envTemplate, config.URL, cfg.GetString(gather.CK_connector_dockerVersion.String()))
	err = ioutil.WriteFile(archiveDir+"/.env", []byte(env), os.ModePerm)
	if err != nil {
		return err
	}

	configTemplate, err := ioutil.ReadFile("/files/connector/adh-roadrunner.yml")
	if err != nil {
		return err
	}

	rrConfig := fmt.Sprintf(string(configTemplate), cfg.GetString("deploy.domain"), tenantID, zone)
	err = ioutil.WriteFile(archiveDir+"/adh-roadrunner.yml", []byte(rrConfig), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// DownloadRoadrunner - Download Roadrunner package for isntallation
func (tsh *TestDataServiceHandler) DownloadRoadrunner(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	cfg := gather.GetConfig()
	startTime := time.Now()
	logger.Log.Infof("Received DownloadRoadrunner request")

	archiveDir := "/tmp/roadrunnerArchive"
	os.MkdirAll(archiveDir, os.ModePerm)

	httpC := http.DefaultClient
	tr := &http.Transport{}
	httpC.Transport = tr

	accessToken, err := dockerLogin()

	if err != nil {
		msg := fmt.Sprintf("Unable to login to docker: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	imageName := cfg.GetString(gather.CK_connector_dockerImageName.String())
	baseURL := cfg.GetString(gather.CK_connector_dockerRegistry.String()) + imageName + "/"
	manifestURL := baseURL + "manifests/" + cfg.GetString(gather.CK_connector_dockerVersion.String())
	layerURL := baseURL + "blobs/"

	req, err := http.NewRequest("GET", manifestURL, nil)
	if err != nil {
		msg := fmt.Sprintf("Unable to create docker image manifest request: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// fetch manifest from docker registry
	manifestResp, err := httpC.Do(req)

	if err != nil {
		msg := fmt.Sprintf("Unable to fetch docker image manifest: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	manifest, err := ioutil.ReadAll(manifestResp.Body)
	if err != nil {
		msg := fmt.Sprintf("Unable to read docker image manifest: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}
	manifestObj := &ManifestV2{}

	err = json.Unmarshal(manifest, manifestObj)

	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshall docker image manifest: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	// convert from manivestV2 to manifest V1
	manifestV1, _ := convertManifest(manifestObj, cfg.GetString(gather.CK_connector_dockerRegistryPrefix.String())+
		imageName+":"+cfg.GetString(gather.CK_connector_dockerVersion.String()))

	manifestBytes, _ := json.Marshal(manifestV1)

	ioutil.WriteFile(archiveDir+"/manifest.json", manifestBytes, os.ModePerm)

	// Get the config object

	config := manifestObj.Config
	req, err = http.NewRequest("GET", layerURL+config.Digest, nil)
	if err != nil {
		msg := fmt.Sprintf("Unable to create docker image config request: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", config.Mediatype)

	configResp, err := httpC.Do(req)

	if err != nil {
		msg := fmt.Sprintf("Unable to fetch docker image config: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	configBytes, _ := ioutil.ReadAll(configResp.Body)
	ioutil.WriteFile(archiveDir+"/"+config.Digest, configBytes, os.ModePerm)

	// fetch the blobs that make up the docker image
	for _, l := range manifestObj.Layers {
		req, err := http.NewRequest("GET", layerURL+l.Digest, nil)
		if err != nil {
			msg := fmt.Sprintf("Unable to create docker layer request for url %s: %s ", layerURL+l.Digest, err.Error())
			reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Accept", l.Mediatype)

		layerResp, err := httpC.Do(req)

		if err != nil {
			msg := fmt.Sprintf("Unable to fetch docker layer request for url %s: %s ", layerURL+l.Digest, err.Error())
			reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
			return
		}

		name := strings.Split(l.Digest, ":")[1]
		ext := ".tar.gz"
		fullPath := archiveDir + "/" + name + ext

		f, err := os.Create(fullPath)
		if err != nil {
			msg := fmt.Sprintf("Unable to create file %s: %s ", fullPath, err.Error())
			reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
			return
		}

		fileBytes, _ := ioutil.ReadAll(layerResp.Body)
		f.Write(fileBytes)

		f.Close()
	}

	files, err := ioutil.ReadDir(archiveDir)

	if err != nil {
		msg := fmt.Sprintf("Unable to read directory %s: %s ", archiveDir, err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	var filenames []string
	for _, f := range files {
		filenames = append(filenames, archiveDir+"/"+f.Name())
	}

	// create docker image
	err = archiver.Tar.Make(archiveDir+"/roadrunner.docker", filenames)
	if err != nil {
		msg := fmt.Sprintf("Unable to save docker image %s: %s ", archiveDir+"/roadrunner.docker", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	archivePath := archiveDir + "/roadrunner.tar.gz"

	// write env.sh file
	err = tsh.writeConnectorConfigs(archiveDir, queryParams["tenant"][0], queryParams["zone"][0])
	if err != nil {
		msg := fmt.Sprintf("Unable to write env file: %s ", err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}
	// Make arhive for downloading
	err = archiver.Tar.Make(archivePath, []string{archiveDir + "/roadrunner.docker", "/files/connector/run.sh", archiveDir + "/.env", archiveDir + "/adh-roadrunner.yml"})
	if err != nil {
		msg := fmt.Sprintf("Unable to save roadrunner archive  %s: %s ", archivePath, err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	f, err := os.Open(archivePath)

	if err != nil {
		msg := fmt.Sprintf("Unable to open archive for downloading %s: %s ", archivePath, err.Error())
		reportError(w, startTime, "500", downloadRRStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Successfully generate Roadrunner package for downloading, sending to client.")

	w.Header().Add("Content-Disposition", "attachment; filename=DataHubConnector.tar.gz;")
	w.Header().Add("Content-Type", "multipart/form-data")

	io.Copy(w, f)
}

// generateAndSendKafkaMsg - Generates a Kafka message to send metric data to druid.
func generateAndSendKafkaMsg(kafkaProducer *kafka.Writer, ts int64, tenantID string, moName, faux string) (string, error) {
	nts := fmt.Sprintf("%d", ts)
	payload := strings.Replace(faux, tenantWord, tenantID, -1)
	payload = strings.Replace(payload, monObjIDWord, moName, -1)
	payload = strings.Replace(payload, monObjNameWord, moName, -1)
	payload = strings.Replace(payload, timestampWord, nts, -1)
	logger.Log.Debugf("Kafka sending: %s", payload)

	err := kafkaProducer.WriteMessages(context.Background(), kafka.Message{
		Topic: kafkaTopic,
		Value: []byte(payload),
	})
	if err != nil {
		return "", err
	}
	return "", nil
}
