package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TapokGo/TapokURLShortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitApp(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-app-*.log")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tmpFile.Close())
		err = os.Remove(tmpFile.Name())
		require.NoError(t, err)
	}()
	require.Nil(t, err)

	t.Setenv("URL_SHORTENER_ENV", "prod")
	t.Setenv("URL_SHORTENER_LOG_PATH", tmpFile.Name())

	cfg, err := config.LoadConfig("./testdata/valid.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	tmpDir := t.TempDir()
	cfg.StoragePath = filepath.Join(tmpDir, "test.db")

	app, err := New(*cfg)
	defer func() {
		require.NoError(t, app.Close())
		require.Nil(t, app.repoCloser)
		require.Nil(t, app.logCloser)
	}()

	require.NoError(t, err)
	require.NotNil(t, app)

	assert.NotNil(t, app.Logger)
	assert.NotNil(t, app.logCloser)
	assert.NotNil(t, app.repoCloser)

	assert.Equal(t, *cfg, app.cfg)
}
