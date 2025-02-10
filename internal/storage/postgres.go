package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

// reads connection parameters from the environment variables,
// opens a connection to the PostgreSQL database, and runs migrations.
func NewPostgresStorage() (*PostgresStorage, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("database connection parameters missing (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)")
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	// max "parallel" connections
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations using golang-migrate.
	m, err := migrate.New("file://db/migrations", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrations: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// Save inserts a new mapping of shortURL to originalURL.
// If the short URL already exists (unique constraint violation), it returns ErrShortURLExists.
func (p *PostgresStorage) Save(originalURL, shortURL string) error {
	_, err := p.db.Exec(
		`INSERT INTO url_mappings (short_url, original_url) VALUES ($1, $2)`,
		shortURL, originalURL,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation code in PostgreSQL.
				return ErrShortURLExists
			}
		}
		return err
	}
	return nil
}

func (p *PostgresStorage) Get(shortURL string) (string, error) {
	var originalURL string
	err := p.db.QueryRow(
		`SELECT original_url FROM url_mappings WHERE short_url = $1`,
		shortURL,
	).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrURLNotFound
		}
		return "", err
	}
	return originalURL, nil
}

func (p *PostgresStorage) FindByOriginal(originalURL string) (string, error) {
	var shortURL string
	err := p.db.QueryRow(
		`SELECT short_url FROM url_mappings WHERE original_url = $1`,
		originalURL,
	).Scan(&shortURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrURLNotFound
		}
		return "", err
	}
	return shortURL, nil
}
