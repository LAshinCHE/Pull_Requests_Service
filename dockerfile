FROM golang:1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pr-service ./cmd

RUN go install github.com/pressly/goose/v3/cmd/goose@latest


FROM debian:bookworm-slim

WORKDIR /app

# Устанавливаем goose
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Приложение
COPY --from=builder /app/pr-service .
COPY --from=builder /app/migrations ./migrations

ENV DATABASE_URL=postgres://postgres:postgres@db:5432/pr_service?sslmode=disable

# ВАЖНО — shell команда, чтобы работали переменные и &&
ENTRYPOINT ["/bin/sh", "-c", "goose -dir ./migrations postgres \"$DATABASE_URL\" up && ./pr-service"]