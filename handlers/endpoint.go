package handlers

import (
	"twintail/requests"
	"twintail/services"

	"github.com/labstack/echo/v5"
)

type EndpointService interface {
	CheckInstalled() error
	GetServiceByName(name string) (*services.ServiceDetailView, error)
	AddEndpoint(params services.EndpointParams) error
	RemoveEndpoint(params services.EndpointParams) error
	UpdateEndpoint(params services.UpdateEndpointParams) error
}

type EndpointHandler struct {
	tailscale EndpointService
}

func NewEndpointHandler(tailscale EndpointService) *EndpointHandler {
	return &EndpointHandler{
		tailscale: tailscale,
	}
}

func (h *EndpointHandler) Create(ctx *echo.Context) error {
	if err := h.tailscale.CheckInstalled(); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
	}
	name := ctx.Param("name")
	var req requests.StoreEndpointRequest
	return ctx.Render(200, "new_endpoint.html", map[string]any{
		"ServiceName": name,
		"FormData":    req.Default(),
	})
}

func (h *EndpointHandler) Store(ctx *echo.Context) error {
	name := ctx.Param("name")
	var req requests.StoreEndpointRequest
	if err := req.FromContext(ctx); err != nil {
		return ctx.Render(200, "new_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    req,
		})
	}

	if err := h.tailscale.AddEndpoint(req.ToParams(name)); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.Render(200, "new_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    req,
		})
	}

	return ctx.Redirect(303, "/services/"+name)
}

func (h *EndpointHandler) Delete(ctx *echo.Context) error {
	if err := h.tailscale.CheckInstalled(); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
	}
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

func (h *EndpointHandler) Destroy(ctx *echo.Context) error {
	name := ctx.Param("name")
	var req requests.DestroyEndpointRequest
	if err := req.FromContext(ctx); err != nil {
		return ctx.String(500, "Invalid request: "+err.Error())
	}

	if err := h.tailscale.RemoveEndpoint(req.ToParams(name)); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.String(500, "Failed to delete endpoint: "+err.Error())
	}

	svc, _ := h.tailscale.GetServiceByName(name)
	if svc == nil {
		return ctx.Redirect(303, "/")
	}

	return ctx.Redirect(303, "/services/"+name)
}

func (h *EndpointHandler) Edit(ctx *echo.Context) error {
	if err := h.tailscale.CheckInstalled(); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
	}
	name := ctx.Param("name")
	protocol := ctx.QueryParam("protocol")
	exposePort := ctx.QueryParam("port")
	destination := ctx.QueryParam("destination")

	return ctx.Render(200, "edit_endpoint.html", map[string]any{
		"ServiceName": name,
		"FormData": requests.UpdateEndpointRequest{
			Protocol:       protocol,
			ExposePort:     exposePort,
			OldDestination: destination,
			NewDestination: destination,
		},
	})
}

func (h *EndpointHandler) Update(ctx *echo.Context) error {
	name := ctx.Param("name")
	var req requests.UpdateEndpointRequest
	if err := req.FromContext(ctx); err != nil {
		return ctx.Render(200, "edit_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    req,
		})
	}

	if err := h.tailscale.UpdateEndpoint(req.ToParams(name)); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.Render(200, "edit_endpoint.html", map[string]any{
			"ServiceName": name,
			"Error":       err.Error(),
			"FormData":    req,
		})
	}

	return ctx.Redirect(303, "/services/"+name)
}
