package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const DatabaseName = "state.db"

type Store struct{ db *sql.DB }

type Instance struct {
	Name       string    `json:"name"`
	Provider   string    `json:"provider"`
	Version    string    `json:"version"`
	Status     string    `json:"status"`
	SpecDigest string    `json:"specDigest"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Operation struct {
	ID           string     `json:"id"`
	InstanceName string     `json:"instanceName"`
	Kind         string     `json:"kind"`
	Status       string     `json:"status"`
	SpecDigest   string     `json:"specDigest"`
	StartedAt    time.Time  `json:"startedAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	Error        string     `json:"error,omitempty"`
}

func Exists(directory string) bool {
	_, err := os.Stat(filepath.Join(directory, DatabaseName))
	return err == nil
}

func Open(directory string) (*Store, error) {
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return nil, fmt.Errorf("create state directory: %w", err)
	}
	db, err := sql.Open("sqlite", filepath.Join(directory, DatabaseName))
	if err != nil {
		return nil, fmt.Errorf("open state database: %w", err)
	}
	store := &Store{db: db}
	if err := store.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func OpenReadOnly(directory string) (*Store, error) {
	path := filepath.Join(directory, DatabaseName)
	if !Exists(directory) {
		return nil, os.ErrNotExist
	}
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("open read-only state database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("verify read-only state database: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	const schema = `
PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY);
INSERT OR IGNORE INTO schema_version(version) VALUES (1);
CREATE TABLE IF NOT EXISTS instances (
  name TEXT PRIMARY KEY,
  provider TEXT NOT NULL,
  version TEXT NOT NULL,
  status TEXT NOT NULL,
  spec_digest TEXT NOT NULL,
  spec_json BLOB NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS operations (
  id TEXT PRIMARY KEY,
  instance_name TEXT NOT NULL,
  kind TEXT NOT NULL,
  status TEXT NOT NULL,
  spec_digest TEXT NOT NULL,
  started_at TEXT NOT NULL,
  completed_at TEXT,
  error TEXT NOT NULL DEFAULT ''
);`
	if _, err := s.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("migrate state database: %w", err)
	}
	return nil
}

func (s *Store) UpsertInstance(ctx context.Context, instance Instance, normalizedSpec any) error {
	data, err := json.Marshal(normalizedSpec)
	if err != nil {
		return fmt.Errorf("serialize instance specification: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO instances(name, provider, version, status, spec_digest, spec_json, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(name) DO UPDATE SET provider=excluded.provider, version=excluded.version,
status=excluded.status, spec_digest=excluded.spec_digest, spec_json=excluded.spec_json,
updated_at=excluded.updated_at`, instance.Name, instance.Provider, instance.Version,
		instance.Status, instance.SpecDigest, data, instance.UpdatedAt.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("upsert instance: %w", err)
	}
	return nil
}

func (s *Store) Instances(ctx context.Context, name string) ([]Instance, error) {
	query := `SELECT name, provider, version, status, spec_digest, updated_at FROM instances`
	args := []any{}
	if name != "" {
		query += ` WHERE name = ?`
		args = append(args, name)
	}
	query += ` ORDER BY name`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query instances: %w", err)
	}
	defer rows.Close()
	var result []Instance
	for rows.Next() {
		var instance Instance
		var updated string
		if err := rows.Scan(&instance.Name, &instance.Provider, &instance.Version, &instance.Status, &instance.SpecDigest, &updated); err != nil {
			return nil, fmt.Errorf("scan instance: %w", err)
		}
		instance.UpdatedAt, err = time.Parse(time.RFC3339Nano, updated)
		if err != nil {
			return nil, fmt.Errorf("parse instance update time: %w", err)
		}
		result = append(result, instance)
	}
	return result, rows.Err()
}

func (s *Store) StartOperation(ctx context.Context, operation Operation) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO operations(id, instance_name, kind, status, spec_digest, started_at)
VALUES (?, ?, ?, ?, ?, ?)`, operation.ID, operation.InstanceName, operation.Kind,
		operation.Status, operation.SpecDigest, operation.StartedAt.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("start operation: %w", err)
	}
	return nil
}

func (s *Store) FinishOperation(ctx context.Context, id, status, message string, completed time.Time) error {
	result, err := s.db.ExecContext(ctx, `UPDATE operations SET status=?, error=?, completed_at=? WHERE id=?`,
		status, message, completed.UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return fmt.Errorf("finish operation: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("inspect operation update: %w", err)
	}
	if rows != 1 {
		return errors.New("operation not found")
	}
	return nil
}
