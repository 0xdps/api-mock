'use client'

import { useState, useRef, useEffect } from 'react'
import { CodeExample } from './CodeExample'
import { getApiUrl } from '@/lib/api'

const API_URL = getApiUrl()

// Helper to pluralize resource names
function pluralize(word: string): string {
  // Common irregular plurals
  const irregulars: Record<string, string> = {
    'person': 'people',
    'child': 'children',
    'tooth': 'teeth',
    'foot': 'feet',
    'mouse': 'mice',
    'goose': 'geese',
  }
  
  if (irregulars[word.toLowerCase()]) {
    return irregulars[word.toLowerCase()]
  }
  
  // Words ending in 'y' preceded by consonant
  if (word.match(/[^aeiou]y$/i)) {
    return word.slice(0, -1) + 'ies'
  }
  
  // Words ending in 's', 'ss', 'sh', 'ch', 'x', 'z'
  if (word.match(/(s|ss|sh|ch|x|z)$/i)) {
    return word + 'es'
  }
  
  // Default: just add 's'
  return word + 's'
}

interface ResourceDocumentationProps {
  resource: string
  schema: any
  group: string
}

export function ResourceDocumentation({ resource, schema, group }: ResourceDocumentationProps) {
  // Playground state
  const [endpointType, setEndpointType] = useState<'collection' | 'single' | 'meta'>('collection')
  const [count, setCount] = useState(10)
  const [itemId, setItemId] = useState('1')
  const [noCache, setNoCache] = useState(false)
  const [response, setResponse] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [requestTime, setRequestTime] = useState<number | null>(null)
  const [showRawSchema, setShowRawSchema] = useState(false)
  const [selectedPathType, setSelectedPathType] = useState<'direct' | 'group'>('direct')
  const responseRef = useRef<HTMLDivElement>(null)
  
  // Pagination state
  const [page, setPage] = useState(1)
  const [limit, setLimit] = useState(10)
  const [offset, setOffset] = useState('')
  
  // Sorting state
  const [sortField, setSortField] = useState('')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc')
  
  // Search state
  const [searchQuery, setSearchQuery] = useState('')
  const [searchFields, setSearchFields] = useState<string[]>([])
  const [searchParam, setSearchParam] = useState<'q' | 'search'>('q')
  
  // Field filtering state
  const [selectedFields, setSelectedFields] = useState<string[]>([])
  
  // Advanced options state
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [delay, setDelay] = useState(0)
  const [flakyRate, setFlakyRate] = useState(0)
  const [skipCache, setSkipCache] = useState(false)
  
  // Build URL based on endpoint type
  const buildUrl = (usePath?: string) => {
    const routes = schema['x-resource']?.routes || {}
    // Use provided path, or schema path, or pluralized resource name
    const primaryPath = usePath || routes.path || `/${pluralize(resource)}`
    
    let url = `${API_URL}${primaryPath}`
    
    if (endpointType === 'single') {
      url += `/${itemId}`
    } else if (endpointType === 'meta') {
      url += '/meta'
    }
    
    const params = new URLSearchParams()
    
    if (endpointType === 'collection') {
      // Add pagination params
      if (offset) {
        params.append('offset', offset)
      } else {
        params.append('page', page.toString())
      }
      params.append('limit', limit.toString())
      
      // Add sorting params
      if (sortField) {
        params.append('sort', sortField)
        params.append('order', sortOrder)
      }
      
      // Add search params
      if (searchQuery) {
        params.append(searchParam, searchQuery)
        if (searchFields.length > 0) {
          params.append('search_fields', searchFields.join(','))
        }
      }
      
      // Add field filtering
      if (selectedFields.length > 0) {
        params.append('fields', selectedFields.join(','))
      }
      
      // Add delay
      if (delay > 0) {
        params.append('delay', delay.toString())
      }
      
      // Add flaky rate
      if (flakyRate > 0) {
        params.append('flakyRate', (flakyRate / 100).toString())
      }
      
      // Add skip cache
      if (skipCache) {
        params.append('skip_cache', 'true')
      }
      
      // Legacy count param (for old behavior)
      params.append('count', count.toString())
    }
    
    if (noCache && endpointType !== 'meta') {
      params.append('nocache', 'true')
    }
    
    const queryString = params.toString()
    return queryString ? `${url}?${queryString}` : url
  }
  
  // Get paths
  const routes = schema['x-resource']?.routes || {}
  const directPath = routes.path || `/${pluralize(resource)}`
  const groupPath = `/${group}/${pluralize(resource)}`
  const aliases = routes.aliases || []
  
  const directUrl = buildUrl(directPath)
  const groupUrl = buildUrl(groupPath)
  
  // Use selected URL for fetching
  const url = selectedPathType === 'direct' ? directUrl : groupUrl
  
  // Auto-scroll to response section when response/error changes
  useEffect(() => {
    if ((response || error) && responseRef.current) {
      responseRef.current.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    }
  }, [response, error])
  
  // Fetch data
  const handleFetch = async () => {
    setLoading(true)
    setError(null)
    setRequestTime(null)
    
    const startTime = performance.now()
    
    try {
      const res = await fetch(url)
      const endTime = performance.now()
      setRequestTime(Math.round(endTime - startTime))
      
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      }
      
      const data = await res.json()
      const cacheHeader = res.headers.get('X-Cache')
      const requestIdHeader = res.headers.get('X-Request-ID')
      
      setResponse({
        data,
        status: res.status,
        cacheStatus: cacheHeader || (endpointType === 'meta' ? 'N/A' : 'MISS'),
        requestId: requestIdHeader || null,
        headers: {
          cache: cacheHeader,
          requestId: requestIdHeader,
        }
      })
    } catch (err: any) {
      setError(err.message)
      setResponse(null)
    } finally {
      setLoading(false)
    }
  }
  
  const properties = schema.properties || {}
  const requiredFields = schema.required || []
  const resourceName = schema['x-resource']?.name || pluralize(resource)
  
  // Get field names for dropdowns
  const fieldNames = Object.keys(properties)
  
  // Toggle field selection
  const toggleField = (field: string) => {
    if (selectedFields.includes(field)) {
      setSelectedFields(selectedFields.filter(f => f !== field))
    } else {
      setSelectedFields([...selectedFields, field])
    }
  }
  
  // Toggle all fields
  const toggleAllFields = () => {
    if (selectedFields.length === fieldNames.length) {
      setSelectedFields([])
    } else {
      setSelectedFields(fieldNames)
    }
  }
  
  const getGroupIcon = (group: string): string => {
    const icons: Record<string, string> = {
      people: '👥',
      business: '💼',
      commerce: '🛒',
      content: '📝',
      social: '💬',
      media: '🎬',
      travel: '✈️',
      location: '🌍',
      finance: '💰',
      food: '🍔',
      education: '🎓',
      sports: '⚽',
      productivity: '✅',
      reference: '📚',
    }
    return icons[group] || '📦'
  }
  
  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <div className="flex items-center gap-3 mb-4">
          <span className="text-3xl">{getGroupIcon(group)}</span>
          <div>
            <h1 className="text-4xl font-bold text-white capitalize">{resourceName}</h1>
            <p className="text-slate-400 text-sm mt-1">
              <span className="capitalize">{group}</span> resource
            </p>
          </div>
        </div>
        <p className="text-lg text-slate-300">
          {schema.description || `Access and manage ${resourceName} data through our RESTful API.`}
        </p>
      </div>
      
      {/* Interactive Playground */}
      <section className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
        <h2 className="text-2xl font-bold text-white mb-4">🎮 Try it Now</h2>
        
        {/* Request Builder */}
        <div className="space-y-4">
          {/* Endpoint Type */}
          <div>
            <label className="block text-slate-300 mb-2 font-medium text-sm">
              Endpoint Type
            </label>
            <div className="flex gap-2">
              <button
                onClick={() => setEndpointType('collection')}
                className={`flex-1 px-4 py-2 rounded transition text-sm ${
                  endpointType === 'collection'
                    ? 'bg-primary-500 text-white'
                    : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                }`}
              >
                Collection
              </button>
              <button
                onClick={() => setEndpointType('single')}
                className={`flex-1 px-4 py-2 rounded transition text-sm ${
                  endpointType === 'single'
                    ? 'bg-primary-500 text-white'
                    : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                }`}
              >
                Single Item
              </button>
              <button
                onClick={() => setEndpointType('meta')}
                className={`flex-1 px-4 py-2 rounded transition text-sm ${
                  endpointType === 'meta'
                    ? 'bg-primary-500 text-white'
                    : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                }`}
              >
                Schema
              </button>
            </div>
          </div>
          
          {/* Pagination (Collection only) */}
          {endpointType === 'collection' && (
            <div className="space-y-4 p-4 bg-slate-900/50 rounded border border-slate-700">
              <h3 className="text-sm font-semibold text-white">Pagination</h3>
              <div className="grid grid-cols-3 gap-3">
                <div>
                  <label className="block text-slate-300 mb-2 text-xs">
                    Page
                  </label>
                  <input
                    type="number"
                    min="1"
                    value={page}
                    onChange={(e) => setPage(parseInt(e.target.value) || 1)}
                    disabled={!!offset}
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm disabled:opacity-50"
                  />
                </div>
                <div>
                  <label className="block text-slate-300 mb-2 text-xs">
                    Limit
                  </label>
                  <select
                    value={limit}
                    onChange={(e) => setLimit(parseInt(e.target.value))}
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
                  >
                    <option value="5">5</option>
                    <option value="10">10</option>
                    <option value="25">25</option>
                    <option value="50">50</option>
                    <option value="100">100</option>
                  </select>
                </div>
                <div>
                  <label className="block text-slate-300 mb-2 text-xs">
                    Offset (optional)
                  </label>
                  <input
                    type="number"
                    min="0"
                    value={offset}
                    onChange={(e) => setOffset(e.target.value)}
                    placeholder="Auto"
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm placeholder-slate-500"
                  />
                </div>
              </div>
            </div>
          )}
          
          {/* Sorting (Collection only) */}
          {endpointType === 'collection' && (
            <div className="space-y-4 p-4 bg-slate-900/50 rounded border border-slate-700">
              <h3 className="text-sm font-semibold text-white">Sorting</h3>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-slate-300 mb-2 text-xs">
                    Sort by
                  </label>
                  <select
                    value={sortField}
                    onChange={(e) => setSortField(e.target.value)}
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
                  >
                    <option value="">None</option>
                    {fieldNames.map(field => (
                      <option key={field} value={field}>{field}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-slate-300 mb-2 text-xs">
                    Order
                  </label>
                  <select
                    value={sortOrder}
                    onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
                    disabled={!sortField}
                    className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm disabled:opacity-50"
                  >
                    <option value="asc">↑ Ascending</option>
                    <option value="desc">↓ Descending</option>
                  </select>
                </div>
              </div>
            </div>
          )}
          
          {/* Search (Collection only) */}
          {endpointType === 'collection' && (
            <div className="space-y-4 p-4 bg-slate-900/50 rounded border border-slate-700">
              <h3 className="text-sm font-semibold text-white">Search</h3>
              <div className="space-y-3">
                <div>
                  <label className="block text-slate-300 mb-2 text-xs">
                    Query
                  </label>
                  <div className="relative">
                    <input
                      type="text"
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      placeholder="Search..."
                      className="w-full bg-slate-700 text-white px-3 py-2 pr-8 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm placeholder-slate-500"
                    />
                    {searchQuery && (
                      <button
                        onClick={() => setSearchQuery('')}
                        className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                      >
                        ✕
                      </button>
                    )}
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-slate-300 mb-2 text-xs">
                      Search in fields
                    </label>
                    <select
                      multiple
                      value={searchFields}
                      onChange={(e) => setSearchFields(Array.from(e.target.selectedOptions, option => option.value))}
                      className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
                      size={3}
                    >
                      {fieldNames.map(field => (
                        <option key={field} value={field}>{field}</option>
                      ))}
                    </select>
                    <p className="text-xs text-slate-500 mt-1">Cmd/Ctrl+click to select multiple</p>
                  </div>
                  <div>
                    <label className="block text-slate-300 mb-2 text-xs">
                      Parameter name
                    </label>
                    <div className="space-y-2">
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="radio"
                          checked={searchParam === 'q'}
                          onChange={() => setSearchParam('q')}
                          className="text-primary-500"
                        />
                        <span className="text-slate-300 text-sm">q</span>
                      </label>
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="radio"
                          checked={searchParam === 'search'}
                          onChange={() => setSearchParam('search')}
                          className="text-primary-500"
                        />
                        <span className="text-slate-300 text-sm">search</span>
                      </label>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}
          
          {/* Field Filtering (Collection only) */}
          {endpointType === 'collection' && (
            <div className="space-y-4 p-4 bg-slate-900/50 rounded border border-slate-700">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-white">Field Filtering</h3>
                <button
                  onClick={toggleAllFields}
                  className="text-xs text-primary-400 hover:text-primary-300"
                >
                  {selectedFields.length === fieldNames.length ? 'Deselect All' : 'Select All'}
                </button>
              </div>
              <div className="grid grid-cols-3 sm:grid-cols-4 gap-2 max-h-32 overflow-y-auto">
                {fieldNames.map(field => (
                  <label key={field} className="flex items-center gap-2 cursor-pointer text-sm">
                    <input
                      type="checkbox"
                      checked={selectedFields.includes(field)}
                      onChange={() => toggleField(field)}
                      className="rounded border-slate-600 text-primary-500 focus:ring-primary-500"
                    />
                    <span className="text-slate-300 text-xs truncate">{field}</span>
                  </label>
                ))}
              </div>
              {selectedFields.length > 0 && (
                <p className="text-xs text-slate-400">
                  Selected: {selectedFields.join(', ')}
                </p>
              )}
            </div>
          )}
          
          {/* Advanced Options (Collection only) */}
          {endpointType === 'collection' && (
            <div className="border border-slate-700 rounded">
              <button
                onClick={() => setShowAdvanced(!showAdvanced)}
                className="w-full flex items-center justify-between p-4 bg-slate-900/30 hover:bg-slate-900/50 transition"
              >
                <h3 className="text-sm font-semibold text-white">
                  {showAdvanced ? '▼' : '▶'} Advanced Options
                </h3>
                <span className="text-xs text-slate-400">
                  Middleware parameters
                </span>
              </button>
              
              {showAdvanced && (
                <div className="p-4 space-y-4 border-t border-slate-700">
                  {/* Delay */}
                  <div>
                    <label className="block text-slate-300 mb-2 text-xs">
                      Delay: {delay}ms {delay >= 5000 && '⚠️'}
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
                    <div className="flex justify-between text-xs text-slate-500 mt-1">
                      <span>0ms</span>
                      <span>5000ms</span>
                    </div>
                  </div>
                  
                  {/* Flaky Rate */}
                  <div>
                    <label className="block text-slate-300 mb-2 text-xs">
                      Flaky Rate: {flakyRate}% {flakyRate > 50 && '⚠️ High failure rate!'}
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
                    <div className="flex justify-between text-xs text-slate-500 mt-1">
                      <span>0%</span>
                      <span>100%</span>
                    </div>
                  </div>
                  
                  {/* Cache Controls */}
                  <div className="flex gap-4">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={skipCache}
                        onChange={(e) => setSkipCache(e.target.checked)}
                        className="rounded border-slate-600 text-primary-500"
                      />
                      <span className="text-slate-300 text-xs">Skip Cache</span>
                    </label>
                  </div>
                </div>
              )}
            </div>
          )}
          
          {/* Options */}
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
            
            {/* Cache Toggle */}
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
          
          {/* URL Preview with Selection */}
          <div className="space-y-3">
            <div>
              <label className="block text-slate-300 mb-2 font-medium text-sm">
                Request URL
              </label>
              
              <div className="space-y-2">
                {/* Direct Path - Clickable Box */}
                <button
                  type="button"
                  onClick={() => setSelectedPathType('direct')}
                  className={`w-full text-left p-3 rounded transition cursor-pointer ${
                    selectedPathType === 'direct'
                      ? 'bg-slate-900 border-2 border-primary-500'
                      : 'bg-slate-900 border border-slate-600 hover:border-slate-500'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs text-slate-400">Direct Access</span>
                    {selectedPathType === 'direct' && (
                      <span className="text-xs text-primary-400 font-medium">✓ Active</span>
                    )}
                  </div>
                  <code className={`text-sm break-all ${
                    selectedPathType === 'direct' ? 'text-primary-400' : 'text-slate-500'
                  }`}>
                    {directUrl}
                  </code>
                </button>
                
                {/* Group Path - Clickable Box */}
                <button
                  type="button"
                  onClick={() => setSelectedPathType('group')}
                  className={`w-full text-left p-3 rounded transition cursor-pointer ${
                    selectedPathType === 'group'
                      ? 'bg-slate-900 border-2 border-primary-500'
                      : 'bg-slate-900 border border-slate-600 hover:border-slate-500'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs text-slate-400">Via Group Path</span>
                    {selectedPathType === 'group' && (
                      <span className="text-xs text-primary-400 font-medium">✓ Active</span>
                    )}
                  </div>
                  <code className={`text-sm break-all ${
                    selectedPathType === 'group' ? 'text-primary-400' : 'text-slate-500'
                  }`}>
                    {groupUrl}
                  </code>
                </button>
              </div>
              
              {/* Aliases */}
              {aliases.length > 0 && (
                <div className="mt-3 pt-3 border-t border-slate-700">
                  <div className="text-xs text-slate-400 mb-2">Additional aliases:</div>
                  <div className="flex flex-wrap gap-2">
                    {aliases.map((alias: string) => (
                      <code key={alias} className="text-xs bg-slate-900 text-slate-500 px-2 py-1 rounded border border-slate-700">
                        {alias}
                      </code>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
          
          {/* Send Button */}
          <button
            onClick={handleFetch}
            disabled={loading}
            className="w-full bg-primary-500 hover:bg-primary-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-6 py-3 rounded font-semibold transition"
          >
            {loading ? 'Loading...' : '▶ Send Request'}
          </button>
          
          {/* Response */}
          {(response || error) && (
            <div ref={responseRef} className="border-t border-slate-700 pt-4 mt-4">
              {error ? (
                <div className="bg-red-900/20 border border-red-500/50 p-4 rounded">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-red-400 font-semibold">❌ Error</span>
                  </div>
                  <p className="text-red-300 text-sm">{error}</p>
                </div>
              ) : (
                <div>
                  <div className="flex items-center justify-between mb-3 flex-wrap gap-2">
                    <div className="flex items-center gap-3 flex-wrap">
                      <span className="text-green-400 font-semibold">✓ Success</span>
                      <span className="text-slate-400 text-sm">
                        Status: <span className="text-white">{response.status}</span>
                      </span>
                      {response.cacheStatus !== 'N/A' && (
                        <span className="text-slate-400 text-sm">
                          Cache: <span className={
                            response.cacheStatus === 'HIT' ? 'text-green-400' : 
                            response.cacheStatus === 'BYPASS' ? 'text-yellow-400' : 
                            'text-slate-400'
                          }>
                            {response.cacheStatus}
                          </span>
                        </span>
                      )}
                      {requestTime !== null && (
                        <span className="text-slate-400 text-sm">
                          Time: <span className="text-white">{requestTime}ms</span>
                        </span>
                      )}
                      {response.requestId && (
                        <span className="text-slate-400 text-sm">
                          ID: <span className="text-slate-300 font-mono text-xs">{response.requestId}</span>
                        </span>
                      )}
                    </div>
                  </div>
                  
                  {/* Applied Middleware Indicators */}
                  {endpointType === 'collection' && (delay > 0 || flakyRate > 0 || skipCache || selectedFields.length > 0 || sortField || searchQuery) && (
                    <div className="mb-3 p-3 bg-blue-900/20 border border-blue-500/30 rounded">
                      <p className="text-xs text-blue-300 font-semibold mb-2">Applied Middleware:</p>
                      <div className="flex flex-wrap gap-2">
                        {delay > 0 && (
                          <span className="text-xs bg-blue-900/50 text-blue-200 px-2 py-1 rounded border border-blue-500/30">
                            ⏱️ Delayed {delay}ms
                          </span>
                        )}
                        {flakyRate > 0 && (
                          <span className="text-xs bg-yellow-900/50 text-yellow-200 px-2 py-1 rounded border border-yellow-500/30">
                            🎲 Flaky {flakyRate}%
                          </span>
                        )}
                        {skipCache && (
                          <span className="text-xs bg-purple-900/50 text-purple-200 px-2 py-1 rounded border border-purple-500/30">
                            🚫 Cache Bypassed
                          </span>
                        )}
                        {selectedFields.length > 0 && (
                          <span className="text-xs bg-green-900/50 text-green-200 px-2 py-1 rounded border border-green-500/30">
                            🔍 Fields: {selectedFields.join(', ')}
                          </span>
                        )}
                        {sortField && (
                          <span className="text-xs bg-indigo-900/50 text-indigo-200 px-2 py-1 rounded border border-indigo-500/30">
                            🔀 Sort: {sortField} ({sortOrder})
                          </span>
                        )}
                        {searchQuery && (
                          <span className="text-xs bg-pink-900/50 text-pink-200 px-2 py-1 rounded border border-pink-500/30">
                            🔎 Search: "{searchQuery}"
                            {searchFields.length > 0 && ` in ${searchFields.join(', ')}`}
                          </span>
                        )}
                      </div>
                    </div>
                  )}
                  
                  {/* Pagination Metadata */}
                  {response.data?.pagination && (
                    <div className="mb-3 p-3 bg-slate-900 border border-slate-600 rounded">
                      <p className="text-xs text-slate-400 font-semibold mb-2">Pagination:</p>
                      <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs">
                        <div>
                          <span className="text-slate-500">Page:</span>{' '}
                          <span className="text-white">{response.data.pagination.page}</span>
                        </div>
                        <div>
                          <span className="text-slate-500">Limit:</span>{' '}
                          <span className="text-white">{response.data.pagination.limit}</span>
                        </div>
                        <div>
                          <span className="text-slate-500">Total:</span>{' '}
                          <span className="text-white">{response.data.pagination.total}</span>
                        </div>
                        <div>
                          <span className="text-slate-500">Pages:</span>{' '}
                          <span className="text-white">{response.data.pagination.total_pages}</span>
                        </div>
                      </div>
                      <p className="text-xs text-slate-400 mt-2">
                        Showing {((response.data.pagination.page - 1) * response.data.pagination.limit) + 1}-
                        {Math.min(response.data.pagination.page * response.data.pagination.limit, response.data.pagination.total)} of {response.data.pagination.total} items
                      </p>
                    </div>
                  )}
                  
                  <div className="bg-slate-900 p-4 rounded border border-slate-600 overflow-auto max-h-96">
                    <pre className="text-green-400 text-sm">
                      <code>{JSON.stringify(response.data, null, 2)}</code>
                    </pre>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </section>
      
      {/* Endpoints */}
      <section className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
        <h2 className="text-2xl font-bold text-white mb-4">Endpoints</h2>
        
        <div className="space-y-4">
          {/* Collection Endpoint */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <span className="bg-green-500 text-white px-2 py-1 rounded text-xs font-semibold">GET</span>
              <code className="text-primary-500">{directPath}</code>
            </div>
            <p className="text-slate-400 text-sm mb-2">Get a collection of {resourceName}</p>
            <div className="text-xs text-slate-500 mb-1">
              <span className="font-semibold">Parameters:</span> count, seed, nocache
            </div>
            <div className="text-xs text-slate-400">
              Alternative: <code className="text-slate-500">{groupPath}</code>
            </div>
          </div>
          
          {/* Single Item Endpoint */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <span className="bg-green-500 text-white px-2 py-1 rounded text-xs font-semibold">GET</span>
              <code className="text-primary-500">{directPath}/:id</code>
            </div>
            <p className="text-slate-400 text-sm mb-2">Get a single {resource} by ID</p>
            <div className="text-xs text-slate-500 mb-1">
              <span className="font-semibold">Parameters:</span> nocache
            </div>
            <div className="text-xs text-slate-400">
              Alternative: <code className="text-slate-500">{groupPath}/:id</code>
            </div>
          </div>
          
          {/* Meta Endpoint */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <span className="bg-green-500 text-white px-2 py-1 rounded text-xs font-semibold">GET</span>
              <code className="text-primary-500">{directPath}/meta</code>
            </div>
            <p className="text-slate-400 text-sm mb-1">Get resource metadata and schema</p>
            <div className="text-xs text-slate-400">
              Alternative: <code className="text-slate-500">{groupPath}/meta</code>
            </div>
          </div>
          
          {/* Aliases */}
          {aliases.length > 0 && (
            <div className="mt-4 pt-4 border-t border-slate-700">
              <p className="text-sm text-slate-400 mb-2">
                <span className="font-semibold">Additional aliases:</span>
              </p>
              <div className="flex flex-wrap gap-2">
                {aliases.map((alias: string) => (
                  <code key={alias} className="text-xs bg-slate-900 text-slate-400 px-2 py-1 rounded">
                    {alias}
                  </code>
                ))}
              </div>
            </div>
          )}
        </div>
      </section>
      
      {/* Code Examples */}
      <section>
        <h2 className="text-2xl font-bold text-white mb-4">Code Examples</h2>
        
        <div className="grid md:grid-cols-2 gap-4">
          {/* JavaScript */}
          <div>
            <h3 className="text-lg font-semibold text-white mb-3">JavaScript</h3>
            <CodeExample 
              title="Fetch Collection"
              code={`// Get 10 ${resourceName}
fetch('${API_URL}${directPath}?count=10')
  .then(res => res.json())
  .then(data => console.log(data));

// Get single item
fetch('${API_URL}${directPath}/1')
  .then(res => res.json())
  .then(data => console.log(data));`}
              language="javascript"
            />
          </div>
          
          {/* Python */}
          <div>
            <h3 className="text-lg font-semibold text-white mb-3">Python</h3>
            <CodeExample 
              title="Fetch with Requests"
              code={`import requests

# Get collection
response = requests.get(
    '${API_URL}${directPath}?count=10'
)
data = response.json()

# Get single item
item = requests.get(
    '${API_URL}${directPath}/1'
).json()`}
              language="python"
            />
          </div>
          
          {/* cURL */}
          <div>
            <h3 className="text-lg font-semibold text-white mb-3">cURL</h3>
            <CodeExample 
              title="Command Line"
              code={`# Get collection
curl "${API_URL}${directPath}?count=10"

# Get single item
curl "${API_URL}${directPath}/1"

# Get schema
curl "${API_URL}${directPath}/meta"`}
              language="bash"
            />
          </div>
          
          {/* Fresh Data */}
          <div>
            <h3 className="text-lg font-semibold text-white mb-3">Bypass Cache</h3>
            <CodeExample 
              title="Get Fresh Data"
              code={`// Bypass cache for fresh data
fetch('${API_URL}${directPath}?count=10&nocache=true')
  .then(res => {
    console.log('Cache:', res.headers.get('X-Cache'));
    return res.json();
  })
  .then(data => console.log(data));`}
              language="javascript"
            />
          </div>
        </div>
      </section>
      
      {/* Schema Properties */}
      <section className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold text-white">Schema Properties</h2>
        </div>
        
        {Object.keys(properties).length > 0 ? (
          <div className="space-y-4">
            {Object.entries(properties).map(([key, value]: [string, any]) => (
              <div key={key} className="bg-slate-900/50 p-4 rounded-lg border border-slate-700">
                <div className="flex items-start justify-between mb-2">
                  <div className="flex items-baseline gap-2 flex-wrap">
                    <code className="text-primary-500 font-semibold text-base">{key}</code>
                    <span className="text-xs bg-slate-700 text-slate-300 px-2 py-0.5 rounded uppercase font-medium">
                      {value.type || 'any'}
                    </span>
                    {requiredFields.includes(key) && (
                      <span className="text-xs bg-red-500/20 text-red-400 px-2 py-0.5 rounded">
                        required
                      </span>
                    )}
                  </div>
                </div>
                
                {value.description && (
                  <p className="text-sm text-slate-400 mb-2">{value.description}</p>
                )}
                
                <div className="flex flex-wrap gap-2 text-xs">
                  {(value['x-generator'] || value['x-faker']) && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded flex items-center gap-1">
                      <span className="text-slate-500">Generator:</span>
                      <code className="text-primary-400">{value['x-generator'] || value['x-faker']}</code>
                    </span>
                  )}
                  {value['x-generator-params'] && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded flex items-center gap-1">
                      <span className="text-slate-500">Params:</span>
                      <code className="text-primary-400 text-xs">{JSON.stringify(value['x-generator-params'])}</code>
                    </span>
                  )}
                  {value.format && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded flex items-center gap-1">
                      <span className="text-slate-500">Format:</span>
                      <code className="text-primary-400">{value.format}</code>
                    </span>
                  )}
                  {value.enum && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded flex items-center gap-1">
                      <span className="text-slate-500">Enum:</span>
                      <code className="text-primary-400">{value.enum.join(', ')}</code>
                    </span>
                  )}
                  {value.minimum !== undefined && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded">
                      Min: {value.minimum}
                    </span>
                  )}
                  {value.maximum !== undefined && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded">
                      Max: {value.maximum}
                    </span>
                  )}
                  {value.minLength !== undefined && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded">
                      Min Length: {value.minLength}
                    </span>
                  )}
                  {value.maxLength !== undefined && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded">
                      Max Length: {value.maxLength}
                    </span>
                  )}
                  {value.pattern && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded flex items-center gap-1">
                      <span className="text-slate-500">Pattern:</span>
                      <code className="text-primary-400 text-xs">{value.pattern}</code>
                    </span>
                  )}
                  {value.example !== undefined && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-1 rounded flex items-center gap-1">
                      <span className="text-slate-500">Example:</span>
                      <code className="text-primary-400">{JSON.stringify(value.example)}</code>
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="bg-slate-900/50 p-8 rounded-lg border border-slate-700 text-center">
            <p className="text-slate-400 mb-4">No properties defined in schema</p>
            <button
              onClick={() => setShowRawSchema(true)}
              className="inline-block bg-primary-500 hover:bg-primary-600 text-white px-4 py-2 rounded text-sm transition"
            >
              View Raw Schema
            </button>
          </div>
        )}
        
        {/* Raw Schema Toggle */}
        <div className="mt-4 pt-4 border-t border-slate-700">
          <div className="flex items-center justify-between mb-3">
            <button
              onClick={() => setShowRawSchema(!showRawSchema)}
              className="flex items-center gap-2 text-sm text-slate-400 hover:text-white transition"
            >
              <svg 
                className={`w-4 h-4 transition-transform ${showRawSchema ? 'rotate-90' : ''}`}
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
              {showRawSchema ? 'Hide' : 'Show'} Complete JSON Schema
            </button>
            {showRawSchema && (
              <button
                onClick={() => {
                  navigator.clipboard.writeText(JSON.stringify(schema, null, 2))
                }}
                className="text-xs text-primary-500 hover:text-primary-400 transition"
              >
                Copy Schema
              </button>
            )}
          </div>
          
          {showRawSchema && (
            <div className="mt-3 bg-slate-900 p-4 rounded border border-slate-600 overflow-auto max-h-[600px]">
              <pre className="text-green-400 text-xs">
                <code>{JSON.stringify(schema, null, 2)}</code>
              </pre>
            </div>
          )}
        </div>
      </section>
      
      {/* Query Parameters */}
      <section className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
        <h2 className="text-2xl font-bold text-white mb-4">Query Parameters</h2>
        
        <div className="space-y-4">
          <div>
            <code className="text-primary-500 font-semibold">count</code>
            <span className="text-xs text-slate-500 ml-2">integer</span>
            <p className="text-sm text-slate-400 mt-1">
              Number of items to return (default: 10, max: 100)
            </p>
          </div>
          
          <div>
            <code className="text-primary-500 font-semibold">seed</code>
            <span className="text-xs text-slate-500 ml-2">integer</span>
            <p className="text-sm text-slate-400 mt-1">
              Seed for reproducible data generation
            </p>
          </div>
          
          <div>
            <code className="text-primary-500 font-semibold">nocache</code>
            <span className="text-xs text-slate-500 ml-2">boolean</span>
            <p className="text-sm text-slate-400 mt-1">
              Bypass cache and generate fresh data on every request (aliases: fresh, _nocache)
            </p>
          </div>
        </div>
      </section>
      
    </div>
  )
}

