'use client'

import { useState } from 'react'
import { DocsSidebar } from './DocsSidebar'
import { ResourceDocumentation } from './ResourceDocumentation'

interface DocsContentProps {
  schemas: Array<{ name: string; schema: any; group: string }>
  groups: Record<string, string[]>
}

export function DocsContent({ schemas, groups }: DocsContentProps) {
  const [selectedResource, setSelectedResource] = useState(schemas[0]?.name || '')
  
  // Find the selected schema
  const selectedSchema = schemas.find(s => s.name === selectedResource)
  
  return (
    <div className="flex h-[calc(100vh-4rem)] overflow-hidden">
      {/* Sidebar - Fixed width on desktop */}
      <div className="w-80 flex-shrink-0 hidden lg:block">
        <DocsSidebar 
          groups={groups}
          selectedResource={selectedResource}
          onResourceSelect={setSelectedResource}
        />
      </div>
      
      {/* Mobile Sidebar */}
      <div className="lg:hidden">
        <DocsSidebar 
          groups={groups}
          selectedResource={selectedResource}
          onResourceSelect={setSelectedResource}
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

