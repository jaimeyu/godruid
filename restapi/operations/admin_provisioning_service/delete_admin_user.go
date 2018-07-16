// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteAdminUserHandlerFunc turns a function with the right signature into a delete admin user handler
type DeleteAdminUserHandlerFunc func(DeleteAdminUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteAdminUserHandlerFunc) Handle(params DeleteAdminUserParams) middleware.Responder {
	return fn(params)
}

// DeleteAdminUserHandler interface for that can handle valid delete admin user params
type DeleteAdminUserHandler interface {
	Handle(DeleteAdminUserParams) middleware.Responder
}

// NewDeleteAdminUser creates a new http.Handler for the delete admin user operation
func NewDeleteAdminUser(ctx *middleware.Context, handler DeleteAdminUserHandler) *DeleteAdminUser {
	return &DeleteAdminUser{Context: ctx, Handler: handler}
}

/*DeleteAdminUser swagger:route DELETE /v1/admin/{value} AdminProvisioningService deleteAdminUser

Delete a User with Administrative access.

*/
type DeleteAdminUser struct {
	Context *middleware.Context
	Handler DeleteAdminUserHandler
}

func (o *DeleteAdminUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteAdminUserParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
