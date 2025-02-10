package storage

import "testing"

func TestRAMStorage(t *testing.T) {
	s := NewRAMStorage()

	originalURL := "http://example.com"
	shortURL := "abc123XYZ_"

	if err := s.Save(originalURL, shortURL); err != nil {
		t.Fatalf("unexpected error saving URL: %v", err)
	}
	if err := s.Save("http://different.com", shortURL); err != ErrShortURLExists {
		t.Fatalf("expected ErrShortURLExists, got: %v", err)
	}
	got, err := s.Get(shortURL)
	if err != nil {
		t.Fatalf("unexpected error getting URL: %v", err)
	}
	if got != originalURL {
		t.Fatalf("expected %s, got %s", originalURL, got)
	}
	gotShort, err := s.FindByOriginal(originalURL)
	if err != nil {
		t.Fatalf("unexpected error in FindByOriginal: %v", err)
	}
	if gotShort != shortURL {
		t.Fatalf("expected %s, got %s", shortURL, gotShort)
	}
	_, err = s.Get("nonexistent")
	if err != ErrURLNotFound {
		t.Fatalf("expected ErrURLNotFound for nonexistent shortURL, got %v", err)
	}
}
