'use client'

import { useState } from 'react'
import { getApiUrl } from '@/lib/api'

const API_URL = getApiUrl()

const statusCodes = {
  '2xx': [
    { code: 200, name: 'OK', description: 'Standard success response' },
    { code: 201, name: 'Created', description: 'Resource created successfully' },
    { code: 202, name: 'Accepted', description: 'Request accepted for processing' },
    { code: 204, name: 'No Content', description: 'Success with no response body' },
  ],
  '3xx': [
    { code: 301, name: 'Moved Permanently', description: 'Resource permanently moved' },
    { code: 302, name: 'Found', description: 'Temporary redirect' },
    { code: 304, name: 'Not Modified', description: 'Cached version still valid' },
  ],
  '4xx': [
    { code: 400, name: 'Bad Request', description: 'Invalid request syntax' },
    { code: 401, name: 'Unauthorized', description: 'Authentication required' },
    { code: 403, name: 'Forbidden', description: 'Access denied' },
    { code: 404, name: 'Not Found', description: 'Resource not found' },
    { code: 409, name: 'Conflict', description: 'Request conflicts with current state' },
    { code: 422, name: 'Unprocessable Entity', description: 'Validation error' },
    { code: 429, name: 'Too Many Requests', description: 'Rate limit exceeded' },
  ],
  '5xx': [
    { code: 500, name: 'Internal Server Error', description: 'Generic server error' },
    { code: 502, name: 'Bad Gateway', description: 'Invalid response from upstream' },
    { code: 503, name: 'Service Unavailable', description: 'Service temporarily down' },
    { code: 504, name: 'Gateway Timeout', description: 'Upstream timeout' },
  ],
}

