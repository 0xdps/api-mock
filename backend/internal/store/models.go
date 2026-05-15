package store

// SystemUserID is the user_id used for all seeded official (featured) templates.
const SystemUserID = "system"

// User mirrors the users table in MesaHub.
// ID is the nube-auth user ID (usr_...).
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Plan      string `json:"plan"`
	CreatedAt string `json:"created_at"`
}

// Template mirrors the templates table.
type Template struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`  // "public" | "private"
	Type        string `json:"type"`        // "fork" | "custom"
	BaseSchema  string `json:"base_schema"` // nullable — name of built-in schema if fork
	SchemaJSON  string `json:"schema_json"` // raw JSON schema definition
	IsFeatured  bool   `json:"is_featured"` // true for seeded official templates
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// APIKey mirrors the api_keys table.
// KeyHash is SHA-256 of the raw token — never store the plaintext.
type APIKey struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Name       string  `json:"name"`
	KeyHash    string  `json:"key_hash"`
	KeyPrefix  string  `json:"key_prefix"`
	Type       string  `json:"type"`   // "access" | "personal"
	Status     string  `json:"status"` // "active" | "revoked"
	CreatedAt  string  `json:"created_at"`
	LastUsedAt *string `json:"last_used_at"`
}

// UsageLog mirrors the usage_logs table.
type UsageLog struct {
	ID         string `json:"id"`
	APIKeyID   string `json:"api_key_id"`
	TemplateID string `json:"template_id"`
	StatusCode int    `json:"status_code"`
	CreatedAt  string `json:"created_at"`
}

// Plan limits (TBD — change constants when values are decided).
const (
	FreeTierMaxTemplates = 10
	FreeTierRPM          = 60
)
