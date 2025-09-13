package config_test

import (
	"golang/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestLoad_Defaults(t *testing.T) {
	t.Setenv("DB_PASSWORD", "secret") 

	cfg, err := config.Load()
	require.NoError(t, err)

	require.Equal(t, "secret", cfg.Database.Password)
	require.Equal(t, "localhost", cfg.Database.Host)
}

func TestLoad_MissingRequired(t *testing.T) {
	_, err := config.Load()
	require.Error(t, err, "missing DB_PASSWORD")
}