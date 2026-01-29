package controllers

import (
	"twintail/services"

	"github.com/labstack/echo/v5"
)

type TailscaleService interface {
	GetServeStatus() ([]services.ServiceView, error)
	GetServiceByName(name string) (*services.ServiceDetailView, error)
	AdvertiseService(params services.AdvertiseServiceParams) error
	ClearService(name string) error
}

type ServiceController struct {
	tailscale TailscaleService
}

func NewServiceController(tailscale TailscaleService) *ServiceController {
	return &ServiceController{
		tailscale: tailscale,
	}
}

func (c *ServiceController) Index(ctx *echo.Context) error {
	svcs, err := c.tailscale.GetServeStatus()
	if err != nil {
		return ctx.String(500, "Failed to get serve status: "+err.Error())
	}
	return ctx.Render(200, "index.html", map[string]any{
		"Services": svcs,
	})
}

type NewServiceFormData struct {
	ServiceName string `form:"service_name" validate:"required"`
	Protocol    string `form:"protocol" validate:"required,oneof=https http tcp+tls tcp"`
	ExposePort  string `form:"expose_port" validate:"required,numeric"`
	Destination string `form:"destination" validate:"required"`
}

func (c *ServiceController) Create(ctx *echo.Context) error {
	return ctx.Render(200, "new_service.html", map[string]any{
		"FormData": NewServiceFormData{Protocol: "https", ExposePort: "443"},
	})
}

func (c *ServiceController) Store(ctx *echo.Context) error {
	var formData NewServiceFormData
	if err := ctx.Bind(&formData); err != nil {
		return ctx.Render(200, "new_service.html", map[string]any{
			"Error":    err.Error(),
			"FormData": formData,
		})
	}
	if err := ctx.Validate(&formData); err != nil {
		return ctx.Render(200, "new_service.html", map[string]any{
			"Error":    err.Error(),
			"FormData": formData,
		})
	}

	params := services.AdvertiseServiceParams{
		ServiceName: formData.ServiceName,
		Protocol:    formData.Protocol,
		ExposePort:  formData.ExposePort,
		Destination: formData.Destination,
	}

	if err := c.tailscale.AdvertiseService(params); err != nil {
		return ctx.Render(200, "new_service.html", map[string]any{
			"Error":    err.Error(),
			"FormData": formData,
		})
	}

	return ctx.Redirect(303, "/services/"+formData.ServiceName)
}

func (c *ServiceController) Show(ctx *echo.Context) error {
	name := ctx.Param("name")
	svc, err := c.tailscale.GetServiceByName(name)
	if err != nil {
		return ctx.String(500, "Failed to get service: "+err.Error())
	}
	if svc == nil {
		return ctx.String(404, "Service not found")
	}
	return ctx.Render(200, "show_service.html", map[string]any{
		"Service": svc,
	})
}

func (c *ServiceController) Delete(ctx *echo.Context) error {
	name := ctx.Param("name")
	svc, err := c.tailscale.GetServiceByName(name)
	if err != nil {
		return ctx.String(500, "Failed to get service: "+err.Error())
	}
	if svc == nil {
		return ctx.String(404, "Service not found")
	}
	return ctx.Render(200, "confirm_delete.html", map[string]any{
		"Service": svc,
	})
}

func (c *ServiceController) Destroy(ctx *echo.Context) error {
	name := ctx.Param("name")
	if err := c.tailscale.ClearService(name); err != nil {
		return ctx.String(500, "Failed to delete service: "+err.Error())
	}
	return ctx.Redirect(303, "/")
}
