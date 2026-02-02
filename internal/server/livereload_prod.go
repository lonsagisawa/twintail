//go:build prod

package server

import "github.com/labstack/echo/v5"

func SetupLiveReload(e *echo.Echo) {}

func liveReloadScript() string {
	return ""
}
