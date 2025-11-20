# Mockly Frontend

Next.js 14 website with landing page, documentation, and interactive playground featuring server-side rendering.

## ✨ Features

- **Server-Side Rendering (SSR)** - Fast initial page loads with fresh data
- **Automatic API Detection** - Seamlessly works in dev and production
- **Interactive Documentation** - Live API info and schema exploration
- **API Playground** - Test endpoints with real-time responses
- **Modern UI** - Beautiful, responsive design with Tailwind CSS
- **Mobile Friendly** - Fully responsive across all devices
- **Type Safe** - Full TypeScript support with auto-generated types

## 🚀 Quick Start

### Development

```bash
npm install
npm run dev
```

Website runs on **http://localhost:3000**

### Production Build

```bash
npm run build
npm start
```

### Type Checking

```bash
npm run type-check
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

| Environment | API URL |
|-------------|---------|
| Production | `https://api.mockly.codes` |
| Development | `http://localhost:8080` |

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
  if (typeof window !== 'undefined') {
    return window.location.hostname !== 'localhost'
      ? 'https://api.mockly.codes'
      : 'http://localhost:8080'
  }
  
  // Server-side
  return process.env.NEXT_PUBLIC_API_URL || 
    (process.env.NODE_ENV === 'production' 
      ? 'https://api.mockly.codes' 
      : 'http://localhost:8080')
}
```

## 🎯 Server-Side Rendering Strategy

### Pages Using SSR (Server Components)

| Page | Rendering | Cache | Purpose |
|------|-----------|-------|---------|
| `/` (Home) | SSR + ISR | 5 min | Fresh resources list |
| `/docs` | SSR + ISR | 5 min | Live API info + schemas |

**Benefits:**
- ⚡ Faster initial page load
- 🔍 Better SEO (search engine indexable)
- 📊 Fresh data on each visit
- 🌐 Works without JavaScript

**Example:**
```typescript
// app/page.tsx
async function getResources() {
  const res = await fetch('https://api.mockly.codes/', {
    next: { revalidate: 300 } // ISR: 5 minutes
  })
  return res.json()
}

export default async function Home() {
  const resources = await getResources()
  // Server-rendered with fresh data
}
```

### Pages Using CSR (Client Components)

| Page | Rendering | Purpose |
|------|-----------|---------|
| `/playground` | Client-Side | Interactive API testing |

**Benefits:**
- 🎮 Full interactivity
- 📡 Dynamic data fetching
- ⚙️ User-driven actions

**Example:**
```typescript
'use client' // Client component

const [data, setData] = useState(null)

const fetchData = async () => {
  const res = await fetch('https://api.mockly.codes/users')
  setData(await res.json())
}
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
- Test any endpoint with custom parameters
- Live JSON response viewer
- Code examples in cURL, JavaScript, and Python

## 🛠️ Technology Stack

- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **Rendering:** React Server Components + Client Components
- **Deployment:** Vercel
- **API:** https://api.mockly.codes

## 📦 Type Safety

TypeScript types are auto-generated from JSON schemas in `types/api.ts`:

```typescript
// Resource interfaces
export interface User {
  id: number
  username: string
  email: string
  // ... more fields
}

export interface Post {
  id: number
  user_id: number
  title: string
  // ... more fields
}

// API response types
export type CollectionResponse<T extends ResourceName> = ResourceMap[T][]
export type SingleResponse<T extends ResourceName> = ResourceMap[T]

// Root API response
export interface ApiRootResponse {
  message: string
  version: string
  resources: string[]
  docs: string
}
```

## 🚀 Deployment

### Deploy to Vercel

```bash
# Install Vercel CLI
npm install -g vercel

# Deploy
cd frontend
vercel --prod
```

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

Server component:
```typescript
// Fetch on server
async function getData() {
  const res = await fetch('https://api.mockly.codes/users', {
    next: { revalidate: 300 }
  })
  return res.json()
}
```

Client component:
```typescript
'use client'
import { getApiUrl } from '@/lib/api'

const apiUrl = getApiUrl()
// Use apiUrl for fetching
```

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
npm run build
```

### Type Errors

```bash
# Check for type errors
npm run type-check
```

### API Connection Issues

Check that:
1. Backend is running on http://localhost:8080 (dev)
2. Production API is accessible: https://api.mockly.codes
3. CORS is enabled on the API

## 📝 Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm start` | Start production server |
| `npm run lint` | Run ESLint |
| `npm run type-check` | TypeScript type checking |

## 🤝 Contributing

1. Create a new branch
2. Make your changes
3. Run `npm run type-check`
4. Submit a pull request

See [../CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

## 📄 License

MIT License - see [../LICENSE](../LICENSE)
