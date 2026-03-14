import Link from 'next/link'
import Image from 'next/image'
import { Activity, BookOpen, Github } from 'lucide-react'

export function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-white/5 bg-slate-900/60 backdrop-blur-xl">
      <nav className="container mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" data-umami-event="nav_h" data-umami-label="Header Home" className="flex items-center gap-3 text-2xl font-bold text-white group transition-all">
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
          <Link href="/playground" data-umami-event="nav_pg" data-umami-label="Header Playground" className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2">
            <Activity className="w-4 h-4" /> Playground
          </Link>
          <Link href="/resources" data-umami-event="nav_rs" data-umami-label="Header Resources" className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2">
            <BookOpen className="w-4 h-4" /> Resources
          </Link>
          <a 
            href="https://github.com/0xdps/fake-stack" 
            target="_blank" 
            rel="noopener noreferrer"
            data-umami-event="nav_gh"
            data-umami-label="Header GitHub"
            className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2"
          >
            <Github className="w-4 h-4" /> GitHub
          </a>
        </div>
      </nav>
    </header>
  )
}
