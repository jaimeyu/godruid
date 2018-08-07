// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// DeleteTenantOKCode is the HTTP code returned for type DeleteTenantOK
const DeleteTenantOKCode int = 200

/*DeleteTenantOK delete tenant o k

swagger:response deleteTenantOK
*/
type DeleteTenantOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenant `json:"body,omitempty"`
}

// NewDeleteTenantOK creates DeleteTenantOK with default headers values
func NewDeleteTenantOK() *DeleteTenantOK {

	return &DeleteTenantOK{}
}

// WithPayload adds the payload to the delete tenant o k response
func (o *DeleteTenantOK) WithPayload(payload *swagmodels.JSONAPITenant) *DeleteTenantOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete tenant o k response
func (o *DeleteTenantOK) SetPayload(payload *swagmodels.JSONAPITenant) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTenantOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteTenantForbiddenCode is the HTTP code returned for type DeleteTenantForbidden
const DeleteTenantForbiddenCode int = 403

/*DeleteTenantForbidden Requestor does not have authorization to perform this action

swagger:response deleteTenantForbidden
*/
type DeleteTenantForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewDeleteTenantForbidden creates DeleteTenantForbidden with default headers values
func NewDeleteTenantForbidden() *DeleteTenantForbidden {

	return &DeleteTenantForbidden{}
}

// WithPayload adds the payload to the delete tenant forbidden response
func (o *DeleteTenantForbidden) WithPayload(payload string) *DeleteTenantForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete tenant forbidden response
func (o *DeleteTenantForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTenantForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// DeleteTenantInternalServerErrorCode is the HTTP code returned for type DeleteTenantInternalServerError
const DeleteTenantInternalServerErrorCode int = 500

/*DeleteTenantInternalServerError Unexpected error processing request

swagger:response deleteTenantInternalServerError
*/
type DeleteTenantInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewDeleteTenantInternalServerError creates DeleteTenantInternalServerError with default headers values
func NewDeleteTenantInternalServerError() *DeleteTenantInternalServerError {

	return &DeleteTenantInternalServerError{}
}

// WithPayload adds the payload to the delete tenant internal server error response
func (o *DeleteTenantInternalServerError) WithPayload(payload string) *DeleteTenantInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete tenant internal server error response
func (o *DeleteTenantInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTenantInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
