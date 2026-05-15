'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'
import { listPublicTemplates } from '@/lib/mockly-api'
import type { MocklyTemplate } from '@/types/mockly'
import { Globe, FileCode2, Search, ArrowRight, Loader2, Plus, Star, ExternalLink } from 'lucide-react'
import { schemasManifest } from '@/lib/schemas-manifest'

const API_BASE = 'https://api.mockly.codes/v1'
const PAGE_SIZE = 24

// Maps plural API slug (e.g. "carts") → singular schema name (e.g. "cart")
// Used to build correct /resources/[name] links for featured templates
const resourceSlugToSchemaName = new Map(
  schemasManifest.map(s => [s.schema?.['x-resource']?.name as string, s.name])
)

export function ExploreTemplatesClient() {
  const [templates, setTemplates] = useState<MocklyTemplate[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [offset, setOffset] = useState(0)
  const [hasMore, setHasMore] = useState(false)
  const [loadingMore, setLoadingMore] = useState(false)
  const [total, setTotal] = useState<number | null>(null)

  useEffect(() => {
    setLoading(true)
    setOffset(0)
    listPublicTemplates(PAGE_SIZE, 0)
      .then(r => {
        setTemplates(r.templates ?? [])
        setTotal(r.total ?? null)
        setHasMore((r.templates?.length ?? 0) === PAGE_SIZE)
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  async function loadMore() {
    const next = offset + PAGE_SIZE
    setLoadingMore(true)
    try {
      const r = await listPublicTemplates(PAGE_SIZE, next)
      setTemplates(prev => [...prev, ...(r.templates ?? [])])
      setOffset(next)
      setHasMore((r.templates?.length ?? 0) === PAGE_SIZE)
    } finally {
      setLoadingMore(false)
    }
  }

  const filtered = templates.filter(t =>
    !search ||
    t.name.toLowerCase().includes(search.toLowerCase()) ||
    t.slug.toLowerCase().includes(search.toLowerCase()) ||
    (t.description && t.description.toLowerCase().includes(search.toLowerCase()))
  )

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col">
      <Header />

      <main className="flex-1 container mx-auto px-4 py-12">
        {/* Hero */}
        <div className="text-center mb-12 relative">
          <div className="absolute inset-0 -top-12 pointer-events-none">
            <div className="w-96 h-96 bg-primary-500/5 rounded-full blur-3xl mx-auto" />
          </div>
          <div className="relative">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-primary-500/10 border border-primary-500/20 text-primary-400 text-xs font-medium mb-5">
              <Globe className="w-3.5 h-3.5" /> 100+ Official Templates
            </div>
            <h1 className="text-4xl font-bold text-white mb-3">Explore Templates</h1>
            <p className="text-slate-400 max-w-lg mx-auto text-sm leading-relaxed">
              100+ official schemas and community templates. Pick one, call the endpoint, get realistic data instantly.
            </p>
          </div>
        </div>

        {/* Search */}
        <div className="relative max-w-lg mx-auto mb-10">
          <Search className="absolute left-3.5 top-3 w-4 h-4 text-slate-500 pointer-events-none" />
          <input
            type="text"
            placeholder="Search by name, slug, or description…"
            value={search}
            onChange={e => setSearch(e.target.value)}
            className="w-full pl-10 pr-4 py-3 rounded-xl bg-slate-900/60 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:border-primary-500 transition-colors shadow-lg shadow-black/20"
          />
        </div>

        {/* Grid */}
        {loading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {Array.from({ length: 12 }).map((_, i) => (
              <div key={i} className="h-40 rounded-xl bg-slate-900/50 border border-slate-700 animate-pulse" />
            ))}
          </div>
        ) : error ? (
          <div className="text-center py-20">
            <p className="text-red-400 text-sm">{error}</p>
          </div>
        ) : filtered.length === 0 ? (
          <div className="relative overflow-hidden text-center py-20 rounded-2xl bg-slate-900/30 border border-dashed border-slate-700">
            <div className="absolute inset-0 bg-gradient-to-br from-primary-500/5 via-transparent to-primary-500/5 pointer-events-none" />
            <FileCode2 className="w-10 h-10 text-slate-600 mx-auto mb-3" />
            <p className="text-slate-300 font-semibold mb-1">{search ? 'No templates match your search' : 'No public templates yet'}</p>
            <p className="text-slate-600 text-sm mb-6">Be the first to publish one.</p>
            <Link
              href="/dashboard/templates/new"
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-primary-600 text-white text-sm font-medium hover:bg-primary-500 transition-colors shadow-lg shadow-primary-500/20"
            >
              <Plus className="w-4 h-4" /> Create a template
            </Link>
          </div>
        ) : (
          <>
            {!search && (
              <p className="text-slate-500 text-xs mb-4">
                {total !== null ? total : filtered.length} template{(total ?? filtered.length) !== 1 ? 's' : ''} available
              </p>
            )}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {filtered.map(t => (
                <TemplateCard key={t.id} template={t} />
              ))}
            </div>

            {hasMore && !search && (
              <div className="text-center mt-8">
                <button
                  onClick={loadMore}
                  disabled={loadingMore}
                  className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-slate-900/60 border border-slate-700 text-slate-300 text-sm hover:border-primary-500/30 hover:text-white transition-all disabled:opacity-50"
                >
                  {loadingMore ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                  Load more
                </button>
              </div>
            )}
          </>
        )}
      </main>

      <Footer />
    </div>
  )
}

function TemplateCard({ template: t }: { template: MocklyTemplate }) {
  const featuredEndpoint = t.is_featured ? `${API_BASE}/${t.slug}` : null
  const resourceName = resourceSlugToSchemaName.get(t.slug) ?? t.slug
  const detailHref = t.is_featured ? `/resources/${resourceName}` : `/templates/${t.id}`

  return (
    <Link
      href={detailHref}
      className="group relative bg-slate-900/50 border border-slate-700 rounded-xl p-5 hover:border-primary-500/30 hover:bg-slate-900/80 transition-all flex flex-col gap-3"
    >
      <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-primary-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />

      <div className="flex items-start justify-between gap-2">
        <div className="w-9 h-9 rounded-lg bg-primary-500/10 border border-primary-500/20 flex items-center justify-center shrink-0">
          <FileCode2 className="w-4 h-4 text-primary-400" />
        </div>
        <div className="flex items-center gap-1.5 shrink-0 flex-wrap justify-end">
          {t.is_featured && (
            <span className="flex items-center gap-1 text-xs text-amber-400 bg-amber-500/10 border border-amber-500/20 px-2 py-0.5 rounded-full">
              <Star className="w-3 h-3" /> Official
            </span>
          )}
          <span className="flex items-center gap-1 text-xs text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 rounded-full">
            <Globe className="w-3 h-3" /> Public
          </span>
        </div>
      </div>

      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-white truncate group-hover:text-primary-300 transition-colors">{t.name}</p>
        {t.description && (
          <p className="text-xs text-slate-500 mt-0.5 line-clamp-2">{t.description}</p>
        )}
      </div>

      <div className="flex items-center gap-2 mt-auto" onClick={e => e.stopPropagation()}>
        {featuredEndpoint ? (
          <>
            <div className="flex-1 min-w-0 bg-black/30 border border-slate-700/60 rounded-lg px-2.5 py-1.5">
              <p className="text-xs font-mono text-slate-400 truncate">{featuredEndpoint}</p>
            </div>
            <a
              href={featuredEndpoint + '?limit=3'}
              target="_blank"
              rel="noopener noreferrer"
              className="shrink-0 flex items-center justify-center w-7 h-7 rounded-lg bg-primary-500/10 border border-primary-500/20 text-primary-400 hover:bg-primary-500/20 hover:text-primary-300 transition-colors"
              title="Open live endpoint"
            >
              <ExternalLink className="w-3.5 h-3.5" />
            </a>
          </>
        ) : (
          <div className="flex-1 min-w-0 bg-black/30 border border-slate-700/60 rounded-lg px-2.5 py-1.5 flex items-center gap-1.5">
            <p className="text-xs font-mono text-slate-400 truncate">Open playground</p>
            <ArrowRight className="w-3 h-3 text-slate-500 shrink-0" />
          </div>
        )}
      </div>
    </Link>
  )
}
