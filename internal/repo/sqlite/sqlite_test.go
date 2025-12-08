package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/TapokGo/TapokURLShortener/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlite_New(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "test.db")
	storage, err := New(path)
	defer func() {
		err = storage.Close()
		require.NoError(t, err)
	}()

	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestSqlite_SaveAndGet(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "test.db")
	storage, _ := New(path)
	defer func() {
		_ = storage.Close()
	}()

	err := storage.Save("short_url", "origin_url")
	require.NoError(t, err)

	original, err := storage.Get("short_url")
	require.NoError(t, err)
	assert.Equal(t, "origin_url", original)
}

func TestSqlite_DuplicateSave(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "test.db")
	storage, _ := New(path)
	defer func() {
		_ = storage.Close()
	}()

	err := storage.Save("short_url", "origin_url")
	require.NoError(t, err)

	err = storage.Save("short_url", "origin_url")
	require.Error(t, err)
	assert.ErrorIs(t, err, repo.ErrDuplicate)
}

func TestSqlite_GetNoResult(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "test.db")
	storage, _ := New(path)
	defer func() {
		_ = storage.Close()
	}()

	original, err := storage.Get("short_url")
	require.Error(t, err)
	assert.Empty(t, original)
	assert.ErrorIs(t, err, repo.ErrNotFound)
}
