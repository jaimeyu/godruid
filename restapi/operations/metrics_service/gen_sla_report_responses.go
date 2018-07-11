// Code generated by go-swagger; DO NOT EDIT.

package metrics_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GenSLAReportOKCode is the HTTP code returned for type GenSLAReportOK
const GenSLAReportOKCode int = 200

/*GenSLAReportOK gen Sla report o k

swagger:response genSlaReportOK
*/
type GenSLAReportOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.GathergrpcJSONAPIObject `json:"body,omitempty"`
}

// NewGenSLAReportOK creates GenSLAReportOK with default headers values
func NewGenSLAReportOK() *GenSLAReportOK {

	return &GenSLAReportOK{}
}

// WithPayload adds the payload to the gen Sla report o k response
func (o *GenSLAReportOK) WithPayload(payload *swagmodels.GathergrpcJSONAPIObject) *GenSLAReportOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the gen Sla report o k response
func (o *GenSLAReportOK) SetPayload(payload *swagmodels.GathergrpcJSONAPIObject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenSLAReportOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GenSLAReportBadRequestCode is the HTTP code returned for type GenSLAReportBadRequest
const GenSLAReportBadRequestCode int = 400

/*GenSLAReportBadRequest Request data does not pass validation

swagger:response genSlaReportBadRequest
*/
type GenSLAReportBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGenSLAReportBadRequest creates GenSLAReportBadRequest with default headers values
func NewGenSLAReportBadRequest() *GenSLAReportBadRequest {

	return &GenSLAReportBadRequest{}
}

// WithPayload adds the payload to the gen Sla report bad request response
func (o *GenSLAReportBadRequest) WithPayload(payload string) *GenSLAReportBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the gen Sla report bad request response
func (o *GenSLAReportBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenSLAReportBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GenSLAReportForbiddenCode is the HTTP code returned for type GenSLAReportForbidden
const GenSLAReportForbiddenCode int = 403

/*GenSLAReportForbidden Requestor does not have authorization to perform this action

swagger:response genSlaReportForbidden
*/
type GenSLAReportForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGenSLAReportForbidden creates GenSLAReportForbidden with default headers values
func NewGenSLAReportForbidden() *GenSLAReportForbidden {

	return &GenSLAReportForbidden{}
}

// WithPayload adds the payload to the gen Sla report forbidden response
func (o *GenSLAReportForbidden) WithPayload(payload string) *GenSLAReportForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the gen Sla report forbidden response
func (o *GenSLAReportForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenSLAReportForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GenSLAReportNotFoundCode is the HTTP code returned for type GenSLAReportNotFound
const GenSLAReportNotFoundCode int = 404

/*GenSLAReportNotFound Components of the request parameters were not in the provisioning database

swagger:response genSlaReportNotFound
*/
type GenSLAReportNotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGenSLAReportNotFound creates GenSLAReportNotFound with default headers values
func NewGenSLAReportNotFound() *GenSLAReportNotFound {

	return &GenSLAReportNotFound{}
}

// WithPayload adds the payload to the gen Sla report not found response
func (o *GenSLAReportNotFound) WithPayload(payload string) *GenSLAReportNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the gen Sla report not found response
func (o *GenSLAReportNotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenSLAReportNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GenSLAReportInternalServerErrorCode is the HTTP code returned for type GenSLAReportInternalServerError
const GenSLAReportInternalServerErrorCode int = 500

/*GenSLAReportInternalServerError Unexpected error processing request

swagger:response genSlaReportInternalServerError
*/
type GenSLAReportInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGenSLAReportInternalServerError creates GenSLAReportInternalServerError with default headers values
func NewGenSLAReportInternalServerError() *GenSLAReportInternalServerError {

	return &GenSLAReportInternalServerError{}
}

// WithPayload adds the payload to the gen Sla report internal server error response
func (o *GenSLAReportInternalServerError) WithPayload(payload string) *GenSLAReportInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the gen Sla report internal server error response
func (o *GenSLAReportInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenSLAReportInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
