package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/messaging"

	"github.com/accedian/adh-gather/gather"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
)

const (
	recommendationRequestPath = "/recommendation"
	makeRecommendationAPIStr  = "colt_make_rcm"
)

type ColtMEFHandler struct {
	routes     []server.Route
	httpClient *http.Client

	requestReader *messaging.KafkaConsumer
	pendingReader *messaging.KafkaConsumer

	pendingWriter *messaging.KafkaProducer
	resultWriter  *messaging.KafkaProducer

	server       string
	appID        string
	sharedSecret string
}

func CreateColtMEFHandler() *ColtMEFHandler {
	requestTopic := "colt-mef-requests"
	pendingTopic := "colt-mef-pending"
	resultTopic := "colt-mef-results"
	result := new(ColtMEFHandler)

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	result.httpClient = &http.Client{Transport: tr}

	cfg := gather.GetConfig()
	result.server = cfg.GetString(gather.CK_args_coltmef_server.String())
	result.appID = cfg.GetString(gather.CK_args_coltmef_appid.String())
	result.sharedSecret = cfg.GetString(gather.CK_args_coltmef_secret.String())

	result.requestReader = messaging.CreateKafkaReader(requestTopic, "0")
	result.pendingReader = messaging.CreateKafkaReader(pendingTopic, "0")

	// Start the message readers
	go func() {
		for {
			result.requestReader.ReadMessage(result.handleRecommendationRequest)
		}
	}()
	go func() {
		for {
			result.pendingReader.ReadMessage(result.handleRecommendationStatusCheck)
		}
	}()

	result.pendingWriter = messaging.CreateKafkaWriter(pendingTopic)
	result.resultWriter = messaging.CreateKafkaWriter(resultTopic)

	result.routes = []server.Route{

		server.Route{
			Name:        "MakeRecommendation",
			Method:      "POST",
			Pattern:     "/colt-mef/recommendation",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, result.MakeRecommendation),
		},
	}

	return result
}

// RegisterAPIHandlers - will bind any REST API routes defined in this service
// to the passed in request multiplexor.
func (cmh *ColtMEFHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range cmh.routes {
		logger.Log.Debugf("Registering endpoint: %v", route)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

// TODO: remove this handler when we officially move to just KAFKA
// MakeRecommendation - REST Handler for Recommend a service change API.
func (cmh *ColtMEFHandler) MakeRecommendation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", makeRecommendationAPIStr, msg, http.StatusBadRequest)
		return
	}

	// Deserialize the request
	requestObj := &ColtRecommendation{}
	err = json.Unmarshal(requestBytes, requestObj)
	if err != nil {
		msg := fmt.Sprintf("Unable to read service change data: %s", err.Error())
		reportError(w, startTime, "400", makeRecommendationAPIStr, msg, http.StatusBadRequest)
		return
	}

	responseObj, code, err := cmh.doMakeRecommendation(requestObj)
	if err != nil {
		reportError(w, startTime, string(code), makeRecommendationAPIStr, err.Error(), code)
		return
	}

	logger.Log.Infof("Completed service change: %s", db.HistogramStr, string(requestBytes))
	trackAPIMetrics(startTime, "200", makeRecommendationAPIStr)
	fmt.Fprintf(w, responseObj.RecommendationID)
}

// doMakeRecommendation - Handles the logic to make a call to the Colt POST /api/performance/recommendation API for making a service change recommendation
func (cmh *ColtMEFHandler) doMakeRecommendation(requestObj *ColtRecommendation) (*ColtRecommendationResponse, int, error) {

	// Re-serialize the bytes to ensure we do not have any "extra stuff" in the request
	requestObjBytes, err := json.Marshal(requestObj)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to prepare service change data: %s", err.Error())
	}

	// Setup the request to Colt
	req, err := http.NewRequest("POST", cmh.server, bytes.NewBuffer(requestObjBytes))
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to prepare service change request: %s", err.Error())
	}

	// Fill in necessary request headers
	req.Header.Set("x-colt-app-id", cmh.appID)
	req.Header.Set("x-colt-app-sig", getAuthHeader(requestObjBytes, cmh.sharedSecret, recommendationRequestPath))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Issue request to COlt
	resp, err := cmh.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to issue service change: %s", err.Error())
	}

	defer resp.Body.Close()

	// Read the request
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to read service change response: %s", err.Error())
	}

	logger.Log.Debugf("MAKE RECOMMENDATION RESPONSE: %s", string(respBytes))

	if resp.StatusCode != http.StatusCreated {
		// Request was not successful, format the error response
		responseObj := &ColtError{}
		err = json.Unmarshal(respBytes, responseObj)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal service change response: %s", err.Error())
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("Service change failed: %d - %s", responseObj.Code, responseObj.Message)
	}

	// Request was successful, format the response object
	responseObj := &ColtRecommendationResponse{}
	err = json.Unmarshal(respBytes, responseObj)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal service change response: %s", err.Error())
	}

	return responseObj, http.StatusOK, nil
}

// doCheckRecommendationStatus - Handles the logic to make a call to the Colt GET /api/performance/recommendation/{recommendationID} API for checking the
// status of a service change recommendation
func (cmh *ColtMEFHandler) doCheckRecommendationStatus(recommendationID string) (*ColtRecommendationState, int, error) {
	// Setup the request to Colt
	req, err := http.NewRequest("GET", cmh.server, nil)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to prepare service change status request: %s", err.Error())
	}

	req.Header.Set("x-colt-app-id", cmh.appID)
	req.Header.Set("x-colt-app-sig", getAuthHeader([]byte{}, cmh.sharedSecret, recommendationRequestPath+"/"+recommendationID))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cmh.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to issue service change status request: %s", err.Error())
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to read service change status response: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {

		responseObj := &ColtError{}
		err = json.Unmarshal(respBytes, responseObj)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal service change status response: %s", err.Error())
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("Service change status check failed: %d - %s", responseObj.Code, responseObj.Message)
	}

	responseObj := &ColtRecommendationState{}
	err = json.Unmarshal(respBytes, responseObj)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal service change status response: %s", err.Error())
	}

	return responseObj, http.StatusOK, nil
}

