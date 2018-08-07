// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteConnectorInstanceV2Params creates a new DeleteConnectorInstanceV2Params object
// no default values defined in spec.
func NewDeleteConnectorInstanceV2Params() DeleteConnectorInstanceV2Params {

	return DeleteConnectorInstanceV2Params{}
}

// DeleteConnectorInstanceV2Params contains all the bound params for the delete connector instance v2 operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteConnectorInstanceV2
type DeleteConnectorInstanceV2Params struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ConnectorInstanceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteConnectorInstanceV2Params() beforehand.
func (o *DeleteConnectorInstanceV2Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rConnectorInstanceID, rhkConnectorInstanceID, _ := route.Params.GetOK("connectorInstanceId")
	if err := o.bindConnectorInstanceID(rConnectorInstanceID, rhkConnectorInstanceID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindConnectorInstanceID binds and validates parameter ConnectorInstanceID from path.
func (o *DeleteConnectorInstanceV2Params) bindConnectorInstanceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ConnectorInstanceID = raw

	return nil
}
