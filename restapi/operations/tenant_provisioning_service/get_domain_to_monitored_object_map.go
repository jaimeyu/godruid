// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetDomainToMonitoredObjectMapHandlerFunc turns a function with the right signature into a get domain to monitored object map handler
type GetDomainToMonitoredObjectMapHandlerFunc func(GetDomainToMonitoredObjectMapParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetDomainToMonitoredObjectMapHandlerFunc) Handle(params GetDomainToMonitoredObjectMapParams) middleware.Responder {
	return fn(params)
}

// GetDomainToMonitoredObjectMapHandler interface for that can handle valid get domain to monitored object map params
type GetDomainToMonitoredObjectMapHandler interface {
	Handle(GetDomainToMonitoredObjectMapParams) middleware.Responder
}

// NewGetDomainToMonitoredObjectMap creates a new http.Handler for the get domain to monitored object map operation
func NewGetDomainToMonitoredObjectMap(ctx *middleware.Context, handler GetDomainToMonitoredObjectMapHandler) *GetDomainToMonitoredObjectMap {
	return &GetDomainToMonitoredObjectMap{Context: ctx, Handler: handler}
}

/*GetDomainToMonitoredObjectMap swagger:route POST /v1/tenants/{tenantId}/monitored-object-domain-map TenantProvisioningService getDomainToMonitoredObjectMap

Retrieve a mapping of Monitored Objects that are associated with each Domain.

*/
type GetDomainToMonitoredObjectMap struct {
	Context *middleware.Context
	Handler GetDomainToMonitoredObjectMapHandler
}

func (o *GetDomainToMonitoredObjectMap) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetDomainToMonitoredObjectMapParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
