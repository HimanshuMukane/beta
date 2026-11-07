package main

import (
	"context"
	"net/http"

	"core"
	"core/utils"
	"especa_finance/routes"
)

func main() {
	core.Run(func(pool *utils.DBPool, cfg utils.Config) http.Handler {
		if pool.Primary != nil {
			var companyID string
			ctx := context.Background()
			_ = pool.Primary.GetContext(ctx, &companyID, "SELECT id FROM companies LIMIT 1")
			if companyID == "" {
				companyID = "comp_default_01"
			}
			return routes.SetupRouterWithAuth(pool.Primary, companyID, cfg)
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		})
	})
}
