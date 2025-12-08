// Package repo provides an interface for work with the repository
package repo

import "errors"

var (
	// ErrNotFound returned when if short URL not exist in db
	ErrNotFound = errors.New("URL not found")

	// ErrDuplicate returned if short URl already exists in db
	ErrDuplicate = errors.New("URL already exists")
)

// URLStorage defines the contract for URL storage
type URLStorage interface {
	// Save save pair URL (short, original).
	// Must return err if:
	// -  short already exists (ErrDuplicate)
	// - error from db
	Save(short, original string) error

	// Get return original URL by short URL
	// Must return:
	// - (url, nil) - if record exists
	// - ("", ErrNotFound) - if not found
	// - ("", err) - idb error
	Get(short string) (string, error)

	// Close close db instance.
	// If db cannot close it should return error
	Close() error
}
