//go:build prod

package main

import (
	"embed"
	"io/fs"

	"twintail/services"
)

//go:embed locales/*.json
var localesFS embed.FS

func loadI18n() *services.I18n {
	subFS, err := fs.Sub(localesFS, "locales")
	if err != nil {
		panic("failed to load locales: " + err.Error())
	}
	i18n, err := services.NewI18n(subFS, "en")
	if err != nil {
		panic("failed to load locales: " + err.Error())
	}
	return i18n
}
