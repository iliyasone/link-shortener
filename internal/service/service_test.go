package service

import (
	"testing"

	"link-shortener/internal/storage"
)

func TestSaveURL(t *testing.T) {
	store := storage.NewRAMStorage()
	t.Run("New URL generates valid short URL", func(t *testing.T) {
		originalURL := "https://example.com"
		shortURL, err := SaveURL(store, originalURL)
		if err != nil {
			t.Fatalf("unexpected error when saving URL: %v", err)
		}
		if len(shortURL) != 10 {
			t.Errorf("expected short URL length of 10, got %d", len(shortURL))
		}
	})

	t.Run("Saving duplicate URL returns the same short URL", func(t *testing.T) {
		originalURL := "https://duplicate.com"
		shortURL1, err := SaveURL(store, originalURL)
		if err != nil {
			t.Fatalf("unexpected error on first save: %v", err)
		}
		shortURL2, err := SaveURL(store, originalURL)
		if err != nil {
			t.Fatalf("unexpected error on duplicate save: %v", err)
		}
		if shortURL1 != shortURL2 {
			t.Errorf("expected same short URL for duplicate original URL, got %q and %q", shortURL1, shortURL2)
		}
	})
}
