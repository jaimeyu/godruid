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
	"github.com/accedian/adh-gather/models"

	"github.com/accedian/adh-gather/gather"

	"github.com/accedian/adh-gather/logger"
)

const (
	logPrefix                 = "COLT-MEF: "
	recommendationRequestPath = "/recommendation"
	errorPrefix               = "Recommendation"

	slackURL = "https://hooks.slack.com/services/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"
)

type ColtMEFHandler struct {
	httpClient *http.Client

	requestReader *messaging.KafkaConsumer
	// pendingReader *messaging.KafkaConsumer

	// pendingWriter *messaging.KafkaProducer
	resultWriter *messaging.KafkaProducer

	server           string
	appID            string
	sharedSecret     string
	statusRetryCount int

	pollCheckpoint1 float64
	pollCheckpoint2 float64
	pollCheckpoint3 float64
}

func CreateColtMEFHandler() *ColtMEFHandler {
	requestTopic := "colt-mef-requests"
	// pendingTopic := "colt-mef-pending"
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
	result.statusRetryCount = cfg.GetInt(gather.CK_args_coltmef_statusretrycount.String())

	result.pollCheckpoint1 = cfg.GetFloat64(gather.CK_args_coltmef_checkpoint1.String())
	result.pollCheckpoint2 = cfg.GetFloat64(gather.CK_args_coltmef_checkpoint2.String())
	result.pollCheckpoint3 = cfg.GetFloat64(gather.CK_args_coltmef_checkpoint3.String())

	logger.Log.Infof("%sStarting event handler at %s with app ID %s", logPrefix, result.server, result.appID)

	result.requestReader = messaging.CreateKafkaReader(requestTopic, "0")
	// result.pendingReader = messaging.CreateKafkaReader(pendingTopic, "0")

	// Start the message readers
	go func() {
		for {
			result.requestReader.ReadMessage(result.handleRecommendationRequest)
		}
	}()
	// go func() {
	// 	for {
	// 		result.pendingReader.ReadMessage(result.handleRecommendationStatusCheck)
	// 	}
	// }()

	// result.pendingWriter = messaging.CreateKafkaWriter(pendingTopic)
	result.resultWriter = messaging.CreateKafkaWriter(resultTopic)

	return result
}

// doMakeRecommendation - Handles the logic to make a call to the Colt POST /api/performance/recommendation API for making a service change recommendation
func (cmh *ColtMEFHandler) doMakeRecommendation(requestID string, requestObj *ColtRecommendation) (*ColtRecommendationResponse, int, error) {

	// Re-serialize the bytes to ensure we do not have any "extra stuff" in the request
	requestObjBytes, err := json.Marshal(requestObj)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to prepare service change data for service change request %s: %s", requestID, err.Error())
	}

	// Setup the request to Colt
	req, err := http.NewRequest("POST", cmh.server, bytes.NewBuffer(requestObjBytes))
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to build outgoing service change request for service change request %s: %s", requestID, err.Error())
	}

	// Fill in necessary request headers
	authHeader := getAuthHeader(requestObjBytes, cmh.sharedSecret, recommendationRequestPath)
	req.Header.Set("x-colt-app-id", cmh.appID)
	req.Header.Set("x-colt-app-sig", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("%s Submitting recommendation %s to server %s for app-id %s using auth token %s", logPrefix, string(requestObjBytes), cmh.server, cmh.appID, authHeader)
	}

	// Issue request to COlt
	resp, err := cmh.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to issue service change for service change request %s: %s", requestID, err.Error())
	}

	defer resp.Body.Close()

	// Read the request
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to read service change response for service change request %s: %s", requestID, err.Error())
	}

	logger.Log.Debugf("%sMAKE RECOMMENDATION RESPONSE [SERVICE CHANGE REQ: %s]: %s", logPrefix, requestID, string(respBytes))

	if resp.StatusCode != http.StatusOK {
		// Request was not successful, format the error response
		responseObj := &ColtError{}
		err = json.Unmarshal(respBytes, responseObj)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal recommendation response for service change request %s: %s", requestID, err.Error())
		}

		return nil, resp.StatusCode, fmt.Errorf("Service change failed for service change request %s: %d - %s", requestID, responseObj.Code, responseObj.Message)
	}

	// Request was successful, format the response object
	responseObj := &ColtRecommendationResponse{}
	err = json.Unmarshal(respBytes, responseObj)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal recommendation response for service change request %s: %s", requestID, err.Error())
	}

	return responseObj, http.StatusOK, nil
}

