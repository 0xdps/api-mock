'use client'

import { useState } from 'react'
import { EchoTester } from '@/components/playground/EchoTester'
import { StatusCodeGenerator } from '@/components/playground/StatusCodeGenerator'
import { DelayTester } from '@/components/playground/DelayTester'
import { MiddlewareTester } from '@/components/playground/MiddlewareTester'
import { ChaosTester } from '@/components/playground/ChaosTester'
import { ResourceSelector } from '@/components/ResourceSelector'

type TabType = 'echo' | 'status' | 'delay' | 'middleware' | 'chaos' | 'resources'

const tabs: { id: TabType; label: string; icon: string }[] = [
  { id: 'echo', label: 'Echo', icon: '🔊' },
  { id: 'status', label: 'Status', icon: '📊' },
  { id: 'delay', label: 'Delay', icon: '⏱️' },
  { id: 'middleware', label: 'Middleware', icon: '⚙️' },
  { id: 'chaos', label: 'Chaos', icon: '🎲' },
  { id: 'resources', label: 'Resources', icon: '📦' },
]

export default function PlaygroundPage() {
  const [activeTab, setActiveTab] = useState<TabType>('echo')

  const renderTabContent = () => {
    switch (activeTab) {
      case 'echo':
        return <EchoTester />
      case 'status':
        return <StatusCodeGenerator />
      case 'delay':
        return <DelayTester />
      case 'middleware':
        return <MiddlewareTester />
      case 'chaos':
        return <ChaosTester />
      case 'resources':
        return <ResourceSelector />
      default:
        return null
    }
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-white mb-2">
          API Testing Playground
        </h1>
        <p className="text-slate-400 text-lg">
          Interactive testing interface for all API features
        </p>
      </div>

      {/* Tab Navigation */}
      <div className="mb-8 border-b border-slate-700">
        <div className="flex space-x-1 overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                flex items-center gap-2 px-6 py-3 font-medium transition-all whitespace-nowrap
                border-b-2 
                ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-white bg-slate-800/50'
                    : 'border-transparent text-slate-400 hover:text-slate-300 hover:bg-slate-800/30'
                }
              `}
            >
              <span className="text-xl">{tab.icon}</span>
              <span>{tab.label}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Tab Content */}
      <div className="bg-slate-800/30 rounded-lg border border-slate-700 p-6">
        {renderTabContent()}
      </div>
    </div>
  )
}
