// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// CreateTenantDomainOKCode is the HTTP code returned for type CreateTenantDomainOK
const CreateTenantDomainOKCode int = 200

/*CreateTenantDomainOK create tenant domain o k

swagger:response createTenantDomainOK
*/
type CreateTenantDomainOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantDomain `json:"body,omitempty"`
}

// NewCreateTenantDomainOK creates CreateTenantDomainOK with default headers values
func NewCreateTenantDomainOK() *CreateTenantDomainOK {

	return &CreateTenantDomainOK{}
}

// WithPayload adds the payload to the create tenant domain o k response
func (o *CreateTenantDomainOK) WithPayload(payload *swagmodels.JSONAPITenantDomain) *CreateTenantDomainOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant domain o k response
func (o *CreateTenantDomainOK) SetPayload(payload *swagmodels.JSONAPITenantDomain) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantDomainOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTenantDomainBadRequestCode is the HTTP code returned for type CreateTenantDomainBadRequest
const CreateTenantDomainBadRequestCode int = 400

/*CreateTenantDomainBadRequest Request data does not pass validation

swagger:response createTenantDomainBadRequest
*/
type CreateTenantDomainBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantDomainBadRequest creates CreateTenantDomainBadRequest with default headers values
func NewCreateTenantDomainBadRequest() *CreateTenantDomainBadRequest {

	return &CreateTenantDomainBadRequest{}
}

// WithPayload adds the payload to the create tenant domain bad request response
func (o *CreateTenantDomainBadRequest) WithPayload(payload string) *CreateTenantDomainBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant domain bad request response
func (o *CreateTenantDomainBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantDomainBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantDomainForbiddenCode is the HTTP code returned for type CreateTenantDomainForbidden
const CreateTenantDomainForbiddenCode int = 403

/*CreateTenantDomainForbidden Requestor does not have authorization to perform this action

swagger:response createTenantDomainForbidden
*/
type CreateTenantDomainForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantDomainForbidden creates CreateTenantDomainForbidden with default headers values
func NewCreateTenantDomainForbidden() *CreateTenantDomainForbidden {

	return &CreateTenantDomainForbidden{}
}

// WithPayload adds the payload to the create tenant domain forbidden response
func (o *CreateTenantDomainForbidden) WithPayload(payload string) *CreateTenantDomainForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant domain forbidden response
func (o *CreateTenantDomainForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantDomainForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantDomainInternalServerErrorCode is the HTTP code returned for type CreateTenantDomainInternalServerError
const CreateTenantDomainInternalServerErrorCode int = 500

/*CreateTenantDomainInternalServerError Unexpected error processing request

swagger:response createTenantDomainInternalServerError
*/
type CreateTenantDomainInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantDomainInternalServerError creates CreateTenantDomainInternalServerError with default headers values
func NewCreateTenantDomainInternalServerError() *CreateTenantDomainInternalServerError {

	return &CreateTenantDomainInternalServerError{}
}

// WithPayload adds the payload to the create tenant domain internal server error response
func (o *CreateTenantDomainInternalServerError) WithPayload(payload string) *CreateTenantDomainInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant domain internal server error response
func (o *CreateTenantDomainInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantDomainInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
