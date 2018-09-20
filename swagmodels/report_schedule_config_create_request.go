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

// ReportScheduleConfigCreateRequest Object used to create a new Report Genaration Schedule in Datahub
// swagger:model ReportScheduleConfigCreateRequest
type ReportScheduleConfigCreateRequest struct {

	// data
	// Required: true
	Data *ReportScheduleConfigCreateRequestData `json:"data"`
}

// Validate validates this report schedule config create request
func (m *ReportScheduleConfigCreateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ReportScheduleConfigCreateRequest) validateData(formats strfmt.Registry) error {

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
func (m *ReportScheduleConfigCreateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReportScheduleConfigCreateRequest) UnmarshalBinary(b []byte) error {
	var res ReportScheduleConfigCreateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ReportScheduleConfigCreateRequestData report schedule config create request data
// swagger:model ReportScheduleConfigCreateRequestData
type ReportScheduleConfigCreateRequestData struct {

	// attributes
	// Required: true
	Attributes *ReportScheduleConfigCreateRequestDataAttributes `json:"attributes"`

	// id
	ID string `json:"id,omitempty"`

	// relationships
	Relationships *ReportScheduleConfigRelationships `json:"relationships,omitempty"`

	// type
	// Required: true
	// Enum: [reportScheduleConfigs]
	Type *string `json:"type"`
}

// Validate validates this report schedule config create request data
func (m *ReportScheduleConfigCreateRequestData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
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

func (m *ReportScheduleConfigCreateRequestData) validateAttributes(formats strfmt.Registry) error {

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

func (m *ReportScheduleConfigCreateRequestData) validateRelationships(formats strfmt.Registry) error {

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

var reportScheduleConfigCreateRequestDataTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["reportScheduleConfigs"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		reportScheduleConfigCreateRequestDataTypeTypePropEnum = append(reportScheduleConfigCreateRequestDataTypeTypePropEnum, v)
	}
}

const (

	// ReportScheduleConfigCreateRequestDataTypeReportScheduleConfigs captures enum value "reportScheduleConfigs"
	ReportScheduleConfigCreateRequestDataTypeReportScheduleConfigs string = "reportScheduleConfigs"
)

// prop value enum
func (m *ReportScheduleConfigCreateRequestData) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, reportScheduleConfigCreateRequestDataTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ReportScheduleConfigCreateRequestData) validateType(formats strfmt.Registry) error {

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
func (m *ReportScheduleConfigCreateRequestData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReportScheduleConfigCreateRequestData) UnmarshalBinary(b []byte) error {
	var res ReportScheduleConfigCreateRequestData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ReportScheduleConfigCreateRequestDataAttributes report schedule config create request data attributes
// swagger:model ReportScheduleConfigCreateRequestDataAttributes
type ReportScheduleConfigCreateRequestDataAttributes struct {

	// When true, the report will be generated. When false, the report will not be generated
	Active bool `json:"active,omitempty"`

	// Recurring day of the month when this report should be generated
	DayMonth string `json:"dayMonth,omitempty"`

	// Recurring day of the week when this report should be generated
	DayWeek string `json:"dayWeek,omitempty"`

	// Time period for which individual results should be aggregated
	Granularity string `json:"granularity,omitempty"`

	// Recurring hour when this report should be generated
	Hour string `json:"hour,omitempty"`

	// meta
	Meta map[string][]string `json:"meta,omitempty"`

	// Recurring minute when this report should be generated
	Minute string `json:"minute,omitempty"`

	// Recurring month when this report should be generated
	Month string `json:"month,omitempty"`

	// Identifying name for the report to be generated
	// Required: true
	Name *string `json:"name"`

	// The type of report this config will generate
	ReportType string `json:"reportType,omitempty"`

	// The unique identifier of the Threshold Profile used to generate the report
	ThresholdProfile string `json:"thresholdProfile,omitempty"`

	// Period of time for which the report will be generated
	TimeRangeDuration string `json:"timeRangeDuration,omitempty"`

	// Amount if time, in ms, before which the request to generate the report should be cancelled
	Timeout int64 `json:"timeout,omitempty"`

	// Timezone used to display the results in the generated report
	Timezone string `json:"timezone,omitempty"`
}

// Validate validates this report schedule config create request data attributes
func (m *ReportScheduleConfigCreateRequestDataAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ReportScheduleConfigCreateRequestDataAttributes) validateName(formats strfmt.Registry) error {

	if err := validate.Required("data"+"."+"attributes"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ReportScheduleConfigCreateRequestDataAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReportScheduleConfigCreateRequestDataAttributes) UnmarshalBinary(b []byte) error {
	var res ReportScheduleConfigCreateRequestDataAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
