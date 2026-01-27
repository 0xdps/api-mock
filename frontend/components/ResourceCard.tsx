import Link from 'next/link'

interface ResourceCardProps {
  name: string
  apiUrl: string
  group?: string
}

export function ResourceCard({ name, apiUrl, group }: ResourceCardProps) {
  return (
    <div className="bg-slate-800/40 backdrop-blur-md p-6 rounded-xl border border-white/5 shadow-premium hover:border-primary-500/30 transition-all duration-300 group">
      <h3 className="text-xl font-semibold text-white mb-3 capitalize tracking-tight group-hover:text-primary-400 transition-colors">
        {name}
      </h3>
      
      <div className="space-y-2 text-sm">
        <EndpointLink 
          href={`${apiUrl}/${name}?count=5`}
          label={`GET /${name}`}
        />
        <EndpointLink 
          href={`${apiUrl}/${name}/1`}
          label={`GET /${name}/:id`}
        />
        <EndpointLink 
          href={`${apiUrl}/${name}/meta`}
          label={`GET /${name}/meta`}
        />
      </div>
      
      <Link 
        href={`/resources#${name}`}
        className="mt-6 inline-flex items-center text-primary-500 hover:text-primary-400 text-sm font-semibold transition-colors"
      >
        View Schema <span className="ml-1 opacity-70">→</span>
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
