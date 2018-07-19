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

// GetThresholdCrossingURL generates an URL for the get threshold crossing operation
type GetThresholdCrossingURL struct {
	Direction          []string
	Domain             []string
	Granularity        *string
	Interval           string
	Meta               []string
	Metric             []string
	ObjectType         []string
	Tenant             string
	ThresholdProfileID string
	Timeout            *int32
	Vendor             []string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetThresholdCrossingURL) WithBasePath(bp string) *GetThresholdCrossingURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetThresholdCrossingURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetThresholdCrossingURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/v1/threshold-crossing"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var directionIR []string
	for _, directionI := range o.Direction {
		directionIS := directionI
		if directionIS != "" {
			directionIR = append(directionIR, directionIS)
		}
	}

	direction := swag.JoinByFormat(directionIR, "")

	if len(direction) > 0 {
		qsv := direction[0]
		if qsv != "" {
			qs.Set("direction", qsv)
		}
	}

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

	var metaIR []string
	for _, metaI := range o.Meta {
		metaIS := metaI
		if metaIS != "" {
			metaIR = append(metaIR, metaIS)
		}
	}

	meta := swag.JoinByFormat(metaIR, "")

	if len(meta) > 0 {
		qsv := meta[0]
		if qsv != "" {
			qs.Set("meta", qsv)
		}
	}

	var metricIR []string
	for _, metricI := range o.Metric {
		metricIS := metricI
		if metricIS != "" {
			metricIR = append(metricIR, metricIS)
		}
	}

	metric := swag.JoinByFormat(metricIR, "")

	if len(metric) > 0 {
		qsv := metric[0]
		if qsv != "" {
			qs.Set("metric", qsv)
		}
	}

	var objectTypeIR []string
	for _, objectTypeI := range o.ObjectType {
		objectTypeIS := objectTypeI
		if objectTypeIS != "" {
			objectTypeIR = append(objectTypeIR, objectTypeIS)
		}
	}

	objectType := swag.JoinByFormat(objectTypeIR, "")

	if len(objectType) > 0 {
		qsv := objectType[0]
		if qsv != "" {
			qs.Set("objectType", qsv)
		}
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

	var vendorIR []string
	for _, vendorI := range o.Vendor {
		vendorIS := vendorI
		if vendorIS != "" {
			vendorIR = append(vendorIR, vendorIS)
		}
	}

	vendor := swag.JoinByFormat(vendorIR, "")

	if len(vendor) > 0 {
		qsv := vendor[0]
		if qsv != "" {
			qs.Set("vendor", qsv)
		}
	}

	result.RawQuery = qs.Encode()

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetThresholdCrossingURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetThresholdCrossingURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetThresholdCrossingURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetThresholdCrossingURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetThresholdCrossingURL")
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
func (o *GetThresholdCrossingURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
