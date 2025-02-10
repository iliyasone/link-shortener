package storage

import (
    "errors"
    "sync"
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
        return errors.New("short URL already exists")
    }
    r.store[shortURL] = originalURL
    return nil
}


func (r *RAMStorage) Get(shortURL string) (string, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    original, exists := r.store[shortURL]
    if !exists {
        return "", errors.New("URL not found")
    }
    return original, nil
}
