//go:build !prod

package main

import (
	"os"

	"twintail/services"
)

func loadI18n() *services.I18n {
	i18n, err := services.NewI18n(os.DirFS("locales"), "en")
	if err != nil {
		panic("failed to load locales: " + err.Error())
	}
	return i18n
}
