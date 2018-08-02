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

// ConnectorInstanceUpdateRequest connector instance update request
// swagger:model ConnectorInstanceUpdateRequest
type ConnectorInstanceUpdateRequest struct {

	// data
	// Required: true
	Data *ConnectorInstanceUpdateRequestData `json:"data"`
}

// Validate validates this connector instance update request
func (m *ConnectorInstanceUpdateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConnectorInstanceUpdateRequest) validateData(formats strfmt.Registry) error {

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
func (m *ConnectorInstanceUpdateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConnectorInstanceUpdateRequest) UnmarshalBinary(b []byte) error {
	var res ConnectorInstanceUpdateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ConnectorInstanceUpdateRequestData connector instance update request data
// swagger:model ConnectorInstanceUpdateRequestData
type ConnectorInstanceUpdateRequestData struct {

	// attributes
	// Required: true
	Attributes *ConnectorInstanceUpdateRequestDataAttributes `json:"attributes"`

	// id
	// Required: true
	ID *string `json:"id"`

	// type
	// Required: true
	// Enum: [connectorInstances]
	Type *string `json:"type"`
}

// Validate validates this connector instance update request data
func (m *ConnectorInstanceUpdateRequestData) Validate(formats strfmt.Registry) error {
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

func (m *ConnectorInstanceUpdateRequestData) validateAttributes(formats strfmt.Registry) error {

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

func (m *ConnectorInstanceUpdateRequestData) validateID(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

var connectorInstanceUpdateRequestDataTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["connectorInstances"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		connectorInstanceUpdateRequestDataTypeTypePropEnum = append(connectorInstanceUpdateRequestDataTypeTypePropEnum, v)
	}
}

const (

	// ConnectorInstanceUpdateRequestDataTypeConnectorInstances captures enum value "connectorInstances"
	ConnectorInstanceUpdateRequestDataTypeConnectorInstances string = "connectorInstances"
)

// prop value enum
func (m *ConnectorInstanceUpdateRequestData) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, connectorInstanceUpdateRequestDataTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ConnectorInstanceUpdateRequestData) validateType(formats strfmt.Registry) error {

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
func (m *ConnectorInstanceUpdateRequestData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConnectorInstanceUpdateRequestData) UnmarshalBinary(b []byte) error {
	var res ConnectorInstanceUpdateRequestData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ConnectorInstanceUpdateRequestDataAttributes connector instance update request data attributes
// swagger:model ConnectorInstanceUpdateRequestDataAttributes
type ConnectorInstanceUpdateRequestDataAttributes struct {

	// rev
	// Required: true
	Rev *string `json:"_rev"`

	// created timestamp
	CreatedTimestamp int64 `json:"createdTimestamp,omitempty"`

	// hostname
	Hostname string `json:"hostname,omitempty"`

	// last modified timestamp
	LastModifiedTimestamp int64 `json:"lastModifiedTimestamp,omitempty"`

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this connector instance update request data attributes
func (m *ConnectorInstanceUpdateRequestDataAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRev(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConnectorInstanceUpdateRequestDataAttributes) validateRev(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes"+"."+"_rev", "body", m.Rev); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ConnectorInstanceUpdateRequestDataAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConnectorInstanceUpdateRequestDataAttributes) UnmarshalBinary(b []byte) error {
	var res ConnectorInstanceUpdateRequestDataAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
