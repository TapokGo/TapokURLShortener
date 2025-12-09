// Package dto provides a Data Transfer Object(DTO)
package dto

// CreateShortURL is a model of input JSON
type CreateShortURL struct {
	URL string `json:"url" validate:"required, url"`
}

// ShortURLResponse is a model of response to client
type ShortURLResponse struct {
	ShortURL string `json:"short_url"`
}
