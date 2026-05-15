import Link from 'next/link'
import Image from 'next/image'
import { Activity, Github, Compass, BookOpen } from 'lucide-react'
import { HeaderAuth } from './HeaderAuth'

interface HeaderProps {
  compact?: boolean
}

export function Header({ compact }: HeaderProps) {
  return (
    <header className="sticky top-0 z-50 border-b border-white/5 bg-slate-900/60 backdrop-blur-xl">
      <nav className={`py-4 flex items-center justify-between ${compact ? 'px-4' : 'container mx-auto px-4'}`}>
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
          <Link href="/docs" data-umami-event="nav_docs" data-umami-label="Header Docs" className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2">
            <BookOpen className="w-4 h-4" /> Docs
          </Link>
          <Link href="/templates" data-umami-event="nav_tpl" data-umami-label="Header Explore" className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2">
            <Compass className="w-4 h-4" /> Explore
          </Link>
          <a 
            href="https://github.com/0xdps/api-mockly" 
            target="_blank" 
            rel="noopener noreferrer"
            data-umami-event="nav_gh"
            data-umami-label="Header GitHub"
            className="text-slate-300 hover:text-white transition-colors font-medium text-sm flex items-center gap-2"
          >
            <Github className="w-4 h-4" /> GitHub
          </a>
          <HeaderAuth />
        </div>
      </nav>
    </header>
  )
}
