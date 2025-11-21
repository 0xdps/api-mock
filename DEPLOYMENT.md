# Deployment Guide

## 🚀 Backend Deployment (Fly.io)

### Prerequisites

1. **Install Fly CLI:**
```bash
curl -L https://fly.io/install.sh | sh
```

2. **Login to Fly.io:**
```bash
flyctl auth login
```

### Important: Schema Syncing for Deployment

The backend deployment **automatically includes all schemas** from `shared/schemas/`. Here's how it works:

```
Repository Root
├── shared/schemas/          ← Source schemas
│   └── *.json
└── backend/
    ├── Dockerfile           ← Copies from shared/schemas/
    ├── fly.toml            ← Builds from root directory
    └── internal/schema/
        └── embedded/        ← Schemas copied here during build
```

### Deployment Process

#### Method 1: Using Makefile (Recommended)

```bash
# From backend directory
cd backend
make deploy
```

This automatically:
- ✅ Syncs schemas before deployment
- ✅ Builds from repository root
- ✅ Deploys to Fly.io

#### Method 2: Using Fly CLI Directly

```bash
# From repository root
flyctl deploy --config backend/fly.toml
```

**Important:** Always deploy from the **repository root**, not from the `backend/` directory.

#### Method 3: First-Time Setup

If deploying for the first time:

```bash
# From repository root
cd backend
flyctl launch --no-deploy

# Edit fly.toml if needed
# Then deploy
cd ..
flyctl deploy --config backend/fly.toml
```

### Docker Build Process

The Dockerfile is designed to:

1. **Build from repository root** to access `shared/schemas/`
2. **Copy schemas** into `backend/internal/schema/embedded/`
3. **Build Go binary** with embedded schemas
4. **Create minimal Alpine image** (~20MB)

### How Fly.io Finds Schemas

The `fly.toml` configuration specifies:

```toml
[build]
  dockerfile = "backend/Dockerfile"
  ignorefile = ".dockerignore"
```

The Dockerfile then:

```dockerfile
# Copy shared schemas to backend
COPY shared/schemas/*.json internal/schema/embedded/
```

### Testing Docker Build Locally

```bash
# From repository root
docker build -f backend/Dockerfile -t api-mockly:latest .

# Run locally
docker run -p 8080:8080 api-mockly:latest

# Test
curl http://localhost:8080/
```

Or use the Makefile:

```bash
cd backend
make docker-build    # Builds image
make docker-run      # Runs container
```

### Verifying Schema Inclusion

After deployment, verify all schemas are loaded:

```bash
# Check deployed API
curl https://api-mockly.fly.dev/ | jq '.resources | length'

# Should return: 54
```

### Deployment Checklist

Before deploying:

- ✅ All schemas exist in `shared/schemas/`
- ✅ Run `make sync-schemas` to sync to backend
- ✅ Test locally with `go run cmd/server/main.go`
- ✅ Verify all endpoints work
- ✅ Deploy from **repository root**

### Common Issues

#### Issue: Schemas not found during build

**Problem:** Building from `backend/` directory

```bash
# ❌ Wrong - can't access ../shared/
cd backend
docker build -t api-mockly .
```

**Solution:** Build from repository root

```bash
# ✅ Correct
cd /path/to/api-mock  # Repository root
docker build -f backend/Dockerfile -t api-mockly .
```

#### Issue: Some schemas missing

**Problem:** Schemas not synced before deployment

**Solution:** Sync schemas first

```bash
cd backend
make sync-schemas
make deploy
```

#### Issue: Fly.io deploy fails

**Problem:** Incorrect build context

**Solution:** Check `fly.toml` has correct dockerfile path:

```toml
[build]
  dockerfile = "backend/Dockerfile"
```

And deploy from root:

```bash
flyctl deploy --config backend/fly.toml
```

### Environment Variables

No environment variables needed for schemas - they're embedded at build time!

Optional variables:
```bash
# Set port (default: 8080)
flyctl secrets set PORT=8080

# Set any custom config
flyctl secrets set MY_VAR=value
```

### Updating Schemas in Production

When you add or modify schemas:

```bash
# 1. Update schema in shared/schemas/
vim shared/schemas/new-resource.json

# 2. Redeploy (schemas auto-sync during build)
cd backend
make deploy

# 3. Verify
curl https://api-mockly.fly.dev/ | jq '.resources'
```

---

## 🌐 Frontend Deployment (Vercel)

### Prerequisites

1. **Install Vercel CLI:**
```bash
npm install -g vercel
```

2. **Login to Vercel:**
```bash
vercel login
```

### Important: Monorepo Configuration

This project is a **monorepo** with both backend and frontend. The frontend Next.js app is in the `frontend/` subdirectory.

**Critical:** You MUST configure Vercel's **Root Directory** setting to `frontend` for the deployment to work correctly.

