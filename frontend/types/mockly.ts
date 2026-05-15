// Mockly platform types — NOT auto-generated, do not overwrite

export type TemplateVisibility = 'public' | 'private'
export type TemplateType = 'custom' | 'fork'
export type APIKeyType = 'personal' | 'access'
export type APIKeyStatus = 'active' | 'revoked'

export interface MocklyUser {
  id: string
  nube_id: string
  email: string
  name: string
  created_at: string
}

export interface MocklyTemplate {
  id: string
  user_id: string
  name: string
  slug: string
  description: string
  visibility: TemplateVisibility
  type: TemplateType
  base_schema: string
  schema_json: string
  is_featured?: boolean
  created_at: string
  updated_at: string
}

export interface MocklyAPIKey {
  id: string
  user_id: string
  name: string
  key_prefix: string
  type: APIKeyType
  status: APIKeyStatus
  last_used_at: string | null
  created_at: string
}

export interface CreateTemplatePayload {
  name: string
  slug: string
  description?: string
  visibility: TemplateVisibility
  type: TemplateType
  base_schema?: string
  schema_json: string
}
