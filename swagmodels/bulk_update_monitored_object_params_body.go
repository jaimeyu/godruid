// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// BulkUpdateMonitoredObjectParamsBody bulk update monitored object params body
// swagger:model bulkUpdateMonitoredObjectParamsBody
type BulkUpdateMonitoredObjectParamsBody struct {

	// p0
	// Required: true
	P0 *CreateTenantMonitoredObjectRequest `json:"-"` // custom serializer

}

// UnmarshalJSON unmarshals this tuple type from a JSON array
func (m *BulkUpdateMonitoredObjectParamsBody) UnmarshalJSON(raw []byte) error {
	// stage 1, get the array but just the array
	var stage1 []json.RawMessage
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&stage1); err != nil {
		return err
	}

	// stage 2

	if len(stage1) > 0 {
		buf = bytes.NewBuffer(stage1[0])
		dec := json.NewDecoder(buf)
		dec.UseNumber()
		if err := dec.Decode(m.P0); err != nil {
			return err
		}

	}

	return nil
}

// MarshalJSON marshals this tuple type into a JSON array
func (m BulkUpdateMonitoredObjectParamsBody) MarshalJSON() ([]byte, error) {
	data := []interface{}{
		m.P0,
	}

	return json.Marshal(data)
}

// Validate validates this bulk update monitored object params body
func (m *BulkUpdateMonitoredObjectParamsBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateP0(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BulkUpdateMonitoredObjectParamsBody) validateP0(formats strfmt.Registry) error {

	if err := validate.Required("0", "body", m.P0); err != nil {
		return err
	}

	if m.P0 != nil {

		if err := m.P0.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("0")
			}
			return err
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *BulkUpdateMonitoredObjectParamsBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BulkUpdateMonitoredObjectParamsBody) UnmarshalBinary(b []byte) error {
	var res BulkUpdateMonitoredObjectParamsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
