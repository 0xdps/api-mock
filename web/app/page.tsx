import Link from 'next/link'
import { ResourceCard } from '@/components/ResourceCard'
import { CodeExample } from '@/components/CodeExample'
import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'

const API_URL = process.env.NEXT_PUBLIC_API_URL || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:8080')

async function getResources() {
  // For now, return static list since API might not be deployed yet
  // In production, this will be fetched from the API
  return ['users', 'posts', 'products', 'comments', 'todos', 'reviews']
  
  /* Uncomment when API is deployed:
  try {
    const res = await fetch(`${API_URL}/`, { cache: 'no-store' })
    const data = await res.json()
    return data.resources || []
  } catch (error) {
    console.error('Failed to fetch resources:', error)
    return []
  }
  */
}

export default async function Home() {
  const resources = await getResources()

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <Header />
      
      {/* Hero Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <h1 className="text-6xl font-bold text-white mb-6">
          apimock<span className="text-primary-500">.codes</span>
        </h1>
        <p className="text-xl text-slate-300 mb-8 max-w-2xl mx-auto">
          Schema-driven mock API service with realistic data. 
          Perfect for prototyping, testing, and demos.
        </p>
        <div className="flex gap-4 justify-center">
          <Link 
            href="/docs" 
            className="bg-primary-500 hover:bg-primary-600 text-white px-8 py-3 rounded-lg font-semibold transition"
          >
            View Documentation
          </Link>
          <Link 
            href="/playground" 
            className="bg-slate-700 hover:bg-slate-600 text-white px-8 py-3 rounded-lg font-semibold transition"
          >
            Try API Playground
          </Link>
        </div>
      </section>

      {/* Features */}
      <section className="container mx-auto px-4 py-16">
        <div className="grid md:grid-cols-3 gap-8">
          <FeatureCard 
            icon="⚡️"
            title="Instant Access"
            description="No signup required. Start using our API endpoints immediately with realistic fake data."
          />
          <FeatureCard 
            icon="📝"
            title="Schema-Driven"
            description="All endpoints are auto-generated from JSON schemas. Add new resources without code."
          />
          <FeatureCard 
            icon="🔄"
            title="RESTful API"
            description="Standard REST endpoints with collection, single item, and metadata routes."
          />
          <FeatureCard 
            icon="🎨"
            title="Realistic Data"
            description="50+ data generators produce realistic names, emails, addresses, and more."
          />
          <FeatureCard 
            icon="🚀"
            title="Fast & Reliable"
            description="Built with Go for high performance. CORS enabled for browser access."
          />
          <FeatureCard 
            icon="🆓"
            title="Free Forever"
            description="Completely free to use. No rate limits. Perfect for learning and testing."
          />
        </div>
      </section>

      {/* Quick Example */}
      <section className="container mx-auto px-4 py-16">
        <h2 className="text-4xl font-bold text-white mb-8 text-center">
          Get Started in Seconds
        </h2>
        <div className="max-w-4xl mx-auto">
          <CodeExample 
            title="Fetch Users"
            code={`fetch('${API_URL}/api/users?count=5')
  .then(res => res.json())
  .then(data => console.log(data));`}
            language="javascript"
          />
          
          <div className="mt-8 bg-slate-800 rounded-lg p-6">
            <p className="text-slate-300 mb-4">Response:</p>
            <pre className="text-green-400 overflow-x-auto">
              <code>{`[
  {
    "id": 42,
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "city": "San Francisco",
    "country": "USA"
  },
  // ... 4 more users
]`}</code>
            </pre>
          </div>
        </div>
      </section>

      {/* Available Resources */}
      <section className="container mx-auto px-4 py-16">
        <h2 className="text-4xl font-bold text-white mb-8 text-center">
          Available Resources
        </h2>
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-6xl mx-auto">
          {resources.map((resource: string) => (
            <ResourceCard key={resource} name={resource} />
          ))}
        </div>
        <p className="text-center text-slate-400 mt-8">
          More resources coming soon! Check the{' '}
          <Link href="/docs" className="text-primary-500 hover:underline">
            documentation
          </Link>
          {' '}for details.
        </p>
      </section>

      {/* Use Cases */}
      <section className="container mx-auto px-4 py-16">
        <h2 className="text-4xl font-bold text-white mb-8 text-center">
          Perfect For
        </h2>
        <div className="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto">
          <UseCaseCard 
            title="Frontend Development"
            description="Build UIs without waiting for backend APIs. Test edge cases with controlled data."
          />
          <UseCaseCard 
            title="Learning & Tutorials"
            description="Practice API integration without complex setup. Perfect for coding bootcamps."
          />
          <UseCaseCard 
            title="Prototyping"
            description="Quickly validate ideas with realistic data. Show demos to stakeholders."
          />
          <UseCaseCard 
            title="Testing"
            description="Test applications with consistent, reproducible data. Mock external dependencies."
          />
        </div>
      </section>

      <Footer />
    </div>
  )
}

function FeatureCard({ icon, title, description }: { icon: string; title: string; description: string }) {
  return (
    <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
      <div className="text-4xl mb-4">{icon}</div>
      <h3 className="text-xl font-semibold text-white mb-2">{title}</h3>
      <p className="text-slate-400">{description}</p>
    </div>
  )
}

function UseCaseCard({ title, description }: { title: string; description: string }) {
  return (
    <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
      <h3 className="text-xl font-semibold text-white mb-2">{title}</h3>
      <p className="text-slate-400">{description}</p>
    </div>
  )
}
