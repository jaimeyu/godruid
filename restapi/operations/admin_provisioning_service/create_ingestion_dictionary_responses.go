// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// CreateIngestionDictionaryOKCode is the HTTP code returned for type CreateIngestionDictionaryOK
const CreateIngestionDictionaryOKCode int = 200

/*CreateIngestionDictionaryOK create ingestion dictionary o k

swagger:response createIngestionDictionaryOK
*/
type CreateIngestionDictionaryOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPIIngestionDictionary `json:"body,omitempty"`
}

// NewCreateIngestionDictionaryOK creates CreateIngestionDictionaryOK with default headers values
func NewCreateIngestionDictionaryOK() *CreateIngestionDictionaryOK {

	return &CreateIngestionDictionaryOK{}
}

// WithPayload adds the payload to the create ingestion dictionary o k response
func (o *CreateIngestionDictionaryOK) WithPayload(payload *swagmodels.JSONAPIIngestionDictionary) *CreateIngestionDictionaryOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create ingestion dictionary o k response
func (o *CreateIngestionDictionaryOK) SetPayload(payload *swagmodels.JSONAPIIngestionDictionary) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateIngestionDictionaryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateIngestionDictionaryBadRequestCode is the HTTP code returned for type CreateIngestionDictionaryBadRequest
const CreateIngestionDictionaryBadRequestCode int = 400

/*CreateIngestionDictionaryBadRequest Request data does not pass validation

swagger:response createIngestionDictionaryBadRequest
*/
type CreateIngestionDictionaryBadRequest struct {
}

// NewCreateIngestionDictionaryBadRequest creates CreateIngestionDictionaryBadRequest with default headers values
func NewCreateIngestionDictionaryBadRequest() *CreateIngestionDictionaryBadRequest {

	return &CreateIngestionDictionaryBadRequest{}
}

// WriteResponse to the client
func (o *CreateIngestionDictionaryBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// CreateIngestionDictionaryInternalServerErrorCode is the HTTP code returned for type CreateIngestionDictionaryInternalServerError
const CreateIngestionDictionaryInternalServerErrorCode int = 500

/*CreateIngestionDictionaryInternalServerError Unexpected error processing request

swagger:response createIngestionDictionaryInternalServerError
*/
type CreateIngestionDictionaryInternalServerError struct {
}

// NewCreateIngestionDictionaryInternalServerError creates CreateIngestionDictionaryInternalServerError with default headers values
func NewCreateIngestionDictionaryInternalServerError() *CreateIngestionDictionaryInternalServerError {

	return &CreateIngestionDictionaryInternalServerError{}
}

// WriteResponse to the client
func (o *CreateIngestionDictionaryInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
