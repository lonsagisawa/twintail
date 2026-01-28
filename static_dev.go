//go:build !prod

package main

import (
	"io/fs"
	"os"
)

func getStaticFS() fs.FS {
	return os.DirFS("static")
}
