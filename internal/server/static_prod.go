//go:build prod

package server

import (
	"embed"
	"io/fs"
)

//go:embed static all:static/dist
var staticFS embed.FS

func GetStaticFS() fs.FS {
	sub, _ := fs.Sub(staticFS, "static")
	return sub
}

func noCacheEnabled() bool {
	return false
}
