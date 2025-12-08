package slog

import (
	"os"
	"strings"
	"testing"

	"github.com/TapokGo/TapokURLShortener/internal/config"
	logInterface "github.com/TapokGo/TapokURLShortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlog(t *testing.T) {
	logger, logReader, cleanup := newProdLogger(t)
	defer cleanup()

	require.NotNil(t, logger)

	logger.Info("test log", "key", "value")
	logFileContent := logReader()

	assert.Contains(t, logFileContent, `"msg":"test log"`)
	assert.Contains(t, logFileContent, `"key":"value"`)
	assert.Contains(t, logFileContent, `"level":"INFO"`)
}

func TestSlog_With(t *testing.T) {
	logger, logReader, cleanup := newProdLogger(t)
	defer cleanup()

	newLogger := logger.With("key", "value")
	require.NotNil(t, newLogger)
	require.NotEqual(t, logger, newLogger)

	logger.Info("original")
	newLogger.Info("after with")

	logFileContent := logReader()
	logs := strings.Split(strings.TrimSpace(logFileContent), "\n")

	assert.NotContains(t, logs[0], `"key":"value"`)
	assert.Contains(t, logs[1], `"key":"value"`)
}

func newProdLogger(t *testing.T) (logger logInterface.Logger, readFileFunc func() string, cleanUp func()) {
	t.Setenv("URL_SHORTENER_ENV", "prod")

	tmpFile, err := os.CreateTemp("", "test-app-*.log")
	require.NoError(t, err)

	t.Setenv("URL_SHORTENER_LOG_PATH", tmpFile.Name())

	cfg, err := config.LoadConfig("")
	require.NoError(t, err)

	logger, err = New(cfg)
	require.NoError(t, err)
	assert.NotNil(t, logger)

	return logger,
		func() string {
			logContent, err := os.ReadFile(tmpFile.Name())
			require.NoError(t, err)
			return string(logContent)
		},
		func() {
			if logger != nil {
				_ = logger.Close()
				_ = os.Remove(tmpFile.Name())
			}
		}
}
