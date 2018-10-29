// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// IngestionProfileMetricList Provides an array of objects which outline the vendor/monitoredObjectType/metrics that are actively being stroed in Datahub
// swagger:model IngestionProfileMetricList
type IngestionProfileMetricList []*IngestionProfileMetricListItems0

// Validate validates this ingestion profile metric list
func (m IngestionProfileMetricList) Validate(formats strfmt.Registry) error {
	var res []error

	for i := 0; i < len(m); i++ {
		if swag.IsZero(m[i]) { // not required
			continue
		}

		if m[i] != nil {
			if err := m[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName(strconv.Itoa(i))
				}
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// IngestionProfileMetricListItems0 ingestion profile metric list items0
// swagger:model IngestionProfileMetricListItems0
type IngestionProfileMetricListItems0 struct {

	// Provides data properties by which the Metric may be filtered and/or aggregated
	Dimensions interface{} `json:"dimensions,omitempty"`

	// Describes the direction of the test in case a Threshold needs to be different for one direction (i.e. actuator to reflector) versus another (i.e. round trip)
	Direction string `json:"direction,omitempty"`

	// When true, this metric will be recorded by Datahub. When false, this metric is ommitted.
	Enabled bool `json:"enabled,omitempty"`

	// The name of the Metric
	Metric string `json:"metric,omitempty"`

	// The name of the type of Monitored Object for which this Metric is applicable
	MonitoredObjectType string `json:"monitoredObjectType,omitempty"`

	// The name of the Vendor from which this Metric originates
	Vendor string `json:"vendor,omitempty"`
}

// Validate validates this ingestion profile metric list items0
func (m *IngestionProfileMetricListItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *IngestionProfileMetricListItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IngestionProfileMetricListItems0) UnmarshalBinary(b []byte) error {
	var res IngestionProfileMetricListItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
