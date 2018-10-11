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

// CardMetric card metric
// swagger:model CardMetric
type CardMetric struct {

	// directions
	Directions []string `json:"directions,omitempty"`

	// metric
	Metric string `json:"metric,omitempty"`

	// monitored object types
	MonitoredObjectTypes []string `json:"monitoredObjectTypes"`

	// options
	Options *CardMetricOptions `json:"options,omitempty"`

	// unit
	Unit string `json:"unit,omitempty"`

	// vendor
	Vendor string `json:"vendor,omitempty"`
}

// Validate validates this card metric
func (m *CardMetric) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateOptions(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CardMetric) validateOptions(formats strfmt.Registry) error {

	if swag.IsZero(m.Options) { // not required
		return nil
	}

	if m.Options != nil {
		if err := m.Options.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("options")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CardMetric) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CardMetric) UnmarshalBinary(b []byte) error {
	var res CardMetric
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// CardMetricOptions card metric options
// swagger:model CardMetricOptions
type CardMetricOptions struct {

	// aggregation
	// Enum: [none sum]
	Aggregation string `json:"aggregation,omitempty"`

	// bins
	Bins []float64 `json:"bins"`

	// buckets
	Buckets []interface{} `json:"buckets,omitempty"`

	// directions
	Directions []string `json:"directions,omitempty"`

	// format unit
	FormatUnit string `json:"formatUnit,omitempty"`

	// series
	Series []string `json:"series"`

	// type
	// Enum: [measure events bins]
	Type string `json:"type,omitempty"`

	// use bins
	UseBins bool `json:"useBins,omitempty"`

	// use explicit series
	UseExplicitSeries bool `json:"useExplicitSeries,omitempty"`
}

// Validate validates this card metric options
func (m *CardMetricOptions) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAggregation(formats); err != nil {
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

var cardMetricOptionsTypeAggregationPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["none","sum"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		cardMetricOptionsTypeAggregationPropEnum = append(cardMetricOptionsTypeAggregationPropEnum, v)
	}
}

const (

	// CardMetricOptionsAggregationNone captures enum value "none"
	CardMetricOptionsAggregationNone string = "none"

	// CardMetricOptionsAggregationSum captures enum value "sum"
	CardMetricOptionsAggregationSum string = "sum"
)

// prop value enum
func (m *CardMetricOptions) validateAggregationEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, cardMetricOptionsTypeAggregationPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CardMetricOptions) validateAggregation(formats strfmt.Registry) error {

	if swag.IsZero(m.Aggregation) { // not required
		return nil
	}

	// value enum
	if err := m.validateAggregationEnum("options"+"."+"aggregation", "body", m.Aggregation); err != nil {
		return err
	}

	return nil
}

var cardMetricOptionsTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["measure","events","bins"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		cardMetricOptionsTypeTypePropEnum = append(cardMetricOptionsTypeTypePropEnum, v)
	}
}

const (

	// CardMetricOptionsTypeMeasure captures enum value "measure"
	CardMetricOptionsTypeMeasure string = "measure"

	// CardMetricOptionsTypeEvents captures enum value "events"
	CardMetricOptionsTypeEvents string = "events"

	// CardMetricOptionsTypeBins captures enum value "bins"
	CardMetricOptionsTypeBins string = "bins"
)

// prop value enum
func (m *CardMetricOptions) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, cardMetricOptionsTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CardMetricOptions) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("options"+"."+"type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CardMetricOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CardMetricOptions) UnmarshalBinary(b []byte) error {
	var res CardMetricOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
