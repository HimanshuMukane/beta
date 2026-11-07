package core

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

	"core/utils"
)

type RouterBuilder func(pool *utils.DBPool, cfg utils.Config) http.Handler

func Run(setupRouter RouterBuilder) {
	if err := utils.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	cfg := utils.AppConfig
	pool, err := utils.InitDBPool(cfg)
	if err != nil {
		slog.Error("Failed to initialize database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	router := setupRouter(pool, cfg)
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server running", "port", port, "env", cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdownChan
	fmt.Println("\nShutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}
}
