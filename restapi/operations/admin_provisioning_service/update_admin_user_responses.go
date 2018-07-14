// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// UpdateAdminUserOKCode is the HTTP code returned for type UpdateAdminUserOK
const UpdateAdminUserOKCode int = 200

/*UpdateAdminUserOK update admin user o k

swagger:response updateAdminUserOK
*/
type UpdateAdminUserOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPIAdminUser `json:"body,omitempty"`
}

// NewUpdateAdminUserOK creates UpdateAdminUserOK with default headers values
func NewUpdateAdminUserOK() *UpdateAdminUserOK {

	return &UpdateAdminUserOK{}
}

// WithPayload adds the payload to the update admin user o k response
func (o *UpdateAdminUserOK) WithPayload(payload *swagmodels.JSONAPIAdminUser) *UpdateAdminUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update admin user o k response
func (o *UpdateAdminUserOK) SetPayload(payload *swagmodels.JSONAPIAdminUser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateAdminUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateAdminUserBadRequestCode is the HTTP code returned for type UpdateAdminUserBadRequest
const UpdateAdminUserBadRequestCode int = 400

/*UpdateAdminUserBadRequest Request data does not pass validation

swagger:response updateAdminUserBadRequest
*/
type UpdateAdminUserBadRequest struct {
}

// NewUpdateAdminUserBadRequest creates UpdateAdminUserBadRequest with default headers values
func NewUpdateAdminUserBadRequest() *UpdateAdminUserBadRequest {

	return &UpdateAdminUserBadRequest{}
}

// WriteResponse to the client
func (o *UpdateAdminUserBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// UpdateAdminUserInternalServerErrorCode is the HTTP code returned for type UpdateAdminUserInternalServerError
const UpdateAdminUserInternalServerErrorCode int = 500

/*UpdateAdminUserInternalServerError Unexpected error processing request

swagger:response updateAdminUserInternalServerError
*/
type UpdateAdminUserInternalServerError struct {
}

// NewUpdateAdminUserInternalServerError creates UpdateAdminUserInternalServerError with default headers values
func NewUpdateAdminUserInternalServerError() *UpdateAdminUserInternalServerError {

	return &UpdateAdminUserInternalServerError{}
}

// WriteResponse to the client
func (o *UpdateAdminUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
