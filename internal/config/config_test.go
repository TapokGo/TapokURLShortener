package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_YAML(t *testing.T) {
	t.Parallel()
	cfg, err := LoadConfig("testdata/valid.yaml")
	require.NoError(t, err)

	assert.Equal(t, "dev", cfg.Env)
	assert.Equal(t, "./storage/storage.db", cfg.StoragePath)
	assert.Equal(t, "localhost", cfg.HTTPServer.Address)
	assert.Equal(t, 8082, cfg.HTTPServer.Port)
	assert.Equal(t, 4*time.Second, cfg.HTTPServer.Timeout)
	assert.Equal(t, 60*time.Second, cfg.HTTPServer.IdleTimeout)
}

func TestLoadConfig_OverrideENV(t *testing.T) {
	t.Setenv("URL_SHORTENER_ADDRESS", "babai")
	t.Setenv("URL_SHORTENER_PORT", "3000")

	cfg, err := LoadConfig("testdata/valid.yaml")
	require.NoError(t, err)

	assert.Equal(t, "babai", cfg.HTTPServer.Address)
	assert.Equal(t, 3000, cfg.HTTPServer.Port)
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	t.Setenv("URL_SHORTENER_PORT", "0")

	cfg, err := LoadConfig("testdata/valid.yaml")
	require.Error(t, err)

	assert.Contains(t, err.Error(), "invalid port")
	assert.Nil(t, cfg)
}

func TestLoadConfig_InvalidENV(t *testing.T) {
	t.Setenv("URL_SHORTENER_ENV", "stage")

	cfg, err := LoadConfig("testdata/valid.yaml")
	require.Error(t, err)

	assert.Contains(t, err.Error(), "invalid env")
	assert.Nil(t, cfg)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	t.Parallel()
	cfg, err := LoadConfig("testdata/non_valid.yaml")
	require.Error(t, err)

	assert.Contains(t, err.Error(), "failed to read config")
	assert.Nil(t, cfg)
}
