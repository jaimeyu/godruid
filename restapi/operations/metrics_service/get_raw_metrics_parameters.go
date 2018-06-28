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

// NewGetRawMetricsParams creates a new GetRawMetricsParams object
// no default values defined in spec.
func NewGetRawMetricsParams() GetRawMetricsParams {

	return GetRawMetricsParams{}
}

// GetRawMetricsParams contains all the bound params for the get raw metrics operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetRawMetrics
type GetRawMetricsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	*/
	Direction *string
	/*
	  In: query
	*/
	Granularity *string
	/*ISO-8601 Intervals.
	  In: query
	*/
	Interval *string
	/*
	  In: query
	*/
	Metric []string
	/*
	  In: query
	*/
	MonitoredObjectID []string
	/*
	  In: query
	*/
	ObjectType *string
	/*
	  In: query
	*/
	Tenant *string
	/*
	  In: query
	*/
	Timeout *int32
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetRawMetricsParams() beforehand.
func (o *GetRawMetricsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qDirection, qhkDirection, _ := qs.GetOK("direction")
	if err := o.bindDirection(qDirection, qhkDirection, route.Formats); err != nil {
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

	qMonitoredObjectID, qhkMonitoredObjectID, _ := qs.GetOK("monitoredObjectId")
	if err := o.bindMonitoredObjectID(qMonitoredObjectID, qhkMonitoredObjectID, route.Formats); err != nil {
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

	qTimeout, qhkTimeout, _ := qs.GetOK("timeout")
	if err := o.bindTimeout(qTimeout, qhkTimeout, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetRawMetricsParams) bindDirection(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *GetRawMetricsParams) bindGranularity(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *GetRawMetricsParams) bindInterval(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *GetRawMetricsParams) bindMetric(rawData []string, hasKey bool, formats strfmt.Registry) error {

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

func (o *GetRawMetricsParams) bindMonitoredObjectID(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvMonitoredObjectID string
	if len(rawData) > 0 {
		qvMonitoredObjectID = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	monitoredObjectIDIC := swag.SplitByFormat(qvMonitoredObjectID, "")
	if len(monitoredObjectIDIC) == 0 {
		return nil
	}

	var monitoredObjectIDIR []string
	for _, monitoredObjectIDIV := range monitoredObjectIDIC {
		monitoredObjectIDI := monitoredObjectIDIV

		monitoredObjectIDIR = append(monitoredObjectIDIR, monitoredObjectIDI)
	}

	o.MonitoredObjectID = monitoredObjectIDIR

	return nil
}

func (o *GetRawMetricsParams) bindObjectType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.ObjectType = &raw

	return nil
}

func (o *GetRawMetricsParams) bindTenant(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *GetRawMetricsParams) bindTimeout(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
