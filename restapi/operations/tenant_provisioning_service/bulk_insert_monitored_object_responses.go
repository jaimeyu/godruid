// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// BulkInsertMonitoredObjectOKCode is the HTTP code returned for type BulkInsertMonitoredObjectOK
const BulkInsertMonitoredObjectOKCode int = 200

/*BulkInsertMonitoredObjectOK bulk insert monitored object o k

swagger:response bulkInsertMonitoredObjectOK
*/
type BulkInsertMonitoredObjectOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.BulkOperationResult `json:"body,omitempty"`
}

// NewBulkInsertMonitoredObjectOK creates BulkInsertMonitoredObjectOK with default headers values
func NewBulkInsertMonitoredObjectOK() *BulkInsertMonitoredObjectOK {

	return &BulkInsertMonitoredObjectOK{}
}

// WithPayload adds the payload to the bulk insert monitored object o k response
func (o *BulkInsertMonitoredObjectOK) WithPayload(payload *swagmodels.BulkOperationResult) *BulkInsertMonitoredObjectOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk insert monitored object o k response
func (o *BulkInsertMonitoredObjectOK) SetPayload(payload *swagmodels.BulkOperationResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkInsertMonitoredObjectOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BulkInsertMonitoredObjectBadRequestCode is the HTTP code returned for type BulkInsertMonitoredObjectBadRequest
const BulkInsertMonitoredObjectBadRequestCode int = 400

/*BulkInsertMonitoredObjectBadRequest Request data does not pass validation

swagger:response bulkInsertMonitoredObjectBadRequest
*/
type BulkInsertMonitoredObjectBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewBulkInsertMonitoredObjectBadRequest creates BulkInsertMonitoredObjectBadRequest with default headers values
func NewBulkInsertMonitoredObjectBadRequest() *BulkInsertMonitoredObjectBadRequest {

	return &BulkInsertMonitoredObjectBadRequest{}
}

// WithPayload adds the payload to the bulk insert monitored object bad request response
func (o *BulkInsertMonitoredObjectBadRequest) WithPayload(payload string) *BulkInsertMonitoredObjectBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk insert monitored object bad request response
func (o *BulkInsertMonitoredObjectBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkInsertMonitoredObjectBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// BulkInsertMonitoredObjectForbiddenCode is the HTTP code returned for type BulkInsertMonitoredObjectForbidden
const BulkInsertMonitoredObjectForbiddenCode int = 403

/*BulkInsertMonitoredObjectForbidden Requestor does not have authorization to perform this action

swagger:response bulkInsertMonitoredObjectForbidden
*/
type BulkInsertMonitoredObjectForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewBulkInsertMonitoredObjectForbidden creates BulkInsertMonitoredObjectForbidden with default headers values
func NewBulkInsertMonitoredObjectForbidden() *BulkInsertMonitoredObjectForbidden {

	return &BulkInsertMonitoredObjectForbidden{}
}

// WithPayload adds the payload to the bulk insert monitored object forbidden response
func (o *BulkInsertMonitoredObjectForbidden) WithPayload(payload string) *BulkInsertMonitoredObjectForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk insert monitored object forbidden response
func (o *BulkInsertMonitoredObjectForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkInsertMonitoredObjectForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// BulkInsertMonitoredObjectInternalServerErrorCode is the HTTP code returned for type BulkInsertMonitoredObjectInternalServerError
const BulkInsertMonitoredObjectInternalServerErrorCode int = 500

/*BulkInsertMonitoredObjectInternalServerError Unexpected error processing request

swagger:response bulkInsertMonitoredObjectInternalServerError
*/
type BulkInsertMonitoredObjectInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewBulkInsertMonitoredObjectInternalServerError creates BulkInsertMonitoredObjectInternalServerError with default headers values
func NewBulkInsertMonitoredObjectInternalServerError() *BulkInsertMonitoredObjectInternalServerError {

	return &BulkInsertMonitoredObjectInternalServerError{}
}

// WithPayload adds the payload to the bulk insert monitored object internal server error response
func (o *BulkInsertMonitoredObjectInternalServerError) WithPayload(payload string) *BulkInsertMonitoredObjectInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk insert monitored object internal server error response
func (o *BulkInsertMonitoredObjectInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkInsertMonitoredObjectInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
