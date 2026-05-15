# api-mockly justfile
# Run `just` to see available recipes.

# Default: list all recipes
default:
    @just --list

# ── Setup ─────────────────────────────────────────────────────────────────────

# First-time setup: install all dependencies and copy env files
setup:
    @echo "Installing Node dependencies..."
    pnpm install
    @echo "Downloading Go dependencies..."
    cd backend && go mod download && go mod tidy
    @just env-copy
    @echo ""
    @echo "Done. Edit .env.backend and .env.frontend in the root, then run:"
    @echo "  just dev      — run everything"
    @echo "  just be       — backend only"
    @echo "  just fe       — frontend only"

# Copy root env files to their respective locations
env-copy:
    @echo "Copying .env.backend → backend/.env"
    cp .env.backend backend/.env
    @echo "Copying .env.frontend → frontend/.env.local"
    cp .env.frontend frontend/.env.local

# ── Code generation ───────────────────────────────────────────────────────────

# Generate TypeScript types from shared schemas and sync to backend
gen:
    pnpm --filter @api-mockly/shared generate-types

# ── Local development ─────────────────────────────────────────────────────────

# Run backend only (copies env first)
be: env-copy
    cd backend && go run ./cmd/server

# Run frontend only (copies env first)
fe: env-copy gen
    cd frontend && pnpm dev

# Run backend and frontend together using Turborepo
dev: env-copy gen
    pnpm dev

# ── Build ─────────────────────────────────────────────────────────────────────

# Build backend binary to backend/bin/server
build-be: gen
    cd backend && go build -o bin/server ./cmd/server
    @echo "Built → backend/bin/server"

# Build frontend for production
build-fe: gen
    cd frontend && pnpm build

# Build everything
build: build-be build-fe

# ── Test ──────────────────────────────────────────────────────────────────────

# Run backend tests
test-be:
    cd backend && go test ./...

# Run backend tests with verbose output
test-be-v:
    cd backend && go test -v ./...

# ── Docker (local stack) ──────────────────────────────────────────────────────

# Spin up the full local stack (redis + backend + frontend)
up:
    @echo "Loading .env.frontend for NEXT_PUBLIC_ build args..."
    set -a && . ./.env.frontend && set +a && \
    docker compose up --build

# Same as up but run in the background
up-d:
    @echo "Loading .env.frontend for NEXT_PUBLIC_ build args..."
    set -a && . ./.env.frontend && set +a && \
    docker compose up --build -d

# Stop the local stack
down:
    docker compose down

# Tail logs (all services, or pass a service name: just logs backend)
logs service='':
    docker compose logs -f {{ service }}

# Rebuild a single service without restarting others: just rebuild backend
rebuild service:
    set -a && . ./.env.frontend && set +a && \
    docker compose up --build -d {{ service }}

# ── Docker (image builds only) ────────────────────────────────────────────────

# Build backend Docker image
docker-build:
    docker build -t api-mockly:latest -f Dockerfile .
    @echo "Built → api-mockly:latest"

# Run backend Docker image locally
docker-run:
    docker run -p 8080:8080 --env-file backend/.env api-mockly:latest

# Build and run Docker locally
docker: docker-build docker-run

# ── Deploy ────────────────────────────────────────────────────────────────────

# Deploy backend to Fly.io
deploy-be:
    flyctl deploy

# Deploy frontend to Vercel
deploy-fe:
    cd frontend && vercel --prod

# Deploy backend to Railway
deploy-railway:
    railway up

# View Railway logs
railway-logs:
    railway logs

# ── Clean ─────────────────────────────────────────────────────────────────────

# Remove all build artifacts
clean:
    rm -rf backend/bin/ frontend/.next/ frontend/out/
    @echo "Cleaned build artifacts"
