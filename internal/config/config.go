// Package config provides application configuration.
package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Sentry   SentryConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port int `envconfig:"PORT" default:"3000"`
}

type SentryConfig struct {
	DSN string `envconfig:"SENTRY_DSN"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" default:"app"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("APP", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
