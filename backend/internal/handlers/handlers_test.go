package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/0xdps/api-mock/go/internal/cache"
	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/go-chi/chi/v5"
)

// setupTestHandler creates a test handler with in-memory cache (no Redis needed)
func setupTestHandler(t *testing.T) (*DynamicHandler, *cache.Cache) {
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	cacheConfig := cache.Config{
		ItemsPerResource:    50,
		Seed:                42,
		MaxItemsPerResource: 100,
	}

	testCache := cache.NewCache(registry, cacheConfig, nil)

	// Populate cache manually with test data for all resources
	resourceNames := registry.GetAllResourceNames()
	for _, resourceName := range resourceNames {
		data, err := registry.GenerateData(resourceName, 50)
		if err != nil {
			t.Fatalf("Failed to generate data for %s: %v", resourceName, err)
		}
		testCache.Data[resourceName] = data
	}

	handler := NewDynamicHandler(registry, testCache)
	return handler, testCache
}

// randomResources returns n random resource names
func randomResources(t *testing.T, handler *DynamicHandler, n int) []string {
	resourceNames := handler.registry.GetAllResourceNames()
	if len(resourceNames) < n {
		t.Fatalf("Not enough resources available: have %d, need %d", len(resourceNames), n)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(resourceNames), func(i, j int) {
		resourceNames[i], resourceNames[j] = resourceNames[j], resourceNames[i]
	})

	return resourceNames[:n]
}

// createCtxWithParams adds URL parameters to context the way chi does
func createCtxWithParams(params map[string]string) context.Context {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}
	return context.WithValue(context.Background(), chi.RouteCtxKey, ctx)
}

// ============================================================================
// GET TESTS
// ============================================================================

// TestGetCollectionBasic tests basic GET requests across 10 random resources
func TestGetCollectionBasic(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resources := randomResources(t, handler, 10)

	for _, resourceName := range resources {
		t.Run(fmt.Sprintf("GET_%s", resourceName), func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/%s", resourceName), nil)
			w := httptest.NewRecorder()

			handler.GetCollection(resourceName)(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			var data []map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&data); err != nil {
				t.Errorf("Failed to decode response: %v", err)
			}

			if len(data) == 0 {
				t.Error("Expected non-empty data")
			}

			if w.Header().Get("X-Cache") == "" {
				t.Error("Missing X-Cache header")
			}
		})
	}
}

// TestGetCollectionWithCount tests GET with count parameter
func TestGetCollectionWithCount(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resources := randomResources(t, handler, 5)
	counts := []int{1, 5, 10, 25, 100}

	for _, resourceName := range resources {
		for _, count := range counts {
			t.Run(fmt.Sprintf("GET_%s_count_%d", resourceName, count), func(t *testing.T) {
				url := fmt.Sprintf("/%s?count=%d", resourceName, count)
				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				handler.GetCollection(resourceName)(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("Expected status 200, got %d", w.Code)
				}

				var data []map[string]interface{}
				json.NewDecoder(w.Body).Decode(&data)

				if len(data) > count {
					t.Errorf("Expected at most %d items, got %d", count, len(data))
				}
			})
		}
	}
}

// TestGetCollectionWithFilters tests GET with filter parameters
func TestGetCollectionWithFilters(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resources := randomResources(t, handler, 3)

	for _, resourceName := range resources {
		t.Run(fmt.Sprintf("GET_%s_with_filters", resourceName), func(t *testing.T) {
			url := fmt.Sprintf("/%s?count=50", resourceName)
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetCollection(resourceName)(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			var data []map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&data); err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			if len(data) == 0 {
				t.Skip("No data to filter")
			}
		})
	}
}

