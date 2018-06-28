// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetValidTypesHandlerFunc turns a function with the right signature into a get valid types handler
type GetValidTypesHandlerFunc func(GetValidTypesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetValidTypesHandlerFunc) Handle(params GetValidTypesParams) middleware.Responder {
	return fn(params)
}

// GetValidTypesHandler interface for that can handle valid get valid types params
type GetValidTypesHandler interface {
	Handle(GetValidTypesParams) middleware.Responder
}

// NewGetValidTypes creates a new http.Handler for the get valid types operation
func NewGetValidTypes(ctx *middleware.Context, handler GetValidTypesHandler) *GetValidTypes {
	return &GetValidTypes{Context: ctx, Handler: handler}
}

/*GetValidTypes swagger:route GET /v1/valid-types AdminProvisioningService getValidTypes

Retrieve a Valid Types object.

*/
type GetValidTypes struct {
	Context *middleware.Context
	Handler GetValidTypesHandler
}

func (o *GetValidTypes) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetValidTypesParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
