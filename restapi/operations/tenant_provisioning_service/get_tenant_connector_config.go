// Code generated by go-swagger; DO NOT EDIT.

package tenant_provisioning_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetTenantConnectorConfigHandlerFunc turns a function with the right signature into a get tenant connector config handler
type GetTenantConnectorConfigHandlerFunc func(GetTenantConnectorConfigParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTenantConnectorConfigHandlerFunc) Handle(params GetTenantConnectorConfigParams) middleware.Responder {
	return fn(params)
}

// GetTenantConnectorConfigHandler interface for that can handle valid get tenant connector config params
type GetTenantConnectorConfigHandler interface {
	Handle(GetTenantConnectorConfigParams) middleware.Responder
}

// NewGetTenantConnectorConfig creates a new http.Handler for the get tenant connector config operation
func NewGetTenantConnectorConfig(ctx *middleware.Context, handler GetTenantConnectorConfigHandler) *GetTenantConnectorConfig {
	return &GetTenantConnectorConfig{Context: ctx, Handler: handler}
}

/*GetTenantConnectorConfig swagger:route GET /v1/tenants/{tenantId}/connector-configs/{connectorId} TenantProvisioningService getTenantConnectorConfig

Retrieve a Tenant ConnectorConfig by Id.

*/
type GetTenantConnectorConfig struct {
	Context *middleware.Context
	Handler GetTenantConnectorConfigHandler
}

func (o *GetTenantConnectorConfig) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTenantConnectorConfigParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
