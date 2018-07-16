// Code generated by go-swagger; DO NOT EDIT.

package metrics_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetThresholdCrossingByMonitoredObjectOKCode is the HTTP code returned for type GetThresholdCrossingByMonitoredObjectOK
const GetThresholdCrossingByMonitoredObjectOKCode int = 200

/*GetThresholdCrossingByMonitoredObjectOK get threshold crossing by monitored object o k

swagger:response getThresholdCrossingByMonitoredObjectOK
*/
type GetThresholdCrossingByMonitoredObjectOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.GathergrpcJSONAPIObject `json:"body,omitempty"`
}

// NewGetThresholdCrossingByMonitoredObjectOK creates GetThresholdCrossingByMonitoredObjectOK with default headers values
func NewGetThresholdCrossingByMonitoredObjectOK() *GetThresholdCrossingByMonitoredObjectOK {

	return &GetThresholdCrossingByMonitoredObjectOK{}
}

// WithPayload adds the payload to the get threshold crossing by monitored object o k response
func (o *GetThresholdCrossingByMonitoredObjectOK) WithPayload(payload *swagmodels.GathergrpcJSONAPIObject) *GetThresholdCrossingByMonitoredObjectOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get threshold crossing by monitored object o k response
func (o *GetThresholdCrossingByMonitoredObjectOK) SetPayload(payload *swagmodels.GathergrpcJSONAPIObject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetThresholdCrossingByMonitoredObjectOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetThresholdCrossingByMonitoredObjectBadRequestCode is the HTTP code returned for type GetThresholdCrossingByMonitoredObjectBadRequest
const GetThresholdCrossingByMonitoredObjectBadRequestCode int = 400

/*GetThresholdCrossingByMonitoredObjectBadRequest Request data does not pass validation

swagger:response getThresholdCrossingByMonitoredObjectBadRequest
*/
type GetThresholdCrossingByMonitoredObjectBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetThresholdCrossingByMonitoredObjectBadRequest creates GetThresholdCrossingByMonitoredObjectBadRequest with default headers values
func NewGetThresholdCrossingByMonitoredObjectBadRequest() *GetThresholdCrossingByMonitoredObjectBadRequest {

	return &GetThresholdCrossingByMonitoredObjectBadRequest{}
}

// WithPayload adds the payload to the get threshold crossing by monitored object bad request response
func (o *GetThresholdCrossingByMonitoredObjectBadRequest) WithPayload(payload string) *GetThresholdCrossingByMonitoredObjectBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get threshold crossing by monitored object bad request response
func (o *GetThresholdCrossingByMonitoredObjectBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetThresholdCrossingByMonitoredObjectBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetThresholdCrossingByMonitoredObjectForbiddenCode is the HTTP code returned for type GetThresholdCrossingByMonitoredObjectForbidden
const GetThresholdCrossingByMonitoredObjectForbiddenCode int = 403

/*GetThresholdCrossingByMonitoredObjectForbidden Requestor does not have authorization to perform this action

swagger:response getThresholdCrossingByMonitoredObjectForbidden
*/
type GetThresholdCrossingByMonitoredObjectForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetThresholdCrossingByMonitoredObjectForbidden creates GetThresholdCrossingByMonitoredObjectForbidden with default headers values
func NewGetThresholdCrossingByMonitoredObjectForbidden() *GetThresholdCrossingByMonitoredObjectForbidden {

	return &GetThresholdCrossingByMonitoredObjectForbidden{}
}

// WithPayload adds the payload to the get threshold crossing by monitored object forbidden response
func (o *GetThresholdCrossingByMonitoredObjectForbidden) WithPayload(payload string) *GetThresholdCrossingByMonitoredObjectForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get threshold crossing by monitored object forbidden response
func (o *GetThresholdCrossingByMonitoredObjectForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetThresholdCrossingByMonitoredObjectForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetThresholdCrossingByMonitoredObjectNotFoundCode is the HTTP code returned for type GetThresholdCrossingByMonitoredObjectNotFound
const GetThresholdCrossingByMonitoredObjectNotFoundCode int = 404

/*GetThresholdCrossingByMonitoredObjectNotFound Threshold profile not found

swagger:response getThresholdCrossingByMonitoredObjectNotFound
*/
type GetThresholdCrossingByMonitoredObjectNotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetThresholdCrossingByMonitoredObjectNotFound creates GetThresholdCrossingByMonitoredObjectNotFound with default headers values
func NewGetThresholdCrossingByMonitoredObjectNotFound() *GetThresholdCrossingByMonitoredObjectNotFound {

	return &GetThresholdCrossingByMonitoredObjectNotFound{}
}

// WithPayload adds the payload to the get threshold crossing by monitored object not found response
func (o *GetThresholdCrossingByMonitoredObjectNotFound) WithPayload(payload string) *GetThresholdCrossingByMonitoredObjectNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get threshold crossing by monitored object not found response
func (o *GetThresholdCrossingByMonitoredObjectNotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetThresholdCrossingByMonitoredObjectNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetThresholdCrossingByMonitoredObjectInternalServerErrorCode is the HTTP code returned for type GetThresholdCrossingByMonitoredObjectInternalServerError
const GetThresholdCrossingByMonitoredObjectInternalServerErrorCode int = 500

/*GetThresholdCrossingByMonitoredObjectInternalServerError Unexpected error processing request

swagger:response getThresholdCrossingByMonitoredObjectInternalServerError
*/
type GetThresholdCrossingByMonitoredObjectInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetThresholdCrossingByMonitoredObjectInternalServerError creates GetThresholdCrossingByMonitoredObjectInternalServerError with default headers values
func NewGetThresholdCrossingByMonitoredObjectInternalServerError() *GetThresholdCrossingByMonitoredObjectInternalServerError {

	return &GetThresholdCrossingByMonitoredObjectInternalServerError{}
}

// WithPayload adds the payload to the get threshold crossing by monitored object internal server error response
func (o *GetThresholdCrossingByMonitoredObjectInternalServerError) WithPayload(payload string) *GetThresholdCrossingByMonitoredObjectInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get threshold crossing by monitored object internal server error response
func (o *GetThresholdCrossingByMonitoredObjectInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetThresholdCrossingByMonitoredObjectInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
