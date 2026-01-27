package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func TestDelayMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		delay       string
		minDuration time.Duration
		maxDuration time.Duration
	}{
		{"No delay parameter", "", 0, 10 * time.Millisecond},
		{"Valid delay 100ms", "100", 90 * time.Millisecond, 150 * time.Millisecond},
		{"Valid delay 500ms", "500", 480 * time.Millisecond, 550 * time.Millisecond},
		{"Capped at 30s", "35000", 29 * time.Second, 31 * time.Second},
		{"Invalid delay", "invalid", 0, 10 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := DelayMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/?delay="+tt.delay, nil)
			w := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(w, req)
			duration := time.Since(start)

			if duration < tt.minDuration || duration > tt.maxDuration {
				t.Errorf("Duration %v not in expected range [%v, %v]", duration, tt.minDuration, tt.maxDuration)
			}

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestFlakyMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		rate        string
		expectFails bool
	}{
		{"No rate parameter", "", false},
		{"Success rate 1.0", "1.0", false},
		{"Success rate 0.0", "0.0", true},
		{"Invalid rate", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := FlakyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			req := httptest.NewRequest("GET", "/?flakyRate="+tt.rate, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if tt.expectFails {
				if w.Code != http.StatusServiceUnavailable {
					t.Errorf("Expected status 503 for rate %s, got %d", tt.rate, w.Code)
				}

				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if _, ok := response["error"]; !ok {
					t.Error("Expected error field in response")
				}
			} else if tt.rate == "1.0" {
				if w.Code != http.StatusOK {
					t.Errorf("Expected status 200 for rate 1.0, got %d", w.Code)
				}
			}
		})
	}
}

func TestFieldFilterMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		fields   string
		expected []string
	}{
		{"No fields", "", nil},
		{"Single field", "name", []string{"name"}},
		{"Multiple fields", "id,name,email", []string{"id", "name", "email"}},
		{"Fields with spaces", "id,%20name%20,%20email", []string{"id", "name", "email"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := FieldFilterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fields, ok := GetFields(r.Context())
				if tt.expected == nil {
					if ok {
						t.Error("Expected no fields, but got some")
					}
				} else {
					if !ok {
						t.Fatal("Expected fields in context")
					}
					if len(fields) != len(tt.expected) {
						t.Errorf("Expected %d fields, got %d", len(tt.expected), len(fields))
					}
					for i, field := range fields {
						if field != tt.expected[i] {
							t.Errorf("Expected field %s, got %s", tt.expected[i], field)
						}
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/?fields="+tt.fields, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}

func TestCacheBypassMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		skipCache      string
		expectedBypass bool
	}{
		{"No skip_cache parameter", "", false},
		{"skip_cache=true", "true", true},
		{"skip_cache=1", "1", true},
		{"skip_cache=false", "false", false},
		{"skip_cache=0", "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CacheBypassMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				skip := ShouldSkipCache(r.Context())
				if skip != tt.expectedBypass {
					t.Errorf("Expected skip=%v, got %v", tt.expectedBypass, skip)
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/?skip_cache="+tt.skipCache, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}

func TestTenantMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		tenantID       string
		expectedTenant string
		shouldExist    bool
	}{
		{"No tenant header", "", "", false},
		{"With tenant header", "tenant-123", "tenant-123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := TenantMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tenantID, ok := GetTenantID(r.Context())
				if tt.shouldExist {
					if !ok {
						t.Fatal("Expected tenant ID in context")
					}
					if tenantID != tt.expectedTenant {
						t.Errorf("Expected tenant %s, got %s", tt.expectedTenant, tenantID)
					}
				} else {
					if ok {
						t.Error("Expected no tenant ID in context")
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			if tt.tenantID != "" {
				req.Header.Set("X-Tenant-ID", tt.tenantID)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}

func TestRBACMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		expectedRole string
		shouldExist  bool
	}{
		{"No role header", "", "", false},
		{"With admin role", "admin", "admin", true},
		{"With user role", "user", "user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RBACMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				role, ok := GetRole(r.Context())
				if tt.shouldExist {
					if !ok {
						t.Fatal("Expected role in context")
					}
					if role != tt.expectedRole {
						t.Errorf("Expected role %s, got %s", tt.expectedRole, role)
					}
				} else {
					if ok {
						t.Error("Expected no role in context")
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			if tt.role != "" {
				req.Header.Set("X-Role", tt.role)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}

func TestIdempotencyMiddleware(t *testing.T) {
	cache := NewSimpleIdempotencyCache()

	tests := []struct {
		name             string
		method           string
		idempotencyKey   string
		firstStatusCode  int
		firstBody        string
		secondStatusCode int
		secondBody       string
		shouldUseCached  bool
	}{
		{
			"POST with idempotency key - second request uses cached",
			"POST",
			"key-123",
			201,
			`{"id":1,"status":"created"}`,
			201,
			`{"id":1,"status":"created"}`,
			true,
		},
		{
			"POST without idempotency key - both execute",
			"POST",
			"",
			201,
			`{"id":1}`,
			201,
			`{"id":2}`,
			false,
		},
		{
			"GET ignored by idempotency middleware",
			"GET",
			"key-456",
			200,
			`{"data":"first"}`,
			200,
			`{"data":"second"}`,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			handler := IdempotencyMiddleware(cache)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.Header().Set("Content-Type", "application/json")
				if tt.method == "POST" {
					if requestCount == 1 {
						w.WriteHeader(tt.firstStatusCode)
						w.Write([]byte(tt.firstBody))
					} else {
						// Second request should have different data if not cached
						w.WriteHeader(tt.firstStatusCode)
						w.Write([]byte(`{"id":2,"status":"different"}`))
					}
				} else {
					// GET requests
					w.WriteHeader(http.StatusOK)
					if requestCount == 1 {
						w.Write([]byte(tt.firstBody))
					} else {
						w.Write([]byte(tt.secondBody))
					}
				}
			}))

			// First request
			req1 := httptest.NewRequest(tt.method, "/resource", strings.NewReader("{}"))
			if tt.idempotencyKey != "" {
				req1.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}
			w1 := httptest.NewRecorder()
			handler.ServeHTTP(w1, req1)

			if w1.Code != tt.firstStatusCode {
				t.Errorf("First request: expected status %d, got %d", tt.firstStatusCode, w1.Code)
			}

			// Second request
			requestCount = 2 // Reset for clarity
			req2 := httptest.NewRequest(tt.method, "/resource", strings.NewReader("{}"))
			if tt.idempotencyKey != "" {
				req2.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}
			w2 := httptest.NewRecorder()
			handler.ServeHTTP(w2, req2)

			if w2.Code != tt.secondStatusCode {
				t.Errorf("Second request: expected status %d, got %d", tt.secondStatusCode, w2.Code)
			}

			// Check if response was cached
			if tt.shouldUseCached {
				if w1.Body.String() != w2.Body.String() {
					t.Errorf("Expected cached response, but got different responses:\nFirst: %s\nSecond: %s",
						w1.Body.String(), w2.Body.String())
				}
			}
		})
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		providedRequestID string
		shouldGenerate    bool
	}{
		{"No request ID provided", "", true},
		{"Request ID provided", "custom-request-id-123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := r.Context().Value(chimiddleware.RequestIDKey)
				if requestID == nil {
					t.Fatal("Expected request ID in context")
				}

				requestIDStr, ok := requestID.(string)
				if !ok {
					t.Fatal("Request ID is not a string")
				}

				if tt.shouldGenerate {
					if len(requestIDStr) != 32 { // hex encoded 16 bytes
						t.Errorf("Expected generated request ID length 32, got %d", len(requestIDStr))
					}
				} else {
					if requestIDStr != tt.providedRequestID {
						t.Errorf("Expected request ID %s, got %s", tt.providedRequestID, requestIDStr)
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			if tt.providedRequestID != "" {
				req.Header.Set("X-Request-ID", tt.providedRequestID)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Check response header
			responseRequestID := w.Header().Get("X-Request-ID")
			if tt.shouldGenerate {
				if len(responseRequestID) != 32 {
					t.Errorf("Expected response request ID length 32, got %d", len(responseRequestID))
				}
			} else {
				if responseRequestID != tt.providedRequestID {
					t.Errorf("Expected response request ID %s, got %s", tt.providedRequestID, responseRequestID)
				}
			}
		})
	}
}

func TestFilterFields(t *testing.T) {
	data := map[string]interface{}{
		"id":    1,
		"name":  "John",
		"email": "john@example.com",
		"age":   30,
	}

	tests := []struct {
		name     string
		fields   []string
		expected map[string]interface{}
	}{
		{"No fields - return all", []string{}, data},
		{"Single field", []string{"name"}, map[string]interface{}{"name": "John"}},
		{"Multiple fields", []string{"id", "email"}, map[string]interface{}{"id": 1, "email": "john@example.com"}},
		{"Non-existent field", []string{"address"}, map[string]interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterFields(data, tt.fields)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d fields, got %d", len(tt.expected), len(result))
			}
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("Expected %s=%v, got %v", key, expectedVal, result[key])
				}
			}
		})
	}
}

func TestFilterFieldsSlice(t *testing.T) {
	data := []map[string]interface{}{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
	}

	tests := []struct {
		name     string
		fields   []string
		expected []map[string]interface{}
	}{
		{"No fields - return all", []string{}, data},
		{"Single field", []string{"name"}, []map[string]interface{}{
			{"name": "Alice"},
			{"name": "Bob"},
		}},
		{"Multiple fields", []string{"id", "name"}, []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterFieldsSlice(data, tt.fields)
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d items, got %d", len(tt.expected), len(result))
			}
			for i, expected := range tt.expected {
				if len(result[i]) != len(expected) {
					t.Errorf("Item %d: expected %d fields, got %d", i, len(expected), len(result[i]))
				}
				for key, expectedVal := range expected {
					if result[i][key] != expectedVal {
						t.Errorf("Item %d: expected %s=%v, got %v", i, key, expectedVal, result[i][key])
					}
				}
			}
		})
	}
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()

	// Test GetFields
	ctx = context.WithValue(ctx, FieldsKey, []string{"id", "name"})
	fields, ok := GetFields(ctx)
	if !ok {
		t.Error("Expected fields in context")
	}
	if len(fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fields))
	}

	// Test GetTenantID
	ctx = context.WithValue(ctx, TenantIDKey, "tenant-123")
	tenantID, ok := GetTenantID(ctx)
	if !ok {
		t.Error("Expected tenant ID in context")
	}
	if tenantID != "tenant-123" {
		t.Errorf("Expected tenant-123, got %s", tenantID)
	}

	// Test GetRole
	ctx = context.WithValue(ctx, RoleKey, "admin")
	role, ok := GetRole(ctx)
	if !ok {
		t.Error("Expected role in context")
	}
	if role != "admin" {
		t.Errorf("Expected admin, got %s", role)
	}

	// Test ShouldSkipCache
	ctx = context.WithValue(ctx, SkipCacheKey, true)
	skip := ShouldSkipCache(ctx)
	if !skip {
		t.Error("Expected skip cache to be true")
	}

	// Test GetIdempotencyKey
	ctx = context.WithValue(ctx, IdempotencyKeyKey, "key-456")
	key, ok := GetIdempotencyKey(ctx)
	if !ok {
		t.Error("Expected idempotency key in context")
	}
	if key != "key-456" {
		t.Errorf("Expected key-456, got %s", key)
	}
}
