// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateReportScheduleConfigV2HandlerFunc turns a function with the right signature into a update report schedule config v2 handler
type UpdateReportScheduleConfigV2HandlerFunc func(UpdateReportScheduleConfigV2Params) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateReportScheduleConfigV2HandlerFunc) Handle(params UpdateReportScheduleConfigV2Params) middleware.Responder {
	return fn(params)
}

// UpdateReportScheduleConfigV2Handler interface for that can handle valid update report schedule config v2 params
type UpdateReportScheduleConfigV2Handler interface {
	Handle(UpdateReportScheduleConfigV2Params) middleware.Responder
}

// NewUpdateReportScheduleConfigV2 creates a new http.Handler for the update report schedule config v2 operation
func NewUpdateReportScheduleConfigV2(ctx *middleware.Context, handler UpdateReportScheduleConfigV2Handler) *UpdateReportScheduleConfigV2 {
	return &UpdateReportScheduleConfigV2{Context: ctx, Handler: handler}
}

/*UpdateReportScheduleConfigV2 swagger:route PATCH /v2/report-schedule-configs/{configId} TenantProvisioningServiceV2 updateReportScheduleConfigV2

Update a Report Schedule Configuration for a Tenant.

*/
type UpdateReportScheduleConfigV2 struct {
	Context *middleware.Context
	Handler UpdateReportScheduleConfigV2Handler
}

func (o *UpdateReportScheduleConfigV2) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateReportScheduleConfigV2Params()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
