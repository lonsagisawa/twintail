//go:build !prod

package main

import (
	"html/template"
	"os"
)

func ViteTags(entry string) template.HTML {
	origin := os.Getenv("VITE_DEV_SERVER_URL")
	if origin == "" {
		origin = "http://localhost:5173"
	}

	s := `<script type="module" src="` + origin + `/@vite/client"></script>` + "\n" +
		`<script type="module" src="` + origin + `/` + entry + `"></script>`
	return template.HTML(s)
}
