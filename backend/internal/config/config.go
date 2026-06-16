package config

import "os"

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	QRFrontendURL  string
}

func Load() *Config {
	return &Config{
		DatabaseURL:   envOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/restofnb?sslmode=disable"),
		JWTSecret:     envOrDefault("JWT_SECRET", "dev-secret-change-in-production"),
		QRFrontendURL: envOrDefault("QR_FRONTEND_URL", "http://localhost:3000"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
