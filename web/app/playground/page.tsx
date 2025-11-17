'use client'

import { useState, useEffect } from 'react'
import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'

const API_URL = process.env.NEXT_PUBLIC_API_URL || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:8080')

export default function PlaygroundPage() {
  const [resources, setResources] = useState<string[]>([])
  const [selectedResource, setSelectedResource] = useState('')
  const [count, setCount] = useState(5)
  const [itemId, setItemId] = useState('')
  const [endpointType, setEndpointType] = useState<'collection' | 'single' | 'meta'>('collection')
  const [response, setResponse] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  
  useEffect(() => {
    fetch(`${API_URL}/`)
      .then(res => res.json())
      .then(data => {
        setResources(data.resources || [])
        if (data.resources?.length > 0) {
          setSelectedResource(data.resources[0])
        }
      })
      .catch(err => console.error('Failed to fetch resources:', err))
  }, [])
  
  const buildUrl = () => {
    if (!selectedResource) return ''
    
    switch (endpointType) {
      case 'collection':
        return `${API_URL}/api/${selectedResource}?count=${count}`
      case 'single':
        return `${API_URL}/api/${selectedResource}/${itemId || '1'}`
      case 'meta':
        return `${API_URL}/api/${selectedResource}/meta`
      default:
        return ''
    }
  }
  
  const handleFetch = async () => {
    setLoading(true)
    setError(null)
    
    try {
      const url = buildUrl()
      const res = await fetch(url)
      const data = await res.json()
      setResponse(data)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }
  
  const url = buildUrl()
  
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <Header />
      
      <main className="container mx-auto px-4 py-12 max-w-6xl">
        <h1 className="text-5xl font-bold text-white mb-6">API Playground</h1>
        <p className="text-xl text-slate-300 mb-12">
          Test API endpoints interactively
        </p>
        
        <div className="grid lg:grid-cols-2 gap-8">
          {/* Request Builder */}
          <div className="space-y-6">
            <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
              <h2 className="text-2xl font-bold text-white mb-6">Request</h2>
              
              {/* Resource Selection */}
              <div className="mb-4">
                <label className="block text-slate-300 mb-2 font-medium">
                  Resource
                </label>
                <select
                  value={selectedResource}
                  onChange={(e) => setSelectedResource(e.target.value)}
                  className="w-full bg-slate-700 text-white px-4 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none"
                >
                  {resources.map(resource => (
                    <option key={resource} value={resource}>
                      {resource}
                    </option>
                  ))}
                </select>
              </div>
              
              {/* Endpoint Type */}
              <div className="mb-4">
                <label className="block text-slate-300 mb-2 font-medium">
                  Endpoint Type
                </label>
                <div className="flex gap-2">
                  <button
                    onClick={() => setEndpointType('collection')}
                    className={`flex-1 px-4 py-2 rounded transition ${
                      endpointType === 'collection'
                        ? 'bg-primary-500 text-white'
                        : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                    }`}
                  >
                    Collection
                  </button>
                  <button
                    onClick={() => setEndpointType('single')}
                    className={`flex-1 px-4 py-2 rounded transition ${
                      endpointType === 'single'
                        ? 'bg-primary-500 text-white'
                        : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                    }`}
                  >
                    Single
                  </button>
                  <button
                    onClick={() => setEndpointType('meta')}
                    className={`flex-1 px-4 py-2 rounded transition ${
                      endpointType === 'meta'
                        ? 'bg-primary-500 text-white'
                        : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                    }`}
                  >
                    Meta
                  </button>
                </div>
              </div>
              
              {/* Count (for collection) */}
              {endpointType === 'collection' && (
                <div className="mb-4">
                  <label className="block text-slate-300 mb-2 font-medium">
                    Count (1-100)
                  </label>
                  <input
                    type="number"
                    min="1"
                    max="100"
                    value={count}
                    onChange={(e) => setCount(parseInt(e.target.value) || 1)}
                    className="w-full bg-slate-700 text-white px-4 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none"
                  />
                </div>
              )}
              
              {/* Item ID (for single) */}
              {endpointType === 'single' && (
                <div className="mb-4">
                  <label className="block text-slate-300 mb-2 font-medium">
                    Item ID
                  </label>
                  <input
                    type="text"
                    value={itemId}
                    onChange={(e) => setItemId(e.target.value)}
                    placeholder="1"
                    className="w-full bg-slate-700 text-white px-4 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none"
                  />
                </div>
              )}
              
              {/* URL Preview */}
              <div className="mb-6">
                <label className="block text-slate-300 mb-2 font-medium">
                  Request URL
                </label>
                <div className="bg-slate-900 p-3 rounded border border-slate-600 overflow-x-auto">
                  <code className="text-primary-500 text-sm">{url}</code>
                </div>
              </div>
              
              {/* Send Button */}
              <button
                onClick={handleFetch}
                disabled={loading || !selectedResource}
                className="w-full bg-primary-500 hover:bg-primary-600 disabled:bg-slate-700 text-white px-6 py-3 rounded font-semibold transition"
              >
                {loading ? 'Loading...' : 'Send Request'}
              </button>
            </div>
            
            {/* Code Examples */}
            <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
              <h3 className="text-xl font-bold text-white mb-4">Code Examples</h3>
              
              <div className="space-y-4">
                <CodeSnippet title="cURL" code={`curl "${url}"`} />
                <CodeSnippet 
                  title="JavaScript" 
                  code={`fetch('${url}')
  .then(res => res.json())
  .then(data => console.log(data));`} 
                />
                <CodeSnippet 
                  title="Python" 
                  code={`import requests
response = requests.get('${url}')
print(response.json())`} 
                />
              </div>
            </div>
          </div>
          
          {/* Response */}
          <div className="space-y-6">
            <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700 min-h-[600px]">
              <h2 className="text-2xl font-bold text-white mb-6">Response</h2>
              
              {error && (
                <div className="bg-red-900/20 border border-red-500 text-red-300 p-4 rounded">
                  <strong>Error:</strong> {error}
                </div>
              )}
              
              {response && !error && (
                <div className="space-y-4">
                  <div className="flex items-center justify-between mb-4">
                    <span className="text-green-400 font-semibold">Status: 200 OK</span>
                    <button
                      onClick={() => navigator.clipboard.writeText(JSON.stringify(response, null, 2))}
                      className="text-sm bg-slate-700 hover:bg-slate-600 text-slate-300 px-3 py-1 rounded transition"
                    >
                      Copy JSON
                    </button>
                  </div>
                  
                  <pre className="bg-slate-900 p-4 rounded border border-slate-600 overflow-auto max-h-[500px]">
                    <code className="text-sm text-green-400">
                      {JSON.stringify(response, null, 2)}
                    </code>
                  </pre>
                </div>
              )}
              
              {!response && !error && !loading && (
                <div className="text-center text-slate-400 py-20">
                  <p className="text-lg">Send a request to see the response</p>
                </div>
              )}
              
              {loading && (
                <div className="text-center text-slate-400 py-20">
                  <div className="animate-pulse text-lg">Loading...</div>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
      
      <Footer />
    </div>
  )
}

function CodeSnippet({ title, code }: { title: string; code: string }) {
  return (
    <div>
      <div className="text-slate-400 text-sm mb-2">{title}</div>
      <pre className="bg-slate-900 p-3 rounded border border-slate-600 overflow-x-auto">
        <code className="text-xs text-slate-300">{code}</code>
      </pre>
    </div>
  )
}
