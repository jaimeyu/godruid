// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// UpdateTenantOKCode is the HTTP code returned for type UpdateTenantOK
const UpdateTenantOKCode int = 200

/*UpdateTenantOK update tenant o k

swagger:response updateTenantOK
*/
type UpdateTenantOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenant `json:"body,omitempty"`
}

// NewUpdateTenantOK creates UpdateTenantOK with default headers values
func NewUpdateTenantOK() *UpdateTenantOK {

	return &UpdateTenantOK{}
}

// WithPayload adds the payload to the update tenant o k response
func (o *UpdateTenantOK) WithPayload(payload *swagmodels.JSONAPITenant) *UpdateTenantOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update tenant o k response
func (o *UpdateTenantOK) SetPayload(payload *swagmodels.JSONAPITenant) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTenantOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTenantBadRequestCode is the HTTP code returned for type UpdateTenantBadRequest
const UpdateTenantBadRequestCode int = 400

/*UpdateTenantBadRequest Request data does not pass validation

swagger:response updateTenantBadRequest
*/
type UpdateTenantBadRequest struct {
}

// NewUpdateTenantBadRequest creates UpdateTenantBadRequest with default headers values
func NewUpdateTenantBadRequest() *UpdateTenantBadRequest {

	return &UpdateTenantBadRequest{}
}

// WriteResponse to the client
func (o *UpdateTenantBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// UpdateTenantInternalServerErrorCode is the HTTP code returned for type UpdateTenantInternalServerError
const UpdateTenantInternalServerErrorCode int = 500

/*UpdateTenantInternalServerError Unexpected error processing request

swagger:response updateTenantInternalServerError
*/
type UpdateTenantInternalServerError struct {
}

// NewUpdateTenantInternalServerError creates UpdateTenantInternalServerError with default headers values
func NewUpdateTenantInternalServerError() *UpdateTenantInternalServerError {

	return &UpdateTenantInternalServerError{}
}

// WriteResponse to the client
func (o *UpdateTenantInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
