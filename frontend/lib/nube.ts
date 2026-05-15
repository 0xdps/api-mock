import { NubeAuthClient } from '@nube-auth/client'

const GATEWAY_URL = process.env.NEXT_PUBLIC_NUBE_GATEWAY_URL || 'https://api.nubeauth.com'
const APP_ID = process.env.NEXT_PUBLIC_NUBE_APP_ID || ''
const TOKEN_KEY = 'mockly_session_token'
const VERIFIER_KEY = 'mockly_pkce_verifier'

/** Return stored session token (client-side only). */
export function getStoredToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem(TOKEN_KEY)
}

/** Persist a session token. */
export function setStoredToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

/** Remove stored token (logout). */
export function clearStoredToken(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(VERIFIER_KEY)
}

/** Build an unauthenticated client (for buildOAuthUrl / exchangeCode). */
function buildAnonClient(): NubeAuthClient {
  return new NubeAuthClient({ gatewayUrl: GATEWAY_URL, appId: APP_ID })
}

/** Build an authenticated client using the stored session token. */
export function buildAuthedClient(token: string): NubeAuthClient {
  return new NubeAuthClient({ gatewayUrl: GATEWAY_URL, appId: APP_ID, sessionToken: token })
}

/** Build an authenticated client or throw if no token is stored. */
export function requireAuthedClient(): NubeAuthClient {
  const token = getStoredToken()
  if (!token) throw new Error('Not authenticated')
  return buildAuthedClient(token)
}

/**
 * Initiate the Google OAuth flow via nube-auth.
 * Stores the PKCE verifier and redirects the browser to the OAuth URL.
 */
export async function startGoogleLogin(returnTo: string): Promise<void> {
  const client = buildAnonClient()
  const { url, codeVerifier } = await client.app.buildOAuthUrl({ returnTo })
  localStorage.setItem(VERIFIER_KEY, codeVerifier)
  window.location.href = url
}

/**
 * Exchange the OAuth callback code for a session token.
 * Returns the token on success, null on failure.
 */
export async function exchangeOAuthCode(code: string): Promise<{ token: string } | { error: string }> {
  const codeVerifier = localStorage.getItem(VERIFIER_KEY) || undefined
  if (!codeVerifier) {
    console.warn('[nube] PKCE verifier missing from localStorage — was the login flow started in this browser?')
  }
  const client = buildAnonClient()
  try {
    const result = await client.app.exchangeCode(code, { codeVerifier })
    setStoredToken(result.sessionToken)
    localStorage.removeItem(VERIFIER_KEY)
    return { token: result.sessionToken }
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err)
    console.error('[nube] exchangeCode failed:', msg)
    return { error: msg }
  }
}
