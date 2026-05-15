import { schemasManifest } from './schemas-manifest'

export interface SitemapUrl {
  loc: string
  lastmod?: string
  changefreq?: 'always' | 'hourly' | 'daily' | 'weekly' | 'monthly' | 'yearly' | 'never'
  priority?: number
}

export function getAllResources(): Array<{ name: string; group: string }> {
  return schemasManifest.map(({ name, group }) => ({ name, group }))
}

export function getAllGroups(): Record<string, string[]> {
  const groups: Record<string, string[]> = {}
  for (const { name, group } of schemasManifest) {
    if (!groups[group]) {
      groups[group] = []
    }
    groups[group].push(name)
  }
  return groups
}

export function getBaseUrl(): string {
  if (process.env.NEXT_PUBLIC_SITE_URL) {
    return process.env.NEXT_PUBLIC_SITE_URL
  }
  return process.env.NODE_ENV === 'production'
    ? 'https://www.mockly.codes'
    : 'http://localhost:3000'
}

