// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetAllReportScheduleConfigsV2OKCode is the HTTP code returned for type GetAllReportScheduleConfigsV2OK
const GetAllReportScheduleConfigsV2OKCode int = 200

/*GetAllReportScheduleConfigsV2OK get all report schedule configs v2 o k

swagger:response getAllReportScheduleConfigsV2OK
*/
type GetAllReportScheduleConfigsV2OK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.ReportScheduleConfigListResponse `json:"body,omitempty"`
}

// NewGetAllReportScheduleConfigsV2OK creates GetAllReportScheduleConfigsV2OK with default headers values
func NewGetAllReportScheduleConfigsV2OK() *GetAllReportScheduleConfigsV2OK {

	return &GetAllReportScheduleConfigsV2OK{}
}

// WithPayload adds the payload to the get all report schedule configs v2 o k response
func (o *GetAllReportScheduleConfigsV2OK) WithPayload(payload *swagmodels.ReportScheduleConfigListResponse) *GetAllReportScheduleConfigsV2OK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all report schedule configs v2 o k response
func (o *GetAllReportScheduleConfigsV2OK) SetPayload(payload *swagmodels.ReportScheduleConfigListResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReportScheduleConfigsV2OK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAllReportScheduleConfigsV2ForbiddenCode is the HTTP code returned for type GetAllReportScheduleConfigsV2Forbidden
const GetAllReportScheduleConfigsV2ForbiddenCode int = 403

/*GetAllReportScheduleConfigsV2Forbidden Requestor does not have authorization to perform this action

swagger:response getAllReportScheduleConfigsV2Forbidden
*/
type GetAllReportScheduleConfigsV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAllReportScheduleConfigsV2Forbidden creates GetAllReportScheduleConfigsV2Forbidden with default headers values
func NewGetAllReportScheduleConfigsV2Forbidden() *GetAllReportScheduleConfigsV2Forbidden {

	return &GetAllReportScheduleConfigsV2Forbidden{}
}

// WithPayload adds the payload to the get all report schedule configs v2 forbidden response
func (o *GetAllReportScheduleConfigsV2Forbidden) WithPayload(payload string) *GetAllReportScheduleConfigsV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all report schedule configs v2 forbidden response
func (o *GetAllReportScheduleConfigsV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReportScheduleConfigsV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetAllReportScheduleConfigsV2NotFoundCode is the HTTP code returned for type GetAllReportScheduleConfigsV2NotFound
const GetAllReportScheduleConfigsV2NotFoundCode int = 404

/*GetAllReportScheduleConfigsV2NotFound No Report Schedule Configurations are provisioned

swagger:response getAllReportScheduleConfigsV2NotFound
*/
type GetAllReportScheduleConfigsV2NotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAllReportScheduleConfigsV2NotFound creates GetAllReportScheduleConfigsV2NotFound with default headers values
func NewGetAllReportScheduleConfigsV2NotFound() *GetAllReportScheduleConfigsV2NotFound {

	return &GetAllReportScheduleConfigsV2NotFound{}
}

// WithPayload adds the payload to the get all report schedule configs v2 not found response
func (o *GetAllReportScheduleConfigsV2NotFound) WithPayload(payload string) *GetAllReportScheduleConfigsV2NotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all report schedule configs v2 not found response
func (o *GetAllReportScheduleConfigsV2NotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReportScheduleConfigsV2NotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetAllReportScheduleConfigsV2InternalServerErrorCode is the HTTP code returned for type GetAllReportScheduleConfigsV2InternalServerError
const GetAllReportScheduleConfigsV2InternalServerErrorCode int = 500

/*GetAllReportScheduleConfigsV2InternalServerError Unexpected error processing request

swagger:response getAllReportScheduleConfigsV2InternalServerError
*/
type GetAllReportScheduleConfigsV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAllReportScheduleConfigsV2InternalServerError creates GetAllReportScheduleConfigsV2InternalServerError with default headers values
func NewGetAllReportScheduleConfigsV2InternalServerError() *GetAllReportScheduleConfigsV2InternalServerError {

	return &GetAllReportScheduleConfigsV2InternalServerError{}
}

// WithPayload adds the payload to the get all report schedule configs v2 internal server error response
func (o *GetAllReportScheduleConfigsV2InternalServerError) WithPayload(payload string) *GetAllReportScheduleConfigsV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all report schedule configs v2 internal server error response
func (o *GetAllReportScheduleConfigsV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllReportScheduleConfigsV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
