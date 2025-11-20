# Build stage
# Simple build - schemas are already committed to git via pre-commit hook

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code (including committed schemas in internal/schema/embedded/)
COPY . .

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
