# Mockly API - Comprehensive Documentation

**Base URL:** `https://api.mockly.codes`  
**Version:** 1.0.0  
**Last Updated:** December 2025

## Table of Contents

- [Introduction](#introduction)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Base Endpoints](#base-endpoints)
- [Resource Endpoints](#resource-endpoints)
- [Query Parameters](#query-parameters)
  - [Pagination](#pagination)
  - [Sorting](#sorting)
  - [Search](#search)
  - [Filtering](#filtering)
  - [Field Selection](#field-selection)
  - [Locale](#locale)
- [Request Headers](#request-headers)
- [Middleware Features](#middleware-features)
  - [Delay Simulation](#delay-simulation)
  - [Chaos Engineering](#chaos-engineering)
  - [Cache Control](#cache-control)
  - [Idempotency](#idempotency)
  - [Multi-tenancy](#multi-tenancy)
  - [Role-Based Access Control](#role-based-access-control)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Rate Limits](#rate-limits)
- [CORS](#cors)
- [Code Examples](#code-examples)
- [Available Resources](#available-resources)

---

## Introduction

Mockly is a free, schema-driven mock API service that provides realistic test data for over 100 resources across 14 categories. Perfect for:

- **Frontend Development** - Build UIs without backend dependencies
- **Testing** - Consistent, reproducible test data
- **Prototyping** - Quick demos with realistic data
- **Learning** - Practice API integration

### Key Features

✅ **No Authentication Required** - Start using immediately  
✅ **100+ Resources** - Products, users, orders, articles, and more  
✅ **Realistic Data** - 50+ data generators (gofakeit)  
✅ **Full REST Support** - GET, POST, PUT, PATCH, DELETE  
✅ **Advanced Querying** - Pagination, sorting, search, filtering  
✅ **CORS Enabled** - Works in browsers  
✅ **Free Forever** - No rate limits  

---

## Quick Start

### Get a list of products

```bash
curl https://api.mockly.codes/products?limit=5
```

### Get a single product

```bash
curl https://api.mockly.codes/products/42
```

### Search for products

```bash
curl "https://api.mockly.codes/products?q=laptop&limit=10"
```

### Get products with sorting

```bash
curl "https://api.mockly.codes/products?sort=price&order=asc&limit=10"
```

---

## Authentication

**No authentication required!** All endpoints are publicly accessible. This makes Mockly perfect for frontend development, testing, and learning.

---

## Base Endpoints

### Health Check

Get API information and available resources.

```
GET /
```

**Response:**
```json
{
  "message": "Mockly API",
  "version": "1.0.0",
  "docs": "https://www.mockly.codes/docs",
  "resources": ["products", "users", "orders", ...],
  "groups": {
    "commerce": ["products", "orders", "carts", ...],
    "people": ["users", "contacts", "employees", ...],
    ...
  }
}
```

### Static Assets

```
GET /favicon.svg    # SVG favicon
GET /favicon.ico    # ICO favicon
```

### Admin Endpoints

```
GET  /admin/cache/stats     # Get cache statistics
POST /admin/cache/refresh   # Refresh cache
```

---

## Resource Endpoints

All resources follow RESTful conventions with three main endpoints:

### 1. Collection Endpoint

Get a list of items for a resource.

```
GET /{resource}
```

**Supported Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `offset` - Alternative to page (default: calculated from page)
- `sort` - Field name to sort by
- `order` - Sort order (`asc` or `desc`, default: `asc`)
- `q` or `search` - Search query
- `search_fields` - Comma-separated fields to search in
- `fields` - Comma-separated fields to return
- `{field}={value}` - Filter by field value (see [Filtering](#filtering))
- `locale` - BCP-47 locale tag for locale-aware data generation (e.g. `en-IN`, `ja-JP`)
- `count` - Legacy parameter for item count (backwards compatibility)

**Example:**
```bash
curl "https://api.mockly.codes/products?page=1&limit=20&sort=price&order=asc"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Laptop",
      "price": 999.99,
      "category": "Electronics",
      "brand": "TechCorp",
      "description": "High-performance laptop"
    },
    ...
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### 2. Single Item Endpoint

Get a specific item by ID.

```
GET /{resource}/{id}
```

**Example:**
```bash
curl https://api.mockly.codes/products/42
```

**Response:**
```json
{
  "id": 42,
  "name": "Laptop",
  "price": 999.99,
  "category": "Electronics",
  "brand": "TechCorp",
  "description": "High-performance laptop",
  "in_stock": true,
  "rating": 4.5
}
```

### 3. Metadata Endpoint

Get JSON schema and metadata for a resource.

```
GET /{resource}/meta
```

**Example:**
```bash
curl https://api.mockly.codes/products/meta
```

**Response:**
```json
{
  "resource": "products",
  "schema": {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "Product",
    "type": "object",
    "properties": {
      "id": { "type": "integer" },
      "name": { "type": "string" },
      "price": { "type": "number" },
      ...
    }
  }
}
```

---

## Query Parameters

### Pagination

Control the number of items returned and navigate through pages.

**Parameters:**
- `page` (integer, default: 1) - Page number (1-based)
- `limit` (integer, default: 10, max: 100) - Items per page
- `offset` (integer) - Manual offset (overrides page calculation)

**Examples:**

```bash
# Get first page with 20 items
curl "https://api.mockly.codes/products?page=1&limit=20"

# Get second page
curl "https://api.mockly.codes/products?page=2&limit=20"

# Use offset directly
curl "https://api.mockly.codes/products?offset=40&limit=20"
```

**Response includes pagination metadata:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Sorting

Sort results by any field in ascending or descending order.

**Parameters:**
- `sort` (string) - Field name to sort by
- `order` (string, default: `asc`) - Sort order: `asc` or `desc`

**Examples:**

```bash
# Sort by price (ascending)
curl "https://api.mockly.codes/products?sort=price&order=asc"

# Sort by name (descending)
curl "https://api.mockly.codes/products?sort=name&order=desc"

# Combine with pagination
curl "https://api.mockly.codes/products?sort=created_at&order=desc&limit=10"
```

### Search

Full-text search across resource fields.

**Parameters:**
- `q` or `search` (string) - Search query
- `search_fields` (string) - Comma-separated list of fields to search in (optional)

**Search Behavior:**
- Case-insensitive
- Partial matching
- Searches all text fields by default
- Specify `search_fields` for targeted search

**Examples:**

```bash
# Search all fields for "laptop"
curl "https://api.mockly.codes/products?q=laptop"

# Search only in name and description
curl "https://api.mockly.codes/products?q=laptop&search_fields=name,description"

# Alternative parameter name
curl "https://api.mockly.codes/products?search=laptop"

# Combine with pagination and sorting
curl "https://api.mockly.codes/products?q=laptop&sort=price&order=asc&limit=20"
```

### Filtering

Filter results by exact field values.

**Format:** `?{field}={value}`

**Supported Operators:**
- `=` (default) - Equals
- `!=` - Not equals
- `>` - Greater than
- `<` - Less than
- `>=` - Greater than or equal
- `<=` - Less than or equal
- `~contains` - Contains substring
- `~startsWith` - Starts with
- `~endsWith` - Ends with

**Examples:**

```bash
# Filter by category
curl "https://api.mockly.codes/products?category=Electronics"

# Filter by price range
curl "https://api.mockly.codes/products?price>=100&price<=1000"

# Filter with contains
curl "https://api.mockly.codes/products?name~contains=laptop"

# Multiple filters
curl "https://api.mockly.codes/products?category=Electronics&in_stock=true&price<500"

# Combine with sorting
curl "https://api.mockly.codes/products?category=Electronics&sort=price&order=asc"
```

**Reserved Parameters** (not treated as filters):
- `page`, `limit`, `offset` - Pagination
- `sort`, `order` - Sorting
- `q`, `search`, `search_fields` - Search
- `fields` - Field selection
- `count`, `nocache`, `fresh` - Legacy/cache control
- `delay`, `flakyRate`, `skip_cache` - Middleware
- `locale` - Locale-aware data generation

### Field Selection

Return only specific fields to reduce response size.

**Parameter:**
- `fields` (string) - Comma-separated list of field names

**Examples:**

```bash
# Get only id, name, and price
curl "https://api.mockly.codes/products?fields=id,name,price"

# Combine with other parameters
curl "https://api.mockly.codes/products?fields=id,name&limit=50&sort=name"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Laptop",
      "price": 999.99
    },
    ...
  ],
  "pagination": {...}
}
```

### Locale

Generate locale-aware data for personal and location fields (names, cities, phone numbers, postal codes, etc.).

**Parameter:** `locale` (string, BCP-47 tag)
- Overrides `first_name`, `last_name`, `name`, `city`, `state`, `country`, `phone`, and `zip` generators
- Bypasses cache automatically — always returns fresh data
- Unsupported locale tags silently fall back to the default (English) generators

**Supported Locales:**

| Tag | Country |
|---|---|
| `en-IN` | India |
| `ja-JP` | Japan |
| `de-DE` | Germany |
| `fr-FR` | France |
| `zh-CN` | China |
| `pt-BR` | Brazil |
| `es-ES` | Spain |
| `ko-KR` | South Korea |
| `ar-SA` | Saudi Arabia |
| `en-GB` | United Kingdom |

**Examples:**

```bash
# Indian users
curl "https://api.mockly.codes/users?locale=en-IN&limit=5"

# Japanese employees
curl "https://api.mockly.codes/employees?locale=ja-JP&limit=10"

# German customers, sorted by name
curl "https://api.mockly.codes/customers?locale=de-DE&sort=first_name&order=asc"

# Works on any people resource
curl "https://api.mockly.codes/contacts?locale=pt-BR&limit=20"
```

**Example response for `?locale=en-IN`:**
```json
{
  "id": 4821,
  "first_name": "Priya",
  "last_name": "Sharma",
  "email": "priya.sharma@example.com",
  "phone": "+919876543210",
  "city": "Mumbai",
  "state": "Maharashtra",
  "country": "India"
}
```

---

## Request Headers

### Standard Headers

**X-Request-ID** (optional)
- Custom request identifier
- Useful for tracing and debugging
- Auto-generated if not provided
- Returned in response headers

```bash
curl -H "X-Request-ID: req-12345" https://api.mockly.codes/products
```

**X-Tenant-ID** (optional)
- Multi-tenancy support
- Isolate data by tenant
- For testing multi-tenant applications

```bash
curl -H "X-Tenant-ID: tenant-abc" https://api.mockly.codes/products
```

**X-Role** (optional)
- Role-based access simulation
- Test different permission levels
- Values: `admin`, `user`, `guest`, or custom

```bash
curl -H "X-Role: admin" https://api.mockly.codes/products
```

**Idempotency-Key** (optional)
- Ensures idempotent POST/PUT/PATCH operations
- Prevents duplicate requests
- Response cached for 24 hours

```bash
curl -X POST -H "Idempotency-Key: unique-key-123" \
  https://api.mockly.codes/products \
  -d '{"name": "New Product"}'
```

### Response Headers

**X-Request-ID**
- Echoes or provides the request ID

**X-Cache**
- Cache status: `HIT`, `MISS`, or `BYPASS`

**X-Total-Count** (collection endpoints)
- Total number of items available

---

## Middleware Features

### Delay Simulation

Simulate network latency or slow API responses.

**Parameter:** `delay` (integer, milliseconds)
- Adds artificial delay to response
- Max: 30,000ms (30 seconds)
- Useful for testing timeouts, loading states

**Examples:**

```bash
# 1 second delay
curl "https://api.mockly.codes/products?delay=1000"

# 5 second delay
curl "https://api.mockly.codes/products?delay=5000"

# Test timeout handling
curl --max-time 2 "https://api.mockly.codes/products?delay=3000"
```

### Chaos Engineering

Simulate random failures for resilience testing.

**Parameter:** `flakyRate` (float, 0.0 to 1.0)
- Probability of request failure
- 0.0 = never fails, 1.0 = always fails
- Returns 503 Service Unavailable on simulated failure

**Examples:**

```bash
# 10% chance of failure
curl "https://api.mockly.codes/products?flakyRate=0.1"

# 50% chance of failure
curl "https://api.mockly.codes/products?flakyRate=0.5"

# 100% failure (testing error handling)
curl "https://api.mockly.codes/products?flakyRate=1.0"
```

**Failure Response:**
```json
{
  "error": "Service Unavailable",
  "message": "Simulated failure from flakyRate parameter",
  "rate": 0.5
}
```

### Cache Control

Control caching behavior.

**Parameter:** `skip_cache=true` or `skip_cache=1`
- Bypass cache and generate fresh data
- Useful for testing with different data sets
- Response header: `X-Cache: BYPASS`

**Examples:**

```bash
# Normal request (uses cache)
curl "https://api.mockly.codes/products?limit=10"
# Response: X-Cache: HIT

# Bypass cache (fresh data)
curl "https://api.mockly.codes/products?limit=10&skip_cache=true"
# Response: X-Cache: BYPASS
```

### Idempotency

Ensure POST/PUT/PATCH requests are idempotent.

**Header:** `Idempotency-Key: {unique-key}`
- First request: Executes and caches response
- Subsequent requests: Returns cached response
- Cache TTL: 24 hours
- Only applies to POST, PUT, PATCH methods

**Example:**

```bash
# First request - executes
curl -X POST \
  -H "Idempotency-Key: order-12345" \
  -H "Content-Type: application/json" \
  -d '{"product": "laptop", "quantity": 1}' \
  https://api.mockly.codes/orders

# Second request - returns cached response
curl -X POST \
  -H "Idempotency-Key: order-12345" \
  -H "Content-Type: application/json" \
  -d '{"product": "laptop", "quantity": 1}' \
  https://api.mockly.codes/orders
```

### Multi-tenancy

Test multi-tenant applications.

**Header:** `X-Tenant-ID: {tenant-identifier}`
- Isolate data by tenant
- Each tenant gets consistent data
- Useful for testing tenant isolation

**Example:**

```bash
# Tenant A
curl -H "X-Tenant-ID: tenant-a" https://api.mockly.codes/products

# Tenant B (different data)
curl -H "X-Tenant-ID: tenant-b" https://api.mockly.codes/products
```

### Role-Based Access Control

Simulate different user roles.

**Header:** `X-Role: {role-name}`
- Test different permission levels
- Common roles: `admin`, `user`, `guest`
- Custom roles supported

**Example:**

```bash
# Admin role
curl -H "X-Role: admin" https://api.mockly.codes/products

# User role
curl -H "X-Role: user" https://api.mockly.codes/products

# Guest role
curl -H "X-Role: guest" https://api.mockly.codes/products
```

---

## Response Format

### Success Response (Collection)

```json
{
  "data": [
    {
      "id": 1,
      "name": "Item 1",
      ...
    },
    ...
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### Success Response (Single Item)

```json
{
  "id": 42,
  "name": "Item 42",
  "field": "value",
  ...
}
```

### Error Response

```json
{
  "error": "Error Type",
  "message": "Detailed error message"
}
```

### HTTP Status Codes

- **200 OK** - Successful GET request
- **201 Created** - Successful POST request
- **400 Bad Request** - Invalid parameters
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Simulated failure (flakyRate)

---

## Error Handling

### Common Errors

**Invalid Resource**
```bash
curl https://api.mockly.codes/invalid-resource
# 404 Not Found
{
  "error": "Resource not found"
}
```

**Invalid Parameter**
```bash
curl "https://api.mockly.codes/products?limit=invalid"
# Uses default value (10)
```

**Simulated Failure**
```bash
curl "https://api.mockly.codes/products?flakyRate=1.0"
# 503 Service Unavailable
{
  "error": "Service Unavailable",
  "message": "Simulated failure from flakyRate parameter",
  "rate": 1.0
}
```

---

## Rate Limits

**No rate limits!** Mockly is completely free to use with unlimited requests. Perfect for:
- High-volume testing
- Load testing
- Continuous integration
- Development environments

---

## CORS

Cross-Origin Resource Sharing (CORS) is enabled for all origins.

**Allowed:**
- All origins (`*`)
- All methods (GET, POST, PUT, PATCH, DELETE, OPTIONS)
- All headers

**Perfect for:**
- Frontend development
- Browser-based applications
- Local development

---

## Code Examples

### JavaScript / Fetch

```javascript
// Get products with pagination
fetch('https://api.mockly.codes/products?page=1&limit=20')
  .then(res => res.json())
  .then(data => {
    console.log(`Total: ${data.pagination.total} products`);
    console.log(data.data);
  });

// Search with error handling
async function searchProducts(query) {
  try {
    const response = await fetch(
      `https://api.mockly.codes/products?q=${encodeURIComponent(query)}&limit=10`
    );
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Search failed:', error);
    return [];
  }
}

// POST with idempotency
fetch('https://api.mockly.codes/orders', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Idempotency-Key': 'order-12345'
  },
  body: JSON.stringify({
    product: 'laptop',
    quantity: 1
  })
})
  .then(res => res.json())
  .then(data => console.log(data));
```

### Python / Requests

```python
import requests

# Get products
response = requests.get(
    'https://api.mockly.codes/products',
    params={
        'page': 1,
        'limit': 20,
        'sort': 'price',
        'order': 'asc'
    }
)
data = response.json()
print(f"Total: {data['pagination']['total']} products")

# Search products
def search_products(query, fields=None):
    params = {'q': query, 'limit': 10}
    if fields:
        params['search_fields'] = ','.join(fields)
    
    response = requests.get(
        'https://api.mockly.codes/products',
        params=params
    )
    return response.json()['data']

results = search_products('laptop', ['name', 'description'])

# Multi-tenant request
response = requests.get(
    'https://api.mockly.codes/products',
    headers={'X-Tenant-ID': 'tenant-123'}
)
```

### cURL

```bash
# Basic GET
curl https://api.mockly.codes/products

# With pagination and sorting
curl "https://api.mockly.codes/products?page=1&limit=20&sort=price&order=asc"

# Search
curl "https://api.mockly.codes/products?q=laptop&search_fields=name,description"

# Filter
curl "https://api.mockly.codes/products?category=Electronics&price<1000"

# With headers
curl -H "X-Tenant-ID: tenant-123" \
     -H "X-Role: admin" \
     https://api.mockly.codes/products

# POST with idempotency
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: unique-key-123" \
  -d '{"name": "New Product", "price": 99.99}' \
  https://api.mockly.codes/products

# Test with delay
curl "https://api.mockly.codes/products?delay=2000"

# Chaos testing
curl "https://api.mockly.codes/products?flakyRate=0.5"
```

---

## Available Resources

Mockly provides 100+ resources across 14 categories:

### Commerce (🛒)
- `products`, `orders`, `carts`, `payments`, `invoices`, `shipments`, `customers`, `reviews`, `coupons`, `categories`, `wishlists`, `returns`

### People (👥)
- `users`, `profiles`, `contacts`, `employees`, `customers`, `authors`, `students`, `teachers`

### Business (💼)
- `companies`, `organizations`, `departments`, `jobs`, `contracts`, `meetings`, `proposals`, `reports`, `budgets`, `vendors`, `clients`

### Content (📝)
- `articles`, `posts`, `pages`, `blogs`, `news`, `tutorials`, `documentation`, `faqs`, `announcements`

### Social (💬)
- `comments`, `likes`, `shares`, `followers`, `notifications`, `messages`, `mentions`, `tags`, `hashtags`

### Media (🎬)
- `images`, `videos`, `audio`, `albums`, `playlists`, `galleries`, `podcasts`, `streams`

### Travel (✈️)
- `flights`, `hotels`, `bookings`, `destinations`, `attractions`, `tours`, `transportation`, `itineraries`

### Location (🌍)
- `addresses`, `countries`, `cities`, `states`, `locations`, `places`, `landmarks`, `coordinates`

### Finance (💰)
- `transactions`, `accounts`, `cards`, `invoices`, `subscriptions`, `crypto`, `stocks`, `exchanges`

### Food (🍔)
- `recipes`, `ingredients`, `restaurants`, `menus`, `cuisines`, `nutrition`, `meals`

### Education (🎓)
- `courses`, `lessons`, `assignments`, `grades`, `schools`, `programs`, `certifications`, `exams`

### Sports (⚽)
- `teams`, `players`, `matches`, `scores`, `tournaments`, `leagues`, `stats`, `fixtures`

### Productivity (✅)
- `tasks`, `projects`, `notes`, `todos`, `events`, `reminders`, `calendars`, `workspaces`

### Reference (📚)
- `books`, `authors`, `quotes`, `definitions`, `translations`, `weather`, `currencies`, `timezones`

**Full list:** Visit https://www.mockly.codes/docs

---

## Advanced Usage

### Combining Features

Combine multiple features for complex testing scenarios:

```bash
# Pagination + Sorting + Search + Field Selection
curl "https://api.mockly.codes/products?\
q=laptop&\
search_fields=name,description&\
sort=price&\
order=asc&\
page=1&\
limit=20&\
fields=id,name,price,rating"

# Multi-tenant + RBAC + Delay + Cache Bypass
curl -H "X-Tenant-ID: tenant-a" \
     -H "X-Role: admin" \
     "https://api.mockly.codes/products?\
skip_cache=true&\
delay=1000&\
limit=10"

# Chaos Testing + Idempotency
curl -X POST \
  -H "Idempotency-Key: order-123" \
  -H "Content-Type: application/json" \
  -d '{"product": "laptop"}' \
  "https://api.mockly.codes/orders?flakyRate=0.3"
```

### Testing Scenarios

**Load Testing**
```bash
# Simulate slow network
for i in {1..100}; do
  curl "https://api.mockly.codes/products?delay=500&limit=50" &
done
wait
```

**Chaos Engineering**
```bash
# Random failures
curl "https://api.mockly.codes/products?flakyRate=0.5"

# Retry logic testing
for i in {1..10}; do
  curl "https://api.mockly.codes/products?flakyRate=0.7" || echo "Failed attempt $i"
done
```

**Cache Warming**
```bash
# Bypass cache to test generation
curl "https://api.mockly.codes/products?skip_cache=true&limit=100"
```

---

## Best Practices

### 1. Use Pagination
```bash
# ✅ Good
curl "https://api.mockly.codes/products?page=1&limit=50"

# ❌ Avoid (uses default limit=10)
curl "https://api.mockly.codes/products"
```

### 2. Specify Fields
```bash
# ✅ Good (smaller response)
curl "https://api.mockly.codes/products?fields=id,name,price"

# ❌ Avoid (large response with all fields)
curl "https://api.mockly.codes/products"
```

### 3. Use Search Wisely
```bash
# ✅ Good (targeted search)
curl "https://api.mockly.codes/products?q=laptop&search_fields=name"

# ⚠️ OK but slower (searches all fields)
curl "https://api.mockly.codes/products?q=laptop"
```

### 4. Combine Filters
```bash
# ✅ Good (specific results)
curl "https://api.mockly.codes/products?category=Electronics&price<500&in_stock=true"
```

### 5. Handle Errors
```javascript
// ✅ Good
try {
  const response = await fetch('https://api.mockly.codes/products');
  if (!response.ok) throw new Error(`HTTP ${response.status}`);
  const data = await response.json();
} catch (error) {
  console.error('Request failed:', error);
}
```

---

## Support & Resources

- **Documentation:** https://www.mockly.codes/docs
- **Playground:** https://www.mockly.codes/playground
- **GitHub:** https://github.com/0xdps/api-mockly
- **Website:** https://www.mockly.codes

---

## Changelog

### Version 1.0.0 (December 2025)
- ✅ 100+ resources across 14 categories
- ✅ Full pagination support
- ✅ Advanced search with field targeting
- ✅ Sorting and filtering
- ✅ Field selection
- ✅ Chaos engineering (delay, flaky)
- ✅ Idempotency support
- ✅ Multi-tenancy and RBAC
- ✅ Cache control
- ✅ Request tracing

---

## License

MIT License - Free to use for any purpose.

---

**Happy mocking! 🎭**

For more examples and interactive testing, visit the [Playground](https://www.mockly.codes/playground).
