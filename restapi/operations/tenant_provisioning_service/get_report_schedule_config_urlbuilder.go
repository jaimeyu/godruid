// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"
)

// GetReportScheduleConfigURL generates an URL for the get report schedule config operation
type GetReportScheduleConfigURL struct {
	ConfigID string
	TenantID string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetReportScheduleConfigURL) WithBasePath(bp string) *GetReportScheduleConfigURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetReportScheduleConfigURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetReportScheduleConfigURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/v1/tenants/{tenantId}/report-schedule-configs/{configId}"

	configID := o.ConfigID
	if configID != "" {
		_path = strings.Replace(_path, "{configId}", configID, -1)
	} else {
		return nil, errors.New("ConfigID is required on GetReportScheduleConfigURL")
	}

	tenantID := o.TenantID
	if tenantID != "" {
		_path = strings.Replace(_path, "{tenantId}", tenantID, -1)
	} else {
		return nil, errors.New("TenantID is required on GetReportScheduleConfigURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetReportScheduleConfigURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetReportScheduleConfigURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetReportScheduleConfigURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetReportScheduleConfigURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetReportScheduleConfigURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetReportScheduleConfigURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
