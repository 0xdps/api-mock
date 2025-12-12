# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2024-12-12

### Added - Backend API
- 🎯 **Advanced Query Parameters**
  - Pagination with page, limit, offset
  - Sorting by any field with asc/desc order
  - Full-text search with field targeting (q, search, search_fields)
  - Field filtering (exact match, range, contains, startsWith, endsWith)
  - Field selection for response optimization (fields parameter)
- 🛠️ **Global Middleware Suite** (11+ middleware)
  - Delay simulation (max 30s) for testing timeouts
  - Chaos engineering with flakyRate (0.0-1.0)
  - Cache control with skip_cache parameter
  - Field filtering middleware
  - Multi-tenancy support (X-Tenant-ID header)
  - Role-Based Access Control (X-Role header)
  - Idempotency for POST/PUT/PATCH (24h TTL)
  - Request tracing (X-Request-ID header)
- 📚 **Comprehensive API Documentation**
  - Complete API_DOCUMENTATION.md (400+ lines)
  - All query parameters with examples
  - All middleware features documented
  - Code examples in JavaScript, Python, cURL
  - Response format specifications
  - Error handling guide
  - Best practices section
- 🔍 **Enhanced Cache System**
  - Redis integration with in-memory fallback
  - 10,000 pre-cached items (100 per resource)
  - X-Cache header (HIT/MISS/BYPASS)
  - Cache bypass via query parameter
  - Admin endpoints (/admin/cache/stats, /admin/cache/refresh)
- ⚙️ **Environment Configuration**
  - .env file support for backend development
  - Configurable Redis TLS (REDIS_TLS_ENABLED)
  - Local cache mode for development (CACHE_MODE=local)

### Added - Frontend
- 🎮 **Enhanced Playground - Major UX Redesign**
  - **50/50 split-screen layout**: Request builder (left) and response viewer (right)
  - **Tabbed navigation**: Quick, Pagination, Filters, Middleware, Path tabs
  - **Compact header**: Reduced from text-4xl to text-2xl for more screen space
  - **Full-width layout**: Removed max-width constraints for wider content area
  - **Independent scrolling**: Both panels scroll independently with matching heights
  - Middleware tester with all 11 middleware
  - Pagination controls (page, limit, offset)
  - Sorting controls (field, order)
  - Search with field targeting
  - Field selection with preview
  - Request/response headers display
  - Request timing metrics
  - Real-time URL preview with path options
- 📄 **Sitemap Updates**
  - Added playground utility tools to sitemap.xml
  - Added playground section to sitemap.html
  - 130+ total URLs indexed
- 🎨 **UI Improvements**
  - Search section with comprehensive guide
  - Favicon in header and homepage
  - Active middleware indicators
  - Live query preview
  - Better responsive design
- ⚛️ **Framework Upgrades**
  - Next.js 16 with App Router
  - React 19 with server components
  - TypeScript 5.6+ with improved types
  - Better SSR/ISR performance

### Changed
- 🔧 **Backend Improvements**
  - Fixed pagination bug (reserved params no longer treated as filters)
  - Enhanced filter parser to skip 10+ reserved parameters
  - Improved middleware ordering and registration
  - Better error handling and validation
- 🎨 **Frontend Enhancements**
  - Updated documentation components with modern styling
  - Improved code examples with syntax highlighting
  - Enhanced playground with better UX
  - Better TypeScript type safety
- 📊 **Performance**
  - Faster response times with optimized cache
  - Better pagination performance
  - Improved search algorithm

### Fixed
- ✅ Pagination returning empty results with query params
- ✅ Search UI styling inconsistencies
- ✅ TypeScript build errors in components
- ✅ Missing favicon in navigation
- ✅ Sitemap missing playground pages
- ✅ Filter parser treating pagination params as filters

## [1.1.0] - 2025-11-20

### Added
- ✨ Server-Side Rendering (SSR) for home and documentation pages
- 🔄 Incremental Static Regeneration (ISR) with 5-minute cache
- 📊 Live API status display on documentation page
- 🎯 Enhanced TypeScript types with Review interface
- 📝 Comprehensive frontend documentation

### Changed
- 🚀 Updated frontend to use production API (https://api.mockly.codes)
- ⚡ Improved page load performance with SSR
- 📚 Consolidated documentation into single comprehensive README per module
- 🎨 Enhanced ResourceCard component for SSR compatibility
- 🔧 Simplified environment variable configuration

### Fixed
- ✅ Fixed API URL endpoints (removed incorrect /api prefix)
- 🐛 Resolved cache configuration warnings in Next.js build
- 📖 Corrected framework references (chi instead of Gin)

## [1.0.0] - 2025-11-17

### Added
- 🚀 Initial release of Mockly
- Schema-driven API with automatic endpoint generation
- 6 default resources: users, posts, products, comments, todos, reviews
- Go backend with chi router
- Next.js 14 website with App Router
- Interactive API playground
- Comprehensive documentation page
- Support for 50+ data generators via gofakeit
- CORS enabled for all endpoints
- Fly.io deployment for backend
- Dynamic API URL configuration based on deployment environment
- Embedded schemas for reliable deployment
- Custom route paths and aliases support
- Metadata endpoints for all resources

### API Features
- RESTful endpoints: collection, single item, metadata
- Query parameter support: `?count=N` (max 100)
- ID-based lookups: `/{resource}/{id}`
- Schema metadata: `/{resource}/meta`
- Health check endpoint: `/`
- Comprehensive error handling and validation
- Production endpoint: https://api.mockly.codes

### Website Features
- Server-side rendered pages for better SEO
- Modern landing page with live resource data
- Interactive playground with live API testing
- Documentation with code examples in multiple languages
- Real-time API status and version display
- Resource cards with direct API links
- Responsive design with Tailwind CSS
- Dark theme optimized for developers

### Technical Stack
- **Backend:** Go 1.23+, chi v5, gofakeit v7
- **Frontend:** Next.js 14 (SSR/ISR), TypeScript, Tailwind CSS
- **Deployment:** Backend on Fly.io, Frontend on Vercel
- **Development:** Hot reload, local dev servers

### Documentation
- README with quick start guide
- Contributing guidelines
- MIT License
- Deployment instructions
- Schema creation guide

## [Unreleased]

### Planned
- Additional resource schemas (articles, events, categories)
- Rate limiting for production
- API versioning support
- OpenAPI/Swagger documentation
- GraphQL endpoint
- WebSocket support for real-time data
- Custom domain setup instructions
- Docker support for self-hosting
- Authentication examples (optional)
- Pagination support
- Filtering and sorting capabilities

---

## Version History

### Versioning Scheme
- **Major (X.0.0):** Breaking changes
- **Minor (0.X.0):** New features, backward compatible
- **Patch (0.0.X):** Bug fixes, backward compatible

### Links
- [Website](https://mockly.codes)
- [API](https://api.mockly.codes)
- [Documentation](https://mockly.codes/docs)
- [Playground](https://mockly.codes/playground)
- [GitHub Repository](https://github.com/0xdps/api-mock)
- [Issue Tracker](https://github.com/0xdps/api-mock/issues)
