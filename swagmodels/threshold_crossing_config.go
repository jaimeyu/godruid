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

// ThresholdCrossingConfig The necessary request parameters for the metric api call
// swagger:model ThresholdCrossingConfig
type ThresholdCrossingConfig struct {

	// the granularity for timeseries in ISO-8601 duration format, or ALL
	Granularity string `json:"granularity,omitempty"`

	// Time boundary for the metrics under consideration using the ISO-8601 standard
	// Required: true
	Interval *string `json:"interval"`

	// An object that allows filtering on arbitrary metadata criteria and their values. Refer to the MetaFilter object for additional details
	Meta MetaFilter `json:"meta,omitempty"`

	// limits the results to include only metrics in the whitelist
	// Required: true
	Metrics []*MetricIdentifierFilter `json:"metrics"`

	// ID of the threshold profile that is used to select metrics and events
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`

	// query timeout in milliseconds
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this threshold crossing config
func (m *ThresholdCrossingConfig) Validate(formats strfmt.Registry) error {
	var res []error

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

func (m *ThresholdCrossingConfig) validateInterval(formats strfmt.Registry) error {

	if err := validate.Required("interval", "body", m.Interval); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdCrossingConfig) validateMeta(formats strfmt.Registry) error {

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

func (m *ThresholdCrossingConfig) validateMetrics(formats strfmt.Registry) error {

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
func (m *ThresholdCrossingConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThresholdCrossingConfig) UnmarshalBinary(b []byte) error {
	var res ThresholdCrossingConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
