// Package url provides utilities for create short URL by origin and get original URL by short
package url

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"net/url"

	"github.com/TapokGo/TapokURLShortener/internal/repo"
	"github.com/TapokGo/TapokURLShortener/internal/service"
)

// URLService is a model of url-shortener service layer
type urlService struct {
	repo repo.URLStorage
}

// New create the new URLService
func New(s repo.URLStorage) service.URLService {
	return &urlService{
		repo: s,
	}
}

// CreateShortURL create short URL by origin URL
func (u *urlService) CreateShortURL(originURL string) (string, error) {
	// Validate URL
	url, err := url.ParseRequestURI(originURL)
	if err != nil {
		return "", service.ErrInvalidURL
	}

	if url.Host == "" {
		return "", service.ErrInvalidURL
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return "", service.ErrInvalidURL
	}

	// Try to create shortURL until we get unique shortURL
	const maxAttempts = 5
	for range maxAttempts {
		shortURL, err := generateAlias(8)
		if err != nil {
			return "", err
		}

		err = u.repo.Save(shortURL, originURL)
		if err == nil {
			return shortURL, nil
		}

		if !errors.Is(err, repo.ErrDuplicate) {
			return "", fmt.Errorf("failed to create short URL: %w", err)
		}
	}

	return "", service.ErrAliasGenFailed
}

// ResolveShortURL return shortURL original pair from db
func (u *urlService) ResolveShortURL(shortURL string) (string, error) {
	originalURL, err := u.repo.Get(shortURL)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return "", service.ErrNotFound
		}
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}

	return originalURL, nil
}

func generateAlias(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(bytes)
	return s[:n], nil
}
