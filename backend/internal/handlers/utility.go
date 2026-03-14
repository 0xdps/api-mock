package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// UtilityHandlers provides testing utility endpoints
type UtilityHandlers struct{}

// NewUtilityHandlers creates a new utility handler
func NewUtilityHandlers() *UtilityHandlers {
	return &UtilityHandlers{}
}

// Echo returns the full request details
func (h *UtilityHandlers) Echo(w http.ResponseWriter, r *http.Request) {
	var body interface{}

	contentType := r.Header.Get("Content-Type")

	// Read body first as raw bytes (this is more reliable)
	var bodyBytes []byte
	if r.Body != nil {
		defer r.Body.Close()
		bodyBytes, _ = io.ReadAll(r.Body)
	}

	// Handle different content types
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		// Parse form data from body bytes
		if len(bodyBytes) > 0 {
			values, err := url.ParseQuery(string(bodyBytes))
			if err == nil {
				formData := make(map[string]interface{})
				for key, vals := range values {
					if len(vals) == 1 {
						formData[key] = vals[0]
					} else {
						formData[key] = vals
					}
				}
				body = formData
			}
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// For multipart, we need to recreate the body reader
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse multipart form data (limit to 10MB)
		if err := r.ParseMultipartForm(10 << 20); err == nil {
			formData := make(map[string]interface{})

			// Add form values
			for key, values := range r.MultipartForm.Value {
				if len(values) == 1 {
					formData[key] = values[0]
				} else {
					formData[key] = values
				}
			}

			// Add file information (not the actual file content)
			if len(r.MultipartForm.File) > 0 {
				files := make(map[string]interface{})
				for key, fileHeaders := range r.MultipartForm.File {
					fileInfos := make([]map[string]interface{}, len(fileHeaders))
					for i, fh := range fileHeaders {
						fileInfos[i] = map[string]interface{}{
							"filename": fh.Filename,
							"size":     fh.Size,
							"header":   fh.Header,
						}
					}
					if len(fileInfos) == 1 {
						files[key] = fileInfos[0]
					} else {
						files[key] = fileInfos
					}
				}
				formData["_files"] = files
			}

			body = formData
		}
	} else if len(bodyBytes) > 0 {
		// Try JSON first
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			// If not JSON, return as raw string
			body = string(bodyBytes)
		}
	}

	// Collect headers
	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// Collect query parameters
	query := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	response := map[string]interface{}{
		"method":      r.Method,
		"path":        r.URL.Path,
		"headers":     headers,
		"query":       query,
		"body":        body,
		"contentType": contentType,
		"timestamp":   time.Now().Unix(),
		"host":        r.Host,
		"remote":      r.RemoteAddr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delay adds artificial latency to the response
func (h *UtilityHandlers) Delay(w http.ResponseWriter, r *http.Request) {
	// Get delay from URL parameter
	delayStr := chi.URLParam(r, "ms")
	delayMs, err := strconv.Atoi(delayStr)
	if err != nil || delayMs < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid delay parameter, must be a positive integer (milliseconds)",
		})
		return
	}

	// Cap at 30 seconds to prevent abuse
	if delayMs > 30000 {
		delayMs = 30000
	}

	// Sleep for the specified duration
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"delayed":  true,
		"delay_ms": delayMs,
		"message":  fmt.Sprintf("Delayed for %d milliseconds", delayMs),
	})
}

// DelayRandom adds random latency within a range
func (h *UtilityHandlers) DelayRandom(w http.ResponseWriter, r *http.Request) {
	// Get min and max from query parameters
	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")

	min := 100  // default 100ms
	max := 2000 // default 2s

	if minStr != "" {
		if m, err := strconv.Atoi(minStr); err == nil && m >= 0 {
			min = m
		}
	}

	if maxStr != "" {
		if m, err := strconv.Atoi(maxStr); err == nil && m > min {
			max = m
		}
	}

	// Cap at 30 seconds
	if max > 30000 {
		max = 30000
	}

	// Random delay between min and max
	delayMs := min + rand.Intn(max-min+1)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"delayed":  true,
		"delay_ms": delayMs,
		"min":      min,
		"max":      max,
		"message":  fmt.Sprintf("Randomly delayed for %d milliseconds", delayMs),
	})
}

