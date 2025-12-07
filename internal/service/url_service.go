// Package service provides utilities for create short URL by origin and get original URL by short
package service

import "github.com/Tapok-Go/TestURLShortener/internal/repo"

// URLService is a model of url-shortener service layer
type URLService struct {
	repo repo.URLStorage
}

// NewURLService create the new URLService
func NewURLService(s repo.URLStorage) *URLService {
	return &URLService{
		repo: s,
	}
}

// CreateShortURL create short URL by origin URL
func (u *URLService) CreateShortURL(originURL string) (string, error) {
	// Validate URL

	// Create alias

	// Save pair to the db (and check what that short URL not already exist)

	// return short url
	return "", nil
}

// ResolveShortURL return shortURL original pair from db
func (u *URLService) ResolveShortURL(shortURL string) (string, error) {
	// Get short URL from db
	return "", nil
}

func (u *URLService) generateAlias(originURL string) (string, error) {
	// Create unique short URL base on original
	return "", nil
}
