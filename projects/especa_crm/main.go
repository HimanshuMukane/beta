package main

import (
	"context"
	"net/http"

	"core"
	"core/utils"
	crmDB "especa_crm/db"
	"especa_crm/routes"
)

func main() {
	core.Run(func(pool *utils.DBPool, cfg utils.Config) http.Handler {
		if pool.Primary != nil {
			ctx := context.Background()
			_ = crmDB.RunMigrations(ctx, pool.Primary, "projects/especa_crm/db/migrations")

			var companyID string
			_ = pool.Primary.GetContext(ctx, &companyID, "SELECT id FROM companies LIMIT 1")
			if companyID == "" {
				companyID = "comp_default_01"
			}
			_ = crmDB.SeedInitialCRM(ctx, pool.Primary, companyID)

			return routes.SetupRouterWithAuth(pool.Primary, companyID, cfg)
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		})
	})
}
