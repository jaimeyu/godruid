// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// UpdateReportScheduleConfigOKCode is the HTTP code returned for type UpdateReportScheduleConfigOK
const UpdateReportScheduleConfigOKCode int = 200

/*UpdateReportScheduleConfigOK update report schedule config o k

swagger:response updateReportScheduleConfigOK
*/
type UpdateReportScheduleConfigOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantReportScheduleConfig `json:"body,omitempty"`
}

// NewUpdateReportScheduleConfigOK creates UpdateReportScheduleConfigOK with default headers values
func NewUpdateReportScheduleConfigOK() *UpdateReportScheduleConfigOK {

	return &UpdateReportScheduleConfigOK{}
}

// WithPayload adds the payload to the update report schedule config o k response
func (o *UpdateReportScheduleConfigOK) WithPayload(payload *swagmodels.JSONAPITenantReportScheduleConfig) *UpdateReportScheduleConfigOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update report schedule config o k response
func (o *UpdateReportScheduleConfigOK) SetPayload(payload *swagmodels.JSONAPITenantReportScheduleConfig) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateReportScheduleConfigOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateReportScheduleConfigBadRequestCode is the HTTP code returned for type UpdateReportScheduleConfigBadRequest
const UpdateReportScheduleConfigBadRequestCode int = 400

/*UpdateReportScheduleConfigBadRequest Request data does not pass validation

swagger:response updateReportScheduleConfigBadRequest
*/
type UpdateReportScheduleConfigBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateReportScheduleConfigBadRequest creates UpdateReportScheduleConfigBadRequest with default headers values
func NewUpdateReportScheduleConfigBadRequest() *UpdateReportScheduleConfigBadRequest {

	return &UpdateReportScheduleConfigBadRequest{}
}

// WithPayload adds the payload to the update report schedule config bad request response
func (o *UpdateReportScheduleConfigBadRequest) WithPayload(payload string) *UpdateReportScheduleConfigBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update report schedule config bad request response
func (o *UpdateReportScheduleConfigBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateReportScheduleConfigBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateReportScheduleConfigForbiddenCode is the HTTP code returned for type UpdateReportScheduleConfigForbidden
const UpdateReportScheduleConfigForbiddenCode int = 403

/*UpdateReportScheduleConfigForbidden Requestor does not have authorization to perform this action

swagger:response updateReportScheduleConfigForbidden
*/
type UpdateReportScheduleConfigForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateReportScheduleConfigForbidden creates UpdateReportScheduleConfigForbidden with default headers values
func NewUpdateReportScheduleConfigForbidden() *UpdateReportScheduleConfigForbidden {

	return &UpdateReportScheduleConfigForbidden{}
}

// WithPayload adds the payload to the update report schedule config forbidden response
func (o *UpdateReportScheduleConfigForbidden) WithPayload(payload string) *UpdateReportScheduleConfigForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update report schedule config forbidden response
func (o *UpdateReportScheduleConfigForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateReportScheduleConfigForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateReportScheduleConfigInternalServerErrorCode is the HTTP code returned for type UpdateReportScheduleConfigInternalServerError
const UpdateReportScheduleConfigInternalServerErrorCode int = 500

/*UpdateReportScheduleConfigInternalServerError Unexpected error processing request

swagger:response updateReportScheduleConfigInternalServerError
*/
type UpdateReportScheduleConfigInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateReportScheduleConfigInternalServerError creates UpdateReportScheduleConfigInternalServerError with default headers values
func NewUpdateReportScheduleConfigInternalServerError() *UpdateReportScheduleConfigInternalServerError {

	return &UpdateReportScheduleConfigInternalServerError{}
}

// WithPayload adds the payload to the update report schedule config internal server error response
func (o *UpdateReportScheduleConfigInternalServerError) WithPayload(payload string) *UpdateReportScheduleConfigInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update report schedule config internal server error response
func (o *UpdateReportScheduleConfigInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateReportScheduleConfigInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
