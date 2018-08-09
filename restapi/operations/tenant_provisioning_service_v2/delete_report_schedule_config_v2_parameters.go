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

// NewDeleteReportScheduleConfigV2Params creates a new DeleteReportScheduleConfigV2Params object
// no default values defined in spec.
func NewDeleteReportScheduleConfigV2Params() DeleteReportScheduleConfigV2Params {

	return DeleteReportScheduleConfigV2Params{}
}

// DeleteReportScheduleConfigV2Params contains all the bound params for the delete report schedule config v2 operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteReportScheduleConfigV2
type DeleteReportScheduleConfigV2Params struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ConfigID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteReportScheduleConfigV2Params() beforehand.
func (o *DeleteReportScheduleConfigV2Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rConfigID, rhkConfigID, _ := route.Params.GetOK("configId")
	if err := o.bindConfigID(rConfigID, rhkConfigID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindConfigID binds and validates parameter ConfigID from path.
func (o *DeleteReportScheduleConfigV2Params) bindConfigID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ConfigID = raw

	return nil
}
