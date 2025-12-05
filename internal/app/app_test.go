package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitApp(t *testing.T) {
	t.Setenv("URL_SHORTENER_ENV", "prod")

	// Create temp file for logs
	tmpFile, err := os.CreateTemp("", "test-app-*.log")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	t.Setenv("URL_SHORTENER_LOG_PATH", tmpFile.Name())

	cfg, err := config.LoadConfig("./testdata/valid.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)


	tmpDir := t.TempDir()
	cfg.StoragePath = filepath.Join(tmpDir, "test.db")

	app, err := New(*cfg)
	require.NoError(t, err)
	require.NotNil(t, app)

	assert.NotNil(t, app.Logger)
	assert.NotNil(t, app.logFile)
	assert.Equal(t, *cfg, app.Cfg)

	assert.Equal(t, cfg.LogPath, app.logFile.Name())
}

func TestInitApp_CloseFile(t *testing.T) {
	t.Setenv("URL_SHORTENER_ENV", "prod")

	// Create temp file for logs
	tmpFile, err := os.CreateTemp("", "test-app-*.log")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	t.Setenv("URL_SHORTENER_LOG_PATH", tmpFile.Name())

	cfg, err := config.LoadConfig("./testdata/valid.yaml")
	require.NoError(t, err)

	tmpDir := t.TempDir()
	cfg.StoragePath = filepath.Join(tmpDir, "test.db")

	app, err := New(*cfg)
	require.NoError(t, err)
	require.NotNil(t, app)

	err = app.Close()
	require.NoError(t, err)
	assert.Nil(t, app.logFile)
}

func TestIniApp_DevMode(t *testing.T) {
	cfg, err := config.LoadConfig("./testdata/valid.yaml")
	require.NoError(t, err)

	tmpDir := t.TempDir()
	cfg.StoragePath = filepath.Join(tmpDir, "test.db")

	app, err := New(*cfg)
	require.NoError(t, err)
	require.NotNil(t, app)

	assert.Nil(t, app.logFile)
	assert.NotNil(t, app.Logger)
}
