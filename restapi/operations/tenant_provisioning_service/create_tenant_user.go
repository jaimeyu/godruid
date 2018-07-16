// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// CreateTenantUserHandlerFunc turns a function with the right signature into a create tenant user handler
type CreateTenantUserHandlerFunc func(CreateTenantUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateTenantUserHandlerFunc) Handle(params CreateTenantUserParams) middleware.Responder {
	return fn(params)
}

// CreateTenantUserHandler interface for that can handle valid create tenant user params
type CreateTenantUserHandler interface {
	Handle(CreateTenantUserParams) middleware.Responder
}

// NewCreateTenantUser creates a new http.Handler for the create tenant user operation
func NewCreateTenantUser(ctx *middleware.Context, handler CreateTenantUserHandler) *CreateTenantUser {
	return &CreateTenantUser{Context: ctx, Handler: handler}
}

/*CreateTenantUser swagger:route POST /v1/tenants/{tenantId}/users TenantProvisioningService createTenantUser

Create a User for a Tenant.

*/
type CreateTenantUser struct {
	Context *middleware.Context
	Handler CreateTenantUserHandler
}

func (o *CreateTenantUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateTenantUserParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
