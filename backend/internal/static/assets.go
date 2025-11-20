package static

import (
	"net/http"

	_ "embed"
)

//go:embed assets/icon.svg
var IconSVG []byte

//go:embed assets/logo.svg
var LogoSVG []byte

// ServeIconSVG writes the embedded SVG as image/svg+xml
func ServeIconSVG(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "image/svg+xml")
    w.Write(IconSVG)
}

// ServeIconICO writes the embedded SVG but uses the ico content-type.
// Some browsers will accept SVG as a favicon; if you prefer a true .ico
// file generate it and embed that instead.
func ServeIconICO(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "image/x-icon")
    w.Write(IconSVG)
}
