package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
)

const (
	unableToReadRequestStr     = "Unable to read Test Data content"
	missingTenantNameStr       = "Missing tenantName field"
	missingDomainSetStr        = "Missing domainSet field"
	missingMonObjID            = "Missing Monitored Object id field"
	missingMonObjActuatorName  = "Missing Monitored Object actuatorName field"
	missingMonObjReflectorName = "Missing Monitored Object reflectorName field"
	missingMonObjObjectName    = "Missing Monitored Object objectName field"
	moIDStr                    = "id"
	moActuatorNameStr          = "actuatorName"
	moReflectorNameStr         = "reflectorName"
	moObjectNameStr            = "objectName"
)

// TestDataServiceHandler - handler for all APIs related to test data provisioning.
type TestDataServiceHandler struct {
	routes []server.Route
	// adminDB  db.AdminServiceDatastore
	// tenantDB db.TenantServiceDatastore
	pouchDB db.PouchDBPluginServiceDatastore
	grpcSH  *GRPCServiceHandler
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

	result.grpcSH = CreateCoordinator()

	return result
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

	// Get the contents of the request:
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read %s content: %s", "Test Data", err.Error()), http.StatusBadRequest)
		return
	}

	err = validatePopulateTestDataRequest(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a Tenant metadata using the provided name:
	tenantName := requestBody["tenantName"].(string)
	tenantDescriptor, err := tsh.grpcSH.CreateTenant(nil, (generateTenantDescriptor(tenantName)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
		return
	}

	// Now the TenantDB exisits....add a default user.
	_, err = tsh.grpcSH.CreateTenantUser(nil, generateTenantUser(tenantName, tenantDescriptor.GetXId()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
		return
	}

	// Add in the Domain Objects.
	tenantDomains := requestBody["domainSet"].([]interface{})
	createdDomainIDSet := []string{}
	for _, domainName := range tenantDomains {
		domain, err := tsh.grpcSH.CreateTenantDomain(nil, generateTenantDomain(domainName.(string), tenantDescriptor.GetXId(), tenantDescriptor.GetData().GetDefaultThresholdProfile()))
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
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

				_, err := tsh.grpcSH.CreateMonitoredObject(nil, generateMonitoredObject(obj[moIDStr].(string), tenantDescriptor.GetXId(), obj[moActuatorNameStr].(string), obj[moReflectorNameStr].(string), obj[moObjectNameStr].(string), createdDomainIDSet))
				if err != nil {
					http.Error(w, fmt.Sprintf("Unable to provision Tenant %s content: %s", tenantName, err.Error()), http.StatusInternalServerError)
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
	fmt.Fprintf(w, string(response))
}

// PurgeDB - deletes all records from a DB but does not destroy the DB.
func (tsh *TestDataServiceHandler) PurgeDB(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, 2)

	// Get a list of documents from the DB:
	docs, err := tsh.pouchDB.GetAllDBDocs(dbName, map[string]interface{}{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to purge DB %s: %s", dbName, err.Error()), http.StatusInternalServerError)
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
			return
		}
	}

	// Purge complete, send back success details:
	fmt.Fprintf(w, "Successfully purged all records from DB "+dbName)
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

func generateTenantDescriptor(name string) *pb.TenantDescriptorRequest {
	result := pb.TenantDescriptorRequest{}

	tenantStr := string(db.TenantDescriptorType)
	result.Data = &pb.TenantDescriptor{}
	result.Data.Name = name
	result.Data.UrlSubdomain = strings.ToLower(name) + ".npav.accedian.net"
	result.Data.State = 2
	result.Data.Datatype = tenantStr
	result.Data.DefaultThresholdProfile = ""
	result.Data.CreatedTimestamp = time.Now().Unix()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	result.XId = db.GenerateID(result.Data, tenantStr)

	return &result
}

func generateTenantUser(name string, tenantID string) *pb.TenantUserRequest {
	result := pb.TenantUserRequest{}

	tenantUserStr := string(db.TenantUserType)

	result.Data = &pb.TenantUser{}
	result.Data.Datatype = tenantUserStr
	result.Data.Username = "admin@" + strings.ToLower(name) + ".com"
	result.Data.TenantId = tenantID
	result.Data.UserVerified = true
	result.Data.State = 2
	result.Data.SendOnboardingEmail = false
	result.Data.Password = "admin"
	result.Data.OnboardingToken = "anonboardingtokenforthedefaulttenantuser"
	result.Data.CreatedTimestamp = time.Now().Unix()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	result.XId = db.GenerateID(result.Data, tenantUserStr)

	return &result
}

func generateTenantDomain(name string, tenantID string, defaultThreshPrf string) *pb.TenantDomainRequest {
	result := pb.TenantDomainRequest{}

	tenantDomainStr := string(db.TenantDomainType)
	result.Data = &pb.TenantDomain{}
	result.Data.Datatype = tenantDomainStr
	result.Data.Name = name
	result.Data.TenantId = tenantID
	result.Data.Color = "#0000FF"

	// Only use the hash portion of the id for the reference since the thresholdProfile_2_ is all
	// relational pouch tagging.
	thrPrfIDParts := strings.Split(defaultThreshPrf, "_")
	result.Data.ThresholdProfileSet = append(result.Data.GetThresholdProfileSet(), thrPrfIDParts[len(thrPrfIDParts)-1])
	result.Data.CreatedTimestamp = time.Now().Unix()
	result.Data.LastModifiedTimestamp = result.GetData().GetCreatedTimestamp()

	result.XId = db.GenerateID(result.Data, tenantDomainStr)

	return &result
}

func generateMonitoredObject(id string, tenantID string, actuatorName string, reflectorName string, objectName string, domainIDSet []string) *pb.MonitoredObjectRequest {
	result := pb.MonitoredObjectRequest{}

	tenantMonObjStr := string(db.TenantMonitoredObjectType)
	result.Data = &pb.MonitoredObject{}
	result.Data.Datatype = tenantMonObjStr
	result.Data.Id = id
	result.Data.TenantId = tenantID
	result.Data.ActuatorName = actuatorName
	result.Data.ActuatorType = pb.MonitoredObject_ACCEDIAN_VNID
	result.Data.ReflectorName = reflectorName
	result.Data.ReflectorType = pb.MonitoredObject_ACCEDIAN_VNID
	result.Data.ObjectName = objectName
	result.Data.ObjectType = pb.MonitoredObject_TWAMP

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
		case 0:
			// Do nothing
		case 1:
			// Just take a randome one.
			result.GetData().DomainSet = append(result.GetData().GetDomainSet(), domainIDSet[rand.Intn(lengthOfDomainSet)])
		case lengthOfDomainSet:
			// Take them all
			result.GetData().DomainSet = append(result.GetData().GetDomainSet(), domainIDSet...)
		default:
			// Take a subset.
			indextToStopAt := rand.Intn(lengthOfDomainSet)
			result.GetData().DomainSet = append(result.GetData().GetDomainSet(), domainIDSet[:indextToStopAt]...)
		}
	}

	result.XId = db.GenerateID(result.Data, tenantMonObjStr)

	return &result
}
