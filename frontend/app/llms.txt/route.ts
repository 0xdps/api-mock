import { NextResponse } from 'next/server'

const CONTENT = `# Mockly — Free Mock API Service

> Schema-driven mock REST API providing realistic fake data for 100+ resources across 14 categories. No auth required for public endpoints. Free forever, no rate limits, CORS enabled.

## Links

- Website: https://www.mockly.codes
- Playground: https://www.mockly.codes/playground
- Docs: https://www.mockly.codes/docs
- Explore Templates: https://www.mockly.codes/templates
- Full API Reference: https://www.mockly.codes/llms-full.txt

## API Base URL

https://api.mockly.codes

## Authentication

Two tiers:

1. **Public built-in resources** — No authentication required
   \`GET https://api.mockly.codes/{resource}\`

2. **User-created templates** — API key required
   \`GET https://api.mockly.codes/t/{templateId}\`
   Header: \`Authorization: Bearer <api-key>\`

API keys start with \`mk_\` (personal) or \`mak_\` (public template access).
Get API keys at: https://www.mockly.codes/dashboard/api-keys

## Quick Start

\`\`\`bash
# List products (no auth required)
curl "https://api.mockly.codes/products?limit=10"

# Search users
curl "https://api.mockly.codes/users?q=john&limit=5"

# Filter with operators
curl "https://api.mockly.codes/products?category=Electronics&price<500"

# Sort descending
curl "https://api.mockly.codes/products?sort=price&order=desc"

# Select fields only
curl "https://api.mockly.codes/products?fields=id,name,price&limit=20"

# Locale-aware data (Indian users)
curl "https://api.mockly.codes/users?locale=en-IN&limit=5"

# User template (API key required)
curl -H "Authorization: Bearer mk_your_key" "https://api.mockly.codes/t/template-id"
\`\`\`

## Endpoint Patterns

- \`GET /{resource}\` — paginated list of items
- \`GET /{resource}/{id}\` — single item by numeric ID
- \`GET /{resource}/meta\` — JSON Schema definition for the resource
- \`GET /t/{templateId}\` — user template endpoint (requires API key)
- \`GET /v1/{slug}\` — featured template shorthand (same as /{resource})

## Query Parameters

| Parameter       | Type        | Default | Description                            |
|-----------------|-------------|---------|----------------------------------------|
| \`limit\`         | int 1–100   | 10      | Items per page                         |
| \`page\`          | int         | 1       | Page number (1-based)                  |
| \`offset\`        | int         | —       | Manual offset, overrides page          |
| \`sort\`          | string      | —       | Field name to sort by                  |
| \`order\`         | string      | asc     | Sort direction: \`asc\` or \`desc\`        |
| \`q\`             | string      | —       | Full-text search across all fields     |
| \`search_fields\` | string      | —       | Comma-separated fields to search in    |
| \`fields\`        | string      | —       | Comma-separated fields to return       |
| \`{field}=val\`   | string      | —       | Filter by field value                  |
| \`locale\`        | string      | —       | BCP-47 locale tag for locale-aware names, cities, phones (e.g. \`en-IN\`, \`ja-JP\`, \`de-DE\`) |
| \`delay\`         | int (ms)    | —       | Simulate network latency               |
| \`flakyRate\`     | float 0–1   | —       | Probability of 503 error (chaos)       |
| \`skip_cache\`    | bool        | false   | Bypass response cache                  |

## Filtering Operators

Append operator to field name: \`?{field}{op}={value}\`

- \`=\` — equals (default, e.g. \`?category=Electronics\`)
- \`!=\` — not equals
- \`>\`, \`>=\`, \`<\`, \`<=\` — numeric comparisons
- \`~contains\` — substring match (e.g. \`?name~contains=laptop\`)
- \`~startsWith\` — prefix match
- \`~endsWith\` — suffix match

Examples:
\`\`\`
?price>=100&price<=500&category=Electronics
?name~contains=pro&in_stock=true
?sort=price&order=asc&page=2&limit=20
\`\`\`

## Request Headers

| Header            | Description                                      |
|-------------------|--------------------------------------------------|
| \`X-Tenant-ID\`     | Isolate data by tenant (multi-tenant testing)    |
| \`X-Role\`          | Simulate RBAC: admin, user, guest, or custom     |
| \`Idempotency-Key\` | Make POST/PUT/PATCH idempotent (24h cache)       |
| \`X-Request-ID\`    | Custom trace ID (echoed in response headers)     |

## Response Format

Collection endpoints:
\`\`\`json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
\`\`\`

Single item (\`GET /{resource}/{id}\`): returns the object directly, no wrapper.

Response headers: \`X-Cache\` (HIT/MISS/BYPASS), \`X-Total-Count\`, \`X-Request-ID\`

## Available Resources (100+)

**commerce** — products, orders, carts, payments, invoices, categories, wishlists, returns, refunds, discounts, shipping, promotions, inventory, coupons
**people** — users, profiles, contacts, employees, authors, customers, instructors, mentors, students, players
**business** — companies, organizations, departments, jobs, contracts, meetings, proposals, reports, budgets, vendors, clients, subscriptions
**content** — articles, posts, blogs, news, tutorials, documents, guides, podcasts
**social** — comments, likes, shares, followers, notifications, messages, mentions, reviews, testimonials
**media** — images, videos, audio, albums, playlists, movies, books, photos, songs, streams
**travel** — flights, hotels, bookings, destinations, attractions, tours, restaurants, properties, cars, travelguides
**location** — countries, cities, states, regions, coordinates, weather
**finance** — transactions, accounts, budgets, crypto, stocks, currencies
**food** — recipes, ingredients, dishes
**education** — courses
**sports** — teams, matches
**productivity** — tasks, projects, notes, todos, events, tickets
**reference** — faqs, quotes, languages

## What Mockly Cannot Do

- Data is not persisted — POST/PUT/PATCH/DELETE accept requests but don't modify stored state
- No querying across multiple resources in one call
- No GraphQL or WebSocket endpoints
- No real authentication on the API itself (dashboard auth is separate)
- Generated data is deterministic per cache key, not truly random every call

## Code Examples

JavaScript:
\`\`\`javascript
const res = await fetch('https://api.mockly.codes/products?limit=10&sort=price&order=asc')
const { data, pagination } = await res.json()
console.log(\`Got \${data.length} of \${pagination.total} products\`)
\`\`\`

Python:
\`\`\`python
import requests
r = requests.get('https://api.mockly.codes/products', params={'q': 'laptop', 'limit': 5})
print(r.json()['data'])
\`\`\`

Go:
\`\`\`go
resp, _ := http.Get("https://api.mockly.codes/products?limit=10")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
\`\`\`
`

export async function GET() {
  return new NextResponse(CONTENT, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      'Cache-Control': 'public, max-age=86400, s-maxage=86400',
    },
  })
}
