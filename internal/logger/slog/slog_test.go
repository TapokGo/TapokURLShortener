package slog

import (
	"os"
	"strings"
	"testing"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	logInterface "github.com/Tapok-Go/TestURLShortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlog(t *testing.T) {
	logger, logContent, cleanup := newProdLogger(t)
	defer cleanup()

	assert.NotNil(t, logger)
	logger.Info("test prod message", "key", "value")

	assert.Contains(t, logContent(), `"msg":"test prod message"`)
	assert.Contains(t, logContent(), `"key":"value"`)
}

func TestSlog_With(t *testing.T) {
	logger, logContent, cleanup := newProdLogger(t)
	defer cleanup()

	newLogger := logger.With("user_id", "1234")
	require.NotNil(t, newLogger)
	require.NotEqual(t, logger, newLogger)

	newLogger.Info("babai")
	assert.Contains(t, logContent(), `"user_id":"1234"`)

	logger.Info("original")
	newLogger.Info("with context")

	lines := strings.Split(strings.TrimSpace(logContent()), "\n")

	assert.NotContains(t, lines[1], `"user_id"`, "must not contain a new fields")
	assert.Contains(t, lines[2], `"user_id":"1234"`)
}

func newProdLogger(t *testing.T) (logger logInterface.Logger, readFileFunc func() string, cleanUp func()) {
	t.Setenv("URL_SHORTENER_ENV", "prod")

	tmpFile, err := os.CreateTemp("", "test-app-*.log")
	require.NoError(t, err)

	t.Setenv("URL_SHORTENER_LOG_PATH", tmpFile.Name())

	cfg, err := config.LoadConfig("")
	require.NoError(t, err)

	logger, logFile, err := NewSlogLogger(cfg)
	assert.NotNil(t, logFile)

	return logger,
		func() string {
			logContent, err := os.ReadFile(logFile.Name())
			require.NoError(t, err)
			return string(logContent)
		},
		func() {
			if logFile != nil {
				_ = logFile.Close()
				_ = os.Remove(tmpFile.Name())
			}
		}
}
