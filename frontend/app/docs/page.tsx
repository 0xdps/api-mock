import type { Metadata } from 'next'
import Link from 'next/link'
import {
  Zap,
  Lock,
  Globe,
  ArrowRight,
  Terminal,
  Code2,
  BookOpen,
  Layers,
  Settings,
  Search,
  Filter,
  Database,
  AlertTriangle,
  Clock,
  ShieldCheck,
  Copy,
  ExternalLink,
  ChevronRight,
} from 'lucide-react'
import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'

export const metadata: Metadata = {
  title: 'API Documentation | Mockly',
  description:
    'Complete API reference for Mockly — 100+ mock data resources, query parameters, filtering, sorting, middleware, and user templates. Everything developers need to integrate.',
  openGraph: {
    title: 'API Documentation | Mockly',
    description:
      'Complete reference for the Mockly mock API — endpoints, query parameters, filtering, and more.',
    url: 'https://www.mockly.codes/docs',
    siteName: 'Mockly',
    type: 'website',
  },
}

// ── Shared primitives ─────────────────────────────────────────────────────────

function Section({
  id,
  title,
  icon: Icon,
  children,
}: {
  id: string
  title: string
  icon: React.ElementType
  children: React.ReactNode
}) {
  return (
    <section id={id} className="scroll-mt-20">
      <div className="flex items-center gap-3 mb-5">
        <div className="w-8 h-8 rounded-lg bg-primary-500/10 border border-primary-500/20 flex items-center justify-center shrink-0">
          <Icon className="w-4 h-4 text-primary-400" />
        </div>
        <h2 className="text-xl font-bold text-white">{title}</h2>
      </div>
      {children}
    </section>
  )
}

function CodeBlock({ code, lang = 'bash' }: { code: string; lang?: string }) {
  return (
    <pre className="bg-slate-950 border border-slate-700/60 rounded-xl p-4 overflow-x-auto text-xs text-slate-300 font-mono leading-relaxed">
      <code>{code.trim()}</code>
    </pre>
  )
}

