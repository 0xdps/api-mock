import Link from 'next/link'
import Image from 'next/image'
import { Activity, BookOpen, Github } from 'lucide-react'

export function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-white/5 bg-slate-900/60 backdrop-blur-xl">
      <nav className="container mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-3 text-2xl font-bold text-white group transition-all">
          <Image 
            src="/logo.svg" 
            alt="Mockly Logo" 
            width={32} 
            height={32}
            className="w-8 h-8 transition-transform group-hover:scale-110"
          />
          <span className="tracking-tight">Mockly</span>
        </Link>
        
        <div className="flex items-center gap-8">
          <Link href="/playground" className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2">
            <Activity className="w-4 h-4" /> Playground
          </Link>
          <Link href="/resources" className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2">
            <BookOpen className="w-4 h-4" /> Resources
          </Link>
          <a 
            href="https://github.com/0xdps/fake-stack" 
            target="_blank" 
            rel="noopener noreferrer"
            className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2"
          >
            <Github className="w-4 h-4" /> GitHub
          </a>
        </div>
      </nav>
    </header>
  )
}
