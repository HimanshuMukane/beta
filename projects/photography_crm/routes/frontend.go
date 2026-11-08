package routes

import (
	"github.com/go-chi/chi/v5"

	"core/utils"
)

func RegisterFrontendRoutes(r *chi.Mux, _ utils.Config) {
	spa := utils.ServeSPA("frontend/apps/photography_crm/dist")
	r.Handle("/*", spa)
}
