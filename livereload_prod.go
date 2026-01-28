//go:build prod

package main

import "github.com/labstack/echo/v5"

func setupLiveReload(e *echo.Echo) {}

func liveReloadScript() string {
	return ""
}
