# Mockly API

Go serverless function handler for Vercel deployment with dynamic route registration from JSON schemas.

## Note

This directory contains the Vercel serverless function handler (`index.go`). 

For **local development**, use the `local/` directory instead:

```bash
cd ../local
go run main.go
```

Server runs on http://localhost:8080

## Stack

- Go 1.22+
- Gin v1.11.0
- gofakeit/v7
- CORS enabled

## Documentation

See [../README.md](../README.md) for:
- Complete API documentation
- Schema-driven development guide
- Deployment instructions
- Available endpoints and resources
