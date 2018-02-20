package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/satori/go.uuid"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/datastore/inMemory"
	"github.com/accedian/adh-gather/gather"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/server"
	wr "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/gorilla/mux"
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

	populateTestDataStr  = "populate_test_data"
	purgeDBStr           = "purge_db"
	generateSLAReportStr = "gen_sla_report"
	getDocsByTypeStr     = "get_docs_by_type"
)

// TestDataServiceHandler - handler for all APIs related to test data provisioning.
type TestDataServiceHandler struct {
	routes []server.Route
	// adminDB  db.AdminServiceDatastore
	// tenantDB db.TenantServiceDatastore
	pouchDB db.PouchDBPluginServiceDatastore
	grpcSH  *GRPCServiceHandler
	testDB  db.TestDataServiceDatastore
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
			HandlerFunc: result.PopulateTestData,
		},

		server.Route{
			Name:        "PurgeDB",
			Method:      "DELETE",
			Pattern:     "/test-data/{dbname}",
			HandlerFunc: result.PurgeDB,
		},

		server.Route{
			Name:        "GenerateHistoricalDomainSLAReports",
			Method:      "POST",
			Pattern:     "/test-data/domain-sla-reports",
			HandlerFunc: result.GenerateHistoricalDomainSLAReports,
		},

		server.Route{
			Name:        "GetAllDocsByType",
			Method:      "GET",
			Pattern:     "/test-data/{dbname}/{datatype}",
			HandlerFunc: result.GetAllDocsByType,
		},
	}

	// Wire up the datastore impls
	// admindb, err := getAdminServiceDatastore()
	// if err != nil {
	// 	logger.Log.Fatalf("Unable to instantiate TestDataServiceHandler: %s", err.Error())
	// }
	// result.adminDB = admindb

	// tenantdb, err := getTenantServiceDatastore()
	// if err != nil {
	// 	logger.Log.Fatalf("Unable to instantiate TestDataServiceHandler: %s", err.Error())
	// }
	// result.tenantDB = tenantdb

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
		http.Error(w, fmt.Sprintf("Unable to read %s content: %s", "Test Data", err.Error()), http.StatusBadRequest)
		mon.TrackAPITimeMetricInSeconds(startTime, "400", populateTestDataStr)
		return
	}

	err = validatePopulateTestDataRequest(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		mon.TrackAPITimeMetricInSeconds(startTime, "400", populateTestDataStr)
		return
	}

	// Create a Tenant metadata using the provided name:
	tenantName := requestBody["tenantName"].(string)
	desc, err := generateTenantDescriptor(tenantName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
		return
	}
	tenantDescriptor, err := tsh.grpcSH.CreateTenant(nil, desc)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
		return
	}

	// Now the TenantDB exisits....add a default user.
	dataIDForTenant := tenantDescriptor.GetXId()
	user, err := generateTenantUser(tenantName, dataIDForTenant)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
		return
	}
	_, err = tsh.grpcSH.CreateTenantUser(nil, user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
		return
	}

	// Fetch the tenant Meta so that teh default threshold profile id is available:
	tenantMeta, err := tsh.grpcSH.GetTenantMeta(nil, &wr.StringValue{Value: tenantDescriptor.GetXId()})
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
		return
	}

	// Add in the Domain Objects.
	tenantDomains := requestBody["domainSet"].([]interface{})
	createdDomainIDSet := []string{}
	for _, domainName := range tenantDomains {
		dom, err := generateTenantDomain(domainName.(string), dataIDForTenant, tenantMeta.GetData().GetDefaultThresholdProfile())
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
			mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
			return
		}
		domain, err := tsh.grpcSH.CreateTenantDomain(nil, dom)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
			mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
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
					http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
					mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
					return
				}
				_, err = tsh.grpcSH.CreateMonitoredObject(nil, mo)
				if err != nil {
					http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
					mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
					return
				}
			}
		}
	}

	// All operations successful, just return the Tenant Descriptor as success:
	response, err := json.Marshal(tenantDescriptor)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating response: %s", err.Error()), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Unable to purge DB %s of documents of type %s: %s", dbName, datatypeFilter, err.Error()), http.StatusInternalServerError)
			mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
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
				http.Error(w, fmt.Sprintf("Unable to purge DB %s: %s", dbName, err.Error()), http.StatusInternalServerError)
				mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
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
		http.Error(w, fmt.Sprintf("Unable to purge DB %s: %s", dbName, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
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
			http.Error(w, fmt.Sprintf("Unable to purge DB %s: %s", dbName, err.Error()), http.StatusInternalServerError)
			mon.TrackAPITimeMetricInSeconds(startTime, "500", populateTestDataStr)
			return
		}
	}

	// Purge complete, send back success details:
	mon.TrackAPITimeMetricInSeconds(startTime, "200", populateTestDataStr)
	fmt.Fprintf(w, "Successfully purged all records from DB "+dbName)
}

