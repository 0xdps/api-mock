import { Metadata } from 'next'
import { ResourceDocumentation } from '@/components/ResourceDocumentation'
import { schemasManifest } from '@/lib/schemas-manifest'

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
    const description = schema.schema.description || `API documentation for ${title} resource. Free mock data API for testing and development.`
    
    return {
      title: `${title} API Documentation | Mockly`,
      description,
      openGraph: {
        title: `${title} API Documentation | Mockly`,
        description,
        type: 'website',
      },
      twitter: {
        card: 'summary',
        title: `${title} API Documentation | Mockly`,
        description,
      },
    }
  }
  
  return {
    title: 'API Documentation | Mockly',
    description: 'API documentation for Mockly resources.',
  }
}

export default async function ResourcePage({
  params,
}: {
  params: Promise<{ resource: string }>
}) {
  const { resource } = await params
  const schemas = await getSchemas()
  const schema = schemas.find(s => s.name === resource)
  
  if (!schema) {
    return (
      <div className="mx-auto p-4">
        <div className="text-center text-slate-400 mt-20">
          <p className="text-xl">Resource not found</p>
        </div>
      </div>
    )
  }
  
  return (
    <div className="mx-auto p-4">
      <ResourceDocumentation 
        resource={schema.name}
        schema={schema.schema}
        group={schema.group}
      />
    </div>
  )
}

