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
	require.NotNil(t, cfg)

	assert.Equal(t, "dev", cfg.Env)
	assert.Equal(t, "./storage/storage.db", cfg.StoragePath)
	assert.Equal(t, "localhost", cfg.HTTPServer.Address)
	assert.Equal(t, 8082, cfg.HTTPServer.Port)
	assert.Equal(t, 4*time.Second, cfg.HTTPServer.Timeout)
	assert.Equal(t, 60*time.Second, cfg.HTTPServer.IdleTimeout)
}

func TestLoadConfig_OverrideENV(t *testing.T) {
	t.Setenv("URL_SHORTENER_ADDRESS", "test_url_address")
	t.Setenv("URL_SHORTENER_PORT", "3000")

	cfg, err := LoadConfig("testdata/valid.yaml")
	require.NotNil(t, cfg)
	require.NoError(t, err)

	assert.Equal(t, "test_url_address", cfg.HTTPServer.Address)
	assert.Equal(t, 3000, cfg.HTTPServer.Port)
}

func TestLoadConfig_InvalidData(t *testing.T) {
	t.Run("Invalid port", func(t *testing.T) {
		t.Setenv("URL_SHORTENER_PORT", "0")

		cfg, err := LoadConfig("testdata/valid.yaml")
		require.Error(t, err)
		require.Nil(t, cfg)
		assert.Contains(t, err.Error(), "invalid port")
	})

	t.Run("Invalid ENV", func(t *testing.T) {
		t.Setenv("URL_SHORTENER_ENV", "stage")

		cfg, err := LoadConfig("testdata/valid.yaml")
		require.Error(t, err)
		require.Nil(t, cfg)
		assert.Contains(t, err.Error(), "invalid env")
	})
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	t.Parallel()
	cfg, err := LoadConfig("testdata/non_valid.yaml")
	require.Error(t, err)
	require.Nil(t, cfg)

	assert.Contains(t, err.Error(), "failed to read config")
	assert.Nil(t, cfg)
}
