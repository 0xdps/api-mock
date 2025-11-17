# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-17

### Added
- 🚀 Initial release of Mockly
- Schema-driven API with automatic endpoint generation
- 6 default resources: users, posts, products, comments, todos, reviews
- Go backend with Gin framework
- Next.js 14 website with App Router
- Interactive API playground
- Comprehensive documentation page
- Support for 50+ data generators via gofakeit
- CORS enabled for all endpoints
- Vercel deployment configuration for monorepo
- Dynamic API URL configuration based on deployment environment
- Embedded schemas for reliable deployment
- Custom route paths and aliases support
- Metadata endpoints for all resources

### API Features
- RESTful endpoints: collection, single item, metadata
- Query parameter support: `?count=N` (max 1000)
- ID-based lookups: `/api/{resource}/:id`
- Schema metadata: `/api/{resource}/meta`
- Health check endpoint: `/api`
- Comprehensive error handling and validation

### Website Features
- Modern landing page with feature highlights
- Interactive playground with live API testing
- Documentation with code examples in multiple languages
- Resource cards with direct API links
- Responsive design with Tailwind CSS
- Dark theme optimized for developers

### Technical Stack
- **Backend:** Go 1.22+, Gin v1.11.0, gofakeit v7
- **Frontend:** Next.js 14, TypeScript, Tailwind CSS
- **Deployment:** Vercel (serverless functions)
- **Development:** Hot reload, local dev server

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
- [Live Site](https://apimock02.vercel.app)
- [GitHub Repository](https://github.com/0xdps/api-mock)
- [Issue Tracker](https://github.com/0xdps/api-mock/issues)
