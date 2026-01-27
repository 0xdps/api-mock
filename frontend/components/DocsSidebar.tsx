'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Search, X, ChevronRight, Menu } from 'lucide-react'
import { getGroupIcon } from '@/lib/icons'

interface DocsSidebarProps {
  groups: Record<string, string[]>
  selectedResource: string
}

export function DocsSidebar({ groups, selectedResource }: DocsSidebarProps) {
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
  
  // Auto-expand group containing selected resource
  useEffect(() => {
    if (selectedResource && !search) {
      const groupContainingResource = Object.entries(groups).find(([_, resources]) => 
        resources.includes(selectedResource)
      )?.[0]
      
      if (groupContainingResource) {
        setExpandedGroups(prev => {
          const next = new Set(prev)
          next.add(groupContainingResource)
          return next
        })
      }
    }
  }, [selectedResource, groups, search])
  
  const toggleGroup = (groupName: string) => {
    const newExpanded = new Set(expandedGroups)
    if (newExpanded.has(groupName)) {
      newExpanded.delete(groupName)
    } else {
      newExpanded.add(groupName)
    }
    setExpandedGroups(newExpanded)
  }
  
  const handleResourceClick = () => {
    setIsSidebarOpen(false) // Close on mobile
  }
  
  const sidebarContent = (
    <>
      <div className="p-4 border-b border-white/5">
        <div className="relative">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search resources..."
            className="w-full bg-slate-900 border border-white/10 rounded-lg px-9 py-2 text-sm text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-primary-500 transition-all font-medium"
          />
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
          {search && (
            <button
              onClick={() => setSearch('')}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
            >
              <X className="w-4 h-4" />
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
                  className="w-full flex items-center justify-between px-3 py-2 rounded-lg hover:bg-white/5 transition-colors text-left group"
                >
                  <div className="flex items-center gap-2.5">
                    {(() => {
                      const Icon = getGroupIcon(groupName)
                      return <Icon className="w-4 h-4 text-slate-400 group-hover:text-primary-400 transition-colors" />
                    })()}
                    <span className="text-white font-semibold capitalize text-sm tracking-tight">
                      {groupName}
                    </span>
                    <span className="text-[10px] text-slate-500 font-bold bg-slate-900 px-1.5 py-0.5 rounded border border-white/5">
                      {resources.length}
                    </span>
                  </div>
                  <ChevronRight
                    className={`w-4 h-4 text-slate-500 transition-transform duration-200 ${
                      expandedGroups.has(groupName) ? 'rotate-90 text-primary-500' : ''
                    }`}
                  />
                </button>
                
                {/* Resources in Group */}
                {expandedGroups.has(groupName) && (
                  <div className="ml-4 mt-1 space-y-1">
                    {resources.map((resource) => (
                      <Link
                        key={resource}
                        href={`/docs/${resource}`}
                        onClick={handleResourceClick}
                        scroll={false}
                        className={`block w-full text-left px-3 py-1.5 rounded text-sm transition ${
                          selectedResource === resource
                            ? 'bg-primary-500 text-white'
                            : 'text-slate-300 hover:bg-slate-700 hover:text-white'
                        }`}
                      >
                        {resource}
                      </Link>
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
        className="lg:hidden fixed bottom-6 right-6 z-50 bg-primary-500 hover:bg-primary-600 text-white p-4 rounded-full shadow-2xl shadow-primary-500/40 transition-all active:scale-95"
      >
        <Menu className="w-6 h-6" />
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
        <div className="p-4 border-b border-white/5 flex items-center justify-between">
          <h2 className="text-white font-bold tracking-tight">Resources</h2>
          <button
            onClick={() => setIsSidebarOpen(false)}
            className="text-slate-400 hover:text-white transition-colors"
          >
            <X className="w-6 h-6" />
          </button>
        </div>
        {sidebarContent}
      </aside>
    </>
  )
}

