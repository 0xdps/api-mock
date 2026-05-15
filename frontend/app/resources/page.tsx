import { redirect } from 'next/navigation'

export const metadata = {
  title: 'API Documentation | Mockly - Free Mock API for Developers',
  description: 'Comprehensive API documentation for 100+ mock data resources across 14 categories. Free schema-driven mock API service for testing and development.',
}

export default function DocsPage() {
  redirect('/templates')
}
