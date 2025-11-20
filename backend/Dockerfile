# Build stage
# NOTE: This Dockerfile should be built from the repository root:
#   docker build -f backend/Dockerfile .
#
# This allows access to shared/schemas/ directory

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install make for schema syncing
RUN apk add --no-cache make

# Copy go mod files from backend/
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ .

# Copy shared schemas to backend
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
