// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// MetricBaselineDataWrapper metric baseline data wrapper
// swagger:model MetricBaselineDataWrapper
type MetricBaselineDataWrapper struct {

	// attributes
	// Required: true
	Attributes *MetricBaselineData `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Required: true
	// Enum: [metricBaselineData]
	Type *string `json:"type"`
}

// Validate validates this metric baseline data wrapper
func (m *MetricBaselineDataWrapper) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MetricBaselineDataWrapper) validateAttributes(formats strfmt.Registry) error {

	if err := validate.Required("attributes", "body", m.Attributes); err != nil {
		return err
	}

	if m.Attributes != nil {
		if err := m.Attributes.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("attributes")
			}
			return err
		}
	}

	return nil
}

var metricBaselineDataWrapperTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["metricBaselineData"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		metricBaselineDataWrapperTypeTypePropEnum = append(metricBaselineDataWrapperTypeTypePropEnum, v)
	}
}

const (

	// MetricBaselineDataWrapperTypeMetricBaselineData captures enum value "metricBaselineData"
	MetricBaselineDataWrapperTypeMetricBaselineData string = "metricBaselineData"
)

// prop value enum
func (m *MetricBaselineDataWrapper) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, metricBaselineDataWrapperTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *MetricBaselineDataWrapper) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MetricBaselineDataWrapper) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MetricBaselineDataWrapper) UnmarshalBinary(b []byte) error {
	var res MetricBaselineDataWrapper
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
