// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetTenantIDByAliasOKCode is the HTTP code returned for type GetTenantIDByAliasOK
const GetTenantIDByAliasOKCode int = 200

/*GetTenantIDByAliasOK get tenant Id by alias o k

swagger:response getTenantIdByAliasOK
*/
type GetTenantIDByAliasOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetTenantIDByAliasOK creates GetTenantIDByAliasOK with default headers values
func NewGetTenantIDByAliasOK() *GetTenantIDByAliasOK {

	return &GetTenantIDByAliasOK{}
}

// WithPayload adds the payload to the get tenant Id by alias o k response
func (o *GetTenantIDByAliasOK) WithPayload(payload string) *GetTenantIDByAliasOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tenant Id by alias o k response
func (o *GetTenantIDByAliasOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTenantIDByAliasOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetTenantIDByAliasInternalServerErrorCode is the HTTP code returned for type GetTenantIDByAliasInternalServerError
const GetTenantIDByAliasInternalServerErrorCode int = 500

/*GetTenantIDByAliasInternalServerError Unexpected error processing request

swagger:response getTenantIdByAliasInternalServerError
*/
type GetTenantIDByAliasInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetTenantIDByAliasInternalServerError creates GetTenantIDByAliasInternalServerError with default headers values
func NewGetTenantIDByAliasInternalServerError() *GetTenantIDByAliasInternalServerError {

	return &GetTenantIDByAliasInternalServerError{}
}

// WithPayload adds the payload to the get tenant Id by alias internal server error response
func (o *GetTenantIDByAliasInternalServerError) WithPayload(payload string) *GetTenantIDByAliasInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get tenant Id by alias internal server error response
func (o *GetTenantIDByAliasInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTenantIDByAliasInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
