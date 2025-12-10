package middleware

import (
"net/http"
"net/http/httptest"
"testing"

"github.com/go-chi/chi/v5"
)

func TestSetupCORS_ConfigurationExists(t *testing.T) {
	corsHandler := SetupCORS()
	if corsHandler == nil {
		t.Fatal("Expected CORS handler to be non-nil")
	}
}

func TestCORS_AllowedOrigins(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
w.Write([]byte("OK"))
})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	acoHeader := w.Header().Get("Access-Control-Allow-Origin")
	if acoHeader != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin: *, got %s", acoHeader)
	}
}

func TestCORS_PreflightRequest(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Errorf("Expected status 200 or 204 for OPTIONS, got %d", w.Code)
	}

	acamHeader := w.Header().Get("Access-Control-Allow-Methods")
	if acamHeader == "" {
		t.Error("Expected Access-Control-Allow-Methods header to be set")
	}
}

func TestCORS_AllowedMethods(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Post("/test", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
})

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
req := httptest.NewRequest("OPTIONS", "/test", nil)
req.Header.Set("Origin", "http://example.com")
req.Header.Set("Access-Control-Request-Method", method)
w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			acamHeader := w.Header().Get("Access-Control-Allow-Methods")
			if acamHeader == "" {
				t.Errorf("Expected Access-Control-Allow-Methods header for %s", method)
			}
		})
	}
}

func TestCORS_AllowedHeaders(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
})

	headers := []string{"Content-Type", "Authorization", "X-No-Cache"}
	for _, header := range headers {
		t.Run(header, func(t *testing.T) {
req := httptest.NewRequest("OPTIONS", "/test", nil)
req.Header.Set("Origin", "http://example.com")
req.Header.Set("Access-Control-Request-Method", "GET")
req.Header.Set("Access-Control-Request-Headers", header)
w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			acahHeader := w.Header().Get("Access-Control-Allow-Headers")
			if acahHeader == "" {
				t.Errorf("Expected Access-Control-Allow-Headers for %s", header)
			}
		})
	}
}

func TestCORS_MaxAge(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	maxAgeHeader := w.Header().Get("Access-Control-Max-Age")
	if maxAgeHeader == "" {
		t.Error("Expected Access-Control-Max-Age header to be set")
	}
}

func TestCORS_ExposedHeaders(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Link", "test-link")
		w.Header().Set("X-Cache", "HIT")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Header().Get("Link") != "test-link" {
		t.Error("Expected Link header to be preserved")
	}
	if w.Header().Get("X-Cache") != "HIT" {
		t.Error("Expected X-Cache header to be preserved")
	}
}

func TestCORS_CredentialsNotAllowed(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	acacHeader := w.Header().Get("Access-Control-Allow-Credentials")
	if acacHeader == "true" {
		t.Error("Expected credentials to not be allowed")
	}
}

func TestCORS_MultipleOrigins(t *testing.T) {
	router := chi.NewRouter()
	router.Use(SetupCORS().Handler)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
})

	origins := []string{
		"http://localhost:3000",
		"https://example.com",
		"https://test.org",
	}

	for _, origin := range origins {
		t.Run(origin, func(t *testing.T) {
req := httptest.NewRequest("GET", "/test", nil)
req.Header.Set("Origin", origin)
w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			acoHeader := w.Header().Get("Access-Control-Allow-Origin")
			if acoHeader != "*" {
				t.Errorf("Expected wildcard origin, got %s", acoHeader)
			}
		})
	}
}
