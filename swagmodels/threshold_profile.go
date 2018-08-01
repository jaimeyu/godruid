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

// ThresholdProfile threshold profile
// swagger:model ThresholdProfile
type ThresholdProfile struct {

	// attributes
	// Required: true
	Attributes *ThresholdProfileAttributes `json:"attributes"`

	// id
	// Required: true
	ID *string `json:"id"`

	// type
	// Required: true
	// Enum: [thresholdProfiles]
	Type *string `json:"type"`
}

// Validate validates this threshold profile
func (m *ThresholdProfile) Validate(formats strfmt.Registry) error {
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

func (m *ThresholdProfile) validateAttributes(formats strfmt.Registry) error {

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

func (m *ThresholdProfile) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

var thresholdProfileTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["thresholdProfiles"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		thresholdProfileTypeTypePropEnum = append(thresholdProfileTypeTypePropEnum, v)
	}
}

const (

	// ThresholdProfileTypeThresholdProfiles captures enum value "thresholdProfiles"
	ThresholdProfileTypeThresholdProfiles string = "thresholdProfiles"
)

// prop value enum
func (m *ThresholdProfile) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, thresholdProfileTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ThresholdProfile) validateType(formats strfmt.Registry) error {

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
func (m *ThresholdProfile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThresholdProfile) UnmarshalBinary(b []byte) error {
	var res ThresholdProfile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ThresholdProfileAttributes threshold profile attributes
// swagger:model ThresholdProfileAttributes
type ThresholdProfileAttributes struct {

	// id
	// Required: true
	ID *string `json:"_id"`

	// Value used to ensure updates to this object are handled in order.
	// Required: true
	Rev *string `json:"_rev"`

	// Time since epoch at which this object was instantiated.
	// Required: true
	CreatedTimestamp *int64 `json:"createdTimestamp"`

	// Name used to identify this type of record in Datahub
	// Required: true
	Datatype *string `json:"datatype"`

	// Time since epoch at which this object was last altered.
	// Required: true
	LastModifiedTimestamp *int64 `json:"lastModifiedTimestamp"`

	// Identifying name of a Threshold Profile
	// Required: true
	Name *string `json:"name"`

	// Unique identifier of the Tenant in Datahub
	// Required: true
	TenantID *string `json:"tenantId"`

	// threshold list
	ThresholdList ThresholdList `json:"thresholdList"`

	// Thresholds will be deprecated in the next API version. Please use the 'thresholdList' property instead
	// Required: true
	Thresholds *ThresholdsObject `json:"thresholds"`
}

// Validate validates this threshold profile attributes
func (m *ThresholdProfileAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRev(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDatatype(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastModifiedTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTenantID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateThresholdList(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateThresholds(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ThresholdProfileAttributes) validateID(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"_id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateRev(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"_rev", "body", m.Rev); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateCreatedTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"createdTimestamp", "body", m.CreatedTimestamp); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateDatatype(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"datatype", "body", m.Datatype); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateLastModifiedTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"lastModifiedTimestamp", "body", m.LastModifiedTimestamp); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateName(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateTenantID(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"tenantId", "body", m.TenantID); err != nil {
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateThresholdList(formats strfmt.Registry) error {

	if swag.IsZero(m.ThresholdList) { // not required
		return nil
	}

	if err := m.ThresholdList.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("attributes" + "." + "thresholdList")
		}
		return err
	}

	return nil
}

func (m *ThresholdProfileAttributes) validateThresholds(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"thresholds", "body", m.Thresholds); err != nil {
		return err
	}

	if m.Thresholds != nil {
		if err := m.Thresholds.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("attributes" + "." + "thresholds")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ThresholdProfileAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThresholdProfileAttributes) UnmarshalBinary(b []byte) error {
	var res ThresholdProfileAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
