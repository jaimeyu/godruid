// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// JSONAPIAdminUserList Json Api admin user list
// swagger:model JsonApiAdminUserList
type JSONAPIAdminUserList struct {

	// data
	Data JSONAPIAdminUserListData `json:"data"`
}

// Validate validates this Json Api admin user list
func (m *JSONAPIAdminUserList) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *JSONAPIAdminUserList) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JSONAPIAdminUserList) UnmarshalBinary(b []byte) error {
	var res JSONAPIAdminUserList
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
