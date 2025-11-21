# Build stage
# Schemas are copied from shared/schemas/ during build (single source of truth)

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files from backend directory
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Copy schemas from shared/ (source of truth) to embedded directory (flattened)
COPY shared/schemas/*/*.json ./internal/schema/embedded/

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

