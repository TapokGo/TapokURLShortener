// Package sqlite temp
package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Tapok-Go/TestURLShortener/internal/config"
	"github.com/Tapok-Go/TestURLShortener/internal/repo"
	_ "modernc.org/sqlite" // Register sqlite driver
)

// TODO: Rewrite with global errors
var ()

type storage struct {
	db       *sql.DB
	saveStmt *sql.Stmt
	getStmt  *sql.Stmt
}

// New init the db instance
func New(cfg *config.Config) (repo.URLStorage, error) {
	// Check a dir is exists
	dbDir := filepath.Dir(cfg.StoragePath)
	if err := os.MkdirAll(dbDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create db dir %q: %w", dbDir, err)
	}

	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	saveStmt, getStmt, err := createStmts(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create statements: %w", err)
	}

	s := &storage{
		db:       db,
		saveStmt: saveStmt,
		getStmt:  getStmt,
	}

	if err = s.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return s, nil
}

// Save function save URL to db
func (s *storage) Save(short, origin string) error {
	// TODO: add check existence of short URL
	_, err := s.saveStmt.Exec(short, origin)
	if err != nil {
		return fmt.Errorf("failed to save pair URL to db: %w", err)
	}

	return nil
}

// Get function get original URL by short
func (s *storage) Get(short string) (string, error) {
	var originURL string
	err := s.getStmt.QueryRow(short).Scan(&originURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("result not found: %w", err)
		}
		return "", fmt.Errorf("failed to save pair URL to db: %w", err)
	}

	return originURL, nil
}

// Close db instance
func (s *storage) Close() error {
	err := s.saveStmt.Close()
	if err != nil {
		return fmt.Errorf("failed to close save statement: %w", err)
	}

	err = s.getStmt.Close()
	if err != nil {
		return fmt.Errorf("failed to close get statement: %w", err)
	}

	err = s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db instance: %w", err)
	}

	return nil
}

func (s *storage) initSchema() error {
	q := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_url TEXT NOT NULL,
		origin_url TEXT NOT NULL
	)
	`

	if _, err := s.db.Exec(q); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func createStmts(db *sql.DB) (save *sql.Stmt, get *sql.Stmt, err error) {
	saveQuery := `
		INSERT INTO urls (
			short_url, 
			origin_url
		) VALUES (
			?, ?
		)
	`
	getQuery := `
		SELECT 
			origin_url FROM urls
		WHERE short_url = ? 
	`

	saveStmt, err := db.Prepare(saveQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create save statements: %w", err)
	}

	getStmt, err := db.Prepare(getQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create get statements: %w", err)
	}

	return saveStmt, getStmt, nil
}
