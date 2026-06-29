package config

import (
	"fmt"
	"os"
)

// Config holds application configuration.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load reads configuration from environment variables.
// DATABASE_URL can be provided as a full URL or built from DB_HOST, DB_USER,
// DB_PASSWORD, DB_NAME, DB_PORT and DB_SSLMODE variables.
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := buildDatabaseURL()
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required (or DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE)")
	}

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
	}, nil
}

func buildDatabaseURL() string {
	if fullURL := os.Getenv("DATABASE_URL"); fullURL != "" {
		return fullURL
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" || user == "" || name == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslMode)
}
