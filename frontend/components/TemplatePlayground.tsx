'use client'

import { useState, useEffect, useMemo } from 'react'
import {
  Play,
  Eye,
  EyeOff,
  Copy,
  Check,
  Key,
  Loader2,
  ChevronRight,
  Terminal,
  Code2,
  Clock,
  AlertCircle,
  CheckCircle2,
  ArrowRight,
} from 'lucide-react'
import { useAuthContext } from '@/components/AuthProvider'
import { listAPIKeys } from '@/lib/mockly-api'
import { getApiUrl } from '@/lib/api'
import { CodeExample } from '@/components/CodeExample'

const PUBLIC_API_BASE = 'https://api.mockly.codes'

// ── Types ─────────────────────────────────────────────────────────────────────

interface TemplatePlaygroundProps {
  templateId: string
  schemaProperties?: Record<string, { type?: string; description?: string; format?: string }>
}

type RequestTab = 'quick' | 'sort' | 'search' | 'fields'
type CodeTab = 'curl' | 'js' | 'python' | 'go' | 'node'

// ── Helper: format JSON ───────────────────────────────────────────────────────

function formatJSON(value: unknown): string {
  return JSON.stringify(value, null, 2)
}

// ── Helper: build params ──────────────────────────────────────────────────────

function buildParams(opts: {
  limit: number
  page: number
  sortField: string
  sortOrder: 'asc' | 'desc'
  searchQuery: string
  selectedFields: string[]
}): URLSearchParams {
  const p = new URLSearchParams()
  p.set('limit', String(opts.limit))
  p.set('page', String(opts.page))
  if (opts.sortField) {
    p.set('sort', opts.sortField)
    p.set('order', opts.sortOrder)
  }
  if (opts.searchQuery) p.set('q', opts.searchQuery)
  if (opts.selectedFields.length > 0) p.set('fields', opts.selectedFields.join(','))
  return p
}

// ── Code generators ───────────────────────────────────────────────────────────

function genCurl(url: string, apiKey: string): string {
  const k = apiKey || 'YOUR_API_KEY'
  return `curl -X GET \\
  "${url}" \\
  -H "Authorization: Bearer ${k}"`
}

function genJS(url: string, apiKey: string): string {
  const k = apiKey || 'YOUR_API_KEY'
  return `const response = await fetch('${url}', {
  headers: {
    'Authorization': 'Bearer ${k}'
  }
})

const data = await response.json()
console.log(data)`
}

