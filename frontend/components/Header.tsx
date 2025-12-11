import Link from 'next/link'
import Image from 'next/image'

export function Header() {
  return (
    <header className="border-b border-slate-700 bg-slate-900/50 backdrop-blur">
      <nav className="container mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-3 text-2xl font-bold text-white hover:opacity-90 transition">
          <Image 
            src="/favicon.svg" 
            alt="Mockly Logo" 
            width={32} 
            height={32}
            className="w-8 h-8"
          />
          Mockly
        </Link>
        
        <div className="flex items-center gap-6">
          <Link href="/playground" className="text-slate-300 hover:text-white transition">
            Playground
          </Link>
          <Link href="/docs" className="text-slate-300 hover:text-white transition">
            Documentation
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
