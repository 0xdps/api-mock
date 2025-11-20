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
- 🎯 **54 Resources** - From users to weather, stocks to movies
- 📂 **14 Categories** - Resources organized into logical groups (people, commerce, content, etc.)
- 💡 **Realistic Data** - Powered by gofakeit with 200+ generators
- ⚡ **Fast & Reliable** - Go backend with chi router
- 🌐 **CORS Enabled** - Ready for frontend development
- 🎨 **Modern UI** - Next.js website with SSR and interactive playground
- 📦 **No Database** - Generates data on-the-fly
- 🆓 **Free Forever** - Open source and self-hostable
- ⚡ **Server-Side Rendering** - Fast page loads with fresh data
- 🤖 **Automated Build** - TypeScript types auto-generated from schemas

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

This automatically:
- ✅ Generates TypeScript types from schemas
- ✅ Starts API server on http://localhost:8080
- ✅ Starts website on http://localhost:3000

**Note:** Types are auto-generated on every build. No manual steps needed!

**Or run separately:**
```bash
# API only
npm run api:dev

# Website only  
npm run web:dev
```

### Using the API

```bash
# Browse by group (metadata only)
curl http://localhost:8080/people

# Get resources via group path
curl 'http://localhost:8080/people/users?count=10'
curl http://localhost:8080/people/users/123

# Or access resources directly
curl 'http://localhost:8080/users?count=10'
curl http://localhost:8080/users/123

# Get resource metadata
curl http://localhost:8080/users/meta

# Commerce group examples
curl http://localhost:8080/commerce
curl 'http://localhost:8080/commerce/products?count=5'

# Use production API
curl 'https://api.mockly.codes/people/users?count=10'
```

## 📚 Available Resources (54 Endpoints!)

### 📂 Resource Groups

Resources are organized into 14 logical categories:

| Group | Resources | Description |
|-------|-----------|-------------|
| 👥 **People** | users, contacts, students, players | User profiles and people |
| 💼 **Business** | companies, organizations, jobs, meetings, invoices, subscriptions | Business entities |
| 🛒 **Commerce** | products, orders, payments, coupons, categories, tags | E-commerce |
| 📝 **Content** | articles, posts, news, podcasts | Written content |
| 💬 **Social** | comments, reviews, messages, notifications, testimonials | Social interactions |
| 🎬 **Media** | movies, books, albums, videos, images | Entertainment media |
| ✈️ **Travel** | hotels, flights, restaurants, properties, cars | Travel & hospitality |
| 🌍 **Location** | countries, cities, weather | Geographic data |
| 💰 **Finance** | currencies, stocks, crypto | Financial data |
| 🍔 **Food** | recipes | Food & cooking |
| 🎓 **Education** | courses | Educational content |
| ⚽ **Sports** | matches, teams | Sports data |
| ✅ **Productivity** | todos, notes, projects, tasks, tickets, events | Task management |
| 📚 **Reference** | faqs, quotes, languages | Reference data |

### API Endpoints

**Group Endpoints:**
```bash
GET /{group}                    # Group info (metadata only)
GET /{group}/{resource}         # Collection via group path
GET /{group}/{resource}/:id     # Single item via group path
```

**Direct Resource Endpoints:**
```bash
GET /{resource}?count=N         # Collection (max 100)
GET /{resource}/:id             # Single item
GET /{resource}/meta            # Resource metadata
```

**Examples:**
```bash
# Browse by category
curl https://api.mockly.codes/people
# Returns: {"group": "people", "resources": ["users", "contacts", ...], "count": 4}

# Get data via category path
curl 'https://api.mockly.codes/people/users?count=5'

# Or access directly (backwards compatible)
curl 'https://api.mockly.codes/users?count=5'
```

**Production API:** https://api.mockly.codes

## 🎯 Schema-Driven Development

### Adding a New Resource (Zero Code!)

**1. Create a schema** in `shared/schemas/`:

