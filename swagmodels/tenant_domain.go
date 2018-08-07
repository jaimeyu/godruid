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

// TenantDomain tenant domain
// swagger:model TenantDomain
type TenantDomain struct {

	// attributes
	Attributes *TenantDomainAttr `json:"attributes,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Enum: [domains]
	Type string `json:"type,omitempty"`
}

// Validate validates this tenant domain
func (m *TenantDomain) Validate(formats strfmt.Registry) error {
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

func (m *TenantDomain) validateAttributes(formats strfmt.Registry) error {

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

var tenantDomainTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["domains"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		tenantDomainTypeTypePropEnum = append(tenantDomainTypeTypePropEnum, v)
	}
}

const (

	// TenantDomainTypeDomains captures enum value "domains"
	TenantDomainTypeDomains string = "domains"
)

// prop value enum
func (m *TenantDomain) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, tenantDomainTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *TenantDomain) validateType(formats strfmt.Registry) error {

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
func (m *TenantDomain) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantDomain) UnmarshalBinary(b []byte) error {
	var res TenantDomain
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
