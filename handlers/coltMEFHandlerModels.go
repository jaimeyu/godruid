package handlers

// ColtRecommendation - direct model of the object sent to the Colt POST /api/performance/recommendation API
type ColtRecommendation struct {
	ServiceID       string `json:"service_id"`
	Action          string `json:"action"`
	BandwidthChange int    `json:"bandwidth_change,omitempty"`
}

// ColtError - direct model of the object returned by Colt APIs when there is an error
type ColtError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ColtRecommendationResponse - direct model of the successful response form the Colt POST /api/performance/recommendation API
type ColtRecommendationResponse struct {
	RecommendationID string `json:"recommendation_id"`
}

// ColtRecommendationState - direct model from the Colt GET /api/performance/recommendation/{recommendationID} API
type ColtRecommendationState struct {
	State string `json:"state"`
}

// ServiceChangeRequest - Datahub  model used to communicate a service change request
type ServiceChangeRequest struct {
	RequestID     string              `json:"requestId"`
	ServiceChange *ColtRecommendation `json:"serviceChange"`
}

// ServiceChangeCheckStatusRequest - Datahub  model used to communicate a request to check the status of a pending service change
type ServiceChangeCheckStatusRequest struct {
	RequestID        string `json:"requestId"`
	RecommendationID string `json:"recommendationID"`
}

// ServiceChangeResult - Datahub  model used to communicate the result of a Service Change request
type ServiceChangeResult struct {
	RequestID        string `json:"requestId"`
	RecommendationID string `json:"recommendationID"`
	Status           string `json:"status"`
	ErrorMessage     string `json:"errorMessage"`
}
