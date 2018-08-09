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

// ThresholdCrossingByMOTopNAPIRequestObject threshold crossing by m o top n API request object
// swagger:model ThresholdCrossingByMOTopNAPIRequestObject
type ThresholdCrossingByMOTopNAPIRequestObject struct {

	// the granularity for timeseries in ISO-8601 duration format, or ALL
	Granularity string `json:"granularity,omitempty"`

	// ISO-8601 interval
	// Required: true
	Interval *string `json:"interval"`

	// set of domains identifiers to use for filtering
	Meta map[string][]string `json:"meta,omitempty"`

	// the metric to be used for the top N query
	Metric *MetricIdentifierObject `json:"metric,omitempty"`

	// query timeout in milliseconds
	NumResults int32 `json:"numResults,omitempty"`

	// the tenant identifier
	// Required: true
	TenantID *string `json:"tenantId"`

	// ID of the threshold profile that is used to select metrics and events
	ThresholdProfileID string `json:"thresholdProfileId,omitempty"`

	// query timeout in milliseconds
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this threshold crossing by m o top n API request object
func (m *ThresholdCrossingByMOTopNAPIRequestObject) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInterval(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetric(formats); err != nil {
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

func (m *ThresholdCrossingByMOTopNAPIRequestObject) validateInterval(formats strfmt.Registry) error {

	if err := validate.Required("interval", "body", m.Interval); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdCrossingByMOTopNAPIRequestObject) validateMetric(formats strfmt.Registry) error {

	if swag.IsZero(m.Metric) { // not required
		return nil
	}

	if m.Metric != nil {
		if err := m.Metric.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("metric")
			}
			return err
		}
	}

	return nil
}

func (m *ThresholdCrossingByMOTopNAPIRequestObject) validateTenantID(formats strfmt.Registry) error {

	if err := validate.Required("tenantId", "body", m.TenantID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ThresholdCrossingByMOTopNAPIRequestObject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThresholdCrossingByMOTopNAPIRequestObject) UnmarshalBinary(b []byte) error {
	var res ThresholdCrossingByMOTopNAPIRequestObject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
