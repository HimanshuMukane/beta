package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string
	AppEnv     string
	ServerPort string

	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
	DatabaseURL string

	IdentityDBName string
	IdentityDBURL  string

	AllowTestDBHeader bool
	TestDBHeaderKey   string

	SessionCookieName   string
	SessionCookieDomain string
	SessionTTLHours     int
	SessionSecure       bool
	SessionSecret       string
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
		// Automatically check if caller belongs to a specific project directory with a local .env
		if projectDir := findCallerProjectDir(); projectDir != "" {
			localEnv := filepath.Join(projectDir, ".env")
			if info, err := os.Stat(localEnv); err == nil {
				// Avoid reloading the same root .env if running from that directory
				if rootInfo, err := os.Stat(".env"); err != nil || !os.SameFile(info, rootInfo) {
					_ = godotenv.Overload(localEnv)
				}
			}
		} else {
			_ = godotenv.Overload(".env")
		}
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
	if dbPassword == "" {
		dbPassword = "especa_root_password"
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" && dbName != "" {
		databaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName)
	} else if databaseURL == "" {
		databaseURL = fmt.Sprintf("root:%s@tcp(%s:%s)/platform_identity?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)
	}

	identityDBName := strings.TrimSpace(getEnv("IDENTITY_DB_NAME", ""))
	identityDBURL := strings.TrimSpace(os.Getenv("IDENTITY_DATABASE_URL"))
	if identityDBURL == "" && identityDBName != "" {
		identityDBURL = fmt.Sprintf("root:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			dbPassword, dbHost, dbPort, identityDBName)
	}

	appName := strings.TrimSpace(os.Getenv("APP_NAME"))
	if appName == "" {
		if projectDir := findCallerProjectDir(); projectDir != "" {
			appName = filepath.Base(projectDir)
		} else {
			appName = "default"
		}
	}

	cfg := Config{
		AppName:             appName,
		AppEnv:              strings.ToLower(strings.TrimSpace(getEnv("APP_ENV", "development"))),
		ServerPort:          strings.TrimSpace(getEnv("SERVER_PORT", "8080")),
		DBHost:              dbHost,
		DBPort:              dbPort,
		DBName:              dbName,
		DBUser:              dbUser,
		DBPassword:          dbPassword,
		DatabaseURL:         databaseURL,
		IdentityDBName:      identityDBName,
		IdentityDBURL:       identityDBURL,
		TestDBHeaderKey:     strings.TrimSpace(getEnv("TEST_DB_HEADER_KEY", "X-Database-Env")),
		SessionCookieName:   strings.TrimSpace(getEnv("SESSION_COOKIE_NAME", "app_session")),
		SessionCookieDomain: strings.TrimSpace(os.Getenv("SESSION_COOKIE_DOMAIN")),
		SessionSecret:       os.Getenv("SESSION_SECRET"),
		SessionTTLHours:     24,
		SessionSecure:       false,
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

func findCallerProjectDir() string {
	goroot := runtime.GOROOT()
	for i := 1; i < 20; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if goroot != "" && strings.HasPrefix(file, goroot) {
			continue
		}
		dir := filepath.Dir(file)
		if strings.Contains(dir, filepath.Join("projects", "core")) {
			continue
		}
		curr := dir
		for {
			if _, err := os.Stat(filepath.Join(curr, "go.mod")); err == nil {
				return curr
			}
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}
	return ""
}
