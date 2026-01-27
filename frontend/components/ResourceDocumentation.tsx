'use client'

import { useState, useRef, useEffect } from 'react'
import { CodeExample } from './CodeExample'
import { getApiUrl, apiClient } from '@/lib/api'

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
  
  // Active tab for options
  const [activeTab, setActiveTab] = useState<'quick' | 'pagination' | 'filters' | 'middleware' | 'path'>('quick')
  
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
  const [delay, setDelay] = useState(0)
  const [flakyRate, setFlakyRate] = useState(0)
  const [skipCache, setSkipCache] = useState(false)
  
  // Code examples dropdown
  const [showExamples, setShowExamples] = useState(false)
  
  // Advanced options collapse state
  const [showAdvanced, setShowAdvanced] = useState(false)
  
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
      const res = await apiClient.get(url)
      const endTime = performance.now()
      setRequestTime(Math.round(endTime - startTime))

      if (res.isError()) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      }

      const data = res.json()
      const headers = res.headers
      const cacheHeader = headers['x-cache'] || headers['X-Cache']
      const requestIdHeader = headers['x-request-id'] || headers['X-Request-ID']
      
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
    <div className="space-y-4 px-3 sm:px-4 lg:px-6">
      {/* Compact Header */}
      <div className="flex items-center gap-2">
        <span className="text-2xl">{getGroupIcon(group)}</span>
        <div>
          <h1 className="text-2xl font-bold text-white capitalize">{resourceName}</h1>
          <p className="text-slate-400 text-xs">
            {schema.description || `Access and manage ${resourceName} data through our RESTful API.`}
          </p>
        </div>
      </div>
      
      {/* Two-Column Layout: Request Builder (50%) + Response (50%) */}
      <div className="grid lg:grid-cols-2 gap-4">
        {/* LEFT COLUMN: Request Builder (50%) */}
        <div>
            <div className="bg-slate-800/50 backdrop-blur p-4 rounded-lg border border-slate-700 h-[calc(100vh-10rem)] overflow-y-auto">
              <h2 className="text-base font-bold text-white mb-3 flex items-center gap-2">
                <span>🎯</span> Request Builder
              </h2>
              
              {/* Endpoint Type Selector */}
              <div className="mb-3">
                <label className="block text-slate-300 mb-2 text-xs font-medium">
                  Endpoint Type
                </label>
                <div className="grid grid-cols-3 gap-2">
                  <button
                    onClick={() => setEndpointType('collection')}
                    className={`px-3 py-2 rounded text-sm font-medium transition ${
                      endpointType === 'collection'
                        ? 'bg-primary-500 text-white'
                        : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                    }`}
                  >
                    Collection
                  </button>
                  <button
                    onClick={() => setEndpointType('single')}
                    className={`px-3 py-2 rounded text-sm font-medium transition ${
                      endpointType === 'single'
                        ? 'bg-primary-500 text-white'
                        : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                    }`}
                  >
                    Single
                  </button>
                  <button
                    onClick={() => setEndpointType('meta')}
                    className={`px-3 py-2 rounded text-sm font-medium transition ${
                      endpointType === 'meta'
                        ? 'bg-primary-500 text-white'
                        : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
                    }`}
                  >
                    Schema
                  </button>
                </div>
              </div>
              
              {/* Tab Navigation */}
              {endpointType === 'collection' && (
                <div className="mb-3">
                  <div className="flex flex-wrap gap-1 bg-slate-900/50 p-1 rounded">
                    <button
                      onClick={() => setActiveTab('quick')}
                      className={`flex-1 px-2 py-1.5 rounded text-xs font-medium transition ${
                        activeTab === 'quick'
                          ? 'bg-slate-700 text-white'
                          : 'text-slate-400 hover:text-slate-200'
                      }`}
                    >
                      Quick
                    </button>
                    <button
                      onClick={() => setActiveTab('pagination')}
                      className={`flex-1 px-2 py-1.5 rounded text-xs font-medium transition ${
                        activeTab === 'pagination'
                          ? 'bg-slate-700 text-white'
                          : 'text-slate-400 hover:text-slate-200'
                      }`}
                    >
                      Pagination
                    </button>
                    <button
                      onClick={() => setActiveTab('filters')}
                      className={`flex-1 px-2 py-1.5 rounded text-xs font-medium transition ${
                        activeTab === 'filters'
                          ? 'bg-slate-700 text-white'
                          : 'text-slate-400 hover:text-slate-200'
                      }`}
                    >
                      Filters
                    </button>
                    <button
                      onClick={() => setActiveTab('middleware')}
                      className={`flex-1 px-2 py-1.5 rounded text-xs font-medium transition ${
                        activeTab === 'middleware'
                          ? 'bg-slate-700 text-white'
                          : 'text-slate-400 hover:text-slate-200'
                      }`}
                    >
                      Middleware
                    </button>
                    <button
                      onClick={() => setActiveTab('path')}
                      className={`flex-1 px-2 py-1.5 rounded text-xs font-medium transition ${
                        activeTab === 'path'
                          ? 'bg-slate-700 text-white'
                          : 'text-slate-400 hover:text-slate-200'
                      }`}
                    >
                      Path
                    </button>
                  </div>
                </div>
              )}
              
              {/* Tab Content */}
              <div className="space-y-3">
                {/* QUICK TAB */}
                {endpointType === 'collection' && activeTab === 'quick' && (
                  <div className="space-y-3 p-3 bg-slate-900/50 rounded border border-slate-700">
                    <div>
                      <label className="block text-slate-300 mb-1.5 text-xs font-medium">
                        Count <span className="text-slate-500">(1-100)</span>
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
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={skipCache}
                        onChange={(e) => setSkipCache(e.target.checked)}
                        className="rounded border-slate-600 text-primary-500 focus:ring-primary-500"
                      />
                      <span className="text-slate-300 text-xs">Skip Cache</span>
                    </label>
                  </div>
                )}
                
                {/* SINGLE ITEM ID */}
                {endpointType === 'single' && (
                  <div className="space-y-3 p-3 bg-slate-900/50 rounded border border-slate-700">
                    <div>
                      <label className="block text-slate-300 mb-1.5 text-xs font-medium">
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
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={noCache}
                        onChange={(e) => setNoCache(e.target.checked)}
                        className="rounded border-slate-600 text-primary-500 focus:ring-primary-500"
                      />
                      <span className="text-slate-300 text-xs">Bypass Cache</span>
                    </label>
                  </div>
                )}
                
                {/* PAGINATION TAB */}
                {endpointType === 'collection' && activeTab === 'pagination' && (
                  <div className="space-y-3 p-3 bg-slate-900/50 rounded border border-slate-700">
                    <div className="grid grid-cols-3 gap-2">
                      <div>
                        <label className="block text-slate-300 mb-1.5 text-xs">
                          Page
                        </label>
                        <input
                          type="number"
                          min="1"
                          value={page}
                          onChange={(e) => setPage(parseInt(e.target.value) || 1)}
                          disabled={!!offset}
                          className="w-full bg-slate-700 text-white px-2 py-1.5 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm disabled:opacity-50"
                        />
                      </div>
                      <div>
                        <label className="block text-slate-300 mb-1.5 text-xs">
                          Limit
                        </label>
                        <select
                          value={limit}
                          onChange={(e) => setLimit(parseInt(e.target.value))}
                          className="w-full bg-slate-700 text-white px-2 py-1.5 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
                        >
                          <option value="5">5</option>
                          <option value="10">10</option>
                          <option value="25">25</option>
                          <option value="50">50</option>
                          <option value="100">100</option>
                        </select>
                      </div>
                      <div>
                        <label className="block text-slate-300 mb-1.5 text-xs">
                          Offset
                        </label>
                        <input
                          type="number"
                          min="0"
                          value={offset}
                          onChange={(e) => setOffset(e.target.value)}
                          placeholder="Auto"
                          className="w-full bg-slate-700 text-white px-2 py-1.5 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm placeholder-slate-500"
                        />
                      </div>
                    </div>
                  </div>
                )}
                
                {/* FILTERS TAB */}
                {endpointType === 'collection' && activeTab === 'filters' && (
                  <div className="space-y-3 p-3 bg-slate-900/50 rounded border border-slate-700 max-h-96 overflow-y-auto">
                    {/* Search */}
                    <div>
                      <label className="block text-slate-300 mb-1.5 text-xs font-medium">
                        Search Query
                      </label>
                      <input
                        type="text"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        placeholder="Search..."
                        className="w-full bg-slate-700 text-white px-3 py-1.5 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm placeholder-slate-500"
                      />
                    </div>
                    
                    {searchQuery && (
                      <div>
                        <label className="block text-slate-300 mb-1.5 text-xs font-medium">
                          Search Fields (optional)
                        </label>
                        <select
                          multiple
                          value={searchFields}
                          onChange={(e) => setSearchFields(Array.from(e.target.selectedOptions, option => option.value))}
                          className="w-full bg-slate-700 text-white px-2 py-1.5 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-xs"
                          size={3}
                        >
                          {fieldNames.map(field => (
                            <option key={field} value={field}>
                              {field}
                            </option>
                          ))}
                        </select>
                      </div>
                    )}
                    
                    {/* Sort */}
                    <div className="grid grid-cols-2 gap-2">
                      <div>
                        <label className="block text-slate-300 mb-1.5 text-xs">
                          Sort by
                        </label>
                        <select
                          value={sortField}
                          onChange={(e) => setSortField(e.target.value)}
                          className="w-full bg-slate-700 text-white px-2 py-1.5 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
                        >
                          <option value="">None</option>
                          {fieldNames.map(field => (
                            <option key={field} value={field}>{field}</option>
                          ))}
                        </select>
                      </div>
                      <div>
                        <label className="block text-slate-300 mb-1.5 text-xs">
                          Order
                        </label>
                        <select
                          value={sortOrder}
                          onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
                          disabled={!sortField}
                          className="w-full bg-slate-700 text-white px-2 py-1.5 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm disabled:opacity-50"
                        >
                          <option value="asc">↑ Asc</option>
                          <option value="desc">↓ Desc</option>
                        </select>
                      </div>
                    </div>
{/* 
bg-slate-800/50 backdrop-blur p-8 rounded-lg border border-slate-700 text-center h-full flex items-center justify-center              

bg-slate-810/50 backdrop-blur p-4 rounded-lg border border-slate-700

*/}
                    {/* Select Fields */}
                    <div>
                      <div className="flex items-center justify-between mb-1.5">
                        <label className="text-slate-300 text-xs font-medium">
                          Select Fields
                        </label>
                        <button
                          onClick={toggleAllFields}
                          className="text-xs text-primary-400 hover:text-primary-300"
                        >
                          {selectedFields.length === fieldNames.length ? 'None' : 'All'}
                        </button>
                      </div>
                      <div className="grid grid-cols-2 gap-1.5 max-h-32 overflow-y-auto p-2 bg-slate-800 rounded">
                        {fieldNames.map(field => (
                          <label key={field} className="flex items-center gap-1.5 cursor-pointer text-xs">
                            <input
                              type="checkbox"
                              checked={selectedFields.includes(field)}
                              onChange={() => toggleField(field)}
                              className="rounded border-slate-600 text-primary-500 focus:ring-primary-500 w-3 h-3"
                            />
                            <span className="text-slate-300 truncate">{field}</span>
                          </label>
                        ))}
                      </div>
                    </div>
                  </div>
                )}
                
                {/* MIDDLEWARE TAB */}
                {endpointType === 'collection' && activeTab === 'middleware' && (
                  <div className="space-y-3 p-3 bg-slate-900/50 rounded border border-slate-700">
                    <div>
                      <label className="block text-slate-300 mb-1.5 text-xs">
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
                      <div className="flex justify-between text-xs text-slate-500 mt-0.5">
                        <span>0ms</span>
                        <span>5000ms</span>
                      </div>
                    </div>
                    
                    <div>
                      <label className="block text-slate-300 mb-1.5 text-xs">
                        Flaky Rate: {flakyRate}%
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
                      <div className="flex justify-between text-xs text-slate-500 mt-0.5">
                        <span>0%</span>
                        <span>100%</span>
                      </div>
                    </div>
                    
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={skipCache}
                        onChange={(e) => setSkipCache(e.target.checked)}
                        className="rounded border-slate-600 text-primary-500 focus:ring-primary-500"
                      />
                      <span className="text-slate-300 text-xs">Skip Cache</span>
                    </label>
                  </div>
                )}
                
                {/* PATH TAB */}
                {endpointType === 'collection' && activeTab === 'path' && (
                  <div className="space-y-3 p-3 bg-slate-900/50 rounded border border-slate-700">
                    <div className="space-y-2">
                      <label className="flex items-start gap-2 cursor-pointer p-2 bg-slate-800 rounded hover:bg-slate-750">
                        <input
                          type="radio"
                          checked={selectedPathType === 'direct'}
                          onChange={() => setSelectedPathType('direct')}
                          className="mt-0.5"
                        />
                        <div className="flex-1">
                          <div className="text-slate-300 text-xs font-medium mb-1">Direct Access</div>
                          <code className="text-xs text-slate-500 break-all">{directPath}</code>
                        </div>
                      </label>
                      
                      <label className="flex items-start gap-2 cursor-pointer p-2 bg-slate-800 rounded hover:bg-slate-750">
                        <input
                          type="radio"
                          checked={selectedPathType === 'group'}
                          onChange={() => setSelectedPathType('group')}
                          className="mt-0.5"
                        />
                        <div className="flex-1">
                          <div className="text-slate-300 text-xs font-medium mb-1">Via Group</div>
                          <code className="text-xs text-slate-500 break-all">{groupPath}</code>
                        </div>
                      </label>
                    </div>
                    
                    {aliases.length > 0 && (
                      <div className="pt-2 border-t border-slate-700">
                        <div className="text-xs text-slate-400 mb-1.5">Aliases:</div>
                        <div className="flex flex-wrap gap-1">
                          {aliases.map((alias: string) => (
                            <code key={alias} className="text-xs bg-slate-800 text-slate-500 px-1.5 py-0.5 rounded">
                              {alias}
                            </code>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                )}
              </div>
              
              {/* Request URL Preview */}
              <div className="mt-3 pt-3 border-t border-slate-700">
                <label className="block text-slate-300 mb-1.5 text-xs font-medium">
                  Request URL
                </label>
                <div className="flex items-start gap-2">
                  <code className="flex-1 bg-slate-900 text-primary-400 px-3 py-2 rounded border border-slate-600 text-xs break-all">
                    {url}
                  </code>
                  <button
                    onClick={() => navigator.clipboard.writeText(url)}
                    className="px-3 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded border border-slate-600 text-xs transition"
                    title="Copy URL"
                  >
                    📋
                  </button>
                </div>
              </div>
              
              {/* Send Request Button */}
              <button
                onClick={handleFetch}
                disabled={loading}
                className="w-full mt-3 bg-primary-500 hover:bg-primary-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-4 py-2.5 rounded font-semibold transition text-sm"
              >
                {loading ? 'Loading...' : '▶ Send Request'}
              </button>
            </div>
          </div>
        
        {/* RIGHT COLUMN: Response (50%) */}
        <div>
          <div className="h-[calc(100vh-10rem)] overflow-y-auto">
          {(response || error) ? (
            <div ref={responseRef} className="bg-slate-800/50 backdrop-blur p-4 rounded-lg border border-slate-700">
              <h2 className="text-base font-bold text-white mb-3 flex items-center gap-2">
                <span>📊</span> Response
              </h2>
              
              {error ? (
                <div className="bg-red-900/20 border border-red-500/50 p-4 rounded">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="text-red-400 font-semibold">❌ Error</span>
                  </div>
                  <p className="text-red-300 text-sm">{error}</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {/* Response Headers */}
                  <div className="flex items-center justify-between flex-wrap gap-2 p-3 bg-slate-900/50 rounded border border-slate-700">
                    <div className="flex items-center gap-3 flex-wrap text-xs">
                      <span className="text-green-400 font-semibold">✓ {response.status}</span>
                      {response.cacheStatus !== 'N/A' && (
                        <span className="text-slate-400">
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
                        <span className="text-slate-400">
                          Time: <span className="text-white">{requestTime}ms</span>
                        </span>
                      )}
                    </div>
                  </div>
                  
                  {/* Applied Middleware Indicators */}
                  {endpointType === 'collection' && (delay > 0 || flakyRate > 0 || skipCache || selectedFields.length > 0 || sortField || searchQuery) && (
                    <div className="p-3 bg-blue-900/20 border border-blue-500/30 rounded">
                      <p className="text-xs text-blue-300 font-semibold mb-2">Applied Middleware:</p>
                      <div className="flex flex-wrap gap-1.5">
                        {delay > 0 && (
                          <span className="text-xs bg-blue-900/50 text-blue-200 px-2 py-0.5 rounded">
                            ⏱️ {delay}ms
                          </span>
                        )}
                        {flakyRate > 0 && (
                          <span className="text-xs bg-yellow-900/50 text-yellow-200 px-2 py-0.5 rounded">
                            🎲 {flakyRate}%
                          </span>
                        )}
                        {skipCache && (
                          <span className="text-xs bg-purple-900/50 text-purple-200 px-2 py-0.5 rounded">
                            🚫 Cache
                          </span>
                        )}
                        {selectedFields.length > 0 && (
                          <span className="text-xs bg-green-900/50 text-green-200 px-2 py-0.5 rounded">
                            🔍 {selectedFields.length} fields
                          </span>
                        )}
                        {sortField && (
                          <span className="text-xs bg-indigo-900/50 text-indigo-200 px-2 py-0.5 rounded">
                            🔀 {sortField}
                          </span>
                        )}
                        {searchQuery && (
                          <span className="text-xs bg-pink-900/50 text-pink-200 px-2 py-0.5 rounded">
                            🔎 "{searchQuery}"
                          </span>
                        )}
                      </div>
                    </div>
                  )}
                  
                  {/* Pagination Metadata */}
                  {response.data?.pagination && (
                    <div className="p-3 bg-slate-900/50 border border-slate-600 rounded">
                      <p className="text-xs text-slate-400 font-semibold mb-2">Pagination:</p>
                      <div className="grid grid-cols-4 gap-2 text-xs">
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
                    </div>
                  )}
                  
                  {/* Response Body */}
                  <div className="bg-slate-900 p-4 rounded border border-slate-600 overflow-auto" style={{ maxHeight: 'calc(100vh - 362px)' }}>
                    <pre className="text-green-400 text-xs">
                      <code>{JSON.stringify(response.data, null, 2)}</code>
                    </pre>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="bg-slate-800/50 backdrop-blur p-8 rounded-lg border border-slate-700 text-center h-full flex items-center justify-center">
              <div>
                <div className="text-6xl mb-4">📡</div>
                <p className="text-slate-400 text-sm">
                  Configure your request and click <strong className="text-white">Send Request</strong> to see the response here.
                </p>
              </div>
            </div>
          )}
          </div>
        </div>
      </div>
      
      {/* Endpoints Section */}
      <section className="bg-slate-800/50 backdrop-blur p-5 rounded-lg border border-slate-700">
        <h2 className="text-xl font-bold text-white mb-3">Endpoints</h2>
        
        <div className="space-y-3">
          {/* Collection Endpoint */}
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="bg-green-500 text-white px-2 py-1 rounded text-xs font-semibold">GET</span>
              <code className="text-primary-500 text-sm">{directPath}</code>
            </div>
            <p className="text-slate-400 text-xs mb-1">Get a collection of {resourceName}</p>
            <div className="text-xs text-slate-500">
              Alternative: <code className="text-slate-500">{groupPath}</code>
            </div>
          </div>
          
          {/* Single Item Endpoint */}
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="bg-green-500 text-white px-2 py-1 rounded text-xs font-semibold">GET</span>
              <code className="text-primary-500 text-sm">{directPath}/:id</code>
            </div>
            <p className="text-slate-400 text-xs mb-1">Get a single {resource} by ID</p>
            <div className="text-xs text-slate-500">
              Alternative: <code className="text-slate-500">{groupPath}/:id</code>
            </div>
          </div>
          
          {/* Meta Endpoint */}
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="bg-green-500 text-white px-2 py-1 rounded text-xs font-semibold">GET</span>
              <code className="text-primary-500 text-sm">{directPath}/meta</code>
            </div>
            <p className="text-slate-400 text-xs mb-1">Get resource metadata and schema</p>
            <div className="text-xs text-slate-500">
              Alternative: <code className="text-slate-500">{groupPath}/meta</code>
            </div>
          </div>
          
          {/* Aliases */}
          {aliases.length > 0 && (
            <div className="mt-3 pt-3 border-t border-slate-700">
              <p className="text-xs text-slate-400 mb-2 font-semibold">Additional aliases:</p>
              <div className="flex flex-wrap gap-1.5">
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
      
      {/* Search Guide */}
      <section className="bg-gradient-to-br from-blue-900/20 to-purple-900/20 backdrop-blur p-5 rounded-lg border border-blue-700/30">
        <div className="flex items-center gap-2 mb-3">
          <span className="text-xl">🔍</span>
          <h2 className="text-xl font-bold text-white">How to Use Search</h2>
        </div>
        
        <div className="space-y-4">
          {/* Overview */}
          <div>
            <p className="text-slate-300 text-sm leading-relaxed">
              The search feature allows you to filter results by searching for text across all or specific fields. 
              Choose between <code className="px-1.5 py-0.5 bg-slate-800 text-blue-300 rounded text-xs">?q=</code> or{' '}
              <code className="px-1.5 py-0.5 bg-slate-800 text-blue-300 rounded text-xs">?search=</code> parameter names.
            </p>
          </div>

          {/* Basic Examples */}
          <div>
            <h3 className="text-base font-semibold text-white mb-2 flex items-center gap-2">
              <span className="text-blue-400">1.</span> Basic Search
            </h3>
            <div className="space-y-2">
              <div className="bg-slate-900/50 p-3 rounded border border-slate-700">
                <p className="text-xs text-slate-400 mb-2">Search across all text fields:</p>
                <CodeExample 
                  title="Basic Search"
                  code={`GET ${API_URL}${directPath}?q=laptop`}
                  language="bash"
                />
              </div>
            </div>
          </div>

          {/* Field-Specific Search */}
          <div>
            <h3 className="text-base font-semibold text-white mb-2 flex items-center gap-2">
              <span className="text-blue-400">2.</span> Search Specific Fields
            </h3>
            <div className="bg-slate-900/50 p-3 rounded border border-slate-700">
              <CodeExample 
                title="Field-Specific Search"
                code={`GET ${API_URL}${directPath}?q=laptop&search_fields=name,description`}
                language="bash"
              />
            </div>
          </div>

          {/* Tips */}
          <div className="bg-blue-500/10 border border-blue-500/30 rounded p-3">
            <h4 className="text-xs font-semibold text-blue-300 mb-2 flex items-center gap-2">
              <span>💡</span> Pro Tips
            </h4>
            <ul className="space-y-1 text-xs text-slate-300">
              <li className="flex items-start gap-2">
                <span className="text-blue-400">•</span>
                <span>Search is <strong>case-insensitive</strong> and performs partial matching</span>
              </li>
              <li className="flex items-start gap-2">
                <span className="text-blue-400">•</span>
                <span>Without <code className="px-1 py-0.5 bg-slate-800 text-blue-300 rounded">search_fields</code>, all text fields are searched</span>
              </li>
              <li className="flex items-start gap-2">
                <span className="text-blue-400">•</span>
                <span>Combine with pagination to handle large result sets efficiently</span>
              </li>
            </ul>
          </div>
        </div>
      </section>
      
      {/* Code Examples */}
      <section>
        <h2 className="text-xl font-bold text-white mb-3">Code Examples</h2>
        
        <div className="grid md:grid-cols-2 gap-3">
          {/* JavaScript */}
          <div>
            <h3 className="text-base font-semibold text-white mb-2">JavaScript</h3>
            <CodeExample 
              title="Fetch Collection"
              code={`// Get 10 ${resourceName}
fetch('${API_URL}${directPath}?count=10')
  .then(res => res.json())
  .then(data => console.log(data));`}
              language="javascript"
            />
          </div>
          
          {/* Python */}
          <div>
            <h3 className="text-base font-semibold text-white mb-2">Python</h3>
            <CodeExample 
              title="Fetch with Requests"
              code={`import requests

response = requests.get(
    '${API_URL}${directPath}?count=10'
)
data = response.json()`}
              language="python"
            />
          </div>
          
          {/* cURL */}
          <div>
            <h3 className="text-base font-semibold text-white mb-2">cURL</h3>
            <CodeExample 
              title="Command Line"
              code={`curl "${API_URL}${directPath}?count=10"`}
              language="bash"
            />
          </div>
          
          {/* Fresh Data */}
          <div>
            <h3 className="text-base font-semibold text-white mb-2">Bypass Cache</h3>
            <CodeExample 
              title="Get Fresh Data"
              code={`fetch('${API_URL}${directPath}?nocache=true')`}
              language="javascript"
            />
          </div>
        </div>
      </section>
      
      {/* Schema Properties */}
      <section className="bg-slate-800/50 backdrop-blur p-5 rounded-lg border border-slate-700">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-xl font-bold text-white">Schema Properties</h2>
        </div>
        
        {Object.keys(properties).length > 0 ? (
          <div className="space-y-3">
            {Object.entries(properties).map(([key, value]: [string, any]) => (
              <div key={key} className="bg-slate-900/50 p-3 rounded border border-slate-700">
                <div className="flex items-start justify-between mb-2">
                  <div className="flex items-baseline gap-2 flex-wrap">
                    <code className="text-primary-500 font-semibold text-sm">{key}</code>
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
                  <p className="text-xs text-slate-400 mb-2">{value.description}</p>
                )}
                
                <div className="flex flex-wrap gap-1.5 text-xs">
                  {(value['x-generator'] || value['x-faker']) && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-0.5 rounded flex items-center gap-1">
                      <span className="text-slate-500">Gen:</span>
                      <code className="text-primary-400">{value['x-generator'] || value['x-faker']}</code>
                    </span>
                  )}
                  {value.format && (
                    <span className="bg-slate-800 text-slate-400 px-2 py-0.5 rounded">
                      {value.format}
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="bg-slate-900/50 p-6 rounded border border-slate-700 text-center">
            <p className="text-slate-400 text-sm mb-3">No properties defined in schema</p>
            <button
              onClick={() => setShowRawSchema(true)}
              className="inline-block bg-primary-500 hover:bg-primary-600 text-white px-3 py-2 rounded text-sm transition"
            >
              View Raw Schema
            </button>
          </div>
        )}
        
        {/* Raw Schema Toggle */}
        <div className="mt-3 pt-3 border-t border-slate-700">
          <div className="flex items-center justify-between mb-2">
            <button
              onClick={() => setShowRawSchema(!showRawSchema)}
              className="flex items-center gap-2 text-xs text-slate-400 hover:text-white transition"
            >
              <svg 
                className={`w-3 h-3 transition-transform ${showRawSchema ? 'rotate-90' : ''}`}
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
                Copy
              </button>
            )}
          </div>
          
          {showRawSchema && (
            <div className="mt-2 bg-slate-900 p-3 rounded border border-slate-600 overflow-auto max-h-96">
              <pre className="text-green-400 text-xs">
                <code>{JSON.stringify(schema, null, 2)}</code>
              </pre>
            </div>
          )}
        </div>
      </section>
      
      {/* Query Parameters */}
      <section className="bg-slate-800/50 backdrop-blur p-5 rounded-lg border border-slate-700">
        <h2 className="text-xl font-bold text-white mb-3">Query Parameters</h2>
        
        <div className="space-y-3">
          <div>
            <code className="text-primary-500 font-semibold text-sm">count</code>
            <span className="text-xs text-slate-500 ml-2">integer</span>
            <p className="text-xs text-slate-400 mt-1">
              Number of items to return (default: 10, max: 100)
            </p>
          </div>
          
          <div>
            <code className="text-primary-500 font-semibold text-sm">seed</code>
            <span className="text-xs text-slate-500 ml-2">integer</span>
            <p className="text-xs text-slate-400 mt-1">
              Seed for reproducible data generation
            </p>
          </div>
          
          <div>
            <code className="text-primary-500 font-semibold text-sm">nocache</code>
            <span className="text-xs text-slate-500 ml-2">boolean</span>
            <p className="text-xs text-slate-400 mt-1">
              Bypass cache and generate fresh data on every request
            </p>
          </div>
        </div>
      </section>
      
    </div>
  )
}

