// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetAllTenantThresholdProfilesHandlerFunc turns a function with the right signature into a get all tenant threshold profiles handler
type GetAllTenantThresholdProfilesHandlerFunc func(GetAllTenantThresholdProfilesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAllTenantThresholdProfilesHandlerFunc) Handle(params GetAllTenantThresholdProfilesParams) middleware.Responder {
	return fn(params)
}

// GetAllTenantThresholdProfilesHandler interface for that can handle valid get all tenant threshold profiles params
type GetAllTenantThresholdProfilesHandler interface {
	Handle(GetAllTenantThresholdProfilesParams) middleware.Responder
}

// NewGetAllTenantThresholdProfiles creates a new http.Handler for the get all tenant threshold profiles operation
func NewGetAllTenantThresholdProfiles(ctx *middleware.Context, handler GetAllTenantThresholdProfilesHandler) *GetAllTenantThresholdProfiles {
	return &GetAllTenantThresholdProfiles{Context: ctx, Handler: handler}
}

/*GetAllTenantThresholdProfiles swagger:route GET /v1/tenants/{tenantId}/threshold-profile-list TenantProvisioningService getAllTenantThresholdProfiles

Retrieve all Threshold Profiles for the specified Tenant.

*/
type GetAllTenantThresholdProfiles struct {
	Context *middleware.Context
	Handler GetAllTenantThresholdProfilesHandler
}

func (o *GetAllTenantThresholdProfiles) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetAllTenantThresholdProfilesParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
