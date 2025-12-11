package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestUtility_Echo(t *testing.T) {
	handler := NewUtilityHandlers()

	t.Run("GET request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/echo?foo=bar&count=5", nil)
		req.Header.Set("User-Agent", "test-client")
		w := httptest.NewRecorder()

		handler.Echo(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["method"] != "GET" {
			t.Errorf("Expected method GET, got %v", response["method"])
		}

		query := response["query"].(map[string]interface{})
		if query["foo"] != "bar" {
			t.Error("Expected query param foo=bar")
		}
	})

	t.Run("POST request with JSON body", func(t *testing.T) {
		body := strings.NewReader(`{"name": "test", "value": 123}`)
		req := httptest.NewRequest("POST", "/test/echo", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Echo(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["method"] != "POST" {
			t.Errorf("Expected method POST, got %v", response["method"])
		}

		bodyMap := response["body"].(map[string]interface{})
		if bodyMap["name"] != "test" {
			t.Error("Expected body to contain name=test")
		}
	})

	t.Run("POST request with form-urlencoded", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("username", "john_doe")
		formData.Set("email", "john@example.com")
		formData.Set("age", "25")
		
		req := httptest.NewRequest("POST", "/test/echo", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.Echo(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["method"] != "POST" {
			t.Errorf("Expected method POST, got %v", response["method"])
		}

		if response["contentType"] != "application/x-www-form-urlencoded" {
			t.Errorf("Expected contentType application/x-www-form-urlencoded, got %v", response["contentType"])
		}

		bodyMap := response["body"].(map[string]interface{})
		if bodyMap["username"] != "john_doe" {
			t.Errorf("Expected username=john_doe, got %v", bodyMap["username"])
		}
		if bodyMap["email"] != "john@example.com" {
			t.Errorf("Expected email=john@example.com, got %v", bodyMap["email"])
		}
		if bodyMap["age"] != "25" {
			t.Errorf("Expected age=25, got %v", bodyMap["age"])
		}
	})

	t.Run("POST request with multipart form data", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		
		// Add form fields
		writer.WriteField("title", "Test Document")
		writer.WriteField("description", "This is a test")
		writer.WriteField("tags", "test")
		writer.WriteField("tags", "demo")
		
		// Add a fake file
		part, err := writer.CreateFormFile("document", "test.txt")
		if err != nil {
			t.Fatal(err)
		}
		part.Write([]byte("This is test file content"))
		
		writer.Close()
		
		req := httptest.NewRequest("POST", "/test/echo", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		handler.Echo(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["method"] != "POST" {
			t.Errorf("Expected method POST, got %v", response["method"])
		}

		bodyMap := response["body"].(map[string]interface{})
		if bodyMap["title"] != "Test Document" {
			t.Errorf("Expected title='Test Document', got %v", bodyMap["title"])
		}
		if bodyMap["description"] != "This is a test" {
			t.Errorf("Expected description='This is a test', got %v", bodyMap["description"])
		}

		// Check files info
		files := bodyMap["_files"].(map[string]interface{})
		if files == nil {
			t.Error("Expected _files to be present")
		}
		
		docFile := files["document"].(map[string]interface{})
		if docFile["filename"] != "test.txt" {
			t.Errorf("Expected filename=test.txt, got %v", docFile["filename"])
		}
	})

	t.Run("POST request with plain text", func(t *testing.T) {
		body := strings.NewReader("This is plain text content")
		req := httptest.NewRequest("POST", "/test/echo", body)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		handler.Echo(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["method"] != "POST" {
			t.Errorf("Expected method POST, got %v", response["method"])
		}

		bodyStr := response["body"].(string)
		if bodyStr != "This is plain text content" {
			t.Errorf("Expected plain text body, got %v", bodyStr)
		}
	})
}

func TestUtility_Delay(t *testing.T) {
	handler := NewUtilityHandlers()

	t.Run("valid delay", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/delay/100", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("ms", "100")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		start := time.Now()
		handler.Delay(w, req)
		duration := time.Since(start)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		if duration < 100*time.Millisecond {
			t.Errorf("Expected delay >= 100ms, got %v", duration)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["delayed"] != true {
			t.Error("Expected delayed=true")
		}
	})

	t.Run("invalid delay", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/delay/invalid", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("ms", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.Delay(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("delay capped at 30 seconds", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/delay/60000", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("ms", "60000")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		start := time.Now()
		handler.Delay(w, req)
		duration := time.Since(start)

		// Should be capped at 30s, so shouldn't take 60s
		if duration > 31*time.Second {
			t.Error("Delay should be capped at 30 seconds")
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		// Response should show capped value
		if response["delay_ms"].(float64) != 30000 {
			t.Errorf("Expected capped delay_ms=30000, got %v", response["delay_ms"])
		}
	})
}

func TestUtility_DelayRandom(t *testing.T) {
	handler := NewUtilityHandlers()

	t.Run("random delay with defaults", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/delay-random", nil)
		w := httptest.NewRecorder()

		start := time.Now()
		handler.DelayRandom(w, req)
		duration := time.Since(start)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		// Default is 100-2000ms
		if duration < 100*time.Millisecond || duration > 2100*time.Millisecond {
			t.Logf("Delay was %v, expected between 100ms-2000ms (with tolerance)", duration)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["delayed"] != true {
			t.Error("Expected delayed=true")
		}
	})

	t.Run("random delay with custom range", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/delay-random?min=50&max=150", nil)
		w := httptest.NewRecorder()

		start := time.Now()
		handler.DelayRandom(w, req)
		duration := time.Since(start)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		if duration < 50*time.Millisecond || duration > 200*time.Millisecond {
			t.Logf("Delay was %v, expected between 50ms-150ms (with tolerance)", duration)
		}
	})
}

func TestUtility_Status(t *testing.T) {
	handler := NewUtilityHandlers()

	testCases := []struct {
		code     string
		expected int
	}{
		{"200", 200},
		{"201", 201},
		{"400", 400},
		{"404", 404},
		{"500", 500},
		{"503", 503},
	}

	for _, tc := range testCases {
		t.Run("status_"+tc.code, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test/status/"+tc.code, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("code", tc.code)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.Status(w, req)

			if w.Code != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, w.Code)
			}

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			if response["code"].(float64) != float64(tc.expected) {
				t.Errorf("Expected code %d in response, got %v", tc.expected, response["code"])
			}
		})
	}

	t.Run("invalid status code", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/status/999", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("code", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.Status(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 for invalid code, got %d", w.Code)
		}
	})
}

