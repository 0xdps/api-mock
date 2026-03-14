'use client'

import { useState, useRef, useEffect } from 'react'
import { 
  Terminal, 
  Database, 
  Play, 
  RefreshCw, 
  ChevronDown, 
  Search,
  CheckCircle2,
  Settings,
  ShieldCheck,
  Zap,
  Globe,
  Monitor,
  ClipboardList
} from 'lucide-react'
import { getApiUrl, apiClient } from '@/lib/api'
import { trackUmamiEvent } from '@/lib/analytics'

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
    trackUmamiEvent('pg_req_s', {
      resource: selectedResource,
      group: resourceGroup || null,
      endpointType,
      count: endpointType === 'collection' ? count : null,
      hasNoCache: noCache,
    })

    setLoading(true)
    setError(null)
    setRequestTime(null)
    
    const startTime = performance.now()
    
    try {
      // Always use direct URL for the actual API call
      const res = await apiClient.get(directUrl)
      const endTime = performance.now()
      const duration = Math.round(endTime - startTime)
      
      if (res.isError()) {
        throw new Error(`HTTP ${res.status}`)
      }
      
      const data = res.json()
      setResponse(data)
      setRequestTime(duration)

      trackUmamiEvent('pg_req_ok', {
        resource: selectedResource,
        endpointType,
        duration,
      })
    } catch (err: any) {
      const endTime = performance.now()
      const duration = Math.round(endTime - startTime)
      
      setError(err.message)
      setRequestTime(duration)

      trackUmamiEvent('pg_req_er', {
        resource: selectedResource,
        endpointType,
        duration,
        message: err.message,
      })
    } finally {
      setLoading(false)
    }
  }

  const handleCopyResponse = () => {
    if (!response) return

    navigator.clipboard.writeText(JSON.stringify(response, null, 2))
    trackUmamiEvent('pg_rsp_cp', {
      resource: selectedResource,
      endpointType,
    })
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
        <div className="bg-slate-800/40 backdrop-blur-md p-6 rounded-xl border border-white/5 shadow-premium">
          <h2 className="text-sm font-bold text-white mb-6 flex items-center gap-2 uppercase tracking-wider text-slate-300">
            <Terminal className="w-4 h-4 text-primary-400" /> Request Builder
          </h2>
          
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
                  className="w-full bg-slate-900 border border-white/10 rounded-lg text-white px-3 py-2 text-sm pr-8 focus:outline-none focus:ring-1 focus:ring-primary-500 transition-all"
                />
                <button
                  type="button"
                  onClick={() => setIsGroupDropdownOpen(!isGroupDropdownOpen)}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                >
                  <ChevronDown className={`w-4 h-4 transition-transform ${isGroupDropdownOpen ? 'rotate-180' : ''}`} />
                </button>
              </div>
              
              {/* Dropdown List */}
              {isGroupDropdownOpen && (
                <div className="absolute z-10 w-full mt-1 bg-slate-800 border border-white/10 rounded-lg shadow-2xl max-h-60 overflow-auto backdrop-blur-xl">
                  {/* All Resources Option */}
                  <button
                    type="button"
                    onClick={() => handleGroupSelect('')}
                    className={`w-full text-left px-3 py-2 text-sm hover:bg-white/5 transition ${
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
                        className={`w-full text-left px-3 py-2 text-sm hover:bg-white/5 transition ${
                          group === selectedGroup ? 'bg-primary-500/20 text-primary-400' : 'text-white'
                        }`}
                      >
                        {group} <span className="text-slate-500 text-xs ml-1">({groups[group].length})</span>
                      </button>
                    ))
                  ) : (
                    <div className="px-3 py-2 text-sm text-slate-500">
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
                  <span className="ml-1 text-[10px] text-slate-500 font-bold uppercase tracking-tight">
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
                  className="w-full bg-slate-900 border border-white/10 rounded-lg text-white px-3 py-2 text-sm pr-8 focus:outline-none focus:ring-1 focus:ring-primary-500 transition-all"
                />
                <button
                  type="button"
                  onClick={() => setIsResourceDropdownOpen(!isResourceDropdownOpen)}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                >
                  <ChevronDown className={`w-4 h-4 transition-transform ${isResourceDropdownOpen ? 'rotate-180' : ''}`} />
                </button>
              </div>
              
              {/* Dropdown List */}
              {isResourceDropdownOpen && (
                <div className="absolute z-10 w-full mt-1 bg-slate-800 border border-white/10 rounded-lg shadow-2xl max-h-60 overflow-auto backdrop-blur-xl">
                  {filteredResources.length > 0 ? (
                    filteredResources.map(resource => (
                      <button
                        key={resource}
                        type="button"
                        onClick={() => handleResourceSelect(resource)}
                        className={`w-full text-left px-3 py-2 text-sm hover:bg-white/5 transition ${
                          resource === selectedResource ? 'bg-primary-500/20 text-primary-400' : 'text-white'
                        }`}
                      >
                        {resource}
                      </button>
                    ))
                  ) : (
                    <div className="px-3 py-2 text-sm text-slate-500">
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
                className={`flex-1 px-3 py-2 rounded-lg transition-all text-sm font-bold ${
                  endpointType === 'collection'
                    ? 'bg-primary-500 text-white shadow-lg shadow-primary-500/20'
                    : 'bg-white/5 text-slate-400 hover:bg-white/10'
                }`}
              >
                Collection
              </button>
              <button
                onClick={() => setEndpointType('single')}
                className={`flex-1 px-3 py-2 rounded-lg transition-all text-sm font-bold ${
                  endpointType === 'single'
                    ? 'bg-primary-500 text-white shadow-lg shadow-primary-500/20'
                    : 'bg-white/5 text-slate-400 hover:bg-white/10'
                }`}
              >
                Single
              </button>
              <button
                onClick={() => setEndpointType('meta')}
                className={`flex-1 px-3 py-2 rounded-lg transition-all text-sm font-bold ${
                  endpointType === 'meta'
                    ? 'bg-primary-500 text-white shadow-lg shadow-primary-500/20'
                    : 'bg-white/5 text-slate-400 hover:bg-white/10'
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
                    className="w-full bg-slate-900 border border-white/10 text-white px-3 py-2 rounded-lg focus:outline-none focus:ring-1 focus:ring-primary-500 text-sm transition-all"
                  />
                </div>
              )}
              
              {endpointType === 'collection' && (
                <div className="flex-1">
                  <label className="block text-slate-300 mb-2 font-medium text-sm">
                    Count <span className="text-slate-500 text-[10px] font-bold uppercase ml-1">(1-100)</span>
                  </label>
                  <input
                    type="number"
                    min="1"
                    max="100"
                    value={count}
                    onChange={(e) => setCount(parseInt(e.target.value) || 1)}
                    className="w-full bg-slate-900 border border-white/10 text-white px-3 py-2 rounded-lg focus:outline-none focus:ring-1 focus:ring-primary-500 text-sm transition-all"
                  />
                </div>
              )}
              
              {/* Cache Toggle - Inline */}
              {endpointType !== 'meta' && (
                <div className="flex-1">
                  <label className="block text-slate-300 mb-2 font-medium text-sm">
                    Cache
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer bg-white/5 px-3 py-2 rounded-lg border border-white/5 hover:bg-white/10 transition-all group">
                    <input
                      type="checkbox"
                      checked={noCache}
                      onChange={(e) => setNoCache(e.target.checked)}
                      className="w-4 h-4 rounded-md border-white/10 bg-slate-900 text-primary-500 focus:ring-primary-500/20"
                    />
                    <span className="text-slate-200 text-sm font-medium">
                      Bypass
                    </span>
                  </label>
                </div>
              )}
            </div>
          </div>
          
          {/* Send Button */}
          <button
            onClick={handleFetch}
            disabled={loading || !selectedResource}
            className="w-full bg-primary-500 hover:bg-primary-600 disabled:bg-slate-800 disabled:text-slate-500 disabled:cursor-not-allowed text-white px-6 py-3 rounded-xl font-semibold text-base transition-all mb-6 shadow-lg shadow-primary-500/20 flex items-center justify-center gap-3 group active:scale-[0.98]"
          >
            {loading ? <RefreshCw className="w-5 h-5 animate-spin" /> : <Play className="w-5 h-5 fill-current" />}
            {loading ? 'Sending...' : 'Send Request'}
          </button>
          
          {/* URL Preview */}
          <div className="border-t border-white/5 pt-6">
            <label className="text-sm font-bold text-white mb-4 flex items-center gap-2 uppercase tracking-wider text-slate-300">
              <ClipboardList className="w-4 h-4 text-primary-400" /> Request URL
            </label>
            
            {/* Direct URL */}
            <div className="mb-3">
              <div className="flex items-center justify-between mb-1.5">
                <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Direct Access</span>
                <span className="text-[10px] font-bold text-emerald-500 uppercase tracking-wider bg-emerald-500/5 px-1.5 py-0.5 rounded border border-emerald-500/10">Active</span>
              </div>
              <div className="bg-slate-900/50 p-4 rounded-xl border border-emerald-500/20 overflow-x-auto shadow-inner">
                <code className="text-emerald-400 text-xs font-mono break-all">{directUrl}</code>
              </div>
            </div>
            
            {/* Group URL */}
            {groupUrl && (
              <div>
                <div className="flex items-center justify-between mb-1.5">
                  <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Via Group Path</span>
                  <span className="text-[10px] font-bold text-slate-600 uppercase tracking-wider">Alternative</span>
                </div>
                <div className="bg-slate-900/30 p-4 rounded-xl border border-white/5 overflow-x-auto">
                  <code className="text-slate-400 text-xs font-mono break-all">{groupUrl}</code>
                </div>
              </div>
            )}
          </div>
        </div>
        
        {/* Code Examples */}
        <div className="bg-slate-800/40 backdrop-blur-md p-6 rounded-xl border border-white/5 shadow-premium">
          <h3 className="text-sm font-bold text-white mb-4 uppercase tracking-wider text-slate-300">Code Examples</h3>
          
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
        <div className="bg-slate-800/40 backdrop-blur-md p-6 rounded-xl border border-white/5 shadow-premium">
          <h3 className="text-sm font-bold text-white mb-4 uppercase tracking-wider text-slate-300">Supported Features</h3>
          
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs font-medium">
            <div className="flex items-start gap-3 p-3 bg-white/5 rounded-lg border border-white/5">
              <Settings className="w-4 h-4 text-primary-400 mt-0.5" />
              <div>
                <strong className="text-slate-200 block mb-0.5">Query Parameters</strong>
                <span className="text-slate-500">?count=N (max 100)</span>
              </div>
            </div>
            <div className="flex items-start gap-3 p-3 bg-white/5 rounded-lg border border-white/5">
              <RefreshCw className="w-4 h-4 text-emerald-400 mt-0.5" />
              <div>
                <strong className="text-slate-200 block mb-0.5">Cache Bypass</strong>
                <span className="text-slate-500">?nocache=true</span>
              </div>
            </div>
            <div className="flex items-start gap-3 p-3 bg-white/5 rounded-lg border border-white/5">
              <Globe className="w-4 h-4 text-blue-400 mt-0.5" />
              <div>
                <strong className="text-slate-200 block mb-0.5">Group Paths</strong>
                <span className="text-slate-500">/{resourceGroup}/{selectedResource}</span>
              </div>
            </div>
            <div className="flex items-start gap-3 p-3 bg-white/5 rounded-lg border border-white/5">
              <Zap className="w-4 h-4 text-amber-400 mt-0.5" />
              <div>
                <strong className="text-slate-200 block mb-0.5">CORS Enabled</strong>
                <span className="text-slate-500">Access from any origin</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      {/* Response */}
      <div className="space-y-6">
        <div className="bg-slate-800/40 backdrop-blur-md p-6 rounded-xl border border-white/5 shadow-premium min-h-[600px]">
          <h2 className="text-sm font-bold text-white mb-6 flex items-center gap-2 uppercase tracking-wider text-slate-300">
            <Database className="w-4 h-4 text-primary-400" /> Response
          </h2>
          
          {error && (
            <div className="space-y-4">
              {requestTime !== null && (
                <div className="text-slate-400 text-[10px] font-bold uppercase tracking-wider flex items-center gap-1.5 mb-2">
                  <Zap className="w-3.5 h-3.5 text-amber-500" />
                  Request took {requestTime}ms
                </div>
              )}
              <div className="bg-red-500/10 border border-red-500/20 p-5 rounded-xl">
                <div className="flex items-center gap-2 mb-2 text-red-500">
                  <ShieldCheck className="w-4 h-4" />
                  <strong className="text-sm">Error Detected</strong>
                </div>
                <p className="text-red-300/80 text-sm font-medium">{error}</p>
              </div>
            </div>
          )}
          
          {response && !error && (
            <div className="space-y-4">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                  <span className="text-[10px] font-bold bg-emerald-500/10 text-emerald-500 px-2 py-0.5 rounded border border-emerald-500/20 uppercase tracking-tight flex items-center gap-1">
                    <CheckCircle2 className="w-3 h-3" /> Status: 200 OK
                  </span>
                  {requestTime !== null && (
                    <span className="text-slate-500 text-[10px] font-bold uppercase tracking-tight flex items-center gap-1">
                      <Zap className="w-3 h-3 text-amber-500" /> {requestTime}ms
                    </span>
                  )}
                </div>
                <button
                  onClick={handleCopyResponse}
                  className="text-[10px] font-bold uppercase bg-white/5 hover:bg-white/10 text-slate-400 px-3 py-1.5 rounded-lg border border-white/5 transition-all active:scale-95"
                >
                  Copy JSON
                </button>
              </div>
              
              <div className="relative group">
                <pre className="bg-slate-900/50 p-5 rounded-2xl border border-white/5 overflow-auto max-h-[700px] shadow-inner font-mono text-xs leading-relaxed custom-scrollbar">
                  <code className="text-emerald-400/90 [text-shadow:0_0_10px_rgba(52,211,153,0.2)]">
                    {JSON.stringify(response, null, 2)}
                  </code>
                </pre>
              </div>
            </div>
          )}
          
          {!response && !error && !loading && (
            <div className="flex flex-col items-center justify-center py-32 group">
              <div className="w-16 h-16 bg-white/5 rounded-2xl flex items-center justify-center mb-4 group-hover:bg-primary-500/10 transition-colors border border-dashed border-white/10">
                <Monitor className="w-8 h-8 text-slate-600 group-hover:text-primary-400 transition-colors" />
              </div>
              <p className="text-slate-500 font-medium text-sm">Send a request to see the response</p>
            </div>
          )}
          
          {loading && (
            <div className="flex flex-col items-center justify-center py-32">
              <RefreshCw className="w-10 h-10 text-primary-500 animate-spin mb-4" />
              <div className="text-slate-400 font-bold text-sm uppercase tracking-widest animate-pulse">Processing...</div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function CodeSnippet({ title, code }: { title: string; code: string }) {
  const [copied, setCopied] = useState(false)
  
  const handleCopy = () => {
    navigator.clipboard.writeText(code)
    trackUmamiEvent('pg_code_cp', {
      snippet: title,
    })
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <div className="text-slate-500 text-[10px] font-bold uppercase tracking-wider">{title}</div>
        <button 
          onClick={handleCopy}
          className="text-[10px] font-bold text-slate-500 hover:text-white transition-colors"
        >
          {copied ? 'COPIED!' : 'COPY'}
        </button>
      </div>
      <div className="relative group">
        <pre className="bg-slate-900 border border-white/5 p-3 rounded-lg overflow-x-auto group-hover:border-primary-500/30 transition-colors">
          <code className="text-xs text-slate-300 font-mono">{code}</code>
        </pre>
      </div>
    </div>
  )
}
