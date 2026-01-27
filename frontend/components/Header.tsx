import Link from 'next/link'
import Image from 'next/image'

export function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-white/5 bg-slate-900/60 backdrop-blur-xl">
      <nav className="container mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-3 text-2xl font-bold text-white group transition-all">
          <Image 
            src="/favicon.svg" 
            alt="Mockly Logo" 
            width={32} 
            height={32}
            className="w-8 h-8 transition-transform group-hover:scale-110"
          />
          <span className="tracking-tight">Mockly</span>
        </Link>
        
        <div className="flex items-center gap-6">
          <Link href="/playground" className="text-slate-300 hover:text-white transition-colors font-medium text-sm">
            Playground
          </Link>
          <Link href="/docs" className="text-slate-300 hover:text-white transition-colors font-medium text-sm">
            Documentation
          </Link>
          <a 
            href="https://github.com/0xdps/fake-stack" 
            target="_blank" 
            rel="noopener noreferrer"
            className="text-slate-300 hover:text-white transition-colors font-medium text-sm"
          >
            GitHub
          </a>
        </div>
      </nav>
    </header>
  )
}
