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

// HistogramConfig The necessary request parameters for the metric api call
// swagger:model HistogramConfig
type HistogramConfig struct {

	// An array of metric dimensions that filter-in metrics that adhere to those dimensions. Refer to the DimensionFilter object for further information
	Dimensions DimensionFilter `json:"dimensions,omitempty"`

	// The granularity for timeseries in ISO-8601 duration format, or ALL
	Granularity string `json:"granularity,omitempty"`

	// A value of true will have the aggregation request execute on all data regardless of whether it has been cleaned or not
	IgnoreCleaning bool `json:"ignoreCleaning,omitempty"`

	// Time boundary for the metrics under consideration using the ISO-8601 standard
	// Required: true
	Interval *string `json:"interval"`

	// meta
	Meta MetaFilter `json:"meta,omitempty"`

	// A list of the requested metric identifiers and the histogram buckets associated with those identifiers
	// Required: true
	Metrics []*HistogramConfigMetricsItems0 `json:"metrics"`

	// Query timeout in milliseconds
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this histogram config
func (m *HistogramConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDimensions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInterval(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMeta(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetrics(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HistogramConfig) validateDimensions(formats strfmt.Registry) error {

	if swag.IsZero(m.Dimensions) { // not required
		return nil
	}

	if err := m.Dimensions.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("dimensions")
		}
		return err
	}

	return nil
}

func (m *HistogramConfig) validateInterval(formats strfmt.Registry) error {

	if err := validate.Required("interval", "body", m.Interval); err != nil {
		return err
	}

	return nil
}

func (m *HistogramConfig) validateMeta(formats strfmt.Registry) error {

	if swag.IsZero(m.Meta) { // not required
		return nil
	}

	if err := m.Meta.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("meta")
		}
		return err
	}

	return nil
}

func (m *HistogramConfig) validateMetrics(formats strfmt.Registry) error {

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

// MarshalBinary interface implementation
func (m *HistogramConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HistogramConfig) UnmarshalBinary(b []byte) error {
	var res HistogramConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// HistogramConfigMetricsItems0 histogram config metrics items0
// swagger:model HistogramConfigMetricsItems0
type HistogramConfigMetricsItems0 struct {
	MetricIdentifierFilter

	BucketFilter
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *HistogramConfigMetricsItems0) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 MetricIdentifierFilter
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.MetricIdentifierFilter = aO0

	// AO1
	var aO1 BucketFilter
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.BucketFilter = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m HistogramConfigMetricsItems0) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.MetricIdentifierFilter)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.BucketFilter)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this histogram config metrics items0
func (m *HistogramConfigMetricsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with MetricIdentifierFilter
	if err := m.MetricIdentifierFilter.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with BucketFilter
	if err := m.BucketFilter.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0) UnmarshalBinary(b []byte) error {
	var res HistogramConfigMetricsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
