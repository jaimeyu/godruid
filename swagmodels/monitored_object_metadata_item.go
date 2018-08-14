// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// MonitoredObjectMetadataItem monitored object metadata item
// swagger:model MonitoredObjectMetadataItem
type MonitoredObjectMetadataItem struct {

	// key name
	KeyName string `json:"keyName,omitempty"`

	// metadata
	Metadata map[string]string `json:"metadata,omitempty"`

	// metadata key
	MetadataKey string `json:"metadataKey,omitempty"`
}

// Validate validates this monitored object metadata item
func (m *MonitoredObjectMetadataItem) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MonitoredObjectMetadataItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MonitoredObjectMetadataItem) UnmarshalBinary(b []byte) error {
	var res MonitoredObjectMetadataItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
