package routes

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	coreHandlers "core/handlers"
	"core/middleware"
	"core/utils"
	"superadmin/handlers"
)

func RegisterBackendRoutes(
	r *chi.Mux,
	pool *utils.DBPool,
	cfg utils.Config,
) {
	r.Get("/api/health", coreHandlers.HealthHandler(pool))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		r.Post("/signup", coreHandlers.SignupAPIHandler(cfg, "superadmin"))
		r.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))

		r.With(middleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
	})

	// Superadmin Governance over Tenant Features
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "especa_root_password"
	}
	cmsDSN := fmt.Sprintf("root:%s@tcp(%s:%s)/cms_db?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)
	cmsDB, err := sqlx.Open("mysql", cmsDSN)
	if err == nil {
		r.Route("/api/tenants/{tenant_id}/cms-features", func(tr chi.Router) {
			tr.Use(middleware.RequireAuth(cfg))
			tr.Get("/", handlers.GetTenantCMSFeaturesHandler(cmsDB))
			tr.Post("/", handlers.ToggleTenantCMSFeatureHandler(cmsDB))
		})
	}
}
