<div align="center">
  <img src="assets/logo.svg" alt="Mockly Logo" width="120" height="120">
  <h1>Mockly</h1>
  <p><strong>Free mock API service for testing and development</strong></p>
  <p>A modern, schema-driven alternative to JSONPlaceholder</p>
</div>

---

🌐 **Live:** https://www.mockly.codes  
📚 **Docs:** https://www.mockly.codes/docs  
🎮 **Playground:** https://www.mockly.codes/playground  
🤖 **AI-readable:** https://www.mockly.codes/llms.txt

## ✨ Features

### API Features

- 🚀 **Schema-Driven** — Add new endpoints by creating a JSON schema. Zero code changes needed.
- 🎯 **100+ Resources** — From users to weather, stocks to movies
- 📂 **14 Categories** — People, Commerce, Business, Content, Social, Media, Travel, Location, Finance, Food, Education, Sports, Productivity, Reference
- 💡 **Realistic Data** — Powered by gofakeit with 50+ generators
- ⚡ **Fast** — Redis + in-memory cache (<1ms responses)
- 🔄 **Consistent** — Fixed seed ensures reproducible results across restarts
- 🆓 **Free Forever** — No rate limits, no account required for built-in resources

### Query Features

- 📄 **Pagination** — `page`, `limit`, `offset`
- 🔀 **Sorting** — `sort` + `order` on any field
- 🔍 **Full-text search** — `q` with optional `search_fields` targeting
- 🎯 **Filtering** — Exact, range (`>`, `<`, `>=`, `<=`), and substring (`~contains`, `~startsWith`, `~endsWith`) operators
- ✂️ **Field selection** — `fields=id,name,price`

### Middleware & Testing

- ⏱️ **Delay** — `?delay=2000` simulates slow responses (max 30s)
- 🎲 **Chaos** — `?flakyRate=0.3` randomly returns 503s
- 🚫 **Cache bypass** — `?skip_cache=true`
- 🏢 **Multi-tenancy** — `X-Tenant-ID` header
- 🔐 **RBAC** — `X-Role` header
- 🔄 **Idempotency** — `Idempotency-Key` header
- 🔍 **Tracing** — `X-Request-ID` header

### Frontend

- 🎨 **Next.js 16** + React 19, TypeScript, Tailwind CSS
- 🎮 **Interactive Playground** — Live request/response, 5 utility tools
- 📚 **Docs at `/docs`** — Complete API reference with code examples
- 🤖 **AI-friendly** — `/llms.txt` and `/llms-full.txt` for LLM tools
- 🌐 **Hosted on Cloudflare** via OpenNext

## 📦 Project Structure

```
api-mockly/
├── backend/              # Go API server (chi router, gofakeit)
│   ├── cmd/server/       # Entry point
│   ├── internal/         # Handlers, middleware, schema loader, cache
│   └── Dockerfile        # Multi-stage Docker build
├── frontend/             # Next.js 16 website
│   ├── app/              # Pages: /, /docs, /playground, /templates, /resources
│   ├── components/       # React components
│   ├── lib/              # API client, schemas manifest
│   └── wrangler.jsonc    # Cloudflare Workers / Pages config
├── shared/               # Source of truth
│   ├── schemas/          # JSON Schema definitions (one file per resource)
│   └── scripts/          # generate-types.js (auto-generates frontend types)
├── Dockerfile            # Root-level Docker build (used by Railway)
├── railway.toml          # Railway deployment config
└── Makefile              # Dev helpers
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Node.js 20+
- pnpm 9+

### 1. Clone & install

```bash
git clone https://github.com/0xdps/api-mockly.git
cd api-mockly
pnpm install
```

### 2. Start both services

```bash
pnpm run dev
```

This runs:
- ✅ Schema generation (`shared/schemas/` → `frontend/lib/schemas-manifest.ts`)
- ✅ Go API at http://localhost:8080
- ✅ Next.js site at http://localhost:3000

Or run individually:

```bash
# Backend only
make api-dev

