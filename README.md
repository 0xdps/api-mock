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

### API Features
- 🚀 **Schema-Driven** - Add new endpoints by creating JSON schemas (zero code!)
- 🎯 **100+ Resources** - From users to weather, stocks to movies
- 📂 **14 Categories** - Resources organized into logical groups (people, commerce, content, etc.)
- 💡 **Realistic Data** - Powered by gofakeit with 50+ generators
- ⚡ **Blazing Fast** - Redis + in-memory cache (<1ms responses)
- 🔄 **Consistent Data** - Fixed seed ensures reproducible results
- 🆓 **Free Forever** - No rate limits, open source and self-hostable

### Advanced Query Features
- 📄 **Pagination** - page, limit, offset parameters
- 🔀 **Sorting** - Sort by any field, ascending or descending
- 🔍 **Search** - Full-text search with field targeting
- 🎯 **Filtering** - Filter by field values (exact, range, contains)
- ✂️ **Field Selection** - Return only specific fields

### Middleware & Testing
- ⏱️ **Delay Simulation** - Test timeouts (max 30s)
- 🎲 **Chaos Engineering** - Random failures with flakyRate
- 🚫 **Cache Control** - Bypass cache for fresh data
- 🏢 **Multi-tenancy** - Tenant isolation via X-Tenant-ID
- 🔐 **RBAC** - Role-based testing via X-Role
- 🔄 **Idempotency** - Safe retries with Idempotency-Key
- 🔍 **Request Tracing** - Track requests with X-Request-ID

### Frontend
- 🎨 **Modern UI** - Next.js 16 + React 19 with SSR
- 🎮 **Interactive Playground** - Test all API features
- 📚 **Comprehensive Docs** - Live examples and code snippets
- 📊 **Cache Management** - Admin endpoints for stats
- 🤖 **Automated Build** - TypeScript types auto-generated
- 🌐 **CORS Enabled** - Ready for frontend development

## 📦 Project Structure

```
mockly/
├── backend/          # Go API service (chi + gofakeit)
│   ├── cmd/          # Application entry points
│   ├── internal/     # Handlers, middleware, schema
│   ├── Dockerfile    # Docker build config
│   └── fly.toml      # Fly.io deployment config
├── frontend/         # Next.js 16 website (React 19 + TypeScript + Tailwind)
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

**Basic Usage:**
```bash
# Get resources with pagination
curl 'http://localhost:8080/products?page=1&limit=20'

# Get single item
curl http://localhost:8080/products/42

# Get resource metadata
curl http://localhost:8080/products/meta
```

**Advanced Queries:**
```bash
# Search products
curl 'http://localhost:8080/products?q=laptop&search_fields=name,description'

# Sort products by price
curl 'http://localhost:8080/products?sort=price&order=asc&limit=10'

# Filter by category and price
curl 'http://localhost:8080/products?category=Electronics&price<1000'

# Select specific fields
curl 'http://localhost:8080/products?fields=id,name,price&limit=50'
```

**Testing with Middleware:**
```bash
# Add 2-second delay
curl 'http://localhost:8080/products?delay=2000'

# Test with 50% failure rate
curl 'http://localhost:8080/products?flakyRate=0.5'

# Bypass cache for fresh data
curl 'http://localhost:8080/products?skip_cache=true'

# Multi-tenant request
curl -H "X-Tenant-ID: tenant-123" http://localhost:8080/products
```

**Production API:**
```bash
curl 'https://api.mockly.codes/products?page=1&limit=20'
```

📚 **Full API Documentation:** [backend/API_DOCUMENTATION.md](./backend/API_DOCUMENTATION.md)

## 📚 Available Resources (100 Endpoints!)

### 📂 Resource Groups

Resources are organized into 14 logical categories:

| Group | Count | Resources | Description |
|-------|-------|-----------|-------------|
| 🛒 **Commerce** | 14 | products, orders, payments, coupons, categories, tags, carts, wishlists, promotions, discounts, returns, refunds, shipping, inventory | E-commerce |
| 💼 **Business** | 12 | companies, organizations, jobs, meetings, invoices, subscriptions, clients, contracts, proposals, departments, vendors, reports | Business entities |
| ✈️ **Travel** | 10 | hotels, flights, restaurants, properties, cars, tours, attractions, bookings, destinations, travelguides | Travel & hospitality |
| 👥 **People** | 10 | users, contacts, students, players, employees, customers, profiles, authors, instructors, mentors | User profiles and people |
| 🎬 **Media** | 10 | movies, books, albums, videos, images, songs, playlists, photos, audios, streams | Entertainment media |
| 💬 **Social** | 9 | comments, reviews, messages, notifications, testimonials, likes, shares, followers, mentions | Social interactions |
| 📝 **Content** | 8 | articles, posts, news, podcasts, blogs, tutorials, guides, documents | Written content |
| ✅ **Productivity** | 6 | todos, notes, projects, tasks, tickets, events | Task management |
| 🌍 **Location** | 6 | countries, cities, weather, states, regions, coordinates | Geographic data |
| 💰 **Finance** | 6 | currencies, stocks, crypto, transactions, accounts, budgets | Financial data |
| 📚 **Reference** | 3 | faqs, quotes, languages | Reference data |
| 🍔 **Food** | 3 | recipes, ingredients, dishes | Food & cooking |
| ⚽ **Sports** | 2 | matches, teams | Sports data |
| 🎓 **Education** | 1 | courses | Educational content |

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
      "path": "/todos",
      "aliases": ["/tasks", "/todo-items"],
      "methods": ["GET", "POST"]
    }
  }
}
```

