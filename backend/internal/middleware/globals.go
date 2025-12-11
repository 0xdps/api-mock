package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Context keys for global parameters
type contextKey string

const (
	FieldsKey         contextKey = "fields"
	TenantIDKey       contextKey = "tenant_id"
	RoleKey           contextKey = "role"
	SkipCacheKey      contextKey = "skip_cache"
	IdempotencyKeyKey contextKey = "idempotency_key"
)

// DelayMiddleware adds artificial latency based on delay query parameter
func DelayMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delayMs := r.URL.Query().Get("delay"); delayMs != "" {
			if ms, err := strconv.Atoi(delayMs); err == nil {
				// Cap at 30 seconds to prevent abuse
				if ms > 30000 {
					ms = 30000
				}
				if ms > 0 {
					time.Sleep(time.Duration(ms) * time.Millisecond)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// FlakyMiddleware randomly fails requests based on flakyRate parameter
func FlakyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rateStr := r.URL.Query().Get("flakyRate"); rateStr != "" {
			if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
				// Rate should be between 0 and 1
				if rate < 0 {
					rate = 0
				}
				if rate > 1 {
					rate = 1
				}

				// Generate random number between 0 and 1
				b := make([]byte, 8)
				rand.Read(b)
				randomValue := float64(uint64(b[0])%100) / 100.0

				// If random value exceeds the success rate, fail the request
				if randomValue > rate {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusServiceUnavailable)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":   "Service Unavailable",
						"message": "Simulated failure from flakyRate parameter",
						"rate":    rate,
					})
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// FieldFilterMiddleware extracts fields parameter and stores in context
func FieldFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fields := r.URL.Query().Get("fields"); fields != "" {
			// Split by comma and trim spaces
			fieldList := strings.Split(fields, ",")
			for i, field := range fieldList {
				fieldList[i] = strings.TrimSpace(field)
			}
			ctx := context.WithValue(r.Context(), FieldsKey, fieldList)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// CacheBypassMiddleware sets skip_cache flag in context
func CacheBypassMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipCache := r.URL.Query().Get("skip_cache"); skipCache == "true" || skipCache == "1" {
			ctx := context.WithValue(r.Context(), SkipCacheKey, true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// TenantMiddleware extracts X-Tenant-ID header and stores in context
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID != "" {
			ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// RBACMiddleware extracts X-Role header and stores in context
func RBACMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-Role")
		if role != "" {
			ctx := context.WithValue(r.Context(), RoleKey, role)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// CachedResponse represents a cached idempotent response
type CachedResponse struct {
	Status  int
	Headers http.Header
	Body    []byte
}

// IdempotencyCache interface for storing idempotency responses
type IdempotencyCache interface {
	Get(key string) (*CachedResponse, bool)
	Set(key string, response *CachedResponse, ttl time.Duration)
}

// SimpleIdempotencyCache is a simple in-memory implementation
type SimpleIdempotencyCache struct {
	data map[string]*cacheEntry
}

type cacheEntry struct {
	response  *CachedResponse
	expiresAt time.Time
}

// NewSimpleIdempotencyCache creates a new in-memory cache
func NewSimpleIdempotencyCache() *SimpleIdempotencyCache {
	cache := &SimpleIdempotencyCache{
		data: make(map[string]*cacheEntry),
	}
	// Start cleanup goroutine
	go cache.cleanup()
	return cache
}

func (c *SimpleIdempotencyCache) Get(key string) (*CachedResponse, bool) {
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		delete(c.data, key)
		return nil, false
	}
	return entry.response, true
}

func (c *SimpleIdempotencyCache) Set(key string, response *CachedResponse, ttl time.Duration) {
	c.data[key] = &cacheEntry{
		response:  response,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *SimpleIdempotencyCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.expiresAt) {
				delete(c.data, key)
			}
		}
	}
}

// IdempotencyMiddleware handles idempotency keys for POST/PUT/PATCH requests
func IdempotencyMiddleware(cache IdempotencyCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only apply to POST/PUT/PATCH
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Store key in context for handlers to access
			ctx := context.WithValue(r.Context(), IdempotencyKeyKey, key)
			r = r.WithContext(ctx)

			// Check if request with this key was already processed
			if cached, found := cache.Get(key); found {
				// Replay cached response
				for k, v := range cached.Headers {
					w.Header()[k] = v
				}
				w.WriteHeader(cached.Status)
				w.Write(cached.Body)
				return
			}

			// Capture response to cache it
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			// Cache the response (TTL: 24 hours)
			cached := &CachedResponse{
				Status:  rec.Code,
				Headers: rec.Header().Clone(),
				Body:    rec.Body.Bytes(),
			}
			cache.Set(key, cached, 24*time.Hour)

			// Write response to client
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
		})
	}
}

// RequestIDMiddleware generates or extracts X-Request-ID
// Note: chi already has middleware.RequestID, but we provide this for completeness
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate new request ID
			b := make([]byte, 16)
			rand.Read(b)
			requestID = hex.EncodeToString(b)
		}

		// Store in context (chi's middleware.RequestID uses "RequestID" key)
		ctx := context.WithValue(r.Context(), middleware.RequestIDKey, requestID)
		r = r.WithContext(ctx)

		// Set response header
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

// Helper functions to extract values from context

// GetFields retrieves the fields filter from context
func GetFields(ctx context.Context) ([]string, bool) {
	fields, ok := ctx.Value(FieldsKey).([]string)
	return fields, ok
}

// GetTenantID retrieves the tenant ID from context
func GetTenantID(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	return tenantID, ok
}

// GetRole retrieves the role from context
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}

// ShouldSkipCache checks if cache should be bypassed
func ShouldSkipCache(ctx context.Context) bool {
	skip, ok := ctx.Value(SkipCacheKey).(bool)
	return ok && skip
}

// GetIdempotencyKey retrieves the idempotency key from context
func GetIdempotencyKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(IdempotencyKeyKey).(string)
	return key, ok
}

// FilterFields filters a map to include only specified fields
func FilterFields(data map[string]interface{}, fields []string) map[string]interface{} {
	if len(fields) == 0 {
		return data
	}

	filtered := make(map[string]interface{})
	for _, field := range fields {
		if val, exists := data[field]; exists {
			filtered[field] = val
		}
	}
	return filtered
}

// FilterFieldsSlice filters a slice of maps to include only specified fields
func FilterFieldsSlice(data []map[string]interface{}, fields []string) []map[string]interface{} {
	if len(fields) == 0 {
		return data
	}

	filtered := make([]map[string]interface{}, len(data))
	for i, item := range data {
		filtered[i] = FilterFields(item, fields)
	}
	return filtered
}
