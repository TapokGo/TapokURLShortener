package service

import "github.com/Tapok-Go/TestURLShortener/internal/repo"

type URLService struct {
	repo repo.URLStorage
}

func NewURLService(s repo.URLStorage) *URLService {
	return &URLService{
		repo: s,
	}
}
