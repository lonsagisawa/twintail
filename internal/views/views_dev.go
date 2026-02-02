//go:build !prod

package views

import (
	"html/template"
	"io"

	"twintail/internal/services"

	"github.com/labstack/echo/v5"
)

type TemplateRenderer struct {
	baseTemplate *template.Template
	i18n         *services.I18n
}

func NewTemplateRenderer(i18n *services.I18n) *TemplateRenderer {
	base := template.Must(template.New("").
		Funcs(template.FuncMap{
			"viteTags": ViteTags,
			"t":        func(key string) string { return key },
		}).
		ParseGlob("internal/views/views/layouts/*.html"))
	template.Must(base.ParseGlob("internal/views/views/partials/*.html"))

	return &TemplateRenderer{baseTemplate: base, i18n: i18n}
}

func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	lang := "en"
	if l, ok := c.Get("lang").(string); ok {
		lang = l
	}

	translator := t.i18n.GetTranslator(lang)

	tmpl := template.Must(template.Must(t.baseTemplate.Clone()).Funcs(template.FuncMap{
		"t": translator,
	}).ParseFiles("internal/views/views/" + name))

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