This creates:
- ✅ `GET /todos`
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

## ⚡ Cache & Performance

The backend uses an **in-memory cache with warmup on startup** for blazing fast responses:

```bash
🔥 Starting cache warmup (seed: 42, items per resource: 100)...
✅ Cache warmup completed in 2.5s
   📊 Resources: 100, Total items: 10,000
```

### Benefits

- **<1ms response times** - Pre-generated data served from memory
- **Consistent data** - Same ID always returns same object (perfect for testing)
- **Reproducible** - Fixed seed (42) ensures same data across restarts
- **Low CPU usage** - Generate once, serve thousands of times

### Configuration

Configure via environment variables:

```bash
# Number of items to cache per resource (default: 100)
CACHE_ITEMS_PER_RESOURCE=100

# Random seed for reproducible data (default: 42)
CACHE_SEED=42
```

### Admin Endpoints

```bash
# View cache statistics
GET /admin/cache/stats

# Refresh cache (regenerate all data)
POST /admin/cache/refresh
```

### Cache Bypass

Get fresh data without using cache (useful for testing/demos):

```bash
# Via query parameter
GET /users?nocache=true
GET /products?fresh=true

# Via header
curl -H "X-No-Cache: true" http://localhost:8080/users
curl -H "Cache-Control: no-cache" http://localhost:8080/products
```

All responses include `X-Cache` header:
- `X-Cache: HIT` - Served from cache
- `X-Cache: BYPASS` - Cache bypassed
- `X-Cache: MISS` - Cache miss

📖 **Full cache documentation:** [backend/CACHE.md](./backend/CACHE.md)  
📖 **Cache bypass guide:** [backend/CACHE_BYPASS.md](./backend/CACHE_BYPASS.md)

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
- **Framework:** Next.js 16 (App Router with SSR)
- **Language:** TypeScript
- **UI Library:** React 19
- **Styling:** Tailwind CSS
- **HTTP Client:** pingpong-fetch (universal, fast, type-safe)
- **Rendering:** Server-Side with ISR (5-minute cache)
- **Deployment:** Vercel
- **API Integration:** https://api.mockly.codes

## 📖 Documentation

### API Documentation

**Comprehensive API Reference:** [backend/API_DOCUMENTATION.md](./backend/API_DOCUMENTATION.md)
- 📋 Complete endpoint reference
- 🔍 All query parameters (pagination, sorting, search, filtering)
- 🛠️ All 11 middleware features documented
- 💻 Code examples (JavaScript, Python, cURL)
- 📊 Response formats and error handling
- ✨ Best practices and advanced usage

**Online Documentation:** Visit https://mockly.codes/docs for:
- Quick start guide (4 languages)
- Interactive schema explorer
- Live "Try It" buttons
- Real-time API status

### Interactive Playground
Visit https://mockly.codes/playground to:
- 🎮 Test all API features
- 📄 Pagination, sorting, search controls
- 🛠️ Middleware testing (delay, flaky, cache)
- 🔍 Request/response inspection
- 📋 Copy code examples
- 5 utility tools: Echo, Status, Delay, Middleware, Chaos

## 🤝 Contributing

Contributions welcome! The easiest way to contribute is to add new resource schemas:

1. Create `shared/schemas/your-resource.json`
2. Test locally with `npm run dev` (types auto-generate!)
3. Submit a PR

See [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

## 📚 Documentation

- **[backend/API_DOCUMENTATION.md](./backend/API_DOCUMENTATION.md)** - Comprehensive API reference
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Deployment guide
- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Contribution guidelines
- **[CHANGELOG.md](./CHANGELOG.md)** - Version history
- **[backend/README.md](./backend/README.md)** - Backend setup
- **[frontend/README.md](./frontend/README.md)** - Frontend setup

## � License

MIT License - see [LICENSE](./LICENSE)

## 🙏 Acknowledgments

Built with:
- [gofakeit](https://github.com/brianvoe/gofakeit) - Fake data generation
- [chi](https://github.com/go-chi/chi) - Lightweight Go router
- [Next.js 16](https://nextjs.org/) - React framework with SSR
- [React 19](https://react.dev/) - UI library

---

**Made with ❤️ by [Devendra Pratap](https://github.com/0xdps)**
