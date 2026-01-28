//go:build prod

package main

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"sync"
)

type manifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css"`
}

var (
	manifestOnce sync.Once
	manifestMap  map[string]manifestEntry
)

func loadManifest() {
	b, err := fs.ReadFile(getStaticFS(), "dist/.vite/manifest.json")
	if err != nil {
		log.Fatalf("failed to read vite manifest: %v", err)
	}
	m := map[string]manifestEntry{}
	if err := json.Unmarshal(b, &m); err != nil {
		log.Fatalf("failed to parse vite manifest: %v", err)
	}
	manifestMap = m
}

func ViteTags(entry string) template.HTML {
	manifestOnce.Do(loadManifest)

	me, ok := manifestMap[entry]
	if !ok {
		return ""
	}

	out := ""
	for _, css := range me.CSS {
		out += `<link rel="stylesheet" href="/static/dist/` + css + `">` + "\n"
	}
	out += `<script type="module" src="/static/dist/` + me.File + `"></script>`
	return template.HTML(out)
}
