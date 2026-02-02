package server

import (
	"twintail/internal/handlers"

	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Container) {
	e.GET("/", h.Service.Index)
	e.GET("/services/new", h.Service.Create)
	e.POST("/services/new", h.Service.Store)
	e.GET("/services/:name", h.Service.Show)
	e.GET("/services/:name/delete", h.Service.Delete)
	e.POST("/services/:name/delete", h.Service.Destroy)
	e.GET("/services/:name/endpoints/new", h.Endpoint.Create)
	e.POST("/services/:name/endpoints/new", h.Endpoint.Store)
	e.GET("/services/:name/endpoints/edit", h.Endpoint.Edit)
	e.POST("/services/:name/endpoints/edit", h.Endpoint.Update)
	e.GET("/services/:name/endpoints/delete", h.Endpoint.Delete)
	e.POST("/services/:name/endpoints/delete", h.Endpoint.Destroy)

	e.GET("/settings", h.Settings.Show)
	e.POST("/settings", h.Settings.Update)

	// Static files
	e.StaticFS("/static", GetStaticFS())
}
