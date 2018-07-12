// Code generated by go-swagger; DO NOT EDIT.

package admin_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteIngestionDictionaryHandlerFunc turns a function with the right signature into a delete ingestion dictionary handler
type DeleteIngestionDictionaryHandlerFunc func(DeleteIngestionDictionaryParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteIngestionDictionaryHandlerFunc) Handle(params DeleteIngestionDictionaryParams) middleware.Responder {
	return fn(params)
}

// DeleteIngestionDictionaryHandler interface for that can handle valid delete ingestion dictionary params
type DeleteIngestionDictionaryHandler interface {
	Handle(DeleteIngestionDictionaryParams) middleware.Responder
}

// NewDeleteIngestionDictionary creates a new http.Handler for the delete ingestion dictionary operation
func NewDeleteIngestionDictionary(ctx *middleware.Context, handler DeleteIngestionDictionaryHandler) *DeleteIngestionDictionary {
	return &DeleteIngestionDictionary{Context: ctx, Handler: handler}
}

/*DeleteIngestionDictionary swagger:route DELETE /v1/ingestion-dictionaries AdminProvisioningService deleteIngestionDictionary

Delete an Ingestion Dictionary.

*/
type DeleteIngestionDictionary struct {
	Context *middleware.Context
	Handler DeleteIngestionDictionaryHandler
}

func (o *DeleteIngestionDictionary) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteIngestionDictionaryParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
