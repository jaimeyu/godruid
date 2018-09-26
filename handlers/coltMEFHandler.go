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

	server       string
	appID        string
	sharedSecret string
}

func CreateColtMEFHandler() *ColtMEFHandler {
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

// MakeRecommendation - Recommend a service change.
func (cmh *ColtMEFHandler) MakeRecommendation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", makeRecommendationAPIStr, msg, http.StatusBadRequest)
		return
	}

	responseObj, code, err := cmh.doMakeRecommendation(requestBytes)
	if err != nil {
		reportError(w, startTime, string(code), makeRecommendationAPIStr, err.Error(), code)
		return
	}

	logger.Log.Infof("Completed service change: %s", db.HistogramStr, string(requestBytes))
	trackAPIMetrics(startTime, "200", makeRecommendationAPIStr)
	fmt.Fprintf(w, responseObj.RecommendationID)
}

type ColtRecommendation struct {
	SereviceID      string `json:"service_id"`
	Action          string `json:"action"`
	BandwidthChange int    `json:"bandwidth_change,omitempty"`
}

type ColtError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ColtRecommendationResponse struct {
	RecommendationID string `json:"recommendation_id"`
}

type ColtRecommendationState struct {
	State string `json:"state"`
}

func getAuthHeader(recommendation []byte, key string, path string) string {
	hashedPayload := base64HMACSHA256(recommendation, key)

	timeNow := time.Now().UTC()
	dateY, dateM, DateD := timeNow.Date()
	hour := timeNow.Hour()

	timestamp := fmt.Sprintf("%04d%02d%02d%02d", dateY, dateM, DateD, hour)
	requestData := strings.Join([]string{timestamp, path, hashedPayload}, "")

	return base64HMACSHA256([]byte(requestData), key)
}

func base64HMACSHA256(payload []byte, key string) string {
	hashObj := hmac.New(sha256.New, []byte(key))
	hashObj.Write(payload)
	return base64.StdEncoding.EncodeToString(hashObj.Sum(nil))
}

func (cmh *ColtMEFHandler) doMakeRecommendation(requestBytes []byte) (*ColtRecommendationResponse, int, error) {

	// Deserialize the request
	requestObj := &ColtRecommendation{}
	err := json.Unmarshal(requestBytes, requestObj)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to read service change data: %s", err.Error())
	}

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

func handleRecommendationRequest(recommendationRequest []byte) bool {
	// Deserialize the object
	return true
}
