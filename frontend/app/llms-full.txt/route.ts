import { NextResponse } from 'next/server'
import { readFileSync } from 'fs'
import { join } from 'path'

const TEMPLATE_API_ADDENDUM = `

---

## User Template API

In addition to the 100+ built-in public resources, Mockly supports user-created templates.

### What are Templates?

Templates are custom JSON schemas created by users that define the shape of the mock data returned. Each template gets its own endpoint protected by an API key.

### Creating Templates

1. Sign in at https://www.mockly.codes
2. Go to https://www.mockly.codes/dashboard/templates/new
3. Define a JSON Schema with your desired fields
4. Generate an API key at https://www.mockly.codes/dashboard/api-keys

### Template Endpoint

\`\`\`
GET https://api.mockly.codes/t/{templateId}
Authorization: Bearer <api-key>
\`\`\`

All the same query parameters (pagination, sorting, search, filtering, field selection) work on template endpoints.

### Example

\`\`\`bash
# Create a template with schema:
# { "properties": { "name": { "type": "string" }, "score": { "type": "integer" } } }

# Access the template
curl -H "Authorization: Bearer mk_your_key_here" \\
  "https://api.mockly.codes/t/your-template-id?limit=5&sort=score&order=desc"
\`\`\`

### Template Visibility

- **Private** — Only accessible with your own API key
- **Public** — Accessible with any Mockly API key (\`mak_\` prefix)

### Explore Public Templates

Browse community templates at: https://www.mockly.codes/templates

---

*This document was generated from the Mockly API source. For interactive documentation, visit https://www.mockly.codes/docs*
`

export async function GET() {
  let content: string
  try {
    // Read the backend API_DOCUMENTATION.md at request time
    const docPath = join(process.cwd(), '..', 'backend', 'API_DOCUMENTATION.md')
    content = readFileSync(docPath, 'utf-8')
  } catch {
    content = '# Mockly API — Full Documentation\n\nSee https://www.mockly.codes/docs for the full reference.\n'
  }

  return new NextResponse(content + TEMPLATE_API_ADDENDUM, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      'Cache-Control': 'public, max-age=3600, s-maxage=3600',
    },
  })
}
