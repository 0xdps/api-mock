.PHONY: help dev build test clean sync-schemas docker-build deploy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

sync-schemas: ## Sync schemas from shared directory
	@echo "Syncing schemas from shared/schemas..."
	@cp ../shared/schemas/*.json internal/schema/embedded/
	@echo "✓ Synced schemas to backend"

dev: sync-schemas ## Run development server with schema sync
	@echo "Starting development server..."
	@go run cmd/server/main.go

build: sync-schemas ## Build production binary with schema sync
	@echo "Building production binary..."
	@go build -o bin/server ./cmd/server
	@echo "✓ Built binary: bin/server"

docker-build: ## Build Docker image (run from repository root)
	@echo "Building Docker image with schemas..."
	@cd .. && docker build -f backend/Dockerfile -t api-mockly:latest .
	@echo "✓ Built Docker image: api-mockly:latest"

docker-run: ## Run Docker container locally
	@echo "Running Docker container..."
	@docker run -p 8080:8080 api-mockly:latest

test: ## Run tests
	@go test ./...

clean: ## Clean build artifacts
	@rm -rf bin/
	@echo "✓ Cleaned build artifacts"

deps: ## Download dependencies
	@go mod download
	@go mod tidy

deploy: sync-schemas ## Deploy to Fly.io
	@echo "✓ Schemas synced to backend"
	@echo "Deploying to Fly.io..."
	@flyctl deploy
	@echo "✓ Deployed to Fly.io"
