package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	crmRoutes "chi-mysql-server/routes"
	cmsRoutes "cms/routes"
	"core/utils"
	"github.com/joho/godotenv"
	superadminRoutes "superadmin/routes"
)

type ServiceApp struct {
	Name    string
	Port    string
	DBPath  string
	Pool    *utils.DBPool
	Server  *http.Server
}

func main() {
	_ = godotenv.Load()

	apps := []struct {
		name          string
		portEnv       string
		defaultPort   string
		dbEnv         string
		defaultDB     string
		cookieName    string
		routerBuilder func(pool *utils.DBPool, cfg utils.Config) http.Handler
	}{
		{
			name:          "CMS",
			portEnv:       "CMS_PORT",
			defaultPort:   "8081",
			dbEnv:         "CMS_DB_PATH",
			defaultDB:     "cms.db",
			cookieName:    "cms_session",
			routerBuilder: cmsRoutes.SetupRouter,
		},
		{
			name:          "Photography CRM",
			portEnv:       "CRM_PORT",
			defaultPort:   "8082",
			dbEnv:         "CRM_DB_PATH",
			defaultDB:     "photography_crm.db",
			cookieName:    "photography_session",
			routerBuilder: crmRoutes.SetupRouter,
		},
		{
			name:          "Superadmin",
			portEnv:       "SUPERADMIN_PORT",
			defaultPort:   "8083",
			dbEnv:         "SUPERADMIN_DB_PATH",
			defaultDB:     "superadmin.db",
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
		dbPath := getEnv(a.dbEnv, a.defaultDB)
		dbTestPath := filepath.Join(filepath.Dir(dbPath), "test_"+filepath.Base(dbPath))

		cfg := utils.Config{
			AppEnv:            appEnv,
			ServerPort:        port,
			DBPath:            dbPath,
			DBTestPath:        dbTestPath,
			AllowTestDBHeader: appEnv != "production",
			TestDBHeaderKey:   "X-Database-Env",
			SessionCookieName: a.cookieName,
			SessionSecret:     sessionSecret,
			SessionTTLHours:   24,
			SessionSecure:     appEnv == "production",
		}

		pool, err := utils.InitDBPool(cfg)
		if err != nil {
			panic(fmt.Sprintf("[%s] DB initialization failed: %v", a.name, err))
		}

		router := a.routerBuilder(pool, cfg)
		server := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}

		service := &ServiceApp{
			Name:   a.name,
			Port:   port,
			DBPath: dbPath,
			Pool:   pool,
			Server: server,
		}
		activeServices = append(activeServices, service)
	}

	for _, s := range activeServices {
		go func(srv *ServiceApp) {
			fmt.Printf("➜ [%s] listening on http://localhost:%s (db: %s)\n", srv.Name, srv.Port, srv.DBPath)
			if err := srv.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("[%s] server error: %v\n", srv.Name, err)
			}
		}(s)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down multi-listener servers...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, s := range activeServices {
		if err := s.Server.Shutdown(ctx); err != nil {
			fmt.Printf("[%s] shutdown error: %v\n", s.Name, err)
		}
	}
	fmt.Println("All services stopped.")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

