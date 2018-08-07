// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// CreateThresholdProfileV2CreatedCode is the HTTP code returned for type CreateThresholdProfileV2Created
const CreateThresholdProfileV2CreatedCode int = 201

/*CreateThresholdProfileV2Created create threshold profile v2 created

swagger:response createThresholdProfileV2Created
*/
type CreateThresholdProfileV2Created struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.ThresholdProfileResponse `json:"body,omitempty"`
}

// NewCreateThresholdProfileV2Created creates CreateThresholdProfileV2Created with default headers values
func NewCreateThresholdProfileV2Created() *CreateThresholdProfileV2Created {

	return &CreateThresholdProfileV2Created{}
}

// WithPayload adds the payload to the create threshold profile v2 created response
func (o *CreateThresholdProfileV2Created) WithPayload(payload *swagmodels.ThresholdProfileResponse) *CreateThresholdProfileV2Created {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create threshold profile v2 created response
func (o *CreateThresholdProfileV2Created) SetPayload(payload *swagmodels.ThresholdProfileResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateThresholdProfileV2Created) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateThresholdProfileV2BadRequestCode is the HTTP code returned for type CreateThresholdProfileV2BadRequest
const CreateThresholdProfileV2BadRequestCode int = 400

/*CreateThresholdProfileV2BadRequest Request data does not pass validation

swagger:response createThresholdProfileV2BadRequest
*/
type CreateThresholdProfileV2BadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateThresholdProfileV2BadRequest creates CreateThresholdProfileV2BadRequest with default headers values
func NewCreateThresholdProfileV2BadRequest() *CreateThresholdProfileV2BadRequest {

	return &CreateThresholdProfileV2BadRequest{}
}

// WithPayload adds the payload to the create threshold profile v2 bad request response
func (o *CreateThresholdProfileV2BadRequest) WithPayload(payload string) *CreateThresholdProfileV2BadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create threshold profile v2 bad request response
func (o *CreateThresholdProfileV2BadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateThresholdProfileV2BadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateThresholdProfileV2ForbiddenCode is the HTTP code returned for type CreateThresholdProfileV2Forbidden
const CreateThresholdProfileV2ForbiddenCode int = 403

/*CreateThresholdProfileV2Forbidden Requestor does not have authorization to perform this action

swagger:response createThresholdProfileV2Forbidden
*/
type CreateThresholdProfileV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateThresholdProfileV2Forbidden creates CreateThresholdProfileV2Forbidden with default headers values
func NewCreateThresholdProfileV2Forbidden() *CreateThresholdProfileV2Forbidden {

	return &CreateThresholdProfileV2Forbidden{}
}

// WithPayload adds the payload to the create threshold profile v2 forbidden response
func (o *CreateThresholdProfileV2Forbidden) WithPayload(payload string) *CreateThresholdProfileV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create threshold profile v2 forbidden response
func (o *CreateThresholdProfileV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateThresholdProfileV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateThresholdProfileV2ConflictCode is the HTTP code returned for type CreateThresholdProfileV2Conflict
const CreateThresholdProfileV2ConflictCode int = 409

/*CreateThresholdProfileV2Conflict The record is already provisioned

swagger:response createThresholdProfileV2Conflict
*/
type CreateThresholdProfileV2Conflict struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateThresholdProfileV2Conflict creates CreateThresholdProfileV2Conflict with default headers values
func NewCreateThresholdProfileV2Conflict() *CreateThresholdProfileV2Conflict {

	return &CreateThresholdProfileV2Conflict{}
}

// WithPayload adds the payload to the create threshold profile v2 conflict response
func (o *CreateThresholdProfileV2Conflict) WithPayload(payload string) *CreateThresholdProfileV2Conflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create threshold profile v2 conflict response
func (o *CreateThresholdProfileV2Conflict) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateThresholdProfileV2Conflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateThresholdProfileV2InternalServerErrorCode is the HTTP code returned for type CreateThresholdProfileV2InternalServerError
const CreateThresholdProfileV2InternalServerErrorCode int = 500

/*CreateThresholdProfileV2InternalServerError Unexpected error processing request

swagger:response createThresholdProfileV2InternalServerError
*/
type CreateThresholdProfileV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateThresholdProfileV2InternalServerError creates CreateThresholdProfileV2InternalServerError with default headers values
func NewCreateThresholdProfileV2InternalServerError() *CreateThresholdProfileV2InternalServerError {

	return &CreateThresholdProfileV2InternalServerError{}
}

// WithPayload adds the payload to the create threshold profile v2 internal server error response
func (o *CreateThresholdProfileV2InternalServerError) WithPayload(payload string) *CreateThresholdProfileV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create threshold profile v2 internal server error response
func (o *CreateThresholdProfileV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateThresholdProfileV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