// Status returns a specific HTTP status code
func (h *UtilityHandlers) Status(w http.ResponseWriter, r *http.Request) {
	// Get status code from URL parameter
	codeStr := chi.URLParam(r, "code")
	code, err := strconv.Atoi(codeStr)
	if err != nil || code < 100 || code >= 600 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid status code, must be between 100-599",
		})
		return
	}

	// Build response based on status code category
	var message string
	var details interface{}

	switch {
	case code >= 200 && code < 300:
		message = "Success"
		details = map[string]interface{}{
			"category": "2xx Success",
			"note":     "The request was successfully processed",
		}
	case code >= 300 && code < 400:
		message = "Redirect"
		details = map[string]interface{}{
			"category": "3xx Redirect",
			"note":     "The resource has moved or needs redirection",
		}
	case code >= 400 && code < 500:
		message = "Client Error"
		details = map[string]interface{}{
			"category": "4xx Client Error",
			"note":     "The request contains invalid data or unauthorized",
		}
	case code >= 500 && code < 600:
		message = "Server Error"
		details = map[string]interface{}{
			"category": "5xx Server Error",
			"note":     "The server encountered an error processing the request",
		}
	default:
		message = "Informational"
		details = map[string]interface{}{
			"category": "1xx Informational",
			"note":     "Request received, continuing process",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    code,
		"message": message,
		"details": details,
	})
}

// ErrorValidation returns a 422 with validation errors
func (h *UtilityHandlers) ErrorValidation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "Validation failed",
		"errors": []map[string]interface{}{
			{
				"field":   "email",
				"message": "Email is required and must be valid",
				"code":    "required",
			},
			{
				"field":   "age",
				"message": "Age must be between 18 and 100",
				"code":    "range",
			},
			{
				"field":   "username",
				"message": "Username must be at least 3 characters",
				"code":    "min_length",
			},
		},
	})
}

// Flaky randomly succeeds or fails based on successRate
func (h *UtilityHandlers) Flaky(w http.ResponseWriter, r *http.Request) {
	// Get success rate from query parameter (default 0.5 = 50%)
	successRateStr := r.URL.Query().Get("successRate")
	successRate := 0.5

	if successRateStr != "" {
		if rate, err := strconv.ParseFloat(successRateStr, 64); err == nil {
			if rate >= 0 && rate <= 1 {
				successRate = rate
			}
		}
	}

	// Random decision
	roll := rand.Float64()
	success := roll < successRate

	w.Header().Set("Content-Type", "application/json")

	if success {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":      true,
			"successRate":  successRate,
			"roll":         roll,
			"message":      "Request succeeded",
			"attempt_hint": "This endpoint randomly succeeds based on successRate parameter",
		})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":      false,
			"successRate":  successRate,
			"roll":         roll,
			"error":        "Service temporarily unavailable",
			"message":      "Request failed (random failure for testing)",
			"attempt_hint": "Try again - this endpoint randomly fails based on successRate parameter",
		})
	}
}

// Chaos returns random status codes and response shapes
func (h *UtilityHandlers) Chaos(w http.ResponseWriter, r *http.Request) {
	// Random status codes (weighted towards common ones)
	statusCodes := []int{
		200, 200, 200, 200, 200, // 50% success
		400, 401, 403, 404, // 40% client errors
		500, 502, 503, // 30% server errors
	}
	statusCode := statusCodes[rand.Intn(len(statusCodes))]

	// Random response shapes
	responseShapes := []interface{}{
		map[string]interface{}{
			"data":    "Simple string response",
			"chaos":   true,
			"code":    statusCode,
			"message": "Chaos response - shape 1",
		},
		map[string]interface{}{
			"result": []int{1, 2, 3, 4, 5},
			"meta": map[string]string{
				"chaos": "true",
				"type":  "array",
			},
			"code": statusCode,
		},
		map[string]interface{}{
			"error":   true,
			"code":    statusCode,
			"details": "Something went wrong in chaos mode",
			"hint":    "This is expected - chaos endpoint returns random shapes",
		},
		[]string{"chaos", "mode", "activated"},
		map[string]interface{}{
			"nested": map[string]interface{}{
				"deeply": map[string]interface{}{
					"buried": map[string]interface{}{
						"data":  "Found me!",
						"chaos": true,
						"code":  statusCode,
					},
				},
			},
		},
	}

	response := responseShapes[rand.Intn(len(responseShapes))]

	// Random headers
	if rand.Float64() < 0.3 {
		w.Header().Set("X-Chaos-Mode", "active")
	}
	if rand.Float64() < 0.3 {
		w.Header().Set("X-Random-Header", fmt.Sprintf("chaos-%d", rand.Intn(1000)))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
