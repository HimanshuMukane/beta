package routes

import (
	"github.com/go-chi/chi/v5"

	"core/utils"
)

func RegisterFrontendRoutes(r *chi.Mux, _ utils.Config) {
	// Serve React CSR built assets from frontend/apps/superadmin/dist
	spa := utils.ServeSPA("frontend/apps/superadmin/dist")
	r.Handle("/*", spa)
}
