// Package repo provides an interface for work with the repository
package repo

import "errors"

var (
	ErrNotFound  = errors.New("URL not found")
	ErrDuplicate = errors.New("URL already exists")
)

// URLStorage defines the contract for persistent URL mapping storage
type URLStorage interface {
	// Save function save pair URL (short, original).
	// Must return err if:
	// -  short already exists (ErrDuplicate)
	// - error from db
	Save(short, original string) error

	// Get function return original URL by short URL
	// Must return:
	// - (url, nil) - if record already exists
	// - ("", ErrNotFound) - if not found
	// - ("", err) - idb error
	Get(short string) (string, error)

	// Close function close db instance.
	// If db cannot close it should return error
	Close() error
}
