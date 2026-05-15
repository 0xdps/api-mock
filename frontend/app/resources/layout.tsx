import { Header } from '@/components/Header'
import { DocsSidebarWrapper } from '@/components/DocsSidebarWrapper'
import { schemasManifest } from '@/lib/schemas-manifest'

async function getSchemas() {
  return schemasManifest
}

export default async function DocsLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const schemas = await getSchemas()
  
  // Group schemas by category
  const groups = schemas.reduce((acc, item) => {
    const group = item.group || 'other'
    if (!acc[group]) acc[group] = []
    acc[group].push(item.name)
    return acc
  }, {} as Record<string, string[]>)
  
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col">
      <Header compact />
      
      <div className="flex h-[calc(100vh-4rem)] overflow-hidden">
        {/* Sidebar - Fixed width on desktop */}
        <div className="w-80 flex-shrink-0 hidden lg:block">
          <DocsSidebarWrapper groups={groups} />
        </div>
        
        {/* Mobile Sidebar */}
        <div className="lg:hidden">
          <DocsSidebarWrapper groups={groups} />
        </div>
        
        {/* Main Content - Scrollable */}
        <div className="flex-1 overflow-y-auto">
          {children}
        </div>
      </div>
    </div>
  )
}

