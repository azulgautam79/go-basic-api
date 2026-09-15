.PHONY: build run migrate-up migrate-down migrate-force docs

swagger:
	@swag init -g cmd/api/main.go -o docs --parseInternal

build:
	@go build -o bin/api ./cmd/api

run: build
	@./bin/api

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

# Usage: make migrate-force v=1
migrate-force:
	@if [ -z "$(v)" ]; then echo "Error: Please specify version, e.g., 'make migrate-force v=1'"; exit 1; fi
	@go run ./cmd/migrate force $(v)

# Generate Swagger / OpenAPI documentation
docs:
	@swag init -g cmd/api/main.go --parseDependency -o docs