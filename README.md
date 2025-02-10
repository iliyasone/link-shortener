# Link Shortener

## Test application
```bash
curl -X POST -H "Content-Type: application/json" -d '{"originalURL":"https://example.com2"}' http://localhost:8080/shorten
```

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