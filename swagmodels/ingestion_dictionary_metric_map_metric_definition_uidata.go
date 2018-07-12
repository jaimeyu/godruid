// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// IngestionDictionaryMetricMapMetricDefinitionUidata ingestion dictionary metric map metric definition uidata
// swagger:model IngestionDictionaryMetricMapMetricDefinitionUIData
type IngestionDictionaryMetricMapMetricDefinitionUidata struct {

	// group
	Group string `json:"group,omitempty"`

	// position
	Position string `json:"position,omitempty"`
}

// Validate validates this ingestion dictionary metric map metric definition uidata
func (m *IngestionDictionaryMetricMapMetricDefinitionUidata) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *IngestionDictionaryMetricMapMetricDefinitionUidata) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IngestionDictionaryMetricMapMetricDefinitionUidata) UnmarshalBinary(b []byte) error {
	var res IngestionDictionaryMetricMapMetricDefinitionUidata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
