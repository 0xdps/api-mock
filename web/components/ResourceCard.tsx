import Link from 'next/link'

interface ResourceCardProps {
  name: string
}

export function ResourceCard({ name }: ResourceCardProps) {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
  
  return (
    <div className="bg-slate-800/50 backdrop-blur p-6 rounded-lg border border-slate-700 hover:border-primary-500 transition">
      <h3 className="text-xl font-semibold text-white mb-3 capitalize">
        {name}
      </h3>
      
      <div className="space-y-2 text-sm">
        <EndpointLink 
          href={`${apiUrl}/api/${name}?count=5`}
          label={`GET /${name}`}
        />
        <EndpointLink 
          href={`${apiUrl}/api/${name}/1`}
          label={`GET /${name}/:id`}
        />
        <EndpointLink 
          href={`${apiUrl}/api/${name}/meta`}
          label={`GET /${name}/meta`}
        />
      </div>
      
      <Link 
        href={`/docs#${name}`}
        className="mt-4 inline-block text-primary-500 hover:text-primary-400 text-sm font-medium"
      >
        View Schema →
      </Link>
    </div>
  )
}

function EndpointLink({ href, label }: { href: string; label: string }) {
  return (
    <a 
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="block text-slate-400 hover:text-primary-500 font-mono transition"
    >
      {label}
    </a>
  )
}
