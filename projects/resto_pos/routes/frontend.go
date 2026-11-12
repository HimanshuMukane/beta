package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"core/utils"
)

func RegisterFrontendRoutes(router *chi.Mux, _ utils.Config) {
	var singlePageAppHandler http.Handler = utils.ServeSPA("frontend/apps/resto_pos/dist")
	router.Handle("/*", singlePageAppHandler)
}
