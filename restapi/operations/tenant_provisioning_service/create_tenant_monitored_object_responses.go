// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// CreateTenantMonitoredObjectOKCode is the HTTP code returned for type CreateTenantMonitoredObjectOK
const CreateTenantMonitoredObjectOKCode int = 200

/*CreateTenantMonitoredObjectOK create tenant monitored object o k

swagger:response createTenantMonitoredObjectOK
*/
type CreateTenantMonitoredObjectOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantMonitoredObject `json:"body,omitempty"`
}

// NewCreateTenantMonitoredObjectOK creates CreateTenantMonitoredObjectOK with default headers values
func NewCreateTenantMonitoredObjectOK() *CreateTenantMonitoredObjectOK {

	return &CreateTenantMonitoredObjectOK{}
}

// WithPayload adds the payload to the create tenant monitored object o k response
func (o *CreateTenantMonitoredObjectOK) WithPayload(payload *swagmodels.JSONAPITenantMonitoredObject) *CreateTenantMonitoredObjectOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant monitored object o k response
func (o *CreateTenantMonitoredObjectOK) SetPayload(payload *swagmodels.JSONAPITenantMonitoredObject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantMonitoredObjectOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTenantMonitoredObjectBadRequestCode is the HTTP code returned for type CreateTenantMonitoredObjectBadRequest
const CreateTenantMonitoredObjectBadRequestCode int = 400

/*CreateTenantMonitoredObjectBadRequest Request data does not pass validation

swagger:response createTenantMonitoredObjectBadRequest
*/
type CreateTenantMonitoredObjectBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantMonitoredObjectBadRequest creates CreateTenantMonitoredObjectBadRequest with default headers values
func NewCreateTenantMonitoredObjectBadRequest() *CreateTenantMonitoredObjectBadRequest {

	return &CreateTenantMonitoredObjectBadRequest{}
}

// WithPayload adds the payload to the create tenant monitored object bad request response
func (o *CreateTenantMonitoredObjectBadRequest) WithPayload(payload string) *CreateTenantMonitoredObjectBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant monitored object bad request response
func (o *CreateTenantMonitoredObjectBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantMonitoredObjectBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantMonitoredObjectForbiddenCode is the HTTP code returned for type CreateTenantMonitoredObjectForbidden
const CreateTenantMonitoredObjectForbiddenCode int = 403

/*CreateTenantMonitoredObjectForbidden Requestor does not have authorization to perform this action

swagger:response createTenantMonitoredObjectForbidden
*/
type CreateTenantMonitoredObjectForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantMonitoredObjectForbidden creates CreateTenantMonitoredObjectForbidden with default headers values
func NewCreateTenantMonitoredObjectForbidden() *CreateTenantMonitoredObjectForbidden {

	return &CreateTenantMonitoredObjectForbidden{}
}

// WithPayload adds the payload to the create tenant monitored object forbidden response
func (o *CreateTenantMonitoredObjectForbidden) WithPayload(payload string) *CreateTenantMonitoredObjectForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant monitored object forbidden response
func (o *CreateTenantMonitoredObjectForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantMonitoredObjectForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// CreateTenantMonitoredObjectInternalServerErrorCode is the HTTP code returned for type CreateTenantMonitoredObjectInternalServerError
const CreateTenantMonitoredObjectInternalServerErrorCode int = 500

/*CreateTenantMonitoredObjectInternalServerError Unexpected error processing request

swagger:response createTenantMonitoredObjectInternalServerError
*/
type CreateTenantMonitoredObjectInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateTenantMonitoredObjectInternalServerError creates CreateTenantMonitoredObjectInternalServerError with default headers values
func NewCreateTenantMonitoredObjectInternalServerError() *CreateTenantMonitoredObjectInternalServerError {

	return &CreateTenantMonitoredObjectInternalServerError{}
}

// WithPayload adds the payload to the create tenant monitored object internal server error response
func (o *CreateTenantMonitoredObjectInternalServerError) WithPayload(payload string) *CreateTenantMonitoredObjectInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create tenant monitored object internal server error response
func (o *CreateTenantMonitoredObjectInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTenantMonitoredObjectInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
