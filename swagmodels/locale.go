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

// Locale locale
// swagger:model Locale
type Locale struct {

	// attributes
	// Required: true
	Attributes *LocaleAttributes `json:"attributes"`

	// id
	// Required: true
	ID *string `json:"id"`

	// type
	// Required: true
	// Enum: [locales]
	Type *string `json:"type"`
}

// Validate validates this locale
func (m *Locale) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
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

func (m *Locale) validateAttributes(formats strfmt.Registry) error {

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

func (m *Locale) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

var localeTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["locales"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		localeTypeTypePropEnum = append(localeTypeTypePropEnum, v)
	}
}

const (

	// LocaleTypeLocales captures enum value "locales"
	LocaleTypeLocales string = "locales"
)

// prop value enum
func (m *Locale) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, localeTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Locale) validateType(formats strfmt.Registry) error {

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
func (m *Locale) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Locale) UnmarshalBinary(b []byte) error {
	var res Locale
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// LocaleAttributes locale attributes
// swagger:model LocaleAttributes
type LocaleAttributes struct {

	// Value used to ensure updates to this object are handled in order.
	// Required: true
	Rev *string `json:"_rev"`

	// Time since epoch at which this object was instantiated.
	// Required: true
	CreatedTimestamp *int64 `json:"createdTimestamp"`

	// datatype
	// Required: true
	Datatype *string `json:"datatype"`

	// The short-form code for the internationalization region
	// Required: true
	Intl *string `json:"intl"`

	// Time since epoch at which this object was last altered.
	// Required: true
	LastModifiedTimestamp *int64 `json:"lastModifiedTimestamp"`

	// moment
	// Required: true
	Moment *string `json:"moment"`

	// tenant Id
	TenantID string `json:"tenantId,omitempty"`

	// Timezone used to coordinate timestamps for the specified region
	// Required: true
	Timezone *string `json:"timezone"`
}

// Validate validates this locale attributes
func (m *LocaleAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRev(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDatatype(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIntl(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastModifiedTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMoment(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimezone(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LocaleAttributes) validateRev(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"_rev", "body", m.Rev); err != nil {
		return err
	}

	return nil
}

func (m *LocaleAttributes) validateCreatedTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"createdTimestamp", "body", m.CreatedTimestamp); err != nil {
		return err
	}

	return nil
}

func (m *LocaleAttributes) validateDatatype(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"datatype", "body", m.Datatype); err != nil {
		return err
	}

	return nil
}

func (m *LocaleAttributes) validateIntl(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"intl", "body", m.Intl); err != nil {
		return err
	}

	return nil
}

func (m *LocaleAttributes) validateLastModifiedTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"lastModifiedTimestamp", "body", m.LastModifiedTimestamp); err != nil {
		return err
	}

	return nil
}

func (m *LocaleAttributes) validateMoment(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"moment", "body", m.Moment); err != nil {
		return err
	}

	return nil
}

func (m *LocaleAttributes) validateTimezone(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"timezone", "body", m.Timezone); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *LocaleAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LocaleAttributes) UnmarshalBinary(b []byte) error {
	var res LocaleAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
