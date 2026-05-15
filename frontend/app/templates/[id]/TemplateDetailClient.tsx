'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import {
  ArrowLeft,
  FileCode2,
  Globe,
  Lock,
  Star,
  Calendar,
  Hash,
  Loader2,
  AlertCircle,
  ExternalLink,
} from 'lucide-react'
import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'
import { TemplatePlayground } from '@/components/TemplatePlayground'
import { getPublicTemplate } from '@/lib/mockly-api'
import type { MocklyTemplate } from '@/types/mockly'

const PUBLIC_API_BASE = 'https://api.mockly.codes'

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

function SchemaFieldsTable({
  properties,
  required,
}: {
  properties: Record<string, any>
  required: string[]
}) {
  const entries = Object.entries(properties)
  if (entries.length === 0) return null

  return (
    <div className="overflow-x-auto rounded-xl border border-slate-700">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-slate-700 bg-slate-900/60">
            <th className="text-left px-4 py-3 text-xs font-semibold text-slate-400 uppercase tracking-wider">Field</th>
            <th className="text-left px-4 py-3 text-xs font-semibold text-slate-400 uppercase tracking-wider">Type</th>
            <th className="text-left px-4 py-3 text-xs font-semibold text-slate-400 uppercase tracking-wider hidden md:table-cell">Description</th>
            <th className="text-left px-4 py-3 text-xs font-semibold text-slate-400 uppercase tracking-wider">Required</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-700/50">
          {entries.map(([key, val]) => (
            <tr key={key} className="hover:bg-slate-800/30 transition-colors">
              <td className="px-4 py-3"><code className="text-primary-400 font-mono text-sm">{key}</code></td>
              <td className="px-4 py-3">
                <span className="text-xs bg-slate-700 text-slate-300 px-2 py-0.5 rounded font-mono uppercase">
                  {val?.type || val?.format || 'any'}
                </span>
              </td>
              <td className="px-4 py-3 text-slate-400 text-xs hidden md:table-cell">
                {val?.description || val?.['x-faker'] || '—'}
              </td>
              <td className="px-4 py-3">
                {required.includes(key)
                  ? <span className="text-xs text-emerald-400 font-medium">yes</span>
                  : <span className="text-xs text-slate-600">—</span>}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export function TemplateDetailClient({ id }: { id: string }) {
  const router = useRouter()

  const [template, setTemplate] = useState<MocklyTemplate | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getPublicTemplate(id)
      .then(r => {
        if (r.template.is_featured) {
          router.replace(`/resources/${r.template.slug}`)
          return
        }
        setTemplate(r.template)
      })
      .catch(e => setError(e.message))
      .finally(() => setLoading(false))
  }, [id, router])

  const schemaObj = (() => {
    if (!template?.schema_json) return null
    try { return JSON.parse(template.schema_json) } catch { return null }
  })()

  const schemaProperties: Record<string, any> = schemaObj?.properties ?? {}
  const schemaRequired: string[] = schemaObj?.required ?? []
  const endpointUrl = template ? `${PUBLIC_API_BASE}/t/${template.id}` : ''

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col">
      <Header />
      <main className="flex-1 container mx-auto px-4 py-10 max-w-5xl">
        <Link href="/templates" className="inline-flex items-center gap-1.5 text-sm text-slate-400 hover:text-white transition-colors mb-6">
          <ArrowLeft className="w-4 h-4" /> Back to Explore
        </Link>

        {loading && (
          <div className="flex items-center justify-center py-32">
            <Loader2 className="w-8 h-8 text-primary-400 animate-spin" />
          </div>
        )}

        {!loading && error && (
          <div className="flex items-center gap-2 text-red-400 bg-red-500/10 border border-red-500/20 rounded-xl p-4">
            <AlertCircle className="w-5 h-5 shrink-0" />
            <p>{error}</p>
          </div>
        )}

        {!loading && !error && template && (
          <div className="space-y-8">
            {/* Header */}
            <div className="bg-slate-900/50 border border-white/5 rounded-2xl p-6">
              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-xl bg-primary-500/10 border border-primary-500/20 flex items-center justify-center shrink-0">
                  <FileCode2 className="w-6 h-6 text-primary-400" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 flex-wrap mb-1">
                    <h1 className="text-2xl font-bold text-white">{template.name}</h1>
                    {template.is_featured && (
                      <span className="flex items-center gap-1 text-xs text-amber-400 bg-amber-500/10 border border-amber-500/20 px-2 py-0.5 rounded-full">
                        <Star className="w-3 h-3" /> Official
                      </span>
                    )}
                    {template.visibility === 'public' ? (
                      <span className="flex items-center gap-1 text-xs text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 rounded-full">
                        <Globe className="w-3 h-3" /> Public
                      </span>
                    ) : (
                      <span className="flex items-center gap-1 text-xs text-slate-400 bg-slate-700/50 border border-slate-600/50 px-2 py-0.5 rounded-full">
                        <Lock className="w-3 h-3" /> Private
                      </span>
                    )}
                  </div>
                  {template.description && (
                    <p className="text-slate-400 text-sm mt-1 max-w-2xl">{template.description}</p>
                  )}
                  <div className="flex items-center gap-4 mt-3 text-xs text-slate-500">
                    <span className="flex items-center gap-1">
                      <Hash className="w-3.5 h-3.5" />
                      <code className="font-mono">{template.slug}</code>
                    </span>
                    <span className="flex items-center gap-1">
                      <Calendar className="w-3.5 h-3.5" />
                      {formatDate(template.created_at)}
                    </span>
                  </div>
                </div>
                <div className="hidden sm:flex shrink-0">
                  <div className="flex items-center gap-2 bg-black/30 border border-slate-700/60 rounded-lg px-3 py-1.5">
                    <code className="text-xs font-mono text-slate-400 max-w-[240px] truncate">{endpointUrl}</code>
                    <a href={endpointUrl} target="_blank" rel="noopener noreferrer"
                      className="text-primary-400 hover:text-primary-300 transition-colors" title="Open endpoint (requires API key)">
                      <ExternalLink className="w-3.5 h-3.5" />
                    </a>
                  </div>
                </div>
              </div>
            </div>

            {/* Schema fields */}
            {Object.keys(schemaProperties).length > 0 && (
              <section>
                <h2 className="text-lg font-bold text-white mb-3">Schema Fields</h2>
                <SchemaFieldsTable properties={schemaProperties} required={schemaRequired} />
              </section>
            )}

            {/* Playground */}
            <section>
              <h2 className="text-lg font-bold text-white mb-3">Playground</h2>
              <TemplatePlayground templateId={template.id} schemaProperties={schemaProperties} />
            </section>
          </div>
        )}
      </main>
      <Footer />
    </div>
  )
}
