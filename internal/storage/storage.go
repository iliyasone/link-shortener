package storage

import "errors"

var (
	ErrShortURLExists = errors.New("short URL already exists")
	ErrURLNotFound    = errors.New("URL not found")
)

type Storage interface {
	// Unsafe save a new short URL in the storage
	//
	// Note: does not check for uniqueness in reverse direction
	// return ErrShortURLExists if shortURL alredy used
	Save(originalURL, shortURL string) error
	// Get original URL by short URL
	Get(shortURL string) (string, error)
	// Checks if the URL already saved in storage
	FindByOriginal(originalURL string) (string, error)
}
