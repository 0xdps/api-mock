import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'
import { PlaygroundClient } from '@/components/PlaygroundClient'
import { getApiUrl } from '@/lib/api'

const API_URL = getApiUrl()

async function getApiData() {
  try {
    const res = await fetch(`${API_URL}/`, { 
      next: { revalidate: 300 } // Revalidate every 5 minutes (ISR)
    })
    
    if (!res.ok) {
      throw new Error(`API returned ${res.status}`)
    }
    
    const data = await res.json()
    return {
      resources: data.resources || [],
      groups: data.groups || {}
    }
  } catch (error) {
    console.error('Failed to fetch API data:', error)
    // Fallback to known resources
    return {
      resources: ['users', 'posts', 'products', 'comments'],
      groups: {}
    }
  }
}

export default async function PlaygroundPage() {
  const { resources, groups } = await getApiData()
  
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <Header />
      
      <main className="container mx-auto px-4 py-12 max-w-6xl">
        <h1 className="text-5xl font-bold text-white mb-6">API Playground</h1>
        <p className="text-xl text-slate-300 mb-12">
          Test {resources.length} API endpoints interactively across {Object.keys(groups).length} categories
        </p>
        
        <PlaygroundClient resources={resources} groups={groups} />
      </main>
      
      <Footer />
    </div>
  )
}
