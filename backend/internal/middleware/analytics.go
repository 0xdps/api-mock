package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

type umamiBackendTracker struct {
	enabled    bool
	websiteID  string
	hostname   string
	collectURL string
	client     *http.Client

	events chan []byte
}

type umamiEnvelope struct {
	Type    string     `json:"type"`
	Payload umamiEvent `json:"payload"`
}

type umamiEvent struct {
	Website  string                 `json:"website"`
	Hostname string                 `json:"hostname"`
	Language string                 `json:"language"`
	Referrer string                 `json:"referrer"`
	URL      string                 `json:"url"`
	Name     string                 `json:"name"`
	Data     map[string]interface{} `json:"data"`
}

func NewBackendAnalyticsMiddleware(enabled bool, websiteID, collectURL, hostname string) func(http.Handler) http.Handler {
	tracker := &umamiBackendTracker{
		enabled:    enabled && websiteID != "" && collectURL != "",
		websiteID:  websiteID,
		hostname:   hostname,
		collectURL: collectURL,
		client: &http.Client{
			Timeout: 1500 * time.Millisecond,
		},
		events: make(chan []byte, 1024),
	}

	if tracker.enabled {
		go tracker.runWorker()
		log.Printf("📈 Backend analytics enabled (endpoint: %s)", tracker.collectURL)
	} else {
		log.Printf("📈 Backend analytics disabled")
	}

	return func(next http.Handler) http.Handler {
		if !tracker.enabled {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rec, r)

			tracker.trackRequest(r, rec.status, time.Since(start).Milliseconds())
		})
	}
}

func (t *umamiBackendTracker) runWorker() {
	for payload := range t.events {
		req, err := http.NewRequest(http.MethodPost, t.collectURL, bytes.NewReader(payload))
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "api-mockly-backend-analytics/1.0")

		resp, err := t.client.Do(req)
		if err != nil {
			continue
		}

		_ = resp.Body.Close()
	}
}

func (t *umamiBackendTracker) trackRequest(r *http.Request, status int, durationMs int64) {
	if shouldSkipAnalyticsPath(r.URL.Path) {
		return
	}

	routePattern := ""
	if routeCtx := chi.RouteContext(r.Context()); routeCtx != nil {
		routePattern = routeCtx.RoutePattern()
	}

	payload := umamiEnvelope{
		Type: "event",
		Payload: umamiEvent{
			Website:  t.websiteID,
			Hostname: t.hostname,
			Language: "en-US",
			Referrer: r.Referer(),
			URL:      requestURL(r),
			Name:     "be_req",
			Data: map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"route":       routePattern,
				"status":      status,
				"duration_ms": durationMs,
				"has_query":   r.URL.RawQuery != "",
				"source":      "backend",
			},
		},
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return
	}

	select {
	case t.events <- encoded:
	default:
		// Drop event when buffer is full to avoid impacting API latency.
	}
}

func requestURL(r *http.Request) string {
	host := r.Host
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("https://%s%s", host, r.URL.RequestURI())
}

func shouldSkipAnalyticsPath(path string) bool {
	switch path {
	case "/favicon.ico", "/favicon.svg", "/":
		return true
	default:
		return false
	}
}
