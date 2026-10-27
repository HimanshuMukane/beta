package static

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed css/*
var staticFS embed.FS

// FS returns an http.FileSystem for serving static assets embedded in the binary.
func FS() http.FileSystem {
	sub, err := fs.Sub(staticFS, ".")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}
