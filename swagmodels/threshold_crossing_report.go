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

// ThresholdCrossingReport threshold crossing report
// swagger:model ThresholdCrossingReport
type ThresholdCrossingReport struct {

	// metric
	Metric []*ThresholdCrossingReportMetric `json:"metric"`
}

// Validate validates this threshold crossing report
func (m *ThresholdCrossingReport) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMetric(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ThresholdCrossingReport) validateMetric(formats strfmt.Registry) error {

	if swag.IsZero(m.Metric) { // not required
		return nil
	}

	for i := 0; i < len(m.Metric); i++ {
		if swag.IsZero(m.Metric[i]) { // not required
			continue
		}

		if m.Metric[i] != nil {
			if err := m.Metric[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("metric" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ThresholdCrossingReport) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThresholdCrossingReport) UnmarshalBinary(b []byte) error {
	var res ThresholdCrossingReport
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
