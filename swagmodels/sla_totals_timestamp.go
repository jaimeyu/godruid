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

// SLATotalsTimestamp SLA totals timestamp
// swagger:model SLATotalsTimestamp
type SLATotalsTimestamp struct {

	// timestamp
	// Required: true
	Timestamp *string `json:"timestamp"`

	// total duration
	// Required: true
	TotalDuration *int64 `json:"totalDuration"`

	// total violation count
	// Required: true
	TotalViolationCount *int64 `json:"totalViolationCount"`

	// total violation duration
	// Required: true
	TotalViolationDuration *int64 `json:"totalViolationDuration"`
}

// Validate validates this SLA totals timestamp
func (m *SLATotalsTimestamp) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTotalDuration(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTotalViolationCount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTotalViolationDuration(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SLATotalsTimestamp) validateTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("timestamp", "body", m.Timestamp); err != nil {
		return err
	}

	return nil
}

func (m *SLATotalsTimestamp) validateTotalDuration(formats strfmt.Registry) error {

	if err := validate.Required("totalDuration", "body", m.TotalDuration); err != nil {
		return err
	}

	return nil
}

func (m *SLATotalsTimestamp) validateTotalViolationCount(formats strfmt.Registry) error {

	if err := validate.Required("totalViolationCount", "body", m.TotalViolationCount); err != nil {
		return err
	}

	return nil
}

func (m *SLATotalsTimestamp) validateTotalViolationDuration(formats strfmt.Registry) error {

	if err := validate.Required("totalViolationDuration", "body", m.TotalViolationDuration); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SLATotalsTimestamp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SLATotalsTimestamp) UnmarshalBinary(b []byte) error {
	var res SLATotalsTimestamp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
