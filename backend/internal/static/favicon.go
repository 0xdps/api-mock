package static

import (
	"net/http"
)

// Fallback handler if favicon not embedded
func ServeFaviconFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
