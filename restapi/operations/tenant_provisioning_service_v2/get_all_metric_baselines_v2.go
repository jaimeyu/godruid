// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetAllMetricBaselinesV2HandlerFunc turns a function with the right signature into a get all metric baselines v2 handler
type GetAllMetricBaselinesV2HandlerFunc func(GetAllMetricBaselinesV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAllMetricBaselinesV2HandlerFunc) Handle(params GetAllMetricBaselinesV2Params) middleware.Responder {
	return fn(params)
}

// GetAllMetricBaselinesV2Handler interface for that can handle valid get all metric baselines v2 params
type GetAllMetricBaselinesV2Handler interface {
	Handle(GetAllMetricBaselinesV2Params) middleware.Responder
}

// NewGetAllMetricBaselinesV2 creates a new http.Handler for the get all metric baselines v2 operation
func NewGetAllMetricBaselinesV2(ctx *middleware.Context, handler GetAllMetricBaselinesV2Handler) *GetAllMetricBaselinesV2 {
	return &GetAllMetricBaselinesV2{Context: ctx, Handler: handler}
}

/*GetAllMetricBaselinesV2 swagger:route GET /v2/metric-baselines TenantProvisioningServiceV2 getAllMetricBaselinesV2

Get all Tenant Metric Baselines

*/
type GetAllMetricBaselinesV2 struct {
	Context *middleware.Context
	Handler GetAllMetricBaselinesV2Handler
}

func (o *GetAllMetricBaselinesV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetAllMetricBaselinesV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
