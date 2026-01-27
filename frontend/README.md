# Mockly Frontend

Next.js 14 website with landing page, documentation, and interactive playground featuring server-side rendering.

## ✨ Features

- **Server-Side Rendering (SSR)** - Fast initial page loads with fresh data
- **Automatic API Detection** - Seamlessly works in dev and production
- **Interactive Documentation** - Live API info and schema exploration with search guide
- **Enhanced Playground** - Modern 50/50 split-screen layout:
  - **Tabbed Navigation**: Quick, Pagination, Filters, Middleware, Path
  - **Request Builder** (left): Endpoint type selector, all controls organized in tabs
  - **Response Viewer** (right): Live JSON response with syntax highlighting
  - **Pagination controls**: Page, limit, offset with visual feedback
  - **Sorting controls**: Field selector with asc/desc order
  - **Search**: Full-text search with field targeting
  - **Field selection**: Multi-select with preview
  - **Middleware**: 11+ options (delay, flaky, cache bypass, etc.)
  - **Headers & Metrics**: Request timing, cache status, request ID
  - **5 utility tools**: Echo, Status, Delay, Middleware, Chaos
- **Modern UI** - Beautiful, responsive design with Tailwind CSS
- **Comprehensive Sitemaps** - 130+ URLs indexed (XML + HTML)
- **Mobile Friendly** - Fully responsive across all devices
- **Type Safe** - Full TypeScript support with auto-generated types

## 🚀 Quick Start

### Development

```bash
pnpm install
pnpm run dev
```

Website runs on **http://localhost:3000**

### Production Build

```bash
pnpm run build
pnpm start
```

### Type Checking

```bash
pnpm run type-check
```

## 📁 Project Structure

```
frontend/
├── app/                    # Next.js pages (App Router)
│   ├── page.tsx           # Home page (SSR)
│   ├── docs/              # Documentation (SSR)
│   └── playground/        # Interactive tester (CSR)
├── components/            # React components
│   ├── Header.tsx
│   ├── Footer.tsx
│   ├── ResourceCard.tsx
│   └── CodeExample.tsx
├── lib/
│   └── api.ts            # API configuration
├── types/
│   └── api.ts            # TypeScript definitions
└── shared/
    └── schemas/          # JSON schemas
```

## 🔧 Configuration

### Environment Variables

Environment variables are **optional**. The app automatically detects the API URL:

| Environment | API URL                    |
| ----------- | -------------------------- |
| Production  | `https://api.mockly.codes` |
| Development | `http://localhost:8080`    |

#### Override Default Behavior

Create `.env.local` in the frontend directory:

```bash
NEXT_PUBLIC_API_URL=https://api.mockly.codes
```

### API Integration

The `lib/api.ts` file automatically detects the environment:

```typescript
export const getApiUrl = () => {
  // Browser
  if (typeof window !== "undefined") {
    return window.location.hostname !== "localhost"
      ? "https://api.mockly.codes"
      : "http://localhost:8080";
  }

  // Server-side
  return (
    process.env.NEXT_PUBLIC_API_URL ||
    (process.env.NODE_ENV === "production"
      ? "https://api.mockly.codes"
      : "http://localhost:8080")
  );
};
```

## 🎯 Server-Side Rendering Strategy

### Pages Using SSR (Server Components)

| Page       | Rendering | Cache | Purpose                 |
| ---------- | --------- | ----- | ----------------------- |
| `/` (Home) | SSR + ISR | 5 min | Fresh resources list    |
| `/docs`    | SSR + ISR | 5 min | Live API info + schemas |

**Benefits:**

- ⚡ Faster initial page load
- 🔍 Better SEO (search engine indexable)
- 📊 Fresh data on each visit
- 🌐 Works without JavaScript

**Example:**

```typescript
// app/page.tsx
async function getResources() {
  const res = await fetch("https://api.mockly.codes/", {
    next: { revalidate: 300 }, // ISR: 5 minutes
  });
  return res.json();
}

export default async function Home() {
  const resources = await getResources();
  // Server-rendered with fresh data
}
```

### Pages Using CSR (Client Components)

| Page          | Rendering   | Purpose                 |
| ------------- | ----------- | ----------------------- |
| `/playground` | Client-Side | Interactive API testing |

**Benefits:**

- 🎮 Full interactivity
- 📡 Dynamic data fetching
- ⚙️ User-driven actions

**Example:**

```typescript
"use client"; // Client component

const [data, setData] = useState(null);

const fetchData = async () => {
  const res = await fetch("https://api.mockly.codes/users");
  setData(await res.json());
};
```

## 📊 Pages Overview

### Home Page (`/`)

- Server-side rendered with ISR
- Fetches available resources from API
- Displays feature cards and quick examples
- Falls back to static resources if API unavailable

### Documentation (`/docs`)

- Server-side rendered with ISR
- Shows live API status, version, and resource count
- Displays all resource schemas with properties
- Interactive "Try It" buttons for each endpoint

### Playground (`/playground`)

- Client-side rendered for full interactivity
- **50/50 split layout**: Request builder (left) and response viewer (right)
- **Tabbed navigation**: Quick, Pagination, Filters, Middleware, Path tabs
- **Compact header**: Reduced size for more screen space
- **Full-width layout**: No max-width constraints for wider content area
- Test any endpoint with custom parameters
- Live JSON response viewer with syntax highlighting
- Real-time URL preview with path options (direct/group)
- Code examples dropdown with cURL, JavaScript, and Python

## 🛠️ Technology Stack

