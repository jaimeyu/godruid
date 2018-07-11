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

// ValidTypes valid types
// swagger:model ValidTypes
type ValidTypes struct {

	// attributes
	Attributes *ValidTypesAttr `json:"attributes,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this valid types
func (m *ValidTypes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ValidTypes) validateAttributes(formats strfmt.Registry) error {

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

var validTypesTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["validTypes"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		validTypesTypeTypePropEnum = append(validTypesTypeTypePropEnum, v)
	}
}

const (

	// ValidTypesTypeValidTypes captures enum value "validTypes"
	ValidTypesTypeValidTypes string = "validTypes"
)

// prop value enum
func (m *ValidTypes) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, validTypesTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ValidTypes) validateType(formats strfmt.Registry) error {

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
func (m *ValidTypes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ValidTypes) UnmarshalBinary(b []byte) error {
	var res ValidTypes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
