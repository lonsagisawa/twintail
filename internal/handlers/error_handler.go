package handlers

import (
	"net/http"

	"twintail/internal/services"

	"github.com/labstack/echo/v5"
)

func HTTPErrorHandler(c *echo.Context, err error) {
	if services.IsTailscaleNotInstalledError(err) {
		if err := c.Render(http.StatusOK, "tailscale_not_installed.html", nil); err != nil {
			c.Logger().Error("render error", "error", err)
		}
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error("http error", "error", err)

	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}
	if c.Request().Method == http.MethodHead {
		_ = c.NoContent(code)
	} else {
		_ = c.String(code, err.Error())
	}
}
