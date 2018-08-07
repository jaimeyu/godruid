// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetReportScheduleConfigHandlerFunc turns a function with the right signature into a get report schedule config handler
type GetReportScheduleConfigHandlerFunc func(GetReportScheduleConfigParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetReportScheduleConfigHandlerFunc) Handle(params GetReportScheduleConfigParams) middleware.Responder {
	return fn(params)
}

// GetReportScheduleConfigHandler interface for that can handle valid get report schedule config params
type GetReportScheduleConfigHandler interface {
	Handle(GetReportScheduleConfigParams) middleware.Responder
}

// NewGetReportScheduleConfig creates a new http.Handler for the get report schedule config operation
func NewGetReportScheduleConfig(ctx *middleware.Context, handler GetReportScheduleConfigHandler) *GetReportScheduleConfig {
	return &GetReportScheduleConfig{Context: ctx, Handler: handler}
}

/*GetReportScheduleConfig swagger:route GET /v1/tenants/{tenantId}/report-schedule-configs/{configId} TenantProvisioningService getReportScheduleConfig

Retrieve a report schedule configuration for a Tenant by configuration Id.

*/
type GetReportScheduleConfig struct {
	Context *middleware.Context
	Handler GetReportScheduleConfigHandler
}

func (o *GetReportScheduleConfig) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetReportScheduleConfigParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
