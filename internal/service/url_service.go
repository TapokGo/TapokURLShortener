// Package service temp
package service

import "github.com/Tapok-Go/TestURLShortener/internal/repo"

// URLService temp
type URLService struct {
	repo repo.URLStorage
}

// NewURLService temp
func NewURLService(s repo.URLStorage) *URLService {
	return &URLService{
		repo: s,
	}
}
