// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// CreateTenantThresholdProfileOKCode is the HTTP code returned for type CreateTenantThresholdProfileOK
const CreateTenantThresholdProfileOKCode int = 200

/*CreateTenantThresholdProfileOK create tenant threshold profile o k

swagger:response createTenantThresholdProfileOK
*/
type CreateTenantThresholdProfileOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantThresholdProfile `json:"body,omitempty"`
}

// NewCreateTenantThresholdProfileOK creates CreateTenantThresholdProfileOK with default headers values
func NewCreateTenantThresholdProfileOK() *CreateTenantThresholdProfileOK {

	return &CreateTenantThresholdProfileOK{}
}

// WithPayload adds the payload to the create tenant threshold profile o k response
func (o *CreateTenantThresholdProfileOK) WithPayload(payload *swagmodels.JSONAPITenantThresholdProfile) *CreateTenantThresholdProfileOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant threshold profile o k response
func (o *CreateTenantThresholdProfileOK) SetPayload(payload *swagmodels.JSONAPITenantThresholdProfile) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantThresholdProfileOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTenantThresholdProfileBadRequestCode is the HTTP code returned for type CreateTenantThresholdProfileBadRequest
const CreateTenantThresholdProfileBadRequestCode int = 400

/*CreateTenantThresholdProfileBadRequest Request data does not pass validation

swagger:response createTenantThresholdProfileBadRequest
*/
type CreateTenantThresholdProfileBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantThresholdProfileBadRequest creates CreateTenantThresholdProfileBadRequest with default headers values
func NewCreateTenantThresholdProfileBadRequest() *CreateTenantThresholdProfileBadRequest {

	return &CreateTenantThresholdProfileBadRequest{}
}

// WithPayload adds the payload to the create tenant threshold profile bad request response
func (o *CreateTenantThresholdProfileBadRequest) WithPayload(payload string) *CreateTenantThresholdProfileBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant threshold profile bad request response
func (o *CreateTenantThresholdProfileBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantThresholdProfileBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantThresholdProfileForbiddenCode is the HTTP code returned for type CreateTenantThresholdProfileForbidden
const CreateTenantThresholdProfileForbiddenCode int = 403

/*CreateTenantThresholdProfileForbidden Requestor does not have authorization to perform this action

swagger:response createTenantThresholdProfileForbidden
*/
type CreateTenantThresholdProfileForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantThresholdProfileForbidden creates CreateTenantThresholdProfileForbidden with default headers values
func NewCreateTenantThresholdProfileForbidden() *CreateTenantThresholdProfileForbidden {

	return &CreateTenantThresholdProfileForbidden{}
}

// WithPayload adds the payload to the create tenant threshold profile forbidden response
func (o *CreateTenantThresholdProfileForbidden) WithPayload(payload string) *CreateTenantThresholdProfileForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant threshold profile forbidden response
func (o *CreateTenantThresholdProfileForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantThresholdProfileForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantThresholdProfileInternalServerErrorCode is the HTTP code returned for type CreateTenantThresholdProfileInternalServerError
const CreateTenantThresholdProfileInternalServerErrorCode int = 500

/*CreateTenantThresholdProfileInternalServerError Unexpected error processing request

swagger:response createTenantThresholdProfileInternalServerError
*/
type CreateTenantThresholdProfileInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantThresholdProfileInternalServerError creates CreateTenantThresholdProfileInternalServerError with default headers values
func NewCreateTenantThresholdProfileInternalServerError() *CreateTenantThresholdProfileInternalServerError {

	return &CreateTenantThresholdProfileInternalServerError{}
}

// WithPayload adds the payload to the create tenant threshold profile internal server error response
func (o *CreateTenantThresholdProfileInternalServerError) WithPayload(payload string) *CreateTenantThresholdProfileInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant threshold profile internal server error response
func (o *CreateTenantThresholdProfileInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantThresholdProfileInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