- **Framework:** Next.js 16 (App Router)
- **UI Library:** React 19 (Server Components + Client Components)
- **Language:** TypeScript 5.6+
- **Styling:** Tailwind CSS
- **HTTP Client:** pingpong-fetch (Universal HTTP client)
- **Rendering:** SSR + ISR (5-minute revalidation)
- **Deployment:** Vercel
- **API:** https://api.mockly.codes
- **Type Generation:** Automated from JSON schemas

## 📦 Type Safety

TypeScript types are auto-generated from JSON schemas in `types/api.ts`:

```typescript
// Resource interfaces
export interface User {
  id: number;
  username: string;
  email: string;
  // ... more fields
}

export interface Post {
  id: number;
  user_id: number;
  title: string;
  // ... more fields
}

// API response types
export type CollectionResponse<T extends ResourceName> = ResourceMap[T][];
export type SingleResponse<T extends ResourceName> = ResourceMap[T];

// Root API response
export interface ApiRootResponse {
  message: string;
  version: string;
  resources: string[];
  docs: string;
}
```

## 🚀 Deployment

### Deploy to Vercel (Monorepo Setup)

This project is a monorepo with backend and frontend. The frontend is deployed separately to Vercel.

**Important:** Since the Next.js app is in the `frontend/` subdirectory, you must configure Vercel's Root Directory.

#### Vercel Dashboard Configuration:

1. Go to your project Settings → **General**
2. Under **Root Directory**, click **Edit**
3. Set to: `frontend`
4. Under **Build & Development Settings**:
   - **Build Command**: Auto-detect (Next.js) or `pnpm run build`
   - **Install Command**: Auto-detect or `pnpm install`
   - **Output Directory**: `.next` (auto-detected)
5. Redeploy

#### Or use Vercel CLI:

```bash
# From project root
vercel --prod

# When prompted, set Root Directory to: frontend
```

#### Why Root Directory = frontend?

- Your git repo stays at the project root (✅ backend/, frontend/, shared/)
- Vercel treats `frontend/` as the deployment root
- This is standard for monorepos - one repo, multiple deployable apps
- Backend deploys separately (fly.io), frontend deploys to Vercel

### Environment Variables (Vercel)

Set in Vercel dashboard or CLI:

```bash
vercel env add NEXT_PUBLIC_API_URL
# Enter: https://api.mockly.codes
```

**Note:** Environment variables are optional. The frontend defaults to `https://api.mockly.codes` in production.

## 🎨 Customization

### Adding New Pages

Create a new file in `app/`:

```typescript
// app/about/page.tsx
export default function AboutPage() {
  return <div>About Mockly</div>
}
```

### Creating Components

Add components in `components/`:

```typescript
// components/MyComponent.tsx
export function MyComponent() {
  return <div>My Component</div>
}
```

### Using API in Components

We use **pingpong-fetch**, a universal HTTP client that works in both browser and Node.js environments.

#### Server Components (SSR)

```typescript
import { createApiClient } from '@/lib/api'

async function getData() {
  const client = createApiClient()
  const res = await client.get('/users')
  return res.json()
}

export default async function Page() {
  const data = await getData()
  return <div>{JSON.stringify(data)}</div>
}
```

#### Client Components

```typescript
'use client'
import { apiClient } from '@/lib/api'

const MyComponent = () => {
  const [data, setData] = useState(null)

  useEffect(() => {
    const fetchData = async () => {
      const res = await apiClient.get('/users')
      setData(res.json())
    }
    fetchData()
  }, [])

  return <div>{JSON.stringify(data)}</div>
}
```

#### Benefits of pingpong-fetch

- ✅ **Universal** - Works in Node.js (SSR) and browsers
- ✅ **Fast** - Uses undici in Node.js, native fetch in browsers
- ✅ **Type-Safe** - Full TypeScript support
- ✅ **Retry Logic** - Automatic retry with exponential backoff
- ✅ **Chainable** - `.json()`, `.text()`, `.ok()`, `.isError()` methods
- ✅ **Small** - ~5-6KB gzipped

#### API Client Configuration

The `apiClient` is pre-configured with:

- Base URL (auto-detected: production or localhost)
- 30-second timeout
- Automatic JSON content-type
- Retry logic (2 retries with exponential backoff)

See [lib/api.ts](lib/api.ts) for configuration.

## 📚 Related Documentation

- **Project Overview:** [../README.md](../README.md)
- **API Documentation:** https://mockly.codes/docs
- **Backend Setup:** [../backend/README.md](../backend/README.md)
- **Contributing:** [../CONTRIBUTING.md](../CONTRIBUTING.md)

## 🐛 Troubleshooting

### Build Errors

```bash
# Clear cache and rebuild
rm -rf .next
pnpm run build
```

### Type Errors

```bash
# Check for type errors
pnpm run type-check
```

### API Connection Issues

Check that:

1. Backend is running on http://localhost:8080 (dev)
2. Production API is accessible: https://api.mockly.codes
3. CORS is enabled on the API

## 📝 Scripts

| Command               | Description              |
| --------------------- | ------------------------ |
| `pnpm run dev`        | Start development server |
| `pnpm run build`      | Build for production     |
| `pnpm start`          | Start production server  |
| `pnpm run lint`       | Run ESLint               |
| `pnpm run type-check` | TypeScript type checking |

## 🤝 Contributing

1. Create a new branch
2. Make your changes
3. Run `pnpm run build`
4. Run `pnpm run type-check`
5. Submit a pull request

See [../CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

## 📄 License

MIT License - see [../LICENSE](../LICENSE)
