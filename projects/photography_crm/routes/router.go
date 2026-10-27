package routes

import (
	"net/http"

	"chi-mysql-server/utils"
)

func SetupRouter(pool *utils.DBPool, cfg utils.Config) http.Handler {
	mux := http.NewServeMux()

	RegisterFrontendRoutes(mux)
	RegisterBackendRoutes(mux, pool)

	return Chain(mux, DatabaseRoutingMiddleware(pool, cfg))
}
