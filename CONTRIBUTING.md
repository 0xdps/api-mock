# Contributing to Mockly

Thank you for your interest in contributing to Mockly! 🎉

## How to Contribute

There are many ways to contribute to this project:

### 1. Add New Resource Schemas 📝

The easiest way to contribute is by adding new resource schemas. All endpoints are auto-generated from JSON schemas.

**Steps:**

1. **Fork the repository**
2. **Create a new schema** in `shared/schemas/your-resource.json`
3. **Test locally:**
   ```bash
   cd local
   go run main.go
   # Test your endpoint: curl http://localhost:8080/api/your-resource?count=5
   ```
4. **Submit a Pull Request**

**Schema Example:**

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Book",
  "type": "object",
  "x-resource": {
    "name": "books",
    "singular": "book",
    "description": "Library books",
    "routes": {
      "path": "/books",
      "methods": ["GET"]
    }
  },
  "properties": {
    "id": {
      "type": "integer",
      "x-generator": "random_int",
      "x-generator-params": { "min": 1, "max": 10000 }
    },
    "title": {
      "type": "string",
      "x-generator": "sentence"
    },
    "author": {
      "type": "string",
      "x-generator": "name"
    },
    "isbn": {
      "type": "string",
      "x-generator": "uuid"
    },
    "published_date": {
      "type": "string",
      "x-generator": "past_date"
    }
  }
}
```

### 2. Report Bugs 🐛

If you find a bug, please create an issue with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Screenshots if applicable

### 3. Suggest Features 💡

Have an idea? Open an issue with:
- Clear description of the feature
- Use case / why it's needed
- Possible implementation approach

### 4. Improve Documentation 📚

Help improve:
- README clarity
- Code comments
- API documentation
- Examples and tutorials

### 5. Fix Issues 🔧

Look for issues labeled:
- `good first issue` - Great for beginners
- `help wanted` - We need help with these
- `bug` - Something isn't working

## Development Setup

### Prerequisites

- Go 1.22+
- Node.js 18+
- npm or yarn

### Local Setup

1. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/api-mock.git
   cd api-mock
   ```

2. **Run the API:**
   ```bash
   cd local
   go run main.go
   ```

3. **Run the website:**
   ```bash
   cd web
   npm install
   npm run dev
   ```

### Project Structure

```
api-mock/
├── api/              # Vercel serverless API handler
│   └── index.go      # Entry point for Vercel deployment
├── local/            # Local development API server
│   └── main.go       # Entry point for local testing
├── lib/              # Shared Go libraries
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # CORS, etc.
│   ├── schema/       # Schema loader & generator
│   └── shared/       # Embedded schemas for deployment
├── shared/           # Source of truth
│   └── schemas/      # JSON schema definitions
├── web/              # Next.js website
│   ├── app/          # App router pages
│   ├── components/   # React components
│   └── lib/          # Utility functions
└── vercel.json       # Deployment config
```

## Code Style Guidelines

### Go Code

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

### TypeScript/React

- Use TypeScript for type safety
- Follow React best practices
- Use functional components with hooks
- Keep components small and focused

### JSON Schemas

- Follow JSON Schema Draft 7 spec
- Include clear descriptions
- Use appropriate data generators
- Test locally before submitting

## Pull Request Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Write clear, concise commits
   - Add tests if applicable
   - Update documentation

3. **Test your changes:**
   ```bash
   # Test API
   cd local && go run main.go
   
   # Test website
   cd web && npm run dev
   ```

4. **Submit PR:**
   - Clear title and description
   - Reference related issues
   - Include screenshots if UI changes
   - Wait for review

## Commit Message Guidelines

Use clear, descriptive commit messages:

```
feat: add books resource schema
fix: correct user email generator
docs: update contributing guide
refactor: simplify schema loader
```

Prefixes:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code refactoring
- `test:` - Tests
- `chore:` - Maintenance

## Available Data Generators

When creating schemas, you can use these generators:

**Personal:**
- `name`, `first_name`, `last_name`, `email`, `username`, `password`

**Address:**
- `address`, `city`, `state`, `country`, `zip_code`, `latitude`, `longitude`

**Company:**
- `company`, `job`, `catch_phrase`

**Internet:**
- `url`, `domain_name`, `ipv4`, `ipv6`, `uuid`, `mac_address`

**Dates:**
- `date`, `date_time`, `past_date`, `future_date`

**Text:**
- `word`, `sentence`, `paragraph`, `text`

**Numbers:**
- `random_int`, `random_digit`, `random_number`

**Other:**
- `phone_number`, `boolean`, `user_agent`

See [gofakeit documentation](https://github.com/brianvoe/gofakeit) for more.

## Testing

Before submitting:

1. **Test the API endpoint:**
   ```bash
   curl http://localhost:8080/api/your-resource?count=5
   ```

2. **Verify JSON response:**
   - Check all fields are present
   - Data types are correct
   - Generators work as expected

3. **Test metadata endpoint:**
   ```bash
   curl http://localhost:8080/api/your-resource/meta
   ```

## Questions?

- Open an issue for questions
- Check existing issues/PRs first
- Be respectful and constructive

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Provide constructive feedback
- Focus on the issue, not the person

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Mockly! 🚀
