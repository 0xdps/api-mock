# Deployment Guide

Mockly has two independently deployed services:

| Service | Platform | Config file |
|---|---|---|
| **Backend** (Go API) | Railway | `railway.toml`, `Dockerfile` |
| **Frontend** (Next.js) | Cloudflare Workers | `frontend/wrangler.jsonc` |

---

## Backend → Railway

The backend is a Dockerized Go binary. Railway builds it from the repo root using `Dockerfile` and `railway.toml`.

### Prerequisites

```bash
npm install -g @railway/cli
railway login
```

### First-time setup

1. **Create a project on [railway.com](https://railway.com)**
2. **Add a Redis plugin** (Dashboard → New → Database → Redis). Railway auto-injects `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`.
3. **Link your local clone:**

```bash
railway link    # select your project
```

4. **Add a Redis database** in the Railway dashboard: New > Database > Redis. This injects `${{Redis.REDIS_URL}}` automatically.

5. **Set environment variables:**

```bash
# Reference the Redis plugin (Railway resolves ${{...}} in variable values)
railway variables set 'REDIS_URL=${{Redis.REDIS_URL}}'
railway variables set CACHE_MODE=all
railway variables set CACHE_ITEMS_PER_RESOURCE=100
railway variables set CACHE_SEED=42
railway variables set MAX_ITEMS_PER_RESOURCE=1000
# Required for user/template features:
railway variables set MESAHUB_URL=<your-mesahub-url>
railway variables set NUBE_GATEWAY_URL=https://api.nubeauth.com
```

6. **Deploy:**

```bash
railway up
```

### Auto-deploy on push

Connect your GitHub repo in the Railway dashboard (Settings → Source). Railway will auto-deploy on every push to `trunk`.

### Local Docker build (same as Railway)

```bash
# From repo root — mirrors what Railway builds
docker build -t api-mockly-backend .
docker run -p 8080:8080 \
  -e CACHE_MODE=local \
  -e CACHE_ITEMS_PER_RESOURCE=100 \
  -e CACHE_SEED=42 \
  api-mockly-backend
```

### Environment variables reference

| Variable | Default | Description |
|---|---|---|
| `PORT` | set by Railway | HTTP listen port |
| `REDIS_HOST` | — | Auto-set by Redis plugin |
| `REDIS_PORT` | — | Auto-set by Redis plugin |
| `REDIS_PASSWORD` | — | Auto-set by Redis plugin |
| `REDIS_DB` | `0` | Redis database index |
| `REDIS_TLS_ENABLED` | `false` | Enable TLS for external Redis |
| `CACHE_MODE` | `off` | `off` / `local` / `remote` / `all` |
| `CACHE_ITEMS_PER_RESOURCE` | `100` | Items pre-generated per resource |
| `CACHE_SEED` | `42` | Seed for reproducible data |
| `MAX_ITEMS_PER_RESOURCE` | `1000` | Max items per request |
| `MESAHUB_URL` | — | MesaHub DB URL (for user templates) |
| `MESAHUB_TOKEN` | — | MesaHub auth token |
| `NUBE_JWKS_URL` | — | JWKS endpoint for auth token validation |

### Verify

```bash
railway logs               # tail logs
railway domain             # get your public URL

curl https://your-app.up.railway.app/
# → {"api":"mockly","version":"...","resources":100}
```

---

## Frontend → Cloudflare

The frontend is a Next.js 16 app deployed to **Cloudflare Workers** using [OpenNext for Cloudflare](https://github.com/opennextjs/opennextjs-cloudflare).

### Prerequisites

```bash
pnpm add -g wrangler
wrangler login
```

### Build

```bash
cd frontend
pnpm build
```

This runs (in order):
1. `generate-types` — generates `lib/schemas-manifest.ts` from `shared/schemas/`
2. `next build` — compiles the Next.js app
3. `wrangler build` (via OpenNext) — outputs `.open-next/`

The output is:
- `.open-next/worker.js` — Cloudflare Worker entry point
- `.open-next/assets/` — static assets

### Deploy to Cloudflare Pages

```bash
cd frontend
pnpm build
wrangler pages deploy .open-next/assets --project-name api-mockly
```

Or via the Workers route (configured in `wrangler.jsonc`):

```bash
cd frontend
pnpm build
wrangler deploy
```

### Preview locally

```bash
cd frontend
pnpm build
wrangler pages dev .open-next/assets
```

### CI/CD via Cloudflare Pages dashboard

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com) → **Workers & Pages** → **Create** → **Pages** → **Connect to Git**
2. Select the `api-mockly` repo
3. Set:
   - **Build command:** `cd frontend && pnpm build`
   - **Build output directory:** `frontend/.open-next/assets`
   - **Root directory:** `/` (repo root)
4. Add environment variables if needed (see below)
5. Every push to `trunk` triggers a deployment automatically.

### wrangler.jsonc overview

```jsonc
{
  "name": "api-mockly",
  "compatibility_date": "2026-05-02",
  "compatibility_flags": ["nodejs_compat"],
  "main": ".open-next/worker.js",
  "assets": {
    "directory": ".open-next/assets",
    "binding": "ASSETS"
  }
}
```

### Environment variables (frontend)

Set these in the Cloudflare Pages dashboard under Settings → Environment Variables:

| Variable | Description |
|---|---|
| `NEXT_PUBLIC_API_URL` | Public API base URL (defaults to `https://api.mockly.codes`) |

---

## Production URLs

| Service | URL |
|---|---|
| API | https://api.mockly.codes |
| Website | https://www.mockly.codes |
| Docs | https://www.mockly.codes/docs |
| Playground | https://www.mockly.codes/playground |
| LLMs.txt | https://www.mockly.codes/llms.txt |

---

## Troubleshooting

**Backend not starting**
- Check `railway logs` for Go panic output
- Ensure `CACHE_MODE=local` if Redis is not configured yet

**Frontend 500 on Cloudflare**
- Run `wrangler pages dev` locally to reproduce — errors surface in terminal
- Ensure `nodejs_compat` flag is set in `wrangler.jsonc`
- Dynamic routes (e.g. `/templates/[id]`) require the Worker entry, not just static assets

**`wrangler deploy` fails with "No D1 binding"**
- Mockly does not use D1 — ignore this if it appears in old configs; `wrangler.jsonc` is the authoritative config

**Schema changes not reflected after deploy**
- Re-run `pnpm build` from `frontend/` — schema generation is a prebuild step
- Schemas are embedded into the binary at Docker build time for the backend
