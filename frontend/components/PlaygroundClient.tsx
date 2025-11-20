'use client'

import { useState } from 'react'
import { getApiUrl } from '@/lib/api'

const API_URL = getApiUrl()

interface PlaygroundClientProps {
  resources: string[]
  groups: Record<string, string[]>
}

export function PlaygroundClient({ resources, groups }: PlaygroundClientProps) {
  const [selectedResource, setSelectedResource] = useState(resources[0] || '')
  const [selectedGroup, setSelectedGroup] = useState<string>('')
  const [useGroupPath, setUseGroupPath] = useState(false)
  const [count, setCount] = useState(5)
  const [itemId, setItemId] = useState('')
  const [endpointType, setEndpointType] = useState<'collection' | 'single' | 'meta'>('collection')
  const [response, setResponse] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  
  // Filter resources based on selected group
  const availableResources = selectedGroup 
    ? (groups[selectedGroup] || [])
    : resources
  
  // Get the group for the selected resource
  const resourceGroup = Object.entries(groups).find(([_, resources]) => 
    resources.includes(selectedResource)
  )?.[0] || ''
  
  const buildUrl = (withGroup: boolean = false) => {
    if (!selectedResource) return ''
    
    // Build base path based on whether to use group path
    const basePath = withGroup && resourceGroup
      ? `${API_URL}/${resourceGroup}/${selectedResource}`
      : `${API_URL}/${selectedResource}`
    
    switch (endpointType) {
      case 'collection':
        return `${basePath}?count=${count}`
      case 'single':
        return `${basePath}/${itemId || '1'}`
      case 'meta':
        return `${basePath}/meta`
      default:
        return ''
    }
  }
  
  const directUrl = buildUrl(false)
  const groupUrl = resourceGroup ? buildUrl(true) : null
  
  const handleFetch = async () => {
    setLoading(true)
    setError(null)
    
    try {
      // Always use direct URL for the actual API call
      const res = await fetch(directUrl)
      const data = await res.json()
      setResponse(data)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }
  
  // Update selected resource when group changes
  const handleGroupChange = (newGroup: string) => {
    setSelectedGroup(newGroup)
    // Reset to first resource in the new group
    if (newGroup && groups[newGroup]?.length > 0) {
      setSelectedResource(groups[newGroup][0])
    } else if (!newGroup) {
      setSelectedResource(resources[0] || '')
    }
  }
  
  return (
    <div className="grid lg:grid-cols-2 gap-8">
      {/* Request Builder */}
      <div className="space-y-6">
        <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
          <h2 className="text-2xl font-bold text-white mb-6">Request</h2>
          
          {/* Group Selection (Optional) */}
          <div className="mb-4">
            <label className="block text-slate-300 mb-2 font-medium">
              Filter by Group (Optional)
            </label>
            <select
              value={selectedGroup}
              onChange={(e) => handleGroupChange(e.target.value)}
              className="w-full bg-slate-700 text-white px-4 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none"
            >
              <option value="">All Resources</option>
              {Object.keys(groups).sort().map(group => (
                <option key={group} value={group}>
                  {group} ({groups[group].length} resources)
                </option>
              ))}
            </select>
            {selectedGroup && (
              <p className="text-xs text-slate-400 mt-1">
                Showing only resources from {selectedGroup} group
              </p>
            )}
          </div>
          
          {/* Resource Selection */}
          <div className="mb-4">
            <label className="block text-slate-300 mb-2 font-medium">
              Resource
              {resourceGroup && !selectedGroup && (
                <span className="ml-2 text-xs text-slate-400">
                  (in {resourceGroup} group)
                </span>
              )}
            </label>
            <select
              value={selectedResource}
              onChange={(e) => setSelectedResource(e.target.value)}
              className="w-full bg-slate-700 text-white px-4 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none"
            >
              {availableResources.map(resource => (
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
          
          {/* Query Parameters */}
          {endpointType === 'collection' && (
            <div className="mb-4">
              <label className="block text-slate-300 mb-2 font-medium">
                Query Parameters
              </label>
              <div className="bg-slate-900/50 p-4 rounded border border-slate-600 space-y-3">
                <div>
                  <label className="block text-slate-400 text-sm mb-1">
                    count <span className="text-slate-500">(1-100)</span>
                  </label>
                  <input
                    type="number"
                    min="1"
                    max="100"
                    value={count}
                    onChange={(e) => setCount(parseInt(e.target.value) || 1)}
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none"
                  />
                  <p className="text-xs text-slate-500 mt-1">
                    Number of items to return
                  </p>
                </div>
              </div>
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
              Available URLs
            </label>
            
            {/* Direct URL */}
            <div className="mb-3">
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-slate-400">Direct Access</span>
                <span className="text-xs text-green-400">✓ Using this</span>
              </div>
              <div className="bg-slate-900 p-3 rounded border border-green-500/30 overflow-x-auto">
                <code className="text-green-400 text-sm break-all">{directUrl}</code>
              </div>
            </div>
            
            {/* Group URL */}
            {groupUrl && (
              <div>
                <div className="flex items-center justify-between mb-1">
                  <span className="text-xs text-slate-400">Via Group Path</span>
                  <span className="text-xs text-slate-500">(alternative)</span>
                </div>
                <div className="bg-slate-900 p-3 rounded border border-slate-600 overflow-x-auto">
                  <code className="text-slate-400 text-sm break-all">{groupUrl}</code>
                </div>
              </div>
            )}
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
            <CodeSnippet title="cURL" code={`curl "${directUrl}"`} />
            <CodeSnippet 
              title="JavaScript" 
              code={`fetch('${directUrl}')
  .then(res => res.json())
  .then(data => console.log(data));`} 
            />
            <CodeSnippet 
              title="Python" 
              code={`import requests
response = requests.get('${directUrl}')
print(response.json())`} 
            />
          </div>
        </div>
        
        {/* Supported Features */}
        <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
          <h3 className="text-xl font-bold text-white mb-4">Supported Features</h3>
          
          <div className="space-y-3 text-sm">
            <div className="flex items-start gap-2">
              <span className="text-green-400 mt-0.5">✓</span>
              <div>
                <strong className="text-slate-200">Query Parameters:</strong>
                <span className="text-slate-400"> ?count=N (max 100)</span>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <span className="text-green-400 mt-0.5">✓</span>
              <div>
                <strong className="text-slate-200">Group Paths:</strong>
                <span className="text-slate-400"> /{resourceGroup}/{selectedResource}</span>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <span className="text-green-400 mt-0.5">✓</span>
              <div>
                <strong className="text-slate-200">Direct Access:</strong>
                <span className="text-slate-400"> /{selectedResource}</span>
              </div>
            </div>
            <div className="flex items-start gap-2">
              <span className="text-green-400 mt-0.5">✓</span>
              <div>
                <strong className="text-slate-200">CORS:</strong>
                <span className="text-slate-400"> Enabled for all origins</span>
              </div>
            </div>
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

