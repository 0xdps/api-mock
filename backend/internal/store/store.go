package store

import (
	"fmt"
	"log"

	"github.com/mesahub-db/mesahub-pkg-go/mesahub"
)

// Store wraps a MesaHub database for all persistent user data.
type Store struct {
	db *mesahub.DatabaseHandle
}

// New creates a Store and bootstraps all required tables.
// connectionURL must be a MesaHub connection string:
//
//	mh://shs_apikey@host/dbname
func New(connectionURL string) (*Store, error) {
	info, err := mesahub.ParseURL(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("invalid MESAHUB_URL: %w", err)
	}

	client := mesahub.New(mesahub.Config{
		APIKey:      info.APIKey,
		APIURL:      info.APIURL,
		RoutePrefix: info.RoutePrefix,
	})
	db := client.DB(info.DBName)

	s := &Store{db: db}
	if err := s.bootstrap(); err != nil {
		return nil, fmt.Errorf("store bootstrap: %w", err)
	}
	return s, nil
}

// bootstrap creates all tables if they do not already exist.
func (s *Store) bootstrap() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         TEXT PRIMARY KEY,
			email      TEXT NOT NULL,
			name       TEXT NOT NULL DEFAULT '',
			plan       TEXT NOT NULL DEFAULT 'free',
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS templates (
			id          TEXT PRIMARY KEY,
			user_id     TEXT NOT NULL,
			name        TEXT NOT NULL,
			slug        TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			visibility  TEXT NOT NULL DEFAULT 'private',
			type        TEXT NOT NULL DEFAULT 'custom',
			base_schema TEXT,
			schema_json TEXT NOT NULL,
			is_featured INTEGER NOT NULL DEFAULT 0,
			created_at  TEXT NOT NULL,
			updated_at  TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS api_keys (
			id           TEXT PRIMARY KEY,
			user_id      TEXT NOT NULL,
			name         TEXT NOT NULL,
			key_hash     TEXT NOT NULL UNIQUE,
			key_prefix   TEXT NOT NULL,
			type         TEXT NOT NULL DEFAULT 'personal',
			status       TEXT NOT NULL DEFAULT 'active',
			created_at   TEXT NOT NULL,
			last_used_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS usage_logs (
			id          TEXT PRIMARY KEY,
			api_key_id  TEXT NOT NULL,
			template_id TEXT NOT NULL,
			status_code INTEGER NOT NULL,
			created_at  TEXT NOT NULL
		)`,
	}

	for _, stmt := range stmts {
		if _, err := s.db.Exec(stmt, nil); err != nil {
			return fmt.Errorf("exec bootstrap DDL: %w", err)
		}
	}

	// Run schema migrations for existing databases (errors are intentionally ignored
	// since the column may already exist from a previous deployment).
	migrations := []string{
		`ALTER TABLE templates ADD COLUMN is_featured INTEGER NOT NULL DEFAULT 0`,
	}
	for _, m := range migrations {
		s.db.Exec(m, nil) //nolint:errcheck
	}

	log.Println("✅ MesaHub tables bootstrapped")
	return nil
}

// DB returns the underlying MesaHub DatabaseHandle (used by sub-stores).
func (s *Store) DB() *mesahub.DatabaseHandle {
	return s.db
}
