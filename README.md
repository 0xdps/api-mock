<div align="center">
  <img src="assets/logo.svg" alt="Mockly Logo" width="120" height="120">
  <h1>Mockly</h1>
  <p><strong>Free mock API service for testing and development</strong></p>
  <p>A modern, schema-driven alternative to JSONPlaceholder</p>
</div>

---

🌐 **Live:** https://mockly.codes  
📚 **Docs:** https://mockly.codes/docs  
🎮 **Playground:** https://mockly.codes/playground

## ✨ Features

- 🚀 **Schema-Driven** - Add new endpoints by creating JSON schemas (zero code!)
- 🎯 **RESTful API** - Standard REST endpoints for all resources
- 💡 **Realistic Data** - Powered by gofakeit with 50+ generators
- ⚡ **Fast & Reliable** - Go backend with chi router
- 🌐 **CORS Enabled** - Ready for frontend development
- 🎨 **Modern UI** - Next.js website with interactive playground
- 📦 **No Database** - Generates data on-the-fly
- 🆓 **Free Forever** - Open source and self-hostable

## 📦 Project Structure

```
mockly/
├── backend/          # Go API service (chi + gofakeit)
│   ├── cmd/          # Application entry points
│   ├── internal/     # Handlers, middleware, schema
│   ├── Dockerfile    # Docker build config
│   └── fly.toml      # Fly.io deployment config
├── frontend/         # Next.js 14 website (TypeScript + Tailwind)
│   ├── app/          # Pages (landing, docs, playground)
│   ├── components/   # React components
│   └── lib/          # Utility functions
├── shared/           # Source of truth for schemas & scripts
│   ├── schemas/      # Resource definitions (user, post, etc.)
│   └── scripts/      # Build scripts (generate-types.js)
└── package.json      # Root scripts for dev workflow
```

## 🚀 Quick Start

### Development Setup

**1. Clone the repository:**
```bash
git clone https://github.com/0xdps/api-mockly.git
cd api-mockly
```

**2. Install dependencies:**
```bash
# Install root dependencies (concurrently)
npm install

# Install web dependencies
npm run web:install
```

**3. Start both API and website:**
```bash
npm run dev
```

This starts:
- API server on http://localhost:8080
- Website on http://localhost:3000

**Or run separately:**
```bash
# API only
npm run api:dev

# Website only  
npm run web:dev
```

### Using the API

```bash
# Get 10 users
curl 'http://localhost:8080/api/users?count=10'

# Get single user by ID
curl http://localhost:8080/api/users/123

# Get resource metadata
curl http://localhost:8080/api/users/meta

# Get products
curl 'http://localhost:8080/api/products?count=50'
```

## 📚 Available Resources

| Resource | Endpoint | Fields |
|----------|----------|--------|
| Users | `/api/users` | id, username, email, name, avatar, bio, etc. |
| Posts | `/api/posts` | id, user_id, title, content, published_at, etc. |
| Products | `/api/products` | id, name, description, price, category, etc. |
| Comments | `/api/comments` | id, post_id, user_id, content, created_at |
| Todos | `/api/todos` | id, user_id, title, completed, due_date |
| Reviews | `/api/reviews` | id, product_id, user_id, rating, comment |

All endpoints support:
- **Collection:** `GET /api/{resource}?count=N` (max 100)
- **Single Item:** `GET /api/{resource}/:id`
- **Metadata:** `GET /api/{resource}/meta`

## 🎯 Schema-Driven Development

### Adding a New Resource (Zero Code!)

