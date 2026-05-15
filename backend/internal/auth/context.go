package auth

import "context"

type contextKey string

const (
	nubeUserKey  contextKey = "nube_user"
	apiKeyCtxKey contextKey = "api_key_ctx"
)

// APIKeyContext holds the resolved API key and its owner user ID.
type APIKeyContext struct {
	KeyID  string
	UserID string
	Type   string // "access" | "personal"
}

// WithNubeUser stores a NubeUser in the context.
func WithNubeUser(ctx context.Context, u *NubeUser) context.Context {
	return context.WithValue(ctx, nubeUserKey, u)
}

// GetNubeUser retrieves the NubeUser from the context. Returns nil if absent.
func GetNubeUser(ctx context.Context) *NubeUser {
	u, _ := ctx.Value(nubeUserKey).(*NubeUser)
	return u
}

// WithAPIKey stores an APIKeyContext in the context.
func WithAPIKey(ctx context.Context, k *APIKeyContext) context.Context {
	return context.WithValue(ctx, apiKeyCtxKey, k)
}

// GetAPIKey retrieves the APIKeyContext from the context. Returns nil if absent.
func GetAPIKey(ctx context.Context) *APIKeyContext {
	k, _ := ctx.Value(apiKeyCtxKey).(*APIKeyContext)
	return k
}
