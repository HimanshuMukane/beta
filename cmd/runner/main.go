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

type AppDefinition struct {
	name          string
	projectID     string
	envFile       string
	portEnv       string
	defaultPort   string
	dbName        string
	cookieName    string
	routerBuilder func(pool *utils.DBPool, cfg utils.Config) http.Handler
}

func main() {
	_ = godotenv.Load()

	var appDefinitions []AppDefinition = []AppDefinition{
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

	var sessionSecret string = os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "0123456789012345678901234567890123456789012345678901234567890123"
	}

	var appEnv string = os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	var dbHost string = os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	var dbPort string = os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	var dbPassword string = os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "especa_root_password"
	}

	var activeServices []*ServiceApp = make([]*ServiceApp, 0, len(appDefinitions))
	defer func() {
		var serviceIndex int
		for serviceIndex = 0; serviceIndex < len(activeServices); serviceIndex++ {
			var serviceInstance *ServiceApp = activeServices[serviceIndex]
			if serviceInstance.Pool != nil {
				_ = serviceInstance.Pool.Close()
			}
		}
	}()

	var appIndex int
	for appIndex = 0; appIndex < len(appDefinitions); appIndex++ {
		var appDef AppDefinition = appDefinitions[appIndex]
		var localEnv map[string]string
		localEnv, _ = godotenv.Read(filepath.Clean(appDef.envFile))

		var port string = os.Getenv(appDef.portEnv)
		if port == "" && localEnv != nil {
			port = localEnv["SERVER_PORT"]
			if port == "" {
				port = localEnv["PORT"]
			}
		}
		if port == "" {
			port = appDef.defaultPort
		}

		var dbName string = appDef.dbName
		if localEnv != nil && localEnv["DB_NAME"] != "" {
			dbName = localEnv["DB_NAME"]
		}

		var dbUser string = "root"
		if localEnv != nil && localEnv["DB_USER"] != "" {
			dbUser = localEnv["DB_USER"]
		}

		var serviceDbPassword string = dbPassword
		if localEnv != nil && localEnv["DB_PASSWORD"] != "" {
			serviceDbPassword = localEnv["DB_PASSWORD"]
		}

		var cookieName string = appDef.cookieName
		if localEnv != nil && localEnv["SESSION_COOKIE_NAME"] != "" {
			cookieName = localEnv["SESSION_COOKIE_NAME"]
		}

		var databaseURL string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbUser, serviceDbPassword, dbHost, dbPort, dbName)

		var config utils.Config = utils.Config{
			AppName:           appDef.projectID,
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

		var dbPool *utils.DBPool
		var initErr error
		dbPool, initErr = utils.InitDBPool(config)
		if initErr != nil {
			slog.Error("Database pool initialization failed", "service", appDef.name, "error", initErr)
			continue
		}

		var router http.Handler = appDef.routerBuilder(dbPool, config)
		var serverInstance *http.Server = &http.Server{
			Addr:              ":" + port,
			Handler:           router,
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		}

		var serviceInstance *ServiceApp = &ServiceApp{
			Name:   appDef.name,
			Port:   port,
			Pool:   dbPool,
			Server: serverInstance,
		}
		activeServices = append(activeServices, serviceInstance)
	}

	var activeIndex int
	for activeIndex = 0; activeIndex < len(activeServices); activeIndex++ {
		var serviceItem *ServiceApp = activeServices[activeIndex]
		go func(srv *ServiceApp) {
			slog.Info("Service listening", "service", srv.Name, "port", srv.Port)
			var listenErr error = srv.Server.ListenAndServe()
			if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
				slog.Error("Service error", "service", srv.Name, "error", listenErr)
			}
		}(serviceItem)
	}

	var shutdownSignalChannel chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(shutdownSignalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignalChannel

	slog.Info("Shutting down multi-listener services...")
	var shutdownContext context.Context
	var cancelShutdown context.CancelFunc
	shutdownContext, cancelShutdown = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	var shutdownIndex int
	for shutdownIndex = 0; shutdownIndex < len(activeServices); shutdownIndex++ {
		var serviceInstance *ServiceApp = activeServices[shutdownIndex]
		var shutdownErr error = serviceInstance.Server.Shutdown(shutdownContext)
		if shutdownErr != nil {
			slog.Error("Shutdown error", "service", serviceInstance.Name, "error", shutdownErr)
		}
	}
	slog.Info("All services stopped.")
}
