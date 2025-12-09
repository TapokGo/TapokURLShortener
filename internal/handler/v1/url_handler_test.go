package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TapokGo/TapokURLShortener/internal/handler/v1/dto"
	"github.com/TapokGo/TapokURLShortener/internal/logger"
	"github.com/TapokGo/TapokURLShortener/internal/service"
)

type mockService struct {
	CreatShortURlFunc   func(originalURL string) (string, error)
	ResolveShortURLFunc func(shortURL string) (string, error)
}

func (m *mockService) CreateShortURL(originalURL string) (string, error) {
	if m.CreatShortURlFunc != nil {
		return m.CreatShortURlFunc(originalURL)
	}
	return "", nil
}

func (m *mockService) ResolveShortURL(shortURL string) (string, error) {
	if m.ResolveShortURLFunc != nil {
		return m.ResolveShortURLFunc(shortURL)
	}
	return "", nil
}

type dummyLogger struct{}

func (d dummyLogger) Error(msg string, keysAndValues ...interface{}) {}
func (d dummyLogger) Warn(msg string, keysAndValues ...interface{}) {}
func (d dummyLogger) Info(msg string, keysAndValues ...interface{}) {}
func (d dummyLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (d dummyLogger) With(keysAndValues ...interface{}) (logger.Logger) {
	return dummyLogger{}
}
func (d dummyLogger) Close() error {
	return nil
}


func TestCreateShortURL_Success(t *testing.T) {
	mockSvc := &mockService{
		CreatShortURlFunc: func(originalURL string) (string, error) {
			return "abc123", nil
		},
	}

	handler := New(mockSvc, dummyLogger{}, "http://localhost:8080")
	r := chi.NewRouter()
	handler.Register(r)

	body := `{"url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp dto.ShortURLResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "http://localhost:8080/abc123", resp.ShortURL)
}

func TestCreateShortURL_InvalidJSON(t *testing.T) {
	mockSvc := &mockService{}

	handler := New(mockSvc, dummyLogger{}, "http://localhost:8080")
	r := chi.NewRouter()
	handler.Register(r)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateShortURL_EmptyURL(t *testing.T) {
	mockSvc := &mockService{
		CreatShortURlFunc: func(originalURL string) (string, error) {
			return "", service.ErrInvalidURL
		},
	}

	handler := New(mockSvc, dummyLogger{}, "http://localhost:8080")
	r := chi.NewRouter()
	handler.Register(r)

	body := `{"url": ""}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateShortURL_AliasGenFailed(t *testing.T) {
	mockSvc := &mockService{
		CreatShortURlFunc: func(originalURL string) (string, error) {
			return "", service.ErrAliasGenFailed
		},
	}

	handler := New(mockSvc, dummyLogger{}, "http://localhost:8080")
	r := chi.NewRouter()
	handler.Register(r)

	body := `{"url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRedirect_Success(t *testing.T) {
	mockSvc := &mockService{
		ResolveShortURLFunc: func(shortURL string) (string, error) {
			return "https://real-site.com", nil
		},
	}

	handler := New(mockSvc, dummyLogger{}, "http://localhost:8080")
	r := chi.NewRouter()
	handler.Register(r)

	req := httptest.NewRequest(http.MethodGet, "/xyz789", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://real-site.com", w.Header().Get("Location"))
}

func TestRedirect_NotFound(t *testing.T) {
	mockSvc := &mockService{
		ResolveShortURLFunc: func(shortURL string) (string, error) {
			return "", service.ErrNotFound
		},
	}

	handler := New(mockSvc, dummyLogger{}, "http://localhost:8080")
	r := chi.NewRouter()
	handler.Register(r)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}