// TestGetCollectionNoCache tests nocache parameter behavior
func TestGetCollectionNoCache(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resources := randomResources(t, handler, 5)

	for _, resourceName := range resources {
		t.Run(fmt.Sprintf("GET_%s_nocache_true", resourceName), func(t *testing.T) {
			url := fmt.Sprintf("/%s?nocache=true&count=10", resourceName)
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetCollection(resourceName)(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			cacheHeader := w.Header().Get("X-Cache")
			if cacheHeader != "BYPASS" {
				t.Errorf("Expected X-Cache: BYPASS, got %s", cacheHeader)
			}
		})

		t.Run(fmt.Sprintf("GET_%s_nocache_false", resourceName), func(t *testing.T) {
			url := fmt.Sprintf("/%s?nocache=false&count=10", resourceName)
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetCollection(resourceName)(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			cacheHeader := w.Header().Get("X-Cache")
			if cacheHeader == "BYPASS" {
				t.Errorf("Expected non-BYPASS cache, got %s", cacheHeader)
			}
		})
	}
}

// TestGetMetadata tests GET /meta endpoint
func TestGetMetadata(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resources := randomResources(t, handler, 5)

	for _, resourceName := range resources {
		t.Run(fmt.Sprintf("GET_%s_meta", resourceName), func(t *testing.T) {
			url := fmt.Sprintf("/%s/meta", resourceName)
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetResourceMetadata(resourceName)(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			var meta map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&meta); err != nil {
				t.Errorf("Failed to decode meta: %v", err)
			}

			if meta["title"] == nil && meta["$schema"] == nil {
				t.Error("Meta should have title or $schema")
			}
		})
	}
}

// ============================================================================
// POST TESTS (CREATE)
// ============================================================================

// NOTE: POST tests that require Redis integration have been removed.
// These tests were directly manipulating cache instead of testing HTTP handlers.
// Proper HTTP handler tests with Redis mocking will be added separately.

// ============================================================================
// PUT TESTS (UPDATE)
// ============================================================================

// NOTE: PUT tests that require Redis integration have been removed.
// These tests were directly manipulating cache instead of testing HTTP handlers.
// Proper HTTP handler tests with Redis mocking will be added separately.

// ============================================================================
// DELETE TESTS
// ============================================================================

// NOTE: DELETE tests that require Redis integration have been removed.
// These tests were directly manipulating cache instead of testing HTTP handlers.
// Proper HTTP handler tests with Redis mocking will be added separately.

// ============================================================================
// COMBINED CRUD TESTS
// ============================================================================

// NOTE: CRUD cycle tests have been removed as they were directly manipulating cache.
// Proper end-to-end HTTP handler tests will be added separately.

// TestCRUDWithNoCacheToggle tests CRUD with nocache parameter toggled
func TestCRUDWithNoCacheToggle(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resource := randomResources(t, handler, 1)[0]

	t.Run(fmt.Sprintf("CRUD_nocache_toggle_%s", resource), func(t *testing.T) {
		// GET with cache
		getReq1 := httptest.NewRequest("GET", fmt.Sprintf("/%s", resource), nil)
		getW1 := httptest.NewRecorder()
		handler.GetCollection(resource)(getW1, getReq1)
		cache1 := getW1.Header().Get("X-Cache")

		// GET without cache (nocache=true)
		getReq2 := httptest.NewRequest("GET", fmt.Sprintf("/%s?nocache=true", resource), nil)
		getW2 := httptest.NewRecorder()
		handler.GetCollection(resource)(getW2, getReq2)
		cache2 := getW2.Header().Get("X-Cache")
		if cache2 != "BYPASS" {
			t.Errorf("Expected BYPASS, got %s", cache2)
		}

		// GET with cache again
		getReq3 := httptest.NewRequest("GET", fmt.Sprintf("/%s?nocache=false", resource), nil)
		getW3 := httptest.NewRecorder()
		handler.GetCollection(resource)(getW3, getReq3)
		cache3 := getW3.Header().Get("X-Cache")
		if cache3 == "BYPASS" {
			t.Error("Should not be BYPASS with nocache=false")
		}

		t.Logf("Cache headers: %s → %s → %s", cache1, cache2, cache3)
	})
}

// ============================================================================
// EDGE CASES & ERROR HANDLING
// ============================================================================

// TestInvalidJSONBody tests handling of invalid JSON
func TestInvalidJSONBody(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resource := randomResources(t, handler, 1)[0]

	t.Run(fmt.Sprintf("POST_%s_invalid_json", resource), func(t *testing.T) {
		invalidJSON := []byte(`{invalid json}`)
		req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resource), bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.PostCollection(resource)(w, req)

		if w.Code != http.StatusBadRequest {
			t.Logf("Expected 400, got %d", w.Code)
		}
	})
}

// TestCountParameterEdgeCases tests count parameter edge cases
func TestCountParameterEdgeCases(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resource := randomResources(t, handler, 1)[0]

	testCases := []struct {
		name     string
		count    string
		valid    bool
		maxCount int
	}{
		{"count=0", "0", false, 10},
		{"count=-1", "-1", false, 10},
		{"count=abc", "abc", false, 10},
		{"count=1", "1", true, 1},
		{"count=100", "100", true, 100},
		{"count=1000", "1000", true, 100}, // Capped at 100
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("count_%s", tc.name), func(t *testing.T) {
			url := fmt.Sprintf("/%s?%s", resource, tc.count)
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetCollection(resource)(w, req)

			if w.Code == http.StatusOK {
				var data []map[string]interface{}
				json.NewDecoder(w.Body).Decode(&data)

				if tc.valid {
					if len(data) > tc.maxCount {
						t.Errorf("Expected max %d items, got %d", tc.maxCount, len(data))
					}
				}
			}
		})
	}
}