# Frontend only
pnpm --filter frontend dev
```

### 3. Try it

```bash
# Built-in resource — no auth
curl "http://localhost:8080/products?limit=5&sort=price&order=asc"

# Search
curl "http://localhost:8080/users?q=alice&limit=5"

# Single item
curl "http://localhost:8080/products/42"

# Schema metadata
curl "http://localhost:8080/products/meta"
```

## 🏗️ Build & Deploy

### Backend → Railway

The backend is a Dockerized Go binary deployed on [Railway](https://railway.com).

**Local build:**

```bash
# From the repo root (Railway builds from here)
docker build -t api-mockly-backend .
docker run -p 8080:8080 api-mockly-backend
```

**Deploy to Railway:**

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and link
railway login
railway link          # Select your project

# Set environment variables (once)
railway variables set CACHE_MODE=all
railway variables set CACHE_ITEMS_PER_RESOURCE=100
railway variables set CACHE_SEED=42
railway variables set MAX_ITEMS_PER_RESOURCE=1000
# Redis variables (REDIS_HOST, REDIS_PORT, REDIS_PASSWORD) are injected
# automatically when you add a Redis plugin in the Railway dashboard.

# Deploy
railway up
```

Railway uses `railway.toml` at the repo root which points to `Dockerfile`.
On every push to `trunk`, Railway auto-deploys if GitHub integration is enabled.

**Required environment variables:**

| Variable | Default | Description |
|---|---|---|
| `PORT` | set by Railway | HTTP listen port |
| `REDIS_HOST` | — | Auto-set by Redis plugin |
| `REDIS_PORT` | — | Auto-set by Redis plugin |
| `REDIS_PASSWORD` | — | Auto-set by Redis plugin |
| `REDIS_DB` | `0` | Redis database index |
| `REDIS_TLS_ENABLED` | `false` | Set `true` for external Redis |
| `CACHE_MODE` | `off` | `off` / `local` / `remote` / `all` |
| `CACHE_ITEMS_PER_RESOURCE` | `100` | Items to pre-generate per resource |
| `CACHE_SEED` | `42` | Seed for reproducible data |
| `MAX_ITEMS_PER_RESOURCE` | `1000` | Max items requestable per resource |

### Frontend → Cloudflare

