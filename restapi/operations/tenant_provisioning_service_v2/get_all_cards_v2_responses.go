// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetAllCardsV2OKCode is the HTTP code returned for type GetAllCardsV2OK
const GetAllCardsV2OKCode int = 200

/*GetAllCardsV2OK get all cards v2 o k

swagger:response getAllCardsV2OK
*/
type GetAllCardsV2OK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.CardListResponse `json:"body,omitempty"`
}

// NewGetAllCardsV2OK creates GetAllCardsV2OK with default headers values
func NewGetAllCardsV2OK() *GetAllCardsV2OK {

	return &GetAllCardsV2OK{}
}

// WithPayload adds the payload to the get all cards v2 o k response
func (o *GetAllCardsV2OK) WithPayload(payload *swagmodels.CardListResponse) *GetAllCardsV2OK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all cards v2 o k response
func (o *GetAllCardsV2OK) SetPayload(payload *swagmodels.CardListResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCardsV2OK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAllCardsV2ForbiddenCode is the HTTP code returned for type GetAllCardsV2Forbidden
const GetAllCardsV2ForbiddenCode int = 403

/*GetAllCardsV2Forbidden Requestor does not have authorization to perform this action

swagger:response getAllCardsV2Forbidden
*/
type GetAllCardsV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAllCardsV2Forbidden creates GetAllCardsV2Forbidden with default headers values
func NewGetAllCardsV2Forbidden() *GetAllCardsV2Forbidden {

	return &GetAllCardsV2Forbidden{}
}

// WithPayload adds the payload to the get all cards v2 forbidden response
func (o *GetAllCardsV2Forbidden) WithPayload(payload string) *GetAllCardsV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all cards v2 forbidden response
func (o *GetAllCardsV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCardsV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetAllCardsV2NotFoundCode is the HTTP code returned for type GetAllCardsV2NotFound
const GetAllCardsV2NotFoundCode int = 404

/*GetAllCardsV2NotFound No Cardss are provisioned

swagger:response getAllCardsV2NotFound
*/
type GetAllCardsV2NotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAllCardsV2NotFound creates GetAllCardsV2NotFound with default headers values
func NewGetAllCardsV2NotFound() *GetAllCardsV2NotFound {

	return &GetAllCardsV2NotFound{}
}

// WithPayload adds the payload to the get all cards v2 not found response
func (o *GetAllCardsV2NotFound) WithPayload(payload string) *GetAllCardsV2NotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all cards v2 not found response
func (o *GetAllCardsV2NotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCardsV2NotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetAllCardsV2InternalServerErrorCode is the HTTP code returned for type GetAllCardsV2InternalServerError
const GetAllCardsV2InternalServerErrorCode int = 500

/*GetAllCardsV2InternalServerError Unexpected error processing request

swagger:response getAllCardsV2InternalServerError
*/
type GetAllCardsV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetAllCardsV2InternalServerError creates GetAllCardsV2InternalServerError with default headers values
func NewGetAllCardsV2InternalServerError() *GetAllCardsV2InternalServerError {

	return &GetAllCardsV2InternalServerError{}
}

// WithPayload adds the payload to the get all cards v2 internal server error response
func (o *GetAllCardsV2InternalServerError) WithPayload(payload string) *GetAllCardsV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all cards v2 internal server error response
func (o *GetAllCardsV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCardsV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
