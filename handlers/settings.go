package handlers

import (
	"net/http"
	"time"

	"twintail/services"

	"github.com/labstack/echo/v5"
)

type SettingsHandler struct{}

func NewSettingsHandler() *SettingsHandler {
	return &SettingsHandler{}
}

func (h *SettingsHandler) Show(ctx *echo.Context) error {
	currentLang := ctx.Get("lang").(string)
	return ctx.Render(200, "settings.html", map[string]any{
		"CurrentLang": currentLang,
		"Languages":   services.GetSupportedLanguages(),
	})
}

type SettingsFormData struct {
	Lang string `form:"lang" validate:"required,oneof=en ja"`
}

func (h *SettingsHandler) Update(ctx *echo.Context) error {
	var formData SettingsFormData
	if err := ctx.Bind(&formData); err != nil {
		return ctx.Redirect(303, "/settings")
	}
	if err := ctx.Validate(&formData); err != nil {
		return ctx.Redirect(303, "/settings")
	}

	cookie := &http.Cookie{
		Name:     "lang",
		Value:    formData.Lang,
		Path:     "/",
		MaxAge:   int(365 * 24 * time.Hour / time.Second),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	ctx.SetCookie(cookie)

	return ctx.Redirect(303, "/settings")
}
