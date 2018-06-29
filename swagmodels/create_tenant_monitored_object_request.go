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

// CreateTenantMonitoredObjectRequest create tenant monitored object request
// swagger:model CreateTenantMonitoredObjectRequest
type CreateTenantMonitoredObjectRequest struct {

	// id
	ID string `json:"_id,omitempty"`

	// rev
	Rev string `json:"_rev,omitempty"`

	// actuator name
	ActuatorName string `json:"actuatorName,omitempty"`

	// actuator type
	ActuatorType string `json:"actuatorType,omitempty"`

	// domain set
	DomainSet []string `json:"domainSet"`

	// object Id
	ObjectID string `json:"objectId,omitempty"`

	// object name
	ObjectName string `json:"objectName,omitempty"`

	// object type
	ObjectType string `json:"objectType,omitempty"`

	// reflector name
	ReflectorName string `json:"reflectorName,omitempty"`

	// reflector type
	ReflectorType string `json:"reflectorType,omitempty"`

	// tenant Id
	TenantID string `json:"tenantId,omitempty"`
}

// Validate validates this create tenant monitored object request
func (m *CreateTenantMonitoredObjectRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActuatorType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateDomainSet(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateObjectType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateReflectorType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var createTenantMonitoredObjectRequestTypeActuatorTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","accedian-nid","accedian-vnid"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		createTenantMonitoredObjectRequestTypeActuatorTypePropEnum = append(createTenantMonitoredObjectRequestTypeActuatorTypePropEnum, v)
	}
}

const (

	// CreateTenantMonitoredObjectRequestActuatorTypeUnknown captures enum value "unknown"
	CreateTenantMonitoredObjectRequestActuatorTypeUnknown string = "unknown"

	// CreateTenantMonitoredObjectRequestActuatorTypeAccedianNid captures enum value "accedian-nid"
	CreateTenantMonitoredObjectRequestActuatorTypeAccedianNid string = "accedian-nid"

	// CreateTenantMonitoredObjectRequestActuatorTypeAccedianVnid captures enum value "accedian-vnid"
	CreateTenantMonitoredObjectRequestActuatorTypeAccedianVnid string = "accedian-vnid"
)

// prop value enum
func (m *CreateTenantMonitoredObjectRequest) validateActuatorTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, createTenantMonitoredObjectRequestTypeActuatorTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CreateTenantMonitoredObjectRequest) validateActuatorType(formats strfmt.Registry) error {

	if swag.IsZero(m.ActuatorType) { // not required
		return nil
	}

	// value enum
	if err := m.validateActuatorTypeEnum("actuatorType", "body", m.ActuatorType); err != nil {
		return err
	}

	return nil
}

func (m *CreateTenantMonitoredObjectRequest) validateDomainSet(formats strfmt.Registry) error {

	if swag.IsZero(m.DomainSet) { // not required
		return nil
	}

	return nil
}

var createTenantMonitoredObjectRequestTypeObjectTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","flowmeter","twamp-pe","twamp-sf","twamp-sl"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		createTenantMonitoredObjectRequestTypeObjectTypePropEnum = append(createTenantMonitoredObjectRequestTypeObjectTypePropEnum, v)
	}
}

const (

	// CreateTenantMonitoredObjectRequestObjectTypeUnknown captures enum value "unknown"
	CreateTenantMonitoredObjectRequestObjectTypeUnknown string = "unknown"

	// CreateTenantMonitoredObjectRequestObjectTypeFlowmeter captures enum value "flowmeter"
	CreateTenantMonitoredObjectRequestObjectTypeFlowmeter string = "flowmeter"

	// CreateTenantMonitoredObjectRequestObjectTypeTwampPe captures enum value "twamp-pe"
	CreateTenantMonitoredObjectRequestObjectTypeTwampPe string = "twamp-pe"

	// CreateTenantMonitoredObjectRequestObjectTypeTwampSf captures enum value "twamp-sf"
	CreateTenantMonitoredObjectRequestObjectTypeTwampSf string = "twamp-sf"

	// CreateTenantMonitoredObjectRequestObjectTypeTwampSl captures enum value "twamp-sl"
	CreateTenantMonitoredObjectRequestObjectTypeTwampSl string = "twamp-sl"
)

// prop value enum
func (m *CreateTenantMonitoredObjectRequest) validateObjectTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, createTenantMonitoredObjectRequestTypeObjectTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CreateTenantMonitoredObjectRequest) validateObjectType(formats strfmt.Registry) error {

	if swag.IsZero(m.ObjectType) { // not required
		return nil
	}

	// value enum
	if err := m.validateObjectTypeEnum("objectType", "body", m.ObjectType); err != nil {
		return err
	}

	return nil
}

var createTenantMonitoredObjectRequestTypeReflectorTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["unknown","accedian-nid","accedian-vnid"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		createTenantMonitoredObjectRequestTypeReflectorTypePropEnum = append(createTenantMonitoredObjectRequestTypeReflectorTypePropEnum, v)
	}
}

const (

	// CreateTenantMonitoredObjectRequestReflectorTypeUnknown captures enum value "unknown"
	CreateTenantMonitoredObjectRequestReflectorTypeUnknown string = "unknown"

	// CreateTenantMonitoredObjectRequestReflectorTypeAccedianNid captures enum value "accedian-nid"
	CreateTenantMonitoredObjectRequestReflectorTypeAccedianNid string = "accedian-nid"

	// CreateTenantMonitoredObjectRequestReflectorTypeAccedianVnid captures enum value "accedian-vnid"
	CreateTenantMonitoredObjectRequestReflectorTypeAccedianVnid string = "accedian-vnid"
)

// prop value enum
func (m *CreateTenantMonitoredObjectRequest) validateReflectorTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, createTenantMonitoredObjectRequestTypeReflectorTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CreateTenantMonitoredObjectRequest) validateReflectorType(formats strfmt.Registry) error {

	if swag.IsZero(m.ReflectorType) { // not required
		return nil
	}

	// value enum
	if err := m.validateReflectorTypeEnum("reflectorType", "body", m.ReflectorType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CreateTenantMonitoredObjectRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateTenantMonitoredObjectRequest) UnmarshalBinary(b []byte) error {
	var res CreateTenantMonitoredObjectRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
