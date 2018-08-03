// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	swagmodels "github.com/accedian/adh-gather/swagmodels"
)

// NewUpdateTenantDomainParams creates a new UpdateTenantDomainParams object
// no default values defined in spec.
func NewUpdateTenantDomainParams() UpdateTenantDomainParams {

	return UpdateTenantDomainParams{}
}

// UpdateTenantDomainParams contains all the bound params for the update tenant domain operation
// typically these are obtained from a http.Request
//
// swagger:parameters UpdateTenantDomain
type UpdateTenantDomainParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Body *swagmodels.JSONAPITenantDomain
	/*
	  Required: true
	  In: path
	*/
	DomainID string
	/*
	  Required: true
	  In: path
	*/
	TenantID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdateTenantDomainParams() beforehand.
func (o *UpdateTenantDomainParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body swagmodels.JSONAPITenantDomain
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body"))
	}
	rDomainID, rhkDomainID, _ := route.Params.GetOK("domainId")
	if err := o.bindDomainID(rDomainID, rhkDomainID, route.Formats); err != nil {
		res = append(res, err)
	}

	rTenantID, rhkTenantID, _ := route.Params.GetOK("tenantId")
	if err := o.bindTenantID(rTenantID, rhkTenantID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDomainID binds and validates parameter DomainID from path.
func (o *UpdateTenantDomainParams) bindDomainID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.DomainID = raw

	return nil
}

// bindTenantID binds and validates parameter TenantID from path.
func (o *UpdateTenantDomainParams) bindTenantID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.TenantID = raw

	return nil
}
