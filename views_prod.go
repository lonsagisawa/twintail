//go:build prod

package main

import (
	"embed"
	"html/template"
)

//go:embed views/*.html
var viewsFS embed.FS

func parseTemplates() *template.Template {
	return template.Must(template.ParseFS(viewsFS, "views/*.html"))
}
