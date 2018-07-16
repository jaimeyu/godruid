// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteReportScheduleConfigHandlerFunc turns a function with the right signature into a delete report schedule config handler
type DeleteReportScheduleConfigHandlerFunc func(DeleteReportScheduleConfigParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteReportScheduleConfigHandlerFunc) Handle(params DeleteReportScheduleConfigParams) middleware.Responder {
	return fn(params)
}

// DeleteReportScheduleConfigHandler interface for that can handle valid delete report schedule config params
type DeleteReportScheduleConfigHandler interface {
	Handle(DeleteReportScheduleConfigParams) middleware.Responder
}

// NewDeleteReportScheduleConfig creates a new http.Handler for the delete report schedule config operation
func NewDeleteReportScheduleConfig(ctx *middleware.Context, handler DeleteReportScheduleConfigHandler) *DeleteReportScheduleConfig {
	return &DeleteReportScheduleConfig{Context: ctx, Handler: handler}
}

/*DeleteReportScheduleConfig swagger:route DELETE /v1/tenants/{tenantId}/report-schedule-configs/{configId} TenantProvisioningService deleteReportScheduleConfig

Delete a report schedule configuration for a Tenant by configuration Id.

*/
type DeleteReportScheduleConfig struct {
	Context *middleware.Context
	Handler DeleteReportScheduleConfigHandler
}

func (o *DeleteReportScheduleConfig) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteReportScheduleConfigParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
