package routes

import (
	"net/http"
	"strings"

	"chi-mysql-server/utils"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func DatabaseRoutingMiddleware(pool *utils.DBPool, cfg utils.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			useTest := false

			if cfg.AllowTestDBHeader && cfg.AppEnv != "production" {
				headerVal := strings.TrimSpace(r.Header.Get(cfg.TestDBHeaderKey))
				if strings.EqualFold(headerVal, "test") {
					useTest = true
				}
			}

			db := pool.Resolve(useTest)
			ctx := utils.WithDB(r.Context(), db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
