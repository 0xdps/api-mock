package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootResponseIncludesPlatform(t *testing.T) {
	t.Setenv("PLATFORM", "test")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message":   "Mockly API",
			"version":   "1.0.0",
			"docs":      "https://www.mockly.codes/docs",
			"resources": []string{"users"},
			"groups":    map[string][]string{"people": []string{"users"}},
			"platform":  detectPlatform(),
		}
		json.NewEncoder(w).Encode(response)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["platform"] != "test" {
		t.Fatalf("expected platform=test, got %v", payload["platform"])
	}
}
