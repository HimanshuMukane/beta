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

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DatabaseURL string

	AllowTestDBHeader bool
	TestDBHeaderKey   string

	SessionCookieName string
	SessionTTLHours   int
	SessionSecure     bool
	SessionSecret     string
	LimenSecret       string
}

var AppConfig Config

func LoadConfig(filenames ...string) error {
	for _, rootPath := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(rootPath); err == nil {
			_ = godotenv.Load(rootPath)
			break
		}
	}

	if len(filenames) > 0 {
		for _, f := range filenames {
			_ = godotenv.Overload(f)
		}
	} else {
		_ = godotenv.Overload(".env")
	}

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))
	dbPassword := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	dbUser := strings.TrimSpace(os.Getenv("DB_USER"))

	if dbUser == "" {
		if dbName != "" {
			dbUser = dbName
		} else {
			dbUser = "especa"
		}
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" && dbName != "" {
		databaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName)
	} else if databaseURL == "" {
		databaseURL = fmt.Sprintf("especa:especa@tcp(%s:%s)/especa?parseTime=true&loc=Local", dbHost, dbPort)
	}

	cfg := Config{
		AppEnv:            strings.ToLower(strings.TrimSpace(getEnv("APP_ENV", "development"))),
		ServerPort:        strings.TrimSpace(getEnv("SERVER_PORT", "8080")),
		DBHost:            dbHost,
		DBPort:            dbPort,
		DBName:            dbName,
		DBUser:            dbUser,
		DBPassword:        dbPassword,
		DatabaseURL:       databaseURL,
		TestDBHeaderKey:   strings.TrimSpace(getEnv("TEST_DB_HEADER_KEY", "X-Database-Env")),
		SessionCookieName: strings.TrimSpace(getEnv("SESSION_COOKIE_NAME", "app_session")),
		SessionSecret:     os.Getenv("SESSION_SECRET"),
		LimenSecret:       getEnv("LIMEN_SECRET", os.Getenv("SESSION_SECRET")),
		SessionTTLHours:   24,
		SessionSecure:     false,
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
		cfg.AllowTestDBHeader = true
	}

	if value := os.Getenv("SESSION_SECURE"); value != "" {
		secure, err := strconv.ParseBool(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("SESSION_SECURE must be true or false")
		}
		cfg.SessionSecure = secure
	} else if cfg.AppEnv == "production" {
		cfg.SessionSecure = true
	}

	if value := os.Getenv("SESSION_TTL_HOURS"); value != "" {
		hours, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || hours <= 0 {
			return fmt.Errorf("SESSION_TTL_HOURS must be a positive integer")
		}
		cfg.SessionTTLHours = hours
	}

	if cfg.SessionSecret == "" {
		if cfg.AppEnv == "production" {
			return fmt.Errorf("SESSION_SECRET is required in production")
		}
		cfg.SessionSecret = "dev-insecure-secret-key-must-be-at-least-64-bytes-long-for-hmac-sha256-safety"
	}
	if len(cfg.SessionSecret) < 32 {
		return fmt.Errorf("SESSION_SECRET must be at least 32 bytes")
	}

	if cfg.LimenSecret == "" {
		cfg.LimenSecret = cfg.SessionSecret
	}

	AppConfig = cfg
	return nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists || strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