The frontend is a Next.js 16 app deployed to **Cloudflare Workers** using [OpenNext](https://github.com/opennextjs/opennextjs-cloudflare).

**Local build:**

```bash
cd frontend
pnpm build                    # Compiles Next.js
pnpm exec wrangler pages dev  # Preview locally via Wrangler
```

**Deploy to Cloudflare Pages:**

```bash
# Install Wrangler (if not already)
pnpm add -g wrangler

# Login
wrangler login

# Build + deploy (run from frontend/)
cd frontend
pnpm build
pnpm exec wrangler pages deploy .open-next/assets --project-name api-mockly
```

Or using the Workers deploy path (configured in `wrangler.jsonc`):

```bash
cd frontend
pnpm build
pnpm exec wrangler deploy
```

`wrangler.jsonc` is already configured:
- **entry:** `.open-next/worker.js`
- **assets:** `.open-next/assets`
- **compatibility:** `nodejs_compat` flag, date `2026-05-02`

**CI/CD:** Connect the repo to Cloudflare Pages in the dashboard for automatic deployments on push. Set the build command to:

```
cd frontend && pnpm build
```

And the output directory to:

```
frontend/.open-next/assets
```

## 📚 Available Resources (100 Endpoints)

Resources are organized into 14 categories:

| Group | Count | Resources |
|---|---|---|
| 🛒 Commerce | 14 | products, orders, payments, coupons, categories, tags, carts, wishlists, promotions, discounts, returns, refunds, shipping, inventory |
| 💼 Business | 12 | companies, organizations, jobs, meetings, invoices, subscriptions, clients, contracts, proposals, departments, vendors, reports |
| ✈️ Travel | 10 | hotels, flights, restaurants, properties, cars, tours, attractions, bookings, destinations, travelguides |
| 👥 People | 10 | users, contacts, students, players, employees, customers, profiles, authors, instructors, mentors |
| 🎬 Media | 10 | movies, books, albums, videos, images, songs, playlists, photos, audios, streams |
| 💬 Social | 9 | comments, reviews, messages, notifications, testimonials, likes, shares, followers, mentions |
| 📝 Content | 8 | articles, posts, news, podcasts, blogs, tutorials, guides, documents |
| ✅ Productivity | 6 | todos, notes, projects, tasks, tickets, events |
| 🌍 Location | 6 | countries, cities, weather, states, regions, coordinates |
| 💰 Finance | 6 | currencies, stocks, crypto, transactions, accounts, budgets |
| 📚 Reference | 3 | faqs, quotes, languages |
| 🍔 Food | 3 | recipes, ingredients, dishes |
| ⚽ Sports | 2 | matches, teams |
| 🎓 Education | 1 | courses |

All accessible at `GET https://api.mockly.codes/{resource}`.

## 🎯 Adding a New Resource

Create a JSON schema in `shared/schemas/{group}/your-resource.json`:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Widget",
  "type": "object",
  "x-resource": {
    "name": "widgets",
    "singular": "widget",
    "description": "Example widget resource",
    "group": "commerce"
  },
  "properties": {
    "id": {
      "type": "integer",
      "x-generator": "random_int",
      "x-generator-params": { "min": 1, "max": 10000 }
    },
    "name": { "type": "string", "x-generator": "word" },
    "price": { "type": "number", "x-generator": "price" },
    "in_stock": { "type": "boolean", "x-generator": "bool" }
  }
}
```

Then restart the API — `GET /widgets` is live. Types are auto-generated on the next `pnpm build`.

See the full generator list in [`backend/internal/schema/loader.go`](./backend/internal/schema/loader.go).

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| **API language** | Go 1.23 |
| **API router** | chi v5 |
| **Data generation** | gofakeit/v7 |
| **Database** | MesaHub (SQLite-compatible, for user templates) |
| **Cache** | Redis + in-memory |
| **API deployment** | Railway (Docker) |
| **Frontend framework** | Next.js 16 / React 19 |
| **Frontend language** | TypeScript |
| **Styling** | Tailwind CSS |
| **Frontend adapter** | OpenNext for Cloudflare |
| **Frontend deployment** | Cloudflare Workers / Pages |

## 📖 Documentation

| File | Contents |
|---|---|
| [mockly.codes/docs](https://www.mockly.codes/docs) | Human-readable full API reference |
| [mockly.codes/llms.txt](https://www.mockly.codes/llms.txt) | Concise AI/LLM-readable reference |
| [mockly.codes/llms-full.txt](https://www.mockly.codes/llms-full.txt) | Full API reference for AI code generation |
| [backend/API_DOCUMENTATION.md](./backend/API_DOCUMENTATION.md) | Raw backend API reference |
| [DEPLOYMENT.md](./DEPLOYMENT.md) | Detailed deployment guide |
| [CONTRIBUTING.md](./CONTRIBUTING.md) | Contribution guidelines |
| [CHANGELOG.md](./CHANGELOG.md) | Version history |

## 🤝 Contributing

Contributions welcome! The easiest contribution is a new schema — no Go knowledge needed.

1. Fork the repo
2. Add `shared/schemas/{group}/your-resource.json`
3. Run `pnpm run dev` and test `http://localhost:8080/your-resource`
4. Open a PR

See [CONTRIBUTING.md](./CONTRIBUTING.md) for full guidelines.

## 📄 License

MIT — see [LICENSE](./LICENSE)

## 🙏 Acknowledgments

- [gofakeit](https://github.com/brianvoe/gofakeit) — fake data generation
- [chi](https://github.com/go-chi/chi) — Go HTTP router
- [Next.js](https://nextjs.org/) — React framework
- [OpenNext for Cloudflare](https://github.com/opennextjs/opennextjs-cloudflare) — Cloudflare adapter

---

**Made with ❤️ by [Devendra Pratap](https://github.com/0xdps)**


