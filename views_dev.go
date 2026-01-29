//go:build !prod

package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v5"
)

type TemplateRenderer struct {
	baseTemplate *template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	base := template.Must(template.New("").
		Funcs(template.FuncMap{
			"viteTags": ViteTags,
		}).
		ParseGlob("views/layouts/*.html"))
	template.Must(base.ParseGlob("views/partials/*.html"))

	return &TemplateRenderer{baseTemplate: base}
}

func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	if m, ok := data.(map[string]any); ok {
		m["LiveReloadScript"] = c.Get("liveReloadScript")
	}
	tmpl := template.Must(template.Must(t.baseTemplate.Clone()).ParseFiles("views/" + name))
	return tmpl.ExecuteTemplate(w, "base", data)
}

func parseTemplates() *TemplateRenderer {
	return NewTemplateRenderer()
}