function genPython(url: string, apiKey: string): string {
  const k = apiKey || 'YOUR_API_KEY'
  const [base, qs] = url.split('?')
  const params = Object.fromEntries(new URLSearchParams(qs).entries())
  const paramsStr = JSON.stringify(params, null, 4)
    .replace(/"/g, "'")
    .replace(/^\{/, '{')
  return `import requests

response = requests.get(
    '${base}',
    headers={'Authorization': 'Bearer ${k}'},
    params=${paramsStr}
)

data = response.json()
print(data)`
}

function genGo(url: string, apiKey: string): string {
  const k = apiKey || 'YOUR_API_KEY'
  return `package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    req, _ := http.NewRequest("GET", "${url}", nil)
    req.Header.Set("Authorization", "Bearer ${k}")

    resp, _ := http.DefaultClient.Do(req)
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}`
}

function genNode(url: string, apiKey: string): string {
  const k = apiKey || 'YOUR_API_KEY'
  return `const axios = require('axios')

const { data } = await axios.get('${url}', {
  headers: {
    Authorization: 'Bearer ${k}'
  }
})

console.log(data)`
}

// ── Component ─────────────────────────────────────────────────────────────────

export function TemplatePlayground({ templateId, schemaProperties }: TemplatePlaygroundProps) {
  const { token } = useAuthContext()

  // API key
  const [apiKey, setApiKey] = useState('')
  const [showKey, setShowKey] = useState(false)
  const [keyHints, setKeyHints] = useState<{ prefix: string; type: string }[]>([])

  // Request builder
  const [activeTab, setActiveTab] = useState<RequestTab>('quick')
  const [limit, setLimit] = useState(10)
  const [page, setPage] = useState(1)
  const [sortField, setSortField] = useState('')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedFields, setSelectedFields] = useState<string[]>([])

  // Response
  const [response, setResponse] = useState<unknown>(null)
  const [responseStatus, setResponseStatus] = useState<number | null>(null)
  const [responseTime, setResponseTime] = useState<number | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Code tab
  const [codeTab, setCodeTab] = useState<CodeTab>('curl')
  const [copied, setCopied] = useState(false)

  // Load key hints from user's account
  useEffect(() => {
    if (!token) return
    listAPIKeys(token)
      .then(({ keys }) => {
        const hints = keys
          .filter(k => k.status === 'active')
          .map(k => ({ prefix: k.key_prefix, type: k.type }))
        setKeyHints(hints)
      })
      .catch(() => {})
  }, [token])

  // Build URLs
  const params = useMemo(
    () => buildParams({ limit, page, sortField, sortOrder, searchQuery, selectedFields }),
    [limit, page, sortField, sortOrder, searchQuery, selectedFields],
  )

  const publicUrl = `${PUBLIC_API_BASE}/t/${templateId}?${params}`
  const localUrl = `${getApiUrl()}/t/${templateId}?${params}`

  // Code snippets (all dynamic)
  const code = useMemo(() => {
    const k = apiKey
    return {
      curl: genCurl(publicUrl, k),
      js: genJS(publicUrl, k),
      python: genPython(publicUrl, k),
      go: genGo(publicUrl, k),
      node: genNode(publicUrl, k),
    }
  }, [publicUrl, apiKey])

  const fieldNames = useMemo(() => Object.keys(schemaProperties ?? {}), [schemaProperties])

  // Send request
  const handleSend = async () => {
    if (!apiKey.trim()) {
      setError('An API key is required. Enter your personal (mk_…) or access key (mak_…).')
      return
    }
    setLoading(true)
    setError(null)
    setResponse(null)
    setResponseStatus(null)
    setResponseTime(null)

    const start = performance.now()
    try {
      const res = await fetch(localUrl, {
        headers: {
          Authorization: `Bearer ${apiKey.trim()}`,
          Accept: 'application/json',
        },
      })
      const elapsed = Math.round(performance.now() - start)
      setResponseStatus(res.status)
      setResponseTime(elapsed)

      const json = await res.json()
      if (!res.ok) {
        setError((json as any).error || `HTTP ${res.status}`)
      } else {
        setResponse(json)
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Request failed')
    } finally {
      setLoading(false)
    }
  }

  const handleCopyCode = async () => {
    await navigator.clipboard.writeText(code[codeTab])
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const toggleField = (f: string) => {
    setSelectedFields(prev =>
      prev.includes(f) ? prev.filter(x => x !== f) : [...prev, f],
    )
  }

  const REQUEST_TABS: { id: RequestTab; label: string }[] = [
    { id: 'quick', label: 'Quick' },
    { id: 'sort', label: 'Sort' },
    { id: 'search', label: 'Search' },
    ...(fieldNames.length > 0 ? [{ id: 'fields' as RequestTab, label: 'Fields' }] : []),
  ]

  const CODE_TABS: { id: CodeTab; label: string; lang: string }[] = [
    { id: 'curl', label: 'cURL', lang: 'bash' },
    { id: 'js', label: 'JavaScript', lang: 'javascript' },
    { id: 'python', label: 'Python', lang: 'python' },
    { id: 'go', label: 'Go', lang: 'go' },
    { id: 'node', label: 'Node.js', lang: 'javascript' },
  ]

  return (
    <div className="space-y-4">
      {/* ── URL preview bar ───────────────────────────────────────────── */}
      <div className="bg-slate-900/80 rounded-xl border border-slate-700 overflow-hidden">
        <div className="flex items-center gap-3 px-4 py-3 border-b border-slate-700/60">
          <span className="shrink-0 text-xs font-bold text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 rounded font-mono">
            GET
          </span>
          <code className="flex-1 text-xs text-slate-300 font-mono truncate">{publicUrl}</code>
        </div>
        <div className="flex items-center gap-3 px-4 py-2.5">
          <span className="shrink-0 text-xs text-slate-500 font-mono w-36">Authorization:</span>
          <span className="text-xs text-slate-500 font-mono">Bearer</span>
          <span className="text-xs font-mono text-slate-400 truncate">
            {apiKey ? (showKey ? apiKey : '•'.repeat(Math.min(apiKey.length, 24))) : <span className="text-slate-600 italic">not set</span>}
          </span>
        </div>
      </div>

      {/* ── API key input ─────────────────────────────────────────────── */}
      <div className="bg-slate-800/40 rounded-xl border border-slate-700 p-4 space-y-3">
        <div className="flex items-center gap-2">
          <Key className="w-4 h-4 text-primary-400 shrink-0" />
          <span className="text-sm font-semibold text-white">API Key</span>
          <span className="text-xs text-slate-500 ml-auto">Required for all requests</span>
        </div>

        <div className="flex items-center gap-2">
          <div className="flex-1 relative">
            <input
              type={showKey ? 'text' : 'password'}
              value={apiKey}
              onChange={e => setApiKey(e.target.value)}
              placeholder="mk_… or mak_…"
              autoComplete="off"
              className="w-full bg-slate-900/60 border border-slate-600 rounded-lg px-3 py-2 text-sm text-white font-mono placeholder-slate-600 focus:outline-none focus:border-primary-500 transition-colors"
            />
          </div>
          <button
            onClick={() => setShowKey(v => !v)}
            className="p-2 rounded-lg bg-slate-700 hover:bg-slate-600 text-slate-400 transition-colors"
            title={showKey ? 'Hide key' : 'Show key'}
          >
            {showKey ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
          </button>
        </div>

        {/* Key hints (logged-in users) */}
        {keyHints.length > 0 && (
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-xs text-slate-500">Your keys:</span>
            {keyHints.map(k => (
              <span
                key={k.prefix}
                className="text-xs font-mono text-slate-400 bg-slate-900/60 border border-slate-700 px-2 py-0.5 rounded"
              >
                {k.prefix}… <span className="text-slate-600">({k.type})</span>
              </span>
            ))}
          </div>
        )}

        {!token && (
          <p className="text-xs text-slate-500">
            <a href="/dashboard" className="text-primary-400 hover:text-primary-300 underline-offset-2 hover:underline">
              Sign in
            </a>{' '}
            to see your API keys, or enter one manually.
          </p>
        )}
      </div>

      {/* ── Two-column: request builder | response ────────────────────── */}
      <div className="grid lg:grid-cols-2 gap-4">
        {/* Left: request builder */}
        <div className="bg-slate-800/40 rounded-xl border border-slate-700 flex flex-col">
          {/* Tab bar */}
          <div className="flex items-center gap-1 px-4 pt-4 pb-0">
            {REQUEST_TABS.map(t => (
              <button
                key={t.id}
                onClick={() => setActiveTab(t.id)}
                className={`px-3 py-1.5 rounded-t text-xs font-medium transition-colors ${
                  activeTab === t.id
                    ? 'bg-slate-700 text-white'
                    : 'text-slate-500 hover:text-slate-300'
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>

          <div className="p-4 flex-1 space-y-4">
            {activeTab === 'quick' && (
              <div className="space-y-3">
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs text-slate-400 mb-1">Limit</label>
                    <input
                      type="number"
                      min={1}
                      max={100}
                      value={limit}
                      onChange={e => setLimit(Math.max(1, Math.min(100, +e.target.value)))}
                      className="w-full bg-slate-900/60 border border-slate-600 rounded-lg px-3 py-1.5 text-sm text-white focus:outline-none focus:border-primary-500"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-slate-400 mb-1">Page</label>
                    <input
                      type="number"
                      min={1}
                      value={page}
                      onChange={e => setPage(Math.max(1, +e.target.value))}
                      className="w-full bg-slate-900/60 border border-slate-600 rounded-lg px-3 py-1.5 text-sm text-white focus:outline-none focus:border-primary-500"
                    />
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'sort' && (
              <div className="space-y-3">
                <div>
                  <label className="block text-xs text-slate-400 mb-1">Sort field</label>
                  {fieldNames.length > 0 ? (
                    <select
                      value={sortField}
                      onChange={e => setSortField(e.target.value)}
                      className="w-full bg-slate-900/60 border border-slate-600 rounded-lg px-3 py-1.5 text-sm text-white focus:outline-none focus:border-primary-500"
                    >
                      <option value="">— none —</option>
                      {fieldNames.map(f => (
                        <option key={f} value={f}>{f}</option>
                      ))}
                    </select>
                  ) : (
                    <input
                      type="text"
                      placeholder="field_name"
                      value={sortField}
                      onChange={e => setSortField(e.target.value)}
                      className="w-full bg-slate-900/60 border border-slate-600 rounded-lg px-3 py-1.5 text-sm text-white focus:outline-none focus:border-primary-500"
                    />
                  )}
                </div>
                <div>
                  <label className="block text-xs text-slate-400 mb-1">Order</label>
                  <div className="flex gap-2">
                    {(['asc', 'desc'] as const).map(o => (
                      <button
                        key={o}
                        onClick={() => setSortOrder(o)}
                        className={`flex-1 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                          sortOrder === o
                            ? 'bg-primary-500 text-white'
                            : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                        }`}
                      >
                        {o.toUpperCase()}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'search' && (
              <div>
                <label className="block text-xs text-slate-400 mb-1">Search query</label>
                <input
                  type="text"
                  placeholder="e.g. John"
                  value={searchQuery}
                  onChange={e => setSearchQuery(e.target.value)}
                  className="w-full bg-slate-900/60 border border-slate-600 rounded-lg px-3 py-1.5 text-sm text-white focus:outline-none focus:border-primary-500"
                />
                <p className="text-xs text-slate-500 mt-1">Sent as <code className="text-slate-400">?q=…</code></p>
              </div>
            )}

            {activeTab === 'fields' && fieldNames.length > 0 && (
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <label className="text-xs text-slate-400">Fields to return</label>
                  <button
                    onClick={() => setSelectedFields(prev => prev.length === fieldNames.length ? [] : fieldNames)}
                    className="text-xs text-primary-400 hover:text-primary-300 transition-colors"
                  >
                    {selectedFields.length === fieldNames.length ? 'Deselect all' : 'Select all'}
                  </button>
                </div>
                <div className="grid grid-cols-2 gap-1.5 max-h-40 overflow-y-auto">
                  {fieldNames.map(f => (
                    <label key={f} className="flex items-center gap-2 cursor-pointer group">
                      <input
                        type="checkbox"
                        checked={selectedFields.includes(f)}
                        onChange={() => toggleField(f)}
                        className="rounded border-slate-600 bg-slate-900 text-primary-500 focus:ring-primary-500"
                      />
                      <span className="text-xs text-slate-300 font-mono group-hover:text-white truncate">{f}</span>
                    </label>
                  ))}
                </div>
                {selectedFields.length > 0 && (
                  <p className="text-xs text-slate-500">
                    <code className="text-slate-400">fields={selectedFields.join(',')}</code>
                  </p>
                )}
              </div>
            )}
          </div>

          {/* Send button */}
          <div className="p-4 border-t border-slate-700">
            <button
              onClick={handleSend}
              disabled={loading}
              className="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-primary-600 hover:bg-primary-500 disabled:opacity-50 text-white text-sm font-semibold transition-colors shadow-lg shadow-primary-500/10"
            >
              {loading ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  Sending…
                </>
              ) : (
                <>
                  <Play className="w-4 h-4" />
                  Send Request
                </>
              )}
            </button>
          </div>
        </div>

        {/* Right: response */}
        <div className="bg-slate-800/40 rounded-xl border border-slate-700 flex flex-col min-h-[300px]">
          <div className="flex items-center gap-3 px-4 py-3 border-b border-slate-700/60">
            <span className="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
              <Terminal className="w-3.5 h-3.5 text-primary-400" /> Response
            </span>
            {responseStatus !== null && (
              <div className="ml-auto flex items-center gap-2">
                <span
                  className={`text-xs font-bold px-2 py-0.5 rounded font-mono ${
                    responseStatus < 300
                      ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                      : 'bg-red-500/10 text-red-400 border border-red-500/20'
                  }`}
                >
                  {responseStatus}
                </span>
                {responseTime !== null && (
                  <span className="text-xs text-slate-500 flex items-center gap-1">
                    <Clock className="w-3 h-3" /> {responseTime}ms
                  </span>
                )}
              </div>
            )}
          </div>

          <div className="flex-1 overflow-auto p-4">
            {loading && (
              <div className="flex items-center justify-center h-32">
                <Loader2 className="w-8 h-8 text-primary-400 animate-spin" />
              </div>
            )}
            {!loading && error && (
              <div className="flex items-start gap-2 text-red-400">
                <AlertCircle className="w-4 h-4 mt-0.5 shrink-0" />
                <p className="text-sm">{error}</p>
              </div>
            )}
            {!loading && !error && response === null && (
              <div className="flex flex-col items-center justify-center h-32 text-slate-600 text-sm">
                <ArrowRight className="w-6 h-6 mb-2" />
                Hit <strong className="text-slate-500 mx-1">Send Request</strong> to see the response
              </div>
            )}
            {!loading && !error && response !== null && (
              <pre className="text-xs text-slate-300 font-mono whitespace-pre-wrap break-words leading-relaxed">
                {formatJSON(response)}
              </pre>
            )}
          </div>
        </div>
      </div>

      {/* ── Code examples ─────────────────────────────────────────────── */}
      <div className="bg-slate-800/40 rounded-xl border border-slate-700 overflow-hidden">
        {/* Tab bar */}
        <div className="flex items-center gap-1 px-4 py-3 border-b border-slate-700/60 bg-slate-900/40">
          <Code2 className="w-4 h-4 text-primary-400 mr-1" />
          {CODE_TABS.map(t => (
            <button
              key={t.id}
              onClick={() => setCodeTab(t.id)}
              className={`px-3 py-1 rounded text-xs font-medium transition-colors ${
                codeTab === t.id
                  ? 'bg-primary-500/20 text-primary-300 border border-primary-500/30'
                  : 'text-slate-500 hover:text-slate-300'
              }`}
            >
              {t.label}
            </button>
          ))}
          <button
            onClick={handleCopyCode}
            className="ml-auto flex items-center gap-1.5 text-xs text-slate-400 hover:text-white bg-slate-700 hover:bg-slate-600 px-3 py-1 rounded transition-colors"
          >
            {copied ? <Check className="w-3 h-3" /> : <Copy className="w-3 h-3" />}
            {copied ? 'Copied!' : 'Copy'}
          </button>
        </div>

        <pre className="p-4 overflow-x-auto text-xs text-slate-300 font-mono leading-relaxed">
          <code>{code[codeTab]}</code>
        </pre>
      </div>
    </div>
  )
}
