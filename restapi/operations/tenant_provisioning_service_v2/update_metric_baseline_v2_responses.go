// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// UpdateMetricBaselineV2OKCode is the HTTP code returned for type UpdateMetricBaselineV2OK
const UpdateMetricBaselineV2OKCode int = 200

/*UpdateMetricBaselineV2OK update metric baseline v2 o k

swagger:response updateMetricBaselineV2OK
*/
type UpdateMetricBaselineV2OK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.MetricBaselineResponse `json:"body,omitempty"`
}

// NewUpdateMetricBaselineV2OK creates UpdateMetricBaselineV2OK with default headers values
func NewUpdateMetricBaselineV2OK() *UpdateMetricBaselineV2OK {

	return &UpdateMetricBaselineV2OK{}
}

// WithPayload adds the payload to the update metric baseline v2 o k response
func (o *UpdateMetricBaselineV2OK) WithPayload(payload *swagmodels.MetricBaselineResponse) *UpdateMetricBaselineV2OK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update metric baseline v2 o k response
func (o *UpdateMetricBaselineV2OK) SetPayload(payload *swagmodels.MetricBaselineResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMetricBaselineV2OK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMetricBaselineV2BadRequestCode is the HTTP code returned for type UpdateMetricBaselineV2BadRequest
const UpdateMetricBaselineV2BadRequestCode int = 400

/*UpdateMetricBaselineV2BadRequest Request data does not pass validation

swagger:response updateMetricBaselineV2BadRequest
*/
type UpdateMetricBaselineV2BadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateMetricBaselineV2BadRequest creates UpdateMetricBaselineV2BadRequest with default headers values
func NewUpdateMetricBaselineV2BadRequest() *UpdateMetricBaselineV2BadRequest {

	return &UpdateMetricBaselineV2BadRequest{}
}

// WithPayload adds the payload to the update metric baseline v2 bad request response
func (o *UpdateMetricBaselineV2BadRequest) WithPayload(payload string) *UpdateMetricBaselineV2BadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update metric baseline v2 bad request response
func (o *UpdateMetricBaselineV2BadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMetricBaselineV2BadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateMetricBaselineV2ForbiddenCode is the HTTP code returned for type UpdateMetricBaselineV2Forbidden
const UpdateMetricBaselineV2ForbiddenCode int = 403

/*UpdateMetricBaselineV2Forbidden Requestor does not have authorization to perform this action

swagger:response updateMetricBaselineV2Forbidden
*/
type UpdateMetricBaselineV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateMetricBaselineV2Forbidden creates UpdateMetricBaselineV2Forbidden with default headers values
func NewUpdateMetricBaselineV2Forbidden() *UpdateMetricBaselineV2Forbidden {

	return &UpdateMetricBaselineV2Forbidden{}
}

// WithPayload adds the payload to the update metric baseline v2 forbidden response
func (o *UpdateMetricBaselineV2Forbidden) WithPayload(payload string) *UpdateMetricBaselineV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update metric baseline v2 forbidden response
func (o *UpdateMetricBaselineV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMetricBaselineV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateMetricBaselineV2NotFoundCode is the HTTP code returned for type UpdateMetricBaselineV2NotFound
const UpdateMetricBaselineV2NotFoundCode int = 404

/*UpdateMetricBaselineV2NotFound The specified Metric Baseline is not provisioned

swagger:response updateMetricBaselineV2NotFound
*/
type UpdateMetricBaselineV2NotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateMetricBaselineV2NotFound creates UpdateMetricBaselineV2NotFound with default headers values
func NewUpdateMetricBaselineV2NotFound() *UpdateMetricBaselineV2NotFound {

	return &UpdateMetricBaselineV2NotFound{}
}

// WithPayload adds the payload to the update metric baseline v2 not found response
func (o *UpdateMetricBaselineV2NotFound) WithPayload(payload string) *UpdateMetricBaselineV2NotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update metric baseline v2 not found response
func (o *UpdateMetricBaselineV2NotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMetricBaselineV2NotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateMetricBaselineV2ConflictCode is the HTTP code returned for type UpdateMetricBaselineV2Conflict
const UpdateMetricBaselineV2ConflictCode int = 409

/*UpdateMetricBaselineV2Conflict Incorrect revision provided for th update request

swagger:response updateMetricBaselineV2Conflict
*/
type UpdateMetricBaselineV2Conflict struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateMetricBaselineV2Conflict creates UpdateMetricBaselineV2Conflict with default headers values
func NewUpdateMetricBaselineV2Conflict() *UpdateMetricBaselineV2Conflict {

	return &UpdateMetricBaselineV2Conflict{}
}

// WithPayload adds the payload to the update metric baseline v2 conflict response
func (o *UpdateMetricBaselineV2Conflict) WithPayload(payload string) *UpdateMetricBaselineV2Conflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update metric baseline v2 conflict response
func (o *UpdateMetricBaselineV2Conflict) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMetricBaselineV2Conflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateMetricBaselineV2InternalServerErrorCode is the HTTP code returned for type UpdateMetricBaselineV2InternalServerError
const UpdateMetricBaselineV2InternalServerErrorCode int = 500

/*UpdateMetricBaselineV2InternalServerError Unexpected error processing request

swagger:response updateMetricBaselineV2InternalServerError
*/
type UpdateMetricBaselineV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateMetricBaselineV2InternalServerError creates UpdateMetricBaselineV2InternalServerError with default headers values
func NewUpdateMetricBaselineV2InternalServerError() *UpdateMetricBaselineV2InternalServerError {

	return &UpdateMetricBaselineV2InternalServerError{}
}

// WithPayload adds the payload to the update metric baseline v2 internal server error response
func (o *UpdateMetricBaselineV2InternalServerError) WithPayload(payload string) *UpdateMetricBaselineV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update metric baseline v2 internal server error response
func (o *UpdateMetricBaselineV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMetricBaselineV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
