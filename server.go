package main

import (
	"html/template"
	"io"
	"os"

	"twintail/controllers"
	"twintail/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(c *echo.Context, w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func liveReloadMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("liveReloadScript", template.HTML(liveReloadScript()))
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

	setupLiveReload(e)

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	tailscaleSvc := services.NewTailscaleService()
	serviceCtrl := controllers.NewServiceController(tailscaleSvc)

	e.GET("/", serviceCtrl.Index)

	e.StaticFS("/static", getStaticFS())

	if err := e.Start(":" + port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
