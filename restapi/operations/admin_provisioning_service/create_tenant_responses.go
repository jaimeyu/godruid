// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// CreateTenantOKCode is the HTTP code returned for type CreateTenantOK
const CreateTenantOKCode int = 200

/*CreateTenantOK create tenant o k

swagger:response createTenantOK
*/
type CreateTenantOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenant `json:"body,omitempty"`
}

// NewCreateTenantOK creates CreateTenantOK with default headers values
func NewCreateTenantOK() *CreateTenantOK {

	return &CreateTenantOK{}
}

// WithPayload adds the payload to the create tenant o k response
func (o *CreateTenantOK) WithPayload(payload *swagmodels.JSONAPITenant) *CreateTenantOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant o k response
func (o *CreateTenantOK) SetPayload(payload *swagmodels.JSONAPITenant) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTenantBadRequestCode is the HTTP code returned for type CreateTenantBadRequest
const CreateTenantBadRequestCode int = 400

/*CreateTenantBadRequest Request data does not pass validation

swagger:response createTenantBadRequest
*/
type CreateTenantBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantBadRequest creates CreateTenantBadRequest with default headers values
func NewCreateTenantBadRequest() *CreateTenantBadRequest {

	return &CreateTenantBadRequest{}
}

// WithPayload adds the payload to the create tenant bad request response
func (o *CreateTenantBadRequest) WithPayload(payload string) *CreateTenantBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant bad request response
func (o *CreateTenantBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantForbiddenCode is the HTTP code returned for type CreateTenantForbidden
const CreateTenantForbiddenCode int = 403

/*CreateTenantForbidden Requestor does not have authorization to perform this action

swagger:response createTenantForbidden
*/
type CreateTenantForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantForbidden creates CreateTenantForbidden with default headers values
func NewCreateTenantForbidden() *CreateTenantForbidden {

	return &CreateTenantForbidden{}
}

// WithPayload adds the payload to the create tenant forbidden response
func (o *CreateTenantForbidden) WithPayload(payload string) *CreateTenantForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant forbidden response
func (o *CreateTenantForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantConflictCode is the HTTP code returned for type CreateTenantConflict
const CreateTenantConflictCode int = 409

/*CreateTenantConflict The Tenant being provisioned already exists

swagger:response createTenantConflict
*/
type CreateTenantConflict struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantConflict creates CreateTenantConflict with default headers values
func NewCreateTenantConflict() *CreateTenantConflict {

	return &CreateTenantConflict{}
}

// WithPayload adds the payload to the create tenant conflict response
func (o *CreateTenantConflict) WithPayload(payload string) *CreateTenantConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant conflict response
func (o *CreateTenantConflict) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantInternalServerErrorCode is the HTTP code returned for type CreateTenantInternalServerError
const CreateTenantInternalServerErrorCode int = 500

/*CreateTenantInternalServerError Unexpected error processing request

swagger:response createTenantInternalServerError
*/
type CreateTenantInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantInternalServerError creates CreateTenantInternalServerError with default headers values
func NewCreateTenantInternalServerError() *CreateTenantInternalServerError {

	return &CreateTenantInternalServerError{}
}

// WithPayload adds the payload to the create tenant internal server error response
func (o *CreateTenantInternalServerError) WithPayload(payload string) *CreateTenantInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant internal server error response
func (o *CreateTenantInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
