// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ThresholdCrossingReportMetric threshold crossing report metric
// swagger:model ThresholdCrossingReportMetric
type ThresholdCrossingReportMetric struct {
	MetricIdentifierFilter

	// critical
	Critical []*ThresholdCrossingViolations `json:"critical"`

	// major
	Major []*ThresholdCrossingViolations `json:"major"`

	// minor
	Minor []*ThresholdCrossingViolations `json:"minor"`

	// sla
	SLA []*ThresholdCrossingViolations `json:"sla"`

	// warning
	Warning []*ThresholdCrossingViolations `json:"warning"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *ThresholdCrossingReportMetric) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 MetricIdentifierFilter
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.MetricIdentifierFilter = aO0

	// AO1
	var dataAO1 struct {
		Critical []*ThresholdCrossingViolations `json:"critical,omitempty"`

		Major []*ThresholdCrossingViolations `json:"major,omitempty"`

		Minor []*ThresholdCrossingViolations `json:"minor,omitempty"`

		SLA []*ThresholdCrossingViolations `json:"sla,omitempty"`

		Warning []*ThresholdCrossingViolations `json:"warning,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Critical = dataAO1.Critical

	m.Major = dataAO1.Major

	m.Minor = dataAO1.Minor

	m.SLA = dataAO1.SLA

	m.Warning = dataAO1.Warning

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m ThresholdCrossingReportMetric) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.MetricIdentifierFilter)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	var dataAO1 struct {
		Critical []*ThresholdCrossingViolations `json:"critical,omitempty"`

		Major []*ThresholdCrossingViolations `json:"major,omitempty"`

		Minor []*ThresholdCrossingViolations `json:"minor,omitempty"`

		SLA []*ThresholdCrossingViolations `json:"sla,omitempty"`

		Warning []*ThresholdCrossingViolations `json:"warning,omitempty"`
	}

	dataAO1.Critical = m.Critical

	dataAO1.Major = m.Major

	dataAO1.Minor = m.Minor

	dataAO1.SLA = m.SLA

	dataAO1.Warning = m.Warning

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this threshold crossing report metric
func (m *ThresholdCrossingReportMetric) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with MetricIdentifierFilter
	if err := m.MetricIdentifierFilter.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCritical(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMajor(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMinor(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSLA(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWarning(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ThresholdCrossingReportMetric) validateCritical(formats strfmt.Registry) error {

	if swag.IsZero(m.Critical) { // not required
		return nil
	}

	for i := 0; i < len(m.Critical); i++ {
		if swag.IsZero(m.Critical[i]) { // not required
			continue
		}

		if m.Critical[i] != nil {
			if err := m.Critical[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("critical" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ThresholdCrossingReportMetric) validateMajor(formats strfmt.Registry) error {

	if swag.IsZero(m.Major) { // not required
		return nil
	}

	for i := 0; i < len(m.Major); i++ {
		if swag.IsZero(m.Major[i]) { // not required
			continue
		}

		if m.Major[i] != nil {
			if err := m.Major[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("major" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ThresholdCrossingReportMetric) validateMinor(formats strfmt.Registry) error {

	if swag.IsZero(m.Minor) { // not required
		return nil
	}

	for i := 0; i < len(m.Minor); i++ {
		if swag.IsZero(m.Minor[i]) { // not required
			continue
		}

		if m.Minor[i] != nil {
			if err := m.Minor[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("minor" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ThresholdCrossingReportMetric) validateSLA(formats strfmt.Registry) error {

	if swag.IsZero(m.SLA) { // not required
		return nil
	}

	for i := 0; i < len(m.SLA); i++ {
		if swag.IsZero(m.SLA[i]) { // not required
			continue
		}

		if m.SLA[i] != nil {
			if err := m.SLA[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("sla" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ThresholdCrossingReportMetric) validateWarning(formats strfmt.Registry) error {

	if swag.IsZero(m.Warning) { // not required
		return nil
	}

	for i := 0; i < len(m.Warning); i++ {
		if swag.IsZero(m.Warning[i]) { // not required
			continue
		}

		if m.Warning[i] != nil {
			if err := m.Warning[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("warning" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ThresholdCrossingReportMetric) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThresholdCrossingReportMetric) UnmarshalBinary(b []byte) error {
	var res ThresholdCrossingReportMetric
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
