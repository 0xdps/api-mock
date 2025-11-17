import Link from 'next/link'

export function Header() {
  return (
    <header className="border-b border-slate-700 bg-slate-900/50 backdrop-blur">
      <nav className="container mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="text-2xl font-bold text-white">
          apimock<span className="text-primary-500">.codes</span>
        </Link>
        
        <div className="flex items-center gap-6">
          <Link href="/docs" className="text-slate-300 hover:text-white transition">
            Docs
          </Link>
          <Link href="/playground" className="text-slate-300 hover:text-white transition">
            Playground
          </Link>
          <a 
            href="https://github.com/0xdps/fake-stack" 
            target="_blank" 
            rel="noopener noreferrer"
            className="text-slate-300 hover:text-white transition"
          >
            GitHub
          </a>
        </div>
      </nav>
    </header>
  )
}
