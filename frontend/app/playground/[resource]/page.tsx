import { Metadata } from 'next'
import { ResourceDocumentation } from '@/components/ResourceDocumentation'
import fs from 'fs'
import path from 'path'

function toTitleCase(str: string): string {
  return str
    .split(/[\s_-]+/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}

async function getSchemas() {
  const schemasDir = path.join(process.cwd(), '../shared/schemas')
  const schemas: Array<{ name: string; schema: any; group: string }> = []
  
  const entries = fs.readdirSync(schemasDir, { withFileTypes: true })
  
  for (const entry of entries) {
    if (entry.isDirectory()) {
      const groupDir = path.join(schemasDir, entry.name)
      const files = fs.readdirSync(groupDir).filter((f: string) => f.endsWith('.json'))
      
      for (const file of files) {
        const content = fs.readFileSync(path.join(groupDir, file), 'utf-8')
        schemas.push({
          name: file.replace('.json', ''),
          schema: JSON.parse(content),
          group: entry.name
        })
      }
    }
  }
  
  return schemas.sort((a, b) => a.name.localeCompare(b.name))
}

export async function generateStaticParams() {
  const schemas = await getSchemas()
  return schemas.map((schema) => ({
    resource: schema.name,
  }))
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ resource: string }>
}): Promise<Metadata> {
  const { resource } = await params
  const schemas = await getSchemas()
  const schema = schemas.find(s => s.name === resource)
  
  if (schema) {
    const title = toTitleCase(resource)
    const description = schema.schema.description || `API testing playground for ${title} resource. Test pagination, sorting, search, and middleware features.`
    
    return {
      title: `${title} Testing | Mockly Playground`,
      description,
      openGraph: {
        title: `${title} Testing | Mockly Playground`,
        description,
        type: 'website',
      },
      twitter: {
        card: 'summary',
        title: `${title} Testing | Mockly Playground`,
        description,
      },
    }
  }
  
  return {
    title: 'Resource Testing | Mockly Playground',
    description: 'Test API resources in the Mockly playground.',
  }
}

export default async function ResourcePlaygroundPage({
  params,
}: {
  params: Promise<{ resource: string }>
}) {
  const { resource } = await params
  const schemas = await getSchemas()
  const schema = schemas.find(s => s.name === resource)
  
  if (!schema) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center text-slate-400 mt-20">
          <p className="text-xl">Resource not found</p>
          <a 
            href="/playground"
            className="mt-4 inline-block text-blue-400 hover:text-blue-300 underline"
          >
            ← Back to Playground
          </a>
        </div>
      </div>
    )
  }
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Breadcrumb */}
      <div className="mb-4 flex items-center gap-2 text-sm text-slate-400">
        <a href="/playground" className="hover:text-slate-300">
          Playground
        </a>
        <span>/</span>
        <a href="/playground?tab=resources" className="hover:text-slate-300">
          Resources
        </a>
        <span>/</span>
        <span className="text-white">{toTitleCase(resource)}</span>
      </div>

      <ResourceDocumentation 
        resource={schema.name}
        schema={schema.schema}
        group={schema.group}
      />
    </div>
  )
}
