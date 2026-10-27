package routes

import (
	"net/http"

	"chi-mysql-server/handlers"
	"chi-mysql-server/ui/static"
)

func RegisterFrontendRoutes(mux *http.ServeMux) {
	fileServer(mux, "/static/", static.FS())

	mux.HandleFunc("GET /{$}", handlers.HomeHandler)
	mux.HandleFunc("GET /signin", handlers.SignInHandler)
	mux.HandleFunc("GET /signup", handlers.SignUpHandler)
	mux.HandleFunc("GET /forgot-password", handlers.ForgotPasswordHandler)
	mux.HandleFunc("/", handlers.NotFoundHandler)
}

func fileServer(mux *http.ServeMux, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))
	mux.Handle("GET "+path, fs)
}
