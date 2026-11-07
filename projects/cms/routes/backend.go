package routes

import (
	"github.com/go-chi/chi/v5"

	"cms/handlers"
	"cms/repositories"
	coreHandlers "core/handlers"
	"core/middleware"
	"core/utils"
)

func RegisterBackendRoutes(
	r *chi.Mux,
	pool *utils.DBPool,
	cfg utils.Config,
) {
	r.Get("/api/health", coreHandlers.HealthHandler(pool))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		r.Post("/signup", coreHandlers.SignupAPIHandler(cfg, "editor"))
		r.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))

		r.With(middleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
	})

	// Headless Developer Content Delivery API (v1)
	if pool != nil && pool.Primary != nil {
		cmsRepo := repositories.NewCMSRepository(pool.Primary)

		r.Route("/api/v1", func(v1 chi.Router) {
			v1.Use(handlers.DeveloperAPIKeyMiddleware(cmsRepo))

			// Custom Content Types & Entries
			v1.Group(func(contentRouter chi.Router) {
				contentRouter.Use(handlers.RequireFeatureMiddleware(cmsRepo, "content_types"))
				contentRouter.Get("/content/{type_slug}", handlers.GetContentEntriesHandler(cmsRepo))
				contentRouter.Get("/content/{type_slug}/{entry_slug}", handlers.GetContentEntryBySlugHandler(cmsRepo))
			})

			// Navigation Menus
			v1.Group(func(navRouter chi.Router) {
				navRouter.Use(handlers.RequireFeatureMiddleware(cmsRepo, "navigation"))
				navRouter.Get("/navigation/{menu_key}", handlers.GetMenuHandler(cmsRepo))
			})

			// Form Submissions from Custom Sites
			v1.Group(func(formRouter chi.Router) {
				formRouter.Use(handlers.RequireFeatureMiddleware(cmsRepo, "forms"))
				formRouter.Post("/forms/{form_key}", handlers.SubmitFormHandler(cmsRepo))
			})
		})
	}
}
