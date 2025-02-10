package storage

import (
	"sync"
)

type RAMStorage struct {
	mu sync.RWMutex
	// key: shortURL, value: originalURL
	forwardMap map[string]string
	// key: originalURL, value: shortURL
	reverseMap map[string]string
}

func NewRAMStorage() *RAMStorage {
	return &RAMStorage{
		forwardMap: make(map[string]string),
		reverseMap: make(map[string]string),
	}
}

func (r *RAMStorage) Save(originalURL, shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.forwardMap[shortURL]; exists {
		return ErrShortURLExists
	}
	r.forwardMap[shortURL] = originalURL
	r.reverseMap[originalURL] = shortURL
	return nil
}

func (r *RAMStorage) Get(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	original, exists := r.forwardMap[shortURL]
	if !exists {
		return "", ErrURLNotFound
	}
	return original, nil
}

func (r *RAMStorage) FindByOriginal(originalURL string) (string, error) {
	short, exists := r.reverseMap[originalURL]
	if !exists {
		return "", ErrURLNotFound
	}
	return short, nil
}
