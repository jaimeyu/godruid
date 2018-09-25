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
	"github.com/accedian/adh-gather/models"
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

	// Turn the query Params into the request object:
	requestObj := &ColtRecommendation{}
	err = json.Unmarshal(requestBytes, requestObj)
	if err != nil {
		msg := fmt.Sprintf("Unable to read service change data: %s", err.Error())
		reportError(w, startTime, "400", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Issuing Service Change request for: %s", models.AsJSONString(requestObj))

	requestObjBytes, err := json.Marshal(requestObj)
	if err != nil {
		msg := fmt.Sprintf("Unable to prepare service change data: %s", err.Error())
		reportError(w, startTime, "400", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
		return
	}

	// Setup the request to Colt
	req, err := http.NewRequest("POST", cmh.server, bytes.NewBuffer(requestObjBytes))
	if err != nil {
		msg := fmt.Sprintf("Unable to prepare service change request: %s", err.Error())
		reportError(w, startTime, "400", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
		return
	}

	req.Header.Set("x-colt-app-id", cmh.appID)
	req.Header.Set("x-colt-app-sig", getAuthHeader(requestObjBytes, cmh.sharedSecret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cmh.httpClient.Do(req)
	if err != nil {
		msg := fmt.Sprintf("Unable to issue service change: %s", err.Error())
		reportError(w, startTime, "500", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("Unable to read service change response: %s", err.Error())
		reportError(w, startTime, "500", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusCreated {

		responseObj := &ColtError{}
		err = json.Unmarshal(respBytes, responseObj)
		if err != nil {
			msg := fmt.Sprintf("Unable to unmarshal service change response: %s", err.Error())
			reportError(w, startTime, "500", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf("Service change failed: %d - %s", responseObj.Code, responseObj.Message)
		reportError(w, startTime, "500", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
		return
	}

	responseObj := &ColtRecommendationResponse{}
	err = json.Unmarshal(respBytes, responseObj)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal service change response: %s", err.Error())
		reportError(w, startTime, "500", makeRecommendationAPIStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Completed service change: %s", db.HistogramStr, models.AsJSONString(requestObj))
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

func getAuthHeader(recommendation []byte, key string) string {
	hashedPayload := base64HMACSHA256(recommendation, key)

	timeNow := time.Now().UTC()
	dateY, dateM, DateD := timeNow.Date()
	hour := timeNow.Hour()

	timestamp := fmt.Sprintf("%04d%02d%02d%02d", dateY, dateM, DateD, hour)
	requestData := strings.Join([]string{timestamp, recommendationRequestPath, hashedPayload}, "")

	return base64HMACSHA256([]byte(requestData), key)
}

func base64HMACSHA256(payload []byte, key string) string {
	hashObj := hmac.New(sha256.New, []byte(key))
	hashObj.Write(payload)
	return base64.StdEncoding.EncodeToString(hashObj.Sum(nil))
}
