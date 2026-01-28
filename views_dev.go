//go:build !prod

package main

import (
	"html/template"
)

func parseTemplates() *template.Template {
	return template.Must(template.New("").
		Funcs(template.FuncMap{
			"viteTags": ViteTags,
		}).
		ParseGlob("views/*.html"))
}
