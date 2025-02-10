package generator

import (
	"strings"
	"testing"
)

func TestDefaultGenerator(t *testing.T) {
	gen := NewGenerator()
	url, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() returned an error: %v", err)
	}
	if len(url) != gen.URLLength {
		t.Errorf("expected length %d, got %d", gen.URLLength, len(url))
	}
	for i, ch := range url {
		if !strings.ContainsRune(gen.Charset, ch) {
			t.Errorf("character at index %d, %q, is not in charset %q", i, ch, gen.Charset)
		}
	}
}

func TestGenerateMultiple(t *testing.T) {
	gen := NewGenerator()
	seen := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		url, err := gen.Generate()
		if err != nil {
			t.Fatalf("Generate() returned an error on iteration %d: %v", i, err)
		}
		if len(url) != gen.URLLength {
			t.Errorf("iteration %d: expected length %d, got %d", i, gen.URLLength, len(url))
		}
		for _, ch := range url {
			if !strings.ContainsRune(gen.Charset, ch) {
				t.Errorf("iteration %d: character %q not in allowed charset %q", i, ch, gen.Charset)
			}
		}
		if seen[url] {
			t.Logf("duplicate generated on iteration %d: %q", i, url)
		}
		seen[url] = true
	}
}

// TestCustomGenerator verifies that different length and symbols work properly.
func TestCustomGenerator(t *testing.T) {
	customCharset := "abc123!@#"
	customLength := 8

	gen := &Generator{
		Charset:   customCharset,
		URLLength: customLength,
	}
	url, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() returned an error: %v", err)
	}
	if len(url) != customLength {
		t.Errorf("expected length %d, got %d", customLength, len(url))
	}
	for i, ch := range url {
		if !strings.ContainsRune(customCharset, ch) {
			t.Errorf("character at index %d, %q, is not in charset %q", i, ch, customCharset)
		}
	}
}
