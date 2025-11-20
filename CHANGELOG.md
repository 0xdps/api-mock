# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