// GenerateHistoricalDomainSLAReports - generates weekly SLA reports for a Tenant based on the Domains provisioned for that Tenant.
func (tsh *TestDataServiceHandler) GenerateHistoricalDomainSLAReports(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	// Get the contents of the request:
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read %s content: %s", "Test Data", err.Error()), http.StatusBadRequest)
		mon.TrackAPITimeMetricInSeconds(startTime, "400", generateSLAReportStr)
		return
	}

	err = validateGenerateDomainSLAReportRequest(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		mon.TrackAPITimeMetricInSeconds(startTime, "400", generateSLAReportStr)
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
		http.Error(w, fmt.Sprintf("Unable to generate %s content for Tenant %s: No Domain data provisioned", db.DomainSlaReportStr, tenantID), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", generateSLAReportStr)
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
				http.Error(w, fmt.Sprintf("Unable to generate %s content for Tenant %s: %s", db.DomainSlaReportStr, tenantID, err.Error()), http.StatusInternalServerError)
				mon.TrackAPITimeMetricInSeconds(startTime, "500", generateSLAReportStr)
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
	tenantIDForUpdate := db.PrependToDataID(tenantID, string(db.TenantDescriptorType))
	_, err = tsh.pouchDB.BulkDBUpdate(tenantIDForUpdate, insertBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to insert %ss for Tenant %s: %s", db.DomainSlaReportStr, tenantID, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", generateSLAReportStr)
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
		http.Error(w, fmt.Sprintf("Unable to retrieve documents without a DB name and Document type"), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", getDocsByTypeStr)
		return
	}

	allDocsByType, err := tsh.testDB.GetAllDocsByDatatype(dbName, docType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve documents of type %s from DB %s: %s", docType, dbName, err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", getDocsByTypeStr)
		return
	}

	// All operations successful, just return the Tenant Descriptor as success:
	response, err := json.Marshal(allDocsByType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating response: %s", err.Error()), http.StatusInternalServerError)
		mon.TrackAPITimeMetricInSeconds(startTime, "500", getDocsByTypeStr)
		return
	}
	mon.TrackAPITimeMetricInSeconds(startTime, "200", getDocsByTypeStr)
	fmt.Fprintf(w, string(response))
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

	tenantStr := string(db.TenantDescriptorType)
	result.Data = &pb.TenantDescriptorData{}
	result.Data.Name = name
	result.Data.UrlSubdomain = strings.ToLower(name) + ".npav.accedian.net"
	result.Data.State = 2
	result.Data.Datatype = tenantStr
	result.Data.CreatedTimestamp = db.MakeTimestamp()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	result.XId = db.GenerateID(result.Data, tenantStr)

	return &result, nil
}

func generateTenantUser(name string, tenantID string) (*pb.TenantUser, error) {
	result := pb.TenantUser{}

	tenantUserStr := string(db.TenantUserType)

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

	tenantDomainStr := string(db.TenantDomainType)
	result.Data = &pb.TenantDomainData{}
	result.Data.Datatype = tenantDomainStr
	result.Data.Name = name
	result.Data.TenantId = tenantID
	result.Data.Color = "#4EC5C1"

	// Only use the hash portion of the id for the reference since the thresholdProfile_2_ is all
	// relational pouch tagging.
	thrPrfIDParts := strings.Split(defaultThreshPrf, "_")
	result.Data.ThresholdProfileSet = append(result.Data.GetThresholdProfileSet(), thrPrfIDParts[len(thrPrfIDParts)-1])
	result.Data.CreatedTimestamp = db.MakeTimestamp()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	result.XId = db.GenerateID(result.Data, tenantDomainStr)

	return &result, nil
}

func generateMonitoredObject(id string, tenantID string, actuatorName string, reflectorName string, objectName string, domainIDSet []string) (*pb.MonitoredObject, error) {
	result := pb.MonitoredObject{}

	tenantMonObjStr := string(db.TenantMonitoredObjectType)
	result.Data = &pb.MonitoredObjectData{}
	result.Data.Datatype = tenantMonObjStr

	result.Data.Id = id
	result.Data.TenantId = tenantID
	result.Data.ActuatorName = actuatorName
	result.Data.ActuatorType = string(AccedianVNID)
	result.Data.ReflectorName = reflectorName
	result.Data.ReflectorType = string(AccedianVNID)
	result.Data.ObjectName = objectName
	result.Data.ObjectType = string(TwampPE)

	// To provision the DomainSet, need to obtain a subset of the passed in domain set.

	lengthOfDomainSet := len(domainIDSet)
	if lengthOfDomainSet == 1 {
		// Only 1 domain, coinflip to see if it is selected.
		rand.Seed(time.Now().Unix())
		if rand.Intn(10)%2 == 0 {
			result.GetData().DomainSet = append(result.GetData().GetDomainSet(), domainIDSet...)
		}
	} else if lengthOfDomainSet != 0 {
		numDomainsToSelect := rand.Intn(lengthOfDomainSet) + 1
		switch numDomainsToSelect {
		case 1:
			// Just take a random one.
			result.GetData().DomainSet = append(result.GetData().GetDomainSet(), domainIDSet[rand.Intn(lengthOfDomainSet)])
		case lengthOfDomainSet:
			// Take them all
			result.GetData().DomainSet = append(result.GetData().GetDomainSet(), domainIDSet...)
		default:
			// Take a subset.
			indextToStopAt := rand.Intn(lengthOfDomainSet) + 1
			result.GetData().DomainSet = append(result.GetData().GetDomainSet(), domainIDSet[:indextToStopAt]...)
		}
	}

	result.XId = db.GenerateID(result.Data, tenantMonObjStr)

	return &result, nil
}
