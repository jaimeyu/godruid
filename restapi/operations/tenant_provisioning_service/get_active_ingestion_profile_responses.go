// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetActiveIngestionProfileOKCode is the HTTP code returned for type GetActiveIngestionProfileOK
const GetActiveIngestionProfileOKCode int = 200

/*GetActiveIngestionProfileOK get active ingestion profile o k

swagger:response getActiveIngestionProfileOK
*/
type GetActiveIngestionProfileOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantIngestionProfile `json:"body,omitempty"`
}

// NewGetActiveIngestionProfileOK creates GetActiveIngestionProfileOK with default headers values
func NewGetActiveIngestionProfileOK() *GetActiveIngestionProfileOK {

	return &GetActiveIngestionProfileOK{}
}

// WithPayload adds the payload to the get active ingestion profile o k response
func (o *GetActiveIngestionProfileOK) WithPayload(payload *swagmodels.JSONAPITenantIngestionProfile) *GetActiveIngestionProfileOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get active ingestion profile o k response
func (o *GetActiveIngestionProfileOK) SetPayload(payload *swagmodels.JSONAPITenantIngestionProfile) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActiveIngestionProfileOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetActiveIngestionProfileInternalServerErrorCode is the HTTP code returned for type GetActiveIngestionProfileInternalServerError
const GetActiveIngestionProfileInternalServerErrorCode int = 500

/*GetActiveIngestionProfileInternalServerError Unexpected error processing request

swagger:response getActiveIngestionProfileInternalServerError
*/
type GetActiveIngestionProfileInternalServerError struct {
}

// NewGetActiveIngestionProfileInternalServerError creates GetActiveIngestionProfileInternalServerError with default headers values
func NewGetActiveIngestionProfileInternalServerError() *GetActiveIngestionProfileInternalServerError {

	return &GetActiveIngestionProfileInternalServerError{}
}

// WriteResponse to the client
func (o *GetActiveIngestionProfileInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
