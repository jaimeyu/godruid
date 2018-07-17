// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service_v2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateDataCleaningProfileHandlerFunc turns a function with the right signature into a update data cleaning profile handler
type UpdateDataCleaningProfileHandlerFunc func(UpdateDataCleaningProfileParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateDataCleaningProfileHandlerFunc) Handle(params UpdateDataCleaningProfileParams) middleware.Responder {
	return fn(params)
}

// UpdateDataCleaningProfileHandler interface for that can handle valid update data cleaning profile params
type UpdateDataCleaningProfileHandler interface {
	Handle(UpdateDataCleaningProfileParams) middleware.Responder
}

// NewUpdateDataCleaningProfile creates a new http.Handler for the update data cleaning profile operation
func NewUpdateDataCleaningProfile(ctx *middleware.Context, handler UpdateDataCleaningProfileHandler) *UpdateDataCleaningProfile {
	return &UpdateDataCleaningProfile{Context: ctx, Handler: handler}
}

/*UpdateDataCleaningProfile swagger:route PATCH /v2/data-cleaning-profiles TenantProvisioningServiceV2 updateDataCleaningProfile

Provides ability to Update a Tenant Data Cleaning Profile

*/
type UpdateDataCleaningProfile struct {
	Context *middleware.Context
	Handler UpdateDataCleaningProfileHandler
}

func (o *UpdateDataCleaningProfile) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateDataCleaningProfileParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
