package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"core/utils"
)

func RegisterFrontendRoutes(router *chi.Mux, _ utils.Config) {
	var singlePageAppHandler http.Handler = utils.ServeSPA("frontend/apps/photography_crm/dist")
	router.Handle("/*", singlePageAppHandler)
}
