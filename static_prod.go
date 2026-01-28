//go:build prod

package main

import (
	"embed"
	"io/fs"
)

//go:embed static all:static/dist
var staticFS embed.FS

func getStaticFS() fs.FS {
	sub, _ := fs.Sub(staticFS, "static")
	return sub
}

func noCacheEnabled() bool {
	return false
}
