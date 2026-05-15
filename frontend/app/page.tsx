import Link from 'next/link'
import Image from 'next/image'
import { 
  Zap, 
  FileText, 
  Repeat, 
  Palette, 
  Rocket, 
  ShieldCheck, 
  ArrowRight, 
  Github,
  Layout,
  GraduationCap,
  Laptop,
  Play,
  Compass
} from 'lucide-react'
import { CodeExample } from '@/components/CodeExample'
import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'
import { getApiUrl } from '@/lib/api'

const API_URL = getApiUrl()

export default async function Home() {

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <Header />
      
      {/* Hero Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <div className="flex items-center justify-center gap-4 mb-6">
          <Image 
            src="/logo.svg" 
            alt="Mockly Logo" 
            width={64} 
            height={64}
            className="w-16 h-16"
          />
          <h1 className="text-6xl font-bold text-white">
            Mockly
          </h1>
        </div>
        <p className="text-xl text-slate-300 mb-8 max-w-2xl mx-auto">
          Schema-driven mock API service with realistic data. 
          Perfect for prototyping, testing, and demos.
        </p>
        <div className="flex flex-wrap gap-4 justify-center">
          <Link 
            href="/playground" 
            data-umami-event="hero_pg"
            data-umami-label="Hero Explore API"
            className="group bg-primary-500 hover:bg-primary-600 text-white px-8 py-3 rounded-xl font-semibold transition-all flex items-center gap-2 shadow-lg shadow-primary-500/20"
          >
            <Play className="w-4 h-4" /> Try the Playground <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-1" />
          </Link>
          <Link
            href="/templates"
            data-umami-event="hero_tpl"
            data-umami-label="Hero Explore Templates"
            className="group bg-slate-800 hover:bg-slate-700 border border-slate-600 hover:border-primary-500/40 text-white px-8 py-3 rounded-xl font-semibold transition-all flex items-center gap-2"
          >
            <Compass className="w-4 h-4 text-primary-400" /> Explore Templates
          </Link>
          <a 
            href="https://github.com/0xdps/api-mockly" 
            target="_blank"
            rel="noopener noreferrer"
            data-umami-event="hero_gh"
            data-umami-label="Hero GitHub"
            className="bg-slate-700 hover:bg-slate-600 text-white px-8 py-3 rounded-xl font-semibold transition-all flex items-center gap-2"
          >
            <Github className="w-4 h-4" /> View on GitHub
          </a>
        </div>
      </section>

      {/* Features */}
      <section className="container mx-auto px-4 py-16">
        <div className="grid md:grid-cols-3 gap-8">
          <FeatureCard 
            icon={<Zap className="w-6 h-6 text-yellow-400" />}
            title="Instant Access"
            description="No signup required. Start using our API endpoints immediately with realistic fake data."
          />
          <FeatureCard 
            icon={<FileText className="w-6 h-6 text-primary-400" />}
            title="Schema-Driven"
            description="All endpoints are auto-generated from JSON schemas. Add new resources without code."
          />
          <FeatureCard 
            icon={<Repeat className="w-6 h-6 text-green-400" />}
            title="RESTful API"
            description="Standard REST endpoints with collection, single item, and metadata routes."
          />
          <FeatureCard 
            icon={<Palette className="w-6 h-6 text-purple-400" />}
            title="Realistic Data"
            description="50+ data generators produce realistic names, emails, addresses, and more."
          />
          <FeatureCard 
            icon={<Rocket className="w-6 h-6 text-orange-400" />}
            title="Fast & Reliable"
            description="Built with Go for high performance. CORS enabled for browser access."
          />
          <FeatureCard 
            icon={<ShieldCheck className="w-6 h-6 text-blue-400" />}
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
            title="Fetch Users (via group path)"
            code={`// Browse by category
fetch('${API_URL}/people')
  .then(res => res.json())
  .then(data => console.log(data.resources));

// Get users via group path
fetch('${API_URL}/people/users?count=5')
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

      {/* Available Resources - Hidden */}
      {/* <section className="container mx-auto px-4 py-16">
        <h2 className="text-4xl font-bold text-white mb-8 text-center">
          Available Resources
        </h2>
        <p className="text-center text-slate-300 mb-12 max-w-2xl mx-auto">
          Browse {resources.length} resources organized into {Object.keys(groups).length} categories
        </p>
        
        {Object.keys(groups).length > 0 ? (
          <div className="space-y-12 max-w-7xl mx-auto">
            {Object.entries(groups).sort(([a], [b]) => a.localeCompare(b)).map(([groupName, groupData]: [string, any]) => {
              // Handle both array format (direct) and object format (with resources property)
              const resourceList = Array.isArray(groupData) ? groupData : (groupData?.resources || []);
              
              return (
                <div key={groupName}>
                  <div className="flex items-center gap-3 mb-6">
                    <h3 className="text-2xl font-bold text-white capitalize">
                      {getGroupIcon(groupName)} {groupName}
                    </h3>
                    <span className="text-sm text-slate-400 bg-slate-800 px-3 py-1 rounded-full">
                      {resourceList.length} resources
                    </span>
                  </div>
                  <div className="grid md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                    {resourceList.map((resource: string) => (
                      <ResourceCard key={resource} name={resource} apiUrl={API_URL} group={groupName} />
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        ) : (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-6xl mx-auto">
            {resources.map((resource: string) => (
              <ResourceCard key={resource} name={resource} apiUrl={API_URL} />
            ))}
          </div>
        )}
        
        <p className="text-center text-slate-400 mt-12">
          Check the{' '}
          <Link href="/resources" className="text-primary-500 hover:underline">
            resources
          </Link>
          {' '}for detailed API usage and examples.
        </p>
      </section> */}

      {/* Use Cases */}
      <section className="container mx-auto px-4 py-16">
        <h2 className="text-4xl font-bold text-white mb-8 text-center">
          Perfect For
        </h2>
        <div className="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto">
          <UseCaseCard 
            icon={<Layout className="w-5 h-5 text-primary-400" />}
            title="Frontend Development"
            description="Build UIs without waiting for backend APIs. Test edge cases with controlled data."
          />
          <UseCaseCard 
            icon={<GraduationCap className="w-5 h-5 text-green-400" />}
            title="Learning & Tutorials"
            description="Practice API integration without complex setup. Perfect for coding bootcamps."
          />
          <UseCaseCard 
            icon={<Rocket className="w-5 h-5 text-orange-400" />}
            title="Prototyping"
            description="Quickly validate ideas with realistic data. Show demos to stakeholders."
          />
          <UseCaseCard 
            icon={<Laptop className="w-5 h-5 text-blue-400" />}
            title="Testing"
            description="Test applications with consistent, reproducible data. Mock external dependencies."
          />
        </div>
      </section>

      <Footer />
    </div>
  )
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="bg-slate-800/40 backdrop-blur-md p-6 rounded-xl border border-white/5 shadow-premium hover:border-primary-500/30 transition-all duration-300 group">
      <div className="mb-4 p-3 bg-white/5 rounded-lg w-fit transition-transform group-hover:scale-110 group-hover:bg-white/10">{icon}</div>
      <h3 className="text-xl font-semibold text-white mb-2">{title}</h3>
      <p className="text-slate-400 leading-relaxed text-sm">{description}</p>
    </div>
  )
}

function UseCaseCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="bg-slate-800/40 backdrop-blur-md p-6 rounded-xl border border-white/5 shadow-premium hover:border-primary-500/30 transition-all duration-300 group flex gap-4 items-start">
      <div className="mt-1 p-2 bg-white/5 rounded-lg transition-transform group-hover:scale-110 group-hover:bg-white/10">{icon}</div>
      <div>
        <h3 className="text-lg font-semibold text-white mb-1">{title}</h3>
        <p className="text-slate-400 leading-relaxed text-sm">{description}</p>
      </div>
    </div>
  )
}
