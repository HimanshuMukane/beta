package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chi-mysql-server/handlers"
	"chi-mysql-server/ui/static"
)

func RegisterFrontendRoutes(r chi.Router) {
	fileServer(r, "/static", static.FS())

	r.Get("/", handlers.HomeHandler)
	r.Get("/signin", handlers.SignInHandler)
	r.Get("/signup", handlers.SignUpHandler)
	r.Get("/forgot-password", handlers.ForgotPasswordHandler)
	r.NotFound(handlers.NotFoundHandler)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := rctx.RoutePattern()
		if len(pathPrefix) > 2 && pathPrefix[len(pathPrefix)-2:] == "/*" {
			pathPrefix = pathPrefix[:len(pathPrefix)-2]
		}
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
