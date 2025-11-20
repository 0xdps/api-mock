# Build stage
# This Dockerfile works with build context = backend/ directory
# Schemas must be synced BEFORE deployment via: make sync-schemas

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy all backend source code
COPY . .

# Note: Schemas should already be in internal/schema/embedded/
# They are synced before deployment via: make sync-schemas
# Verify they exist
RUN ls -la internal/schema/embedded/*.json || echo "WARNING: No schemas found!"

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
