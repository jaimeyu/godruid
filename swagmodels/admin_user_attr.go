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

// AdminUserAttr admin user attr
// swagger:model AdminUserAttr
type AdminUserAttr struct {

	// id
	ID string `json:"_id,omitempty"`

	// rev
	Rev string `json:"_rev,omitempty"`

	// created timestamp
	CreatedTimestamp int64 `json:"createdTimestamp,omitempty"`

	// datatype
	Datatype string `json:"datatype,omitempty"`

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

	// user verified
	UserVerified bool `json:"userVerified,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this admin user attr
func (m *AdminUserAttr) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateState(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var adminUserAttrTypeStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["USER_UNKNOWN","INVITED","ACTIVE","SUSPENDED","PENDING_DELETE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		adminUserAttrTypeStatePropEnum = append(adminUserAttrTypeStatePropEnum, v)
	}
}

const (

	// AdminUserAttrStateUSERUNKNOWN captures enum value "USER_UNKNOWN"
	AdminUserAttrStateUSERUNKNOWN string = "USER_UNKNOWN"

	// AdminUserAttrStateINVITED captures enum value "INVITED"
	AdminUserAttrStateINVITED string = "INVITED"

	// AdminUserAttrStateACTIVE captures enum value "ACTIVE"
	AdminUserAttrStateACTIVE string = "ACTIVE"

	// AdminUserAttrStateSUSPENDED captures enum value "SUSPENDED"
	AdminUserAttrStateSUSPENDED string = "SUSPENDED"

	// AdminUserAttrStatePENDINGDELETE captures enum value "PENDING_DELETE"
	AdminUserAttrStatePENDINGDELETE string = "PENDING_DELETE"
)

// prop value enum
func (m *AdminUserAttr) validateStateEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, adminUserAttrTypeStatePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *AdminUserAttr) validateState(formats strfmt.Registry) error {

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
func (m *AdminUserAttr) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AdminUserAttr) UnmarshalBinary(b []byte) error {
	var res AdminUserAttr
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
