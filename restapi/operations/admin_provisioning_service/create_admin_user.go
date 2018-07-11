// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// CreateAdminUserHandlerFunc turns a function with the right signature into a create admin user handler
type CreateAdminUserHandlerFunc func(CreateAdminUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateAdminUserHandlerFunc) Handle(params CreateAdminUserParams) middleware.Responder {
	return fn(params)
}

// CreateAdminUserHandler interface for that can handle valid create admin user params
type CreateAdminUserHandler interface {
	Handle(CreateAdminUserParams) middleware.Responder
}

// NewCreateAdminUser creates a new http.Handler for the create admin user operation
func NewCreateAdminUser(ctx *middleware.Context, handler CreateAdminUserHandler) *CreateAdminUser {
	return &CreateAdminUser{Context: ctx, Handler: handler}
}

/*CreateAdminUser swagger:route POST /v1/admin AdminProvisioningService createAdminUser

Create a User with Administrative access.

*/
type CreateAdminUser struct {
	Context *middleware.Context
	Handler CreateAdminUserHandler
}

func (o *CreateAdminUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateAdminUserParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
