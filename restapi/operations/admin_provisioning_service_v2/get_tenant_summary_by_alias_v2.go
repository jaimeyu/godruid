// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetTenantSummaryByAliasV2HandlerFunc turns a function with the right signature into a get tenant summary by alias v2 handler
type GetTenantSummaryByAliasV2HandlerFunc func(GetTenantSummaryByAliasV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTenantSummaryByAliasV2HandlerFunc) Handle(params GetTenantSummaryByAliasV2Params) middleware.Responder {
	return fn(params)
}

// GetTenantSummaryByAliasV2Handler interface for that can handle valid get tenant summary by alias v2 params
type GetTenantSummaryByAliasV2Handler interface {
	Handle(GetTenantSummaryByAliasV2Params) middleware.Responder
}

// NewGetTenantSummaryByAliasV2 creates a new http.Handler for the get tenant summary by alias v2 operation
func NewGetTenantSummaryByAliasV2(ctx *middleware.Context, handler GetTenantSummaryByAliasV2Handler) *GetTenantSummaryByAliasV2 {
	return &GetTenantSummaryByAliasV2{Context: ctx, Handler: handler}
}

/*GetTenantSummaryByAliasV2 swagger:route GET /v2/tenant-summary-by-alias/{value} AdminProvisioningServiceV2 getTenantSummaryByAliasV2

Returns a summary of the Tenant that matches the provided alias.

*/
type GetTenantSummaryByAliasV2 struct {
	Context *middleware.Context
	Handler GetTenantSummaryByAliasV2Handler
}

func (o *GetTenantSummaryByAliasV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTenantSummaryByAliasV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
