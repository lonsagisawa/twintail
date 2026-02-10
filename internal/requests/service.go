package requests

import (
	"twintail/internal/services"

	"github.com/labstack/echo/v5"
)

type StoreServiceRequest struct {
	ServiceName string `form:"service_name" validate:"required,excludesall=; \n\r\x60\x00"`
	Protocol    string `form:"protocol" validate:"required,oneof=https http tcp+tls tcp"`
	ExposePort  string `form:"expose_port" validate:"required,numeric"`
	Destination string `form:"destination" validate:"required,excludesall=; \n\r\x60\x00"`
}

func (r *StoreServiceRequest) FromContext(ctx *echo.Context) error {
	if err := ctx.Bind(r); err != nil {
		return err
	}
	return ctx.Validate(r)
}

func (r *StoreServiceRequest) ToParams() services.AdvertiseServiceParams {
	return services.AdvertiseServiceParams{
		ServiceName: r.ServiceName,
		Protocol:    r.Protocol,
		ExposePort:  r.ExposePort,
		Destination: r.Destination,
	}
}

func (r *StoreServiceRequest) Default() StoreServiceRequest {
	return StoreServiceRequest{
		Protocol:   "https",
		ExposePort: "443",
	}
}
