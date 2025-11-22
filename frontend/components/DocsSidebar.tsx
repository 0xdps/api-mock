'use client'

import { useState, useRef, useEffect } from 'react'

interface DocsSidebarProps {
  groups: Record<string, string[]>
  selectedResource: string
  onResourceSelect: (resource: string) => void
}

export function DocsSidebar({ groups, selectedResource, onResourceSelect }: DocsSidebarProps) {
  const [search, setSearch] = useState('')
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(new Set(Object.keys(groups)))
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  
  // Filter groups and resources based on search
  const filteredGroups = Object.entries(groups).reduce((acc, [groupName, resources]) => {
    const filteredResources = resources.filter(resource =>
      resource.toLowerCase().includes(search.toLowerCase())
    )
    
    if (filteredResources.length > 0 || groupName.toLowerCase().includes(search.toLowerCase())) {
      acc[groupName] = search 
        ? resources.filter(r => r.toLowerCase().includes(search.toLowerCase()))
        : resources
    }
    
    return acc
  }, {} as Record<string, string[]>)
  
  // Auto-expand groups when searching
  useEffect(() => {
    if (search) {
      setExpandedGroups(new Set(Object.keys(filteredGroups)))
    }
  }, [search, filteredGroups])
  
  const toggleGroup = (groupName: string) => {
    const newExpanded = new Set(expandedGroups)
    if (newExpanded.has(groupName)) {
      newExpanded.delete(groupName)
    } else {
      newExpanded.add(groupName)
    }
    setExpandedGroups(newExpanded)
  }
  
  const handleResourceClick = (resource: string) => {
    onResourceSelect(resource)
    setIsSidebarOpen(false) // Close on mobile
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
  
  const sidebarContent = (
    <>
      {/* Search */}
      <div className="p-4 border-b border-slate-700">
        <div className="relative">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search resources..."
            className="w-full bg-slate-700 text-white px-3 py-2 pr-8 rounded border border-slate-600 focus:border-primary-500 focus:outline-none text-sm"
          />
          {search && (
            <button
              onClick={() => setSearch('')}
              className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          )}
        </div>
      </div>
      
      {/* Resource List */}
      <div className="overflow-y-auto flex-1">
        {Object.keys(filteredGroups).length > 0 ? (
          <div className="p-2">
            {Object.entries(filteredGroups).sort(([a], [b]) => a.localeCompare(b)).map(([groupName, resources]) => (
              <div key={groupName} className="mb-2">
                {/* Group Header */}
                <button
                  onClick={() => toggleGroup(groupName)}
                  className="w-full flex items-center justify-between px-3 py-2 rounded hover:bg-slate-700 transition text-left"
                >
                  <div className="flex items-center gap-2">
                    <span className="text-lg">{getGroupIcon(groupName)}</span>
                    <span className="text-white font-medium capitalize text-sm">
                      {groupName}
                    </span>
                    <span className="text-xs text-slate-500">
                      ({resources.length})
                    </span>
                  </div>
                  <svg
                    className={`w-4 h-4 text-slate-400 transition-transform ${
                      expandedGroups.has(groupName) ? 'rotate-90' : ''
                    }`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </button>
                
                {/* Resources in Group */}
                {expandedGroups.has(groupName) && (
                  <div className="ml-4 mt-1 space-y-1">
                    {resources.map((resource) => (
                      <button
                        key={resource}
                        onClick={() => handleResourceClick(resource)}
                        className={`w-full text-left px-3 py-1.5 rounded text-sm transition ${
                          selectedResource === resource
                            ? 'bg-primary-500 text-white'
                            : 'text-slate-300 hover:bg-slate-700 hover:text-white'
                        }`}
                      >
                        {resource}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="p-4 text-center text-slate-400 text-sm">
            No resources found
          </div>
        )}
      </div>
    </>
  )
  
  return (
    <>
      {/* Mobile Toggle Button */}
      <button
        onClick={() => setIsSidebarOpen(!isSidebarOpen)}
        className="lg:hidden fixed bottom-4 right-4 z-50 bg-primary-500 text-white p-3 rounded-full shadow-lg"
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      
      {/* Mobile Overlay */}
      {isSidebarOpen && (
        <div
          className="lg:hidden fixed inset-0 bg-black/50 z-40"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}
      
      {/* Sidebar - Desktop */}
      <aside className="hidden lg:flex lg:flex-col h-full bg-slate-800/50 backdrop-blur border-r border-slate-700">
        {sidebarContent}
      </aside>
      
      {/* Sidebar - Mobile */}
      <aside
        className={`lg:hidden fixed top-0 left-0 bottom-0 w-80 bg-slate-800 border-r border-slate-700 z-40 transform transition-transform ${
          isSidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } flex flex-col`}
      >
        {/* Close Button */}
        <div className="p-4 border-b border-slate-700 flex items-center justify-between">
          <h2 className="text-white font-bold">Resources</h2>
          <button
            onClick={() => setIsSidebarOpen(false)}
            className="text-slate-400 hover:text-white"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        {sidebarContent}
      </aside>
    </>
  )
}

