'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useAuthContext } from '@/components/AuthProvider'
import { listMyTemplates, getAccessKeyInfo } from '@/lib/mockly-api'
import type { MocklyTemplate, MocklyAPIKey } from '@/types/mockly'
import { FileCode2, Key, Plus, Globe, Lock, ArrowRight, Zap, BookOpen, ChevronRight } from 'lucide-react'

export default function DashboardPage() {
  const { token, user } = useAuthContext()
  const [templates, setTemplates] = useState<MocklyTemplate[]>([])
  const [accessKey, setAccessKey] = useState<MocklyAPIKey | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) return
    Promise.all([
      listMyTemplates(token).then(r => setTemplates(r.templates)),
      getAccessKeyInfo(token).then(k => setAccessKey(k)).catch(() => {}),
    ]).finally(() => setLoading(false))
  }, [token])

  const firstName = user?.name?.split(' ')[0] ?? 'there'

  return (
    <div className="p-6 max-w-4xl">

      {/* Welcome */}
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-white">Hey, {firstName} 👋</h1>
        <p className="text-slate-400 text-sm mt-1">Your mock API workspace.</p>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        {/* Templates */}
        <div className="relative overflow-hidden bg-slate-900/50 border border-slate-700 rounded-xl p-5 group hover:border-primary-500/30 transition-colors">
          <div className="absolute inset-0 bg-gradient-to-br from-primary-500/5 to-transparent pointer-events-none" />
          <div className="flex items-center gap-2.5 mb-4">
            <div className="w-8 h-8 rounded-lg bg-primary-500/10 border border-primary-500/20 flex items-center justify-center">
              <FileCode2 className="w-4 h-4 text-primary-400" />
            </div>
            <span className="text-sm font-medium text-slate-300">Templates</span>
          </div>
          <p className="text-4xl font-bold text-white tabular-nums">{loading ? '–' : templates.length}</p>
          <p className="text-xs text-slate-500 mt-1">of 10 free-tier slots</p>
          <Link href="/dashboard/templates/new" className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1 text-xs text-primary-400 hover:text-primary-300">
            New <ChevronRight className="w-3 h-3" />
          </Link>
        </div>

        {/* Access Key */}
        <div className="relative overflow-hidden bg-slate-900/50 border border-slate-700 rounded-xl p-5 group hover:border-emerald-500/30 transition-colors">
          <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent pointer-events-none" />
          <div className="flex items-center gap-2.5 mb-4">
            <div className="w-8 h-8 rounded-lg bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center">
              <Key className="w-4 h-4 text-emerald-400" />
            </div>
            <span className="text-sm font-medium text-slate-300">Access Key</span>
          </div>
          {accessKey ? (
            <p className="text-base font-mono text-emerald-300 truncate">{accessKey.key_prefix}…</p>
          ) : (
            <p className="text-sm text-slate-600">{loading ? '–' : 'Not generated'}</p>
          )}
          <p className="text-xs text-slate-500 mt-1">public <code className="text-slate-400">mak_</code> key</p>
          <Link href="/dashboard/api-keys" className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1 text-xs text-emerald-400 hover:text-emerald-300">
            Manage <ChevronRight className="w-3 h-3" />
          </Link>
        </div>
      </div>

      {/* Quick actions */}
      <div className="grid grid-cols-3 gap-3 mb-8">
        {[
          { href: '/dashboard/templates/new', icon: Plus, label: 'New Template', desc: 'Create a custom schema', color: 'text-primary-400', bg: 'bg-primary-500/10 border-primary-500/20 hover:border-primary-500/40' },
          { href: '/playground', icon: Zap, label: 'Playground', desc: 'Test & explore APIs', color: 'text-yellow-400', bg: 'bg-yellow-500/10 border-yellow-500/20 hover:border-yellow-500/40' },
          { href: '/templates', icon: BookOpen, label: 'Explore', desc: 'Browse public templates', color: 'text-purple-400', bg: 'bg-purple-500/10 border-purple-500/20 hover:border-purple-500/40' },
        ].map(({ href, icon: Icon, label, desc, color, bg }) => (
          <Link key={href} href={href} className={`flex items-center gap-3 p-4 rounded-xl border ${bg} transition-colors group`}>
            <div className={`w-8 h-8 rounded-lg ${bg.split(' ')[0]} flex items-center justify-center shrink-0`}>
              <Icon className={`w-4 h-4 ${color}`} />
            </div>
            <div className="min-w-0">
              <p className="text-sm font-semibold text-white group-hover:text-slate-100">{label}</p>
              <p className="text-xs text-slate-500 truncate">{desc}</p>
            </div>
          </Link>
        ))}
      </div>

      {/* Access Key usage */}
      {accessKey && (
        <div className="bg-slate-900/50 border border-slate-700 rounded-xl p-5 mb-6">
          <p className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3">Quick Start</p>
          <p className="text-xs text-slate-500 mb-2">Use your access key to call any public template:</p>
          <code className="block text-xs font-mono text-emerald-300 bg-black/30 rounded-lg px-3 py-2 break-all">
            GET https://api.mockly.codes/v1/<span className="text-primary-400">{'{slug}'}</span>?key={accessKey.key_prefix}…
          </code>
        </div>
      )}

      {/* Recent templates */}
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-sm font-semibold text-slate-300 uppercase tracking-wider">Recent Templates</h2>
        <Link href="/dashboard/templates/new" className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-primary-500/10 border border-primary-500/20 text-primary-400 text-xs font-medium hover:bg-primary-500/20 transition-colors">
          <Plus className="w-3.5 h-3.5" /> New
        </Link>
      </div>

      {loading ? (
        <div className="space-y-2">
          {[1,2,3].map(i => <div key={i} className="h-14 rounded-xl bg-slate-900/50 animate-pulse" />)}
        </div>
      ) : templates.length === 0 ? (
        <div className="relative overflow-hidden bg-slate-900/50 border border-dashed border-slate-700 rounded-xl p-10 text-center">
          <div className="absolute inset-0 bg-gradient-to-br from-primary-500/5 via-transparent to-purple-500/5 pointer-events-none" />
          <div className="w-12 h-12 rounded-xl bg-primary-500/10 border border-primary-500/20 flex items-center justify-center mx-auto mb-4">
            <FileCode2 className="w-6 h-6 text-primary-400" />
          </div>
          <p className="text-white font-semibold mb-1">No templates yet</p>
          <p className="text-slate-500 text-sm mb-5">Create your first schema-driven mock API endpoint.</p>
          <Link href="/dashboard/templates/new" className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-600 text-white text-sm font-medium hover:bg-primary-500 transition-colors">
            <Plus className="w-4 h-4" /> Create your first template
          </Link>
        </div>
      ) : (
        <div className="space-y-2">
          {templates.slice(0, 5).map(t => (
            <Link key={t.id} href={`/dashboard/templates/${t.id}`}
              className="flex items-center gap-4 bg-slate-900/50 border border-slate-700 rounded-xl p-4 hover:border-slate-600 hover:bg-slate-900/80 transition-all group"
            >
              <div className="w-8 h-8 rounded-lg bg-slate-800 border border-slate-700 flex items-center justify-center shrink-0">
                <FileCode2 className="w-4 h-4 text-slate-400 group-hover:text-primary-400 transition-colors" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-white group-hover:text-primary-300 transition-colors truncate">{t.name}</p>
                <p className="text-xs text-slate-500 font-mono truncate mt-0.5">/{t.slug}</p>
              </div>
              <div className="flex items-center gap-3 shrink-0">
                {t.visibility === 'public' ? (
                  <span className="flex items-center gap-1 text-xs text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 rounded-full">
                    <Globe className="w-3 h-3" /> Public
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-xs text-slate-400 bg-slate-500/10 border border-slate-500/20 px-2 py-0.5 rounded-full">
                    <Lock className="w-3 h-3" /> Private
                  </span>
                )}
                <ArrowRight className="w-4 h-4 text-slate-600 group-hover:text-slate-400 transition-colors" />
              </div>
            </Link>
          ))}
          {templates.length > 5 && (
            <Link href="/dashboard/templates" className="flex items-center justify-center gap-1 text-xs text-slate-500 hover:text-primary-400 py-3 transition-colors">
              View all {templates.length} templates <ArrowRight className="w-3 h-3" />
            </Link>
          )}
        </div>
      )}
    </div>
  )
}
