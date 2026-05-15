import type { Metadata } from 'next'
import { getPublicTemplate } from '@/lib/mockly-api'
import { TemplateDetailClient } from './TemplateDetailClient'

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>
}): Promise<Metadata> {
  const { id } = await params
  try {
    const { template } = await getPublicTemplate(id)
    return {
      title: `${template.name} | Mockly Templates`,
      description:
        template.description ||
        `Explore the ${template.name} template on Mockly — free mock API for testing and development.`,
      openGraph: {
        title: `${template.name} | Mockly Templates`,
        description:
          template.description ||
          `Explore the ${template.name} template on Mockly. Free mock API for developers.`,
        url: `https://www.mockly.codes/templates/${id}`,
        siteName: 'Mockly',
        type: 'website',
      },
    }
  } catch {
    return { title: 'Template | Mockly' }
  }
}

export default async function Page({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  return <TemplateDetailClient id={id} />
}
