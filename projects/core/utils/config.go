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
	DatabaseURL string

	AllowTestDBHeader bool
	TestDBHeaderKey   string

	SessionCookieName string
	SessionTTLHours   int
	SessionSecure     bool
	SessionSecret     string

	SupabaseURL       string
	SupabaseAnonKey   string
	SupabaseJWTSecret string
	StorageEndpoint   string
	StorageBucket     string
	StorageAccessKey  string
	StorageSecretKey  string
}

var AppConfig Config

func LoadConfig(filenames ...string) error {
	if len(filenames) > 0 {
		if err := godotenv.Load(filenames...); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("load .env: %w", err)
		}
	} else {
		_ = godotenv.Load(".env", "projects/photography_crm/.env", "projects/cms/.env", "projects/superadmin/.env")
	}


	cfg := Config{
		AppEnv:            strings.ToLower(strings.TrimSpace(getEnv("APP_ENV", "development"))),
		ServerPort:        strings.TrimSpace(getEnv("SERVER_PORT", "8080")),
		DBPath:            strings.TrimSpace(getEnv("DB_PATH", "app.db")),
		DBTestPath:        strings.TrimSpace(getEnv("DB_TEST_PATH", "app_test.db")),
		TestDBHeaderKey:   strings.TrimSpace(getEnv("TEST_DB_HEADER_KEY", "X-Database-Env")),
		SessionCookieName: strings.TrimSpace(getEnv("SESSION_COOKIE_NAME", "app_session")),
		SessionSecret:     os.Getenv("SESSION_SECRET"),
		SessionTTLHours:   24,
		SessionSecure:     true,

		DatabaseURL:       strings.TrimSpace(os.Getenv("DATABASE_URL")),
		SupabaseURL:       strings.TrimSpace(getEnv("SUPABASE_URL", "http://localhost:8000")),
		SupabaseAnonKey:   strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY")),
		SupabaseJWTSecret: strings.TrimSpace(os.Getenv("SUPABASE_JWT_SECRET")),
		StorageEndpoint:   strings.TrimSpace(getEnv("STORAGE_ENDPOINT", "http://localhost:5000")),
		StorageBucket:     strings.TrimSpace(getEnv("STORAGE_BUCKET", "crm-storage")),
		StorageAccessKey:  strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID")),
		StorageSecretKey:  strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY")),
	}

	switch cfg.AppEnv {
	case "development", "test", "production":
	default:
		return fmt.Errorf("APP_ENV must be development, test, or production")
	}

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

	if value := os.Getenv("SESSION_TTL_HOURS"); value != "" {
		hours, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || hours < 1 || hours > 720 {
			return fmt.Errorf("SESSION_TTL_HOURS must be between 1 and 720")
		}
		cfg.SessionTTLHours = hours
	}

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

	if strings.TrimSpace(cfg.SessionSecret) == "" {
		return fmt.Errorf("SESSION_SECRET is required")
	}

	if len(cfg.SessionSecret) < 64 {
		return fmt.Errorf("SESSION_SECRET must contain at least 64 characters")
	}

	if cfg.DatabaseURL == "" {
		if cfg.DBPath == "" {
			return fmt.Errorf("DB_PATH is required when DATABASE_URL is not set")
		}

		if cfg.DBTestPath == "" {
			return fmt.Errorf("DB_TEST_PATH is required")
		}

		if cfg.DBPath == cfg.DBTestPath {
			return fmt.Errorf("DB_PATH and DB_TEST_PATH must be different")
		}
	}

	AppConfig = cfg
	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

