# Mockly API - Go Service

Modern, schema-driven mock API service built with Go and chi router.

## Tech Stack

- **Go 1.23+**
- **chi router** - Lightweight, idiomatic HTTP router
- **gofakeit v7** - Realistic fake data generation
- **Fly.io** - Production deployment

## Project Structure

```
go/
├── cmd/
│   └── server/          # Main application entry point
│       └── main.go
├── internal/            # Private application code
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware (CORS, etc.)
│   └── schema/          # Schema loading and data generation
│       └── embedded/    # Embedded JSON schemas
├── Dockerfile           # Multi-stage Docker build
├── fly.toml             # Fly.io configuration
└── go.mod               # Go module dependencies
```

## Development

### Run Locally

```bash
cd go
go run cmd/server/main.go
```

Server runs on `http://localhost:8080`

### Build

```bash
go build -o bin/server ./cmd/server
./bin/server
```

### Test Endpoints

```bash
# Health check
curl http://localhost:8080/

# Get users
curl http://localhost:8080/api/users?count=5

# Get single user
curl http://localhost:8080/api/users/123

# Get resource metadata
curl http://localhost:8080/api/users/meta
```

## Deployment

### Deploy to Fly.io

**1. Install Fly CLI:**
```bash
curl -L https://fly.io/install.sh | sh
```

**2. Login to Fly.io:**
```bash
flyctl auth login
```

**3. Launch the app (first time):**
```bash
cd go
flyctl launch
```

**4. Deploy updates:**
```bash
flyctl deploy
```

**5. Open the app:**
```bash
flyctl open
```

### Environment Variables

- `PORT` - Server port (default: 8080)

## Router: Chi vs Gin

We use **chi** instead of Gin for several advantages:

### Why Chi?

- **Lightweight** - Minimal dependencies, just standard library
- **Idiomatic** - Uses standard `http.Handler` interface
- **Better Routing** - More flexible routing patterns and middleware composition
- **No Hidden Magic** - Explicit, straightforward code
- **Standard Library** - Uses `net/http` directly, better compatibility

### Migration Notes

**Gin:**
```go
r.GET("/users/:id", handler)
c.Param("id")
c.JSON(200, data)
```

**Chi:**
```go
r.Get("/users/{id}", handler)
chi.URLParam(r, "id")
json.NewEncoder(w).Encode(data)
```

## Adding New Resources

1. Add JSON schema to `internal/schema/embedded/`
2. Restart server - routes are auto-generated!

No code changes needed.

## License

MIT
