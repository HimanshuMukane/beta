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
		AppEnv:            strings.ToLower(strings.TrimSpace(getEnv("APP_ENV", "development"))),
		ServerPort:        strings.TrimSpace(getEnv("SERVER_PORT", "8080")),
		DBPath:            strings.TrimSpace(getEnv("DB_PATH", "photography_crm.db")),
		DBTestPath:        strings.TrimSpace(getEnv("DB_TEST_PATH", "photography_crm_test.db")),
		TestDBHeaderKey:   strings.TrimSpace(getEnv("TEST_DB_HEADER_KEY", "X-Database-Env")),
		SessionCookieName: strings.TrimSpace(getEnv("SESSION_COOKIE_NAME", "photography_session")),
		SessionSecret:     os.Getenv("SESSION_SECRET"),
		SessionTTLHours:   24,
		SessionSecure:     true,
	}

	// Validate the application environment.
	switch cfg.AppEnv {
	case "development", "test", "production":
	default:
		return fmt.Errorf("APP_ENV must be development, test, or production")
	}

	// Allow selecting the test database only outside production.
	if value := os.Getenv("ALLOW_TEST_DB_HEADER"); value != "" {
		allowed, err := strconv.ParseBool(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("ALLOW_TEST_DB_HEADER must be true or false")
		}
		cfg.AllowTestDBHeader = allowed
	} else {
		cfg.AllowTestDBHeader = cfg.AppEnv != "production"
	}

	if cfg.AppEnv == "production" && cfg.AllowTestDBHeader {
		return fmt.Errorf("ALLOW_TEST_DB_HEADER cannot be enabled in production")
	}

	// Validate the session lifetime.
	if value := os.Getenv("SESSION_TTL_HOURS"); value != "" {
		hours, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || hours < 1 || hours > 720 {
			return fmt.Errorf("SESSION_TTL_HOURS must be between 1 and 720")
		}
		cfg.SessionTTLHours = hours
	}

	// Configure whether session cookies require HTTPS.
	if value := os.Getenv("SESSION_COOKIE_SECURE"); value != "" {
		secure, err := strconv.ParseBool(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("SESSION_COOKIE_SECURE must be true or false")
		}
		cfg.SessionSecure = secure
	}

	if cfg.AppEnv == "production" && !cfg.SessionSecure {
		return fmt.Errorf("SESSION_COOKIE_SECURE must be true in production")
	}

	// A session secret is required for authentication.
	if strings.TrimSpace(cfg.SessionSecret) == "" {
		return fmt.Errorf("SESSION_SECRET is required")
	}

	if len(cfg.SessionSecret) < 64 {
		return fmt.Errorf("SESSION_SECRET must contain at least 64 characters")
	}

	// Validate database paths.
	if cfg.DBPath == "" {
		return fmt.Errorf("DB_PATH is required")
	}

	if cfg.DBTestPath == "" {
		return fmt.Errorf("DB_TEST_PATH is required")
	}

	if cfg.DBPath == cfg.DBTestPath {
		return fmt.Errorf("DB_PATH and DB_TEST_PATH must differ")
	}

	// Validate the server port.
	port, err := strconv.Atoi(cfg.ServerPort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
	}

	AppConfig = cfg
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
