'use client'

import { useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import Link from 'next/link'
import { useAuthContext } from '@/components/AuthProvider'
import { Header } from '@/components/Header'
import { LayoutDashboard, FileCode2, Key, Globe, LogOut } from 'lucide-react'

const NAV_ITEMS = [
  {
    href: '/dashboard',
    label: 'Overview',
    icon: LayoutDashboard,
    iconColor: 'text-primary-400',
    activeBg: 'bg-primary-500/15',
    activeText: 'text-primary-300',
  },
  {
    href: '/dashboard/templates',
    label: 'My Templates',
    icon: FileCode2,
    iconColor: 'text-violet-400',
    activeBg: 'bg-violet-500/15',
    activeText: 'text-violet-300',
  },
  {
    href: '/dashboard/api-keys',
    label: 'API Keys',
    icon: Key,
    iconColor: 'text-amber-400',
    activeBg: 'bg-amber-500/15',
    activeText: 'text-amber-300',
  },
  {
    href: '/templates',
    label: 'Explore',
    icon: Globe,
    iconColor: 'text-emerald-400',
    activeBg: 'bg-emerald-500/15',
    activeText: 'text-emerald-300',
  },
]

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { user, isLoading, logout } = useAuthContext()
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    if (!isLoading && !user) {
      router.replace('/login')
    }
  }, [isLoading, user, router])

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center">
        <svg className="w-8 h-8 animate-spin text-primary-400" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
      </div>
    )
  }

  if (!user) return null

  return (
    <div className="h-screen overflow-hidden bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col">
      <Header compact />

      <div className="flex flex-1 min-h-0">
        {/* Sidebar — fixed height, no scroll */}
        <aside className="w-56 shrink-0 border-r border-white/5 bg-slate-900/60 flex flex-col">
          <div className="px-4 py-3 border-b border-white/5">
            <p className="text-xs text-slate-500 truncate">{user.email}</p>
          </div>

          <nav className="flex-1 p-3 space-y-0.5">
            {NAV_ITEMS.map(({ href, label, icon: Icon, iconColor, activeBg, activeText }) => {
              const active = pathname === href || (href !== '/dashboard' && pathname.startsWith(href))
              return (
                <Link
                  key={href}
                  href={href}
                  className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all ${
                    active
                      ? `${activeBg} ${activeText} font-medium`
                      : 'text-slate-400 hover:text-white hover:bg-white/5'
                  }`}
                >
                  <Icon className={`w-4 h-4 shrink-0 ${active ? iconColor : iconColor + ' opacity-60'}`} />
                  {label}
                </Link>
              )
            })}
          </nav>

          <div className="p-3 border-t border-white/5">
            <button
              onClick={logout}
              className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-500 hover:text-red-400 hover:bg-red-500/10 transition-all w-full"
            >
              <LogOut className="w-4 h-4 shrink-0" />
              Sign out
            </button>
          </div>
        </aside>

        {/* Main content — scrollable */}
        <main className="flex-1 overflow-auto">
          {children}
        </main>
      </div>
    </div>
  )
}
