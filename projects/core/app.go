package core

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := utils.AppConfig
	pool, err := utils.InitDBPool(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database pool: %v", err)
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
		fmt.Printf("Server running at http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-shutdownChan
	fmt.Println("\nShutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
