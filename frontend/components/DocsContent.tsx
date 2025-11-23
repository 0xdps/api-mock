'use client'

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { DocsSidebar } from './DocsSidebar'
import { ResourceDocumentation } from './ResourceDocumentation'

interface DocsContentProps {
  schemas: Array<{ name: string; schema: any; group: string }>
  groups: Record<string, string[]>
  initialResource: string
}

export function DocsContent({ schemas, groups, initialResource }: DocsContentProps) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [selectedResource, setSelectedResource] = useState(initialResource)
  
  // Sync with URL on mount (for direct loads)
  useEffect(() => {
    const resourceFromUrl = searchParams.get('resource')
    if (resourceFromUrl && schemas.some(s => s.name === resourceFromUrl)) {
      setSelectedResource(resourceFromUrl)
    }
  }, [searchParams, schemas])
  
  // Update URL when resource changes (client-side navigation)
  const handleResourceSelect = (resource: string) => {
    setSelectedResource(resource)
    // Update URL without page reload (shallow routing)
    const params = new URLSearchParams(searchParams.toString())
    params.set('resource', resource)
    router.push(`/docs?${params.toString()}`, { scroll: false })
  }
  
  // Find the selected schema
  const selectedSchema = schemas.find(s => s.name === selectedResource)
  
  return (
    <div className="flex h-[calc(100vh-4rem)] overflow-hidden">
      {/* Sidebar - Fixed width on desktop */}
      <div className="w-80 flex-shrink-0 hidden lg:block">
        <DocsSidebar 
          groups={groups}
          selectedResource={selectedResource}
          onResourceSelect={handleResourceSelect}
        />
      </div>
      
      {/* Mobile Sidebar */}
      <div className="lg:hidden">
        <DocsSidebar 
          groups={groups}
          selectedResource={selectedResource}
          onResourceSelect={handleResourceSelect}
        />
      </div>
      
      {/* Main Content - Scrollable */}
      <div className="flex-1 overflow-y-auto">
        <div className="max-w-5xl mx-auto p-8">
          {selectedSchema ? (
            <ResourceDocumentation 
              resource={selectedSchema.name}
              schema={selectedSchema.schema}
              group={selectedSchema.group}
            />
          ) : (
            <div className="text-center text-slate-400 mt-20">
              <p className="text-xl">Select a resource to view documentation</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

