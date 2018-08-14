// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DashboardUpdateRequest dashboard update request
// swagger:model DashboardUpdateRequest
type DashboardUpdateRequest struct {

	// data
	// Required: true
	Data *DashboardUpdateRequestData `json:"data"`
}

// Validate validates this dashboard update request
func (m *DashboardUpdateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DashboardUpdateRequest) validateData(formats strfmt.Registry) error {

	if err := validate.Required("data", "body", m.Data); err != nil {
		return err
	}

	if m.Data != nil {
		if err := m.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DashboardUpdateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DashboardUpdateRequest) UnmarshalBinary(b []byte) error {
	var res DashboardUpdateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DashboardUpdateRequestData dashboard update request data
// swagger:model DashboardUpdateRequestData
type DashboardUpdateRequestData struct {

	// attributes
	// Required: true
	Attributes *DashboardUpdateRequestDataAttributes `json:"attributes"`

	// id
	// Required: true
	ID *string `json:"id"`

	// relationships
	Relationships *DashboardRelationships `json:"relationships,omitempty"`

	// type
	// Required: true
	// Enum: [dashboards]
	Type *string `json:"type"`
}

// Validate validates this dashboard update request data
func (m *DashboardUpdateRequestData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRelationships(formats); err != nil {
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

func (m *DashboardUpdateRequestData) validateAttributes(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes", "body", m.Attributes); err != nil {
		return err
	}

	if m.Attributes != nil {
		if err := m.Attributes.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data" + "." + "attributes")
			}
			return err
		}
	}

	return nil
}

func (m *DashboardUpdateRequestData) validateID(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *DashboardUpdateRequestData) validateRelationships(formats strfmt.Registry) error {

	if swag.IsZero(m.Relationships) { // not required
		return nil
	}

	if m.Relationships != nil {
		if err := m.Relationships.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data" + "." + "relationships")
			}
			return err
		}
	}

	return nil
}

var dashboardUpdateRequestDataTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["dashboards"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		dashboardUpdateRequestDataTypeTypePropEnum = append(dashboardUpdateRequestDataTypeTypePropEnum, v)
	}
}

const (

	// DashboardUpdateRequestDataTypeDashboards captures enum value "dashboards"
	DashboardUpdateRequestDataTypeDashboards string = "dashboards"
)

// prop value enum
func (m *DashboardUpdateRequestData) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, dashboardUpdateRequestDataTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *DashboardUpdateRequestData) validateType(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("data"+"."+"type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DashboardUpdateRequestData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DashboardUpdateRequestData) UnmarshalBinary(b []byte) error {
	var res DashboardUpdateRequestData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DashboardUpdateRequestDataAttributes dashboard update request data attributes
// swagger:model DashboardUpdateRequestDataAttributes
type DashboardUpdateRequestDataAttributes struct {

	// rev
	// Required: true
	Rev *string `json:"_rev"`

	// card positions
	CardPositions CardPositions `json:"cardPositions,omitempty"`

	// category
	Category string `json:"category,omitempty"`

	// metadata filters
	MetadataFilters []*MetadataFilter `json:"metadataFilters"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this dashboard update request data attributes
func (m *DashboardUpdateRequestDataAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRev(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCardPositions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetadataFilters(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DashboardUpdateRequestDataAttributes) validateRev(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes"+"."+"_rev", "body", m.Rev); err != nil {
		return err
	}

	return nil
}

func (m *DashboardUpdateRequestDataAttributes) validateCardPositions(formats strfmt.Registry) error {

	if swag.IsZero(m.CardPositions) { // not required
		return nil
	}

	if err := m.CardPositions.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("data" + "." + "attributes" + "." + "cardPositions")
		}
		return err
	}

	return nil
}

func (m *DashboardUpdateRequestDataAttributes) validateMetadataFilters(formats strfmt.Registry) error {

	if swag.IsZero(m.MetadataFilters) { // not required
		return nil
	}

	for i := 0; i < len(m.MetadataFilters); i++ {
		if swag.IsZero(m.MetadataFilters[i]) { // not required
			continue
		}

		if m.MetadataFilters[i] != nil {
			if err := m.MetadataFilters[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("data" + "." + "attributes" + "." + "metadataFilters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *DashboardUpdateRequestDataAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DashboardUpdateRequestDataAttributes) UnmarshalBinary(b []byte) error {
	var res DashboardUpdateRequestDataAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
