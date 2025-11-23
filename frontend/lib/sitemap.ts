import fs from 'fs'
import path from 'path'

export interface SitemapUrl {
  loc: string
  lastmod?: string
  changefreq?: 'always' | 'hourly' | 'daily' | 'weekly' | 'monthly' | 'yearly' | 'never'
  priority?: number
}

export function getAllResources(): Array<{ name: string; group: string }> {
  const schemasDir = path.join(process.cwd(), '../shared/schemas')
  const resources: Array<{ name: string; group: string }> = []
  
  if (!fs.existsSync(schemasDir)) {
    return resources
  }
  
  const entries = fs.readdirSync(schemasDir, { withFileTypes: true })
  
  for (const entry of entries) {
    if (entry.isDirectory()) {
      const groupDir = path.join(schemasDir, entry.name)
      const files = fs.readdirSync(groupDir).filter((f: string) => f.endsWith('.json'))
      
      for (const file of files) {
        resources.push({
          name: file.replace('.json', ''),
          group: entry.name
        })
      }
    }
  }
  
  return resources.sort((a, b) => a.name.localeCompare(b.name))
}

export function getAllGroups(): Record<string, string[]> {
  const resources = getAllResources()
  const groups: Record<string, string[]> = {}
  
  for (const resource of resources) {
    if (!groups[resource.group]) {
      groups[resource.group] = []
    }
    groups[resource.group].push(resource.name)
  }
  
  return groups
}

export function getBaseUrl(): string {
  if (process.env.NEXT_PUBLIC_SITE_URL) {
    return process.env.NEXT_PUBLIC_SITE_URL
  }
  return process.env.NODE_ENV === 'production'
    ? 'https://mockly.codes'
    : 'http://localhost:3000'
}

