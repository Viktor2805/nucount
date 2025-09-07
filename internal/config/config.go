// Package config handles the application's configuration, including environment
// variable loading and database connection setup.
package config

import "os"

type Config struct {
	Port      string
	SentryDSN string
}


type SentryConfig struct {
	DSN string `envconfig:"SENTRY_DSN"`
}

func LoadConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		SentryDSN: os.Getenv("SENTRY_DSN"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

