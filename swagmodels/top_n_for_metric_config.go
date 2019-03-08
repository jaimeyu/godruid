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

// TopNForMetricConfig The necessary request parameters for the metric api call
// swagger:model TopNForMetricConfig
type TopNForMetricConfig struct {
	BucketFilter

	TopNForMetricConfigAllOf1

	// The type of aggregation (avg/min/max)
	// Enum: [min max avg]
	Aggregator string `json:"aggregator,omitempty"`

	// dimensions
	Dimensions DimensionFilter `json:"dimensions,omitempty"`

	// A value of true will have the aggregation request execute on all data regardless of whether it has been cleaned or not
	IgnoreCleaning bool `json:"ignoreCleaning,omitempty"`

	// Time boundary for the metrics under consideration using the ISO-8601 standard
	Interval string `json:"interval,omitempty"`

	// meta
	Meta MetaFilter `json:"meta,omitempty"`

	// metric
	Metric *MetricIdentifierFilter `json:"metric,omitempty"`

	// metrics view
	MetricsView []*MetricView `json:"metricsView,omitempty"`

	// An optional array of monitored objects that we want to retrieve specific topn against. This attribute cannot be used if the meta attribute is also present in the request.
	MonitoredObjects []string `json:"monitoredObjects,omitempty"`

	// Number of results to return
	NumResults int64 `json:"numResults,omitempty"`

	// Indicates whether the response should return the topn in ascending or descending order. The default value is descending
	// Enum: [asc desc]
	Sorted string `json:"sorted,omitempty"`

	// Query timeout in milliseconds
	Timeout int64 `json:"timeout,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *TopNForMetricConfig) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 BucketFilter
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.BucketFilter = aO0

	// AO1
	var aO1 TopNForMetricConfigAllOf1
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.TopNForMetricConfigAllOf1 = aO1

	// AO2
	var dataAO2 struct {
		Aggregator string `json:"aggregator,omitempty"`

		Dimensions DimensionFilter `json:"dimensions,omitempty"`

		IgnoreCleaning bool `json:"ignoreCleaning,omitempty"`

		Interval string `json:"interval,omitempty"`

		Meta MetaFilter `json:"meta,omitempty"`

		Metric *MetricIdentifierFilter `json:"metric,omitempty"`

		MetricsView []*MetricView `json:"metricsView,omitempty"`

		MonitoredObjects []string `json:"monitoredObjects,omitempty"`

		NumResults int64 `json:"numResults,omitempty"`

		Sorted string `json:"sorted,omitempty"`

		Timeout int64 `json:"timeout,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO2); err != nil {
		return err
	}

	m.Aggregator = dataAO2.Aggregator

	m.Dimensions = dataAO2.Dimensions

	m.IgnoreCleaning = dataAO2.IgnoreCleaning

	m.Interval = dataAO2.Interval

	m.Meta = dataAO2.Meta

	m.Metric = dataAO2.Metric

	m.MetricsView = dataAO2.MetricsView

	m.MonitoredObjects = dataAO2.MonitoredObjects

	m.NumResults = dataAO2.NumResults

	m.Sorted = dataAO2.Sorted

	m.Timeout = dataAO2.Timeout

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m TopNForMetricConfig) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 3)

	aO0, err := swag.WriteJSON(m.BucketFilter)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.TopNForMetricConfigAllOf1)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)

	var dataAO2 struct {
		Aggregator string `json:"aggregator,omitempty"`

		Dimensions DimensionFilter `json:"dimensions,omitempty"`

		IgnoreCleaning bool `json:"ignoreCleaning,omitempty"`

		Interval string `json:"interval,omitempty"`

		Meta MetaFilter `json:"meta,omitempty"`

		Metric *MetricIdentifierFilter `json:"metric,omitempty"`

		MetricsView []*MetricView `json:"metricsView,omitempty"`

		MonitoredObjects []string `json:"monitoredObjects,omitempty"`

		NumResults int64 `json:"numResults,omitempty"`

		Sorted string `json:"sorted,omitempty"`

		Timeout int64 `json:"timeout,omitempty"`
	}

	dataAO2.Aggregator = m.Aggregator

	dataAO2.Dimensions = m.Dimensions

	dataAO2.IgnoreCleaning = m.IgnoreCleaning

	dataAO2.Interval = m.Interval

	dataAO2.Meta = m.Meta

	dataAO2.Metric = m.Metric

	dataAO2.MetricsView = m.MetricsView

	dataAO2.MonitoredObjects = m.MonitoredObjects

	dataAO2.NumResults = m.NumResults

	dataAO2.Sorted = m.Sorted

	dataAO2.Timeout = m.Timeout

	jsonDataAO2, errAO2 := swag.WriteJSON(dataAO2)
	if errAO2 != nil {
		return nil, errAO2
	}
	_parts = append(_parts, jsonDataAO2)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this top n for metric config
