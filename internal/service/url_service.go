// Package service provides utilities for create short URL by origin and get original URL by short
package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"net/url"

	"github.com/Tapok-Go/TestURLShortener/internal/repo"
)

var (
	// ErrInvalidURL returned if original URL is invalid
	ErrInvalidURL = errors.New("invalid original URL")

	// ErrAliasGenFailed returned if can't create unique short URL
	ErrAliasGenFailed = errors.New("failed to create inique short URL")

	// ErrNotFound returned if short URl not exists
	ErrNotFound = errors.New("record not found")
)

// URLService is a model of url-shortener service layer
type URLService struct {
	repo repo.URLStorage
}

// NewURLService create the new URLService
func New(s repo.URLStorage) *URLService {
	return &URLService{
		repo: s,
	}
}

// CreateShortURL create short URL by origin URL
func (u *URLService) CreateShortURL(originURL string) (string, error) {
	// Validate URL
	_, err := url.Parse(originURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	// Try to create shortURL untill we get unique shortURL
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

	return "", ErrAliasGenFailed
}

// ResolveShortURL return shortURL original pair from db
func (u *URLService) ResolveShortURL(shortURL string) (string, error) {
	originalURL, err := u.repo.Get(shortURL)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return "", ErrNotFound
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