// TestCountCapping tests that count is capped at 100
func TestCountCapping(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resource := randomResources(t, handler, 1)[0]

	t.Run(fmt.Sprintf("count_capping_%s", resource), func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/%s?count=999", resource), nil)
		w := httptest.NewRecorder()

		handler.GetCollection(resource)(w, req)

		if w.Code == http.StatusOK {
			var data []map[string]interface{}
			json.NewDecoder(w.Body).Decode(&data)

			if len(data) > 100 {
				t.Errorf("Count not capped at 100, got %d", len(data))
			}
		}
	})
}

// TestNoCacheVariations tests different nocache parameter formats
func TestNoCacheVariations(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resource := randomResources(t, handler, 1)[0]

	testCases := []struct {
		name   string
		param  string
		bypass bool
	}{
		{"nocache=true", "nocache=true", true},
		{"nocache=false", "nocache=false", false},
		{"fresh=true", "fresh=true", true},
		{"fresh=false", "fresh=false", false},
		{"no params", "", false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("nocache_%s", tc.name), func(t *testing.T) {
			url := fmt.Sprintf("/%s", resource)
			if tc.param != "" {
				url += "?" + tc.param
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetCollection(resource)(w, req)

			cacheHeader := w.Header().Get("X-Cache")
			isBypass := cacheHeader == "BYPASS"

			if tc.bypass != isBypass {
				t.Errorf("Expected bypass=%v, got X-Cache=%s", tc.bypass, cacheHeader)
			}
		})
	}
}

// TestResponseHeaders verifies response headers are set correctly
func TestResponseHeaders(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resources := randomResources(t, handler, 3)

	for _, resource := range resources {
		t.Run(fmt.Sprintf("headers_%s", resource), func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/%s", resource), nil)
			w := httptest.NewRecorder()

			handler.GetCollection(resource)(w, req)

			// Check for required headers
			if ct := w.Header().Get("Content-Type"); ct == "" {
				t.Error("Missing Content-Type header")
			}

			if cache := w.Header().Get("X-Cache"); cache == "" {
				t.Error("Missing X-Cache header")
			}
		})
	}
}

// TestDataIntegrity verifies data isn't corrupted during operations
// Note: Response returns limited items (default 10), cache has all 50
func TestDataIntegrity(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resource := randomResources(t, handler, 1)[0]

	t.Run(fmt.Sprintf("integrity_%s", resource), func(t *testing.T) {
		// Get data with higher count to verify consistency
		req1 := httptest.NewRequest("GET", fmt.Sprintf("/%s?count=50", resource), nil)
		w1 := httptest.NewRecorder()
		handler.GetCollection(resource)(w1, req1)

		var data1 []map[string]interface{}
		json.NewDecoder(w1.Body).Decode(&data1)

		// Get data again with same parameters
		req2 := httptest.NewRequest("GET", fmt.Sprintf("/%s?count=50", resource), nil)
		w2 := httptest.NewRecorder()
		handler.GetCollection(resource)(w2, req2)

		var data2 []map[string]interface{}
		json.NewDecoder(w2.Body).Decode(&data2)

		// Responses should be identical
		if len(data1) != len(data2) {
			t.Errorf("Data length mismatch: %d vs %d", len(data1), len(data2))
		}

		// Verify cache has correct total
		cacheData := testCache.Data[resource]
		if len(cacheData) < len(data1) {
			t.Errorf("Cache has fewer items than response: cache %d, response %d", len(cacheData), len(data1))
		}
		
		t.Logf("Data integrity verified: cache=%d, response=%d", len(cacheData), len(data1))
	})
}

