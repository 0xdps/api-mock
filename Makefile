.PHONY: help dev build test clean sync-schemas docker-build deploy api-dev web-dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Schema and Type Management
sync-schemas: ## Sync schemas to backend and generate frontend types
	@echo "Syncing schemas and generating types..."
	@npm run sync-schemas
	@echo "✓ Synced schemas and generated types"

generate-types: ## Generate TypeScript types from schemas
	@echo "Generating TypeScript types..."
	@npm run generate-types
	@echo "✓ Generated types"

# Backend
api-dev: sync-schemas ## Run backend development server
	@echo "Starting backend development server..."
	@if [ -f backend/.env ]; then \
		echo "Loading environment variables from backend/.env"; \
		export $$(cat backend/.env | grep -v '^#' | xargs) && \
		echo "Cache config: ITEMS_PER_RESOURCE=$${CACHE_ITEMS_PER_RESOURCE:-100}, SEED=$${CACHE_SEED:-42}" && \
		cd backend && go run cmd/server/main.go; \
	else \
		echo "No .env file found, using defaults"; \
		echo "Cache config: ITEMS_PER_RESOURCE=$${CACHE_ITEMS_PER_RESOURCE:-100}, SEED=$${CACHE_SEED:-42}" && \
		cd backend && go run cmd/server/main.go; \
	fi

api-build: sync-schemas ## Build backend binary
	@echo "Building backend binary..."
	@cd backend && go build -o bin/server ./cmd/server
	@echo "✓ Built binary: backend/bin/server"

api-test: ## Run backend tests
	@cd backend && go test ./...

# Frontend
web-dev: generate-types ## Run frontend development server
	@echo "Starting frontend development server..."
	@cd frontend && NODE_OPTIONS='--no-warnings' npm run dev

web-build: generate-types ## Build frontend
	@echo "Building frontend..."
	@cd frontend && npm run build
	@echo "✓ Built frontend"

web-install: ## Install frontend dependencies
	@cd frontend && npm install

# Combined Development
dev: ## Run both API and web in parallel (requires npm concurrently)
	@npm run dev

# Docker
docker-build: ## Build Docker image for backend
	@echo "Building Docker image..."
	@docker build -t api-mockly:latest -f Dockerfile .
	@echo "✓ Built Docker image: api-mockly:latest"

docker-run: ## Run Docker container locally
	@echo "Running Docker container..."
	@docker run -p 8080:8080 api-mockly:latest

docker-test: docker-build docker-run ## Build and run Docker locally

# Deployment
deploy-api: ## Deploy backend to Fly.io
	@echo "Deploying backend to Fly.io..."
	@echo "Note: Schemas are committed to git via pre-commit hook"
	@flyctl deploy
	@echo "✓ Deployed to Fly.io"

deploy-railway: ## Deploy backend to Railway
	@echo "Deploying backend to Railway..."
	@railway up
	@echo "✓ Deployed to Railway"

railway-init: ## Initialize Railway project and add Redis
	@echo "Initializing Railway project..."
	@railway init
	@echo "Adding Redis service..."
	@railway add redis
	@echo "Setting environment variables..."
	@railway variables set REDIS_DB=0
	@railway variables set REDIS_TLS_ENABLED=false
	@railway variables set CACHE_MODE=all
	@railway variables set CACHE_ITEMS_PER_RESOURCE=100
	@railway variables set CACHE_SEED=42
	@railway variables set MAX_ITEMS_PER_RESOURCE=1000
	@echo "✓ Railway project initialized"
	@echo "Deploy with: make deploy-railway"

railway-logs: ## View Railway deployment logs
	@railway logs

railway-status: ## Check Railway deployment status
	@railway status

deploy-web: ## Deploy frontend to Vercel
	@echo "Deploying frontend to Vercel..."
	@cd frontend && vercel --prod
	@echo "✓ Deployed to Vercel"

deploy: deploy-api deploy-web ## Deploy both backend (Fly.io) and frontend

# Build
build: api-build web-build ## Build both backend and frontend

# Clean
clean: ## Clean build artifacts
	@rm -rf backend/bin/
	@rm -rf frontend/.next/
	@rm -rf frontend/out/
	@echo "✓ Cleaned build artifacts"

# Dependencies
deps: ## Install all dependencies
	@echo "Installing root dependencies..."
	@npm install
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "Downloading backend dependencies..."
	@cd backend && go mod download && go mod tidy
	@echo "✓ Installed all dependencies"

