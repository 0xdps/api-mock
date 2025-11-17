# apimock.codes

**Free mock API service for testing and development** - A modern, schema-driven alternative to JSONPlaceholder.

🌐 **Live:** https://apimock-opyavp46a-0xdps-team.vercel.app  
📚 **Docs:** https://apimock-opyavp46a-0xdps-team.vercel.app/docs  
🎮 **Playground:** https://apimock-opyavp46a-0xdps-team.vercel.app/playground

## ✨ Features

- 🚀 **Schema-Driven** - Add new endpoints by creating JSON schemas (zero code!)
- 🎯 **RESTful API** - Standard REST endpoints for all resources
- 💡 **Realistic Data** - Powered by gofakeit with 50+ generators
- ⚡ **Fast & Reliable** - Go backend with Gin framework
- 🌐 **CORS Enabled** - Ready for frontend development
- 🎨 **Modern UI** - Next.js website with interactive playground
- 📦 **No Database** - Generates data on-the-fly
- 🆓 **Free Forever** - Open source and self-hostable

## 📦 Project Structure

```
apimock/
├── api/              # Go API server (Gin + gofakeit)
│   ├── main.go       # Entry point, dynamic route registration
│   ├── schema/       # Schema loader and data generator
│   ├── handlers/     # Generic resource handlers
│   └── shared/       # JSON schemas (copied for deployment)
├── web/              # Next.js 14 website (TypeScript + Tailwind)
│   ├── app/          # Pages (landing, docs, playground)
│   ├── components/   # React components
│   └── shared/       # JSON schemas (copied for deployment)
├── shared/           # Source of truth for schemas
│   └── schemas/      # Resource definitions (user, post, etc.)
└── vercel.json       # Monorepo deployment config
```

## 🚀 Quick Start

### Development Setup

**1. Clone the repository:**
```bash
git clone https://github.com/0xdps/fake-stack.git
cd fake-stack/apimock
```

**2. Start the API:**
```bash
cd api
go run main.go
```
API runs on http://localhost:8080

**3. Start the website:**
```bash
cd web
npm install
npm run dev
```
Website runs on http://localhost:3000

### Using the API

```bash
# Get 10 users
curl http://localhost:8080/api/users?count=10

# Get single user by ID
curl http://localhost:8080/api/users/123

# Get resource metadata
curl http://localhost:8080/api/users/meta

# Get products
curl http://localhost:8080/api/products?count=50
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
- **Collection:** `GET /api/{resource}?count=N` (max 1000)
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

**2. Restart the server:**
```bash
cd api && go run main.go
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

### Deploy to Vercel (Monorepo)

**1. Install Vercel CLI:**
```bash
npm install -g vercel
```

**2. Link project:**
```bash
cd apimock
vercel link
```

**3. Deploy:**
```bash
vercel --prod
```

The deployment configuration in `vercel.json` automatically:
- Builds the Next.js website
- Creates serverless functions for the Go API
- Routes `/api/*` to the API
- Routes everything else to the website

### Environment Variables

Set in Vercel dashboard:
```env
NEXT_PUBLIC_API_URL=https://your-api-domain.vercel.app
```

### Custom Domains

Configure in Vercel project settings:
- **Website:** `apimock.codes` → `/`
- **API:** `api.apimock.codes` → `/api`

### Disable Deployment Protection

1. Go to https://vercel.com/[your-team]/apimock/settings/deployment-protection
2. Select "Disabled" or "Standard" mode
3. Redeploy if necessary

## 🛠️ Technology Stack

### Backend (API)
- **Language:** Go 1.22+
- **Framework:** Gin v1.11.0
- **Data Generation:** gofakeit/v7
- **CORS:** gin-contrib/cors

### Frontend (Website)
- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **Deployment:** Vercel Serverless

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
2. Test locally with `cd api && go run main.go`
3. Submit a PR

See [../CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

## 📁 Related Projects

- **[fakestack](../README.md)** - Database population tool (Python, Node.js, Go)
- Parent project that powers the data generation

## 📄 License

MIT License - see [../LICENSE](../LICENSE)

## 🙏 Acknowledgments

Built with:
- [gofakeit](https://github.com/brianvoe/gofakeit) - Fake data generation
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Next.js](https://nextjs.org/) - React framework
- Part of the [fakestack](https://github.com/0xdps/fake-stack) ecosystem

---

**Made with ❤️ by [Devendra Pratap](https://github.com/0xdps)**
