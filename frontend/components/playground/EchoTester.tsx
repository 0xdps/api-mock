'use client'

import { useState } from 'react'
import { Play, RefreshCw, CheckCircle2, AlertCircle, Zap, Terminal, ClipboardList, X } from 'lucide-react'
import { getApiUrl, apiClient } from '@/lib/api'
import { RequestResponseLayout } from './RequestResponseLayout'

const API_URL = getApiUrl()

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
type ContentType = 'application/json' | 'application/x-www-form-urlencoded' | 'multipart/form-data'

export function EchoTester() {
  const [method, setMethod] = useState<HttpMethod>('GET')
  const [contentType, setContentType] = useState<ContentType>('application/json')
  const [queryParams, setQueryParams] = useState<Array<{ key: string; value: string }>>([{ key: '', value: '' }])
  const [headers, setHeaders] = useState<Array<{ key: string; value: string }>>([{ key: '', value: '' }])
  const [jsonBody, setJsonBody] = useState('{\n  "message": "Hello, API!"\n}')
  const [formData, setFormData] = useState<Array<{ key: string; value: string }>>([{ key: '', value: '' }])
  const [response, setResponse] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [requestTime, setRequestTime] = useState<number | null>(null)

  const addQueryParam = () => setQueryParams([...queryParams, { key: '', value: '' }])
  const removeQueryParam = (index: number) => setQueryParams(queryParams.filter((_, i) => i !== index))
  const updateQueryParam = (index: number, field: 'key' | 'value', value: string) => {
    const updated = [...queryParams]
    updated[index][field] = value
    setQueryParams(updated)
  }

  const addHeader = () => setHeaders([...headers, { key: '', value: '' }])
  const removeHeader = (index: number) => setHeaders(headers.filter((_, i) => i !== index))
  const updateHeader = (index: number, field: 'key' | 'value', value: string) => {
    const updated = [...headers]
    updated[index][field] = value
    setHeaders(updated)
  }

  const addFormField = () => setFormData([...formData, { key: '', value: '' }])
  const removeFormField = (index: number) => setFormData(formData.filter((_, i) => i !== index))
  const updateFormField = (index: number, field: 'key' | 'value', value: string) => {
    const updated = [...formData]
    updated[index][field] = value
    setFormData(updated)
  }

  const buildUrl = () => {
    let url = `${API_URL}/echo`
    const params = queryParams.filter(p => p.key).map(p => `${encodeURIComponent(p.key)}=${encodeURIComponent(p.value)}`)
    if (params.length > 0) {
      url += '?' + params.join('&')
    }
    return url
  }

  const handleSend = async () => {
    setLoading(true)
    setError(null)
    setResponse(null)
    setRequestTime(null)

    const startTime = performance.now()

    try {
      const url = buildUrl()
      const fetchHeaders: Record<string, string> = {}
      
      headers.filter(h => h.key).forEach(h => {
        fetchHeaders[h.key] = h.value
      })

      const options: RequestInit = {
        method,
        headers: fetchHeaders,
      }

      if (['POST', 'PUT', 'PATCH'].includes(method)) {
        if (contentType === 'application/json') {
          options.headers = { ...options.headers, 'Content-Type': contentType }
          options.body = jsonBody
        } else if (contentType === 'application/x-www-form-urlencoded') {
          options.headers = { ...options.headers, 'Content-Type': contentType }
          const params = new URLSearchParams()
          formData.filter(f => f.key).forEach(f => params.append(f.key, f.value))
          options.body = params.toString()
        } else if (contentType === 'multipart/form-data') {
          const form = new FormData()
          formData.filter(f => f.key).forEach(f => form.append(f.key, f.value))
          options.body = form
        }
      }

      const res = await fetch(url, options)
      const endTime = performance.now()
      setRequestTime(Math.round(endTime - startTime))

      if (!res.ok) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      }

      const data = await res.json()
      setResponse({
        data,
        status: res.status,
        headers: Object.fromEntries(res.headers.entries())
      })
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const generateCurl = () => {
    const url = buildUrl()
    let curl = `curl -X ${method} '${url}'`
    
    headers.filter(h => h.key).forEach(h => {
      curl += ` \\\n  -H '${h.key}: ${h.value}'`
    })

    if (['POST', 'PUT', 'PATCH'].includes(method)) {
      if (contentType === 'application/json') {
        curl += ` \\\n  -H 'Content-Type: application/json'`
        curl += ` \\\n  -d '${jsonBody.replace(/\n/g, '')}'`
      } else if (contentType === 'application/x-www-form-urlencoded') {
        curl += ` \\\n  -H 'Content-Type: application/x-www-form-urlencoded'`
        const params = formData.filter(f => f.key).map(f => `${f.key}=${f.value}`).join('&')
        curl += ` \\\n  -d '${params}'`
      }
    }

    return curl
  }

  const requestPanel = (
    <div className="space-y-4">
      {/* Method Selector */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          HTTP Method
        </label>
        <div className="flex gap-2 flex-wrap">
          {(['GET', 'POST', 'PUT', 'PATCH', 'DELETE'] as HttpMethod[]).map(m => (
            <button
              key={m}
              onClick={() => setMethod(m)}
              className={`px-4 py-2 rounded transition text-sm font-semibold ${
                method === m
                  ? 'bg-blue-500 text-white'
                  : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
              }`}
            >
              {m}
            </button>
          ))}
        </div>
      </div>

      {/* Content Type (for POST/PUT/PATCH) */}
      {['POST', 'PUT', 'PATCH'].includes(method) && (
        <div>
          <label className="block text-slate-300 mb-2 font-medium text-sm">
            Content-Type
          </label>
          <select
            value={contentType}
            onChange={(e) => setContentType(e.target.value as ContentType)}
            className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
          >
            <option value="application/json">JSON (application/json)</option>
            <option value="application/x-www-form-urlencoded">Form URL Encoded</option>
            <option value="multipart/form-data">Multipart Form Data</option>
          </select>
        </div>
      )}

      {/* Query Parameters */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <label className="block text-slate-300 font-medium text-sm">
            Query Parameters
          </label>
          <button
            onClick={addQueryParam}
            className="text-xs text-blue-400 hover:text-blue-300"
          >
            + Add
          </button>
        </div>
        <div className="space-y-2">
          {queryParams.map((param, index) => (
            <div key={index} className="flex gap-2">
              <input
                type="text"
                placeholder="Key"
                value={param.key}
                onChange={(e) => updateQueryParam(index, 'key', e.target.value)}
                className="flex-1 bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
              />
              <input
                type="text"
                placeholder="Value"
                value={param.value}
                onChange={(e) => updateQueryParam(index, 'value', e.target.value)}
                className="flex-1 bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
              />
              <button
                onClick={() => removeQueryParam(index)}
                className="px-3 py-2 bg-red-900/30 text-red-400 hover:bg-red-900/50 rounded text-sm"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Headers */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <label className="block text-slate-300 font-medium text-sm">
            Headers
          </label>
          <button
            onClick={addHeader}
            className="text-xs text-blue-400 hover:text-blue-300"
          >
            + Add
          </button>
        </div>
        <div className="space-y-2">
          {headers.map((header, index) => (
            <div key={index} className="flex gap-2">
              <input
                type="text"
                placeholder="Header Name"
                value={header.key}
                onChange={(e) => updateHeader(index, 'key', e.target.value)}
                className="flex-1 bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
              />
              <input
                type="text"
                placeholder="Header Value"
                value={header.value}
                onChange={(e) => updateHeader(index, 'value', e.target.value)}
                className="flex-1 bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
              />
              <button
                onClick={() => removeHeader(index)}
                className="px-3 py-2 bg-red-900/30 text-red-400 hover:bg-red-900/50 rounded text-sm"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Request Body */}
      {['POST', 'PUT', 'PATCH'].includes(method) && (
        <div>
          <label className="block text-slate-300 mb-2 font-medium text-sm">
            Request Body
          </label>
          {contentType === 'application/json' ? (
            <textarea
              value={jsonBody}
              onChange={(e) => setJsonBody(e.target.value)}
              className="w-full h-32 bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm font-mono"
              placeholder='{"key": "value"}'
            />
          ) : (
            <div className="space-y-2">
              {formData.map((field, index) => (
                <div key={index} className="flex gap-2">
                  <input
                    type="text"
                    placeholder="Field Name"
                    value={field.key}
                    onChange={(e) => updateFormField(index, 'key', e.target.value)}
                    className="flex-1 bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
                  />
                  <input
                    type="text"
                    placeholder="Field Value"
                    value={field.value}
                    onChange={(e) => updateFormField(index, 'value', e.target.value)}
                    className="flex-1 bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
                  />
                  <button
                    onClick={() => removeFormField(index)}
                    className="px-3 py-2 bg-red-900/30 text-red-400 hover:bg-red-900/50 rounded text-sm"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </div>
              ))}
              <button
                onClick={addFormField}
                className="text-xs text-blue-400 hover:text-blue-300"
              >
                + Add Field
              </button>
            </div>
          )}
        </div>
      )}

      {/* URL Preview */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Request URL
        </label>
        <div className="bg-slate-900 p-3 rounded border border-slate-600">
          <code className="text-blue-400 text-sm break-all">{buildUrl()}</code>
        </div>
      </div>

      {/* cURL Command */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          cURL Command
        </label>
        <div className="bg-slate-900 p-4 rounded border border-slate-600 overflow-auto">
          <pre className="text-slate-400 text-sm">
            <code>{generateCurl()}</code>
          </pre>
        </div>
      </div>
    </div>
  )

  const actionButton = (
    <button
      onClick={handleSend}
      disabled={loading}
      className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-6 py-3 rounded-xl font-semibold text-base transition flex items-center justify-center gap-3"
    >
      {loading ? <RefreshCw className="w-5 h-5 animate-spin" /> : <Play className="w-5 h-5 fill-current" />}
      {loading ? 'Sending...' : 'Send Request'}
    </button>
  )

  const responsePanel = (
    <div className="flex flex-col h-full space-y-4">
      {!response && !error ? (
        <div className="text-center py-12 text-slate-400">
          <p className="text-lg mb-2">No response yet</p>
          <p className="text-sm">Configure your request and click Send</p>
        </div>
      ) : error ? (
        <div className="bg-red-900/20 border border-red-500/50 p-4 rounded">
          <div className="flex items-center gap-2 mb-2">
            <AlertCircle className="w-5 h-5 text-red-400" />
            <span className="text-red-400 font-bold uppercase tracking-wider text-xs">Error</span>
          </div>
          <p className="text-red-300 text-sm">{error}</p>
        </div>
      ) : (
        <>
          <div className="flex items-center gap-4 flex-wrap">
            <div className="flex items-center gap-1.5">
              <CheckCircle2 className="w-4 h-4 text-emerald-400" />
              <span className="text-emerald-400 font-bold uppercase tracking-wider text-xs">Success</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="text-slate-500 font-bold uppercase tracking-wider text-[10px]">Status:</span>
              <span className="text-white font-mono text-xs bg-emerald-500/10 px-1.5 py-0.5 rounded border border-emerald-500/20">{response.status}</span>
            </div>
            {requestTime !== null && (
              <div className="flex items-center gap-1.5">
                <Zap className="w-3.5 h-3.5 text-amber-500" />
                <span className="text-white font-mono text-xs">{requestTime}ms</span>
              </div>
            )}
          </div>

          <div className="flex-1 flex flex-col min-h-0">
            <h4 className="text-sm font-semibold text-slate-300 mb-2">Response Body</h4>
            <div className="flex-1 bg-slate-900 p-4 rounded border border-slate-600 overflow-auto">
              <pre className="text-green-400 text-sm">
                <code>{JSON.stringify(response.data, null, 2)}</code>
              </pre>
            </div>
          </div>
        </>
      )}
    </div>
  )

  return (
    <RequestResponseLayout
      requestPanel={requestPanel}
      responsePanel={responsePanel}
      actionButton={actionButton}
      isLoading={loading}
    />
  )
}
