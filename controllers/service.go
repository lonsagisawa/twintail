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
	ServiceName string
	Protocol    string
	ExposePort  string
	Destination string
}

func (c *ServiceController) Create(ctx *echo.Context) error {
	return ctx.Render(200, "new_service.html", map[string]any{
		"FormData": NewServiceFormData{Protocol: "https", ExposePort: "443"},
	})
}

func (c *ServiceController) Store(ctx *echo.Context) error {
	formData := NewServiceFormData{
		ServiceName: ctx.FormValue("service_name"),
		Protocol:    ctx.FormValue("protocol"),
		ExposePort:  ctx.FormValue("expose_port"),
		Destination: ctx.FormValue("destination"),
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

	return ctx.Render(200, "new_service.html", map[string]any{
		"Success":  "Service '" + formData.ServiceName + "' has been advertised successfully.",
		"FormData": NewServiceFormData{Protocol: "https", ExposePort: "443"},
	})
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

func (c *ServiceController) ConfirmDelete(ctx *echo.Context) error {
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
