package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmsRoutes "cms/routes"
	"core/utils"
	"github.com/joho/godotenv"
	crmRoutes "photography_crm/routes"
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
		portEnv       string
		defaultPort   string
		dbName        string
		cookieName    string
		routerBuilder func(pool *utils.DBPool, cfg utils.Config) http.Handler
	}{
		{
			name:          "CMS",
			portEnv:       "CMS_PORT",
			defaultPort:   "8081",
			dbName:        "cms_db",
			cookieName:    "cms_session",
			routerBuilder: cmsRoutes.SetupRouter,
		},
		{
			name:          "Photography CRM",
			portEnv:       "CRM_PORT",
			defaultPort:   "8082",
			dbName:        "photography_crm_db",
			cookieName:    "photography_session",
			routerBuilder: crmRoutes.SetupRouter,
		},
		{
			name:          "Superadmin",
			portEnv:       "SUPERADMIN_PORT",
			defaultPort:   "8083",
			dbName:        "superadmin_db",
			cookieName:    "superadmin_session",
			routerBuilder: superadminRoutes.SetupRouter,
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
		port := getEnv(a.portEnv, a.defaultPort)
		databaseURL := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbPassword, dbHost, dbPort, a.dbName)
		identityURL := fmt.Sprintf("root:%s@tcp(%s:%s)/platform_identity?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)

		cfg := utils.Config{
			AppEnv:            appEnv,
			ServerPort:        port,
			DatabaseURL:       databaseURL,
			IdentityDBURL:     identityURL,
			AllowTestDBHeader: appEnv != "production",
			TestDBHeaderKey:   "X-Database-Env",
			SessionCookieName: a.cookieName,
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

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