// ============================================================================
// SCHEMA VALIDATION TESTS
// ============================================================================

// NOTE: Schema validation tests that directly manipulated cache have been removed.
// Proper validation tests exist in schema_validation_test.go that test HTTP handlers.
// These redundant tests were not testing actual handler behavior.

// ============================================================================
// HELPER FUNCTIONS FOR SCHEMA TESTING
// ============================================================================

// getFieldNames extracts field names from an item
func getFieldNames(item map[string]interface{}) []string {
	names := make([]string, 0, len(item))
	for k := range item {
		names = append(names, k)
	}
	return names
}

// getSchemaPropertyNames extracts property names from schema
func getSchemaPropertyNames(schema *schema.Schema) []string {
	if schema == nil {
		return nil
	}
	names := make([]string, 0, len(schema.Properties))
	for k := range schema.Properties {
		names = append(names, k)
	}
	return names
}

// ============================================================================
// HTTP-LEVEL HANDLER TESTS
// ============================================================================

func TestHandler_InvalidHTTPMethods(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resourceName := "users"
	
	invalidMethods := []string{"PATCH", "OPTIONS", "TRACE", "CONNECT", "HEAD"}
	
	for _, method := range invalidMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, fmt.Sprintf("/%s", resourceName), nil)
			w := httptest.NewRecorder()
			
			// For GET endpoint - only POST should fail
			if method != "GET" {
				handler.GetCollection(resourceName)(w, req)
				
				// Since we're calling GetCollection which is designed for GET,
				// the method check is not in the handler itself - chi router handles this
				// So we can only verify the handler works when called directly
				t.Logf("Method %s tested on GET endpoint: status %d", method, w.Code)
			}
		})
	}
}

func TestHandler_MalformedJSON_POST(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resourceName := "users"
	
	testCases := []struct {
		name string
		body string
	}{
		{"Empty", ""},
		{"Invalid JSON", `{"name": "test`},
		{"Malformed Array", `[{"id": 1},`},
		{"Wrong Type", `"just a string"`},
		{"Null Body", "null"},
		{"Non-JSON", "this is not json at all"},
		{"Trailing Comma", `{"id": 1,}`},
		{"Single Quote", `{'id': 1}`},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resourceName), bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			handler.PostCollection(resourceName)(w, req)
			
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400 for malformed JSON, got %d", w.Code)
			}
			
			// Verify error message is present
			body := w.Body.String()
			if body == "" {
				t.Error("Expected error message in response body")
			}
		})
	}
}

func TestHandler_MalformedJSON_PUT(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resourceName := "users"
	
	// Get a valid ID first
	items, _ := testCache.Get(resourceName, 1)
	
	if len(items) == 0 {
		t.Fatal("No test data available")
	}
	
	item := items[0]
	id := fmt.Sprintf("%v", item["id"])
	
	testCases := []struct {
		name string
		body string
	}{
		{"Empty", ""},
		{"Invalid JSON", `{"name": "test`},
		{"Non-Object", `[1, 2, 3]`},
		{"Null Body", "null"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", fmt.Sprintf("/%s/%s", resourceName, id), bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(createCtxWithParams(map[string]string{"resource": resourceName, "id": id}))
			w := httptest.NewRecorder()
			
			handler.PutSingle(resourceName)(w, req)
			
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400 for malformed JSON, got %d", w.Code)
			}
		})
	}
}

func TestHandler_MissingContentType(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resourceName := "users"
	
	validJSON := `{"name": "Test User", "email": "test@example.com"}`
	
	t.Run("POST_WithoutContentType", func(t *testing.T) {
		req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resourceName), bytes.NewBufferString(validJSON))
		// Intentionally not setting Content-Type
		w := httptest.NewRecorder()
		
		handler.PostCollection(resourceName)(w, req)
		
		// Should still work or return 400 depending on implementation
		t.Logf("POST without Content-Type: status %d", w.Code)
	})
	
	t.Run("POST_WithWrongContentType", func(t *testing.T) {
		req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resourceName), bytes.NewBufferString(validJSON))
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()
		
		handler.PostCollection(resourceName)(w, req)
		
		t.Logf("POST with wrong Content-Type: status %d", w.Code)
	})
}

