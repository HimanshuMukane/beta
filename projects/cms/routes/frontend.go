package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"cms/handlers"
	"cms/ui/static"
	"core/middleware"
	"core/utils"
)

func RegisterFrontendRoutes(r *chi.Mux, cfg utils.Config) {
	fs := http.StripPrefix("/static/", http.FileServer(static.FS()))
	r.Handle("/static/*", fs)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuthPage(cfg))
		r.Get("/", handlers.HomeHandler)
	})

	r.Get("/signin", handlers.SignInHandler)
	r.Get("/signup", handlers.SignUpHandler)
	r.Get("/forgot-password", handlers.ForgotPasswordHandler)
}
