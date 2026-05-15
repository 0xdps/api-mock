'use client'

import { useEffect, useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import { useAuthContext } from '@/components/AuthProvider'
import { getTemplate, updateTemplate } from '@/lib/mockly-api'
import { TemplatePlayground } from '@/components/TemplatePlayground'
import type { MocklyTemplate, TemplateVisibility } from '@/types/mockly'
import { ChevronDown } from 'lucide-react'

const SLUG_RE = /^[a-z0-9-]{1,60}$/

export default function EditTemplatePage() {
  const { token } = useAuthContext()
  const router = useRouter()
  const params = useParams<{ id: string }>()

  const [tmpl, setTmpl] = useState<MocklyTemplate | null>(null)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [visibility, setVisibility] = useState<TemplateVisibility>('private')
  const [schemaJSON, setSchemaJSON] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token || !params.id) return
    getTemplate(token, params.id)
      .then(t => {
        setTmpl(t)
        setName(t.name)
        setDescription(t.description)
        setVisibility(t.visibility)
        // Pretty-print the stored JSON
        try { setSchemaJSON(JSON.stringify(JSON.parse(t.schema_json), null, 2)) }
        catch { setSchemaJSON(t.schema_json) }
      })
      .catch(() => router.push('/dashboard/templates'))
      .finally(() => setLoading(false))
  }, [token, params.id, router])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!token || !tmpl) return
    setError(null)

    try { JSON.parse(schemaJSON) } catch {
      setError('Schema JSON is not valid JSON.')
      return
    }

    setSaving(true)
    try {
      await updateTemplate(token, tmpl.id, { name, description, visibility, schema_json: schemaJSON })
      router.push('/dashboard/templates')
    } catch (e: any) {
      setError(e?.message || 'Failed to update template')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="p-8 text-slate-500 text-sm">Loading…</div>
    )
  }

  if (!tmpl) return null

  return (
    <div className="p-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-white">Edit Template</h1>
        <p className="text-slate-400 text-sm mt-1 font-mono">/{tmpl.slug}</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        {error && (
          <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-sm">{error}</div>
        )}

        <div>
          <label className="block text-xs font-medium text-slate-400 mb-1.5">Name</label>
          <input
            type="text" required value={name} onChange={e => setName(e.target.value)}
            className="w-full px-3 py-2 rounded-lg bg-slate-900 border border-white/10 text-white text-sm focus:outline-none focus:border-indigo-500 transition-colors"
          />
        </div>

        <div>
          <label className="block text-xs font-medium text-slate-400 mb-1.5">Description</label>
          <input
            type="text" value={description} onChange={e => setDescription(e.target.value)}
            className="w-full px-3 py-2 rounded-lg bg-slate-900 border border-white/10 text-white text-sm focus:outline-none focus:border-indigo-500 transition-colors"
          />
        </div>

        <div>
          <label className="block text-xs font-medium text-slate-400 mb-1.5">Visibility</label>
          <div className="relative w-40">
            <select
              value={visibility} onChange={e => setVisibility(e.target.value as TemplateVisibility)}
              className="w-full appearance-none px-3 py-2 pr-8 rounded-lg bg-slate-900 border border-white/10 text-white text-sm focus:outline-none focus:border-indigo-500 transition-colors cursor-pointer"
            >
              <option value="private">Private</option>
              <option value="public">Public</option>
            </select>
            <ChevronDown className="absolute right-2 top-2.5 w-4 h-4 text-slate-500 pointer-events-none" />
          </div>
        </div>

        <div>
          <label className="block text-xs font-medium text-slate-400 mb-1.5">Schema JSON</label>
          <textarea
            value={schemaJSON} onChange={e => setSchemaJSON(e.target.value)}
            rows={18} spellCheck={false}
            className="w-full px-3 py-2 rounded-lg bg-slate-900 border border-white/10 text-white text-xs font-mono focus:outline-none focus:border-indigo-500 transition-colors resize-y"
          />
        </div>

        <div className="flex items-center gap-3 pt-2">
          <button
            type="submit" disabled={saving}
            className="px-5 py-2 rounded-lg bg-indigo-600 text-white text-sm font-medium hover:bg-indigo-500 transition-colors disabled:opacity-60"
          >
            {saving ? 'Saving…' : 'Save Changes'}
          </button>
          <button type="button" onClick={() => router.back()} className="px-5 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-white/5 text-sm transition-colors">
            Cancel
          </button>
        </div>
      </form>

      {/* Playground */}
      <div className="mt-12">
        <h2 className="text-lg font-bold text-white mb-1">Playground</h2>
        <p className="text-slate-500 text-xs mb-4">Test your template endpoint with a real API key.</p>
        <TemplatePlayground
          templateId={tmpl.id}
          schemaProperties={(() => {
            try { return JSON.parse(schemaJSON)?.properties ?? {} }
            catch { return {} }
          })()}
        />
      </div>
    </div>
  )
}