func TestUtility_ErrorValidation(t *testing.T) {
	handler := NewUtilityHandlers()

	req := httptest.NewRequest("GET", "/test/error/validation", nil)
	w := httptest.NewRecorder()

	handler.ErrorValidation(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["error"] == "" {
		t.Error("Expected error message")
	}

	errors := response["errors"].([]interface{})
	if len(errors) == 0 {
		t.Error("Expected validation errors array")
	}

	// Check first error has required fields
	firstError := errors[0].(map[string]interface{})
	if firstError["field"] == nil || firstError["message"] == nil {
		t.Error("Expected error to have field and message")
	}
}

func TestUtility_Flaky(t *testing.T) {
	handler := NewUtilityHandlers()

	t.Run("flaky with 100% success rate", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/flaky?successRate=1.0", nil)
		w := httptest.NewRecorder()

		handler.Flaky(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 with successRate=1.0, got %d", w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["success"] != true {
			t.Error("Expected success=true with successRate=1.0")
		}
	})

	t.Run("flaky with 0% success rate", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/flaky?successRate=0.0", nil)
		w := httptest.NewRecorder()

		handler.Flaky(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected 503 with successRate=0.0, got %d", w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if response["success"] != false {
			t.Error("Expected success=false with successRate=0.0")
		}
	})

	t.Run("flaky with default rate (should sometimes fail)", func(t *testing.T) {
		// Run multiple times to ensure randomness
		successCount := 0
		failCount := 0

		for i := 0; i < 20; i++ {
			req := httptest.NewRequest("GET", "/test/flaky", nil)
			w := httptest.NewRecorder()

			handler.Flaky(w, req)

			if w.Code == http.StatusOK {
				successCount++
			} else {
				failCount++
			}
		}

		// With default 50% rate, we should have both successes and failures
		// (with 20 requests, extremely unlikely to get all success or all fail)
		if successCount == 0 || failCount == 0 {
			t.Logf("After 20 requests: %d successes, %d failures", successCount, failCount)
			t.Log("Expected mix of successes and failures with default 50% rate")
		}
	})
}

func TestUtility_Chaos(t *testing.T) {
	handler := NewUtilityHandlers()

	// Run multiple times to see variety
	statusCodes := make(map[int]int)

	for i := 0; i < 50; i++ {
		req := httptest.NewRequest("GET", "/test/chaos", nil)
		w := httptest.NewRecorder()

		handler.Chaos(w, req)

		statusCodes[w.Code]++

		// Should always return valid JSON
		var response interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Errorf("Chaos should return valid JSON, got error: %v", err)
		}
	}

	// Should have generated multiple different status codes
	if len(statusCodes) < 2 {
		t.Errorf("Expected multiple different status codes, got only %v", statusCodes)
	}

	t.Logf("Chaos generated status codes: %v", statusCodes)
}

func TestUtility_ChaosResponseShapes(t *testing.T) {
	handler := NewUtilityHandlers()

	responseTypes := make(map[string]int)

	for i := 0; i < 30; i++ {
		req := httptest.NewRequest("GET", "/test/chaos", nil)
		w := httptest.NewRecorder()

		handler.Chaos(w, req)

		var response interface{}
		json.NewDecoder(w.Body).Decode(&response)

		// Categorize response type
		switch v := response.(type) {
		case map[string]interface{}:
			if _, ok := v["data"]; ok {
				responseTypes["shape1"]++
			} else if _, ok := v["result"]; ok {
				responseTypes["shape2"]++
			} else if _, ok := v["nested"]; ok {
				responseTypes["shape5"]++
			} else {
				responseTypes["shape3"]++
			}
		case []interface{}:
			responseTypes["shape4"]++
		}
	}

	// Should have generated multiple different response shapes
	if len(responseTypes) < 2 {
		t.Logf("Expected multiple response shapes, got: %v", responseTypes)
	}

	t.Logf("Chaos generated response shapes: %v", responseTypes)
}
