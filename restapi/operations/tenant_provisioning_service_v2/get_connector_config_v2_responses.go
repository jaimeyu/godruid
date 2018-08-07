// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetConnectorConfigV2OKCode is the HTTP code returned for type GetConnectorConfigV2OK
const GetConnectorConfigV2OKCode int = 200

/*GetConnectorConfigV2OK get connector config v2 o k

swagger:response getConnectorConfigV2OK
*/
type GetConnectorConfigV2OK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.ConnectorConfigResponse `json:"body,omitempty"`
}

// NewGetConnectorConfigV2OK creates GetConnectorConfigV2OK with default headers values
func NewGetConnectorConfigV2OK() *GetConnectorConfigV2OK {

	return &GetConnectorConfigV2OK{}
}

// WithPayload adds the payload to the get connector config v2 o k response
func (o *GetConnectorConfigV2OK) WithPayload(payload *swagmodels.ConnectorConfigResponse) *GetConnectorConfigV2OK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get connector config v2 o k response
func (o *GetConnectorConfigV2OK) SetPayload(payload *swagmodels.ConnectorConfigResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConnectorConfigV2OK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetConnectorConfigV2ForbiddenCode is the HTTP code returned for type GetConnectorConfigV2Forbidden
const GetConnectorConfigV2ForbiddenCode int = 403

/*GetConnectorConfigV2Forbidden Requestor does not have authorization to perform this action

swagger:response getConnectorConfigV2Forbidden
*/
type GetConnectorConfigV2Forbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetConnectorConfigV2Forbidden creates GetConnectorConfigV2Forbidden with default headers values
func NewGetConnectorConfigV2Forbidden() *GetConnectorConfigV2Forbidden {

	return &GetConnectorConfigV2Forbidden{}
}

// WithPayload adds the payload to the get connector config v2 forbidden response
func (o *GetConnectorConfigV2Forbidden) WithPayload(payload string) *GetConnectorConfigV2Forbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get connector config v2 forbidden response
func (o *GetConnectorConfigV2Forbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConnectorConfigV2Forbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetConnectorConfigV2NotFoundCode is the HTTP code returned for type GetConnectorConfigV2NotFound
const GetConnectorConfigV2NotFoundCode int = 404

/*GetConnectorConfigV2NotFound The specified connector configuration is not provisioned

swagger:response getConnectorConfigV2NotFound
*/
type GetConnectorConfigV2NotFound struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetConnectorConfigV2NotFound creates GetConnectorConfigV2NotFound with default headers values
func NewGetConnectorConfigV2NotFound() *GetConnectorConfigV2NotFound {

	return &GetConnectorConfigV2NotFound{}
}

// WithPayload adds the payload to the get connector config v2 not found response
func (o *GetConnectorConfigV2NotFound) WithPayload(payload string) *GetConnectorConfigV2NotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get connector config v2 not found response
func (o *GetConnectorConfigV2NotFound) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConnectorConfigV2NotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetConnectorConfigV2InternalServerErrorCode is the HTTP code returned for type GetConnectorConfigV2InternalServerError
const GetConnectorConfigV2InternalServerErrorCode int = 500

/*GetConnectorConfigV2InternalServerError Unexpected error processing request

swagger:response getConnectorConfigV2InternalServerError
*/
type GetConnectorConfigV2InternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetConnectorConfigV2InternalServerError creates GetConnectorConfigV2InternalServerError with default headers values
func NewGetConnectorConfigV2InternalServerError() *GetConnectorConfigV2InternalServerError {

	return &GetConnectorConfigV2InternalServerError{}
}

// WithPayload adds the payload to the get connector config v2 internal server error response
func (o *GetConnectorConfigV2InternalServerError) WithPayload(payload string) *GetConnectorConfigV2InternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get connector config v2 internal server error response
func (o *GetConnectorConfigV2InternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConnectorConfigV2InternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
