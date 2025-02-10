package storage

type Storage interface {
    Save(originalURL, shortURL string) error
    Get(shortURL string) (string, error)
    FindByOriginal(originalURL string) (string, error)
}
