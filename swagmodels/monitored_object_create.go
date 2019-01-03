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

// MonitoredObjectCreate monitored object create
// swagger:model MonitoredObjectCreate
type MonitoredObjectCreate struct {

	// attributes
	// Required: true
	Attributes *MonitoredObjectCreateAttributes `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// type
	// Required: true
	// Enum: [monitoredObjects]
	Type *string `json:"type"`
}

// Validate validates this monitored object create
func (m *MonitoredObjectCreate) Validate(formats strfmt.Registry) error {
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

func (m *MonitoredObjectCreate) validateAttributes(formats strfmt.Registry) error {

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

var monitoredObjectCreateTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["monitoredObjects"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		monitoredObjectCreateTypeTypePropEnum = append(monitoredObjectCreateTypeTypePropEnum, v)
	}
}

const (

	// MonitoredObjectCreateTypeMonitoredObjects captures enum value "monitoredObjects"
	MonitoredObjectCreateTypeMonitoredObjects string = "monitoredObjects"
)

// prop value enum
func (m *MonitoredObjectCreate) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, monitoredObjectCreateTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *MonitoredObjectCreate) validateType(formats strfmt.Registry) error {

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
func (m *MonitoredObjectCreate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MonitoredObjectCreate) UnmarshalBinary(b []byte) error {
	var res MonitoredObjectCreate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// MonitoredObjectCreateAttributes monitored object create attributes
// swagger:model MonitoredObjectCreateAttributes
type MonitoredObjectCreateAttributes struct {

	// Name of the origin of the Monitored Object
	ActuatorName string `json:"actuatorName,omitempty"`

	// Type of the origin of the Monitored Object
	// Enum: [unknown accedian-nid accedian-vnid]
	ActuatorType string `json:"actuatorType,omitempty"`

	// Attributes added to a Monitored Object that help identify the Mlnitored Object as well as provide flitering/grouping properties
	Meta map[string]string `json:"meta,omitempty"`

	// Unique identifier of the Monitored Object in Datahub
	// Required: true
	ObjectID *string `json:"objectId"`

	// Name of the Monitored Object
	ObjectName string `json:"objectName,omitempty"`

	// Type of the Monitored Object
	// Enum: [unknown flowmeter twamp-pe twamp-sf twamp-sl]
	ObjectType string `json:"objectType,omitempty"`

	// Name of the target of the Monitored Object
	ReflectorName string `json:"reflectorName,omitempty"`

	// Type of the target of the Monitored Object
	// Enum: [unknown accedian-nid accedian-vnid]
	ReflectorType string `json:"reflectorType,omitempty"`
}

// Validate validates this monitored object create attributes
func (m *MonitoredObjectCreateAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActuatorType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateObjectID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateObjectType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReflectorType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var monitoredObjectCreateAttributesTypeActuatorTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","accedian-nid","accedian-vnid"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		monitoredObjectCreateAttributesTypeActuatorTypePropEnum = append(monitoredObjectCreateAttributesTypeActuatorTypePropEnum, v)
	}
}

const (

	// MonitoredObjectCreateAttributesActuatorTypeUnknown captures enum value "unknown"
	MonitoredObjectCreateAttributesActuatorTypeUnknown string = "unknown"

	// MonitoredObjectCreateAttributesActuatorTypeAccedianNid captures enum value "accedian-nid"
	MonitoredObjectCreateAttributesActuatorTypeAccedianNid string = "accedian-nid"

	// MonitoredObjectCreateAttributesActuatorTypeAccedianVnid captures enum value "accedian-vnid"
	MonitoredObjectCreateAttributesActuatorTypeAccedianVnid string = "accedian-vnid"
)

// prop value enum
func (m *MonitoredObjectCreateAttributes) validateActuatorTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, monitoredObjectCreateAttributesTypeActuatorTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *MonitoredObjectCreateAttributes) validateActuatorType(formats strfmt.Registry) error {

	if swag.IsZero(m.ActuatorType) { // not required
		return nil
	}

	// value enum
	if err := m.validateActuatorTypeEnum("attributes"+"."+"actuatorType", "body", m.ActuatorType); err != nil {
		return err
	}

	return nil
}

func (m *MonitoredObjectCreateAttributes) validateObjectID(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"objectId", "body", m.ObjectID); err != nil {
		return err
	}

	return nil
}

var monitoredObjectCreateAttributesTypeObjectTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","flowmeter","twamp-pe","twamp-sf","twamp-sl", "cisco-interface", "cisco-node-summary"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		monitoredObjectCreateAttributesTypeObjectTypePropEnum = append(monitoredObjectCreateAttributesTypeObjectTypePropEnum, v)
	}
}

const (

	// MonitoredObjectCreateAttributesObjectTypeUnknown captures enum value "unknown"
	MonitoredObjectCreateAttributesObjectTypeUnknown string = "unknown"

	// MonitoredObjectCreateAttributesObjectTypeFlowmeter captures enum value "flowmeter"
	MonitoredObjectCreateAttributesObjectTypeFlowmeter string = "flowmeter"

	// MonitoredObjectCreateAttributesObjectTypeTwampPe captures enum value "twamp-pe"
	MonitoredObjectCreateAttributesObjectTypeTwampPe string = "twamp-pe"

	// MonitoredObjectCreateAttributesObjectTypeTwampSf captures enum value "twamp-sf"
	MonitoredObjectCreateAttributesObjectTypeTwampSf string = "twamp-sf"

	// MonitoredObjectCreateAttributesObjectTypeTwampSl captures enum value "twamp-sl"
	MonitoredObjectCreateAttributesObjectTypeTwampSl string = "twamp-sl"
)

// prop value enum
func (m *MonitoredObjectCreateAttributes) validateObjectTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, monitoredObjectCreateAttributesTypeObjectTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *MonitoredObjectCreateAttributes) validateObjectType(formats strfmt.Registry) error {

	if swag.IsZero(m.ObjectType) { // not required
		return nil
	}

	// value enum
	if err := m.validateObjectTypeEnum("attributes"+"."+"objectType", "body", m.ObjectType); err != nil {
		return err
	}

	return nil
}

var monitoredObjectCreateAttributesTypeReflectorTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","accedian-nid","accedian-vnid"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		monitoredObjectCreateAttributesTypeReflectorTypePropEnum = append(monitoredObjectCreateAttributesTypeReflectorTypePropEnum, v)
	}
}

const (

	// MonitoredObjectCreateAttributesReflectorTypeUnknown captures enum value "unknown"
	MonitoredObjectCreateAttributesReflectorTypeUnknown string = "unknown"

	// MonitoredObjectCreateAttributesReflectorTypeAccedianNid captures enum value "accedian-nid"
	MonitoredObjectCreateAttributesReflectorTypeAccedianNid string = "accedian-nid"

	// MonitoredObjectCreateAttributesReflectorTypeAccedianVnid captures enum value "accedian-vnid"
	MonitoredObjectCreateAttributesReflectorTypeAccedianVnid string = "accedian-vnid"
)

// prop value enum
func (m *MonitoredObjectCreateAttributes) validateReflectorTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, monitoredObjectCreateAttributesTypeReflectorTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *MonitoredObjectCreateAttributes) validateReflectorType(formats strfmt.Registry) error {

	if swag.IsZero(m.ReflectorType) { // not required
		return nil
	}

	// value enum
	if err := m.validateReflectorTypeEnum("attributes"+"."+"reflectorType", "body", m.ReflectorType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MonitoredObjectCreateAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MonitoredObjectCreateAttributes) UnmarshalBinary(b []byte) error {
	var res MonitoredObjectCreateAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
