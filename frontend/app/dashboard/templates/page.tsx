'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useAuthContext } from '@/components/AuthProvider'
import { listMyTemplates, deleteTemplate } from '@/lib/mockly-api'
import type { MocklyTemplate } from '@/types/mockly'
import { Plus, Globe, Lock, Trash2, Pencil, FileCode2, ArrowRight } from 'lucide-react'

const API_BASE = 'https://api.mockly.codes/v1'

export default function TemplatesPage() {
  const { token } = useAuthContext()
  const [templates, setTemplates] = useState<MocklyTemplate[]>([])
  const [loading, setLoading] = useState(true)
  const [deletingId, setDeletingId] = useState<string | null>(null)

  async function load() {
    if (!token) return
    const { templates } = await listMyTemplates(token)
    setTemplates(templates)
  }

  useEffect(() => {
    load().finally(() => setLoading(false))
  }, [token])

  async function handleDelete(id: string) {
    if (!token) return
    setDeletingId(id)
    await deleteTemplate(token, id).catch(() => {})
    setTemplates(prev => prev.filter(t => t.id !== id))
    setDeletingId(null)
  }

  return (
    <div className="p-8 max-w-4xl">
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl bg-violet-500/10 border border-violet-500/20 flex items-center justify-center">
            <FileCode2 className="w-4 h-4 text-violet-400" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">My Templates</h1>
            <p className="text-slate-500 text-xs mt-0.5">{loading ? '…' : `${templates.length} / 10 used`}</p>
          </div>
        </div>
        <Link
          href="/dashboard/templates/new"
          className="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-primary-600 text-white text-sm font-medium hover:bg-primary-500 transition-colors"
        >
          <Plus className="w-4 h-4" /> New Template
        </Link>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[1, 2, 3].map(i => <div key={i} className="h-20 rounded-xl bg-slate-900/50 border border-slate-700 animate-pulse" />)}
        </div>
      ) : templates.length === 0 ? (
        <div className="relative overflow-hidden bg-slate-900/50 border border-dashed border-slate-700 rounded-2xl p-12 text-center">
          <div className="absolute inset-0 bg-gradient-to-br from-violet-500/5 via-transparent to-primary-500/5 pointer-events-none" />
          <div className="w-12 h-12 rounded-2xl bg-violet-500/10 border border-violet-500/20 flex items-center justify-center mx-auto mb-4">
            <FileCode2 className="w-6 h-6 text-violet-400" />
          </div>
          <p className="text-white font-semibold mb-1">No templates yet</p>
          <p className="text-slate-500 text-sm mb-5">Create your first mock API schema to start generating data.</p>
          <Link
            href="/dashboard/templates/new"
            className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-primary-600 text-white text-sm font-medium hover:bg-primary-500 transition-colors"
          >
            <Plus className="w-4 h-4" /> Create your first template
          </Link>
        </div>
      ) : (
        <div className="space-y-2">
          {templates.map(t => (
            <div
              key={t.id}
              className="group relative bg-slate-900/50 border border-slate-700 rounded-xl p-4 flex items-center gap-4 hover:border-violet-500/30 transition-all"
            >
              <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-violet-500/3 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />

              <div className="w-9 h-9 rounded-lg bg-violet-500/10 border border-violet-500/10 flex items-center justify-center shrink-0">
                <FileCode2 className="w-4 h-4 text-violet-400" />
              </div>

              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-0.5">
                  <p className="text-sm font-medium text-white truncate group-hover:text-violet-300 transition-colors">{t.name}</p>
                  {t.visibility === 'public' ? (
                    <span className="flex items-center gap-1 text-xs text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 rounded-full shrink-0">
                      <Globe className="w-3 h-3" /> Public
                    </span>
                  ) : (
                    <span className="flex items-center gap-1 text-xs text-slate-400 bg-slate-500/10 border border-slate-600/30 px-2 py-0.5 rounded-full shrink-0">
                      <Lock className="w-3 h-3" /> Private
                    </span>
                  )}
                </div>
                <p className="text-xs text-slate-500 font-mono">{API_BASE}/{t.slug}</p>
              </div>

              <div className="flex items-center gap-1.5 shrink-0">
                <Link
                  href={`/dashboard/templates/${t.id}`}
                  className="flex items-center gap-1.5 p-2 text-slate-400 hover:text-violet-300 hover:bg-violet-500/10 rounded-lg transition-colors"
                >
                  <Pencil className="w-4 h-4" />
                </Link>
                <a
                  href={`${API_BASE}/${t.slug}?limit=3`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-1.5 p-2 text-slate-400 hover:text-primary-300 hover:bg-primary-500/10 rounded-lg transition-colors"
                >
                  <ArrowRight className="w-4 h-4" />
                </a>
                <button
                  onClick={() => handleDelete(t.id)}
                  disabled={deletingId === t.id}
                  className="p-2 text-slate-500 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors disabled:opacity-50"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