### Deployment Process

#### Method 1: Using Vercel Dashboard (Recommended for First Deploy)

1. Connect your GitHub repo to Vercel
2. **IMPORTANT:** In Settings → General:
   - **Root Directory:** Set to `frontend` ⚠️
   - **Framework Preset:** Next.js (auto-detected)
   - **Build Command:** `npm run build` (auto-detected)
   - **Install Command:** `npm install` (auto-detected)
   - **Output Directory:** `.next` (auto-detected)

3. Every push to `main`/`trunk` auto-deploys!

**Why Root Directory = frontend?**
- Your git repo structure stays at the project root
- Vercel treats `frontend/` as the deployment root
- This is standard for monorepos - one repo, multiple deployable apps
- Backend deploys separately (Fly.io), frontend deploys to Vercel

#### Method 2: Using Vercel CLI

```bash
# From project root (not frontend/)
vercel --prod

# When prompted:
# - Set Root Directory to: frontend
# - Accept other defaults
```

**Note:** Types are auto-generated during build via `prebuild` hook!

### Environment Variables

Optional - frontend auto-detects API URL:

```bash
# Production (optional - auto-detected)
NEXT_PUBLIC_API_URL=https://api.mockly.codes
```

Set in Vercel dashboard or CLI:

```bash
vercel env add NEXT_PUBLIC_API_URL production
# Enter: https://api.mockly.codes
```

### Build Process

The frontend build automatically:

1. ✅ Generates TypeScript types from schemas (via `prebuild`)
2. ✅ Builds Next.js app with SSR
3. ✅ Optimizes for production

### Deployment Checklist

- ✅ Types generated (automatic via `prebuild`)
- ✅ Environment variables set (optional)
- ✅ Test build locally: `npm run build`
- ✅ Deploy to Vercel

---

## 🔄 CI/CD Pipeline

### GitHub Actions

The `.github/workflows/build.yml` automatically:

1. Generates TypeScript types
2. Verifies types are in sync
3. Builds backend with schemas
4. Builds frontend with types
5. Runs tests

### Continuous Deployment

**Backend (Fly.io):**

Add to `.github/workflows/deploy.yml`:

```yaml
name: Deploy Backend

on:
  push:
    branches: [main, trunk]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Fly CLI
        uses: superfly/flyctl-actions/setup-flyctl@master
      
      - name: Deploy to Fly.io
        run: flyctl deploy --config backend/fly.toml
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

**Frontend (Vercel):**

Automatic via Vercel GitHub integration - no config needed!

---

## 📊 Deployment Summary

| Aspect | Backend (Fly.io) | Frontend (Vercel) |
|--------|------------------|-------------------|
| **Git Root** | Repository root | Repository root |
| **Deploy Root** | Repository root | `frontend/` (via Root Directory setting) |
| **Build From** | Repository root | `frontend/` directory |
| **Schemas** | Auto-copied during build | N/A |
| **Types** | N/A | Auto-generated via `prebuild` |
| **Command** | `flyctl deploy --config backend/fly.toml` | `vercel --prod` (from project root) |
| **Auto-Deploy** | Via GitHub Actions | Via GitHub integration |
| **Key Setting** | `fly.toml` dockerfile path | Vercel Root Directory = `frontend` ⚠️ |

---

## ✅ Post-Deployment Verification

### Backend

```bash
# Check API is live
curl https://api-mockly.fly.dev/

# Verify all 54 endpoints
curl https://api-mockly.fly.dev/ | jq '.resources | length'

# Test specific endpoints
curl https://api-mockly.fly.dev/users?count=5
curl https://api-mockly.fly.dev/weather?count=3
```

### Frontend

```bash
# Visit website
open https://mockly.codes

# Check docs page
open https://mockly.codes/docs

# Test playground
open https://mockly.codes/playground
```

---

## 🎯 Key Takeaways

### For Backend Deployment:

1. ✅ **Always build from repository root** (not `backend/`)
2. ✅ **Schemas auto-copy** during Docker build
3. ✅ **Use Makefile** for convenience: `make deploy`
4. ✅ **Verify schemas** after deployment

### For Frontend Deployment:

1. ⚠️ **MUST set Root Directory** to `frontend` in Vercel settings
2. ✅ **Types auto-generate** via `prebuild` hook
3. ✅ **No manual steps** needed after initial setup
4. ✅ **Vercel GitHub integration** is easiest
5. ✅ **API URL auto-detected** in production
6. ✅ **Git repo structure** stays at project root (monorepo)

### For Both:

- ✅ **No manual schema syncing** needed
- ✅ **No manual type generation** needed
- ✅ **Everything is automated** in build process
- ✅ **CI/CD enforces** correctness

---

**Questions? Check the main [README.md](./README.md) or [BUILD_PROCESS.md](./BUILD_PROCESS.md)**

