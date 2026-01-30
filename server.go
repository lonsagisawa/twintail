package main

import (
	"html/template"
	"os"

	"twintail/handlers"
	"twintail/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func i18nMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			acceptLang := c.Request().Header.Get("Accept-Language")
			lang := services.ParseAcceptLanguage(acceptLang)
			c.Set("lang", lang)
			return next(c)
		}
	}
}

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
	e.Use(i18nMiddleware())
	e.Use(liveReloadMiddleware())
	e.Use(noCacheMiddleware())

	setupLiveReload(e)

	e.Renderer = parseTemplates()
	e.Validator = NewCustomValidator()

	tailscaleSvc := services.NewTailscaleService()
	serviceHandler := handlers.NewServiceHandler(tailscaleSvc)
	endpointHandler := handlers.NewEndpointHandler(tailscaleSvc)

	e.GET("/", serviceHandler.Index)
	e.GET("/services/new", serviceHandler.Create)
	e.POST("/services/new", serviceHandler.Store)
	e.GET("/services/:name", serviceHandler.Show)
	e.GET("/services/:name/delete", serviceHandler.Delete)
	e.POST("/services/:name/delete", serviceHandler.Destroy)
	e.GET("/services/:name/endpoints/new", endpointHandler.Create)
	e.POST("/services/:name/endpoints/new", endpointHandler.Store)
	e.GET("/services/:name/endpoints/edit", endpointHandler.Edit)
	e.POST("/services/:name/endpoints/edit", endpointHandler.Update)
	e.GET("/services/:name/endpoints/delete", endpointHandler.Delete)
	e.POST("/services/:name/endpoints/delete", endpointHandler.Destroy)

	e.StaticFS("/static", getStaticFS())

	if err := e.Start(":" + port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
