'use client'

import { useState } from 'react'
import { getApiUrl, apiClient } from '@/lib/api'
import { RequestResponseLayout } from './RequestResponseLayout'

const API_URL = getApiUrl()

export function MiddlewareTester() {
  // Global parameters
  const [delay, setDelay] = useState(0)
  const [flakyRate, setFlakyRate] = useState(0)
  const [skipCache, setSkipCache] = useState(false)
  const [fields, setFields] = useState('id,name,email')
  
  // Headers
  const [requestId, setRequestId] = useState('')
  const [tenantId, setTenantId] = useState('tenant-123')
  const [role, setRole] = useState('admin')
  const [idempotencyKey, setIdempotencyKey] = useState('')
  
  // Query parameters
  const [page, setPage] = useState(1)
  const [limit, setLimit] = useState(10)
  const [sortField, setSortField] = useState('id')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc')
  const [searchQuery, setSearchQuery] = useState('')
  const [searchFields, setSearchFields] = useState('')
  
  // Test resource
  const [resource, setResource] = useState('products')
  
  // Response
  const [response, setResponse] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [requestTime, setRequestTime] = useState<number | null>(null)

  const generateRequestId = () => {
    const id = `req-${Math.random().toString(36).substr(2, 9)}`
    setRequestId(id)
  }

  const generateIdempotencyKey = () => {
    const key = `idem-${Math.random().toString(36).substr(2, 9)}`
    setIdempotencyKey(key)
  }

  const buildUrl = () => {
    let url = `${API_URL}/${resource}`
    const params = new URLSearchParams()
    
    // Pagination
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    
    // Sorting
    if (sortField) {
      params.append('sort', sortField)
      params.append('order', sortOrder)
    }
    
    // Search
    if (searchQuery) {
      params.append('q', searchQuery)
      if (searchFields) {
        params.append('search_fields', searchFields)
      }
    }
    
    // Global parameters
    if (delay > 0) {
      params.append('delay', delay.toString())
    }
    if (flakyRate > 0) {
      params.append('flakyRate', (flakyRate / 100).toString())
    }
    if (skipCache) {
      params.append('skip_cache', 'true')
    }
    if (fields) {
      params.append('fields', fields)
    }
    
    return `${url}?${params.toString()}`
  }

  const handleTest = async () => {
    setLoading(true)
    setError(null)
    setResponse(null)
    setRequestTime(null)

    const startTime = performance.now()

    try {
      const url = buildUrl()
      const headers: Record<string, string> = {}
      
      if (requestId) headers['X-Request-ID'] = requestId
      if (tenantId) headers['X-Tenant-ID'] = tenantId
      if (role) headers['X-Role'] = role
      if (idempotencyKey) headers['Idempotency-Key'] = idempotencyKey

      const res = await apiClient.get(url, { headers })
      const endTime = performance.now()
      setRequestTime(Math.round(endTime - startTime))

      if (res.isError()) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      }

      const data = res.json()
      const responseHeaders = res.headers

      setResponse({
        data,
        status: res.status,
        headers: responseHeaders,
      })
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const requestPanel = (
    <>
      {/* Test Resource */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Test Resource
        </label>
        <select
          value={resource}
          onChange={(e) => setResource(e.target.value)}
          className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
        >
          <option value="products">Products</option>
          <option value="users">Users</option>
          <option value="orders">Orders</option>
          <option value="articles">Articles</option>
        </select>
      </div>

      {/* Global Middleware Parameters */}
      <div className="space-y-4 p-4 bg-slate-900/50 rounded border border-slate-700">
        <h3 className="text-sm font-semibold text-white">Global Middleware Parameters</h3>
        
        <div>
          <label className="block text-slate-300 mb-2 text-xs">
            Delay: {delay}ms
          </label>
          <input
            type="range"
            min="0"
            max="5000"
            step="100"
            value={delay}
            onChange={(e) => setDelay(parseInt(e.target.value))}
            className="w-full"
          />
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-xs">
            Flaky Rate: {flakyRate}% {flakyRate > 50 && '⚠️'}
          </label>
          <input
            type="range"
            min="0"
            max="100"
            step="5"
            value={flakyRate}
            onChange={(e) => setFlakyRate(parseInt(e.target.value))}
            className="w-full"
          />
        </div>
        
        <div>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={skipCache}
              onChange={(e) => setSkipCache(e.target.checked)}
              className="rounded border-slate-600 text-blue-500"
            />
            <span className="text-slate-300 text-sm">Skip Cache</span>
          </label>
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-xs">
            Field Filtering (comma-separated)
          </label>
          <input
            type="text"
            value={fields}
            onChange={(e) => setFields(e.target.value)}
            placeholder="id,name,email"
            className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
          />
        </div>
      </div>

      {/* Request Headers */}
      <div className="space-y-4 p-4 bg-slate-900/50 rounded border border-slate-700">
        <h3 className="text-sm font-semibold text-white">Request Headers</h3>
        
        <div>
          <div className="flex items-center justify-between mb-2">
            <label className="text-slate-300 text-xs">X-Request-ID</label>
            <button
              onClick={generateRequestId}
              className="text-xs text-blue-400 hover:text-blue-300"
            >
              ↻ Generate
            </button>
          </div>
          <input
            type="text"
            value={requestId}
            onChange={(e) => setRequestId(e.target.value)}
            placeholder="Auto-generated if empty"
            className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
          />
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-xs">X-Tenant-ID</label>
          <input
            type="text"
            value={tenantId}
            onChange={(e) => setTenantId(e.target.value)}
            placeholder="tenant-123"
            className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
          />
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-xs">X-Role</label>
          <select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
          >
            <option value="">None</option>
            <option value="admin">Admin</option>
            <option value="user">User</option>
            <option value="guest">Guest</option>
          </select>
        </div>
        
        <div>
          <div className="flex items-center justify-between mb-2">
            <label className="text-slate-300 text-xs">Idempotency-Key</label>
            <button
              onClick={generateIdempotencyKey}
              className="text-xs text-blue-400 hover:text-blue-300"
            >
              ↻ Generate
            </button>
          </div>
          <input
            type="text"
            value={idempotencyKey}
            onChange={(e) => setIdempotencyKey(e.target.value)}
            placeholder="Optional for GET requests"
            className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
          />
        </div>
      </div>

      {/* Query Parameters */}
      <div className="space-y-4 p-4 bg-slate-900/50 rounded border border-slate-700">
        <h3 className="text-sm font-semibold text-white">Query Parameters</h3>
        
        <div className="grid grid-cols-2 gap-3">
          <div>
            <label className="block text-slate-300 mb-2 text-xs">Page</label>
            <input
              type="number"
              min="1"
              value={page}
              onChange={(e) => setPage(parseInt(e.target.value) || 1)}
              className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
            />
          </div>
          <div>
            <label className="block text-slate-300 mb-2 text-xs">Limit</label>
            <input
              type="number"
              min="1"
              max="100"
              value={limit}
              onChange={(e) => setLimit(parseInt(e.target.value) || 10)}
              className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
            />
          </div>
        </div>
        
        <div className="grid grid-cols-2 gap-3">
          <div>
            <label className="block text-slate-300 mb-2 text-xs">Sort by</label>
            <input
              type="text"
              value={sortField}
              onChange={(e) => setSortField(e.target.value)}
              placeholder="id"
              className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
            />
          </div>
          <div>
            <label className="block text-slate-300 mb-2 text-xs">Order</label>
            <select
              value={sortOrder}
              onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
              className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
            >
              <option value="asc">↑ Ascending</option>
              <option value="desc">↓ Descending</option>
            </select>
          </div>
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-xs">Search Query</label>
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search..."
            className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
          />
        </div>
        
        {searchQuery && (
          <div>
            <label className="block text-slate-300 mb-2 text-xs">Search Fields (optional)</label>
            <input
              type="text"
              value={searchFields}
              onChange={(e) => setSearchFields(e.target.value)}
              placeholder="name,description"
              className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
            />
          </div>
        )}
      </div>

      {/* URL Preview */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Request URL
        </label>
        <div className="bg-slate-900 p-3 rounded border border-slate-600 overflow-auto">
          <code className="text-blue-400 text-xs break-all">{buildUrl()}</code>
        </div>
      </div>
    </>
  )

  const actionButton = (
    <button
      onClick={handleTest}
      disabled={loading}
      className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-6 py-3 rounded font-semibold transition"
    >
      {loading ? 'Testing...' : '▶ Run Test'}
    </button>
  )

  const responsePanel = (
    <>
      {(response || error) ? (
        <div className="flex flex-col h-full space-y-4">
          {error ? (
            <div className="bg-red-900/20 border border-red-500/50 p-4 rounded">
              <div className="flex items-center gap-2 mb-2">
                <span className="text-red-400 font-semibold">❌ Error</span>
              </div>
              <p className="text-red-300 text-sm">{error}</p>
            </div>
          ) : (
            <>
              <div className="flex items-center gap-3 flex-wrap">
                <span className="text-green-400 font-semibold">✓ Success</span>
                <span className="text-slate-400 text-sm">
                  Status: <span className="text-white">{response.status}</span>
                </span>
                {requestTime !== null && (
                  <span className="text-slate-400 text-sm">
                    Time: <span className="text-white">{requestTime}ms</span>
                  </span>
                )}
              </div>

              {/* Response Headers */}
              {response.headers && Object.keys(response.headers).length > 0 && (
                <div>
                  <h3 className="text-lg font-semibold text-white mb-2">Response Headers</h3>
                  <div className="bg-slate-900 p-4 rounded border border-slate-600 space-y-1 max-h-48 overflow-auto">
                    {Object.entries(response.headers).map(([key, value]) => (
                      <div key={key} className="text-sm">
                        <span className="text-blue-400">{key}:</span>{' '}
                        <span className="text-slate-300">{value as string}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Response Body */}
              <div className="flex-1 flex flex-col min-h-0">
                <h3 className="text-lg font-semibold text-white mb-2">Response Body</h3>
                <div className="flex-1 bg-slate-900 p-4 rounded border border-slate-600 overflow-auto">
                  <pre className="text-green-400 text-sm">
                    <code>{JSON.stringify(response.data, null, 2)}</code>
                  </pre>
                </div>
              </div>
            </>
          )}
        </div>
      ) : (
        <div className="text-center py-12 text-slate-400">
          <p>Run a test to see the response</p>
        </div>
      )}
    </>
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
