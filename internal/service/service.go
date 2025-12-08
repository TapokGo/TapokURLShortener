// Package service provides an interface for work with the service
package service

import "errors"

var (
	// ErrInvalidURL returned if original URL is invalid
	ErrInvalidURL = errors.New("invalid original URL")

	// ErrAliasGenFailed returned if can't create unique short URL
	ErrAliasGenFailed = errors.New("failed to create unique short URL")

	// ErrNotFound returned if short URl not exists
	ErrNotFound = errors.New("record not found")
)

// URLService defines the contract for service
type URLService interface {
	/* Create short URL by original.
	Must return:
	- ("", ErrInvalidURL) - if URL is invalid
	- ("", err) - failed to create alias, db error
	- ("", ErrAliasGenFailed) - failed to create unique short URL (attempts are over)
	- (shortURL, nil) - successful create shortURL
	*/
	CreateShortURL(originURL string) (string, error)

	/* Get original URL by short
	Must return:
	- ("", ErrNotFound) - record not exist in db
	- ("", err) - unexpected error
	- (originalURL, nil) - successful; get originalURL
	*/
	ResolveShortURL(shortURL string) (string, error)
}
