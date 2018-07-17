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

// DataCleaningProfileCreateRequest data cleaning profile create request
// swagger:model DataCleaningProfileCreateRequest
type DataCleaningProfileCreateRequest struct {

	// attributes
	// Required: true
	Attributes *DataCleaningProfileCreateRequestAttr `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this data cleaning profile create request
func (m *DataCleaningProfileCreateRequest) Validate(formats strfmt.Registry) error {
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

func (m *DataCleaningProfileCreateRequest) validateAttributes(formats strfmt.Registry) error {

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

var dataCleaningProfileCreateRequestTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["dataCleaningProfiles"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		dataCleaningProfileCreateRequestTypeTypePropEnum = append(dataCleaningProfileCreateRequestTypeTypePropEnum, v)
	}
}

const (

	// DataCleaningProfileCreateRequestTypeDataCleaningProfiles captures enum value "dataCleaningProfiles"
	DataCleaningProfileCreateRequestTypeDataCleaningProfiles string = "dataCleaningProfiles"
)

// prop value enum
func (m *DataCleaningProfileCreateRequest) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, dataCleaningProfileCreateRequestTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *DataCleaningProfileCreateRequest) validateType(formats strfmt.Registry) error {

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
func (m *DataCleaningProfileCreateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataCleaningProfileCreateRequest) UnmarshalBinary(b []byte) error {
	var res DataCleaningProfileCreateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
