//go:build prod

package views

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"twintail/internal/services"

	"github.com/labstack/echo/v5"
)

//go:embed views/layouts/*.html views/partials/*.html views/*.html
var viewsFS embed.FS

type TemplateRenderer struct {
	templates map[string]*template.Template
	i18n      *services.I18n
}

func NewTemplateRenderer(i18n *services.I18n) *TemplateRenderer {
	funcs := template.FuncMap{
		"viteTags": ViteTags,
		"t":        func(key string) string { return key },
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

	return &TemplateRenderer{templates: templates, i18n: i18n}
}

func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	lang := "en"
	if l, ok := c.Get("lang").(string); ok {
		lang = l
	}

	translator := t.i18n.GetTranslator(lang)

	tmpl := template.Must(t.templates[name].Clone())
	tmpl = tmpl.Funcs(template.FuncMap{
		"t": translator,
	})

	if m, ok := data.(map[string]any); ok {
		m["LiveReloadScript"] = c.Get("liveReloadScript")
		m["Lang"] = lang
	}

	return tmpl.ExecuteTemplate(w, "base", data)
}

var globalI18n *services.I18n

func ParseTemplates() *TemplateRenderer {
	globalI18n = services.LoadI18n()
	return NewTemplateRenderer(globalI18n)
}
