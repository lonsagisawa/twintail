//go:build !prod

package server

import (
	"io/fs"
	"os"
)

func GetStaticFS() fs.FS {
	return os.DirFS("internal/server/static")
}

func noCacheEnabled() bool {
	return true
}
