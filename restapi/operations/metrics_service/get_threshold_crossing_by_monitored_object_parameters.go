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
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetThresholdCrossingByMonitoredObjectParams creates a new GetThresholdCrossingByMonitoredObjectParams object
// no default values defined in spec.
func NewGetThresholdCrossingByMonitoredObjectParams() GetThresholdCrossingByMonitoredObjectParams {

	return GetThresholdCrossingByMonitoredObjectParams{}
}

// GetThresholdCrossingByMonitoredObjectParams contains all the bound params for the get threshold crossing by monitored object operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetThresholdCrossingByMonitoredObject
type GetThresholdCrossingByMonitoredObjectParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	*/
	Direction []string
	/*Domain ID
	  In: query
	*/
	Domain []string
	/*ISO-8601 period combination.
	  In: query
	*/
	Granularity *string
	/*ISO-8601 Intervals.
	  Required: true
	  In: query
	*/
	Interval string
	/*
	  In: query
	*/
	Metric []string
	/*
	  In: query
	*/
	ObjectType []string
	/*Tenant ID
	  Required: true
	  In: query
	*/
	Tenant string
	/*
	  Required: true
	  In: query
	*/
	ThresholdProfileID string
	/*
	  In: query
	*/
	Timeout *int32
	/*
	  In: query
	*/
	Vendor []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetThresholdCrossingByMonitoredObjectParams() beforehand.
func (o *GetThresholdCrossingByMonitoredObjectParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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

	qInterval, qhkInterval, _ := qs.GetOK("interval")
	if err := o.bindInterval(qInterval, qhkInterval, route.Formats); err != nil {
		res = append(res, err)
	}

	qMetric, qhkMetric, _ := qs.GetOK("metric")
	if err := o.bindMetric(qMetric, qhkMetric, route.Formats); err != nil {
		res = append(res, err)
	}

	qObjectType, qhkObjectType, _ := qs.GetOK("objectType")
	if err := o.bindObjectType(qObjectType, qhkObjectType, route.Formats); err != nil {
		res = append(res, err)
	}

	qTenant, qhkTenant, _ := qs.GetOK("tenant")
	if err := o.bindTenant(qTenant, qhkTenant, route.Formats); err != nil {
		res = append(res, err)
	}

	qThresholdProfileID, qhkThresholdProfileID, _ := qs.GetOK("thresholdProfileId")
	if err := o.bindThresholdProfileID(qThresholdProfileID, qhkThresholdProfileID, route.Formats); err != nil {
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

// bindDirection binds and validates array parameter Direction from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetThresholdCrossingByMonitoredObjectParams) bindDirection(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvDirection string
	if len(rawData) > 0 {
		qvDirection = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	directionIC := swag.SplitByFormat(qvDirection, "")
	if len(directionIC) == 0 {
		return nil
	}

	var directionIR []string
	for _, directionIV := range directionIC {
		directionI := directionIV

		directionIR = append(directionIR, directionI)
	}

	o.Direction = directionIR

	return nil
}

// bindDomain binds and validates array parameter Domain from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetThresholdCrossingByMonitoredObjectParams) bindDomain(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvDomain string
	if len(rawData) > 0 {
		qvDomain = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	domainIC := swag.SplitByFormat(qvDomain, "")
	if len(domainIC) == 0 {
		return nil
	}

	var domainIR []string
	for _, domainIV := range domainIC {
		domainI := domainIV

		domainIR = append(domainIR, domainI)
	}

	o.Domain = domainIR

	return nil
}

// bindGranularity binds and validates parameter Granularity from query.
func (o *GetThresholdCrossingByMonitoredObjectParams) bindGranularity(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindInterval binds and validates parameter Interval from query.
func (o *GetThresholdCrossingByMonitoredObjectParams) bindInterval(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("interval", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("interval", "query", raw); err != nil {
		return err
	}

	o.Interval = raw

	return nil
}

// bindMetric binds and validates array parameter Metric from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetThresholdCrossingByMonitoredObjectParams) bindMetric(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvMetric string
	if len(rawData) > 0 {
		qvMetric = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	metricIC := swag.SplitByFormat(qvMetric, "")
	if len(metricIC) == 0 {
		return nil
	}

	var metricIR []string
	for _, metricIV := range metricIC {
		metricI := metricIV

		metricIR = append(metricIR, metricI)
	}

	o.Metric = metricIR

	return nil
}

// bindObjectType binds and validates array parameter ObjectType from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetThresholdCrossingByMonitoredObjectParams) bindObjectType(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvObjectType string
	if len(rawData) > 0 {
		qvObjectType = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	objectTypeIC := swag.SplitByFormat(qvObjectType, "")
	if len(objectTypeIC) == 0 {
		return nil
	}

	var objectTypeIR []string
	for _, objectTypeIV := range objectTypeIC {
		objectTypeI := objectTypeIV

		objectTypeIR = append(objectTypeIR, objectTypeI)
	}

	o.ObjectType = objectTypeIR

	return nil
}

// bindTenant binds and validates parameter Tenant from query.
func (o *GetThresholdCrossingByMonitoredObjectParams) bindTenant(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("tenant", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("tenant", "query", raw); err != nil {
		return err
	}

	o.Tenant = raw

	return nil
}

// bindThresholdProfileID binds and validates parameter ThresholdProfileID from query.
func (o *GetThresholdCrossingByMonitoredObjectParams) bindThresholdProfileID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("thresholdProfileId", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("thresholdProfileId", "query", raw); err != nil {
		return err
	}

	o.ThresholdProfileID = raw

	return nil
}

// bindTimeout binds and validates parameter Timeout from query.
func (o *GetThresholdCrossingByMonitoredObjectParams) bindTimeout(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindVendor binds and validates array parameter Vendor from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetThresholdCrossingByMonitoredObjectParams) bindVendor(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvVendor string
	if len(rawData) > 0 {
		qvVendor = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	vendorIC := swag.SplitByFormat(qvVendor, "")
	if len(vendorIC) == 0 {
		return nil
	}

	var vendorIR []string
	for _, vendorIV := range vendorIC {
		vendorI := vendorIV

		vendorIR = append(vendorIR, vendorI)
	}

	o.Vendor = vendorIR

	return nil
}
