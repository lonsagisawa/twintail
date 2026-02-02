//go:build !prod

package services

import (
	"os"
)

func LoadI18n() *I18n {
	i18n, err := NewI18n(os.DirFS("internal/services/locales"), "en")
	if err != nil {
		panic("failed to load locales: " + err.Error())
	}
	return i18n
}
