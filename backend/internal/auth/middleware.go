package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/0xdps/api-mock/go/internal/store"
	"github.com/redis/go-redis/v9"
)

// SessionMiddleware validates a nube-auth Bearer token on every request.
// Unauthenticated requests receive a 401.
func SessionMiddleware(nube *NubeClient, rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearer(r)
			if token == "" {
				respondUnauth(w, "missing Authorization header")
				return
			}

			user, err := resolveSession(r.Context(), nube, rdb, token)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "auth service error"})
				return
			}
			if user == nil {
				respondUnauth(w, "invalid or expired token")
				return
			}

			next.ServeHTTP(w, r.WithContext(WithNubeUser(r.Context(), user)))
		})
	}
}

// APIKeyMiddleware validates an API key sent via X-API-Key header or ?api_key= query param.
// Unauthenticated requests receive a 401.
func APIKeyMiddleware(keys *store.APIKeys, rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := extractAPIKey(r)
			if raw == "" {
				respondUnauth(w, "missing API key — provide X-API-Key header or ?api_key= param")
				return
			}

			apiKey, err := resolveAPIKey(r.Context(), keys, rdb, raw)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "key lookup error"})
				return
			}
			if apiKey == nil {
				respondUnauth(w, "invalid or revoked API key")
				return
			}

			// Touch last_used_at asynchronously so it doesn't add latency
			go keys.TouchLastUsed(apiKey.KeyID)

			next.ServeHTTP(w, r.WithContext(WithAPIKey(r.Context(), apiKey)))
		})
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func extractBearer(r *http.Request) string {
	v := r.Header.Get("Authorization")
	if strings.HasPrefix(v, "Bearer ") {
		return strings.TrimPrefix(v, "Bearer ")
	}
	return ""
}

func extractAPIKey(r *http.Request) string {
	if k := r.Header.Get("X-API-Key"); k != "" {
		return k
	}
	return r.URL.Query().Get("api_key")
}

// resolveSession checks Redis first (30 s TTL) then calls nube-auth.
func resolveSession(ctx context.Context, nube *NubeClient, rdb *redis.Client, token string) (*NubeUser, error) {
	cacheKey := "session:" + token[:min(len(token), 64)]

	if rdb != nil {
		if cached, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
			var u NubeUser
			if json.Unmarshal([]byte(cached), &u) == nil {
				return &u, nil
			}
		}
	}

	user, err := nube.ValidateSession(ctx, token)
	if err != nil || user == nil {
		return user, err
	}

	if rdb != nil {
		if b, err := json.Marshal(user); err == nil {
			_ = rdb.Set(ctx, cacheKey, b, 30*time.Second).Err()
		}
	}
	return user, nil
}

// resolveAPIKey checks Redis first (60 s TTL) then queries MesaHub.
func resolveAPIKey(ctx context.Context, keys *store.APIKeys, rdb *redis.Client, raw string) (*APIKeyContext, error) {
	cacheKey := "apikey:" + raw[:min(len(raw), 64)]

	if rdb != nil {
		if cached, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
			var k APIKeyContext
			if json.Unmarshal([]byte(cached), &k) == nil {
				return &k, nil
			}
		}
	}

	k, err := keys.GetByRawKey(raw)
	if err != nil || k == nil {
		return nil, err
	}

	kctx := &APIKeyContext{KeyID: k.ID, UserID: k.UserID, Type: k.Type}

	if rdb != nil {
		if b, err := json.Marshal(kctx); err == nil {
			_ = rdb.Set(ctx, cacheKey, b, 60*time.Second).Err()
		}
	}
	return kctx, nil
}

func respondUnauth(w http.ResponseWriter, msg string) {
	respondJSON(w, http.StatusUnauthorized, map[string]string{"error": msg})
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
