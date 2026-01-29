package main

import (
	"html/template"
	"os"

	"twintail/controllers"
	"twintail/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func liveReloadMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("liveReloadScript", template.HTML(liveReloadScript()))
			return next(c)
		}
	}
}

func noCacheMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if noCacheEnabled() {
				c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				c.Response().Header().Set("Pragma", "no-cache")
				c.Response().Header().Set("Expires", "0")
			}
			return next(c)
		}
	}
}

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8077"
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(liveReloadMiddleware())
	e.Use(noCacheMiddleware())

	setupLiveReload(e)

	e.Renderer = parseTemplates()

	tailscaleSvc := services.NewTailscaleService()
	serviceCtrl := controllers.NewServiceController(tailscaleSvc)

	e.GET("/", serviceCtrl.Index)
	e.GET("/services/new", serviceCtrl.Create)
	e.POST("/services/new", serviceCtrl.Store)
	e.GET("/services/:name", serviceCtrl.Show)
	e.GET("/services/:name/delete", serviceCtrl.ConfirmDelete)
	e.POST("/services/:name/delete", serviceCtrl.Destroy)

	e.StaticFS("/static", getStaticFS())

	if err := e.Start(":" + port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
