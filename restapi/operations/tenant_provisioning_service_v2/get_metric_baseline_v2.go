// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetMetricBaselineV2HandlerFunc turns a function with the right signature into a get metric baseline v2 handler
type GetMetricBaselineV2HandlerFunc func(GetMetricBaselineV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn GetMetricBaselineV2HandlerFunc) Handle(params GetMetricBaselineV2Params) middleware.Responder {
	return fn(params)
}

// GetMetricBaselineV2Handler interface for that can handle valid get metric baseline v2 params
type GetMetricBaselineV2Handler interface {
	Handle(GetMetricBaselineV2Params) middleware.Responder
}

// NewGetMetricBaselineV2 creates a new http.Handler for the get metric baseline v2 operation
func NewGetMetricBaselineV2(ctx *middleware.Context, handler GetMetricBaselineV2Handler) *GetMetricBaselineV2 {
	return &GetMetricBaselineV2{Context: ctx, Handler: handler}
}

/*GetMetricBaselineV2 swagger:route GET /v2/metric-baselines/{metricBaselineId} TenantProvisioningServiceV2 getMetricBaselineV2

Retrieve a Tenant Metric Baseline by id.

*/
type GetMetricBaselineV2 struct {
	Context *middleware.Context
	Handler GetMetricBaselineV2Handler
}

func (o *GetMetricBaselineV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetMetricBaselineV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
