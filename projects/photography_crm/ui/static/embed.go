package static

import (
	"embed"
	"io/fs"
	"net/http"
)

var staticFS embed.FS

func FS() http.FileSystem {
	sub, err := fs.Sub(staticFS, ".")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}
