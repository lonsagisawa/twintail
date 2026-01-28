package controllers

import (
	"twintail/services"

	"github.com/labstack/echo/v5"
)

type ServiceController struct {
	tailscale *services.TailscaleService
}

func NewServiceController(tailscale *services.TailscaleService) *ServiceController {
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
		"Services":         svcs,
		"LiveReloadScript": ctx.Get("liveReloadScript"),
	})
}
