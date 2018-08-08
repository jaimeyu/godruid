// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetMonitoredObjectV2HandlerFunc turns a function with the right signature into a get monitored object v2 handler
type GetMonitoredObjectV2HandlerFunc func(GetMonitoredObjectV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn GetMonitoredObjectV2HandlerFunc) Handle(params GetMonitoredObjectV2Params) middleware.Responder {
	return fn(params)
}

// GetMonitoredObjectV2Handler interface for that can handle valid get monitored object v2 params
type GetMonitoredObjectV2Handler interface {
	Handle(GetMonitoredObjectV2Params) middleware.Responder
}

// NewGetMonitoredObjectV2 creates a new http.Handler for the get monitored object v2 operation
func NewGetMonitoredObjectV2(ctx *middleware.Context, handler GetMonitoredObjectV2Handler) *GetMonitoredObjectV2 {
	return &GetMonitoredObjectV2{Context: ctx, Handler: handler}
}

/*GetMonitoredObjectV2 swagger:route GET /v2/monitored-objects/{monObjId} TenantProvisioningServiceV2 getMonitoredObjectV2

Retrieve a Tenant Monitored Object by id.

*/
type GetMonitoredObjectV2 struct {
	Context *middleware.Context
	Handler GetMonitoredObjectV2Handler
}

func (o *GetMonitoredObjectV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetMonitoredObjectV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
