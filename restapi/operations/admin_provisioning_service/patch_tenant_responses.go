// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// PatchTenantOKCode is the HTTP code returned for type PatchTenantOK
const PatchTenantOKCode int = 200

/*PatchTenantOK patch tenant o k

swagger:response patchTenantOK
*/
type PatchTenantOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenant `json:"body,omitempty"`
}

// NewPatchTenantOK creates PatchTenantOK with default headers values
func NewPatchTenantOK() *PatchTenantOK {

	return &PatchTenantOK{}
}

// WithPayload adds the payload to the patch tenant o k response
func (o *PatchTenantOK) WithPayload(payload *swagmodels.JSONAPITenant) *PatchTenantOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the patch tenant o k response
func (o *PatchTenantOK) SetPayload(payload *swagmodels.JSONAPITenant) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PatchTenantOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PatchTenantBadRequestCode is the HTTP code returned for type PatchTenantBadRequest
const PatchTenantBadRequestCode int = 400

/*PatchTenantBadRequest Request data does not pass validation

swagger:response patchTenantBadRequest
*/
type PatchTenantBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPatchTenantBadRequest creates PatchTenantBadRequest with default headers values
func NewPatchTenantBadRequest() *PatchTenantBadRequest {

	return &PatchTenantBadRequest{}
}

// WithPayload adds the payload to the patch tenant bad request response
func (o *PatchTenantBadRequest) WithPayload(payload string) *PatchTenantBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the patch tenant bad request response
func (o *PatchTenantBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PatchTenantBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// PatchTenantForbiddenCode is the HTTP code returned for type PatchTenantForbidden
const PatchTenantForbiddenCode int = 403

/*PatchTenantForbidden Requestor does not have authorization to perform this action

swagger:response patchTenantForbidden
*/
type PatchTenantForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPatchTenantForbidden creates PatchTenantForbidden with default headers values
func NewPatchTenantForbidden() *PatchTenantForbidden {

	return &PatchTenantForbidden{}
}

// WithPayload adds the payload to the patch tenant forbidden response
func (o *PatchTenantForbidden) WithPayload(payload string) *PatchTenantForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the patch tenant forbidden response
func (o *PatchTenantForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PatchTenantForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// PatchTenantInternalServerErrorCode is the HTTP code returned for type PatchTenantInternalServerError
const PatchTenantInternalServerErrorCode int = 500

/*PatchTenantInternalServerError Unexpected error processing request

swagger:response patchTenantInternalServerError
*/
type PatchTenantInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPatchTenantInternalServerError creates PatchTenantInternalServerError with default headers values
func NewPatchTenantInternalServerError() *PatchTenantInternalServerError {

	return &PatchTenantInternalServerError{}
}

// WithPayload adds the payload to the patch tenant internal server error response
func (o *PatchTenantInternalServerError) WithPayload(payload string) *PatchTenantInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the patch tenant internal server error response
func (o *PatchTenantInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PatchTenantInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
