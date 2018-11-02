// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteMetricBaselineV2HandlerFunc turns a function with the right signature into a delete metric baseline v2 handler
type DeleteMetricBaselineV2HandlerFunc func(DeleteMetricBaselineV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteMetricBaselineV2HandlerFunc) Handle(params DeleteMetricBaselineV2Params) middleware.Responder {
	return fn(params)
}

// DeleteMetricBaselineV2Handler interface for that can handle valid delete metric baseline v2 params
type DeleteMetricBaselineV2Handler interface {
	Handle(DeleteMetricBaselineV2Params) middleware.Responder
}

// NewDeleteMetricBaselineV2 creates a new http.Handler for the delete metric baseline v2 operation
func NewDeleteMetricBaselineV2(ctx *middleware.Context, handler DeleteMetricBaselineV2Handler) *DeleteMetricBaselineV2 {
	return &DeleteMetricBaselineV2{Context: ctx, Handler: handler}
}

/*DeleteMetricBaselineV2 swagger:route DELETE /v2/metric-baselines/{metricBaselineId} TenantProvisioningServiceV2 deleteMetricBaselineV2

Delete a Tenant Metric Baseline specified by the provided Metric Baseline Id.

*/
type DeleteMetricBaselineV2 struct {
	Context *middleware.Context
	Handler DeleteMetricBaselineV2Handler
}

func (o *DeleteMetricBaselineV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteMetricBaselineV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
