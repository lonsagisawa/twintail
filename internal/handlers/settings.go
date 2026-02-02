package handlers

import (
	"net/http"
	"time"

	"twintail/internal/requests"
	"twintail/internal/services"

	"github.com/labstack/echo/v5"
)

type SettingsHandler struct{}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

func (h *SettingsHandler) Show(ctx *echo.Context) error {
	currentLang := ctx.Get("lang").(string)
	return ctx.Render(http.StatusOK, "settings.html", map[string]any{
		"CurrentLang": currentLang,
		"Languages":   services.GetSupportedLanguages(),
	})
}

func (h *SettingsHandler) Update(ctx *echo.Context) error {
	var req requests.UpdateSettingsRequest
	if err := req.FromContext(ctx); err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/settings")
	}

	cookie := &http.Cookie{
		Name:     "lang",
		Value:    req.Lang,
		Path:     "/",
		MaxAge:   int(365 * 24 * time.Hour / time.Second),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	ctx.SetCookie(cookie)

	return ctx.Redirect(http.StatusSeeOther, "/settings")
}
