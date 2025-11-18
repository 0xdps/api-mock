.PHONY: help build run dev test clean deploy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the server binary
	go build -o bin/server ./cmd/server

run: build ## Build and run the server
	./bin/server

dev: ## Run the server in development mode
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

clean: ## Remove build artifacts
	rm -rf bin/
	go clean

deploy: ## Deploy to Fly.io
	flyctl deploy

logs: ## View Fly.io logs
	flyctl logs

status: ## Check Fly.io app status
	flyctl status

docker-build: ## Build Docker image locally
	docker build -t mockly-api .

docker-run: docker-build ## Build and run Docker image locally
	docker run -p 8080:8080 mockly-api

fmt: ## Format Go code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

tidy: ## Tidy Go modules
	go mod tidy
