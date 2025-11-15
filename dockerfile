FROM golang:1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pr-service ./cmd/server

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM golang:1

WORKDIR /app

COPY --from=builder /app/pr-service .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY --from=builder /app/migrations ./migrations

ENV DATABASE_URL=postgres://postgres:postgres@db:5432/pr_service?sslmode=disable
ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=$DATABASE_URL

ENTRYPOINT goose -dir ./migrations postgres "$DATABASE_URL" up && ./pr-service