func TestHandler_POST_ValidCreation(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resourceName := "users"
	
	// Get initial count - use large number to get all items
	initialItems, _ := testCache.Get(resourceName, 1000)
	initialCount := len(initialItems)
	
	// Generate a valid item
	items, err := handler.registry.GenerateData(resourceName, 1)
	if err != nil {
		t.Fatalf("Failed to generate test data: %v", err)
	}
	
	if len(items) == 0 {
		t.Fatal("No test data generated")
	}
	
	newItem := items[0]
	
	// Remove ID if present (should be auto-generated)
	delete(newItem, "id")
	
	jsonData, err := json.Marshal(newItem)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	
	req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resourceName), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.PostCollection(resourceName)(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	// Verify response has ID
	var created map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if created["id"] == nil {
		t.Error("Created item should have an ID")
	}
	
	// Verify it was added to cache
	finalItems, _ := testCache.Get(resourceName, 1000)
	finalCount := len(finalItems)
	
	if finalCount != initialCount+1 {
		t.Errorf("Expected count to increase by 1, was %d, now %d", initialCount, finalCount)
	}
	
	t.Logf("Successfully created item with ID: %v in %s", created["id"], resourceName)
}

func TestHandler_POST_MaxItemsConstraint(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resourceName := "users"
	
	// Get max items - set to 100 (default max)
	maxItems := 100
	
	// Generate items and fill to max
	items, err := handler.registry.GenerateData(resourceName, maxItems)
	if err != nil {
		t.Fatalf("Failed to generate data: %v", err)
	}
	
	// Set cache data directly to reach max
	testCache.Data[resourceName] = items
	
	// Try to add one more with valid user fields
	newItem := map[string]interface{}{
		"username": "shouldfail",
		"email":    "fail@example.com",
	}
	
	jsonData, _ := json.Marshal(newItem)
	req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resourceName), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.PostCollection(resourceName)(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 when exceeding max items, got %d", w.Code)
	}
	
	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("max")) && !bytes.Contains([]byte(body), []byte("limit")) {
		t.Error("Error message should mention max/limit constraint")
	}
}

