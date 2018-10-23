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

// JSONAPISLAReportRequest JSON API compliant wrapper for the SLA report query
// swagger:model JsonApiSLAReportRequest
type JSONAPISLAReportRequest struct {

	// data
	Data *JSONAPISLAReportRequestData `json:"data,omitempty"`
}

// Validate validates this Json Api SLA report request
func (m *JSONAPISLAReportRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JSONAPISLAReportRequest) validateData(formats strfmt.Registry) error {

	if swag.IsZero(m.Data) { // not required
		return nil
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
func (m *JSONAPISLAReportRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JSONAPISLAReportRequest) UnmarshalBinary(b []byte) error {
	var res JSONAPISLAReportRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// JSONAPISLAReportRequestData JSON API SLA report request data
// swagger:model JSONAPISLAReportRequestData
type JSONAPISLAReportRequestData struct {

	// attributes
	// Required: true
	Attributes *SLAReportConfig `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Required: true
	// Enum: [slaReports]
	Type *string `json:"type"`
}

// Validate validates this JSON API SLA report request data
func (m *JSONAPISLAReportRequestData) Validate(formats strfmt.Registry) error {
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

func (m *JSONAPISLAReportRequestData) validateAttributes(formats strfmt.Registry) error {

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

var jsonApiSlaReportRequestDataTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["slaReports"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		jsonApiSlaReportRequestDataTypeTypePropEnum = append(jsonApiSlaReportRequestDataTypeTypePropEnum, v)
	}
}

const (

	// JSONAPISLAReportRequestDataTypeSLAReports captures enum value "slaReports"
	JSONAPISLAReportRequestDataTypeSLAReports string = "slaReports"
)

// prop value enum
func (m *JSONAPISLAReportRequestData) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, jsonApiSlaReportRequestDataTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *JSONAPISLAReportRequestData) validateType(formats strfmt.Registry) error {

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
func (m *JSONAPISLAReportRequestData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JSONAPISLAReportRequestData) UnmarshalBinary(b []byte) error {
	var res JSONAPISLAReportRequestData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
