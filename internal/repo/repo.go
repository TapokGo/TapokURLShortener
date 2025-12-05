// Package repo temp
package repo

// URLStorage is an interface that provides a SaveURL method
// allowing saves a URL to the db
type URLStorage interface {
	// Save function save pair URL (short, original).
	// Must return err if:
	// - short already exists
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
