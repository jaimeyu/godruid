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

// TenantUserAttr tenant user attr
// swagger:model TenantUserAttr
type TenantUserAttr struct {

	// id
	ID string `json:"_id,omitempty"`

	// rev
	Rev string `json:"_rev,omitempty"`

	// created timestamp
	CreatedTimestamp int64 `json:"createdTimestamp,omitempty"`

	// datatype
	Datatype string `json:"datatype,omitempty"`

	// domains
	Domains []string `json:"domains"`

	// last modified timestamp
	LastModifiedTimestamp int64 `json:"lastModifiedTimestamp,omitempty"`

	// onboarding token
	OnboardingToken string `json:"onboardingToken,omitempty"`

	// password
	Password string `json:"password,omitempty"`

	// send onboarding email
	SendOnboardingEmail bool `json:"sendOnboardingEmail,omitempty"`

	// state
	State string `json:"state,omitempty"`

	// tenant Id
	TenantID string `json:"tenantId,omitempty"`

	// user verified
	UserVerified bool `json:"userVerified,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this tenant user attr
func (m *TenantUserAttr) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDomains(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateState(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TenantUserAttr) validateDomains(formats strfmt.Registry) error {

	if swag.IsZero(m.Domains) { // not required
		return nil
	}

	return nil
}

var tenantUserAttrTypeStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["USER_UNKNOWN","INVITED","ACTIVE","SUSPENDED","PENDING_DELETE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		tenantUserAttrTypeStatePropEnum = append(tenantUserAttrTypeStatePropEnum, v)
	}
}

const (

	// TenantUserAttrStateUSERUNKNOWN captures enum value "USER_UNKNOWN"
	TenantUserAttrStateUSERUNKNOWN string = "USER_UNKNOWN"

	// TenantUserAttrStateINVITED captures enum value "INVITED"
	TenantUserAttrStateINVITED string = "INVITED"

	// TenantUserAttrStateACTIVE captures enum value "ACTIVE"
	TenantUserAttrStateACTIVE string = "ACTIVE"

	// TenantUserAttrStateSUSPENDED captures enum value "SUSPENDED"
	TenantUserAttrStateSUSPENDED string = "SUSPENDED"

	// TenantUserAttrStatePENDINGDELETE captures enum value "PENDING_DELETE"
	TenantUserAttrStatePENDINGDELETE string = "PENDING_DELETE"
)

// prop value enum
func (m *TenantUserAttr) validateStateEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, tenantUserAttrTypeStatePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *TenantUserAttr) validateState(formats strfmt.Registry) error {

	if swag.IsZero(m.State) { // not required
		return nil
	}

	// value enum
	if err := m.validateStateEnum("state", "body", m.State); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TenantUserAttr) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TenantUserAttr) UnmarshalBinary(b []byte) error {
	var res TenantUserAttr
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