// handleRecommendationRequest - method to be used to handle messages pulled off of the service change request topic
func (cmh *ColtMEFHandler) handleRecommendationRequest(requestBytes []byte) bool {

	requestObj := &ServiceChangeRequest{}
	err := json.Unmarshal(requestBytes, requestObj)
	if err != nil {
		logger.Log.Errorf("Unable to read service change data: %s", err.Error())
		return true
	}

	// Issue the Recommendation API to Colt
	responseObj, code, err := cmh.doMakeRecommendation(requestObj.ServiceChange)
	if err != nil {
		msg := err.Error()
		logger.Log.Error(msg)
		cmh.writeResult(requestObj.RequestID, "", "FAILED", msg)
		return true
	}

	if code != http.StatusOK {
		msg := fmt.Sprintf("Unable to complete Service Change Request. Response code was %d", code)
		logger.Log.Error(msg)
		cmh.writeResult(requestObj.RequestID, responseObj.RecommendationID, "FAILED", msg)
		return true
	}

	// Write a record to the Pending Topic to continue the process
	logger.Log.Infof("Service Change Recommendation %s added to pending queue", responseObj.RecommendationID)
	cmh.writePending(requestObj.RequestID, responseObj.RecommendationID)

	// Recommendation was completed successfully
	return true
}

// handleRecommendationStatusCheck - method to be used to handle messages pulled off of the service change pending topic
func (cmh *ColtMEFHandler) handleRecommendationStatusCheck(recommendationStatusRequest []byte) bool {

	// Deserialize the request
	requestObj := &ServiceChangeCheckStatusRequest{}
	err := json.Unmarshal(recommendationStatusRequest, requestObj)
	if err != nil {
		logger.Log.Errorf("Unable to read service change status data: %s", err.Error())
		return true
	}

	// Poll the status API until we get a successful response
	maxPollCount := 5
	pollCount := 0
	var pollResp *ColtRecommendationState
	for {
		time.Sleep(10 * time.Second)

		var code int
		var err error
		pollResp, code, err = cmh.doCheckRecommendationStatus(requestObj.RecommendationID)
		if err != nil || code != http.StatusOK {
			msg := fmt.Sprintf("Unable to check status of Recommendation %s: %d - %s", requestObj.RecommendationID, code, err.Error())
			logger.Log.Errorf(msg)
			cmh.writeResult(requestObj.RequestID, requestObj.RecommendationID, "FAILED", msg)
			return true
		}

		if pollResp.State == "PENDING" {
			continue
		}

		if pollCount >= maxPollCount {
			// Too many attempts, just fail this request
			msg := fmt.Sprintf("Unable to check status of Recommendation %s: Request timed out waiting for Recommendation completion", requestObj.RecommendationID)
			logger.Log.Errorf(msg)
			cmh.writeResult(requestObj.RequestID, requestObj.RecommendationID, "FAILED", msg)
			return true
		}

		break
	}

	logger.Log.Infof("Service Change Recommendation %s completed with result %s", requestObj.RecommendationID, pollResp.State)
	cmh.writeResult(requestObj.RequestID, requestObj.RecommendationID, pollResp.State, "")
	return true
}

// writeResult - helper to write a result to the service change result topic
func (cmh *ColtMEFHandler) writeResult(reqID string, recID string, state string, err string) error {
	result := ServiceChangeResult{
		RequestID:        reqID,
		RecommendationID: recID,
		Status:           state,
		ErrorMessage:     err,
	}

	msgBytes, msgErr := json.Marshal(result)
	if msgErr != nil {
		return fmt.Errorf("Error marshalling result for recommendation %s: %s", reqID, msgErr.Error())
	}

	return cmh.resultWriter.WriteMessage("Result:"+reqID, msgBytes)
}

// writePending - helper to write a result to the service change pending topic
func (cmh *ColtMEFHandler) writePending(reqID string, recID string) error {
	result := ServiceChangeCheckStatusRequest{
		RequestID:        reqID,
		RecommendationID: recID,
	}

	msgBytes, msgErr := json.Marshal(result)
	if msgErr != nil {
		return fmt.Errorf("Error marshalling pending object for recommendation %s: %s", reqID, msgErr.Error())
	}

	return cmh.pendingWriter.WriteMessage("Pending:"+reqID, msgBytes)
}

// getAuthHeader - helper to build the necessary auth token for REST calls to Colt APIs
func getAuthHeader(recommendation []byte, key string, path string) string {
	hashedPayload := base64HMACSHA256(recommendation, key)

	timeNow := time.Now().UTC()
	dateY, dateM, DateD := timeNow.Date()
	hour := timeNow.Hour()

	timestamp := fmt.Sprintf("%04d%02d%02d%02d", dateY, dateM, DateD, hour)
	requestData := strings.Join([]string{timestamp, path, hashedPayload}, "")

	return base64HMACSHA256([]byte(requestData), key)
}

// base64HMACSHA256 - helper that builds a base64 encoded string out of a sha256 HMAC encoded value
func base64HMACSHA256(payload []byte, key string) string {
	hashObj := hmac.New(sha256.New, []byte(key))
	hashObj.Write(payload)
	return base64.StdEncoding.EncodeToString(hashObj.Sum(nil))
}
