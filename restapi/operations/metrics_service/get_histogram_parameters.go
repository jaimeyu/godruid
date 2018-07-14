// Code generated by go-swagger; DO NOT EDIT.

package metrics_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetHistogramParams creates a new GetHistogramParams object
// no default values defined in spec.
func NewGetHistogramParams() GetHistogramParams {

	return GetHistogramParams{}
}

// GetHistogramParams contains all the bound params for the get histogram operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetHistogram
type GetHistogramParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	*/
	Direction *string
	/*
	  In: query
	*/
	Domain *string
	/*ISO-8601 period combination.
	  In: query
	*/
	Granularity *string
	/*
	  In: query
	*/
	GranularityBuckets *int32
	/*ISO-8601 Intervals.
	  In: query
	*/
	Interval *string
	/*
	  In: query
	*/
	Metric *string
	/*
	  In: query
	*/
	Resolution *int32
	/*
	  In: query
	*/
	Tenant *string
	/*
	  In: query
	*/
	Timeout *int32
	/*
	  In: query
	*/
	Vendor *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetHistogramParams() beforehand.
func (o *GetHistogramParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qDirection, qhkDirection, _ := qs.GetOK("direction")
	if err := o.bindDirection(qDirection, qhkDirection, route.Formats); err != nil {
		res = append(res, err)
	}

	qDomain, qhkDomain, _ := qs.GetOK("domain")
	if err := o.bindDomain(qDomain, qhkDomain, route.Formats); err != nil {
		res = append(res, err)
	}

	qGranularity, qhkGranularity, _ := qs.GetOK("granularity")
	if err := o.bindGranularity(qGranularity, qhkGranularity, route.Formats); err != nil {
		res = append(res, err)
	}

	qGranularityBuckets, qhkGranularityBuckets, _ := qs.GetOK("granularityBuckets")
	if err := o.bindGranularityBuckets(qGranularityBuckets, qhkGranularityBuckets, route.Formats); err != nil {
		res = append(res, err)
	}

	qInterval, qhkInterval, _ := qs.GetOK("interval")
	if err := o.bindInterval(qInterval, qhkInterval, route.Formats); err != nil {
		res = append(res, err)
	}

	qMetric, qhkMetric, _ := qs.GetOK("metric")
	if err := o.bindMetric(qMetric, qhkMetric, route.Formats); err != nil {
		res = append(res, err)
	}

	qResolution, qhkResolution, _ := qs.GetOK("resolution")
	if err := o.bindResolution(qResolution, qhkResolution, route.Formats); err != nil {
		res = append(res, err)
	}

	qTenant, qhkTenant, _ := qs.GetOK("tenant")
	if err := o.bindTenant(qTenant, qhkTenant, route.Formats); err != nil {
		res = append(res, err)
	}

	qTimeout, qhkTimeout, _ := qs.GetOK("timeout")
	if err := o.bindTimeout(qTimeout, qhkTimeout, route.Formats); err != nil {
		res = append(res, err)
	}

	qVendor, qhkVendor, _ := qs.GetOK("vendor")
	if err := o.bindVendor(qVendor, qhkVendor, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetHistogramParams) bindDirection(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Direction = &raw

	return nil
}

func (o *GetHistogramParams) bindDomain(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Domain = &raw

	return nil
}

func (o *GetHistogramParams) bindGranularity(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Granularity = &raw

	return nil
}

func (o *GetHistogramParams) bindGranularityBuckets(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("granularityBuckets", "query", "int32", raw)
	}
	o.GranularityBuckets = &value

	return nil
}

func (o *GetHistogramParams) bindInterval(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Interval = &raw

	return nil
}

func (o *GetHistogramParams) bindMetric(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Metric = &raw

	return nil
}

func (o *GetHistogramParams) bindResolution(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("resolution", "query", "int32", raw)
	}
	o.Resolution = &value

	return nil
}

func (o *GetHistogramParams) bindTenant(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Tenant = &raw

	return nil
}

func (o *GetHistogramParams) bindTimeout(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("timeout", "query", "int32", raw)
	}
	o.Timeout = &value

	return nil
}

func (o *GetHistogramParams) bindVendor(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Vendor = &raw

	return nil
}
