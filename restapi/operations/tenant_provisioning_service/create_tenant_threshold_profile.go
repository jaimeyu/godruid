// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// CreateTenantThresholdProfileHandlerFunc turns a function with the right signature into a create tenant threshold profile handler
type CreateTenantThresholdProfileHandlerFunc func(CreateTenantThresholdProfileParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateTenantThresholdProfileHandlerFunc) Handle(params CreateTenantThresholdProfileParams) middleware.Responder {
	return fn(params)
}

// CreateTenantThresholdProfileHandler interface for that can handle valid create tenant threshold profile params
type CreateTenantThresholdProfileHandler interface {
	Handle(CreateTenantThresholdProfileParams) middleware.Responder
}

// NewCreateTenantThresholdProfile creates a new http.Handler for the create tenant threshold profile operation
func NewCreateTenantThresholdProfile(ctx *middleware.Context, handler CreateTenantThresholdProfileHandler) *CreateTenantThresholdProfile {
	return &CreateTenantThresholdProfile{Context: ctx, Handler: handler}
}

/*CreateTenantThresholdProfile swagger:route POST /v1/tenants/{tenantId}/threshold-profiles TenantProvisioningService createTenantThresholdProfile

Create a Threshold Profile for a Tenant.

*/
type CreateTenantThresholdProfile struct {
	Context *middleware.Context
	Handler CreateTenantThresholdProfileHandler
}

func (o *CreateTenantThresholdProfile) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateTenantThresholdProfileParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
