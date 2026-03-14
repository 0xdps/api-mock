'use client'

import { usePathname, useRouter } from 'next/navigation'
import Link from 'next/link'
import { 
  Volume2, 
  BarChart3, 
  Clock, 
  Layers, 
  Zap, 
  Package,
  Activity
} from 'lucide-react'

type TabType = 'echo' | 'status' | 'delay' | 'middleware' | 'chaos'

const tabs = [
  { id: 'echo', label: 'Echo', icon: Volume2, description: 'Test request/response', path: '/playground/echo' },
  { id: 'status', label: 'Status', icon: BarChart3, description: 'Generate HTTP status codes', path: '/playground/status' },
  { id: 'delay', label: 'Delay', icon: Clock, description: 'Test response delays', path: '/playground/delay' },
  { id: 'middleware', label: 'Middleware', icon: Layers, description: 'Test middleware parameters', path: '/playground/middleware' },
  { id: 'chaos', label: 'Chaos', icon: Zap, description: 'Chaos engineering tests', path: '/playground/chaos' },
] as const

export function PlaygroundSidebar() {
  const pathname = usePathname()
  const router = useRouter()

  const activeTab = tabs.find(tab => pathname === tab.path)

  return (
    <div className="h-full flex">
      {/* Sidebar Navigation */}
      <div className="w-64 bg-slate-900/50 border-r border-slate-700 flex flex-col">
        {/* Utility Tools */}
        <div className="flex-1 overflow-y-auto p-3">
          <div className="space-y-1">
            <div className="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Utilities
            </div>
            {tabs.map((tab) => (
              <Link
                key={tab.id}
                href={tab.path}
                data-umami-event="pg_tab"
                data-umami-label={`Playground Tab ${tab.label}`}
                className={`
                  w-full flex items-start gap-3 px-3 py-2.5 rounded-xl transition-all text-left group
                  ${
                    pathname === tab.path
                      ? 'bg-primary-500/10 text-primary-400 border border-primary-500/20 shadow-lg shadow-primary-500/5'
                      : 'text-slate-400 hover:text-slate-200 hover:bg-white/5 border border-transparent'
                  }
                `}
              >
                <tab.icon className={`w-5 h-5 flex-shrink-0 mt-0.5 transition-colors ${pathname === tab.path ? 'text-primary-400' : 'text-slate-500 group-hover:text-primary-400'}`} />
                <div className="flex-1 min-w-0">
                  <div className="font-bold text-sm tracking-tight">{tab.label}</div>
                  <div className="text-[10px] text-slate-500 font-medium uppercase tracking-wider mt-0.5">{tab.description}</div>
                </div>
              </Link>
            ))}
          </div>

          {/* Resources Link */}
          <div className="mt-6 space-y-1">
            <div className="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Resources
            </div>
            <Link
              href="/resources"
              data-umami-event="pg_rs"
              data-umami-label="Playground Browse Resources"
              className="w-full flex items-start gap-3 px-3 py-2.5 rounded-xl transition-all text-left text-slate-400 hover:text-slate-200 hover:bg-white/5 border border-transparent group"
            >
              <Package className="w-5 h-5 flex-shrink-0 mt-0.5 text-slate-500 group-hover:text-primary-400 transition-colors" />
              <div className="flex-1 min-w-0">
                <div className="font-bold text-sm tracking-tight">Browse Resources</div>
                <div className="text-[10px] text-slate-500 font-medium uppercase tracking-wider mt-0.5">Test 100+ API endpoints</div>
              </div>
            </Link>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Content Header */}
        {activeTab && (
          <div className="bg-slate-900/40 backdrop-blur-md border-b border-white/5 px-8 py-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-primary-500/10 rounded-xl border border-primary-500/20 shadow-lg shadow-primary-500/5">
                <activeTab.icon className="w-8 h-8 text-primary-400" />
              </div>
              <div>
                <h1 className="text-2xl font-bold text-white tracking-tight">{activeTab.label}</h1>
                <p className="text-sm text-slate-400 font-medium mt-0.5">{activeTab.description}</p>
              </div>
            </div>
          </div>
        )}

        {/* Tab Content - passed as children */}
        <div className="flex-1 min-h-0 p-6">
          {/* Children will be rendered here */}
        </div>
      </div>
    </div>
  )
}
