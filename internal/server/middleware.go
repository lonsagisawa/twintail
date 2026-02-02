package server

import (
	"html/template"

	"twintail/internal/services"

	"github.com/labstack/echo/v5"
)

func I18nMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			var lang string
			if cookie, err := c.Cookie("lang"); err == nil && cookie.Value != "" {
				lang = services.NormalizeLang(cookie.Value)
			} else {
				acceptLang := c.Request().Header.Get("Accept-Language")
				lang = services.ParseAcceptLanguage(acceptLang)
			}
			c.Set("lang", lang)
			return next(c)
		}
	}
}

func LiveReloadMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("liveReloadScript", template.HTML(liveReloadScript()))
			return next(c)
		}
	}
}

func NoCacheMiddleware() echo.MiddlewareFunc {
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