// doCheckRecommendationStatus - Handles the logic to make a call to the Colt GET /api/performance/recommendation/{recommendationID} API for checking the
// status of a service change recommendation
func (cmh *ColtMEFHandler) doCheckRecommendationStatus(requestID string, recommendationID string) (*ColtRecommendationState, int, error) {
	// Setup the request to Colt
	req, err := http.NewRequest("GET", cmh.server+"/"+recommendationID, nil)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to prepare service change status request for service change request %s and recommendation %s: %s", requestID, recommendationID, err.Error())
	}

	req.Header.Set("x-colt-app-id", cmh.appID)
	req.Header.Set("x-colt-app-sig", getAuthHeader([]byte{}, cmh.sharedSecret, recommendationRequestPath+"/"+recommendationID))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cmh.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to issue service change status request for service change request %s and recommendation %s: %s", requestID, recommendationID, err.Error())
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to read service change status response for service change request %s and recommendation %s: %s", requestID, recommendationID, err.Error())
	}

	logger.Log.Debugf("%sCHECK RECOMMENDATION STATE RESPONSE [SERVICE CHANGE REQ: %s, RECOMMENDATION REQ: %s]: %s", logPrefix, requestID, recommendationID, string(respBytes))

	if resp.StatusCode != http.StatusOK {

		responseObj := &ColtError{}
		err = json.Unmarshal(respBytes, responseObj)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal service change status response for service change request %s and recommendation %s: %s", requestID, recommendationID, err.Error())
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("Service change status check failed for service change request %s and recommendation %s: %d - %s", requestID, recommendationID, responseObj.Code, responseObj.Message)
	}

	responseObj := &ColtRecommendationState{}
	err = json.Unmarshal(respBytes, responseObj)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unable to unmarshal service change status response for service change request %s and recommendation %s: %s", requestID, recommendationID, err.Error())
	}

	return responseObj, http.StatusOK, nil
}

// handleRecommendationRequest - method to be used to handle messages pulled off of the service change request topic
func (cmh *ColtMEFHandler) handleRecommendationRequest(requestBytes []byte) bool {

	requestObj := &ServiceChangeRequest{}
	err := json.Unmarshal(requestBytes, requestObj)
	if err != nil {
		logger.Log.Errorf("%sUnable to read service change data: %s", logPrefix, err.Error())
		return true
	}

	logger.Log.Infof("%sReceived service change request %s: %s", logPrefix, requestObj.RequestID, models.AsJSONString(requestObj))

	// Issue the Recommendation API to Colt
	responseObj, code, err := cmh.doMakeRecommendation(requestObj.RequestID, requestObj.ServiceChange)
	if err != nil {
		// Handle a duplicate service change request for the same service:
		if code == http.StatusConflict {
			recommendationID := getRecommendationIDFromConflictMessage(err.Error())

			// Existing request for the serice, poll for the result
			logger.Log.Warningf("%sService Change recommendation %s is already in progress for service %s. Initiating status check", logPrefix, recommendationID, requestObj.ServiceChange.ServiceID)
			return cmh.pollRecommendationStatus(requestObj.RequestID, recommendationID)
		}

		msg := err.Error()
		logger.Log.Errorf("%s%s", logPrefix, msg)
		cmh.writeResult(requestObj.RequestID, "", "FAILED", msg)
		return true
	}

	if code != http.StatusOK {
		// Handle any other error a permaent failure
		msg := fmt.Sprintf("Unable to complete service change request %s. Response code was %d", requestObj.RequestID, code)
		logger.Log.Errorf("%s%s", logPrefix, msg)
		cmh.writeResult(requestObj.RequestID, responseObj.RecommendationID, "FAILED", msg)
		return true
	}

	logger.Log.Infof("%sRecommendation %s: successfully submitted", logPrefix, responseObj.RecommendationID)

	// The result will depend on the status of the polling result for the service change request
	return cmh.pollRecommendationStatus(requestObj.RequestID, responseObj.RecommendationID)
}

