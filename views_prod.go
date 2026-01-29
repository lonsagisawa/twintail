//go:build prod

package main

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/labstack/echo/v5"
)

//go:embed views/layouts/*.html views/partials/*.html views/*.html
var viewsFS embed.FS

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	funcs := template.FuncMap{
		"viteTags": ViteTags,
	}

	base := template.Must(template.New("").Funcs(funcs).ParseFS(viewsFS, "views/layouts/*.html", "views/partials/*.html"))

	templates := make(map[string]*template.Template)

	entries, err := fs.ReadDir(viewsFS, "views")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}
		page := entry.Name()
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
