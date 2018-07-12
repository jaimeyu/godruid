// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// TopNMetricsAPIRequestObject top n metrics API request object
// swagger:model TopNMetricsAPIRequestObject
type TopNMetricsAPIRequestObject struct {

	// The type of aggregation (avg/min/max)
	// Required: true
	Aggregator *string `json:"aggregator"`

	// set of domains identifiers to use for filtering
	Domains []string `json:"domains"`

	// ISO-8601 interval
	// Required: true
	Interval *string `json:"interval"`

	// metrics
	// Required: true
	Metrics TopNMetricsAPIRequestObjectMetrics `json:"metrics"`

	// metrics view
	MetricsView TopNMetricsAPIRequestObjectMetricsView `json:"metricsView"`

	// set of monitored objects identifiers to use for filtering
	MonitoredObjects []string `json:"monitoredObjects"`

	// Number of results to return
	NumResults int64 `json:"numResults,omitempty"`

	// the tenant identifier
	// Required: true
	TenantID *string `json:"tenantId"`

	// query timeout in milliseconds
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this top n metrics API request object
func (m *TopNMetricsAPIRequestObject) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAggregator(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateDomains(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateInterval(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateMetrics(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateMonitoredObjects(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTenantID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TopNMetricsAPIRequestObject) validateAggregator(formats strfmt.Registry) error {

	if err := validate.Required("aggregator", "body", m.Aggregator); err != nil {
		return err
	}

	return nil
}

func (m *TopNMetricsAPIRequestObject) validateDomains(formats strfmt.Registry) error {

	if swag.IsZero(m.Domains) { // not required
		return nil
	}

	return nil
}

func (m *TopNMetricsAPIRequestObject) validateInterval(formats strfmt.Registry) error {

	if err := validate.Required("interval", "body", m.Interval); err != nil {
		return err
	}

	return nil
}

func (m *TopNMetricsAPIRequestObject) validateMetrics(formats strfmt.Registry) error {

	if err := validate.Required("metrics", "body", m.Metrics); err != nil {
		return err
	}

	if err := m.Metrics.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("metrics")
		}
		return err
	}

	return nil
}

func (m *TopNMetricsAPIRequestObject) validateMonitoredObjects(formats strfmt.Registry) error {

	if swag.IsZero(m.MonitoredObjects) { // not required
		return nil
	}

	return nil
}

func (m *TopNMetricsAPIRequestObject) validateTenantID(formats strfmt.Registry) error {

	if err := validate.Required("tenantId", "body", m.TenantID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TopNMetricsAPIRequestObject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TopNMetricsAPIRequestObject) UnmarshalBinary(b []byte) error {
	var res TopNMetricsAPIRequestObject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
