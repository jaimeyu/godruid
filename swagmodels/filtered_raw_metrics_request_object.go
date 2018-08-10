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

// FilteredRawMetricsRequestObject filtered raw metrics request object
// swagger:model FilteredRawMetricsRequestObject
type FilteredRawMetricsRequestObject struct {

	// directions
	// Required: true
	Directions []string `json:"directions"`

	// the granularity for timeseries in ISO-8601 duration format, or ALL
	Granularity string `json:"granularity,omitempty"`

	// ISO-8601 interval
	// Required: true
	Interval *string `json:"interval"`

	// set of meta keys and list of values for the purposes of filtering
	// Required: true
	Meta map[string][]string `json:"meta"`

	// metrics
	// Required: true
	Metrics []string `json:"metrics"`

	// the type of monitored object
	// Required: true
	ObjectType *string `json:"objectType"`

	// the tenant identifier
	// Required: true
	TenantID *string `json:"tenantId"`

	// query timeout in milliseconds
	Timeout int64 `json:"timeout,omitempty"`
}

// Validate validates this filtered raw metrics request object
func (m *FilteredRawMetricsRequestObject) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDirections(formats); err != nil {
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

	if err := m.validateObjectType(formats); err != nil {
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

func (m *FilteredRawMetricsRequestObject) validateDirections(formats strfmt.Registry) error {

	if err := validate.Required("directions", "body", m.Directions); err != nil {
		return err
	}

	return nil
}

func (m *FilteredRawMetricsRequestObject) validateInterval(formats strfmt.Registry) error {

	if err := validate.Required("interval", "body", m.Interval); err != nil {
		return err
	}

	return nil
}

func (m *FilteredRawMetricsRequestObject) validateMeta(formats strfmt.Registry) error {

	return nil
}

func (m *FilteredRawMetricsRequestObject) validateMetrics(formats strfmt.Registry) error {

	if err := validate.Required("metrics", "body", m.Metrics); err != nil {
		return err
	}

	return nil
}

func (m *FilteredRawMetricsRequestObject) validateObjectType(formats strfmt.Registry) error {

	if err := validate.Required("objectType", "body", m.ObjectType); err != nil {
		return err
	}

	return nil
}

func (m *FilteredRawMetricsRequestObject) validateTenantID(formats strfmt.Registry) error {

	if err := validate.Required("tenantId", "body", m.TenantID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *FilteredRawMetricsRequestObject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FilteredRawMetricsRequestObject) UnmarshalBinary(b []byte) error {
	var res FilteredRawMetricsRequestObject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
