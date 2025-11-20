# Build stage
# IMPORTANT: This Dockerfile must be built from the repository root:
#   From repo root: docker build -f backend/Dockerfile .
#   OR: flyctl deploy --config backend/fly.toml --dockerfile backend/Dockerfile
#
# This allows access to shared/schemas/ directory

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files from backend/
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ .

# Copy shared schemas to backend (from repository root)
COPY shared/schemas/*.json internal/schema/embedded/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /bin/server .

EXPOSE 8080

CMD ["./server"]
