package handlers

import (
	"twintail/requests"
	"twintail/services"

	"github.com/labstack/echo/v5"
)

type TailscaleService interface {
	CheckInstalled() error
	GetServeStatus() ([]services.ServiceView, error)
	GetServiceByName(name string) (*services.ServiceDetailView, error)
	AdvertiseService(params services.AdvertiseServiceParams) error
	ClearService(name string) error
}

type ServiceHandler struct {
	tailscale TailscaleService
}

func NewServiceHandler(tailscale TailscaleService) *ServiceHandler {
	return &ServiceHandler{
		tailscale: tailscale,
	}
}

func (h *ServiceHandler) Index(ctx *echo.Context) error {
	svcs, err := h.tailscale.GetServeStatus()
	if err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.String(500, "Failed to get serve status: "+err.Error())
	}
	return ctx.Render(200, "index.html", map[string]any{
		"Services": svcs,
	})
}

func (h *ServiceHandler) Create(ctx *echo.Context) error {
	if err := h.tailscale.CheckInstalled(); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
	}
	var req requests.StoreServiceRequest
	return ctx.Render(200, "new_service.html", map[string]any{
		"FormData": req.Default(),
	})
}

func (h *ServiceHandler) Store(ctx *echo.Context) error {
	var req requests.StoreServiceRequest
	if err := req.FromContext(ctx); err != nil {
		return ctx.Render(200, "new_service.html", map[string]any{
			"Error":    err.Error(),
			"FormData": req,
		})
	}

	if err := h.tailscale.AdvertiseService(req.ToParams()); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.Render(200, "new_service.html", map[string]any{
			"Error":    err.Error(),
			"FormData": req,
		})
	}

	return ctx.Redirect(303, "/services/"+req.ServiceName)
}

func (h *ServiceHandler) Show(ctx *echo.Context) error {
	name := ctx.Param("name")
	svc, err := h.tailscale.GetServiceByName(name)
	if err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.String(500, "Failed to get service: "+err.Error())
	}
	if svc == nil {
		return ctx.String(404, "Service not found")
	}
	return ctx.Render(200, "show_service.html", map[string]any{
		"Service": svc,
	})
}

func (h *ServiceHandler) Delete(ctx *echo.Context) error {
	name := ctx.Param("name")
	svc, err := h.tailscale.GetServiceByName(name)
	if err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.String(500, "Failed to get service: "+err.Error())
	}
	if svc == nil {
		return ctx.String(404, "Service not found")
	}
	return ctx.Render(200, "confirm_delete.html", map[string]any{
		"Service": svc,
	})
}

func (h *ServiceHandler) Destroy(ctx *echo.Context) error {
	name := ctx.Param("name")
	if err := h.tailscale.ClearService(name); err != nil {
		if services.IsTailscaleNotInstalledError(err) {
			return ctx.Render(200, "tailscale_not_installed.html", nil)
		}
		return ctx.String(500, "Failed to delete service: "+err.Error())
	}
	return ctx.Redirect(303, "/")
}
