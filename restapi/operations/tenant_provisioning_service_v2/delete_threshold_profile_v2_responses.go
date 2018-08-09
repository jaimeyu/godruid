// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// DeleteThresholdProfileV2OKCode is the HTTP code returned for type DeleteThresholdProfileV2OK
const DeleteThresholdProfileV2OKCode int = 200

/*DeleteThresholdProfileV2OK delete threshold profile v2 o k

swagger:response deleteThresholdProfileV2OK
*/
type DeleteThresholdProfileV2OK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.ThresholdProfileResponse `json:"body,omitempty"`
}

// NewDeleteThresholdProfileV2OK creates DeleteThresholdProfileV2OK with default headers values
func NewDeleteThresholdProfileV2OK() *DeleteThresholdProfileV2OK {

	return &DeleteThresholdProfileV2OK{}
}

// WithPayload adds the payload to the delete threshold profile v2 o k response
func (o *DeleteThresholdProfileV2OK) WithPayload(payload *swagmodels.ThresholdProfileResponse) *DeleteThresholdProfileV2OK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete threshold profile v2 o k response
func (o *DeleteThresholdProfileV2OK) SetPayload(payload *swagmodels.ThresholdProfileResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteThresholdProfileV2OK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteThresholdProfileV2ForbiddenCode is the HTTP code returned for type DeleteThresholdProfileV2Forbidden
const DeleteThresholdProfileV2ForbiddenCode int = 403

/*DeleteThresholdProfileV2Forbidden Requestor does not have authorization to perform this action

swagger:response deleteThresholdProfileV2Forbidden
*/
type DeleteThresholdProfileV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewDeleteThresholdProfileV2Forbidden creates DeleteThresholdProfileV2Forbidden with default headers values
func NewDeleteThresholdProfileV2Forbidden() *DeleteThresholdProfileV2Forbidden {

	return &DeleteThresholdProfileV2Forbidden{}
}

// WithPayload adds the payload to the delete threshold profile v2 forbidden response
func (o *DeleteThresholdProfileV2Forbidden) WithPayload(payload string) *DeleteThresholdProfileV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete threshold profile v2 forbidden response
func (o *DeleteThresholdProfileV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteThresholdProfileV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// DeleteThresholdProfileV2NotFoundCode is the HTTP code returned for type DeleteThresholdProfileV2NotFound
const DeleteThresholdProfileV2NotFoundCode int = 404

/*DeleteThresholdProfileV2NotFound The requested Threshold Profile is not provisioned

swagger:response deleteThresholdProfileV2NotFound
*/
type DeleteThresholdProfileV2NotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewDeleteThresholdProfileV2NotFound creates DeleteThresholdProfileV2NotFound with default headers values
func NewDeleteThresholdProfileV2NotFound() *DeleteThresholdProfileV2NotFound {

	return &DeleteThresholdProfileV2NotFound{}
}

// WithPayload adds the payload to the delete threshold profile v2 not found response
func (o *DeleteThresholdProfileV2NotFound) WithPayload(payload string) *DeleteThresholdProfileV2NotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete threshold profile v2 not found response
func (o *DeleteThresholdProfileV2NotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteThresholdProfileV2NotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// DeleteThresholdProfileV2InternalServerErrorCode is the HTTP code returned for type DeleteThresholdProfileV2InternalServerError
const DeleteThresholdProfileV2InternalServerErrorCode int = 500

/*DeleteThresholdProfileV2InternalServerError Unexpected error processing request

swagger:response deleteThresholdProfileV2InternalServerError
*/
type DeleteThresholdProfileV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewDeleteThresholdProfileV2InternalServerError creates DeleteThresholdProfileV2InternalServerError with default headers values
func NewDeleteThresholdProfileV2InternalServerError() *DeleteThresholdProfileV2InternalServerError {

	return &DeleteThresholdProfileV2InternalServerError{}
}

// WithPayload adds the payload to the delete threshold profile v2 internal server error response
func (o *DeleteThresholdProfileV2InternalServerError) WithPayload(payload string) *DeleteThresholdProfileV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete threshold profile v2 internal server error response
func (o *DeleteThresholdProfileV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteThresholdProfileV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
