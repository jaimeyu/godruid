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

// CardCreateRequest Object used to create a Card in Datahub
// swagger:model CardCreateRequest
type CardCreateRequest struct {

	// data
	// Required: true
	Data *CardCreateRequestData `json:"data"`
}

// Validate validates this card create request
func (m *CardCreateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CardCreateRequest) validateData(formats strfmt.Registry) error {

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
func (m *CardCreateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CardCreateRequest) UnmarshalBinary(b []byte) error {
	var res CardCreateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// CardCreateRequestData card create request data
// swagger:model CardCreateRequestData
type CardCreateRequestData struct {

	// attributes
	// Required: true
	Attributes *CardCreateRequestDataAttributes `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Required: true
	// Enum: [cards]
	Type *string `json:"type"`
}

// Validate validates this card create request data
func (m *CardCreateRequestData) Validate(formats strfmt.Registry) error {
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

func (m *CardCreateRequestData) validateAttributes(formats strfmt.Registry) error {

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

var cardCreateRequestDataTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["cards"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		cardCreateRequestDataTypeTypePropEnum = append(cardCreateRequestDataTypeTypePropEnum, v)
	}
}

const (

	// CardCreateRequestDataTypeCards captures enum value "cards"
	CardCreateRequestDataTypeCards string = "cards"
)

// prop value enum
func (m *CardCreateRequestData) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, cardCreateRequestDataTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CardCreateRequestData) validateType(formats strfmt.Registry) error {

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
func (m *CardCreateRequestData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CardCreateRequestData) UnmarshalBinary(b []byte) error {
	var res CardCreateRequestData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// CardCreateRequestDataAttributes card create request data attributes
// swagger:model CardCreateRequestDataAttributes
type CardCreateRequestDataAttributes struct {

	// description
	Description string `json:"description,omitempty"`

	// metrics
	Metrics []*CardMetric `json:"metrics"`

	// name
	// Required: true
	Name *string `json:"name"`

	// state
	// Required: true
	// Enum: [active pending]
	State *string `json:"state"`

	// visualization
	Visualization *CardVisualization `json:"visualization,omitempty"`
}

// Validate validates this card create request data attributes
func (m *CardCreateRequestDataAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMetrics(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateState(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVisualization(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CardCreateRequestDataAttributes) validateMetrics(formats strfmt.Registry) error {

	if swag.IsZero(m.Metrics) { // not required
		return nil
	}

	for i := 0; i < len(m.Metrics); i++ {
		if swag.IsZero(m.Metrics[i]) { // not required
			continue
		}

		if m.Metrics[i] != nil {
			if err := m.Metrics[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("data" + "." + "attributes" + "." + "metrics" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *CardCreateRequestDataAttributes) validateName(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

var cardCreateRequestDataAttributesTypeStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["active","pending"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		cardCreateRequestDataAttributesTypeStatePropEnum = append(cardCreateRequestDataAttributesTypeStatePropEnum, v)
	}
}

const (

	// CardCreateRequestDataAttributesStateActive captures enum value "active"
	CardCreateRequestDataAttributesStateActive string = "active"

	// CardCreateRequestDataAttributesStatePending captures enum value "pending"
	CardCreateRequestDataAttributesStatePending string = "pending"
)

// prop value enum
func (m *CardCreateRequestDataAttributes) validateStateEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, cardCreateRequestDataAttributesTypeStatePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CardCreateRequestDataAttributes) validateState(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes"+"."+"state", "body", m.State); err != nil {
		return err
	}

	// value enum
	if err := m.validateStateEnum("data"+"."+"attributes"+"."+"state", "body", *m.State); err != nil {
		return err
	}

	return nil
}

func (m *CardCreateRequestDataAttributes) validateVisualization(formats strfmt.Registry) error {

	if swag.IsZero(m.Visualization) { // not required
		return nil
	}

	if m.Visualization != nil {
		if err := m.Visualization.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data" + "." + "attributes" + "." + "visualization")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CardCreateRequestDataAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CardCreateRequestDataAttributes) UnmarshalBinary(b []byte) error {
	var res CardCreateRequestDataAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
