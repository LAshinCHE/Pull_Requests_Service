version: "3.8"
services:
  pdb:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: pr
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data

  app:
    build: .
    depends_on:
      - pdb
    environment:
      DATABASE_URL: postgres://postgres:postgres@db:5432/pr?sslmode=disable
      PORT: 8080
    ports:
      - "8080:8080"
    command: ["make run"]   
volumes:
  db_data: