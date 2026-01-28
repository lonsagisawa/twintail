//go:build !prod

package main

import (
	"html/template"
)

func parseTemplates() *template.Template {
	return template.Must(template.ParseGlob("views/*.html"))
}
