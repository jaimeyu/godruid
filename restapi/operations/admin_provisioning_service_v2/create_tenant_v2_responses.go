// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// CreateTenantV2CreatedCode is the HTTP code returned for type CreateTenantV2Created
const CreateTenantV2CreatedCode int = 201

/*CreateTenantV2Created create tenant v2 created

swagger:response createTenantV2Created
*/
type CreateTenantV2Created struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.TenantResponse `json:"body,omitempty"`
}

// NewCreateTenantV2Created creates CreateTenantV2Created with default headers values
func NewCreateTenantV2Created() *CreateTenantV2Created {

	return &CreateTenantV2Created{}
}

// WithPayload adds the payload to the create tenant v2 created response
func (o *CreateTenantV2Created) WithPayload(payload *swagmodels.TenantResponse) *CreateTenantV2Created {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant v2 created response
func (o *CreateTenantV2Created) SetPayload(payload *swagmodels.TenantResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantV2Created) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTenantV2BadRequestCode is the HTTP code returned for type CreateTenantV2BadRequest
const CreateTenantV2BadRequestCode int = 400

/*CreateTenantV2BadRequest Request data does not pass validation

swagger:response createTenantV2BadRequest
*/
type CreateTenantV2BadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantV2BadRequest creates CreateTenantV2BadRequest with default headers values
func NewCreateTenantV2BadRequest() *CreateTenantV2BadRequest {

	return &CreateTenantV2BadRequest{}
}

// WithPayload adds the payload to the create tenant v2 bad request response
func (o *CreateTenantV2BadRequest) WithPayload(payload string) *CreateTenantV2BadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant v2 bad request response
func (o *CreateTenantV2BadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantV2BadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantV2ForbiddenCode is the HTTP code returned for type CreateTenantV2Forbidden
const CreateTenantV2ForbiddenCode int = 403

/*CreateTenantV2Forbidden Requestor does not have authorization to perform this action

swagger:response createTenantV2Forbidden
*/
type CreateTenantV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantV2Forbidden creates CreateTenantV2Forbidden with default headers values
func NewCreateTenantV2Forbidden() *CreateTenantV2Forbidden {

	return &CreateTenantV2Forbidden{}
}

// WithPayload adds the payload to the create tenant v2 forbidden response
func (o *CreateTenantV2Forbidden) WithPayload(payload string) *CreateTenantV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant v2 forbidden response
func (o *CreateTenantV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantV2ConflictCode is the HTTP code returned for type CreateTenantV2Conflict
const CreateTenantV2ConflictCode int = 409

/*CreateTenantV2Conflict The Tenant being provisioned already exists

swagger:response createTenantV2Conflict
*/
type CreateTenantV2Conflict struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantV2Conflict creates CreateTenantV2Conflict with default headers values
func NewCreateTenantV2Conflict() *CreateTenantV2Conflict {

	return &CreateTenantV2Conflict{}
}

// WithPayload adds the payload to the create tenant v2 conflict response
func (o *CreateTenantV2Conflict) WithPayload(payload string) *CreateTenantV2Conflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant v2 conflict response
func (o *CreateTenantV2Conflict) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantV2Conflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantV2InternalServerErrorCode is the HTTP code returned for type CreateTenantV2InternalServerError
const CreateTenantV2InternalServerErrorCode int = 500

/*CreateTenantV2InternalServerError Unexpected error processing request

swagger:response createTenantV2InternalServerError
*/
type CreateTenantV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantV2InternalServerError creates CreateTenantV2InternalServerError with default headers values
func NewCreateTenantV2InternalServerError() *CreateTenantV2InternalServerError {

	return &CreateTenantV2InternalServerError{}
}

// WithPayload adds the payload to the create tenant v2 internal server error response
func (o *CreateTenantV2InternalServerError) WithPayload(payload string) *CreateTenantV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant v2 internal server error response
func (o *CreateTenantV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
