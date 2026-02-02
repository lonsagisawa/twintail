package requests

import (
	"twintail/internal/services"

	"github.com/labstack/echo/v5"
)

type StoreEndpointRequest struct {
	Protocol    string `form:"protocol" validate:"required,oneof=https http tcp+tls tcp"`
	ExposePort  string `form:"expose_port" validate:"required,numeric"`
	Destination string `form:"destination" validate:"required"`
}

func (r *StoreEndpointRequest) FromContext(ctx *echo.Context) error {
	if err := ctx.Bind(r); err != nil {
		return err
	}
	return ctx.Validate(r)
}

func (r *StoreEndpointRequest) ToParams(serviceName string) services.EndpointParams {
	return services.EndpointParams{
		ServiceName: serviceName,
		Protocol:    r.Protocol,
		ExposePort:  r.ExposePort,
		Destination: r.Destination,
	}
}

func (r *StoreEndpointRequest) Default() StoreEndpointRequest {
	return StoreEndpointRequest{
		Protocol:   "https",
		ExposePort: "443",
	}
}

type DestroyEndpointRequest struct {
	Protocol    string `form:"protocol" validate:"required"`
	ExposePort  string `form:"expose_port" validate:"required"`
	Destination string `form:"destination" validate:"required"`
}

func (r *DestroyEndpointRequest) FromContext(ctx *echo.Context) error {
	if err := ctx.Bind(r); err != nil {
		return err
	}
	return ctx.Validate(r)
}

func (r *DestroyEndpointRequest) ToParams(serviceName string) services.EndpointParams {
	return services.EndpointParams{
		ServiceName: serviceName,
		Protocol:    r.Protocol,
		ExposePort:  r.ExposePort,
		Destination: r.Destination,
	}
}

type UpdateEndpointRequest struct {
	Protocol       string `form:"protocol" validate:"required,oneof=https http tcp+tls tcp"`
	ExposePort     string `form:"expose_port" validate:"required,numeric"`
	OldDestination string `form:"old_destination" validate:"required"`
	NewDestination string `form:"new_destination" validate:"required"`
}

func (r *UpdateEndpointRequest) FromContext(ctx *echo.Context) error {
	if err := ctx.Bind(r); err != nil {
		return err
	}
	return ctx.Validate(r)
}

func (r *UpdateEndpointRequest) ToParams(serviceName string) services.UpdateEndpointParams {
	return services.UpdateEndpointParams{
		ServiceName:    serviceName,
		Protocol:       r.Protocol,
		ExposePort:     r.ExposePort,
		OldDestination: r.OldDestination,
		NewDestination: r.NewDestination,
	}
}
