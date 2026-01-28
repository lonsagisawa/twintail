//go:build prod

package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v5"
)

//go:embed views/layouts/*.html views/*.html
var viewsFS embed.FS

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	funcs := template.FuncMap{
		"viteTags": ViteTags,
	}

	base := template.Must(template.New("").Funcs(funcs).ParseFS(viewsFS, "views/layouts/*.html"))

	templates := make(map[string]*template.Template)
	pages := []string{"index.html", "new_service.html"}

	for _, page := range pages {
		tmpl := template.Must(template.Must(base.Clone()).ParseFS(viewsFS, "views/"+page))
		templates[page] = tmpl
	}

	return &TemplateRenderer{templates: templates}
}

func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	if m, ok := data.(map[string]any); ok {
		m["LiveReloadScript"] = c.Get("liveReloadScript")
	}
	return t.templates[name].ExecuteTemplate(w, "base", data)
}

func parseTemplates() *TemplateRenderer {
	return NewTemplateRenderer()
}
