// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
)

// NewGetAllMetadataConfigsV2Params creates a new GetAllMetadataConfigsV2Params object
// no default values defined in spec.
func NewGetAllMetadataConfigsV2Params() GetAllMetadataConfigsV2Params {

	return GetAllMetadataConfigsV2Params{}
}

// GetAllMetadataConfigsV2Params contains all the bound params for the get all metadata configs v2 operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetAllMetadataConfigsV2
type GetAllMetadataConfigsV2Params struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetAllMetadataConfigsV2Params() beforehand.
func (o *GetAllMetadataConfigsV2Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
