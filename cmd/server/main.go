package main

import (
	"twintail/internal/config"
	"twintail/internal/handlers"
	"twintail/internal/server"
	"twintail/internal/services"
	"twintail/internal/validator"
	"twintail/internal/views"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(server.I18nMiddleware())
	e.Use(server.LiveReloadMiddleware())
	e.Use(server.NoCacheMiddleware())

	// Set custom HTTP error handler
	e.HTTPErrorHandler = handlers.HTTPErrorHandler

	server.SetupLiveReload(e)

	e.Renderer = views.ParseTemplates()
	e.Validator = validator.NewCustomValidator()

	tailscaleSvc := services.NewTailscaleService()
	container := handlers.NewContainer(tailscaleSvc)

	server.RegisterRoutes(e, container)

	if err := e.Start(":" + cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
