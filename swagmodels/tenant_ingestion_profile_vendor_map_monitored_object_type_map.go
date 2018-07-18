// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// TenantIngestionProfileVendorMapMonitoredObjectTypeMap tenant ingestion profile vendor map monitored object type map
// swagger:model TenantIngestionProfileVendorMapMonitoredObjectTypeMap
type TenantIngestionProfileVendorMapMonitoredObjectTypeMap struct {

	// monitored object type map
	MonitoredObjectTypeMap map[string]TenantIngestionProfileVendorMapMonitoredObjectTypeMapMetricMap `json:"monitoredObjectTypeMap,omitempty"`
}

// Validate validates this tenant ingestion profile vendor map monitored object type map
func (m *TenantIngestionProfileVendorMapMonitoredObjectTypeMap) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMonitoredObjectTypeMap(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TenantIngestionProfileVendorMapMonitoredObjectTypeMap) validateMonitoredObjectTypeMap(formats strfmt.Registry) error {

	if swag.IsZero(m.MonitoredObjectTypeMap) { // not required
		return nil
	}

	for k := range m.MonitoredObjectTypeMap {

		if swag.IsZero(m.MonitoredObjectTypeMap[k]) { // not required
			continue
		}

		if val, ok := m.MonitoredObjectTypeMap[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *TenantIngestionProfileVendorMapMonitoredObjectTypeMap) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantIngestionProfileVendorMapMonitoredObjectTypeMap) UnmarshalBinary(b []byte) error {
	var res TenantIngestionProfileVendorMapMonitoredObjectTypeMap
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
