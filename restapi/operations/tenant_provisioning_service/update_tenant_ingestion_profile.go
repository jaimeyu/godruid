// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateTenantIngestionProfileHandlerFunc turns a function with the right signature into a update tenant ingestion profile handler
type UpdateTenantIngestionProfileHandlerFunc func(UpdateTenantIngestionProfileParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateTenantIngestionProfileHandlerFunc) Handle(params UpdateTenantIngestionProfileParams) middleware.Responder {
	return fn(params)
}

// UpdateTenantIngestionProfileHandler interface for that can handle valid update tenant ingestion profile params
type UpdateTenantIngestionProfileHandler interface {
	Handle(UpdateTenantIngestionProfileParams) middleware.Responder
}

// NewUpdateTenantIngestionProfile creates a new http.Handler for the update tenant ingestion profile operation
func NewUpdateTenantIngestionProfile(ctx *middleware.Context, handler UpdateTenantIngestionProfileHandler) *UpdateTenantIngestionProfile {
	return &UpdateTenantIngestionProfile{Context: ctx, Handler: handler}
}

/*UpdateTenantIngestionProfile swagger:route PUT /v1/tenants/{tenantId}/ingestion-profiles TenantProvisioningService updateTenantIngestionProfile

Update a Tenant Ingestion Profile

*/
type UpdateTenantIngestionProfile struct {
	Context *middleware.Context
	Handler UpdateTenantIngestionProfileHandler
}

func (o *UpdateTenantIngestionProfile) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateTenantIngestionProfileParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
