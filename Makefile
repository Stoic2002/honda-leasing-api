.PHONY: build run lint tidy wire

APP_NAME = honda-leasing-api
MAIN_FILE = cmd/api/main.go

build:
	go build -o bin/$(APP_NAME) $(MAIN_FILE)

run:
	APP_ENV=dev go run $(MAIN_FILE)

run-dev:
	APP_ENV=dev go run $(MAIN_FILE)

run-staging:
	APP_ENV=staging go run $(MAIN_FILE)

run-prod:
	APP_ENV=prod go run $(MAIN_FILE)

seed:
	APP_ENV=dev go run cmd/seed/main.go

seed-dev:
	APP_ENV=dev go run cmd/seed/main.go

seed-staging:
	APP_ENV=staging go run cmd/seed/main.go

seed-prod:
	APP_ENV=prod go run cmd/seed/main.go

lint:
	golangci-lint run ./...

vet:
	go vet ./...

tidy:
	go mod tidy

wire:
	wire ./cmd/api
