'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthContext } from '@/components/AuthProvider'
import { createTemplate } from '@/lib/mockly-api'
import { schemasManifest } from '@/lib/schemas-manifest'
import type { TemplateVisibility, TemplateType } from '@/types/mockly'
import { ChevronDown, Globe, Lock, FileCode2, Zap, Copy, Check } from 'lucide-react'

const SLUG_RE = /^[a-z0-9-]{1,60}$/
const API_BASE = 'https://api.mockly.codes/v1'

const GENERATORS = [
  { annotation: 'x-generator: name.full', description: 'Full name', example: 'Alice Johnson' },
  { annotation: 'x-generator: name.first', description: 'First name', example: 'Alice' },
  { annotation: 'x-generator: email', description: 'Email address', example: 'alice@example.com' },
  { annotation: 'x-generator: internet.uuid', description: 'UUID v4', example: '550e8400-e29b...' },
  { annotation: 'x-generator: number.int', description: 'Random integer', example: '42' },
  { annotation: 'x-generator: date.past', description: 'Past ISO date', example: '2023-04-12T...' },
  { annotation: 'x-generator: lorem.sentence', description: 'Lorem sentence', example: 'Lorem ipsum...' },
  { annotation: 'x-generator: internet.url', description: 'URL', example: 'https://...' },
  { annotation: 'x-generator: address.city', description: 'City name', example: 'San Francisco' },
  { annotation: 'x-generator: color.rgb', description: 'Hex colour', example: '#4a90e2' },
]

const EXAMPLE_SCHEMA = `{
  "properties": {
    "id": {
      "type": "string",
      "x-generator": "internet.uuid"
    },
    "name": {
      "type": "string",
      "x-generator": "name.full"
    },
    "email": {
      "type": "string",
      "x-generator": "email"
    },
    "createdAt": {
      "type": "string",
      "x-generator": "date.past"
    }
  }
}`

function schemaToJSON(schema: Record<string, any>): string {
  return JSON.stringify(schema, null, 2)
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)
  function copy() {
    navigator.clipboard.writeText(text).catch(() => {})
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }
  return (
    <button onClick={copy} className="p-1 rounded text-slate-500 hover:text-slate-300 transition-colors">
      {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
    </button>
  )
}

