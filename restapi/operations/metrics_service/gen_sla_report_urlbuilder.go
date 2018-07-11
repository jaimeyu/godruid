// Code generated by go-swagger; DO NOT EDIT.

package metrics_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"

	"github.com/go-openapi/swag"
)

// GenSLAReportURL generates an URL for the gen SLA report operation
type GenSLAReportURL struct {
	Domain             []string
	Granularity        *string
	Interval           string
	Tenant             string
	ThresholdProfileID string
	Timeout            *int32
	Timezone           *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GenSLAReportURL) WithBasePath(bp string) *GenSLAReportURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GenSLAReportURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GenSLAReportURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/v1/generate-sla-report"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var domainIR []string
	for _, domainI := range o.Domain {
		domainIS := domainI
		if domainIS != "" {
			domainIR = append(domainIR, domainIS)
		}
	}

	domain := swag.JoinByFormat(domainIR, "")

	if len(domain) > 0 {
		qsv := domain[0]
		if qsv != "" {
			qs.Set("domain", qsv)
		}
	}

	var granularity string
	if o.Granularity != nil {
		granularity = *o.Granularity
	}
	if granularity != "" {
		qs.Set("granularity", granularity)
	}

	interval := o.Interval
	if interval != "" {
		qs.Set("interval", interval)
	}

	tenant := o.Tenant
	if tenant != "" {
		qs.Set("tenant", tenant)
	}

	thresholdProfileID := o.ThresholdProfileID
	if thresholdProfileID != "" {
		qs.Set("thresholdProfileId", thresholdProfileID)
	}

	var timeout string
	if o.Timeout != nil {
		timeout = swag.FormatInt32(*o.Timeout)
	}
	if timeout != "" {
		qs.Set("timeout", timeout)
	}

	var timezone string
	if o.Timezone != nil {
		timezone = *o.Timezone
	}
	if timezone != "" {
		qs.Set("timezone", timezone)
	}

	result.RawQuery = qs.Encode()

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GenSLAReportURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GenSLAReportURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GenSLAReportURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GenSLAReportURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GenSLAReportURL")
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
func (o *GenSLAReportURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
