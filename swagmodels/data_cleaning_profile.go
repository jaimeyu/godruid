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

// DataCleaningProfile data cleaning profile
// swagger:model DataCleaningProfile
type DataCleaningProfile struct {

	// attributes
	// Required: true
	Attributes *DataCleaningProfileAttributes `json:"attributes"`

	// id
	// Required: true
	ID *string `json:"id"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this data cleaning profile
func (m *DataCleaningProfile) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
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

func (m *DataCleaningProfile) validateAttributes(formats strfmt.Registry) error {

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

func (m *DataCleaningProfile) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

var dataCleaningProfileTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["dataCleaningProfiles"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		dataCleaningProfileTypeTypePropEnum = append(dataCleaningProfileTypeTypePropEnum, v)
	}
}

const (

	// DataCleaningProfileTypeDataCleaningProfiles captures enum value "dataCleaningProfiles"
	DataCleaningProfileTypeDataCleaningProfiles string = "dataCleaningProfiles"
)

// prop value enum
func (m *DataCleaningProfile) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, dataCleaningProfileTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *DataCleaningProfile) validateType(formats strfmt.Registry) error {

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
func (m *DataCleaningProfile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataCleaningProfile) UnmarshalBinary(b []byte) error {
	var res DataCleaningProfile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
