// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateDashboardV2HandlerFunc turns a function with the right signature into a update dashboard v2 handler
type UpdateDashboardV2HandlerFunc func(UpdateDashboardV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateDashboardV2HandlerFunc) Handle(params UpdateDashboardV2Params) middleware.Responder {
	return fn(params)
}

// UpdateDashboardV2Handler interface for that can handle valid update dashboard v2 params
type UpdateDashboardV2Handler interface {
	Handle(UpdateDashboardV2Params) middleware.Responder
}

// NewUpdateDashboardV2 creates a new http.Handler for the update dashboard v2 operation
func NewUpdateDashboardV2(ctx *middleware.Context, handler UpdateDashboardV2Handler) *UpdateDashboardV2 {
	return &UpdateDashboardV2{Context: ctx, Handler: handler}
}

/*UpdateDashboardV2 swagger:route PATCH /v2/dashboards/{dashboardId} TenantProvisioningServiceV2 updateDashboardV2

Update a Tenant Dashboard specified by the provided Dashboard Id.

*/
type UpdateDashboardV2 struct {
	Context *middleware.Context
	Handler UpdateDashboardV2Handler
}

func (o *UpdateDashboardV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateDashboardV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