func getRecommendationIDFromConflictMessage(errorMessage string) string {

	errorParts := strings.Split(errorMessage, " ")
	for i, val := range errorParts {
		if val == errorPrefix {
			return errorParts[i+1]
		}
	}

	return ""
}

func (cmh *ColtMEFHandler) pollRecommendationStatus(requestID string, recommendationID string) bool {
	pendingBytes, err := createPendingPayload(requestID, recommendationID)
	if err != nil {
		msg := fmt.Sprintf("Unable to add service change recommendation %s for service change request %s to pending queue: %s", recommendationID, requestID, err.Error())
		logger.Log.Errorf("%s%s", logPrefix, msg)
		cmh.writeResult(requestID, recommendationID, "FAILED", msg)
		return true
	}

	// Poll status of update until complete:
	return cmh.handleRecommendationStatusCheck(pendingBytes)
}

// handleRecommendationStatusCheck - method to be used to handle messages pulled off of the service change pending topic
func (cmh *ColtMEFHandler) handleRecommendationStatusCheck(recommendationStatusRequest []byte) bool {

	// Deserialize the request
	requestObj := &ServiceChangeCheckStatusRequest{}
	err := json.Unmarshal(recommendationStatusRequest, requestObj)
	if err != nil {
		logger.Log.Errorf("%sUnable to read service change status data: %s", logPrefix, err.Error())
		return true
	}

	logger.Log.Debugf("%sPulled status request for recommendation %s for service change request %s", logPrefix, requestObj.RecommendationID, requestObj.RequestID)

	// Poll the status API until we get a successful response
	pollStart := time.Now()
	messageCount := 0
	pollCount := 0
	var pollResp *ColtRecommendationState
	for {
		time.Sleep(10 * time.Second)

		duration := time.Since(pollStart).Seconds()
		cmh.postUpdateState(duration, &messageCount, requestObj.RequestID, requestObj.RecommendationID)

		var code int
		var err error
		pollResp, code, err = cmh.doCheckRecommendationStatus(requestObj.RequestID, requestObj.RecommendationID)
		if err != nil || code != http.StatusOK {
			msg := fmt.Sprintf("Unable to check status of recommendation %s for service change request %s: %d - %s", requestObj.RecommendationID, requestObj.RequestID, code, err.Error())
			logger.Log.Errorf("%s%s", logPrefix, msg)
			cmh.writeResult(requestObj.RequestID, requestObj.RecommendationID, "FAILED", msg)
			return true
		}

		if pollResp.State == "PENDING" || pollResp.State == "INPROGRESS" {
			continue
		}

		if pollCount >= cmh.statusRetryCount {
			// Too many attempts, just fail this request
			msg := fmt.Sprintf("Unable to check status of recommendation %s for service change request %s: Request timed out waiting for Recommendation completion", requestObj.RecommendationID, requestObj.RequestID)
			logger.Log.Errorf("%s%s", logPrefix, msg)
			cmh.writeResult(requestObj.RequestID, requestObj.RecommendationID, "FAILED", msg)
			return true
		}

		break
	}

	logger.Log.Infof("%sService change status check completed for recommendation %s for service change request %s with result %s", logPrefix, requestObj.RecommendationID, requestObj.RequestID, pollResp.State)
	cmh.writeResult(requestObj.RequestID, requestObj.RecommendationID, pollResp.State, "")
	return true
}