**1. Create a schema** in `shared/schemas/`:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Order",
  "type": "object",
  "x-resource": {
    "name": "orders",
    "singular": "order",
    "description": "E-commerce orders",
    "routes": {
      "path": "/orders",
      "methods": ["GET"],
      "aliases": ["/purchases"]
    }
  },
  "properties": {
    "id": {
      "type": "integer",
      "x-generator": "random_int",
      "x-generator-params": { "min": 1, "max": 10000 }
    },
    "total": {
      "type": "number",
      "x-generator": "random_int",
      "x-generator-params": { "min": 10, "max": 500 }
    },
    "status": {
      "type": "string",
      "x-generator": "word"
    }
  }
}
```

**2. Restart the API server:**
```bash
cd backend && make dev
# or: go run cmd/server/main.go
```

**That's it!** Your new endpoint is live at `/api/orders` 🎉

### Custom Routes

Define custom paths, aliases, and methods in schemas:

```json
{
  "x-resource": {
    "name": "todos",
    "routes": {
      "path": "/v1/todos",
      "aliases": ["/tasks", "/todo-items"],
      "methods": ["GET", "POST"]
    }
  }
}
```

This creates:
- ✅ `GET /api/v1/todos`
- ✅ `GET /api/tasks` (alias)
- ✅ `GET /api/todo-items` (alias)

### Supported Generators

- **Personal:** `name`, `first_name`, `email`, `username`, `password`
- **Address:** `address`, `city`, `country`, `zip_code`, `latitude`
- **Company:** `company`, `job`, `catch_phrase`
- **Internet:** `url`, `domain_name`, `ipv4`, `uuid`, `mac_address`
- **Dates:** `date`, `date_time`, `past_date`, `future_date`
- **Text:** `word`, `sentence`, `paragraph`, `text`
- **Numbers:** `random_int`, `random_digit`, `random_number`
- **Other:** `phone_number`, `boolean`, `user_agent`

## 🚀 Deployment

### API: Deploy to Fly.io

**1. Install Fly CLI:**
```bash
curl -L https://fly.io/install.sh | sh
```

**2. Login:**
```bash
flyctl auth login
```

**3. Launch (first time):**
```bash
cd backend
flyctl launch
```

**4. Deploy updates:**
```bash
flyctl deploy
```

**5. View logs:**
```bash
flyctl logs
```

The API will be available at: `https://mockly-api.fly.dev`

### Website: Deploy to Vercel

**1. Install Vercel CLI:**
```bash
npm install -g vercel
```

**2. Deploy:**
```bash
cd frontend
vercel --prod
```

**3. Set environment variable:**
```bash
vercel env add NEXT_PUBLIC_API_URL
# Enter: https://mockly-api.fly.dev
```

### Custom Domain Setup

**For API (Fly.io):**
```bash
flyctl certs add api.mockly.codes
```

**For Website (Vercel):**
- Add `mockly.codes` in Vercel project settings
- Configure DNS to point to Vercel

## 🛠️ Technology Stack

### Backend (API)
- **Language:** Go 1.23+
- **Router:** chi v5 (lightweight, idiomatic)
- **Data Generation:** gofakeit/v7
- **CORS:** go-chi/cors
- **Deployment:** Fly.io

### Frontend (Website)
- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **Deployment:** Vercel

## 📖 Documentation

### API Documentation
Visit the `/docs` page for:
- Quick start guide (4 languages: JavaScript, cURL, Python, Node.js)
- Complete endpoint reference
- Schema documentation for all resources
- Live "Try It" buttons

### Interactive Playground
Visit the `/playground` page to:
- Test endpoints without writing code
- Select resources and endpoint types
- Adjust parameters (count, ID)
- See live JSON responses
- Copy code examples (cURL, JavaScript, Python)

## 🤝 Contributing

Contributions welcome! The easiest way to contribute is to add new resource schemas:

1. Create `shared/schemas/your-resource.json`
2. Test locally with `cd local && go run main.go`
3. Submit a PR

See [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

## � License

MIT License - see [LICENSE](./LICENSE)

## 🙏 Acknowledgments

Built with:
- [gofakeit](https://github.com/brianvoe/gofakeit) - Fake data generation
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Next.js](https://nextjs.org/) - React framework

---

**Made with ❤️ by [Devendra Pratap](https://github.com/0xdps)**
