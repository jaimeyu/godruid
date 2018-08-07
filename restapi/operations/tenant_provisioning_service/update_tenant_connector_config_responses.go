// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// UpdateTenantConnectorConfigOKCode is the HTTP code returned for type UpdateTenantConnectorConfigOK
const UpdateTenantConnectorConfigOKCode int = 200

/*UpdateTenantConnectorConfigOK update tenant connector config o k

swagger:response updateTenantConnectorConfigOK
*/
type UpdateTenantConnectorConfigOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantConnectorConfig `json:"body,omitempty"`
}

// NewUpdateTenantConnectorConfigOK creates UpdateTenantConnectorConfigOK with default headers values
func NewUpdateTenantConnectorConfigOK() *UpdateTenantConnectorConfigOK {

	return &UpdateTenantConnectorConfigOK{}
}

// WithPayload adds the payload to the update tenant connector config o k response
func (o *UpdateTenantConnectorConfigOK) WithPayload(payload *swagmodels.JSONAPITenantConnectorConfig) *UpdateTenantConnectorConfigOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update tenant connector config o k response
func (o *UpdateTenantConnectorConfigOK) SetPayload(payload *swagmodels.JSONAPITenantConnectorConfig) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTenantConnectorConfigOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTenantConnectorConfigBadRequestCode is the HTTP code returned for type UpdateTenantConnectorConfigBadRequest
const UpdateTenantConnectorConfigBadRequestCode int = 400

/*UpdateTenantConnectorConfigBadRequest Request data does not pass validation

swagger:response updateTenantConnectorConfigBadRequest
*/
type UpdateTenantConnectorConfigBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateTenantConnectorConfigBadRequest creates UpdateTenantConnectorConfigBadRequest with default headers values
func NewUpdateTenantConnectorConfigBadRequest() *UpdateTenantConnectorConfigBadRequest {

	return &UpdateTenantConnectorConfigBadRequest{}
}

// WithPayload adds the payload to the update tenant connector config bad request response
func (o *UpdateTenantConnectorConfigBadRequest) WithPayload(payload string) *UpdateTenantConnectorConfigBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update tenant connector config bad request response
func (o *UpdateTenantConnectorConfigBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTenantConnectorConfigBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateTenantConnectorConfigForbiddenCode is the HTTP code returned for type UpdateTenantConnectorConfigForbidden
const UpdateTenantConnectorConfigForbiddenCode int = 403

/*UpdateTenantConnectorConfigForbidden Requestor does not have authorization to perform this action

swagger:response updateTenantConnectorConfigForbidden
*/
type UpdateTenantConnectorConfigForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateTenantConnectorConfigForbidden creates UpdateTenantConnectorConfigForbidden with default headers values
func NewUpdateTenantConnectorConfigForbidden() *UpdateTenantConnectorConfigForbidden {

	return &UpdateTenantConnectorConfigForbidden{}
}

// WithPayload adds the payload to the update tenant connector config forbidden response
func (o *UpdateTenantConnectorConfigForbidden) WithPayload(payload string) *UpdateTenantConnectorConfigForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update tenant connector config forbidden response
func (o *UpdateTenantConnectorConfigForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTenantConnectorConfigForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateTenantConnectorConfigInternalServerErrorCode is the HTTP code returned for type UpdateTenantConnectorConfigInternalServerError
const UpdateTenantConnectorConfigInternalServerErrorCode int = 500

/*UpdateTenantConnectorConfigInternalServerError Unexpected error processing request

swagger:response updateTenantConnectorConfigInternalServerError
*/
type UpdateTenantConnectorConfigInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewUpdateTenantConnectorConfigInternalServerError creates UpdateTenantConnectorConfigInternalServerError with default headers values
func NewUpdateTenantConnectorConfigInternalServerError() *UpdateTenantConnectorConfigInternalServerError {

	return &UpdateTenantConnectorConfigInternalServerError{}
}

// WithPayload adds the payload to the update tenant connector config internal server error response
func (o *UpdateTenantConnectorConfigInternalServerError) WithPayload(payload string) *UpdateTenantConnectorConfigInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update tenant connector config internal server error response
func (o *UpdateTenantConnectorConfigInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTenantConnectorConfigInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
