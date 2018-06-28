// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// DeleteTenantMonitoredObjectOKCode is the HTTP code returned for type DeleteTenantMonitoredObjectOK
const DeleteTenantMonitoredObjectOKCode int = 200

/*DeleteTenantMonitoredObjectOK delete tenant monitored object o k

swagger:response deleteTenantMonitoredObjectOK
*/
type DeleteTenantMonitoredObjectOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPITenantMonitoredObject `json:"body,omitempty"`
}

// NewDeleteTenantMonitoredObjectOK creates DeleteTenantMonitoredObjectOK with default headers values
func NewDeleteTenantMonitoredObjectOK() *DeleteTenantMonitoredObjectOK {

	return &DeleteTenantMonitoredObjectOK{}
}

// WithPayload adds the payload to the delete tenant monitored object o k response
func (o *DeleteTenantMonitoredObjectOK) WithPayload(payload *swagmodels.JSONAPITenantMonitoredObject) *DeleteTenantMonitoredObjectOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete tenant monitored object o k response
func (o *DeleteTenantMonitoredObjectOK) SetPayload(payload *swagmodels.JSONAPITenantMonitoredObject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTenantMonitoredObjectOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteTenantMonitoredObjectInternalServerErrorCode is the HTTP code returned for type DeleteTenantMonitoredObjectInternalServerError
const DeleteTenantMonitoredObjectInternalServerErrorCode int = 500

/*DeleteTenantMonitoredObjectInternalServerError Unexpected error processing request

swagger:response deleteTenantMonitoredObjectInternalServerError
*/
type DeleteTenantMonitoredObjectInternalServerError struct {
}

// NewDeleteTenantMonitoredObjectInternalServerError creates DeleteTenantMonitoredObjectInternalServerError with default headers values
func NewDeleteTenantMonitoredObjectInternalServerError() *DeleteTenantMonitoredObjectInternalServerError {

	return &DeleteTenantMonitoredObjectInternalServerError{}
}

// WriteResponse to the client
func (o *DeleteTenantMonitoredObjectInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
