'use client'

import Link from 'next/link'
import { useAuthContext } from './AuthProvider'
import { LayoutDashboard, LogIn } from 'lucide-react'

export function HeaderAuth() {
  const { user, isLoading } = useAuthContext()

  if (isLoading) {
    return <div className="w-20 h-8 rounded-lg bg-white/5 animate-pulse" />
  }

  if (user) {
    return (
      <Link
        href="/dashboard"
        className="flex items-center gap-2 text-slate-300 hover:text-white transition-colors font-medium text-sm"
      >
        <LayoutDashboard className="w-4 h-4" />
        Dashboard
      </Link>
    )
  }

  return (
    <Link
      href="/login"
      className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-indigo-600 text-white text-sm font-medium hover:bg-indigo-500 transition-colors"
    >
      <LogIn className="w-4 h-4" />
      Sign in
    </Link>
  )
}