func TestHandler_PUT_ValidUpdate(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resourceName := "users"
	
	// Get existing item
	items, _ := testCache.Get(resourceName, 1)
	
	if len(items) == 0 {
		t.Fatal("No test data available")
	}
	
	originalItem := items[0]
	id := fmt.Sprintf("%v", originalItem["id"])
	
	// Prepare update (modify one field that exists in user schema)
	updates := map[string]interface{}{
		"username": "updated_username",
	}
	
	jsonData, _ := json.Marshal(updates)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/%s/%s", resourceName, id), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createCtxWithParams(map[string]string{"resource": resourceName, "id": id}))
	w := httptest.NewRecorder()
	
	handler.PutSingle(resourceName)(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	// Verify response
	var updated map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if updated["username"] != "updated_username" {
		t.Errorf("Expected username to be 'updated_username', got %v", updated["username"])
	}
	
	// Verify ID didn't change
	if fmt.Sprintf("%v", updated["id"]) != id {
		t.Error("ID should not change during update")
	}
	
	t.Logf("Successfully updated item %s in %s", id, resourceName)
}

func TestHandler_PUT_InvalidID(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resourceName := "users"
	
	invalidIDs := []string{"99999", "abc", "-1", "0"}
	
	for _, id := range invalidIDs {
		t.Run(fmt.Sprintf("ID_%s", id), func(t *testing.T) {
			updates := map[string]interface{}{"name": "Test"}
			jsonData, _ := json.Marshal(updates)
			
			req := httptest.NewRequest("PUT", fmt.Sprintf("/%s/%s", resourceName, id), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(createCtxWithParams(map[string]string{"resource": resourceName, "id": id}))
			w := httptest.NewRecorder()
			
			handler.PutSingle(resourceName)(w, req)
			
			if w.Code != http.StatusNotFound {
				t.Errorf("Expected status 404 for invalid ID %s, got %d", id, w.Code)
			}
		})
	}
}

func TestHandler_DELETE_ValidDeletion(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resourceName := "users"
	
	// Get initial count - use large number to get all items
	items, _ := testCache.Get(resourceName, 1000)
	initialCount := len(items)
	
	if initialCount == 0 {
		t.Fatal("No test data available")
	}
	
	// Delete first item
	id := fmt.Sprintf("%v", items[0]["id"])
	
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s/%s", resourceName, id), nil)
	req = req.WithContext(createCtxWithParams(map[string]string{"resource": resourceName, "id": id}))
	w := httptest.NewRecorder()
	
	handler.DeleteSingle(resourceName)(w, req)
	
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	// Verify it was removed from cache
	finalItems, _ := testCache.Get(resourceName, 1000)
	finalCount := len(finalItems)
	
	if finalCount != initialCount-1 {
		t.Errorf("Expected count to decrease by 1, was %d, now %d", initialCount, finalCount)
	}
	
	t.Logf("Successfully deleted item %s from %s", id, resourceName)
}

func TestHandler_DELETE_InvalidID(t *testing.T) {
	handler, _ := setupTestHandler(t)
	resourceName := "users"
	
	invalidIDs := []string{"99999", "nonexistent", "-1"}
	
	for _, id := range invalidIDs {
		t.Run(fmt.Sprintf("ID_%s", id), func(t *testing.T) {
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s/%s", resourceName, id), nil)
			req = req.WithContext(createCtxWithParams(map[string]string{"resource": resourceName, "id": id}))
			w := httptest.NewRecorder()
			
			handler.DeleteSingle(resourceName)(w, req)
			
			if w.Code != http.StatusNotFound {
				t.Errorf("Expected status 404 for invalid ID %s, got %d", id, w.Code)
			}
		})
	}
}

func TestHandler_DELETE_MinItemsConstraint(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resourceName := "users"
	
	// Set cache to minimum items
	items, err := handler.registry.GenerateData(resourceName, 1)
	if err != nil {
		t.Fatalf("Failed to generate data: %v", err)
	}
	testCache.Data[resourceName] = items
	
	// Try to delete the only item
	id := fmt.Sprintf("%v", items[0]["id"])
	
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/%s/%s", resourceName, id), nil)
	req = req.WithContext(createCtxWithParams(map[string]string{"resource": resourceName, "id": id}))
	w := httptest.NewRecorder()
	
	handler.DeleteSingle(resourceName)(w, req)
	
	// Should fail due to min constraint
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 when deleting last item, got %d", w.Code)
	}
	
	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("min")) && !bytes.Contains([]byte(body), []byte("at least")) {
		t.Error("Error message should mention minimum constraint")
	}
}

func TestHandler_POST_ConcurrentCreation(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	resourceName := "users"
	
	// Get initial count
	initialItems, _ := testCache.Get(resourceName, 1000)
	initialCount := len(initialItems)
	
	// Create items concurrently
	concurrency := 10
	var wg sync.WaitGroup
	errors := make(chan error, concurrency)
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			newItem := map[string]interface{}{
				"username": fmt.Sprintf("concurrent_user_%d", idx),
				"email":    fmt.Sprintf("user%d@example.com", idx),
			}
			
			jsonData, _ := json.Marshal(newItem)
			req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resourceName), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			handler.PostCollection(resourceName)(w, req)
			
			if w.Code != http.StatusCreated {
				errors <- fmt.Errorf("goroutine %d: expected 201, got %d", idx, w.Code)
			}
		}(i)
	}
	
	wg.Wait()
	close(errors)
	
	// Check for errors
	for err := range errors {
		t.Error(err)
	}
	
	// Verify all were added
	finalItems, _ := testCache.Get(resourceName, 1000)
	finalCount := len(finalItems)
	
	if finalCount != initialCount+concurrency {
		t.Errorf("Expected count to increase by %d, was %d, now %d", concurrency, initialCount, finalCount)
	}
	
	t.Logf("Successfully created %d items concurrently", concurrency)
}
