// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// GetAllAdminUsersOKCode is the HTTP code returned for type GetAllAdminUsersOK
const GetAllAdminUsersOKCode int = 200

/*GetAllAdminUsersOK get all admin users o k

swagger:response getAllAdminUsersOK
*/
type GetAllAdminUsersOK struct {

	/*
	  In: Body
	*/
	Payload *swagmodels.JSONAPIAdminUserList `json:"body,omitempty"`
}

// NewGetAllAdminUsersOK creates GetAllAdminUsersOK with default headers values
func NewGetAllAdminUsersOK() *GetAllAdminUsersOK {

	return &GetAllAdminUsersOK{}
}

// WithPayload adds the payload to the get all admin users o k response
func (o *GetAllAdminUsersOK) WithPayload(payload *swagmodels.JSONAPIAdminUserList) *GetAllAdminUsersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all admin users o k response
func (o *GetAllAdminUsersOK) SetPayload(payload *swagmodels.JSONAPIAdminUserList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllAdminUsersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