> **Note:** TypeScript types are automatically generated from schemas during build. No manual steps required!

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Order",
  "type": "object",
  "x-resource": {
    "name": "orders",
    "singular": "order",
    "description": "E-commerce orders",
    "group": "commerce",
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

**That's it!** Your new endpoint is live at `/orders` 🎉

The build process automatically:
- ✅ Syncs schemas to backend
- ✅ Generates TypeScript types for frontend
- ✅ Loads new routes in API server

See [BUILD_PROCESS.md](./BUILD_PROCESS.md) for details.

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
- ✅ `GET /v1/todos`
- ✅ `GET /tasks` (alias)
- ✅ `GET /todo-items` (alias)

### Supported Generators (200+)

We support over 200 generators via gofakeit. Popular ones:

**Personal:** `name`, `first_name`, `email`, `username`, `password`, `avatar`  
**Address:** `address`, `city`, `country`, `zip_code`, `latitude`, `longitude`  
**Company:** `company`, `job`, `catch_phrase`, `company_name`, `job_title`  
**Internet:** `url`, `domain_name`, `ipv4`, `uuid`, `mac_address`, `image_url`  
**Dates:** `date`, `date_time`, `past_date`, `future_date`, `time_zone`  
**Text:** `word`, `sentence`, `paragraph`, `text`, `quote_text`  
**Numbers:** `random_int`, `random_digit`, `random_number`, `float32`  
**Weather:** `weather_temperature`, `weather_description`, `weather_humidity`  
**Finance:** `currency_code`, `exchange_rate`, `stock_symbol`, `crypto_name`  
**Geography:** `country_name`, `city_name`, `language_name`, `capital_city`  
**Media:** `movie_title`, `book_title`, `album_title`, `video_title`  
**Travel:** `flight_number`, `hotel_name`, `restaurant_name`, `recipe_name`  
**Business:** `invoice_number`, `order_id`, `payment_status`, `ticket_id`  
**Other:** `phone_number`, `boolean`, `user_agent`, `car_model`

See all generators in [`backend/internal/schema/loader.go`](./backend/internal/schema/loader.go)

## 🚀 Deployment

### Quick Deploy

**Backend (Fly.io):**
```bash
cd backend
make deploy  # Auto-syncs schemas and deploys
```

**Frontend (Vercel):**
```bash
cd frontend
vercel --prod  # Auto-generates types and deploys
```

### Important: Schema Syncing

The deployment process **automatically includes** all schemas:

- ✅ Backend: Schemas copied from `shared/` during Docker build
- ✅ Frontend: Types auto-generated via `prebuild` hook
- ✅ No manual steps required!

**See [DEPLOYMENT.md](./DEPLOYMENT.md) for complete deployment guide.**

### Production URLs

- **API:** https://api.mockly.codes
- **Website:** https://mockly.codes
- **Docs:** https://mockly.codes/docs
- **Playground:** https://mockly.codes/playground

## 🛠️ Technology Stack

### Backend (API)
- **Language:** Go 1.23+
- **Router:** chi v5 (lightweight, idiomatic)
- **Data Generation:** gofakeit/v7
- **CORS:** go-chi/cors
- **Deployment:** Fly.io

### Frontend (Website)
- **Framework:** Next.js 14 (App Router with SSR)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **Rendering:** Server-Side with ISR (5-minute cache)
- **Deployment:** Vercel
- **API Integration:** https://api.mockly.codes

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
2. Test locally with `npm run dev` (types auto-generate!)
3. Submit a PR

See [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

## 📚 Documentation

- **[QUICK_START.md](./QUICK_START.md)** - Quick start guide
- **[BUILD_PROCESS.md](./BUILD_PROCESS.md)** - Build automation
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Deployment guide
- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Contribution guidelines
- **[CHANGELOG.md](./CHANGELOG.md)** - Version history

## � License

MIT License - see [LICENSE](./LICENSE)

## 🙏 Acknowledgments

Built with:
- [gofakeit](https://github.com/brianvoe/gofakeit) - Fake data generation
- [chi](https://github.com/go-chi/chi) - Lightweight Go router
- [Next.js](https://nextjs.org/) - React framework with SSR

---

**Made with ❤️ by [Devendra Pratap](https://github.com/0xdps)**
