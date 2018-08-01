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

// TenantCreationObject tenant creation object
// swagger:model TenantCreationObject
type TenantCreationObject struct {

	// attributes
	// Required: true
	Attributes *TenantCreationObjectAttributes `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Required: true
	// Enum: [tenants]
	Type *string `json:"type"`
}

// Validate validates this tenant creation object
func (m *TenantCreationObject) Validate(formats strfmt.Registry) error {
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

func (m *TenantCreationObject) validateAttributes(formats strfmt.Registry) error {

	if err := validate.Required("attributes", "body", m.Attributes); err != nil {
		return err
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

var tenantCreationObjectTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["tenants"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		tenantCreationObjectTypeTypePropEnum = append(tenantCreationObjectTypeTypePropEnum, v)
	}
}

const (

	// TenantCreationObjectTypeTenants captures enum value "tenants"
	TenantCreationObjectTypeTenants string = "tenants"
)

// prop value enum
func (m *TenantCreationObject) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, tenantCreationObjectTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *TenantCreationObject) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TenantCreationObject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantCreationObject) UnmarshalBinary(b []byte) error {
	var res TenantCreationObject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TenantCreationObjectAttributes tenant creation object attributes
// swagger:model TenantCreationObjectAttributes
type TenantCreationObjectAttributes struct {

	// The name of the Tenant
	// Required: true
	Name *string `json:"name"`

	// The subdomain used in the URL for accessing the Tenant's portal in Datahub
	// Required: true
	URLSubdomain *string `json:"urlSubdomain"`
}

// Validate validates this tenant creation object attributes
func (m *TenantCreationObjectAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURLSubdomain(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TenantCreationObjectAttributes) validateName(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *TenantCreationObjectAttributes) validateURLSubdomain(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"urlSubdomain", "body", m.URLSubdomain); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TenantCreationObjectAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantCreationObjectAttributes) UnmarshalBinary(b []byte) error {
	var res TenantCreationObjectAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
