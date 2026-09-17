.PHONY: run build test lint tidy \
	docker-up docker-down docker-logs \
	migrate-up migrate-down migrate-create \
	migrate-sqlite-up migrate-sqlite-down migrate-sqlite-create \
	ui-install ui-dev ui-build ui-lint

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

lint:
	go vet ./...

tidy:
	go mod tidy

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

DB_NAME ?= corelog
DB_URL ?= postgres://corelog:corelog@localhost:5432/$(DB_NAME)?sslmode=disable

migrate-up:
	migrate -path migrations/postgres -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations/postgres -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations/postgres -seq $(name)

SQLITE_DB_PATH ?= data/$(DB_NAME).db

migrate-sqlite-up:
	migrate -path migrations/sqlite -database "sqlite3://$(SQLITE_DB_PATH)" up

migrate-sqlite-down:
	migrate -path migrations/sqlite -database "sqlite3://$(SQLITE_DB_PATH)" down 1

migrate-sqlite-create:
	migrate create -ext sql -dir migrations/sqlite -seq $(name)

ui-install:
	cd ui && npm install

ui-dev:
	cd ui && npm run dev

ui-build:
	cd ui && npm run build

ui-lint:
	cd ui && npm run lint
