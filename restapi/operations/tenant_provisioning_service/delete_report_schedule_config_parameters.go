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

// NewDeleteReportScheduleConfigParams creates a new DeleteReportScheduleConfigParams object
// no default values defined in spec.
func NewDeleteReportScheduleConfigParams() DeleteReportScheduleConfigParams {

	return DeleteReportScheduleConfigParams{}
}

// DeleteReportScheduleConfigParams contains all the bound params for the delete report schedule config operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteReportScheduleConfig
type DeleteReportScheduleConfigParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ConfigID string
	/*
	  Required: true
	  In: path
	*/
	TenantID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteReportScheduleConfigParams() beforehand.
func (o *DeleteReportScheduleConfigParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rConfigID, rhkConfigID, _ := route.Params.GetOK("configId")
	if err := o.bindConfigID(rConfigID, rhkConfigID, route.Formats); err != nil {
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

// bindConfigID binds and validates parameter ConfigID from path.
func (o *DeleteReportScheduleConfigParams) bindConfigID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ConfigID = raw

	return nil
}

// bindTenantID binds and validates parameter TenantID from path.
func (o *DeleteReportScheduleConfigParams) bindTenantID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.TenantID = raw

	return nil
}
