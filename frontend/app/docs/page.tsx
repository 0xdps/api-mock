import { Header } from '@/components/Header'
import { DocsContent } from '@/components/DocsContent'
import fs from 'fs'
import path from 'path'

async function getSchemas() {
  // Schemas are in the parent directory (repository root), organized by group
  const schemasDir = path.join(process.cwd(), '../shared/schemas')
  const schemas: Array<{ name: string; schema: any; group: string }> = []
  
  // Read all subdirectories (groups)
  const entries = fs.readdirSync(schemasDir, { withFileTypes: true })
  
  for (const entry of entries) {
    if (entry.isDirectory()) {
      const groupDir = path.join(schemasDir, entry.name)
      const files = fs.readdirSync(groupDir).filter((f: string) => f.endsWith('.json'))
      
      for (const file of files) {
        const content = fs.readFileSync(path.join(groupDir, file), 'utf-8')
        schemas.push({
          name: file.replace('.json', ''),
          schema: JSON.parse(content),
          group: entry.name
        })
      }
    }
  }
  
  return schemas.sort((a, b) => a.name.localeCompare(b.name))
}

export default async function DocsPage() {
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
      <Header />
      
      <DocsContent schemas={schemas} groups={groups} />
    </div>
  )
}
