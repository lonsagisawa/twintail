//go:build prod

package services

import (
	"embed"
	"io/fs"
)

//go:embed locales/*.json
var localesFS embed.FS

func LoadI18n() *I18n {
	subFS, err := fs.Sub(localesFS, "locales")
	if err != nil {
		panic("failed to load locales: " + err.Error())
	}
	i18n, err := NewI18n(subFS, "en")
	if err != nil {
		panic("failed to load locales: " + err.Error())
	}
	return i18n
}
