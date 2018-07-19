// Code generated by go-swagger; DO NOT EDIT.

package swagmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// TenantThresholdProfileMetricMetricMap tenant threshold profile metric metric map
// swagger:model tenantThresholdProfileMetricMetricMap
type TenantThresholdProfileMetricMetricMap map[string]TenantThresholdProfileUIEventAttrMap

// Validate validates this tenant threshold profile metric metric map
func (m TenantThresholdProfileMetricMetricMap) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.Required("", "body", TenantThresholdProfileMetricMetricMap(m)); err != nil {
		return err
	}

	for k := range m {

		if val, ok := m[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
