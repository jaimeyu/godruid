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

// TenantConnectorInstance tenant connector instance
// swagger:model TenantConnectorInstance
type TenantConnectorInstance struct {

	// attributes
	Attributes *TenantConnectorInstanceAttr `json:"attributes,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Enum: [connectorInstances]
	Type string `json:"type,omitempty"`
}

// Validate validates this tenant connector instance
func (m *TenantConnectorInstance) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
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

func (m *TenantConnectorInstance) validateAttributes(formats strfmt.Registry) error {

	if swag.IsZero(m.Attributes) { // not required
		return nil
	}

	if m.Attributes != nil {
		if err := m.Attributes.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("attributes")
			}
			return err
		}
	}

	return nil
}

var tenantConnectorInstanceTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["connectorInstances"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		tenantConnectorInstanceTypeTypePropEnum = append(tenantConnectorInstanceTypeTypePropEnum, v)
	}
}

const (

	// TenantConnectorInstanceTypeConnectorInstances captures enum value "connectorInstances"
	TenantConnectorInstanceTypeConnectorInstances string = "connectorInstances"
)

// prop value enum
func (m *TenantConnectorInstance) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, tenantConnectorInstanceTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *TenantConnectorInstance) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TenantConnectorInstance) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantConnectorInstance) UnmarshalBinary(b []byte) error {
	var res TenantConnectorInstance
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
