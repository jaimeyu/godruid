// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// BulkUpsertMonitoredObjectMetaOKCode is the HTTP code returned for type BulkUpsertMonitoredObjectMetaOK
const BulkUpsertMonitoredObjectMetaOKCode int = 200

/*BulkUpsertMonitoredObjectMetaOK bulk upsert monitored object meta o k

swagger:response bulkUpsertMonitoredObjectMetaOK
*/
type BulkUpsertMonitoredObjectMetaOK struct {

	/*
	  In: Body
	*/
	Payload swagmodels.BulkOperationResponse `json:"body,omitempty"`
}

// NewBulkUpsertMonitoredObjectMetaOK creates BulkUpsertMonitoredObjectMetaOK with default headers values
func NewBulkUpsertMonitoredObjectMetaOK() *BulkUpsertMonitoredObjectMetaOK {

	return &BulkUpsertMonitoredObjectMetaOK{}
}

// WithPayload adds the payload to the bulk upsert monitored object meta o k response
func (o *BulkUpsertMonitoredObjectMetaOK) WithPayload(payload swagmodels.BulkOperationResponse) *BulkUpsertMonitoredObjectMetaOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk upsert monitored object meta o k response
func (o *BulkUpsertMonitoredObjectMetaOK) SetPayload(payload swagmodels.BulkOperationResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkUpsertMonitoredObjectMetaOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make(swagmodels.BulkOperationResponse, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// BulkUpsertMonitoredObjectMetaBadRequestCode is the HTTP code returned for type BulkUpsertMonitoredObjectMetaBadRequest
const BulkUpsertMonitoredObjectMetaBadRequestCode int = 400

/*BulkUpsertMonitoredObjectMetaBadRequest Request data does not pass validation

swagger:response bulkUpsertMonitoredObjectMetaBadRequest
*/
type BulkUpsertMonitoredObjectMetaBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewBulkUpsertMonitoredObjectMetaBadRequest creates BulkUpsertMonitoredObjectMetaBadRequest with default headers values
func NewBulkUpsertMonitoredObjectMetaBadRequest() *BulkUpsertMonitoredObjectMetaBadRequest {

	return &BulkUpsertMonitoredObjectMetaBadRequest{}
}

// WithPayload adds the payload to the bulk upsert monitored object meta bad request response
func (o *BulkUpsertMonitoredObjectMetaBadRequest) WithPayload(payload string) *BulkUpsertMonitoredObjectMetaBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk upsert monitored object meta bad request response
func (o *BulkUpsertMonitoredObjectMetaBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkUpsertMonitoredObjectMetaBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// BulkUpsertMonitoredObjectMetaForbiddenCode is the HTTP code returned for type BulkUpsertMonitoredObjectMetaForbidden
const BulkUpsertMonitoredObjectMetaForbiddenCode int = 403

/*BulkUpsertMonitoredObjectMetaForbidden Requestor does not have authorization to perform this action

swagger:response bulkUpsertMonitoredObjectMetaForbidden
*/
type BulkUpsertMonitoredObjectMetaForbidden struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewBulkUpsertMonitoredObjectMetaForbidden creates BulkUpsertMonitoredObjectMetaForbidden with default headers values
func NewBulkUpsertMonitoredObjectMetaForbidden() *BulkUpsertMonitoredObjectMetaForbidden {

	return &BulkUpsertMonitoredObjectMetaForbidden{}
}

// WithPayload adds the payload to the bulk upsert monitored object meta forbidden response
func (o *BulkUpsertMonitoredObjectMetaForbidden) WithPayload(payload string) *BulkUpsertMonitoredObjectMetaForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk upsert monitored object meta forbidden response
func (o *BulkUpsertMonitoredObjectMetaForbidden) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkUpsertMonitoredObjectMetaForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// BulkUpsertMonitoredObjectMetaInternalServerErrorCode is the HTTP code returned for type BulkUpsertMonitoredObjectMetaInternalServerError
const BulkUpsertMonitoredObjectMetaInternalServerErrorCode int = 500

/*BulkUpsertMonitoredObjectMetaInternalServerError Unexpected error processing request

swagger:response bulkUpsertMonitoredObjectMetaInternalServerError
*/
type BulkUpsertMonitoredObjectMetaInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewBulkUpsertMonitoredObjectMetaInternalServerError creates BulkUpsertMonitoredObjectMetaInternalServerError with default headers values
func NewBulkUpsertMonitoredObjectMetaInternalServerError() *BulkUpsertMonitoredObjectMetaInternalServerError {

	return &BulkUpsertMonitoredObjectMetaInternalServerError{}
}

// WithPayload adds the payload to the bulk upsert monitored object meta internal server error response
func (o *BulkUpsertMonitoredObjectMetaInternalServerError) WithPayload(payload string) *BulkUpsertMonitoredObjectMetaInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the bulk upsert monitored object meta internal server error response
func (o *BulkUpsertMonitoredObjectMetaInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BulkUpsertMonitoredObjectMetaInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
