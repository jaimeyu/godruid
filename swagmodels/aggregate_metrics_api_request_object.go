// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

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

	// the granularity for timeseries in ISO-8601 duration format, or ALL
	Granularity string `json:"granularity,omitempty"`

	// ISO-8601 interval
	// Required: true
	Interval *string `json:"interval"`

	// set of meta keys and list of values for the purposes of filtering
	Meta map[string][]string `json:"meta,omitempty"`

	// metrics
	// Required: true
	Metrics []*MetricIdentifierObject `json:"metrics"`

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
		res = append(res, err)
	}

	if err := m.validateInterval(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetrics(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTenantID(formats); err != nil {
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

	for i := 0; i < len(m.Metrics); i++ {
		if swag.IsZero(m.Metrics[i]) { // not required
			continue
		}

		if m.Metrics[i] != nil {
			if err := m.Metrics[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("metrics" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

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

// AggregateMetricsAPIRequestObjectAggregation the aggregation function
// swagger:model AggregateMetricsAPIRequestObjectAggregation
type AggregateMetricsAPIRequestObjectAggregation struct {

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this aggregate metrics API request object aggregation
func (m *AggregateMetricsAPIRequestObjectAggregation) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AggregateMetricsAPIRequestObjectAggregation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AggregateMetricsAPIRequestObjectAggregation) UnmarshalBinary(b []byte) error {
	var res AggregateMetricsAPIRequestObjectAggregation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