export default function NewTemplatePage() {
  const { token } = useAuthContext()
  const router = useRouter()

  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [description, setDescription] = useState('')
  const [visibility, setVisibility] = useState<TemplateVisibility>('private')
  const [type, setType] = useState<TemplateType>('custom')
  const [baseSchema, setBaseSchema] = useState('')
  const [schemaJSON, setSchemaJSON] = useState(EXAMPLE_SCHEMA)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  function handleTypeChange(t: TemplateType) {
    setType(t)
    if (t === 'fork' && baseSchema) {
      const found = schemasManifest.find(s => s.name === baseSchema)
      if (found) setSchemaJSON(schemaToJSON(found.schema))
    } else if (t === 'custom') {
      setSchemaJSON(EXAMPLE_SCHEMA)
    }
  }

  function handleBaseSchemaChange(n: string) {
    setBaseSchema(n)
    const found = schemasManifest.find(s => s.name === n)
    if (found) setSchemaJSON(schemaToJSON(found.schema))
  }

  function toSlug(s: string) {
    return s.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '').slice(0, 60)
  }

  function handleNameChange(v: string) {
    setName(v)
    if (!slug || slug === toSlug(name)) setSlug(toSlug(v))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!token) return
    setError(null)
    if (!SLUG_RE.test(slug)) {
      setError('Slug must be 1–60 lowercase letters, numbers, or hyphens.')
      return
    }
    try { JSON.parse(schemaJSON) } catch {
      setError('Schema JSON is not valid JSON.')
      return
    }
    if (type === 'fork' && !baseSchema) {
      setError('Please select a base schema to fork from.')
      return
    }
    setSaving(true)
    try {
      await createTemplate(token, { name, slug, description, visibility, type, base_schema: type === 'fork' ? baseSchema : undefined, schema_json: schemaJSON })
      router.push('/dashboard/templates')
    } catch (e: any) {
      setError(e?.message || 'Failed to create template')
    } finally {
      setSaving(false)
    }
  }

  const apiUrl = slug ? `${API_BASE}/${slug}` : `${API_BASE}/{slug}`
  const curlExample = `curl "${apiUrl}?limit=5"`

  return (
    <div className="flex min-h-full">
      {/* ── Left: form ─────────────────────────────────────── */}
      <div className="flex-1 p-6 min-w-0 max-w-2xl">
        <div className="mb-6">
          <h1 className="text-xl font-bold text-white">New Template</h1>
          <p className="text-slate-400 text-sm mt-1">Define a custom schema to generate mock data.</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-sm">{error}</div>
          )}

          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2">
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Name</label>
              <input type="text" required value={name} onChange={e => handleNameChange(e.target.value)}
                placeholder="My Users Schema"
                className="w-full px-3 py-2 rounded-lg bg-slate-900/60 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:border-primary-500 transition-colors"
              />
            </div>

            <div className="col-span-2">
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Slug</label>
              <input type="text" required value={slug} onChange={e => setSlug(toSlug(e.target.value))}
                placeholder="my-users-schema"
                className="w-full px-3 py-2 rounded-lg bg-slate-900/60 border border-slate-700 text-white placeholder-slate-500 text-sm font-mono focus:outline-none focus:border-primary-500 transition-colors"
              />
              <p className="text-xs text-slate-500 mt-1">Lowercase letters, numbers, hyphens only.</p>
            </div>

            <div className="col-span-2">
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Description <span className="text-slate-600">(optional)</span></label>
              <input type="text" value={description} onChange={e => setDescription(e.target.value)}
                placeholder="A schema for user profiles"
                className="w-full px-3 py-2 rounded-lg bg-slate-900/60 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:border-primary-500 transition-colors"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Visibility</label>
              <div className="relative">
                <select value={visibility} onChange={e => setVisibility(e.target.value as TemplateVisibility)}
                  className="w-full appearance-none px-3 py-2 pr-8 rounded-lg bg-slate-900/60 border border-slate-700 text-white text-sm focus:outline-none focus:border-primary-500 transition-colors cursor-pointer"
                >
                  <option value="private">Private</option>
                  <option value="public">Public</option>
                </select>
                <ChevronDown className="absolute right-2 top-2.5 w-4 h-4 text-slate-500 pointer-events-none" />
              </div>
            </div>

            <div>
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Type</label>
              <div className="relative">
                <select value={type} onChange={e => handleTypeChange(e.target.value as TemplateType)}
                  className="w-full appearance-none px-3 py-2 pr-8 rounded-lg bg-slate-900/60 border border-slate-700 text-white text-sm focus:outline-none focus:border-primary-500 transition-colors cursor-pointer"
                >
                  <option value="custom">Custom</option>
                  <option value="fork">Fork built-in</option>
                </select>
                <ChevronDown className="absolute right-2 top-2.5 w-4 h-4 text-slate-500 pointer-events-none" />
              </div>
            </div>
          </div>

          {type === 'fork' && (
            <div>
              <label className="block text-xs font-medium text-slate-400 mb-1.5">Base Schema</label>
              <div className="relative">
                <select value={baseSchema} onChange={e => handleBaseSchemaChange(e.target.value)} required
                  className="w-full appearance-none px-3 py-2 pr-8 rounded-lg bg-slate-900/60 border border-slate-700 text-white text-sm focus:outline-none focus:border-primary-500 transition-colors cursor-pointer"
                >
                  <option value="">Select a built-in schema…</option>
                  {schemasManifest.map(s => (
                    <option key={s.name} value={s.name}>{s.name} ({s.group})</option>
                  ))}
                </select>
                <ChevronDown className="absolute right-2 top-2.5 w-4 h-4 text-slate-500 pointer-events-none" />
              </div>
            </div>
          )}

          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1.5">Schema JSON</label>
            <textarea
              value={schemaJSON} onChange={e => setSchemaJSON(e.target.value)}
              rows={18} spellCheck={false}
              className="w-full px-3 py-2 rounded-lg bg-slate-900/60 border border-slate-700 text-white text-xs font-mono focus:outline-none focus:border-primary-500 transition-colors resize-y"
            />
          </div>

          <div className="flex items-center gap-3 pt-2">
            <button type="submit" disabled={saving}
              className="px-5 py-2 rounded-lg bg-primary-600 text-white text-sm font-medium hover:bg-primary-500 transition-colors disabled:opacity-60"
            >
              {saving ? 'Creating…' : 'Create Template'}
            </button>
            <button type="button" onClick={() => router.back()}
              className="px-5 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-white/5 text-sm transition-colors"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>

      {/* ── Right: live preview + docs ──────────────────────── */}
      <div className="flex-1 min-w-0 border-l border-slate-700 bg-slate-900/30 p-6 overflow-y-auto flex flex-col gap-6">

        {/* Live URL preview */}
        <div className="mb-6">
          <div className="flex items-center gap-2 mb-3">
            <div className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
            <p className="text-xs font-semibold text-slate-300 uppercase tracking-wider">Live URL</p>
          </div>
          <div className="bg-black/40 border border-slate-700 rounded-lg p-3 mb-2">
            <p className="text-xs font-mono text-slate-400 break-all leading-relaxed">
              <span className="text-emerald-400">GET</span>{' '}
              <span className="text-primary-300">{apiUrl}</span>
            </p>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="flex-1 bg-black/40 border border-slate-700 rounded-lg px-3 py-2">
              <p className="text-xs font-mono text-slate-400 truncate">{curlExample}</p>
            </div>
            <CopyButton text={curlExample} />
          </div>
          <p className="text-xs text-slate-600 mt-2">Add <code className="text-slate-500">?limit=N</code> to control how many records are returned.</p>
        </div>

        {/* Visibility indicator */}
        <div className="flex items-center gap-2 p-3 rounded-lg bg-slate-900/50 border border-slate-700">
          {visibility === 'public' ? (
            <>
              <Globe className="w-4 h-4 text-emerald-400 shrink-0" />
              <p className="text-xs text-slate-400"><span className="text-emerald-400 font-medium">Public</span> — anyone can call this endpoint with their access key.</p>
            </>
          ) : (
            <>
              <Lock className="w-4 h-4 text-slate-400 shrink-0" />
              <p className="text-xs text-slate-400"><span className="text-slate-300 font-medium">Private</span> — only your personal <code className="text-slate-400">mk_</code> keys can access this.</p>
            </>
          )}
        </div>

        {/* x-generator reference */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Zap className="w-3.5 h-3.5 text-yellow-400" />
            <p className="text-xs font-semibold text-slate-300 uppercase tracking-wider">Generator Annotations</p>
          </div>
          <p className="text-xs text-slate-500 mb-3">Add <code className="text-slate-400">x-generator</code> to any property to control what data is generated.</p>
          <div className="grid grid-cols-2 gap-x-4 gap-y-2">
            {GENERATORS.map(g => (
              <div key={g.annotation} className="bg-black/30 border border-slate-800 rounded-lg px-2.5 py-2">
                <code className="block text-xs text-primary-300 mb-0.5">{g.annotation.split(': ')[1]}</code>
                <span className="text-slate-600 text-xs">{g.description}</span>
                <span className="block text-slate-500 text-xs font-mono truncate">{g.example}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Schema structure hint */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <FileCode2 className="w-3.5 h-3.5 text-primary-400" />
            <p className="text-xs font-semibold text-slate-300 uppercase tracking-wider">Schema Format</p>
          </div>
          <pre className="text-xs font-mono text-slate-400 bg-black/40 border border-slate-700 rounded-lg p-4 overflow-x-auto leading-relaxed whitespace-pre-wrap">{`{
  "properties": {
    "fieldName": {
      "type": "string",
      "x-generator": "email"
    },
    "count": {
      "type": "integer",
      "x-generator": "number.int",
      "minimum": 1,
      "maximum": 100
    }
  }
}`}</pre>
        </div>
      </div>
    </div>
  )
}
