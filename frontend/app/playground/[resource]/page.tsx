import { Metadata } from 'next'
import { ResourceDocumentation } from '@/components/ResourceDocumentation'
import { schemasManifest } from '@/lib/schemas-manifest'
import Link from 'next/link'

function toTitleCase(str: string): string {
  return str
    .split(/[\s_-]+/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}

async function getSchemas() {
  return schemasManifest
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
  
  // Try exact match first, then try removing 's' for plural
  let schema = schemas.find(s => s.name === resource)
  if (!schema && resource.endsWith('s')) {
    schema = schemas.find(s => s.name === resource.slice(0, -1))
  }
  
  if (!schema) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center text-slate-400 mt-20">
          <p className="text-xl">Resource not found</p>
          <Link 
            href="/playground"
            className="mt-4 inline-block text-blue-400 hover:text-blue-300 underline"
          >
            ← Back to Playground
          </Link>
        </div>
      </div>
    )
  }
  
  return (
    <div className="mx-auto py-8">
      {/* Breadcrumb */}
      <div className="mb-4 flex items-center gap-2 text-sm text-slate-400 px-3 sm:px-4 lg:px-6">
        <Link href="/playground" className="hover:text-slate-300">
          Playground
        </Link>
        <span>/</span>
        <Link href="/playground?tab=resources" className="hover:text-slate-300">
          Resources
        </Link>
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
