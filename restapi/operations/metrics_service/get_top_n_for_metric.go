// Code generated by go-swagger; DO NOT EDIT.

package metrics_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetTopNForMetricHandlerFunc turns a function with the right signature into a get top n for metric handler
type GetTopNForMetricHandlerFunc func(GetTopNForMetricParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTopNForMetricHandlerFunc) Handle(params GetTopNForMetricParams) middleware.Responder {
	return fn(params)
}

// GetTopNForMetricHandler interface for that can handle valid get top n for metric params
type GetTopNForMetricHandler interface {
	Handle(GetTopNForMetricParams) middleware.Responder
}

// NewGetTopNForMetric creates a new http.Handler for the get top n for metric operation
func NewGetTopNForMetric(ctx *middleware.Context, handler GetTopNForMetricHandler) *GetTopNForMetric {
	return &GetTopNForMetric{Context: ctx, Handler: handler}
}

/*GetTopNForMetric swagger:route POST /v1/topn-metrics MetricsService getTopNForMetric

Retrieve a Top-N for a given sets of metrics

*/
type GetTopNForMetric struct {
	Context *middleware.Context
	Handler GetTopNForMetricHandler
}

func (o *GetTopNForMetric) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTopNForMetricParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
