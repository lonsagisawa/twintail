//go:build prod

package main

import (
	"embed"
	"html/template"
)

//go:embed views/*.html
var viewsFS embed.FS

func parseTemplates() *template.Template {
	return template.Must(template.New("").
		Funcs(template.FuncMap{
			"viteTags": ViteTags,
		}).
		ParseFS(viewsFS, "views/*.html"))
}
