// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// UpdateIngestionProfileV2OKCode is the HTTP code returned for type UpdateIngestionProfileV2OK
const UpdateIngestionProfileV2OKCode int = 200

/*UpdateIngestionProfileV2OK update ingestion profile v2 o k

swagger:response updateIngestionProfileV2OK
*/
type UpdateIngestionProfileV2OK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.IngestionProfileResponse `json:"body,omitempty"`
}

// NewUpdateIngestionProfileV2OK creates UpdateIngestionProfileV2OK with default headers values
func NewUpdateIngestionProfileV2OK() *UpdateIngestionProfileV2OK {

	return &UpdateIngestionProfileV2OK{}
}

// WithPayload adds the payload to the update ingestion profile v2 o k response
func (o *UpdateIngestionProfileV2OK) WithPayload(payload *swagmodels.IngestionProfileResponse) *UpdateIngestionProfileV2OK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update ingestion profile v2 o k response
func (o *UpdateIngestionProfileV2OK) SetPayload(payload *swagmodels.IngestionProfileResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateIngestionProfileV2OK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateIngestionProfileV2BadRequestCode is the HTTP code returned for type UpdateIngestionProfileV2BadRequest
const UpdateIngestionProfileV2BadRequestCode int = 400

/*UpdateIngestionProfileV2BadRequest Request data does not pass validation

swagger:response updateIngestionProfileV2BadRequest
*/
type UpdateIngestionProfileV2BadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateIngestionProfileV2BadRequest creates UpdateIngestionProfileV2BadRequest with default headers values
func NewUpdateIngestionProfileV2BadRequest() *UpdateIngestionProfileV2BadRequest {

	return &UpdateIngestionProfileV2BadRequest{}
}

// WithPayload adds the payload to the update ingestion profile v2 bad request response
func (o *UpdateIngestionProfileV2BadRequest) WithPayload(payload string) *UpdateIngestionProfileV2BadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update ingestion profile v2 bad request response
func (o *UpdateIngestionProfileV2BadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateIngestionProfileV2BadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateIngestionProfileV2ForbiddenCode is the HTTP code returned for type UpdateIngestionProfileV2Forbidden
const UpdateIngestionProfileV2ForbiddenCode int = 403

/*UpdateIngestionProfileV2Forbidden Requestor does not have authorization to perform this action

swagger:response updateIngestionProfileV2Forbidden
*/
type UpdateIngestionProfileV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateIngestionProfileV2Forbidden creates UpdateIngestionProfileV2Forbidden with default headers values
func NewUpdateIngestionProfileV2Forbidden() *UpdateIngestionProfileV2Forbidden {

	return &UpdateIngestionProfileV2Forbidden{}
}

// WithPayload adds the payload to the update ingestion profile v2 forbidden response
func (o *UpdateIngestionProfileV2Forbidden) WithPayload(payload string) *UpdateIngestionProfileV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update ingestion profile v2 forbidden response
func (o *UpdateIngestionProfileV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateIngestionProfileV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateIngestionProfileV2NotFoundCode is the HTTP code returned for type UpdateIngestionProfileV2NotFound
const UpdateIngestionProfileV2NotFoundCode int = 404

/*UpdateIngestionProfileV2NotFound The specified Ingestion Profile is not provisioned

swagger:response updateIngestionProfileV2NotFound
*/
type UpdateIngestionProfileV2NotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateIngestionProfileV2NotFound creates UpdateIngestionProfileV2NotFound with default headers values
func NewUpdateIngestionProfileV2NotFound() *UpdateIngestionProfileV2NotFound {

	return &UpdateIngestionProfileV2NotFound{}
}

// WithPayload adds the payload to the update ingestion profile v2 not found response
func (o *UpdateIngestionProfileV2NotFound) WithPayload(payload string) *UpdateIngestionProfileV2NotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update ingestion profile v2 not found response
func (o *UpdateIngestionProfileV2NotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateIngestionProfileV2NotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateIngestionProfileV2ConflictCode is the HTTP code returned for type UpdateIngestionProfileV2Conflict
const UpdateIngestionProfileV2ConflictCode int = 409

/*UpdateIngestionProfileV2Conflict Incorrect revision provided for th update request

swagger:response updateIngestionProfileV2Conflict
*/
type UpdateIngestionProfileV2Conflict struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateIngestionProfileV2Conflict creates UpdateIngestionProfileV2Conflict with default headers values
func NewUpdateIngestionProfileV2Conflict() *UpdateIngestionProfileV2Conflict {

	return &UpdateIngestionProfileV2Conflict{}
}

// WithPayload adds the payload to the update ingestion profile v2 conflict response
func (o *UpdateIngestionProfileV2Conflict) WithPayload(payload string) *UpdateIngestionProfileV2Conflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update ingestion profile v2 conflict response
func (o *UpdateIngestionProfileV2Conflict) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateIngestionProfileV2Conflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateIngestionProfileV2InternalServerErrorCode is the HTTP code returned for type UpdateIngestionProfileV2InternalServerError
const UpdateIngestionProfileV2InternalServerErrorCode int = 500

/*UpdateIngestionProfileV2InternalServerError Unexpected error processing request

swagger:response updateIngestionProfileV2InternalServerError
*/
type UpdateIngestionProfileV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateIngestionProfileV2InternalServerError creates UpdateIngestionProfileV2InternalServerError with default headers values
func NewUpdateIngestionProfileV2InternalServerError() *UpdateIngestionProfileV2InternalServerError {

	return &UpdateIngestionProfileV2InternalServerError{}
}

// WithPayload adds the payload to the update ingestion profile v2 internal server error response
func (o *UpdateIngestionProfileV2InternalServerError) WithPayload(payload string) *UpdateIngestionProfileV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update ingestion profile v2 internal server error response
func (o *UpdateIngestionProfileV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateIngestionProfileV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
