package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	cmsRoutes "cms/routes"
	"core/utils"
	crmRoutes "especa_crm/routes"
	financeRoutes "especa_finance/routes"
	"github.com/joho/godotenv"
	hrmsRoutes "hrms/routes"
	photoCrmRoutes "photography_crm/routes"
	superadminRoutes "superadmin/routes"
)

type ServiceApp struct {
	Name   string
	Port   string
	Pool   *utils.DBPool
	Server *http.Server
}

func main() {
	_ = godotenv.Load()

	apps := []struct {
		name          string
		projectID     string
		envFile       string
		portEnv       string
		defaultPort   string
		dbName        string
		cookieName    string
		routerBuilder func(pool *utils.DBPool, cfg utils.Config) http.Handler
	}{
		{
			name:          "CMS",
			projectID:     "cms",
			envFile:       "projects/cms/.env",
			portEnv:       "CMS_PORT",
			defaultPort:   "8081",
			dbName:        "cms_db",
			cookieName:    "cms_session",
			routerBuilder: cmsRoutes.SetupRouter,
		},
		{
			name:          "Photography CRM",
			projectID:     "photography_crm",
			envFile:       "projects/photography_crm/.env",
			portEnv:       "CRM_PORT",
			defaultPort:   "8082",
			dbName:        "photography_crm_db",
			cookieName:    "photography_session",
			routerBuilder: photoCrmRoutes.SetupRouter,
		},
		{
			name:          "Superadmin",
			projectID:     "superadmin",
			envFile:       "projects/superadmin/.env",
			portEnv:       "SUPERADMIN_PORT",
			defaultPort:   "8083",
			dbName:        "superadmin_db",
			cookieName:    "superadmin_session",
			routerBuilder: superadminRoutes.SetupRouter,
		},
		{
			name:          "Finance",
			projectID:     "especa_finance",
			envFile:       "projects/especa_finance/.env",
			portEnv:       "FINANCE_PORT",
			defaultPort:   "8084",
			dbName:        "finance_db",
			cookieName:    "finance_session",
			routerBuilder: func(pool *utils.DBPool, cfg utils.Config) http.Handler {
				if pool.Primary != nil {
					return financeRoutes.SetupRouterWithAuth(pool.Primary, "comp_default_01", cfg)
				}
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
				})
			},
		},
		{
			name:          "Especa CRM",
			projectID:     "especa_crm",
			envFile:       "projects/especa_crm/.env",
			portEnv:       "ESPECA_CRM_PORT",
			defaultPort:   "8086",
			dbName:        "especa_crm_db",
			cookieName:    "crm_session",
			routerBuilder: func(pool *utils.DBPool, cfg utils.Config) http.Handler {
				if pool.Primary != nil {
					return crmRoutes.SetupRouterWithAuth(pool.Primary, "comp_default_01", cfg)
				}
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
				})
			},
		},
		{
			name:          "HRMS",
			projectID:     "hrms",
			envFile:       "projects/hrms/.env",
			portEnv:       "HRMS_PORT",
			defaultPort:   "8087",
			dbName:        "hrms_db",
			cookieName:    "hrms_session",
			routerBuilder: hrmsRoutes.SetupRouter,
		},
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "0123456789012345678901234567890123456789012345678901234567890123"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

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

	activeServices := make([]*ServiceApp, 0, len(apps))
	defer func() {
		for _, s := range activeServices {
			if s.Pool != nil {
				_ = s.Pool.Close()
			}
		}
	}()

	for _, a := range apps {
		// Read service-specific .env file if available
		localEnv, _ := godotenv.Read(filepath.Clean(a.envFile))

		port := os.Getenv(a.portEnv)
		if port == "" && localEnv != nil {
			port = localEnv["SERVER_PORT"]
			if port == "" {
				port = localEnv["PORT"]
			}
		}
		if port == "" {
			port = a.defaultPort
		}

		dbName := a.dbName
		if localEnv != nil && localEnv["DB_NAME"] != "" {
			dbName = localEnv["DB_NAME"]
		}

		dbUser := "root"
		if localEnv != nil && localEnv["DB_USER"] != "" {
			dbUser = localEnv["DB_USER"]
		}

		serviceDbPassword := dbPassword
		if localEnv != nil && localEnv["DB_PASSWORD"] != "" {
			serviceDbPassword = localEnv["DB_PASSWORD"]
		}

		cookieName := a.cookieName
		if localEnv != nil && localEnv["SESSION_COOKIE_NAME"] != "" {
			cookieName = localEnv["SESSION_COOKIE_NAME"]
		}

		databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbUser, serviceDbPassword, dbHost, dbPort, dbName)
		identityURL := fmt.Sprintf("root:%s@tcp(%s:%s)/platform_identity?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)

		cfg := utils.Config{
			AppName:           a.projectID,
			AppEnv:            appEnv,
			ServerPort:        port,
			DatabaseURL:       databaseURL,
			IdentityDBURL:     identityURL,
			AllowTestDBHeader: appEnv != "production",
			TestDBHeaderKey:   "X-Database-Env",
			SessionCookieName: cookieName,
			SessionSecret:     sessionSecret,
			SessionTTLHours:   24,
			SessionSecure:     appEnv == "production",
		}

		pool, err := utils.InitDBPool(cfg)
		if err != nil {
			slog.Error("Database pool initialization failed", "service", a.name, "error", err)
			continue
		}

		router := a.routerBuilder(pool, cfg)
		server := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}

		service := &ServiceApp{
			Name:   a.name,
			Port:   port,
			Pool:   pool,
			Server: server,
		}
		activeServices = append(activeServices, service)
	}

	for _, s := range activeServices {
		go func(srv *ServiceApp) {
			slog.Info("Service listening", "service", srv.Name, "port", srv.Port)
			if err := srv.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Service error", "service", srv.Name, "error", err)
			}
		}(s)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down multi-listener services...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, s := range activeServices {
		if err := s.Server.Shutdown(ctx); err != nil {
			slog.Error("Shutdown error", "service", s.Name, "error", err)
		}
	}
	slog.Info("All services stopped.")
}
