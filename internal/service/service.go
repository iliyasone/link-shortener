package service

import (
	"errors"
	"link-shortener/internal/storage"
	"link-shortener/pkg/generator"
)

var ErrTooManyAttemps = errors.New("failed to generate a unique short URL after many attempts")

func SaveURL(store storage.Storage, originalURL string) (string, error) {
	if existing, err := store.FindByOriginal(originalURL); err == nil {
		return existing, nil
	}
	const maxRetries = 1000
	defaultGenerator := generator.NewGenerator()
	for i := 0; i < maxRetries; i++ {
		candidate, err := defaultGenerator.Generate()
		if err != nil {
			return "", err
		}
		err = store.Save(originalURL, candidate)
		if err == nil {
			return candidate, nil
		}
		if errors.Is(err, storage.ErrShortURLExists) {
			continue
		}
		// another error, return
		return "", err
	}
	return "", ErrTooManyAttemps
}
