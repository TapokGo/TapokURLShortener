package repo

// URLStorage is an interface that provides a SaveURL method
// allowwig saves a URL to the db
type URLStorage interface {
	SaveURL(temp string) error
	Close() error
}
