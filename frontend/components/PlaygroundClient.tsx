'use client'

import { useState, useRef, useEffect } from 'react'
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
  const [noCache, setNoCache] = useState(false)
  const [response, setResponse] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [requestTime, setRequestTime] = useState<number | null>(null)
  const [resourceSearch, setResourceSearch] = useState('')
  const [isResourceDropdownOpen, setIsResourceDropdownOpen] = useState(false)
  const resourceDropdownRef = useRef<HTMLDivElement>(null)
  const [groupSearch, setGroupSearch] = useState('')
  const [isGroupDropdownOpen, setIsGroupDropdownOpen] = useState(false)
  const groupDropdownRef = useRef<HTMLDivElement>(null)
  
  // Filter resources based on selected group
  const availableResources = selectedGroup 
    ? (groups[selectedGroup] || [])
    : resources
  
  // Filter resources based on search
  const filteredResources = availableResources.filter(resource =>
    resource.toLowerCase().includes(resourceSearch.toLowerCase())
  )
  
  // Filter groups based on search
  const groupNames = Object.keys(groups).sort()
  const filteredGroups = groupNames.filter(group =>
    group.toLowerCase().includes(groupSearch.toLowerCase())
  )
  
  // Get the group for the selected resource
  const resourceGroup = Object.entries(groups).find(([_, resources]) => 
    resources.includes(selectedResource)
  )?.[0] || ''
  
  // Close dropdowns when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (resourceDropdownRef.current && !resourceDropdownRef.current.contains(event.target as Node)) {
        setIsResourceDropdownOpen(false)
      }
      if (groupDropdownRef.current && !groupDropdownRef.current.contains(event.target as Node)) {
        setIsGroupDropdownOpen(false)
      }
    }
    
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])
  
  const buildUrl = (withGroup: boolean = false) => {
    if (!selectedResource) return ''
    
    // Build base path based on whether to use group path
    const basePath = withGroup && resourceGroup
      ? `${API_URL}/${resourceGroup}/${selectedResource}`
      : `${API_URL}/${selectedResource}`
    
    let url = ''
    switch (endpointType) {
      case 'collection':
        url = `${basePath}?count=${count}`
        break
      case 'single':
        url = `${basePath}/${itemId || '1'}`
        break
      case 'meta':
        url = `${basePath}/meta`
        break
      default:
        return ''
    }
    
    // Add nocache parameter if enabled (for collection and single endpoints)
    if (noCache && endpointType !== 'meta') {
      const separator = url.includes('?') ? '&' : '?'
      url += `${separator}nocache=true`
    }
    
    return url
  }
  
  const directUrl = buildUrl(false)
  const groupUrl = resourceGroup ? buildUrl(true) : null
  
  const handleFetch = async () => {
    setLoading(true)
    setError(null)
    setRequestTime(null)
    
    const startTime = performance.now()
    
    try {
      // Always use direct URL for the actual API call
      const res = await fetch(directUrl)
      const data = await res.json()
      const endTime = performance.now()
      const duration = Math.round(endTime - startTime)
      
      setResponse(data)
      setRequestTime(duration)
    } catch (err: any) {
      const endTime = performance.now()
      const duration = Math.round(endTime - startTime)
      
      setError(err.message)
      setRequestTime(duration)
    } finally {
      setLoading(false)
    }
  }
  
  // Update selected resource when group changes
  const handleGroupChange = (newGroup: string) => {
    setSelectedGroup(newGroup)
    setResourceSearch('') // Reset search
    // Reset to first resource in the new group
    if (newGroup && groups[newGroup]?.length > 0) {
      setSelectedResource(groups[newGroup][0])
    } else if (!newGroup) {
      setSelectedResource(resources[0] || '')
    }
  }
  
  // Handle group selection
  const handleGroupSelect = (group: string) => {
    handleGroupChange(group)
    setGroupSearch('')
    setIsGroupDropdownOpen(false)
  }
  
  // Handle resource selection
  const handleResourceSelect = (resource: string) => {
    setSelectedResource(resource)
    setResourceSearch('')
    setIsResourceDropdownOpen(false)
  }
  
  return (
    <div className="grid lg:grid-cols-2 gap-8">
      {/* Request Builder */}
      <div className="space-y-6">
        <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
          <h2 className="text-2xl font-bold text-white mb-6">Request</h2>
          
          {/* Group & Resource Selection - Single Row */}
          <div className="grid grid-cols-2 gap-4 mb-4">
            {/* Group Selection (Optional) - Searchable */}
            <div ref={groupDropdownRef} className="relative">
              <label className="block text-slate-300 mb-2 font-medium text-sm">
                Filter by Group (Optional)
              </label>
              
              {/* Selected Group Display / Search Input */}
              <div className="relative">
                <input
                  type="text"
                  value={isGroupDropdownOpen ? groupSearch : (selectedGroup || 'All Resources')}
                  onChange={(e) => {
                    setGroupSearch(e.target.value)
                    setIsGroupDropdownOpen(true)
                  }}
                  onFocus={() => setIsGroupDropdownOpen(true)}
                  placeholder="Search groups..."
                  className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm pr-8"
                />
                <button
                  type="button"
                  onClick={() => setIsGroupDropdownOpen(!isGroupDropdownOpen)}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>
              
              {/* Dropdown List */}
              {isGroupDropdownOpen && (
                <div className="absolute z-10 w-full mt-1 bg-slate-700 border border-slate-600 rounded shadow-lg max-h-60 overflow-auto">
                  {/* All Resources Option */}
                  <button
                    type="button"
                    onClick={() => handleGroupSelect('')}
                    className={`w-full text-left px-3 py-2 text-sm hover:bg-slate-600 transition ${
                      !selectedGroup ? 'bg-primary-500/20 text-primary-400' : 'text-white'
                    }`}
                  >
                    All Resources
                  </button>
                  
                  {/* Filtered Groups */}
                  {filteredGroups.length > 0 ? (
                    filteredGroups.map(group => (
                      <button
                        key={group}
                        type="button"
                        onClick={() => handleGroupSelect(group)}
                        className={`w-full text-left px-3 py-2 text-sm hover:bg-slate-600 transition ${
                          group === selectedGroup ? 'bg-primary-500/20 text-primary-400' : 'text-white'
                        }`}
                      >
                        {group} <span className="text-slate-400">({groups[group].length})</span>
                      </button>
                    ))
                  ) : (
                    <div className="px-3 py-2 text-sm text-slate-400">
                      No groups found
                    </div>
                  )}
                </div>
              )}
            </div>
            
            {/* Resource Selection - Searchable */}
            <div ref={resourceDropdownRef} className="relative">
              <label className="block text-slate-300 mb-2 font-medium text-sm">
                Resource
                {resourceGroup && !selectedGroup && (
                  <span className="ml-1 text-xs text-slate-400">
                    ({resourceGroup})
                  </span>
                )}
              </label>
              
              {/* Selected Resource Display / Search Input */}
              <div className="relative">
                <input
                  type="text"
                  value={isResourceDropdownOpen ? resourceSearch : selectedResource}
                  onChange={(e) => {
                    setResourceSearch(e.target.value)
                    setIsResourceDropdownOpen(true)
                  }}
                  onFocus={() => setIsResourceDropdownOpen(true)}
                  placeholder="Search resources..."
                  className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm pr-8"
                />
                <button
                  type="button"
                  onClick={() => setIsResourceDropdownOpen(!isResourceDropdownOpen)}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>
              
              {/* Dropdown List */}
              {isResourceDropdownOpen && (
                <div className="absolute z-10 w-full mt-1 bg-slate-700 border border-slate-600 rounded shadow-lg max-h-60 overflow-auto">
                  {filteredResources.length > 0 ? (
                    filteredResources.map(resource => (
                      <button
                        key={resource}
                        type="button"
                        onClick={() => handleResourceSelect(resource)}
                        className={`w-full text-left px-3 py-2 text-sm hover:bg-slate-600 transition ${
                          resource === selectedResource ? 'bg-primary-500/20 text-primary-400' : 'text-white'
                        }`}
                      >
                        {resource}
                      </button>
                    ))
                  ) : (
                    <div className="px-3 py-2 text-sm text-slate-400">
                      No resources found
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
          
          {/* Endpoint Type - Compact */}
          <div className="mb-4">
            <label className="block text-slate-300 mb-2 font-medium text-sm">
              Endpoint Type
            </label>
            <div className="flex gap-2">
              <button
                onClick={() => setEndpointType('collection')}
                className={`flex-1 px-3 py-2 rounded transition text-sm ${
                  endpointType === 'collection'
                    ? 'bg-primary-500 text-white'
                    : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                }`}
              >
                Collection
              </button>
              <button
                onClick={() => setEndpointType('single')}
                className={`flex-1 px-3 py-2 rounded transition text-sm ${
                  endpointType === 'single'
                    ? 'bg-primary-500 text-white'
                    : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                }`}
              >
                Single
              </button>
              <button
                onClick={() => setEndpointType('meta')}
                className={`flex-1 px-3 py-2 rounded transition text-sm ${
                  endpointType === 'meta'
                    ? 'bg-primary-500 text-white'
                    : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                }`}
              >
                Meta
              </button>
            </div>
          </div>
          
          {/* Options Row - Single Line for Count/ID and Cache */}
          <div className="mb-6">
            <div className="flex gap-4 items-end">
              {/* Item ID (for single) or Count (for collection) */}
              {endpointType === 'single' && (
                <div className="flex-1">
                  <label className="block text-slate-300 mb-2 font-medium text-sm">
                    Item ID
                  </label>
                  <input
                    type="text"
                    value={itemId}
                    onChange={(e) => setItemId(e.target.value)}
                    placeholder="1"
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
                  />
                </div>
              )}
              
              {endpointType === 'collection' && (
                <div className="flex-1">
                  <label className="block text-slate-300 mb-2 font-medium text-sm">
                    Count <span className="text-slate-500 text-xs">(1-100)</span>
                  </label>
                  <input
                    type="number"
                    min="1"
                    max="100"
                    value={count}
                    onChange={(e) => setCount(parseInt(e.target.value) || 1)}
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
                  />
                </div>
              )}
              
              {/* Cache Toggle - Inline */}
              {endpointType !== 'meta' && (
                <div className="flex-1">
                  <label className="block text-slate-300 mb-2 font-medium text-sm">
                    Cache
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer bg-slate-700 px-3 py-2 rounded border border-slate-600 hover:bg-slate-600 transition">
                    <input
                      type="checkbox"
                      checked={noCache}
                      onChange={(e) => setNoCache(e.target.checked)}
                      className="w-4 h-4 rounded border-slate-600 text-primary-500 focus:ring-primary-500 focus:ring-offset-slate-900"
                    />
                    <span className="text-slate-200 text-sm">
                      Bypass
                    </span>
                  </label>
                </div>
              )}
            </div>
          </div>
          
          {/* Send Button - Positioned at top for easy access */}
          <button
            onClick={handleFetch}
            disabled={loading || !selectedResource}
            className="w-full bg-primary-500 hover:bg-primary-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-6 py-3 rounded font-semibold transition mb-6"
          >
            {loading ? 'Loading...' : 'Send Request'}
          </button>
          
          {/* URL Preview - Reference section */}
          <div className="border-t border-slate-700 pt-6">
            <label className="text-slate-300 mb-3 font-medium flex items-center gap-2">
              <span>📋</span>
              Request URL
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
                <strong className="text-slate-200">Cache Bypass:</strong>
                <span className="text-slate-400"> ?nocache=true or ?fresh=true</span>
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
            <div className="space-y-3">
              {requestTime !== null && (
                <div className="text-slate-400 text-sm">
                  ⚡ Request took {requestTime}ms
                </div>
              )}
              <div className="bg-red-900/20 border border-red-500 text-red-300 p-4 rounded">
                <strong>Error:</strong> {error}
              </div>
            </div>
          )}
          
          {response && !error && (
            <div className="space-y-4">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                  <span className="text-green-400 font-semibold">Status: 200 OK</span>
                  {requestTime !== null && (
                    <span className="text-slate-400 text-sm">
                      ⚡ {requestTime}ms
                    </span>
                  )}
                </div>
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

