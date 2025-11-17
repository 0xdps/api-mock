import { Header } from '@/components/Header'
import { Footer } from '@/components/Footer'
import { CodeExample } from '@/components/CodeExample'
import fs from 'fs'
import path from 'path'

const API_URL = process.env.NEXT_PUBLIC_API_URL || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:8080')

async function getSchemas() {
  const schemasDir = path.join(process.cwd(), 'shared/schemas')
  const files = fs.readdirSync(schemasDir).filter((f: string) => f.endsWith('.json'))
  
  return files.map((file: string) => {
    const content = fs.readFileSync(path.join(schemasDir, file), 'utf-8')
    return {
      name: file.replace('.json', ''),
      schema: JSON.parse(content)
    }
  })
}

export default async function DocsPage() {
  const schemas = await getSchemas()
  
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <Header />
      
      <main className="container mx-auto px-4 py-12 max-w-6xl">
        <h1 className="text-5xl font-bold text-white mb-6">Documentation</h1>
        <p className="text-xl text-slate-300 mb-12">
          Complete API reference and usage examples
        </p>
        
        {/* Quick Start */}
        <section className="mb-16">
          <h2 className="text-3xl font-bold text-white mb-6">Quick Start</h2>
          
          <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700 mb-6">
            <h3 className="text-xl font-semibold text-white mb-4">Base URL</h3>
            <code className="text-primary-500 text-lg">{API_URL}/api</code>
          </div>
          
          <div className="grid md:grid-cols-2 gap-6 mb-6">
            <div>
              <h3 className="text-xl font-semibold text-white mb-4">JavaScript/Fetch</h3>
              <CodeExample 
                title="Fetch Users"
                code={`fetch('${API_URL}/api/users?count=10')
  .then(res => res.json())
  .then(data => console.log(data));`}
              />
            </div>
            
            <div>
              <h3 className="text-xl font-semibold text-white mb-4">cURL</h3>
              <CodeExample 
                title="Fetch Users"
                code={`curl "${API_URL}/api/users?count=10"`}
              />
            </div>
            
            <div>
              <h3 className="text-xl font-semibold text-white mb-4">Python</h3>
              <CodeExample 
                title="Fetch Users"
                code={`import requests

response = requests.get('${API_URL}/api/users?count=10')
users = response.json()
print(users)`}
              />
            </div>
            
            <div>
              <h3 className="text-xl font-semibold text-white mb-4">Node.js/Axios</h3>
              <CodeExample 
                title="Fetch Users"
                code={`const axios = require('axios');

const response = await axios.get('${API_URL}/api/users?count=10');
console.log(response.data);`}
              />
            </div>
          </div>
        </section>
        
        {/* Endpoints */}
        <section className="mb-16">
          <h2 className="text-3xl font-bold text-white mb-6">Common Endpoints</h2>
          
          <div className="space-y-4">
            <EndpointDoc 
              method="GET"
              path="/{resource}"
              description="Get a collection of items"
              params={[
                { name: 'count', type: 'integer', description: 'Number of items to return (default: 10, max: 100)' },
                { name: 'seed', type: 'integer', description: 'Seed for reproducible data' }
              ]}
            />
            
            <EndpointDoc 
              method="GET"
              path="/{resource}/:id"
              description="Get a single item by ID"
              params={[
                { name: 'id', type: 'integer', description: 'Item ID' }
              ]}
            />
            
            <EndpointDoc 
              method="GET"
              path="/{resource}/meta"
              description="Get resource metadata and schema"
              params={[]}
            />
          </div>
        </section>
        
        {/* Available Resources */}
        <section className="mb-16">
          <h2 className="text-3xl font-bold text-white mb-6">Available Resources</h2>
          
          <div className="space-y-8">
            {schemas.map(({ name, schema }) => (
              <ResourceDoc key={name} name={name} schema={schema} />
            ))}
          </div>
        </section>
        
        {/* Features */}
        <section className="mb-16">
          <h2 className="text-3xl font-bold text-white mb-6">Features</h2>
          
          <div className="space-y-6">
            <FeatureDoc 
              title="CORS Enabled"
              description="All endpoints support CORS, making them accessible from browser applications."
            />
            
            <FeatureDoc 
              title="No Authentication"
              description="No API keys or authentication required. Start using immediately."
            />
            
            <FeatureDoc 
              title="Schema-Driven"
              description="All resources are defined by JSON schemas. View schemas using the /meta endpoint."
            />
            
            <FeatureDoc 
              title="Custom Routes"
              description="Resources can define custom paths and aliases in their schemas."
            />
            
            <FeatureDoc 
              title="Realistic Data"
              description="Uses 50+ generators for names, emails, addresses, dates, and more."
            />
            
            <FeatureDoc 
              title="Reproducible"
              description="Use the seed parameter to get the same data across requests."
            />
          </div>
        </section>
      </main>
      
      <Footer />
    </div>
  )
}

function EndpointDoc({ method, path, description, params }: any) {
  return (
    <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
      <div className="flex items-center gap-3 mb-3">
        <span className="bg-primary-500 text-white px-3 py-1 rounded text-sm font-semibold">
          {method}
        </span>
        <code className="text-slate-300 text-lg">{path}</code>
      </div>
      <p className="text-slate-400 mb-4">{description}</p>
      
      {params.length > 0 && (
        <div>
          <h4 className="text-white font-semibold mb-2">Parameters:</h4>
          <ul className="space-y-2">
            {params.map((param: any) => (
              <li key={param.name} className="text-sm">
                <code className="text-primary-500">{param.name}</code>
                <span className="text-slate-500"> ({param.type})</span>
                <span className="text-slate-400"> - {param.description}</span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}

function ResourceDoc({ name, schema }: any) {
  const properties = schema.properties || {}
  const resourceName = schema['x-resource']?.name || name
  
  return (
    <div id={name} className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700">
      <h3 className="text-2xl font-bold text-white mb-3 capitalize">{resourceName}</h3>
      <p className="text-slate-400 mb-4">{schema.description || `${resourceName} resource`}</p>
      
      <div className="mb-4">
        <h4 className="text-white font-semibold mb-3">Properties:</h4>
        <div className="space-y-2">
          {Object.entries(properties).map(([key, value]: [string, any]) => (
            <div key={key} className="text-sm flex items-start gap-2">
              <code className="text-primary-500">{key}</code>
              <span className="text-slate-500">({value.type || 'any'})</span>
              {value.description && (
                <span className="text-slate-400">- {value.description}</span>
              )}
            </div>
          ))}
        </div>
      </div>
      
      <div className="flex gap-2">
        <a 
          href={`${API_URL}/api/${name}?count=3`}
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm bg-primary-500 hover:bg-primary-600 text-white px-4 py-2 rounded transition"
        >
          Try It →
        </a>
        <a 
          href={`${API_URL}/api/${name}/meta`}
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm bg-slate-700 hover:bg-slate-600 text-white px-4 py-2 rounded transition"
        >
          View Schema
        </a>
      </div>
    </div>
  )
}

function FeatureDoc({ title, description }: { title: string; description: string }) {
  return (
    <div className="bg-slate-800/50 backdrop-blur p-4 rounded-lg border border-slate-700">
      <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
      <p className="text-slate-400">{description}</p>
    </div>
  )
}