func (cmh *ColtMEFHandler) postUpdateState(duration float64, messageCount *int, requestID string, recommendationID string) {
	if duration > cmh.pollCheckpoint1 && *messageCount == 0 {
		msg := fmt.Sprintf("Service Change %s for recommendation %s has been in progress for over %.1f seconds", requestID, recommendationID, cmh.pollCheckpoint1)
		postSlackUpdate(cmh.httpClient, requestID, msg)
		*messageCount = *messageCount + 1
	} else if duration > cmh.pollCheckpoint2 && *messageCount == 1 {
		msg := fmt.Sprintf("Service Change %s for recommendation %s has been in progress for over %.1f seconds and may be stuck", requestID, recommendationID, cmh.pollCheckpoint2)
		postSlackUpdate(cmh.httpClient, requestID, msg)
		*messageCount = *messageCount + 1
	} else if duration > cmh.pollCheckpoint3 && *messageCount == 2 {
		msg := fmt.Sprintf("Service Change %s for recommendation %s has been in progress for over %.1f seconds and is certainly stuck, contact novitasdevops@colt.net", requestID, recommendationID, cmh.pollCheckpoint3)
		postSlackUpdate(cmh.httpClient, requestID, msg)
		*messageCount = *messageCount + 1
	}
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
		return fmt.Errorf("Error marshalling result for service change request %s and recommendation %s: %s", reqID, recID, msgErr.Error())
	}

	logger.Log.Debugf("%sCompleted Change Service Request %s: Recommendation %s Status %s Message %s", logPrefix, result.RequestID, result.RecommendationID, result.Status, result.ErrorMessage)

	return cmh.resultWriter.WriteMessage("Result:"+reqID, msgBytes)
}

// writePending - helper to write a result to the service change pending topic
// func (cmh *ColtMEFHandler) writePending(reqID string, recID string) error {
// 	msgBytes, msgErr := createPendingPayload(reqID, recID)
// 	if msgErr != nil {
// 		return fmt.Errorf("Error marshalling pending object for service change request %s and recommendation ID %s: %s", reqID, recID, msgErr.Error())
// 	}

// 	logger.Log.Debugf("%sService change request %s moving to Pending state for recommendation ID %s", logPrefix, reqID, recID)

// 	return cmh.pendingWriter.WriteMessage("Pending:"+reqID, msgBytes)
// }

func createPendingPayload(reqID string, recID string) ([]byte, error) {
	result := ServiceChangeCheckStatusRequest{
		RequestID:        reqID,
		RecommendationID: recID,
	}

	return json.Marshal(result)
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

func postSlackUpdate(client *http.Client, serviceChangeID string, payload string) {

	reqPaylod := map[string]interface{}{}
	reqPaylod["text"] = payload

	payloadBytes, err := json.Marshal(reqPaylod)
	if err != nil {
		logger.Log.Errorf("%sUnable to unmarshal payload for slack message for service change request %s: %s", logPrefix, serviceChangeID, err.Error())
	}

	req, err := http.NewRequest("POST", slackURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Log.Errorf("%sUnable to build slack message for service change request %s: %s", logPrefix, serviceChangeID, err.Error())
	}

	// Fill in necessary request headers
	req.Header.Set("Content-Type", "application/json")

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("%s Submitting slack message of status for service change request %s: %s", logPrefix, serviceChangeID, payload)
	}

	// Issue request to COlt
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Errorf("%sUnable to issue slack update post for service change request %s: %s", logPrefix, serviceChangeID, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("%sError completing slack update post for service change request %s: Response Code %d", logPrefix, serviceChangeID, resp.StatusCode)
	}
}
