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

	"core/utils"
	"github.com/joho/godotenv"
	photoCrmRoutes "photography_crm/routes"
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
			name:          "Photography CRM",
			projectID:     "photography_crm",
			envFile:       "projects/photography_crm/.env",
			portEnv:       "CRM_PORT",
			defaultPort:   "8082",
			dbName:        "photography_crm_db",
			cookieName:    "photography_session",
			routerBuilder: photoCrmRoutes.SetupRouter,
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

		cfg := utils.Config{
			AppName:           a.projectID,
			AppEnv:            appEnv,
			ServerPort:        port,
			DatabaseURL:       databaseURL,
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
