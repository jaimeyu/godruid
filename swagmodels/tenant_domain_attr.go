// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// TenantDomainAttr tenant domain attr
// swagger:model TenantDomainAttr
type TenantDomainAttr struct {

	// id
	ID string `json:"_id,omitempty"`

	// rev
	Rev string `json:"_rev,omitempty"`

	// color
	Color string `json:"color,omitempty"`

	// created timestamp
	CreatedTimestamp int64 `json:"createdTimestamp,omitempty"`

	// datatype
	Datatype string `json:"datatype,omitempty"`

	// last modified timestamp
	LastModifiedTimestamp int64 `json:"lastModifiedTimestamp,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// tenant Id
	TenantID string `json:"tenantId,omitempty"`
}

// Validate validates this tenant domain attr
func (m *TenantDomainAttr) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TenantDomainAttr) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantDomainAttr) UnmarshalBinary(b []byte) error {
	var res TenantDomainAttr
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
