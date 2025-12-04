package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
)

type Storage struct {
	db *sql.DB
}

// Init the db instance
func New(cfg *config.Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db instance: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

// Save URL to db
func (s *Storage) SaveURL(temp string) error {
	return nil
}

// Close db instance
func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db instance")
	}
	
	return nil
}
