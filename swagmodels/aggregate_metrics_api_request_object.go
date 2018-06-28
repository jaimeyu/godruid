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

// AggregateMetricsAPIRequestObject aggregate metrics API request object
// swagger:model AggregateMetricsAPIRequestObject
type AggregateMetricsAPIRequestObject struct {

	// aggregation
	// Required: true
	Aggregation *AggregateMetricsAPIRequestObjectAggregation `json:"aggregation"`

	// set of domains identifiers to use for filtering
	DomainIds []string `json:"domainIds"`

	// the granularity for timeseries in ISO-8601 duration format, or ALL
	Granularity string `json:"granularity,omitempty"`

	// ISO-8601 interval
	// Required: true
	Interval *string `json:"interval"`

	// metrics
	// Required: true
	Metrics AggregateMetricsAPIRequestObjectMetrics `json:"metrics"`

	// the tenant identifier
	// Required: true
	TenantID *string `json:"tenantId"`

	// query timeout in milliseconds
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this aggregate metrics API request object
func (m *AggregateMetricsAPIRequestObject) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAggregation(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateDomainIds(formats); err != nil {
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

	if err := m.validateTenantID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AggregateMetricsAPIRequestObject) validateAggregation(formats strfmt.Registry) error {

	if err := validate.Required("aggregation", "body", m.Aggregation); err != nil {
		return err
	}

	if m.Aggregation != nil {

		if err := m.Aggregation.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("aggregation")
			}
			return err
		}

	}

	return nil
}

func (m *AggregateMetricsAPIRequestObject) validateDomainIds(formats strfmt.Registry) error {

	if swag.IsZero(m.DomainIds) { // not required
		return nil
	}

	return nil
}

func (m *AggregateMetricsAPIRequestObject) validateInterval(formats strfmt.Registry) error {

	if err := validate.Required("interval", "body", m.Interval); err != nil {
		return err
	}

	return nil
}

func (m *AggregateMetricsAPIRequestObject) validateMetrics(formats strfmt.Registry) error {

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

func (m *AggregateMetricsAPIRequestObject) validateTenantID(formats strfmt.Registry) error {

	if err := validate.Required("tenantId", "body", m.TenantID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AggregateMetricsAPIRequestObject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AggregateMetricsAPIRequestObject) UnmarshalBinary(b []byte) error {
	var res AggregateMetricsAPIRequestObject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
