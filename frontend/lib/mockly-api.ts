import { getApiUrl } from './api'
import type { MocklyUser, MocklyTemplate, MocklyAPIKey, CreateTemplatePayload } from '../types/mockly'

function authHeaders(token: string): HeadersInit {
  return { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }
}

async function apiFetch<T>(path: string, token: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${getApiUrl()}${path}`, {
    ...init,
    headers: { ...authHeaders(token), ...(init?.headers ?? {}) },
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error((body as any)?.error || `HTTP ${res.status}`)
  }
  return res.json()
}

// ── User ─────────────────────────────────────────────────────────────────────

export async function syncUser(token: string): Promise<{ user: MocklyUser; created: boolean; access_key?: string }> {
  return apiFetch('/auth/sync', token, { method: 'POST' })
}

export async function getMe(token: string): Promise<MocklyUser> {
  return apiFetch('/me', token)
}

// ── Templates ─────────────────────────────────────────────────────────────────

export async function listMyTemplates(token: string): Promise<{ templates: MocklyTemplate[]; count: number }> {
  return apiFetch('/me/templates', token)
}

export async function createTemplate(token: string, payload: CreateTemplatePayload): Promise<MocklyTemplate> {
  return apiFetch('/me/templates', token, { method: 'POST', body: JSON.stringify(payload) })
}

export async function getTemplate(token: string, id: string): Promise<MocklyTemplate> {
  return apiFetch(`/me/templates/${id}`, token)
}

export async function updateTemplate(token: string, id: string, patch: Partial<MocklyTemplate>): Promise<MocklyTemplate> {
  return apiFetch(`/me/templates/${id}`, token, { method: 'PUT', body: JSON.stringify(patch) })
}

export async function deleteTemplate(token: string, id: string): Promise<void> {
  await apiFetch(`/me/templates/${id}`, token, { method: 'DELETE' })
}

export async function listPublicTemplates(limit = 20, offset = 0): Promise<{ templates: MocklyTemplate[]; limit: number; offset: number; total: number }> {
  const res = await fetch(`${getApiUrl()}/templates?limit=${limit}&offset=${offset}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getPublicTemplate(id: string): Promise<{ template: MocklyTemplate }> {
  const res = await fetch(`${getApiUrl()}/templates/${id}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

// ── API Keys ──────────────────────────────────────────────────────────────────

export async function listAPIKeys(token: string): Promise<{ keys: MocklyAPIKey[]; count: number }> {
  return apiFetch('/me/api-keys', token)
}

export async function createAPIKey(token: string, name: string): Promise<{ key: MocklyAPIKey; raw_key: string }> {
  return apiFetch('/me/api-keys', token, { method: 'POST', body: JSON.stringify({ name }) })
}

export async function revokeAPIKey(token: string, id: string): Promise<void> {
  await apiFetch(`/me/api-keys/${id}`, token, { method: 'DELETE' })
}

export async function getAccessKeyInfo(token: string): Promise<MocklyAPIKey> {
  return apiFetch('/me/access-key', token)
}

export async function regenerateAccessKey(token: string): Promise<{ key: MocklyAPIKey; raw_key: string }> {
  return apiFetch('/me/access-key/regenerate', token, { method: 'POST' })
}
