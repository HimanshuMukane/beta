package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	ServerPort string

	DBPath     string
	DBTestPath string

	AllowTestDBHeader bool
	TestDBHeaderKey   string

	SessionCookieName string
	SessionTTLHours   int
	SessionSecure     bool
	SessionSecret     string
}

var AppConfig Config

func LoadConfig() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load .env: %w", err)
	}

	cfg := Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DBPath:            getEnv("DB_PATH", "photography_crm.db"),
		DBTestPath:        getEnv("DB_TEST_PATH", "photography_crm_test.db"),
		TestDBHeaderKey:   getEnv("TEST_DB_HEADER_KEY", "X-Database-Env"),
		SessionCookieName: getEnv("SESSION_COOKIE_NAME", "photography_session"),
		SessionSecret:     os.Getenv("SESSION_SECRET"),
		SessionTTLHours:   24,
		SessionSecure:     true,
	}

	if value := os.Getenv("ALLOW_TEST_DB_HEADER"); value != "" {
		allowed, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("ALLOW_TEST_DB_HEADER must be true or false")
		}
		cfg.AllowTestDBHeader = allowed
	} else {
		// Enabled by default in non-production environments
		cfg.AllowTestDBHeader = (cfg.AppEnv != "production")
	}

	if cfg.AppEnv == "production" && cfg.AllowTestDBHeader {
		return fmt.Errorf("ALLOW_TEST_DB_HEADER cannot be enabled in production environment")
	}

	if value := os.Getenv("SESSION_TTL_HOURS"); value != "" {
		hours, err := strconv.Atoi(value)
		if err != nil || hours < 1 || hours > 720 {
			return fmt.Errorf("SESSION_TTL_HOURS must be between 1 and 720")
		}
		cfg.SessionTTLHours = hours
	}

	if value := os.Getenv("SESSION_COOKIE_SECURE"); value != "" {
		secure, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("SESSION_COOKIE_SECURE must be true or false")
		}
		cfg.SessionSecure = secure
	}

	if len(cfg.SessionSecret) > 0 && len(cfg.SessionSecret) < 64 {
		return fmt.Errorf("SESSION_SECRET must contain at least 64 characters if set")
	}

	if cfg.AppEnv == "production" && !cfg.SessionSecure {
		return fmt.Errorf("SESSION_COOKIE_SECURE must be true in production")
	}

	if strings.TrimSpace(cfg.DBPath) == "" {
		return fmt.Errorf("DB_PATH is required")
	}

	if err := validatePort("SERVER_PORT", cfg.ServerPort); err != nil {
		return err
	}

	AppConfig = cfg
	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func validatePort(name, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("%s must be a valid port between 1 and 65535", name)
	}
	return nil
}
