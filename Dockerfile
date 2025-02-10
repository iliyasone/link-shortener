FROM golang:1.23.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o link-shortener ./cmd/link-shortener

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/link-shortener .

# Copy the migration files
COPY --from=builder /app/db/migrations /root/db/migrations

EXPOSE 8080
CMD ["./link-shortener"]
