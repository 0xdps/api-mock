'use client'

import { createContext, useContext, useEffect, useState, useCallback, ReactNode } from 'react'
import { NubeAuthProvider } from '@nube-auth/react'
import { getStoredToken, clearStoredToken, buildAuthedClient } from '@/lib/nube'
import { syncUser } from '@/lib/mockly-api'
import type { MocklyUser } from '@/types/mockly'

const GATEWAY_URL = process.env.NEXT_PUBLIC_NUBE_GATEWAY_URL || 'https://api.nubeauth.com'
const APP_ID = process.env.NEXT_PUBLIC_NUBE_APP_ID || ''

// ── Context ──────────────────────────────────────────────────────────────────

interface AuthContextValue {
  token: string | null
  user: MocklyUser | null
  isLoading: boolean
  logout: () => void
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue>({
  token: null,
  user: null,
  isLoading: true,
  logout: () => {},
  refreshUser: async () => {},
})

export function useAuthContext() {
  return useContext(AuthContext)
}

// ── Provider ──────────────────────────────────────────────────────────────────

interface AuthProviderProps {
  children: ReactNode
}

function InnerAuthProvider({ children }: AuthProviderProps) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<MocklyUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const loadUser = useCallback(async (tok: string) => {
    try {
      const { user: syncedUser } = await syncUser(tok)
      setUser(syncedUser)
    } catch {
      // Token may be expired — clear it
      clearStoredToken()
      setToken(null)
      setUser(null)
    }
  }, [])

  const refreshUser = useCallback(async () => {
    const tok = getStoredToken()
    if (tok) await loadUser(tok)
  }, [loadUser])

  const logout = useCallback(() => {
    clearStoredToken()
    setToken(null)
    setUser(null)
    // Also sign out from nube-auth if possible
    const tok = getStoredToken()
    if (tok) {
      buildAuthedClient(tok).auth.logout().catch(() => {})
    }
  }, [])

  useEffect(() => {
    const tok = getStoredToken()
    if (tok) {
      setToken(tok)
      loadUser(tok).finally(() => setIsLoading(false))
    } else {
      setIsLoading(false)
    }
  }, [loadUser])

  const nubeConfig = {
    gatewayUrl: GATEWAY_URL,
    appId: APP_ID,
    ...(token ? { sessionToken: token } : {}),
  }

  return (
    <NubeAuthProvider config={nubeConfig}>
      <AuthContext.Provider value={{ token, user, isLoading, logout, refreshUser }}>
        {children}
      </AuthContext.Provider>
    </NubeAuthProvider>
  )
}

export function AuthProvider({ children }: AuthProviderProps) {
  return <InnerAuthProvider>{children}</InnerAuthProvider>
}
