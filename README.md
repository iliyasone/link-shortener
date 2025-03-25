# Link Shortener

A simple URL shortener service written in Go. It provides an HTTP API to create and resolve shortened URLs. The service is distributed as a Docker image and supports two storage backends:
- **In-Memory (RAM)** – for fast, ephemeral storage.
- **PostgreSQL** – for persistent storage, complete with auto migration using [golang-migrate](https://github.com/golang-migrate/migrate).

---

## Overview

The Link Shortener service accepts HTTP requests to shorten long URLs and retrieve the original URLs from their shortened counterparts. It guarantees that:
- Each original URL corresponds to one unique short URL.
- Short URLs are exactly 10 characters long.
- Allowed characters include lowercase and uppercase Latin letters, digits, and the underscore (`_`).

## Key Idea

For any given original URL, the algorithm generates a cryptographically secure, random short-hand URL. The storage system also maintains a reverse mapping, which allows for quick retrieval and verification of whether the URL has been stored previously. This mechanism is implemented consistently in both PostgreSQL and in-memory (RAM) storage solutions.

## Project's file structure
```
.
├── Dockerfile
├── README.md
├── go.mod
├── go.sum
├── docker-compose.db.yaml
├── docker-compose.ram.yml
├── cmd
│   └── link-shortener
│       └── main.go
├── db
│   └── migrations
│       ├── 0001_init.up.sql
│       └── 0001_init.down.sql
├── internal
│   ├── handlers
│   │   └── handlers.go
│   ├── service
│   │   └── service.go
│   └── storage
│       ├── storage.go
│       ├── ram.go
│       └── postgres.go
├── pkg
│   └── generator
│       └── generator.go
└── .github
    └── workflows
        └── go-tests.yml
```


## API Endpoints

### POST `/shorten`
- **Description:** Accepts a JSON payload with an original URL, saves it, and returns a shortened URL.
- **Example Request:**
  ```bash
  curl -X POST -H "Content-Type: application/json" \
       -d '{"originalURL":"https://example.com"}' \
       http://localhost:8080/shorten
  ```

### GET `/{shortURL}`
- **Description:** Redirects to the original URL associated with the provided shortened URL.
- **Example Request:**
  ```bash
  curl http://localhost:8080/abcDEF123_
  ```

---

## How to Run

### Run Tests
Make sure your tests pass before deploying:
```bash
go test -v ./...
```

### Run the In-Memory (RAM) Version
For quick testing with the in-memory storage backend:
```bash
docker compose -f docker-compose.ram.yaml up --build
```

### Run the PostgreSQL (Database) Version
To run with a persistent PostgreSQL backend (ensure your Docker environment is properly configured):
```bash
docker compose -f docker-compose.db.yaml up --build -d
```

---

## Dev Tools

### Installing golang-migrate

1. **Install:**
   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

2. **Add to PATH:**
   ```bash
   echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **Verify Installation:**
   ```bash
   migrate -version
   ```

The auto migration is set up to run on application startup when using the PostgreSQL storage. It checks for pending migrations in the `db/migrations` directory and applies them automatically.