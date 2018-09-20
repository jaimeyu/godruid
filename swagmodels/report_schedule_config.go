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

// ReportScheduleConfig report schedule config
// swagger:model ReportScheduleConfig
type ReportScheduleConfig struct {

	// attributes
	// Required: true
	Attributes *ReportScheduleConfigAttributes `json:"attributes"`

	// id
	// Required: true
	ID *string `json:"id"`

	// relationships
	Relationships *ReportScheduleConfigRelationships `json:"relationships,omitempty"`

	// type
	// Required: true
	// Enum: [reportScheduleConfigs]
	Type *string `json:"type"`
}

// Validate validates this report schedule config
func (m *ReportScheduleConfig) Validate(formats strfmt.Registry) error {
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

func (m *ReportScheduleConfig) validateAttributes(formats strfmt.Registry) error {

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

func (m *ReportScheduleConfig) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfig) validateRelationships(formats strfmt.Registry) error {

	if swag.IsZero(m.Relationships) { // not required
		return nil
	}

	if m.Relationships != nil {
		if err := m.Relationships.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("relationships")
			}
			return err
		}
	}

	return nil
}

var reportScheduleConfigTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["reportScheduleConfigs"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		reportScheduleConfigTypeTypePropEnum = append(reportScheduleConfigTypeTypePropEnum, v)
	}
}

const (

	// ReportScheduleConfigTypeReportScheduleConfigs captures enum value "reportScheduleConfigs"
	ReportScheduleConfigTypeReportScheduleConfigs string = "reportScheduleConfigs"
)

// prop value enum
func (m *ReportScheduleConfig) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, reportScheduleConfigTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ReportScheduleConfig) validateType(formats strfmt.Registry) error {

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
func (m *ReportScheduleConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReportScheduleConfig) UnmarshalBinary(b []byte) error {
	var res ReportScheduleConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ReportScheduleConfigAttributes report schedule config attributes
// swagger:model ReportScheduleConfigAttributes
type ReportScheduleConfigAttributes struct {

	// Value used to ensure updates to this object are handled in order.
	// Required: true
	Rev *string `json:"_rev"`

	// When true, the report will be generated. When false, the report will not be generated
	// Required: true
	Active *bool `json:"active"`

	// Time since epoch at which this object was instantiated.
	// Required: true
	CreatedTimestamp *int64 `json:"createdTimestamp"`

	// Name used to identify this type of record in Datahub
	// Required: true
	Datatype *string `json:"datatype"`

	// Recurring day of the month when this report should be generated
	// Required: true
	DayMonth *string `json:"dayMonth"`

	// Recurring day of the week when this report should be generated
	// Required: true
	DayWeek *string `json:"dayWeek"`

	// Time period for which individual results should be aggregated
	// Required: true
	Granularity *string `json:"granularity"`

	// Recurring hour when this report should be generated
	// Required: true
	Hour *string `json:"hour"`

	// Time since epoch at which this object was last altered.
	// Required: true
	LastModifiedTimestamp *int64 `json:"lastModifiedTimestamp"`

	// meta
	// Required: true
	Meta map[string][]string `json:"meta"`

	// Recurring minute when this report should be generated
	// Required: true
	Minute *string `json:"minute"`

	// Recurring month when this report should be generated
	// Required: true
	Month *string `json:"month"`

	// Identifying name for the report to be generated
	// Required: true
	Name *string `json:"name"`

	// The type of report this config will generate
	// Required: true
	ReportType *string `json:"reportType"`

	// Unique identifier of the Tenant in Datahub
	// Required: true
	TenantID *string `json:"tenantId"`

	// Period of time for which the report will be generated
	// Required: true
	TimeRangeDuration *string `json:"timeRangeDuration"`

	// Amount if time, in ms, before which the request to generate the report should be cancelled
	// Required: true
	Timeout *int64 `json:"timeout"`

	// Timezone used to display the results in the generated report
	// Required: true
	Timezone *string `json:"timezone"`
}

// Validate validates this report schedule config attributes
func (m *ReportScheduleConfigAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRev(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateActive(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDatatype(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDayMonth(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDayWeek(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGranularity(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHour(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastModifiedTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMeta(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMinute(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMonth(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReportType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTenantID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeRangeDuration(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeout(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimezone(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ReportScheduleConfigAttributes) validateRev(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"_rev", "body", m.Rev); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateActive(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"active", "body", m.Active); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateCreatedTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"createdTimestamp", "body", m.CreatedTimestamp); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateDatatype(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"datatype", "body", m.Datatype); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateDayMonth(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"dayMonth", "body", m.DayMonth); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateDayWeek(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"dayWeek", "body", m.DayWeek); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateGranularity(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"granularity", "body", m.Granularity); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateHour(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"hour", "body", m.Hour); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateLastModifiedTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"lastModifiedTimestamp", "body", m.LastModifiedTimestamp); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateMeta(formats strfmt.Registry) error {

	return nil
}

func (m *ReportScheduleConfigAttributes) validateMinute(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"minute", "body", m.Minute); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateMonth(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"month", "body", m.Month); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateName(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateReportType(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"reportType", "body", m.ReportType); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateTenantID(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"tenantId", "body", m.TenantID); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateTimeRangeDuration(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"timeRangeDuration", "body", m.TimeRangeDuration); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateTimeout(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"timeout", "body", m.Timeout); err != nil {
		return err
	}

	return nil
}

func (m *ReportScheduleConfigAttributes) validateTimezone(formats strfmt.Registry) error {

	if err := validate.Required("attributes"+"."+"timezone", "body", m.Timezone); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ReportScheduleConfigAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReportScheduleConfigAttributes) UnmarshalBinary(b []byte) error {
	var res ReportScheduleConfigAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
