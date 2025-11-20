package static

import (
	"fmt"
	"net/http"

	_ "embed"
)

//go:embed assets/icon.svg
var IconSVG []byte

//go:embed assets/logo.svg
var LogoSVG []byte

func init() {
	fmt.Printf("[DEBUG] IconSVG loaded: %d bytes\n", len(IconSVG))
	fmt.Printf("[DEBUG] LogoSVG loaded: %d bytes\n", len(LogoSVG))
}

// ServeIconSVG writes the embedded SVG as image/svg+xml with caching headers.
func ServeIconSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(IconSVG)
}

// ServeIconICO writes the embedded icon SVG with proper icon content-type.
// Modern browsers support SVG favicons; this serves the same SVG with appropriate headers.
func ServeIconICO(w http.ResponseWriter, r *http.Request) {
	// Serve SVG with appropriate cache headers for favicon
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(IconSVG)
}
