package utils

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

var AppConfig Config

func init() {
	LoadConfig()
}

func LoadConfig() {
	_ = godotenv.Load()

	if os.Getenv("DB_HOST") == "" {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			projectRoot := filepath.Join(filepath.Dir(filename), "..")
			_ = godotenv.Load(filepath.Join(projectRoot, ".env"))
		}
	}

	AppConfig = Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}
}