func (m *TopNForMetricConfig) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with BucketFilter
	if err := m.BucketFilter.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with TopNForMetricConfigAllOf1

	if err := m.validateAggregator(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDimensions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMeta(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetric(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetricsView(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSorted(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var topNForMetricConfigTypeAggregatorPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["min","max","avg"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		topNForMetricConfigTypeAggregatorPropEnum = append(topNForMetricConfigTypeAggregatorPropEnum, v)
	}
}

// property enum
func (m *TopNForMetricConfig) validateAggregatorEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, topNForMetricConfigTypeAggregatorPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *TopNForMetricConfig) validateAggregator(formats strfmt.Registry) error {

	if swag.IsZero(m.Aggregator) { // not required
		return nil
	}

	// value enum
	if err := m.validateAggregatorEnum("aggregator", "body", m.Aggregator); err != nil {
		return err
	}

	return nil
}

func (m *TopNForMetricConfig) validateDimensions(formats strfmt.Registry) error {

	if swag.IsZero(m.Dimensions) { // not required
		return nil
	}

	if err := m.Dimensions.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("dimensions")
		}
		return err
	}

	return nil
}

func (m *TopNForMetricConfig) validateMeta(formats strfmt.Registry) error {

	if swag.IsZero(m.Meta) { // not required
		return nil
	}

	if err := m.Meta.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("meta")
		}
		return err
	}

	return nil
}

func (m *TopNForMetricConfig) validateMetric(formats strfmt.Registry) error {

	if swag.IsZero(m.Metric) { // not required
		return nil
	}

	if m.Metric != nil {
		if err := m.Metric.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("metric")
			}
			return err
		}
	}

	return nil
}

func (m *TopNForMetricConfig) validateMetricsView(formats strfmt.Registry) error {

	if swag.IsZero(m.MetricsView) { // not required
		return nil
	}

	for i := 0; i < len(m.MetricsView); i++ {
		if swag.IsZero(m.MetricsView[i]) { // not required
			continue
		}

		if m.MetricsView[i] != nil {
			if err := m.MetricsView[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("metricsView" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

var topNForMetricConfigTypeSortedPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["asc","desc"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		topNForMetricConfigTypeSortedPropEnum = append(topNForMetricConfigTypeSortedPropEnum, v)
	}
}

// property enum
func (m *TopNForMetricConfig) validateSortedEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, topNForMetricConfigTypeSortedPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *TopNForMetricConfig) validateSorted(formats strfmt.Registry) error {

	if swag.IsZero(m.Sorted) { // not required
		return nil
	}

	// value enum
	if err := m.validateSortedEnum("sorted", "body", m.Sorted); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TopNForMetricConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TopNForMetricConfig) UnmarshalBinary(b []byte) error {
	var res TopNForMetricConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// TopNForMetricConfigAllOf1 top n for metric config all of1
// swagger:model TopNForMetricConfigAllOf1
type TopNForMetricConfigAllOf1 interface{}
