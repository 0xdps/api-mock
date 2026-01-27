'use client'

import { useState, useMemo } from 'react'
import { useRouter } from 'next/navigation'
import { Search, X, ChevronDown } from 'lucide-react'
import { getGroupIcon } from '@/lib/icons'

interface ResourceSelectorProps {
  // Resources and groups will be fetched from API or passed as props
}

// Mock data - in production, this would come from API or server props
const mockResources = {
  people: ['user', 'profile', 'contact', 'employee', 'author'],
  business: ['company', 'invoice', 'contract', 'proposal', 'meeting', 'client', 'vendor', 'department', 'job', 'organization', 'report', 'subscription'],
  commerce: ['product', 'order', 'cart', 'payment', 'category', 'coupon', 'discount', 'inventory', 'promotion', 'refund', 'return', 'shipping', 'tag', 'wishlist'],
  content: ['article', 'blog', 'document', 'page', 'post'],
  social: ['comment', 'message', 'notification', 'feed'],
  media: ['image', 'video', 'audio', 'playlist'],
  travel: ['flight', 'hotel', 'booking', 'destination'],
  location: ['address', 'city', 'country', 'place'],
  finance: ['transaction', 'account', 'budget', 'expense'],
  food: ['recipe', 'menu', 'restaurant', 'ingredient'],
  sports: ['team', 'player', 'match', 'tournament'],
  education: ['course', 'lesson', 'student', 'instructor'],
  productivity: ['task', 'project', 'calendar', 'note'],
  reference: ['country', 'language', 'currency', 'timezone'],
}

export function ResourceSelector({}: ResourceSelectorProps) {
  const [search, setSearch] = useState('')
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(
    new Set(['people', 'business', 'commerce'])
  )
  const router = useRouter()

  // Filter resources based on search
  const filteredResources = useMemo(() => {
    if (!search) return mockResources

    const filtered: Record<string, string[]> = {}
    Object.entries(mockResources).forEach(([category, resources]) => {
      const matchedResources = resources.filter(resource =>
        resource.toLowerCase().includes(search.toLowerCase())
      )
      if (matchedResources.length > 0) {
        filtered[category] = matchedResources
      }
    })
    return filtered
  }, [search])

  const toggleCategory = (category: string) => {
    const newExpanded = new Set(expandedCategories)
    if (newExpanded.has(category)) {
      newExpanded.delete(category)
    } else {
      newExpanded.add(category)
    }
    setExpandedCategories(newExpanded)
  }

  const handleResourceClick = (resource: string) => {
    router.push(`/playground/${resource}`)
  }

  const totalResources = Object.values(filteredResources).reduce(
    (sum, resources) => sum + resources.length,
    0
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">
          Select a Resource to Test
        </h2>
        <p className="text-slate-400">
          Choose from {totalResources} available resources across {Object.keys(filteredResources).length} categories
        </p>
      </div>

      {/* Search */}
      <div className="relative">
        <input
          type="text"
          placeholder="Search resources..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full px-10 py-3 bg-slate-900 border border-white/10 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-primary-500 transition-all font-medium"
        />
        <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500" />
        {search && (
          <button
            onClick={() => setSearch('')}
            className="absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        )}
      </div>

      {/* Categories */}
      <div className="space-y-4">
        {Object.entries(filteredResources).map(([category, resources]) => {
          const isExpanded = expandedCategories.has(category)
          const Icon = getGroupIcon(category)

          return (
            <div
              key={category}
              className="bg-slate-800/40 backdrop-blur-md border border-white/5 rounded-xl overflow-hidden shadow-premium transition-all"
            >
              {/* Category Header */}
              <button
                onClick={() => toggleCategory(category)}
                className="w-full flex items-center justify-between p-4 hover:bg-white/5 transition-colors group"
              >
                <div className="flex items-center gap-4">
                  <div className="p-2 bg-primary-500/10 rounded-lg group-hover:bg-primary-500/20 transition-colors">
                    <Icon className="w-5 h-5 text-primary-400" />
                  </div>
                  <div className="text-left">
                    <span className="block text-lg font-bold text-white capitalize leading-none mb-1">
                      {category}
                    </span>
                    <span className="text-xs text-slate-500 font-medium">
                      {resources.length} resource{resources.length !== 1 ? 's' : ''}
                    </span>
                  </div>
                </div>
                <div className={`text-slate-500 transition-transform duration-200 ${isExpanded ? 'rotate-180' : ''}`}>
                  <ChevronDown className="w-5 h-5" />
                </div>
              </button>

              {/* Resources List */}
              {isExpanded && (
                <div className="border-t border-white/5 p-4 bg-slate-900/30">
                  <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
                    {resources.map((resource) => (
                      <button
                        key={resource}
                        onClick={() => handleResourceClick(resource)}
                        className="px-3 py-2 bg-slate-800/50 hover:bg-primary-500 text-slate-300 hover:text-white rounded-lg border border-white/5 text-left transition-all text-sm font-medium active:scale-95"
                      >
                        {resource}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>

      {/* No Results */}
      {totalResources === 0 && (
        <div className="text-center py-12">
          <p className="text-slate-400 text-lg">
            No resources found matching &quot;{search}&quot;
          </p>
          <button
            onClick={() => setSearch('')}
            className="mt-4 text-primary-400 hover:text-primary-300 underline font-medium"
          >
            Clear search
          </button>
        </div>
      )}
    </div>
  )
}
