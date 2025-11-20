# Mockly API - Go Backend

Modern, schema-driven mock API service built with Go and chi router.

**Production API:** https://api.mockly.codes

## Tech Stack

- **Go 1.23+**
- **chi router** - Lightweight, idiomatic HTTP router
- **gofakeit v7** - Realistic fake data generation (50+ generators)
- **Embedded schemas** - No external files needed
- **Fly.io** - Production deployment

## Project Structure

```
backend/
├── cmd/
│   └── server/              # Main application entry point
│       └── main.go
├── internal/                # Private application code
│   ├── handlers/
│   │   └── dynamic.go       # Dynamic route handlers
│   ├── middleware/
│   │   └── cors.go          # CORS middleware
│   ├── schema/
│   │   ├── loader.go        # Schema registry
│   │   └── embedded/        # Embedded JSON schemas
│   │       ├── user.json
│   │       ├── post.json
│   │       ├── product.json
│   │       └── ...
│   └── static/
│       ├── assets.go        # Embedded static assets
│       └── favicon.go       # Favicon handlers
├── Dockerfile               # Multi-stage Docker build
├── fly.toml                 # Fly.io deployment config
├── Makefile                 # Build commands
├── go.mod                   # Go dependencies
└── README.md                # This file
```

## Quick Start

### Run Locally

```bash
cd backend
go run cmd/server/main.go
```

Server runs on **http://localhost:8080**

### Using Makefile

```bash
cd backend

# Run in development mode
make dev

# Build binary
make build

# Run tests
make test
```

### Build Binary

```bash
cd backend
go build -o bin/server ./cmd/server
./bin/server
```

### Test Endpoints

```bash
# Health check / API info
curl http://localhost:8080/

# Get users
curl http://localhost:8080/users?count=5

# Get single user
curl http://localhost:8080/users/123

# Get resource metadata
curl http://localhost:8080/users/meta

# Production API
curl https://api.mockly.codes/users?count=5
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

## Features

### Schema-Driven Architecture
- Add JSON schemas to auto-generate endpoints
- No code changes needed for new resources
- Custom routes and aliases support
- Meta endpoint for each resource

### Data Generation
- 50+ realistic data generators via gofakeit
- Personal: name, email, username, password
- Location: address, city, country, coordinates
- Internet: URL, domain, IP, UUID
- Dates: past, future, date-time
- Text: word, sentence, paragraph
- And much more!

### CORS Support
- Enabled by default for all origins
- Perfect for frontend development
- Configurable headers and methods

## API Endpoints

All resources support these endpoints:

- `GET /{resource}?count=N` - Get collection (max 100)
- `GET /{resource}/{id}` - Get single item by ID
- `GET /{resource}/meta` - Get schema metadata

### Available Resources

| Resource | Endpoint |
|----------|----------|
| Users | `/users` |
| Posts | `/posts` |
| Products | `/products` |
| Comments | `/comments` |
| Todos | `/todos` |
| Reviews | `/reviews` |

## Adding New Resources

1. Create JSON schema in `internal/schema/embedded/`
2. Restart server - routes auto-generated!
3. No code changes required

**Example schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Order",
  "x-resource": {
    "name": "orders",
    "singular": "order"
  },
  "properties": {
    "id": {
      "type": "integer",
      "x-generator": "random_int"
    }
  }
}
```

## License

MIT
