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

// ConnectorInstanceCreateRequest connector instance create request
// swagger:model ConnectorInstanceCreateRequest
type ConnectorInstanceCreateRequest struct {

	// data
	// Required: true
	Data *ConnectorInstanceCreateRequestData `json:"data"`
}

// Validate validates this connector instance create request
func (m *ConnectorInstanceCreateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConnectorInstanceCreateRequest) validateData(formats strfmt.Registry) error {

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
func (m *ConnectorInstanceCreateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConnectorInstanceCreateRequest) UnmarshalBinary(b []byte) error {
	var res ConnectorInstanceCreateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ConnectorInstanceCreateRequestData connector instance create request data
// swagger:model ConnectorInstanceCreateRequestData
type ConnectorInstanceCreateRequestData struct {

	// attributes
	// Required: true
	Attributes *ConnectorInstanceCreateRequestDataAttributes `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Required: true
	// Enum: [connectorInstances]
	Type *string `json:"type"`
}

// Validate validates this connector instance create request data
func (m *ConnectorInstanceCreateRequestData) Validate(formats strfmt.Registry) error {
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

func (m *ConnectorInstanceCreateRequestData) validateAttributes(formats strfmt.Registry) error {

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

var connectorInstanceCreateRequestDataTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["connectorInstances"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		connectorInstanceCreateRequestDataTypeTypePropEnum = append(connectorInstanceCreateRequestDataTypeTypePropEnum, v)
	}
}

const (

	// ConnectorInstanceCreateRequestDataTypeConnectorInstances captures enum value "connectorInstances"
	ConnectorInstanceCreateRequestDataTypeConnectorInstances string = "connectorInstances"
)

// prop value enum
func (m *ConnectorInstanceCreateRequestData) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, connectorInstanceCreateRequestDataTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ConnectorInstanceCreateRequestData) validateType(formats strfmt.Registry) error {

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
func (m *ConnectorInstanceCreateRequestData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConnectorInstanceCreateRequestData) UnmarshalBinary(b []byte) error {
	var res ConnectorInstanceCreateRequestData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ConnectorInstanceCreateRequestDataAttributes connector instance create request data attributes
// swagger:model ConnectorInstanceCreateRequestDataAttributes
type ConnectorInstanceCreateRequestDataAttributes struct {

	// hostname
	// Required: true
	Hostname *string `json:"hostname"`

	// status
	// Required: true
	Status *string `json:"status"`
}

// Validate validates this connector instance create request data attributes
func (m *ConnectorInstanceCreateRequestDataAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHostname(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConnectorInstanceCreateRequestDataAttributes) validateHostname(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes"+"."+"hostname", "body", m.Hostname); err != nil {
		return err
	}

	return nil
}

func (m *ConnectorInstanceCreateRequestDataAttributes) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes"+"."+"status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ConnectorInstanceCreateRequestDataAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConnectorInstanceCreateRequestDataAttributes) UnmarshalBinary(b []byte) error {
	var res ConnectorInstanceCreateRequestDataAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
