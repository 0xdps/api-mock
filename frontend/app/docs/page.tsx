import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'API Documentation | Mockly - Free Mock API for Developers',
  description: 'Comprehensive API documentation for 100+ mock data resources across 14 categories. Free schema-driven mock API service for testing and development.',
}

export default function DocsPage() {
  return (
    <div className="max-w-5xl mx-auto p-8">
      <div className="text-center text-slate-400 mt-20">
        <p className="text-xl">Select a resource to view documentation</p>
      </div>
    </div>
  )
}
