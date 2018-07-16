// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateIngestionDictionaryHandlerFunc turns a function with the right signature into a update ingestion dictionary handler
type UpdateIngestionDictionaryHandlerFunc func(UpdateIngestionDictionaryParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateIngestionDictionaryHandlerFunc) Handle(params UpdateIngestionDictionaryParams) middleware.Responder {
	return fn(params)
}

// UpdateIngestionDictionaryHandler interface for that can handle valid update ingestion dictionary params
type UpdateIngestionDictionaryHandler interface {
	Handle(UpdateIngestionDictionaryParams) middleware.Responder
}

// NewUpdateIngestionDictionary creates a new http.Handler for the update ingestion dictionary operation
func NewUpdateIngestionDictionary(ctx *middleware.Context, handler UpdateIngestionDictionaryHandler) *UpdateIngestionDictionary {
	return &UpdateIngestionDictionary{Context: ctx, Handler: handler}
}

/*UpdateIngestionDictionary swagger:route PUT /v1/ingestion-dictionaries AdminProvisioningService updateIngestionDictionary

Update an Ingestion Dictionary.

*/
type UpdateIngestionDictionary struct {
	Context *middleware.Context
	Handler UpdateIngestionDictionaryHandler
}

func (o *UpdateIngestionDictionary) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateIngestionDictionaryParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
