// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetTenantIDByAliasV2HandlerFunc turns a function with the right signature into a get tenant Id by alias v2 handler
type GetTenantIDByAliasV2HandlerFunc func(GetTenantIDByAliasV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTenantIDByAliasV2HandlerFunc) Handle(params GetTenantIDByAliasV2Params) middleware.Responder {
	return fn(params)
}

// GetTenantIDByAliasV2Handler interface for that can handle valid get tenant Id by alias v2 params
type GetTenantIDByAliasV2Handler interface {
	Handle(GetTenantIDByAliasV2Params) middleware.Responder
}

// NewGetTenantIDByAliasV2 creates a new http.Handler for the get tenant Id by alias v2 operation
func NewGetTenantIDByAliasV2(ctx *middleware.Context, handler GetTenantIDByAliasV2Handler) *GetTenantIDByAliasV2 {
	return &GetTenantIDByAliasV2{Context: ctx, Handler: handler}
}

/*GetTenantIDByAliasV2 swagger:route GET /v2/tenant-by-alias/{value} AdminProvisioningServiceV2 getTenantIdByAliasV2

Returns the Id of a Tenant that matches the provided alias.

*/
type GetTenantIDByAliasV2 struct {
	Context *middleware.Context
	Handler GetTenantIDByAliasV2Handler
}

func (o *GetTenantIDByAliasV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTenantIDByAliasV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
