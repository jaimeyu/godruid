// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetTenantThresholdProfileParams creates a new GetTenantThresholdProfileParams object
// no default values defined in spec.
func NewGetTenantThresholdProfileParams() GetTenantThresholdProfileParams {

	return GetTenantThresholdProfileParams{}
}

// GetTenantThresholdProfileParams contains all the bound params for the get tenant threshold profile operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetTenantThresholdProfile
type GetTenantThresholdProfileParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	TenantID string
	/*
	  Required: true
	  In: path
	*/
	ThrPrfID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetTenantThresholdProfileParams() beforehand.
func (o *GetTenantThresholdProfileParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rTenantID, rhkTenantID, _ := route.Params.GetOK("tenantId")
	if err := o.bindTenantID(rTenantID, rhkTenantID, route.Formats); err != nil {
		res = append(res, err)
	}

	rThrPrfID, rhkThrPrfID, _ := route.Params.GetOK("thrPrfId")
	if err := o.bindThrPrfID(rThrPrfID, rhkThrPrfID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTenantThresholdProfileParams) bindTenantID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.TenantID = raw

	return nil
}

func (o *GetTenantThresholdProfileParams) bindThrPrfID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ThrPrfID = raw

	return nil
}
