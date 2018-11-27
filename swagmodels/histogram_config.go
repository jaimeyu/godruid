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

	// An object that allows filtering on arbitrary metadata criteria and their values. Refer to the MetaFilter object for additional details
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

	// An ordered set of histogram buckets that should be filled with the appropriate metric data
	Buckets []*HistogramConfigMetricsItems0BucketsItems0 `json:"buckets"`
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
	var dataAO1 struct {
		Buckets []*HistogramConfigMetricsItems0BucketsItems0 `json:"buckets,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Buckets = dataAO1.Buckets

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

	var dataAO1 struct {
		Buckets []*HistogramConfigMetricsItems0BucketsItems0 `json:"buckets,omitempty"`
	}

	dataAO1.Buckets = m.Buckets

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this histogram config metrics items0
func (m *HistogramConfigMetricsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with MetricIdentifierFilter
	if err := m.MetricIdentifierFilter.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBuckets(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HistogramConfigMetricsItems0) validateBuckets(formats strfmt.Registry) error {

	if swag.IsZero(m.Buckets) { // not required
		return nil
	}

	for i := 0; i < len(m.Buckets); i++ {
		if swag.IsZero(m.Buckets[i]) { // not required
			continue
		}

		if m.Buckets[i] != nil {
			if err := m.Buckets[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("buckets" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

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

// HistogramConfigMetricsItems0BucketsItems0 histogram config metrics items0 buckets items0
// swagger:model HistogramConfigMetricsItems0BucketsItems0
type HistogramConfigMetricsItems0BucketsItems0 struct {

	// lower
	Lower *HistogramConfigMetricsItems0BucketsItems0Lower `json:"lower,omitempty"`

	// upper
	Upper *HistogramConfigMetricsItems0BucketsItems0Upper `json:"upper,omitempty"`
}

// Validate validates this histogram config metrics items0 buckets items0
func (m *HistogramConfigMetricsItems0BucketsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLower(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpper(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HistogramConfigMetricsItems0BucketsItems0) validateLower(formats strfmt.Registry) error {

	if swag.IsZero(m.Lower) { // not required
		return nil
	}

	if m.Lower != nil {
		if err := m.Lower.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("lower")
			}
			return err
		}
	}

	return nil
}

func (m *HistogramConfigMetricsItems0BucketsItems0) validateUpper(formats strfmt.Registry) error {

	if swag.IsZero(m.Upper) { // not required
		return nil
	}

	if m.Upper != nil {
		if err := m.Upper.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("upper")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0BucketsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0BucketsItems0) UnmarshalBinary(b []byte) error {
	var res HistogramConfigMetricsItems0BucketsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// HistogramConfigMetricsItems0BucketsItems0Lower The specification for the lower boundary of the bucket
// swagger:model HistogramConfigMetricsItems0BucketsItems0Lower
type HistogramConfigMetricsItems0BucketsItems0Lower struct {

	// If set to true, then the lower value is assumed to be exclusive. Otherwise a value of false or the absence of this value assumes that the lower value is to be taken inclusively
	Strict bool `json:"strict,omitempty"`

	// The lower, positive number to be used to describe the lowest value of the bucket. Omitting this value assumes that this bucket includes anything lower than the defined "upper" value
	// Required: true
	Value *float32 `json:"value"`
}

// Validate validates this histogram config metrics items0 buckets items0 lower
func (m *HistogramConfigMetricsItems0BucketsItems0Lower) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HistogramConfigMetricsItems0BucketsItems0Lower) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("lower"+"."+"value", "body", m.Value); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0BucketsItems0Lower) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0BucketsItems0Lower) UnmarshalBinary(b []byte) error {
	var res HistogramConfigMetricsItems0BucketsItems0Lower
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// HistogramConfigMetricsItems0BucketsItems0Upper The specification for the upper boundary of the bucket
// swagger:model HistogramConfigMetricsItems0BucketsItems0Upper
type HistogramConfigMetricsItems0BucketsItems0Upper struct {

	// If set to true, then the upper value is assumed to be exclusive. Otherwise a value of false or the absence of this value assumes that the upper value is to be taken inclusively
	Strict bool `json:"strict,omitempty"`

	// The upper, positive number to be used to describe the highest value of the bucket. Omitting this value assumes that this bucket includes anything higher than the defined "lower" value
	// Required: true
	Value *float32 `json:"value"`
}

// Validate validates this histogram config metrics items0 buckets items0 upper
func (m *HistogramConfigMetricsItems0BucketsItems0Upper) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HistogramConfigMetricsItems0BucketsItems0Upper) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("upper"+"."+"value", "body", m.Value); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0BucketsItems0Upper) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HistogramConfigMetricsItems0BucketsItems0Upper) UnmarshalBinary(b []byte) error {
	var res HistogramConfigMetricsItems0BucketsItems0Upper
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
