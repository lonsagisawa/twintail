package controllers

import (
	"twintail/services"

	"github.com/labstack/echo/v5"
)

type EndpointService interface {
	GetServiceByName(name string) (*services.ServiceDetailView, error)
	AddEndpoint(params services.EndpointParams) error
	RemoveEndpoint(params services.EndpointParams) error
}

type EndpointController struct {
	tailscale EndpointService
}

func NewEndpointController(tailscale EndpointService) *EndpointController {
	return &EndpointController{
		tailscale: tailscale,
	}
}

type EndpointFormData struct {
	Protocol    string
	ExposePort  string
	Destination string
}

func (c *EndpointController) Create(ctx *echo.Context) error {
	name := ctx.Param("name")
	return ctx.Render(200, "new_endpoint.html", map[string]any{
		"ServiceName": name,
		"FormData":    EndpointFormData{Protocol: "https", ExposePort: "443"},
	})
}

func (c *EndpointController) Store(ctx *echo.Context) error {
	name := ctx.Param("name")
	formData := EndpointFormData{
		Protocol:    ctx.FormValue("protocol"),
		ExposePort:  ctx.FormValue("expose_port"),
		Destination: ctx.FormValue("destination"),
	}

	params := services.EndpointParams{
		ServiceName: name,
		Protocol:    formData.Protocol,
		ExposePort:  formData.ExposePort,
		Destination: formData.Destination,
	}

	if err := c.tailscale.AddEndpoint(params); err != nil {
		return ctx.Render(200, "new_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    formData,
		})
	}

	return ctx.Redirect(303, "/services/"+name)
}

func (c *EndpointController) ConfirmDelete(ctx *echo.Context) error {
	name := ctx.Param("name")
	protocol := ctx.QueryParam("protocol")
	exposePort := ctx.QueryParam("port")
	destination := ctx.QueryParam("destination")

	return ctx.Render(200, "confirm_delete_endpoint.html", map[string]any{
		"ServiceName": name,
		"Protocol":    protocol,
		"ExposePort":  exposePort,
		"Destination": destination,
	})
}

func (c *EndpointController) Destroy(ctx *echo.Context) error {
	name := ctx.Param("name")
	params := services.EndpointParams{
		ServiceName: name,
		Protocol:    ctx.FormValue("protocol"),
		ExposePort:  ctx.FormValue("expose_port"),
		Destination: ctx.FormValue("destination"),
	}

	if err := c.tailscale.RemoveEndpoint(params); err != nil {
		return ctx.String(500, "Failed to delete endpoint: "+err.Error())
	}

	svc, _ := c.tailscale.GetServiceByName(name)
	if svc == nil {
		return ctx.Redirect(303, "/")
	}

	return ctx.Redirect(303, "/services/"+name)
}
