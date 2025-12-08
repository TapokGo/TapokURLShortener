package url_service

import (
	"testing"

	"github.com/Tapok-Go/TestURLShortener/internal/repo"
	"github.com/Tapok-Go/TestURLShortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	SaveFunc  func(shortURL, originalURL string) error
	GetFunc   func(shortURL string) (string, error)
	CloseFunc func() error
}

func (m *mockStorage) Save(shortURL, originalURL string) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(shortURL, originalURL)
	}
	return nil
}

func (m *mockStorage) Get(shortURL string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(shortURL)
	}
	return "", nil
}

func (m *mockStorage) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestService_CreateShortURL(t *testing.T) {
	mockRepo := &mockStorage{
		SaveFunc: func(shortURL, originalURL string) error {
			return nil
		},
	}

	svc := New(mockRepo)
	shortURL, err := svc.CreateShortURL("https://example.com")
	require.NoError(t, err)
	require.NotEmpty(t, shortURL)

	assert.Len(t, shortURL, 8)
}

func TestService_CreateShortURL_Validation(t *testing.T) {
	mockRepo := &mockStorage{}
	svc := New(mockRepo)

	shortURL, err := svc.CreateShortURL("/test")
	require.Empty(t, shortURL)
	require.Error(t, err)

	assert.ErrorIs(t, err, service.ErrInvalidURL)

}

func TestService_CreateShortURL_WithDuplicate(t *testing.T) {
	attempt := 0
	mockRepo := &mockStorage{
		SaveFunc: func(shortURL, originalURL string) error {
			attempt++
			if attempt == 1 {
				return repo.ErrDuplicate
			}
			return nil
		},
	}

	svc := New(mockRepo)
	shortURL, err := svc.CreateShortURL("https://example.com")
	require.NoError(t, err)
	require.NotEmpty(t, shortURL)

	assert.Len(t, shortURL, 8)
	assert.Equal(t, 2, attempt)
}

func TestService_CreateShortURL_MaxAttemptExceeded(t *testing.T) {
	attempt := 0
	mockRepo := &mockStorage{
		SaveFunc: func(shortURL, originalURL string) error {
			attempt++
			return repo.ErrDuplicate
		},
	}

	svc := New(mockRepo)
	shortURL, err := svc.CreateShortURL("https://example.com")
	require.Error(t, err)
	require.Empty(t, shortURL)

	assert.Equal(t, 5, attempt)
	assert.ErrorIs(t, err, service.ErrAliasGenFailed)
}

func TestService_ResolveShortURL(t *testing.T) {
	mockRepo := &mockStorage{
		GetFunc: func(shortURL string) (string, error) {
			return "original test", nil
		},
	}

	svc := New(mockRepo)
	originURL, err := svc.ResolveShortURL("short test")
	require.NotEmpty(t, originURL)
	require.NoError(t, err)

	assert.Equal(t, "original test", originURL)
}

func TestService_ResolveShortURL_NotFound(t *testing.T) {
	mockRepo := &mockStorage{
		GetFunc: func(shortURL string) (string, error) {
			return "", repo.ErrNotFound
		},
	}

	svc := New(mockRepo)
	originURL, err := svc.ResolveShortURL("short test")
	require.Error(t, err)
	require.Empty(t, originURL)

	assert.ErrorIs(t, err, service.ErrNotFound)
}
