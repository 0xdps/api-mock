package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0xdps/api-mock/go/internal/cache"
	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/go-chi/chi/v5"
)

// ============================================================================
// SETUP HELPERS
// ============================================================================

func setupValidationTest(t *testing.T) *DynamicHandler {
	t.Helper()
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	cacheConfig := cache.Config{
		ItemsPerResource:    10,
		Seed:                42,
		MaxItemsPerResource: 100,
	}

	testCache := cache.NewCache(registry, cacheConfig, nil)
	handler := NewDynamicHandler(registry, testCache)
	return handler
}

// ============================================================================
// TYPE VALIDATION TESTS
// ============================================================================

func TestValidation_TypeChecking(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		resource       string
		data           map[string]interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "valid user with correct types",
			resource:       "users",
			data:           map[string]interface{}{"username": "john", "email": "john@example.com", "age": 25, "gender": "male"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid - age should be number",
			resource:       "users",
			data:           map[string]interface{}{"username": "john", "email": "john@example.com", "age": "25", "gender": "male"},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "age",
		},
		{
			name:           "invalid - username should be string",
			resource:       "users",
			data:           map[string]interface{}{"username": 123, "email": "john@example.com", "age": 25, "gender": "male"},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.data)
			req := httptest.NewRequest("POST", "/"+tt.resource, bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection(tt.resource)(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// REQUIRED FIELDS VALIDATION TESTS
// ============================================================================

func TestValidation_RequiredFields(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		resource       string
		data           map[string]interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "valid user with all required fields",
			resource:       "users",
			data:           map[string]interface{}{"username": "john", "email": "john@example.com", "age": 25, "gender": "male"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "missing required field - username",
			resource:       "users",
			data:           map[string]interface{}{"email": "john@example.com", "age": 25, "gender": "male"},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "username",
		},
		{
			name:           "missing required field - email",
			resource:       "users",
			data:           map[string]interface{}{"username": "john", "age": 25, "gender": "male"},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.data)
			req := httptest.NewRequest("POST", "/"+tt.resource, bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection(tt.resource)(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// FORMAT VALIDATION TESTS
// ============================================================================

func TestValidation_EmailFormat(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		email          string
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "valid email",
			email:          "john@example.com",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid email with subdomain",
			email:          "john.doe@mail.example.com",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid email - missing @",
			email:          "johnexample.com",
			expectedStatus: http.StatusBadRequest,
			errorContains:  "email",
		},
		{
			name:           "invalid email - plain text",
			email:          "notanemail",
			expectedStatus: http.StatusBadRequest,
			errorContains:  "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]interface{}{
				"username": "john",
				"email":    tt.email,
				"age":      25,
				"gender":   "male",
			}
			body, _ := json.Marshal(data)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection("users")(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// ENUM VALIDATION TESTS
// ============================================================================

func TestValidation_EnumConstraint(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		gender         string
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "valid enum value - male",
			gender:         "male",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid enum value - female",
			gender:         "female",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid enum value",
			gender:         "other",
			expectedStatus: http.StatusBadRequest,
			errorContains:  "gender",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]interface{}{
				"username": "john",
				"email":    "john@example.com",
				"age":      25,
				"gender":   tt.gender,
			}
			body, _ := json.Marshal(data)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection("users")(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// NUMERIC CONSTRAINTS VALIDATION TESTS
// ============================================================================

func TestValidation_NumericConstraints(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		age            interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "valid age - minimum boundary",
			age:            0,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid age - normal value",
			age:            25,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid age - maximum boundary",
			age:            150,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid age - below minimum",
			age:            -1,
			expectedStatus: http.StatusBadRequest,
			errorContains:  "age",
		},
		{
			name:           "invalid age - above maximum",
			age:            151,
			expectedStatus: http.StatusBadRequest,
			errorContains:  "age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]interface{}{
				"username": "john",
				"email":    "john@example.com",
				"age":      tt.age,
				"gender":   "male",
			}
			body, _ := json.Marshal(data)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection("users")(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// STRING CONSTRAINTS VALIDATION TESTS
// ============================================================================

func TestValidation_StringConstraints(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		username       string
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "valid username - minimum length",
			username:       "abc",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid username - normal length",
			username:       "john_doe",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid username - maximum length",
			username:       "abcdefghijklmnopqrst",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid username - too short",
			username:       "ab",
			expectedStatus: http.StatusBadRequest,
			errorContains:  "username",
		},
		{
			name:           "invalid username - too long",
			username:       "abcdefghijklmnopqrstu",
			expectedStatus: http.StatusBadRequest,
			errorContains:  "username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]interface{}{
				"username": tt.username,
				"email":    "john@example.com",
				"age":      25,
				"gender":   "male",
			}
			body, _ := json.Marshal(data)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection("users")(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// ARRAY CONSTRAINTS VALIDATION TESTS
// ============================================================================

func TestValidation_ArrayConstraints(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		tags           []interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "valid tags - minimum items",
			tags:           []interface{}{"tag1"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid tags - normal count",
			tags:           []interface{}{"tag1", "tag2", "tag3"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid tags - below minimum",
			tags:           []interface{}{},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "tags",
		},
		{
			name:           "invalid tags - above maximum",
			tags:           []interface{}{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "tags",
		},
		{
			name:           "invalid tags - duplicate items",
			tags:           []interface{}{"tag1", "tag1"},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "tags",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]interface{}{
				"title":   "Test Article",
				"content": "Test content",
				"tags":    tt.tags,
			}
			body, _ := json.Marshal(data)
			req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection("articles")(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// ADDITIONAL PROPERTIES VALIDATION TESTS
// ============================================================================

func TestValidation_AdditionalProperties(t *testing.T) {
	handler := setupValidationTest(t)

	tests := []struct {
		name           string
		data           map[string]interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name: "valid - only allowed properties",
			data: map[string]interface{}{
				"username": "john",
				"email":    "john@example.com",
				"age":      25,
				"gender":   "male",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid - additional property not allowed",
			data: map[string]interface{}{
				"username":   "john",
				"email":      "john@example.com",
				"age":        25,
				"gender":     "male",
				"extraField": "not allowed",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "additional",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.data)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.PostCollection("users")(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// UPDATE (PUT) VALIDATION TESTS
// ============================================================================

func TestValidation_UpdateOperation(t *testing.T) {
	handler := setupValidationTest(t)

	// First, create a user via POST
	createData := map[string]interface{}{
		"username": "john",
		"email":    "john@example.com",
		"age":      25,
		"gender":   "male",
	}
	body, _ := json.Marshal(createData)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.PostCollection("users")(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create user: %d - %s", w.Code, w.Body.String())
	}

	// Extract the created user ID
	var createdUser map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createdUser)
	
	var userID string
	switch id := createdUser["id"].(type) {
	case string:
		userID = id
	case int:
		userID = fmt.Sprintf("%d", id)
	case float64:
		userID = fmt.Sprintf("%.0f", id)
	default:
		t.Fatalf("Created user has no ID or ID is wrong type: %T", createdUser["id"])
	}

	tests := []struct {
		name           string
		updateData     map[string]interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name: "valid update - change email",
			updateData: map[string]interface{}{
				"email": "newemail@example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update - change age",
			updateData: map[string]interface{}{
				"age": 30,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid update - bad email format",
			updateData: map[string]interface{}{
				"email": "notanemail",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "email",
		},
		{
			name: "invalid update - age below minimum",
			updateData: map[string]interface{}{
				"age": -5,
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "age",
		},
		{
			name: "invalid update - invalid gender enum",
			updateData: map[string]interface{}{
				"gender": "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "gender",
		},
		{
			name: "invalid update - username too short",
			updateData: map[string]interface{}{
				"username": "ab",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateBody, _ := json.Marshal(tt.updateData)
			req := httptest.NewRequest("PUT", "/users/"+userID, bytes.NewBuffer(updateBody))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			handler.PutSingle("users")(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.errorContains != "" && !strings.Contains(w.Body.String(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, w.Body.String())
			}
		})
	}
}

// ============================================================================
// DUPLICATE ID TESTS  
// ============================================================================

func TestValidation_DuplicateID(t *testing.T) {
	handler := setupValidationTest(t)

	// First, create a user to get an ID
	createData := map[string]interface{}{
		"username": "john",
		"email":    "john@example.com",
		"age":      25,
		"gender":   "male",
	}
	body, _ := json.Marshal(createData)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.PostCollection("users")(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create user: %d - %s", w.Code, w.Body.String())
	}

	t.Run("reject POST with ID field", func(t *testing.T) {
		// Try to create a user with an ID field
		dataWithID := map[string]interface{}{
			"id":       12345,
			"username": "jane",
			"email":    "jane@example.com",
			"age":      30,
			"gender":   "female",
		}
		body, _ := json.Marshal(dataWithID)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.PostCollection("users")(w, req)

		// Should return 400 Bad Request
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %d. Body: %s", w.Code, w.Body.String())
		}

		if !strings.Contains(w.Body.String(), "Cannot specify 'id' field") {
			t.Errorf("Expected error about ID field not allowed, got: %s", w.Body.String())
		}
	})

	t.Run("auto-generate ID when not provided", func(t *testing.T) {
		// Create user without ID field
		dataWithoutID := map[string]interface{}{
			"username": "alice",
			"email":    "alice@example.com",
			"age":      28,
			"gender":   "female",
		}
		body, _ := json.Marshal(dataWithoutID)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.PostCollection("users")(w, req)

		// Should succeed
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201 Created, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Verify ID was auto-generated
		var createdUser map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createdUser)
		if createdUser["id"] == nil {
			t.Errorf("Expected ID to be auto-generated, but it's missing")
		}
	})
}
