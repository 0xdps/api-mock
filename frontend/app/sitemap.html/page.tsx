import Link from 'next/link'
import { getAllResources, getAllGroups, getBaseUrl } from '@/lib/sitemap'
import { Metadata } from 'next'
import { Users, Briefcase, ShoppingCart, FileText, MessageCircle, Film, Plane, Globe, DollarSign, UtensilsCrossed, GraduationCap, Trophy, CheckSquare, Book, Package, LucideIcon } from 'lucide-react'

export const metadata: Metadata = {
  title: 'Sitemap | Mockly',
  description: 'Complete sitemap of all pages and resources available on Mockly.',
}

function getGroupIcon(group: string): LucideIcon {
  const icons: Record<string, LucideIcon> = {
    people: Users,
    business: Briefcase,
    commerce: ShoppingCart,
    content: FileText,
    social: MessageCircle,
    media: Film,
    travel: Plane,
    location: Globe,
    finance: DollarSign,
    food: UtensilsCrossed,
    education: GraduationCap,
    sports: Trophy,
    productivity: CheckSquare,
    reference: Book,
  }
  return icons[group] || Package
}

export default function SitemapPage() {
  const baseUrl = getBaseUrl()
  const resources = getAllResources()
  const groups = getAllGroups()
  
  // Static pages
  const staticPages = [
    { path: '/', label: 'Home', description: 'Mockly homepage' },
    { path: '/docs', label: 'Documentation', description: 'Full API reference, query parameters, and code examples' },
    { path: '/templates', label: 'Explore Templates', description: 'Browse 100+ official and community mock data templates' },
    { path: '/playground', label: 'Playground', description: 'Interactive API testing tools' },
  ]
  
  // Playground utility pages
  const playgroundPages = [
    { path: '/playground/echo', label: 'Echo Tester', description: 'Inspect request headers and body' },
    { path: '/playground/status', label: 'Status Code Generator', description: 'Test HTTP status codes' },
    { path: '/playground/delay', label: 'Delay Tester', description: 'Test API latency and timeouts' },
    { path: '/playground/middleware', label: 'Middleware Testing', description: 'Test global middleware parameters' },
    { path: '/playground/chaos', label: 'Chaos Engineering', description: 'Test fault tolerance and error handling' },
  ]
  
  // Resource pages grouped by group
  const resourcesByGroup = Object.entries(groups).sort(([a], [b]) => a.localeCompare(b))
  
  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-950 via-slate-900 to-slate-950">
      <div className="container mx-auto px-4 py-16 max-w-6xl">
        <div className="mb-12">
          <h1 className="text-4xl font-bold text-white mb-4">Sitemap</h1>
          <p className="text-slate-400 text-lg">
            Complete list of all pages and resources available on Mockly
          </p>
        </div>
        
        {/* Static Pages */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">Main Pages</h2>
          <div className="grid md:grid-cols-3 gap-4">
            {staticPages.map((page) => (
              <Link
                key={page.path}
                href={page.path}
                className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
              >
                <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                  {page.label}
                </h3>
                <p className="text-slate-400 text-sm">{page.description}</p>
                <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  Visit →
                </p>
              </Link>
            ))}
          </div>
        </section>
        
        {/* Playground Pages */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">Playground Tools</h2>
          <p className="text-slate-400 mb-4">
            Interactive testing tools for exploring API features
          </p>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            {playgroundPages.map((page) => (
              <Link
                key={page.path}
                href={page.path}
                className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
              >
                <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                  {page.label}
                </h3>
                <p className="text-slate-400 text-sm">{page.description}</p>
                <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  Try it →
                </p>
              </Link>
            ))}
          </div>
        </section>
        
        {/* Resource Documentation Pages */}
        <section>
          <h2 className="text-2xl font-bold text-white mb-6">API Resources</h2>
          <div className="space-y-8">
            {resourcesByGroup.map(([group, groupResources]) => (
              <div key={group}>
                <div className="flex items-center gap-3 mb-4">
                  <h3 className="text-xl font-semibold text-white flex items-center gap-2">
                    {(() => {
                      const Icon = getGroupIcon(group)
                      return <Icon className="w-5 h-5" />
                    })()}
                    {group.charAt(0).toUpperCase() + group.slice(1)}
                  </h3>
                  <span className="text-sm text-slate-400 bg-slate-800 px-3 py-1 rounded-full">
                    {groupResources.length} resources
                  </span>
                </div>
                <div className="grid md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
                  {groupResources.sort().map((resource) => (
                    <Link
                      key={resource}
                      href={`/resources/${resource}`}
                      className="block p-4 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
                    >
                      <span className="text-slate-300 font-medium group-hover:text-blue-400 transition-colors">
                        {resource}
                      </span>
                      <p className="text-blue-400 text-xs mt-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        View docs →
                      </p>
                    </Link>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </section>
        
        {/* External Links */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-white mb-4 mt-4">External Links</h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            <a
              href="https://github.com/0xdps/api-mockly"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
            >
              <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                GitHub
              </h3>
              <p className="text-slate-400 text-sm">Source code and contributions</p>
              <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                Visit →
              </p>
            </a>
            <a
              href="https://www.mockly.codes/"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
            >
              <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                Mockly
              </h3>
              <p className="text-slate-400 text-sm">Main project website</p>
              <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                Visit →
              </p>
            </a>
            <a
              href="https://pinboard-gpt.dps.codes/"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
            >
              <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                Pinboard GPT
              </h3>
              <p className="text-slate-400 text-sm">Related project</p>
              <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                Visit →
              </p>
            </a>
            <a
              href="https://devutil.dps.codes/"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
            >
              <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                DevUtil
              </h3>
              <p className="text-slate-400 text-sm">Developer utilities</p>
              <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                Visit →
              </p>
            </a>
            <a
              href="https://www.pingpong.codes/"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
            >
              <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                PingPong
              </h3>
              <p className="text-slate-400 text-sm">Related project</p>
              <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                Visit →
              </p>
            </a>
            <a
              href="https://fake-stack.readthedocs.io/"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
            >
              <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                Fake Stack Docs
              </h3>
              <p className="text-slate-400 text-sm">Documentation</p>
              <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                Visit →
              </p>
            </a>
            <a
              href="https://dps.codes"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg hover:border-blue-500/50 transition-colors group"
            >
              <h3 className="text-white font-semibold mb-2 group-hover:text-blue-400 transition-colors">
                0xdps
              </h3>
              <p className="text-slate-400 text-sm">Powered by 0xdps</p>
              <p className="text-blue-400 text-sm mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                Visit →
              </p>
            </a>
          </div>
        </section>
        
        {/* XML Sitemap Link */}
        <div className="mt-12 p-6 bg-slate-800/50 border border-slate-700/50 rounded-lg text-center">
          <p className="text-slate-400 mb-2">
            For search engines, view the{' '}
            <Link href="/sitemap.xml" className="text-blue-400 hover:text-blue-300 underline">
              XML sitemap
            </Link>
          </p>
          <p className="text-slate-500 text-sm">
            Total: {staticPages.length + playgroundPages.length + resources.length} pages
          </p>
        </div>
      </div>
    </div>
  )
}

