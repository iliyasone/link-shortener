your task is implement `internal/storage/postgres.go` for this Golang web app
analyze what else i need to add to make it works

I already has implemented the RAM version, and some basic Docker setup for postgres database

If it possible, make the auto migration using golang-migrate

Am I right that it would be check for correct applyied migrations each run?

Instructions:
always follow the best industrial practices about Golang. Give me deep intuition about what is going on.

Code files:
cmd/link-shortener/main.go
```go
package main

import (
	"flag"
	"log"

	"link-shortener/internal/handlers"
	"link-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	storageType := flag.String("storage", "ram", "Storage backend to use: 'ram' or 'postgres'")
	flag.Parse()

	var store storage.Storage
	if *storageType == "postgres" {
		log.Fatal("Postgres storage not implemented yet")
	} else if *storageType == "ram" {
		store = storage.NewRAMStorage()
	} else {
		log.Fatal("Wrong storage type")
	}

	r := gin.Default()
	r.POST("/shorten", handlers.PostShortenURL(store))
	r.GET("/:shortURL", handlers.GetOriginalURL(store))

	log.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
```

internal/storage/storage.go
```go
package storage

import "errors"

var (
	ErrShortURLExists = errors.New("short URL already exists")
	ErrURLNotFound = errors.New("URL not found")
)

type Storage interface {
    // Unsafe save a new short URL in the storage
    //
    // Note: does not check for uniqueness in either direction
    Save(originalURL, shortURL string) error
    // Get original URL by short URL
    Get(shortURL string) (string, error)
    // Checks if the URL already saved in storage
    FindByOriginal(originalURL string) (string, error)
}
```

internal/storage/ram.go
```go
package storage

import (
	"sync"
)

type RAMStorage struct {
    mu    sync.RWMutex
    // key: shortURL, value: originalURL
    forwardMap map[string]string
    // key: originalURL, value: shortURL
	reverseMap  map[string]string
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

```

internal/service/service.go
```go
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

```

docker-compose.db.yaml
```yaml
version: '3.8'
services:
  db:
    image: postgres:13
    restart: always
    environment:
      POSTGRES_USER: shortener
      POSTGRES_PASSWORD: password
      POSTGRES_DB: shortener_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
  
  app:
    build:
      context: .
      dockerfile: Dockerfile
    depends_on:
      - db
    ports:
      - "8080:8080"
    command: ["--storage=postgres"]
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=shortener
      - DB_PASSWORD=password
      - DB_NAME=shortener_db

volumes:
  postgres_data:
```

db/migrations/0001_init.up.sql
```sql
CREATE TABLE IF NOT EXISTS url_mappings (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create an index for fast lookups of original URLs. 
-- This is breaking normalization for speed up
CREATE INDEX idx_original_url ON url_mappings(original_url);

```

db/migrations/0001_init.down.sql
```sql
DROP INDEX IF EXISTS idx_original_url;
DROP TABLE IF EXISTS url_mappings;

```

README.md
```
# Link Shortener

## Run tests
```bash
go test -v ./...
```

## Run the RAM version
```bash
docker-compose -f docker-compose.ram.yml up --build
```
## Run the database version
```bash
docker-compose -f docker-compose.db.yml up --build
```

## Dev tools:

1. Install golang-migrate
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2. Add to the path
```bash
echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

3. Check
```bash
migrate -version
```
```

full task text:
```md
Задание (Стажер-разработчик)

Укорачиватель ссылок

Необходимо реализовать сервис, который должен предоставлять API по созданию сокращенных ссылок следующего формата:
- Ссылка должна быть уникальной и на один оригинальный URL должна ссылаться только одна сокращенная ссылка.
- Ссылка должна быть длинной 10 символов
- Ссылка должна состоять из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание)

Сервис должен быть написан на Go и принимать следующие запросы по http:
1. Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый
2. Метод Get, который будет принимать сокращённый URL и возвращать оригинальный URL

Решение должно быть предоставлено в «конечном виде», а именно:
- Сервис должен быть распространён в виде Docker-образа 
- В качестве хранилища ожидается использовать две реализации. Какое хранилище использовать, указывается параметром при запуске сервиса.  
    - Первое это postgresql.
    - Второе - самостоятельно написать пакет для хранения ссылок в памяти приложения.
- Покрыть реализованный функционал Unit-тестами

Результат предоставить в виде публичного репозитория на github.com

В процессе собеседования-ревью посмотрим:
- Как генерируются ссылки и почему предложенный алгоритм будет работать; насколько он соответствует заданию и прост в понимании.
- Как раскиданы типы по файлам, файлики по пакетам, пакеты по приложению: структуру проекта.
- Как обрабатываются ошибки в разных сценариях использования
- Насколько удобен и логичен сервис в использовании
- Как сервис будет себя вести, если им будут пользоваться одновременно сотни людей (как например youtu.be / ya.cc)
- Что будет, если сервис оставить работать на очень долгое время
- Общую чистоту кода
```