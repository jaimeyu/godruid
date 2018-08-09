// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetSLAReportV2OKCode is the HTTP code returned for type GetSLAReportV2OK
const GetSLAReportV2OKCode int = 200

/*GetSLAReportV2OK get Sla report v2 o k

swagger:response getSlaReportV2OK
*/
type GetSLAReportV2OK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.GathergrpcJSONAPIObject `json:"body,omitempty"`
}

// NewGetSLAReportV2OK creates GetSLAReportV2OK with default headers values
func NewGetSLAReportV2OK() *GetSLAReportV2OK {

	return &GetSLAReportV2OK{}
}

// WithPayload adds the payload to the get Sla report v2 o k response
func (o *GetSLAReportV2OK) WithPayload(payload *swagmodels.GathergrpcJSONAPIObject) *GetSLAReportV2OK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Sla report v2 o k response
func (o *GetSLAReportV2OK) SetPayload(payload *swagmodels.GathergrpcJSONAPIObject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSLAReportV2OK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetSLAReportV2ForbiddenCode is the HTTP code returned for type GetSLAReportV2Forbidden
const GetSLAReportV2ForbiddenCode int = 403

/*GetSLAReportV2Forbidden Requestor does not have authorization to perform this action

swagger:response getSlaReportV2Forbidden
*/
type GetSLAReportV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetSLAReportV2Forbidden creates GetSLAReportV2Forbidden with default headers values
func NewGetSLAReportV2Forbidden() *GetSLAReportV2Forbidden {

	return &GetSLAReportV2Forbidden{}
}

// WithPayload adds the payload to the get Sla report v2 forbidden response
func (o *GetSLAReportV2Forbidden) WithPayload(payload string) *GetSLAReportV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Sla report v2 forbidden response
func (o *GetSLAReportV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSLAReportV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetSLAReportV2NotFoundCode is the HTTP code returned for type GetSLAReportV2NotFound
const GetSLAReportV2NotFoundCode int = 404

/*GetSLAReportV2NotFound The requested Report is not provisioned in Datahub

swagger:response getSlaReportV2NotFound
*/
type GetSLAReportV2NotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetSLAReportV2NotFound creates GetSLAReportV2NotFound with default headers values
func NewGetSLAReportV2NotFound() *GetSLAReportV2NotFound {

	return &GetSLAReportV2NotFound{}
}

// WithPayload adds the payload to the get Sla report v2 not found response
func (o *GetSLAReportV2NotFound) WithPayload(payload string) *GetSLAReportV2NotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Sla report v2 not found response
func (o *GetSLAReportV2NotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSLAReportV2NotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetSLAReportV2InternalServerErrorCode is the HTTP code returned for type GetSLAReportV2InternalServerError
const GetSLAReportV2InternalServerErrorCode int = 500

/*GetSLAReportV2InternalServerError Unexpected error processing request

swagger:response getSlaReportV2InternalServerError
*/
type GetSLAReportV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetSLAReportV2InternalServerError creates GetSLAReportV2InternalServerError with default headers values
func NewGetSLAReportV2InternalServerError() *GetSLAReportV2InternalServerError {

	return &GetSLAReportV2InternalServerError{}
}

// WithPayload adds the payload to the get Sla report v2 internal server error response
func (o *GetSLAReportV2InternalServerError) WithPayload(payload string) *GetSLAReportV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Sla report v2 internal server error response
func (o *GetSLAReportV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSLAReportV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
