package config

import (
	"fmt"
	"os"
)

const defaultMaxSteps = 10

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	AppPort           string
	AppEnv            string
	FrontendDevOrigin string
	FrontendDist      string
	MigrationsPath    string
	MaxSteps          int
}

func Load() Config {
	return Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "profil_math"),
		AppPort:           getEnv("APP_PORT", "8080"),
		AppEnv:            getEnv("APP_ENV", "development"),
		FrontendDevOrigin: getEnv("FRONTEND_DEV_ORIGIN", "http://localhost:5173"),
		FrontendDist:      getEnv("FRONTEND_DIST", "frontend/dist"),
		MigrationsPath:    getEnv("MIGRATIONS_PATH", "migrations"),
		MaxSteps:          defaultMaxSteps,
	}
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
