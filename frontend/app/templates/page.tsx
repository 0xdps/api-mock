import type { Metadata } from 'next'
import { ExploreTemplatesClient } from './ExploreTemplatesClient'

export const metadata: Metadata = {
  title: 'Explore Templates | Mockly — Free Mock API',
  description:
    'Browse 100+ official schemas and community templates. Call any endpoint and get realistic fake data instantly — no setup required.',
  openGraph: {
    title: 'Explore Templates | Mockly',
    description: 'Browse 100+ official schemas and community templates. Free mock API for developers.',
    url: 'https://www.mockly.codes/templates',
    siteName: 'Mockly',
    type: 'website',
  },
}

export default function Page() {
  return <ExploreTemplatesClient />
}
