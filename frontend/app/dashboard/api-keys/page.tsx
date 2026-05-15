'use client'

import { useEffect, useState } from 'react'
import { useAuthContext } from '@/components/AuthProvider'
import { listAPIKeys, createAPIKey, revokeAPIKey, getAccessKeyInfo, regenerateAccessKey } from '@/lib/mockly-api'
import type { MocklyAPIKey } from '@/types/mockly'
import { Copy, Trash2, RefreshCw, Plus, Key, Globe, ShieldCheck, Check } from 'lucide-react'

function KeyBadge({ type }: { type: string }) {
  return type === 'access'
    ? <span className="text-xs bg-emerald-500/10 text-emerald-400 px-2 py-0.5 rounded-full border border-emerald-500/20">mak_</span>
    : <span className="text-xs bg-primary-500/10 text-primary-400 px-2 py-0.5 rounded-full border border-primary-500/20">mk_</span>
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)
  function copy() {
    navigator.clipboard?.writeText(text).catch(() => {})
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }
  return (
    <button onClick={copy} className="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-white/5 transition-colors">
      {copied ? <Check className="w-4 h-4 text-emerald-400" /> : <Copy className="w-4 h-4" />}
    </button>
  )
}

export default function APIKeysPage() {
  const { token } = useAuthContext()
  const [keys, setKeys] = useState<MocklyAPIKey[]>([])
  const [accessKey, setAccessKey] = useState<MocklyAPIKey | null>(null)
  const [loading, setLoading] = useState(true)
  const [newKeyName, setNewKeyName] = useState('')
  const [creating, setCreating] = useState(false)
  const [newRawKey, setNewRawKey] = useState<string | null>(null)
  const [newAccessRaw, setNewAccessRaw] = useState<string | null>(null)
  const [showRaw, setShowRaw] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [revokingId, setRevokingId] = useState<string | null>(null)

  async function load() {
    if (!token) return
    const [keysRes] = await Promise.all([
      listAPIKeys(token),
      getAccessKeyInfo(token).then(k => setAccessKey(k)).catch(() => {}),
    ])
    setKeys(keysRes.keys.filter(k => k.type === 'personal' && k.status === 'active'))
  }

  useEffect(() => {
    load().finally(() => setLoading(false))
  }, [token])

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!token || !newKeyName.trim()) return
    setCreating(true)
    setError(null)
    try {
      const { raw_key } = await createAPIKey(token, newKeyName.trim())
      setNewRawKey(raw_key)
      setNewKeyName('')
      await load()
    } catch (e: any) {
      setError(e?.message || 'Failed to create key')
    } finally {
      setCreating(false)
    }
  }

  async function handleRevoke(id: string) {
    if (!token) return
    setRevokingId(id)
    await revokeAPIKey(token, id).catch(() => {})
    setKeys(prev => prev.filter(k => k.id !== id))
    setRevokingId(null)
  }

  async function handleRegenerateAccess() {
    if (!token) return
    try {
      const { raw_key, key } = await regenerateAccessKey(token)
      setNewAccessRaw(raw_key)
      setAccessKey(key)
    } catch (e: any) {
      setError(e?.message || 'Failed to regenerate access key')
    }
  }

  return (
    <div className="p-8 max-w-3xl">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <div className="w-9 h-9 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center">
            <Key className="w-4 h-4 text-amber-400" />
          </div>
          <h1 className="text-2xl font-bold text-white">API Keys</h1>
        </div>
        <p className="text-slate-400 text-sm ml-12">Manage authentication keys for your templates.</p>
      </div>

      {error && (
        <div className="mb-6 p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm">{error}</div>
      )}

      {/* One-time display for newly created key */}
      {newRawKey && (
        <div className="mb-6 p-4 rounded-xl bg-emerald-900/20 border border-emerald-500/20">
          <div className="flex items-center gap-2 mb-2">
            <ShieldCheck className="w-4 h-4 text-emerald-400" />
            <p className="text-xs text-emerald-400 font-medium">New personal key — copy it now, it won&apos;t be shown again</p>
          </div>
          <div className="flex items-center gap-2">
            <code className="flex-1 text-sm font-mono text-white bg-black/30 px-3 py-2 rounded-lg truncate border border-emerald-500/10">{newRawKey}</code>
            <CopyButton text={newRawKey} />
          </div>
          <button onClick={() => setNewRawKey(null)} className="mt-2 text-xs text-slate-500 hover:text-slate-300 transition-colors">Dismiss</button>
        </div>
      )}

      {/* Access key section */}
      <section className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <div className="w-7 h-7 rounded-lg bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center">
            <Globe className="w-3.5 h-3.5 text-emerald-400" />
          </div>
          <h2 className="text-sm font-semibold text-white">Public Access Key</h2>
          <KeyBadge type="access" />
        </div>
        <p className="text-xs text-slate-500 mb-4 ml-10">Share this key to let others access your public templates. Regenerating invalidates the previous key immediately.</p>

        {newAccessRaw && (
          <div className="mb-4 p-4 rounded-xl bg-emerald-900/20 border border-emerald-500/20">
            <div className="flex items-center gap-2 mb-2">
              <ShieldCheck className="w-4 h-4 text-emerald-400" />
              <p className="text-xs text-emerald-400 font-medium">New access key — copy it now</p>
            </div>
            <div className="flex items-center gap-2">
              <code className="flex-1 text-sm font-mono text-white bg-black/30 px-3 py-2 rounded-lg truncate border border-emerald-500/10">{newAccessRaw}</code>
              <CopyButton text={newAccessRaw} />
            </div>
            <button onClick={() => setNewAccessRaw(null)} className="mt-2 text-xs text-slate-500 hover:text-slate-300 transition-colors">Dismiss</button>
          </div>
        )}

        {loading ? (
          <div className="h-16 rounded-xl bg-slate-900/50 border border-slate-700 animate-pulse" />
        ) : accessKey ? (
          <div className="group relative bg-slate-900/50 border border-slate-700 rounded-xl p-4 hover:border-emerald-500/30 transition-all">
            <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-emerald-500/3 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
            <div className="flex items-center gap-3">
              <code className="flex-1 text-sm font-mono text-slate-300 min-w-0 truncate">{accessKey.key_prefix}{'•'.repeat(24)}</code>
              <button
                onClick={handleRegenerateAccess}
                className="flex items-center gap-1.5 text-xs text-slate-400 hover:text-white bg-white/5 hover:bg-white/10 px-3 py-1.5 rounded-lg transition-colors shrink-0"
              >
                <RefreshCw className="w-3.5 h-3.5" /> Regenerate
              </button>
            </div>
          </div>
        ) : (
          <div className="text-slate-500 text-sm p-4 rounded-xl bg-slate-900/30 border border-slate-800">No access key found.</div>
        )}
      </section>

      {/* Personal keys */}
      <section>
        <div className="flex items-center gap-3 mb-2">
          <div className="w-7 h-7 rounded-lg bg-primary-500/10 border border-primary-500/20 flex items-center justify-center">
            <Key className="w-3.5 h-3.5 text-primary-400" />
          </div>
          <h2 className="text-sm font-semibold text-white">Personal Keys</h2>
          <KeyBadge type="personal" />
        </div>
        <p className="text-xs text-slate-500 mb-5 ml-10">For accessing your own private templates from your apps.</p>

        {/* Create new key */}
        <form onSubmit={handleCreate} className="flex gap-2 mb-5">
          <input
            type="text" value={newKeyName} onChange={e => setNewKeyName(e.target.value)}
            placeholder="Key name (e.g. my-app)"
            className="flex-1 px-3 py-2.5 rounded-lg bg-slate-900/60 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:border-primary-500 transition-colors"
          />
          <button
            type="submit" disabled={creating || !newKeyName.trim()}
            className="flex items-center gap-1.5 px-4 py-2.5 rounded-lg bg-primary-600 text-white text-sm font-medium hover:bg-primary-500 transition-colors disabled:opacity-60"
          >
            <Plus className="w-4 h-4" /> {creating ? 'Creating…' : 'Create'}
          </button>
        </form>

        {loading ? (
          <div className="space-y-2">
            {[1, 2].map(i => <div key={i} className="h-16 rounded-xl bg-slate-900/50 border border-slate-700 animate-pulse" />)}
          </div>
        ) : keys.length === 0 ? (
          <div className="text-center p-8 rounded-xl bg-slate-900/30 border border-dashed border-slate-700">
            <Key className="w-8 h-8 text-slate-600 mx-auto mb-2" />
            <p className="text-slate-400 text-sm font-medium mb-1">No personal keys yet</p>
            <p className="text-slate-600 text-xs">Create a key above to start calling your private templates.</p>
          </div>
        ) : (
          <div className="space-y-2">
            {keys.map(k => (
              <div key={k.id} className="group relative bg-slate-900/50 border border-slate-700 rounded-xl p-4 hover:border-primary-500/20 transition-all flex items-center gap-4">
                <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-primary-500/3 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
                <div className="w-8 h-8 rounded-lg bg-primary-500/10 border border-primary-500/10 flex items-center justify-center shrink-0">
                  <Key className="w-3.5 h-3.5 text-primary-400" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-white truncate">{k.name}</p>
                  <code className="text-xs font-mono text-slate-500">{k.key_prefix}{'•'.repeat(24)}</code>
                </div>
                <button
                  onClick={() => handleRevoke(k.id)}
                  disabled={revokingId === k.id}
                  className="p-2 text-slate-500 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors disabled:opacity-50 shrink-0"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
