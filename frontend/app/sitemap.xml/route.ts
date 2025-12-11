import { NextResponse } from 'next/server'
import { getAllResources, getAllGroups, getBaseUrl } from '@/lib/sitemap'

export async function GET() {
  const baseUrl = getBaseUrl()
  const resources = getAllResources()
  const now = new Date().toISOString()
  
  // Static pages
  const staticPages = [
    { loc: '', changefreq: 'daily', priority: 1.0 },
    { loc: '/docs', changefreq: 'weekly', priority: 0.9 },
    { loc: '/playground', changefreq: 'weekly', priority: 0.9 },
    { loc: '/playground/echo', changefreq: 'monthly', priority: 0.8 },
    { loc: '/playground/status', changefreq: 'monthly', priority: 0.8 },
    { loc: '/playground/delay', changefreq: 'monthly', priority: 0.8 },
    { loc: '/playground/middleware', changefreq: 'monthly', priority: 0.8 },
    { loc: '/playground/chaos', changefreq: 'monthly', priority: 0.8 },
  ]
  
  // Resource documentation pages
  const resourcePages = resources.map(resource => ({
    loc: `/docs/${resource.name}`,
    changefreq: 'monthly' as const,
    priority: 0.7,
  }))
  
  // Playground resource pages
  const playgroundResourcePages = resources.map(resource => ({
    loc: `/playground/${resource.name}`,
    changefreq: 'monthly' as const,
    priority: 0.6,
  }))
  
  // Build XML
  const urls = [
    ...staticPages,
    ...resourcePages,
    ...playgroundResourcePages,
  ]
  
  const xml = `<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/xsl" href="/sitemap.xsl"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls.map(url => `  <url>
    <loc>${baseUrl}${url.loc}</loc>
    <lastmod>${now}</lastmod>
    <changefreq>${url.changefreq}</changefreq>
    <priority>${url.priority}</priority>
  </url>`).join('\n')}
</urlset>`
  
  return new NextResponse(xml, {
    headers: {
      'Content-Type': 'application/xml',
      'Cache-Control': 'public, max-age=3600, s-maxage=3600',
    },
  })
}

