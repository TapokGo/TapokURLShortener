// Package v1 provides the 1st version of the API
package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/TapokGo/TapokURLShortener/internal/dto"
	"github.com/TapokGo/TapokURLShortener/internal/handler/v1/httperror"
	"github.com/TapokGo/TapokURLShortener/internal/logger"
	"github.com/TapokGo/TapokURLShortener/internal/service"
)

// TODO: add logging

// URLHandler is a model of the handler layer
type URLHandler struct {
	urlService service.URLService
	baseURL    string
	logger     logger.Logger
}

// New returned a URLHandler pointer
func New(urlService service.URLService, logger logger.Logger, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		baseURL:    baseURL,
		logger:     logger,
	}
}

func (h *URLHandler) Register(r chi.Router) {
	r.Post("/shorten", h.CreateShortURL)
	r.Get("/{code}", h.Redirect)
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateShortURL
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, httperror.InvalidRequest("invalid request"))
		return
	}

	code, err := h.urlService.CreateShortURL(req.URL)
	if err != nil {
		if err == service.ErrInvalidURL {
			writeError(w, httperror.InvalidRequest("invalid URL"))
			return
		}

		if err == service.ErrAliasGenFailed {
			writeError(w, httperror.InternalServerError("attempt over"))
			return
		}

		writeError(w, httperror.InternalServerError("failed to create short URL"))
	}

	response := dto.ShortURLResponse{
		ShortURL: h.baseURL + "/" + code,
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeError(w, httperror.InvalidRequest("invalid URL parameter"))
		return
	}

	originalURL, err := h.urlService.ResolveShortURL(code)
	if err != nil {
		if err == service.ErrNotFound {
			writeError(w, httperror.NotFound("original URL not found"))
			return
		}

		writeError(w, httperror.InternalServerError("failed to get original URL"))
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err httperror.HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
