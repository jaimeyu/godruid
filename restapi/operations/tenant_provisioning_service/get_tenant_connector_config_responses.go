// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetTenantConnectorConfigOKCode is the HTTP code returned for type GetTenantConnectorConfigOK
const GetTenantConnectorConfigOKCode int = 200

/*GetTenantConnectorConfigOK get tenant connector config o k

swagger:response getTenantConnectorConfigOK
*/
type GetTenantConnectorConfigOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantConnectorConfig `json:"body,omitempty"`
}

// NewGetTenantConnectorConfigOK creates GetTenantConnectorConfigOK with default headers values
func NewGetTenantConnectorConfigOK() *GetTenantConnectorConfigOK {

	return &GetTenantConnectorConfigOK{}
}

// WithPayload adds the payload to the get tenant connector config o k response
func (o *GetTenantConnectorConfigOK) WithPayload(payload *swagmodels.JSONAPITenantConnectorConfig) *GetTenantConnectorConfigOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tenant connector config o k response
func (o *GetTenantConnectorConfigOK) SetPayload(payload *swagmodels.JSONAPITenantConnectorConfig) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTenantConnectorConfigOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetTenantConnectorConfigInternalServerErrorCode is the HTTP code returned for type GetTenantConnectorConfigInternalServerError
const GetTenantConnectorConfigInternalServerErrorCode int = 500

/*GetTenantConnectorConfigInternalServerError Unexpected error processing request

swagger:response getTenantConnectorConfigInternalServerError
*/
type GetTenantConnectorConfigInternalServerError struct {
}

// NewGetTenantConnectorConfigInternalServerError creates GetTenantConnectorConfigInternalServerError with default headers values
func NewGetTenantConnectorConfigInternalServerError() *GetTenantConnectorConfigInternalServerError {

	return &GetTenantConnectorConfigInternalServerError{}
}

// WriteResponse to the client
func (o *GetTenantConnectorConfigInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
