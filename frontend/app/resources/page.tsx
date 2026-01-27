import { redirect } from 'next/navigation'
import { getAllResources } from '@/lib/sitemap'

export const metadata = {
  title: 'API Documentation | Mockly - Free Mock API for Developers',
  description: 'Comprehensive API documentation for 100+ mock data resources across 14 categories. Free schema-driven mock API service for testing and development.',
}

export default function DocsPage() {
  // Get all available resources
  const resources = getAllResources()
  
  // Pick a random resource
  const randomResource = resources[Math.floor(Math.random() * resources.length)]
  
  // Redirect to random resource
  redirect(`/resources/${randomResource.name}`)
}
