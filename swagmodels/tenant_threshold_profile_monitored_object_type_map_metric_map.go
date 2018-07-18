// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// TenantThresholdProfileMonitoredObjectTypeMapMetricMap tenant threshold profile monitored object type map metric map
// swagger:model TenantThresholdProfileMonitoredObjectTypeMapMetricMap
type TenantThresholdProfileMonitoredObjectTypeMapMetricMap struct {

	// metric map
	MetricMap TenantThresholdProfileMonitoredObjectTypeMapMetricMapMetricMap `json:"metricMap,omitempty"`
}

// Validate validates this tenant threshold profile monitored object type map metric map
func (m *TenantThresholdProfileMonitoredObjectTypeMapMetricMap) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMetricMap(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TenantThresholdProfileMonitoredObjectTypeMapMetricMap) validateMetricMap(formats strfmt.Registry) error {

	if swag.IsZero(m.MetricMap) { // not required
		return nil
	}

	if err := m.MetricMap.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("metricMap")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TenantThresholdProfileMonitoredObjectTypeMapMetricMap) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantThresholdProfileMonitoredObjectTypeMapMetricMap) UnmarshalBinary(b []byte) error {
	var res TenantThresholdProfileMonitoredObjectTypeMapMetricMap
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