export function StatusCodeGenerator() {
  const [selectedCode, setSelectedCode] = useState(200)
  const [customMessage, setCustomMessage] = useState('')
  const [response, setResponse] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [requestTime, setRequestTime] = useState<number | null>(null)

  const handleGenerate = async () => {
    setLoading(true)
    setError(null)
    setResponse(null)
    setRequestTime(null)

    const startTime = performance.now()

    try {
      let url = `${API_URL}/status/${selectedCode}`
      if (customMessage) {
        url += `?message=${encodeURIComponent(customMessage)}`
      }

      const res = await fetch(url)
      const endTime = performance.now()
      setRequestTime(Math.round(endTime - startTime))

      const data = await res.json().catch(() => null)
      
      setResponse({
        data,
        status: res.status,
        statusText: res.statusText,
      })
    } catch (err: any) {
      const endTime = performance.now()
      setRequestTime(Math.round(endTime - startTime))
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const getColorForCode = (code: number) => {
    if (code >= 200 && code < 300) return 'text-green-400 border-green-500 bg-green-900/20'
    if (code >= 300 && code < 400) return 'text-blue-400 border-blue-500 bg-blue-900/20'
    if (code >= 400 && code < 500) return 'text-yellow-400 border-yellow-500 bg-yellow-900/20'
    if (code >= 500) return 'text-red-400 border-red-500 bg-red-900/20'
    return 'text-slate-400 border-slate-500 bg-slate-900/20'
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">
          Status Code Generator
        </h2>
        <p className="text-slate-400">
          Generate specific HTTP status codes for testing error handling
        </p>
      </div>

      {/* Quick Select */}
      <div>
        <label className="block text-slate-300 mb-3 font-medium text-sm">
          Quick Select
        </label>
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
          <button
            onClick={() => { setSelectedCode(200); setCustomMessage('') }}
            className="px-4 py-2 rounded transition text-sm font-semibold bg-green-900/30 text-green-400 border border-green-500/50 hover:bg-green-900/50"
          >
            200 OK
          </button>
          <button
            onClick={() => { setSelectedCode(201); setCustomMessage('') }}
            className="px-4 py-2 rounded transition text-sm font-semibold bg-green-900/30 text-green-400 border border-green-500/50 hover:bg-green-900/50"
          >
            201 Created
          </button>
          <button
            onClick={() => { setSelectedCode(400); setCustomMessage('') }}
            className="px-4 py-2 rounded transition text-sm font-semibold bg-yellow-900/30 text-yellow-400 border border-yellow-500/50 hover:bg-yellow-900/50"
          >
            400 Bad Request
          </button>
          <button
            onClick={() => { setSelectedCode(401); setCustomMessage('') }}
            className="px-4 py-2 rounded transition text-sm font-semibold bg-yellow-900/30 text-yellow-400 border border-yellow-500/50 hover:bg-yellow-900/50"
          >
            401 Unauthorized
          </button>
          <button
            onClick={() => { setSelectedCode(404); setCustomMessage('') }}
            className="px-4 py-2 rounded transition text-sm font-semibold bg-yellow-900/30 text-yellow-400 border border-yellow-500/50 hover:bg-yellow-900/50"
          >
            404 Not Found
          </button>
          <button
            onClick={() => { setSelectedCode(500); setCustomMessage('') }}
            className="px-4 py-2 rounded transition text-sm font-semibold bg-red-900/30 text-red-400 border border-red-500/50 hover:bg-red-900/50"
          >
            500 Error
          </button>
        </div>
      </div>

      {/* Status Code Categories */}
      <div className="space-y-3">
        <label className="block text-slate-300 font-medium text-sm">
          All Status Codes
        </label>
        {Object.entries(statusCodes).map(([category, codes]) => (
          <div key={category} className="bg-slate-900/50 rounded border border-slate-700 p-3">
            <h3 className="text-sm font-semibold text-slate-300 mb-2">{category} - {
              category === '2xx' ? 'Success' :
              category === '3xx' ? 'Redirection' :
              category === '4xx' ? 'Client Error' :
              'Server Error'
            }</h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {codes.map(({ code, name, description }) => (
                <button
                  key={code}
                  onClick={() => setSelectedCode(code)}
                  className={`text-left p-3 rounded border transition ${
                    selectedCode === code
                      ? 'border-blue-500 bg-blue-900/30'
                      : 'border-slate-600 hover:border-slate-500 bg-slate-800/50'
                  }`}
                >
                  <div className="flex items-center gap-2 mb-1">
                    <span className={`font-semibold ${getColorForCode(code).split(' ')[0]}`}>
                      {code}
                    </span>
                    <span className="text-slate-300 text-sm">{name}</span>
                  </div>
                  <p className="text-xs text-slate-500">{description}</p>
                </button>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* Custom Message */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Custom Message (Optional)
        </label>
        <input
          type="text"
          value={customMessage}
          onChange={(e) => setCustomMessage(e.target.value)}
          placeholder="Enter custom error message..."
          className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm"
        />
      </div>

      {/* URL Preview */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Request URL
        </label>
        <div className="bg-slate-900 p-3 rounded border border-slate-600">
          <code className="text-blue-400 text-sm break-all">
            {API_URL}/status/{selectedCode}{customMessage ? `?message=${encodeURIComponent(customMessage)}` : ''}
          </code>
        </div>
      </div>

      {/* Generate Button */}
      <button
        onClick={handleGenerate}
        disabled={loading}
        className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-6 py-3 rounded font-semibold transition"
      >
        {loading ? 'Generating...' : '▶ Generate Response'}
      </button>

      {/* Response */}
      {(response || error) && (
        <div className="border-t border-slate-700 pt-6 space-y-4">
          <div className="flex items-center gap-3 flex-wrap">
            <span className={`font-semibold ${
              response?.status >= 200 && response?.status < 300 ? 'text-green-400' :
              response?.status >= 300 && response?.status < 400 ? 'text-blue-400' :
              response?.status >= 400 && response?.status < 500 ? 'text-yellow-400' :
              'text-red-400'
            }`}>
              {response?.status || 'Error'}
            </span>
            {response?.statusText && (
              <span className="text-slate-300">{response.statusText}</span>
            )}
            {requestTime !== null && (
              <span className="text-slate-400 text-sm">
                Time: <span className="text-white">{requestTime}ms</span>
              </span>
            )}
          </div>

          {error && (
            <div className="bg-red-900/20 border border-red-500/50 p-4 rounded">
              <p className="text-red-300 text-sm">{error}</p>
            </div>
          )}

          {response?.data && (
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Response Body</h3>
              <div className="bg-slate-900 p-4 rounded border border-slate-600 overflow-auto">
                <pre className={`text-sm ${getColorForCode(response.status).split(' ')[0]}`}>
                  <code>{JSON.stringify(response.data, null, 2)}</code>
                </pre>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
