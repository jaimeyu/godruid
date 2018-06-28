// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetTenantSummaryByAliasOKCode is the HTTP code returned for type GetTenantSummaryByAliasOK
const GetTenantSummaryByAliasOKCode int = 200

/*GetTenantSummaryByAliasOK get tenant summary by alias o k

swagger:response getTenantSummaryByAliasOK
*/
type GetTenantSummaryByAliasOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.TenantSummary `json:"body,omitempty"`
}

// NewGetTenantSummaryByAliasOK creates GetTenantSummaryByAliasOK with default headers values
func NewGetTenantSummaryByAliasOK() *GetTenantSummaryByAliasOK {

	return &GetTenantSummaryByAliasOK{}
}

// WithPayload adds the payload to the get tenant summary by alias o k response
func (o *GetTenantSummaryByAliasOK) WithPayload(payload *swagmodels.TenantSummary) *GetTenantSummaryByAliasOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tenant summary by alias o k response
func (o *GetTenantSummaryByAliasOK) SetPayload(payload *swagmodels.TenantSummary) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTenantSummaryByAliasOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetTenantSummaryByAliasInternalServerErrorCode is the HTTP code returned for type GetTenantSummaryByAliasInternalServerError
const GetTenantSummaryByAliasInternalServerErrorCode int = 500

/*GetTenantSummaryByAliasInternalServerError Unexpected error processing request

swagger:response getTenantSummaryByAliasInternalServerError
*/
type GetTenantSummaryByAliasInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetTenantSummaryByAliasInternalServerError creates GetTenantSummaryByAliasInternalServerError with default headers values
func NewGetTenantSummaryByAliasInternalServerError() *GetTenantSummaryByAliasInternalServerError {

	return &GetTenantSummaryByAliasInternalServerError{}
}

// WithPayload adds the payload to the get tenant summary by alias internal server error response
func (o *GetTenantSummaryByAliasInternalServerError) WithPayload(payload string) *GetTenantSummaryByAliasInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tenant summary by alias internal server error response
func (o *GetTenantSummaryByAliasInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTenantSummaryByAliasInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
