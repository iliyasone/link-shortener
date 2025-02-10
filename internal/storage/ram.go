package storage

import (
    "errors"
    "sync"
)

var (
	ErrShortURLExists = errors.New("short URL already exists")
	ErrURLNotFound = errors.New("URL not found")
)

type RAMStorage struct {
    mu    sync.RWMutex
    store map[string]string // key: shortURL, value: originalURL
}

func NewRAMStorage() *RAMStorage {
    return &RAMStorage{
        store: make(map[string]string),
    }
}

func (r *RAMStorage) Save(originalURL, shortURL string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

	if _, exists := r.store[shortURL]; exists {
        return ErrShortURLExists
    }
    r.store[shortURL] = originalURL
    return nil
}


func (r *RAMStorage) Get(shortURL string) (string, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    original, exists := r.store[shortURL]
    if !exists {
        return "", ErrURLNotFound
    }
    return original, nil
}
