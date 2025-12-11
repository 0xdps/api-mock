'use client'

import { usePathname } from 'next/navigation'
import Link from 'next/link'
import { ReactNode } from 'react'

type TabType = 'echo' | 'status' | 'delay' | 'middleware' | 'chaos'

const tabs: { id: TabType; label: string; icon: string; description: string; path: string }[] = [
  { id: 'echo', label: 'Echo', icon: '🔊', description: 'Test request/response echoing', path: '/playground/echo' },
  { id: 'status', label: 'Status', icon: '📊', description: 'Generate HTTP status codes', path: '/playground/status' },
  { id: 'delay', label: 'Delay', icon: '⏱️', description: 'Test response delays', path: '/playground/delay' },
  { id: 'middleware', label: 'Middleware', icon: '⚙️', description: 'Test middleware parameters', path: '/playground/middleware' },
  { id: 'chaos', label: 'Chaos', icon: '🎲', description: 'Chaos engineering tests', path: '/playground/chaos' },
]

export default function PlaygroundTemplate({ children }: { children: ReactNode }) {
  const pathname = usePathname()
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
                className={`
                  w-full flex items-start gap-3 px-3 py-2.5 rounded-lg transition-all text-left
                  ${
                    pathname === tab.path
                      ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
                      : 'text-slate-400 hover:text-slate-300 hover:bg-slate-800/50'
                  }
                `}
              >
                <span className="text-xl flex-shrink-0">{tab.icon}</span>
                <div className="flex-1 min-w-0">
                  <div className="font-medium text-sm">{tab.label}</div>
                  <div className="text-xs text-slate-500 mt-0.5">{tab.description}</div>
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
              href="/docs"
              className="w-full flex items-start gap-3 px-3 py-2.5 rounded-lg transition-all text-left text-slate-400 hover:text-slate-300 hover:bg-slate-800/50"
            >
              <span className="text-xl flex-shrink-0">📦</span>
              <div className="flex-1 min-w-0">
                <div className="font-medium text-sm">Browse Resources</div>
                <div className="text-xs text-slate-500 mt-0.5">Test 100+ API endpoints</div>
              </div>
            </Link>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Content Header */}
        {activeTab && (
          <div className="bg-slate-900/30 border-b border-slate-700 px-6 py-4">
            <div className="flex items-center gap-3">
              <span className="text-2xl">{activeTab.icon}</span>
              <div>
                <h1 className="text-xl font-bold text-white">{activeTab.label}</h1>
                <p className="text-sm text-slate-400">{activeTab.description}</p>
              </div>
            </div>
          </div>
        )}

        {/* Tab Content */}
        <div className="flex-1 min-h-0 p-6">
          {children}
        </div>
      </div>
    </div>
  )
}
