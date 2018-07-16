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

// MetricViewObject metric view object
// swagger:model MetricViewObject
type MetricViewObject struct {

	// aggregator
	// Required: true
	Aggregator *string `json:"aggregator"`

	// metric
	// Required: true
	Metric *string `json:"metric"`

	// name
	// Required: true
	Name *string `json:"name"`
}

// Validate validates this metric view object
func (m *MetricViewObject) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAggregator(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateMetric(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MetricViewObject) validateAggregator(formats strfmt.Registry) error {

	if err := validate.Required("aggregator", "body", m.Aggregator); err != nil {
		return err
	}

	return nil
}

func (m *MetricViewObject) validateMetric(formats strfmt.Registry) error {

	if err := validate.Required("metric", "body", m.Metric); err != nil {
		return err
	}

	return nil
}

func (m *MetricViewObject) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MetricViewObject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MetricViewObject) UnmarshalBinary(b []byte) error {
	var res MetricViewObject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