function Table({ head, rows }: { head: string[]; rows: string[][] }) {
  return (
    <div className="overflow-x-auto rounded-xl border border-slate-700">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-slate-700 bg-slate-900/60">
            {head.map(h => (
              <th key={h} className="text-left px-4 py-3 text-xs font-semibold text-slate-400 uppercase tracking-wider whitespace-nowrap">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-700/50">
          {rows.map((row, i) => (
            <tr key={i} className="hover:bg-slate-800/30 transition-colors">
              {row.map((cell, j) => (
                <td key={j} className="px-4 py-3 text-slate-300 text-xs align-top">
                  {j === 0 ? (
                    <code className="text-primary-400 font-mono bg-slate-800/60 px-1.5 py-0.5 rounded text-xs">{cell}</code>
                  ) : (
                    cell
                  )}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Note({ children, variant = 'info' }: { children: React.ReactNode; variant?: 'info' | 'warning' | 'tip' }) {
  const styles = {
    info: 'bg-blue-500/10 border-blue-500/30 text-blue-300',
    warning: 'bg-amber-500/10 border-amber-500/30 text-amber-300',
    tip: 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300',
  }
  return (
    <div className={`border rounded-xl p-4 text-sm ${styles[variant]}`}>
      {children}
    </div>
  )
}

// ── Sidebar TOC ───────────────────────────────────────────────────────────────

const TOC = [
  { id: 'overview', label: 'Overview' },
  { id: 'quick-start', label: 'Quick Start' },
  { id: 'endpoints', label: 'Endpoints' },
  { id: 'auth', label: 'Authentication' },
  { id: 'query-params', label: 'Query Parameters' },
  { id: 'filtering', label: 'Filtering' },
  { id: 'search', label: 'Search' },
  { id: 'middleware', label: 'Middleware & Headers' },
  { id: 'response-format', label: 'Response Format' },
  { id: 'code-examples', label: 'Code Examples' },
  { id: 'resources', label: 'All Resources' },
  { id: 'limits', label: 'Limits & CORS' },
  { id: 'ai', label: 'Use with AI' },
]

// ── Page ──────────────────────────────────────────────────────────────────────

export default function DocsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col">
      <Header />

      <div className="flex-1 container mx-auto px-4 py-10 max-w-7xl">
        <div className="flex gap-10">
          {/* ── Sidebar ── */}
          <aside className="hidden lg:block w-52 shrink-0">
            <div className="sticky top-24 space-y-0.5">
              <p className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3 px-2">On this page</p>
              {TOC.map(item => (
                <a
                  key={item.id}
                  href={`#${item.id}`}
                  className="block px-2 py-1.5 text-sm text-slate-400 hover:text-white rounded transition-colors"
                >
                  {item.label}
                </a>
              ))}
              <div className="pt-4 border-t border-slate-700/50 mt-4 space-y-2">
                <a
                  href="/llms.txt"
                  className="flex items-center gap-1.5 text-xs text-primary-400 hover:text-primary-300 transition-colors px-2"
                >
                  <Code2 className="w-3 h-3" /> llms.txt
                </a>
                <a
                  href="/llms-full.txt"
                  className="flex items-center gap-1.5 text-xs text-primary-400 hover:text-primary-300 transition-colors px-2"
                >
                  <BookOpen className="w-3 h-3" /> llms-full.txt
                </a>
              </div>
            </div>
          </aside>

          {/* ── Content ── */}
          <main className="flex-1 min-w-0 space-y-14">

            {/* Header */}
            <div>
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary-500/10 border border-primary-500/20 text-primary-400 text-xs font-medium mb-4">
                <BookOpen className="w-3.5 h-3.5" /> API Reference
              </div>
              <h1 className="text-4xl font-bold text-white mb-3">Documentation</h1>
              <p className="text-slate-400 text-lg max-w-2xl">
                Everything you need to use Mockly — built-in resources, user templates, query parameters, and middleware.
              </p>
              <div className="flex gap-3 mt-5">
                <Link href="/playground" className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white text-sm rounded-lg transition-colors font-medium">
                  <Zap className="w-3.5 h-3.5" /> Try Playground
                </Link>
                <Link href="/templates" className="inline-flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 text-sm rounded-lg transition-colors border border-slate-700">
                  <Layers className="w-3.5 h-3.5" /> Explore Templates
                </Link>
              </div>
            </div>

            {/* ── Overview ── */}
            <Section id="overview" title="Overview" icon={Globe}>
              <div className="grid md:grid-cols-2 gap-4 mb-6">
                {[
                  { icon: Zap, title: 'No Setup', desc: '100+ endpoints ready instantly — no account, no config.' },
                  { icon: Globe, title: 'CORS Enabled', desc: 'Works directly from browsers and frontend apps.' },
                  { icon: Database, title: 'Realistic Data', desc: 'Powered by 50+ data generators for authentic-looking data.' },
                  { icon: ShieldCheck, title: 'Free Forever', desc: 'No rate limits, no plans, no credit card required.' },
                ].map(f => (
                  <div key={f.title} className="bg-slate-900/50 border border-white/5 rounded-xl p-4 flex gap-3">
                    <f.icon className="w-5 h-5 text-primary-400 shrink-0 mt-0.5" />
                    <div>
                      <p className="text-sm font-semibold text-white">{f.title}</p>
                      <p className="text-xs text-slate-500 mt-0.5">{f.desc}</p>
                    </div>
                  </div>
                ))}
              </div>

              <Note variant="tip">
                <strong>Base URL:</strong>{' '}
                <code className="font-mono text-emerald-300">https://api.mockly.codes</code>
                <br />
                All endpoints listed in this reference are relative to this base.
              </Note>
            </Section>

            {/* ── Quick Start ── */}
            <Section id="quick-start" title="Quick Start" icon={Zap}>
              <p className="text-slate-400 text-sm mb-4">No account needed. Copy, paste, run.</p>
              <CodeBlock code={`# Get 10 products
curl "https://api.mockly.codes/products?limit=10"

# Search users
curl "https://api.mockly.codes/users?q=john&limit=5"

# Sort products by price (descending)
curl "https://api.mockly.codes/products?sort=price&order=desc"

# Filter by field + paginate
curl "https://api.mockly.codes/products?category=Electronics&limit=20&page=2"

# Select specific fields only
curl "https://api.mockly.codes/products?fields=id,name,price&limit=50"`} />
            </Section>

            {/* ── Endpoints ── */}
            <Section id="endpoints" title="Endpoints" icon={Layers}>
              <Table
                head={['Pattern', 'Description']}
                rows={[
                  ['GET /{resource}', 'Paginated list of items — supports all query parameters'],
                  ['GET /{resource}/{id}', 'Single item by numeric ID (1–100)'],
                  ['GET /{resource}/meta', 'JSON Schema definition for the resource'],
                  ['GET /t/{templateId}', 'User template endpoint — requires API key header'],
                  ['GET /v1/{slug}', 'Featured template shorthand (same as /{resource})'],
                  ['GET /', 'API info: version, available resources, categories'],
                ]}
              />
              <p className="text-xs text-slate-500 mt-3">
                Replace <code className="text-slate-400 font-mono">&#123;resource&#125;</code> with any resource slug, e.g. <code className="text-slate-400 font-mono">products</code>, <code className="text-slate-400 font-mono">users</code>, <code className="text-slate-400 font-mono">orders</code>. See the{' '}
                <a href="#resources" className="text-primary-400 hover:underline">full resource list</a> below.
              </p>
            </Section>

            {/* ── Auth ── */}
            <Section id="auth" title="Authentication" icon={Lock}>
              <div className="grid md:grid-cols-2 gap-4 mb-5">
                <div className="bg-slate-900/50 border border-emerald-500/20 rounded-xl p-5">
                  <div className="flex items-center gap-2 mb-3">
                    <Globe className="w-4 h-4 text-emerald-400" />
                    <span className="text-sm font-semibold text-white">Built-in Resources</span>
                    <span className="text-xs text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded-full border border-emerald-500/20">No auth</span>
                  </div>
                  <p className="text-xs text-slate-400 mb-3">All 100+ built-in resources are completely public. No header required.</p>
                  <CodeBlock code={`curl "https://api.mockly.codes/products"`} />
                </div>
                <div className="bg-slate-900/50 border border-primary-500/20 rounded-xl p-5">
                  <div className="flex items-center gap-2 mb-3">
                    <Lock className="w-4 h-4 text-primary-400" />
                    <span className="text-sm font-semibold text-white">User Templates</span>
                    <span className="text-xs text-primary-400 bg-primary-500/10 px-2 py-0.5 rounded-full border border-primary-500/20">API key</span>
                  </div>
                  <p className="text-xs text-slate-400 mb-3">Custom templates require an API key in the Authorization header.</p>
                  <CodeBlock code={`curl -H "Authorization: Bearer mk_your_key" \\
  "https://api.mockly.codes/t/{templateId}"`} />
                </div>
              </div>
              <Note variant="info">
                Get an API key at{' '}
                <Link href="/dashboard/api-keys" className="underline">
                  mockly.codes/dashboard/api-keys
                </Link>{' '}
                after signing in. Keys starting with <code className="font-mono">mk_</code> are personal; <code className="font-mono">mak_</code> keys work for public templates.
              </Note>
            </Section>

            {/* ── Query Parameters ── */}
            <Section id="query-params" title="Query Parameters" icon={Settings}>
              <Table
                head={['Parameter', 'Type', 'Default', 'Description']}
                rows={[
                  ['limit', 'int 1–100', '10', 'Number of items to return per page'],
                  ['page', 'int', '1', 'Page number (1-based)'],
                  ['offset', 'int', '—', 'Manual offset; overrides page if set'],
                  ['sort', 'string', '—', 'Field name to sort results by'],
                  ['order', 'string', 'asc', 'Sort direction: asc or desc'],
                  ['q', 'string', '—', 'Full-text search query across all fields'],
                  ['search_fields', 'string', '—', 'Comma-separated fields to restrict search to'],
                  ['fields', 'string', '—', 'Comma-separated fields to include in response'],
                  ['{field}=value', 'string', '—', 'Filter by exact or operator-qualified value'],
                  ['delay', 'int (ms)', '—', 'Simulate network latency (max 30,000 ms)'],
                  ['flakyRate', 'float 0–1', '—', 'Probability of returning a 503 error'],
                  ['skip_cache', 'bool', 'false', 'Bypass the response cache for fresh data'],
                ]}
              />
              <div className="mt-4">
                <CodeBlock code={`# Combine multiple parameters
curl "https://api.mockly.codes/products?\\
  q=laptop&search_fields=name,description&\\
  sort=price&order=asc&\\
  page=1&limit=20&\\
  fields=id,name,price,rating"`} />
              </div>
            </Section>

            {/* ── Filtering ── */}
            <Section id="filtering" title="Filtering" icon={Filter}>
              <p className="text-slate-400 text-sm mb-4">
                Filter results by appending operators to field names: <code className="text-slate-300 font-mono">?&#123;field&#125;&#123;operator&#125;=&#123;value&#125;</code>
              </p>
              <Table
                head={['Operator', 'Meaning', 'Example']}
                rows={[
                  ['=', 'Equals (default)', '?category=Electronics'],
                  ['!=', 'Not equals', '?status!=inactive'],
                  ['>', 'Greater than', '?price>100'],
                  ['>=', 'Greater than or equal', '?rating>=4'],
                  ['<', 'Less than', '?price<500'],
                  ['<=', 'Less than or equal', '?stock<=10'],
                  ['~contains', 'Contains substring', '?name~contains=pro'],
                  ['~startsWith', 'Starts with', '?slug~startsWith=us'],
                  ['~endsWith', 'Ends with', '?email~endsWith=.com'],
                ]}
              />
              <div className="mt-4">
                <CodeBlock code={`# Price range + in-stock filter
curl "https://api.mockly.codes/products?price>=100&price<=500&in_stock=true"

# Name contains + sorted
curl "https://api.mockly.codes/products?name~contains=pro&sort=price&order=asc"`} />
              </div>
            </Section>

            {/* ── Search ── */}
            <Section id="search" title="Search" icon={Search}>
              <p className="text-slate-400 text-sm mb-4">
                Use <code className="text-slate-300 font-mono">?q=</code> for full-text search. By default all string fields are searched. Narrow it with <code className="text-slate-300 font-mono">search_fields</code>.
              </p>
              <CodeBlock code={`# Search all fields
curl "https://api.mockly.codes/products?q=laptop"

# Search specific fields only
curl "https://api.mockly.codes/products?q=laptop&search_fields=name,description"

# Combine with pagination and sort
curl "https://api.mockly.codes/products?q=phone&sort=price&order=asc&page=1&limit=10"`} />
              <div className="mt-4">
                <Note variant="info">
                  Search is <strong>case-insensitive</strong> and does <strong>partial matching</strong>.
                  Providing <code className="font-mono">search_fields</code> is recommended for better performance on large schemas.
                </Note>
              </div>
            </Section>

            {/* ── Middleware ── */}
            <Section id="middleware" title="Middleware & Headers" icon={Code2}>
              <p className="text-slate-400 text-sm mb-5">
                Mockly includes built-in middleware for testing resilience, multi-tenancy, and access control.
              </p>

              <div className="space-y-5">
                {/* Delay */}
                <div className="bg-slate-900/40 border border-slate-700/60 rounded-xl p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <Clock className="w-4 h-4 text-primary-400" />
                    <span className="font-semibold text-white text-sm">Delay Simulation</span>
                    <code className="text-xs font-mono text-slate-400 bg-slate-800 px-2 py-0.5 rounded">?delay=ms</code>
                  </div>
                  <p className="text-xs text-slate-500 mb-3">Simulate slow network responses. Max 30,000 ms.</p>
                  <CodeBlock code={`curl "https://api.mockly.codes/products?delay=2000"   # 2 second delay`} />
                </div>

                {/* Chaos */}
                <div className="bg-slate-900/40 border border-slate-700/60 rounded-xl p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <AlertTriangle className="w-4 h-4 text-amber-400" />
                    <span className="font-semibold text-white text-sm">Chaos Engineering</span>
                    <code className="text-xs font-mono text-slate-400 bg-slate-800 px-2 py-0.5 rounded">?flakyRate=0.0–1.0</code>
                  </div>
                  <p className="text-xs text-slate-500 mb-3">Randomly returns 503 errors at the given probability. Test your retry logic.</p>
                  <CodeBlock code={`curl "https://api.mockly.codes/products?flakyRate=0.3"  # 30% chance of failure`} />
                </div>

                {/* Request headers */}
                <div className="bg-slate-900/40 border border-slate-700/60 rounded-xl p-5">
                  <div className="flex items-center gap-2 mb-3">
                    <Settings className="w-4 h-4 text-primary-400" />
                    <span className="font-semibold text-white text-sm">Request Headers</span>
                  </div>
                  <Table
                    head={['Header', 'Description']}
                    rows={[
                      ['X-Tenant-ID', 'Isolate data by tenant — same tenant always gets same data'],
                      ['X-Role', 'Simulate RBAC: admin, user, guest, or any custom role'],
                      ['Idempotency-Key', 'Make POST/PUT/PATCH idempotent — response cached 24h'],
                      ['X-Request-ID', 'Custom trace ID echoed back in response headers'],
                    ]}
                  />
                  <div className="mt-3">
                    <CodeBlock code={`# Multi-tenant + role simulation
curl -H "X-Tenant-ID: tenant-a" -H "X-Role: admin" \\
  "https://api.mockly.codes/products?limit=5"

# Idempotent POST
curl -X POST -H "Idempotency-Key: order-123" \\
  -H "Content-Type: application/json" \\
  -d '{"product": "laptop"}' \\
  "https://api.mockly.codes/orders"`} />
                  </div>
                </div>
              </div>
            </Section>

            {/* ── Response Format ── */}
            <Section id="response-format" title="Response Format" icon={Database}>
              <div className="grid md:grid-cols-2 gap-4">
                <div>
                  <p className="text-sm font-semibold text-slate-300 mb-2">Collection (<code className="font-mono text-xs">GET /&#123;resource&#125;</code>)</p>
                  <CodeBlock code={`{
  "data": [
    { "id": 1, "name": "Laptop", "price": 999.99 },
    { "id": 2, "name": "Phone", "price": 499.99 }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}`} lang="json" />
                </div>
                <div>
                  <p className="text-sm font-semibold text-slate-300 mb-2">Single item (<code className="font-mono text-xs">GET /&#123;resource&#125;/42</code>)</p>
                  <CodeBlock code={`{
  "id": 42,
  "name": "Laptop",
  "price": 999.99,
  "category": "Electronics",
  "in_stock": true,
  "rating": 4.5
}`} lang="json" />
                </div>
              </div>
              <div className="mt-4">
                <p className="text-sm font-semibold text-slate-300 mb-2">Response Headers</p>
                <Table
                  head={['Header', 'Value', 'Description']}
                  rows={[
                    ['X-Cache', 'HIT / MISS / BYPASS', 'Cache status for the response'],
                    ['X-Total-Count', 'integer', 'Total number of items available (collection endpoints)'],
                    ['X-Request-ID', 'string', 'Request trace ID (echoed or generated)'],
                  ]}
                />
              </div>
            </Section>

            {/* ── Code Examples ── */}
            <Section id="code-examples" title="Code Examples" icon={Terminal}>
              <div className="space-y-6">
                <div>
                  <p className="text-sm font-semibold text-slate-300 mb-2">JavaScript / TypeScript</p>
                  <CodeBlock code={`// Fetch with pagination
const res = await fetch('https://api.mockly.codes/products?page=1&limit=20&sort=price&order=asc')
const { data, pagination } = await res.json()
console.log(\`Got \${data.length} of \${pagination.total} items\`)

// Search
const search = await fetch(\`https://api.mockly.codes/products?q=\${encodeURIComponent('laptop')}&limit=10\`)
const results = await search.json()

// User template (API key)
const tmpl = await fetch('https://api.mockly.codes/t/your-template-id', {
  headers: { Authorization: 'Bearer mk_your_key' }
})
const { data: items } = await tmpl.json()`} lang="javascript" />
                </div>

                <div>
                  <p className="text-sm font-semibold text-slate-300 mb-2">Python</p>
                  <CodeBlock code={`import requests

# Paginated + sorted
r = requests.get('https://api.mockly.codes/products', params={
    'page': 1, 'limit': 20, 'sort': 'price', 'order': 'asc'
})
data = r.json()
print(f"Total: {data['pagination']['total']}")

# Search specific fields
r = requests.get('https://api.mockly.codes/products', params={
    'q': 'laptop', 'search_fields': 'name,description', 'limit': 5
})

# User template
r = requests.get('https://api.mockly.codes/t/your-template-id',
    headers={'Authorization': 'Bearer mk_your_key'}
)`} lang="python" />
                </div>

                <div>
                  <p className="text-sm font-semibold text-slate-300 mb-2">Go</p>
                  <CodeBlock code={`package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    req, _ := http.NewRequest("GET",
        "https://api.mockly.codes/products?limit=10&sort=price&order=desc", nil)

    // For user templates, add auth header:
    // req.Header.Set("Authorization", "Bearer mk_your_key")

    resp, _ := http.DefaultClient.Do(req)
    defer resp.Body.Close()

    var result map[string]any
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Println(result)
}`} lang="go" />
                </div>
              </div>
            </Section>

            {/* ── All Resources ── */}
            <Section id="resources" title="All Resources" icon={Layers}>
              <p className="text-slate-400 text-sm mb-5">
                100 resources across 14 categories. All accessible at <code className="text-slate-300 font-mono">GET https://api.mockly.codes/&#123;resource&#125;</code>.
              </p>
              <div className="grid sm:grid-cols-2 lg:grid-cols-2 gap-3">
                {([
                  {
                    cat: 'commerce', emoji: '🛒',
                    items: ['products', 'orders', 'carts', 'payments', 'invoices', 'categories', 'wishlists', 'returns', 'refunds', 'discounts', 'shipping', 'promotions', 'inventory', 'coupons'],
                  },
                  {
                    cat: 'people', emoji: '👥',
                    items: ['users', 'profiles', 'contacts', 'employees', 'authors', 'customers', 'instructors', 'mentors', 'students', 'players'],
                  },
                  {
                    cat: 'business', emoji: '💼',
                    items: ['companies', 'organizations', 'departments', 'jobs', 'contracts', 'meetings', 'proposals', 'reports', 'budgets', 'vendors', 'clients', 'subscriptions'],
                  },
                  {
                    cat: 'content', emoji: '📝',
                    items: ['articles', 'posts', 'blogs', 'news', 'tutorials', 'documents', 'guides', 'podcasts'],
                  },
                  {
                    cat: 'social', emoji: '💬',
                    items: ['comments', 'likes', 'shares', 'followers', 'notifications', 'messages', 'mentions', 'reviews', 'testimonials'],
                  },
                  {
                    cat: 'media', emoji: '🎬',
                    items: ['images', 'videos', 'audio', 'albums', 'playlists', 'movies', 'books', 'photos', 'songs', 'streams'],
                  },
                  {
                    cat: 'travel', emoji: '✈️',
                    items: ['flights', 'hotels', 'bookings', 'destinations', 'attractions', 'tours', 'restaurants', 'properties', 'cars', 'travelguides'],
                  },
                  {
                    cat: 'location', emoji: '🌍',
                    items: ['countries', 'cities', 'states', 'regions', 'coordinates', 'weather'],
                  },
                  {
                    cat: 'finance', emoji: '💰',
                    items: ['transactions', 'accounts', 'budgets', 'crypto', 'stocks', 'currencies'],
                  },
                  {
                    cat: 'food', emoji: '🍔',
                    items: ['recipes', 'ingredients', 'dishes'],
                  },
                  {
                    cat: 'education', emoji: '🎓',
                    items: ['courses'],
                  },
                  {
                    cat: 'sports', emoji: '⚽',
                    items: ['teams', 'matches'],
                  },
                  {
                    cat: 'productivity', emoji: '✅',
                    items: ['tasks', 'projects', 'notes', 'todos', 'events', 'tickets'],
                  },
                  {
                    cat: 'reference', emoji: '📚',
                    items: ['faqs', 'quotes', 'languages'],
                  },
                ] as const).map(({ cat, emoji, items }) => (
                  <div key={cat} className="bg-slate-900/40 border border-slate-700/60 rounded-xl p-4">
                    <p className="text-sm font-semibold text-white mb-2">
                      {emoji} {cat}
                    </p>
                    <div className="flex flex-wrap gap-1.5">
                      {items.map(r => (
                        <Link
                          key={r}
                          href={`/resources/${r}`}
                          className="text-xs font-mono text-slate-400 hover:text-primary-300 bg-slate-800/60 hover:bg-slate-800 border border-slate-700 hover:border-primary-500/30 px-2 py-0.5 rounded transition-all"
                        >
                          {r}
                        </Link>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
              <p className="text-xs text-slate-600 mt-4">
                Click any resource to open its interactive documentation.
              </p>
            </Section>

            {/* ── Limits ── */}
            <Section id="limits" title="Limits & CORS" icon={ShieldCheck}>
              <div className="grid md:grid-cols-2 gap-4">
                <Note variant="tip">
                  <strong className="block mb-1">No Rate Limits</strong>
                  Mockly is free with unlimited requests — no throttling, no quotas. Use it in CI, load tests, or high-frequency dev scenarios.
                </Note>
                <Note variant="tip">
                  <strong className="block mb-1">CORS Fully Open</strong>
                  All origins, methods, and headers are allowed. Works from any browser, Postman, or cURL without extra config.
                </Note>
              </div>
              <div className="mt-4 bg-slate-900/40 border border-slate-700/60 rounded-xl p-5">
                <p className="text-sm font-semibold text-white mb-3">What Mockly cannot do</p>
                <ul className="space-y-1.5 text-xs text-slate-400">
                  {[
                    'Data is not persisted — POST/PUT/PATCH/DELETE requests are accepted but state is never modified',
                    'Cannot JOIN or query across multiple resources in one request',
                    'No GraphQL, WebSockets, or streaming endpoints',
                    'Generated data is deterministic per cache key — bypass with ?skip_cache=true for different data',
                  ].map(item => (
                    <li key={item} className="flex gap-2">
                      <span className="text-slate-600 shrink-0 mt-0.5">—</span>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            </Section>

            {/* ── AI ── */}
            <Section id="ai" title="Use with AI" icon={Code2}>
              <p className="text-slate-400 text-sm mb-5">
                Mockly provides machine-readable documentation specifically for LLM-based tools (GitHub Copilot, ChatGPT, Claude, Cursor, etc.).
              </p>
              <div className="grid md:grid-cols-2 gap-4 mb-6">
                <div className="bg-slate-900/40 border border-primary-500/20 rounded-xl p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <Code2 className="w-4 h-4 text-primary-400" />
                    <span className="font-semibold text-white text-sm">llms.txt</span>
                  </div>
                  <p className="text-xs text-slate-500 mb-3">Concise structured reference — paste this URL into any AI tool to give it full context on Mockly.</p>
                  <a
                    href="/llms.txt"
                    target="_blank"
                    className="inline-flex items-center gap-1.5 text-xs text-primary-400 hover:text-primary-300 font-mono transition-colors"
                  >
                    mockly.codes/llms.txt <ExternalLink className="w-3 h-3" />
                  </a>
                </div>
                <div className="bg-slate-900/40 border border-slate-700/60 rounded-xl p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <BookOpen className="w-4 h-4 text-slate-400" />
                    <span className="font-semibold text-white text-sm">llms-full.txt</span>
                  </div>
                  <p className="text-xs text-slate-500 mb-3">Complete API reference including all examples — best for deep code generation tasks.</p>
                  <a
                    href="/llms-full.txt"
                    target="_blank"
                    className="inline-flex items-center gap-1.5 text-xs text-slate-400 hover:text-white font-mono transition-colors"
                  >
                    mockly.codes/llms-full.txt <ExternalLink className="w-3 h-3" />
                  </a>
                </div>
              </div>

              <p className="text-sm font-semibold text-white mb-3">Example AI prompts</p>
              <div className="space-y-3">
                {[
                  {
                    label: 'React hook',
                    prompt: 'Using the Mockly API (https://api.mockly.codes), write a React hook that fetches paginated products with search and sorting. Use https://www.mockly.codes/llms.txt for the full API reference.',
                  },
                  {
                    label: 'Python data loader',
                    prompt: 'Using the Mockly API (base URL: https://api.mockly.codes), write a Python function that loads all users across pages and returns a pandas DataFrame. No auth needed.',
                  },
                  {
                    label: 'Testing harness',
                    prompt: 'Write a Jest test suite that uses the Mockly API (https://api.mockly.codes/products) to verify pagination, sorting, and filtering work correctly in a frontend component.',
                  },
                ].map(({ label, prompt }) => (
                  <div key={label} className="bg-slate-950 border border-slate-700/60 rounded-xl p-4">
                    <p className="text-xs font-semibold text-slate-400 mb-2">{label}</p>
                    <p className="text-xs text-slate-300 leading-relaxed font-mono">{prompt}</p>
                  </div>
                ))}
              </div>
            </Section>

          </main>
        </div>
      </div>
      <Footer />
    </div>
  )
}
